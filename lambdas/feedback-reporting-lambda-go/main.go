package lambda

import (
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
	userDAO   = daosUser.NewDAOByTableName(D, namespace, schema.AdaptiveUsers.Name)
	connGen   = daosCommon.CreateConnectionGenFromEnv()
)

func HandleRequest(ctx context.Context, engage models.UserEngage) (coachings []Coaching, err error) {
	logger = logger.WithLambdaContext(ctx)
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("error in feedback-reporting-lambda %v", err2)
		}
	}()
	userID := engage.UserID
	targetID := engage.TargetID
	date := engage.Date
	threadTs := engage.ThreadTs
	channel := engage.Channel
	conn := connGen.ForPlatformID(engage.TeamID.ToPlatformID())
	var reportFor = userID
	var sendTo = userID
	// A user can request a report and it can also be requested by a user in a community
	if targetID != "" {
		reportFor = targetID
	}
	logger.Infof("Got the user and target")
	// When request comes from a channel, we should respond back to the channel
	// We treat this channel as a user, as in we have profile information for this channel
	// if engage.Channel != "" {
	//	sendTo = engage.Channel
	// }
	var t time.Time
	if date != "" {
		fmt.Printf("Date is present in UserEngage.Date=%s", date)
		t, err = core.ISODateLayout.Parse(date)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse %s as date", date))
	} else {
		t = time.Now()
		fmt.Printf("Date not present in UserEngage, using the date of current time %v", t)
	}
	var y, m, d = t.Date()
	bt := business_time.NewDate(y, int(m), d)
	logger.Infof("Date %v", bt)
	quarter := bt.GetPreviousQuarter()
	year := bt.GetPreviousQuarterYear()
	fmt.Println(fmt.Sprintf("### quarter: %d, year: %d", quarter, year))

	engs := selectReceivedFeedbackUnsafe(reportFor, quarter, year)

	for _, each := range engs {
		coachings = append(coachings, convertUserFeedbackToCoaching(each))
	}
	coachings = apr.ResolveCompetencies(coachings, apr.GetCompetencyImpl(conn))
	
	// We post the generation status only if the request is from a community. In that case, target is not empty
	postCondition := targetID != "" && threadTs != ""

	if len(coachings) > 0 {
		filepath := fmt.Sprintf("/tmp/%s.pdf", userID)
		user := userDAO.ReadUnsafe(userID)
		_, err = apr.BuildReportWithCustomValuesTyped(coachings, []apr.Coaching{}, user.DisplayName, quarter, year, filepath,
			fetch_dialog.NewDAO(dns.Dynamo, dialogTable), logger)
		if err == nil {
			err = s.AddFile(filepath, reportBucket, fmt.Sprintf("%s/%d/%d/performance_report.pdf", reportFor, year,
				quarter))
			if err == nil {
				if postCondition {
					publish(models.PlatformSimpleNotification{UserId: sendTo, Channel: channel, Message: fmt.Sprintf(
						"_<@%s>'s performance report for quarter `%d` of year `%d` has been generated._", reportFor,
						quarter, year), ThreadTs: threadTs})
				}
			}
			deleteFile(filepath)
		}
	} else if postCondition {
		publish(models.PlatformSimpleNotification{UserId: sendTo, Channel: channel, Message: fmt.Sprintf(
			"_Report not generated. <@%s> did not receive any feedback for quarter `%d` of year `%d`_",
			reportFor, quarter, year), ThreadTs: threadTs})
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

func main() {
	ls.Start(HandleRequest)
}
