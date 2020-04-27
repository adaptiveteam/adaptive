package lambda

import (
	"regexp"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/communityUser"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsPlatform "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	// daosUser "github.com/adaptiveteam/adaptive/daos/user"
)

var (
	allComms = []string{string(community.HR), string(community.Admin), string(community.Coaching), string(community.User),
		string(community.Strategy), string(community.Competency)}
	region = utils.NonEmptyEnv("AWS_REGION")

	botMentionRegex = regexp.MustCompile(`(?m)<@([a-zA-Z0-9]+)> ([a-zA-Z\s\d]+)`)
	// `<@UEFF123> test report <@REFF123>`
	requestForUserRegex             = regexp.MustCompile(`(?m)<@([a-zA-Z0-9]+)> ([a-zA-Z\s\d]+) <@([a-zA-Z0-9]+)>`)
	userProfileLambda               = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	usersTable                      = utils.NonEmptyEnv("USERS_TABLE_NAME")
	engagementTable                 = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	orgCommunitiesTable             = utils.NonEmptyEnv("USER_COMMUNITIES_TABLE")
	namespace                       = utils.NonEmptyEnv("LOG_NAMESPACE")
	FeedbackReportPostingLambdaName = utils.NonEmptyEnv("FEEDBACK_REPORT_POSTING_LAMBDA_NAME")
	FeedbackReportingLambdaName     = utils.NonEmptyEnv("FEEDBACK_REPORTING_LAMBDA_NAME")
	platformNotificationTopic       = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	engScriptingLambda              = utils.NonEmptyEnv("USER_ENGAGEMENT_SCRIPTING_LAMBDA_NAME")
	coachingRelationshipsTable      = "" // utils.NonEmptyEnv("COACHING_RELATIONSHIPS_TABLE_NAME")
	// coachingRejectionsTable       = utils.NonEmptyEnv("COACHING_REJECTIONS_TABLE_NAME")
	userObjectivesTable             = utils.NonEmptyEnv("USER_OBJECTIVES_TABLE_NAME")
	userObjectivesProgressTable     = utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_TABLE")
	communityUsersTable             = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	engagementSchedulerLambda       = utils.NonEmptyEnv("ENGAGEMENT_SCHEDULER_LAMBDA_NAME")
	userSetupLambda                 = utils.NonEmptyEnv("USER_SETUP_LAMBDA_NAME")

	strategyCommunitiesTable            = utils.NonEmptyEnv("STRATEGY_COMMUNITIES_TABLE_NAME")
	capabilityCommunitiesTable          = utils.NonEmptyEnv("CAPABILITY_COMMUNITIES_TABLE_NAME")
	strategyInitiativeCommunitiesTable  = utils.NonEmptyEnv("INITIATIVE_COMMUNITIES_TABLE_NAME")
	strategyCommunitiesChannelIndex     = "StrategyCommunitiesChannelIndex"
	strategyVisionTableName             = utils.NonEmptyEnv("VISION_TABLE_NAME")

	adHocHolidaysTable             = utils.NonEmptyEnv("HOLIDAYS_AD_HOC_TABLE")

	d         = awsutils.NewDynamo(region, "", namespace)
	lambdaAPI = awsutils.NewLambda(region, "", namespace)
	s         = awsutils.NewSNS(region, "", namespace)

	dns = common.DynamoNamespace{Dynamo: d, Namespace: namespace}

	// Deprecated:
	userLambda = utils.UserProfileLambda{
		Region:            region,
		Namespace:         namespace,
		ProfileLambdaName: userProfileLambda,
	}

	clientID         = utils.NonEmptyEnv("CLIENT_ID")
	schema           = models.SchemaForClientID(clientID)
	// userDAO          = daosUser.NewDAOByTableName(d, namespace, schema.AdaptiveUsers.Name)
	communityUserDAO = communityUser.NewDAOFromSchema(d, namespace, schema)
	connGen          = daosCommon.DynamoDBConnectionGen{
		Dynamo:          d,
		TableNamePrefix: clientID,
	}
)

func userTokenSyncUnsafe(userID string) string {
	token, err2 := utilsPlatform.GetTokenForUser(d, clientID, userID)
	core.ErrorHandler(err2, "userTokenSyncUnsafe", "GetTokenForUser")
	return token
}
