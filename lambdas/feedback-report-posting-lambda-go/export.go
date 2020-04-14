package lambda

import (
	"encoding/json"
	"time"

	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// PostReportToUserAsync sends the previously generated report to user
func PostReportToUserAsync(teamID models.TeamID, userID string, date time.Time) (err error) {
	defer core.RecoverToErrorVar("PostReportToUser", &err)
	engage := models.UserEngage{
		UserID: userID,
		Date:   core.ISODateLayout.Format(date),
		TeamID: teamID,
		// Update: true,
	}
	var userEngageBytes []byte
	userEngageBytes, err = json.Marshal(engage)
	if err == nil {
		namespace := utils.NonEmptyEnv("LOG_NAMESPACE")
		region    := utils.NonEmptyEnv("AWS_REGION")
		l         := awsutils.NewLambda(region, "", namespace)
		feedbackReportingLambdaName := utils.NonEmptyEnv("FEEDBACK_REPORTING_LAMBDA_NAME")

		_, err = l.InvokeFunction(feedbackReportingLambdaName, userEngageBytes, true)
	}
	return
}
