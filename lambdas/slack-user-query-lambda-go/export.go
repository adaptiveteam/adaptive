package lambda

import (
	"encoding/json"
	"fmt"

	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// RefreshCommunitiesForPlatform calls this lambda from outside
func RefreshCommunitiesForPlatform(l awsutils.LambdaRequest, slackLambdaName string, teamID models.TeamID) {
	// Invoke a lambda that handles slack users
	clientReq := models.ClientPlatformRequest{TeamID: teamID}
	reqBytes, err2 := json.Marshal(clientReq)
	core.ErrorHandler(err2, namespace, "Could not marshal to ClientPlatformRequest")
	_, err3 := l.InvokeFunction(slackLambdaName, reqBytes, true)
	core.ErrorHandler(err3, namespace, fmt.Sprintf("Could not invoke %s", slackLambdaName))
}
