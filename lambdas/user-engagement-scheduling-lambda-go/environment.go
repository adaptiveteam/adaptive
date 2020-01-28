package lambda

import (
	"encoding/json"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	// workflows "github.com/adaptiveteam/adaptive/workflows"
	// "github.com/adaptiveteam/adaptive/daos/common"
)

type Config struct {
	namespace               string
	engScriptingLambdaArn   string
	usersTable              string
	communityUsersTable     string
	usersScheduledTimeIndex string
	usersZoneOffsetIndex    string
	communityUsersUserIndex string
	clientConfigTable       string
	region                  string
	cw                      *awsutils.CloudWatchRequest
	l                       *awsutils.LambdaRequest
	d                       *awsutils.DynamoRequest
	clientID                string
}

func readConfigFromEnvironment() Config {
	region := utils.NonEmptyEnv("AWS_REGION")
	namespace := utils.NonEmptyEnv("LOG_NAMESPACE")
	return Config{
		namespace:               namespace,
		engScriptingLambdaArn:   utils.NonEmptyEnv("USER_ENGAGEMENT_SCRIPTING_LAMBDA_ARN"),
		usersTable:              utils.NonEmptyEnv("USERS_TABLE_NAME"),
		communityUsersTable:     utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME"),
		usersScheduledTimeIndex: utils.NonEmptyEnv("USERS_SCHEDULED_TIME_INDEX"),
		usersZoneOffsetIndex:    utils.NonEmptyEnv("USERS_TIMEZONE_OFFSET_INDEX"),
		communityUsersUserIndex: utils.NonEmptyEnv("COMMUNITY_USERS_USER_INDEX"),
		clientConfigTable:       utils.NonEmptyEnv("CLIENT_CONFIG_TABLE_NAME"),
		region:                  region,
		cw:                      awsutils.NewCloudWatch(region, "", namespace),
		l:                       awsutils.NewLambda(region, "", namespace),
		d:                       awsutils.NewDynamo(region, "", namespace),
		clientID: 				 utils.NonEmptyEnv("CLIENT_ID"),
	}
}

func invokeScriptingLambda(engage models.UserEngage, config Config) (err error) {
	payloadJSONBytes, _ := json.Marshal(engage)
	_, err = config.l.InvokeFunction(config.engScriptingLambdaArn, payloadJSONBytes, true)
	return
}

// func triggerPostponedEvents(engage models.UserEngage, config Config) (err error) {
// 	logger.WithField("userID", engage.UserId).Infof("triggerPostponedEvents")
// 	conn := common.DynamoDBConnection{
// 		Dynamo: config.d,
// 		ClientID: config.clientID,
// 		PlatformID: engage.PlatformID,
// 	}
// 	err = workflows.TriggerAllPostponedEvents(engage.PlatformID, engage.UserId)(conn)
// 	return
// }
