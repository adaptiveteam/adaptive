package lambda

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	"github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
	"github.com/sirupsen/logrus"
	utilsPlatform "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
)

var (
	region                      = utils.NonEmptyEnv("AWS_REGION")
	namespace                   = utils.NonEmptyEnv("LOG_NAMESPACE")
	l                           = awsutils.NewLambda(region, "", namespace)
	feedbackLambda              = utils.NonEmptyEnv("USER_FEEDBACK_SETUP_LAMBDA_NAME")
	feedbackReportingLambda     = utils.NonEmptyEnv("FEEDBACK_REPORTING_LAMBDA_NAME")
	feedbackReportPostingLambda = utils.NonEmptyEnv("FEEDBACK_REPORT_POSTING_LAMBDA_NAME")
	platformNotificationTopic   = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	engagementsTable            = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	sns                         = awsutils.NewSNS(region, "", namespace)
	d                           = awsutils.NewDynamo(region, "", namespace)

	communitiesTable             = utils.NonEmptyEnv("USER_COMMUNITY_TABLE")
	communityPlatformIndex       = utils.NonEmptyEnv("USER_COMMUNITY_PLATFORM_INDEX")
	communityUsersTable          = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE")
	communityUsersCommunityIndex = utils.NonEmptyEnv("COMMUNITY_USERS_COMMUNITY_INDEX")

	clientID = utils.NonEmptyEnv("CLIENT_ID")

	userObjectivesTable = utils.NonEmptyEnv("USER_OBJECTIVES_TABLE_NAME")
	userObjectiveDAO    = userObjective.NewDAOByTableName(d, namespace, userObjectivesTable)

	userObjectiveProgressTable = utils.NonEmptyEnv("USER_OBJECTIVES_PROGRESS_TABLE")
	userObjectiveProgressDAO = userObjectiveProgress.NewDAOByTableName(d, namespace, userObjectiveProgressTable)

	schema  = models.SchemaForClientID(clientID)
	userDao = utilsUser.NewDAOFromSchema(d, namespace, schema)

	logger = alog.LambdaLogger(logrus.InfoLevel)
	platformDAO         = utilsPlatform.NewDAOFromSchema(d, namespace, schema)
)

func filterObjectivesByObjectiveType(objectives []userObjective.UserObjective, objectiveType userObjective.DevelopmentObjectiveType) (res []userObjective.UserObjective) {
	for _, objective := range objectives {
		if objective.ObjectiveType == objectiveType {
			res = append(res, objective)
		}
	}
	return
}
