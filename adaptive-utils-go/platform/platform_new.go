package platform

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// ConversationID is the identifier of channel where to put new message
// if ChannelID is empty, then UserID is being used and the message is sent to user directly.
type ConversationID string

// TargetMessageID is the identifier of the message
// if ConversationID is either UserID or ChannelID
type TargetMessageID struct {
	ConversationID ConversationID `json:"conversation_id"`
	Ts             string         `json:"ts"`
}

// ResponseURLMessageID identifies message using response url
type ResponseURLMessageID struct {
	ResponseURL string
}

// ThreadID is the identifier of thread where to put new message
type ThreadID struct {
	ConversationID ConversationID `json:"conversation_id"`
	ThreadTs       string         `json:"thread_ts"`
}

// MessageContent is a message we send back to Slack. It only contains the contents.
type MessageContent struct {
	Message     ui.RichText        `json:"message,omitempty"`
	Attachments []model.Attachment `json:"attachments,omitempty"`
}

// ResponseType is a enum that describes different kinds of messages that we send to Slack
type ResponseType string

const (
	// OverrideTargetMessageType overrides the original message with the new one
	OverrideTargetMessageType ResponseType = "override-target"
	OverrideMessageByURLType  ResponseType = "override-by-url"
	// PostToConversationType is posting message to a conversation
	PostToConversationType ResponseType = "post-to-channel"
	// PostToUserPrivatelyInConversationType is posting message to user
	// in a specific way - privetly but directly in a conversation
	PostToUserPrivatelyInConversationType ResponseType = "post-to-user-secretly-in-channel"
	PostToThreadType                      ResponseType = "post-to-thread"
	DeleteTargetMessageType               ResponseType = "delete-target"
	DeleteMessageByURLType                ResponseType = "delete-by-url"
)

// OverrideTargetMessage is the subtype of Response that contains information about
// the target message to override and new message content.
type OverrideTargetMessage struct {
	TargetMessageID TargetMessageID `json:"target_message_id"`
	Body            MessageContent  `json:"body"`
}

// OverrideMessageByURL is the subtype of Response that contains information about
// the target message to override and new message content.
type OverrideMessageByURL struct {
	ResponseURLMessageID ResponseURLMessageID `json:"response_url_message_id"`
	Body                 MessageContent       `json:"body"`
}

// PostToConversation is posting message to conversation
type PostToConversation struct {
	ConversationID ConversationID `json:"conversation_id"`
	Body           MessageContent `json:"body"`
}

// PostToUserPrivatelyInConversation is posting message to user (UserID)
// in a specific way - privately but directly in a conversation
type PostToUserPrivatelyInConversation struct {
	UserID         string         `json:"user_id"`
	ConversationID ConversationID `json:"conversation_id"`
	Body           MessageContent `json:"body"`
}

// PostToThreadResponsePart is posting message to channel
type PostToThreadResponsePart struct {
	ThreadID ThreadID       `json:"thread_id"`
	Body     MessageContent `json:"body"`
}

// DeleteTargetMessage deletes the message
type DeleteTargetMessage struct {
	TargetMessageID TargetMessageID `json:"target_message_id"`
}

// DeleteMessageByURL deletes the ephemeral message
type DeleteMessageByURL struct {
	ResponseURLMessageID ResponseURLMessageID `json:"response_url_message_id"`
}

// Response is a message we send back to Slack.
// It is a replacement for PlatformSimpleNotification, restructuring it according to it's usage.
// There are a few different message types. The sum type is emulated using Golang structs.
// Type is used to distinguish between different
type Response struct {
	Type                              ResponseType                       `json:"type"`
	PostToConversation                *PostToConversation                `json:"post_to_conversation,omitempty"`
	PostToUserPrivatelyInConversation *PostToUserPrivatelyInConversation `json:"post_to_user_privately_in_conversation,omitempty"`
	PostToThreadResponsePart          *PostToThreadResponsePart          `json:"post_to_thread_response_part,omitempty"`
	OverrideTargetMessage             *OverrideTargetMessage             `json:"override_target_message,omitempty"`
	OverrideMessageByURL              *OverrideMessageByURL              `json:"override_message_by_url,omitempty"`

	DeleteTargetMessage *DeleteTargetMessage `json:"delete_target_message,omitempty"`
	DeleteMessageByURL  *DeleteMessageByURL  `json:"delete_message_by_url,omitempty"`
}

// TeamResponse is the lower level message that works at the level
// of communication of Slack and Adaptive. That's why it has PlatformID
// which enables such communication.
type TeamResponse struct {
	TeamID   models.TeamID `json:"platform_id"`
	Response Response      `json:"response"`
}

// Message constructs platform message without information about where to post it.
func Message(text ui.RichText, attachements ...model.Attachment) MessageContent {
	return MessageContent{
		Message:     text,
		Attachments: attachements,
	}
}

// MessageID constructs TargetMessageID from given conversation id and timestamp
func MessageID(conversationID ConversationID, ts string) TargetMessageID {
	return TargetMessageID{
		Ts:             ts,
		ConversationID: conversationID,
	}
}

// Delete constructs Response that will delete the target message.
func Delete(targetMessageID TargetMessageID) Response {
	return Response{
		Type: DeleteTargetMessageType,
		DeleteTargetMessage: &DeleteTargetMessage{
			TargetMessageID: targetMessageID,
		},
	}
}

// DeleteByResponseURL deletes message given response url
func DeleteByResponseURL(responseURL string) Response {
	return Response{
		Type: DeleteMessageByURLType,
		DeleteMessageByURL: &DeleteMessageByURL{
			ResponseURLMessageID: ResponseURLMessageID{ResponseURL: responseURL},
		},
	}
}

// Override constructs Response that will override the target message.
func Override(targetMessageID TargetMessageID, body MessageContent) Response {
	return Response{
		Type: OverrideTargetMessageType,
		OverrideTargetMessage: &OverrideTargetMessage{
			TargetMessageID: targetMessageID,
			Body:            body,
		},
	}
}

// OverrideByURL constructs Response that will override the message identified
// by responseURL.
func OverrideByURL(responseURLMessageID ResponseURLMessageID, body MessageContent) Response {
	return Response{
		Type: OverrideMessageByURLType,
		OverrideMessageByURL: &OverrideMessageByURL{
			ResponseURLMessageID: responseURLMessageID,
			Body:                 body,
		},
	}
}

// Post constructs Response that will simply post a new message into the conversation.
func Post(conversationID ConversationID, body MessageContent) Response {
	return Response{
		Type: PostToConversationType,
		PostToConversation: &PostToConversation{
			ConversationID: conversationID,
			Body:           body,
		},
	}
}

// PostEphemeral constructs Response that will post an ephemeral message
func PostEphemeral(userID string, conversationID ConversationID, body MessageContent) Response {
	return Response{
		Type: PostToUserPrivatelyInConversationType,
		PostToUserPrivatelyInConversation: &PostToUserPrivatelyInConversation{
			UserID:         userID,
			ConversationID: conversationID,
			Body:           body,
		},
	}
}

// PostToThread constructs Response that will post the message into the thread.
func PostToThread(threadID ThreadID, body MessageContent) Response {
	return Response{
		Type: PostToThreadType,
		PostToThreadResponsePart: &PostToThreadResponsePart{
			ThreadID: threadID,
			Body:     body,
		},
	}
}
