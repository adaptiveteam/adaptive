package common

import (
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

// CreateConnectionGenFromEnv constructs connection generator reading environment variables
func CreateConnectionGenFromEnv() (connGen DynamoDBConnectionGen) {
	namespace := core.NonEmptyEnv("LOG_NAMESPACE")
	region    := core.NonEmptyEnv("AWS_REGION")
	d         := awsutils.NewDynamo(region, "", namespace)
	clientID  := core.NonEmptyEnv("CLIENT_ID")

	connGen = DynamoDBConnectionGen{
		Dynamo: d,
		TableNamePrefix: clientID,
	}
	return
}

// CreateConnectionFromEnv reads env vars and creates connection for the given platform id
func CreateConnectionFromEnv(platformID PlatformID) (conn DynamoDBConnection) {
	return CreateConnectionGenFromEnv().ForPlatformID(platformID)
}
