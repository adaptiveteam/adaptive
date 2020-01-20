package issues

import (
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/pkg/errors"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	issuesUtils "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
)

func (w workflowImpl) showDialog(ctx wf.EventHandlingContext, issue Issue, dialogSituationID DialogSituationIDWithoutIssueType) (out wf.EventOutput, err error) {
	ctx.Data[dialogSituationIDKey] = string(dialogSituationID)

	itype := IssueType(ctx.Data[issueTypeKey])
	tc := getTypeClass(itype)
	var survey ebm.AttachmentActionSurvey
	survey, err = tc.CreateDialog(w, ctx, issue)
	out.Interaction = wf.OpenSurvey(survey)
	return
}

func (w workflowImpl) OnEdit(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	defer w.recoverToErrorVar("workflowImpl.OnEdit", &err) // because there are panics downstream
	issueID := ctx.Data[issueIDKey]
	itype := IssueType(ctx.Data[issueTypeKey])
	w.AdaptiveLogger.WithField("issueID", issueID).Info("OnEdit")
	var issue Issue
	issue, err = issuesUtils.Read(itype, issueID)(w.DynamoDBConnection)
	out, err = w.showDialog(ctx, issue, UpdateContext)
	out.Interaction.KeepOriginal = true
	return
}

func (w workflowImpl) OnDialogSubmitted(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	defer w.recoverToErrorVar("workflowImpl.OnDialogSubmitted", &err) // because there are panics downstream
	// issueID := ctx.Data[issueIDKey]
	itype := IssueType(ctx.Data[issueTypeKey])
	tc := getTypeClass(itype)
	var oldIssue Issue
	issueID, updated := ctx.Data[issueIDKey]
	w.AdaptiveLogger.WithField("issueID", issueID).Info("OnDialogSubmitted")
	if updated {
		oldIssue, err = issuesUtils.Read(itype, issueID)(w.DynamoDBConnection)
		if err != nil {
			err = errors.Wrapf(err, "Error reading issue %s:%s", itype, issueID)
			return
		}
	} else {
		oldIssue = tc.Empty()
	}
	newIssue := tc.ExtractFromContext(ctx, issueID, updated, oldIssue)
	(&newIssue).NormalizeIssueDateTimes()
	//issuesUtils.NormalizeIssueDateTimes(&newIssue)
	w.prefetch(ctx, &newIssue)
	if updated {
		w.prefetch(ctx, &oldIssue)
	} else {
		oldIssue = newIssue
	}

	newAndOldIssues := NewAndOldIssues{
		NewIssue: newIssue,
		OldIssue: oldIssue,
		Updated:  updated,
	}
	w.AdaptiveLogger.Infof("OnDialogSubmitted: Saving %v\n", newIssue)
	err = issuesUtils.Save(newIssue)(w.DynamoDBConnection)
	if err != nil {
		err = errors.Wrapf(err, "OnDialogSubmitted: Saving")
		return
	}
	ctx.Data[issueIDKey] = newIssue.UserObjective.ID
	if newIssue.UserObjective.ID == "" && !updated {
		w.AdaptiveLogger.Warnf("INVALID(2): issueID is empty %v\n", newIssue)
	}
	if err == nil {
		dialogSituationID := DialogSituationIDWithoutIssueType(ctx.Data[dialogSituationIDKey])
		eventDescription := ObjectiveCreatedUpdatedStatusTemplate(updated, ctx.Request.User.ID)
		out, err = w.onNewOrUpdatedItemAvailable(ctx, tc, newAndOldIssues, dialogSituationID, eventDescription, false)
	} else {
		w.AdaptiveLogger.WithError(err).Error("OnDialogSubmitted: Couldn't create an "+ui.RichText(itype.Template()))
		out = ctx.Reply("Couldn't create an " + ui.RichText(itype.Template()))
		err = nil // we want to show error interaction and we have logged the error
	}

	err = errors.Wrap(err, "{OnDialogSubmitted}")
	return
}

// func CreateObjectiveWorkflow_OnDialogSubmitted(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
// 	logger.Infof("CreateObjectiveWorkflow_OnDialogSubmitted")
// 	reply := simpleReply(ctx)
// 	capCommID := ctx.Data[capCommIDKey]
// 	var item models.StrategyObjective
// 	var oldItem models.StrategyObjective
// 	var updated bool
// 	ExtractFromContext(ctx wf.EventHandlingContext, oldIssue Issue) (newIssue Issue, updated bool) {
// 	item, oldItem, updated, err = getItemAndOldItem(ctx)
// 	err = saveItem(ctx.PlatformID, item, capCommID)
// 	if err == nil {
// 		out = onNewItemAvailable(ctx, item, oldItem, updated, capCommID)
// 	} else {
// 		logger.WithField("error", err).Errorf("CreateObjectiveWorkflow_OnDialogSubmitted error: %+v", err)
// 		out = ctx.Reply("Couldn't create an objective")
// 		err = nil // we want to show error interaction and we have logged the error
// 	}
// 	out.KeepOriginal = true
// 	out.RuntimeData = runtimeData(oldItem) // keeping the old item so that we'll be able to show it again after analysis.

// 	return
// }

func (w workflowImpl) OnDialogCancelled(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	issueID := ctx.Data[issueIDKey]
	w.AdaptiveLogger.WithField("issueID", issueID).Info("OnDialogCancelled")
	return
}

// func extractTypedObjectiveFromContext(ctx wf.EventHandlingContext) (item models.StrategyObjective, updated bool, err error) {
// 	item.ID, updated = ctx.Data[itemIDKey]
// 	item.Name = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveName]
// 	item.Type = models.StrategyObjectiveType(ctx.Request.DialogSubmissionCallback.Submission[SObjectiveType])
// 	item.Description = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveDescription]
// 	item.AsMeasuredBy = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveMeasures]
// 	item.Targets = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveTargets]
// 	item.Advocate = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveAdvocate]
// 	item.ExpectedEndDate = ctx.Request.DialogSubmissionCallback.Submission[SObjectiveEndDate]

// 	item.PlatformID = string(ctx.PlatformID)
// 	return
// }

func findStrategyCommunityConversation(w workflowImpl, ctx wf.EventHandlingContext) (platform.ConversationID, error) {
	comm, err2 := AdaptiveCommunityReadByID(community.Strategy)(w.DynamoDBConnection)
	err2 = errors.Wrap(err2, "{findStrategyCommunityConversation}")
	return platform.ConversationID(comm.Channel), err2
}
