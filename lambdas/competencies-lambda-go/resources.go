package competencies

import (
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	// "github.com/adaptiveteam/adaptive/engagement-builder/say"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"strings"
)

// Subject is the typeof the key for accessing templates
type Subject string

const (
	// AdaptiveValuesContext is the key to fetch dialog messages
	AdaptiveValuesContext = "engagements/adaptiveValue"

	InviteToListOfAdaptiveValuesSubject Subject = "InviteToListOfAdaptiveValues"
	AdaptiveValuesListTitleSubject      Subject = "AdaptiveValuesListTitle"
	AdaptiveValueItemSubject            Subject = "AdaptiveValueItem"
	// NewAdaptiveValueNoticeSubject       Subject = "NewAdaptiveValueNotice"
	CreatedAdaptiveValueNoticeSubject Subject = "CreatedAdaptiveValueNotice"
	EditButtonSubject                 Subject = "EditButton"
	// EditingAdaptiveValueNoticeSubject   Subject = "EditingAdaptiveValueNotice"
	UpdatedAdaptiveValueNoticeSubject Subject = "UpdatedAdaptiveValueNotice"
	DeleteButtonSubject               Subject = "DeleteButton"
	DeletedAdaptiveValueNoticeSubject Subject = "DeletedAdaptiveValueNotice"
	CancelDeleteSubject               Subject = "CancelDeleteAction"
	CreateButtonSubject               Subject = "CreateButton"
	ModifyListButtonSubject           Subject = "ModifyListButton"
)

const (
	inviteToListOfAdaptiveValuesTemplate = "Please find the list of Competencies in the thread below :arrow_down:"
	adaptiveValuesListTitleTemplate      = "Competencies"
	adaptiveValueItemTemplate            = "• *{AdaptiveValue.Name}*\n\t{AdaptiveValue.Description}\n"
	// newAdaptiveValueNoticeTemplate       = "Creating a new adaptive value..."
	createdAdaptiveValueNoticeTemplate = "Competency has been created"
	editButtonLabel                    = "Edit" // "🖉"
	// editingAdaptiveValueNoticeTemplate   = "Updating the adaptive value..."
	updatedAdaptiveValueNoticeTemplate = "Competency has been updated"
	deleteButtonLabel                  = "Close" // "🗑"
	deletedAdaptiveValueNoticeTemplate = "Closed the Competency"
	cancelDeleteText                   = "Competency close has been cancelled"
	createButtonLabel                  = "Add another"
	modifyListButtonLabel              = "Modify/Close one"

	valueSurveyLabel = "Competency"
)

var (
	templates = map[Subject]string{
		InviteToListOfAdaptiveValuesSubject: inviteToListOfAdaptiveValuesTemplate,
		AdaptiveValuesListTitleSubject:      adaptiveValuesListTitleTemplate,
		AdaptiveValueItemSubject:            adaptiveValueItemTemplate,
		// NewAdaptiveValueNoticeSubject:       newAdaptiveValueNoticeTemplate,
		CreatedAdaptiveValueNoticeSubject: createdAdaptiveValueNoticeTemplate,
		EditButtonSubject:                 editButtonLabel,
		// EditingAdaptiveValueNoticeSubject:   editingAdaptiveValueNoticeTemplate,
		UpdatedAdaptiveValueNoticeSubject: updatedAdaptiveValueNoticeTemplate,
		DeleteButtonSubject:               deleteButtonLabel,
		DeletedAdaptiveValueNoticeSubject: deletedAdaptiveValueNoticeTemplate,
		CancelDeleteSubject:               cancelDeleteText,
		CreateButtonSubject:               createButtonLabel,
		ModifyListButtonSubject:           modifyListButtonLabel,
	}
)

// RetrieveTemplate returns a dialog message template for
// a given context and subject
func RetrieveTemplate(context string, subject Subject) string {
	if context == AdaptiveValuesContext {
		return templates[subject]
	}
	return "No template for context " + context
}

// AdaptiveValuesTemplate returns a dialog message template for
// a given subject
func AdaptiveValuesTemplate(subject Subject) string {
	return RetrieveTemplate(AdaptiveValuesContext, subject)
}

// RenderAdaptiveValueItem returns a formatted string based on the provided data
func RenderAdaptiveValueItem(adaptiveValue models.AdaptiveValue) string {
	t := RetrieveTemplate(AdaptiveValuesContext, AdaptiveValueItemSubject)
	t = strings.Replace(t, "{AdaptiveValue.Name}", adaptiveValue.Name, -1)
	t = strings.Replace(t, "{AdaptiveValue.ValueType}", adaptiveValue.ValueType, -1)
	t = strings.Replace(t, "{AdaptiveValue.Description}", adaptiveValue.Description, -1)
	return t
}

// field creates a field view from provided label and value.
func field(label ui.PlainText, value ui.RichText) ebm.AttachmentField {
	return ebm.AttachmentField{
		Title: string(label),
		Value: string(value),
		Short: false,
	}
}

// ShowAdaptiveValueItemAsFields converts adaptive value.
// Don't forget to enable markdown everywhere!
// MarkDownIn([]ebm.MarkdownField{ebm.MarkdownFieldText, ebm.MarkdownFieldPretext}).
func ShowAdaptiveValueItemAsFields(adaptiveValue models.AdaptiveValue) []ebm.AttachmentField {
	return []ebm.AttachmentField{
		field("Name", ui.RichText(adaptiveValue.Name).Bold()),
		field("Type", ui.RichText(adaptiveValue.ValueType)),
		field("Description", ui.RichText(adaptiveValue.Description)),
	}
}
