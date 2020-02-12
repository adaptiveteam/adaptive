package issues

import (
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	issues "github.com/adaptiveteam/adaptive/adaptive-engagements/issues"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
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
	if newAndOldIssues.NewIssue.UserObjective.ID == "" {
		w.AdaptiveLogger.Warnf("INVALID(3): issueID is empty %v\n", newAndOldIssues.NewIssue)
	}
	viewState := issues.ViewState{
		IsShowingDetails: ctx.GetFlag(isShowingDetailsKey),
		IsShowingProgress: ctx.GetFlag(isShowingProgressKey),
		IsWritable: true,
	}
	view := issues.GetInteractiveMessage(newAndOldIssues, viewState)

	view.OverrideOriginal = true
	out.Interaction = wf.Interaction{
		Messages: []wf.InteractiveMessage{view},
	}
	err = errors.Wrap(err, "{standardView}")
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
