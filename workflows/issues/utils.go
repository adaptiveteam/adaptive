package issues

import (
	"time"
	"reflect"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"

	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

func runtimeData(d interface{}) *interface{} { return &d }

func caption(trueCaption ui.PlainText, falseCaption ui.PlainText) func(bool) ui.PlainText {
	return func(flag bool) (res ui.PlainText) {
		if flag {
			res = trueCaption
		} else {
			res = falseCaption
		}
		return
	}
}

// GetFieldString returns the string field value from an interface
func GetFieldString(i interface{}, field string) string {
	// Create a value for the slice.
	v := reflect.ValueOf(i)
	// Get the field of the slice element that we want to set.
	f := v.FieldByName(field)
	// Get value
	return f.String()
}

func channelizeID(msgID mapper.MessageID) (messageID chan mapper.MessageID) {
	messageID = make(chan mapper.MessageID, 1)
	messageID <- msgID
	return
}

func toMapperMessageID(id platform.TargetMessageID) mapper.MessageID {
	return mapper.MessageID{
		ConversationID: id.ConversationID,
		Ts:             id.Ts,
	}
}

func filterIssues(issues []Issue, p IssuePredicate) (res []Issue) {
	for _, i := range issues {
		if p(i) {
			res = append(res, i)
		}
	}
	return
}

func AsDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// ParseDateOrTimestamp handles the date/timestamp field value from database.
// It tries to parse 2006-01-02 or timestamp.
// if failed, logs failure and returns empty time with false flag.
func ParseDateOrTimestamp(dateOrTimestampOrEmpty string) (t time.Time, isDefined bool) {
	return core.ParseDateOrTimestamp(dateOrTimestampOrEmpty)
}
