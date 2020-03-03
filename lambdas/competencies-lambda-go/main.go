package competencies

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/pkg/errors"

	"time"

	evalues "github.com/adaptiveteam/adaptive/adaptive-engagements/values"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	plat "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
)

const (
	AdaptiveValuesNamespace          = "adaptive_values"
	AdaptiveValuesListMenuItem       = AdaptiveValuesNamespace + ":list"
	AdaptiveValuesSimpleListMenuItem = AdaptiveValuesNamespace + ":simple_list"
	AdaptiveValuesCreateNewMenuItem  = AdaptiveValuesNamespace + ":create_new"

	AdaptiveValuesDialogContext = "dialog/competency/language-coaching"
)

var (
	newAdaptiveValueTemplate = models.AdaptiveValue{
		ID:          core.Uuid(),
		Name:        "",
		Description: "",
		ValueType:   "",
	}
)

type SlackConversationKind string

const (
	InThread SlackConversationKind = "in-thread"
	InChat   SlackConversationKind = "in-chat"
)
const (
	ModifyListAdaptiveValueAction = "modify-list-adaptive-value"
)

func EditAdaptiveValueAction(slackConversationKind SlackConversationKind) string {
	return "edit-adaptive-value-" + string(slackConversationKind)
}

func SubmitNewAdaptiveValueAction(slackConversationKind SlackConversationKind) string {
	return "submit-new-adaptive-value-" + string(slackConversationKind)
}

func SubmitUpdatedAdaptiveValueAction(slackConversationKind SlackConversationKind) string {
	return "submit-updated-adaptive-value-" + string(slackConversationKind)
}

func DeleteAdaptiveValueAction(slackConversationKind SlackConversationKind) string {
	return "delete-adaptive-value-" + string(slackConversationKind)
}

func CreateAdaptiveValueAction(slackConversationKind SlackConversationKind) string {
	return "create-adaptive-value-" + string(slackConversationKind)
}

// func ModifyAdaptiveValueAction(slackConversationKind SlackConversationKind) string {
// 	return "modify-adaptive-value-" + string(slackConversationKind)
// }

// func main() {
// 	lambda.Start(HandleRequest)
// }

// HandleRequest receives lambda json event
func HandleRequest(ctx context.Context, e events.SNSEvent) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("Error in values-lambda %v", err2)
		}
	}()
	fmt.Println("adaptiveValues/main.go/HandleRequest entered")
	for _, record := range e.Records {
		fmt.Println(record)
		if record.SNS.Message == "warmup" {
			// Ignoring warmup messages
		} else {
			np := models.UnmarshalNamespacePayload4JSONUnsafe(record.SNS.Message)
			if np.Namespace == AdaptiveValuesNamespace {
				err = HandleNamespacePayload4(np)
			}
		}
	}
	return // we do not have handlable errors. Only panics
}

// HandleNamespacePayload4 - handle all logic
func HandleNamespacePayload4(np models.NamespacePayload4) (err error) {
	defer core.RecoverToErrorVar("Competencies", &err)
	teamID := np.TeamID
	switch np.SlackRequest.Type {
	case models.InteractionSlackRequestType:
		dispatchSlackInteractionCallback(np.SlackRequest.InteractionCallback, teamID)
	case models.DialogSubmissionSlackRequestType:
		dispatchSlackDialogSubmissionCallback(np.SlackRequest.InteractionCallback, np.SlackRequest.DialogSubmissionCallback, teamID)
	default:
		err = errors.Errorf("Unknown request of type %s: %v\n", np.SlackRequest.Type, np)
	}
	return
}

// noMessageOverrideTs is a predefined constant that allows to skip original message overriding
const noMessageOverrideTs = ""

func dispatchSlackInteractionCallback(request slack.InteractionCallback, teamID models.TeamID) {
	// defer deleteOriginalEng(request.User.ID, request.Channel.ID, request.MessageTs)
	// defer platform.RecoverGracefully(request)

	action := request.ActionCallback.AttachmentActions[0]
	platform.Debug(request, "Action: "+action.Name)
	notes := responses()
	// 'menu_list' is for the options that are presented to the userID
	switch action.Name {
	case "menu_list":
		selectedOption := action.SelectedOptions[0].Value
		switch selectedOption {
		case AdaptiveValuesListMenuItem:
			notes = detailedListMenuItemFunc(request, teamID)
			// notes = append(notes, clearRequestMessage(request)...) // do not clear original message. It'll be replaced with invite to thread.
		case AdaptiveValuesSimpleListMenuItem:
			notes = listMenuItemFunc(request, teamID)
			notes = append(notes, clearRequestMessage(request)...) // clear original message
		case AdaptiveValuesCreateNewMenuItem:
			notes = createNewAdaptiveValueMenuItemAndClearOriginalMessageAfterwards(request, teamID)
			//			notes = append(notes, clearRequestMessage(request)) // we don't clear original message immediately. Only on dialog submission
		default:
			panic(errors.New("Unknown competencies command " + selectedOption))
		}
	// case ModifyAdaptiveValueAction:
	case DeleteAdaptiveValueAction(InThread):
		notes = onDeleteButtonClicked(request, action.Value)
	case DeleteAdaptiveValueAction(InChat):
		notes = onDeleteButtonClicked(request, action.Value)
	case EditAdaptiveValueAction(InThread):
		notes = onEditButtonClicked(request, teamID, action.Value, InThread)
	case EditAdaptiveValueAction(InChat):
		notes = onEditButtonClicked(request, teamID, action.Value, InChat)
	case CreateAdaptiveValueAction(InThread):
		notes = onCreateAction(request, teamID, noMessageOverrideTs, InThread)
	case CreateAdaptiveValueAction(InChat):
		notes = onCreateAction(request, teamID, noMessageOverrideTs, InChat)
	case ModifyListAdaptiveValueAction:
		notes = detailedListMenuItemFunc(request, teamID)
	default:
		platform.Debug(request, "unknown action "+action.Name)
	}
	platform.PublishAll(notes)
}

func dispatchSlackDialogSubmissionCallback(request slack.InteractionCallback,
	dialog slack.DialogSubmissionCallback, teamID models.TeamID) {
	platform.Debug(request, "Got filled form")
	mc := utils.MessageCallbackParseUnsafe(request.CallbackID, AdaptiveValuesNamespace)
	notes := responses()
	switch mc.Action {
	case SubmitNewAdaptiveValueAction(InThread):
		notes = onSubmitNewButtonClicked(request, dialog, teamID, InThread)
	case SubmitNewAdaptiveValueAction(InChat):
		notes = onSubmitNewButtonClicked(request, dialog, teamID, InChat)
	case SubmitUpdatedAdaptiveValueAction(InThread):
		notes = onSubmitUpdatedButtonClicked(request, dialog, teamID, InThread)
	case SubmitUpdatedAdaptiveValueAction(InChat):
		notes = onSubmitUpdatedButtonClicked(request, dialog, teamID, InChat)
	default:
		platform.Debug(request, "Couldn't handle dialog submission "+request.CallbackID)
	}
	platform.PublishAll(notes)
}

// listMenuItemFunc renders the list of adaptiveValues directly in chat.
func listMenuItemFunc(request slack.InteractionCallback, teamID models.TeamID) []models.PlatformSimpleNotification {
	platform.Debug(request, "Listing adaptive values 2 ...")
	adaptiveValues := adaptiveValuesTableDao.ForPlatformID(teamID.ToPlatformID()).AllUnsafe()
	sort.Slice(adaptiveValues, func(i, j int) bool {
		return adaptiveValues[i].Name < adaptiveValues[j].Name
	})
	adaptiveValueItems := mapAdaptiveValueString(adaptiveValues, RenderAdaptiveValueItem)
	listOfItems := strings.Join(adaptiveValueItems, "\n")
	mc := callbackID(request.User.ID, ModifyListAdaptiveValueAction)

	attachment := eb.NewAttachmentBuilder().
		Title(AdaptiveValuesTemplate(AdaptiveValuesListTitleSubject)).
		Text(listOfItems).
		Actions([]ebm.AttachmentAction{
			modifyListAdaptiveValueAction(request),
			addAnotherAdaptiveValueAction(request, InChat),
		}).
		CallbackId(mc.ToCallbackID()).
		ToAttachment()
	p := utils.InteractionCallbackSimpleResponse(request, "") // RenderAdaptiveValueItem(h)
	p.Attachments = []ebm.Attachment{attachment}

	return responses(p)
}

func mapAdaptiveValueString(vs []models.AdaptiveValue, f func(models.AdaptiveValue) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func mapAdaptiveValuePlatformSimpleNotification(vs []models.AdaptiveValue,
	f func(models.AdaptiveValue) models.PlatformSimpleNotification) []models.PlatformSimpleNotification {
	vsm := make([]models.PlatformSimpleNotification, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

var (
	deleteAdaptiveValueAction = func(request slack.InteractionCallback, h models.AdaptiveValue, slackConversationKind SlackConversationKind) ebm.AttachmentAction {
		return eb.NewButtonDanger(DeleteAdaptiveValueAction(slackConversationKind), h.ID,
			ui.PlainText(AdaptiveValuesTemplate(DeleteButtonSubject)),
			ui.PlainText(AdaptiveValuesTemplate(CancelDeleteSubject)),
		)
	}

	editAdaptiveValueAction = func(request slack.InteractionCallback, h models.AdaptiveValue, slackConversationKind SlackConversationKind) ebm.AttachmentAction {
		return eb.NewButton(EditAdaptiveValueAction(slackConversationKind), h.ID,
			ui.PlainText(AdaptiveValuesTemplate(EditButtonSubject)))
	}

	addAnotherAdaptiveValueAction = func(request slack.InteractionCallback, slackConversationKind SlackConversationKind) ebm.AttachmentAction {
		return eb.NewButton(CreateAdaptiveValueAction(slackConversationKind), "",
			ui.PlainText(AdaptiveValuesTemplate(CreateButtonSubject)))
	}

	modifyListAdaptiveValueAction = func(request slack.InteractionCallback) ebm.AttachmentAction {
		return eb.NewButton(ModifyListAdaptiveValueAction, "",
			ui.PlainText(AdaptiveValuesTemplate(ModifyListButtonSubject)))
	}
)

func clearRequestMessage(request slack.InteractionCallback) []models.PlatformSimpleNotification {
	return responses(utils.InteractionCallbackOverrideRequestMessage(request, ""))
}

func clearMessageAtTs(request slack.InteractionCallback, ts string) (clearMessage models.PlatformSimpleNotification) {
	clearMessage = utils.InteractionCallbackSimpleResponse(request, "")
	clearMessage.Ts = ts
	return
}

// adaptiveValueInlineViewAttachment creates an inline representation of the adaptiveValue with edit/delete buttons
func adaptiveValueInlineViewAttachment(request slack.InteractionCallback, h models.AdaptiveValue, slackConversationKind SlackConversationKind) ebm.Attachment {
	mc := callbackID(request.User.ID, EditAdaptiveValueAction(slackConversationKind))

	actions := []ebm.AttachmentAction{}
	if DoesUserHaveWriteAccessToCompetencies(request.User.ID) {
		actions = []ebm.AttachmentAction{
			editAdaptiveValueAction(request, h, slackConversationKind),
			deleteAdaptiveValueAction(request, h, slackConversationKind),
			addAnotherAdaptiveValueAction(request, slackConversationKind),
		}
	}
	attachment := eb.NewAttachmentBuilder().
		Fields(ShowAdaptiveValueItemAsFields(h)).
		Actions(actions).
		CallbackId(mc.ToCallbackID()).
		ToAttachment()
	return attachment
}

// adaptiveValueInlineView creates an inline representation of the adaptiveValue with edit/delete buttons
func adaptiveValueInlineView(request slack.InteractionCallback, h models.AdaptiveValue, slackConversationKind SlackConversationKind) models.PlatformSimpleNotification {
	p := utils.InteractionCallbackSimpleResponse(request, "") // RenderAdaptiveValueItem(h)
	p.Attachments = []ebm.Attachment{adaptiveValueInlineViewAttachment(request, h, slackConversationKind)}
	return p
}

func adaptiveValueEditChatMessage(request slack.InteractionCallback, slackConversationKind SlackConversationKind) func(models.AdaptiveValue) models.PlatformSimpleNotification {
	return func(h models.AdaptiveValue) models.PlatformSimpleNotification {
		p := adaptiveValueInlineView(request, h, slackConversationKind)
		p.ThreadTs = utils.TimeStamp(request)
		return p
	}
}

// detailedListMenuItemFunc sends the list of adaptiveValues to thread
func detailedListMenuItemFunc(request slack.InteractionCallback, teamID models.TeamID) []models.PlatformSimpleNotification {
	platform.Debug(request, "Listing competencies...")
	adaptiveValues := adaptiveValuesTableDao.ForPlatformID(teamID.ToPlatformID()).
		AllUnsafe()
	sort.Slice(adaptiveValues, func(i, j int) bool {
		return adaptiveValues[i].Name < adaptiveValues[j].Name
	})
	adaptiveValueItems := mapAdaptiveValuePlatformSimpleNotification(adaptiveValues,
		adaptiveValueEditChatMessage(request, InThread))
	replacementMessage := utils.InteractionCallbackOverrideRequestMessage(request,
		AdaptiveValuesTemplate(InviteToListOfAdaptiveValuesSubject))
	ts := time.Now()
	for _, item := range adaptiveValueItems {
		ts = ts.Add(time.Millisecond)
		item.Ts = core.TimestampLayout.Format(ts)
	}

	threadMessages := append(responses(
		replacementMessage,
		//	title,
	),
		adaptiveValueItems...)
	return threadMessages
}

func convertFormToAdaptiveValue(request slack.InteractionCallback, teamID models.TeamID, form map[string]string) models.AdaptiveValue {
	return models.AdaptiveValue{
		ID:          core.Uuid(),
		Name:        form["Name"],
		Description: form["Description"],
		ValueType:   form["ValueType"],
		PlatformID:  teamID.ToPlatformID(),
	}
}

func createNewAdaptiveValueMenuItemAndClearOriginalMessageAfterwards(request slack.InteractionCallback, teamID models.TeamID) []models.PlatformSimpleNotification {
	tsToOverride := request.MessageTs // we can reuse this message id to override menu if needed
	return onCreateAction(request, teamID, tsToOverride, InChat)
}

func onCreateAction(request slack.InteractionCallback, teamID models.TeamID, ts string, slackConversationKind SlackConversationKind) (resp []models.PlatformSimpleNotification) {
	resp, ok := ensureUserHasWriteAccessToValues(request)
	if ok {
		mc := callbackID(request.User.ID, SubmitNewAdaptiveValueAction(slackConversationKind))
		survey := utils.AttachmentSurvey(valueSurveyLabel, evalues.EditAdaptiveValueForm(&newAdaptiveValueTemplate))
		survey2 := ebm.AttachmentActionSurvey2{
			TriggerID: request.TriggerID,
			CallbackID: mc.ToCallbackID(),
			AttachmentActionSurvey: survey,
			State: ts,
		}
		// dialog, err := utils.ConvertSurveyToSlackDialog(survey, request.TriggerID, mc.ToCallbackID(), ts)
		// if err == nil {
			// Open a survey associated with the engagement
			conn:= daosCommon.DynamoDBConnection{
				Dynamo: d,
				ClientID: clientID,
				PlatformID: teamID.ToPlatformID(),
			}
			err2 := mapper.SlackAdapterForTeamID(conn).ShowDialog(survey2) //request.TriggerID, dialog)
			platform.Debug(request, "onCreateAction: OpenDialog - done")
			platform.ErrorHandler(request,
				fmt.Sprintf("Could not open dialog from %s survey", request.CallbackID),
				err2,
			)
			resp = responses()
		// }
	}
	return
}

func conversationID(request slack.InteractionCallback) (conversationID plat.ConversationID) {
	if request.Channel.ID == "" {
		conversationID = plat.ConversationID(request.User.ID)
	} else {
		conversationID = plat.ConversationID(request.Channel.ID)
	}
	return
}

func conversationContext(request slack.InteractionCallback, msgID mapper.MessageID) utils.ConversationContext {
	platform.Debug(request, "Channel.ID = "+request.Channel.ID)
	platform.Debug(request, "User.ID = "+request.User.ID)

	ctx := utils.ConversationContext{
		UserID:            request.User.ID,
		ConversationID:    string(conversationID(request)),
		OriginalMessageTs: msgID.Ts,
		ThreadTs:          msgID.Ts,
	}
	platform.Debug(request, fmt.Sprintf("ctx = %v", ctx))
	return ctx
}

func onSubmitNewButtonClicked(
	request slack.InteractionCallback,
	dialog slack.DialogSubmissionCallback,
	teamID models.TeamID,
	slackConversationKind SlackConversationKind,
) (resp []models.PlatformSimpleNotification) {
	resp, ok := ensureUserHasWriteAccessToValues(request)
	if ok {
		platform.Debug(request, "Creating adaptiveValue "+dialog.Submission["Name"])
		adaptiveValue := convertFormToAdaptiveValue(request, teamID, dialog.Submission)
		adaptiveValuesTableDao.ForPlatformID(teamID.ToPlatformID()).Create(adaptiveValue)
		platform.Debug(request, "Created adaptiveValue "+adaptiveValue.Name)

		convID := conversationID(request)
		platform.Debug(request, "teamID = "+teamID.ToString())
		slackAdapter := slackAPI(teamID)

		clearMessageOptionalTs := dialog.State
		if clearMessageOptionalTs != "" {
			slackAdapter.PostAsync(plat.Delete(plat.MessageID(convID, clearMessageOptionalTs)))
		}

		valueViewAttachment := adaptiveValueInlineViewAttachment(request, adaptiveValue, slackConversationKind)
		messageID := slackAdapter.PostAsync(
			plat.Post(convID, plat.Message("", valueViewAttachment)),
		)
		// meanwhile we'll perform analysis of the new competency description
		resp = analyseMessage(request, messageID, utils.TextAnalysisInput{
			Text:                       adaptiveValue.Description,
			OriginalMessageAttachments: []ebm.Attachment{valueViewAttachment},
			Namespace:                  namespace,
			Context:                    AdaptiveValuesDialogContext,
		})
	}
	return
}

func onDeleteButtonClicked(request slack.InteractionCallback, id string) (resp []models.PlatformSimpleNotification) {
	resp, ok := ensureUserHasWriteAccessToValues(request)
	if ok {
		err := adaptiveValuesTableDao.Delete(id)
		platform.ErrorHandler(request, fmt.Sprintf("Could not delete %s", id), err)
		resp = responses(
			utils.InteractionCallbackOverrideRequestMessage(request, AdaptiveValuesTemplate(DeletedAdaptiveValueNoticeSubject)),
		)
	}
	return
}

// func getOverrideMessageTs(request slack.InteractionCallback,
// 	slackConversationKind SlackConversationKind,
// ) (overrideMessageTs string) {
// 	switch(slackConversationKind){
// 	case InChat:
// 		overrideMessageTs = request.MessageTs
// 	case InThread:
// 		overrideMessageTs = noMessageOverrideTs // When in the thread, we do not override the message
// 	}
// 	return
// }

func onEditButtonClicked(
	request slack.InteractionCallback,
	teamID models.TeamID,
	id string,
	slackConversationKind SlackConversationKind,
) (resp []models.PlatformSimpleNotification) {
	resp, ok := ensureUserHasWriteAccessToValues(request)
	if ok {
		// ut := retrieveUserToken(request)
		adaptiveValue := adaptiveValuesTableDao.ReadUnsafe(id)
		platform.Debug(request, "Found adaptiveValue: "+adaptiveValue.Name)
		mc := callbackID(request.User.ID, SubmitUpdatedAdaptiveValueAction(slackConversationKind))
		mc.Target = id

		survey := utils.AttachmentSurvey(valueSurveyLabel, evalues.EditAdaptiveValueForm(&adaptiveValue))
		survey2 := ebm.AttachmentActionSurvey2{
			TriggerID: request.TriggerID,
			CallbackID: mc.ToCallbackID(),
			AttachmentActionSurvey: survey,
			State: request.MessageTs,
		}

		// Open a survey associated with the engagement
		err2 := slackAPI(teamID).ShowDialog(survey2)
		platform.Debug(request, "onEditButtonClicked: OpenDialog - done")
		platform.ErrorHandler(request,
			fmt.Sprintf("Could not open dialog from %s survey", request.CallbackID),
			err2,
		)
		resp = responses()
	}
	// utils.InteractionCallbackOverrideRequestMessage(request, AdaptiveValuesTemplate(EditingAdaptiveValueNoticeSubject))
	// When just pressing Edit button we don't want to update anything.
	return
}

func globalConnection(teamID models.TeamID) daosCommon.DynamoDBConnection {
	return daosCommon.DynamoDBConnection{
		Dynamo: d,
		ClientID: clientID,
		PlatformID: teamID.ToPlatformID(),
	}
}

func onSubmitUpdatedButtonClicked(
	request slack.InteractionCallback,
	dialog slack.DialogSubmissionCallback,
	teamID models.TeamID,
	slackConversationKind SlackConversationKind,
) (resp []models.PlatformSimpleNotification) {
	resp, ok := ensureUserHasWriteAccessToValues(request)
	if ok {
		platform.Debug(request, "Updating adaptiveValue "+dialog.Submission["Name"])
		mc := utils.MessageCallbackParseUnsafe(request.CallbackID, AdaptiveValuesNamespace)
		adaptiveValue := convertFormToAdaptiveValue(request, teamID, dialog.Submission)
		adaptiveValue.ID = mc.Target // we saved id in Target
		err := adaptiveValuesTableDao.Update(adaptiveValue)
		platform.ErrorHandler(request, "Updating adaptiveValue", err)
		platform.Debug(request, "Updated adaptiveValue "+dialog.Submission["Name"])
		conn := globalConnection(teamID)
		slackAdapter := mapper.SlackAdapterForTeamID(conn)
		convID := conversationID(request)
		messageForThread := plat.Message("", adaptiveValueInlineViewAttachment(request, adaptiveValue, InThread))
		messageForChatSpace := plat.Message("", adaptiveValueInlineViewAttachment(request, adaptiveValue, InChat))
		overrideMessageOptionalTs := dialog.State
		oldMsgID := plat.TargetMessageID{Ts: overrideMessageOptionalTs, ConversationID: convID}
		var msg plat.Response
		// we always post to the chat-space regardless of whether we started in the thread
		// the analysis will be posted to the new thread near this message.
		switch slackConversationKind {
		case InThread: // if we were in thread, we create new message in chat space with a separate new thread
			msg = plat.Post(convID, messageForChatSpace)
			// we also override the old message in the thread
			slackAdapter.PostAsync(plat.Override(oldMsgID, messageForThread))
		case InChat: // if we were in chat, we override the old message and use it's id for analysis thread
			msg = plat.Override(oldMsgID, messageForChatSpace)
		}
		// Anyway, we are posting to the chat space (probably overriding the old message)
		messageID := slackAdapter.PostAsync(msg)

		// If we are in thread, we also post the updated message into the thread,
		// overriding the message
		analysisInput := utils.TextAnalysisInput{
			Text:                       adaptiveValue.Description,
			OriginalMessageAttachments: []ebm.Attachment{adaptiveValueInlineViewAttachment(request, adaptiveValue, InChat)},
			Namespace:                  namespace,
			Context:                    AdaptiveValuesDialogContext,
		}
		resp = analyseMessage(request, messageID, analysisInput)
	}
	return
}

func ensureUserHasWriteAccessToValues(request slack.InteractionCallback) (resp []models.PlatformSimpleNotification, ok bool) {
	ok = DoesUserHaveWriteAccessToCompetencies(request.User.ID)
	if ok {
		log.Printf("Write access allowed %s", request.User.ID)
	} else {
		log.Printf("Write access denied %s", request.User.ID)
		resp = responses(
			utils.InteractionCallbackOverrideRequestMessage(request, "It looks like you don't have write access to Values at the moment"),
		)
	}
	return
}
