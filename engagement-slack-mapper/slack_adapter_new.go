package engagement_slack_mapper

import (
	"github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/pkg/errors"
	"log"

	"github.com/adaptiveteam/adaptive/daos/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/slack-go/slack"
)
// MessageID identifies message that has been posted.
type MessageID struct {
	ConversationID platform.ConversationID
	Ts string
}

// // PlatformAdapter2 is the main interface to communicate with Slack
// type PlatformAdapter2 interface {
// 	PlatformTokenUnsafe(teamID models.TeamID) string
// 	PostSync(teamResponse platform.TeamResponse) (id MessageID, err error)
// 	PostAsync(teamResponse platform.TeamResponse) chan MessageID
// 	// PostSyncUnsafe(teamResponse platform.TeamResponse) (id MessageID)
// 	ForPlatformID(teamID models.TeamID) PlatformAPI
// }

// PlatformAPI is an API for sending messages to Slack with a known endpoint.
type PlatformAPI interface {
	PostAsync(response platform.Response) chan MessageID
	PostSync(response platform.Response) (id MessageID, err error)
	PostSyncUnsafe(response platform.Response) MessageID
	ShowDialog(survey ebm.AttachmentActionSurvey2) error
}

// // PlatformAdapter2Impl implements PlatformAdapter2
// // NB! Mutable data structure with cache. Should be passed around via reference
// type PlatformAdapter2Impl struct {
// 	conn             common.DynamoDBConnection
// 	tokenCache       map[models.TeamID]string
// }

// // SlackAdapter2 constructs adapter for sending messages to Slack
// func SlackAdapter2(conn common.DynamoDBConnection, teamID models.TeamID) PlatformAPI {
// 	return slack.New(a.PlatformTokenUnsafe(teamID)) &PlatformAdapter2Impl{
// 		conn:       conn,
// 		tokenCache: make(map[models.TeamID]string),
// 	}
// }

// // PlatformTokenUnsafe returns platform token for the given platform id
// func (a *PlatformAdapter2Impl) PlatformTokenUnsafe(teamID models.TeamID) string {
// 	if teamID.IsEmpty() {
// 		panic(errors.New("in PlatformAdapter2 teamID is empty"))
// 	}
// 	if a.tokenCache == nil {
// 		panic(errors.New("PlatformAdapter2 a.tokenCache == nil"))
// 	}
// 	res, ok := a.tokenCache[teamID]
// 	if !ok {
// 		res = a.PlatformTokenDao.GetPlatformTokenUnsafe(teamID)
// 		a.tokenCache[teamID] = res
// 	}
// 	return res
// }

// // PostSync sends message to Slack and returns the identifier of the message.
// func (a *PlatformAdapter2Impl) PostSync(teamResponse platform.TeamResponse) (id MessageID, err error) {
// 	b := a.ForPlatformID(teamResponse.PlatformID)
// 	return b.PostSync(teamResponse.Response)
// }

// // PostAsync sends message to Slack asyncronously.
// func (a *PlatformAdapter2Impl) PostAsync(teamResponse platform.TeamResponse) chan MessageID {
// 	b := a.ForPlatformID(teamResponse.PlatformID)
// 	return b.PostAsync(teamResponse.Response)
// }

// PlatformAPIImpl implements PlatformAPI interface for sending messages to Slack
type PlatformAPIImpl struct {
	API *slack.Client
}

// SlackAdapterForTeamID constructs API for a given teamID
func SlackAdapterForTeamID(conn common.DynamoDBConnection) PlatformAPI {
	teamID := models.ParseTeamID(conn.PlatformID)
	token, err2 := platform.GetToken(teamID)(conn)
	if err2 != nil {
		panic(errors.Wrapf(err2, "SlackAdapterForTeamID"))
	}
	return PlatformAPIImpl{
		API: slack.New(token),
	}
}

// PostSync sends message to Slack and returns the identifier of the message.
func (b PlatformAPIImpl) PostSync(response platform.Response) (id MessageID, err error) {
	switch response.Type {
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
		return MessageID{}, errors.New("Unknown response type: " + string(response.Type))
	}
}

// PostSyncUnsafe posts message and panics in case of errors
func (b PlatformAPIImpl) PostSyncUnsafe(response platform.Response) (id MessageID) {
	id, err2 := b.PostSync(response)
	if err2 != nil {
		log.Panicf("Error posting Response(type = %s) %v to Slack: %+v\n", response.Type, response, err2)
	}
	return id
}

// PostAsync posts message asynchronously and logs error in case of errors
func (b PlatformAPIImpl) PostAsync(response platform.Response) (messageID chan MessageID) {
	messageID = make(chan MessageID, 1)
	core_utils_go.Go("PostAsync", func() {
		msgID := &MessageID{}
		defer func() {
			messageID <- *msgID
			close(messageID)
			if err := recover(); err != nil {
				log.Printf("error engagement-slack-mapper, PlatformAdapter2 %v\n", err)
			}
		}()
		*msgID = b.PostSyncUnsafe(response)
	})
	return
}

// DeleteMessage sends delete message to Slack
func (b PlatformAPIImpl) DeleteMessage(msg platform.DeleteTargetMessage) (id MessageID, err error) {
	channel, ts, err := b.API.DeleteMessage(
		string(msg.TargetMessageID.ConversationID),
		msg.TargetMessageID.Ts,
	)
	return MessageID{ConversationID: platform.ConversationID(channel), Ts: ts}, err
}

// PostToConversation sends a message to Slack
func (b PlatformAPIImpl) PostToConversation(msg platform.PostToConversation) (id MessageID, err error) {
	channel, ts, err := b.API.PostMessage(
		string(msg.ConversationID),
		MessageContentToMsgOptions(msg.Body)...,
	)
	return MessageID{ConversationID: platform.ConversationID(channel), Ts: ts}, err
}

// PostEphemeralToConversation sends delete message to Slack
func (b PlatformAPIImpl) PostEphemeralToConversation(msg platform.PostToUserPrivatelyInConversation) (id MessageID, err error) {
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
func (b PlatformAPIImpl) PostToThread(msg platform.PostToThreadResponsePart) (id MessageID, err error) {
	msgParams := MessageContentToMsgOptions(msg.Body)
	msgParams = append(msgParams, slack.MsgOptionTS(msg.ThreadID.ThreadTs))
	channel, ts, err := b.API.PostMessage(
		string(msg.ThreadID.ConversationID),
		msgParams...,
	)
	return MessageID{ConversationID: platform.ConversationID(channel), Ts: ts}, err
}

// OverrideTargetMessage sends message to Slack that overrides another existing message.
func (b PlatformAPIImpl) OverrideTargetMessage(msg platform.OverrideTargetMessage) (id MessageID, err error) {
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
func (b PlatformAPIImpl) ShowDialog(survey ebm.AttachmentActionSurvey2) (err error) {
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
func (b PlatformAPIImpl) DeleteMessageByURL(msg platform.DeleteMessageByURL) (err error) {
	_, err = b.API.DeleteEphemeral(msg.ResponseURLMessageID.ResponseURL)
	// urlOpt := slack.MsgOptionResponseURL(msg.ResponseURLMessageID.ResponseURL)
	// deleteOriginal := slack.MsgOptionDelete("")
	// _, _, err = b.API.PostMessage("", urlOpt, deleteOriginal)// DeleteEphemeral(msg.ResponseURLMessageID.ResponseURL)
	return
}

// OverrideMessageByURL overrides message using ResponseURL
func (b PlatformAPIImpl) OverrideMessageByURL(msg platform.OverrideMessageByURL) (err error) {
	// urlOpt := slack.MsgOptionResponseURL(msg.ResponseURLMessageID.ResponseURL)
	// _, err = b.API.SendMessage("", urlOpt)
	// 	msg.ResponseURLMessageID.ResponseURL,
	_, err = b.API.SendResponse(msg.ResponseURLMessageID.ResponseURL,
		slack.Msg{
			Attachments:     Attachments(msg.Body.Attachments),
			Text:            string(msg.Body.Message),
			ReplaceOriginal: true,
		})
	return
}
