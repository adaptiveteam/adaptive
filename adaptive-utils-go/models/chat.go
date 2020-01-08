package models

import (
	"fmt"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"reflect"
	"strconv"
)

type AttachActionName string

// Constant values that are used in engagements
const (
	SlackInChannel        = "in_channel"
	SubmitLabel           = "Submit"
	CancelLabel           = "Cancel"
	EditLabel             = "Edit"
	YesLabel              = "Yes"
	NoLabel               = "No"
	BlueColorHex          = "#3AA3E3"
	DefaultAttachmentType = "default"
	MenuSelectType        = "select"
	ButtonType            = "button"
	OkColor               = "good"
	WarningColor          = "warning"
	DangerColor           = "danger"

	// attachment action variables
	Now    AttachActionName = "now"
	Update AttachActionName = "update"
	Ignore AttachActionName = "ignore"
	Back   AttachActionName = "back"
	Ok     AttachActionName = "ok"
	Select AttachActionName = "select"
	No     AttachActionName = "no"
	Cancel AttachActionName = "cancel"
)

type MessageCallback struct {
	// Module in context
	Module string `json:"module"`
	// Source for the engagement
	Source string `json:"source"`
	// Topic for chat
	Topic string `json:"topic"`
	// Action for the chat
	Action string `json:"action"`
	// Who would be affected by the engagement
	Target string `json:"target"`
	// Current month
	Month string `json:"month"`
	// Current year
	Year string `json:"year"`
}

// DummyMessageCallback creates a callback message that could be customized further.
func DummyMessageCallback(source string) MessageCallback {
	year, month := core.CurrentYearMonth()
	return MessageCallback{
		Module: "dummy_namespace",
		Source: source, Topic: "dummy", Action: "noaction",
		Month: strconv.Itoa(int(month)),
		Year:  strconv.Itoa(year)}
}

// This method sets a field value of MessageCallback
func (m *MessageCallback) Set(field string, value string) {
	v := reflect.ValueOf(m).Elem().FieldByName(field)
	if v.IsValid() {
		v.SetString(value)
	}
}

func (m *MessageCallback) WithModule(module string) *MessageCallback {
	m.Module = module
	return m
}

func (m *MessageCallback) WithSource(source string) *MessageCallback {
	m.Source = source
	return m
}

func (m *MessageCallback) WithTopic(topic string) *MessageCallback {
	m.Topic = topic
	return m
}

func (m *MessageCallback) WithAction(action string) *MessageCallback {
	m.Action = action
	return m
}

func (m *MessageCallback) WithTarget(target string) *MessageCallback {
	m.Target = target
	return m
}

// Sprint deprecated. Use ToCallbackID instead
func (m *MessageCallback) Sprint() string {
	return m.ToCallbackID()
}

// ToCallbackID serializes this object to string that contains all parts.
// The name of this method communicates semantics better than Sprint.
func (m *MessageCallback) ToCallbackID() string {
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", m.Module, m.Source, m.Topic, m.Action, m.Target, m.Month, m.Year)
}

func PublishToSNS(sns *awsutils.SnsRequest, msg PlatformSimpleNotification, platformNotificationTopic, namespace string) {
	_, err := sns.Publish(msg, platformNotificationTopic)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not pusblish message to %s topic", platformNotificationTopic))
}
