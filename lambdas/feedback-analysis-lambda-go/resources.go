package lambda

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
)

const (
	FeedbackFieldTitle    = "Feedback"
	ConfidenceFactorTitle = "Confidence Factor"
)

func ReportAuthorTemplate(mc models.MessageCallback, value models.AdaptiveValue) ebm.AttachmentAuthor {
	return ebm.AttachmentAuthor{
		Name: fmt.Sprintf("<@%s>'s %s", mc.Target, value.Name),
	}
}

func EditActionTemplate(mc models.MessageCallback, uf models.UserFeedback) ebm.AttachmentAction {
	attachAction, _ := eb.NewAttachmentActionBuilder().
		Name(fmt.Sprintf("ask_%s", uf.ValueID)).
		Text(models.EditLabel).
		ActionType(models.ButtonType).
		Value(mc.ToCallbackID()).
		Build()
	return *attachAction
}

func FeedbackAttachmentTemplate(mc models.MessageCallback, uf models.UserFeedback,
	value models.AdaptiveValue) (attach ebm.Attachment, err error) {
	attach1, _ := eb.NewAttachmentBuilder().
		CallbackId(mc.ToCallbackID()).
		Author(ReportAuthorTemplate(mc, value)).
		Color(models.BlueColorHex).
		Fields([]ebm.AttachmentField{
			{Title: FeedbackFieldTitle, Value: uf.Feedback},
			{Title: ConfidenceFactorTitle, Value: coaching.Feedback360RatingMap[uf.ConfidenceFactor]},
		}).
		Actions([]ebm.AttachmentAction{EditActionTemplate(mc, uf)}).
		Build()
	attach = *attach1
	return
}
