package issues

import (
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

type IssueProperty = func(issue Issue) ui.PlainText

func getIssueTypeFromContext(ctx wf.EventHandlingContext) IssueType {
	return IssueType(ctx.Data[issueTypeKey])
}


