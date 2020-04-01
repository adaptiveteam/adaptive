package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

var (
	namespace                    = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                       = utils.NonEmptyEnv("AWS_REGION")
	engagementTable              = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	userTable                    = utils.NonEmptyEnv("USERS_TABLE_NAME")
	usersPlatformIndex           = utils.NonEmptyEnv("USERS_PLATFORM_INDEX")
	communityUsersTable          = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	communityUsersCommunityIndex = string(adaptiveCommunityUser.PlatformIDCommunityIDIndex)
	// coachingRelationshipsTable        = utils.NonEmptyEnv("COACHING_RELATIONSHIPS_TABLE_NAME")
	coachingRelationshipsCoacheeIndex = "CoacheeQuarterYearIndex"
	coachQuarterYearIndex             = utils.NonEmptyEnv("COACHING_RELATIONSHIPS_COACH_QUARTER_YEAR_INDEX")
	coachingRelationshipsQYIndex      = utils.NonEmptyEnv("COACHING_RELATIONSHIPS_QUARTER_YEAR_INDEX")
	userObjectivesTable               = utils.NonEmptyEnv("USER_OBJECTIVES_TABLE")
	userObjectivesUserIdIndex         = utils.NonEmptyEnv("USER_OBJECTIVES_USER_ID_INDEX")
	userObjectivesPartnerIndex        = utils.NonEmptyEnv("USER_OBJECTIVES_PARTNER_INDEX")
	userProfileLambda                 = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	platformNotificationTopic         = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	communitiesTable                  = utils.NonEmptyEnv("ADAPTIVE_COMMUNITIES_TABLE")

	feedbackTableName              = utils.NonEmptyEnv("USER_FEEDBACK_TABLE_NAME")
	feedbackSourceQuarterYearIndex = utils.NonEmptyEnv("USER_FEEDBACK_SOURCE_QUARTER_YEAR_INDEX")
	feedbackReportLambda           = utils.NonEmptyEnv("FEEDBACK_REPORTING_LAMBDA_NAME")

	userObjectivesProgressTable   = utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_TABLE")
	userObjectivesProgressIdIndex = utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_ID_INDEX")

	adHocHolidaysTable             = utils.NonEmptyEnv("HOLIDAYS_AD_HOC_TABLE")
	adHocHolidaysPlatformDateIndex = utils.NonEmptyEnv("HOLIDAYS_PLATFORM_DATE_INDEX")

	strategyInitiativesTableName     = utils.NonEmptyEnv("STRATEGY_INITIATIVES_TABLE")
	strategyInitiativesPlatformIndex = "PlatformIDIndex"
	strategyObjectivesTableName      = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_TABLE")
	strategyObjectivesPlatformIndex  = "PlatformIDIndex"

	d       = awsutils.NewDynamo(region, "", namespace)
	l       = awsutils.NewLambda(region, "", namespace)
	s       = awsutils.NewSNS(region, "", namespace)
	dns     = common.DynamoNamespace{Dynamo: d, Namespace: namespace}
	bharath = "UE48A5TC0"
)
