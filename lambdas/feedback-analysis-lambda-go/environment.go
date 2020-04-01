package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	evalues "github.com/adaptiveteam/adaptive/adaptive-engagements/values"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

const (
	FeedbackDialogContext = "dialog/feedback/language-coaching"
)

var (
	namespace                 = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                    = utils.NonEmptyEnv("AWS_REGION")
	dialogTable               = utils.NonEmptyEnv("DIALOG_TABLE")
	platformNotificationTopic = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	sns                       = awsutils.NewSNS(region, "", namespace)

	clientID  = utils.NonEmptyEnv("CLIENT_ID")
	d         = awsutils.NewDynamo(utils.NonEmptyEnv("AWS_REGION"), "", namespace)
	dns       = common.DynamoNamespace{Dynamo: d, Namespace: namespace}
	schema    = models.SchemaForClientID(clientID)
	valuesDao = evalues.NewDAOFromSchema(&dns, schema)
)
