package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	"github.com/adaptiveteam/adaptive/lambdas/feedback-report-posting-lambda-go"
	daosUser "github.com/adaptiveteam/adaptive/daos/user"
	"github.com/adaptiveteam/adaptive/lambdas/feedback-reporting-lambda-go"
	"context"
	"encoding/json"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
	"time"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
)

const (
	DateFormat = string(core.ISODateLayout)
)

const (
	GenerateReportHR       = "generate_report"
	FetchReportHR          = "fetch_report"
	SimulateDateFieldID    = "simulate_date"
	SimulateUserFieldID    = "simulate_user"
	CurrentQuarterSchedule = "current_quarter_schedule"
	NextQuarterSchedule    = "next_quarter_schedule"
	SelectCoachee          = "select_coachee"
	RequestCoach           = "request_coach"
	AdaptiveAccess         = "adaptive_access"
	CoachConfirm           = "coach_confirm"

	format = "Monday, January _2"
)

var (
	commentsSurvey = CommentsSurvey(CoachingLabel, CoachRejectionReasonLabel, CommentsName)
	logger         = alog.LambdaLogger(logrus.InfoLevel)
)

type MsgState struct {
	ThreadTs    string `json:"thread_ts"`
	ObjectiveId string `json:"objective_id,omitempty"`
}

func surveyState(message slack.InteractionCallback, target string) func() string {
	return func() string {
		// When the original message is from a thread, we need to post to the same thread
		// Below logic checks if the incoming message is from a thread
		var ts string
		if message.OriginalMessage.ThreadTimestamp == "" {
			ts = message.MessageTs
		} else {
			ts = message.OriginalMessage.ThreadTimestamp
		}
		fmt.Println("### dialog callback id: " + message.CallbackID)
		msgStateBytes, err := json.Marshal(MsgState{ThreadTs: ts, ObjectiveId: target})
		core.ErrorHandler(err, namespace, "Could not marshal MsgState")
		return string(msgStateBytes)
	}
}

func updateCommentsAttachmentActions(mc models.MessageCallback) []ebm.AttachmentAction {
	return []ebm.AttachmentAction{*models.GenAttachAction(mc,
		models.Update, "I would like to change this",
		models.EmptyActionConfirm(), true)}
}

func viewCommentsAttachment(mc models.MessageCallback, title, comments ui.RichText) []ebm.Attachment {
	attach := utils.ChatAttachment(string(title), "", "", mc.ToCallbackID(), updateCommentsAttachmentActions(mc),
		[]ebm.AttachmentField{
			{
				Title: string(CommentsLabel),
				Value: string(comments),
				Short: true,
			},
		}, time.Now().Unix())
	return []ebm.Attachment{*attach}
}

func communitiesMapToStrings(comms []models.AdaptiveCommunity) []string {
	var communitiesMapToStrings []string
	for _, each := range comms {
		communitiesMapToStrings = append(communitiesMapToStrings, each.ID)
	}
	return communitiesMapToStrings
}

func HandleRequest(ctx context.Context, e events.SNSEvent) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("Error in community-slack-processor-lambda %v", err2)
		}
	}()
	logger = logger.WithLambdaContext(ctx)
	for _, record := range e.Records {
		fmt.Printf("HandleRequest: %v\n", record)
		message := record.SNS.Message
		if message == "warmup" {
			log.Println("Warmed up...")
		} else {
			// models.ParseEventsAPIEventAsSlackRequestUnsafe()
			var np models.NamespacePayload4
			err := json.Unmarshal([]byte(message), &np)
			core.ErrorHandler(err, namespace, fmt.Sprintf("Could not unmarshal sns record to NamespacePayload4"))
			conn := connGen.ForPlatformID(np.TeamID.ToPlatformID())

			// This module only looks for payload with 'app-mention' namespace
			if np.Namespace == "community" {
				// Handling event_callback messages
				if np.SlackRequest.Type == models.EventsAPIEventSlackRequestType {
					logger.WithField("app_mention_event", np).Info()
					dispatchAppMentionSlackEvent(np.ToEventsAPIEventUnsafe(), conn)
				} else {
					request := np.SlackRequest.InteractionCallback
					if np.TeamID.IsEmpty() {
						logger.WithField("namespace", namespace).Errorf("Platform id empty in incoming message: %s", message)
					} else {
						switch np.SlackRequest.Type {
						case models.InteractionSlackRequestType:
							logger.WithField("interaction_callback_event", np).Info()
							dispatchCommunityInteractionCallback(request, np.TeamID, conn)
						case models.DialogSubmissionSlackRequestType:
							logger.WithField("dialog_submission_event", np).Info()
							// Handling dialog submission for each answer
							dispatchCommunityDialogSubmission(request, np.TeamID)
						}
					}
				}
			} else if np.Namespace == "adaptive-channel" {
				logger.WithField("adaptive_channel_event", np).Info()
				adaptiveChannelNamespaceEventHandler(np.ToEventsAPIEventUnsafe(), 
					np.TeamID, conn)
			}
		}
	}
	return err
}

func dispatchCommunityMenuAction2(request slack.InteractionCallback,
	teamID models.TeamID,
	selectedMenuItem string,
	conn daosCommon.DynamoDBConnection,
) {
	mc := callback(request.User.ID, "init", "select")
	switch selectedMenuItem {
	case RequestCoach:
		response := onRequestCoachClicked(request, mc, conn)
		respond(teamID, response)
	case GenerateReportHR:
		generateReportMenuHandler(request, mc, conn)
	case FetchReportHR:
		fetchReportMenuHandler(request, mc, conn)
	case SimulateCurrentQuarterAction:
		simulateCurrentQuarterMenuHandler(request, mc, conn)
	case SimulateNextQuarterAction:
		simulateNextQuarterMenuHandler(request, mc, conn)
	case CurrentQuarterSchedule:
		currentQuarterScheduleMenuHandler(request, conn)
	case NextQuarterSchedule:
		nextQuarterScheduleMenuHandler(request, conn)
	case CommunitySubscribeAction:
		response := onCommunitySubscribeClicked(request, teamID)
		respond(teamID, response)
	case CommunityUnsubscribeAction:
		response := onCommunityUnsubscribeClicked(request, teamID, conn)
		respond(teamID, response)
	}
}
func dispatchCommunityInteractionCallback(request slack.InteractionCallback,
	teamID models.TeamID, 
	conn daosCommon.DynamoDBConnection,
) {
	action := *request.ActionCallback.AttachmentActions[0]
	logger.WithField("event", "interaction_callback").WithField("action", action.Name).
		WithField("platform", teamID).Info("Handling Callback event")
	// This is for the options presented to the user
	if action.Name == CommunityMenuActionName {
		selected := action.SelectedOptions[0]
		selectedValue := selected.Value
		fmt.Printf("Interaction callback handling 4. selectedValue=%s\n", selectedValue)
		// community:<user>:init:select::4:2019
		dispatchCommunityMenuAction2(request, teamID, selectedValue, conn)
	} else {
		// Parse callback Id to messageCallback
		mc, err := utils.ParseToCallback(request.CallbackID)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))
		fmt.Printf("Interaction callback handling 2. mc.Module=%s, mc.Topic=%s, mc.Action=%s\n", mc.Module, mc.Topic, mc.Action)
		if mc.Module == "community" {
			if mc.Topic == "init" {
				if mc.Action == "select" {
					selected := action.SelectedOptions[0]
					selectedValue := selected.Value
					fmt.Printf("Interaction callback handling 3. selectedValue=%s\n", selectedValue)
					switch selectedValue {
					case CommunitySubscribeAction:
						onCommunitySubscribeClicked(request, teamID)
					case CommunityUnsubscribeAction:
						response := onCommunityUnsubscribeClicked(request, teamID, conn)
						respond(teamID, response)
					}
				}
			} else if mc.Topic == "reports" {
				action := *request.ActionCallback.AttachmentActions[0]
				if strings.HasPrefix(mc.Action, GenerateReportHR) {
					// generate report actions: now and cancel
					suffixAction := strings.TrimPrefix(action.Name, GenerateReportHR+"_")
					communityNamespaceReportsGenerateReportCallback(request, suffixAction, action, *mc)
				} else if strings.HasPrefix(mc.Action, FetchReportHR) {
					// fetch report actions: now and cancel
					suffixAction := strings.TrimPrefix(action.Name, FetchReportHR+"_")
					communityNamespaceReportsFetchReportCallback(request, teamID, suffixAction, action, *mc)
				}
			} else if mc.Topic == CoachingName {
				action := *request.ActionCallback.AttachmentActions[0]

				if strings.HasPrefix(mc.Action, CoachConfirm) {

					suffixAction := strings.TrimPrefix(action.Name, CoachConfirm+"_")
					communityNamespaceCoachingCoachConfirmCallback(request, suffixAction, action, *mc)

				} else if strings.HasPrefix(mc.Action, RequestCoach) {
					suffixAction := strings.TrimPrefix(action.Name, RequestCoach+"_")

					communityNamespaceCoachingRequestCoachCallback(request, suffixAction, action, *mc, teamID)
				}
			} else if mc.Topic == string(community.Admin) {
				if strings.HasPrefix(mc.Action, AdaptiveAccess) {
					suffixAction := strings.TrimPrefix(action.Name, AdaptiveAccess+"_")
					communityNamespaceAdminAccessCallback(request, suffixAction, conn)
				}
			} else if mc.Action == "select" {
				action := *request.ActionCallback.AttachmentActions[0]
				if mc.Topic == "subscription" {
					if action.Name == "select" {
						value := action.SelectedOptions[0].Value
						logger.Infof("Subscribing %s channel with Adaptive", value)
						onCommunitySubscribeCommunityClicked(request, value, *mc, teamID, conn)
					} else if action.Name == "cancel" {
						dispatchCommunityNamespaceMenuSelectSubscriptionCancel(request)
					}
				} else if mc.Topic == "unsubscription" {
					if action.Name == "select" {
						message := onCommunityUnsubscribeCommunityClicked(request, action.SelectedOptions[0].Value,
							*mc, conn)
						replyReplace(request, teamID, message)
					} else if action.Name == "cancel" {
						message := onCommunityUnsubscribeCancelled(request, *mc)
						replyReplace(request, teamID, message)
					}
				}
			}
		}
	}
}

func dispatchCommunityDialogSubmission(dialog slack.InteractionCallback, teamID models.TeamID) {
	// Parse callback Id to messageCallback
	mc, err := utils.ParseToCallback(dialog.CallbackID)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))

	var msgState MsgState
	err = json.Unmarshal([]byte(dialog.State), &msgState)

	// if mc.Topic == CoachingName {
	// 	communityNamespaceCoachingDialogSubmissionHandler(dialog, msgState, *mc, dialog.Submission)
	// } else 
	if mc.Topic == "init" {
		dispatchCommunitySimulateDialogSubmission(dialog, msgState, *mc, dialog.Submission)
	} else {
		logger.Errorf("Unhandled mc.Topic=%s", mc.Topic)
	}
}

func communityNamespaceReportsGenerateReportCallback(request slack.InteractionCallback, suffixAction string, action slack.AttachmentAction, mc models.MessageCallback) {
	userID := request.User.ID
	channelID := request.Channel.ID
	switch suffixAction {
	case string(models.Now):
		target := action.SelectedOptions[0].Value
		// Posting message to the channel in which user requested this
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
			Message: string(GeneratingReportForUserNotification(target)),
			Ts:      request.MessageTs, Attachments: models.EmptyAttachs(), AsUser: true})
		err2 := feedbackReportingLambda.GeneratePerformanceReportAndPostToUserAsync(userID, time.Now())

		core.ErrorHandler(err2, namespace, "Could not invoke GeneratePerformanceReportAndPostToUserAsync")
	case string(models.Cancel):
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Message: "", AsUser: true, Ts: request.MessageTs})
	}

}

func communityNamespaceReportsFetchReportCallback(request slack.InteractionCallback, teamID models.TeamID, suffixAction string, action slack.AttachmentAction, mc models.MessageCallback) {
	userID := request.User.ID
	channelID := request.Channel.ID
	switch suffixAction {
	case string(models.Now):
		targetUserID := action.SelectedOptions[0].Value
		// Posting message to the channel in which user requested this
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
			Message: string(FetchingReportForUserNotification(targetUserID)),
			Ts:      request.OriginalMessage.Timestamp, Attachments: models.EmptyAttachs(), AsUser: true})
		err2 := feedbackReportPostingLambda.DeliverReportToUserAsync(teamID, userID, time.Now())
		core.ErrorHandler(err2, namespace, "Could not DeliverReportToUserAsync")
	case string(models.Cancel):
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Message: "", AsUser: true, Ts: request.OriginalMessage.Timestamp})
	}
}

func communityNamespaceCoachingCoachConfirmCallback(request slack.InteractionCallback, suffixAction string, action slack.AttachmentAction, mc models.MessageCallback) {
	userId := request.User.ID
	channelId := request.Channel.ID

	switch suffixAction {
	case string(models.Now):
		// Posting message to the channel in which user requested this
		// Write coaching relationship to table
		year, quarter := core.CurrentYearQuarter()
		cqy := fmt.Sprintf("%s:%d:%d", mc.Source, quarter, year)
		ceqy := fmt.Sprintf("%s:%d:%d", mc.Target, quarter, year)
		coachRel := models.CoachingRelationship{CoachQuarterYear: cqy, CoacheeQuarterYear: ceqy, Coachee: mc.Target, Quarter: quarter, Year: year, CoachRequested: true}
		err := d.PutTableEntry(coachRel, coachingRelationshipsTable)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not write entry to %s table", coachingRelationshipsTable))
		// send a direct notification to the coach
		publish(models.PlatformSimpleNotification{UserId: mc.Source, AsUser: true,
			Message: string(CoachingRequestConfirmToCoachNotification(mc.Target, quarter, year))})
		// send notification to the coachee
		publish(models.PlatformSimpleNotification{UserId: userId, Channel: channelId, AsUser: true,
			Message: string(CoachingRequestConfirmToCoacheeNotification(mc.Source, quarter, year))})
		// Update engagement as answered. This engagement has been posted to target for confirmation.
		utils.UpdateEngAsAnswered(mc.Target, mc.ToCallbackID(), engagementTable, d, namespace)

		// Delete the original engagement from coachee's chat
		DeleteOriginalEng(userId, channelId, request.MessageTs)
	case string(models.Ignore):
		// Coachee decided not to accept the current coach
		id := action.Value // callbackId
		val := commentsSurvey
		ut := userTokenSyncUnsafe(request.User.ID)
		api := slack.New(ut)
		err := utils.SlackSurvey(api, request, val, id, surveyState(request, mc.Target))
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+request.CallbackID))
		// Mark engagement as ignored
		utils.UpdateEngAsIgnored(mc.Target, mc.ToCallbackID(), engagementTable, d, namespace)
	case string(models.Update):
		// This case is when a coachee entered survey information for not accepting the coach
		id := action.Value
		// User chose to update the text enter through earlier dialog interaction
		val := models.FillSurvey(commentsSurvey, map[string]string{
			CommentsName: utils.SlackFieldValue(request.OriginalMessage.Attachments[0], string(CommentsLabel))})
		ut := userTokenSyncUnsafe(request.User.ID)
		api := slack.New(ut)
		// Open a survey associated with the engagement
		err := utils.SlackSurvey(api, request, val, id, surveyState(request, mc.Target))
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+request.CallbackID))
	}

}

var (
	coachingSurvey = ebm.AttachmentActionSurvey{
		Title:       string(CoachingDialogTitle),
		SubmitLabel: models.SubmitLabel,
		Elements: []ebm.AttachmentActionTextElement{
			ebm.NewTextArea(CommentsName, CoachingRejectionExplanationPrompt, CommentsPlaceholder, ""),
		},
	}
)

func communityNamespaceCoachingRequestCoachCallback(request slack.InteractionCallback,
	suffixAction string, action slack.AttachmentAction,
	mc models.MessageCallback, teamID models.TeamID) {
	switch suffixAction {
	case string(models.Now):
		requestedCoach := action.SelectedOptions[0].Value
		// Send an engagement to coach asking if he is ok to coach
		ConfirmCoachRequestEngagement(teamID, *mc.WithTarget(requestedCoach), requestedCoach, "",
			string(CoachingRequestTitle(request.User.ID)), "", "", true)

		publish(models.PlatformSimpleNotification{UserId: request.User.ID, Channel: request.Channel.ID,
			Message: string(CoachingRequestSentNotificationToCoacheeTitle(requestedCoach)),
			AsUser:  true})
		// Delete original engagement
		DeleteOriginalEng(request.User.ID, request.Channel.ID, request.MessageTs)
		// Mark engagement as answered
		utils.UpdateEngAsAnswered(request.User.ID, mc.WithTarget("").ToCallbackID(), engagementTable, d, namespace)
	case string(Confirm):
		// Requested coach has agreed to coach the coachee. Add an entry in relationships table
		year, quarter := core.CurrentYearQuarter()
		// coach is from mc.Target and coachee is from mc.Source
		cqy := fmt.Sprintf("%s:%d:%d", mc.Target, quarter, year)
		ceqy := fmt.Sprintf("%s:%d:%d", mc.Source, quarter, year)
		coachRel := models.CoachingRelationship{CoachQuarterYear: cqy, CoacheeQuarterYear: ceqy, Coachee: mc.Source, Quarter: quarter, Year: year, CoacheeRequested: true}
		err := d.PutTableEntry(coachRel, coachingRelationshipsTable)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not write entry to %s table", coachingRelationshipsTable))
		// Send notification to coachee
		publish(models.PlatformSimpleNotification{UserId: mc.Source,
			Message: string(CoachingRequestConfirmationToCoachee(request.User.ID, quarter, year)),
			AsUser:  true})
		// Send notification to coach
		publish(models.PlatformSimpleNotification{UserId: mc.Target,
			Message: string(CoachingRequestConfirmationToCoach(mc.Source, quarter, year)),
			AsUser:  true})
		// Update engagement as answered
		utils.UpdateEngAsAnswered(mc.Target, mc.ToCallbackID(), engagementTable, d, namespace)
		DeleteOriginalEng(request.User.ID, request.Channel.ID, request.MessageTs)
	case string(models.Cancel):
		DeleteOriginalEng(request.User.ID, request.Channel.ID, request.MessageTs)
		// Update engagement as ignored
		utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
	case string(No):
		id := mc.ToCallbackID()
		ut := userTokenSyncUnsafe(request.User.ID)
		api := slack.New(ut)

		// Open a survey associated with the engagement
		err := dialogFromSurvey(api, request, coachingSurvey, id, mc.Target)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+request.CallbackID))
	case string(models.Update):
		// This case is when coach entered comments for not accepting coachee
		id := mc.ToCallbackID()
		val := models.FillSurvey(coachingSurvey, map[string]string{CommentsName: utils.SlackFieldValue(request.OriginalMessage.Attachments[0], string(CommentsLabel))})
		ut := userTokenSyncUnsafe(request.User.ID)
		api := slack.New(ut)
		// Open a survey associated with the engagement
		err := dialogFromSurvey(api, request, val, id, mc.Target)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+request.CallbackID))
	}

}

func communityNamespaceAdminAccessCallback(request slack.InteractionCallback, suffixAction string, conn daosCommon.DynamoDBConnection) {
	switch suffixAction {
	case string(models.Now):
		publish(models.PlatformSimpleNotification{UserId: request.User.ID, Channel: request.Channel.ID,
			Message: string(AdminRequestSentAcknowledgement)})
		userID := request.User.ID
		user, err2 := daosUser.Read(userID)(conn)
		core.ErrorHandler(err2, "communityNamespaceAdminAccessCallback", "userDAO.Read")
		teamID := models.ParseTeamID(user.PlatformID)
		// ut := userTokenSyncUnsafe(request.User.ID)
		//communityDAO.ReadByIDUnsafe(ut.PlatformID, string(community.Admin))
		adminComm := adaptiveCommunity.ReadUnsafe(teamID.ToPlatformID(), string(community.Admin))(conn)
		userComm := adaptiveCommunity.ReadUnsafe(teamID.ToPlatformID(), string(community.User))(conn)
		thisUser := ui.RichText("<@" + request.User.ID + ">")
		commInfo := core.IfThenElse(userComm.ChannelID == core.EmptyString,
			NoUserCommunityErrorMessage(thisUser),
			NoUserInviteToUserCommunityErrorMessage(thisUser)).(ui.RichText)
		publish(models.PlatformSimpleNotification{UserId: request.User.ID, Channel: adminComm.ChannelID,
			Message: string(UserRequestAdaptiveAccessNotifitationToAdminCommunity(userID, string(commInfo)))})
	case string(models.Ignore):
		publish(models.PlatformSimpleNotification{UserId: request.User.ID, Channel: request.Channel.ID,
			Message: string(FarewellNotification)})
	}
	publish(models.PlatformSimpleNotification{UserId: request.User.ID, Channel: request.Channel.ID, AsUser: true, Ts: request.MessageTs})
}
func dispatchCommunityNamespaceMenuSelectSubscriptionCancel(request slack.InteractionCallback) {
	// User has decided to cancel this. This means, we remove this engagement
	// For now, we mark this as answered
	// TODO: We need to handle cases where a user ignores an engagement.
	// This is different from not being reminded any more
	publish(overrideRequestMessageAsUser(request, SubscriptionCancellationConfirmation))
}

const (
	StrategyDevelopmentObjective models.DevelopmentObjectiveType = "strategy"
)

func dispatchCommunitySimulateDialogSubmission(dialog slack.InteractionCallback, msgState MsgState, mc models.MessageCallback, form map[string]string) {
	switch mc.Action {
	case SimulateNextQuarterAction, SimulateCurrentQuarterAction:
		emulUser := form[SimulateUserFieldID]
		emulDate := form[SimulateDateFieldID]
		date, err := time.Parse(format, emulDate)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse date"))
		// TODO: Remove the year hard-coding
		parsedDate := business_time.NewDate(2019, int(date.Month()), date.Day()).DateToString(DateFormat)
		publish(models.PlatformSimpleNotification{
			UserId:  dialog.User.ID,
			Channel: dialog.Channel.ID,
			Ts:      msgState.ThreadTs,
			Message: string(SimulationUserDateNotificatoin(emulDate, emulUser))},
		)
		engage := models.UserEngage{
			UserID: emulUser,
			Date: parsedDate,
		}
		engScheduleBytes, _ := json.Marshal(engage)
		_, _ = lambdaAPI.InvokeFunction(engagementSchedulerLambda, engScheduleBytes, true)
	}
}

func adaptiveChannelNamespaceEventHandler(eventsAPIEvent slackevents.EventsAPIEvent, 
	teamID models.TeamID,
	conn daosCommon.DynamoDBConnection,
) {
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		// For invite and remove events
		eventType := eventsAPIEvent.InnerEvent.Type
		logger.WithField("callback_event_type", eventsAPIEvent.InnerEvent.Type).Info()
		switch eventType {
		case slackevents.MemberJoinedChannel:
			slackMsg := *eventsAPIEvent.InnerEvent.Data.(*slackevents.MemberJoinedChannelEvent)
			onMemberJoinedChannel(slackMsg, conn)
		case "member_left_channel":
			slackMsg := *eventsAPIEvent.InnerEvent.Data.(*slack.MemberLeftChannelEvent)
			onMemberLeftChannel(teamID, slackMsg)
		case "group_left":
			// This is when Adaptive leaves a private channel
			cbEvent := *eventsAPIEvent.Data.(*slackevents.EventsAPICallbackEvent)
			// slack.GroupLeftEvent doesn't populate user id. Do not use that field.
			onGroupLeftEvent(teamID, cbEvent, conn)
		case "channel_deleted": // docs: https://api.slack.com/events/channel_deleted
			channelDeletedEvent := *eventsAPIEvent.InnerEvent.Data.(*slack.ChannelDeletedEvent)
			channelUnsubscribeUnsafe(channelDeletedEvent.Channel, conn)
		case "channel_archive":
			channelDeletedEvent := *eventsAPIEvent.InnerEvent.Data.(*slack.ChannelArchiveEvent)
			channelUnsubscribeUnsafe(channelDeletedEvent.Channel, conn)
		case "group_deleted":
			channelDeletedEvent := *eventsAPIEvent.InnerEvent.Data.(*slack.GroupCloseEvent)
			channelUnsubscribeUnsafe(channelDeletedEvent.Channel, conn)
		case "group_archive":
			channelDeletedEvent := *eventsAPIEvent.InnerEvent.Data.(*slack.GroupArchiveEvent)
			channelUnsubscribeUnsafe(channelDeletedEvent.Channel, conn)
		default:
			logger.Warnf("Unhandled %s event type", eventType)
		}
		logger.Infof("Handling of %s completed", eventType)
	} else {
		logger.Warnf("Unsupported eventsAPIEvent.Type %s", eventsAPIEvent.Type)
	}
}

func dialogFromSurvey(api *slack.Client, message slack.InteractionCallback, survey ebm.AttachmentActionSurvey, id string, objectiveId string) error {
	survState := func() string {
		// When the original message is from a thread, we need to post to the same thread
		// Below logic checks if the incoming message is from a thread
		var ts string
		if message.OriginalMessage.ThreadTimestamp == "" {
			ts = message.MessageTs
		} else {
			ts = message.OriginalMessage.ThreadTimestamp
		}
		fmt.Println("### dialog callback id: " + message.CallbackID)
		msgStateBytes, err := json.Marshal(MsgState{ThreadTs: ts, ObjectiveId: objectiveId})
		core.ErrorHandler(err, namespace, "Could not marshal MsgState")
		return string(msgStateBytes)
	}

	return utils.SlackSurvey(api, message, survey, id, survState)
}

func ConfirmCoachRequestEngagement(teamID models.TeamID, mc models.MessageCallback, coach, title, text,
	fallback, learnLink string, urgent bool) {
	utils.AddChatEngagement(mc, title, text, fallback, coach, requestCoachAttachmentActions(mc, learnLink),
		[]ebm.AttachmentField{}, teamID, urgent, engagementTable, d, namespace, time.Now().Unix(),
		models.UserEngagementCheckWithValue{})
}

func requestCoachAttachmentActions(mc models.MessageCallback, learnTrailPath string) []ebm.AttachmentAction {
	return models.AppendOptionalAction(
		[]ebm.AttachmentAction{
			*models.GenAttachAction(mc, Confirm, "Yes", ebm.AttachmentActionConfirm{
				OkText:      models.YesLabel,
				DismissText: models.CancelLabel,
			}, false),
			*models.GenAttachAction(mc, No, "No", models.EmptyActionConfirm(), true),
		},
		models.LearnMoreAction(models.ConcatPrefixOpt("docs/general/", learnTrailPath)))
}
