package engagement_slack_mapper

import (
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"errors"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/nlopes/slack"
	"log"
)

// PlatformAdapter2 is the main interface to communicate with Slack
type PlatformAdapter2 interface {
	PlatformTokenUnsafe(platformID models.PlatformID) string
	PostSync(platformResponse platform.PlatformResponse) (id MessageID, err error)
	PostAsync(platformResponse platform.PlatformResponse) chan MessageID
	// PostSyncUnsafe(platformResponse platform.PlatformResponse) (id MessageID)
	ForPlatformID(platformID models.PlatformID) PlatformAPI 
}

// PlatformAPI is an API for sending messages to Slack with a known endpoint.
type PlatformAPI interface {
	PostAsync(response platform.Response) chan MessageID
	PostSync(response platform.Response) (id MessageID, err error)
	PostSyncUnsafe(response platform.Response) MessageID
	ShowDialog(survey ebm.AttachmentActionSurvey2) error
}
// PlatformAdapter2Impl implements PlatformAdapter2
// NB! Mutable data structure with cache. Should be passed around via reference
type PlatformAdapter2Impl struct {
	PlatformTokenDao platform.DAO
	tokenCache map[models.PlatformID] string
}

// SlackAdapter2 constructs adapter for sending messages to Slack
func SlackAdapter2(platformTokenDao platform.DAO) PlatformAdapter2 {
	return &PlatformAdapter2Impl{
		PlatformTokenDao: platformTokenDao,
		tokenCache: make(map[models.PlatformID] string),
	}
}

// PlatformTokenUnsafe returns platform token for the given platform id
func (a* PlatformAdapter2Impl)PlatformTokenUnsafe(platformID models.PlatformID) string {
	if platformID == "" {
		panic("in PlatformAdapter2 platformID == ''")
	}
	if a.tokenCache == nil {
		panic("PlatformAdapter2 a.tokenCache == nil")
	}
	res, ok := a.tokenCache[platformID]
	if !ok {
		res = a.PlatformTokenDao.GetPlatformTokenUnsafe(platformID)
		a.tokenCache[platformID] = res
	}
	return res
}

// PostSync sends message to Slack and returns the identifier of the message.
func (a* PlatformAdapter2Impl)PostSync(platformResponse platform.PlatformResponse) (id MessageID, err error) {
	b := a.ForPlatformID(platformResponse.PlatformID)
	return b.PostSync(platformResponse.Response)
}
// PostAsync sends message to Slack asyncronously.
func (a* PlatformAdapter2Impl)PostAsync(platformResponse platform.PlatformResponse) chan MessageID {
	b := a.ForPlatformID(platformResponse.PlatformID)
	return b.PostAsync(platformResponse.Response)
}

// PlatformAPIImpl implements PlatformAPI interface for sending messages to Slack
type PlatformAPIImpl struct {
	API *slack.Client
}
// ForPlatformID constructs API for given platform
func (a* PlatformAdapter2Impl)ForPlatformID(platformID models.PlatformID) PlatformAPI {
	return PlatformAPIImpl{
		API: slack.New(a.PlatformTokenUnsafe(platformID)),
	}
}

// PostSync sends message to Slack and returns the identifier of the message.
func (b PlatformAPIImpl) PostSync(response platform.Response) (id MessageID, err error) {
	switch(response.Type){
	case platform.DeleteTargetMessageType:
		return b.DeleteMessage(*response.DeleteTargetMessage)
	case platform.DeleteMessageByURLType:
		err = b.DeleteMessageByURL(*response.DeleteMessageByURL)
		return
	case platform.PostToConversationType:
		return b.PostToConversation(*response.PostToConversation)
	case platform.PostToUserPrivatelyInConversationType:
		return b.PostEphemeralToConversation(*response.PostToUserPrivatelyInConversation)
	case platform.PostToThreadType:
		return b.PostToThread(*response.PostToThreadResponsePart)
	case platform.OverrideTargetMessageType:
		return b.OverrideTargetMessage(*response.OverrideTargetMessage)
	case platform.OverrideMessageByURLType:
		err = b.OverrideMessageByURL(*response.OverrideMessageByURL)
		return 
	default:
		return MessageID{}, errors.New("Unknown response type: "+string(response.Type))
	}
}
// PostSyncUnsafe posts message and panics in case of errors
func (b PlatformAPIImpl) PostSyncUnsafe(response platform.Response) (id MessageID) {
	id, err := b.PostSync(response)
	if err != nil {
		log.Printf("Error posting message %v to Slack: %v\n", response, err)
		panic(err)
	}
	return id
}

// PostAsync posts message asynchronously and logs error in case of errors
func (b PlatformAPIImpl) PostAsync(response platform.Response) (messageID chan MessageID) {
	messageID = make(chan MessageID, 1)
	go func (){
		msgID := &MessageID{}
		defer func(){ messageID <- *msgID
			close(messageID)
			if err := recover(); err != nil {
				log.Printf("error engagement-slack-mapper, PlatformAdapter2 %v\n", err)
			}
		}()
		*msgID = b.PostSyncUnsafe(response)
	}()
	return
}

// DeleteMessage sends delete message to Slack
func (b PlatformAPIImpl)DeleteMessage(msg platform.DeleteTargetMessage) (id MessageID, err error) {
	channel, ts, err := b.API.DeleteMessage(
		string(msg.TargetMessageID.ConversationID), 
		msg.TargetMessageID.Ts,
	)
	return MessageID{ConversationID: platform.ConversationID(channel), Ts: ts}, err
}

// PostToConversation sends a message to Slack
func (b PlatformAPIImpl)PostToConversation(msg platform.PostToConversation) (id MessageID, err error) {
	channel, ts, err := b.API.PostMessage(
		string(msg.ConversationID), 
		MessageContentToMsgOptions(msg.Body)...,
	)
	return MessageID{ConversationID: platform.ConversationID(channel), Ts: ts}, err
}

// PostEphemeralToConversation sends delete message to Slack
func (b PlatformAPIImpl)PostEphemeralToConversation(msg platform.PostToUserPrivatelyInConversation) (id MessageID, err error) {
	options := MessageContentToMsgOptions(msg.Body)
	ephemeral := slack.MsgOptionPostEphemeral(msg.UserID)
	options = append(options, ephemeral)
	channel, ts, err := b.API.PostMessage(
		string(msg.ConversationID), 
		options...,
	)
	return MessageID{ConversationID: platform.ConversationID(channel), Ts: ts}, err
}

// PostToThread sends delete message to Slack
func (b PlatformAPIImpl)PostToThread(msg platform.PostToThreadResponsePart) (id MessageID, err error) {
	msgParams := MessageContentToMsgOptions(msg.Body)
	msgParams = append(msgParams, slack.MsgOptionTS(msg.ThreadID.ThreadTs))
	channel, ts, err := b.API.PostMessage(
		string(msg.ThreadID.ConversationID), 
		msgParams...,
	)
	return MessageID{ConversationID: platform.ConversationID(channel), Ts: ts}, err
}

// OverrideTargetMessage sends message to Slack that overrides another existing message.
func (b PlatformAPIImpl)OverrideTargetMessage(msg platform.OverrideTargetMessage) (id MessageID, err error) {
	msgParams := MessageContentToMsgOptions(msg.Body)
	channel, ts, _, err := // text is ignored
	b.API.UpdateMessage(
		string(msg.TargetMessageID.ConversationID), 
		string(msg.TargetMessageID.Ts), 
		msgParams...,
	)
	return MessageID{ConversationID: platform.ConversationID(channel), Ts: ts}, err
}

// ShowDialog shows dialog in Slack
func (b PlatformAPIImpl)ShowDialog(survey ebm.AttachmentActionSurvey2) (err error) {
	dialog, err := utils.ConvertSurveyToSlackDialog(survey.AttachmentActionSurvey, survey.TriggerID, survey.CallbackID, survey.State)
	if err == nil {
		err = b.API.OpenDialog(survey.TriggerID, dialog)
	}
	return
}

// MessageContentToMsgOptions converts MessageContent to Slack options
func MessageContentToMsgOptions(content platform.MessageContent) []slack.MsgOption {
	return []slack.MsgOption{
		slack.MsgOptionText(string(content.Message), false),
		slack.MsgOptionAsUser(true), 
		slack.MsgOptionAttachments(Attachments(content.Attachments)...),
	}
}

// DeleteMessageByURL sends delete message to Slack using ResponseURL
func (b PlatformAPIImpl)DeleteMessageByURL(msg platform.DeleteMessageByURL) (err error) {
	_, err = b.API.DeleteEphemeral(msg.ResponseURLMessageID.ResponseURL)
	// urlOpt := slack.MsgOptionResponseURL(msg.ResponseURLMessageID.ResponseURL)
	// deleteOriginal := slack.MsgOptionDelete("")
	// _, _, err = b.API.PostMessage("", urlOpt, deleteOriginal)// DeleteEphemeral(msg.ResponseURLMessageID.ResponseURL)
	return
}

// OverrideMessageByURL overrides message using ResponseURL
func (b PlatformAPIImpl)OverrideMessageByURL(msg platform.OverrideMessageByURL) (err error) {
	// urlOpt := slack.MsgOptionResponseURL(msg.ResponseURLMessageID.ResponseURL)
	// _, err = b.API.SendMessage("", urlOpt)
	// 	msg.ResponseURLMessageID.ResponseURL, 
	_, err = b.API.SendResponse(msg.ResponseURLMessageID.ResponseURL, 
			slack.Msg{
			Attachments: Attachments(msg.Body.Attachments),
			Text: string(msg.Body.Message),
			ReplaceOriginal: true,
		})
	return
}
