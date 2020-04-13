package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"
	"context"
	"fmt"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	ls "github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

var (
	logger = alog.LambdaLogger(logrus.InfoLevel)
)

func HandleRequest(ctx context.Context, uf models.UserFeedback) (err error) {
	logger = logger.WithLambdaContext(ctx)
	defer core.RecoverToErrorVar("feedback-analysis-lambda-go", &err)
	// Parse callback to MessageCallback
	var mc *models.MessageCallback
	mc, err = utils.ParseToCallback(uf.ID)
	core.ErrorHandler(err, namespace, "Could not parse to callback")

	// We post the analysis back to the user with an 'Edit' button. This edit is treated the same way as the original
	// feedback edit. So, we append 'ask_' as prefix to this attachment
	mc.Set("Action", fmt.Sprintf("ask_%s", uf.ValueID))

	conn := connGen.ForPlatformID(uf.PlatformID)
	var values [] models.AdaptiveValue
	values, err = adaptiveValue.ReadOrEmpty(uf.ValueID)(conn)
	values = adaptiveValue.AdaptiveValueFilterActive(values)
	found := len(values) > 0
	if err == nil && found {
		value := values[0]
		attach, err := FeedbackAttachmentTemplate(*mc, uf, value)
		if err == nil {
			utils.ECAnalysis(uf.Feedback, FeedbackDialogContext, "Feedback", dialogTable,
				mc.ToCallbackID(), uf.Source, uf.ChannelID, uf.MsgTimestamp, uf.MsgTimestamp,
				[]ebm.Attachment{attach},
				sns, platformNotificationTopic, namespace)
		} else {
			logger.WithField("error", err).Errorf("Could not analyze the feedback text from %s user", uf.Source)
		}
	} else if err != nil {
		logger.WithField("error", err).Errorf("Could not read value with id %s", uf.ValueID)
	} else {
		logger.Errorf("Value with id %s not found", uf.ValueID)
	}
	return
}

func main() {
	ls.Start(HandleRequest)
}
