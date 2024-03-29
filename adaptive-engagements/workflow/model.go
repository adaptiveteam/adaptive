package workflow

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

// A link to a next state should be specified outside of the step
// handler. As this link is encoded inside `mc` message callback,
// we should have an encapsulated version of this `mc`. 
// (1) For instance,
// we can pass `path` to the handler, so that it'll be able to append
// it's internal details. And handler might take a few paths, if
// there is a branching in the step.
// (2) A better alternative is to avoid `path`s inside the handler and
// keep routing logic completely outside of the ordinary handler.
// Routing could be implemented by a function from dialog submission to
// state name. Each state will have an associated handler.

// State - is the state name
type State string
// InitialState might be used as the very first state of the workflow
const InitialState State = ""
// DoneState can be used to represent the final state without events
const DoneState State = "done"

// Event is a high level representation of an event.
// This is parsed from slack input.
type Event string
// MenuDummyEvent is a marker for event parser.
// In this case the real event is extracted from a different place.
const MenuDummyEvent Event = "menu"
// DialogDummyEvent is a marker for event parser.
// In this case the real event is either `submit` or `cancel`
const DialogDummyEvent Event = "dialog"
// SelectorDummyEvent is a marker for event parser.
// The actual event is determined by the name of the control
const SelectorDummyEvent Event = "selector"
// DialogSubmittedEvent is an event when dialog is actually submitted
const DialogSubmittedEvent Event = "submit"
// DialogCancelledEvent is an event when dialog is cancelled
const DialogCancelledEvent Event = "cancel"
// MessageLevelDummyEvent is an event that is attahed to the message as a whole.
// It's then converted to events from individual buttons.
const MessageLevelDummyEvent Event = "message-level-event"
//ImmediateDummyEvent denotes an event that happens immediately after another handler
const ImmediateDummyEvent Event = "auto"
// ImmediateErrorEvent is triggered when there is an error during event handling
const ImmediateErrorEvent Event = "error"
// MenuDummyActionName is an inner event
const MenuDummyActionName Event = "menu-action"
// DefaultEvent is the event that starts the workflow
const DefaultEvent Event = ""

// Data - workflow data. It'll be added to query.
// NB! Do not use State and Event keys, because they are reserved.
type Data map[string]string
// StateFieldName -
const StateFieldName = "State"
// EventFieldName -
const EventFieldName = "Event"
// DataFieldName -this is actually not used but is reserved for future needs.
const DataFieldName = "Data"
// DialogMessageTsFieldName is being used to save time stamp of the message that triggered the dialog.
// It'll be used as message id.
const DialogMessageTsFieldName = "DialogMessageTs"
// IsPermanentMessageFieldName is used to prevent the deletion of the message.
const IsPermanentMessageFieldName = "IsPermanentMessage"

// ReservedKeys is a list of field names that are used by workflow engine
// easy to test contains
var ReservedKeys = map[string]struct{}{
	StateFieldName: struct{}{}, 
	EventFieldName: struct{}{},
	DataFieldName: struct{}{},
	DialogMessageTsFieldName: struct{}{},
	IsPermanentMessageFieldName: struct{}{},
}

// TemplateID identifies the template of the workflow. 
// We use it to select the correct workflow
type TemplateID string

// Instance describes a running instance of a workflow.
// It's very well serializable, doesn't contain any code.
// not used at the moment. Currently we save `data` directly in callback id
type Instance struct {
	ID string
	TemplateID
}
// EventData contains the detailed information about event
type EventData struct{
	State
	Event
	Data
	DialogMessageTs string
	IsOriginalPermanent bool
}
// Arrow is a data structure that captures a connection between states.
type Arrow struct {
	State
	Event
	Next State
}

// Handler handles the incoming event.
type Handler = func (ctx EventHandlingContext) (EventOutput, error)

// Router allows to switch to a different state based on arbitrary conditions
type Router  = func (ctx EventHandlingContext) State

// Template contains the complete description of the workflow
type Template struct {
	Init State
	FSA map[struct{State; Event}] Handler
	Parser SlackEventParser
}

// SlackEventParser analyses the slack input 
type SlackEventParser = func (np models.NamespacePayload4) (EventData, error)

// SimpleHandler overrides NextState with predefined state
func SimpleHandler(handler Handler, nextState State) Handler {
	return func(ctx EventHandlingContext) (out EventOutput, err error) {
		out, err = handler(ctx)
		out.NextState = nextState
		return
	}
}

// Nop does nothing
var Nop Handler = func(ctx EventHandlingContext) (out EventOutput, err error) {
	return
}

// NoOpHandler just moves to the provided state
func NoOpHandler(nextState State) Handler {
	return func(ctx EventHandlingContext) (out EventOutput, err error) {
		out.NextState = nextState
		return
	}
}

// ShowData renders Data as a compact string. For logging purposes.
func ShowData(d Data) (res string) {
	res = "{"
	for key, value := range d {
		res = res + fmt.Sprintf("%s:%s;", key, value)
	}
	res = res + "}"
	return
}
