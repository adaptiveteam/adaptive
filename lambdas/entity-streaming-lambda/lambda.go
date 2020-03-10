package lambda

import (
	"context"
	"encoding/json"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	streamhandler "github.com/adaptiveteam/adaptive/lambdas/stream-handler"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
)

var (
	logger            = alog.LambdaLogger(logrus.InfoLevel)
	streamEventMapper = utils.NonEmptyEnv("STREAM_EVENT_MAPPER_LAMBDA")
)

func marshalStreamImageToInterfaceUnsafe(change map[string]events.DynamoDBAttributeValue) (u interface{}) {
	err := awsutils.UnmarshalStreamImage(change, &u)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal %v to interface", change)
	}
	return
}

func HandleRequest(ctx context.Context, e events.DynamoDBEvent) {
	for _, record := range e.Records {
		eventSourceARN := record.EventSourceArn
		recordChange := record.Change
		recordEventName := record.EventName
		recordEventID := record.EventID

		logger.Infof("Processing request data for event ID %s, type %s, source ARN %s",
			recordEventID, recordEventName, eventSourceARN)

		var event model.StreamEntity

		switch recordEventName {
		case string(events.DynamoDBOperationTypeInsert):
			newIface := marshalStreamImageToInterfaceUnsafe(recordChange.NewImage)
			event = model.StreamEntity{
				TableName: eventSourceARN,
				NewEntity: newIface,
				EventType: model.StreamEventAdd,
			}
		case string(events.DynamoDBOperationTypeModify):
			oldIface := marshalStreamImageToInterfaceUnsafe(recordChange.OldImage)
			newIface := marshalStreamImageToInterfaceUnsafe(recordChange.NewImage)
			event = model.StreamEntity{
				TableName: eventSourceARN,
				OldEntity: oldIface,
				NewEntity: newIface,
				EventType: model.StreamEventEdit,
			}
		case string(events.DynamoDBOperationTypeRemove):
			oldIface := marshalStreamImageToInterfaceUnsafe(recordChange.OldImage)
			event = model.StreamEntity{
				TableName: eventSourceARN,
				OldEntity: oldIface,
				EventType: model.StreamEventDelete,
			}
		default:
			logger.Warnf("Event %s did not match any case", recordEventName)
		}

		if event.EventType != "" {
			byt, _ := json.Marshal(event)
			_, err := streamhandler.LambdaClient.InvokeFunction(streamEventMapper, byt, false)
			if err != nil {
				logger.WithField("error", err).Errorf("Could not invoke stream mapper lambda for %s event",
					recordEventName)
			}
		} else {
			logger.Errorf("Could not match %s event to any of the handlers", recordEventName)
		}
	}
}
