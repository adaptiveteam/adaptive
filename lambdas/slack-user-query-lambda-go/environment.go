package lambda

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	plat "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/daos/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/sirupsen/logrus"
)

var (
	region    = utils.NonEmptyEnv("AWS_REGION")
	namespace = utils.NonEmptyEnv("LOG_NAMESPACE")

	clientID = utils.NonEmptyEnv("CLIENT_ID")

	d      = awsutils.NewDynamo(region, "", namespace)
	schema = models.SchemaForClientID(clientID)

	platformTokenDao = plat.NewDAOFromSchema(d, namespace, schema)
	// instead of profile lambda
	userDao = user.NewDAOByTableName(d, namespace, schema.AdaptiveUsers.Name)

	logger = alog.LambdaLogger(logrus.InfoLevel)

	userCommunityTable           = utils.NonEmptyEnv("USER_COMMUNITY_TABLE_NAME")
	userCommunityPlatformIndex   = utils.NonEmptyEnv("USER_COMMUNITY_PLATFORM_INDEX")
	communityUsersTable          = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	communityUsersCommunityIndex = utils.NonEmptyEnv("COMMUNITY_USERS_COMMUNITY_INDEX")
)
