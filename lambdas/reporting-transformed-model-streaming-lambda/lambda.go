package reporting_transformed_model_streaming_lambda

import (
	"context"

	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	sqlconnector "github.com/adaptiveteam/adaptive/adaptive-utils-go/sql-connector"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/mapping"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	streamhandler "github.com/adaptiveteam/adaptive/lambdas/stream-handler"
	"github.com/sirupsen/logrus"
)

var (
	logger       = alog.LambdaLogger(logrus.InfoLevel)
	clientIDFunc = func() string { return utils.NonEmptyEnv("CLIENT_ID") }
)

func HandleRequest(ctx context.Context, e1 model.StreamEntity) {
	logger = logger.WithLambdaContext(ctx)
	logger.WithField("event", &e1).Info("Incoming request")

	conn1, err2 := sqlconnector.ReadRDSConfigFromEnv().GormOpen()

	if err2 == nil {
		defer func() {
			_ = conn1.Close()
		}()
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
