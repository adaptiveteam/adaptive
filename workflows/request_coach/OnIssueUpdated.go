package request_coach

import (
	"time"

	common "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	engCommon "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	engIssues "github.com/adaptiveteam/adaptive/adaptive-engagements/issues"
	issues "github.com/adaptiveteam/adaptive/adaptive-engagements/issues"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	utilsIssues "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	wfCommon "github.com/adaptiveteam/adaptive/workflows/common"
	"github.com/pkg/errors"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/adaptiveteam/adaptive/workflows/exchange"
)

// OnIssueUpdated - an event from issue owner when something is changed.
// `data` will contain `exchange.DialogSituationKey` with one of the dialog situations.
func (w workflowImpl) OnIssueUpdated() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		dialogSituation := DialogSituationIDWithoutIssueType(ctx.Data[exchange.DialogSituationKey])

		isShowingProgress := dialogSituation == utilsIssues.ProgressUpdateContext
		ctx.SetFlag(exchange.IsShowingProgressKey, isShowingProgress)
		out, err = w.standardView(ctx)
		out = out.WithNextState(UpdateShownState)
		return
	}
}

func (w workflowImpl) OnProvideFeedback() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[IssueIDKey]
		issueType := issues.IssueType(ctx.Data[IssueTypeKey])
		dialogSituation := DialogSituationIDWithoutIssueType(ctx.Data[exchange.DialogSituationKey])
		log := w.AdaptiveLogger.
			WithField("issueID", issueID).
			WithField("issueType", issueType).
			WithField("dialogSituation", dialogSituation).
			WithField("Handler", "OnProvideFeedback")
		log.Info("Start")
		ctx.SetFlag(exchange.IsShowingProgressKey, true) // enable show progress. This will make sure that progress is prefetched
		var newAndOldIssues issues.NewAndOldIssues
		newAndOldIssues, err = wfCommon.WorkflowContext(w).GetNewAndOldIssues(ctx)
		if err != nil {
			return
		}
		issue := newAndOldIssues.NewIssue
		var oldPartnerReportedProgress string
		var oldPartnerComments ui.PlainText
		if len(issue.Progress) > 0 {
			p := newAndOldIssues.NewIssue.Progress[0]
			oldPartnerReportedProgress = p.PartnerReportedProgress
			oldPartnerComments = ui.PlainText(p.PartnerComments)
		}
		switch dialogSituation {
		case utilsIssues.ProgressUpdateContext:
			out = out.WithSurvey(utils.AttachmentSurvey(
				string("Feedback on the changes"),
				progressCommentSurveyElements(
					ui.PlainText(issue.UserObjective.Name),
					issue.UserObjective.CreatedDate,
					oldPartnerReportedProgress,
					oldPartnerComments,
				),
			))
		default:
			// no text
		}
		return
	}
}

func (w workflowImpl) OnDismiss() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[IssueIDKey]
		issueType := issues.IssueType(ctx.Data[IssueTypeKey])
		dialogSituation := DialogSituationIDWithoutIssueType(ctx.Data[exchange.DialogSituationKey])
		log := w.AdaptiveLogger.
			WithField("issueID", issueID).
			WithField("issueType", issueType).
			WithField("dialogSituation", dialogSituation).
			WithField("Handler", "OnDismiss")
		log.Info("Dismiss")
		return
	}
}

func progressCommentSurveyElements(objName ui.PlainText, startDate string,
	oldPartnerReportedProgress string,
	oldPartnerComments ui.PlainText,
) []ebm.AttachmentActionTextElement {
	nameConstrained := engIssues.ObjectiveCommentsTitle(objName)
	today := core.ISODateLayout.Format(time.Now())
	elapsedDays := common.DurationDays(startDate, today, core.ISODateLayout, "progressCommentSurveyElements")
	return []ebm.AttachmentActionTextElement{
		{
			Label:    string(engIssues.ObjectiveStatusLabel(elapsedDays, startDate)),
			Name:     ObjectiveStatusColor,
			ElemType: models.MenuSelectType,
			Options:  utils.AttachActionElementOptions(models.ObjectiveStatusColorKeyValues),
			Value:    oldPartnerReportedProgress,
		},
		ebm.NewTextArea(ObjectiveProgressComments, nameConstrained, ObjectiveProgressCommentsPlaceholder, oldPartnerComments),
	}

}

const (
	ObjectiveProgressComments                         = "objective_progress"
	ObjectiveProgressCommentsPlaceholder ui.PlainText = ebm.EmptyPlaceholder
)

const (
	ObjectiveStatusColor       = "objective_status_color"
	ObjectiveCloseoutComment   = "objective_closeout_comment"
	ObjectiveNoCloseoutComment = "objective_no_closeout_comment"
)

// OnCommentsSubmitted starts the following sequence of steps
// - start analysis in go-routine;
// - show message and request immediate event when message id is ready;
// - when message id is received (OnNewOrUpdatedCoachCommentAvailableOnMessageIDAvailableEvent)
//   it waits for analysis results at most 2 seconds;
// - overrides message with the colored message and puts analysis into thread.
// So two processes runs asynchronously in parallel - publishing to Slack and analysis.
func (w workflowImpl) OnCommentsSubmitted() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := exchange.GetIssueID(ctx)
		issueType := exchange.GetIssueType(ctx)
		dialogSituation := DialogSituationIDWithoutIssueType(ctx.Data[exchange.DialogSituationKey])
		log := w.AdaptiveLogger.
			WithField("issueID", issueID).
			WithField("issueType", issueType).
			WithField("dialogSituation", dialogSituation).
			WithField("Handler", "OnCommentsSubmitted")
		log.Info("Start")
		ctx.SetFlag(exchange.IsShowingProgressKey, true) // enable show progress. This will make sure that progress is prefetched
		var newAndOldIssues issues.NewAndOldIssues
		newAndOldIssues, err = wfCommon.WorkflowContext(w).GetNewAndOldIssues(ctx)
		if err != nil {
			return
		}
		if len(newAndOldIssues.NewIssue.Progress) > 0 {
			p := newAndOldIssues.NewIssue.Progress[0]
			p.PartnerID = ctx.Request.User.ID
			p.PartnerReportedProgress = ctx.Request.Submission[ObjectiveStatusColor]
			p.PartnerComments = ctx.Request.Submission[ObjectiveProgressComments]
			err = userObjectiveProgress.CreateOrUpdate(p)(w.DynamoDBConnection)
			if err != nil {
				return
			}
			out, err = w.onNewOrUpdatedCoachCommentAvailable(ctx, newAndOldIssues)
		}
		return
	}
}

func (w workflowImpl) OnDetails(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	ctx.ToggleFlag(exchange.IsShowingDetailsKey)
	return w.standardView(ctx)
}

func (w workflowImpl) OnProgressShow(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	ctx.ToggleFlag(exchange.IsShowingProgressKey)
	return w.standardView(ctx)
}

func (w workflowImpl) standardView(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	w.AdaptiveLogger.Info("standardView")

	issueID := ctx.Data[IssueIDKey]
	issueType := issues.IssueType(ctx.Data[IssueTypeKey])
	dialogSituation := DialogSituationIDWithoutIssueType(ctx.Data[exchange.DialogSituationKey])

	log := w.AdaptiveLogger.
		WithField("issueID", issueID).
		WithField("issueType", issueType).
		WithField("dialogSituation", dialogSituation).
		WithField("Handler", "standardView")
	log.Info("Start")
	var newAndOldIssues issues.NewAndOldIssues
	newAndOldIssues, err = wfCommon.WorkflowContext(w).GetNewAndOldIssues(ctx)
	if err != nil {
		return
	}
	issue := newAndOldIssues.NewIssue
	var notificationText ui.RichText
	switch dialogSituation {
	case utilsIssues.UpdateContext:
		notificationText = ui.Sprintf("%s has updated the below %s. "+
			"You might want to provide some valuable feedback on this update.",
			engCommon.TaggedUser(issue.UserObjective.UserID),
			issue.GetIssueType().Template())
	// case utilsIssues.ProgressUpdateContext:
	// 	notificationText = ui.Sprintf("%s has updated progress on the below %s. "+
	// 		"You might want to provide some valuable feedback on this update.",
	// 		engCommon.TaggedUser(issue.UserObjective.UserID),
	// 		issue.GetIssueType().Template())
	default:
		// no text
	}
	if notificationText != "" {
		viewState := issues.ViewState{
			IsShowingDetails:  exchange.IsShowingDetails(ctx),
			IsShowingProgress: exchange.IsShowingProgress(ctx),
			IsWritable:        false,
		}

		view := engIssues.GetInteractiveMessage(newAndOldIssues, viewState)
		view.AttachmentText = notificationText
		view.OverrideOriginal = true

		var feedbackButtonCaption ui.PlainText
		progresses := newAndOldIssues.NewIssue.PrefetchedData.Progress
		if len(progresses) == 0 || progresses[0].PartnerComments == "" {
			feedbackButtonCaption = "Provide feedback"
		} else {
			feedbackButtonCaption = "Update feedback"
		}
		view.InteractiveElements = append(
			view.InteractiveElements,
			wf.Button(ConfirmedEvent, feedbackButtonCaption),
			// wf.Button(RejectedEvent, "Dismiss"),
		)
		out = out.WithInteractiveMessage(view)
	}

	err = errors.Wrap(err, "{standardView}")
	return
}
