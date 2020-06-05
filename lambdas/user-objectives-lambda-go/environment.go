package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	dialogFetcher "github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/sirupsen/logrus"
)

var (
	platformNotificationTopic = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	namespace                 = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                    = utils.NonEmptyEnv("AWS_REGION")

	engagementTable                  = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	userTable                        = utils.NonEmptyEnv("USERS_TABLE_NAME")
	userProfileLambda                = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	userObjectivesTable              = utils.NonEmptyEnv("USER_OBJECTIVES_TABLE_NAME")        // needed
	userObjectivesProgressTable      = utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_TABLE")    // needed
	objectiveCloseoutPath            = ""//utils.NonEmptyEnv("USER_OBJECTIVES_CLOSEOUT_LEARN_MORE_PATH")
	dialogTableName                  = utils.NonEmptyEnv("DIALOG_TABLE")
	communityUsersTable              = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	communitiesTable                 = utils.NonEmptyEnv("USER_COMMUNITIES_TABLE")
	// partnershipRejectionsTable       = utils.NonEmptyEnv("ACCOUNTABILITY_PARTNERSHIP_REJECTION_TABLE")

	capabilityCommunitiesTableName     = utils.NonEmptyEnv("CAPABILITY_COMMUNITIES_TABLE_NAME")
	strategyInitiativeCommunitiesTable = utils.NonEmptyEnv("INITIATIVE_COMMUNITIES_TABLE_NAME")

	strategyInitiativesTableName                = utils.NonEmptyEnv("STRATEGY_INITIATIVES_TABLE_NAME")
	strategyObjectivesTableName                 = utils.NonEmptyEnv("STRATEGY_OBJECTIVES_TABLE")
)

var (
	d   = awsutils.NewDynamo(region, "", namespace)
	s   = awsutils.NewSNS(region, "", namespace)
	dns = common.DynamoNamespace{Dynamo: d, Namespace: namespace}

	clientID = utils.NonEmptyEnv("CLIENT_ID")
	schema   = models.SchemaForClientID(clientID)
	
	sns                   = awsutils.NewSNS(region, "", namespace)
	valueDao              = adaptiveValue.NewDAO(d, namespace, clientID)

	logger = alog.LambdaLogger(logrus.InfoLevel)

	platformInstance = utils.Platform{
		Sns:                       *sns,
		PlatformNotificationTopic: platformNotificationTopic,
		Namespace:                 namespace,
		IsInteractiveDebugEnabled: false,
	}
	// platformDAO      = utilsPlatform.NewDAOFromSchema(d, namespace, schema)
	dialogFetcherDAO = dialogFetcher.NewDAO(d, dialogTableName)
	connGen          = daosCommon.DynamoDBConnectionGen{
		Dynamo:          d,
		TableNamePrefix: clientID,
	}
)
