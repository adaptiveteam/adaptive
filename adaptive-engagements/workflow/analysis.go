package workflow

import (
	dialogFetcher "github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/pkg/errors"

	// "fmt"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/nlopes/slack"
)

// AnalyseMessage performs all typical operations to post standard analysis.
func AnalyseMessage(
	dialogFetcherDAO dialogFetcher.DAO,
	request slack.InteractionCallback,
	messageID chan mapper.MessageID,
	input utils.TextAnalysisInput,
	originalMessage InteractiveMessage,
) (message InteractiveMessage, err error) {
	// Once we receive the analysis from Meaning Cloud on the user's feedback, we post that result to the original message's thread
	analysis, errors1 := utils.AnalyzeText(dialogFetcherDAO, input)
	if len(errors1) > 0 {
		err = errors1[0]
	} else {
		message = originalMessage
		if analysis.Summary == "" { // if Summary is empty, we don't show it.
			err = errors.New("Analysis summary is empty")
		} else {
			msgID := <-messageID // waiting for message id of the original message to become available
			ctx := conversationContext(request, msgID)
			message = PresentTextAnalysisResults(ctx, analysis, originalMessage)
		}
	}
	return
}

func conversationContext(request slack.InteractionCallback, msgID mapper.MessageID) utils.ConversationContext {
	ctx := utils.ConversationContext{
		UserID:            request.User.ID,
		ConversationID:    string(GetConversationID(request)),
		OriginalMessageTs: msgID.Ts,
		ThreadTs:          msgID.Ts,
	}
	return ctx
}

// PresentTextAnalysisResults represents text analysis results to user in the given conversation context.
func PresentTextAnalysisResults(conversationContext utils.ConversationContext,
	analysisResults utils.TextAnalysisResults,
	originalMessage InteractiveMessage) (message InteractiveMessage) {
	color := utils.ColorStatusByIsGoodAndLength(analysisResults.IsGood, analysisResults.RecommendationsCount)
	message = originalMessage
	message.Color = color // Update the original attachments with the new color

	note := utils.RecommendationsMessage(analysisResults.TextAnalysisInput.Text, analysisResults.Summary, color)
	attach := note.Attachments[0]
	message.Thread = []InteractiveMessage{
		{
			PassiveMessage: PassiveMessage{
				Color:          color,
				Pretext:        ui.RichText(attach.Pretext),
				AttachmentText: analysisResults.Summary,
			},
		},
	}
	message.OverrideOriginal = true
	// colorCodedOriginalMessageOverrideNote := models.PlatformSimpleNotification{
	// 	Attachments: utils.RepaintAttachmentsWithColor(analysisResults.TextAnalysisInput.OriginalMessageAttachments, color)}
	return
	// []models.PlatformSimpleNotification{
	// 	SendNoteToUserThread(conversationContext, note),
	// 	UpdateOriginalMessageInUserChannel(conversationContext, colorCodedOriginalMessageOverrideNote),
	// }
}

// SendNoteToUserThread updates the message to send it to specific thread and user.
func SendNoteToUserThread(conversationContext utils.ConversationContext, note models.PlatformSimpleNotification) models.PlatformSimpleNotification {
	note.AsUser = true
	note.UserId = conversationContext.UserID
	note.Channel = conversationContext.ConversationID
	note.ThreadTs = conversationContext.ThreadTs
	return note
}

// UpdateOriginalMessageInUserChannel updates the message to send it to specific channel and user to override message with given ts.
func UpdateOriginalMessageInUserChannel(conversationContext utils.ConversationContext, note models.PlatformSimpleNotification) models.PlatformSimpleNotification {
	note.AsUser = true
	note.UserId = conversationContext.UserID
	note.Channel = conversationContext.ConversationID
	note.Ts = conversationContext.OriginalMessageTs
	return note
}
