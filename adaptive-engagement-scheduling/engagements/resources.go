package engagements

import (
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

func mapUserIDToDisplayName(userIDs []string) (dns []ui.RichText) {
	for _, each := range userIDs {
		dns = append(dns, ui.Sprintf("<@%s>", each))
	}
	return
}

func renderUserIDList(userIDs []string) ui.RichText {
	dns := mapUserIDToDisplayName(userIDs)
	return ui.Join(dns, ", ")
}

// ProvidedFeedbackConfirmationAndSuggestProvidingMoreTemplate is a text template for the case
// when we want to acknowledge that the user has provided feedback to a few peers and
// we recommend to provide more feedback.
func ProvidedFeedbackConfirmationAndSuggestProvidingMoreTemplate(userIDs []string) (text ui.RichText) {
	if len(userIDs) > 0 {
		text = ui.Sprintf(
			"You have provided feedback to these colleagues: %s. "+
				"Would you like to update your feedback for them or give feedback to anyone else?",
			renderUserIDList(core.Distinct(userIDs)))
	} else {
		text = "You haven't provided feedback to anyone yet. Do you want to do that?"
	}
	return
}

// RequestedFeedbackConfirmationAndSuggestRequestingMoreTemplate - 
func RequestedFeedbackConfirmationAndSuggestRequestingMoreTemplate(userIDs []string) (text ui.RichText) {
	if len(userIDs) > 0 {
		text = ui.Sprintf(
			"You have requested feedback from these peers: %s. "+
				"Would you like to request feedback from anyone else?",
			renderUserIDList(core.Distinct(userIDs)))
	} else {
		text = "You may request feedback from anyone. Would you like to do that?"
	}

	return
}
