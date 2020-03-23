package common

import (
	"github.com/adaptiveteam/adaptive/daos/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	
)

var globalDns *DynamoNamespace
var globalS3  *awsutils.S3Request
var connGen   *common.DynamoDBConnectionGen

func initGlobals() {
	Namespace := utils.NonEmptyEnv("LOG_NAMESPACE")
	Region    := utils.NonEmptyEnv("AWS_REGION")
	d         := awsutils.NewDynamo(Region, "", Namespace)
	globalDns = &DynamoNamespace{Dynamo: d, Namespace: Namespace}
	globalS3  = awsutils.NewS3(Region, "", Namespace)
	clientID  := utils.NonEmptyEnv("CLIENT_ID")

	connGen   = &common.DynamoDBConnectionGen{Dynamo: d, TableNamePrefix: clientID}
}
// DeprecatedGetGlobalDns reads environment variables and creates a commection to Dynamo
// Deprecated: Shouldn't be used
func DeprecatedGetGlobalDns() DynamoNamespace {
	if globalDns == nil {
		initGlobals()
	}
	return  *globalDns
}

// DeprecatedGetGlobalS3 reads environment variables and creates a commection to S3
// Deprecated: Shouldn't be used
func DeprecatedGetGlobalS3() *awsutils.S3Request {
	if globalS3 == nil {
		initGlobals()
	}
	return  globalS3
}

// DeprecatedGetGlobalConnectionGen reads environment variables and creates a commection to DynamoDB
// Deprecated: Shouldn't be used. Pass connection from lambdas instead
func DeprecatedGetGlobalConnectionGen() common.DynamoDBConnectionGen {
	if connGen == nil {
		initGlobals()
	}
	return *connGen
}
