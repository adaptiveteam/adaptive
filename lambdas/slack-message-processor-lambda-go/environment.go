package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-reports/utilities"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	dialogFetcher "github.com/adaptiveteam/adaptive/dialog-fetcher"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/sirupsen/logrus"
)

var (
	SayHelloMenuItem = "say hello"

	namespace                 = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                    = utils.NonEmptyEnv("AWS_REGION")
	sns                       = awsutils.NewSNS(region, "", namespace)
	payloadTopicArn           = utils.NonEmptyEnv("NAMESPACE_PAYLOAD_TOPIC_ARN")
	platformNotificationTopic = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	helpPage                  = utils.NonEmptyEnv("ADAPTIVE_HELP_PAGE")
	settingsCommands          = []string{user.UpdateSettings, user.AskForEngagements}
	feedbackCommands          = []string{coaching.GiveFeedback, user.FetchReport, user.GenerateReport, coaching.ViewCoachees}

	d                                = awsutils.NewDynamo(utils.NonEmptyEnv("AWS_REGION"), "", namespace)
	l                                = awsutils.NewLambda(utils.NonEmptyEnv("AWS_REGION"), "", namespace)
	_                                = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	userObjectivesTableName          = utils.NonEmptyEnv("USER_OBJECTIVES_TABLE_NAME")
	userObjectivesPartnerIndex       = utils.NonEmptyEnv("USER_OBJECTIVES_PARTNER_INDEX")
	userObjectivesUserIndex          = utils.NonEmptyEnv("USER_OBJECTIVES_USER_ID_INDEX")
	userObjectivesTypeIndex          = "UserIDTypeIndex"
	userCommunitiesTable             = utils.NonEmptyEnv("USER_COMMUNITIES_TABLE")
	communityUsersTable              = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	communityUsersUserCommunityIndex = utils.NonEmptyEnv("COMMUNITY_USERS_USER_COMMUNITY_INDEX")
	communityUsersUserIndex          = utils.NonEmptyEnv("COMMUNITY_USERS_USER_INDEX")
	dialogTableName                  = utils.NonEmptyEnv("DIALOG_TABLE")
	visionTable                      = utils.NonEmptyEnv("VISION_TABLE_NAME")

	strategyObjectivesTableName                    = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_TABLE_NAME")
	strategyObjectivesPlatformIndex                = "StrategyObjectivesPlatformIndex"
	capabilityCommunitiesTable                     = utils.NonEmptyEnv("CAPABILITY_COMMUNITIES_TABLE_NAME")     // required
	capabilityCommunitiesPlatformIndex             = "CapabilityCommunitiesPlatformIndex" // required
	initiativeCommunitiesTable                     = utils.NonEmptyEnv("INITIATIVE_COMMUNITIES_TABLE_NAME")     // required
	initiativeCommunitiesPlatformIndex             = utils.NonEmptyEnv("INITIATIVE_COMMUNITIES_PLATFORM_INDEX") // required
	_                                              = utils.NonEmptyEnv("STRATEGY_COMMUNITIES_TABLE_NAME")
	strategyCommunitiesPlatformChannelCreatedIndex = "StrategyCommunityPlatformChannelCreatedIndex"

	engScriptingLambda = utils.NonEmptyEnv("USER_ENGAGEMENT_SCRIPTING_LAMBDA_NAME")

	dialogFetcherDAO = dialogFetcher.NewDAO(d, dialogTableName)

	clientID = utils.NonEmptyEnv("CLIENT_ID")
	schema   = models.SchemaForClientID(clientID)

	userDAO = utilsUser.NewDAOFromSchema(d, namespace, schema)

	usersTable           = utils.NonEmptyEnv("USERS_TABLE_NAME")
	profileLambdaName    = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	slackProcessorSuffix = utils.NonEmptyEnv("SLACK_MESSAGE_PROCESSOR_SUFFIX")

	logger = alog.LambdaLogger(logrus.InfoLevel)

	platformTokenDAO = platform.NewDAOFromSchema(d, namespace, schema)
	platformAdapter  = mapper.SlackAdapter2(platformTokenDAO)
)

type RDSConfig struct {
	Driver           string
	ConnectionString string
}

// ReadRDSConfigFromEnv read config from env
func ReadRDSConfigFromEnv() RDSConfig {
	rdsHost := utils.NonEmptyEnv("RDS_HOST")
	GlobalRDSConfig := RDSConfig{Driver: "mysql", ConnectionString: utilities.ConnectionString(
		rdsHost,
		utils.NonEmptyEnv("RDS_USER"),
		utils.NonEmptyEnv("RDS_PASSWORD"),
		utils.NonEmptyEnv("RDS_PORT"),
		utils.NonEmptyEnv("RDS_DB_NAME"),
	)}
	return GlobalRDSConfig

}
