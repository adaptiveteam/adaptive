package models

import (
	"encoding/json"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"net/url"
	"strings"
)

// LambdaRequest is a struct for string payload to invoke a lambda
// Deprecated: Not used anymore.
type LambdaRequest struct {
	Platform ClientPlatform `json:"platform,omitempty"`
	Payload  string         `json:"payload"`
}

// NamespacePayload is a message routed from Slack
// Deprecated: Use NamespacePayload4 instead.
type NamespacePayload struct {
	Id        string `json:"id"`
	Namespace string `json:"namespace"`
	Payload   string `json:"payload"`
}

// NamespacePayload4 is a message routed from Slack.
// two main differences - PlatformID and EventsAPIEvent in deserialized form
type NamespacePayload4 struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	PlatformRequest
}

// ToNamespacePayload converts to the old version
func (np2 NamespacePayload4) ToNamespacePayload() NamespacePayload {
	return NamespacePayload{
		Id:        np2.ID,
		Namespace: np2.Namespace,
		Payload:   np2.EventsAPIEventSerialized,
	}
}

// Marshall converts NamespacePayload4 to json and then to string
func (np2 NamespacePayload4) Marshall() string {
	byt, err := json.Marshal(np2)
	core.ErrorHandler(err, np2.Namespace, "Couldn't marshal event np2.ID="+np2.ID)
	return string(byt)
}

// GetActionPathUnsafe parses callback id as action path
func (np2 NamespacePayload4)GetActionPathUnsafe() ActionPath {
	return ParseActionPath(np2.InteractionCallback.CallbackID)
}

// UnmarshalNamespacePayloadJSON parses NamespacePayload
func UnmarshalNamespacePayloadJSON(jsMessage string) (NamespacePayload, error) {
	var res NamespacePayload
	err := json.Unmarshal([]byte(jsMessage), &res)
	return res, err
}

// UnmarshalNamespacePayloadJSONUnsafe unmarshals NamespacePayload and panics in case of errors.
func UnmarshalNamespacePayloadJSONUnsafe(jsMessage string) NamespacePayload {
	res, err := UnmarshalNamespacePayloadJSON(jsMessage)
	core.ErrorHandler(err, "request", "Could not unmarshal sns record to NamespacePayload")
	return res
}

// UnmarshalNamespacePayload4JSON parses NamespacePayload4
func UnmarshalNamespacePayload4JSON(jsMessage string) (np2 NamespacePayload4, err error) {
	err = json.Unmarshal([]byte(jsMessage), &np2)
	return np2, err
}

// UnmarshalNamespacePayload4JSONUnsafe unmarshals NamespacePayload4 and panics in case of errors.
func UnmarshalNamespacePayload4JSONUnsafe(jsMessage string) NamespacePayload4 {
	res, err := UnmarshalNamespacePayload4JSON(jsMessage)
	core.ErrorHandler(err, "request", "Could not unmarshal sns record to NamespacePayload4")
	return res
}

// SlackRequestType is enum of two possible slack events
type SlackRequestType string

const (
	InteractionSlackRequestType        SlackRequestType = "slack.InteractionCallback"
	DialogSubmissionSlackRequestType   SlackRequestType = "slack.DialogSubmissionCallback"
	DialogCancellationSlackRequestType SlackRequestType = "slack.DialogCancellationCallback"
	EventsAPIEventSlackRequestType     SlackRequestType = "slackevents.EventsAPIEvent"
)

// SlackRequest contains Slack message
type SlackRequest struct {
	Type SlackRequestType `json:"type"`
	// eventsAPIEventSerialized always contains serialized EventsAPIEvent
	EventsAPIEventSerialized string `json:"events_api_event"`
	// InteractionCallback contains data if the event is InteractionCallback
	InteractionCallback slack.InteractionCallback `json:"interaction_callback"`
	// DialogSubmissionCallback contains data if it's also DialogSubmissionCallback
	DialogSubmissionCallback slack.DialogSubmissionCallback `json:"dialog_submission_callback"`
}

// PlatformRequest contains both platform id and Slack message
type PlatformRequest struct {
	TeamID TeamID
	SlackRequest
}

// ToEventsAPIEventUnsafe parses serialized version of EventsAPIEvent
func (sr SlackRequest) ToEventsAPIEventUnsafe() (evt slackevents.EventsAPIEvent) {
	evt, err := slackevents.ParseEvent(
		json.RawMessage(sr.EventsAPIEventSerialized),
		slackevents.OptionNoVerifyToken(),
	)
	core.ErrorHandler(err, "slack-request", "Couldn't parse EventsAPIEvent in "+sr.EventsAPIEventSerialized)
	return
}

// EventsAPIEvent constructs SlackRequest from serialized EventsAPIEvent
func EventsAPIEvent(serialized string) (res SlackRequest) {
	return ParseEventsAPIEventAsSlackRequestUnsafe(serialized)
}

// ParseBodyAsSlackRequestUnsafe parses body string either using `QueryUnescape`
// when there is `payload=` prefix, or without unescape otherwise.
func ParseBodyAsSlackRequestUnsafe(body string) (res SlackRequest) {
	serialized := ""
	payloadPrefix := "payload="
	if strings.HasPrefix(body, "payload=") {
		var err error
		serialized, err = url.QueryUnescape(strings.TrimPrefix(body, payloadPrefix))
		core.ErrorHandler(err, "SlackRequest", "Couldn't unescape "+body)
	} else {
		serialized = body
	}
	return ParseEventsAPIEventAsSlackRequestUnsafe(serialized)
}

// ParseEventsAPIEventAsSlackRequest constructs SlackRequest from serialized EventsAPIEvent
func ParseEventsAPIEventAsSlackRequest(serialized string) (res SlackRequest) {
	return ParseEventsAPIEventAsSlackRequestUnsafe(serialized)
}

// ParseEventsAPIEventAsSlackRequestUnsafe constructs SlackRequest from serialized EventsAPIEvent
func ParseEventsAPIEventAsSlackRequestUnsafe(serialized string) (res SlackRequest) {
	res.EventsAPIEventSerialized = serialized
	evt, err := slackevents.ParseEvent(
		json.RawMessage(serialized),
		slackevents.OptionNoVerifyToken(),
	)
	core.ErrorHandler(err, "slack-request", "Couldn't parse EventsAPIEvent in "+serialized)
	switch evt.Type {
	case string(slack.InteractionTypeDialogSubmission):
		res.Type = DialogSubmissionSlackRequestType
		err = json.Unmarshal([]byte(serialized), &res.InteractionCallback)
		core.ErrorHandler(err, "slack-request", "Couldn't parse InteractionCallback and DialogSubmissionCallback in "+serialized)
		res.DialogSubmissionCallback = res.InteractionCallback.DialogSubmissionCallback
	case string(slack.InteractionTypeInteractionMessage):
		res.Type = InteractionSlackRequestType
		err = json.Unmarshal([]byte(serialized), &res.InteractionCallback)
		core.ErrorHandler(err, "slack-request", "Couldn't parse InteractionCallback in "+serialized)
	case string(slack.InteractionTypeDialogCancellation):
		res.Type = DialogCancellationSlackRequestType
		err = json.Unmarshal([]byte(serialized), &res.InteractionCallback)
		core.ErrorHandler(err, "slack-request", "Couldn't parse InteractionCallback and DialogSubmissionCallback in "+serialized)
	default:
		res.Type = EventsAPIEventSlackRequestType
	}
	return
}
