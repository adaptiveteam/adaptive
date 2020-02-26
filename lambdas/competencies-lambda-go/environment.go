package competencies

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	evalues "github.com/adaptiveteam/adaptive/adaptive-engagements/values"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	// plat "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"

	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	dialogFetcher "github.com/adaptiveteam/adaptive/dialog-fetcher"

	// "github.com/nlopes/slack"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/communityUser"

)

var (
	platformNotificationTopic = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	namespace                 = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                    = utils.NonEmptyEnv("AWS_REGION")
	engagementTable           = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	dialogTableName           = utils.NonEmptyEnv("DIALOG_TABLE")
	clientID                  = utils.NonEmptyEnv("CLIENT_ID")

	sns              = awsutils.NewSNS(region, "", namespace)
	d                = awsutils.NewDynamo(region, "", namespace)
	dns              = common.DynamoNamespace{Dynamo: d, Namespace: namespace}
	schema           = models.SchemaForClientID(clientID)
	valuesDao        = evalues.NewDAOFromSchema(&dns, schema)
	dialogFetcherDao = dialogFetcher.NewDAO(d, dialogTableName)

	adaptiveValuesTableDao = valuesDao

	platform = utils.Platform{
		Sns:                       *sns,
		PlatformNotificationTopic: platformNotificationTopic,
		Namespace:                 namespace,
		IsInteractiveDebugEnabled: false,
	}

	userDAO          = utilsUser.NewDAOFromSchema(d, namespace, schema)
	// community
	communityUsersTable              = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	communityUsersUserCommunityIndex = utils.NonEmptyEnv("COMMUNITY_USERS_USER_COMMUNITY_INDEX")
	communityUsersUserIndex          = utils.NonEmptyEnv("COMMUNITY_USERS_USER_INDEX")
	// TODO: use DAO for the query
	communityUserDAO = communityUser.NewDAOFromSchema(d, namespace, schema)
)

func slackAPI(teamID models.TeamID) mapper.PlatformAPI {
	conn:= daosCommon.DynamoDBConnection{
		Dynamo: d,
		ClientID: clientID,
		PlatformID: teamID.ToPlatformID(),
	}
	return mapper.SlackAdapterForTeamID(conn)
}
