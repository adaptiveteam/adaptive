package exchange

import (
	"time"

	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
)

// NotifyAboutUpdatesForIssue creates a postponed event
func NotifyAboutUpdatesForIssue(newAndOldIssues NewAndOldIssues, dialogSituationID DialogSituationIDWithoutIssueType) (evt wf.PostponeEventForAnotherUser) {
	evt = wf.PostponeEventForAnotherUser{
		ActionPath: wf.ExternalActionPathWithData(
			RequestCoachPath,
			"init",
			IssueUpdatedEvent,
			map[string]string{
				IssueIDKey:   newAndOldIssues.NewIssue.GetIssueID(),
				IssueTypeKey: string(newAndOldIssues.NewIssue.GetIssueType()),
				DialogSituationKey: string(dialogSituationID),
			},
			false, // IsOriginalPermanent
		),
		UserID:       newAndOldIssues.NewIssue.UserObjective.AccountabilityPartner,
		ValidThrough: time.Now().Add(DefaultCoachRequestValidityDuration),
	}
	return
}
// NotifyOwnerAboutFeedbackOnUpdatesForIssue -
func NotifyOwnerAboutFeedbackOnUpdatesForIssue(issue Issue) (evt wf.PostponeEventForAnotherUser) {
	evt = wf.PostponeEventForAnotherUser{
		ActionPath: wf.ExternalActionPathWithData(
			IssuesPath,
			"init",
			IssueFeedbackOnUpdatesEvent,
			map[string]string{
				IssueIDKey:   issue.GetIssueID(),
				IssueTypeKey: string(issue.GetIssueType()),
			},
			false, // IsOriginalPermanent
		),
		UserID:       issue.UserObjective.UserID,
		ValidThrough: time.Now().Add(DefaultCoachRequestValidityDuration),
	}
	return
}
