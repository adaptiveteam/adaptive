package models

import (
	"encoding/json"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/nlopes/slack/slackevents"
)

// LambdaRequestJSONUnmarshalUnsafe - unmarshalls LambdaRequest
func LambdaRequestJSONUnmarshalUnsafe(jsMessage string, namespace string) LambdaRequest {
	res, err := LambdaRequestJSONUnmarshal(jsMessage)
	core.ErrorHandler(err, namespace, "LambdaRequest not unmarshaled from '" + jsMessage + "'")
	return res
}

// LambdaRequestJSONUnmarshal - unmarshalls LambdaRequest
func LambdaRequestJSONUnmarshal(jsMessage string) (res LambdaRequest, err error) {
	err = json.Unmarshal([]byte(jsMessage), &res)
	return res, err
}

// ToJSON returns json string
func (l *LambdaRequest) ToJSON() (string, error) {
	b, err := json.Marshal(&l)
	return string(b), err
}

// ToJSONUnsafe returns json string and panics in case of any errors
func (l *LambdaRequest) ToJSONUnsafe(namespace string) string {
	str, err := l.ToJSON()
	core.ErrorHandler(err, namespace, "LambdaRequest failed to marshal")
	return str
}

// NamespacePayloadJSONUnmarshalUnsafe - unmarshalls NamespacePayload
func NamespacePayloadJSONUnmarshalUnsafe(jsMessage string, namespace string) NamespacePayload {
	res, err := NamespacePayloadJSONUnmarshal(jsMessage)
	core.ErrorHandler(err, namespace, "NamespacePayload not unmarshaled from '" + jsMessage + "'")
	return res
}

// NamespacePayloadJSONUnmarshal - unmarshalls NamespacePayload
func NamespacePayloadJSONUnmarshal(jsMessage string) (res NamespacePayload, err error) {
	err = json.Unmarshal([]byte(jsMessage), &res)
	return res, err
}

// ToJSON returns json string
func (np *NamespacePayload) ToJSON() (string, error) {
	b, err := json.Marshal(&np)
	return string(b), err
}

// ToJSONUnsafe returns json string and panics in case of any errors
func (np *NamespacePayload) ToJSONUnsafe(namespace string) string {
	str, err := np.ToJSON()
	core.ErrorHandler(err, namespace, "NamespacePayload failed to marshal")
	return str
}

// ParseGwSNSRequest parses EventsAPIEvent
// deprecated. Use NamespacePayload.ParseEventsAPIEventUnsafe as it is named more consistently
func (np NamespacePayload)ParseGwSNSRequest(namespace string) slackevents.EventsAPIEvent {
	return np.ParseEventsAPIEventUnsafe(namespace)
}

// ParseEventsAPIEventUnsafe - extracts and parses the payload from NamespacePayload
func (np NamespacePayload)ParseEventsAPIEventUnsafe(namespace string) slackevents.EventsAPIEvent {
	res, err := slackevents.ParseEvent(json.RawMessage(np.Payload), slackevents.OptionNoVerifyToken()) 
	core.ErrorHandler(err, namespace, "Could not parse the event: " + np.Payload)
	return res
}
