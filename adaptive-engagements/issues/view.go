package issues

import (
	"github.com/adaptiveteam/adaptive/workflows/exchange"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	utilsIssues "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// Issue -
type Issue = utilsIssues.Issue
type IssueType = utilsIssues.IssueType
type NewAndOldIssues = utilsIssues.NewAndOldIssues

// ViewState the state of view
type ViewState struct {
	IsShowingDetails   bool
	IsShowingProgress  bool
	IsWritable         bool // has Edit/Cancel buttons and progress
	HasCommentsButtons bool // has Add comment, Comment latest updates buttons
}

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

// OrdinaryViewFields - 
func OrdinaryViewFields(newAndOldIssues NewAndOldIssues, viewState ViewState) (fields []ebm.AttachmentField) {
	itype := newAndOldIssues.NewIssue.GetIssueType()
	view := GetView(itype)
	fields = view.GetMainFields(newAndOldIssues)
	if viewState.IsShowingDetails {
		fields = append(fields, view.GetDetailsFields(newAndOldIssues)...)
	}
	if viewState.IsShowingProgress {
		fields = append(fields, view.GetProgressFields(newAndOldIssues)...)
	}
	return
}
const (
	// MessageIDAvailableEvent wf.Event = "MessageIDAvailableEvent"
	EditEvent wf.Event = "EditEvent"
	// AddAnotherEvent wf.Event = "AddAnotherEvent"
	DetailsEvent wf.Event = "DetailsEvent"
	CancelEvent wf.Event = "CancelEvent"
	ProgressShowEvent wf.Event = "ProgressShowEvent"
	ProgressIntermediateEvent wf.Event = "ProgressIntermediateEvent"
	ProgressCloseoutEvent wf.Event = "ProgressCloseoutEvent"
)

func caption(trueCaption ui.PlainText, falseCaption ui.PlainText) func(bool) ui.PlainText {
	return func(flag bool) (res ui.PlainText) {
		if flag {
			res = trueCaption
		} else {
			res = falseCaption
		}
		return
	}
}

func OrdinaryInteractiveElements(newAndOldIssues NewAndOldIssues, viewState ViewState) (buttons []wf.InteractiveElement) {
	isCompleted := newAndOldIssues.NewIssue.Completed == 1 && newAndOldIssues.NewIssue.PartnerVerifiedCompletion

	isWritable := !isCompleted && viewState.IsWritable
	details := wf.Button(DetailsEvent, caption("Show less", "Show more")(viewState.IsShowingDetails))
	buttons = append(buttons, details)
	progressShow := wf.MenuOption(ProgressShowEvent, caption("Hide", "Show")(viewState.IsShowingProgress))
	progressOptions := []wf.SimpleAction{progressShow}
	// addAnother := wf.Button("add-another", "Add another?")
	if isWritable {
		buttons = append(buttons, wf.Button(EditEvent, "Edit"))
		progressIntermediate := wf.MenuOption(ProgressIntermediateEvent, "Add/Update progress")
		progressCloseout := wf.MenuOption(ProgressCloseoutEvent, "Closeout")
		progressOptions = append(progressOptions, progressIntermediate, progressCloseout)
	}
	progress := wf.InlineMenu("Progress", progressOptions...)
	buttons = append(buttons, progress)
	if isWritable {
		buttons = append(buttons, wf.AckButton(CancelEvent, "Cancel"))
	}
	return
}
// GetInteractiveMessage - returns an interactive message that represents the given issue.
// View might be in different states.
func GetInteractiveMessage(newAndOldIssues NewAndOldIssues, viewState ViewState) (view wf.InteractiveMessage) {
	itype := newAndOldIssues.NewIssue.GetIssueType()

	fields := OrdinaryViewFields(newAndOldIssues, viewState)
	interactiveElements := OrdinaryInteractiveElements(newAndOldIssues, viewState)
	createdDate := core.ParseDateOrElseToday(newAndOldIssues.NewIssue.UserObjective.CreatedDate)
	
	view = wf.InteractiveMessage{
		PassiveMessage: wf.PassiveMessage{
			Fields:             ebm.OmitEmpty(fields),
			IsPermanentMessage: true, // we don't ever want to delete the form view message
			Footer:             ebm.AttachmentFooter{Text: "Created at", Timestamp: createdDate},
		},
		InteractiveElements: interactiveElements,
		DataOverride: wf.Data{
			exchange.IssueIDKey:   newAndOldIssues.NewIssue.UserObjective.ID,
			exchange.IssueTypeKey: string(itype), //probably we don't need this because it's available
		},
	}
	return
}
