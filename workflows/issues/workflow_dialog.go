package issues

import (
	// "github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	// daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	issuesUtils "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	// eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	ex "github.com/adaptiveteam/adaptive/workflows/exchange"
	"github.com/pkg/errors"
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
	err = errors.Wrapf(err, "OnEdit/issuesUtils.Read(issueID=%s)", issue.GetIssueID())
	if err != nil {
		return
	}
	out, err = w.showDialog(ctx, issue, UpdateContext)
	out.Interaction.KeepOriginal = true
	return
}

func (w workflowImpl) getFromContext(ctx wf.EventHandlingContext) (newAndOldIssues NewAndOldIssues, err error) {
	itype := IssueType(ctx.Data[issueTypeKey])
	tc := getTypeClass(itype)
	issueID, updated := ctx.Data[issueIDKey]
	var oldIssue Issue
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
	err = w.prefetch(ctx, &newIssue)
	if err == nil {
		if updated {
			err = w.prefetch(ctx, &oldIssue)
		} else {
			oldIssue = newIssue
		}

		newAndOldIssues = NewAndOldIssues{
			NewIssue: newIssue,
			OldIssue: oldIssue,
			Updated:  updated,
		}
	}
	return
}

func (w workflowImpl) OnDialogSubmitted(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	defer w.recoverToErrorVar("workflowImpl.OnDialogSubmitted", &err) // because there are panics downstream
	issueID, _ := ctx.Data[issueIDKey]
	log := w.AdaptiveLogger.WithField("issueID", issueID).WithField("Handler", "OnDialogSubmitted")
	log.Info("Start")

	var newAndOldIssues NewAndOldIssues
	newAndOldIssues, err = w.getFromContext(ctx)
	if err != nil {
		return
	}
	newAP := newAndOldIssues.NewIssue.UserObjective.AccountabilityPartner
	oldAP := newAndOldIssues.OldIssue.UserObjective.AccountabilityPartner
	isCoachRequestNeeded :=
		!utilsUser.IsSpecialOrEmptyUserID(newAP) &&
			!newAndOldIssues.Updated ||
			(newAP != oldAP)

	var postponedEvents []wf.PostponeEventForAnotherUser
	if isCoachRequestNeeded {
		postponedEvents = w.requestCoach(ctx, newAndOldIssues)
	}
	var responses []platform.Response
	newAdvocate := newAndOldIssues.NewIssue.UserObjective.UserID
	oldAdvocate := newAndOldIssues.OldIssue.UserObjective.UserID
	isCoacheeRequestNeeded := 
		!utilsUser.IsSpecialOrEmptyUserID(newAdvocate) &&
			!newAndOldIssues.Updated ||
			(newAdvocate != oldAdvocate)
	if isCoacheeRequestNeeded {
		postponedEvents = append(postponedEvents, w.requestCoachee(ctx, newAndOldIssues)...)
	}
	shouldNotifyOldCoach := !utilsUser.IsSpecialOrEmptyUserID(oldAP) && newAndOldIssues.Updated && newAP != oldAP
	if shouldNotifyOldCoach {
		responses = append(responses, 
			platform.Post(platform.ConversationID(oldAP),
				platform.MessageContent{
					Message: ui.Sprintf("<@%s> has requested a different accountability partner for the %s:\n%s\n%s",
						newAndOldIssues.NewIssue.UserObjective.UserID, 
						newAndOldIssues.NewIssue.GetIssueType().Template(),
						newAndOldIssues.NewIssue.UserObjective.Name,
						newAndOldIssues.NewIssue.UserObjective.Description,
					),
				},
			),
		)
	}
	w.AdaptiveLogger.Infof("OnDialogSubmitted: Saving %v\n", newAndOldIssues.NewIssue)
	err = issuesUtils.Save(newAndOldIssues.NewIssue)(w.DynamoDBConnection)
	err = errors.Wrapf(err, "OnDialogSubmitted: Saving")
	if err != nil {
		return
	}
	issueID = newAndOldIssues.NewIssue.GetIssueID()
	ctx.Data[issueIDKey] = issueID
	if newAndOldIssues.NewIssue.UserObjective.ID == "" && !newAndOldIssues.Updated {
		w.AdaptiveLogger.Warnf("INVALID(2): issueID is empty %v\n", newAndOldIssues.NewIssue)
	}
	itype := IssueType(ctx.Data[issueTypeKey])
	tc := getTypeClass(itype)
	if err == nil {
		dialogSituationID := DialogSituationIDWithoutIssueType(ctx.Data[dialogSituationIDKey])
		eventDescription := ObjectiveCreatedUpdatedStatusTemplate(itype, newAndOldIssues.Updated, ctx.Request.User.ID)
		out, err = w.onNewOrUpdatedItemAvailable(ctx, tc, newAndOldIssues, dialogSituationID, eventDescription, false)
	} else {
		log.WithError(err).Error("OnDialogSubmitted: Couldn't create an " + ui.RichText(itype.Template()))
		out = ctx.Reply("Couldn't create an " + ui.RichText(itype.Template()))
		err = nil // we want to show error interaction and we have logged the error
	}
	out.PostponedEvents = postponedEvents
	out.Responses = append(out.Responses, responses...)
	if utilsUser.UserID_Requested == newAP {
		
		// responses = w.requestCoachViaCoachingCommunity(ctx, newAndOldIssues)
		out.ImmediateEvents = append(out.ImmediateEvents, 
			ex.RequestCoachViaCommunity(itype, issueID),
		)
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
	return platform.ConversationID(comm.ChannelID), err2
}

func (w workflowImpl) requestCoach(ctx wf.EventHandlingContext, newAndOldIssues NewAndOldIssues) (postponedEvents []wf.PostponeEventForAnotherUser) {
	newAP := newAndOldIssues.NewIssue.UserObjective.AccountabilityPartner
	if !utilsUser.IsSpecialOrEmptyUserID(newAP) && 
	newAP != ctx.Request.User.ID {// it's a different user
		postponedEvents = []wf.PostponeEventForAnotherUser{
			ex.RequestCoach(
				newAndOldIssues.NewIssue.GetIssueType(),
				newAndOldIssues.NewIssue.UserObjective.ID,
				newAndOldIssues.NewIssue.UserObjective.AccountabilityPartner,
			),
		}
	}
	return
}

func (w workflowImpl) requestCoachee(ctx wf.EventHandlingContext, newAndOldIssues NewAndOldIssues) (postponedEvents []wf.PostponeEventForAnotherUser) {
	newAdvocate := newAndOldIssues.NewIssue.UserObjective.UserID
	if !utilsUser.IsSpecialOrEmptyUserID(newAdvocate) && 
	newAdvocate != ctx.Request.User.ID {// it's a different user
		postponedEvents = []wf.PostponeEventForAnotherUser{
			ex.RequestCoachee(
				newAndOldIssues.NewIssue.GetIssueType(),
				newAndOldIssues.NewIssue.UserObjective.ID,
				newAdvocate,
			),
		}
	}
	return
}
