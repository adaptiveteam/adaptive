package workflow

import (
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
)
// EventOutput contains the result of event handling
type EventOutput struct {
	// Interaction is the response to the original user.
	Interaction
	// TargetedInteractions are interactions that can be sent to another user or community
	TargetedInteractions []TargetedInteraction
	// DataOverride override some of the keys of the context data.
	DataOverride Data
	NextState State
	// ImmediateEvent is a flag/event that allows immediate processing of the next state
	// with the same input. The only difference will be the state and event:
	// state = NextState, event = ImmediateEvent
	// if it's empty, then no immediate processing is triggered
	ImmediateEvent Event
	// RuntimeData could be used to pass information between immediate state handlers
	RuntimeData map[string]interface{}
}
// WithPostponedEvent - appends a PostponedEvent to output
func (eo EventOutput) WithPostponedEvent(events ...PostponeEventForAnotherUser) (out EventOutput) {
	out = eo
	out.PostponedEvents = append(out.PostponedEvents, events...)
	return
}

// WithRuntimeData - sets a RuntimeData to output
func (eo EventOutput) WithRuntimeData(name string, rd interface {}) (out EventOutput) {
	out = eo
	if out.RuntimeData == nil {
		out.RuntimeData = make(map[string]interface{})
	}
	out.RuntimeData[name] = rd
	return
}

// WithNextState - sets the NextState to output
func (eo EventOutput) WithNextState(nextState State) (out EventOutput) {
	out = eo
	out.NextState = nextState
	return
}

// WithInteractiveMessage - adds InteractiveMessages to output
func (eo EventOutput) WithInteractiveMessage(messages ... InteractiveMessage) (out EventOutput) {
	out = eo
	out.Interaction.Messages = append(out.Interaction.Messages, messages ...)
	return
}

// WithCommunityInteraction - adds targeted message to output
func (eo EventOutput) WithCommunityInteraction(targetCommunityID string, messages ... InteractiveMessage) (out EventOutput) {
	out = eo
	out.TargetedInteractions = append(out.TargetedInteractions, 
		TargetedInteraction{
			InteractionTarget: InteractionTarget{
				CommunityID: targetCommunityID,
			},
			Interaction: Interaction{
				Messages: messages,
			},
		},
	)
	return
}


// WithPrependInteractiveMessage - adds InteractiveMessages to output before other messages
func (eo EventOutput) WithPrependInteractiveMessage(messages ... InteractiveMessage) (out EventOutput) {
	out = eo
	out.Interaction.Messages = append(messages, out.Interaction.Messages ...)
	return
}

// WithSurvey appends the survey to the output
func (eo EventOutput) WithSurvey(aaSurvey ebm.AttachmentActionSurvey) (out EventOutput) {
	out = eo
	out.Interaction.OptionalSurvey = append(out.Interaction.OptionalSurvey, Survey{AttachmentActionSurvey: aaSurvey})
	return
}
// WithResponse - appends a Response to output
func (eo EventOutput) WithResponse(responses ... platform.Response) (out EventOutput) {
	out = eo
	out.Responses = append(eo.Responses, responses...)
	return
}
