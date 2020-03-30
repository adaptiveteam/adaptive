package lambda

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

var (
	RemoveEngagementLabel = ui.PlainText("Remove from here")
)

func FeedbackRequestedAskIfYouWantToProvideTemplate(userID string) string {
	return fmt.Sprintf("<@%s> is requesting feedback from you. Do you have feedback to provide that will help improve your colleague's growth and performance?", userID)
}

func ConfirmFeedbackRequestedTemplate(userID string) ui.RichText {
	return ui.RichText(fmt.Sprintf(
		"Ok, I have scheduled a notfication to <@%s> "+
			"about your request for feedback from them "+
			"if they haven't already provided feedback to you.",
		userID),
	).Italics()
}
