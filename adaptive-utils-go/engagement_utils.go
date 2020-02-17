package adaptive_utils_go

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
)

// AttachActionElementOptions converts slice of key-value pairs into attachment action element options
func AttachActionElementOptions(elems []models.KvPair) []ebm.AttachmentActionElementOption {
	var elemOptions []ebm.AttachmentActionElementOption
	for _, each := range elems {
		elemOptions = append(elemOptions, ebm.AttachmentActionElementOption{
			Label: each.Key,
			Value: each.Value,
		})
	}
	return elemOptions
}

// AttachmentSurvey creates a survey instance
func AttachmentSurvey(title string, elems []ebm.AttachmentActionTextElement) ebm.AttachmentActionSurvey {
	return ebm.AttachmentActionSurvey{
		Title:       core.ClipString(title, 24, "..."),
		SubmitLabel: models.SubmitLabel,
		Elements:    elems,
	}
}

// AttachmentConfirm instantiates confirmation action
func AttachmentConfirm(title, text string) ebm.AttachmentActionConfirm {
	op := ebm.AttachmentActionConfirm{}
	op.OkText = models.YesLabel
	op.DismissText = models.CancelLabel
	if title != core.EmptyString {
		op.Title = title
	}
	if text != core.EmptyString {
		op.Text = text
	}
	return op
}

type publish func(models.PlatformSimpleNotification)

// DeleteOriginalEng deletes the original engagement from chat space based on the timestamp provided
// Based on the user, it's going to look for the chat token
func DeleteOriginalEng(userID, channel, ts string, publishCallback publish) {
	publishCallback(MakeCommandToDeleteOldText(userID, channel, ts))
}

// MakeCommandToDeleteOldText creates a PlatformSimpleNotification that will delete the old text
func MakeCommandToDeleteOldText(userID, channel, ts string) models.PlatformSimpleNotification {
	return MakeCommandToReplaceChatMessage(userID, channel, ts, "")
}

// MakeCommandToReplaceChatMessage creates a PlatformSimpleNotification that will 
// replace old message with the new one
func MakeCommandToReplaceChatMessage(userID, channel, ts, text string) models.PlatformSimpleNotification {
	return models.PlatformSimpleNotification{
		UserId: userID, 
		Channel: channel, 
		Ts: ts,
		Message: text, 
		Attachments: models.EmptyAttachs()}
}
