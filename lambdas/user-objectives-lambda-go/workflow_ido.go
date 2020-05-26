package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	"encoding/json"
	"errors"
	"fmt"
	common2 "github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/common"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	daosUser "github.com/adaptiveteam/adaptive/daos/user"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"log"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
)

const itemIDKey = "itemID"
const progressDateKey = "progressDate"
const isShowingDetailsKey = "isShowingDetails"
const isShowingProgressKey = "isShowingProgress"

// This file contains a generic mechanism for handling the creation of strategy objectives.

// CreateIDOWorkflow is a workflow template for creating IDOs
var CreateIDOWorkflow = wf.NamedTemplate{
	Name: "create-ido", Template: CreateIDOWorkflow_Workflow(),
}

// CreateIDOWorkflow_Workflow creates workflow template for creating IDO
func CreateIDOWorkflow_Workflow() wf.Template {
	log.Println("CreateIDOWorkflow_Workflow") // this should only appear once when lambda is started and `allRoutes` are constructed/
	return wf.Template{
		Init: "init", // initial state is "init". This is used when the user first triggers the workflow
		FSA: map[struct {
			wf.State
			wf.Event
		}]wf.Handler{
			{State: "init", Event: ""}:                               CreateIDOWorkflow_OnInit(true),
			{State: "form-shown", Event: "submit"}:                   wf.SimpleHandler(CreateIDOWorkflow_OnDialogSubmitted, "message-posted"),
			{State: "form-shown", Event: "cancel"}:                   wf.SimpleHandler(CreateIDOWorkflow_OnDialogCancelled, "done"),
			{State: "message-posted", Event: "message-id-available"}: wf.SimpleHandler(CreateIDOWorkflow_OnFieldsShown, "message-posted"), // returning to the same state for other events to trigger
			// // the following events are for buttons. Will be invoked not immediately
			{State: "message-posted", Event: "details"}:              wf.SimpleHandler(CreateIDOWorkflow_OnDetails, "message-posted"),
			{State: "message-posted", Event: "edit"}:                 wf.SimpleHandler(CreateIDOWorkflow_OnEdit, "form-shown"),
			{State: "message-posted", Event: "cancel"}:               wf.SimpleHandler(CreateIDOWorkflow_OnProgressCancel, "done"),
			{State: "message-posted", Event: "progressShow"}:         wf.SimpleHandler(CreateIDOWorkflow_OnProgressShow, "message-posted"),
			{State: "message-posted", Event: "progressIntermediate"}: wf.SimpleHandler(CreateIDOWorkflow_OnProgressIntermediate, "progress-form-shown"),
			{State: "progress-form-shown", Event: "submit"}:          wf.SimpleHandler(CreateIDOWorkflow_OnProgressFormSubmitted, "message-posted"),
			{State: "progress-form-shown", Event: "cancel"}:          wf.SimpleHandler(CreateIDOWorkflow_OnDialogCancelled, "done"), // NB! we handle on cancel using the same method
			{State: "message-posted", Event: "progressCloseout"}:     wf.SimpleHandler(CreateIDOWorkflow_OnProgressCloseout, "message-posted"),
			{State: "message-posted", Event: "add-another"}:          CreateIDOWorkflow_OnInit(false),
			{State: "init", Event: "view-idos"}:                      wf.SimpleHandler(CreateIDOWorkflow_OnViewIDOs, "message-posted"),
		},
		Parser: wf.Parser,
	}
}

func CreateIDOWorkflow_OnInit(isFromMainMenu bool) func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		log.Println("CreateIDOWorkflow_OnInit")
		out, err = CreateIDOWorkflow_ShowDialog(models.UserObjective{})(ctx)
		out.NextState = "form-shown"
		out.KeepOriginal = !isFromMainMenu
		return
	}
}

// CreateIDOWorkflow_ShowDialog handles an event and returns interaction for displaying the dialog box.
func CreateIDOWorkflow_ShowDialog(item models.UserObjective) func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		out.Interaction = CreateIDODialog(ctx, item)
		out.KeepOriginal = true // we have already deleted the message.
		return
	}
}

// CreateIDODialog returns interaction for displaying the dialog box.
func CreateIDODialog(ctx wf.EventHandlingContext, item models.UserObjective) wf.Interaction {
	conn := connGen.ForPlatformID(ctx.TeamID.ToPlatformID())
	coaches, dates, initiativesAndObjectives := LoadObjectiveDialogDictionaries(ctx.Request.User.ID, ctx.TeamID, item, conn)
	return wf.OpenSurvey(ObjectiveSurvey(item, coaches, dates, initiativesAndObjectives))
}

func convertKvPairToPlainTextOption(pairs []models.KvPair) (out []ebm.AttachmentActionElementPlainTextOption) {
	for _, p := range pairs {
		out = append(out, ebm.AttachmentActionElementPlainTextOption{Value: p.Value, Label: ui.PlainText(p.Key)})
	}
	return
}

func getCommentsFromProgress(objectiveProgress []models.UserObjectiveProgress) (comments []ui.RichText) {
	for _, each := range objectiveProgress {
		comments = append(comments, ui.Sprintf("%s (%s percent, [%s] status)", each.Comments, each.PercentTimeLapsed, models.ObjectiveStatusColorLabels[each.StatusColor]))
	}
	return
}

func getAccountabilityPartner(conn daosCommon.DynamoDBConnection) func (item models.UserObjective) ui.PlainText {
	return func (item models.UserObjective) ui.PlainText {
		return readUserDisplayName(conn, item.AccountabilityPartner)
	}
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
	return ui.PlainText(JoinRichText(comments, "\n"))
}

func readUserDisplayName(conn daosCommon.DynamoDBConnection, userID string) (displayName ui.PlainText) {
	if userID == "none" { 
		displayName = "Not needed" 
	} else {
		accountabilityPartners, err2 := daosUser.ReadOrEmpty(userID)(conn)

		if err2 == nil && len(accountabilityPartners) > 0{
			displayName = ui.PlainText(accountabilityPartners[0].DisplayName)
		} else {
			displayName = "Unknown"
			logger.Infof("Couldn't find AccountabilityPartner @%s, %+v\n", userID, err2)
		}
	}
	return
}

func attachmentFieldNewOld(label ui.PlainText, prop func(models.UserObjective) ui.PlainText, newItem, oldItem models.UserObjective) ebm.AttachmentField {
	return ebm.AttachmentField{
		Title: string(label),
		Value: strategy.NewAndOld(string(prop(newItem)), string(prop(oldItem))),
	}
}

func objectiveToFields(newItem, oldItem models.UserObjective, teamID models.TeamID) (fields []ebm.AttachmentField) {
	newTypeLabel, newAlignment := objectiveType(teamID)(newItem)
	oldTypeLabel, oldAlignment := objectiveType(teamID)(oldItem)
	getName := func(item models.UserObjective) ui.PlainText { return ui.PlainText(item.Name) }
	getDescription := func(item models.UserObjective) ui.PlainText { return ui.PlainText(item.Description) }
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(NameLabel, getName, newItem, oldItem),
		attachmentFieldNewOld(DescriptionLabel, getDescription, newItem, oldItem),
		{Title: string("Type"), Value: strategy.NewAndOld(string(newTypeLabel), string(oldTypeLabel))},
		{Title: string(StrategyAssociationFieldLabel), Value: strategy.NewAndOld(string(newAlignment), string(oldAlignment))},
		attachmentFieldNewOld(TimelineLabel, renderObjectiveViewDate, newItem, oldItem),
	}
	return
}

func objectiveToFieldsDetails(newItem, oldItem models.UserObjective, teamID models.TeamID, 
	conn daosCommon.DynamoDBConnection,
) (fields []ebm.AttachmentField) {
	// For ViewMore action, we only need the latest comment
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(AccountabilityPartnerLabel, getAccountabilityPartner(conn), newItem, oldItem),
		attachmentFieldNewOld(StatusLabel, getStatus, newItem, oldItem),
		attachmentField(LastReportedProgressLabel, getLatestComments(newItem)),
	}
	return
}

func CreateIDOWorkflow_OnDialogSubmitted(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	logger.Infof("CreateIDOWorkflow_OnDialogSubmitted")
	var item models.UserObjective
	var updated bool
	item, updated, err = extractObjectiveFromContext(ctx)
	logger.Infof("CreateIDOWorkflow_OnDialogSubmitted item=%v, isUpdating=%v", item, updated)
	var oldItem models.UserObjective

	if err == nil {
		if updated {
			oldItem = userObjectiveByID(item.ID)
			item.PlatformID = ctx.TeamID.ToPlatformID()
			item.Quarter = oldItem.Quarter
			item.Year = oldItem.Year
			item.CreatedDate = oldItem.CreatedDate
		} else {
			// Quarter, Year, CreatedDate have already been updated by extractObjectiveFromContext
			oldItem = item
		}
	}

	conn := connGen.ForPlatformID(ctx.TeamID.ToPlatformID())
	err = d.PutTableEntry(item, userObjectivesTable)

	if err == nil {
		view := viewObjectiveWritable(ctx, item, oldItem, conn)
		view.OverrideOriginal = updated
		out.Interaction = wf.Interaction{
			Messages: []wf.InteractiveMessage{view},
		}
		out.ImmediateEvent = "message-id-available" // this is needed to post analysis
		out.DataOverride = map[string]string{itemIDKey: item.ID}
	} else {
		logger.WithField("error", err).Errorf("CreateIDOWorkflow_OnDialogSubmitted error: %+v", err)
		out.Interaction = wf.SimpleResponses(
			platform.Post(platform.ConversationID(ctx.Request.User.ID),
				platform.MessageContent{Message: ui.Sprintf("Couldn't create an objective")},
			),
		)
		err = nil // we want to show error interaction
	}
	out.KeepOriginal = true
	out = out.WithRuntimeData(UserObjectiveKey, oldItem) // keeping the old item so that we'll be able to show it again after analysis.
	out.DataOverride[itemIDKey] = item.ID

	// Coach notifications
	if !updated {
		err = onUserObjectivePartnerSelection(conn, item, ctx)
	} else if item.AccountabilityPartner != oldItem.AccountabilityPartner {
		// There is a change in the coach
		err = onUserObjectivePartnerSelection(conn, item, ctx)
		if err == nil {
			// Send a notification to the old coach
			publish(models.PlatformSimpleNotification{UserId: oldItem.AccountabilityPartner, Channel: ctx.Request.Channel.ID,
				Message: fmt.Sprintf("%s chose to opt for a new accountability partner for the following objective: `%s`",
					common.TaggedUser(item.UserID), item.Name)})
		}
	}
	return
}

func onUserObjectivePartnerSelection(conn daosCommon.DynamoDBConnection, item models.UserObjective, ctx wf.EventHandlingContext) (err error) {
	mc := models.MessageCallback{
		Module: "objectives",
		Source: ctx.Request.User.ID,
		Topic:  "init",
		Action: "ask",
	}
	if item.AccountabilityPartner == "none" {
		// ignoring when coach is not needed
	} else if item.AccountabilityPartner == "requested" {
		// A user requested a coach. Post an engagement to coaching community.
		comm := community.CommunityById("coaching", ctx.TeamID, communitiesTable)
		publish(models.PlatformSimpleNotification{UserId: ctx.Request.User.ID, Channel: comm.ChannelID,
			Attachments: coachingCommAttachs(mc, item)})
		publish(models.PlatformSimpleNotification{UserId: ctx.Request.User.ID, Channel: ctx.Request.Channel.ID,
			Message: core.TextWrap(fmt.Sprintf(
				"I have sent a notification to the coaching community about your request for a coach for the objective: `%s`", item.Name),
				core.Underscore)})
	} else if item.AccountabilityPartner != core.EmptyString {
		// Send a notification to accountability partner if that person is willing to partner with you
		AskForPartnershipEngagement(ctx.TeamID, *mc.WithTopic("coaching").WithTarget(item.ID),
			item.AccountabilityPartner, fmt.Sprintf(
				"%s is requesting your coaching for the below Individual Development Objective. Are you available to partner with and guide your colleague with this effort?",
				common.TaggedUser(item.UserID)), fmt.Sprintf("*%s*: %s\n*%s*: %s", NameLabel, item.Name,
				DescriptionLabel, core.TextWrap(item.Description, core.Underscore)), "", "", false)

		var coachUts [] models.User
		coachUts, err = daosUser.ReadOrEmpty(item.AccountabilityPartner)(conn)
		if err == nil {
			if len(coachUts) > 0 {
				publish(models.PlatformSimpleNotification{UserId: ctx.Request.User.ID, Channel: ctx.Request.Channel.ID,
					Message: fmt.Sprintf("I have also sent a notification to your selected coach, %s, for confirmation",
						common.TaggedUser(item.AccountabilityPartner))})
			}
		}
	}
	return
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
	details := wf.Button("details", caption("Show less", "Show more")(isShowingDetails))
	progressShow := wf.MenuOption("progressShow", caption("Hide", "Show")(isShowingProgress))
	// addAnother := wf.Button("add-another", "Add another?")
	if isCompleted {
		buttons = wf.InteractiveElements(details, wf.InlineMenu("Progress", progressShow))
	} else {
		edit := wf.Button("edit", "Edit")
		cancel := wf.Button("cancel", "Cancel")
		cancel.Button.RequiresConfirmation = true
		progressIntermediate := wf.MenuOption("progressIntermediate", "Add/Update progress")
		progressCloseout := wf.MenuOption("progressCloseout", "Closeout")
		progress := wf.InlineMenu("Progress", progressShow, progressIntermediate, progressCloseout)
		buttons = wf.InteractiveElements(details, edit, progress, cancel)
	}
	return
}

func viewObjectiveWritable(ctx wf.EventHandlingContext, newItem models.UserObjective, oldItem models.UserObjective, conn daosCommon.DynamoDBConnection) wf.InteractiveMessage {
	isShowingDetails := ctx.GetFlag(isShowingDetailsKey)
	isShowingProgress := ctx.GetFlag(isShowingProgressKey)
	fields := objectiveToFields(newItem, oldItem, ctx.TeamID)
	if isShowingDetails {
		fields = append(fields, objectiveToFieldsDetails(newItem, oldItem, ctx.TeamID, conn)...)
	}
	if isShowingProgress {
		fields = append(fields, userObjectiveProgressField(newItem))
	}
	isCompleted := newItem.Completed == 1 && newItem.PartnerVerifiedCompletion
	return wf.InteractiveMessage{
		PassiveMessage: wf.PassiveMessage{
			Fields: fields,
		},
		InteractiveElements: objectiveWritableOperations(isShowingDetails, isShowingProgress, isCompleted),
	}
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

const UserObjectiveKey = "UserObjective"
// CreateIDOWorkflow_OnFieldsShown is triggered when
// the message with objective information has been shown.
// Now we want to run an analysis.
// TODO: start analysis earlier, in go-routine
func CreateIDOWorkflow_OnFieldsShown(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	logger.Infof("CreateIDOWorkflow_OnFieldsShown")
	// meanwhile we'll perform analysis of the new description
	messageID := channelizeID(toMapperMessageID(ctx.TargetMessageID))

	itemID := ctx.Data[itemIDKey]
	if itemIDKey == "" {
		err = errors.New("itemIDKey is not defined")
		return
	}
	item := userObjectiveByID(itemID)
	oldItem := item
	oldItemI, ok := ctx.TryGetRuntimeData(UserObjectiveKey)
	if ok {
		oldItem = oldItemI.(models.UserObjective)
	} else {
		logger.Infof("runtime data is empty")
	}
	conn := connGen.ForPlatformID(ctx.TeamID.ToPlatformID())

	// item, err2 := UserObjectiveDAO.Read(itemID)
	if err == nil {
		viewItem := viewObjectiveWritable(ctx, item, oldItem, conn)
		var resp wf.InteractiveMessage
		resp, err = wf.AnalyseMessage(dialogFetcherDAO, ctx.Request, messageID, utils.TextAnalysisInput{
			Text:                       item.Description,
			OriginalMessageAttachments: []ebm.Attachment{},
			Namespace:                  namespace,
			Context:                    IDODescriptionContext,
		},
			viewItem,
		)
		resp.OverrideOriginal = true
		if err == nil {
			out.Interaction.Messages = append(out.Interaction.Messages, resp)
		}
	}
	out.NextState = "done"
	out.KeepOriginal = true // we want to override it, so, not to delete
	return // we do not show anything else to the user
}

func extractObjectiveFromContext(ctx wf.EventHandlingContext) (item models.UserObjective, updated bool, err error) {
	form := ctx.Request.DialogSubmissionCallback.Submission
	var objectiveID string
	objectiveID, updated = ctx.Data[itemIDKey]
	if !updated {
		objectiveID = core.Uuid()
	}
	userID := ctx.Request.User.ID
	objName := form[objectives.ObjectiveName]
	objDescription := form[objectives.ObjectiveDescription]
	partner := form[objectives.ObjectiveAccountabilityPartner]
	endDate := form[objectives.ObjectiveEndDate]
	strategyEntityID := form[objectives.ObjectiveStrategyAlignment]
	// Get the alignment type for the aligned objective
	alignment, alignmentID := getAlignedStrategyTypeFromAlignmentID(strategyEntityID)
	year, quarter := core.CurrentYearQuarter()

	item = models.UserObjective{
		ID:                          objectiveID,
		UserID:                      userID,
		Name:                        objName,
		Description:                 objDescription,
		AccountabilityPartner:       partner,
		ObjectiveType:               models.IndividualDevelopmentObjective,
		StrategyAlignmentEntityID:   alignmentID,
		StrategyAlignmentEntityType: alignment,
		Quarter:                     quarter,
		Year:                        year,
		CreatedDate:                 core.ISODateLayout.Format(time.Now()),
		ExpectedEndDate:             endDate,
		PlatformID:                  ctx.TeamID.ToPlatformID(),
	}
	return
}

// CreateIDOWorkflow_OnDialogCancelled is triggered when the user pressed `cancel`
// in the dialog.
func CreateIDOWorkflow_OnDialogCancelled(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	logger.Infof("CreateIDOWorkflow_OnDialogCancelled")
	return
}

// LoadObjectiveDialogDictionaries loads dictionaries that are needed for objective dialog
func LoadObjectiveDialogDictionaries(
	userID string, teamID models.TeamID, 
	item models.UserObjective, conn daosCommon.DynamoDBConnection,
) (
	coaches, 
	dates []ebm.AttachmentActionElementPlainTextOption, 
	initiativesAndObjectives []ebm.AttachmentActionElementOptionGroup,
) {
	pValues := platformValues(teamID)
	logger.Infof("Retrieved values for %s platform: %v", teamID, pValues)
	initiativesAndObjectives = append(InitsAndObjectives(userID, teamID),
		pValues...)
	allMembers := IDOCoaches(userID, teamID, conn)
	allDates := objectives.DevelopmentObjectiveDates(namespace, item.ExpectedEndDate)
	coaches = convertKvPairToPlainTextOption(allMembers)
	dates = convertKvPairToPlainTextOption(allDates)
	return
}

// ObjectiveSurvey shows a form to create or modify an objective
func ObjectiveSurvey(item models.UserObjective,
	coaches, dates []ebm.AttachmentActionElementPlainTextOption,
	initiativesAndObjectives []ebm.AttachmentActionElementOptionGroup) ebm.AttachmentActionSurvey {
	alignment := objectives.AlignmentIDFromAlignedStrategyType(item.StrategyAlignmentEntityType, item.StrategyAlignmentEntityID)
	return ebm.AttachmentActionSurvey{
		Title: "Objective",
		Elements: []ebm.AttachmentActionTextElement{
			ebm.NewTextBox(objectives.ObjectiveName, "Name", objectives.ObjectiveNamePlaceholder, ui.PlainText(item.Name)),
			ebm.NewTextArea(objectives.ObjectiveDescription, "Description", objectives.ObjectiveDescriptionPlaceholder, ui.PlainText(item.Description)),
			ebm.NewSimpleOptionsSelect(objectives.ObjectiveAccountabilityPartner, "Coach", ebm.EmptyPlaceholder, string(item.AccountabilityPartner), coaches...),
			ebm.NewSimpleOptionsSelect(objectives.ObjectiveEndDate, "Expected end date", ebm.EmptyPlaceholder, item.ExpectedEndDate, dates...),
			ebm.NewSimpleOptionGroupsSelect(objectives.ObjectiveStrategyAlignment, "Strategy Alignment", ebm.EmptyPlaceholder, alignment,
				initiativesAndObjectives...),
		},
	}
}

// CreateIDOWorkflow_OnEdit -
func CreateIDOWorkflow_OnEdit(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			logger.Errorf("CreateIDOWorkflow_OnEdit panic recovered: %+v", err2)
			err = err2.(error)
		}
	}()

	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateIDOWorkflow_OnEdit itemID:%s", itemID)
	item := userObjectiveByID(itemID)
	out, err = CreateIDOWorkflow_ShowDialog(item)(ctx)
	return
}

func standardView(ctx wf.EventHandlingContext, item models.UserObjective) (out wf.EventOutput, err error) {
	conn := connGen.ForPlatformID(ctx.TeamID.ToPlatformID())

	view := viewObjectiveWritable(ctx, item, item, conn)
	view.OverrideOriginal = true
	out.Interaction = wf.Interaction{
		Messages: wf.InteractiveMessages(view),
	}
	out.KeepOriginal = true // we want to override it, so, not to delete
	return
}

func CreateIDOWorkflow_OnDetails(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateIDOWorkflow_OnDetails itemID:%s", itemID)
	ctx.ToggleFlag(isShowingDetailsKey)
	item := userObjectiveByID(itemID)

	return standardView(ctx, item)
}

func CreateIDOWorkflow_OnProgressCancel(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateIDOWorkflow_OnProgressCancel itemID:%s", itemID)
	item := userObjectiveByID(itemID)
	SetObjectiveField(item, "cancelled", 1)
	out.Interaction = wf.Interaction{
		Messages: wf.InteractiveMessages(wf.InteractiveMessage{
			PassiveMessage: wf.PassiveMessage{ // publish the message to the user
				Text: ui.Sprintf("Ok, cancelled the following objective: `%s`", item.Name),
			},
		}),
	}
	if item.Accepted == 1 { // post only if the objective has a coach
		out.Responses = append(out.Responses,
			platform.Post(platform.ConversationID(item.AccountabilityPartner),
				platform.MessageContent{
					Message: ui.Sprintf("<@%s> has cancelled the following objective: `%s`", item.UserID, item.Name),
				},
			),
		)
	}
	out.KeepOriginal = true // we want to keep it, because it might contain thread
	return
}

func CreateIDOWorkflow_OnProgressShow(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateIDOWorkflow_OnProgressShow itemID:%s", itemID)
	ctx.ToggleFlag(isShowingProgressKey)

	item := userObjectiveByID(itemID)
	return standardView(ctx, item)
}

func mapObjectiveProgressToRichText(ops []models.UserObjectiveProgress, f func(models.UserObjectiveProgress) ui.RichText) (texts []ui.RichText) {
	for _, each := range ops {
		texts = append(texts, f(each))
	}
	return
}

func getObjectiveProgressComment(op models.UserObjectiveProgress) ui.RichText {
	return ui.Sprintf("[%s] %s (%s)", models.ObjectiveStatusColorLabels[op.StatusColor], op.Comments, op.CreatedOn)
}

func userObjectiveProgressField(item models.UserObjective) (field ebm.AttachmentField) {
	ops, err2 := userObjectiveProgressByID(item.ID, -1)
	comments := mapObjectiveProgressToRichText(ops, getObjectiveProgressComment)
	var commentsJoined ui.RichText
	if err2 == nil {
		commentsJoined = ui.ListItems(comments...)
	} else {
		logger.Errorf("An error occurred while obtaining progress: %+v\n", err2)
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

// "made some progress"
func CreateIDOWorkflow_OnProgressIntermediate(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateIDOWorkflow_OnProgressIntermediate itemID:%s", itemID)

	var comments string
	var status models.ObjectiveStatusColor
	objectiveProgress := LatestProgressUpdateByObjectiveID(itemID)
	if len(objectiveProgress) > 0 {
		comments = objectiveProgress[0].Comments
		status = objectiveProgress[0].StatusColor
	}

	today := time.Now().Format(DateFormat)
	item := userObjectiveByID(itemID)
	label := ObjectiveProgressText2(item, today)

	survey := utils.AttachmentSurvey(string(label),
		progressCommentSurveyElements(ui.PlainText(item.Name), item.CreatedDate))
	surveyWithValues := fillCommentsSurveyValues(survey, comments, status)
	out.Interaction = wf.OpenSurvey(surveyWithValues)
	out.KeepOriginal = true
	return
}

func CreateIDOWorkflow_OnProgressFormSubmitted(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateIDOWorkflow_OnProgressFormSubmitted itemID:%s", itemID)
	item := userObjectiveByID(itemID)
	var progress models.UserObjectiveProgress
	progress, err = extractObjectiveProgressFromContext(ctx, item)
	err = d.PutTableEntry(progress, userObjectivesProgressTable)

	// attachs := viewProgressAttachment(mc,
	// 	ui.PlainText(Sprintf("This is your reported progress for the below %s", typLabel)),
	// 	"",
	// 	comments,
	// 	statusColor, item, models.Update)
	// publish(models.PlatformSimpleNotification{UserId: dialog.User.ID, Channel: dialog.Channel.ID, Ts: msgState.ThreadTs, Attachments: attachs})
	ctx.SetFlag(isShowingProgressKey, true) // enable show progress
	if err == nil {
		out, err = standardView(ctx, item)

		// Once submitted, post a view engagement to the coachee
		mc := models.MessageCallback{
			Source: ctx.Request.User.ID,
			Module: "objectives",
			Topic:  "Comments",
			Action: "confirm",
			Target: item.ID,
		}
		attachs := viewProgressAttachment(mc,
			ui.PlainText(ui.Sprintf("This is your reported progress for the below Development Objective")),
			"",
			ui.PlainText(progress.Comments),
			progress.StatusColor, item, models.Update)

		// Posting update view to the coachee
		conn := connGen.ForPlatformID(ctx.TeamID.ToPlatformID())
		slackAdapter := mapper.SlackAdapterForTeamID(conn)
		messageID, _ := slackAdapter.PostSync(
			platform.Post(platform.ConversationID(ctx.Request.User.ID), platform.Message("", attachs...)),
		)

		// Post an engagement to the coach about this update
		// Add an engagement for partner to review the progress. Coach is retrieved directly from objective
		PartnerReviewUserObjectiveProgressEngagement(ctx.TeamID, mc, item.AccountabilityPartner,
			progress.CreatedOn, ui.PlainText(progress.Comments), progress.StatusColor, item, false)

		msgState := MsgState{
			Ts:          messageID.Ts,
			ThreadTs:    messageID.Ts,
			ObjectiveId: item.ID,
		}
		// setting the message state to be used by legacy code
		msgStateBytes, _ := json.Marshal(msgState)
		ctx.Request.DialogSubmissionCallback.State = string(msgStateBytes)

		// doing analysis
		utils.ECAnalysis(progress.Comments, progressUpdateContext(item), "Progress update",
			dialogTableName, mc.ToCallbackID(), ctx.Request.User.ID, ctx.Request.Channel.ID, messageID.Ts,
			messageID.Ts, attachs, s, platformNotificationTopic, namespace)
	} else {
		logger.WithField("error", err).Errorf("CreateIDOWorkflow_OnProgressFormSubmitted error: %+v", err)
		out.Interaction = wf.SimpleResponses(
			platform.Post(platform.ConversationID(ctx.Request.User.ID),
				platform.MessageContent{Message: ui.Sprintf("Couldn't create an objective")},
			),
		)
		err = nil // we want to show error interaction
	}
	return
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
		PartnerID:         item.AccountabilityPartner,
		PlatformID:        item.PlatformID,
		PercentTimeLapsed: IntToString(percentTimeLapsed(today, item.CreatedDate, item.ExpectedEndDate)),
		StatusColor:       models.ObjectiveStatusColor(statusColor)}
	return
}

// CreateIDOWorkflow_OnProgressCloseout = "finished the objective"
func CreateIDOWorkflow_OnProgressCloseout(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	logger.Infof("CreateIDOWorkflow_OnProgressCloseout itemID:%s", itemID)

	item := userObjectiveByID(itemID)
	// If there is no partner assigned, send a message to the user that issue can't be closed-out until there is a coach
	if item.AccountabilityPartner == "requested" || item.AccountabilityPartner == "none" {
		out.Interaction.Messages = wf.InteractiveMessages(wf.InteractiveMessage{
			PassiveMessage: wf.PassiveMessage{
				Text: ui.Sprintf("You do not have a coach for the objective: `%s`. Please get a coach before attemping to close out.", item.Name),
			},
		})
	} else {
		item := userObjectiveByID(itemID)
		typLabel := objectiveTypeLabel(item)
		mc := models.MessageCallback{ // TODO: generate the correct MessageCallback for closeoutEng
			Module: "objectives",
			Target: item.ID,
			Source: item.UserID,
			Action: "ask_closeout", // will be replaced with `closeout`
			Topic:  "init",
		}
		// send a notification to the partner
		objectives.ObjectiveCloseoutEng(engagementTable, mc, item.AccountabilityPartner,
			fmt.Sprintf("%s wants to close the following %s. Are you ok with that?", common.TaggedUser(ctx.Request.User.ID), typLabel),
			fmt.Sprintf("*%s*: %s \n *%s*: %s", NameLabel, item.Name, DescriptionLabel, item.Description),
			string(closeoutLabel(item.ID)), objectiveCloseoutPath, false, dns, common2.EngagementEmptyCheck, ctx.TeamID)
		// Mark objective as closed
		SetObjectiveField(item, "completed", 1)
		SetObjectiveField(item, "completed_date", core.ISODateLayout.Format(time.Now()))

		// send a notification to the coachee that partner has been notified
		out.Interaction.Messages = wf.InteractiveMessages(wf.InteractiveMessage{
			PassiveMessage: wf.PassiveMessage{
				Text: ui.Sprintf("Awesome! I’ll schedule time with %s to close out the %s: `%s`",
					common.TaggedUser(item.AccountabilityPartner), typLabel, item.Name),
			},
		})
	}
	out.KeepOriginal = true // we want to override it, so, not to delete
	return
}

func CreateIDOWorkflow_OnViewIDOs(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	conn := connGen.ForPlatformID(ctx.TeamID.ToPlatformID())
	userID := ctx.Request.User.ID
	typ := models.IndividualDevelopmentObjective
	out.KeepOriginal = true // we want to override it, so, not to delete
	// Times in AWS are in UTC
	allObjs := objectives.AllUserObjectives(userID, userObjectivesTable, string(userObjective.UserIDTypeIndex), typ, 0)
	threadMessages := wf.InteractiveMessages()
	for _, obj := range allObjs {
		view := viewObjectiveWritable(ctx, obj, obj, conn)
		view.DataOverride = wf.Data{itemIDKey: obj.ID}
		threadMessages = append(threadMessages, view)
	}
	var msg ui.RichText
	if len(threadMessages) == 0 {
		msg = "You do not have any Individual Development Objectives yet."
	} else {
		msg = "You can find the list of your Individual Development Objectives in the thread. :point_down:"
	}
	out.Interaction.Messages = wf.InteractiveMessages(wf.InteractiveMessage{
		PassiveMessage: wf.PassiveMessage{
			Text:             msg,
			OverrideOriginal: true, // we override main menu message
		},
		Thread: threadMessages,
	})
	return
}
