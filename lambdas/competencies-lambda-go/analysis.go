package competencies

import (
	"fmt"

	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/slack-go/slack"
)

func analyseMessage(request slack.InteractionCallback, messageID chan mapper.MessageID, input utils.TextAnalysisInput) []models.PlatformSimpleNotification {
	// Once we receive the analysis from Meaning Cloud on the user's feedback, we post that result to the original message's thread
	analysis, errors := utils.AnalyzeText(dialogFetcherDao, input)
	platform.Debug(request, "Analyzed")
	if len(errors) > 0 {
		platform.Debug(request, fmt.Sprintf("Errors in text analysis: %v", errors))
	}

	notes := responses()
	if analysis.Summary != "" {
		msgID := <-messageID // waiting for message id of the original message to become available
		ctx := conversationContext(request, msgID)
		notes = ctx.PresentTextAnalysisResults(analysis)
	}
	return notes
}
