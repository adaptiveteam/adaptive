package common

import (
	"fmt"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"

	acfn "github.com/adaptiveteam/adaptive/adaptive-check-function-names"
	ac "github.com/adaptiveteam/adaptive/adaptive-checks"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	aug "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/checks"
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
	ProductionProfile = checks.CheckFunctionMap{
		// Feedback
		acfn.FeedbackGivenThisQuarter: ac.FeedbackGivenForTheQuarter,
		acfn.FeedbackForThePreviousQuarterExists: ac.FeedbackForThePreviousQuarterExists,
		acfn.InLastMonthOfQuarter: func(userID string, date business_time.Date) (rv bool) {
			return date.GetMonth()%3 == 0
		},

		// Community membership
		acfn.InCapabilityCommunity: ac.InCapabilityCommunity,
		acfn.InValuesCommunity:     ac.InCompetenciesCommunity,
		acfn.InHRCommunity:         ac.InHRCommunity,
		acfn.InStrategyCommunity:   ac.InStrategyCommunity,
		acfn.InInitiativeCommunity: ac.InitiativeCommunityExistsForMe,

		// Component existence

		// Miscellaneous
		acfn.UserSettingsExist: func(userID string, date business_time.Date) (rv bool) {
			return true
		},
		acfn.HolidaysExist:                                     ac.HolidaysExist,
		acfn.CoacheesExist:                                     ac.CoacheesExist,
		acfn.AdvocatesExist:                                    ac.AdvocatesExist,
		acfn.CollaborationReportExists:                         ac.ReportExists,
		acfn.UndeliveredEngagementsExistForMe:                  ac.UndeliveredEngagementsExistForMe,
		acfn.UndeliveredEngagementsOrPostponedEventsExistForMe: ac.UndeliveredEngagementsOrPostponedEventsExistForMe,
		acfn.CanBeNudgedForIDO:                                 ac.CanBeNudgedForIDOCreation,

		// Strategy component existence independent of the user
		acfn.TeamValuesExist:     ac.TeamValuesExist,
		acfn.CompanyVisionExists: ac.CompanyVisionExists,
		acfn.ObjectivesExist:     ac.ObjectivesExist,
		acfn.InitiativesExist:    ac.InitiativesExistInMyCapabilityCommunities,

		// Strategy component existence for a given user
		acfn.IDOsExistForMe:        ac.IDOsExistForMe,
		acfn.ObjectivesExistForMe:  ac.ObjectivesExistForMe,
		acfn.InitiativesExistForMe: ac.InitiativesExistForMe,

		// Stale components exist for a specific individual
		acfn.StaleIDOsExistForMe:        ac.StaleIDOsExist,
		acfn.StaleInitiativesExistForMe: ac.StaleInitiativesExistForMe,
		acfn.StaleObjectivesExistForMe:  ac.StaleObjectivesExistForMe,

		// Community existence
		acfn.CapabilityCommunityExists: ac.CapabilityCommunityExists,
		// TODO: A doubt here
		acfn.MultipleCapabilityCommunitiesExists: ac.MultipleCapabilityCommunitiesExists,
		acfn.InitiativeCommunityExists:           ac.InitiativeCommunityExistsForMe,
		// TODO: Implement this
		acfn.MultipleInitiativeCommunitiesExists: func(userID string, date business_time.Date) (rv bool) {
			return false
		},

		// State of community
		acfn.ObjectivesExistInMyCapabilityCommunities:  ac.ObjectivesExistInMyCapabilityCommunities,
		acfn.InitiativesExistInMyCapabilityCommunities: ac.InitiativesExistInMyCapabilityCommunities,
		acfn.InitiativesExistInMyInitiativeCommunities: ac.InitiativesExistInMyInitiativeCommunities,
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
