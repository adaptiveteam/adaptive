package feedbackReportingLambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adaptiveteam/adaptive/daos/userFeedback"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	apr "github.com/adaptiveteam/adaptive/adaptive-reports/performance-report"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	_ "github.com/adaptiveteam/adaptive/daos"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	daosUser "github.com/adaptiveteam/adaptive/daos/user"
	fetch_dialog "github.com/adaptiveteam/adaptive/dialog-fetcher"
	ls "github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

type Coaching = apr.Coaching

func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func publish(msg models.PlatformSimpleNotification) {
	_, err := sns.Publish(msg, platformNotificationTopic)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not pusblish message to %s topic",
		platformNotificationTopic))
}

func publishAll(notes []models.PlatformSimpleNotification) {
	for _, note := range notes {
		publish(note)
	}
}

var (
	namespace                 = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                    = utils.NonEmptyEnv("AWS_REGION")
	D                         = awsutils.NewDynamo(region, "", namespace)
	s                         = awsutils.NewS3(region, "", namespace)
	clientID                  = utils.NonEmptyEnv("CLIENT_ID")
	table                     = userFeedback.TableName(clientID)
	reportBucket              = utils.NonEmptyEnv("FEEDBACK_REPORTS_BUCKET_NAME")
	userProfileLambda         = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	feedbackTargetIndex       = string(userFeedback.QuarterYearTargetIndex)
	platformNotificationTopic = utils.NonEmptyEnv("PLATFORM_NOTIFICATION_TOPIC")
	sns                       = awsutils.NewSNS(region, "", namespace)
	dialogTable               = utils.NonEmptyEnv("DIALOG_TABLE")

	dns       = common.DynamoNamespace{Dynamo: D, Namespace: namespace}
	schema    = models.SchemaForClientID(clientID)
	logger    = alog.LambdaLogger(logrus.InfoLevel)
	connGen   = daosCommon.CreateConnectionGenFromEnv()
)

func HandleRequest(ctx context.Context, engage models.UserEngage) (coachings []Coaching, err error) {
	defer core.RecoverToErrorVar("feedback-reporting-lambda", &err)
	logger = logger.WithLambdaContext(ctx)
	// userID := engage.UserID
	targetID := engage.TargetID
	date := engage.Date
	threadTs := engage.ThreadTs
	channel := engage.Channel
	// A user can request a report and it can also be requested by a user in a community
	var reportFor = coaching.ReportFor(engage.UserID, engage.TargetID)
	// When request comes from a channel, we should respond back to the channel
	// We treat this channel as a user, as in we have profile information for this channel
	// if engage.Channel != "" {
	//	sendTo = engage.Channel
	// }
	quarter, year := getQuarterYearForDateOrElseNow(date)
	engs := selectReceivedFeedbackUnsafe(reportFor, quarter, year)

	for _, each := range engs {
		coachings = append(coachings, convertUserFeedbackToCoaching(each))
	}
	conn := connGen.ForPlatformID(engage.TeamID.ToPlatformID())
	coachings = apr.ResolveCompetencies(coachings, apr.GetCompetencyImpl(conn))
	
	// We post the generation status only if the request is from a community. In that case, target is not empty
	postCondition := targetID != "" && threadTs != ""

	notes := []models.PlatformSimpleNotification{}
	msg := ""
	if len(coachings) > 0 {
		filepath := fmt.Sprintf("/tmp/%s.pdf", reportFor)
		user := daosUser.ReadUnsafe(conn.PlatformID, reportFor)(conn)
		_, err = apr.BuildReportWithCustomValuesTyped(
			coachings, user.DisplayName, quarter, year,
			filepath,
			fetch_dialog.NewDAO(dns.Dynamo, dialogTable), logger,
		)
		if err == nil {
			defer deleteFile(filepath)
			s3Key := fmt.Sprintf("%s/%d/%d/performance_report.pdf", reportFor, year, quarter)
			err = s.AddFile(filepath, reportBucket, s3Key)
			if err == nil {
				msg = fmt.Sprintf("_<@%s>'s performance report for quarter `%d` of year `%d` has been generated._", reportFor, quarter, year)
			}
		}
	} else {
		msg = fmt.Sprintf("_Report not generated. <@%s> did not receive any feedback for quarter `%d` of year `%d`_", reportFor, quarter, year)
	}
	if postCondition && msg != "" {
		notes = append(notes, models.PlatformSimpleNotification{
			Message: msg,
			UserId: engage.UserID, 
			Channel: channel, 
			ThreadTs: threadTs,
		})
		publishAll(notes)
	}
	if err != nil {
		logger.WithError(err).Errorf("Error with collaboration report generation for %s user", reportFor)
	}
	return coachings, nil
}

func selectReceivedFeedbackUnsafe(reportFor string, quarter, year int) (received []models.UserFeedback) {
	err2 := D.QueryTableWithIndex(table, awsutils.DynamoIndexExpression{
		IndexName: feedbackTargetIndex,
		// there is no != operator for ConditionExpression
		Condition: "quarter_year = :qy AND target = :t",
		Attributes: map[string]interface{}{
			":t":  reportFor,
			":qy": fmt.Sprintf("%d:%d", quarter, year),
		},
	}, map[string]string{}, true, -1, &received)
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not query %s index for feedback", feedbackTargetIndex))
	return
}

func convertUserFeedbackToCoaching(uf models.UserFeedback) apr.Coaching {
	qy := strings.Split(uf.QuarterYear, ":")
	q, _ := strconv.Atoi(qy[0])
	y, _ := strconv.Atoi(qy[1])
	cf, _ := strconv.Atoi(uf.ConfidenceFactor)
	return Coaching{
		Source:   uf.Source,
		Target:   uf.Target,
		Topic:    uf.ValueID,
		// Type:     getCompetency(uf.ValueID),
		Rating:   float64(cf),
		Comments: uf.Feedback,
		Quarter:  q,
		Year:     y,
	}
}

func getQuarterYearForDateOrElseNow(date string) (quarter, year int) {
	var t time.Time
	if date == "" {
		t = time.Now()
		fmt.Printf("Date not present in UserEngage, using the date of current time %v", t)
	} else {
		fmt.Printf("Date is present in UserEngage.Date=%s", date)
		var err2 error
		t, err2 = core.ISODateLayout.Parse(date)
		core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not parse %s as date", date))
	}
	var y, m, d = t.Date()
	bt := business_time.NewDate(y, int(m), d)
	logger.Infof("Date %v", bt)
	quarter = bt.GetPreviousQuarter()
	year = bt.GetPreviousQuarterYear()
	fmt.Println(fmt.Sprintf("### quarter: %d, year: %d", quarter, year))
	return
}

func main() {
	ls.Start(HandleRequest)
}
