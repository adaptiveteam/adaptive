package workflows

import (
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/nlopes/slack"
)

// TriggerAllPostponedEvents reads all postponed events and run them sequentially
// Needs connection
// conn := common.DynamoDBConnection{
// 	Dynamo: d,
// 	ClientID: clientID,
// 	PlatformID: teamID,
// }
// returns the number of posted events
func TriggerAllPostponedEvents(teamID models.TeamID, userID string)func (conn common.DynamoDBConnection)(count int, err error) {
	return wf.ForeachActionPathForUserID(userID, func (ap models.ActionPath, conn common.DynamoDBConnection) (err error) {
		callbackID := ap.Encode()
		logger.WithField("userID", userID).WithField("callbackID", callbackID).
			Infof("invokeWorkflow")
		np := models.NamespacePayload4{
			Namespace: "triggerAllPostponedEvents",
			PlatformRequest: models.PlatformRequest{
				TeamID: teamID,
				SlackRequest: models.SlackRequest{
					Type: models.InteractionSlackRequestType,
					InteractionCallback: slack.InteractionCallback{
						User: slack.User{
							ID: userID,
						},
						CallbackID: callbackID,
					},
				},
			},
		}
		err = InvokeWorkflow(np, conn)
		return
	})
}
