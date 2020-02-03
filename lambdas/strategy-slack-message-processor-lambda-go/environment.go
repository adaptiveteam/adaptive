package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
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
	strategyInitiativesInitiativeCommunityIndex = utils.NonEmptyEnv("STRATEGY_INITIATIVES_INITIATIVE_COMMUNITY_ID_INDEX")
	strategyCommunitiesTable                    = utils.NonEmptyEnv("STRATEGY_COMMUNITIES_TABLE_NAME")

	strategyObjectivesPlatformIndex            = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_PLATFORM_INDEX") // required
	strategyObjectivesCapabilityCommunityIndex = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_CAPABILITY_COMMUNITY_INDEX")
	capabilityCommunitiesPlatformIndex         = "CapabilityCommunitiesPlatformIndex"
	strategyInitiativesPlatformIndex           = utils.NonEmptyEnv("STRATEGY_INITIATIVES_PLATFORM_INDEX")
	strategyInitiativeCommunitiesPlatformIndex = utils.NonEmptyEnv("STRATEGY_INITIATIVE_COMMUNITIES_PLATFORM_INDEX")
	strategyCommunitiesPlatformIndex           = utils.NonEmptyEnv("STRATEGY_COMMUNITIES_PLATFORM_INDEX")
	communityUsersUserCommunityIndex           = utils.NonEmptyEnv("COMMUNITY_USERS_USER_COMMUNITY_INDEX")

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
)
