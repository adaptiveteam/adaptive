package feedbackReportingLambda

import (
	"encoding/json"
	"time"

	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// GeneratePerformanceReportAndPostToUserAsync invokes FEEDBACK_REPORTING_LAMBDA_NAME
// in order to generate performance report based on received feedback.
// After generation this lambda will post a notification that the report is ready.
// Probably it shouldn't.
func GeneratePerformanceReportAndPostToUserAsync(
	// teamID models.TeamID,
	reportFor string,
	date time.Time,
) (err error) {
	defer core.RecoverToErrorVar("GeneratePerformanceReportAndPostToUserAsync", &err)
	engage := models.UserEngage{
		UserID: reportFor,
		Date:   core.ISODateLayout.Format(date), // Date: date.Format(time.RFC3339)
		// TeamID: teamID,
		// Update: true, not used!
	}
	var userEngageBytes []byte
	userEngageBytes, err = json.Marshal(engage)

	if err == nil {
		namespace := utils.NonEmptyEnv("LOG_NAMESPACE")
		region := utils.NonEmptyEnv("AWS_REGION")
		l := awsutils.NewLambda(region, "", namespace)
		feedbackReportingLambdaName := utils.NonEmptyEnv("FEEDBACK_REPORTING_LAMBDA_NAME")

		_, err = l.InvokeFunction(feedbackReportingLambdaName, userEngageBytes, true)
	}
	return
}

// func GeneratePerformanceReportAndPostToChannelAsync(
// 	reportFor string,
// 	date time.Time,
// 	channel, threadMs string,
// ) (err error) {
// 	defer core.RecoverToErrorVar("GeneratePerformanceReportAndPostToChannelAsync", &err)
// 	userEngageByt, _ := json.Marshal(models.UserEngage{
// 		UserID: reportFor,
// 		IsNew:  false, Update: true, Date: date.Format(time.RFC3339)})
// 	_, err := L.InvokeFunction(FeedbackReportingLambdaName, userEngageByt, true)
// 	engage := models.UserEngage{
// 		UserID: userID,
// 		Date:   core.ISODateLayout.Format(date),
// 		TeamID: teamID,
// 		// Update: true,
// 	}
// 	var userEngageBytes []byte
// 	userEngageBytes, err = json.Marshal(engage)
// 	if err == nil {
// 		namespace := utils.NonEmptyEnv("LOG_NAMESPACE")
// 		region := utils.NonEmptyEnv("AWS_REGION")
// 		l := awsutils.NewLambda(region, "", namespace)
// 		feedbackReportingLambdaName := utils.NonEmptyEnv("FEEDBACK_REPORTING_LAMBDA_NAME")

// 		_, err = l.InvokeFunction(feedbackReportingLambdaName, userEngageBytes, true)
// 	}
// 	return
// }
