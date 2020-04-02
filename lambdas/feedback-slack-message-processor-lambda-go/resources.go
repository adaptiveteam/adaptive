package lambda

import (
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

const CoacheeListPretext ui.RichText = "The list of coachees:"
const AdvocateListPretext ui.RichText = "The list of advocates:"

var (
	GiveFeedbackAction                             = "select"
	RequestFeedbackAction                          = "request"
	GiveFeedbackMessage                            = "Whom would you like to give feedback to?"
	RequestFeedbackMessage                         = "Whom would you like to request feedback from?"
	GiveFeedbackNoUsersExistMessage    ui.RichText = "There are no other Adaptive associated users to give feedback to"
	RequestFeedbackNoUsersExistMessage ui.RichText = "There are no other Adaptive associated users to request feedback from"

	FetchingReportMessage   = ui.RichText("Hang tight, fetching the report for you...")
	GeneratingReportMessage = ui.RichText("Generating report in the background. Will notify you once done. Please go on with your work.")

	InternalErrorMessage = ui.RichText("Uh oh! There was some issue on my side and I am looking into it.")
)
