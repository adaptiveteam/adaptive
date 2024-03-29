package lambda

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

var (
	namespace = utils.NonEmptyEnv("LOG_NAMESPACE")
	region    = utils.NonEmptyEnv("AWS_REGION")
	d         = awsutils.NewDynamo(region, "", namespace)
	userTable = utils.NonEmptyEnv("USER_TABLE_NAME")
	confTable = utils.NonEmptyEnv("CLIENT_CONFIG_TABLE_NAME")

	clientID  = utils.NonEmptyEnv("CLIENT_ID")

	schema           = models.SchemaForClientID(clientID)

	connGen          = daosCommon.DynamoDBConnectionGen{
		Dynamo:          d,
		TableNamePrefix: clientID,
	}
)
