package main

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

type Config struct {
	namespace  string
	region     string
	clientID   string
	platformID models.PlatformID

	d         *awsutils.DynamoRequest
}

func readConfigFromEnvVars() (config Config) {
	region    := utils.NonEmptyEnv("AWS_REGION")
	namespace := utils.NonEmptyEnv("LOG_NAMESPACE")
	return Config{
		namespace : namespace,
		region    : region,
		clientID  : utils.NonEmptyEnv("CLIENT_ID"),
		platformID: models.PlatformID(utils.NonEmptyEnv("PLATFORM_ID")),
		d         : awsutils.NewDynamo(region, "", namespace),
	}
}

var userFeedbackTableName = func (clientID string) string { return clientID + "_adaptive_user_feedback" }
const userFeedbackSourceQYIndex = "SourceQuarterYear"
