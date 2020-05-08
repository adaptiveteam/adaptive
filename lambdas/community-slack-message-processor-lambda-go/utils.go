package lambda

import (
	"github.com/adaptiveteam/adaptive/pagination"
	"github.com/pkg/errors"
	"encoding/json"
	"fmt"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	esmodels "github.com/adaptiveteam/adaptive/engagement-scheduling-models"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"strconv"
)

func channelMembersWithPagination(api *slack.Client, channelID string) (users []string, err error) {
	var all pagination.InterfaceSlice
	params := slack.GetUsersInConversationParameters{ChannelID: channelID}
	all, err = SlackGetUsersInConversationPager(api, params).
		Drain()
	users = all.AsStringSlice()
	return 
}

func slackChannelMembers(channelID string, teamID models.TeamID) (users []string) {
	platformToken, err2 := platform.GetToken(teamID)(connGen.ForPlatformID(teamID.ToPlatformID()))
	core.ErrorHandler(err2, "channelMembers", "Could not obtain token")
	api := slack.New(platformToken)
	users, err2 = channelMembersWithPagination(api, channelID)
	core.ErrorHandler(err2, "channelMembers", "Could not GetUsersInConversationParameters for channelID=" + channelID)
	return
}

func writeEngagement(eng models.UserEngagement) {
	err := d.PutTableEntry(eng, engagementTable)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not write to %s table for a new user", engagementTable))
}

func ParseAsAppMentionEvent(apiEvent slackevents.EventsAPIEvent) *slackevents.AppMentionEvent {
	switch apiEvent.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		slackMsg := apiEvent.InnerEvent.Data.(*slackevents.AppMentionEvent)
		return slackMsg
	}
	panic(errors.New("ParseAsAppMentionEvent: Couldn't cast to slackevents.AppMentionEvent"))
}

func DeleteOriginalEng(userId, channel, ts string) {
	utils.DeleteOriginalEng(userId, channel, ts, func(notification models.PlatformSimpleNotification) {
		publish(notification)
	})
}

func concatAppend(slices ...[]esmodels.CrossWalk) []esmodels.CrossWalk {
	var tmp []esmodels.CrossWalk
	for _, s := range slices {
		tmp = append(tmp, s...)
	}
	return tmp
}

func option(name string, label ui.PlainText) ebm.MenuOption {
	return ebm.MenuOption{
		Text:  string(label),
		Value: name,
	}
}

func simpleOption(o ui.PlainText) ebm.MenuOption {
	return simpleOptionStr(string(o))
}

func simpleOptionStr(o string) ebm.MenuOption {
	return ebm.MenuOption{
		Text:  o,
		Value: o,
	}
}

func optionGroup(title ui.PlainText, options ...ebm.MenuOption) ebm.MenuOptionGroup {
	return ebm.MenuOptionGroup{
		Text:    string(title),
		Options: options,
	}
}

func selectControl(name string, label ui.PlainText, options []models.KvPair) ebm.AttachmentActionTextElement {
	return ebm.AttachmentActionTextElement{
		Label:    string(label),
		Name:     name,
		ElemType: models.MenuSelectType,
		Options:  utils.AttachActionElementOptions(options),
		Value:    name,
	}
}

func callback(userId, topic, action string) models.MessageCallback {
	// We are writing month rather than year in engagement because quarter can always be inferred from month
	year, month := core.CurrentYearMonth()
	mc := models.MessageCallback{Module: "community", Source: userId, Topic: topic, Action: action, Target: "", Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	return mc
}

func publish(msg models.PlatformSimpleNotification) {
	models.PublishToSNS(s, msg, platformNotificationTopic, namespace)
}

func respond(teamID models.TeamID, response platform.Response) {
	fmt.Printf("Respond(,%v): %s", response, response.Type)
	presp := platform.TeamResponse{
		TeamID: teamID,
		Response:   response,
	}
	_, err := s.Publish(presp, platformNotificationTopic)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not publish message to %s topic", platformNotificationTopic))
}

func directMessageToUser(teamID models.TeamID, userID string, message platform.MessageContent) {
	respond(teamID,
		platform.Post(platform.ConversationID(userID), message),
	)
}

func simpleMessage(text ui.RichText) platform.MessageContent {
	return platform.MessageContent{Message: text}
}
func replyReplace(request slack.InteractionCallback, teamID models.TeamID, message platform.MessageContent) {
	response := platform.OverrideByURL(
		platform.ResponseURLMessageID{ResponseURL: request.ResponseURL},
		message)
	respond(teamID, response)
}
func deleteOriginalMessage(request slack.InteractionCallback, teamID models.TeamID) {
	response := platform.DeleteByResponseURL(request.ResponseURL)
	respond(teamID, response)
}

func replyInThread(inputMessage slackevents.AppMentionEvent, teamID models.TeamID, body platform.MessageContent) {
	response := platform.PostToThread(platform.ThreadID{ThreadTs: inputMessage.TimeStamp}, body)
	respond(teamID, response)
}

func deleteRequestMessage(request slack.InteractionCallback) platform.Response {
	return platform.Delete(platform.TargetMessageID{
		ConversationID: platform.ConversationID(request.Channel.ID),
		Ts:             request.MessageTs,
	})
}

// COMMON //
func CommentsSurvey(title ui.PlainText, elemLabel ui.PlainText, elemName string) ebm.AttachmentActionSurvey {
	return utils.AttachmentSurvey(string(title), []ebm.AttachmentActionTextElement{
		ebm.NewTextArea(elemName, elemLabel, CommentsSurveyPlaceholder, ""),
	})
}

// Responses

func overrideRequestMessageAsUser(request slack.InteractionCallback, text ui.RichText) models.PlatformSimpleNotification {
	message := utils.InteractionCallbackOverrideRequestMessage(request, string(text))
	return message
}
func simpleResponseAsUser(request slack.InteractionCallback, text ui.RichText) models.PlatformSimpleNotification {
	message := utils.InteractionCallbackSimpleResponse(request, string(text))
	message.AsUser = true
	return message
}
func overrideOriginalMessageAsUser(request slack.InteractionCallback, text ui.RichText) models.PlatformSimpleNotification {
	message := overrideRequestMessageAsUser(request, text)
	message.Ts = request.OriginalMessage.Timestamp
	return message
}

func wrapError(err error, name string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("{%s: %v}", name, err)
}

func invokeUserSetupLambdaUnsafe(userEngage models.UserEngage) {
	invokeLambdaUnsafe(userSetupLambda, userEngage)
}

func invokeLambdaUnsafe(lambdaName string, userEngage models.UserEngage) {
	engageBytes, err := json.Marshal(userEngage)
	core.ErrorHandler(err, namespace, "Could not marshal UserEngage")
	_, err = lambdaAPI.InvokeFunction(lambdaName, engageBytes, false)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not invoke %s", lambdaName))
}

// Lift turns function that works on items into a function that converts slices.
func Lift(f func(interface{}) interface{}) func([]interface{}) []interface{} {
	return func(input []interface{}) (output []interface{}) {
		for _, i := range input {
			output = append(output, f(i))
		}
		return
	}
}

// constructs message with a few options and cancel button.
func selectOptionsMessage(mc models.MessageCallback,
	title ui.PlainText,
	text ui.PlainText,
	fallbackText ui.PlainText,
	opts []ebm.MenuOption,
) platform.MessageContent {
	callbackID := mc.ToCallbackID()

	// Get all subscribed communities
	attachAction1, _ := eb.NewAttachmentActionBuilder().
		Name(mc.Action).
		Text(string(text)).
		ActionType(ebm.AttachmentActionTypeButton).
		Value(callbackID).
		ActionType(ebm.AttachmentActionTypeSelect).
		Confirm(confirm).
		Options(opts).
		Build()

	attachAction2, _ := eb.NewAttachmentActionBuilder().
		Name(ActionCancellationAction).
		Text(string(ActionCancellationText)).
		ActionType(ebm.AttachmentActionTypeButton).
		Value(callbackID).
		Style(models.DangerColor).
		ActionType(models.ButtonType).
		Confirm(confirm).
		Build()

	attach, _ := eb.NewAttachmentBuilder().
		Title(string(title)).
		Text("").
		Fallback(string(fallbackText)).
		Color(models.BlueColorHex).
		AttachmentType(models.DefaultAttachmentType).
		CallbackId(callbackID).
		Actions([]ebm.AttachmentAction{*attachAction1, *attachAction2}).
		Build()
	return platform.MessageContent{
		Attachments: []ebm.Attachment{*attach},
	}
}
