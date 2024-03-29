package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/sql-connector"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	dialogFetcher "github.com/adaptiveteam/adaptive/dialog-fetcher"
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
	communityUsersUserCommunityIndex = "UserIDCommunityIDIndex"
	communityUsersUserIndex          = utils.NonEmptyEnv("COMMUNITY_USERS_USER_INDEX")
	dialogTableName                  = utils.NonEmptyEnv("DIALOG_TABLE")
	visionTable                      = utils.NonEmptyEnv("VISION_TABLE_NAME")

	strategyObjectivesTableName                    = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_TABLE_NAME")
	strategyObjectivesPlatformIndex                = "PlatformIDIndex"
	capabilityCommunitiesTable                     = utils.NonEmptyEnv("CAPABILITY_COMMUNITIES_TABLE_NAME")     // required
	capabilityCommunitiesPlatformIndex             = "PlatformIDIndex"                                          // required
	initiativeCommunitiesTable                     = utils.NonEmptyEnv("INITIATIVE_COMMUNITIES_TABLE_NAME")     // required
	initiativeCommunitiesPlatformIndex             = utils.NonEmptyEnv("INITIATIVE_COMMUNITIES_PLATFORM_INDEX") // required
	_                                              = utils.NonEmptyEnv("STRATEGY_COMMUNITIES_TABLE_NAME")
	strategyCommunitiesPlatformChannelCreatedIndex = "PlatformIDChannelCreatedIndex"

	engScriptingLambda = utils.NonEmptyEnv("USER_ENGAGEMENT_SCRIPTING_LAMBDA_NAME")

	dialogFetcherDAO = dialogFetcher.NewDAO(d, dialogTableName)

	clientID = utils.NonEmptyEnv("CLIENT_ID")
	connGen  = daosCommon.DynamoDBConnectionGen{
		Dynamo:     d,
		TableNamePrefix:   clientID,
	}
	schema   = models.SchemaForClientID(clientID)

	usersTable           = utils.NonEmptyEnv("USERS_TABLE_NAME")
	profileLambdaName    = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	slackProcessorSuffix = utils.NonEmptyEnv("SLACK_MESSAGE_PROCESSOR_SUFFIX")

	logger = alog.LambdaLogger(logrus.InfoLevel)

	// platformAdapter  = mapper.SlackAdapter2(platformTokenDAO)
)

type RDSConfig = sqlconnector.RDSConfig

var ReadRDSConfigFromEnv = sqlconnector.ReadRDSConfigFromEnv

func globalConnection(teamID models.TeamID) daosCommon.DynamoDBConnection {
	return connGen.ForPlatformID(teamID.ToPlatformID())
}
