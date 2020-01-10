package lambda

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	plat "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	// "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	// dialogFetcher "github.com/adaptiveteam/adaptive/dialog-fetcher"
)

var (
	region    = utils.NonEmptyEnv("AWS_REGION")
	namespace = utils.NonEmptyEnv("LOG_NAMESPACE")

	platformNotificationTopic = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	clientID                  = utils.NonEmptyEnv("CLIENT_ID")

	d      = awsutils.NewDynamo(region, "", namespace)
	schema = models.SchemaForClientID(clientID)

	platformTokenDao = plat.NewDAOFromSchema(d, namespace, schema)
	platformAdapter  = mapper.SlackAdapter2(platformTokenDao)
	// instead of profile lambda
	userDao = user.NewDAOFromSchema(d, namespace, schema)
)
