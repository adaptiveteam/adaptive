package exchange

import (
	"time"

	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
)

// DefaultCoachRequestValidityDuration -
const DefaultCoachRequestValidityDuration = 14 * 24 * time.Hour

var RequestCoachPath = CommunityPath.Append(RequestCoachNamespace)

const IssueUpdatedEvent wf.Event = "updated"

// RequestCoach constructs a request coach postponed event.
func RequestCoach(issueType IssueType, issueID string, coachID string) wf.PostponeEventForAnotherUser {
	actionPath := wf.ExternalActionPathWithData(
		RequestCoachPath,
		"init",
		"",
		map[string]string{
			IssueIDKey:   issueID,
			IssueTypeKey: string(issueType),
		},
		false, // IsOriginalPermanent
	)
	return wf.PostponeEventForAnotherUser{
		UserID:       coachID,
		ActionPath:   actionPath,
		ValidThrough: time.Now().Add(DefaultCoachRequestValidityDuration),
	}
}

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
