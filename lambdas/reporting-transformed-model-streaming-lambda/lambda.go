package reporting_transformed_model_streaming_lambda

import (
	"context"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	sqc "github.com/adaptiveteam/adaptive/adaptive-utils-go/sql-connector"
	streamhandler "github.com/adaptiveteam/adaptive/lambdas/stream-handler"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/sirupsen/logrus"
)

var (
	logger   = alog.LambdaLogger(logrus.InfoLevel)
	clientIDFunc = func() string { return utils.NonEmptyEnv("CLIENT_ID") }
)

func HandleRequest(ctx context.Context, e1 model.StreamEntity) {
	logger = logger.WithLambdaContext(ctx)
	logger.WithField("event", &e1).Info("Incoming request")

	conn1, err := sqc.NewMySqlConnection()
	defer func() {
		_ = conn1.Close()
	}()

	if err == nil {
		e2 := streamhandler.StreamEntityHandler(e1, clientIDFunc(), conn1, logger)
		logger.WithField("mapped_event", &e2).Info("Transformed request")
	} else {
		logger.WithField("error", &err).Errorf("Could not establish a database connection")
	}
}
