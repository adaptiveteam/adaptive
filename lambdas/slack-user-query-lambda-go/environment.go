package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/sirupsen/logrus"
)

var (
	region    = utils.NonEmptyEnv("AWS_REGION")
	namespace = utils.NonEmptyEnv("LOG_NAMESPACE")

	clientID = utils.NonEmptyEnv("CLIENT_ID")

	d      = awsutils.NewDynamo(region, "", namespace)
	schema = models.SchemaForClientID(clientID)

	logger = alog.LambdaLogger(logrus.InfoLevel)

	userCommunityTable           = utils.NonEmptyEnv("USER_COMMUNITY_TABLE_NAME")
	userCommunityPlatformIndex   = string(adaptiveCommunity.PlatformIDIndex)
	communityUsersTable          = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	communityUsersCommunityIndex = adaptiveCommunityUser.PlatformIDCommunityIDIndex

)
