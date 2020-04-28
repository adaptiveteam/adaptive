package feedbackSetupLambda

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	evalues "github.com/adaptiveteam/adaptive/adaptive-engagements/values"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	daosUser "github.com/adaptiveteam/adaptive/daos/user"
	"github.com/adaptiveteam/adaptive/daos/userFeedback"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	feedbackReportPostingLambda "github.com/adaptiveteam/adaptive/lambdas/feedback-report-posting-lambda-go"
	ls "github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

const (
	Request     = "request"
	Ask         = "ask"
	MoreDetails = "more_details"
	LessDetails = "less_details"
)

var (
	// Dialog element names
	Feedback         = "response"
	ConfidenceFactor = "confidence_factor"

	// Take from adaptive engagements
	ViewCollaborationReport = "view_collaboration_report"
	logger                  = alog.LambdaLogger(logrus.InfoLevel)
)

func surveyElems(dimension, dimensionDesc string, rating, comments string) []ebm.AttachmentActionTextElement {
	var confidenceValues []ebm.AttachmentActionElementOption
	for i := 5; i >= 1; i-- {
		confidenceValues = append(confidenceValues, ebm.AttachmentActionElementOption{
			Label: fmt.Sprintf("%s", coaching.Feedback360RatingMap[strconv.Itoa(i)]),
			Value: strconv.Itoa(i),
		})
	}
	actionTextElements := []ebm.AttachmentActionTextElement{
		// Placeholder text is limited to 150 characters
		// TODO: Move this to engagement-builder
		ebm.NewTextArea(Feedback, "Advice for the upcoming quarter",
			ui.PlainText(core.ClipString(dimensionDesc, 150, "...")), ui.PlainText(comments)),
		{
			Label:    fmt.Sprintf("%s assessment for the current quarter", strings.Title(dimension)),
			Name:     ConfidenceFactor,
			ElemType: models.MenuSelectType,
			Options:  confidenceValues,
			Value:    rating,
		},
	}

	// Survey box should consist of a menu option to select rating and a text area for a user to enter the feedback
	return actionTextElements
}

func DimensionSurvey(value models.AdaptiveValue, title string, rating, comments string) ebm.AttachmentActionSurvey {
	return ebm.AttachmentActionSurvey{
		Title:       title,
		SubmitLabel: models.SubmitLabel,
		Elements:    surveyElems(value.Name, value.Description, rating, comments),
	}
}

func engAttachActions(mc models.MessageCallback, details bool, feedbackGiven bool) []ebm.AttachmentAction {
	callbackId := mc.ToCallbackID()

	answerLabel := core.IfThenElse(feedbackGiven, "Edit", "I can answer now").(string)
	attachAction1, _ := eb.NewAttachmentActionBuilder().
		Name(mc.Action).
		Text(answerLabel).
		ActionType(ebm.AttachmentActionTypeButton).
		Style(ebm.AttachmentActionStylePrimary).
		Value(callbackId).
		Build()

	attachAction2, _ := eb.NewAttachmentActionBuilder().
		Name(fmt.Sprintf("%s_%s", Ask, models.Ignore)).
		Text(string(RemoveEngagementLabel)).
		ActionType(ebm.AttachmentActionTypeButton).
		Style(ebm.AttachmentActionStyleDanger).
		Value(callbackId).
		Build()

	var detailsName, detailsText string
	if details {
		detailsName = LessDetails
		detailsText = "Less Details"
	} else {
		detailsName = MoreDetails
		detailsText = "More Details"
	}

	attachAction3, _ := eb.NewAttachmentActionBuilder().
		Name(fmt.Sprintf("%s_%s", Ask, detailsName)).
		Text(detailsText).
		ActionType(ebm.AttachmentActionTypeButton).
		Style(ebm.AttachmentActionStylePrimary).
		Value(callbackId).
		Build()

	// TODO: Update this when you enable optional engagements
	return []ebm.AttachmentAction{*attachAction1, *attachAction2, *attachAction3}
}

func confirmNotificationAttachmentActions(mc models.MessageCallback, noText ui.PlainText, learnTrailPath string) []ebm.AttachmentAction {
	return models.AppendOptionalAction(
		[]ebm.AttachmentAction{
			*models.SimpleAttachAction(mc, models.Now, "Yes"),
			*models.SimpleAttachAction(mc, models.Ignore, noText),
		},
		models.LearnMoreAction(models.ConcatPrefixOpt("docs/general/", learnTrailPath)),
	)
}

func ScheduleFeedback(target string, mc models.MessageCallback, userId, attachLabel string,
	teamID models.TeamID) {
	feedbackRequestEngagement(target, *mc.WithAction(Ask).WithTarget(target).WithSource(userId),
		userId, attachLabel, "", teamID)
}

func confirmNotificationAttachment(mc models.MessageCallback, title string, noText ui.PlainText, learnTrailPath string) []ebm.Attachment {
	callbackId := mc.ToCallbackID()
	attach, _ := eb.NewAttachmentBuilder().
		Title(title).
		Fallback("Are you sure you would like to give feedback?").
		Color(models.BlueColorHex).
		AttachmentType(models.DefaultAttachmentType).
		CallbackId(callbackId).
		Actions(confirmNotificationAttachmentActions(mc, noText, learnTrailPath)).
		Build()
	return []ebm.Attachment{*attach}
}

func feedbackRequestEngagement(target string, mc models.MessageCallback, userId, attachLabel, trailPath string,
	teamID models.TeamID) {
	attachs := &confirmNotificationAttachment(mc, attachLabel, "No, thank you", trailPath)[0]
	callbackId := mc.ToCallbackID()
	engagement := eb.NewEngagementBuilder().
		Id(callbackId).
		WithResponseType(models.SlackInChannel).
		WithAttachment(attachs).
		Build()

	bytes, err := engagement.ToJson()
	core.ErrorHandler(err, namespace, "Could not convert engagement to JSON")
	// TODO: Think about managing the priority here
	eng := models.UserEngagement{UserID: userId, TargetID: target, ID: callbackId,
		PlatformID: teamID.ToPlatformID(),
		Script:     string(bytes), Priority: models.UrgentPriority, Answered: 0, CreatedAt: core.CurrentRFCTimestamp()}
	utils.AddEng(eng, engagementTable, d, namespace)
}

func feedbackEngagementAttachment(value models.AdaptiveValue,
	mc models.MessageCallback,
	details bool,
	conn daosCommon.DynamoDBConnection) *ebm.Attachment {
	var existingFeedback string
	op, err := existingFeedbackOnDimension(mc, value, conn)
	if err == nil {
		if op.Feedback != "" {
			existingFeedback = fmt.Sprintf("[%s] %s", coaching.Feedback360RatingMap[op.ConfidenceFactor], op.Feedback)
		}
	} else {
		logger.WithError(err).
			Errorf("Could not get existing feedback for %s dimension for %s user", value.Name, mc.Target)
	}
	actions := engAttachActions(mc, details, op.Feedback != "")
	baseAttachmentBuilder := eb.NewAttachmentBuilder().
		Title(fmt.Sprintf("%s's %s", common.TaggedUser(mc.Target), value.Name)).
		Fallback(value.Name).
		Color(models.BlueColorHex).
		AttachmentType(models.DefaultAttachmentType).
		CallbackId(mc.ToCallbackID()).
		Actions(actions)

	if details {
		// Attach existing response to the attachment if present
		baseAttachmentBuilder.Text(fmt.Sprintf("%s\n%s", value.Description, existingFeedback))
	}

	attach, _ := baseAttachmentBuilder.Build()
	return attach
}

func feedbackEngagement(value models.AdaptiveValue, mc models.MessageCallback, urgent bool,
	conn daosCommon.DynamoDBConnection) {
	callbackId := mc.ToCallbackID()

	// Determining priority for the engagement
	var priority = models.HighPriority
	if urgent {
		priority = models.UrgentPriority
	}

	attach := feedbackEngagementAttachment(value, mc, false, conn)

	engagement := eb.NewEngagementBuilder().
		Id(callbackId).
		WithResponseType(models.SlackInChannel).
		WithAttachment(attach).
		Build()
	bytes, err := engagement.ToJson()
	core.ErrorHandler(err, namespace, "Could not convert engagement to JSON")
	eng := models.UserEngagement{UserID: mc.Source, TargetID: mc.Target, ID: callbackId,
		PlatformID: conn.PlatformID,
		Script:     string(bytes), Priority: priority, Answered: 0, CreatedAt: core.CurrentRFCTimestamp()}
	utils.AddEng(eng, engagementTable, d, namespace)
}

type MsgState struct {
	Id       string `json:"id"`
	ThreadTs string `json:"thread_ts"`
}

func dialogFromSurvey1(api *slack.Client, message slack.InteractionCallback, survey ebm.AttachmentActionSurvey, id string) error {
	survState := func() string {
		// When the original message is from a thread, we need to post to the same thread
		// Below logic checks if the incoming message is from a thread
		var ts string
		if message.OriginalMessage.ThreadTimestamp == "" {
			ts = message.MessageTs
		} else {
			ts = message.OriginalMessage.ThreadTimestamp
		}
		msgStateBytes, err := json.Marshal(MsgState{Id: id, ThreadTs: ts})
		core.ErrorHandler(err, namespace, "Could not marshal MsgState")
		return string(msgStateBytes)
	}

	return utils.SlackSurvey(api, message, survey, id, survState)
}

func EmptyAttachs() []ebm.Attachment {
	return []ebm.Attachment{}
}

// Triggering on user attributes table
func HandleRequest(ctx context.Context, np models.NamespacePayload4) {
	logger = logger.WithLambdaContext(ctx)
	defer core.RecoverAsLogError("feedback-setup-lambda")

	// if request.Payload == "warmup" {
	// 	return nil
	// }
	// Parsing incoming payload
	slackRequest := np.PlatformRequest.SlackRequest
	teamID := np.TeamID
	if !teamID.IsEmpty() {
		switch slackRequest.Type {
		case models.InteractionSlackRequestType:
			dispatchSlackInteractionCallback(slackRequest.InteractionCallback, teamID)
		case models.DialogSubmissionSlackRequestType:
			request := slackRequest.InteractionCallback
			dialog := slackRequest.DialogSubmissionCallback
			logger.Infof("Got dialog submission " + dialog.State)
			dispatchSlackDialogSubmissionCallback(request, dialog, teamID)
		}
	} else {
		logger.Errorf("Platform id is empty for %s user", slackRequest.InteractionCallback.User.ID)
	}

	return
}

const (
	selfCoachingMessage = "*_I think you are awesome too, but you can’t coach yourself.  Please choose someone else to coach to be awesome._*"
)

func overrideOriginalMessage(request slack.InteractionCallback, message string) models.PlatformSimpleNotification {
	// return utils.InteractionCallbackOverrideOriginalMessage(request, message)
	return models.PlatformSimpleNotification{
		UserId:  request.User.ID,
		Channel: request.Channel.ID,
		Message: message,
		Ts:      request.OriginalMessage.Timestamp,
	}
}

func dispatchSlackInteractionCallback(request slack.InteractionCallback, teamID models.TeamID) {
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	// Parse callback Id to messageCallback
	mc := utils.MessageCallbackParseUnsafe(request.CallbackID, namespace)
	log.Println("Callback ID: " + mc.Sprint())
	var notes []models.PlatformSimpleNotification
	// MessageCallback is formed like this: Module:Source:Topic:Action:Target:Month:Year
	if mc.Module == "coaching" {
		if mc.Topic == "user_feedback" {
			values := evalues.ReadAndSortAllAdaptiveValues(conn)
			action := request.ActionCallback.AttachmentActions[0]
			if strings.HasPrefix(action.Name, Request) {
				// Parse callback Id to messageCallback
				act := strings.TrimPrefix(action.Name, fmt.Sprintf("%s_", Request))
				switch act {
				case string(models.Now):
					selectedUser := action.SelectedOptions[0].Value
					notes = scheduleFeedbackForUserHandler(request, mc, selectedUser, teamID)
				case string(models.Ignore), string(models.Cancel):
					notes = cancelEngagementHandler(request, mc)
				}
			} else if strings.HasPrefix(action.Name, Ask) {
				// Parse callback Id to messageCallback
				mc := utils.MessageCallbackParseUnsafe(action.Value, namespace)
				act := strings.TrimPrefix(action.Name, fmt.Sprintf("%s_", Ask))
				switch act {
				case string(models.Now):
					notes = feedbackRequestNowHandler(request, mc, values, conn)
				case string(models.Ignore):
					notes = cancelEngagementHandler(request, mc)
				case MoreDetails, LessDetails:
					notes = onShowDetailsToggle(teamID, mc, request, act)(conn)
				default:
					values, err2 := adaptiveValue.ReadOrEmpty(act)(conn)
					values = adaptiveValue.AdaptiveValueFilterActive(values)
					if err2 == nil && len(values) > 0 {
						value := values[0]
						logger.Infof("Retrieved value with id %s: %v", act, value)
						// this corresponds to the engagements for each of the dimensions
						notes = feedbackDimensionHandler(request, mc, action.Value, value, conn)
					} else if err2 != nil {
						logger.Errorf("Could not read value with id %s: %v", act, err2)
					} else {
						logger.Errorf("Could not find value with id %s", act)
					}
				}
			} else if strings.HasPrefix(action.Name, "confirm") {
				act := strings.TrimPrefix(action.Name, fmt.Sprintf("%s_", "confirm"))
				logger.WithField("action", &action).WithField("act", &act).Info()
				switch act {
				case string(models.Now):
					feedbackShowSelectUserHandler(request, conn)
					notes = []models.PlatformSimpleNotification{
						{UserId: request.User.ID,
							Channel: request.Channel.ID, Message: "", Ts: request.OriginalMessage.Timestamp}}
				case string(models.Ignore):
					notes = cancelEngagementHandler(request, mc)
				}
			} else if mc.Action == "select" {
				if action.Name == "cancel" {
					notes = cancelEngagementHandler(request, mc)
				} else if action.Name == fmt.Sprintf("%s_now", mc.Action) {
					selectedUser := action.SelectedOptions[0].Value
					if selectedUser == request.User.ID {
						// Checking for self-feedback
						notes = []models.PlatformSimpleNotification{
							overrideOriginalMessage(request, selfCoachingMessage)}
					} else {
						notes = feedbackNowEngagementHandler(request, mc, selectedUser, values, conn)
					}
				} else if action.Name == fmt.Sprintf("%s_later", mc.Action) {
					notes = postponeEngagementHandler(request, mc, values, conn)
				} else if action.Name == fmt.Sprintf("%s_ignore", mc.Action) ||
					action.Name == fmt.Sprintf("%s_cancel", mc.Action) {
					notes = cancelEngagementHandler(request, mc)
				}
			}
		}
	} else if mc.Module == "feedback" && mc.Topic == "report" {
		action := request.ActionCallback.AttachmentActions[0]
		if strings.HasPrefix(action.Name, ViewCollaborationReport) {
			act := strings.TrimPrefix(action.Name, fmt.Sprintf("%s_", ViewCollaborationReport))
			notes = viewCollaborationReportHandler(request, teamID, act, mc)
		}
	}
	platform.PublishAll(notes)
}

func quarterYear(date business_time.Date) string {
	year := date.GetYear()
	quarter := date.GetQuarter()
	return fmt.Sprintf("%d:%d", quarter, year)
}

func feedbackShowSelectUserHandler(request slack.InteractionCallback, conn daosCommon.DynamoDBConnection) {
	year, month, _ := time.Now().Date()
	mc := models.MessageCallback{Module: "coaching", Source: request.User.ID,
		Topic: "user_feedback", Action: "select", Target: "", Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	UserSelectEngagement(request.User.ID, mc, []string{}, []string{request.User.ID},
		"Whom would you like to give feedback to?", "coaching-feedback", conn)
}

func UserSelectEngagement(userID string, mc models.MessageCallback, users,
	filter []string, text, context string, conn daosCommon.DynamoDBConnection) {
	user.UserSelectEng(userID, engagementTable, conn, mc,
		users, filter, text, context, models.UserEngagementCheckWithValue{})
}

func scheduleFeedbackForUserHandler(request slack.InteractionCallback, mc models.MessageCallback,
	selectedUser string, teamID models.TeamID) []models.PlatformSimpleNotification {
	// Update feedback request engagement as answered
	utils.UpdateEngAsAnswered(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
	date := business_time.Today(time.UTC)
	quarterYear := quarterYear(date)

	if !feedbackDAO.IsThereFeedbackFromTo(selectedUser, request.User.ID, quarterYear) {
		// Source and target would be reversed for the feedback notification
		ScheduleFeedback(request.User.ID, mc, selectedUser,
			FeedbackRequestedAskIfYouWantToProvideTemplate(request.User.ID), teamID)
	}
	// Send a notification to the current user and delete the original message
	// We do the same action regardless of the fact that user has already provided feedback
	// just to not disclose this information.
	override := overrideOriginalMessage(request, "")
	return []models.PlatformSimpleNotification{override,
		utils.InteractionCallbackSimpleResponse(request,
			string(ConfirmFeedbackRequestedTemplate(selectedUser)))}
}

// sends one request per adaptive value.
func feedbackRequestNowHandler(
	request slack.InteractionCallback,
	mc models.MessageCallback,
	values []models.AdaptiveValue,
	conn daosCommon.DynamoDBConnection,
	) []models.PlatformSimpleNotification {
	target := mc.Target
	utils.UpdateEngAsAnswered(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
	// for the user, write feedback engagements with non-urgent priority
	for _, value := range values {
		// We add prefix 'ask_' to each of the dimension engagement for a user
		mc.WithAction(fmt.Sprintf("ask_%s", value.ID)).WithTarget(target)
		feedbackEngagement(value, mc, true, conn) // TODO: return PlatformSimpleNotification-s
	}
	return []models.PlatformSimpleNotification{overrideOriginalMessage(request, "")}
}

func existingFeedbackOnDimension(mc models.MessageCallback, 
	value models.AdaptiveValue,
	conn daosCommon.DynamoDBConnection) (op models.UserFeedback, err error) {
	key := mc.WithAction(value.ID).ToCallbackID()
	var feedbacks []userFeedback.UserFeedback
	feedbacks, err = userFeedback.ReadOrEmpty(key)(conn)
	if len(feedbacks) > 0 {
		op = feedbacks[0]
	}
	return
}

func feedbackDimensionHandler(
	request slack.InteractionCallback,
	mc models.MessageCallback,
	actionValue string,
	value models.AdaptiveValue,
	conn daosCommon.DynamoDBConnection) []models.PlatformSimpleNotification {

	ut := userTokenSyncUnsafe(request.User.ID)
	tut := daosUser.ReadUnsafe(mc.Target)(conn)
	api := slack.New(ut)
	// key := mc.WithAction(value.ID).Sprint()
	// // Query the feedback table. If this has already been answered, get the confidence factor and script associated with the id
	// params := map[string]*dynamodb.AttributeValue{
	// 	"id": {
	// 		S: aws.String(key),
	// 	},
	// }
	// var op models.UserFeedback
	op, err := existingFeedbackOnDimension(mc, value, conn)
	// err := d.QueryTable(feedbackTable, params, &op)
	// core.ErrorHandler(err, namespace, fmt.Sprintf("Could not query %s table for default values", feedbackTable))
	// Open a survey associated with the engagement
	err = dialogFromSurvey1(api, request, DimensionSurvey(value, tut.DisplayName, op.ConfidenceFactor, op.Feedback), actionValue)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", actionValue))

	return []models.PlatformSimpleNotification{}
}

func feedbackNowEngagementHandler(
	request slack.InteractionCallback,
	mc models.MessageCallback,
	selectedUser string,
	values []models.AdaptiveValue,
	conn daosCommon.DynamoDBConnection) []models.PlatformSimpleNotification {
	// Add engagements with urgent priority
	// We have now added feedback for a coaching engagement. We can now update the original engagement as answered.
	utils.UpdateEngAsAnswered(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)

	// for each user, write feedback engagements with non-urgent priority
	for _, value := range values {
		// We add prefix 'ask_' to each of the dimension engagement for a user
		mc.WithAction(fmt.Sprintf("ask_%s", value.ID)).WithTarget(selectedUser)
		feedbackEngagement(value, mc, true, conn)
	}
	// Delete original engagement
	return []models.PlatformSimpleNotification{overrideOriginalMessage(request, "")}
}

func postponeEngagementHandler(
	request slack.InteractionCallback,
	mc models.MessageCallback,
	values []models.AdaptiveValue,
	conn daosCommon.DynamoDBConnection) []models.PlatformSimpleNotification {
	target := mc.Target
	// Add engagements with non-urgent priority
	// We have now added feedback for a coaching engagement. We can now update the original engagement as answered.
	utils.UpdateEngAsAnswered(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)

	// for each user, write feedback engagements with non-urgent priority
	for _, value := range values {
		// We add prefix 'ask_' to each of the dimension engagement for a user
		mc.Set("Action", fmt.Sprintf("ask_%s", value.ID))
		mc.Set("Target", target)
		feedbackEngagement(value, mc, false, conn)
	}
	return []models.PlatformSimpleNotification{models.PlatformSimpleNotification{UserId: request.User.ID, Channel: request.Channel.ID, Message: fmt.Sprintf(
		"_Ok, you can provide feedback to %s during the next window._", common.TaggedUser(target)), Ts: request.OriginalMessage.Timestamp, Attachments: EmptyAttachs()}}
}

func cancelEngagementHandler(request slack.InteractionCallback, mc models.MessageCallback) []models.PlatformSimpleNotification {
	// Update engagement as answered and don't do anything more
	utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
	// Delete original engagement
	return []models.PlatformSimpleNotification{models.PlatformSimpleNotification{UserId: request.User.ID,
		Channel: request.Channel.ID, Message: "", Ts: request.OriginalMessage.Timestamp}}
}

func viewCollaborationReportHandler(request slack.InteractionCallback,
	teamID models.TeamID, act string, mc models.MessageCallback,
) []models.PlatformSimpleNotification {
	switch act {
	case string(models.Now):
		// These are for the simulated date
		m, err := strconv.Atoi(mc.Month)
		y, err := strconv.Atoi(mc.Year)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse string to int"))
		bt := business_time.NewDate(y, m, 1)
		// view report now
		// engageBytes, _ := json.Marshal(models.UserEngage{
		// 	UserID: request.User.ID,
		// 	IsNew: false,
		// 	Update: true,
		// 	Channel: request.Channel.ID, ThreadTs: request.MessageTs,
		// 	Date: bt.DateToString(string(core.ISODateLayout))})
		// Update original message
		platform.Publish(models.PlatformSimpleNotification{UserId: request.User.ID, Channel: request.Channel.ID,
			Message: fmt.Sprintf("_Hang tight, fetching your collaboration report for quarter `%d`, year `%d` :point_down:_",
				bt.GetPreviousQuarter(), bt.GetPreviousQuarterYear()), Ts: request.MessageTs})
		// This is used to add an engagement on who to give feedback to
		err = feedbackReportPostingLambda.DeliverReportToUserAsync(teamID, request.User.ID, bt.DateToTimeMidnight())
		// _, err = l.InvokeFunction(collaborationReportPostingLambda, engageBytes, false)
		// core.ErrorHandler(err, namespace, fmt.Sprintf("Could not invoke %s from feedback-setup-lambda",
		// 	collaborationReportPostingLambda))
		utils.UpdateEngAsAnswered(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
	case string(models.Ignore):
		utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
		utils.InteractionCallbackSimpleResponse(request, "")
		DeleteOriginalEng(request.User.ID, request.Channel.ID, request.MessageTs)
	}
	return []models.PlatformSimpleNotification{}
}

func dispatchSlackDialogSubmissionCallback(
	request slack.InteractionCallback,
	dialog slack.DialogSubmissionCallback,
	teamID models.TeamID) {
	mc := utils.MessageCallbackParseUnsafe(request.CallbackID, namespace)

	form := dialog.Submission
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	var notes []models.PlatformSimpleNotification
	if strings.HasPrefix(mc.Action, "ask") {
		notes = askDialogSubmissionHandler(request, dialog, form, teamID)(conn)
	} else {
		notes = noOpDialogSubmissionHandler(request, dialog, form, teamID)
	}
	platform.PublishAll(notes)
}

func createFeedbackMessage(request slack.InteractionCallback,
	editAction string,
	targetID string,
	form map[string]string,
	timestamp string,
) func(conn daosCommon.DynamoDBConnection) models.PlatformSimpleNotification {
	return func(conn daosCommon.DynamoDBConnection) models.PlatformSimpleNotification {
		competencyID := strings.TrimPrefix(editAction, "ask_")
		values, err2 := adaptiveValue.ReadOrEmpty(competencyID)(conn)
		values = adaptiveValue.AdaptiveValueFilterActive(values)
		found := len(values) > 0
		var attachNotification models.PlatformSimpleNotification
		if err2 == nil && found {
			value := values[0]
			logger.WithField("value", value).Infof("Retrieved value with id=%s", competencyID)
			confFactor := form[ConfidenceFactor]
			response := form[Feedback]

			attachAction, _ := eb.NewAttachmentActionBuilder().
				Name(editAction).
				Text(models.EditLabel).
				ActionType(models.ButtonType).
				Value(request.CallbackID).
				Build()

			attach, _ := eb.NewAttachmentBuilder().
				CallbackId(request.CallbackID).
				Author(ebm.AttachmentAuthor{Name: fmt.Sprintf("<@%s>'s %s", targetID, value.Name)}).
				Color(models.BlueColorHex).
				Actions([]ebm.AttachmentAction{*attachAction}).
				Fields([]ebm.AttachmentField{
					{
						Title: "Feedback",
						Value: response,
					},
					{
						Title: "Confidence Factor",
						Value: fmt.Sprintf("%s", coaching.Feedback360RatingMap[confFactor]),
						Short: false,
					},
				}).
				Build()
			attachNotification = utils.InteractionCallbackSimpleResponse(request, "")
			attachNotification.Ts = timestamp
			attachNotification.Attachments = []ebm.Attachment{*attach}
		} else if err2 != nil {
			logger.Errorf("Could not retrieve value for id %s: %v", competencyID, err2)
			attachNotification = utils.InteractionCallbackSimpleResponse(request, "Apologies, something has gone wrong")
		} else {
			logger.Errorf("Could not find value for id %s", competencyID)
			attachNotification = utils.InteractionCallbackSimpleResponse(request, "Apologies, something has gone wrong")
		}

		return attachNotification
	}
}

func askDialogSubmissionHandler(
	request slack.InteractionCallback,
	dialog slack.DialogSubmissionCallback,
	form map[string]string,
	teamID models.TeamID) func(conn daosCommon.DynamoDBConnection) []models.PlatformSimpleNotification {
	return func(conn daosCommon.DynamoDBConnection) (notes []models.PlatformSimpleNotification) {
		mc := utils.MessageCallbackParseUnsafe(request.CallbackID, namespace)

		var msgState MsgState
		err := json.Unmarshal([]byte(dialog.State), &msgState)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse msgState %s", dialog.State))

		// Update the original attachment with the received feedback
		updateTheOriginalAttachmentWithTheReceivedFeedback := createFeedbackMessage(request, mc.Action, mc.Target, form, msgState.ThreadTs)(conn)

		// we have now added feedback for a coaching engagement. We can now update the original engagement as answered
		utils.UpdateEngAsAnswered(mc.Source, request.CallbackID, engagementTable, d, namespace)

		// Collecting responses from dialog submission
		confFactor := form[ConfidenceFactor]
		response := form[Feedback]
		value := strings.TrimPrefix(mc.Action, "ask_")

		mc.Set("Action", value)
		// Storing feedback
		feedback := models.UserFeedback{
			ID:               mc.ToCallbackID(),
			Source:           mc.Source,
			Target:           mc.Target,
			ValueID:          value,
			ConfidenceFactor: confFactor,
			Feedback:         response,
			QuarterYear:      fmt.Sprintf("%d:%s", core.MonthStrToQuarter(mc.Month), mc.Year),
			ChannelID:        request.Channel.ID,
			MsgTimestamp:     msgState.ThreadTs,
			PlatformID:       teamID.ToPlatformID(),
		}
		err = d.PutTableEntry(feedback, feedbackTable)
		if err == nil {
			core.ErrorHandler(err, namespace, fmt.Sprintf("Could not write to %s table", feedbackTable))
			byt, _ := json.Marshal(feedback)
			_, err = l.InvokeFunction(feedbackAnalysisLambda, byt, true)
			notes = []models.PlatformSimpleNotification{
				updateTheOriginalAttachmentWithTheReceivedFeedback,
			}
		} else {
			logger.WithField("error", err).Errorf("Could not write to %s table", feedbackTable)
		}
		return
	}
}

func onShowDetailsToggle(teamID models.TeamID, mc models.MessageCallback, request slack.InteractionCallback, act string) func(conn daosCommon.DynamoDBConnection) []models.PlatformSimpleNotification {
	return func(conn daosCommon.DynamoDBConnection) []models.PlatformSimpleNotification {
		details := act == MoreDetails
		valueID := strings.TrimPrefix(mc.Action, fmt.Sprintf("%s_", Ask))
		value := adaptiveValue.ReadUnsafe(valueID)(conn)
		attach := feedbackEngagementAttachment(value, mc, details, conn)
		return []models.PlatformSimpleNotification{
			{UserId: request.User.ID,
				Channel:     request.Channel.ID,
				Ts:          request.OriginalMessage.Timestamp,
				Attachments: []ebm.Attachment{*attach}},
		}
	}
}

func noOpDialogSubmissionHandler(
	request slack.InteractionCallback,
	dialog slack.DialogSubmissionCallback,
	form map[string]string,
	teamID models.TeamID) []models.PlatformSimpleNotification {
	return []models.PlatformSimpleNotification{}
}

func DeleteOriginalEng(userId, channel, ts string) {
	utils.DeleteOriginalEng(userId, channel, ts, func(notification models.PlatformSimpleNotification) {
		platform.Publish(notification)
	})
}

func main() {
	ls.Start(HandleRequest)
}
