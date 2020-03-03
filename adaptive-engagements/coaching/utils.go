package coaching

import (
	"fmt"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	// utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// func UserDNPlatform(userId string, userProfileLambda, region string, dns common.DynamoNamespace) (string, models.TeamID) {
// 	ut, err := utils.UserToken(userId, userProfileLambda, region, dns.Namespace)
// 	core.ErrorHandler(err, dns.Namespace, fmt.Sprintf("Could not query for user token"))
// 	return ut.DisplayName, ut.ClientPlatformRequest.TeamID
// }

func Coach(table string, coachee string, quarter, year int, coachingCoacheeIndex string, dns common.DynamoNamespace) (
	[]models.CoachingRelationship, error) {
	var rels []models.CoachingRelationship
	err := dns.Dynamo.QueryTableWithIndex(table, awsutils.DynamoIndexExpression{
		IndexName: coachingCoacheeIndex,
		Condition: "coachee_quarter_year = :cqy",
		Attributes: map[string]interface{}{
			":cqy": fmt.Sprintf("%s:%d:%d", coachee, quarter, year),
		},
	}, map[string]string{}, true, -1, &rels)
	return rels, err
}

func CoachingRelsQuarterYear(table string, quarter, year int, coachingQYIndex string, dns common.DynamoNamespace) (
	[]models.CoachingRelationship, error) {
	var rels []models.CoachingRelationship
	err := dns.Dynamo.QueryTableWithIndex(table, awsutils.DynamoIndexExpression{
		IndexName: coachingQYIndex,
		// 'year' is a reserved keyword in dynamo
		Condition: "quarter = :q and #year = :y",
		Attributes: map[string]interface{}{
			":q": quarter,
			":y": year,
		},
	}, map[string]string{"#year": "year"}, true, -1, &rels)
	return rels, err
}

func Coachees(table string, coach string, quarter, year int, coachQYIndex string, dns common.DynamoNamespace) (
	[]models.CoachingRelationship, error) {
	var rels []models.CoachingRelationship
	err := dns.Dynamo.QueryTableWithIndex(table, awsutils.DynamoIndexExpression{
		IndexName: coachQYIndex,
		Condition: "coach_quarter_year = :cqy",
		Attributes: map[string]interface{}{
			":cqy": fmt.Sprintf("%s:%d:%d", coach, quarter, year),
		},
	}, map[string]string{}, true, -1, &rels)
	return rels, err
}

func FeedbackGivenForTheQuarter(userID string, quarter, year int, feedbackTable,
	feedbackSourceQYIndex string) ([]models.UserFeedback, error) {
	var res []models.UserFeedback
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(feedbackTable, awsutils.DynamoIndexExpression{
		IndexName: feedbackSourceQYIndex,
		Condition: "quarter_year = :qy and #source = :s",
		Attributes: map[string]interface{}{
			":s":  userID,
			":qy": fmt.Sprintf("%d:%d", quarter, year),
		},
	}, map[string]string{"#source": "source"}, true, -1, &res)
	return res, err
}

// ReportExists checks if feedback report exists in a bucket with a key
func ReportExists(bucket, key string) bool {
	return common.DeprecatedGetGlobalS3().ObjectExists(bucket, key)
}

// UserReportIDForPreviousQuarter constructs key to look for in S3 for a user for the last quarter
func UserReportIDForPreviousQuarter(engage models.UserEngage) (key string, err error) {
	t, err := core.ISODateLayout.Parse(engage.Date)
	if err == nil {
		var y, m, d = t.Date()
		bt := business_time.NewDate(y, int(m), d)
		year := bt.GetPreviousQuarterYear()
		quarter := bt.GetPreviousQuarter()
		key = fmt.Sprintf("%s/%d/%d/%s", ReportFor(engage), year, quarter, ReportName)
	}
	return
}

// ReportFor returns the user or channel, where to post report in
func ReportFor(engage models.UserEngage) (reportFor string) {
	if engage.TargetID == "" {
		reportFor = engage.UserID
	} else {
		reportFor = engage.TargetID
	}
	return
}
