package issues

import (
	issuesUtils "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

type Issue = issuesUtils.Issue
type IssueType = issuesUtils.IssueType
type IssueProgressID = issuesUtils.IssueProgressID
type NewAndOldIssues = issuesUtils.NewAndOldIssues
type IssuePredicate = issuesUtils.IssuePredicate

const IDO = issuesUtils.IDO
const SObjective = issuesUtils.SObjective
const Initiative = issuesUtils.Initiative

type IssueProperty = func(issue Issue) ui.PlainText

func getIssueTypeFromContext(ctx wf.EventHandlingContext) IssueType {
	return IssueType(ctx.Data[issueTypeKey])
}


