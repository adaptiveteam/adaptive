package exchange

import (
	"strconv"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
)

const IssuesNamespace = "issues"
var IssuesPath = CommunityPath.Append(IssuesNamespace)
const InitState wf.State = "init"
// EventByType constructs an event name by issue type
func EventByType(name string, itype IssueType) wf.Event {
	return wf.Event(name + "(" + string(itype) + ")")
}

// PromptStaleIssuesEvent starts a workflow
var PromptStaleIssuesEvent = func (issueType IssueType) wf.Event { return EventByType("prompt_stale_issues", issueType) }
// PromptStaleIssues shows a prompt if user wants to see the list of stale 
// issues of the given type
func PromptStaleIssues(userID string, issueType IssueType, days int) wf.TriggerImmediateEventForAnotherUser {
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
		UserID: userID,
		ActionPath: actionPath,
	}
}