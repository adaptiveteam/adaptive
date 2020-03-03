package EngagementBuilder

import (
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/say"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// AttachmentAction builder pattern code
type AttachmentActionBuilder struct {
	attachmentAction *model.AttachmentAction
}

func NewAttachmentActionBuilder() *AttachmentActionBuilder {
	attachmentAction := &model.AttachmentAction{}
	b := &AttachmentActionBuilder{attachmentAction: attachmentAction}
	return b
}

func (b *AttachmentActionBuilder) Name(name string) *AttachmentActionBuilder {
	b.attachmentAction.Name = name
	return b
}

func (b *AttachmentActionBuilder) Text(text string) *AttachmentActionBuilder {
	b.attachmentAction.Text = text
	return b
}

func (b *AttachmentActionBuilder) Plain(text ui.PlainText) *AttachmentActionBuilder {
	b.attachmentAction.Text = string(text)
	return b
}

func (b *AttachmentActionBuilder) ActionType(actionType model.AttachmentActionType) *AttachmentActionBuilder {
	b.attachmentAction.ActionType = actionType
	return b
}

func (b *AttachmentActionBuilder) Value(value string) *AttachmentActionBuilder {
	b.attachmentAction.Value = value
	return b
}

func (b *AttachmentActionBuilder) Confirm(confirm model.AttachmentActionConfirm) *AttachmentActionBuilder {
	b.attachmentAction.Confirm = confirm
	return b
}

func (b *AttachmentActionBuilder) Style(style model.AttachmentActionStyle) *AttachmentActionBuilder {
	b.attachmentAction.Style = style
	return b
}

func (b *AttachmentActionBuilder) Options(options []model.MenuOption) *AttachmentActionBuilder {
	b.attachmentAction.Options = options
	return b
}

func (b *AttachmentActionBuilder) SelectedOptions(option model.MenuOption) *AttachmentActionBuilder {
	b.attachmentAction.SelectedOptions = option
	return b
}

func (b *AttachmentActionBuilder) OptionGroups(optionGroups []model.MenuOptionGroup) *AttachmentActionBuilder {
	b.attachmentAction.OptionGroups = optionGroups
	return b
}

func (b *AttachmentActionBuilder) Url(url string) *AttachmentActionBuilder {
	b.attachmentAction.Url = url
	return b
}

func (b *AttachmentActionBuilder) DataSource(source model.AttachmentActionDataSource) *AttachmentActionBuilder {
	b.attachmentAction.DataSource = source
	return b
}

// Build converts builder to AttachmentAction.
// Deprecated: Use ToAttachmentAction
// This method is misleading as it always returns `nil` as the second
func (b *AttachmentActionBuilder) Build() (*model.AttachmentAction, error) {
	return b.attachmentAction, nil
}

func (b AttachmentActionBuilder) ToAttachmentAction() model.AttachmentAction {
	return *b.attachmentAction
}

const (
	MenuSelectType = "select"
	ButtonType     = "button"
	OkColor        = "good"
	WarningColor   = "warning"
	DangerColor    = "danger"
)

// NewButton constructs a simple button
func NewButton(name string, value string, text ui.PlainText) model.AttachmentAction {
	return NewAttachmentActionBuilder().
		Name(name).
		Text(string(text)).
		ActionType(ButtonType).
		Value(value).
		ToAttachmentAction()
}

// NewButtonDanger constructs a button with confirmation. The action is only triggered when
// user confirms. Otherwise the cancelText is displayed.
func NewButtonDanger(name string, value string, text, cancelText ui.PlainText) model.AttachmentAction {
	return NewAttachmentActionBuilder().
		Name(name).
		Text(string(text)).
		ActionType(ButtonType).
		Value(value).
		Style(DangerColor).
		Confirm(model.AttachmentActionConfirm{
			OkText:      string(say.Plain(say.Yes)),
			DismissText: string(say.Plain(say.Cancel)),
		}).
		ToAttachmentAction()
}
