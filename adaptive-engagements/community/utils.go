package community

import (
	"log"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CommunityMembers(table string, ID string, platformID models.PlatformID, userCommIndex string) []models.AdaptiveCommunityUser2 {
	var commUsers []models.AdaptiveCommunityUser2
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(table, awsutils.DynamoIndexExpression{
		IndexName: userCommIndex,
		Condition: "platform_id = :pi AND community_id = :c",
		Attributes: map[string]interface{}{
			":c":  ID,
			":pi": platformID,
		},
	}, map[string]string{}, true, -1, &commUsers)
	core.ErrorHandler(err, common.DeprecatedGetGlobalDns().Namespace, fmt.Sprintf("Could not query %s table on %s index",
		table, userCommIndex))
	return commUsers
}

func CommunityById(communityId string, platformId models.PlatformID, communitiesTable string) models.AdaptiveCommunity {
	params := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(communityId),
		},
		"platform_id": {
			S: aws.String(string(platformId)),
		},
	}
	var comm models.AdaptiveCommunity
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTable(communitiesTable, params, &comm)
	core.ErrorHandler(err, common.DeprecatedGetGlobalDns().Namespace, fmt.Sprintf("Could not query %s table", communitiesTable))
	return comm
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
		userID,    communityID,    communityUsersTable,    communityUsersUserCommunityIndex)

	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(communityUsersTable, awsutils.DynamoIndexExpression{
		IndexName: communityUsersUserCommunityIndex,
		Condition: "user_id = :u and community_id = :c",
		Attributes: map[string]interface{}{
			":u": userID,
			":c": communityID,
		},
	}, map[string]string{}, true, -1, &rels)
	log.Printf("queryUserCommIndex(userID=%s, communityID=%s, communityUsersTable=%s, communityUsersUserCommunityIndex=%s) result=%+v", 
		userID,    communityID,    communityUsersTable,    communityUsersUserCommunityIndex, err)
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
	defer func(){
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
func PlatformCommunities(platformID models.PlatformID, communitiesTable, communityPlatformIndex string) (comms []models.AdaptiveCommunity, err error) {
	err = common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(communitiesTable, awsutils.DynamoIndexExpression{
		IndexName: communityPlatformIndex,
		Condition: "platform_id = :pi",
		Attributes: map[string]interface{}{
			":pi": platformID,
		},
	}, map[string]string{}, true, -1, &comms)
	return
}

// PlatformCommunityMembers returns distinct user ids from all the communities based on a platform id
func PlatformCommunityMemberIDs(platformID models.PlatformID, communitiesTable, communityPlatformIndex,
	communityUsersTable, communityUsersCommunityIndex string) (memberIDs []string, err error) {
	communities, err := PlatformCommunities(platformID, communitiesTable, communityPlatformIndex)
	var allCommunitiesMemberIDs []string
	if err == nil {
		for _, each := range communities {
			// Get community members by querying community users table based on platform id and community id
			members := CommunityMembers(communityUsersTable, each.ID, platformID, communityUsersCommunityIndex)
			for _, member := range members {
				allCommunitiesMemberIDs = append(allCommunitiesMemberIDs, member.UserId)
			}
		}
	}
	memberIDs = core.Distinct(allCommunitiesMemberIDs)
	return
}
