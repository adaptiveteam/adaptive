package lambda

import (
	"github.com/pkg/errors"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/strategyCommunity"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func availableCommunities(teamID models.TeamID) []string {
	// Get all used communities
	comms, err := communityDAO.ReadAll(teamID)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not scan %s table", orgCommunitiesTable))
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

func getCommUsers(channelID string) (commUsers []models.AdaptiveCommunityUser3, err error) {
	commUsers, err = communityUserDAO.ReadCommunityUsers(channelID)
	return
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

func deleteCommunityTableEntry(ID string, teamID models.TeamID) {
	communityDAO.DeleteUnsafe(teamID, ID)
}

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

func subscribedCommunityIDs(teamID models.TeamID, channel string) (commIDs []string) {
	comms := subscribedCommunities(teamID, channel)
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

func subscribedCommunities(teamID models.TeamID, channel string) (comms []models.AdaptiveCommunity) {
	var commsForAllPlatforms []models.AdaptiveCommunity
	commsForAllPlatforms, err2 := communityDAO.ReadByChannelID(channel)
	err2 = wrapError(err2, "subscribedCommunities")
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not get subscribed communities for %s channel", channel))

	comms = FilterCommunitiesByPlatformID(commsForAllPlatforms, teamID.ToPlatformID())
	return
}

func createCommunityFromCreatorUser(creatorUserID string, channelID string, communityName string) (err error) {
	// Let's add this channel as a new user
	// the information about the user who initiated this
	var creators []models.User
	creators, err = userDAO.ReadOrEmpty(creatorUserID)
	creator := models.User{}
	if len(creators) > 0 {
		creator = creators[0]
	} else {
		log.Printf("Not found in users id=%s", creatorUserID)
		err = nil
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
		err = userDAO.Create(item)
	}
	return
}

func addUserToAllCommunities(teamID models.TeamID, userID string, subscribedCommunityIDs []models.AdaptiveCommunity) (res []models.AdaptiveCommunityUser3) {
	for _, each := range subscribedCommunityIDs {
		// For each subscribed community, add an entry in community users table
		commUser := models.AdaptiveCommunityUser3{
			ChannelID:   each.ChannelID,
			UserID:      userID,
			CommunityID: each.ID,
			PlatformID:  teamID.ToPlatformID(),
		}
		communityUserDAO.CreateUnsafe(commUser)
		res = append(res, commUser)
	}
	return
}

func addUsersToCommunity(teamID models.TeamID, channelID string, communityID string, userIDs []string) (res []models.AdaptiveCommunityUser3) {
	// Adding existing channel members
	for _, each := range userIDs {
		commUser := models.AdaptiveCommunityUser3{
			CommunityID: communityID,
			UserID:      each,
			ChannelID:   channelID,
			PlatformID:  teamID.ToPlatformID(),
		}
		communityUserDAO.CreateUnsafe(commUser)
		res = append(res, commUser)
	}
	return
}

// removeChannel remove all subscriptions to the channel
func removeChannel(userID, channelID string, teamID models.TeamID) {
	logger.Infof("Removing channel %s because user=%s left channel", channelID, userID)
	// Adaptive bot is removed from the channel
	comms := subscribedCommunities(teamID, channelID)
	logger.Infof("There where %d communities associated with the channel", len(comms))
	// We should delete this channel from users table and deactivate the community
	for _, each := range comms {
		// Delete users from community users table
		communityUserDAO.DeleteAllCommunityMembersUnsafe(each.ChannelID)
		// Delete entry from communities table
		communityDAO.DeleteUnsafe(models.ParseTeamID(each.PlatformID), each.ID)
		// Deleting channel user
		userDAO.DeactivateUnsafe(channelID)
		// Unset channel for strategy communities
		unsetStrategyCommunities(channelID)
		// Post confirmation to Admin about the removal
		postSubscriptionRemovalToAdmin(teamID, each.ID, userID)
	}
}

// TODO: Update this to remove by community id instead of channel id
// This is assuming that there is only one community per channel
func deleteCommunityMembersByCommunityID(communityID string, channelID string) (err error) {
	return communityUserDAO.DeleteAllCommunityMembers(channelID)
}

// channelUnsubscribe removes the channel association with a community.
// Also removes all users from the community.
func channelUnsubscribe(channelID string, teamID models.TeamID) (err error) {
	var subComms []adaptiveCommunity.AdaptiveCommunity
	// Delete the entry from user table only if this is the only unsubscribed community
	subComms, err = adaptiveCommunity.ReadByChannel(channelID)(connGen.ForPlatformID(teamID.ToPlatformID()))
	logger.Infof("Subscribed communities for %s channel in %s platform: %v", channelID, teamID, subComms)
	if err == nil {
		// We should delete this channel from users table and deactivate the community
		// Only when channel has one unsubscribed community and that is indeed the chosen one to unsubscribed, then delete the user

		for _, eachComm := range subComms {
			// Delete entry from user table
			err = userDAO.Deactivate(channelID)
			if err == nil {
				// Delete users from user communities table for the community
				err = deleteCommunityMembersByCommunityID(eachComm.ID, eachComm.ChannelID)
				if err == nil {
					logger.Infof("Removed all community members in %s community for %s platform", eachComm.ID, teamID)
					commParams := idAndPlatformIDParams(eachComm.ID, teamID)
					// Delete entry from communities table
					err = d.DeleteEntry(orgCommunitiesTable, commParams)
					err = errors.Wrapf(err, "Could not delete from %s table in %s platform", orgCommunitiesTable, teamID)
					if err == nil {
						logger.Infof("Removed %v community for %s platform", eachComm, teamID)
					}
				}
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

func channelUnsubscribeUnsafe(channelID string, teamID models.TeamID) {
	err := channelUnsubscribe(channelID, teamID)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not handle channel_deleted event for channel %s", channelID))
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
