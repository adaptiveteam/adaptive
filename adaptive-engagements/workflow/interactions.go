package workflow

import (
	"time"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

// This file contains some interactions that workflow can use to
// request information from user.
// This model is more limited than we have for arbitrary slack interactions.
// The core difference is that these interactions do
// not include callback id.
// The callback id is appended by workflow engine.

// SurveySubmitted event of a dialog being submitted by user
const SurveySubmitted Event = "submit"
// SurveyCancelled when user pressed "Cancel"
const SurveyCancelled Event = "cancel"
// Timeout in case we don't get response back in a reasonable time.
const Timeout   Event = "timeout"
// PassiveMessage represents informational part of the interaction. 
// It doesn't include any actions or interactive elements.
type PassiveMessage struct {
	Text ui.RichText
	AttachmentText ui.RichText
	Color string
	Pretext ui.RichText
	Fields []ebm.AttachmentField
	Footer ebm.AttachmentFooter
	// OverrideOriginal is a flag to override message that triggered the event.
	// If two messages have this flag, then all but the last one will be overridden 
	// by the last one.
	OverrideOriginal bool
	// IsPermanentMessage replaces the old mechanism for not-deleting messaged.
	// Now we can prevent the deletion during creation of the message.
	IsPermanentMessage bool
}

// SimpleAction is a specialized version of menu option to support events 
// directly.
// The same structure is being used for buttons.
type SimpleAction struct {
	Label ui.PlainText
	Event
	RequiresConfirmation bool // only works for buttons
}

// Survey shows dialog to the user
// There are two events from Survey - Submitted and Cancelled
type Survey struct {
	ebm.AttachmentActionSurvey
}

// OptionGroup is part of an advanced menu
type OptionGroup struct {
	Title ui.PlainText
	Options []SimpleAction
}

// DefaultDropDownTitle -
const DefaultDropDownTitle ui.PlainText = "Choose ..." 

// SimpleMenu contains menu definition that will be shown to 
// user directly in the chat space. Each menu option will be a separate event
// if OptionGroups are not empty, then `options` are being ignored
type SimpleMenu struct {
	DropDownTitle ui.PlainText // If empty -> "Choose..."
	OptionGroups []OptionGroup
	Options []SimpleAction
	// EnableIgnoreAction adds a predefined button to ignore the menu
	EnableIgnoreAction bool
}

// SelectorOption is one option. It's value will be available in `Value`
type SelectorOption struct {
	Label ui.PlainText
	Value string
}

// Selector is a control that has a few options all associated with the same event.
type Selector struct {
	DropDownTitle ui.PlainText // If empty -> "Choose..."
	Event
	Options []SelectorOption
}
// type ThreadInteraction struct {
// 	ChatSpaceInteraction Interaction
// 	InThreadInteractions []Interaction
// }
type InteractiveElementType string
const ButtonElementType InteractiveElementType = "button"
const DataSelectionElementType InteractiveElementType = "select"
const MenuElementType InteractiveElementType = "menu"

// InteractiveElement is one of button, data selector, or menu
// the difference between data selector and menu is that 
// each item in menu results in a separate event, like a separate button,
// while selector emits only one event and it's possible to extract the data value.
type InteractiveElement struct {
	InteractiveElementType
	Selector
	SimpleMenu
	Button SimpleAction
}
// InteractiveMessage is a single message that contains some interactive elements
// Thread contains messages that should be placed in the thread associated with this 
// message.
// Nested thread is ignored.
type InteractiveMessage struct {
	PassiveMessage
	InteractiveElements []InteractiveElement
	Thread []InteractiveMessage
	// DataOverride might be used to override some of the values of EventOutput.Data
	DataOverride Data
}

// TriggerImmediateEventForAnotherUser creates an event for that user
type TriggerImmediateEventForAnotherUser struct {
	UserID string
	ActionPath models.ActionPath // workflow that will start for that user.
}
// PostponeEventForAnotherUser is a mechanism to pass some information to another user
// when it is convenient to that user.
type PostponeEventForAnotherUser struct {
	UserID string
	ActionPath models.ActionPath // workflow that will start for that user.
	ValidThrough time.Time // Last moment when this event is still valid.
}
// Interaction represents possible interactions that a workflow step
// can use.
type Interaction struct {
	// Responses should only be used to sends 
	// notifications to user/users/channels/threads.
	// it shouldn't be used for any sort of interactions.
	// These responses are sent regardless of the interaction type.
	Responses []platform.Response
	// OptionalSurvey may contain 0..1 surveys. 
	OptionalSurvey []Survey
	// Messages - a few interactive messages.
	Messages []InteractiveMessage
	// A collection of immediate events that could be sent to another users
	ImmediateEvents []TriggerImmediateEventForAnotherUser
	// A collection of postponed events that will be triggered later when
	// it is convenient to that user.
	PostponedEvents []PostponeEventForAnotherUser
	// effectively - do not delete original. Default is false that means - delete original.
	KeepOriginal bool
}

// MenuOption constructs SimpleAction
func MenuOption(event Event, label ui.PlainText) SimpleAction {
	return SimpleAction{
		Label: label,
		Event: event,
	}
}

// Group constructs OptionGrOptionGroupoup
func Group(title ui.PlainText, options ... SimpleAction) (out OptionGroup) {
	return OptionGroup{
		Title: title,
		Options: options,
	}
}

// MenuMessage shows a menu to the user as a standalone message.
// It also enables "ignore" button
func MenuMessage(title ui.RichText, options ... SimpleAction) (out Interaction) {
	return Interaction{
		Messages: []InteractiveMessage{
			{
				PassiveMessage: PassiveMessage{
					AttachmentText: title,
				},
				InteractiveElements: InteractiveElements(
					InteractiveElement{
						InteractiveElementType: MenuElementType,
						SimpleMenu: SimpleMenu{
							Options: options,
							EnableIgnoreAction: true,
						},
					},
				),
			},
		},
	}
}

// SimpleResponses wraps the list of responses into an Interaction
// This is only used as a last step 
func SimpleResponses(responses ...platform.Response) (out Interaction) {
	return Interaction{
		Responses: responses,
	}
}

// OpenSurvey interacts with user using a survey
func OpenSurvey(survey ebm.AttachmentActionSurvey, responses ...platform.Response)(out Interaction) {
	return Interaction{
		OptionalSurvey: []Survey{{AttachmentActionSurvey: survey}},
		Responses: responses,
	}
}

// Button is an interactive element
func Button(event Event, label ui.PlainText) InteractiveElement {
	return InteractiveElement{
		InteractiveElementType: ButtonElementType,
		Button: MenuOption(event, label),
	}
}
// AckButton is a button with confirmation.
func AckButton(event Event, label ui.PlainText) InteractiveElement {
	return InteractiveElement{
		InteractiveElementType: ButtonElementType,
		Button: SimpleAction{
			Label: label,
			Event: event,
			RequiresConfirmation: true,
		},
	}
}
// DataSelector is an interactive element
func DataSelector(event Event, options ... SelectorOption) InteractiveElement {
	return InteractiveElement{
		InteractiveElementType: DataSelectionElementType,
		Selector: Selector{
			Event: event,
		},
	}
}

// InlineMenu creates a drop down control that is placed among buttons.
func InlineMenu(dropDownTitle ui.PlainText, options ... SimpleAction) InteractiveElement {
	return InteractiveElement{
		InteractiveElementType: MenuElementType,
		SimpleMenu: SimpleMenu{
			DropDownTitle: dropDownTitle,
			Options: options,
		},
	}
}

// Selectors is a helper DSL function for CommandButtons
func Selectors(selects ... Selector) (res []InteractiveElement) {
	for _, s := range selects {
		res = append(res, InteractiveElement{
			InteractiveElementType: DataSelectionElementType,
			Selector: s,
		})
	}
	return 
}
// Buttons constructs CommandButtons interaction
func Buttons(title ui.RichText, actions ... InteractiveElement) (out Interaction) {
	return TextWithButtons(title, actions...)
}
// TextWithButtons constructs interaction with a simple text and a few controls
func TextWithButtons(title ui.RichText, actions ... InteractiveElement) (out Interaction) {
	return Interaction{
		Messages: []InteractiveMessage{
			{
				PassiveMessage: PassiveMessage{
					AttachmentText: title,
				},
				InteractiveElements: actions,
			},
		},
	}
}

// InteractiveElements is a helper constructor of a slice of InteractiveElement s
func InteractiveElements(interactiveElements ... InteractiveElement) []InteractiveElement {
	return interactiveElements
}

// InteractiveMessages - slice constructor
func InteractiveMessages(messages ... InteractiveMessage) []InteractiveMessage {
	return messages
}
