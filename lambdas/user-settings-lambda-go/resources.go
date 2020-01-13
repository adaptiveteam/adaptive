package lambda

import (
	"github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

var (
	YourEngagementHasBeenCancelled ui.RichText = ui.RichText("Your engagement has been canceled").Italics()
	LetMeTakeCare                  ui.RichText = ui.RichText("_Let me take care of this for you..._").Italics()
)

const (
	MeetingTimeUserAttributeID = "meeting_time"
)

// MeetingIsScheduledFor is a template for notifying user about meeting time.
func MeetingIsScheduledFor(meetingTime business_time.LocalTime) ui.RichText {
	return ui.RichText("Awesome! Our meeting is scheduled for " + meetingTime.ToUserFriendly()).Italics()
}
