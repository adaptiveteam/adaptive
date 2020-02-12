package request_coach

import (
	"time"

	common "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	engCommon "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	engIssues "github.com/adaptiveteam/adaptive/adaptive-engagements/issues"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	issues "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	// "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	// "github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/adaptiveteam/adaptive/workflows/exchange"
)

// OnIssueUpdated - an event from issue owner when something is changed.
// `data` will contain `exchange.DialogSituationKey` with one of the dialog situations.
func (w workflowImpl) OnIssueUpdated() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[IssueIDKey]
		issueType := issues.IssueType(ctx.Data[IssueTypeKey])
		dialogSituation := DialogSituationIDWithoutIssueType(ctx.Data[exchange.DialogSituationKey])
		log := w.AdaptiveLogger.
			WithField("issueID", issueID).
			WithField("issueType", issueType).
			WithField("dialogSituation", dialogSituation).
			WithField("Handler", "OnIssueUpdated")
		log.Info("Start")
		var issue issues.Issue
		issue, err = issues.Read(issueType, issueID)(w.DynamoDBConnection)
		if err != nil {
			return
		}
		var notificationText ui.RichText
		switch dialogSituation {
		case issues.UpdateContext:
			notificationText = ui.Sprintf("%s has updated the below %s. "+
				"You might want to provide some valuable feedback on this update.",
				engCommon.TaggedUser(issue.UserObjective.UserID),
				issue.GetIssueType().Template())
		case issues.ProgressUpdateContext:
			notificationText = ui.Sprintf("%s has updated progress on the below %s. "+
				"You might want to provide some valuable feedback on this update.",
				engCommon.TaggedUser(issue.UserObjective.UserID),
				issue.GetIssueType().Template())
		default:
			// no text
		}
		// TODO: show progress/show details
		if notificationText != "" {
			out = out.WithInteractiveMessage(wf.InteractiveMessage{
				PassiveMessage: wf.PassiveMessage{
					AttachmentText: notificationText,
					Fields:         shortViewFields(issue),
				},
				InteractiveElements: []wf.InteractiveElement{
					wf.Button(ConfirmedEvent, "Provide feedback"),
					wf.Button(RejectedEvent, "Dismiss"),
				},
			})
		}
		out = out.WithNextState(FormShownState)
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
		var issue issues.Issue
		issue, err = issues.Read(issueType, issueID)(w.DynamoDBConnection)
		if err != nil {
			return
		}
		switch dialogSituation {
		case issues.ProgressUpdateContext:
			survey := utils.AttachmentSurvey(
				string("Feedback on the recent changes"),
				progressCommentSurveyElements(ui.PlainText(issue.UserObjective.Name),
					issue.UserObjective.CreatedDate))

			out = out.WithSurvey(survey)
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
func progressCommentSurveyElements(objName ui.PlainText, startDate string) []ebm.AttachmentActionTextElement {
	nameConstrained := engIssues.ObjectiveCommentsTitle(objName)
	today := core.ISODateLayout.Format(time.Now())
	elapsedDays := common.DurationDays(startDate, today, core.ISODateLayout, "progressCommentSurveyElements")
	return []ebm.AttachmentActionTextElement{
		{
			Label:    string(engIssues.ObjectiveStatusLabel(elapsedDays, startDate)),
			Name:     ObjectiveStatusColor,
			ElemType: models.MenuSelectType,
			Options:  utils.AttachActionElementOptions(models.ObjectiveStatusColorKeyValues),
		},
		ebm.NewTextArea(ObjectiveProgressComments, nameConstrained, ObjectiveProgressCommentsPlaceholder, ""),
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
	ReviewUserProgressSelect   = "review_user_progress_select"
	UberCoach                  = "uber_coach"
)

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
		// var newAndOldIssues NewAndOldIssues
		ctx.SetFlag(exchange.IsShowingProgressKey, true) // enable show progress. This will make sure that progress is prefetched
		// newAndOldIssues, err = w.getNewAndOldIssues(ctx)
		// ctx.RuntimeData = runtimeData(newAndOldIssues)
		// uo := newAndOldIssues.NewIssue.UserObjective
		// var progress userObjectiveProgress.UserObjectiveProgress
		// progress, err = extractObjectiveProgressFromContext(ctx, uo)
		// if err != nil {
		// 	return
		// }
		// isProgressAvailableForToday := false
		// if len(newAndOldIssues.NewIssue.Progress) > 0 {
		// 	isProgressAvailableForToday = newAndOldIssues.NewIssue.Progress[0].CreatedOn == progress.CreatedOn
		// }
		// err = UserObjectiveProgressSave(progress)(w.DynamoDBConnection)
		// if err != nil {
		// 	return
		// }
		// err = w.prefetch(ctx, &newAndOldIssues.NewIssue)
		// if err != nil {
		// 	return
		// }
		// ctx.RuntimeData = runtimeData(newAndOldIssues)
		// // attachs := viewProgressAttachment(mc,
		// // 	ui.PlainText(Sprintf("This is your reported progress for the below %s", typLabel)),
		// // 	"",
		// // 	comments,
		// // 	statusColor, item, models.Update)
		// // publish(models.PlatformSimpleNotification{UserId: dialog.User.ID, Channel: dialog.Channel.ID, Ts: msgState.ThreadTs, Attachments: attachs})
		// itype := newAndOldIssues.NewIssue.GetIssueType()
		// tc := getTypeClass(itype)
		// if err == nil {
		// 	ctx.SetFlag(isShowingProgressKey, true)

		// 	eventDescription := ObjectiveProgressCreatedUpdatedStatusTemplate(isProgressAvailableForToday, ctx.Request.User.ID)
		// 	out, err = w.onNewOrUpdatedItemAvailable(ctx, tc, newAndOldIssues, ProgressUpdateContext,
		// 		eventDescription, true)
		// 	// if err == nil {
		// 	// 	out.ImmediateEvent = "ProgressFormShown"
		// 	// }
		// } else {
		// 	w.AdaptiveLogger.WithError(err).Error("OnProgressFormSubmitted error")
		// 	out.Interaction = wf.SimpleResponses(
		// 		platform.Post(platform.ConversationID(ctx.Request.User.ID),
		// 			platform.MessageContent{Message: ui.Sprintf("Couldn't save the entered %s progress", tc.IssueTypeName())},
		// 		),
		// 	)
		// 	err = nil // we want to show error interaction
		// }
		// out.RuntimeData = runtimeData(newAndOldIssues)
		return
	}
}
