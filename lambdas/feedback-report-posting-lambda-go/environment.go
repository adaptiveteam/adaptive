package feedbackReportPostingLambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	utilsPlatform "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/nlopes/slack"
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
)

func userTokenSyncUnsafe(userID string) string {
	token, err2 := utilsPlatform.GetTokenForUser(D, clientID, userID)
	core.ErrorHandler(err2, "userTokenSyncUnsafe", "GetTokenForUser")
	return token
}

func getSlackClient(userID string) *slack.Client {
	ut := userTokenSyncUnsafe(userID)
	return slack.New(ut)
}
