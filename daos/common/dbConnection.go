package common

import (
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)
// DynamoDBConnectionGen - a connection to the database with
// the common table prefix. 
type DynamoDBConnectionGen struct {
	Dynamo          *awsutils.DynamoRequest
	TableNamePrefix string
}
// DynamoDBConnection has just what is needed for connection to Dynamo
// ClientID allows to get table names.
// PlatformID is the sharding key that is required in all queries.
type DynamoDBConnection struct {
	Dynamo     *awsutils.DynamoRequest
	ClientID   string
	PlatformID PlatformID
}

func (dgen DynamoDBConnectionGen)ForPlatformID(platformID PlatformID) DynamoDBConnection {
	return DynamoDBConnection{
		Dynamo: dgen.Dynamo,
		ClientID: dgen.TableNamePrefix,
		PlatformID: platformID,
	}
}
