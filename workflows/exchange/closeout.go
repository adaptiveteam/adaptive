package exchange

import (
	"time"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
)
// RequestCloseoutNamespace -
const RequestCloseoutNamespace = "request_closeout"

// RequestCloseoutForIssue constructs a postponed event to request closeout
func RequestCloseoutForIssue(issue Issue) wf.PostponeEventForAnotherUser {
	return RequestCloseout(issue.GetIssueType(), issue.GetIssueID(), issue.UserObjective.AccountabilityPartner)
}

// RequestCloseout constructs a request to coach about closeout of the issue.
// A postponed event is used to communicate the request.
func RequestCloseout(issueType IssueType, issueID string, coachID string) wf.PostponeEventForAnotherUser {
	actionPath := wf.ExternalActionPathWithData(
		models.ParsePath("/community/" + RequestCloseoutNamespace), 
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
