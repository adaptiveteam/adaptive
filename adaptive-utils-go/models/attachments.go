package models

import (
	"fmt"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/say"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

func EmptyAttachs() []model.Attachment {
	return []model.Attachment{}
}

func EmptyActionConfirm() ebm.AttachmentActionConfirm {
	return ebm.AttachmentActionConfirm{}
}

func EmptyActionMenuOptions() []ebm.MenuOption {
	return []ebm.MenuOption{}
}

func EmptyActionMenuOptionGroups() []ebm.MenuOptionGroup {
	return []ebm.MenuOptionGroup{}
}

// SimpleAttachAction constructs action with callback
func SimpleAttachAction(mc MessageCallback, actionName AttachActionName, text ui.PlainText) *model.AttachmentAction {
	return GenAttachAction(mc, actionName, string(text), EmptyActionConfirm(), false)
}

// DialogAttachAction is an action that shows dialog
func DialogAttachAction(mc MessageCallback, actionName AttachActionName, text ui.PlainText) *model.AttachmentAction {
	return GenAttachAction(mc, actionName, string(text), EmptyActionConfirm(), false)
}

// ConfirmAttachAction constructs an action that requires confirmation.
// NB. Default style for actions that require confirmation is "danger"
func ConfirmAttachAction(mc MessageCallback, actionName AttachActionName, text ui.PlainText,
	confirm ebm.AttachmentActionConfirm) *model.AttachmentAction {
	return GenAttachAction(mc, actionName, string(text), confirm, true)
}

func GenAttachAction(mc MessageCallback, actionName AttachActionName, text string, confirm ebm.AttachmentActionConfirm,
	danger bool) *model.AttachmentAction {
	actionStyle := core.IfThenElse(danger, ebm.AttachmentActionStyleDanger, ebm.AttachmentActionStylePrimary).(ebm.AttachmentActionStyle)
	builder := eb.NewAttachmentActionBuilder().
		Name(fmt.Sprintf("%s_%s", mc.Action, actionName)).
		Text(text).
		ActionType(ebm.AttachmentActionTypeButton).
		Style(actionStyle).
		Value(mc.ToCallbackID())
	attachAction, _ := core.IfThenElse(confirm.IsEmpty(), builder, builder.Confirm(confirm)).(*eb.AttachmentActionBuilder).Build()
	return attachAction
}

// ConcatPrefixOpt - concatenates prefix and text if text is nonempty.
// In case when text is empty returns empty result.
// It's equivalent to Scala: textOpt.map(prefix + _)
func ConcatPrefixOpt(prefix, text string) (res string) {
	if text == "" {
		res = ""
	} else {
		res = prefix + text
	}
	return
}

// AppendOptionalAction appends action to given array if it's not empty.
func AppendOptionalAction(actions []ebm.AttachmentAction, actionOpt *ebm.AttachmentAction) (res []ebm.AttachmentAction) {
	if actionOpt == nil {
		res = actions
	} else {
		res = append(actions, *actionOpt)
	}
	return
}

// LearnMoreBaseLink is the path to documentation
const LearnMoreBaseLink = "https://adaptiveteam.github.io/"

// LearnMoreAction creates an attachment action for 'Learn More'.
// it only returns valid value when trailPath is nonempty. Otherwise - nil.
func LearnMoreAction(trailPath string) *ebm.AttachmentAction {
	if trailPath == "" {
		return nil
	}
	attachAction, _ := eb.NewAttachmentActionBuilder().
		Text(string(say.Plain("Learn more..."))).
		ActionType(ebm.AttachmentActionTypeButton).
		Url(LearnMoreBaseLink + trailPath).
		Build()
	return attachAction
}

// NowAttachAction Attachment action for 'Now'
// deprecated. Use SimpleAttachAction with models.Now
func NowAttachAction(mc MessageCallback, text string, confirm ebm.AttachmentActionConfirm) *model.AttachmentAction {
	return GenAttachAction(mc, Now, text, confirm, false)
}

// UpdateAttachAction Attachment action for 'Update'
// deprecated. Use SimpleAttachAction with models.Update
func UpdateAttachAction(mc MessageCallback, text string, confirm ebm.AttachmentActionConfirm) *ebm.AttachmentAction {
	return GenAttachAction(mc, Update, text, confirm, true)
}

// IgnoreAttachAction Attachment action for 'Ignore'
// deprecated. Use SimpleAttachAction with models.Ignore
func IgnoreAttachAction(mc MessageCallback, text string, confirm ebm.AttachmentActionConfirm) *ebm.AttachmentAction {
	return GenAttachAction(mc, Ignore, text, confirm, true)
}

// BackAttachAction Attachment action for 'Back'
// deprecated. Use SimpleAttachAction with models.Back
func BackAttachAction(mc MessageCallback, text string, confirm ebm.AttachmentActionConfirm) *ebm.AttachmentAction {
	return GenAttachAction(mc, Back, text, confirm, false)
}

// Attachment action for `Select`
func SelectAttachAction(mc MessageCallback, actionName AttachActionName, text string, options []ebm.MenuOption, optionGroups []ebm.MenuOptionGroup) *ebm.AttachmentAction {
	callbackId := mc.ToCallbackID()
	baseAttachAction := eb.NewAttachmentActionBuilder().
		Name(fmt.Sprintf("%s_%s", mc.Action, actionName)).
		Text(text).
		ActionType(ebm.AttachmentActionTypeSelect).
		Value(callbackId)
	if len(options) > 0 {
		baseAttachAction.Options(options)
	}
	if len(optionGroups) > 0 {
		baseAttachAction.OptionGroups(optionGroups)
	}
	action, _ := baseAttachAction.Build()
	return action
}

// Fills current survey with provided values. Keys in 'attribs' should match the dialog element names
func FillSurvey(sur ebm.AttachmentActionSurvey, attribs map[string]string) ebm.AttachmentActionSurvey {
	temp := sur
	for i := range temp.Elements {
		// Check if the map contains current dialog element name
		if val, ok := attribs[temp.Elements[i].Name]; ok {
			// If it contains, replace it's value from map
			temp.Elements[i].Value = val
		}
	}
	return temp
}

func AttachmentFields(elems []KvPair) []ebm.AttachmentField {
	var attachFields []ebm.AttachmentField

	for _, each := range elems {
		attachFields = append(attachFields,
			ebm.AttachmentField{
				Title: each.Key,
				Value: each.Value,
				Short: false,
			},
		)
	}

	return attachFields
}
