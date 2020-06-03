package exchange

import (
	"strconv"

	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

const IssuesNamespace = "issues"
var IssuesPath = CommunityPath.Append(IssuesNamespace)
const InitState wf.State = "init"
const MessagePostedState wf.State = "MessagePostedState"

// EventByType constructs an event name by issue type
func EventByType(name string, itype IssueType) wf.Event {
	return wf.Event(name + "(" + string(itype) + ")")
}

// PromptStaleIssuesEvent starts a workflow
var PromptStaleIssuesEvent = func (issueType IssueType) wf.Event { return EventByType("prompt_stale_issues", issueType) }
// PromptStaleIssues shows a prompt if user wants to see the list of stale 
// issues of the given type
func PromptStaleIssues(teamID models.TeamID, userID string, issueType IssueType, days int) wf.TriggerImmediateEventForAnotherUser {
	actionPath := wf.ExternalActionPathWithData(
		IssuesPath,
		InitState, 
		PromptStaleIssuesEvent(issueType),
		map[string]string{
			"days": strconv.Itoa(days),
			// IssueTypeKey: string(issueType), // duplicate information
		},
		false, // IsOriginalPermanent
	)
	return wf.TriggerImmediateEventForAnotherUser{
		TeamID:     teamID,
		UserID:     userID,
		ActionPath: actionPath,
	}
}