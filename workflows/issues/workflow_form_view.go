package issues

import (
	"time"

	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/pkg/errors"
)

type TextAnalysisInputFactory = func(ctx wf.EventHandlingContext) (analysisInput utils.TextAnalysisInput, err error)

type TextExtractor = func(newAndOldIssues NewAndOldIssues) (text string)

func ExtractDescription(newAndOldIssues NewAndOldIssues) (text string) {
	return newAndOldIssues.NewIssue.UserObjective.Description
}

func (w workflowImpl) textAnalysisInput(ctx *wf.EventHandlingContext, textExtractor TextExtractor, dialogSituationID DialogSituationIDWithoutIssueType) (analysisInput utils.TextAnalysisInput, err error) {
	var newAndOldIssues NewAndOldIssues
	newAndOldIssues, err = w.getNewAndOldIssues(*ctx) // userObjectiveByID(itemID)
	err = errors.Wrapf(err, "textAnalysisInput/getNewAndOldIssues")
	ctx.RuntimeData = runtimeData(newAndOldIssues)

	itype := newAndOldIssues.NewIssue.GetIssueType()

	if itype == "" {
		err = errors.Wrapf(err, "textAnalysisInput/itype is empty")
		return
	}
	if err != nil {
		return
	}
	context := w.GetDialogContext(dialogSituationID, itype)

	text := textExtractor(newAndOldIssues)
	analysisInput = utils.TextAnalysisInput{
		Text:                       text,
		OriginalMessageAttachments: []ebm.Attachment{},
		Namespace:                  "OnFieldsShown",
		Context:                    context,
	}
	return
}
func (w workflowImpl) OnFieldsShown(textExtractor TextExtractor, dialogSituationID DialogSituationIDWithoutIssueType) wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[issueIDKey]
		w.AdaptiveLogger.WithField("issueID", issueID).Info("OnFieldsShown")
		// meanwhile we'll perform analysis of the new description
		messageID := channelizeID(toMapperMessageID(ctx.TargetMessageID))

		out, err = w.standardView(ctx)
		viewItem := out.Messages[0]

		var textAnalysisInput utils.TextAnalysisInput
		textAnalysisInput, err = w.textAnalysisInput(&ctx, textExtractor, dialogSituationID)
		if err == nil {
			var resp wf.InteractiveMessage
			resp, err = wf.AnalyseMessage(w.DialogFetcherDAO, ctx.Request, messageID,
				textAnalysisInput, viewItem,
			)
			resp.OverrideOriginal = true
			if err == nil {
				out.Interaction.Messages = wf.InteractiveMessages(resp)
			}
		}
		out.NextState = "done"
		out.KeepOriginal = true // we want to override it, so, not to delete
		return
	}
}

func (w workflowImpl) standardView(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	w.AdaptiveLogger.Info("standardView")
	var newAndOldIssues NewAndOldIssues
	newAndOldIssues, err = w.getNewAndOldIssues(ctx)
	if err != nil {
		return
	}
	itype := IssueType(ctx.Data[issueTypeKey])
	tc := getTypeClass(itype)
	isShowingDetails := ctx.GetFlag(isShowingDetailsKey)
	isShowingProgress := ctx.GetFlag(isShowingProgressKey)

	fields := tc.View(w, isShowingDetails, isShowingProgress, newAndOldIssues)
	interactiveElements := objectiveWritableOperations(ctx, newAndOldIssues)
	createdAt := w.parseDateOrNow(newAndOldIssues.NewIssue.UserObjective.CreatedDate)
	if newAndOldIssues.NewIssue.UserObjective.ID == "" {
		w.AdaptiveLogger.Warnf("INVALID(3): issueID is empty %v\n", newAndOldIssues.NewIssue)
	}

	view := wf.InteractiveMessage{
		PassiveMessage: wf.PassiveMessage{
			Fields:             omitEmpty(fields),
			OverrideOriginal:   true,
			IsPermanentMessage: true, // we don't ever want to delete form view, the message to disappear from the thread
			Footer: ebm.AttachmentFooter{Text: "Created at", Timestamp: createdAt},
		},
		InteractiveElements: interactiveElements,
		DataOverride: wf.Data{
			issueIDKey: newAndOldIssues.NewIssue.UserObjective.ID,
			issueTypeKey: string(itype), //probably we don't need this because it's available
		},
	}
	out.Interaction = wf.Interaction{
		Messages: []wf.InteractiveMessage{view},
	}
	err = errors.Wrap(err, "{standardView}")
	return
}

func (w workflowImpl) parseTimestampOrNow(timestamp string) int64 {
	t, err2 := core.ISODateLayout.Parse(timestamp)
	if err2 != nil {
		w.AdaptiveLogger.WithError(err2).WithField("timestamp", timestamp).
			Error("Couldn't parse timestamp")
		t = time.Now()
	}
	return t.Unix()
}

func (w workflowImpl) parseDateOrNow(date string) int64 {
	t, err2 := time.Parse("2006-01-02", date)
	if err2 != nil {
		w.AdaptiveLogger.WithError(err2).WithField("date", date).
			Error("Couldn't parse date as 2006-01-02")
			t, err2 = core.TimestampLayout.Parse(date)
			if err2 != nil {
				w.AdaptiveLogger.WithError(err2).WithField("date", date).
					Error("Couldn't parse date as timestamp")
				t = time.Now()
			}
	}
	return t.Unix()
}

// omitEmpty removes fields with empty values
func omitEmpty(fields []ebm.AttachmentField) (res []ebm.AttachmentField) {
	for _, f := range fields {
		if f.Value != "" {
			res = append(res, f)
		}
	}
	return
}

func (w workflowImpl) OnDetails(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	ctx.ToggleFlag(isShowingDetailsKey)
	return w.standardView(ctx)
}

func (w workflowImpl) OnProgressShow(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	ctx.ToggleFlag(isShowingProgressKey)
	return w.standardView(ctx)
}

func objectiveWritableOperations(ctx wf.EventHandlingContext, newAndOldIssues NewAndOldIssues) (buttons []wf.InteractiveElement) {
	isShowingDetails := ctx.GetFlag(isShowingDetailsKey)
	isShowingProgress := ctx.GetFlag(isShowingProgressKey)
	isCompleted := newAndOldIssues.NewIssue.Completed == 1 && newAndOldIssues.NewIssue.PartnerVerifiedCompletion

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

func viewObjectiveReadonly(w workflowImpl, tc IssueTypeClass, newAndOldIssues NewAndOldIssues, isShowingProgress bool) platform.MessageContent {
	fields := tc.View(w, false, isShowingProgress, newAndOldIssues)

	return platform.MessageContent{
		Message: "",
		Attachments: []ebm.Attachment{
			{
				Fields: omitEmpty(fields),
			},
		},
	}
}
