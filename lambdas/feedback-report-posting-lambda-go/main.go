package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	// utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ls "github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
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

func loadReport(engage models.UserEngage) (contents []byte, err error) {
	key, err := coaching.UserReportIDForPreviousQuarter(engage)
	if err == nil {
		logger.Infof("Report id for %s user: %s", engage.UserID, key)
		// Check if the report exists
		if !coaching.ReportExists(reportsBucket, key) {
			logger.Infof("Report doesn't exist for %s user, generating now", engage.UserID)
			engBytes, _ := json.Marshal(engage)
			_, err = l.InvokeFunction(feedbackReportingLambda, engBytes, false)
		}
		if err == nil {
			logger.Infof("Report exists for %s user", engage.UserID)
			contents, err = s.GetObject(reportsBucket, key)
		}
	}
	return
}

func sendReport(engage models.UserEngage, contents []byte) (err error) {
	reportFor := coaching.ReportFor(engage)
	postTo := postTo(engage)
	// Upload the file only for non-empty content
	if len(contents) > 0 {
		if err == nil {
			params := slack.FileUploadParameters{
				Title:           string(TitleTemplate(reportFor)),
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
	contents, err := loadReport(engage)
	if err == nil {
		err = sendReport(engage, contents)
	} else {
		logger.WithField("error", err).Errorf("Couldn't load report for %s user", engage.UserID)
	}
	return
}

func main() {
	ls.Start(HandleRequest)
}
