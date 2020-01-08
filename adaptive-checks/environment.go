package adaptive_checks

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	eholidays "github.com/adaptiveteam/adaptive/adaptive-engagements/holidays"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
)

var (
	// vision
	visionTable = utils.NonEmptyEnv("VISION_TABLE_NAME")

	// IDOs
	userObjectivesTable         = utils.NonEmptyEnv("USER_OBJECTIVES_TABLE_NAME")
	userObjectivesUserIndex     = utils.NonEmptyEnv("USER_OBJECTIVES_USER_ID_INDEX")
	userObjectivesProgressTable = utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_TABLE")
	userObjectivesTypeIndex     = utils.NonEmptyEnv("USER_OBJECTIVES_TYPE_INDEX")

	// community
	communityUsersTable              = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	communityUsersUserCommunityIndex = utils.NonEmptyEnv("COMMUNITY_USERS_USER_COMMUNITY_INDEX")
	communityUsersUserIndex          = utils.NonEmptyEnv("COMMUNITY_USERS_USER_INDEX")

	// engagements
	engagementsTable         = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	engagementsAnsweredIndex = utils.NonEmptyEnv("USER_ANSWERED_INDEX")

	// strategy
	initiativesTable                   = utils.NonEmptyEnv("STRATEGY_INITIATIVES_TABLE_NAME")
	initiativesPlatformIndex           = utils.NonEmptyEnv("STRATEGY_INITIATIVES_PLATFORM_INDEX")
	strategyObjectivesTableName        = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_TABLE_NAME")
	strategyObjectivesPlatformIndex    = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_PLATFORM_INDEX")
	capabilityCommunitiesTable         = utils.NonEmptyEnv("CAPABILITY_COMMUNITIES_TABLE_NAME")
	capabilityCommunitiesPlatformIndex = utils.NonEmptyEnv("CAPABILITY_COMMUNITIES_PLATFORM_INDEX")
	initiativeCommunitiesTableName     = utils.NonEmptyEnv("INITIATIVE_COMMUNITIES_TABLE_NAME")
	initiativeCommunitiesPlatformIndex = utils.NonEmptyEnv("INITIATIVE_COMMUNITIES_PLATFORM_INDEX")
	strategyCommunitiesTable           = utils.NonEmptyEnv("STRATEGY_COMMUNITIES_TABLE_NAME")

	strategyObjectivesCapabilityCommunityIndex  = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_CAPABILITY_COMMUNITY_INDEX")
	strategyInitiativesInitiativeCommunityIndex = utils.NonEmptyEnv("STRATEGY_INITIATIVES_INITIATIVE_COMMUNITY_ID_INDEX")

	userFeedbackTable         = utils.NonEmptyEnv("USER_FEEDBACK_TABLE_NAME")
	userFeedbackSourceQYIndex = utils.NonEmptyEnv("USER_FEEDBACK_SOURCE_QUARTER_YEAR_INDEX")
	reportsBucket             = utils.NonEmptyEnv("REPORTS_BUCKET_NAME")

	DateFormat = core.ISODateLayout
)

var (
	namespace = utils.NonEmptyEnv("LOG_NAMESPACE")
	region    = utils.NonEmptyEnv("AWS_REGION")
	d         = awsutils.NewDynamo(region, "", namespace)
	clientID  = utils.NonEmptyEnv("CLIENT_ID")
	schema    = models.SchemaForClientID(clientID)
	userDAO   = utilsUser.NewDAOFromSchema(d, namespace, schema)
	// deprecated. We should change this to just client ID as soon as we rename table
	// coachingRelationshipsTable        = utils.NonEmptyEnv("COACHING_RELATIONSHIPS_TABLE_NAME")
	// coachingRelationshipDAO = coachingRelationship.NewDAO(d, namespace, clientID)
	userObjectiveDAO = userObjective.NewDAOByTableName(d, namespace, userObjectivesTable)

	dns                   = common.DynamoNamespace{Dynamo: d, Namespace: namespace}
	adHocHolidaysTableDao = eholidays.NewDAO(&dns, schema.Holidays.Name, schema.Holidays.PlatformDateIndex)
)

// UserIDToPlatformID converts userID to platformID using
// globally available variables.
func UserIDToPlatformID(userDAO utilsUser.DAO) func(string) models.PlatformID {
	return func(userID string) (platformID models.PlatformID) {
		return models.PlatformID(userDAO.ReadUnsafe(userID).PlatformId)
	}
}
