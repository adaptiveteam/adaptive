package workflows

import (
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/nlopes/slack"

)

// TriggerAllPostponedEvents reads all postponed events and run them sequentially
func TriggerAllPostponedEvents(platformID models.PlatformID, userID string)func (conn common.DynamoDBConnection)(err error) {
	// conn := common.DynamoDBConnection{
	// 	Dynamo: d,
	// 	ClientID: clientID,
	// 	PlatformID: platformID,
	// }
	return wf.ForeachActionPathForUserID(userID, func (ap models.ActionPath) (err error) {
		np := models.NamespacePayload4{
			Namespace: "triggerAllPostponedEvents",
			PlatformRequest: models.PlatformRequest{
				PlatformID: platformID,
				SlackRequest: models.SlackRequest{
					Type: models.InteractionSlackRequestType,
					InteractionCallback: slack.InteractionCallback{
						User: slack.User{
							ID: userID,
						},
						CallbackID: ap.Encode(),
					},
				},
			},
		}
		err = invokeWorkflow(np)
		return
	})
}

func invokeWorkflow(np models.NamespacePayload4) (err error) {
	return
}
