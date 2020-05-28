package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-checks"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/adaptiveteam/adaptive/daos/clientPlatformToken"
	"github.com/adaptiveteam/adaptive/daos/slackTeam"
	"log"
	"net/url"
	"time"
	"github.com/pkg/errors"
	adm "github.com/adaptiveteam/adaptive/adaptive-dynamic-menu"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
	"strings"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/holidays"
	competencies "github.com/adaptiveteam/adaptive/lambdas/competencies-lambda-go"
	holidaysLambda "github.com/adaptiveteam/adaptive/lambdas/holidays-lambda-go"
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
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	HelloWorldNamespace           = "hello-world"
	HolidaysNamespace             = "holidays"
	AdaptiveValuesNamespace       = "adaptive_values"
	AdaptiveChannelNamespace      = "adaptive-channel"
	NoAdaptiveAccessDialogContext = "dialog/engagements/adaptive-access/"
)

// -- Direct Adaptive to create an engagement alerting team members of upcoming holidays
// Enable users to request a list of all upcoming holidays

func loadProfile(conn daosCommon.DynamoDBConnection, userID string) (profile adaptive_checks.TypedProfile) {
	loc, _ := time.LoadLocation("UTC")
	today := business_time.Today(loc)
	profile = adaptive_checks.EvalProfile(conn, userID, today)
	profile.LoadAll()
	return
}

func InitAction(conn daosCommon.DynamoDBConnection, callbackId, userID string) []ebm.Attachment {
	// TODO: Update the timezone here
	// loc, _ := time.LoadLocation("UTC")
	// today := business_time.Today(loc)
	profile := loadProfile(conn, userID)
	mog := adm.AdaptiveDynamicMenu(profile, bindings).Build()
	// dms := adm.AdaptiveDynamicMenu(profile, bindings)
	// mog := dms.Build()//userID, today)
	// dms.StripOutFunctions().Evaluate(userID, business_time.Today(loc))

	logger.Infof("AdaptiveDynamicMenu contains %d groups\n", len(mog))
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
	_, err2 := sns.Publish(msg, platformNotificationTopic)
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not publish message to %s topic", platformNotificationTopic))
}

func helloMessage(userID, channelID string, teamID models.TeamID) {
	keyParams := map[string]*dynamodb.AttributeValue{
		"id": daosCommon.DynS(userID),
	}

	// Check if the user already exists
	var aUser models.User
	found, err2 := d.GetItemOrEmptyFromTable(usersTable, keyParams, &aUser)
	core.ErrorHandler(err2, namespace, "Couldn't find user "+userID)
	// If the user doesn't exist in our tables, add the user first and then proceed to evaluate ADM
	if !found {
		log.Println("User does not exist, adding...")
		// refresh user cache
		engageUser, _ := json.Marshal(models.UserEngage{UserID: userID, TeamID: models.TeamID(teamID)})
		_, err3 := l.InvokeFunction(profileLambdaName, engageUser, false)
		core.ErrorHandler(err3, namespace, "Couldn't add user "+userID)
	}

	rels := strategy.QueryCommunityUserIndex(userID, communityUsersTable, communityUsersUserIndex)
	if len(rels) > 0 {
		// api := slack.New(platformTokenDAO.GetPlatformTokenUnsafe(models.TeamID(teamID)))
		// history, err2 := getChannelHistory(api, channelID)
		// if err2 != nil || !isThereVeryRecentHiResponse(history) {
		conn := connGen.ForPlatformID(teamID.ToPlatformID())
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
			Attachments: InitAction(conn, "init_message", userID)})
		// }
		// if err2 != nil {
		// 	logger.WithError(err2).Errorf("Couldn't GetIMHistory from Slack")
		// }
	} else {
		// get the admin community
		conn := globalConnection(teamID)
		adminComms := adaptiveCommunity.ReadOrEmptyUnsafe(teamID.ToPlatformID(), string(community.Admin))(conn)
		var note models.PlatformSimpleNotification
		if len(adminComms) == 0 {
			// if no admin community, post message to the user about that
			message := "Please ask your Slack administrator to finish setting up Adaptive by creating an Adaptive Admin private channel and then inviting Adaptive to that channel."
			note = models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Message: message}
		} else {
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
			noAccessText, err4 := dialogFetcherDAO.FetchByContextSubject(NoAdaptiveAccessDialogContext, "text")
			core.ErrorHandler(err4, namespace, fmt.Sprintf("Could not get dialog content for %s",
				NoAdaptiveAccessDialogContext))

			attach := utils.ChatAttachment("Sorry, you are not a member of any community yet :disappointed:",
				core.RandomString(noAccessText.Dialog), "", mc.ToCallbackID(),
				actions, []ebm.AttachmentField{}, time.Now().Unix())
			note = models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
				Attachments: []ebm.Attachment{*attach}}
		}
		publish(note)
	}
}

func forwardToNamespaceWithAppID(appID models.TeamID, eventsAPIEvent string) func(string) {
	return func(namespace string) {
		np2 := models.NamespacePayload4{
			ID:        core.Uuid(),
			Namespace: namespace,
			PlatformRequest: models.PlatformRequest{
				TeamID:       appID,
				SlackRequest: models.EventsAPIEvent(eventsAPIEvent),
			},
		}
		logger.WithField("platform_id", appID).Infof("Forwarding to %s", namespace)
		_, err2 := sns.Publish(np2, payloadTopicArn)
		core.ErrorHandler(err2, namespace,
			fmt.Sprintf("2: Could not forward Slack event to topic=%s,  namespace=%s. requestPayload:\n%v",
				payloadTopicArn, namespace, eventsAPIEvent))
	}
}

func invokeLambdaWithAppID(appID models.TeamID, eventsAPIEvent string) func(string) {
	return func(namespace string) {
		np2 := models.NamespacePayload4{
			ID:        core.Uuid(),
			Namespace: namespace,
			PlatformRequest: models.PlatformRequest{
				TeamID:       appID,
				SlackRequest: models.EventsAPIEvent(eventsAPIEvent),
			},
		}
		fmt.Printf("INVOKING LAMBDA %v\n", namespace)
		bytes, err2 := json.Marshal(np2)
		lambdaName := fmt.Sprintf("%s_%s-%s", clientID, namespace, slackProcessorSuffix)
		if err2 == nil {
			_, err2 = l.InvokeFunction(lambdaName, bytes, true)
		}
		core.ErrorHandler(err2, namespace,
			fmt.Sprintf("2: Could not invoke Slack processor lambda for namespace=%s. requestPayload:\n%v", namespace, eventsAPIEvent))
	}
}

// HandleRequest is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	defer core.RecoverAsLogError("HandleRequest.Recover")
	logger = logger.WithLambdaContext(ctx)
	if request.HTTPMethod == "GET" {
		err = HandleRedirectURLGetRequest(globalConnection(models.ParseTeamID("UNKNOWN-PLATFORM-ID")), request)
		if err == nil {
			response = responsePermanentRedirect("https://www.adaptive.team/welcome-to-adaptive")
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
				if err == nil {
					if ahe.Event.Type == AppHomeOpened {
						teamID := models.TeamID{
							TeamID: daosCommon.PlatformID(ahe.TeamID),
							AppID:  daosCommon.PlatformID(ahe.ApiAppID),
						}
						logger.Infof("AppHomeOpened, teamID=%v", teamID)
					}
					response = responseOk
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
					err = errors.Wrap(err, "Could not parse eventsAPIEvent")
					if err == nil {
						response, err = routeEventsAPIEvent(eventsAPIEvent, requestPayload)
					}
				}
			}
		} else {
			err = errors.New("An empty request.Body is not supported")
		}
	}
	err = errors.Wrap(err, "HandleRequest")
	if err != nil {
		logger.WithLambdaContext(ctx).WithError(err).Error("HandleRequest error")
		response = responseServerError(err)
	}
	err = nil
	logger.
		WithLambdaContext(ctx).
		Printf("HandleRequest: normal termination. Returning `err=nil`. response: %v\n", response)
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
) (response events.APIGatewayProxyResponse, err error) {
	logger.Infof("EVENT eventsAPIEvent.Type=%v", eventsAPIEvent.Type)
	response = responseOk
	switch eventsAPIEvent.Type {
	case slackevents.AppHomeOpened:
		appHomeOpened := eventsAPIEvent.Data.(*slackevents.AppHomeOpenedEvent)
		logger.Infof("UNKNOWN-PlatformID")
		helloMessage(appHomeOpened.User, appHomeOpened.Channel,
			models.ParseTeamID("<UNKNOWN-PlatformID>"))
	case slackevents.URLVerification:
		urlVerification := eventsAPIEvent.Data.(*slackevents.EventsAPIURLVerificationEvent)
		response = responseOkBody(urlVerification.Challenge)
	case slackevents.CallbackEvent:
		callbackEvent := eventsAPIEvent.Data.(*slackevents.EventsAPICallbackEvent)
		if callbackEvent.TeamID == "xxxxx" && callbackEvent.APIAppID == "xxxxx" {
			log.Println("Received warmup message with TeamID=xxxxx and AppID=xxxxx ...")
		} else {
			err = routeCallbackEvent(eventsAPIEvent, requestPayload, *callbackEvent)
		}
	default:
		// workaround for slackevents.ParseEvent
		var slackInteractionCallback slack.InteractionCallback
		switch slack.InteractionType(eventsAPIEvent.Type) {
		case slack.InteractionTypeInteractionMessage, // "interactive_message"
			slack.InteractionTypeDialogSubmission,
			slack.InteractionTypeDialogCancellation: 
			slackInteractionCallback = utils.UnmarshallSlackInteractionCallbackUnsafe(requestPayload, namespace)
			data, _ := json.Marshal(slackInteractionCallback)
			fmt.Printf("Parsed slackInteractionCallback:\n%+v", string(data))
			userID := slackInteractionCallback.User.ID
			callbackID := slackInteractionCallback.CallbackID

			var teamID models.TeamID
			teamID, err = ensureTeamID(
				daosCommon.PlatformID(slackInteractionCallback.Team.ID), 
				daosCommon.PlatformID(slackInteractionCallback.APIAppID),
			)
			if err == nil {
				err = routeByCallbackID(slackInteractionCallback, requestPayload, teamID, userID, callbackID)
			}
		default:
			panic(errors.New("Unknown type of Slack message: " + eventsAPIEvent.Type))
		}
	}
	return
}

// ensureTeamID reads tokens and returns correct team id.
// If we have teamID token, then we use it. Otherwise we use appID
func ensureTeamID(teamID, appID daosCommon.PlatformID) (res models.TeamID, err error) {
	var teams []slackTeam.SlackTeam
	if teamID == "" {
		fmt.Printf("TeamID is empty")
	} else {
		teams, err = slackTeam.ReadOrEmpty(teamID)(connGen.ForPlatformID(teamID))
	}
	if err == nil {
		if len(teams) > 0 {
			res = models.TeamID{TeamID: teams[0].TeamID}
		} else {
			var clientConfigs []clientPlatformToken.ClientPlatformToken
			if teamID == "" {
				err = errors.Errorf("AppID is empty: teamID=%s or appID=%s", teamID, appID)
			} else {
				clientConfigs, err = clientPlatformToken.ReadOrEmpty(appID)(connGen.ForPlatformID(appID))
			}
			if err == nil {
				if len(clientConfigs) > 0 {
					res = models.TeamID{AppID: clientConfigs[0].PlatformID}
				} else {
					err = errors.Errorf("Couldn't find teamID=%s or appID=%s", teamID, appID)
				}
			}
		}
	}
	log.Printf("ensureTeamID: %v\n", teamID)
	return
}

func routeCallbackEvent(
	eventsAPIEvent slackevents.EventsAPIEvent,
	requestPayload string,
	callbackEvent slackevents.EventsAPICallbackEvent,
) (err error) {
	var teamID models.TeamID
	teamID, err = ensureTeamID(
		daosCommon.PlatformID(callbackEvent.TeamID),
		daosCommon.PlatformID(callbackEvent.APIAppID),
	)
	if err == nil {
		forwardToNamespace := forwardToNamespaceWithAppID(teamID, requestPayload)
		invokeLambdaWithNamespace := invokeLambdaWithAppID(teamID, requestPayload)
		eventType := eventsAPIEvent.InnerEvent.Type
		logger.WithField("InnerEvent.Type", eventType).Info("routeCallbackEvent")
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
			if slackMsg != nil {
				slackText := core.TrimLower(slackMsg.Text)
				if slackText == "warmup" { // this case is not reachable, because we checked for warmup message earlier
					// Warmed up lambda
					log.Println("Received warmup message (WARN: shouldn't reach this point ever) ...")
					err = nil // because teamID is not correct when it's warmup
				} else if slackText == "hello" || slackText == "hi" {
					logger.Infof("apiAppID: %v", teamID)
					helloMessage(slackMsg.User, slackMsg.Channel, teamID)
				} else if strings.Contains(slackText, "generate") ||
					strings.Contains(slackText, "add to slack") ||
					strings.Contains(slackText, "addtoslack") {
					GenerateAddToSlackURL(slackMsg.User, slackMsg.Channel, teamID)
				} else {
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
							helloMessage(slackMsg.User, slackMsg.Channel, teamID)
							// publish(models.PlatformSimpleNotification{UserId: slackMsg.User, Channel: slackMsg.Channel, Message: "Unable to process your message. Type `help` for instructions."})
						} else {
							logger.WithField("requestPayload", requestPayload).Info("Unknown request. Ignoring")
						}
					}
				}
			}
		}
	}
	return
}

func routeByCallbackID(
	slackInteractionCallback slack.InteractionCallback,
	requestPayload string,
	teamID models.TeamID,
	userID, callbackID string,
) (err error) {
	fmt.Printf("routeByCallbackID(userID=%s,callbackID=%s)\n", userID, callbackID)

	slackRequest := models.EventsAPIEvent(requestPayload)
	np := models.NamespacePayload4{
		ID:        core.Uuid(),
		Namespace: namespace,
		PlatformRequest: models.PlatformRequest{
			TeamID:       teamID,
			SlackRequest: slackRequest,
		},
	}
	forwardToNamespace := forwardToNamespaceWithAppID(teamID, requestPayload)
	invokeLambdaWithNamespace := invokeLambdaWithAppID(teamID, requestPayload)

	if strings.Contains(callbackID, "init_message") {
		// if eventsAPIEvent.Type == string(slack.InteractionTypeInteractionMessage) {
			// var message slack.InteractionCallback
			// message, err = utils.ParseAsInteractionMsg(requestPayload)
			// err = errors.Wrap(err, "Could not parse to interaction type message")
			// logger.Infof("init_message parsed: %v", message)
			// if err != nil {
			// 	return
			// }
			message := slackInteractionCallback
			action := message.ActionCallback.AttachmentActions[0]
			if action.Name == "menu_list" {
				selected := action.SelectedOptions[0]
				menuOption := selected.Value
				err = routeMenuOption(slackInteractionCallback.User.ID, requestPayload, message, teamID,
					menuOption)
			} else if action.Name == "cancel" {
				deleteMessage(message)
			}
		// }
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
	} else if strings.Contains(callbackID, AdaptiveValuesNamespace) {
		competencies.HandleNamespacePayload4(np)
	} else {
		fmt.Printf("Unhandled callbackID=%s", callbackID)
	}
	return
}

func routeMenuOption(
	userID string,
	requestPayload string,
	message slack.InteractionCallback,
	teamID models.TeamID,
	menuOption string,
) (err error) {
	logger.WithField("menuOption", menuOption).Infof("Routing menu option")
	slackRequest := models.EventsAPIEvent(requestPayload)
	np := models.NamespacePayload4{
		ID:        core.Uuid(),
		Namespace: namespace,
		PlatformRequest: models.PlatformRequest{
			TeamID:       teamID,
			SlackRequest: slackRequest,
		},
	}
	forwardToNamespace := forwardToNamespaceWithAppID(teamID, requestPayload)
	invokeLambdaWithNamespace := invokeLambdaWithAppID(teamID, requestPayload)
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	switch menuOption {
	case user.AskForEngagements:
		engage := models.UserEngage{
				UserID: userID, IsNew: true, Update: true, OnDemand: true,
				ThreadTs: message.MessageTs, TeamID: teamID,
		}
		invokeLambdaUnsafe(engScriptingLambda, engage)
		deleteMessage(message)
	case user.UpdateSettings:
		forwardToNamespace("settings")
	case coaching.GiveFeedback, coaching.RequestFeedback, user.GenerateReport,
		user.FetchReport, coaching.ViewCoachees, coaching.ViewAdvocates:
		invokeLambdaWithNamespace("feedback")
	case // workflows that are handled in user-objectives-lambda
		objectives.CreateIDO,
		objectives.CreateIDONow,
		strategy.CreateStrategyObjective, 
		strategy.CreateFinancialObjective,
		strategy.CreateCustomerObjective, 
		strategy.CreateInitiative,

		coaching.SelectCoachee,
		coaching.ReviewCoacheeProgressSelect,
		strategy.ViewCapabilityCommunityInitiatives,
		strategy.ViewAdvocacyInitiatives,
		strategy.ViewInitiativeCommunityInitiatives,
		strategy.ViewCommunityAdvocateObjectives,
		strategy.ViewStrategyObjectives,
		strategy.ViewAdvocacyObjectives,
		user.ViewObjectives,
		user.StaleIDOsForMe,
		user.StaleObjectivesForMe,
		user.StaleInitiativesForMe:
		forwardToNamespace("objectives")
	case coaching.RequestCoach, 
		user.CurrentQuarterSchedule, 
		user.NextQuarterSchedule,
		coaching.GenerateReportHR, coaching.FetchReportHR:
		forwardToNamespace("community")
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
		buf, reportname, err = onStrategyPerformanceReport(ReadRDSConfigFromEnv(), teamID)
		if err == nil {
			err = sendReportToUser(teamID, userID, reportname, buf, conn)
		}
		deleteMessage(message)
		err = errors.Wrap(err, "StrategyPerformanceReport")
	case "IDOPerformanceReport":
		var buf *bytes.Buffer
		var reportname string
		buf, reportname, err = onIDOPerformanceReport(ReadRDSConfigFromEnv(), userID)
		if err == nil {
			err = sendReportToUser(teamID, userID, reportname, buf, conn)
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
		Message: "", Ts: request.MessageTs, 
	})
}
func parseMapUnsafe(input string) (objMap map[string]*json.RawMessage) {
	err2 := json.Unmarshal([]byte(input), &objMap)
	core.ErrorHandler(err2, namespace, "Could not unmarshal json to map: "+input)
	return
}

func getCallbackID(eventsAPIEvent slackevents.EventsAPIEvent) string {
	return (eventsAPIEvent.Data.(slack.InteractionCallback)).CallbackID
}
func getUserID(eventsAPIEvent slackevents.EventsAPIEvent) string {
	return (eventsAPIEvent.Data.(slack.InteractionCallback)).User.ID
}

func invokeLambdaUnsafe(lambdaName string, userEngage models.UserEngage) {
	engageBytes, err2 := json.Marshal(userEngage)
	core.ErrorHandler(err2, namespace, "Could not marshal UserEngage")
	_, err2 = l.InvokeFunction(lambdaName, engageBytes, false)
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not invoke %s", lambdaName))
}
