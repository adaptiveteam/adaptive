package workflow

import (
	"github.com/slack-go/slack"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// EventHandlingContext is the context for handlers. Contains all relevant information
type EventHandlingContext struct {
	TeamID models.TeamID
	Request slack.InteractionCallback
	Instance
	EventData
	// TargetMessageID contains the identifier of a message.
	// If this is the actual invocation from Slack, it is the id of the incoming message.
	// if it's an "immediate move to the next state", then it's the id of the last
	// interactive message sent to Slack.
	// This ID could be used to override the other message.
	// The reason why we cannot delete the old message is that
	// it might have a thread. So we should at most override the title message.
	platform.TargetMessageID
	// RuntimeData could be used to pass information between immediate state handlers
	RuntimeData map[string]interface{}
}

// Reply sends simple text to the requesting user
func (ctx EventHandlingContext)Reply(text ui.RichText) (out EventOutput) {
	out.NextState = ""
	out.Interaction = SimpleResponses(
		platform.Post(platform.ConversationID(ctx.Request.User.ID),
		platform.MessageContent{Message: text}))
	return 
}

// Prompt sends simple text + a few buttons to the requesting user
func (ctx EventHandlingContext)Prompt(text ui.RichText, interactiveElements ... InteractiveElement) (out EventOutput) {
	out.Interaction.Messages = InteractiveMessages(
		InteractiveMessage{
			PassiveMessage: PassiveMessage{
				AttachmentText: text,
			},
			InteractiveElements: interactiveElements,
		},
	)
	return 
}
// ToggleFlag flips the flag
func (ctx EventHandlingContext)ToggleFlag(flag string) {
	_, isOn := ctx.Data[flag]
	if isOn {
		delete(ctx.Data, flag) // removing "flag"
	} else {
		ctx.Data[flag] = "true" // setting "flag"
	}
}
// SetFlag -
func (ctx EventHandlingContext)SetFlag(flag string, value bool) {
	if value {
		ctx.Data[flag] = "true" // setting "flag"
	} else {
		_, isOn := ctx.Data[flag]
		if isOn {
			delete(ctx.Data, flag) // removing "flag"
		}
	}
}

// GetFlag -
func (ctx EventHandlingContext)GetFlag(flag string) (value bool) {
	_, value = ctx.Data[flag]
	return
}

func (ctx EventHandlingContext) WithRuntimeData(name string, rd interface {}) (out EventHandlingContext) {
	out = ctx
	if out.RuntimeData == nil {
		out.RuntimeData = make(map[string]interface{})
	}
	out.RuntimeData[name] = rd
	return
}

func (ctx EventHandlingContext) TryGetRuntimeData(name string) (rd interface{}, found bool) {
	if ctx.RuntimeData != nil {
		rd, found = ctx.RuntimeData[name]
	}
	return
}