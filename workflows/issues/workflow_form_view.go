package issues

import (
	"github.com/adaptiveteam/adaptive/workflows/common"
	issues "github.com/adaptiveteam/adaptive/adaptive-engagements/issues"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
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
	newAndOldIssues, err = w.WorkflowContext.GetNewAndOldIssues(*ctx) // userObjectiveByID(itemID)
	err = errors.Wrapf(err, "textAnalysisInput/WorkflowContext.GetNewAndOldIssues")
	if err != nil {
		return
	}
	ctx.RuntimeData = runtimeData(newAndOldIssues)

	itype := newAndOldIssues.NewIssue.GetIssueType()

	if itype == "" {
		err = errors.New("textAnalysisInput/itype is empty")
		return
	}
	context := common.GetDialogContext(dialogSituationID, itype)

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
		if err == nil {
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
		}
		out.NextState = "done"
		out.KeepOriginal = true // we want to override it, so, not to delete
		return
	}
}

func (w workflowImpl) standardView(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	w.AdaptiveLogger.Info("standardView")
	isShowingProgress := ctx.GetFlag(isShowingProgressKey)
	w.AdaptiveLogger.Infof("standardView: isShowingProgress=%v", isShowingProgress)
	var newAndOldIssues NewAndOldIssues
	newAndOldIssues, err = w.WorkflowContext.GetNewAndOldIssues(ctx)
	if err != nil {
		return
	}
	viewState := issues.ViewState{
		IsShowingDetails:  ctx.GetFlag(isShowingDetailsKey),
		IsShowingProgress: isShowingProgress,
		IsWritable:        true,
	}
	view := issues.GetInteractiveMessage(newAndOldIssues, viewState)

	view.OverrideOriginal = true
	out.Messages = []wf.InteractiveMessage{view}
	err = errors.Wrap(err, "{standardView}")
	return
}

func (w workflowImpl) OnDetails(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	ctx.ToggleFlag(isShowingDetailsKey)
	return w.standardView(ctx)
}

func (w workflowImpl) OnProgressShow(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	oldFlag := ctx.GetFlag(isShowingProgressKey)
	ctx.ToggleFlag(isShowingProgressKey)
	newFlag := ctx.GetFlag(isShowingProgressKey)
	w.AdaptiveLogger.Infof("OnProgressShow: flag %v to %v", oldFlag, newFlag)
	return w.standardView(ctx)
}
