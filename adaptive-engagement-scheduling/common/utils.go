package common

import (
	"fmt"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"

	acfn "github.com/adaptiveteam/adaptive/adaptive-check-function-names"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	aug "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	EngagementEmptyCheck = aug.UserEngagementCheckWithValue{}
	IDOCreateCheck       = aug.UserEngagementCheckWithValue{
		CheckIdentifier: acfn.IDOsExistForMe,
		CheckValue:      false,
	}
)

func PostToCommunity(community community.AdaptiveCommunity, userId, message string) {
	commChannel := communityChannel(userId, community)
	Publish(aug.PlatformSimpleNotification{UserId: userId, Channel: commChannel, Message: core.TextWrap(message, core.Underscore)})
}

func PostToUser(userId string, attachs []ebm.Attachment) {
	Publish(aug.PlatformSimpleNotification{UserId: userId, AsUser: true, Attachments: attachs})
}

func SimpleAttachment(title, text string) *ebm.Attachment {
	return utils.ChatAttachment(title, text, core.EmptyString, core.Uuid(),
		[]ebm.AttachmentAction{}, []ebm.AttachmentField{}, time.Now().Unix())
}

// Publish a message to SNS topic
func Publish(msg aug.PlatformSimpleNotification) {
	_, err := S.Publish(msg, PlatformNotificationTopic)
	core.ErrorHandler(err, Namespace, fmt.Sprintf("Could not pusblish message to %s topic", PlatformNotificationTopic))
}

func globalConnection(teamID models.TeamID) daosCommon.DynamoDBConnection {
	return daosCommon.DynamoDBConnection{
		Dynamo:     D,
		ClientID:   ClientID,
		PlatformID: teamID.ToPlatformID(),
	}
}

func communityChannel(userID string, community community.AdaptiveCommunity) string {
	teamID, err2 := platform.GetTeamIDForUser(D, ClientID, userID)
	core.ErrorHandler(err2, Namespace, fmt.Sprintf("Could not get TeamID for userID %s", userID))
	// Querying for admin community
	params := map[string]*dynamodb.AttributeValue{
		"id":          daosCommon.DynS(string(community)),
		"platform_id": daosCommon.DynS(string(teamID.ToString())),
	}
	var comm aug.AdaptiveCommunity
	err3 := D.GetItemFromTable(CommunitiesTable, params, &comm)
	core.ErrorHandler(err3, Namespace, fmt.Sprintf("Could not query %s table", CommunitiesTable))
	return comm.ChannelID
}

// UserDN renders user id as Slack markup for displaying user name.
func UserDN(userID string) string {
	return "<@" + userID + ">"
}
