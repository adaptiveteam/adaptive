package coaching

import (
	"time"
	// "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/daos/userFeedback"
	"fmt"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	business_time "github.com/adaptiveteam/adaptive/business-time"
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
// FeedbackReceivedForTheQuarter -
func FeedbackReceivedForTheQuarter(userID string, quarter, year int) func (conn daosCommon.DynamoDBConnection) (res []models.UserFeedback, err error) {
	return func (conn daosCommon.DynamoDBConnection) (res []models.UserFeedback, err error) {
		qyr := fmt.Sprintf("%d:%d", quarter, year)
		res, err = userFeedback.ReadByQuarterYearTarget(qyr, userID)(conn)
		return 
	}
}
// ReportExists checks if feedback report exists in a bucket with a key
func ReportExists(bucket, key string) bool {
	return common.DeprecatedGetGlobalS3().ObjectExists(bucket, key)
}

// UserReportIDForPreviousQuarter constructs key to look for in S3 for a user for the last quarter
func UserReportIDForPreviousQuarter(date time.Time, reportForUserID string) (key string, err error) {
	var y, m, d = date.Date()
	bt := business_time.NewDate(y, int(m), d)
	year := bt.GetPreviousQuarterYear()
	quarter := bt.GetPreviousQuarter()
	key = fmt.Sprintf("%s/%d/%d/%s", reportForUserID, year, quarter, ReportName)
	return
}

// ReportFor returns the user or channel, where to post report in
func ReportFor(userID, targetID string) (reportFor string) {
	if targetID == "" {
		reportFor = userID
	} else {
		reportFor = targetID
	}
	return
}
