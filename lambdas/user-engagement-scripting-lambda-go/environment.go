package lambda

import (
	"fmt"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	plat "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/sirupsen/logrus"
)

var (
	namespace                           = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                              = utils.NonEmptyEnv("AWS_REGION")
	clientID                            = utils.NonEmptyEnv("CLIENT_ID")
	userEngagementSchedulerLambdaPrefix = utils.NonEmptyEnv("USER_ENGAGEMENT_SCHEDULER_LAMBDA_PREFIX")
	engagementAnsweredIndex             = "UserIDAnsweredIndex"
	userEngagementSchedulerLambda       = fmt.Sprintf("%s_%s", clientID, userEngagementSchedulerLambdaPrefix)
	d                                   = awsutils.NewDynamo(region, "", namespace)
	l                                   = awsutils.NewLambda(region, "", namespace)
	engagementTable                     = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	platformNotificationTopic           = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	sns                                 = awsutils.NewSNS(region, "", namespace)
	schema                              = models.SchemaForClientID(clientID)
	platformTokenDao                    = plat.NewDAOFromSchema(d, namespace, schema)

	logger = alog.LambdaLogger(logrus.InfoLevel)
)
