package lambda

import (
	"github.com/pkg/errors"
	"time"
	"github.com/adaptiveteam/adaptive/daos/user"
	"bytes"
	"context"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	// utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	ls "github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	feedbackReportingLambda "github.com/adaptiveteam/adaptive/lambdas/feedback-reporting-lambda-go"
)

var (
	logger = alog.LambdaLogger(logrus.InfoLevel)
)

func publish(msg models.PlatformSimpleNotification) {
	_, err := sns.Publish(msg, platformNotificationTopic)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not pusblish message to %s topic", platformNotificationTopic))
}

func postTo(engage models.UserEngage) (postTo string) {
	if engage.Channel == "" {
		postTo = engage.UserID
	} else {
		postTo = engage.Channel
	}
	return
}

func downloadReportContents(reportForUserID string, t time.Time) (contents []byte, err error) {
	var key string
	key, err = coaching.UserReportIDForPreviousQuarter(t, reportForUserID)
	if err == nil {
		logger.Infof("Report id for %s user: %s", reportForUserID, key)
		// Check if the report exists
		if !common.DeprecatedGetGlobalS3().ObjectExists(reportsBucket, key) {
			logger.Infof("Report doesn't exist for %s user, generating now", reportForUserID)
			err = feedbackReportingLambda.GeneratePerformanceReportAndPostToUserAsync(reportForUserID, t)
			err = errors.Wrapf(err, "Could not invoke %s lambda", FeedbackReportingLambdaName)
		}
		if err == nil {
			logger.Infof("Report exists for %s user", reportForUserID)
			contents, err = s.GetObject(reportsBucket, key)
		}
	}
	return
}

func sendReport(engage models.UserEngage, targetUserDisplayName string, contents []byte) (err error) {
	reportFor := coaching.ReportFor(engage.UserID, engage.TargetID)
	postTo := postTo(engage)
	// Upload the file only for non-empty content
	if len(contents) > 0 {
		if err == nil {
			params := slack.FileUploadParameters{
				Title:           string(TitleTemplate(targetUserDisplayName)),
				Filename:        reportName,
				Reader:          bytes.NewBuffer(contents),
				Channels:        []string{postTo},
				ThreadTimestamp: engage.ThreadTs,
			}
			api := getSlackClient(reportFor)
			_, err = api.UploadFile(params)
		}
	} else {
		publish(models.PlatformSimpleNotification{
			UserId: reportFor, Channel: postTo, ThreadTs: engage.ThreadTs,
			Message: string(NoReportTemplate),
		})
	}
	return
}

func HandleRequest(ctx context.Context, engage models.UserEngage) {
	logger = logger.WithLambdaContext(ctx)
	logger.WithField("payload", engage).Infof("Starting...")
	reportForUserID := coaching.ReportFor(engage.UserID, engage.TargetID)
	t, err2 := core.ISODateLayout.Parse(engage.Date)

	if err2 == nil {
		var contents []byte
		contents, err2 = downloadReportContents(reportForUserID, t)
		if err2 == nil {
			var targetUserDisplayName string
			targetUserDisplayName, err2 = getDisplayName(engage.TeamID, reportForUserID)
			if err2 == nil {
				err2 = sendReport(engage, targetUserDisplayName, contents)
			}
		} 
	}
	if err2 != nil {
		logger.WithField("error", err2).Errorf("Couldn't load report for %s user", engage.UserID)
	}
	return
}

func main() {
	ls.Start(HandleRequest)
}

func getDisplayName(teamID models.TeamID, userID string) (displayName string, err error) {
	connGen := daosCommon.CreateConnectionGenFromEnv()
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	var u user.User
	u, err = user.Read(userID)(conn)
	if err == nil {
		displayName = u.DisplayName
	}
	return
}
