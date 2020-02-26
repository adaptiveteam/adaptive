package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
)

var (
	engagementTable              = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	platformNotificationTopic    = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	userProfileLambda            = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	communityUsersTable          = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	communityUsersUserIndex      = utils.NonEmptyEnv("COMMUNITY_USERS_USER_INDEX")
	communityUsersCommunityIndex = utils.NonEmptyEnv("COMMUNITY_USERS_COMMUNITY_INDEX")
	namespace                    = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                       = utils.NonEmptyEnv("AWS_REGION")

	visionTable                                 = utils.NonEmptyEnv("VISION_TABLE_NAME")
	strategyObjectivesTable                     = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_TABLE_NAME") // required
	capabilityCommunitiesTable                  = utils.NonEmptyEnv("CAPABILITY_COMMUNITIES_TABLE_NAME")
	strategyInitiativesTable                    = utils.NonEmptyEnv("STRATEGY_INITIATIVES_TABLE_NAME") // required
	strategyInitiativeCommunitiesTable          = utils.NonEmptyEnv("STRATEGY_INITIATIVE_COMMUNITIES_TABLE_NAME")
	strategyInitiativesInitiativeCommunityIndex = "InitiativeCommunityIDIndex"
	strategyCommunitiesTable                    = utils.NonEmptyEnv("STRATEGY_COMMUNITIES_TABLE_NAME")

	strategyObjectivesPlatformIndex            = "PlatformIDIndex"
	strategyObjectivesCapabilityCommunityIndex = "CapabilityCommunityIDsIndex"
	capabilityCommunitiesPlatformIndex         = "PlatformIDIndex"
	strategyInitiativesPlatformIndex           = "PlatformIDIndex"
	strategyInitiativeCommunitiesPlatformIndex = "PlatformIDIndex"
	strategyCommunitiesPlatformIndex           = "PlatformIDIndex"
	communityUsersUserCommunityIndex           = "UserIDCommunityIDIndex"

	usersTableName     = utils.NonEmptyEnv("USERS_TABLE_NAME")
	usersPlatformIndex = utils.NonEmptyEnv("USERS_PLATFORM_INDEX")

	// IDOs - need to update user objective for the related strategy objectives
	userObjectivesTable         = utils.NonEmptyEnv("USER_OBJECTIVES_TABLE")
	userObjectivesUserIndex     = utils.NonEmptyEnv("USER_OBJECTIVES_USER_ID_INDEX")
	userObjectivesProgressTable = utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_TABLE")
	userObjectivesTypeIndex     = "UserIDTypeIndex"

	communitiesTable = utils.NonEmptyEnv("USER_COMMUNITIES_TABLE")

	d        = awsutils.NewDynamo(utils.NonEmptyEnv("AWS_REGION"), "", namespace)
	s        = awsutils.NewSNS(region, "", namespace)
	dns      = common.DynamoNamespace{Dynamo: d, Namespace: namespace}
	clientID = utils.NonEmptyEnv("CLIENT_ID")
	connGen  = daosCommon.DynamoDBConnectionGen{
		Dynamo:     d,
		TableNamePrefix:   clientID,
	}
)

func slackAPI(teamID models.TeamID) mapper.PlatformAPI {
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	return mapper.SlackAdapterForTeamID(conn)
}
