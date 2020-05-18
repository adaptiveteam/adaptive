package lambda

import (
	daosUser "github.com/adaptiveteam/adaptive/daos/user"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/pkg/errors"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/communityUser"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/strategyCommunity"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func availableCommunities(teamID models.TeamID) []string {
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	// Get all used communities
	comms, err := adaptiveCommunity.ReadByPlatformID(teamID.ToPlatformID())(conn)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not scan adaptiveCommunity table"))
	var b []string
	for _, each := range comms {
		b = append(b, each.ID)
	}

	return core.InAButNotB(allComms, b)
}

func availableStrategyCommunities(teamID models.TeamID, userID string) []models.KvPair {
	var op []models.KvPair
	var strComms []strategy.StrategyCommunity
	err := d.QueryTableWithIndex(strategyCommunitiesTable, awsutils.DynamoIndexExpression{
		IndexName: string(strategyCommunity.PlatformIDChannelCreatedIndex),
		Condition: "platform_id = :pi AND channel_created = :cc",
		Attributes: map[string]interface{}{
			":pi": teamID.ToString(),
			":cc": 0,
		},
	}, map[string]string{}, true, -1, &strComms)
	if err == nil {
		logger.Infof("Available Strategy communities for Adaptive to join: %v", strComms)
		for _, each := range strComms {
			// Return only those communities for which the user is a co-ordinator for
			if each.Advocate == userID {
				var commName string
				switch each.Community {
				case community.Capability:
					params := map[string]*dynamodb.AttributeValue{
						"id":          dynString(each.ID),
						"platform_id": dynString(teamID.ToString()),
					}
					var capComm strategy.CapabilityCommunity
					err = d.GetItemFromTable(capabilityCommunitiesTable, params, &capComm)
					if err != nil {
						logger.WithField("namespace", namespace).WithField("error", err).
							Errorf(fmt.Sprintf("Could not find in %s table: %v", capabilityCommunitiesTable, params))
					} else {
						commName = capComm.Name
					}
				case community.Initiative:
					params := map[string]*dynamodb.AttributeValue{
						"id":          dynString(each.ID),
						"platform_id": dynString(teamID.ToString()),
					}
					var capComm strategy.StrategyInitiativeCommunity
					err = d.GetItemFromTable(strategyInitiativeCommunitiesTable, params, &capComm)
					if err != nil {
						logger.WithField("namespace", namespace).WithField("error", err).
							Errorf(fmt.Sprintf("Could not find in %s table: %v", strategyInitiativeCommunitiesTable, params))
					} else {
						commName = capComm.Name
					}
				}
				op = append(op, models.KvPair{
					Key:   fmt.Sprintf("[%s] %s", string(each.Community), commName),
					Value: fmt.Sprintf("%s:%s", string(each.Community), each.ID),
				})
			}
		}
	} else {
		logger.WithField("namespace", namespace).WithField("error", err).
			Errorf(fmt.Sprintf("Could not query %s table", strategyCommunity.TableNameSuffixVar))
	}
	return op
}

// StrategyCommunityByID finds community by id. panics if not found
func StrategyCommunityByID(ID string) strategy.StrategyCommunity {
	params := map[string]*dynamodb.AttributeValue{
		"id": dynString(ID),
	}
	var stratComm strategy.StrategyCommunity
	err := d.GetItemFromTable(strategyCommunitiesTable, params, &stratComm)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not find %s in %s table", ID, strategyCommunitiesTable))
	return stratComm
}

func unsetStrategyCommunities(channelID string) {
	var strComms []strategy.StrategyCommunity
	err := d.QueryTableWithIndex(strategyCommunitiesTable, awsutils.DynamoIndexExpression{
		IndexName: strategyCommunitiesChannelIndex,
		Condition: "channel_id = :c",
		Attributes: map[string]interface{}{
			":c": channelID,
		},
	}, map[string]string{}, true, -1, &strComms)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not query %s table", communityUsersTable))
	// For each of the strategy community, unset the channel created flag
	exprAttributes := map[string]*dynamodb.AttributeValue{
		":cc": dynNumber(0),
	}
	for _, each := range strComms {
		key := idParams(each.ID)
		updateExpression := "set channel_created = :cc"
		err = d.UpdateTableEntry(exprAttributes, key, updateExpression, strategyCommunitiesTable)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not update channel_created flag in %s table", strategyCommunitiesTable))
	}
}

// func deleteCommunityTableEntry(ID string, teamID models.TeamID) {
// 	communityDAO.DeleteUnsafe(teamID, ID)
// }

// TODO: Simplify this to be more generic
func StrategyCommunityIdTypeName(val string, teamID models.TeamID) (string, string, string) {
	var err error
	res := strings.Split(val, ":")
	commName := res[0]
	parentID := res[1]
	params := idAndPlatformIDParams(parentID, teamID)

	switch community.AdaptiveCommunity(commName) {
	case community.Capability:
		var capComm strategy.CapabilityCommunity
		err = d.GetItemFromTable(capabilityCommunitiesTable, params, &capComm)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not find parentID=%s in %s table", parentID, capabilityCommunitiesTable))
		return parentID, string(commName), capComm.Name
	case community.Initiative:
		var initComm strategy.StrategyInitiativeCommunity
		err = d.GetItemFromTable(strategyInitiativeCommunitiesTable, params, &initComm)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not find parentID=%s in %s table", parentID, strategyInitiativeCommunitiesTable))
		return parentID, string(commName), initComm.Name
	}
	return core.EmptyString, core.EmptyString, core.EmptyString
}

func dynString(str string) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{S: aws.String(str)}
	return &attr
}

func dynNumber(i int) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{N: aws.String(strconv.Itoa(i))}
	return &attr
}

func idAndPlatformIDParams(id string, teamID models.TeamID) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue{
		"id":          dynString(id),
		"platform_id": dynString(teamID.ToString()),
	}
	return params
}

func idParams(id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue{
		"id": dynString(id),
	}
	return params
}

func subscribedCommunityIDs(channel string, conn daosCommon.DynamoDBConnection) (commIDs []string) {
	comms := subscribedCommunities(channel, conn)
	for _, comm := range comms {
		commIDs = append(commIDs, comm.ID)
	}
	return
}

func FilterCommunitiesByPlatformID(commsIn []models.AdaptiveCommunity, platformID daosCommon.PlatformID) (comms []models.AdaptiveCommunity) {
	for _, comm := range commsIn {
		if comm.PlatformID == platformID {
			comms = append(comms, comm)
		}
	}
	return
}

func subscribedCommunities(channel string, conn daosCommon.DynamoDBConnection) (comms []models.AdaptiveCommunity) {
	var commsForAllPlatforms []models.AdaptiveCommunity
	commsForAllPlatforms, err2 := adaptiveCommunity.ReadByChannel(channel)(conn)
	err2 = wrapError(err2, "subscribedCommunities")
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not get subscribed communities for %s channel", channel))

	comms = FilterCommunitiesByPlatformID(commsForAllPlatforms, conn.PlatformID)
	return
}

func createCommunityFromCreatorUser(creatorUserID string, 
	channelID string, communityName string,
	conn daosCommon.DynamoDBConnection) (err error) {
	// Let's add this channel as a new user
	// the information about the user who initiated this
	var creators []models.User
	creators, err = daosUser.ReadOrEmpty(creatorUserID)(conn)
	var creator models.User
	if len(creators) > 0 {
		creator = creators[0]
	} else {
		log.Printf("Not found in users id=%s", creatorUserID)
	}
	if err == nil {
		item := models.User{
			ID:             channelID,
			DisplayName:    fmt.Sprintf("adaptive-%s", communityName),
			FirstName:      "",
			LastName:       "",
			Timezone:       creator.Timezone,
			TimezoneOffset: creator.TimezoneOffset,
			PlatformID:     creator.PlatformID,
			PlatformOrg:    creator.PlatformOrg,
			IsAdmin:        false,
			// Deleted:     false,
			DeactivatedAt: "",
			CreatedAt:     core.CurrentRFCTimestamp(),
			IsShared:      true,
		}
		err = daosUser.Create(item)(conn)
	}
	return
}

func addUserToAllCommunities(userID string, subscribedCommunityIDs []models.AdaptiveCommunity,
	conn daosCommon.DynamoDBConnection) (res []models.AdaptiveCommunityUser3) {
	for _, each := range subscribedCommunityIDs {
		// For each subscribed community, add an entry in community users table
		commUser := adaptiveCommunityUser.AdaptiveCommunityUser{
			ChannelID:   each.ChannelID,
			UserID:      userID,
			CommunityID: each.ID,
			PlatformID:  conn.PlatformID,
		}
		adaptiveCommunityUser.CreateUnsafe(commUser)(conn)
		res = append(res, commUser)
	}
	return
}

func addUsersToCommunity(teamID models.TeamID, channelID string, communityID string, userIDs []string) (res []models.AdaptiveCommunityUser3) {
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	// Adding existing channel members
	for _, each := range userIDs {
		commUser := adaptiveCommunityUser.AdaptiveCommunityUser{
			PlatformID:  teamID.ToPlatformID(),
			CommunityID: communityID,
			ChannelID:   channelID,
			UserID:      each,
		}
		adaptiveCommunityUser.CreateUnsafe(commUser)(conn)
		res = append(res, commUser)
	}
	return
}

// removeChannel remove all subscriptions to the channel
func removeChannel(userID, channelID string, conn daosCommon.DynamoDBConnection) {
	logger.Infof("Removing channel %s because user=%s left channel", channelID, userID)
	teamID := models.ParseTeamID(conn.PlatformID)
	// Adaptive bot is removed from the channel
	comms := subscribedCommunities(channelID, conn)
	logger.Infof("There where %d communities associated with the channel", len(comms))
	// We should delete this channel from users table and deactivate the community
	for _, each := range comms {
		// Delete users from community users table
		communityUser.DeactivateAllCommunityMembersUnsafe(teamID, each.ChannelID)(conn)
		// Delete entry from communities table
		adaptiveCommunity.DeactivateUnsafe(each.PlatformID, each.ID)(conn)
		// Deleting channel user
		daosUser.DeactivateUnsafe(channelID)(conn)
		// Unset channel for strategy communities
		unsetStrategyCommunities(channelID)
		// Post confirmation to Admin about the removal
		postSubscriptionRemovalToAdmin(teamID, each.ID, userID)
	}
}

// TODO: Update this to remove by community id instead of channel id
// This is assuming that there is only one community per channel
func deleteCommunityMembersByCommunityID(teamID models.TeamID, communityID string, channelID string, conn daosCommon.DynamoDBConnection) (err error) {
	return communityUser.DeactivateAllCommunityMembers(teamID, channelID)(conn)
}

// channelUnsubscribe removes the channel association with a community.
// Also removes all users from the community.
func channelUnsubscribe(channelID string, 
	conn daosCommon.DynamoDBConnection) (err error) {
	teamID := models.ParseTeamID(conn.PlatformID)
	var subComms []adaptiveCommunity.AdaptiveCommunity
	// Delete the entry from user table only if this is the only unsubscribed community
	subComms, err = adaptiveCommunity.ReadByChannel(channelID)(connGen.ForPlatformID(teamID.ToPlatformID()))
	logger.Infof("Subscribed communities for %s channel in %s platform: %v", channelID, teamID, subComms)
	if err == nil {
		// We should delete this channel from users table and deactivate the community
		// Only when channel has one unsubscribed community and that is indeed the chosen one to unsubscribed, then delete the user

		for _, eachComm := range subComms {
			// Delete entry from user table
			err = daosUser.Deactivate(channelID)(conn)
			if err == nil {
				// Delete users from user communities table for the community
				err = deleteCommunityMembersByCommunityID(teamID, eachComm.ID, eachComm.ChannelID, conn)
				if err == nil {
					logger.Infof("Removed all community members in %s community for %s platform", eachComm.ID, teamID)
					// Delete entry from communities table
					err = adaptiveCommunity.Deactivate(conn.PlatformID, eachComm.ID)(conn)
					err = errors.Wrapf(err, "Could not delete from adaptiveCommunity table in %s platform", teamID)
					if err == nil {
						logger.Infof("Removed %v community for %s platform", eachComm, teamID)
					}
				}
			} else {
				logger.
					WithField("namespace", namespace).
					WithError(err).
					Errorf("Could not channelUnsubscribe 1(channel=%s, platform=%v) eachComm=%v", channelID, teamID, eachComm)
				}
		}
	}
	if err != nil {
		logger.
			WithField("namespace", namespace).
			WithError(err).
			Errorf("Could not channelUnsubscribe(channel=%s, platform=%v)", channelID, teamID)
	}
	return
}

func channelUnsubscribeUnsafe(channelID string, conn daosCommon.DynamoDBConnection) {
	// teamID := models.ParseTeamID(conn.PlatformID)
	err2 := channelUnsubscribe(channelID, conn)
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not handle channel_deleted event for channel %s", channelID))
}

func updateStrategyCommunity(channelID string, strategyCommunityID string) error {
	// A channel has been created for a objective community. Update strategy communities with the same
	// Set channel_created and channel_id values
	exprAttributes := map[string]*dynamodb.AttributeValue{
		":cc": dynNumber(1),
		":ci": dynString(channelID),
	}
	key := idParams(strategyCommunityID)
	updateExpression := "set channel_created = :cc, channel_id = :ci"
	return d.UpdateTableEntry(exprAttributes, key, updateExpression, strategyCommunitiesTable)
}
