package exchange

import (
	"time"

	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
)

// DefaultCoachRequestValidityDuration -
const DefaultCoachRequestValidityDuration = 14 * 24 * time.Hour

var RequestCoachPath = CommunityPath.Append(RequestCoachNamespace)

const IssueUpdatedEvent wf.Event = "updated"
const IssueFeedbackOnUpdatesEvent wf.Event = "feedback"
const RequestCoacheeEvent wf.Event = "requestCoachee"

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

// RequestCoachee constructs a request coach postponed event.
func RequestCoachee(issueType IssueType, issueID string, coachID string) wf.PostponeEventForAnotherUser {
	actionPath := wf.ExternalActionPathWithData(
		RequestCoachPath,
		"init",
		RequestCoacheeEvent,
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
