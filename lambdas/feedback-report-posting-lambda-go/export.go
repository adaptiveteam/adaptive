package feedbackReportPostingLambda

import (
	"encoding/json"
	"time"

	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// DeliverReportToUserAsync sends the previously generated report to user.
// If the report wasn't generated before, it is regenerated using FeedbackReportingLambda.
func DeliverReportToUserAsync(
	teamID models.TeamID, 
	reportForUserID string, 
	date time.Time) (err error) {
	defer core.RecoverToErrorVar("DeliverReportToUserAsync", &err)
	engage := models.UserEngage{
		UserID: reportForUserID,
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
		feedbackReportPostingLambdaName := utils.NonEmptyEnv("FEEDBACK_REPORT_POSTING_LAMBDA_NAME")

		_, err = l.InvokeFunction(feedbackReportPostingLambdaName, userEngageBytes, true)
	}
	return
}
