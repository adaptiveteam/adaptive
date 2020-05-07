package engagement_slack_mapper

import (
	"encoding/json"
	core_utils_go "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/slack-go/slack"
	"strconv"
)

const (
	// https://api.slack.com/docs/interactive-message-field-guide#option_fields
	SlackMenuOptionTextMaxLength = 30
)

func fields(list []model.AttachmentField) []slack.AttachmentField {
	var op []slack.AttachmentField
	for _, field := range list {
		op = append(op, slack.AttachmentField{
			Title: field.Title,
			Value: field.Value,
			Short: field.Short,
		})
	}
	return op
}

func options(list []model.MenuOption) []slack.AttachmentActionOption {
	var op []slack.AttachmentActionOption
	for _, field := range list {
		op = append(op, slack.AttachmentActionOption{
			Text:        core_utils_go.ClipString(field.Text, SlackMenuOptionTextMaxLength, "..."),
			Value:       field.Value,
			Description: field.Description,
		})
	}
	return op
}

func optionGroups(list []model.MenuOptionGroup) []slack.AttachmentActionOptionGroup {
	var op []slack.AttachmentActionOptionGroup
	for _, field := range list {
		op = append(op, slack.AttachmentActionOptionGroup{
			Text:    core_utils_go.ClipString(field.Text, SlackMenuOptionTextMaxLength, "..."),
			Options: options(field.Options),
		})
	}
	return op
}

func actions(list []model.AttachmentAction) (actions []slack.AttachmentAction) {
	for _, each := range list {
		// If confirm struct is empty, we don't attach to the action
		// When specifying with empty struct, confirm action is taking default values
		// That's why this condition
		attachAction := slack.AttachmentAction{
			Name:         each.Name,
			Text:         each.Text,
			Value:        each.Value,
			Style:        string(each.Style),
			URL:          each.Url,
			Options:      options(each.Options),
			OptionGroups: optionGroups(each.OptionGroups),
			DataSource:   string(each.DataSource),
		}
		if !each.Confirm.IsEmpty() {
			attachAction.Confirm = &slack.ConfirmationField{
				Title:       each.Confirm.Title,
				Text:        each.Confirm.Text,
				OkText:      each.Confirm.OkText,
				DismissText: each.Confirm.DismissText,
			}
		}
		SetButtonOrSelectType(&attachAction, each.ActionType)
		actions = append(actions, attachAction)
	}
	return
}

// SetButtonOrSelectType is a workaround for AttachmentAction.Type field.
// It has type `actionType` which is private to slack.
func SetButtonOrSelectType(action *slack.AttachmentAction, tpe model.AttachmentActionType) {
	switch tpe {
	case "button": action.Type = "button"
	case "select": action.Type = "select"
	}
}

func Attachments(attachs []model.Attachment) []slack.Attachment {
	var slackAttachs []slack.Attachment
	for _, attach := range attachs {
		slackAttachs = append(slackAttachs, slack.Attachment{
			Title:      attach.Title,
			TitleLink:  attach.TitleLink,
			Text:       attach.Text,
			Color:      attach.Color,
			Fallback:   attach.Fallback,
			CallbackID: attach.CallbackId,
			Pretext:    attach.Pretext,
			AuthorName: attach.Author.Name,
			AuthorLink: attach.Author.Link,
			AuthorIcon: attach.Author.Icon,
			Fields:     fields(attach.Fields),
			Actions:    actions(attach.Actions),
			ImageURL:   attach.ImageUrl,
			ThumbURL:   attach.ThumbUrl,
			Footer:     attach.Footer.Text,
			FooterIcon: attach.Footer.Icon,
			Ts:         json.Number(strconv.FormatInt(attach.Footer.Timestamp, 10)),
		})
	}
	return slackAttachs
}

func slackMapper(message model.Message) []slack.MsgOption {
	return []slack.MsgOption{slack.MsgOptionText(message.Text, false),
		slack.MsgOptionAttachments(Attachments(message.Attachments)...), slack.MsgOptionAsUser(true)}
}
