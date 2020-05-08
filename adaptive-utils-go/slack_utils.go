package adaptive_utils_go

import (
	"github.com/pkg/errors"
	"encoding/json"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func ParseApiRequest(payload string) (slackevents.EventsAPIEvent, error) {
	return slackevents.ParseEvent(json.RawMessage(payload), slackevents.OptionNoVerifyToken())
}

// UnsafeUnmarshallSlackEventsAPIEvent unmarshalls EventsAPIEvent. 
// Deprecated: use UnsafeUnmarshallSlackEventsAPIEventUnsafe
func UnsafeUnmarshallSlackEventsAPIEvent(payload, namespace string) slackevents.EventsAPIEvent {
	return UnsafeUnmarshallSlackEventsAPIEventUnsafe(payload, namespace)
}

// UnsafeUnmarshallSlackEventsAPIEventUnsafe unmarshalls EventsAPIEvent. Panics in case of any errors.
func UnsafeUnmarshallSlackEventsAPIEventUnsafe(payload, namespace string) slackevents.EventsAPIEvent {
	return UnmarshallSlackEventsAPIEventUnsafe(payload, namespace)
}

// UnmarshallSlackEventsAPIEventUnsafe unmarshalls EventsAPIEvent. Panics in case of any errors.
func UnmarshallSlackEventsAPIEventUnsafe(payload, namespace string) slackevents.EventsAPIEvent {
	res, err := slackevents.ParseEvent(json.RawMessage(payload), slackevents.OptionNoVerifyToken())
	core.ErrorHandler(err, namespace, "Could not parse to EventsAPIEvent")
	return res
}

// UnmarshallSlackInteractionMsg parses InteractionCallback.
// Deprecated: Use UnmarshallSlackInteractionCallbackUnsafe (which is named consistently)
func UnmarshallSlackInteractionMsg(msg string) slack.InteractionCallback {
	return UnmarshallSlackInteractionCallbackUnsafe(msg, "unknown-namespace")
}

// UnmarshallSlackInteractionCallbackUnsafe parses InteractionCallback. Panics in case of errors
func UnmarshallSlackInteractionCallbackUnsafe(msg, namespace string) (res slack.InteractionCallback) {
	res, err := UnmarshallSlackInteractionCallback(msg, namespace)
	core.ErrorHandler(err, namespace, "InteractionCallback: Could not parse")
	return res
}

// UnmarshallSlackInteractionCallback parses InteractionCallback.
func UnmarshallSlackInteractionCallback(payload, namespace string) (message slack.InteractionCallback, err error) {
	if payload == "" {
		err = fmt.Errorf("%s: empty payload is not a valid slack.InteractionCallback", namespace)
		return
	}
	err = json.Unmarshal([]byte(payload), &message)
	return
}

// UnmarshallSlackDialogSubmissionCallbackUnsafe parses InteractionCallback message as DialogSubmissionCallback
func UnmarshallSlackDialogSubmissionCallbackUnsafe(msg, namespace string) (slack.InteractionCallback, slack.DialogSubmissionCallback) {
	res := UnmarshallSlackInteractionCallbackUnsafe(msg, namespace)
	return res, res.DialogSubmissionCallback
}

// ParseAsInteractionMsg parses InteractionCallback.
// Deprecated: Use UnmarshallSlackInteractionCallback (named consistently)
func ParseAsInteractionMsg(payload string) (slack.InteractionCallback, error) {
	return UnmarshallSlackInteractionCallback(payload, "unknown-namespace-3")
}

// UnsafeUnmarshallSlackInteractionCallback unmarshalls slack.InteractionCallback
func UnsafeUnmarshallSlackInteractionCallback(payload, namespace string) (res slack.InteractionCallback) {
	err := json.Unmarshal([]byte(payload), &res)
	core.ErrorHandler(err, namespace, "Could not parse to interaction type message")
	return res
}

func ParseAsCallbackMsg(apiEvent slackevents.EventsAPIEvent) *slackevents.MessageEvent {
	switch apiEvent.InnerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		slackMsg := apiEvent.InnerEvent.Data.(*slackevents.MessageEvent)
		return slackMsg
	}
	return nil
}

func SlackFieldValue(attach slack.Attachment, name string) string {
	for _, each := range attach.Fields {
		if each.Title == name {
			return each.Value
		}
	}
	return ""
}

// SurveyState is a lazy string
type SurveyState func() string

// SlackSurvey sends dialog to Slack
func SlackSurvey(api *slack.Client, message slack.InteractionCallback,
	survey model.AttachmentActionSurvey, callbackID string,
	state SurveyState) (err error) {
	dialog, err := ConvertSurveyToSlackDialog(survey, message.TriggerID, callbackID, state())
	if err == nil {
		err = api.OpenDialog(message.TriggerID, dialog)
	}
	return
}

// ConvertSurveyToSlackDialog converts our model to slack Dialog
func ConvertSurveyToSlackDialog(survey model.AttachmentActionSurvey,
	triggerID string,
	callbackID string,
	surveyState string,
) (dialog slack.Dialog, err error) {
	var textElems []slack.DialogElement
	var selectElements []slack.DialogElement

	for _, each := range survey.Elements {
		if each.ElemType == string(slack.InputTypeSelect) {
			if len(each.Options) > 0 || len(each.OptionGroups) > 0 {
				selectElement := ConvertAttachmentActionTextElementToDialogElement(each)
				selectElements = append(selectElements, selectElement)
			} else {
				err = errors.New("Element (name=" + each.Name + "; label=" + each.Label + ") of type InputTypeSelect should have at least some options or option groups")
			}
		} else {
			textElement := ConvertAttachmentActionTextElementToDialogElement(each)
			textElems = append(textElems, textElement)
		}
	}

	if err == nil {
		dialog = slack.Dialog{
			TriggerID:      triggerID,
			CallbackID:     callbackID,
			Title:          survey.Title,
			State:          surveyState,
			SubmitLabel:    survey.SubmitLabel,
			NotifyOnCancel: true,
			Elements:       append(selectElements, textElems...),
		}
	}
	return dialog, err
}

// ConvertAttachmentActionTextElementToDialogElement converts to Slack API
func ConvertAttachmentActionTextElementToDialogElement(a model.AttachmentActionTextElement) slack.DialogElement {
	switch a.ElemType {
	case string(slack.InputTypeSelect):
		return convertSelectToDialogInputSelect(a)
	default:
		return convertTextElementToTextInputElement(a)
	}
}

func convertSelectToDialogInputSelect(a model.AttachmentActionTextElement) slack.DialogInputSelect {
	selectElement := slack.DialogInputSelect{}
	// https://api.slack.com/dialogs#text_elements
	selectElement.Name = core.ClipString(a.Name, models.SlackDialogSelectElementNameLimit, "...")
	selectElement.Type = slack.InputTypeSelect
	selectElement.Label = core.ClipString(a.Label, models.SlackDialogSelectElementLabelLimit, "...")
	selectElement.Value = a.Value
	var options []slack.DialogSelectOption
	for _, opt := range a.Options {
		options = append(options, slack.DialogSelectOption{
			Label: core.ClipString(opt.Label, models.SlackDialogSelectElementLabelLimit, "..."),
			Value: opt.Value,
		})
	}
	selectElement.Options = options
	selectElement.OptionGroups = []slack.DialogOptionGroup{}
	for _, optGroup := range a.OptionGroups {
		og := slack.DialogOptionGroup{
			Label:   core.ClipString(string(optGroup.Label), models.SlackDialogSelectElementLabelLimit, "..."),
			Options: []slack.DialogSelectOption{},
		}
		for _, opt := range optGroup.Options {
			og.Options = append(og.Options, slack.DialogSelectOption{
				Label: core.ClipString(opt.Label, models.SlackDialogSelectElementLabelLimit, "..."),
				Value: opt.Value,
			})
		}
		selectElement.OptionGroups = append(selectElement.OptionGroups, og)
	}
	return selectElement
}

func convertTextElementToTextInputElement(a model.AttachmentActionTextElement) slack.TextInputElement {
	textElement := slack.TextInputElement{}
	textElement.Label = core.ClipString(a.Label, models.SlackDialogSelectElementLabelLimit, "...")
	// this is the key that will have entered text as the value
	textElement.Name = core.ClipString(a.Name, models.SlackDialogSelectElementNameLimit, "...")
	textElement.Type = slack.InputType(a.ElemType)
	if string(a.ElemSubtype) != core.EmptyString {
		textElement.Subtype = slack.TextInputSubtype(a.ElemSubtype)
	}
	textElement.Value = a.Value
	textElement.Placeholder = a.Placeholder
	textElement.Hint = a.Hint
	return textElement
}

// InteractionCallback utilities
// - Responses

// InteractionCallbackSimpleResponse creates a simple message 
func InteractionCallbackSimpleResponse(request slack.InteractionCallback, message string) models.PlatformSimpleNotification {
	return models.PlatformSimpleNotification{
		UserId:  request.User.ID,
		Channel: request.Channel.ID,
		Message: message}
}

// InteractionCallbackSimpleNotificationInThread places the message in a thread connected to
func InteractionCallbackSimpleNotificationInThread(request slack.InteractionCallback, message string) models.PlatformSimpleNotification {
	return models.PlatformSimpleNotification{
		UserId:   request.User.ID,
		Channel:  request.Channel.ID,
		Message:  message,
		ThreadTs: TimeStamp(request),
	}
}

// InteractionCallbackOverrideRequestMessage creates a notification that will override the original message from request.
func InteractionCallbackOverrideRequestMessage(request slack.InteractionCallback, message string) models.PlatformSimpleNotification {
	return models.PlatformSimpleNotification{
		UserId:  request.User.ID,
		Channel: request.Channel.ID,
		Message: message,
		Ts:      request.MessageTs,
	}
}

// InteractionCallbackOverrideOriginalMessage creates a notification that will override the original message from request.
func InteractionCallbackOverrideOriginalMessage(request slack.InteractionCallback, message string) models.PlatformSimpleNotification {
	return models.PlatformSimpleNotification{
		UserId:  request.User.ID,
		Channel: request.Channel.ID,
		Message: message,
		Ts:      request.OriginalMessage.Timestamp,
	}
}

// InteractionCallback utilities

// TimeStamp extracts timestamp from the original message
// When the original message is from a thread, we need to post to the same thread
// Below logic checks if the incoming message is from a thread
func TimeStamp(request slack.InteractionCallback) string {
	ts := request.OriginalMessage.ThreadTimestamp
	if ts == "" {
		ts = request.MessageTs
	}
	return ts
}

// Responses is a helper wrapper that allows easier construction of the list of notifications.
func Responses(r ...models.PlatformSimpleNotification) []models.PlatformSimpleNotification {
	return r
}

// ClearOriginalMessage creates a notification that will clear the request message
func ClearOriginalMessage(request slack.InteractionCallback) []models.PlatformSimpleNotification {
	return Responses(InteractionCallbackOverrideRequestMessage(request, ""))
}
