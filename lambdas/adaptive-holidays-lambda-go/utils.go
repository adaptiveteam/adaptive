package lambda

import (
	"encoding/json"
	"fmt"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"net/url"
	"strconv"
	"strings"
)

func parseGwSNSRequest(np models.NamespacePayload) slackevents.EventsAPIEvent {
	ueRequest, err := url.QueryUnescape(np.Payload)
	core.ErrorHandler(err, namespace, "Could not un-escape the request body")

	requestPayload := strings.Replace(ueRequest, "payload=", "", -1)

	res, err := utils.ParseApiRequest(requestPayload)
	core.ErrorHandler(err, namespace, "Could not parse the event")

	return res
}

func publish(msg models.PlatformSimpleNotification) {
	_, err := sns.Publish(msg, platformNotificationTopic)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not publish message to %s topic", platformNotificationTopic))
}
func publishAll(notes []models.PlatformSimpleNotification) {
	for _, note := range notes {
		publish(note)
	}
}
func deleteOriginalEng(userID, channel, ts string) {
	utils.DeleteOriginalEng(userID, channel, ts, publish)
}

func unmarshalNamespacePayloadJSON(jsMessage string) models.NamespacePayload {
	var res models.NamespacePayload
	err := json.Unmarshal([]byte(jsMessage), &res)
	core.ErrorHandler(err, namespace, "Could not unmarshal sns record to NamespacePayload")
	return res
}

func unmarshallSlackInteractionMsg(msg string) slack.InteractionCallback {
	var res slack.InteractionCallback
	res, err := utils.ParseAsInteractionMsg(msg)
	core.ErrorHandler(err, namespace, "InteractionCallback: Could not parse")
	return res
}

func unmarshallSlackDialogSubmissionMsg(msg string) (slack.InteractionCallback, slack.DialogSubmissionCallback) {
	var res slack.InteractionCallback
	res, err := utils.ParseAsInteractionMsg(msg)
	core.ErrorHandler(err, namespace, "DialogSubmissionCallback: Could not parse")
	return res, res.DialogSubmissionCallback
}

func platformSimpleNotification(request slack.InteractionCallback, message string) models.PlatformSimpleNotification {
	return models.PlatformSimpleNotification{
		UserId:  request.User.ID,
		Channel: request.Channel.ID,
		Message: message}
}

func platformSimpleNotificationInThread(request slack.InteractionCallback, message string) models.PlatformSimpleNotification {
	return models.PlatformSimpleNotification{
		UserId:   request.User.ID,
		Channel:  request.Channel.ID,
		Message:  message,
		ThreadTs: timeStamp(request),
	}
}

func platformSimpleNotificationOverride(request slack.InteractionCallback, message string) models.PlatformSimpleNotification {
	return models.PlatformSimpleNotification{
		UserId:  request.User.ID,
		Channel: request.Channel.ID,
		Message: message,
		Ts:      request.MessageTs,
	}
}

func recoverGracefully(request slack.InteractionCallback) {
	if err := recover(); err != nil {
		publish(platformSimpleNotification(request,
			fmt.Sprintf("Error: %s", err)))
	}
}

func debug(request slack.InteractionCallback, message string) {
	msg := "Debug: " + message
	if isInteractiveDebugEnabled {
		publish(platformSimpleNotification(request, msg))
	}
	fmt.Println(msg)
}

func errorHandler(request slack.InteractionCallback, msg string, err error) {
	if err != nil {
		message := fmt.Sprintf("%s while serving request %s \n(%s)", msg, request.CallbackID, err.Error())
		if isInteractiveDebugEnabled {
			publish(platformSimpleNotification(request, "Error: "+message))
		}
		core.ErrorHandler(err, namespace, message)
	}
}

func dummyMessageCallback(source string) models.MessageCallback {
	year, month := core.CurrentYearMonth()
	return models.MessageCallback{Module: HolidaysNamespace,
		Source: source, Topic: "dummy", Action: "noaction",
		Month: strconv.Itoa(int(month)),
		Year:  strconv.Itoa(year)}
}

func urgency(urgent bool) models.PriorityValue {
	return core.IfThenElse(urgent, models.UrgentPriority, models.HighPriority).(models.PriorityValue)
}

func mapAdHocHolidayString(vs []models.AdHocHoliday, f func(models.AdHocHoliday) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func mapAdHocHolidayPlatformSimpleNotification(vs []models.AdHocHoliday,
	f func(models.AdHocHoliday) models.PlatformSimpleNotification) []models.PlatformSimpleNotification {
	vsm := make([]models.PlatformSimpleNotification, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// timeStamp extracts timestamp from the original message
// When the original message is from a thread, we need to post to the same thread
// Below logic checks if the incoming message is from a thread
func timeStamp(request slack.InteractionCallback) string {
	ts := request.OriginalMessage.ThreadTimestamp
	if request.OriginalMessage.ThreadTimestamp == "" {
		ts = request.MessageTs
	}
	return ts
}
