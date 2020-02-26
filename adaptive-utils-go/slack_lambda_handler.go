package adaptive_utils_go

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"fmt"
	ls "github.com/aws/aws-lambda-go/lambda"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/nlopes/slack"
)

// LambdaHandler represents a structured handler for Slack events.
type LambdaHandler struct {
	Namespace string
	DispatchSlackInteractionCallback func (slack.InteractionCallback)
	DispatchSlackDialogSubmissionCallback func (slack.InteractionCallback, slack.DialogSubmissionCallback)
}

// StartHandler starts server
func (l LambdaHandler)StartHandler() {
	ls.Start(l.HandleRequest)
}

// HandleRequest receives lambda json event
func (l LambdaHandler)HandleRequest(ctx context.Context, e events.SNSEvent) error {
	fmt.Println("adaptiveValues/main.go/HandleRequest entered")
	for _, record := range e.Records {
		fmt.Println(record)
		if record.SNS.Message == "warmup" {
			// Ignoring warmup messages
		} else {
			np := models.UnmarshalNamespacePayload4JSONUnsafe(record.SNS.Message)
			if np.Namespace == l.Namespace {
				switch np.PlatformRequest.Type {
				case models.InteractionSlackRequestType:
					l.DispatchSlackInteractionCallback(np.PlatformRequest.InteractionCallback)
				case models.DialogSubmissionSlackRequestType:
					l.DispatchSlackDialogSubmissionCallback(np.PlatformRequest.InteractionCallback, 
						np.PlatformRequest.DialogSubmissionCallback)
				}
			}
		}
	}
	return nil // we do not have handlable errors. Only panics
}

// HandleNamespacePayload4 receives lambda json event
func (l LambdaHandler)HandleNamespacePayload4(np models.NamespacePayload4) error {
	fmt.Println("adaptiveValues/main.go/HandleRequest entered")
	if np.Namespace == l.Namespace {
		switch np.PlatformRequest.Type {
		case models.InteractionSlackRequestType:
			l.DispatchSlackInteractionCallback(np.PlatformRequest.InteractionCallback)
		case models.DialogSubmissionSlackRequestType:
			l.DispatchSlackDialogSubmissionCallback(np.PlatformRequest.InteractionCallback, 
				np.PlatformRequest.DialogSubmissionCallback)
		}
	}
	return nil // we do not have handlable errors. Only panics
}
