package community

import (
	"fmt"
	"log"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommunity "github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
)

func CommunityMembers(table string, ID string, teamID models.TeamID) []models.AdaptiveCommunityUser2 {
	var commUsers []models.AdaptiveCommunityUser2
	err2 := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(table, awsutils.DynamoIndexExpression{
		IndexName: string(adaptiveCommunityUser.PlatformIDCommunityIDIndex),
		Condition: "platform_id = :pi AND community_id = :c",
		Attributes: map[string]interface{}{
			":c":  ID,
			":pi": teamID.ToString(),
		},
	}, map[string]string{}, true, -1, &commUsers)
	err2 = errors.Wrapf(err2, "CommunityMembers(ID=%s, teamID=%s)", ID, teamID.ToString())
	core.ErrorHandler(err2, common.DeprecatedGetGlobalDns().Namespace,
		fmt.Sprintf("CommunityMembers: Could not find community  %s in %s table using %s index",
			ID, table, string(adaptiveCommunityUser.PlatformIDCommunityIDIndex)))
	return commUsers
}

func CommunityById(communityId string, teamID models.TeamID, communitiesTable string) models.AdaptiveCommunity {
	params := map[string]*dynamodb.AttributeValue{
		"id":          daosCommon.DynS(communityId),
		"platform_id": daosCommon.DynS(teamID.ToString()),
	}
	var comm models.AdaptiveCommunity
	err2 := common.DeprecatedGetGlobalDns().Dynamo.GetItemFromTable(communitiesTable, params, &comm)
	core.ErrorHandler(err2, "CommunityById", fmt.Sprintf("Could not find %s in %s table", communityId, communitiesTable))
	return comm
}

func CommunityDAO(conn daosCommon.DynamoDBConnection, communitiesTable string) daosCommunity.DAO {
	return daosCommunity.NewDAOByTableName(conn.Dynamo, "CommunityDAO", communitiesTable)
}

func SubscribedCommunities(channel string, communitiesTable, channelIndex string) (
	[]models.AdaptiveCommunity, error) {
	var comms []models.AdaptiveCommunity
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(communitiesTable, awsutils.DynamoIndexExpression{
		IndexName: channelIndex,
		Condition: "channel = :c",
		Attributes: map[string]interface{}{
			":c": channel,
		},
	}, map[string]string{}, true, -1, &comms)
	return comms, err
}

func queryUserCommIndex(userID, communityID string, communityUsersTable,
	communityUsersUserCommunityIndex string) []interface{} {
	var rels []interface{}
	log.Printf("queryUserCommIndex(userID=%s, communityID=%s, communityUsersTable=%s, communityUsersUserCommunityIndex=%s)",
		userID, communityID, communityUsersTable, communityUsersUserCommunityIndex)

	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(communityUsersTable, awsutils.DynamoIndexExpression{
		IndexName: communityUsersUserCommunityIndex,
		Condition: "user_id = :u and community_id = :c",
		Attributes: map[string]interface{}{
			":u": userID,
			":c": communityID,
		},
	}, map[string]string{}, true, -1, &rels)
	log.Printf("queryUserCommIndex(userID=%s, communityID=%s, communityUsersTable=%s, communityUsersUserCommunityIndex=%s) result=%+v",
		userID, communityID, communityUsersTable, communityUsersUserCommunityIndex, err)
	if err != nil {
		log.Printf("queryUserCommIndex err=%+v", err)
	}
	core.ErrorHandler(err, common.DeprecatedGetGlobalDns().Namespace, fmt.Sprintf("Could not query %s table on %s index",
		communityUsersTable, communityUsersUserCommunityIndex))
	return rels
}

// IsUserInCommunity checks if a user is part of an Adaptive Community
func IsUserInCommunity(userID string, communityUsersTable, communityUsersUserCommunityIndex string,
	community AdaptiveCommunity) (res bool) {
	defer func() {
		if err2 := recover(); err2 != nil {
			log.Printf("IsUserInCommunity got an error %+v", err2)
			res = false
		}
	}()
	rels := queryUserCommIndex(userID, string(community), communityUsersTable, communityUsersUserCommunityIndex)
	res = len(rels) > 0
	return
}

// PlatformCommunities lists all the communities based on the platform id supplied
func PlatformCommunities(teamID models.TeamID, communitiesTable, communityPlatformIndex string) (comms []models.AdaptiveCommunity, err error) {
	err = common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(communitiesTable, awsutils.DynamoIndexExpression{
		IndexName: communityPlatformIndex,
		Condition: "platform_id = :pi",
		Attributes: map[string]interface{}{
			":pi": teamID.ToString(),
		},
	}, map[string]string{}, true, -1, &comms)
	return
}

// PlatformCommunityMembers returns distinct user ids from all the communities based on a platform id
func PlatformCommunityMemberIDs(teamID models.TeamID, communitiesTable, communityPlatformIndex,
	communityUsersTable, communityUsersCommunityIndex string) (memberIDs []string, err error) {
	communities, err := PlatformCommunities(teamID, communitiesTable, communityPlatformIndex)
	var allCommunitiesMemberIDs []string
	if err == nil {
		for _, each := range communities {
			// Get community members by querying community users table based on platform id and community id
			members := CommunityMembers(communityUsersTable, each.ID, teamID)
			for _, member := range members {
				allCommunitiesMemberIDs = append(allCommunitiesMemberIDs, member.UserId)
			}
		}
	}
	memberIDs = core.Distinct(allCommunitiesMemberIDs)
	return
}
