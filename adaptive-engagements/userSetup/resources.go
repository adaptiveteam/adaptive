package userSetup

import (
	"github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

var (
	MeetingTimeUserAttributeID = "meeting_time"
)

var (
	QueryMeetingTime  ui.RichText = "When would you like to meet with me?"
	QueryMeetingTime2 ui.RichText = "Please choose a time to meet with me"

	ChooseMeetingTime  ui.RichText = "Choose meeting time..."
	MeetingTimeUpdated ui.RichText = ui.RichText("Your meeting time has been updated").Italics()

	IntroductionMessage ui.RichText = "Hi! I am Adaptive. In order for us to work together effectively, I need to know a few things about you first."

	UpdateActionName ui.PlainText = "Update"
	CancelActionName ui.PlainText = "Cancel"
	ThankYouMessage  ui.PlainText = "Thank you!"

	PromptToUpdateSettings  ui.RichText = "Would you like to update the settings?"
	PromptToUpdateSettings2 ui.RichText = "Update Settings"

	IncompleteSettings ui.RichText = "Hello. You haven't completed Adaptive configuration."
)

// MeetingIsScheduledFor is a template for notifying user about meeting time.
func MeetingIsScheduledFor(meetingTime business_time.LocalTime) ui.RichText {
	return ui.RichText("Awesome! Our meeting is scheduled for " + meetingTime.ToUserFriendly()).Italics()
}

// MeetingIsScheduledForCalmNotice is a template for notifying user about meeting time.
func MeetingIsScheduledForCalmNotice(meetingTime business_time.LocalTime) ui.RichText {
	return ui.RichText("Currently we are meeting at " + meetingTime.ToUserFriendly()).Italics()
}
