package lambda

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
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
	}
}
