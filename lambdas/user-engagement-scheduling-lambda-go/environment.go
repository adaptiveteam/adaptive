package lambda

import (
	"encoding/json"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

type Config struct {
	namespace               string
	engScriptingLambdaArn   string
	engSchedulerLambdaName  string
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
		engSchedulerLambdaName:  utils.NonEmptyEnv("USER_ENGAGEMENT_SCHEDULER_LAMBDA_NAME"),
		usersTable:              utils.NonEmptyEnv("USERS_TABLE_NAME"),
		communityUsersTable:     utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME"),
		usersScheduledTimeIndex: "PlatformIDAdaptiveScheduledTimeInUTCIndex",
		usersZoneOffsetIndex:    "PlatformIDTimezoneOffsetIndex",
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

func invokeSchedulerLambda(engage models.UserEngage, config Config) (err error) {
	payloadJSONBytes, _ := json.Marshal(engage)
	_, err = config.l.InvokeFunction(config.engSchedulerLambdaName, payloadJSONBytes, true)
	return
}
