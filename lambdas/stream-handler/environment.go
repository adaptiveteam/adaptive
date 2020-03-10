package streamhandler

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

var (
	region       = utils.NonEmptyEnv("AWS_REGION")
	LambdaClient = awsutils.NewLambda(region, "", "")
	DynamoClient = awsutils.NewDynamo(region, "", "")
)
