package lambda

import (
	"encoding/json"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	aug "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/thoas/go-funk"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	// "github.com/adaptiveteam/adaptive/daos/strategyObjective"
	"log"
	"time"
)

const itemIDKey = "itemID"
const capCommIDKey = "capCommID"
const isShowingDetailsKey = "isShowingDetails"
const isShowingProgressKey = "isShowingProgress"

var CreateObjectiveWorkflow = wf.NamedTemplate{
	Name: "create-objective", Template: CreateObjectiveWorkflow_Workflow(),
}

const InitState wf.State = "init"
const (
	DefaultEvent wf.Event = ""
	ViewObjectivesEvent wf.Event = "ViewObjectivesEvent"
	ViewMyObjectivesEvent wf.Event = "ViewMyObjectivesEvent"
)
const MessagePostedState wf.State = "MessagePostedState"
const (
	MessageIDAvailableEvent wf.Event = "MessageIDAvailableEvent"
	EditEvent wf.Event = "EditEvent"
	AddAnotherEvent wf.Event = "AddAnotherEvent"
	DetailsEvent wf.Event = "DetailsEvent"
	CancelEvent wf.Event = "CancelEvent"
	ProgressShowEvent wf.Event = "ProgressShowEvent"
	ProgressIntermediateEvent wf.Event = "ProgressIntermediateEvent"
	ProgressCloseoutEvent wf.Event = "ProgressCloseoutEvent"

)
const FormShownState wf.State = "FormShownState"
const ProgressFormShownState wf.State = "ProgressFormShownState"

const CommunitySelectingState wf.State = "CommunitySelectingState"
const (
	CommunitySelectedEvent wf.Event = "CommunitySelectedEvent"
)
const ObjectiveShownState wf.State = "ObjectiveShownState"
const DoneState wf.State = "DoneState"

// This file contains a generic mechanism for handling the creation of strategy objectives.

func CreateObjectiveWorkflow_Workflow() wf.Template {
	log.Println("CreateObjectiveWorkflow_Workflow")
	return wf.Template{
		Init: InitState, // initial state is "init". This is used when the user first triggers the workflow
		FSA: map[struct {
			wf.State;
			wf.Event
		}]wf.Handler{
			{State: InitState, Event: DefaultEvent}:                         CreateObjectiveWorkflow_OnInit(true),
			{State: CommunitySelectingState, Event: CommunitySelectedEvent}: wf.SimpleHandler(CreateObjectiveWorkflow_OnCommunitySelected, FormShownState),
			{State: FormShownState, Event: wf.SurveySubmitted}:              wf.SimpleHandler(CreateObjectiveWorkflow_OnDialogSubmitted, MessagePostedState),
			{State: MessagePostedState, Event: MessageIDAvailableEvent}:     wf.SimpleHandler(CreateObjectiveWorkflow_OnFieldsShown, MessagePostedState), // returning to the same state for other events to trigger
			// the following events are for buttons. Will be invoked not immediately
			{State: MessagePostedState, Event: EditEvent}:                   wf.SimpleHandler(CreateObjectiveWorkflow_OnEdit, FormShownState),
			{State: MessagePostedState, Event: DetailsEvent}:                wf.SimpleHandler(CreateObjectiveWorkflow_OnDetails, MessagePostedState),
			{State: MessagePostedState, Event: CancelEvent}:                 wf.SimpleHandler(CreateObjectiveWorkflow_OnProgressCancel, DoneState),
			{State: MessagePostedState, Event: ProgressShowEvent}:           wf.SimpleHandler(CreateObjectiveWorkflow_OnProgressShow, MessagePostedState),
			{State: MessagePostedState, Event: ProgressIntermediateEvent}:   wf.SimpleHandler(CreateObjectiveWorkflow_OnProgressIntermediate, ProgressFormShownState),
			{State: MessagePostedState, Event: ProgressCloseoutEvent}:       wf.SimpleHandler(CreateObjectiveWorkflow_OnProgressCloseout, MessagePostedState),
			// {State: MessagePostedState, Event: "delete"}: wf.SimpleHandler(CreateObjectiveWorkflow_OnDelete, DoneState),
			{State: MessagePostedState, Event: AddAnotherEvent}:             CreateObjectiveWorkflow_OnInit(false),
			{State: ProgressFormShownState, Event: wf.SurveySubmitted}:      wf.SimpleHandler(CreateObjectiveWorkflow_OnProgressFormSubmitted, MessagePostedState),
			{State: ProgressFormShownState, Event: wf.SurveyCancelled}:      wf.SimpleHandler(CreateObjectiveWorkflow_OnDialogCancelled, DoneState), // NB! we handle on cancel using the same method
			{State: FormShownState, Event: wf.SurveyCancelled}:              wf.SimpleHandler(CreateObjectiveWorkflow_OnDialogCancelled, DoneState),
			{State: InitState, Event: ViewObjectivesEvent}:                  wf.SimpleHandler(CreateObjectiveWorkflow_OnViewObjectives(unfiltered), ObjectiveShownState),
			{State: InitState, Event: ViewMyObjectivesEvent}:                wf.SimpleHandler(CreateObjectiveWorkflow_OnViewObjectives(filterUserAdvocate), ObjectiveShownState),
			{State: ObjectiveShownState, Event: EditEvent}:                  wf.SimpleHandler(CreateObjectiveWorkflow_OnEdit, FormShownState),
			{State: ObjectiveShownState, Event: DetailsEvent}:               wf.SimpleHandler(CreateObjectiveWorkflow_OnDetails, MessagePostedState),
			{State: ObjectiveShownState, Event: CancelEvent}:                wf.SimpleHandler(CreateObjectiveWorkflow_OnProgressCancel, DoneState),
			{State: ObjectiveShownState, Event: ProgressShowEvent}:          wf.SimpleHandler(CreateObjectiveWorkflow_OnProgressShow, MessagePostedState),
			{State: ObjectiveShownState, Event: ProgressIntermediateEvent}:  wf.SimpleHandler(CreateObjectiveWorkflow_OnProgressIntermediate, ProgressFormShownState),
			{State: ObjectiveShownState, Event: ProgressCloseoutEvent}:      wf.SimpleHandler(CreateObjectiveWorkflow_OnProgressCloseout, MessagePostedState),
			{State: ObjectiveShownState, Event: AddAnotherEvent}:            CreateObjectiveWorkflow_OnInit(false),
		},
		Parser: wf.Parser,
	}
}

func CreateObjectiveWorkflow_OnInit(isFromMainMenu bool) func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		log.Println("CreateObjectiveWorkflow_OnInit")
		reply := simpleReply(ctx)
		if isMemberInCommunity(ctx.Request.User.ID, community.Strategy) {
			// check if the user is in strategy community
			adaptiveAssociatedCapComms := SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated(ctx.PlatformID)

			logger.Infof("Adaptive associated Capability Communities for platform %s: %v", ctx.PlatformID, adaptiveAssociatedCapComms)
			switch len(adaptiveAssociatedCapComms) {
			case 0:
				out = reply("There are no Adaptive associated Objective Communities. " +
						"If you have already created a Objective Community, " +
						"please ask the coordinator to create a *_private_* channel, " +
						"invite Adaptive and associate with the community.")
			case 1: // we already know the community. No need to ask.
				capCommID := adaptiveAssociatedCapComms[0].ID
				out, err = CreateObjectiveWorkflow_ShowDialog(capCommID, models.StrategyObjective{})(ctx)
				out.NextState = FormShownState
			default:
				out.NextState = CommunitySelectingState
				opts := mapCapabilityCommunitiesToOptions(adaptiveAssociatedCapComms, models.PlatformID(ctx.PlatformID))
				// Enable a user to create an objective if user is in strategy community and there are capability communities
				out.Interaction = wf.Buttons(
					"Select a capability community. You can assign the objective to other communities later but you need at least one for now.",
					wf.Selectors(wf.Selector{Event: CommunitySelectedEvent, Options: opts})...) // , wf.MenuOption("ignore", "Not now"))
			}
		} else {
			// send a message that user is not authorized to create objectives
			out = reply("You are not part of the Adaptive Strategy Community or an Objective Community, " +
					"you will not be able to create Capability Objectives.")
		}
		out.KeepOriginal = !isFromMainMenu
		return
	}
}

// CreateObjectiveWorkflow_OnCommunitySelected Is triggered when the user selected a community.
// It returns a dialog to the user.
func CreateObjectiveWorkflow_OnCommunitySelected(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	log.Println("CreateObjectiveWorkflow_OnCommunitySelected")
	capCommID, err := wf.SelectedValue(ctx.Request)
	if err == nil {
		out, err = CreateObjectiveWorkflow_ShowDialog(capCommID, models.StrategyObjective{})(ctx)
	}
	return
}

func simpleReply(ctx wf.EventHandlingContext) func (text ui.RichText) wf.EventOutput {
	return func (text ui.RichText) (out wf.EventOutput) {
		out.NextState = DoneState
		// send a message that user is not authorized to create objectives
		out.Interaction = wf.SimpleResponses(platform.Post(platform.ConversationID(ctx.Request.User.ID),
			platform.MessageContent{Message: text}))
		return 
}}

func mapCapabilityCommunitiesToOptions(comms []strategy.CapabilityCommunity, platformID models.PlatformID) (opts []wf.SelectorOption) {
	for _, each := range comms {
		eachComm := strategy.CapabilityCommunityByID(platformID, each.ID, capabilityCommunitiesTable)
		opts = append(opts, wf.SelectorOption{Label: ui.PlainText(eachComm.Name), Value: eachComm.ID})
	}
	return
}

func CreateObjectiveWorkflow_ShowDialog(capCommID string, item models.StrategyObjective) func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		types, advocates, dates := LoadObjectiveDialogDictionaries(ctx.PlatformID, capCommID, item)
		out.Interaction = wf.OpenSurvey(ObjectiveSurvey(item, types, advocates, dates))
		out.KeepOriginal = false                              // we have already deleted the message.
		out.Data = map[string]string{capCommIDKey: capCommID} // we'll need it when creating the objective
		return
	}
}

func convertKvPairToPlainTextOption(pairs []models.KvPair) (out []ebm.AttachmentActionElementPlainTextOption) {
	for _, p := range pairs {
		out = append(out, ebm.AttachmentActionElementPlainTextOption{Value: p.Value, Label: ui.PlainText(p.Key)})
	}
	return
}

func objectiveToFields(newSo, oldSo *models.StrategyObjective, platformID models.PlatformID) (kvs []ebm.AttachmentField) {
	if oldSo == nil {
		oldSo = newSo
	}
	newDate := formatDate(newSo.ExpectedEndDate, DateFormat, core.USDateLayout)
	oldDate := formatDate(oldSo.ExpectedEndDate, DateFormat, core.USDateLayout)

	kvs = []ebm.AttachmentField{
		{Title: string(SObjectiveTypeLabel), Value: strategy.NewAndOld(string(newSo.ObjectiveType), string(oldSo.ObjectiveType))},
		{Title: string(SObjectiveNameLabel), Value: strategy.NewAndOld(newSo.Name, oldSo.Name)},
		{Title: string(SObjectiveDescriptionLabel), Value: strategy.NewAndOld(newSo.Description, oldSo.Description)},
		{Title: SObjectiveAdvocateLabel, Value: strategy.NewAndOld(common.TaggedUser(newSo.Advocate), common.TaggedUser(oldSo.Advocate))},
		{Title: string(SObjectiveMeasuresLabel), Value: strategy.NewAndOld(newSo.AsMeasuredBy, oldSo.AsMeasuredBy)},
		{Title: string(SObjectiveTargetsLabel), Value: strategy.NewAndOld(newSo.Targets, oldSo.Targets)},
		{Title: SObjectiveEndDateLabel, Value: strategy.NewAndOld(newDate, oldDate)},
	}
	return
}

func attachmentFieldNewOld(label ui.PlainText, prop func(models.UserObjective) ui.PlainText, newItem, oldItem models.UserObjective) ebm.AttachmentField {
	return ebm.AttachmentField{
		Title: string(label),
		Value: strategy.NewAndOld(string(prop(newItem)), string(prop(oldItem))),
	}
}

func attachmentField(label ui.PlainText, value ui.PlainText) ebm.AttachmentField {
	return ebm.AttachmentField{
		Title: string(label), 
		Value: string(value),
	}
}

func objectiveToFieldsDetails(newIssue, oldIssue models.UserObjective, platformID models.PlatformID) (fields []ebm.AttachmentField) {
	// For ViewMore action, we only need the latest comment
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(AccountabilityPartnerLabel, getAccountabilityPartner, newIssue, oldIssue),
		attachmentFieldNewOld(StatusLabel, getStatus, newIssue, oldIssue),
		attachmentField(LastReportedProgressLabel, getLatestComments(newIssue)),
	}
	return
}

func getStatus(item models.UserObjective) (status ui.PlainText) {
	if item.Cancelled == 1 {
		status = StatusCancelled
	} else if item.Completed == 0 {
		status = StatusPending
	} else if item.Completed == 1 && item.PartnerVerifiedCompletion {
		status = StatusCompletedAndPartnerVerifiedCompletion
	} else if item.Completed == 1 && !item.PartnerVerifiedCompletion {
		status = StatusCompletedAndNotPartnerVerifiedCompletion
	}
	return
}

func getLatestComments(item models.UserObjective) (status ui.PlainText) {
	objectiveProgress := LatestProgressUpdateByObjectiveID(item.ID)
	comments := getCommentsFromProgress(objectiveProgress)
	return ui.PlainText(ui.Join(comments, "\n"))
}

func readUserDisplayName(userID string) (displayName ui.PlainText) {
	accountabilityPartner, err := utils.UserToken(userID, userProfileLambda, region, namespace)

	if err == nil {
		displayName = ui.PlainText(accountabilityPartner.DisplayName)
	} else {
		displayName = "Unknown"
		logger.Infof("Couldn't find AccountabilityPartner @" + userID)
	}
	return
}

func getCommentsFromProgress(objectiveProgress []models.UserObjectiveProgress) (comments []ui.RichText) {
	for _, each := range objectiveProgress {
		comments = append(comments, ui.Sprintf("%s (%s percent, [%s] status)", each.Comments, each.PercentTimeLapsed, models.ObjectiveStatusColorLabels[each.StatusColor]))
	}
	return
}

func getAccountabilityPartner(item models.UserObjective) ui.PlainText {
	return readUserDisplayName(item.AccountabilityPartner)
}

func getObjectiveProgressComment(op models.UserObjectiveProgress) ui.RichText {
	return ui.Sprintf("[%s] %s (%s)", models.ObjectiveStatusColorLabels[op.StatusColor], op.Comments, op.CreatedOn)
}

func userObjectiveProgressField(item models.UserObjective) (field ebm.AttachmentField) {
	ops, err := userObjectiveProgressByID(item.ID, -1)
	comments := mapObjectiveProgressToRichText(ops, getObjectiveProgressComment)
	var commentsJoined ui.RichText
	if err == nil {
		commentsJoined = ui.ListItems(comments...)
	} else {
		logger.Errorf("An error occurred while obtaining progress: %+v", err)
		commentsJoined = ""
	}

	progressTitle := "Progress" // ProgressTitle(item)
	// if commentsJoined == "" {
	// 	progressTitle = ProgressAbsentTitle(item)
	// }
	progressBody := commentsJoined
	if commentsJoined == "" {
		progressBody = "No progress"
	}
	return ebm.AttachmentField{
		Title: string(progressTitle),
		Value: string(progressBody),
		Short: true,
	}
}

func mapObjectiveProgressToRichText(ops []models.UserObjectiveProgress, f func(models.UserObjectiveProgress) ui.RichText) (texts []ui.RichText) {
	for _, each := range ops {
		texts = append(texts, f(each))
	}
	return
}

func EditStatusTemplate(updated bool) (text ui.RichText) {
	if updated {
		text = "updated"
	} else {
		text = "created"
	}
	return
}

func readItem(platformID models.PlatformID, itemID string) models.StrategyObjective {
	return strategy.StrategyObjectiveByID(platformID, itemID, strategyObjectivesTable)
}

func saveItem(platformID models.PlatformID, item models.StrategyObjective, capCommID string) (err error) {
	err = d.PutTableEntry(item, strategyObjectivesTable)
	if err == nil {
		uObj := UserObjectiveFromStrategyObjective(&item, capCommID, platformID)
		err = d.PutTableEntry(uObj, userObjectivesTable)
	}
	return
}

func runtimeData(d interface{}) *interface{} {return &d}

func getItemAndOldItem(ctx wf.EventHandlingContext) (item models.StrategyObjective, oldItem models.StrategyObjective, updated bool, err error) {
	item, updated, err = extractTypedObjectiveFromContext(ctx)
	if updated {
		oldItem = readItem(ctx.PlatformID, item.ID)
		item.ID = oldItem.ID
		item.CreatedBy = oldItem.CreatedBy
		item.CreatedAt = oldItem.CreatedAt
		item.CapabilityCommunityIDs = oldItem.CapabilityCommunityIDs
		// item.ModifiedBy = ctx.Request.User.ID
	} else {
		item.ID = core.Uuid()
		item.CreatedBy = ctx.Request.User.ID
		item.CreatedAt = core.CurrentRFCTimestamp()
		item.CapabilityCommunityIDs = []string{ctx.Data[capCommIDKey]}
		oldItem = item
	}
	return
}
func CreateObjectiveWorkflow_OnDialogSubmitted(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	logger.Infof("CreateObjectiveWorkflow_OnDialogSubmitted")
	reply := simpleReply(ctx)
	capCommID := ctx.Data[capCommIDKey]
	var item models.StrategyObjective
	var oldItem models.StrategyObjective
	var updated bool
	item, oldItem, updated, err = getItemAndOldItem(ctx)
	err = saveItem(ctx.PlatformID, item, capCommID)
	if err == nil {
		out = onNewItemAvailable(ctx, item, oldItem, updated, capCommID)
	} else {
		logger.WithField("error", err).Errorf("CreateObjectiveWorkflow_OnDialogSubmitted error: %+v", err)
		out = reply("Couldn't create an objective")
		err = nil // we want to show error interaction and we have logged the error
	}
	out.KeepOriginal = true
	out.RuntimeData = runtimeData(oldItem) // keeping the old item so that we'll be able to show it again after analysis.
	
	return
}

func onNewItemAvailable(ctx wf.EventHandlingContext, item models.StrategyObjective, oldItem models.StrategyObjective, updated bool, capCommID string) (out wf.EventOutput) {
	newIssue := UserObjectiveFromStrategyObjective(&item, item.CapabilityCommunityIDs[0], ctx.PlatformID)
	oldIssue := UserObjectiveFromStrategyObjective(&oldItem, oldItem.CapabilityCommunityIDs[0], ctx.PlatformID)

	view := viewObjectiveWritable(ctx, item, oldItem, *newIssue, *oldIssue)
	view.OverrideOriginal = updated
	out.Interaction = wf.Interaction{
		Messages: []wf.InteractiveMessage{view},
	}
	out.ImmediateEvent = MessageIDAvailableEvent // this is needed to post analysis
	out.Data = map[string]string{capCommIDKey: capCommID, itemIDKey: item.ID}
	msgToStrategyCommunity := viewObjectiveReadonly(ctx, item, oldItem)
	strategyCommunityConversation := findStrategyCommunityConversation(ctx)
	userID := ctx.Request.User.ID
	msgToStrategyCommunity.Message = ui.Sprintf("Below objective has been %s by <@%s>", EditStatusTemplate(updated), userID)
	notification := platform.Post(strategyCommunityConversation, msgToStrategyCommunity)
	logger.Infof("Notification to strategy community (%s): %v", strategyCommunityConversation, msgToStrategyCommunity)
	out.Responses = append(out.Responses, notification)
	return
}
func viewObjectiveWritable(ctx wf.EventHandlingContext, 
	newItem models.StrategyObjective, 
	oldItem models.StrategyObjective,
	newIssue models.UserObjective,
	oldIssue models.UserObjective) wf.InteractiveMessage {
	_, isShowingDetails := ctx.Data[isShowingDetailsKey]
	_, isShowingProgress := ctx.Data[isShowingProgressKey]
	fields := objectiveToFields(&newItem, &oldItem, ctx.PlatformID)
	if isShowingDetails {
		fields = append(fields, objectiveToFieldsDetails(newIssue, oldIssue, ctx.PlatformID)...)
	}
	if isShowingProgress {
		fields = append(fields, userObjectiveProgressField(newIssue))
	}
	isCompleted := newIssue.Completed == 1 && newIssue.PartnerVerifiedCompletion
	return wf.InteractiveMessage{
		PassiveMessage: wf.PassiveMessage{
			Fields: fields,
		},
		InteractiveElements: objectiveWritableOperations(isShowingDetails, isShowingProgress, isCompleted),
	}

	// edit := wf.Button(EditEvent, "Edit")
	// // delete := wf.Button("delete", "Delete")
	// // delete.Button.RequiresConfirmation = true
	// addAnother := wf.Button(AddAnotherEvent, "Add another?")
	// return wf.InteractiveMessage{
	// 	PassiveMessage: wf.PassiveMessage{
	// 		Fields: objectiveToFields(&item, &oldItem,),
	// 	},
	// 	InteractiveElements: []wf.InteractiveElement{edit, addAnother},
	// }
}

func caption(trueCaption ui.PlainText, falseCaption ui.PlainText) func(bool) ui.PlainText {
	return func(flag bool) (res ui.PlainText) {
		if flag {
			res = trueCaption
		} else {
			res = falseCaption
		}
		return
	}
}

func objectiveWritableOperations(isShowingDetails bool, isShowingProgress bool, isCompleted bool) (buttons []wf.InteractiveElement) {
	details := wf.Button(DetailsEvent, caption("Show less", "Show more")(isShowingDetails))
	progressShow := wf.MenuOption(ProgressShowEvent, caption("Hide", "Show")(isShowingProgress))
	// addAnother := wf.Button("add-another", "Add another?")
	if isCompleted {
		buttons = wf.InteractiveElements(details, wf.InlineMenu("Progress", progressShow))
	} else {
		edit := wf.Button(EditEvent, "Edit")
		cancel := wf.Button(CancelEvent, "Cancel")
		cancel.Button.RequiresConfirmation = true
		progressIntermediate := wf.MenuOption(ProgressIntermediateEvent, "Add/Update progress")
		progressCloseout := wf.MenuOption(ProgressCloseoutEvent, "Closeout")
		progress := wf.InlineMenu("Progress", progressShow, progressIntermediate, progressCloseout)
		buttons = wf.InteractiveElements(details, edit, progress, cancel)
	}
	return
}

func viewObjectiveReadonly(ctx wf.EventHandlingContext, item models.StrategyObjective, oldItem models.StrategyObjective) platform.MessageContent {
	return platform.MessageContent{
		Message: "",
		Attachments: []ebm.Attachment{
			{
				Fields: objectiveToFields(&item, &oldItem, ctx.PlatformID),
			},
		},
	}
}

func findStrategyCommunityConversation(ctx wf.EventHandlingContext) platform.ConversationID {
	comm := CommunityById("strategy", ctx.PlatformID)
	return platform.ConversationID(comm.Channel)
}

func channelizeID(msgID mapper.MessageID) (messageID chan mapper.MessageID) {
	messageID = make(chan mapper.MessageID, 1)
	messageID <- msgID
	return
}
func toMapperMessageID(id platform.TargetMessageID) mapper.MessageID {
	return mapper.MessageID{
		ConversationID: id.ConversationID,
		Ts:             id.Ts,
	}
}

// CreateObjectiveWorkflow_OnFieldsShown is triggered when
// the message with objective information has been shown.
// Now we want to run an analysis.
// TODO: start analysis earlier, in go-routine
func CreateObjectiveWorkflow_OnFieldsShown(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	logger.Infof("CreateObjectiveWorkflow_OnFieldsShown")
	// meanwhile we'll perform analysis of the new description
	messageID := channelizeID(toMapperMessageID(ctx.TargetMessageID))

	itemID := ctx.Data[itemIDKey]
	item := strategy.StrategyObjectiveByID(ctx.PlatformID, itemID, strategyObjectivesTable)
	oldItem := item
	newIssue := UserObjectiveFromStrategyObjective(&item, item.CapabilityCommunityIDs[0], ctx.PlatformID)
	if ctx.RuntimeData == nil {
		logger.Infof("runtime data is empty")
	} else {
		oldItem = (*ctx.RuntimeData).(models.StrategyObjective)
	}
	oldIssue := UserObjectiveFromStrategyObjective(&oldItem, oldItem.CapabilityCommunityIDs[0], ctx.PlatformID)
	// item, err := strategyObjectiveDAO.Read(itemID)
	if err == nil {
		viewItem := viewObjectiveWritable(ctx, item, oldItem, *newIssue, *oldIssue)
		var resp wf.InteractiveMessage
		resp, err = wf.AnalyseMessage(dialogFetcherDAO, ctx.Request, messageID, utils.TextAnalysisInput{
			Text:                       item.Description,
			OriginalMessageAttachments: []ebm.Attachment{},
			Namespace:                  namespace,
			Context:                    stratObjDescriptionContext,
		},
			viewItem,
		)
		resp.OverrideOriginal = true
		if err == nil {
			out.Interaction.Messages = append(out.Interaction.Messages, resp)
		}
	}
	out.NextState = DoneState
	out.KeepOriginal = true // we want to override it, so, not to delete
	out.Data = ctx.Data
	return // we do not show anything else to the user
}
func extractTypedObjectiveFromContext(ctx wf.EventHandlingContext) (item models.StrategyObjective, updated bool, err error) {
	item.ID, updated = ctx.Data[itemIDKey]
	item.Name = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveName]
	item.ObjectiveType = models.StrategyObjectiveType(ctx.Request.DialogSubmissionCallback.Submission[SObjectiveType])
	item.Description = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveDescription]
	item.AsMeasuredBy = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveMeasures]
	item.Targets = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveTargets]
	item.Advocate = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveAdvocate]
	item.ExpectedEndDate = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveEndDate]

	item.PlatformID = ctx.PlatformID
	return
}

func CreateObjectiveWorkflow_OnDialogCancelled(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	logger.Infof("CreateObjectiveWorkflow_OnDialogCancelled")
	// out.Interaction = wf.SimpleResponses(
	// 	platform.Post(platform.ConversationID(ctx.Request.User.ID), 
	// 		platform.MessageContent{Message: "Dialog cancelled"},
	// 	),
	// )
	return
}

// LoadObjectiveDialogDictionaries loads dictionaries that are needed for objective dialog
func LoadObjectiveDialogDictionaries(platformID models.PlatformID, capCommID string, item models.StrategyObjective) (types, advocates, dates []ebm.AttachmentActionElementPlainTextOption) {
	allMembers := communityMembersIncludingStrategyMembers(fmt.Sprintf("%s:%s", community.Capability, capCommID), platformID)
	allDates := objectives.StrategyObjectiveDatesWithIndefiniteOption(namespace, item.ExpectedEndDate)
	advocates = convertKvPairToPlainTextOption(allMembers)
	dates = convertKvPairToPlainTextOption(allDates)
	types = convertKvPairToPlainTextOption(ObjectiveTypes())
	return
}

// ObjectiveSurvey shows a form to create or modify an objective
func ObjectiveSurvey(item models.StrategyObjective,
	types, advocates, dates []ebm.AttachmentActionElementPlainTextOption) ebm.AttachmentActionSurvey {
	return ebm.AttachmentActionSurvey{
		Title: "Objective",
		Elements: []ebm.AttachmentActionTextElement{
			ebm.NewSimpleOptionsSelect(SObjectiveType, SObjectiveTypeLabel, ebm.EmptyPlaceholder, string(item.ObjectiveType), types...),
			ebm.NewTextBox(SObjectiveName, SObjectiveNameLabel, ebm.EmptyPlaceholder, ui.PlainText(item.Name)),
			ebm.NewTextArea(SObjectiveDescription, SObjectiveDescriptionLabel, ebm.EmptyPlaceholder, ui.PlainText(item.Description)),
			ebm.NewTextArea(SObjectiveMeasures, SObjectiveMeasuresLabel, ebm.EmptyPlaceholder, ui.PlainText(item.AsMeasuredBy)),
			ebm.NewTextArea(SObjectiveTargets, SObjectiveTargetsLabel, ebm.EmptyPlaceholder, ui.PlainText(item.Targets)),
			ebm.NewSimpleOptionsSelect(SObjectiveAdvocate, SObjectiveAdvocateLabel, ebm.EmptyPlaceholder, item.Advocate, advocates...),
			ebm.NewSimpleOptionsSelect(SObjectiveEndDate, SObjectiveEndDateLabel, ebm.EmptyPlaceholder, item.ExpectedEndDate, dates...),
		},
	}
}

func CreateObjectiveWorkflow_OnEdit(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			logger.Errorf("CreateObjectiveWorkflow_OnEdit panic recovered: %+v", err2)
			err = err2.(error)
		}
	}()

	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateObjectiveWorkflow_OnEdit itemID:%s", itemID)
	item := strategy.StrategyObjectiveByID(ctx.PlatformID, itemID, strategyObjectivesTable)
	bytes, _ := json.Marshal(item)
	logger.Infof("CreateObjectiveWorkflow_OnEdit item:%v", string(bytes))
	if len(item.CapabilityCommunityIDs) < 1 {
		logger.Infof("CreateObjectiveWorkflow_OnEdit CapabilityCommunityIDs is empty")
	} else {
		out, err = CreateObjectiveWorkflow_ShowDialog(item.CapabilityCommunityIDs[0], item)(ctx)
	}
	out.Data = ctx.Data
	return
}

type ObjectivePredicate = func (models.StrategyObjective) bool

type ObjectivePredicateFactory = func (wf.EventHandlingContext) func (models.StrategyObjective) bool

// funk.Filter is used instead of the following function.
// func filterObjectives(objs []models.StrategyObjective, predicate ObjectivePredicate) (filtered []models.StrategyObjective) {
// 	for _, each := range objs {
// 		if predicate(each) {
// 			filtered = append(filtered, each)
// 		}
// 	}
// 	return
// }

func CreateObjectiveWorkflow_OnViewObjectives(objectivesFilter ObjectivePredicateFactory) func (ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	return func (ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		logger.Infof("CreateObjectiveWorkflow_OnViewObjectives")
			// userID := ctx.Request.User.ID
		out.Data = ctx.Data
		out.KeepOriginal = true // we want to override it, so, not to delete
		// Times in AWS are in UTC
		items := strategy.AllOpenStrategyObjectives(ctx.PlatformID, strategyObjectivesTable, strategyObjectivesPlatformIndex,
			userObjectivesTable)
		logger.Infof("CreateObjectiveWorkflow_OnViewObjectives items.len %d", len(items))
		filteredItems := funk.Filter(items, objectivesFilter(ctx)).([]models.StrategyObjective)
		logger.Infof("CreateObjectiveWorkflow_OnViewObjectives filteredItems.len %d", len(filteredItems))
		threadMessages := wf.InteractiveMessages()
		for _, item := range filteredItems {
			newIssue := UserObjectiveFromStrategyObjective(&item, item.CapabilityCommunityIDs[0], ctx.PlatformID)
			view := viewObjectiveWritable(ctx, item, item, *newIssue, *newIssue)
			view.DataOverride = wf.Data{itemIDKey: item.ID}
			threadMessages = append(threadMessages, view)
		}
		var msg ui.RichText
		if len(threadMessages) == 0 {
			msg = "There are no objectives in the strategy yet."
		} else {
			msg = "You can find the list of Strategy Objectives in the thread. :point_down:"
		}
		out.Interaction.Messages = wf.InteractiveMessages(wf.InteractiveMessage{
			PassiveMessage: wf.PassiveMessage{
				Text: msg,
				OverrideOriginal: true,
			},
			Thread:         threadMessages,
		})
		return
	}
}

func filterUserAdvocate(ctx wf.EventHandlingContext) ObjectivePredicate { return func(item models.StrategyObjective) bool {
	return item.Advocate == ctx.Request.User.ID
}}

func unfiltered(wf.EventHandlingContext) ObjectivePredicate { return func(models.StrategyObjective) bool {
	return true
}}

func toggleContextFlag(ctx wf.EventHandlingContext, flag string) {
	_, isOn := ctx.Data[flag]
	if isOn {
		delete(ctx.Data, flag) // removing "flag"
	} else {
		ctx.Data[flag] = "true" // setting "flag"
	}
}

func standardView(ctx wf.EventHandlingContext, item models.StrategyObjective) (out wf.EventOutput, err error) {
	newIssue := UserObjectiveFromStrategyObjective(&item, item.CapabilityCommunityIDs[0], ctx.PlatformID)
	view := viewObjectiveWritable(ctx, item, item, *newIssue, *newIssue)
	view.OverrideOriginal = true
	out.Interaction = wf.Interaction{
		Messages: wf.InteractiveMessages(view),
	}
	out.Data = ctx.Data
	out.KeepOriginal = true // we want to override it, so, not to delete
	return
}

func CreateObjectiveWorkflow_OnDetails(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateObjectiveWorkflow_OnDetails itemID:%s", itemID)
	toggleContextFlag(ctx, isShowingDetailsKey)
	item := readItem(ctx.PlatformID, itemID)
	return standardView(ctx, item)
}

// TODO: just use item.ID
func getIssueID(so models.StrategyObjective) (id string) {
	id = so.ID
	// if len(so.CapabilityCommunityIDs) > 0 {
	// 	commID := so.CapabilityCommunityIDs[0]
	// 	if commID != "" {
	// 		id = fmt.Sprintf("%s_%s", so.ID, commID)
	// 	}
	// } else {
	// }
	return
}

func CreateObjectiveWorkflow_OnProgressCancel(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateIDOWorkflow_OnProgressCancel itemID:%s", itemID)
	item := readItem(ctx.PlatformID, itemID)
	issueID := getIssueID(item)
	issue := userObjectiveByID(issueID)
	//issue.Cancelled  
	SetObjectiveField(issue, "cancelled", 1)
	out.Interaction = wf.Interaction{
		Messages: wf.InteractiveMessages(wf.InteractiveMessage{
			PassiveMessage: wf.PassiveMessage{ // publish the message to the user
				Text: ui.Sprintf("Ok, cancelled the following objective: `%s`", item.Name),
			},
		}),
	}
	if issue.Accepted == 1 { // post only if the objective has a coach
		out.Responses = append(out.Responses,
			platform.Post(platform.ConversationID(issue.AccountabilityPartner),
				platform.MessageContent{
					Message: ui.Sprintf("<@%s> has cancelled the following objective: `%s`", issue.UserID, item.Name),
				},
			),
		)
	}
	out.KeepOriginal = true // we want to keep it, because it might contain thread
	return
}

func CreateObjectiveWorkflow_OnProgressShow(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateIDOWorkflow_OnProgressShow itemID:%s", itemID)
	toggleContextFlag(ctx, isShowingProgressKey)
	item := readItem(ctx.PlatformID, itemID)
	return standardView(ctx, item)
}

func CreateObjectiveWorkflow_OnProgressIntermediate(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateIDOWorkflow_OnProgressIntermediate itemID:%s", itemID)

	comments := ""
	status := models.ObjectiveStatusRedKey
	objectiveProgress := LatestProgressUpdateByObjectiveID(itemID)
	if len(objectiveProgress) > 0 {
		comments = objectiveProgress[0].Comments
		status = objectiveProgress[0].StatusColor
	}

	today := core.ISODateLayout.Format(time.Now())
	item := readItem(ctx.PlatformID, itemID)
	issueID := getIssueID(item)
	issue := userObjectiveByID(issueID)
	label := ObjectiveProgressText2(issue, today)

	survey := utils.AttachmentSurvey(string(label),
		progressCommentSurveyElements(ui.PlainText(item.Name), issue.CreatedDate))
	surveyWithValues := fillCommentsSurveyValues(survey, comments, status)
	out.Interaction = wf.OpenSurvey(surveyWithValues)
	out.KeepOriginal = true
	out.Data = ctx.Data
	return
}

func CreateObjectiveWorkflow_OnProgressFormSubmitted(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateIDOWorkflow_OnProgressFormSubmitted itemID:%s", itemID)
	item := readItem(ctx.PlatformID, itemID)
	issueID := getIssueID(item)
	issue := userObjectiveByID(issueID)
	progress, err := extractObjectiveProgressFromContext(ctx, issue)
	err = d.PutTableEntry(progress, userObjectivesProgressTable)

	// attachs := viewProgressAttachment(mc,
	// 	ui.PlainText(Sprintf("This is your reported progress for the below %s", typLabel)),
	// 	"",
	// 	comments,
	// 	statusColor, item, models.Update)
	// publish(models.PlatformSimpleNotification{UserId: dialog.User.ID, Channel: dialog.Channel.ID, Ts: msgState.ThreadTs, Attachments: attachs})
	ctx.Data[isShowingProgressKey] = "true" // enable show progress
	if err == nil {
		out, err = standardView(ctx, item)

		// Once submitted, post a view engagement to the coachee
		mc := models.MessageCallback{
			Source: ctx.Request.User.ID,
			Module: namespace,
			Topic:  "Comments",
			Action: "confirm",
			Target: item.ID,
		}
		attachs := viewProgressAttachment(mc,
			ui.PlainText(ui.Sprintf("This is your reported progress for the below Development Objective")),
			"",
			ui.PlainText(progress.Comments),
			progress.StatusColor, issue, models.Update)

		// Posting update view to the coachee
		slackAdapter := platformAdapter.ForPlatformID(ctx.PlatformID)
		messageID, _ := slackAdapter.PostSync(
			platform.Post(platform.ConversationID(ctx.Request.User.ID), platform.Message("", attachs...)),
		)

		// Post an engagement to the coach about this update
		// Add an engagement for partner to review the progress. Coach is retrieved directly from objective
		PartnerReviewUserObjectiveProgressEngagement(ctx.PlatformID, mc, issue.AccountabilityPartner,
			progress.CreatedOn, ui.PlainText(progress.Comments), progress.StatusColor, issue, false)

		msgState := MsgState{
			//Ts:          messageID.Ts,
			ThreadTs:    messageID.Ts,
			Id: item.ID,
		}
		// setting the message state to be used by legacy code
		msgStateBytes, _ := json.Marshal(msgState)
		ctx.Request.DialogSubmissionCallback.State = string(msgStateBytes)

		// doing analysis
		utils.ECAnalysis(progress.Comments, progressUpdateContext(issue), "Progress update",
			dialogTableName, mc.ToCallbackID(), ctx.Request.User.ID, ctx.Request.Channel.ID, messageID.Ts,
			messageID.Ts, attachs, s, platformNotificationTopic, namespace)
	} else {
		logger.WithField("error", err).Errorf("CreateObjectiveWorkflow_OnProgressFormSubmitted error: %+v", err)
		out.Interaction = wf.SimpleResponses(
			platform.Post(platform.ConversationID(ctx.Request.User.ID),
				platform.MessageContent{Message: ui.Sprintf("Couldn't create an objective")},
			),
		)
		err = nil // we want to show error interaction
	}
	return
}

func ObjectiveProgressText2(objective models.UserObjective, today string) ui.PlainText {
	var labelText ui.PlainText
	if objective.ExpectedEndDate == common.StrategyIndefiniteDateValue {
		labelText = "Progress"
	} else {
		percentElapsed := percentTimeLapsed(today, objective.CreatedDate, objective.ExpectedEndDate)
		labelText = ui.PlainText(ui.Sprintf("Time used - %d %%", percentElapsed))
	}
	return labelText
}

func percentTimeLapsed(today, start, end string) (percent int) {
	d1 := common.DurationDays(start, today, core.ISODateLayout,
		namespace)
	if end == common.StrategyIndefiniteDateValue {
		percent = 0
	} else {
		d2 := common.DurationDays(start, end, core.ISODateLayout, namespace)
		percent = int(float32(d1) / float32(d2) * float32(100))	
	}
	return 
}

const (
	ObjectiveProgressComments                         = "objective_progress"
	ObjectiveProgressCommentsPlaceholder ui.PlainText = ebm.EmptyPlaceholder
)

const (
	SlackLabelLimit = 48
)
const (
	ObjectiveStatusColor       = "objective_status_color"
	ObjectiveCloseoutComment   = "objective_closeout_comment"
	ObjectiveNoCloseoutComment = "objective_no_closeout_comment"
	ReviewUserProgressSelect   = "review_user_progress_select"
	UberCoach                  = "uber_coach"
)

func progressCommentSurveyElements(objName ui.PlainText, startDate string) []ebm.AttachmentActionTextElement {
	today := core.ISODateLayout.Format(time.Now())
	nameConstrained := ObjectiveCommentsTitle(objName)
	elapsedDays := common.DurationDays(startDate, today, core.ISODateLayout, namespace)
	return []ebm.AttachmentActionTextElement{
		{
			Label:    string(ObjectiveStatusLabel(elapsedDays, startDate)),
			Name:     ObjectiveStatusColor,
			ElemType: models.MenuSelectType,
			Options:  utils.AttachActionElementOptions(models.ObjectiveStatusColorKeyValues),
			// Value:    string(models.ObjectiveStatusColorLabels[statusValue]), // it's necessary to fill afterwards
		},
		ebm.NewTextArea(ObjectiveProgressComments, nameConstrained, ObjectiveProgressCommentsPlaceholder, ""),
	}
}

func limitPlainText(text ui.PlainText, maxLength int) ui.PlainText {
	if len(text) < maxLength {
		return text
	}
	return text[:maxLength-3] + "..."
}

func ObjectiveCommentsTitle(objName ui.PlainText) ui.PlainText {
	nameConstrained := limitPlainText(ui.PlainText("Comments on "+objName), SlackLabelLimit)
	return nameConstrained
}

func ObjectiveStatusLabel(elapsedDays int, startDate string) ui.PlainText {
	return ui.PlainText(ui.Sprintf("Status (%d days since %s)", elapsedDays, startDate))
}

func ObjectiveProgressText(objective models.UserObjective, today string) ui.RichText {
	timeUsed := fmt.Sprintf("%d days elapsed since %s",
		common.DurationDays(objective.CreatedDate, today, core.ISODateLayout, namespace), objective.CreatedDate)
	fmt.Printf("Time used for %s objective: %s", objective.Name, timeUsed)
	if objective.ExpectedEndDate == common.StrategyIndefiniteDateValue {
		return ui.Sprintf("%s", objective.Name)
	} else {
		return ui.Sprintf("%s", objective.Name)
	}
}

func fillCommentsSurveyValues(sur ebm.AttachmentActionSurvey, comments string, status models.ObjectiveStatusColor) ebm.AttachmentActionSurvey {
	return models.FillSurvey(sur, map[string]string{
		ObjectiveProgressComments: comments,
		ObjectiveStatusColor:      string(status),
	})
}
func extractObjectiveProgressFromContext(ctx wf.EventHandlingContext, item models.UserObjective) (progress models.UserObjectiveProgress, err error) {
	form := ctx.Request.DialogSubmissionCallback.Submission

	comments := form[ObjectiveProgressComments]
	statusColor := form[ObjectiveStatusColor]
	today := core.ISODateLayout.Format(time.Now())

	progress = models.UserObjectiveProgress{
		ID:                item.ID,
		CreatedOn:         today,
		UserID:            ctx.Request.User.ID,
		Comments:          comments,
		PlatformID:        item.PlatformID,
		PartnerID:         item.AccountabilityPartner,
		PercentTimeLapsed: IntToString(percentTimeLapsed(today, item.CreatedDate, item.ExpectedEndDate)),
		StatusColor:       models.ObjectiveStatusColor(statusColor)}
	return
}
func viewProgressAttachment(mc models.MessageCallback, title, fallback ui.PlainText, comments ui.PlainText,
	status models.ObjectiveStatusColor, obj models.UserObjective, actionName models.AttachActionName) []ebm.Attachment {
	attach := utils.ChatAttachment(string(title), core.EmptyString, string(fallback), mc.ToCallbackID(),
		updateProgressAttachmentActions(mc, actionName, ui.PlainText(obj.Name),
			obj.CreatedDate, obj.ExpectedEndDate),
		progressFields(comments, status, obj), time.Now().Unix())
	return []ebm.Attachment{*attach}
}

func progressFields(comments ui.PlainText, status models.ObjectiveStatusColor, obj models.UserObjective) []ebm.AttachmentField {
	today := core.ISODateLayout.Format(time.Now())
	timeProgressLabel := ObjectiveProgressText(obj, today)
	return models.AttachmentFields([]models.KvPair{
		{
			Key:   NameLabel,
			Value: string(timeProgressLabel),
		},
		{
			Key:   DescriptionLabel,
			Value: obj.Description,
		},
		// {
		//	Key:   "Strategy Association(s)",
		//	Value: objectiveType(obj),
		// },
		{
			Key:   string(ProgressStatusLabel),
			Value: string(models.ObjectiveStatusColorLabels[status]),
		},
		{
			Key:   string(CommentsLabel),
			Value: string(comments),
		},
	})
}

func updateProgressAttachmentActions(mc models.MessageCallback, actionName models.AttachActionName, objName ui.PlainText, start, end string) []ebm.AttachmentAction {
	return []ebm.AttachmentAction{*models.GenAttachAction(mc, actionName,
		string(ObjectiveProgressChangeCommentsActionLabel), models.EmptyActionConfirm(), true)}
}

// IntToString converts int to string
func IntToString(i int) string {
	return fmt.Sprintf("%d", i)
}

func PartnerReviewUserObjectiveProgressEngagement(platformID models.PlatformID, mc models.MessageCallback, partner string,
	date string, objComments ui.PlainText, statusColor models.ObjectiveStatusColor, uObj models.UserObjective, urgent bool) {
	mc = *mc.WithModule(namespace).WithTopic("coaching").WithAction(ReviewCoacheeProgressAsk).
		WithTarget(fmt.Sprintf("%s_%s", uObj.ID, date)) // .WithTarget(uObj.ID)
	typLabel := objectiveTypeLabel(uObj)
	text := fmt.Sprintf("Here is the progress from %s on the below %s", common.TaggedUser(uObj.UserID), typLabel)

	emptyComments := objComments == ""
	titleString := core.IfThenElse(emptyComments, fmt.Sprintf("No progress reported for the below %s", typLabel), text)

	var actions []ebm.AttachmentAction
	if !emptyComments {
		actions = append(actions,
			*models.SimpleAttachAction(mc, models.Now, "Add my response"),
		)
	}
	actions = append(actions,
		*models.SimpleAttachAction(mc, models.Ignore, "Skip this"))

	utils.AddChatEngagement(mc, titleString.(string), core.EmptyString, "Adaptive at your service", partner, actions,
		progressFields(objComments, statusColor, uObj), platformID, urgent, engagementTable, d, namespace,
		time.Now().Unix(), 	aug.UserEngagementCheckWithValue{},
	)
}

func progressUpdateContext(userObj models.UserObjective) (context string) {
	switch userObj.ObjectiveType {
	case models.IndividualDevelopmentObjective:
		context = IDOProgressUpdateContext
	case models.StrategyDevelopmentObjective:
		switch userObj.StrategyAlignmentEntityType {
		case models.ObjectiveStrategyObjectiveAlignment:
			context = CapabilityObjectiveProgressUpdateContext
		case models.ObjectiveStrategyInitiativeAlignment:
			context = InitiativeProgressUpdateContext
		}
	}
	return
}

// func GetMsgStateUnsafe(request slack.InteractionCallback) (msgState MsgState) {
// 	err := json.Unmarshal([]byte(request.State), &msgState)
// 	core.ErrorHandler(err, namespace, "Couldn't unmarshal MsgState")
// 	return
// }

func CreateObjectiveWorkflow_OnProgressCloseout(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateObjectiveWorkflow_OnProgressCloseout itemID:%s", itemID)

	item := readItem(ctx.PlatformID, itemID)
	issueID := getIssueID(item)
	issue := userObjectiveByID(issueID)

	// If there is no partner assigned, send a message to the user that issue can't be closed-out until there is a coach
	if issue.AccountabilityPartner == "requested" || issue.AccountabilityPartner == "none" {
		out.Interaction.Messages = wf.InteractiveMessages(wf.InteractiveMessage{
			PassiveMessage: wf.PassiveMessage{
				Text: ui.Sprintf("You do not have a coach for the objective: `%s`. Please get a coach before attemping to close out.", item.Name),
			},
		})
	} else {
		typLabel := objectiveTypeLabel(issue)
		mc := models.MessageCallback{ // TODO: generate the correct MessageCallback for closeoutEng
			Module: "objectives",
			Target: item.ID,
			Source: issue.UserID,
			Action: "ask_closeout", // will be replaced with `closeout`
			Topic:  "init",
		}
		// send a notification to the partner
		objectives.ObjectiveCloseoutEng(engagementTable, mc, issue.AccountabilityPartner,
			fmt.Sprintf("<@%s> wants to close the following %s. Are you ok with that?", ctx.Request.User.ID, typLabel),
			fmt.Sprintf("*%s*: %s \n *%s*: %s", NameLabel, item.Name, DescriptionLabel, item.Description),
			string(closeoutLabel(item.ID)), objectiveCloseoutPath, false, dns, aug.UserEngagementCheckWithValue{}, ctx.PlatformID)
		// Mark objective as closed
		SetObjectiveField(issue, "completed", 1)

		// send a notification to the coachee that partner has been notified
		out.Interaction.Messages = wf.InteractiveMessages(wf.InteractiveMessage{
			PassiveMessage: wf.PassiveMessage{
				Text: ui.Sprintf("Awesome! I’ll schedule time with <@%s> to close out the %s: `%s`",
				issue.AccountabilityPartner, typLabel, item.Name),
			},
		})

	}

	out.Data = ctx.Data
	out.KeepOriginal = true // we want to override it, so, not to delete
	return
}


func progressLabel(userObjID string) ui.PlainText {
	return ui.PlainText("Responsibility Progress")
}

func closeoutLabel(userObjID string) ui.PlainText {
	return ui.PlainText("Responsibility Closeout")
}
var objectiveCloseoutPath            = ""// utils.NonEmptyEnv("USER_OBJECTIVES_CLOSEOUT_LEARN_MORE_PATH")

// func CreateObjectiveWorkflow_OnDelete(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
// 	itemID := ctx.Data[itemIDKey]
// 	logger.Infof("CreateObjectiveWorkflow_OnDelete itemID:%s", itemID)
// 	commParams := idAndPlatformIDParams(itemID, ctx.PlatformID)
// 	err = d.DeleteEntry(strategyObjectivesTable, commParams)

// 	out.Interaction = wf.Interaction{ 
// 		Messages: []wf.InteractiveMessage{
// 			{
// 				PassiveMessage: wf.PassiveMessage{
// 					AttachmentText: "Objective deleted",
// 					OverrideOriginal: true,
// 				},
// 			},
// 		},
// 	}
// 	out.KeepOriginal = true
// 	return
// }

// // CreateObjectiveWorkflow_OnAddAnother - switches workflow to the initial state
// func CreateObjectiveWorkflow_OnAddAnother(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
// 	logger.Infof("CreateObjectiveWorkflow_OnAddAnother")
// 	out.ImmediateEvent = "" // creating again
// 	out.Interaction.KeepOriginal = true
// 	// out.
// 	return
// }
const (
	NameLabel                            = "Name"
	DescriptionLabel                     = "Description"
	TimelineLabel                        = "Timeline"
	ProgressCommentsLabel   ui.PlainText = "Comments on Progress"
	ProgressStatusLabel     ui.PlainText = "Current Status"
	ObjectiveProgressLabel  ui.PlainText = "Objective Progress"
	PerceptionOfStatusLabel ui.PlainText = "Your perception of status"
	PerceptionOfStatusName = "perception_of_status"

	CommentsName                  = "Comments"
	CommentsLabel    ui.PlainText = "Comments"
	PercentDoneLabel ui.PlainText = "Percent Done"

	CommentsSurveyPlaceholder ui.PlainText = ebm.EmptyPlaceholder
	CommentsPlaceholder       ui.PlainText = ebm.EmptyPlaceholder

	ReviewCoacheeProgressAsk = "review_coachee_progress_ask"
	ReviewCoachComments      = "review_coach_comments"


	ObjectiveProgressChangeCommentsActionLabel ui.PlainText = "Change my comments"

	ObjectiveProgressChangeCommentsDialogTitle ui.PlainText = "Individual Objectives" // NB! this title might be irrelevant


)


func objectiveTypeLabel(userObj models.UserObjective) string {
	var prefix string
	switch userObj.ObjectiveType {
	case models.IndividualDevelopmentObjective:
		prefix = Individual
	case models.StrategyDevelopmentObjective:
		switch userObj.StrategyAlignmentEntityType {
		case models.ObjectiveStrategyObjectiveAlignment:
			prefix = CapabilityObjective
		case models.ObjectiveStrategyInitiativeAlignment:
			prefix = StrategyInitiative
		}
	}
	return prefix
}


const (
	IDODescriptionContext    = "dialog/ido/language-coaching/description"
	IDOProgressUpdateContext = "dialog/ido/language-coaching/update"

	IDOCloseoutDisagreementContext                 = "dialog/ido/language-coaching/close-out-disagreement"
	InitiativeCloseoutDisagreementContext          = "dialog/strategy/language-coaching/initiative/close-out-disagreement"
	CapabilityObjectiveCloseoutDisagreementContext = "dialog/strategy/language-coaching/objective/close-out-disagreement"

	IDOCloseoutAgreementContext                 = "dialog/ido/language-coaching/close-out-agreement"
	InitiativeCloseoutAgreementContext          = "dialog/strategy/language-coaching/initiative/close-out-agreement"
	CapabilityObjectiveCloseoutAgreementContext = "dialog/strategy/language-coaching/objective/close-out-agreement"

	IDOCoachingRejectionContext       = "dialog/ido/language-coaching/coaching-request-rejection"
	IDOResponseObjectiveUpdateContext = "dialog/ido/language-coaching/update-response"

	CapabilityObjectiveProgressUpdateContext = "dialog/strategy/language-coaching/objective/update"
	InitiativeProgressUpdateContext          = "dialog/strategy/language-coaching/initiative/update"
	CapabilityObjectiveUpdateResponseContext = "dialog/strategy/language-coaching/objective/update-response"
	InitiativeUpdateResponseContext          = "dialog/strategy/language-coaching/initiative/update-response"
	BlueDiamondEmoji                         = ":small_blue_diamond:"
)

const (
	Individual          = "Individual Objective"
	CapabilityObjective = "Capability Objective"
	StrategyInitiative  = "Initiative"
	FinancialObjective  = "Financial Objective"
	CustomerObjective   = "Customer Objective"
)

const (
	AccountabilityPartnerLabel ui.PlainText = "Accountability Partner"
	StatusLabel                ui.PlainText = "Status"
	LastReportedProgressLabel  ui.PlainText = "Last reported progress"
)

const (
	StatusCancelled                                ui.PlainText = "Cancelled"
	StatusPending                                  ui.PlainText = "Pending"
	StatusCompletedAndPartnerVerifiedCompletion    ui.PlainText = "Completed by you and closeout approved by your partner"
	StatusCompletedAndNotPartnerVerifiedCompletion ui.PlainText = "Completed by you and pending closeout approval from your partner"
)
