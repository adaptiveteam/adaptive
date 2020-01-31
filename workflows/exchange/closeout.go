package exchange

import (
	"time"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
)
// RequestCloseoutNamespace -
const RequestCloseoutNamespace = "request_closeout"
var RequestCloseoutPath = CommunityPath.Append(RequestCloseoutNamespace)

// RequestCloseoutForIssue constructs a postponed event to request closeout
func RequestCloseoutForIssue(issue Issue) wf.PostponeEventForAnotherUser {
	return RequestCloseout(issue.GetIssueType(), issue.GetIssueID(), issue.UserObjective.AccountabilityPartner)
}

// RequestCloseout constructs a request to coach about closeout of the issue.
// A postponed event is used to communicate the request.
func RequestCloseout(issueType IssueType, issueID string, coachID string) wf.PostponeEventForAnotherUser {
	actionPath := wf.ExternalActionPathWithData(
		RequestCloseoutPath,
		"init", 
		"",
		map[string]string{
			IssueIDKey: issueID,
			IssueTypeKey: string(issueType),
		},
		false, // IsOriginalPermanent
	)
	return wf.PostponeEventForAnotherUser{
		UserID: coachID,
		ActionPath: actionPath,
		ValidThrough: time.Now().Add(DefaultCoachRequestValidityDuration),
	}
}
