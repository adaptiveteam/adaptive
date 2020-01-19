package workflow

import (
	"fmt"

	"github.com/pkg/errors"

	//"log"
	"encoding/json"
	"net/url"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/nlopes/slack"
)

// LogInfof signature of log.Infof function
type LogInfof = func(format string, args ...interface{})

// PlatformAPIForPlatformID is a function to obtain API by PlatformID
type PlatformAPIForPlatformID = func(platformID models.PlatformID) mapper.PlatformAPI

// PostponeEvent saves an event to a database for further processing.
// The database will be eventually evaluated for a particular user and
// the event will be triggered.
type PostponeEvent = func(platformID models.PlatformID, postponedEvent PostponeEventForAnotherUser) error

// Environment contains mechanisms to deal with external world
type Environment struct {
	// this is provided from outside as a context. When we want to
	// have a callback routed to our instance, we should prepend this prefix.
	Prefix         models.Path
	GetPlatformAPI PlatformAPIForPlatformID
	LogInfof       LogInfof
	PostponeEvent
}

// MaxImmediateSteps is used to limit possible damage in case of errors
const MaxImmediateSteps = 3

func (w Template) recoverToErrorVar(logInfof LogInfof, err *error) {
	if err2 := recover(); err2 != nil {
		if err != nil {
			logInfof("Before recoverToErrorVar err = %+v", err)
		}
		switch err2.(type) {
		case error:
			err3 := err2.(error)
			err4 := errors.Wrap(err3, "Recover from panic in workflow HandleRequest")
			err = &err4
			logInfof("recoverToErrorVar: %+v", err3)
		case string:
			err3 := err2.(string)
			err4 := errors.New("Recover from string-panic in workflow HandleRequest: " + err3)
			err = &err4
			logInfof("recoverToErrorVar: %+v", err)
		}
	}
}

// Validate checks if the template is correct
func (w Template) Validate() (err error) {
	if w.Parser == nil {
		err = errors.New("Parser is missing")
	} else if w.FSA == nil {
		err = errors.New("FSA is missing")
	}
	return
}

// GetRequestHandler returns RequestHandler that will handle
// the incoming request.
func (w Template) GetRequestHandler(env Environment) RequestHandler {
	return func(actionPath models.RelActionPath, np models.NamespacePayload4) (furtherActions []TriggerImmediateEventForAnotherUser, err error) {
		defer w.recoverToErrorVar(env.LogInfof, &err)
		err = w.Validate()
		if err == nil {
			env.LogInfof("Action path: %s", actionPath.Encode())
			var ctx EventHandlingContext
			ctx, err = w.getEventHandlingContext(np)
			if err == nil {
				furtherActions, err = w.handleContext(env, ctx, MaxImmediateSteps)
			}
			env.LogInfof("HandleRequest completed: %+v", err)
		}
		return
	}
}

// handleContext handles context event
// NB! may be called recursively when receive ImmediateEvent from handler.
func (w Template) handleContext(env Environment,
	ctx EventHandlingContext, i int) (furtherActions []TriggerImmediateEventForAnotherUser, err error) { // Support for the immediate execution. NB! Danger
	key := struct {
		State
		Event
	}{State: ctx.State, Event: ctx.Event}
	h, ok := w.FSA[key]
	if ok {
		env.LogInfof("[U: %s]Before handling %s{S:%s,E:%s,D:%s} -> ",
			ctx.Request.User.ID, env.Prefix.Encode(), ctx.State, ctx.Event,
			ShowData(ctx.Data))
		var out EventOutput
		oldData := copyMap(ctx.Data)
		out, err = h(ctx)
		if err == nil {
			bytes, _ := json.Marshal(out)
			env.LogInfof("[U: %s]After handling: %s", ctx.Request.User.ID, string(bytes))
			data := out.Data
			if data == nil {
				data = ctx.Data
			}
			data = overrideData(data, out.DataOverride)
			var lastMessageID platform.TargetMessageID
			furtherActions = out.Interaction.ImmediateEvents
			lastMessageID, err = interact(ctx, env, out.NextState, out.Interaction, oldData, data)
			if err == nil && out.ImmediateEvent != "" {
				newContext := EventHandlingContext{
					PlatformID: ctx.PlatformID,
					Request:    ctx.Request,
					EventData: EventData{
						State:               out.NextState,
						Data:                data,
						Event:               out.ImmediateEvent,
						IsOriginalPermanent: false, // ???
					},
					TargetMessageID: lastMessageID,
					RuntimeData:     out.RuntimeData,
				}
				if i == 0 {
					env.LogInfof("Not enough iteration count for further looping through states.")
				} else {
					env.LogInfof("Looping (iterations left:%d) to state S:%s,E:%s; ", i, out.NextState, out.ImmediateEvent)
					w.handleContext(env, newContext, i-1)
				}
			}
		}
	} else {
		err = errors.New(fmt.Sprintf("[U: %s]Workflow undefined event %s{S:%s,E:%s}",
			ctx.Request.User.ID, env.Prefix.Encode(), ctx.State, ctx.Event))
	}
	return
}

func (w Template) getEventHandlingContext(np models.NamespacePayload4) (ctx EventHandlingContext, err error) {
	var s EventData
	s, err = w.Parser(np)
	if err == nil {
		ctx = EventHandlingContext{
			PlatformID:      np.PlatformID,
			Request:         np.InteractionCallback,
			EventData:       s,
			TargetMessageID: targetMessageID(np.InteractionCallback),
		}
		if s.DialogMessageTs != "" {
			ctx.TargetMessageID.Ts = s.DialogMessageTs
		}
		if ctx.State == "" {
			ctx.State = w.Init
		}
	}
	return
}

// interact sends output information to user/users
func interact(ctx EventHandlingContext,
	env Environment, nextState State,
	interaction Interaction,
	oldData Data,
	data Data) (lastMessageID platform.TargetMessageID, err error) {
	bytes, _ := json.Marshal(interaction)
	env.LogInfof("interact(..., platformID=%s, next state=%s, interaction=%v, data=%s)",
		ctx.PlatformID, nextState, string(bytes), ShowData(data))
	resps := interaction.Responses
	platformAPI := env.GetPlatformAPI(ctx.PlatformID)
	err = sendResponses(platformAPI, resps...)
	isOriginalPermanentF := ctx.IsOriginalPermanent
	isDeletingOriginal := !interaction.KeepOriginal && !isOriginalPermanentF
	if isDeletingOriginal {
		deleteOriginal(ctx, env, interaction)
	}
	if err == nil {
		for _, m := range interaction.Messages {
			lastMessageID, err = sendInteractiveMessage(ctx,
				env,
				platformAPI,
				nextState, data, m,
				isDeletingOriginal)
		}
		if err == nil {
			for _, survey := range interaction.OptionalSurvey {
				ts := ctx.TargetMessageID.Ts
				if isDeletingOriginal {
					ts = ""
				} // if we deleted the message, then no need to save it's id
				nextStatePath := constructActionPath(env.Prefix, nextState, DialogDummyEvent, ts, isOriginalPermanentF, data)
				dialog := ebm.AttachmentActionSurvey2{
					AttachmentActionSurvey: survey.AttachmentActionSurvey,
					TriggerID:              ctx.Request.TriggerID,
					CallbackID:             nextStatePath.Encode(),
					State:                  "", // Not using at the moment
				}
				err = platformAPI.ShowDialog(dialog)
			}
			if err == nil {
				// interaction.ImmediateEvents { // these events will be returned
				for _, evt := range interaction.PostponedEvents {
					err = env.PostponeEvent(ctx.PlatformID, evt)
				}
			}
		}
	} else {
		env.LogInfof("Error while sending responses: %+v", err)
	}
	return
}

func overrideData(d Data, o Data) (res Data) {
	res = d
	if res == nil {
		res = make(Data)
	}
	if o != nil {
		for k, v := range o {
			res[k] = v
		}
	}
	return
}

func sendInteractiveMessage(ctx EventHandlingContext,
	env Environment,
	platformAPI mapper.PlatformAPI,
	nextState State,
	data Data,
	message InteractiveMessage,
	isDeletingOriginal bool) (lastMessageID platform.TargetMessageID, err error) {
	var msg platform.MessageContent
	msg, err = convertMessage(env.Prefix, nextState, data, message)
	if err == nil {
		var response platform.Response
		if message.OverrideOriginal && ctx.TargetMessageID.Ts != "" && !isDeletingOriginal {
			response = platform.Override(ctx.TargetMessageID, msg)
		} else {
			response = platform.Post(platform.ConversationID(ctx.Request.User.ID), msg)
		}
		var mapperMessageID mapper.MessageID
		mapperMessageID, err = platformAPI.PostSync(response)
		lastMessageID = convertToTargetMessageID(mapperMessageID)
		bytes, _ := json.Marshal(lastMessageID)
		env.LogInfof("lastMessageID = (%v)", string(bytes))
		if err == nil {
			threadID := platform.ThreadID{ThreadTs: lastMessageID.Ts, ConversationID: lastMessageID.ConversationID}
			for _, mt := range message.Thread {
				var msgt platform.MessageContent
				msgt, err = convertMessage(env.Prefix, nextState, data, mt)
				if err != nil {
					break
				}
				response := platform.PostToThread(threadID, msgt)
				err = sendResponses(platformAPI, response)
				if err != nil {
					break
				}
			}
		}
	}
	return
}

// deleteOriginal sends message to delete original message
func deleteOriginal(ctx EventHandlingContext,
	env Environment,
	interaction Interaction) {
	var delete platform.Response
	if ctx.Request.ResponseURL != "" {
		delete = platform.DeleteByResponseURL(ctx.Request.ResponseURL)
	} else {
		delete = platform.Delete(ctx.TargetMessageID)
	}
	err2 := sendResponses(env.GetPlatformAPI(ctx.PlatformID), delete)
	if err2 != nil {
		env.LogInfof("couldn't delete old message: %+v", err2)
	}
}

func sendResponses(api mapper.PlatformAPI,
	responses ...platform.Response,
) (err error) {
	for _, response := range responses {
		_, err = api.PostSync(response)
		if err != nil {
			return
		}
	}
	return
}

func ignoreAction(callbackID string) ebm.AttachmentAction {
	return ebm.AttachmentAction{
		Text:       user.SkipActionLabel,
		ActionType: ebm.AttachmentActionTypeButton,
		Value:      callbackID,
		Name:       string(models.Ignore),
	}
}

// convertMessage only converts the message itself, without thread
func convertMessage(
	prefix models.Path,
	nextState State,
	data Data,
	message InteractiveMessage,
) (res platform.MessageContent, err error) {
	overriddenData := overrideData(data, message.DataOverride)
	callback := constructActionPath(prefix, nextState,
		"unchanged-message-event", "", message.IsPermanentMessage, overriddenData) // this should be overridden in inner methods
	res.Message = message.Text
	attachActions := []ebm.AttachmentAction{}
	for _, el := range message.InteractiveElements {
		switch el.InteractiveElementType {
		case DataSelectionElementType:
			action := convertSelector(callback, el.Selector)
			attachActions = append(attachActions, action)
		case ButtonElementType:
			action := convertButton(callback, el.Button)
			attachActions = append(attachActions, action)
		case MenuElementType:
			actions := convertMenu(callback, el.SimpleMenu)
			attachActions = append(attachActions, actions...)
		default:
			err = errors.New("Malformed interactive element")
		}
	}
	nextStatePath := updateEvent(callback, MessageLevelDummyEvent)

	attach := ebm.Attachment{
		Actions:    attachActions,
		Color:      message.Color,
		Pretext:    string(message.Pretext),
		Text:       string(message.AttachmentText),
		Fields:     message.Fields,
		Footer:     message.Footer,
		CallbackId: nextStatePath.Encode(),
	}
	res.Attachments = []ebm.Attachment{attach}
	return
}
func updateEvent(callback models.ActionPath, event Event) (res models.ActionPath) {
	res = callback
	res.Values.Set(EventFieldName, string(event))
	return
}
func convertSelector(callback models.ActionPath, s Selector) (selectAction ebm.AttachmentAction) {
	nextStatePath := updateEvent(callback, s.Event)
	options := []ebm.MenuOption{}
	for _, opt := range s.Options {
		options = append(options,
			ebm.MenuOption{Text: string(opt.Label), Value: opt.Value})
	}
	var text ui.PlainText
	if s.DropDownTitle == "" {
		text = DefaultDropDownTitle
	} else {
		text = s.DropDownTitle
	}
	selectAction = ebm.AttachmentAction{
		Text:       string(text),
		ActionType: ebm.AttachmentActionTypeSelect,
		Value:      nextStatePath.Encode(),
		Name:       string(s.Event),
		Options:    options,
	}
	return
}

var defaultConfirmation = ebm.AttachmentActionConfirm{
	OkText:      models.YesLabel,
	DismissText: models.CancelLabel,
}

func convertButton(callback models.ActionPath, button SimpleAction) (action ebm.AttachmentAction) {
	nextStatePath := updateEvent(callback, button.Event)
	action = ebm.AttachmentAction{
		Text:       string(button.Label),
		ActionType: ebm.AttachmentActionTypeButton,
		Value:      nextStatePath.Encode(),
		Name:       string(button.Event),
	}
	if button.RequiresConfirmation {
		action.Confirm = defaultConfirmation
		action.Style = ebm.AttachmentActionStyleDanger
	}
	return
}

func convertMenu(callback models.ActionPath,
	menu SimpleMenu) (actions []ebm.AttachmentAction) {
	nextStatePath := updateEvent(callback, MenuDummyEvent)
	callbackID := nextStatePath.Encode()
	options := convertOptions(menu.Options)
	optionGroups := convertOptionGroups(menu.OptionGroups)
	var text ui.PlainText
	if menu.DropDownTitle == "" {
		text = DefaultDropDownTitle
	} else {
		text = menu.DropDownTitle
	}
	selectAction := ebm.AttachmentAction{
		Text:         string(text),
		ActionType:   ebm.AttachmentActionTypeSelect,
		Value:        callbackID,
		Name:         string(MenuDummyActionName),
		Options:      options,
		OptionGroups: optionGroups,
	}
	actions = []ebm.AttachmentAction{selectAction}
	if menu.EnableIgnoreAction {
		actions = append(actions, ignoreAction(callbackID))
	}
	return
}

func convertOptions(menuOptions []SimpleAction) (options []ebm.MenuOption) {
	for _, opt := range menuOptions {
		option := ebm.MenuOption{
			Text:  string(opt.Label),
			Value: string(opt.Event),
		}
		options = append(options, option)
	}
	return
}

func convertOptionGroups(menuOptionGroups []OptionGroup) (optionGroups []ebm.MenuOptionGroup) {
	for _, opt := range menuOptionGroups {
		option := ebm.MenuOptionGroup{
			Text:    string(opt.Title),
			Options: convertOptions(opt.Options),
		}
		optionGroups = append(optionGroups, option)
	}
	return
}

// statePath evaluates the new path for the next state
func statePath(prefix models.Path, st State) models.Path {
	return append(prefix, string(st))
}

func actionPath(request slack.InteractionCallback) models.ActionPath {
	return models.ParseActionPath(request.CallbackID)
}

// ExternalActionPath is a constructor of a path to a certain state/event
func ExternalActionPath(prefix models.Path, state State, event Event) models.ActionPath {
	return ExternalActionPathWithData(prefix, state, event, map[string]string{}, true)
}

// ExternalActionPathWithData is a constructor of a path to a certain state/event
func ExternalActionPathWithData(prefix models.Path, state State, event Event, data Data, isOriginalPermanent bool) models.ActionPath {
	return constructActionPath(prefix, state, event, "", isOriginalPermanent, data)
}

func constructActionPath(prefix models.Path, state State, event Event, dialogMessageTs string, isPermanentMessage bool, data Data) models.ActionPath {
	values := url.Values{}
	for key, value := range data {
		if _, isReserved := ReservedKeys[key]; !isReserved {
			values[key] = []string{value}
		}
	}
	// override state and event fields
	values[StateFieldName] = []string{string(state)}
	values[EventFieldName] = []string{string(event)}
	if dialogMessageTs != "" {
		values[DialogMessageTsFieldName] = []string{dialogMessageTs}
	}

	if isPermanentMessage {
		values[IsPermanentMessageFieldName] = []string{"true"}
	}
	return models.ActionPath{
		Path:   prefix,
		Values: values,
	}
}

// Parser extracts state and event from incoming request
func Parser(np models.NamespacePayload4) (res EventData, err error) {
	ap := actionPath(np.InteractionCallback)
	st, ok := ap.Values[StateFieldName]
	if ok && len(st) > 0 {
		res.State = State(st[0])
	}
	ev, ok := ap.Values[EventFieldName]
	if ok && len(st) > 0 {
		res.Event = Event(ev[0])
	}
	ts, ok := ap.Values[DialogMessageTsFieldName]
	if ok && len(ts) > 0 {
		res.DialogMessageTs = ts[0]
	}
	res.Data = make(Data)
	for key, value := range ap.Values {
		if len(value) > 0 {
			if key == IsPermanentMessageFieldName && value[0] == "true" {
				res.IsOriginalPermanent = true
			}
			if _, ok := ReservedKeys[key]; !ok {
				res.Data[key] = value[0]
			}
		}
	}
	if res.Event == DialogDummyEvent { // /strategy/my-workflow?State=myState&Event=dialog
		switch np.SlackRequest.Type {
		case models.DialogSubmissionSlackRequestType:
			res.Event = DialogSubmittedEvent
		case models.DialogCancellationSlackRequestType:
			res.Event = DialogCancelledEvent
		default:
			err = errors.New("Unexpected request type " + string(np.SlackRequest.Type))
		}
	} else {
		if len(np.InteractionCallback.ActionCallback.AttachmentActions) > 0 {
			action := np.InteractionCallback.ActionCallback.AttachmentActions[0]

			if res.Event == SelectorDummyEvent {
				res.Event = Event(action.Name)
			} else if res.Event == MessageLevelDummyEvent && action.Name != string(MenuDummyActionName) {
				res.Event = Event(action.Name)
			} else {
				if res.Event == MenuDummyEvent ||
					(res.Event == MessageLevelDummyEvent && action.Name == string(MenuDummyActionName)) &&
						len(action.SelectedOptions) > 0 {
					res.Event = Event(action.SelectedOptions[0].Value)
				}
				// if len(action.SelectedOptions) > 0 {
				// 	res.Event = Event(action.SelectedOptions[0].Value)
				// 	if res.Event == MenuDummyEvent {
				// 		res.Event = Event(action.SelectedOptions[0].Value)
				// 	} else if res.Event == MessageLevelDummyEvent {

				// 		res.Event = Event(action.SelectedOptions[0].Value)
				// 	}
				// }
			}
		}
	}
	return
}

// GetConversationID returns ConversationID of the request
func GetConversationID(request slack.InteractionCallback) (conversationID platform.ConversationID) {
	if request.Channel.ID == "" {
		conversationID = platform.ConversationID(request.User.ID)
	} else {
		conversationID = platform.ConversationID(request.Channel.ID)
	}
	return
}

func targetMessageID(request slack.InteractionCallback) platform.TargetMessageID {
	ts := request.MessageTs
	ap := actionPath(request)
	dialogTs, ok := ap.Values[DialogMessageTsFieldName]
	if ok && len(dialogTs) > 0 {
		ts = dialogTs[0]
	}
	return platform.TargetMessageID{
		ConversationID: GetConversationID(request),
		Ts:             ts,
	}
}

func (w Template) fixInitState(ctx *EventHandlingContext) {
	if ctx.State == "" {
		ctx.State = w.Init
	}
}

func convertToTargetMessageID(mid mapper.MessageID) (tmid platform.TargetMessageID) {
	return platform.TargetMessageID{
		ConversationID: mid.ConversationID,
		Ts:             mid.Ts,
	}
}

// func isOriginalPermanent(data Data) (isOriginalPermanent bool) {
// 	var isOriginalPermanentV string
// 	isOriginalPermanentV, isOriginalPermanent = data[IsPermanentMessageFieldName]
// 	isOriginalPermanent = isOriginalPermanent && isOriginalPermanentV == "true"
// 	return
// }

func copyMap(d Data) (res Data) {
	res = make(Data)
	for k, v := range d {
		res[k] = v
	}
	return
}
