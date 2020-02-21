package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"time"

	"github.com/pkg/errors"

	adm "github.com/adaptiveteam/adaptive/adaptive-dynamic-menu"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/checks"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"strconv"
	"strings"

	aesc "github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/holidays"
	holidaysLambda "github.com/adaptiveteam/adaptive/lambdas/holidays-lambda-go"
	competencies "github.com/adaptiveteam/adaptive/lambdas/competencies-lambda-go"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/values"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

func GwOk(resp string, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       string(resp),
		StatusCode: 200,
	}, err
}

const (
	HelloWorldNamespace           = "hello-world"
	HolidaysNamespace             = "holidays"
	AdaptiveValuesNamespace       = "adaptive_values"
	AdaptiveChannelNamespace      = "adaptive-channel"
	NoAdaptiveAccessDialogContext = "dialog/engagements/adaptive-access/"
)

// -- Direct Adaptive to create an engagement alerting team members of upcoming holidays
// Enable users to request a list of all upcoming holidays

func checkValues(userID string) checks.CheckResultMap {
	dms := adm.AdaptiveDynamicMenu(aesc.ProductionProfile, bindings)
	loc, _ := time.LoadLocation("UTC")
	today := business_time.Today(loc)
	return dms.StripOutFunctions().Evaluate(userID, today)
}

func InitAction(callbackId, userID string) []ebm.Attachment {
	// TODO: Update the timezone here
	loc, _ := time.LoadLocation("UTC")
	today := business_time.Today(loc)
	dms := adm.AdaptiveDynamicMenu(aesc.ProductionProfile, bindings)
	mog := dms.Build(userID, today)
	dms.StripOutFunctions().Evaluate(userID, business_time.Today(loc))

	attachAction1, _ := eb.NewAttachmentActionBuilder().
		Name("menu_list").
		Text("Pick an option...").
		ActionType(ebm.AttachmentActionTypeSelect).
		OptionGroups(mog).
		Build()

	attachAction2, _ := eb.NewAttachmentActionBuilder().
		Name("cancel").
		Text("I am good for now, thank you!").
		ActionType(models.ButtonType).
		Value(callbackId).
		Build()

	attach, _ := eb.NewAttachmentBuilder().
		Title(user.AdaptiveHiReply).
		Fallback("Adaptive at your service").
		Color(models.BlueColorHex).
		AttachmentType(models.DefaultAttachmentType).
		CallbackId(callbackId).
		Actions([]ebm.AttachmentAction{*attachAction1, *attachAction2}).
		Build()
	return []ebm.Attachment{*attach}
}

func publish(msg models.PlatformSimpleNotification) {
	_, err := sns.Publish(msg, platformNotificationTopic)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not publish message to %s topic", platformNotificationTopic))
}
func globalConnection(platformID models.PlatformID) daosCommon.DynamoDBConnection {
	return daosCommon.DynamoDBConnection{
		Dynamo: d,
		ClientID: clientID,
		PlatformID: platformID,
	}
}
func helloMessage(userID, channelID string, platformID models.PlatformID) {
	keyParams := map[string]*dynamodb.AttributeValue{
		"id": daosCommon.DynS(userID),
	}

	// Check if the user already exists
	var aUser models.User
	found, err := d.GetItemOrEmptyFromTable(usersTable, keyParams, &aUser)
	core.ErrorHandler(err, namespace, "Couldn't find user "+userID)
	// If the user doesn't exist in our tables, add the user first and then proceed to evaluate ADM
	if err == nil {
		if !found {
			log.Println("User does not exist, adding...")
			// refresh user cache
			engageUser, _ := json.Marshal(models.UserEngage{UserId: userID, PlatformID: models.PlatformID(platformID)})
			_, err = l.InvokeFunction(profileLambdaName, engageUser, false)
		}
	}

	rels := strategy.QueryCommunityUserIndex(userID, communityUsersTable, communityUsersUserIndex)
	if len(rels) > 0 {
		// api := slack.New(platformTokenDAO.GetPlatformTokenUnsafe(models.PlatformID(platformID)))
		//history, err2 := getChannelHistory(api, channelID)
		// if err2 != nil || !isThereVeryRecentHiResponse(history) {
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
			Attachments: InitAction("init_message", userID)})
		// }
		// if err2 != nil {
		// 	logger.WithError(err2).Errorf("Couldn't GetIMHistory from Slack")
		// }
	} else {
		// get the admin community
		dao := community.CommunityDAO(globalConnection(platformID), userCommunitiesTable)
		adminComms := dao.ReadOrEmptyUnsafe(platformID, string(community.Admin))
		if len(adminComms) == 0 {
			// if no admin community, post message to the user about that
			message := "Please ask your Slack administrator to finish setting up Adaptive by creating an Adaptive Admin private channel and then invite Adaptive to that channel."
			publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Message: message})
			return
		}
		// publish(models.PlatformSimpleNotification{UserId: slackMsg.User, Channel: slackMsg.Channel,
		//	Message: core.TextWrap("Sorry, you are not a member of user community yet. Please contact your
		//	admin to invite you to a community.", core.Underscore, core.Asterisk), AsUser: true})
		// When user is not a member of community, ask if the user wants to notify admin
		y, m := core.CurrentYearMonth()
		mc := models.MessageCallback{Module: "community", Source: userID, Topic: "admin", Action: "adaptive_access",
			Year: strconv.Itoa(y), Month: strconv.Itoa(m)}
		actions := []ebm.AttachmentAction{
			*models.SimpleAttachAction(mc, models.Now, user.NowActionLabel),
			*models.SimpleAttachAction(mc, models.Ignore, "Nah")}
		noAccessText, err := dialogFetcherDAO.FetchByContextSubject(NoAdaptiveAccessDialogContext, "text")
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not get dialog content for %s",
			NoAdaptiveAccessDialogContext))

		attach := utils.ChatAttachment("Sorry, you are not a member of any community yet :disappointed:",
			core.RandomString(noAccessText.Dialog), "", mc.ToCallbackID(), actions, []ebm.AttachmentField{}, time.Now().Unix())
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
			Attachments: []ebm.Attachment{*attach}})
	}
}

func forwardToNamespaceWithAppID(appID models.PlatformID, eventsAPIEvent string) func(string) {
	return func(namespace string) {
		np2 := models.NamespacePayload4{
			ID:        core.Uuid(),
			Namespace: namespace,
			PlatformRequest: models.PlatformRequest{
				PlatformID:   models.PlatformID(appID),
				SlackRequest: models.EventsAPIEvent(eventsAPIEvent),
			},
		}
		logger.WithField("platform_id", appID).Infof("Forwarding to %s", namespace)
		_, err := sns.Publish(np2, payloadTopicArn)
		core.ErrorHandler(err, namespace,
			fmt.Sprintf("2: Could not forward Slack event to topic=%s,  namespace=%s. requestPayload:\n%v",
				payloadTopicArn, namespace, eventsAPIEvent))
	}
}

func invokeLambdaWithAppID(appID models.PlatformID, eventsAPIEvent string) func(string) {
	return func(namespace string) {
		np2 := models.NamespacePayload4{
			ID:        core.Uuid(),
			Namespace: namespace,
			PlatformRequest: models.PlatformRequest{
				PlatformID:   appID,
				SlackRequest: models.EventsAPIEvent(eventsAPIEvent),
			},
		}
		fmt.Printf("INVOKING LAMBDA %v\n", namespace)
		bytes, err := json.Marshal(np2)
		lambdaName := fmt.Sprintf("%s_%s-%s", clientID, namespace, slackProcessorSuffix)
		_, err = l.InvokeFunction(lambdaName, bytes, true)
		// _, err := sns.Publish(np2, payloadTopicArn)
		core.ErrorHandler(err, namespace,
			fmt.Sprintf("2: Could not invoke Slack processor lambda for namespace=%s. requestPayload:\n%v", namespace, eventsAPIEvent))
	}
}

// HandleRequest is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	// IMPORTANT: It should always return this with the empty body. Else actions won't work.
	response = events.APIGatewayProxyResponse{StatusCode: 200}
	defer core.RecoverAsLogError("HandleRequest.Recover")
	logger = logger.WithLambdaContext(ctx)
	if request.HTTPMethod == "GET" {
		err = HandleRedirectURLGetRequest(globalConnection("UNKNOWN-PLATFORM-ID"), request)
		if err == nil {
			response.StatusCode = 308 // Permanent Redirect
			response.Headers = map[string]string{
				"Location": "https://adaptiveteam.github.io/",
			}
		} else {
			response.StatusCode = 500 // Server error
			response.Body = "Server error. See log for details"

		}
	} else {
		byt, _ := json.Marshal(request)
		logger.WithField("payload", string(byt)).Infof("Incoming gateway request")

		if request.Body != "" {
			// TODO: Refactor this condition. The problem is that slackevents.EventsAPIEvent doesn't have ApiAppId which is required (as PlatformID)
			if strings.Contains(request.Body, AppHomeOpened) {
				var ahe SlackAppHomeEvent
				err = json.Unmarshal([]byte(request.Body), &ahe)
				err = errors.Wrap(err, "Could not parse payload to AppHomeOpened")
				if err != nil {
					return
				}
				if ahe.Event.Type == AppHomeOpened {
					userID := ahe.Event.User
					channelID := ahe.Event.Channel
					// teamID := ahe.TeamID
					helloMessage(userID, channelID, models.PlatformID(ahe.ApiAppID))
				}
			} else {
				var requestPayload string
				requestPayload, err = getRequestPayload(request.Body)
				if err == nil {
					var eventsAPIEvent slackevents.EventsAPIEvent
					eventsAPIEvent, err = slackevents.ParseEvent(
						json.RawMessage(requestPayload),
						slackevents.OptionNoVerifyToken(),
					)
					// eventsAPIEvent.
					err = errors.Wrap(err, "Could not parse eventsAPIEvent")
					if err == nil {
						response, err = routeEventsAPIEvent(eventsAPIEvent, requestPayload)
					}
				}
			}
		}
	}
	err = errors.Wrap(err, "HandleRequest")
	if err != nil {
		logger.WithLambdaContext(ctx).WithError(err).Error("HandleRequest error")
	}
	err = nil
	logger.WithLambdaContext(ctx).Println("HandleRequest: normal termination. Returning `err=nil`")
	return
}

func getRequestPayload(requestBody string) (requestPayload string, err error) {
	if strings.HasPrefix(requestBody, "payload=%7B%22") {
		requestBody, err = url.QueryUnescape(requestBody)
		err = errors.Wrap(err, "Could not unescape gateway request")
	}
	requestPayload = strings.Replace(requestBody, "payload=", "", -1)
	return
}

func routeEventsAPIEvent(eventsAPIEvent slackevents.EventsAPIEvent,
	requestPayload string,
) (response events.APIGatewayProxyResponse, err error)  {
	logger.Infof("EVENT %v", eventsAPIEvent.Type)

	switch eventsAPIEvent.Type {
	case slackevents.AppHomeOpened:
		appHomeOpened := eventsAPIEvent.Data.(*slackevents.AppHomeOpenedEvent)
		helloMessage(appHomeOpened.User, appHomeOpened.Channel, "<UNKNOWN-PlatformID>")
	case slackevents.URLVerification:
		urlVerification := eventsAPIEvent.Data.(*slackevents.EventsAPIURLVerificationEvent)
		response, err = GwOk(urlVerification.Challenge, nil)
	case slackevents.CallbackEvent:
		callbackEvent := eventsAPIEvent.Data.(*slackevents.EventsAPICallbackEvent)
		response = events.APIGatewayProxyResponse{StatusCode: 200}
		err = routeCallbackEvent(eventsAPIEvent, requestPayload, *callbackEvent)
	default:
		logger.Infof("Handling event of type %v", eventsAPIEvent.Type)
		// workaround for slackevents.ParseEvent
		switch slack.InteractionType(eventsAPIEvent.Type) {
		case slack.InteractionTypeInteractionMessage:
			eventsAPIEvent.Data = utils.UnmarshallSlackInteractionCallbackUnsafe(requestPayload, namespace)
		case slack.InteractionTypeDialogSubmission:
			eventsAPIEvent.Data = utils.UnmarshallSlackInteractionCallbackUnsafe(requestPayload, namespace)
		case slack.InteractionTypeDialogCancellation:
			eventsAPIEvent.Data = utils.UnmarshallSlackInteractionCallbackUnsafe(requestPayload, namespace)
		default:
			panic(errors.New("Unknown type of Slack message: " + eventsAPIEvent.Type))
		}
		fmt.Printf("parsed eventsAPIEvent.Data =%v\n", reflect.TypeOf(eventsAPIEvent.Data))

		objMap := parseMapUnsafe(requestPayload)
		if _, ok := objMap["callback_id"]; ok {
			userID := getUserID(eventsAPIEvent)
			callbackID := getCallbackID(eventsAPIEvent)
			fmt.Printf("userID=%v,callbackID=%v\n", userID, callbackID)
			err = routeByCallbackID(eventsAPIEvent, requestPayload, userID, callbackID)
		}
	}
	return
}

func routeCallbackEvent(
	eventsAPIEvent slackevents.EventsAPIEvent,
	requestPayload string, 
	callbackEvent slackevents.EventsAPICallbackEvent,
) (err error)  {
	apiAppID := models.PlatformID(callbackEvent.APIAppID)
	forwardToNamespace := forwardToNamespaceWithAppID(apiAppID, requestPayload)
	invokeLambdaWithNamespace := invokeLambdaWithAppID(apiAppID, requestPayload)
	eventType := eventsAPIEvent.InnerEvent.Type
	fmt.Printf("INNEREVENT %v\n", eventType)
	if eventType == slackevents.AppMention {
		forwardToNamespace("community")
	} else if eventType == slackevents.MemberJoinedChannel {
		forwardToNamespace(AdaptiveChannelNamespace)
	} else if eventType == "member_left_channel" {
		forwardToNamespace(AdaptiveChannelNamespace)
	} else if eventType == "group_left" {
		forwardToNamespace(AdaptiveChannelNamespace)
	} else if eventType == "channel_deleted" {
		forwardToNamespace(AdaptiveChannelNamespace)
	} else if strings.HasPrefix(eventType, "channel_") {
		forwardToNamespace(AdaptiveChannelNamespace)
	} else if strings.HasPrefix(eventType, "app_") {
		forwardToNamespace(AdaptiveChannelNamespace)
	} else if strings.HasPrefix(eventType, "group_") {
		forwardToNamespace(AdaptiveChannelNamespace)
	} else if strings.HasPrefix(eventType, "member_") {
		forwardToNamespace(AdaptiveChannelNamespace)
	} else {
		slackMsg := utils.ParseAsCallbackMsg(eventsAPIEvent)
		slackText := core.TrimLower(slackMsg.Text)
		if slackText == "warmup" {
			// Warmed up lambda
			log.Println("Received warmup message ...")
		} else if slackText == "hello" || slackText == "hi" {
			helloMessage(slackMsg.User, slackMsg.Channel, apiAppID)
		} else if strings.Contains(slackText, "generate") ||
		strings.Contains(slackText, "add to slack") || 
		strings.Contains(slackText, "addtoslack") {
			GenerateAddToSlackURL(slackMsg.User, slackMsg.Channel, apiAppID)
		} else if slackMsg != nil {
			log.Println("### callback event: " + requestPayload)
			// It's not a response to an engagement, but a query
			if strings.Contains(requestPayload, "client_msg_id") {
				// We need to get the token for non-bot messages, so we keep the token retrieval inside the above condition
				if slackText == "help" {
					publish(models.PlatformSimpleNotification{UserId: slackMsg.User, Channel: slackMsg.Channel,
						Message: "Please refer to the adaptive documentation for available commands, " + helpPage, AsUser: true})
				} else if core.ListContainsString(settingsCommands, slackText) {
					forwardToNamespace("settings")
				} else if core.ListContainsString(feedbackCommands, slackText) {
					invokeLambdaWithNamespace("feedback")
				} else if slackText != "" {
					logger.WithField("slackText", slackText).Info("Unknown user command. Showing menu")
					helloMessage(slackMsg.User, slackMsg.Channel, apiAppID)
					// publish(models.PlatformSimpleNotification{UserId: slackMsg.User, Channel: slackMsg.Channel, Message: "Unable to process your message. Type `help` for instructions."})
				} else {
					logger.WithField("requestPayload", requestPayload).Info("Unknown request. Ignoring")
				}
			}
		}
	}
	return 
}
				
func routeByCallbackID(
	eventsAPIEvent slackevents.EventsAPIEvent,
	requestPayload string,
	userID, callbackID string,
) (err error) {
	slackRequest := models.EventsAPIEvent(requestPayload)
	u := userDAO.ReadUnsafe(userID)
	apiAppID := u.PlatformID
	platformID := u.PlatformID
	np := models.NamespacePayload4{
		ID:        core.Uuid(),
		Namespace: namespace,
		PlatformRequest: models.PlatformRequest{
			PlatformID:   apiAppID,
			SlackRequest: slackRequest,
		},
	}	
	forwardToNamespace := forwardToNamespaceWithAppID(apiAppID, requestPayload)
	invokeLambdaWithNamespace := invokeLambdaWithAppID(apiAppID, requestPayload)

	if strings.Contains(callbackID, "init_message") {
		if eventsAPIEvent.Type == string(slack.InteractionTypeInteractionMessage) {
			var message slack.InteractionCallback
			message, err = utils.ParseAsInteractionMsg(requestPayload)
			err = errors.Wrap(err, "Could not parse to interaction type message")
			if err != nil {
				return
			}
			action := message.ActionCallback.AttachmentActions[0]
			if action.Name == "menu_list" {
				selected := action.SelectedOptions[0]
				menuOption := selected.Value
				err = routeMenuOption(eventsAPIEvent, requestPayload, message, platformID,
					menuOption)
			} else if action.Name == "cancel" {
				deleteMessage(message)
			}
		}
	} else if strings.Contains(callbackID, "feedback") {
		invokeLambdaWithNamespace("feedback")
	} else if strings.Contains(callbackID, "user_settings") {
		forwardToNamespace("settings")
	} else if strings.Contains(callbackID, "objectives") {
		forwardToNamespace("objectives")
	} else if strings.Contains(callbackID, "strategy") {
		invokeLambdaWithNamespace("strategy")
	} else if strings.Contains(callbackID, "community") {
		forwardToNamespace("community")
	} else if strings.Contains(callbackID, "holidays") {
		holidaysLambda.LambdaRouting.HandleNamespacePayload4(np)
		// forwardToNamespace(HolidaysNamespace)
	} else if strings.Contains(callbackID, "adaptive_values") {
		forwardToNamespace(AdaptiveValuesNamespace)
	}
	return
}

func routeMenuOption(
	eventsAPIEvent slackevents.EventsAPIEvent, 
	requestPayload string,
	message slack.InteractionCallback,
	platformID models.PlatformID,
	menuOption string,
) (err error) {
	slackRequest := models.EventsAPIEvent(requestPayload)
	userID := getUserID(eventsAPIEvent)
	// callbackID := getCallbackID(eventsAPIEvent)
	// fmt.Printf("userID=%v,callbackID=%v\n", userID, callbackID)
	// u := userDAO.ReadUnsafe(userID)
	// apiAppID := u.PlatformID
	// platformID := u.PlatformID
	np := models.NamespacePayload4{
		ID:        core.Uuid(),
		Namespace: namespace,
		PlatformRequest: models.PlatformRequest{
			PlatformID:   platformID,
			SlackRequest: slackRequest,
		},
	}	
	forwardToNamespace := forwardToNamespaceWithAppID(platformID, requestPayload)
	invokeLambdaWithNamespace := invokeLambdaWithAppID(platformID, requestPayload)
	switch menuOption {
	case user.AskForEngagements:
		engage := models.UserEngageWithCheckValues{
			UserEngage: models.UserEngage{
				UserId: userID, IsNew: true, Update: true, OnDemand: true,
				ThreadTs: message.MessageTs, PlatformID: platformID,
			},
			CheckValues: checkValues(message.User.ID),
		}
		invokeLambdaUnsafe(engScriptingLambda, engage)
		deleteMessage(message)
	case user.UpdateSettings:
		forwardToNamespace("settings")
	case coaching.GiveFeedback, coaching.RequestFeedback, user.GenerateReport,
		user.FetchReport, coaching.ViewCoachees, coaching.ViewAdvocates:
		invokeLambdaWithNamespace("feedback")
	case objectives.CreateIDO, objectives.CreateIDONow, 
		user.StaleIDOsForMe,
		coaching.SelectCoachee, coaching.ReviewCoacheeProgressSelect,
		strategy.ViewCommunityAdvocateObjectives:
		forwardToNamespace("objectives")
	case coaching.RequestCoach, user.CurrentQuarterSchedule, user.NextQuarterSchedule,
		coaching.GenerateReportHR, coaching.FetchReportHR:
		forwardToNamespace("community")
	case strategy.CreateStrategyObjective, strategy.CreateFinancialObjective,
		strategy.CreateCustomerObjective, strategy.ViewStrategyObjectives,
		strategy.ViewAdvocacyObjectives,
		user.ViewObjectives,
		user.StaleObjectivesForMe:
		forwardToNamespace("objectives")
		// invokeLambdaWithNamespace("strategy")
	case strategy.CreateInitiative, 
		strategy.ViewCapabilityCommunityInitiatives,
		strategy.ViewAdvocacyInitiatives, 
		strategy.ViewInitiativeCommunityInitiatives,
		user.StaleInitiativesForMe:
		forwardToNamespace("objectives")
	case strategy.ViewCapabilityCommunityObjectives, 
		strategy.CreateVision, strategy.ViewVision, strategy.ViewEditVision,
		strategy.CreateCapabilityCommunity, strategy.ViewCapabilityCommunities,
		strategy.AssociateStrategyObjectiveToCapabilityCommunity,
		strategy.CreateInitiativeCommunity,
		strategy.AssociateInitiativeWithInitiativeCommunity:
		// forwardToNamespace("strategy")
		invokeLambdaWithNamespace("strategy")
	case SayHelloMenuItem:
		forwardToNamespace(HelloWorldNamespace)
	case 
		holidays.HolidaysListMenuItem, 
		holidays.HolidaysSimpleListMenuItem, 
		holidays.HolidaysCreateNewMenuItem:
		
		holidaysLambda.LambdaRouting.HandleNamespacePayload4(np)
		// forwardToNamespace(HolidaysNamespace)
	case values.AdaptiveValuesListMenuItem, 
		values.AdaptiveValuesSimpleListMenuItem,
		values.AdaptiveValuesCreateNewMenuItem:
		competencies.HandleNamespacePayload4(np)
		// forwardToNamespace(AdaptiveValuesNamespace)
	case "StrategyPerformanceReport":
		var buf *bytes.Buffer
		var reportname string
		buf, reportname, err = onStrategyPerformanceReport(ReadRDSConfigFromEnv(), platformID)
		if err == nil {
			err = sendReportToUser(platformID, userID, reportname, buf)
		}
		deleteMessage(message)
		err = errors.Wrap(err, "StrategyPerformanceReport")
	case "IDOPerformanceReport":
		var buf *bytes.Buffer
		var reportname string
		buf, reportname, err = onIDOPerformanceReport(ReadRDSConfigFromEnv(), userID)
		if err == nil {
			err = sendReportToUser(platformID, userID, reportname, buf)
		}
		deleteMessage(message)
		err = errors.Wrap(err, "IDOPerformanceReport")
	default:
		logger.Infof("Unknown/unhandled menu option '%s'", menuOption)
	}

	return
}

func deleteMessage(request slack.InteractionCallback) {
	publish(models.PlatformSimpleNotification{
		UserId:  request.User.ID,
		Channel: request.Channel.ID,
		Message: "", Ts: request.MessageTs, AsUser: true})
}
func parseMapUnsafe(input string) (objMap map[string]*json.RawMessage) {
	err := json.Unmarshal([]byte(input), &objMap)
	core.ErrorHandler(err, namespace, "Could not unmarshal json to map: "+input)
	return
}

func getCallbackID(eventsAPIEvent slackevents.EventsAPIEvent) string {
	return (eventsAPIEvent.Data.(slack.InteractionCallback)).CallbackID
}
func getUserID(eventsAPIEvent slackevents.EventsAPIEvent) string {
	return (eventsAPIEvent.Data.(slack.InteractionCallback)).User.ID
}

func invokeLambdaUnsafe(lambdaName string, userEngage models.UserEngageWithCheckValues) {
	engageBytes, err := json.Marshal(userEngage)
	core.ErrorHandler(err, namespace, "Could not marshal UserEngage")
	_, err = l.InvokeFunction(lambdaName, engageBytes, false)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not invoke %s", lambdaName))
}
