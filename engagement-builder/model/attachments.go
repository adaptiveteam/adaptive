package model

import (
	"reflect"

	"github.com/adaptiveteam/engagement-builder/ui"
)

type Attachment struct {
	Title      string `json:"title,omitempty"`
	Text       string `json:"text,omitempty"`
	Pretext    string `json:"pretext,omitempty"`
	Fallback   string `json:"fallback"`
	CallbackId string `json:"callback_id,omitempty"`
	// Identifier will act as a unique id for the collection of buttons
	Identifier string `json:"identifier"`
	// Optional: used to visually distinguish an attachment from other messages
	Color   string             `json:"color,omitempty"`
	Fields  []AttachmentField  `json:"fields,omitempty"`
	Actions []AttachmentAction `json:"actions"`
	// value is `default`
	AttachmentType string `json:"type,omitempty"`
	// Footer is optional text to help contextualize and identify an attachment
	Footer AttachmentFooter `json:"footer,omitempty"`
	// Image related fields
	TitleLink string `json:"title_link,omitempty"`
	ImageUrl  string `json:"image_url,omitempty"`
	ThumbUrl  string `json:"thumb_url,omitempty"`
	// Markdown related
	MrkdwnIn []MarkdownField `json:"mrkdwn_in,omitempty"`
	// Author related
	Author AttachmentAuthor `json:"author,omitempty"`
}

type AttachmentAuthor struct {
	Name string `json:"name"`
	Link string `json:"link,omitempty"`
	Icon string `json:"icon,omitempty"`
}

type AttachmentFooter struct {
	Text      string `json:"text"`
	Icon      string `json:"icon,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type AttachmentActionConfirm struct {
	Title       string `json:"title"`
	Text        string `json:"text"`
	OkText      string `json:"ok_text"`
	DismissText string `json:"dismiss_text"`
}

type AttachmentAction struct {
	// Name for the action
	Name string `json:"name"`
	// User-facing label for the message button or menu representing this action.
	Text string `json:"text"`
	// Provide `button` for button or select for menu
	ActionType AttachmentActionType `json:"type"`
	Url        string               `json:"url,omitempty"`
	// A string identifying this specific action
	Value   string                  `json:"value"`
	Confirm AttachmentActionConfirm `json:"confirm,omitempty"`
	// Used only with message buttons, this decorates buttons with extra visual importance
	Style           AttachmentActionStyle      `json:"style"`
	Options         []MenuOption               `json:"options,omitempty"`
	SelectedOptions MenuOption                 `json:"selected_options,omitempty"`
	OptionGroups    []MenuOptionGroup          `json:"option_groups,omitempty"`
	DataSource      AttachmentActionDataSource `json:"data_source,omitempty"`
	// Possible replies for the action. When done, one of the messages from the list will be posted back to user
	Replies []string `json:"replies,omitempty"`
}

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	// An optional flag indicating whether the value is short enough to be displayed side-by-side with other values
	Short bool `json:"short,omitempty"`
}

// AttachmentActionSurvey represents some information
// to show dialog
// depreacated. Use AttachmentActionSurvey2
type AttachmentActionSurvey struct {
	Title       string                        `json:"title,omitempty"`
	SubmitLabel string                        `json:"submit_label,omitempty"`
	Elements    []AttachmentActionTextElement `json:"elements,omitempty"`
}

// AttachmentActionSurvey2 represents all information needed to show dialog
// it should be used instead of AttachmentActionSurvey
type AttachmentActionSurvey2 struct {
	AttachmentActionSurvey
	TriggerID  string
	CallbackID string
	State      string
}

// This method is used to check if the struct is empty
func (a AttachmentActionSurvey) IsEmpty() bool {
	return reflect.DeepEqual(AttachmentActionSurvey{}, a)
}

func (a AttachmentActionConfirm) IsEmpty() bool {
	return reflect.DeepEqual(AttachmentActionConfirm{}, a)
}

// AttachmentActionElementOption - option for select drop down control
// deprecated. Use AttachmentActionElementPlainTextOption
type AttachmentActionElementOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// AttachmentActionElementOptionSortByLabelAsc is a type that is only used to sort menu options by label
type AttachmentActionElementOptionSortByLabelAsc []AttachmentActionElementOption

func (a AttachmentActionElementOptionSortByLabelAsc) Len() int           { return len(a) }
func (a AttachmentActionElementOptionSortByLabelAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AttachmentActionElementOptionSortByLabelAsc) Less(i, j int) bool { return a[i].Label < a[j].Label }

// AttachmentActionElementOptionGroup is a group of options
type AttachmentActionElementOptionGroup struct {
	Label   ui.PlainText                    `json:"label"`
	Options []AttachmentActionElementOption `json:"options"`
}

// AttachmentActionElementPlainTextOption is a structure that represents a single option of a drop-down input
type AttachmentActionElementPlainTextOption struct {
	Label ui.PlainText `json:"label"`
	Value string       `json:"value"`
}

// AttachmentActionTextElement is an element of a dialog
// https://api.slack.com/dialogs#select_elements
type AttachmentActionTextElement struct {
	Label string `json:"label,omitempty"`
	Name  string `json:"name,omitempty"`
	// A default value for this field
	Value string `json:"value"`
	// `text` for text box, `textarea` for multi-line plain text, `select` for multiple options
	ElemType string `json:"type,omitempty"`
	// https://api.slack.com/dialogs#attributes_text_elements
	ElemSubtype  AttachmentActionTextElementSubType   `json:"elem_subtype,omitempty"`
	Options      []AttachmentActionElementOption      `json:"options,omitempty"`
	OptionGroups []AttachmentActionElementOptionGroup `json:"option_groups,omitempty"`
	// A string displayed as needed to help guide users
	Placeholder string `json:"placeholder,omitempty"`
	// Helpful text provided to assist users
	Hint string `json:"hint,omitempty"`
}

type AttachmentActionTextElementSubType string

var (
	AttachmentActionTextElementNumberType AttachmentActionTextElementSubType = "number"
	AttachmentActionTextElementEmailType  AttachmentActionTextElementSubType = "email"
	AttachmentActionTextElementTelType    AttachmentActionTextElementSubType = "tel"
	AttachmentActionTextElementUrlType    AttachmentActionTextElementSubType = "url"
	AttachmentActionTextElementNoType     AttachmentActionTextElementSubType = ""
)

// ElemType is type of an AttachmentActionTextElement
type ElemType string

const (
	ElemTypeTextArea ElemType = "textarea"
	ElemTypeTextBox  ElemType = "text"
	ElemTypeSelect   ElemType = "select"
)

const (
	// EmptyPlainText is the empty text.
	EmptyPlainText ui.PlainText = ""
	// EmptyRichText is the empty text.
	EmptyRichText ui.RichText = ""
	// EmptyPlaceholder is the empty string. It might be used to find all places where a new placeholder is needed
	EmptyPlaceholder ui.PlainText = EmptyPlainText
)

// NewTextArea creates a field for entering paragraphs of text. Doesn't support rich text.
func NewTextArea(name string, label ui.PlainText, placeholderText ui.PlainText, initialValue ui.PlainText) AttachmentActionTextElement {
	return AttachmentActionTextElement{
		Label:       string(label),
		Name:        name,
		ElemType:    string(ElemTypeTextArea),
		Placeholder: string(placeholderText),
		Value:       string(initialValue),
	}
}

// NewTextBox creates a field for entering short text. Doesn't support rich text.
func NewTextBox(name string, label ui.PlainText, placeholderText ui.PlainText, initialValue ui.PlainText) AttachmentActionTextElement {
	return AttachmentActionTextElement{
		Label:       string(label),
		Name:        name,
		ElemType:    string(ElemTypeTextBox),
		Placeholder: string(placeholderText),
		Value:       string(initialValue),
	}
}

// NewDateInput creates a field for entering short text. Doesn't support rich text.
func NewDateInput(name string, label ui.PlainText, placeholderText ui.PlainText, initialValue ui.PlainText) AttachmentActionTextElement {
	return NewTextBox(name, label, placeholderText, initialValue)
}

// NewSelectOption creates a single option for drop-down input
func NewSelectOption(name string, label ui.PlainText) AttachmentActionElementPlainTextOption {
	return AttachmentActionElementPlainTextOption{
		Value: name,
		Label: label,
	}
}

// NewSelectOptionGroup constructs option group
func NewSelectOptionGroup(label ui.PlainText, options ...AttachmentActionElementPlainTextOption) AttachmentActionElementOptionGroup {
	options2 := ConvertPlainTextOptions(options...)
	
	return AttachmentActionElementOptionGroup{
		Label:   label,
		Options: options2,
	}
}

// NewSimpleOptionsSelect creates a drop-down field for selecting one of the options
func NewSimpleOptionsSelect(name string, label ui.PlainText, placeholderText ui.PlainText, initialValue string, options ...AttachmentActionElementPlainTextOption) AttachmentActionTextElement {
	options2 := ConvertPlainTextOptions(options...)
	
	return AttachmentActionTextElement{
		Label:       string(label),
		Name:        name,
		ElemType:    string(ElemTypeSelect),
		Placeholder: string(placeholderText),
		Value:       string(initialValue),
		Options:     options2,
	}
}

// ConvertPlainTextOptions converts one type to another
func ConvertPlainTextOptions(options ...AttachmentActionElementPlainTextOption) (options2 []AttachmentActionElementOption) {
	for _, o := range options {
		options2 = append(options2, AttachmentActionElementOption{
			Value: o.Value,
			Label: string(o.Label),
		})
	}
	return
}

// NewSimpleOptionGroupsSelect creates a drop-down field for selecting one of the options
func NewSimpleOptionGroupsSelect(name string, label ui.PlainText, placeholderText ui.PlainText, initialValue string, optionGroups ...AttachmentActionElementOptionGroup) AttachmentActionTextElement {
	return AttachmentActionTextElement{
		Label:       string(label),
		Name:        name,
		ElemType:    string(ElemTypeSelect),
		Placeholder: string(placeholderText),
		Value:       string(initialValue),
		OptionGroups:     optionGroups,
	}
}
