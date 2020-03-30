package common

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

var (
	Namespace                         = utils.NonEmptyEnv("LOG_NAMESPACE")
	Region                            = utils.NonEmptyEnv("AWS_REGION")
	EngagementTable                   = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	UserTable                         = utils.NonEmptyEnv("USERS_TABLE_NAME")
	UsersPlatformIndex                = utils.NonEmptyEnv("USERS_PLATFORM_INDEX")
	CommunityUsersTable               = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	CommunityUsersCommunityIndex      = string(adaptiveCommunityUser.PlatformIDCommunityIDIndex)
	CoachingRelationshipsTable        = utils.NonEmptyEnv("COACHING_RELATIONSHIPS_TABLE_NAME")
	CoachingRelationshipsCoacheeIndex = "CoacheeQuarterYearIndex"
	CoachQuarterYearIndex             = utils.NonEmptyEnv("COACHING_RELATIONSHIPS_COACH_QUARTER_YEAR_INDEX")
	CoachingRelationshipsQYIndex      = utils.NonEmptyEnv("COACHING_RELATIONSHIPS_QUARTER_YEAR_INDEX")
	UserObjectivesTable               = utils.NonEmptyEnv("USER_OBJECTIVES_TABLE")
	UserObjectivesUserIDIndex         = utils.NonEmptyEnv("USER_OBJECTIVES_USER_ID_INDEX")
	UserObjectivesPartnerIndex        = utils.NonEmptyEnv("USER_OBJECTIVES_PARTNER_INDEX")
	UserProfileLambda                 = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	PlatformNotificationTopic         = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	CommunitiesTable                  = utils.NonEmptyEnv("ADAPTIVE_COMMUNITIES_TABLE")

	FeedbackTableName              = utils.NonEmptyEnv("USER_FEEDBACK_TABLE_NAME")
	FeedbackSourceQuarterYearIndex = utils.NonEmptyEnv("USER_FEEDBACK_SOURCE_QUARTER_YEAR_INDEX")
	FeedbackReportLambda           = utils.NonEmptyEnv("FEEDBACK_REPORTING_LAMBDA_NAME")

	UserObjectivesProgressTable   = utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_TABLE")
	UserObjectivesProgressIDIndex = utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_ID_INDEX")

	AdHocHolidaysTable             = utils.NonEmptyEnv("HOLIDAYS_AD_HOC_TABLE")
	AdHocHolidaysPlatformDateIndex = utils.NonEmptyEnv("HOLIDAYS_PLATFORM_DATE_INDEX")

	StrategyInitiativesTableName     = utils.NonEmptyEnv("STRATEGY_INITIATIVES_TABLE")
	StrategyInitiativesPlatformIndex = "PlatformIDIndex"
	StrategyObjectivesTableName      = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_TABLE")
	StrategyObjectivesPlatformIndex  = "PlatformIDIndex"

	D   = awsutils.NewDynamo(Region, "", Namespace)
	L   = awsutils.NewLambda(Region, "", Namespace)
	S   = awsutils.NewSNS(Region, "", Namespace)
	Dns = common.DynamoNamespace{Dynamo: D, Namespace: Namespace}

	ClientID = utils.NonEmptyEnv("CLIENT_ID")
	Schema   = models.SchemaForClientID(ClientID)
	UserDAO  = utilsUser.NewDAOFromSchema(D, Namespace, Schema)
)

var (
	NowActionLabel  = models.YesLabel
	SkipActionLabel = models.DefaultSkipThisTemplate

	ViewOpenObjectives      = "view_open_objectives"
	ViewCollaborationReport = "view_collaboration_report"
)
