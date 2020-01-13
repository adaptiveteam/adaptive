package models

import (
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/daos/clientPlatformToken"
)

type PlatformName = clientPlatformToken.PlatformName

// Platform names
const (
	SlackPlatform                      PlatformName = "slack"
	MsTeamsPlatform                    PlatformName = "ms-teams"
	SlackDialogSelectElementLabelLimit              = 48
	SlackDialogSelectElementNameLimit               = 300
)

type PlatformSimpleNotification struct {
	// UserId is used for two purposes:
	// - for obtaining platform token
	// - as an alternative conversation id if Channel is empty
	// TODO: replace with PlatformToken. It's available in request and 
	// we should not spend additional time requesting it from the user-profile lambda.
	UserId string `json:"user_id"`
	// Channel to post the message to. This is used when we have information from earlier context
	// If we don't have any context, we can user 'UserId' above to post to a user
	// TODO: rename it to ConversationID and make mandatory. We should not reuse UserID.
	Channel string `json:"channel"`
	Message string `json:"message"`
	// deprecated. 2019-07-31 https://github.com/adaptiveteam/adaptive-core-lambdas/issues/42
	AsUser bool `json:"as_user"`
	// Timestamp if the message has to be updated
	Ts string `json:"ts"`
	// To post in a thread
	ThreadTs    string             `json:"thread_ts"`
	Attachments []model.Attachment `json:"attachments"`
}
