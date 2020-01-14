package fetch_dialog

import (
	"github.com/adaptiveteam/adaptive/daos/contextAliasEntry"
	"github.com/adaptiveteam/adaptive/daos/dialogEntry"
)

// NewDialogEntry creates a new structure properly constructed
func NewDialogEntry(
	context string,
	subject string,
	updated string,
	dialog []string,
	comments []string,
	dialogID string,
	learnMoreLink string,
	learnMoreContent string,
	buildBranch string,
	cultivationBranch string,
	masterBranch string,
	buildID string,
) (rv DialogEntry) {
	if context == "" ||
		subject == "" ||
		updated == "" ||
		dialog == nil ||
		dialogID == "" ||
		buildBranch == "" ||
		cultivationBranch == "" ||
		masterBranch == "" {
		panic("cannot have empty initialization values")
	}
	rv.Context = context
	rv.Subject = subject
	rv.Updated = updated
	rv.Dialog = dialog
	rv.Comments = comments
	rv.DialogID = dialogID
	rv.LearnMoreLink = learnMoreLink
	rv.LearnMoreContent = learnMoreContent
	rv.BuildBranch = buildBranch
	rv.CultivationBranch = cultivationBranch
	rv.MasterBranch = masterBranch
	rv.BuildID = buildID

	return rv
}

type DialogEntry = dialogEntry.DialogEntry
// // DialogEntry stores all of the  relevant information for a piece of dialog including:
// // Context          - This is the context path for the piece of dialog
// // Subject          - This is the dialog subject
// // Updated          - This was when the dialog was last updated
// // Dialog           - These are the dialog options
// // Comments         - Comments to help cultivators understand the dialog intent
// // DialogID         - This is an immutable UUID that developers can use
// // LearnMoreLink    - This the link to the LearnMore page
// // LearnMoreContent - This is the actual content from the LearnMore page
// type DialogEntry struct {
// 	Context           string   `json:"context"`
// 	Subject           string   `json:"subject"`
// 	Updated           string   `json:"updated"`
// 	Dialog            []string `json:"dialog"`
// 	Comments          []string `json:"comments"`
// 	DialogID          string   `json:"dialog_id"`
// 	LearnMoreLink     string   `json:"learn_more_link"`
// 	LearnMoreContent  string   `json:"learn_more_content"`
// 	BuildBranch       string   `json:"build_branch"`
// 	CultivationBranch string   `json:"cultivation_branch"`
// 	MasterBranch      string   `json:"master_branch"`
// 	BuildID           string   `json:"build_id"`
// }
type ContextAliasEntry = contextAliasEntry.ContextAliasEntry
// // ContextAliasEntry contains all of the information needed for a context alias
// // A context alias is a way to alias  a piece of context without spelling out
// // the context path.  If the path changes you can still safely use the alias.
// type ContextAliasEntry struct {
// 	Alias           string   `json:"application_alias"`
// 	Context         string   `json:"context"`
// 	BuildID         string   `json:"build_id"`
// }
