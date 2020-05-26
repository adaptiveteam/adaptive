package strategy

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

var (
	namespace = func() string { return utils.NonEmptyEnv("LOG_NAMESPACE") }
	region    = func() string { return utils.NonEmptyEnv("AWS_REGION") }
	d         = func() *awsutils.DynamoRequest { return awsutils.NewDynamo(region(), "", namespace()) }
	clientID  = func() string { return utils.NonEmptyEnv("CLIENT_ID") }
	schema    = func() models.Schema { return models.SchemaForClientID(clientID()) }
)
