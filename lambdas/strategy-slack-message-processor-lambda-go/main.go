package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ObjectiveEvent          = "objective"
	FinancialObjectiveEvent = "financial_objective"
	CustomerObjectiveEvent  = "customer_objective"
	ObjectiveAdhocEvent     = "objective_adhoc"
	// ObjectiveAssociationAdhoc                      = "objective_association_adhoc"
	VisionEvent                                  = "vision"
	CapabilityCommunityEvent                     = "capability_community"
	InitiativeAdhocEvent                         = "initiative_adhoc"
	InitiativeSelectCommunityEvent               = "initiative_select_community"
	ObjectiveSelectCommunityEvent                = "objective_select_community"
	ObjectiveCommunityAssociationSelectObjective = "objective_community_association_select_objective"
	// ObjectiveCommunityAssociationSelectObjectiveAndCommunity = "objective_community_association_select_objective_community"

	InitiativeCommunityEvent                       = "initiative_community"
	AssociateObjectiveWithCapabilityCommunityEvent = "associate_objective_capability_community"

	CreateStrategyObjective = strategy.CreateStrategyObjective // "create_strategy_objective"
	// CreateStrategyCommunityAssociation         = "create_strategy_community_association"
	CreateCustomerObjective                    = strategy.CreateCustomerObjective  // "create_customer_objective"
	CreateFinancialObjective                   = strategy.CreateFinancialObjective // "create_financial_objective"
	ViewAdvocacyObjectives                     = strategy.ViewAdvocacyObjectives   // "view_strategy_objectives"
	ViewStrategyObjectives                     = strategy.ViewStrategyObjectives   // "view_strategy_objectives"
	ViewCapabilityCommunityObjectives          = strategy.ViewCapabilityCommunityObjectives
	CreateCapabilityCommunity                  = strategy.CreateCapabilityCommunity // "create_capability_community"
	ViewCapabilityCommunities                  = strategy.ViewCapabilityCommunities // "view_capability_communities"
	ViewCapabilityCommunityInitiatives         = strategy.ViewCapabilityCommunityInitiatives
	ViewInitiativeCommunityInitiatives         = strategy.ViewInitiativeCommunityInitiatives
	CreateInitiative                           = strategy.CreateInitiative                                // "create_initiative"
	CreateInitiativeCommunity                  = strategy.CreateInitiativeCommunity                       // "create_initiative_community"
	AssociateInitiativeWithInitiativeCommunity = strategy.AssociateInitiativeWithInitiativeCommunityEvent // "associate_initiative_with_initiative_community"
	InitiativeCommunityAssociationSelectInitiative = "initiative_community_association_select_initiative"

	AssociateStrategyObjectiveToCapabilityCommunity = strategy.AssociateStrategyObjectiveToCapabilityCommunity

	CreateActionName = fmt.Sprintf("%s_%s", strategy.Create, models.Now)

	visionDescriptionContext         = "dialog/strategy/language-coaching/vision"
	stratObjDescriptionContext       = "dialog/strategy/language-coaching/objective/description"
	objAllocationContext             = "dialog/strategy/language-coaching/objective-allocation"
	capCommunityDescriptionContext   = "dialog/strategy/language-coaching/community/capability"
	initiativeDescriptionContext     = "dialog/strategy/language-coaching/initiative/description"
	initiativeCommDescriptionContext = "dialog/strategy/language-coaching/community/initiative"

	logger = alog.LambdaLogger(logrus.InfoLevel)

	isMemberInCommunity = func(userID string, comm community.AdaptiveCommunity) bool {
		return community.IsUserInCommunity(userID, communityUsersTable, communityUsersUserCommunityIndex, comm)
	}
)

type MsgState struct {
	ThreadTs       string `json:"thread_ts"`
	Update         bool   `json:"update"`
	Id             string `json:"id"`
	SelectedOption string `json:"selected_option"`
}

func communityMembers(commID string, teamID models.TeamID) []models.KvPair {
	// Get coaching community members
	commMembers := community.CommunityMembers(communityUsersTable, commID, teamID, communityUsersCommunityIndex)
	logger.Infof("Members in %s Community for %s platform: %s", commID, teamID, commMembers)
	var users []models.KvPair
	for _, each := range commMembers {
		// Self user checking
		u := userDAO.ReadUnsafe(each.UserId)
		if u.DisplayName != "" && !u.IsAdaptiveBot {
			users = append(users, models.KvPair{Key: u.DisplayName, Value: each.UserId})
		}
	}
	logger.Infof("KvPairs from communities for %s community for %s platform: %s", commID, teamID, users)
	return users
}

func CommunityById(communityID string, teamID models.TeamID) models.AdaptiveCommunity {
	return community.CommunityById(communityID, teamID, communitiesTable)
}

// StrategyEntityById - a weird way to read entities.
// Deprecated: should use typesafe DAO.Read
func StrategyEntityById(id string, teamID models.TeamID, table string) interface{} {
	params := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
		"platform_id": {
			S: aws.String(teamID.ToString()),
		},
	}
	var comm interface{}
	err2 := d.GetItemFromTable(table, params, &comm)
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not query %s table", communitiesTable))
	return comm
}

func AllStrategyCommunities(teamID models.TeamID) []strategy.StrategyCommunity {
	var scs []strategy.StrategyCommunity
	err := d.QueryTableWithIndex(strategyCommunitiesTable, awsutils.DynamoIndexExpression{
		IndexName: strategyCommunitiesPlatformIndex,
		Condition: "platform_id = :p",
		Attributes: map[string]interface{}{
			":p": teamID,
		},
	}, map[string]string{}, true, -1, &scs)
	if err != nil {
		logger.Errorf("Could not query %s index on %s table", strategyCommunitiesPlatformIndex,
			strategyCommunitiesTable)
	} else {
		logger.Infof("Queried all strategy communities for %s platform: %v", teamID, scs)
	}
	return scs
}

func dialogFromSurvey(teamID models.TeamID, userID string, message slack.InteractionCallback, survey ebm.AttachmentActionSurvey, callbackID string,
	contextId string, update bool, selected string) (err error) {
	token := platformDAO.GetPlatformTokenUnsafe(teamID)
	if token != "" {
		api := slack.New(token)
		survState := func() string {
			// When the original message is from a thread, we need to post to the same thread
			// Below logic checks if the incoming message is from a thread
			var ts string
			if message.OriginalMessage.ThreadTimestamp == "" {
				ts = message.MessageTs
			} else {
				ts = message.OriginalMessage.ThreadTimestamp
			}
			msgStateBytes, err := json.Marshal(MsgState{ThreadTs: ts, Update: update, Id: contextId, SelectedOption: selected})
			core.ErrorHandler(err, namespace, "Could not marshal MsgState")
			return string(msgStateBytes)
		}
		err = utils.SlackSurvey(api, message, survey, callbackID, survState)
	} else {
		logger.Errorf("Platform token is empty for %s user", userID)
	}
	return
}

func dialogFromSurveyUnsafe(teamID models.TeamID, userID string, message slack.InteractionCallback, survey ebm.AttachmentActionSurvey, callbackID string,
	contextId string, update bool, selected string) {
	err := dialogFromSurvey(teamID, userID, message, survey, callbackID, contextId, update, selected)
	if err != nil {
		logger.Errorf("Could not open dialog from %s survey", callbackID)
		responseOnDialogOpenFailure(userID)
	} else {
		logger.Infof("Opened dialog from survey for %s id", callbackID)
	}
}

func publish(msg models.PlatformSimpleNotification) {
	_, err := s.Publish(msg, platformNotificationTopic)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not publish message to %s topic", platformNotificationTopic))
}

func publishAll(messages []models.PlatformSimpleNotification) {
	for _, msg := range messages {
		publish(msg)
	}
}

func responses(notifications ...models.PlatformSimpleNotification) []models.PlatformSimpleNotification {
	return notifications
}

func respond(teamID models.TeamID, response platform.Response) {
	fmt.Printf("Respond(,%v): %s", response, response.Type)
	presp := platform.TeamResponse{
		TeamID:   teamID,
		Response: response,
	}
	_, err := s.Publish(presp, platformNotificationTopic)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not publish message to %s topic", platformNotificationTopic))
}

// used in many places
func DeleteOriginalEng(userId, channel, ts string) {
	utils.DeleteOriginalEng(userId, channel, ts, func(notification models.PlatformSimpleNotification) {
		publish(notification)
	})
}

func getUpdateParams(actionName string) (string, bool) {
	var act string
	var update bool
	if strings.HasPrefix(actionName, strategy.CreatePrefix) {
		act = strings.TrimPrefix(actionName, strategy.CreatePrefix)
		update = false
	} else if strings.HasPrefix(actionName, strategy.UpdatePrefix) {
		act = strings.TrimPrefix(actionName, strategy.UpdatePrefix)
		update = true
	}
	return act, update
}

func communityMembersIncludingStrategyMembers(commID string, teamID models.TeamID) []models.KvPair {
	// Strategy Community members
	strategyCommMembers := communityMembers(string(community.Strategy), teamID)
	commMembers := communityMembers(commID, teamID)
	return models.DistinctKvPairs(append(strategyCommMembers, commMembers...))
}

func handleObjectiveCreate(mc models.MessageCallback, actionName, actionValue string, userID, channelID string,
	teamID models.TeamID, message slack.InteractionCallback, typ models.StrategyObjectiveType) {
	if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
		// create objective actions
		if strings.HasPrefix(actionName, strategy.CreatePrefix) || strings.HasPrefix(actionName, strategy.UpdatePrefix) {
			act, update := getUpdateParams(actionName)
			switch act {
			case string(models.Now):
				var so models.StrategyObjective
				if typ == models.CapabilityStrategyObjective {
					selected := actionSelected(message.ActionCallback.AttachmentActions)
					capCommID := core.IfThenElse(selected != core.EmptyString, selected, mc.Target).(string)

					if update {
						so = strategy.StrategyObjectiveByID(models.TeamID(teamID), mc.Target, strategyObjectivesTable)
						if !AuthorizedForObjectiveEdit(userID, teamID, &so) {
							logger.Infof("%s user is no longer authorized to edit %s Objective Community", userID, capCommID)
							PostMsgToUser(fmt.Sprintf("You are no longer authorized to edit this Objective Community"),
								userID, channelID, message.MessageTs)
							return
						}
						// During update, original SO id is passed into target. We get cap comm id from existing record
						capCommID = so.CapabilityCommunityIDs[0]
					}
					strategyComm := StrategyCommunityByID(capCommID)
					// Check if the community is still associated with Adaptive
					if strategyComm.ChannelCreated == 1 {
						allMembers := communityMembersIncludingStrategyMembers(fmt.Sprintf("%s:%s", community.Capability, capCommID), teamID)
						if len(allMembers) > 0 {
							objectiveSurveyElements := EditSObjectiveSurveyElems(&so, ObjectiveTypes(), allMembers,
								objectives.StrategyObjectiveDatesWithIndefiniteOption(namespace, so.ExpectedEndDate))
							logger.Infof("Survey elements for Capability Objective for %s user for %s platform: %v.", userID, teamID, objectiveSurveyElements)
							val := utils.AttachmentSurvey(fmt.Sprintf("%s Objective", models.CapabilityStrategyObjective),
								objectiveSurveyElements)
							// Open a survey associated with the engagement
							dialogFromSurveyUnsafe(teamID, userID, message, val, actionValue, mc.Target, update, capCommID)
						} else {
							publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
								Message: "There are no advocates to associate the objective with. It could be that you haven't subscribed Adaptive to a community you are an advocate for."})
						}
					} else {
						logger.Infof("%s Objective Community is no longer associated with Adaptive", capCommID)
						PostMsgToUser(fmt.Sprintf("Selected community is no longer associated with Adaptive"), userID, channelID, message.MessageTs)
					}
				} else if typ == models.FinancialStrategyObjective || typ == models.CustomerStrategyObjective {
					if update {
						so = strategy.StrategyObjectiveByID(models.TeamID(teamID), mc.Target, strategyObjectivesTable)
						// During update, original SO id is passed into target. We get cap comm id from existing record
						// capCommID = so.CapabilityCommunityID
					}
					val := utils.AttachmentSurvey(fmt.Sprintf("%s Objective", typ), EditSObjectiveSurveyElems(&so,
						ObjectiveTypes(), communityMembers(string(community.Strategy), teamID), objectives.StrategyObjectiveDatesWithIndefiniteOption(namespace, so.ExpectedEndDate)))
					// Open a survey associated with the engagement
					dialogFromSurveyUnsafe(teamID, userID, message, val, actionValue, mc.Target, update, core.EmptyString)
				}
			case string(models.Ignore):
				DeleteOriginalEng(userID, channelID, message.MessageTs)
			}
		}
	}
}

func responseOnDialogOpenFailure(userID string) platform.Response {
	return platform.Post(platform.ConversationID(userID),
		platform.MessageContent{Message: ErrorOnDialogOpenMessage})
}

func ViewVisionMissionAttachment(userID string, teamID models.TeamID, vm *models.VisionMission,
	mc models.MessageCallback) []ebm.Attachment {
	// TODO: Currently only strategy community members to update vision
	var members []string
	for _, each := range communityMembers(string(community.Strategy), teamID) {
		members = append(members, each.Value)
	}
	enableActions := core.IfThenElse(core.ListContainsString(members, userID), true, false).(bool)
	return VisionMissionViewAttachment(*mc.WithTopic(VisionEvent), vm, nil, enableActions, false)
}

func handleMissionVisionCreate(mc models.MessageCallback, actionName, actionValue, userID, channelID string,
	teamID models.TeamID, message slack.InteractionCallback) {
	if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
		// create objective actions
		if strings.HasPrefix(actionName, strategy.CreatePrefix) || strings.HasPrefix(actionName, strategy.UpdatePrefix) {
			act, update := getUpdateParams(actionName)
			switch act {
			case string(models.Now):
				var vm models.VisionMission
				if update {
					vm = *StrategyVision(models.TeamID(teamID))
				}
				advocates := communityMembers(string(community.Strategy), teamID)
				val := utils.AttachmentSurvey(string(VisionLabel), EditVisionMissionSurveyElems(&vm, advocates))
				// Open a survey associated with the engagement
				dialogFromSurveyUnsafe(teamID, userID, message, val, actionValue, mc.Target, update, core.EmptyString)
			case string(models.Ignore):
				utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
				DeleteOriginalEng(userID, channelID, message.MessageTs)
			}
		}
	}
}

func handleCapabilityCommunityCreate(mc models.MessageCallback, actionName, actionValue, userID, channelID string,
	message slack.InteractionCallback, teamID models.TeamID) {
	if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
		// create objective actions
		if strings.HasPrefix(actionName, strategy.CreatePrefix) || strings.HasPrefix(actionName, strategy.UpdatePrefix) {
			act, update := getUpdateParams(actionName)
			switch act {
			case string(models.Now):
				var cc strategy.CapabilityCommunity
				if update {
					cc = strategy.CapabilityCommunityByID(teamID, mc.Target, capabilityCommunitiesTable)
				}
				val := utils.AttachmentSurvey("Objective Community", EditCapabilityCommunitySurveyElems(teamID, &cc))
				// Open a survey associated with the engagement
				dialogFromSurveyUnsafe(teamID, userID, message, val, actionValue, mc.Target, update, core.EmptyString)
			case string(models.Ignore):
				utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
				DeleteOriginalEng(userID, channelID, message.MessageTs)
			}
		}
	}
}

func handleInitiativeCreate(mc models.MessageCallback, actionName, actionValue, userID, channelID string,
	teamID models.TeamID, message slack.InteractionCallback) {
	if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
		// create objective actions
		if strings.HasPrefix(actionName, strategy.CreatePrefix) || strings.HasPrefix(actionName, strategy.UpdatePrefix) {
			act, update := getUpdateParams(actionName)
			switch act {
			case string(models.Now):
				var si models.StrategyInitiative
				selected := actionSelected(message.ActionCallback.AttachmentActions)
				initCommID := core.IfThenElse(selected != core.EmptyString, selected, mc.Target).(string)
				if update {
					si = strategy.StrategyInitiativeByID(teamID, mc.Target, strategyInitiativesTable)
					if !AuthorizedForInitiativeAddEdit(userID, teamID, &si) {
						logger.Infof("%s user is no longer authorized to edit %s Initiative Community", userID, initCommID)
						PostMsgToUser(fmt.Sprintf("You are no longer authorized to edit this Initiative Community"),
							userID, channelID, message.MessageTs)
						return
					}
					// During update, original SI id is passed into target. We get init comm id from existing record
					initCommID = si.InitiativeCommunityID
				}
				strategyComm := StrategyCommunityByID(initCommID)
				// Check if the community is still associated with Adaptive
				if strategyComm.ChannelCreated == 1 {
					commMembers := communityMembers(fmt.Sprintf("%s:%s", community.Initiative, initCommID), teamID)
					logger.Infof("Community members for creating an initiative: %s", commMembers)
					objs := UserCommunityObjectives(userID)
					var objKVs []models.KvPair
					for _, eachObj := range objs {
						objKVs = append(objKVs, models.KvPair{Key: eachObj.Name, Value: eachObj.ID})
					}
					initiativeSurveyElements := EditInitiativeSurveyElems(&si, commMembers,
						objectives.StrategyObjectiveDates(namespace, si.ExpectedEndDate), objKVs)
					logger.Infof("Initiative Survey Elements for %s user in %s platform: %s", userID, teamID,
						initiativeSurveyElements)
					val := utils.AttachmentSurvey("Strategy Initiative", initiativeSurveyElements)
					dialogFromSurveyUnsafe(teamID, userID, message, val, actionValue, mc.Target, update, initCommID)
				} else {
					logger.Infof("%s Initiative Community is no longer associated with Adaptive", initCommID)
					PostMsgToUser(fmt.Sprintf("Selected community is no longer associated with Adaptive"), userID, channelID, message.MessageTs)
				}
			case string(models.Ignore):
				utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
				DeleteOriginalEng(userID, channelID, message.MessageTs)
			}
		}
	}
}

func handleInitiativeCommunityCreate(mc models.MessageCallback, actionName, actionValue, userID, channelID string,
	message slack.InteractionCallback, teamID models.TeamID) {
	if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
		// create objective actions
		if strings.HasPrefix(actionName, strategy.CreatePrefix) || strings.HasPrefix(actionName, strategy.UpdatePrefix) {
			act, update := getUpdateParams(actionName)
			switch act {
			case string(models.Now):
				var sic strategy.StrategyInitiativeCommunity
				if update {
					sic = strategy.InitiativeCommunityByID(teamID, mc.Target, strategyInitiativeCommunitiesTable)
				}
				val := utils.AttachmentSurvey("Initiative Community", EditInitiativeCommunitySurveyElems(teamID, &sic,
					userCapabilityCommunities(userID, community.Capability, teamID)))
				// Open a survey associated with the engagement
				dialogFromSurveyUnsafe(teamID, userID, message, val, actionValue, mc.Target, update, core.EmptyString)
			case string(models.Ignore):
				utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
				DeleteOriginalEng(userID, channelID, message.MessageTs)
			}
		}
	}
}

func handleInitiativeCommunityAssociationSelectObjective(mc models.MessageCallback, actionName, actionValue string, userID, channelID string,
	message slack.InteractionCallback, teamID models.TeamID) {
	if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
		// create association
		if strings.HasPrefix(actionName, strategy.CreatePrefix) || strings.HasPrefix(actionName, strategy.UpdatePrefix) {
			act, update := getUpdateParams(actionName)
			switch act {
			case string(models.Now):
				initID := actionSelected(message.ActionCallback.AttachmentActions)
				init := strategy.StrategyInitiativeByID(teamID, initID, strategyInitiativesTable)
				initCommID := mc.Target
				logger.
					WithField("initCommID", initCommID).
					WithField("init.ID", init.ID).
					Infof("Creating/updating %v initiative(ID=%s).InitiativeCommunityID := %s", update, init.ID, initCommID)
				init.InitiativeCommunityID = initCommID

				keyParams := map[string]*dynamodb.AttributeValue{
					"id":          dynString(init.ID),
					"platform_id": dynString(teamID.ToString()),
				}
				exprAttributes := map[string]*dynamodb.AttributeValue{
					":icid": dynString(initCommID),
				}
				updateExpression := "set initiative_community_id = :icid"
				err := d.UpdateTableEntry(exprAttributes, keyParams, updateExpression, strategyInitiativesTable)
				core.ErrorHandler(err, namespace, fmt.Sprintf("Could not update initiative in %s table",
					strategyObjectivesTable))
				PostMsgToUser(fmt.Sprintf("The selected initiative '%s' has been associated with this initiative community", init.Name), 
					userID, channelID, message.MessageTs)
			case string(models.Ignore):
				utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
			}
		}
	}
	DeleteOriginalEng(userID, channelID, message.MessageTs)
}

func handleObjectiveCommunityAssociationSelectObjective(mc models.MessageCallback, actionName, actionValue string, userID, channelID string,
	message slack.InteractionCallback, teamID models.TeamID) {
	if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
		// create association
		if strings.HasPrefix(actionName, strategy.CreatePrefix) || strings.HasPrefix(actionName, strategy.UpdatePrefix) {
			act, update := getUpdateParams(actionName)
			switch act {
			case string(models.Now):
				selected := actionSelected(message.ActionCallback.AttachmentActions)
				capObj := strategy.StrategyObjectiveByID(teamID, selected, strategyObjectivesTable)

				allCapComms := AllCapabilityCommunities(teamID)
				var associatableCapComms []strategy.CapabilityCommunity
				for _, each := range allCapComms {
					if !core.ListContainsString(capObj.CapabilityCommunityIDs, each.ID) {
						associatableCapComms = append(associatableCapComms, each)
					}
				}
				capComms := AsKv(associatableCapComms, "Name", "ID")

				val := utils.AttachmentSurvey("Objective Association",
					EditStrategyAssociation(&capObj, capComms, ui.PlainText("Objective Community"), ui.PlainText("Description")))
				// Open a survey associated with the engagement
				dialogFromSurveyUnsafe(teamID, userID, message, val, actionValue, "Association", update, selected)
				// DeleteOriginalEng(userID, channelID, message.MessageTs)
				if update {
					// delete the original association
				}
			case string(models.Ignore):
				utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
				DeleteOriginalEng(userID, channelID, message.MessageTs)
			}
		}
	}
}

func PostMsgToUser(text, userId, channelId, ts string) {
	publish(models.PlatformSimpleNotification{UserId: userId, Channel: channelId,
		Message: core.TextWrap(text, core.Underscore)})
	if ts != core.EmptyString {
		DeleteOriginalEng(userId, channelId, ts)
	}
}

func PostMsgToCommunity(commID community.AdaptiveCommunity, teamID models.TeamID, text string, attachs []ebm.Attachment) {
	comm := CommunityById(string(commID), teamID)
	publish(models.PlatformSimpleNotification{UserId: comm.ChannelID, Message: core.TextWrap(text, core.Underscore),
		Attachments: attachs})
}

func PostMsgToAdvocate(coordinator, userID string, community community.AdaptiveCommunity, purpose string) {
	if coordinator != userID {
		// Post a notification to the creator of the community that a notification has been sent to coordinator
		text := fmt.Sprintf("I have sent a notification to <@%s> to create a *private channel* and associate it with `%s %s Community` by inviting Adaptive into that channel",
			coordinator, purpose, strings.Title(string(community)))
		publish(models.PlatformSimpleNotification{UserId: userID, Message: text})
	}
	textWithPrefix := fmt.Sprintf("Hello, <@%s>. You have been assigned as a coordinator for the `%s %s Community` by <@%s>. Please create a *private channel*, invite Adaptive into the *private channel* using `/invite @adaptive`, and associate the channel with this new community. Finally, invite the team members to this channel by using the `/invite @user`.",
		coordinator, purpose, strings.Title(string(community)), userID)
	publish(models.PlatformSimpleNotification{UserId: coordinator, Message: textWithPrefix})
}

// Attachment action to add a new Capability Objective
func AddCapabilityObjectiveAttachAction(mc models.MessageCallback) ebm.AttachmentAction {
	return *models.SimpleAttachAction(
		*mc.WithAction(string(strategy.Create)).WithTopic(ObjectiveEvent),
		models.Now, "Create Objective")
}

func formatDate(date string, ipLayout, opLayout core.AdaptiveDateLayout) string {
	res, err := common.FormatDateWithIndefiniteOption(date, ipLayout, opLayout, namespace)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse string to date: %s", date))
	return res
}

func strategyObjectiveToFields(newSo, oldSo *models.StrategyObjective) (kvs []models.KvPair) {
	if oldSo == nil {
		oldSo = newSo
	}
	newDate := formatDate(newSo.ExpectedEndDate, DateFormat, core.USDateLayout)
	oldDate := formatDate(oldSo.ExpectedEndDate, DateFormat, core.USDateLayout)

	kvs = []models.KvPair{
		{Key: string(SObjectiveTypeLabel), Value: strategy.NewAndOld(string(newSo.ObjectiveType), string(oldSo.ObjectiveType))},
		{Key: string(SObjectiveNameLabel), Value: strategy.NewAndOld(newSo.Name, oldSo.Name)},
		{Key: string(SObjectiveDescriptionLabel), Value: strategy.NewAndOld(newSo.Description, oldSo.Description)},
		{Key: SObjectiveAdvocateLabel, Value: strategy.NewAndOld(common.TaggedUser(newSo.Advocate), common.TaggedUser(oldSo.Advocate))},
		{Key: SObjectiveTypeLabel, Value: strategy.NewAndOld(string(newSo.ObjectiveType), string(oldSo.ObjectiveType))},
		{Key: string(SObjectiveMeasuresLabel), Value: strategy.NewAndOld(newSo.AsMeasuredBy, oldSo.AsMeasuredBy)},
		{Key: string(SObjectiveTargetsLabel), Value: strategy.NewAndOld(newSo.Targets, oldSo.Targets)},
		{Key: SObjectiveEndDateLabel, Value: strategy.NewAndOld(newDate, oldDate)},
	}
	return
}

func ObjectiveViewAttachment(userID string, mc models.MessageCallback, newSo, oldSo *models.StrategyObjective, enableActions bool,
	addNewTopic string, initial bool, teamID models.TeamID) []ebm.Attachment {
	var editStatus string
	if oldSo == nil {
		editStatus = "created"
	} else {
		editStatus = "updated"
	}
	var title string
	var actions []ebm.AttachmentAction
	kvs := strategyObjectiveToFields(newSo, oldSo)

	if enableActions && AuthorizedForObjectiveEdit(userID, teamID, newSo) {
		var extraActions []ebm.AttachmentAction
		if initial {
			title = fmt.Sprintf("Here is the strategy objective you %s", editStatus)
		}
		allCapComms := AllCapabilityCommunities(teamID)
		if len(allCapComms) > 0 {
			// extraActions = append(extraActions, AllocateCapabilityObjectiveAttachAction(mc))
		}
		actions = strategy.EditAttachActions(mc, newSo.ID, true, true, false, addNewTopic, extraActions...)
	}
	return strategy.EntityViewAttachment(common.AttachmentEntity{MC: mc, Title: ui.RichText(title), Actions: actions, Fields: kvs})
}

func VisionMissionViewAttachment(mc models.MessageCallback, newVm, oldVm *models.VisionMission, enableActions,
	initial bool) []ebm.Attachment {
	editStatus := "created"
	var title string
	var actions []ebm.AttachmentAction
	var kvs []models.KvPair

	if oldVm != nil {
		editStatus = "updated"
		kvs = []models.KvPair{
			{Key: string(VisionLabel), Value: strategy.NewAndOld(newVm.Vision, oldVm.Vision)},
			{Key: MissionVisionAdvocateLabel,
				Value: strategy.NewAndOld(common.TaggedUser(newVm.Advocate), common.TaggedUser(oldVm.Advocate))},
		}
	} else {
		kvs = []models.KvPair{
			{Key: string(VisionLabel), Value: newVm.Vision}, {Key: MissionVisionAdvocateLabel,
				Value: common.TaggedUser(newVm.Advocate)},
		}
	}

	if enableActions {
		if initial {
			title = fmt.Sprintf("Here is the vision you %s", editStatus)
		}
		// Currently, we don't require adhoc creation of vision
		actions = append(actions, strategy.EditAttachActions(mc, newVm.ID, false, true, false, core.EmptyString,
			AddCapabilityObjectiveAttachAction(mc))...)
	}

	return strategy.EntityViewAttachment(common.AttachmentEntity{MC: mc, Title: ui.RichText(title), Actions: actions, Fields: kvs})
}

// AuthorizedForInitiativeAddEdit returns if the user is authorized to add/edit an initiative
func AuthorizedForInitiativeAddEdit(userID string, teamID models.TeamID, si *models.StrategyInitiative) (authorized bool) {
	// Anyone in an Objective Community
	userCapComms := userCapabilityCommunities(userID, community.Capability, teamID)
	var userCapCommIDs []string
	for _, eachComm := range userCapComms {
		userCapCommIDs = append(userCapCommIDs, eachComm.Value)
	}
	if len(userCapCommIDs) > 0 {
		logger.Infof("%s user belongs these Capability Communities: %v", userID, userCapComms)
		capObjs := AllStrategyObjectives(userID)
		logger.Infof("All Strategy Objectives for user %s: %v", userID, capObjs)
		for _, each := range capObjs {
			if len(core.InAAndB(userCapCommIDs, each.CapabilityCommunityIDs)) > 0 {
				logger.Infof("%s user Objective Community contains Capability Objectives. Authorizing the user to edit.", userID)
				authorized = true
			}
		}
	} else {
		logger.Infof("%s user does not belong to any Objective Community. Not authorized to edit the Initiative.",
			userID)
	}
	return
}

func AuthorizedForObjectiveEdit(userID string, teamID models.TeamID, so *models.StrategyObjective) (authorized bool) {
	// To edit a Capability Objective, a member must be in Strategy Community and should be in the Objective Community associated with the Capability Objective
	isMemberInStrategyCommunity := isMemberInCommunity(userID, community.Strategy)
	if isMemberInStrategyCommunity || (len(so.CapabilityCommunityIDs) > 0 &&
		isMemberInCommunity(userID, community.AdaptiveCommunity(fmt.Sprintf("%s:%s", community.Capability, so.CapabilityCommunityIDs[0])))) {
		logger.Infof("%s user is in Strategy Community and is in Objective Community associated with the Capability Objective. Authorized to edit.", userID)
		authorized = true
	}
	return
}

func InitiativeViewAttachment(userID string, mc models.MessageCallback, newSi, oldSi *models.StrategyInitiative, enableActions, init bool,
	teamID models.TeamID) []ebm.Attachment {
	editStatus := "created"
	var title string
	var actions []ebm.AttachmentAction
	var kvs []models.KvPair

	if oldSi != nil {
		editStatus = "updated"
		kvs = []models.KvPair{
			{Key: string(InitiativeNameLabel), Value: strategy.NewAndOld(newSi.Name, oldSi.Name)},
			{Key: string(InitiativeDescriptionLabel), Value: strategy.NewAndOld(newSi.Description, oldSi.Description)},
			{Key: string(InitiativeVictoryLabel), Value: strategy.NewAndOld(newSi.DefinitionOfVictory, oldSi.DefinitionOfVictory)},
			{Key: InitiativeAdvocateLabel, Value: strategy.NewAndOld(common.TaggedUser(newSi.Advocate), common.TaggedUser(oldSi.Advocate))},
			{Key: InitiativeBudgetLabel, Value: strategy.NewAndOld(newSi.Budget, oldSi.Budget)},
			{Key: InitiateEndDateLabel, Value: strategy.NewAndOld(newSi.ExpectedEndDate, oldSi.ExpectedEndDate)},
			{Key: InitiativeCapabilityObjectiveLabel, Value: strategy.NewAndOld(strategy.StrategyObjectiveByID(teamID,
				newSi.CapabilityObjective, strategyObjectivesTable).Name,
				strategy.StrategyObjectiveByID(teamID, oldSi.CapabilityObjective, strategyObjectivesTable).Name)},
		}
	} else {
		kvs = []models.KvPair{
			{Key: string(InitiativeNameLabel), Value: newSi.Name},
			{Key: string(InitiativeDescriptionLabel), Value: newSi.Description},
			{Key: string(InitiativeVictoryLabel), Value: newSi.DefinitionOfVictory},
			{Key: InitiativeAdvocateLabel, Value: common.TaggedUser(newSi.Advocate)},
			{Key: InitiativeBudgetLabel, Value: newSi.Budget},
			{Key: InitiateEndDateLabel, Value: newSi.ExpectedEndDate},
			{Key: InitiativeCapabilityObjectiveLabel, Value: strategy.StrategyObjectiveByID(teamID, newSi.CapabilityObjective,
				strategyObjectivesTable).Name},
		}
	}

	if enableActions && AuthorizedForInitiativeAddEdit(userID, teamID, newSi) {
		title = core.IfThenElse(init, fmt.Sprintf("This is the initiative you %s", editStatus), core.EmptyString).(string)
		actions = append(actions, strategy.EditAttachActions(mc, newSi.ID, true, true, false, InitiativeAdhocEvent)...)
	}
	return strategy.EntityViewAttachment(common.AttachmentEntity{MC: mc, Title: ui.RichText(title), Actions: actions, Fields: kvs})
}

type StrategyObjectiveCommunityAssociation struct {
	ObjectiveID   string `json:"objective_id"`
	ObjectiveName string `json:"objective_name"`
	CommunityID   string `json:"community_id"`
	CommunityName string `json:"community_name"`
	Description   string `json:"description"`
}

func StrategyObjectiveCommunityAssociationViewAttachment(mc models.MessageCallback, newSa *StrategyObjectiveCommunityAssociation,
	enableActions bool, delete bool) []ebm.Attachment {
	editStatus := "created"
	var title string
	var actions []ebm.AttachmentAction
	kvs := []models.KvPair{
		{Key: "Objective Community", Value: newSa.CommunityName},
		{Key: "Description", Value: newSa.Description},
	}

	if delete {
		editStatus = "removed"
	}
	title = fmt.Sprintf("This is the association you %s for the Capability Objective: %s", editStatus, newSa.ObjectiveName)
	if enableActions {
		actions = append(actions, strategy.EditAttachActions(mc, fmt.Sprintf("%s_%s", newSa.ObjectiveID, newSa.CommunityID), true, false, true, mc.Topic)...)
	}
	return strategy.EntityViewAttachment(common.AttachmentEntity{MC: mc, Title: ui.RichText(title), Actions: actions, Fields: kvs})
}

func handleCreateEvent(topic, text string, userID, channelID string, teamID models.TeamID,
	message slack.InteractionCallback, urgent bool) {
	// Create a strategy objective
	year, month := core.CurrentYearMonth()
	mc := models.MessageCallback{Module: "strategy", Source: userID, Topic: topic, Action: string(strategy.Create),
		Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	if urgent {
		handleCreateEvent1(mc, userID, channelID, teamID, CreateActionName, mc.ToCallbackID(), message)
	} else {
		CreateAskEngagement(engagementTable, teamID, mc, text, "", "", true, dns)
		DeleteOriginalEng(userID, channelID, message.MessageTs)
	}
}

func handleMenuEvent(text, userID string, mc models.MessageCallback, options []ebm.MenuOption) {
	attachActions := []ebm.AttachmentAction{
		*models.SelectAttachAction(mc, models.Now, "Choose ...", options,
			models.EmptyActionMenuOptionGroups()),
		*models.SimpleAttachAction(mc, models.Ignore, user.SkipActionLabel),
	}
	attach := utils.ChatAttachment(core.EmptyString, text, core.EmptyString, mc.ToCallbackID(), attachActions,
		[]ebm.AttachmentField{}, time.Now().Unix())
	publish(models.PlatformSimpleNotification{UserId: userID, Attachments: []ebm.Attachment{*attach}})
}

func initiativeCreateSurveyInitiativeOptions(adaptiveAssociatedComms []strategy.StrategyInitiativeCommunity,
	teamID models.TeamID) []ebm.MenuOption {
	var opts []ebm.MenuOption
	for _, each := range adaptiveAssociatedComms {
		eachComm := strategy.InitiativeCommunityByID(teamID, each.ID, strategyInitiativeCommunitiesTable)
		opts = append(opts, ebm.MenuOption{Text: eachComm.Name, Value: eachComm.ID})
	}
	return opts
}

func allCapabilityCommunitiesOptions(adaptiveAssociatedComms []strategy.CapabilityCommunity, teamID models.TeamID) []ebm.MenuOption {
	var opts []ebm.MenuOption
	for _, each := range adaptiveAssociatedComms {
		eachComm := strategy.CapabilityCommunityByID(teamID, each.ID, capabilityCommunitiesTable)
		opts = append(opts, ebm.MenuOption{Text: eachComm.Name, Value: eachComm.ID})
	}
	return opts
}

func allCapabilityObjectivesAssociatableOptions(userID string, teamID models.TeamID) []ebm.MenuOption {
	var opts []ebm.MenuOption
	objs := AllStrategyCapabilityObjectives(userID)
	capComms := AllCapabilityCommunities(teamID)
	capCommIDs := AsValues(capComms, "ID")
	for _, each := range objs {
		eachObjCapCommIDs := each.CapabilityCommunityIDs
		if len(core.InAButNotB(capCommIDs, eachObjCapCommIDs)) > 0 {
			opts = append(opts, ebm.MenuOption{Text: each.Name, Value: each.ID})
		}
	}
	return opts
}

func allInitiativesAssociatableOptions(userID string, teamID models.TeamID) []ebm.MenuOption {
	var opts []ebm.MenuOption
	inits := AllStrategyInitiatives(teamID)
	initComms := getInitiativeCommunitiesForUserIDUnsafe(userID, models.TeamID(teamID))
	commIDs := AsValues(initComms, "ID")
	for _, each := range inits {
		eachCommIDs := []string{each.InitiativeCommunityID}
		if len(core.InAButNotB(commIDs, eachCommIDs)) > 0 {
			opts = append(opts, ebm.MenuOption{Text: each.Name, Value: each.ID})
		}
	}
	return opts
}

func actionSelected(actions []*slack.AttachmentAction) string {
	var selected string
	if len(actions) > 0 {
		if actions[0].SelectedOptions != nil {
			if len(actions[0].SelectedOptions) > 0 {
				selected = actions[0].SelectedOptions[0].Value
			}
		}
	}
	return selected
}

func handleCreateEvent1(mc models.MessageCallback, userID, channelID string, teamID models.TeamID, actionName,
	actionValue string, message slack.InteractionCallback) {
	logger.WithField("mc.Topic", mc.Topic).Infof("handleCreateEvent1")
	switch mc.Topic {
	case ObjectiveEvent:
		// Once retrieving the cap comm id, set target back to empty. This wil further be used if the cap community is being updated
		// Check for the vision advocate. If no one exists post message of the same
		vm := StrategyVision(models.TeamID(teamID))
		if vm == nil {
			publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
				Message: core.TextWrap("Sorry, I cannot create an objective without having a vision advocate defined.",
					core.Underscore), Ts: message.MessageTs})
		} else {
			handleObjectiveCreate(mc, actionName, actionValue, userID, channelID, teamID,
				message, models.CapabilityStrategyObjective)
		}
	case FinancialObjectiveEvent:
		vm := StrategyVision(models.TeamID(teamID))
		if vm == nil {
			publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
				Message: core.TextWrap("Sorry, I cannot create an objective without having a vision advocate defined.",
					core.Underscore), Ts: message.MessageTs})
		} else {
			handleObjectiveCreate(mc, actionName, actionValue, userID, channelID, teamID,
				message, models.FinancialStrategyObjective)
		}
	case CustomerObjectiveEvent:
		vm := StrategyVision(models.TeamID(teamID))
		if vm == nil {
			publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
				Message: core.TextWrap("Sorry, I cannot create an objective without having a vision advocate defined.",
					core.Underscore), Ts: message.MessageTs})
		} else {
			handleObjectiveCreate(mc, actionName, actionValue, userID, channelID, teamID, message, models.CustomerStrategyObjective)
		}
	case VisionEvent:
		handleMissionVisionCreate(mc, actionName, actionValue, userID, channelID, teamID, message)
	case CapabilityCommunityEvent:
		handleCapabilityCommunityCreate(mc, actionName, actionValue, userID, channelID, message, teamID)
	case strategy.InitiativeEvent:
		handleInitiativeCreate(mc, actionName, actionValue, userID, channelID, teamID, message)
	case InitiativeSelectCommunityEvent:
		handleInitiativeCreate(mc, actionName, actionValue, userID, channelID, teamID, message)
	case ObjectiveSelectCommunityEvent:
		handleObjectiveCreate(mc, actionName, actionValue, userID, channelID, teamID, message, models.CapabilityStrategyObjective)
	case InitiativeCommunityEvent:
		handleInitiativeCommunityCreate(mc, actionName, actionValue, userID, channelID, message, teamID)
	case ObjectiveCommunityAssociationSelectObjective:
		handleObjectiveCommunityAssociationSelectObjective(mc, actionName, actionValue, userID, channelID, message, teamID)
	case InitiativeCommunityAssociationSelectInitiative:
		handleInitiativeCommunityAssociationSelectObjective(mc, actionName, actionValue, userID, channelID, message, teamID)
	// adhoc events
	case ObjectiveAdhocEvent:
		handleMenuObjectiveCreate(userID, channelID, teamID, message, false)
	case AssociateObjectiveWithCapabilityCommunityEvent:
		handleMenuObjectiveAssociationCreate(userID, channelID, message, true, teamID)
	case strategy.CapabilityCommunityAdhocEvent:
		// copied from CreateCapabilityCommunity event in menu_list
		handleCreateEvent(CapabilityCommunityEvent, "Would you like to create a capability community?", userID,
			channelID, teamID, message, true)
	case InitiativeAdhocEvent:
		handleMenuCreateInitiative(userID, channelID, teamID, message, false)
	case strategy.InitiativeCommunityAdhocEvent:
		// copied from CreateInitiativeCommunity event in menu_list
		handleCreateEvent(InitiativeCommunityEvent, "Would you like to create an initiative community?", userID,
			channelID, teamID, message, true)
	case strategy.AssociateInitiativeWithInitiativeCommunityEvent:
		handleMenuInitiativeAssociationCreate(userID, channelID, message, false, teamID,  mc.Target)
	default:
		logger.WithField("mc.Topic", mc.Topic).Infof("Unhandled topic in handleCreateEvent1")
	}
}

func AllStrategyObjectives(userID string) []models.StrategyObjective {
	return strategy.UserStrategyObjectives(userID, strategyObjectivesTable, strategyObjectivesPlatformIndex,
		userObjectivesTable, communityUsersTable, communityUsersUserCommunityIndex)
}

func UserCommunityObjectives(userID string) (objs []models.StrategyObjective) {
	isStrategyUser := isMemberInCommunity(userID, community.Strategy)
	if isStrategyUser {
		objs = AllStrategyObjectives(userID)
	} else {
		objs = strategy.UserCommunityObjectives(userID, strategyObjectivesTable, strategyObjectivesPlatformIndex,
			userObjectivesTable, communityUsersTable, communityUsersUserIndex)
	}
	return
}

func AllStrategyCapabilityObjectives(userID string) []models.StrategyObjective {
	var op []models.StrategyObjective
	for _, each := range AllStrategyObjectives(userID) {
		if len(each.CapabilityCommunityIDs) > 0 {
			op = append(op, each)
		}
	}
	return op
}

func AllCapabilityCommunities(teamID models.TeamID) []strategy.CapabilityCommunity {
	return strategy.AllCapabilityCommunities(teamID, capabilityCommunitiesTable, capabilityCommunitiesPlatformIndex, strategyCommunitiesTable)
}

func AllStrategyInitiatives(teamID models.TeamID) []models.StrategyInitiative {
	return strategy.AllOpenStrategyInitiatives(teamID, strategyInitiativesTable, strategyInitiativesPlatformIndex,
		userObjectivesTable)
}

func CapabilityCommunityInitiatives(userID string) []models.StrategyInitiative {
	return strategy.UserCapabilityCommunityInitiatives(userID, strategyObjectivesTable, strategyObjectivesPlatformIndex,
		strategyInitiativesTable, strategyInitiativesInitiativeCommunityIndex, userObjectivesTable,
		communityUsersTable, communityUsersUserCommunityIndex, communityUsersUserIndex)
}

func InitiativeCommunityInitiatives(userID string) []models.StrategyInitiative {
	return strategy.UserInitiativeCommunityInitiatives(userID,
		strategyInitiativesTable, strategyInitiativesInitiativeCommunityIndex,
		communityUsersTable, communityUsersUserIndex)
}

func StrategyInitiativeCommunitiesForUserID(userID string, teamID models.TeamID) []strategy.StrategyInitiativeCommunity {
	return strategy.UserStrategyInitiativeCommunities(userID, communityUsersTable, communityUsersUserCommunityIndex,
		communityUsersUserIndex, strategyInitiativeCommunitiesTable, strategyInitiativeCommunitiesPlatformIndex,
		strategyCommunitiesTable, teamID)
}

func StrategyVision(teamID models.TeamID) *models.VisionMission {
	return strategy.StrategyVision(teamID, visionTable)
}

// InterfaceSlice converts an interface to list of interfaces if applicable
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic(errors.New("InterfaceSlice() given a non-slice type"))
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// Get a string field value from an interface
func GetFieldString(i interface{}, field string) string {
	// Create a value for the slice.
	v := reflect.ValueOf(i)
	// Get the field of the slice element that we want to set.
	f := v.FieldByName(field)
	// Get value
	return f.String()
}

// AsKv converts a list interface to list of models.KvPair. Key and value fileds are passed and evaluated using reflection
func AsKv(ip interface{}, keyName, valueName string) []models.KvPair {
	var op []models.KvPair
	for _, each := range InterfaceSlice(ip) {
		op = append(op, models.KvPair{Key: GetFieldString(each, keyName), Value: GetFieldString(each, valueName)})
	}
	return op
}

func AsValues(ip interface{}, keyName string) []string {
	var op []string
	for _, each := range InterfaceSlice(ip) {
		op = append(op, GetFieldString(each, keyName))
	}
	return op
}

func callback(source, topic string) models.MessageCallback {
	year, month := core.CurrentYearMonth()
	return models.MessageCallback{Module: "strategy", Source: source, Topic: topic, Month: strconv.Itoa(int(month)),
		Year: strconv.Itoa(year)}
}

func handleMenuCreateInitiative(userID, channelID string, teamID models.TeamID,
	message slack.InteractionCallback, deleteOriginal bool) {
	logger.Infof("In handleMenuCreateInitiative for user %s with platform %s", userID, teamID)
	// Query all the Strategy Initiative communities
	initComms := getInitiativeCommunitiesForUserIDUnsafe(userID, models.TeamID(teamID))

	var adaptiveAssociatedInitComms []strategy.StrategyInitiativeCommunity
	// Get a list of Adaptive associated Initiative communities
	for _, each := range initComms {
		eachStrategyComms := StrategyCommunityByID(each.ID)
		if eachStrategyComms.ChannelCreated == 1 {
			adaptiveAssociatedInitComms = append(adaptiveAssociatedInitComms, each)
		}
	}
	logger.Infof("Adaptive associated Initiative Communities for platform %s: %s", teamID, adaptiveAssociatedInitComms)
	if len(adaptiveAssociatedInitComms) > 0 {
		logger.Infof("Initiatives communities exist for user %s with platform %s", userID, teamID)
		mc := models.MessageCallback{Module: string(community.Strategy), Source: userID, Topic: InitiativeSelectCommunityEvent,
			Action: string(strategy.Create)}
		handleMenuEvent("Select an initiative community", userID, mc,
			initiativeCreateSurveyInitiativeOptions(adaptiveAssociatedInitComms, teamID))
		if deleteOriginal {
			DeleteOriginalEng(userID, channelID, message.MessageTs)
		}
	} else {
		handleCreateEvent(InitiativeCommunityEvent, "There are no Adaptive associated Initiative Communities. If you have already created an Initiative Community, please ask the coordinator to create a *_private_* channel, invite Adaptive and associate with the community.",
			userID, channelID, teamID, message, false)
	}
}

func handleMenuObjectiveCreate(userID, channelID string, teamID models.TeamID, message slack.InteractionCallback, deleteOriginal bool) {
	logger.Infof("Creating Strategy Objective by user %s for platform %s", userID, teamID)
	if isMemberInCommunity(userID, community.Strategy) {
		// check if the user is in strategy community
		adaptiveAssociatedCapComms := SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated(teamID)

		logger.Infof("Adaptive associated Capability Communities for platform %s: %s", teamID, adaptiveAssociatedCapComms)
		if len(adaptiveAssociatedCapComms) > 0 {
			// Enable a user to create an objective if user is in strategy community and there are capability communities
			mc := models.MessageCallback{Module: string(community.Strategy), Source: userID, Topic: ObjectiveSelectCommunityEvent,
				Action: string(strategy.Create)}
			handleMenuEvent("Select a capability community. You can assign the objective to other communities later but you need at least one for now.",
				userID, mc, allCapabilityCommunitiesOptions(adaptiveAssociatedCapComms, models.TeamID(teamID)))
			if deleteOriginal {
				DeleteOriginalEng(userID, channelID, message.MessageTs)
			}
		} else {
			handleCreateEvent(CapabilityCommunityEvent, "There are no Adaptive associated Capability Communities. If you have already created a Objective Community, please ask the coordinator to create a *_private_* channel, invite Adaptive and associate with the community.",
				userID, channelID, teamID, message, false)
		}
	} else {
		// send a message that user is not authorized to create objectives
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
			Message: fmt.Sprintf("You are not part of the Adaptive Strategy Community, you will not be able to create Capability Objectives.")})
	}
}

func handleMenuInitiativeAssociationCreate(userID, channelID string, message slack.InteractionCallback,
	deleteOriginal bool, teamID models.TeamID, initCommID string) {
	mc := models.MessageCallback{
		Module: string(community.Strategy), 
		Source: userID,
		Topic: InitiativeCommunityAssociationSelectInitiative,
		Action: string(strategy.Create),
		Target: initCommID,
	}
	initCommsOpts := allInitiativesAssociatableOptions(userID, teamID)
	if len(initCommsOpts) > 0 {
		// TODO: Check if there are objectives to associate
		handleMenuEvent("Select which initiative you want to associate with the Initiative Community.", userID, mc, initCommsOpts)
	} else {
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
			Message: fmt.Sprintf("There are no initiatives to associate with this Community.")})
	}
	if deleteOriginal {
		DeleteOriginalEng(userID, channelID, message.MessageTs)
	}
}

func handleMenuObjectiveAssociationCreate(userID, channelID string, message slack.InteractionCallback,
	deleteOriginal bool, teamID models.TeamID) {
	mc := models.MessageCallback{Module: string(community.Strategy), Source: userID,
		Topic: ObjectiveCommunityAssociationSelectObjective, Action: string(strategy.Create)}
	objOpts := allCapabilityObjectivesAssociatableOptions(userID, teamID)
	if len(objOpts) > 0 {
		// TODO: Check if there are objectives to associate
		handleMenuEvent("Select which objective you want to associate with a Objective Community.", userID, mc, objOpts)
	} else {
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
			Message: fmt.Sprintf("There are no Capability Objectives to associate with Communities.")})
	}
	if deleteOriginal {
		DeleteOriginalEng(userID, channelID, message.MessageTs)
	}
}

func HandleRequest(ctx context.Context, np models.NamespacePayload4) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			logger.Errorf("Recovering from error in strategy-lambda %+v", err2)
			err = fmt.Errorf("error captured in strategy-lambda %v", err2)
		}
		if err != nil {
			logger.Errorf("Error in strategy-lambda: %+v", err)
		}
	}()
	logger = logger.WithLambdaContext(ctx)

	byt, err := json.Marshal(np)
	logger.WithField("payload", string(byt)).Info()

	if np.ID == "warmup" {
		logger.WithField("warmup", true).Info()
	} else {
		logger.WithField("CallbackID", np.InteractionCallback.CallbackID).Info()
		if strings.HasPrefix(np.InteractionCallback.CallbackID, StrategyPath.Encode()) {
			logger.WithField("CallbackID", np.InteractionCallback.CallbackID).Error("Invocation of the old workflow")
			// err = invokeWorkflow(np)
		} else
		if np.Namespace == "strategy" {
			switch np.SlackRequest.Type {
			case models.InteractionSlackRequestType:
				err = onSlackInteraction(np)
			case models.DialogSubmissionSlackRequestType:
				onDialogSubmission(np)
			case models.DialogCancellationSlackRequestType:
				onDialogCancellation(np)
			}
		}
	}
	return
}

func onSlackInteraction(np models.NamespacePayload4) (err error) {
	logger.WithField("interactive_message_event", np).Info()

	message := np.SlackRequest.InteractionCallback
	request := message
	teamID := np.TeamID

	userID := message.User.ID
	channelID := message.Channel.ID
	action := message.ActionCallback.AttachmentActions[0]
	logger.WithField("action", action).WithField("userID", userID).Info()
	// 'menu_list' is for the options that are presented to the user
	if action.Name == "menu_list" {
		selected := action.SelectedOptions[0]
		logger.
			WithField("name", action.Name).
			WithField("value", action.Value).
			WithField("selected.Value", selected.Value).
			Info("In menu_list")
		switch selected.Value {
		case CreateStrategyObjective:
			// Create a strategy objective
			logger.Error("Not entering Old CreateObjectiveWorkflow")
			// err = enterWorkflow(CreateObjectiveWorkflow, np, "")
		case AssociateInitiativeWithInitiativeCommunity:
			handleMenuObjectiveAssociationCreate(userID, channelID, message, true, models.TeamID(teamID))
		// // case AssociateInitiativeWithInitiativeCommunity:
		// 	handleCreateEvent(strategy.AssociateInitiativeWithInitiativeCommunityEvent, "I see you want to associate an initiative with an initiative community.",
		// 		userID, channelID, teamID, message, true)
		case CreateFinancialObjective:
			// np.InteractionCallback.CallbackID = FirstWorkflowPath.Encode()
			// invokeWorkflow(np) // TODO: This is a temporary invocation of workflow from menu. Just to make sure everything is working.
			handleCreateEvent(FinancialObjectiveEvent, "Would you like to create a financial objective?",
				userID, channelID, teamID, message, true)
		case CreateCustomerObjective:
			handleCreateEvent(CustomerObjectiveEvent, "Would you like to create a customer objective?",
				userID, channelID, teamID, message, true)
		case strategy.CreateVision:
			// Create a strategy objective
			handleCreateEvent(VisionEvent, "Would you like to add vision?", userID, channelID,
				teamID, message, true)
		case strategy.ViewVision, strategy.ViewEditVision:
			onViewEditVision(request, teamID)
		case CreateCapabilityCommunity:
			handleCreateEvent(CapabilityCommunityEvent, "Would you like to create a capability community?",
				userID, channelID, teamID, message, true)
		case AssociateStrategyObjectiveToCapabilityCommunity:
			handleCreateEvent(AssociateObjectiveWithCapabilityCommunityEvent, "I see you want to associate objective with capability community",
				userID, channelID, teamID, message, true)
		case CreateInitiative:
			handleMenuCreateInitiative(userID, channelID, teamID, message, true)
		case CreateInitiativeCommunity:
			handleCreateEvent(InitiativeCommunityEvent, "Would you like to create an initiative community?",
				userID, channelID, teamID, message, true)
		case ViewAdvocacyObjectives:
			logger.Error("Not entering Old CreateObjectiveWorkflow/ViewMyObjectivesEvent")
			// err = enterWorkflow(CreateObjectiveWorkflow, np, ViewMyObjectivesEvent)
		case ViewStrategyObjectives:
			logger.Error("Not entering Old CreateObjectiveWorkflow/ViewObjectivesEvent")
			// err = enterWorkflow(CreateObjectiveWorkflow, np, ViewObjectivesEvent)
			// onViewObjectives(request, teamID, AllStrategyObjectives(userID))
		case ViewCapabilityCommunityObjectives:
			onViewObjectives(request, teamID, UserCommunityObjectives(userID))
		case ViewCapabilityCommunityInitiatives, ViewInitiativeCommunityInitiatives:
			var inits []models.StrategyInitiative
			inStrategyCommunity := community.IsUserInCommunity(userID, communityUsersTable,
				communityUsersUserCommunityIndex, community.Strategy)
			if inStrategyCommunity {
				// User is in Strategy community, show all Initiatives
				inits = AllStrategyInitiatives(teamID)
			} else if selected.Value == ViewCapabilityCommunityInitiatives {
				inits = CapabilityCommunityInitiatives(userID)
			} else if selected.Value == ViewInitiativeCommunityInitiatives {
				inits = InitiativeCommunityInitiatives(userID)
			}
			onViewInitiatives(request, teamID, inits)
		default:
			logger.Infof("Unhandled option %s", selected.Value)
		}
	} else {
		action := message.ActionCallback.AttachmentActions[0]
		var actionValue string

		// Parse callback Id to messageCallback
		// For menu select, take mc from callbackID
		if len(action.SelectedOptions) > 0 {
			mc, err := utils.ParseToCallback(message.CallbackID)
			// For menu options, action value will be empty. Assign it as callbackID
			actionValue = message.CallbackID
			core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))
			handleCreateEvent1(*mc, userID, channelID, teamID, action.Name, actionValue, message)
		} else {
			mc, err := utils.ParseToCallback(action.Value)
			actionValue = action.Value
			core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))
			if mc.Topic == ObjectiveCommunityAssociationSelectObjective {
				switch mc.Action {
				case string(strategy.Delete):
					onObjectiveCommunityAssociationSelectObjectiveDelete(request, teamID)
				case string(strategy.Create):
					// Method to ask user to select another objective to associate
					handleMenuObjectiveAssociationCreate(userID, channelID, message,
						true, models.TeamID(teamID))
				}
			} else {
				handleCreateEvent1(*mc, userID, channelID, teamID, action.Name, actionValue, message)
			}

		}
	}
	return
}

func onObjectiveCommunityAssociationSelectObjectiveDelete(request slack.InteractionCallback,
	teamID models.TeamID) {
	message := request

	userID := message.User.ID
	channelID := message.Channel.ID
	action := message.ActionCallback.AttachmentActions[0]
	mc, err := utils.ParseToCallback(action.Value)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))

	splits := strings.Split(mc.Target, "_")
	capObjID := splits[0]
	capCommID := splits[1]
	capObj := strategy.StrategyObjectiveByID(models.TeamID(teamID), capObjID, strategyObjectivesTable)
	capComm := strategy.CapabilityCommunityByID(models.TeamID(teamID), capCommID, capabilityCommunitiesTable)
	existingCommunities := capObj.CapabilityCommunityIDs
	indexToRemove := SliceIndex(len(existingCommunities), func(i int) bool { return existingCommunities[i] == capCommID })
	updatedCommunityList := RemoveFromSliceOnIndex(existingCommunities, indexToRemove)
	// Updating communities list in the database
	updateObjectiveCommunities(teamID, capObjID, updatedCommunityList)

	target := fmt.Sprintf("%s_%s", capObjID, capCommID)
	newSa := &StrategyObjectiveCommunityAssociation{ObjectiveID: capObjID, ObjectiveName: capObj.Name,
		CommunityID: capCommID, CommunityName: capComm.Name}
	attachs := StrategyObjectiveCommunityAssociationViewAttachment(*mc.WithTopic(AssociateObjectiveWithCapabilityCommunityEvent).WithTarget(target),
		newSa, false, true)
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Attachments: attachs, Ts: message.MessageTs})
}

func onViewInitiatives(request slack.InteractionCallback,
	teamID models.TeamID, inits []models.StrategyInitiative) {
	message := request

	userID := message.User.ID
	channelID := message.Channel.ID
	threadTs := utils.TimeStamp(message)
	mc := callback(userID, InitiativeSelectCommunityEvent)

	if len(inits) > 0 {
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
			Message: "You can find the list of initiatives in the thread."})
		for _, each := range inits {
			attachs := InitiativeViewAttachment(userID, mc, &each, nil, true, false, models.TeamID(teamID))
			publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
				Attachments: attachs, ThreadTs: threadTs})
		}
	} else {
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
			Message: "Currently, there are no strategy initiatives"})
	}
}

func onViewObjectives(request slack.InteractionCallback,
	teamID models.TeamID, objs []models.StrategyObjective) {
	message := request

	userID := message.User.ID
	channelID := message.Channel.ID

	mc := callback(userID, ObjectiveEvent)
	threadTs := utils.TimeStamp(message)
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		Message: "You can find the list of strategy objectives in the thread."})
	// DeleteOriginalEng(userID, channelID, message.MessageTs)
	if len(objs) > 0 {
		for _, each := range objs {
			attachs := ObjectiveViewAttachment(userID, mc, &each, nil, true,
				ObjectiveAdhocEvent, false, models.TeamID(teamID))
			publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
				Attachments: attachs, ThreadTs: threadTs})
		}
	} else {
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
			Message: "Currently, there are no strategy objectives"})
	}
}

func onViewEditVision(request slack.InteractionCallback, teamID models.TeamID) {
	message := request

	userID := message.User.ID
	channelID := message.Channel.ID
	vm := StrategyVision(teamID)
	if vm != nil {
		// Post vision attachment only if it's not nil
		mc := models.MessageCallback{Module: "strategy", Source: userID, Topic: VisionEvent, Target: vm.ID}
		attachs := ViewVisionMissionAttachment(userID, teamID, vm, mc)
		text := "Below is the vision"
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Message: text,
			Attachments: attachs})
	} else {
		PostMsgToUser("There is no vision defined yet", userID, channelID, core.EmptyString)
	}
	DeleteOriginalEng(userID, channelID, message.MessageTs)
}

func onDialogSubmission(np models.NamespacePayload4) {
	logger.WithField("dialog_submission_event", np).Info()
	request := np.SlackRequest.InteractionCallback
	// Handling dialog submission for each answer
	dialog := request
	// Parse callback Id to messageCallback
	mc, err := utils.ParseToCallback(dialog.CallbackID)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))

	var msgState MsgState
	err = json.Unmarshal([]byte(dialog.State), &msgState)
	core.ErrorHandler(err, namespace, "Could not unmarshal to MsgState")

	teamID := models.ParseTeamID(userDAO.ReadUnsafe(dialog.User.ID).PlatformID)
	notes := responses()
	if mc.Topic == ObjectiveSelectCommunityEvent || mc.Topic == ObjectiveEvent {
		if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
			notes = onObjectiveEventCreateOrUpdateDialogSubmission(dialog, msgState, teamID, mc)
		}
	} else if mc.Topic == FinancialObjectiveEvent {
		if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
			notes = onFinancialObjectiveEventCreateOrUpdateDialogSubmission(dialog, msgState, teamID, mc)
		}
	} else if mc.Topic == CustomerObjectiveEvent {
		if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
			notes = onCustomerObjectiveEventCreateOrUpdateDialogSubmission(dialog, msgState, teamID, mc)
		}
	} else if mc.Topic == VisionEvent {
		if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
			notes = onVisionEventCreateOrUpdateDialogSubmission(dialog, msgState, teamID, mc)
		}
	} else if mc.Topic == CapabilityCommunityEvent {
		if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
			notes = onCapabilityCommunityEventCreateOrUpdateDialogSubmission(dialog, msgState, teamID, mc)
		}
	} else if mc.Topic == InitiativeSelectCommunityEvent || mc.Topic == strategy.InitiativeEvent {
		if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
			notes = onInitiativeSelectCommunityEventCreateOrUpdateDialogSubmission(dialog, msgState, teamID, mc)
		}
	} else if mc.Topic == InitiativeCommunityEvent {
		if mc.Action == string(strategy.Create) || mc.Action == string(strategy.Update) {
			notes = onInitiativeCommunityEventCreateOrUpdateDialogSubmission(dialog, msgState, teamID, mc)
		}
	} else if mc.Topic == ObjectiveCommunityAssociationSelectObjective {
		if mc.Action == string(strategy.Create) {
			notes = onObjectiveCommunityAssociationSelectObjectiveCreateDialogSubmission(dialog, msgState, teamID, mc)
		}
	}
	publishAll(notes)
}
func onObjectiveEventCreateOrUpdateDialogSubmission(request slack.InteractionCallback, msgState MsgState, teamID models.TeamID, mc *models.MessageCallback) (resp []models.PlatformSimpleNotification) {
	dialog := request
	userID := dialog.User.ID
	channelID := dialog.Channel.ID
	tpe := dialog.Submission[SObjectiveType]
	advocate := dialog.Submission[SObjectiveAdvocate]
	endDate := dialog.Submission[SObjectiveEndDate]
	name := dialog.Submission[SObjectiveName]
	desc := dialog.Submission[SObjectiveDescription]
	measures := dialog.Submission[SObjectiveMeasures]
	targets := dialog.Submission[SObjectiveTargets]
	capCommID := msgState.SelectedOption

	var editStatus = "created"
	var newSo *models.StrategyObjective
	var oldSo *models.StrategyObjective
	// If update, change existing record. Retrieve id from the dialog state.
	if msgState.Update {
		editStatus = "updated"
		var result models.StrategyObjective
		entity := StrategyEntityById(msgState.Id, teamID, strategyObjectivesTable)
		byt1, _ := json.Marshal(entity)
		err := json.Unmarshal(byt1, &result)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not decode interface to struct"))
		oldSo = &result
		oldComms := oldSo.CapabilityCommunityIDs
		oldComms[0] = capCommID
		newSo = &models.StrategyObjective{ID: msgState.Id, PlatformID: teamID.ToPlatformID(), Name: name,
			Description: desc, AsMeasuredBy: measures, ObjectiveType: models.StrategyObjectiveType(tpe),
			Targets:  targets,
			Advocate: advocate, CapabilityCommunityIDs: oldComms, ExpectedEndDate: endDate,
			CreatedBy: oldSo.CreatedBy, CreatedAt: oldSo.CreatedAt}
	} else {
		id := core.Uuid()
		oldSo = nil
		newSo = &models.StrategyObjective{ID: id, PlatformID: teamID.ToPlatformID(), Name: name, Description: desc,
			AsMeasuredBy: measures, Targets: targets, ObjectiveType: models.StrategyObjectiveType(tpe),
			Advocate: advocate, CapabilityCommunityIDs: []string{capCommID}, ExpectedEndDate: endDate,
			CreatedBy: userID, CreatedAt: time.Now().Format(string(TimestampFormat))}
	}
	// Write entry to table
	err := d.PutTableEntry(*newSo, strategyObjectivesTable)
	core.ErrorHandler(err, namespace, "Could not put strategy objective to DB")
	attachsWithActions := ObjectiveViewAttachment(userID, *mc, newSo, oldSo,
		true, ObjectiveAdhocEvent, true, teamID)
	attachsWithNoActions := ObjectiveViewAttachment(userID, *mc, newSo, oldSo,
		false, ObjectiveAdhocEvent, true, teamID)
	// Publish the updated objective with actions into the user channel
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: msgState.ThreadTs,
		Attachments: attachsWithActions})
	// Do analysis on objective description
	utils.ECAnalysis(desc, stratObjDescriptionContext, "Strategy objective", dialogTableName,
		mc.ToCallbackID(), userID, channelID, msgState.ThreadTs, msgState.ThreadTs, attachsWithActions,
		s, platformNotificationTopic, namespace)
	text := fmt.Sprintf("Below objective has been %s by <@%s>", editStatus, userID)
	// Post objectives with no actions to strategy community
	PostMsgToCommunity(community.Strategy, teamID, text, attachsWithNoActions)
	// Post to associated capability community with no actions
	stratComm := StrategyCommunityByID(capCommID)
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: stratComm.ChannelID,
		Message: text, Attachments: attachsWithNoActions})

	uObj := UserObjectiveFromStrategyObjective(newSo, capCommID, teamID)
	err = d.PutTableEntry(uObj, userObjectivesTable)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not write entry to %s table", userObjectivesTable))

	utils.UpdateEngAsAnswered(mc.Source, mc.WithTarget(core.EmptyString).ToCallbackID(), engagementTable, d, namespace)
	return responses()
}
func onFinancialObjectiveEventCreateOrUpdateDialogSubmission(request slack.InteractionCallback, msgState MsgState, teamID models.TeamID, mc *models.MessageCallback) (resp []models.PlatformSimpleNotification) {
	dialog := request
	userID := dialog.User.ID
	channelID := dialog.Channel.ID
	advocate := dialog.Submission[SObjectiveAdvocate]
	endDate := dialog.Submission[SObjectiveEndDate]
	name := dialog.Submission[SObjectiveName]
	desc := dialog.Submission[SObjectiveDescription]
	measures := dialog.Submission[SObjectiveMeasures]
	targets := dialog.Submission[SObjectiveTargets]

	var editStatus = "created"
	var newSo *models.StrategyObjective
	var oldSo *models.StrategyObjective
	// If update, change existing record. Retrieve id from the dialog state.
	if msgState.Update {
		editStatus = "updated"
		var result models.StrategyObjective
		entity := StrategyEntityById(msgState.Id, teamID, strategyObjectivesTable)
		byt1, _ := json.Marshal(entity)
		err := json.Unmarshal(byt1, &result)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not decode interface to struct"))
		oldSo = &result
		newSo = &models.StrategyObjective{ID: msgState.Id, PlatformID: teamID.ToPlatformID(), Name: name,
			Description: desc, AsMeasuredBy: measures, Targets: targets, ObjectiveType: models.FinancialStrategyObjective,
			Advocate: advocate, ExpectedEndDate: endDate, CreatedBy: userID, CreatedAt: TimestampFormat.Format(time.Now())}
	} else {
		id := core.Uuid()
		oldSo = nil
		newSo = &models.StrategyObjective{ID: id, PlatformID: teamID.ToPlatformID(), Name: name, Description: desc,
			AsMeasuredBy: measures, Targets: targets, ObjectiveType: models.FinancialStrategyObjective,
			Advocate: advocate, ExpectedEndDate: endDate, CreatedBy: userID, CreatedAt: TimestampFormat.Format(time.Now())}
	}
	// Write entry to table
	err := d.PutTableEntry(*newSo, strategyObjectivesTable)
	core.ErrorHandler(err, namespace, "Could not put strategy objective to DB")
	attachsWithActions := ObjectiveViewAttachment(userID, *mc, newSo, oldSo,
		true, FinancialObjectiveEvent, true, teamID)
	attachsWithNoActions := ObjectiveViewAttachment(userID, *mc, newSo, oldSo,
		false, FinancialObjectiveEvent, true, teamID)
	// Publish the updated objective with actions into the user channel
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: msgState.ThreadTs,
		Attachments: attachsWithActions})
	// Do analysis on objective description
	utils.ECAnalysis(desc, stratObjDescriptionContext, "Financial objective", dialogTableName,
		mc.ToCallbackID(), userID, channelID, msgState.ThreadTs, msgState.ThreadTs, attachsWithActions,
		s, platformNotificationTopic, namespace)
	text := fmt.Sprintf("Below objective has been %s by <@%s>", editStatus, userID)
	// Post objectives with no actions to strategy community
	PostMsgToCommunity(community.Strategy, teamID, text, attachsWithNoActions)
	UserObjectiveFromStrategyObjective(newSo, core.EmptyString, teamID)

	utils.UpdateEngAsAnswered(mc.Source, mc.WithTarget(core.EmptyString).ToCallbackID(),
		engagementTable, d, namespace)
	return responses()
}
func onCustomerObjectiveEventCreateOrUpdateDialogSubmission(request slack.InteractionCallback, msgState MsgState, teamID models.TeamID, mc *models.MessageCallback) (resp []models.PlatformSimpleNotification) {
	dialog := request
	userID := dialog.User.ID
	channelID := dialog.Channel.ID
	advocate := dialog.Submission[SObjectiveAdvocate]
	tpe := dialog.Submission[SObjectiveType]
	endDate := dialog.Submission[SObjectiveEndDate]
	name := dialog.Submission[SObjectiveName]
	desc := dialog.Submission[SObjectiveDescription]
	measures := dialog.Submission[SObjectiveMeasures]
	targets := dialog.Submission[SObjectiveTargets]

	var editStatus = "created"
	var newSo *models.StrategyObjective
	var oldSo *models.StrategyObjective
	// If update, change existing record. Retrieve id from the dialog state.
	if msgState.Update {
		editStatus = "updated"
		var result models.StrategyObjective
		entity := StrategyEntityById(msgState.Id, teamID, strategyObjectivesTable)
		byt1, _ := json.Marshal(entity)
		err := json.Unmarshal(byt1, &result)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not decode interface to struct"))
		oldSo = &result
		newSo = &models.StrategyObjective{ID: msgState.Id, PlatformID: teamID.ToPlatformID(), Name: name,
			Description: desc, AsMeasuredBy: measures, Targets: targets, ObjectiveType: models.StrategyObjectiveType(tpe),
			Advocate: advocate, ExpectedEndDate: endDate, CreatedBy: userID, CreatedAt: TimestampFormat.Format(time.Now())}
	} else {
		id := core.Uuid()
		oldSo = nil
		newSo = &models.StrategyObjective{ID: id, PlatformID: teamID.ToPlatformID(), Name: name, Description: desc,
			AsMeasuredBy: measures, Targets: targets, ObjectiveType: models.StrategyObjectiveType(tpe),
			Advocate: advocate, ExpectedEndDate: endDate, CreatedBy: userID, CreatedAt: TimestampFormat.Format(time.Now())}
	}
	// Write entry to table
	err := d.PutTableEntry(*newSo, strategyObjectivesTable)
	core.ErrorHandler(err, namespace, "Could not put strategy objective to DB")
	attachsWithActions := ObjectiveViewAttachment(userID, *mc, newSo, oldSo,
		true, CustomerObjectiveEvent, true, teamID)
	attachsWithNoActions := ObjectiveViewAttachment(userID, *mc, newSo, oldSo,
		false, CustomerObjectiveEvent, true, teamID)
	// Publish the updated objective with actions into the user channel
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: msgState.ThreadTs,
		Attachments: attachsWithActions})
	// Do analysis on objective description
	utils.ECAnalysis(desc, stratObjDescriptionContext, "Customer objective", dialogTableName,
		mc.ToCallbackID(), userID, channelID, msgState.ThreadTs, msgState.ThreadTs, attachsWithActions,
		s, platformNotificationTopic, namespace)
	text := fmt.Sprintf("Below objective has been %s by <@%s>", editStatus, userID)
	// Post objectives with no actions to strategy community
	PostMsgToCommunity(community.Strategy, teamID, text, attachsWithNoActions)

	UserObjectiveFromStrategyObjective(newSo, core.EmptyString, teamID)

	utils.UpdateEngAsAnswered(mc.Source, mc.WithTarget(core.EmptyString).ToCallbackID(),
		engagementTable, d, namespace)
	return responses()
}
func onVisionEventCreateOrUpdateDialogSubmission(request slack.InteractionCallback, msgState MsgState, teamID models.TeamID, mc *models.MessageCallback) (resp []models.PlatformSimpleNotification) {
	dialog := request
	userID := dialog.User.ID
	channelID := dialog.Channel.ID
	vision := dialog.Submission[VisionDescription]
	advocate := dialog.Submission[VisionMissionAdvocate]

	var editStatus = "created"
	var newVm *models.VisionMission
	var oldVm *models.VisionMission
	// If update, change existing record. Retrieve id from the dialog state.
	if msgState.Update {
		editStatus = "updated"
		oldVm = StrategyVision(teamID)
		newVm = &models.VisionMission{ID: msgState.Id, PlatformID: teamID.ToPlatformID(), Vision: vision,
			Advocate: advocate, CreatedBy: userID, CreatedAt: time.Now().Format(string(TimestampFormat))}
	} else {
		id := core.Uuid()
		oldVm = nil
		newVm = &models.VisionMission{ID: id, PlatformID: teamID.ToPlatformID(), Vision: vision,
			Advocate: advocate, CreatedBy: userID, CreatedAt: time.Now().Format(string(TimestampFormat))}
	}
	err := d.PutTableEntry(*newVm, visionTable)
	core.ErrorHandler(err, namespace, "Could not put strategy objective to DB")
	attachsWithActions := VisionMissionViewAttachment(*mc, newVm, oldVm, true, true)
	attachsWithNoActions := VisionMissionViewAttachment(*mc, newVm, oldVm, false, true)
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: msgState.ThreadTs,
		Attachments: attachsWithActions})
	// Do analysis on vision
	utils.ECAnalysis(vision, visionDescriptionContext, string(VisionLabel), dialogTableName,
		mc.ToCallbackID(), userID, channelID, msgState.ThreadTs, msgState.ThreadTs, attachsWithActions,
		s, platformNotificationTopic, namespace)
	message := fmt.Sprintf("Below vision statement has been %s by <@%s>", editStatus, userID)
	PostMsgToCommunity(community.Strategy, teamID, message, attachsWithNoActions)
	stratComms := AllStrategyCommunities(teamID)
	// Posting to all of Strategy related communities (e.g., strategy, capability, initiative)
	for _, each := range stratComms {
		// Posting only to those channels into which Adaptive in invited
		if each.ChannelCreated == 1 {
			publish(models.PlatformSimpleNotification{UserId: userID, Channel: each.ChannelID,
				Message: message, Attachments: attachsWithNoActions})
		}
	}
	utils.UpdateEngAsAnswered(mc.Source, mc.WithTarget(core.EmptyString).ToCallbackID(),
		engagementTable, d, namespace)
	return responses()
}
func onCapabilityCommunityEventCreateOrUpdateDialogSubmission(request slack.InteractionCallback, msgState MsgState, teamID models.TeamID, mc *models.MessageCallback) (resp []models.PlatformSimpleNotification) {
	dialog := request
	userID := dialog.User.ID
	channelID := dialog.Channel.ID
	name := dialog.Submission[strategy.CapabilityCommunityName]
	desc := dialog.Submission[strategy.CapabilityCommunityDescription]
	advocate := dialog.Submission[strategy.CapabilityCommunityCoordinator]

	var editStatus = "created"
	var newCc *strategy.CapabilityCommunity
	var oldCc *strategy.CapabilityCommunity
	if msgState.Update {
		editStatus = "updated"
		var result strategy.CapabilityCommunity
		entity := StrategyEntityById(msgState.Id, teamID, capabilityCommunitiesTable)
		byt, _ := json.Marshal(entity)
		err := json.Unmarshal(byt, &result)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not decode interface to struct"))
		oldCc = &result
		newCc = &strategy.CapabilityCommunity{ID: msgState.Id, PlatformID: teamID.ToPlatformID(), Name: name,
			Description: desc, Advocate: advocate, CreatedBy: userID, CreatedAt: time.Now().Format(string(TimestampFormat))}
	} else {
		id := core.Uuid()
		oldCc = nil
		newCc = &strategy.CapabilityCommunity{ID: id, PlatformID: teamID.ToPlatformID(), Name: name,
			Description: desc, Advocate: advocate, CreatedBy: userID, CreatedAt: time.Now().Format(string(TimestampFormat))}
	}
	// Write entry to table
	err := d.PutTableEntry(*newCc, capabilityCommunitiesTable)
	attachsWithActions := strategy.CapabilityCommunityViewAttachment(*mc, newCc, oldCc, true)
	attachsWithNoActions := strategy.CapabilityCommunityViewAttachment(*mc, newCc, oldCc, false)
	// Publish to the user
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: msgState.ThreadTs,
		Attachments: attachsWithActions})
	// Do coach analysis
	// TODO: This should return an error back
	utils.ECAnalysis(desc, capCommunityDescriptionContext, "Objective Community", dialogTableName,
		mc.ToCallbackID(), userID, channelID, msgState.ThreadTs, msgState.ThreadTs, attachsWithActions,
		s, platformNotificationTopic, namespace)
	// Post to strategy community only if Adaptive is subscribed to it
	stratComm := StrategyCommunityByID(newCc.ID)
	// Check if Adaptive is associated with the community
	if stratComm.ChannelCreated == 1 {
		PostMsgToCommunity(community.Strategy, teamID,
			fmt.Sprintf("Below capability community has been %s by <@%s>",
				editStatus, userID), attachsWithNoActions)
	}
	// Post strategy community entity to the table
	// There is an index on channel id. Hence, it cannot be empty
	strComm := strategy.StrategyCommunity{ID: newCc.ID, PlatformID: teamID.ToPlatformID(), Advocate: advocate,
		Community: community.Capability, ChannelCreated: 0, ChannelID: "none",
		AccountabilityPartner: userID, ParentCommunity: community.Strategy, ParentCommunityChannelID: channelID,
		CreatedAt: time.Now().Format(string(TimestampFormat))}
	err = d.PutTableEntry(strComm, strategyCommunitiesTable)
	core.ErrorHandler(err, name, fmt.Sprintf("Could not add entry to %s table", strategyCommunitiesTable))
	// Post to the advocate only during the creation
	if !msgState.Update {
		PostMsgToAdvocate(advocate, userID, community.Capability, name)
	}
	utils.UpdateEngAsAnswered(mc.Source, mc.WithTarget(core.EmptyString).ToCallbackID(), engagementTable,
		d, namespace)
	return responses()
}
func onInitiativeSelectCommunityEventCreateOrUpdateDialogSubmission(request slack.InteractionCallback, msgState MsgState, teamID models.TeamID, mc *models.MessageCallback) (resp []models.PlatformSimpleNotification) {
	dialog := request
	userID := dialog.User.ID
	channelID := dialog.Channel.ID
	name := dialog.Submission[InitiativeName]
	desc := dialog.Submission[InitiativeDescriptionName]
	victory := dialog.Submission[InitiativeVictoryName]
	advocate := dialog.Submission[InitiativeAdvocateName]
	endDate := dialog.Submission[InitiateEndDateName]
	capObjective := dialog.Submission[InitiativeCapabilityObjectiveName]
	budget := dialog.Submission[InitiativeBudgetName]
	initCommID := msgState.SelectedOption

	var editStatus = "created"
	var newSi *models.StrategyInitiative
	var oldSi *models.StrategyInitiative
	if msgState.Update {
		editStatus = "updated"
		var result models.StrategyInitiative
		entity := StrategyEntityById(msgState.Id, teamID, strategyInitiativesTable)
		byt, _ := json.Marshal(entity)
		err := json.Unmarshal(byt, &result)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not decode interface to struct"))
		oldSi = &result
		newSi = &models.StrategyInitiative{ID: msgState.Id, PlatformID: teamID.ToPlatformID(), Name: name,
			Description: desc, DefinitionOfVictory: victory, Advocate: advocate, InitiativeCommunityID: initCommID,
			ExpectedEndDate: endDate, Budget: budget, CapabilityObjective: capObjective, CreatedBy: userID,
			CreatedAt: time.Now().Format(string(TimestampFormat))}
	} else {
		id := core.Uuid()
		oldSi = nil
		newSi = &models.StrategyInitiative{ID: id, PlatformID: teamID.ToPlatformID(), Name: name, Description: desc,
			DefinitionOfVictory: victory, Advocate: advocate, InitiativeCommunityID: initCommID,
			ExpectedEndDate: endDate, Budget: budget, CapabilityObjective: capObjective, CreatedBy: userID,
			CreatedAt: time.Now().Format(string(TimestampFormat))}
	}
	// Write entry to table
	err := d.PutTableEntry(*newSi, strategyInitiativesTable)
	attachsWithActions := InitiativeViewAttachment(userID, *mc.WithTopic(strategy.InitiativeEvent), newSi,
		oldSi, true, true, models.TeamID(teamID))
	attachsWithNoActions := InitiativeViewAttachment(userID, *mc.WithTopic(strategy.InitiativeEvent), newSi,
		oldSi, false, true, models.TeamID(teamID))
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: msgState.ThreadTs,
		Attachments: attachsWithActions})
	// Do analysis on vision
	utils.ECAnalysis(desc, initiativeDescriptionContext, "Strategy Initiative", dialogTableName,
		mc.ToCallbackID(), userID, channelID, msgState.ThreadTs, msgState.ThreadTs, attachsWithActions,
		s, platformNotificationTopic, namespace)

	text := fmt.Sprintf("Below strategy initiative has been %s by <@%s>", editStatus, userID)
	// Post to the strategy community
	PostMsgToCommunity(community.Strategy, teamID, text, attachsWithNoActions)
	// Post to associated initiative community with no actions
	initComm := StrategyCommunityByID(initCommID)
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: initComm.ChannelID,
		Message: text, Attachments: attachsWithNoActions})
	// Post to associated capability community id
	capObj := strategy.StrategyObjectiveByID(teamID, capObjective, strategyObjectivesTable)
	stratComm := StrategyCommunityByID(capObj.CapabilityCommunityIDs[0])
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: stratComm.ChannelID,
		Message: text, Attachments: attachsWithNoActions})

	// Accountability partner should be the advocate of the objective
	uObj := &models.UserObjective{UserID: newSi.Advocate, Name: newSi.Name, ID: newSi.ID,
		Description: newSi.Description, AccountabilityPartner: capObj.Advocate, Accepted: 1,
		ObjectiveType: StrategyDevelopmentObjective, StrategyAlignmentEntityID: newSi.InitiativeCommunityID,
		StrategyAlignmentEntityType: models.ObjectiveStrategyInitiativeAlignment, PlatformID: teamID.ToPlatformID(),
		CreatedDate:     core.NormalizeDate(newSi.CreatedAt),
		ExpectedEndDate: newSi.ExpectedEndDate}
	err = d.PutTableEntry(uObj, userObjectivesTable)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not write entry to %s table",
		userObjectivesTable))

	utils.UpdateEngAsAnswered(mc.Source, mc.WithTarget(core.EmptyString).ToCallbackID(),
		engagementTable, d, namespace)
	return responses()
}
func onInitiativeCommunityEventCreateOrUpdateDialogSubmission(request slack.InteractionCallback, msgState MsgState, teamID models.TeamID, mc *models.MessageCallback) (resp []models.PlatformSimpleNotification) {
	dialog := request
	userID := dialog.User.ID
	channelID := dialog.Channel.ID
	name := dialog.Submission[strategy.InitiativeCommunityName]
	desc := dialog.Submission[strategy.InitiativeCommunityDescription]
	advocate := dialog.Submission[strategy.InitiativeCommunityCoordinator]
	capCommId := dialog.Submission[strategy.InitiativeCommunityCapabilityCommunity]

	// var editStatus string
	var oldSic *strategy.StrategyInitiativeCommunity
	var id string
	if msgState.Update {
		// editStatus = "updated"
		id = msgState.Id
		result := StrategyInitiativeCommunityByID(msgState.Id, teamID)
		oldSic = &result
	} else {
		// editStatus = "created"
		id = core.Uuid()
		oldSic = nil
	}
	newSic := &strategy.StrategyInitiativeCommunity{ID: id, PlatformID: teamID.ToPlatformID(), Name: name,
		Description: desc, Advocate: advocate, CapabilityCommunityID: capCommId, CreatedBy: userID,
		CreatedAt: time.Now().Format(string(TimestampFormat))}
// Write entry to table
	err2 := d.PutTableEntry(*newSic, strategyInitiativeCommunitiesTable)
	core.ErrorHandler(err2, name, fmt.Sprintf("Could not add entry to %s table: %v", strategyInitiativeCommunitiesTable, *newSic))
	attachsWithActions := strategy.InitiativeCommunityViewAttachmentEditable(*mc, newSic, oldSic,
		capabilityCommunitiesTable, strategyInitiativesTable, strategyInitiativesPlatformIndex)
	attachsWithNoActions := strategy.InitiativeCommunityViewAttachmentReadOnly(*mc, newSic, oldSic,
		capabilityCommunitiesTable)
	// Publish to the user
	publish(models.PlatformSimpleNotification{UserId: request.User.ID, Channel: request.Channel.ID, Ts: msgState.ThreadTs,
		Attachments: attachsWithActions})
	// Do coach analysis
	utils.ECAnalysis(desc, initiativeCommDescriptionContext, "Initiative Community", dialogTableName,
		mc.ToCallbackID(), request.User.ID, request.Channel.ID, msgState.ThreadTs, msgState.ThreadTs, attachsWithActions,
		s, platformNotificationTopic, namespace)
	PostMsgToCommunity(community.Strategy, teamID,
		fmt.Sprintf("Below initiative community has been updated by <@%s>",
			request.User.ID), attachsWithNoActions)
	// if msgState.Update {
	// 	// Post creation to strategy communities only if Adaptive has been associated with this community
	// 	stratComm := StrategyCommunityByID(newSic.ID)
	// 	if stratComm.ChannelCreated == 1 {
	// 		// PostMsgToCommunity(stratComm.ChannelID, teamID,
	// 		// 	fmt.Sprintf("Below initiative community has been updated by <@%s>",
	// 		// 		request.User.ID), attachsWithNoActions)
	// 	}
	// }
	// Post strategy community entity to the table
	strComm := strategy.StrategyCommunity{ID: newSic.ID, PlatformID: teamID.ToPlatformID(), Advocate: advocate,
		Community: community.Initiative, ChannelCreated: 0, ChannelID: "none",
		AccountabilityPartner: request.User.ID, ParentCommunity: community.Strategy, ParentCommunityChannelID: channelID,
		CreatedAt: time.Now().Format(string(TimestampFormat))}
	err3 := d.PutTableEntry(strComm, strategyCommunitiesTable)
	core.ErrorHandler(err3, name, fmt.Sprintf("Could not add entry to %s table: %v", strategyCommunitiesTable, strComm))
	if !msgState.Update {
		// Post to the advocate only during creation
		PostMsgToAdvocate(advocate, request.User.ID, community.Initiative, name)
	}
	utils.UpdateEngAsAnswered(mc.Source, mc.WithTarget(core.EmptyString).ToCallbackID(), engagementTable,
		d, namespace)
	return
}

func onObjectiveCommunityAssociationSelectObjectiveCreateDialogSubmission(dialog slack.InteractionCallback, msgState MsgState, teamID models.TeamID, mc *models.MessageCallback) []models.PlatformSimpleNotification {
	// source := dialog.Submission[StrategyAssociationSourceName]
	communityID := dialog.Submission[StrategyObjectiveAssociationCommunity]
	description := dialog.Submission[StrategyObjectiveCommunityAssociationDescription]

	capObjID := msgState.SelectedOption
	capObj := strategy.StrategyObjectiveByID(teamID, capObjID, strategyObjectivesTable)
	capComm := strategy.CapabilityCommunityByID(teamID, communityID, capabilityCommunitiesTable)
	existingCommunities := capObj.CapabilityCommunityIDs

	if !core.ListContainsString(existingCommunities, communityID) {
		updatedCommunityList := append(existingCommunities, communityID)
		updateObjectiveCommunities(teamID, capObjID, updatedCommunityList)
	}

	// add the association
	if msgState.Update {
		// delete the original association
	}

	newSa := &StrategyObjectiveCommunityAssociation{ObjectiveID: capObjID, ObjectiveName: capObj.Name,
		CommunityID: communityID, CommunityName: capComm.Name, Description: description}
	target := fmt.Sprintf("%s_%s", capObjID, communityID)
	attachs := StrategyObjectiveCommunityAssociationViewAttachment(*mc.WithTarget(target), newSa, true, false)
	return responses(models.PlatformSimpleNotification{UserId: dialog.User.ID,
		Channel: dialog.Channel.ID, Attachments: attachs, Ts: msgState.ThreadTs})
}
func onDialogCancellation(np models.NamespacePayload4) {

}

func updateObjectiveCommunities(teamID models.TeamID, capObjID string, communities []string) {
	keyParams := map[string]*dynamodb.AttributeValue{
		"id":          dynString(capObjID),
		"platform_id": dynString(teamID.ToString()),
	}
	exprAttributes := map[string]*dynamodb.AttributeValue{
		":ccs": dynListString(communities),
	}
	updateExpression := "set capability_community_ids = :ccs"
	err := d.UpdateTableEntry(exprAttributes, keyParams, updateExpression, strategyObjectivesTable)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not update capability communities list in %s table",
		strategyObjectivesTable))
}

func RemoveFromSliceOnIndex(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

var (
	StrategyDevelopmentObjective models.DevelopmentObjectiveType = "strategy"
)

// userCapabilityCommunities queries strategy_community for 
// all communities for the user.
// In case user is in strategy group, returns all communities.
func userCapabilityCommunities(userID string, typ community.AdaptiveCommunity, teamID models.TeamID) (res []models.KvPair) {
	var comms []strategy.CapabilityCommunity
	if isMemberInCommunity(userID, community.Strategy) {
		comms = AllCapabilityCommunities(teamID)
	} else {
		commUsers := strategy.QueryCommunityUserIndex(userID, communityUsersTable, communityUsersUserIndex)
		for _, each := range commUsers {
			splits := strings.Split(each.CommunityId, ":")
			if len(splits) == 2 {
				if splits[0] == string(typ) {
					switch typ {
					case community.Capability:
						capComm := strategy.CapabilityCommunityByID(teamID, splits[1], capabilityCommunitiesTable)
						comms = append(comms, capComm)
					}
				}
			}
		}
	}
	for _, comm := range comms {
		res = append(res, convertCapabilityCommunityToKvPair(comm))
	}
	return res
}

func convertCapabilityCommunityToKvPair(comm strategy.CapabilityCommunity) models.KvPair {
	return models.KvPair{Key: comm.Name, Value: comm.ID}
}

func UserObjectiveFromStrategyObjective(newSo *models.StrategyObjective, commID string,
	teamID models.TeamID) *models.UserObjective {
	logger.Infof("Adding strategy objective as a user objective with id: %s", newSo.ID)
	vision := StrategyVision(teamID)
	id := newSo.ID//core.IfThenElse(commID != core.EmptyString, fmt.Sprintf("%s_%s", newSo.ID, commID), newSo.ID).(string)
	// We are using _ here because `:` will create issues with callback
	createdDate := core.NormalizeDate(newSo.CreatedAt)
	strategyAlignmentEntityID := ""
	if len(newSo.CapabilityCommunityIDs) > 0 {
		strategyAlignmentEntityID = newSo.CapabilityCommunityIDs[0]
	}
	uObj := models.UserObjective{UserID: newSo.Advocate, Name: newSo.Name, ID: id, Description: newSo.Description,
		AccountabilityPartner: vision.Advocate, Accepted: 1,
		ObjectiveType: StrategyDevelopmentObjective, 
		StrategyAlignmentEntityID: strategyAlignmentEntityID,
		StrategyAlignmentEntityType: models.ObjectiveStrategyObjectiveAlignment, 
		PlatformID: teamID.ToPlatformID(),
		CreatedDate: createdDate,
		ExpectedEndDate: newSo.ExpectedEndDate}
	return &uObj
}
