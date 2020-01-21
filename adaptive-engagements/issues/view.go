package issues

import (
	utilsIssues "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// Issue -
type Issue = utilsIssues.Issue
type IssueType = utilsIssues.IssueType
type NewAndOldIssues = utilsIssues.NewAndOldIssues

// View contains some representations of an issue.
type View interface {
	// GetMainFields shows name and description
	GetMainFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField)
	// GetDetailsFields shows all information except name/description
	GetDetailsFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField)
	// GetProgressFields shows only progress summary
	GetProgressFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField)
	GetAlignment(issue Issue) ui.PlainText

	// GetTextView - renders the issue as a simple text. will contain name and description
	GetTextView(issue Issue) ui.RichText
}

// ViewIDO -
type ViewIDO struct{}

// ViewSObjective -
type ViewSObjective struct{}

// ViewInitiative -
type ViewInitiative struct{}

// GetView returns view for issue type
func GetView(issueType IssueType) (view View) {
	switch issueType {
	case utilsIssues.IDO:
		view = ViewIDO{}
	case utilsIssues.SObjective:
		view = ViewSObjective{}
	case utilsIssues.Initiative:
		view = ViewInitiative{}
	}
	return
}
