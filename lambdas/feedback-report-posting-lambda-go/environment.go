package feedbackReportPostingLambda

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/adaptiveteam/adaptive/daos/common"
)

var (
	namespace                 = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                    = utils.NonEmptyEnv("AWS_REGION")
	userProfileLambda         = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	FeedbackReportingLambdaName   = utils.NonEmptyEnv("FEEDBACK_REPORTING_LAMBDA_NAME")
	reportsBucket             = utils.NonEmptyEnv("REPORTS_BUCKET_NAME")
	platformNotificationTopic = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	sns                       = awsutils.NewSNS(region, "", namespace)
	s                         = awsutils.NewS3(region, "", namespace)
	l                         = awsutils.NewLambda(region, "", namespace)
	D                         = awsutils.NewDynamo(region, "", namespace)
	clientID                  = utils.NonEmptyEnv("CLIENT_ID")
	schema    = models.SchemaForClientID(clientID)

	reportName = "performance_report.pdf"
	connGen = common.CreateConnectionGenFromEnv()
)
