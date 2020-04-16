package reporting_transformed_model_streaming_lambda

import (
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/mapping"
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

	conn1, err2 := sqc.NewMySqlConnection()
	defer func() {
		_ = conn1.Close()
	}()

	if err2 == nil {
		if e1.TableName == "" {
			logger.Info("AutoMigrate-ing AllEntities all entities")
			mapping.AutoMigrateAllEntities(conn1)
		} else {
			e2 := streamhandler.StreamEntityHandler(e1, clientIDFunc(), conn1, logger)
			logger.WithField("mapped_event", &e2).Info("Transformed request")
		}
	} else {
		logger.WithField("error", &err2).Errorf("Could not establish a database connection")
	}
}
