package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/userEngagement"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

var (
	MeetingTimeUserAttributeID  = "meeting_time"
	KnownUserSettingsAttributes = []string{MeetingTimeUserAttributeID}
)

var (
	namespace = utils.NonEmptyEnv("LOG_NAMESPACE")
	region    = utils.NonEmptyEnv("AWS_REGION")
	engTable  = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	//attribsTable = utils.NonEmptyEnv("USER_ATTRIBUTES_TABLE_NAME")
	usersTable = utils.NonEmptyEnv("USERS_TABLE_NAME")
	//attribsUserIDIndexName    = utils.NonEmptyEnv("USER_ATTRIBUTES_TABLE_USER_ID_INDEX_NAME")
	platformNotificationTopic = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")

	d        = awsutils.NewDynamo(region, "", namespace)
	sns      = awsutils.NewSNS(region, "", namespace)
	dns      = common.DynamoNamespace{Dynamo: d, Namespace: namespace}
	platform = utils.Platform{
		Sns:                       *sns,
		PlatformNotificationTopic: platformNotificationTopic,
		Namespace:                 namespace,
		IsInteractiveDebugEnabled: false,
	}
	userEngagementDao = userEngagement.NewDAO(dns, engTable)
)
