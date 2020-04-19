package entityBootstrapLambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	streamhandler "github.com/adaptiveteam/adaptive/lambdas/stream-handler"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	aws_utils_go "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	logger            = alog.LambdaLogger(logrus.InfoLevel)
	streamEventMapper = func () string { return utils.NonEmptyEnv("STREAM_EVENT_MAPPER_LAMBDA") }
	clientID          = utils.NonEmptyEnv("CLIENT_ID")
)

func tableWithSuffix(suffix string, list []*string) *string {
	for _, each := range list {
		if strings.Contains(*each, fmt.Sprintf("%s_%s", clientID, suffix)) {
			return each
		}
	}
	return nil
}

func addTableRecordsToDB(tableName string, d *aws_utils_go.DynamoRequest) {
	var op []interface{}
	err := d.ScanTable(tableName, &op)
	if err == nil {
		logger.Infof("Number of records in scan for %s table: %d", tableName, len(op))
		for _, each := range op {
			entity := model.StreamEntity{
				TableName: tableName + "/",
				NewEntity: each,
				// StreamEventEdit updates record if exists and creates if it doesn't
				EventType: model.StreamEventAdd,
			}
			logger.WithField("entity", &entity).Info("Invoke payload")
			byt, _ := json.Marshal(entity)
			io, err := streamhandler.LambdaClient.InvokeFunction(streamEventMapper(), byt, false)
			if err == nil {
				logger.Infof("GoString:"+io.GoString())
			} else {
				logger.WithField("error", err).Errorf("Could not invoke stream mapper lambda")
			} 
		}
	} else {
		logger.WithField("error", err).Errorf("Could not scan %s table", tableName)
	}
}

func HandleRequest(ctx context.Context) {
	logger = logger.WithLambdaContext(ctx)

	tableSuffixes := streamhandler.TableRefKeys()

	d := streamhandler.DynamoClient
	tables, err := d.ListTables(nil)
	if err == nil {
		logger.Infof("Length of table suffixes: %d", len(tableSuffixes))
		for _, tableSuffix := range tableSuffixes {
			table := tableWithSuffix(tableSuffix, tables.TableNames)
			if table != nil {
				logger.Infof("%s suffix exists in the list with table: %s", tableSuffix, *table)
				addTableRecordsToDB(*table, d)
			} else {
				logger.Warnf("No table exists with %s suffix", tableSuffix)
			}
		}
	} else {
		logger.WithField("error", err).Error("Error with retrieving tables")
	}
}
