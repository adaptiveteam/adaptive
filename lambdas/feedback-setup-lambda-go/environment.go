package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	evalues "github.com/adaptiveteam/adaptive/adaptive-engagements/values"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/feedback"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsPlatform "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/userFeedback"
	"github.com/nlopes/slack"
)

var (
	namespace                        = utils.NonEmptyEnv("LOG_NAMESPACE")
	clientID                         = utils.NonEmptyEnv("CLIENT_ID")
	d                                = awsutils.NewDynamo(utils.NonEmptyEnv("AWS_REGION"), "", namespace)
	feedbackTable                    = userFeedback.TableName(clientID)
	engagementTable                  = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	feedbackAnalysisLambda           = utils.NonEmptyEnv("FEEDBACK_ANALYSIS_LAMBDA")
	region                           = utils.NonEmptyEnv("AWS_REGION")
	userProfileLambda                = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	platformNotificationTopic        = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	collaborationReportPostingLambda = utils.NonEmptyEnv("FEEDBACK_REPORT_POSTING_LAMBDA_NAME")
	// feedbackEngagementLambda         = utils.NonEmptyEnv("FEEDBACK_ENGAGEMENT_LAMBDA_NAME")
	sns = awsutils.NewSNS(region, "", namespace)
	l   = awsutils.NewLambda(region, "", namespace)

	dns       = common.DynamoNamespace{Dynamo: d, Namespace: namespace}
	schema    = models.SchemaForClientID(clientID)
	valuesDao = evalues.NewDAOFromSchema(&dns, schema)

	platform = utils.Platform{
		Sns:                       *sns,
		PlatformNotificationTopic: platformNotificationTopic,
		Namespace:                 namespace,
		IsInteractiveDebugEnabled: false,
	}
	feedbackDAO = feedback.NewDAOFromSchema(d, namespace, schema)
	userDao     = utilsUser.NewDAOFromSchema(d, namespace, schema)
	connGen     = daosCommon.DynamoDBConnectionGen{
		Dynamo:          d,
		TableNamePrefix: clientID,
	}
)

func userTokenSyncUnsafe(userID string) string {
	token, err2 := utilsPlatform.GetTokenForUser(d, clientID, userID)
	core.ErrorHandler(err2, "userTokenSyncUnsafe", "GetTokenForUser")
	return token
}

func getSlackClient(userID string) *slack.Client {
	ut := userTokenSyncUnsafe(userID)
	return slack.New(ut)
}

func getTeamID(request slack.InteractionCallback) models.TeamID {
	user := userDao.ReadUnsafe(request.User.ID)
	return models.ParseTeamID(user.PlatformID)
}
