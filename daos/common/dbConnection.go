package common

import (
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)
// DynamoDBConnection has just what is needed for connection to Dynamo
// ClientID allows to get table names.
// PlatformID is the sharding key that is required in all queries.
type DynamoDBConnection struct {
	Dynamo     *awsutils.DynamoRequest
	ClientID   string
	PlatformID PlatformID
}
