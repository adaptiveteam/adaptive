package adaptive_checks

import (
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

type environment struct {
	visionTable string

	// IDOs
	userObjectivesTable string
	userObjectivesUserIndex string
	userObjectivesProgressTable string
	userObjectivesTypeIndex string

	// community
	communityUsersTable string
	communityUsersUserCommunityIndex string
	communityUsersUserIndex string

	// engagements
	engagementsTable string
	engagementsAnsweredIndex string

	// strategy
	initiativesTable string
	initiativesPlatformIndex string
	strategyObjectivesTableName string
	strategyObjectivesPlatformIndex string
	capabilityCommunitiesTable string
	capabilityCommunitiesPlatformIndex string
	initiativeCommunitiesTableName string
	initiativeCommunitiesPlatformIndex string
	strategyCommunitiesTable string
	strategyObjectivesCapabilityCommunityIndex string
	strategyInitiativesInitiativeCommunityIndex string
	userFeedbackTable string
	userFeedbackSourceQYIndex string
	reportsBucket string
	namespace string
	region string
	clientID string
	d         *awsutils.DynamoRequest
	schema    models.Schema
	connGen   daosCommon.DynamoDBConnectionGen

	dns       common.DynamoNamespace
}

func readEnvironment() environment {
	namespace:=  utils.NonEmptyEnv("LOG_NAMESPACE")
	region   :=  utils.NonEmptyEnv("AWS_REGION")
	d        :=  awsutils.NewDynamo(region, "", namespace)
	clientID :=  utils.NonEmptyEnv("CLIENT_ID")
	schema   :=  models.SchemaForClientID(clientID)

	return environment{
		// vision
		visionTable:  utils.NonEmptyEnv("VISION_TABLE_NAME"),

		// IDOs
		userObjectivesTable        :  utils.NonEmptyEnv("USER_OBJECTIVES_TABLE_NAME"),
		userObjectivesUserIndex    :  utils.NonEmptyEnv("USER_OBJECTIVES_USER_ID_INDEX"),
		userObjectivesProgressTable:  utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_TABLE"),
		userObjectivesTypeIndex    :  "UserIDTypeIndex",

		// community
		communityUsersTable             :  utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME"),
		communityUsersUserCommunityIndex:  utils.NonEmptyEnv("COMMUNITY_USERS_USER_COMMUNITY_INDEX"),
		communityUsersUserIndex         :  utils.NonEmptyEnv("COMMUNITY_USERS_USER_INDEX"),

		// engagements
		engagementsTable        :  utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME"),
		engagementsAnsweredIndex:  "UserIDAnsweredIndex",

		// strategy
		initiativesTable                  :  utils.NonEmptyEnv("STRATEGY_INITIATIVES_TABLE_NAME"),
		initiativesPlatformIndex          :  "PlatformIDIndex",
		strategyObjectivesTableName       :  utils.NonEmptyEnv("STRATEGY_OBJECTIVES_TABLE_NAME"),
		strategyObjectivesPlatformIndex   :  "PlatformIDIndex",
		capabilityCommunitiesTable        :  utils.NonEmptyEnv("CAPABILITY_COMMUNITIES_TABLE_NAME"),
		capabilityCommunitiesPlatformIndex:  "PlatformIDIndex",
		initiativeCommunitiesTableName    :  utils.NonEmptyEnv("INITIATIVE_COMMUNITIES_TABLE_NAME"),
		initiativeCommunitiesPlatformIndex:  utils.NonEmptyEnv("INITIATIVE_COMMUNITIES_PLATFORM_INDEX"),
		strategyCommunitiesTable          :  utils.NonEmptyEnv("STRATEGY_COMMUNITIES_TABLE_NAME"),

		strategyObjectivesCapabilityCommunityIndex :  "CapabilityCommunityIDsIndex",
		strategyInitiativesInitiativeCommunityIndex:  "InitiativeCommunityIDIndex",

		userFeedbackTable        :  utils.NonEmptyEnv("USER_FEEDBACK_TABLE_NAME"),
		userFeedbackSourceQYIndex:  utils.NonEmptyEnv("USER_FEEDBACK_SOURCE_QUARTER_YEAR_INDEX"),
		reportsBucket            :  utils.NonEmptyEnv("REPORTS_BUCKET_NAME"),

		namespace:  namespace,
		region   :  region,
		d        :  d,
		clientID :  clientID,
		schema   :  schema,
		connGen  :  daosCommon.CreateConnectionGenFromEnv(),

		dns                  :  common.DynamoNamespace{Dynamo: d, Namespace: namespace},
		// adHocHolidaysTableDao = eholidays.NewDAO(&dns, schema.Holidays.Name, schema.Holidays.PlatformDateIndex)
	}
}

var	DateFormat = core.ISODateLayout

// UserIDToPlatformID converts userID to teamID using
// globally available variables.
func UserIDToPlatformID(userDAO utilsUser.DAO) func(string) daosCommon.PlatformID {
	return func(userID string) (daosCommon.PlatformID) {
		return userDAO.ReadUnsafe(userID).PlatformID
	}
}

// UserIDToTeamID converts userID to teamID using
// globally available variables.
func UserIDToTeamID(userDAO utilsUser.DAO) func(string) models.TeamID {
	return func(userID string) (teamID models.TeamID) {
		return models.ParseTeamID(userDAO.ReadUnsafe(userID).PlatformID)
	}
}
