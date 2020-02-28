package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

var (
	namespace         = utils.NonEmptyEnv("LOG_NAMESPACE")
	clientConfigTable = utils.NonEmptyEnv("CLIENT_CONFIG_TABLE_NAME")
	slackLambdaName   = utils.NonEmptyEnv("SLACK_LAMBDA_FUNCTION_NAME")
	d                 = awsutils.NewDynamo(utils.NonEmptyEnv("AWS_REGION"), "", namespace)
	l                 = awsutils.NewLambda(utils.NonEmptyEnv("AWS_REGION"), "", namespace)
)

func HandleRequest(ctx context.Context) (clientCount int, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("Error in user-query-lambda %v", err2)
		}
	}()
	var scanRes []models.ClientPlatformToken
	err = d.ScanTable(clientConfigTable, &scanRes)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not scan %s table", clientConfigTable))

	clientCount = len(scanRes)
	for _, each := range scanRes {
		if each.PlatformName == models.SlackPlatform {
			// Invoke a lambda that handles slack users
			clientReq := models.ClientPlatformRequest{TeamID: models.ParseTeamID(each.PlatformID), Org: each.Org}
			reqBytes, err := json.Marshal(clientReq)
			core.ErrorHandler(err, namespace, "Could not marshal to ClientPlatformRequest")
			_, err = l.InvokeFunction(slackLambdaName, reqBytes, true)
			core.ErrorHandler(err, namespace, fmt.Sprintf("Could not invoke %s", slackLambdaName))
		} else if each.PlatformName == models.MsTeamsPlatform {
			// Invoke a lambda that handles ms-teams users
		}
	}
	return
}
