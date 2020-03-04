package lambda

import (
	"fmt"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsPlatform "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	dialogFetcher "github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/sirupsen/logrus"
	"github.com/nlopes/slack"
)

var (
	platformNotificationTopic = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	namespace                 = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                    = utils.NonEmptyEnv("AWS_REGION")

	engagementTable                  = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	userTable                        = utils.NonEmptyEnv("USERS_TABLE_NAME")
	usersPlatformIndex               = utils.NonEmptyEnv("USERS_PLATFORM_INDEX")
	userProfileLambda                = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	userObjectivesTable              = utils.NonEmptyEnv("USER_OBJECTIVES_TABLE_NAME")        // needed
	userObjectivesProgressTable      = utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_TABLE")    // needed
	userObjectivesProgressIdIndex    = utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_ID_INDEX") // needed
	userObjectivesUserIdIndex        = utils.NonEmptyEnv("USER_OBJECTIVES_USER_ID_INDEX")     // needed
	userObjectivesAcceptedIndex      = utils.NonEmptyEnv("USER_OBJECTIVES_ACCEPTED_INDEX")
	userObjectivesPartnerIndex       = utils.NonEmptyEnv("USER_OBJECTIVES_PARTNER_INDEX")
	userObjectivesTypeIndex          = utils.NonEmptyEnv("USER_OBJECTIVES_TYPE_INDEX")
	objectiveCloseoutPath            = utils.NonEmptyEnv("USER_OBJECTIVES_CLOSEOUT_LEARN_MORE_PATH")
	dialogTableName                  = utils.NonEmptyEnv("DIALOG_TABLE")
	communityUsersTable              = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	communityUsersCommunityIndex     = utils.NonEmptyEnv("COMMUNITY_USERS_COMMUNITY_INDEX")
	communityUsersUserCommunityIndex = utils.NonEmptyEnv("COMMUNITY_USERS_USER_COMMUNITY_INDEX")
	communityUsersUserIndex          = utils.NonEmptyEnv("COMMUNITY_USERS_USER_INDEX")
	communitiesTable                 = utils.NonEmptyEnv("USER_COMMUNITIES_TABLE")
	partnershipRejectionsTable       = utils.NonEmptyEnv("ACCOUNTABILITY_PARTNERSHIP_REJECTION_TABLE")

	capabilityCommunitiesTableName     = utils.NonEmptyEnv("CAPABILITY_COMMUNITIES_TABLE_NAME")
	strategyInitiativeCommunitiesTable = utils.NonEmptyEnv("INITIATIVE_COMMUNITIES_TABLE_NAME")

	strategyInitiativesTableName                = utils.NonEmptyEnv("STRATEGY_INITIATIVES_TABLE_NAME")
	strategyInitiativesPlatformIndex            = "StrategyInitiativesPlatformIndex"
	strategyObjectivesTableName                 = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_TABLE")
	strategyObjectivesPlatformIndex             = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_PLATFORM_INDEX")
	strategyObjectivesCapabilityCommunityIndex  = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_CAPABILITY_COMMUNITY_INDEX")
	strategyInitiativesInitiativeCommunityIndex = "StrategyInitiativesInitiativeCommunityIndex"
)

var (
	d   = awsutils.NewDynamo(region, "", namespace)
	s   = awsutils.NewSNS(region, "", namespace)
	dns = common.DynamoNamespace{Dynamo: d, Namespace: namespace}

	clientID = utils.NonEmptyEnv("CLIENT_ID")
	schema   = models.SchemaForClientID(clientID)
	userDAO  = utilsUser.NewDAOFromSchema(d, namespace, schema)

	sns                   = awsutils.NewSNS(region, "", namespace)
	valuesTableName       = utils.NonEmptyEnv("VALUES_TABLE_NAME")
	valuesPlatformIdIndex = utils.NonEmptyEnv("VALUES_PLATFORM_ID_INDEX")

	logger = alog.LambdaLogger(logrus.InfoLevel)

	platformInstance = utils.Platform{
		Sns:                       *sns,
		PlatformNotificationTopic: platformNotificationTopic,
		Namespace:                 namespace,
		IsInteractiveDebugEnabled: false,
	}
	platformDAO      = utilsPlatform.NewDAOFromSchema(d, namespace, schema)
	dialogFetcherDAO = dialogFetcher.NewDAO(d, dialogTableName)
	connGen          = daosCommon.DynamoDBConnectionGen{
		Dynamo:          d,
		TableNamePrefix: clientID,
	}
)

// UserIDsToDisplayNames converts a bunch of user ids to their names
// NB! O(n)! TODO: implement a query that returns many users at once.
func UserIDsToDisplayNames(userIDs []string) (res []models.KvPair) {
	if len(userIDs) > 10 {
		fmt.Println("WARN: Very slow user data fetching")
	}
	for _, userID := range userIDs {
		user := userDAO.ReadUnsafe(userID)
		if !user.IsAdaptiveBot {
			res = append(res, models.KvPair{Key: user.DisplayName, Value: userID})
		}
	}
	return
}

func userTokenSyncUnsafe(userID string) string {
	token, err2 := utilsPlatform.GetTokenForUser(d, clientID, userID)
	core.ErrorHandler(err2, "userTokenSyncUnsafe", "GetTokenForUser")
	return token
}

func getSlackClient(request slack.InteractionCallback) *slack.Client {
	ut := userTokenSyncUnsafe(request.User.ID)
	return slack.New(ut)
}
