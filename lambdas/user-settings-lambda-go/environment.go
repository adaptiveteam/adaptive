package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/userEngagement"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	user "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

var (
	namespace                 = utils.NonEmptyEnv("LOG_NAMESPACE")
	//usersTable                = utils.NonEmptyEnv("USERS_TABLE_NAME")
	region                    = utils.NonEmptyEnv("AWS_REGION")
	_                         = utils.NonEmptyEnv("USER_SETUP_LAMBDA_NAME")
	engTable                  = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	platformNotificationTopic = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	d                         = awsutils.NewDynamo(region, "", namespace)
	sns                       = awsutils.NewSNS(region, "", namespace)

	dns      = common.DynamoNamespace{Dynamo: d, Namespace: namespace}
	platform = utils.Platform{
		Sns:                       *sns,
		PlatformNotificationTopic: platformNotificationTopic,
		Namespace:                 namespace,
		IsInteractiveDebugEnabled: false,
	}

	userEngagementDao = userEngagement.NewDAO(dns, engTable)
	clientID                  = utils.NonEmptyEnv("CLIENT_ID")

	schema           = models.SchemaForClientID(clientID)

	userDao = user.NewDAOFromSchema(d, namespace, schema)
)