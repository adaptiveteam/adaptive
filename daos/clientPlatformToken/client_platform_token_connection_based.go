package clientPlatformToken
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/aws/aws-sdk-go/aws"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
)


// Create saves the ClientPlatformToken.
func Create(clientPlatformToken ClientPlatformToken) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := clientPlatformToken.CollectEmptyFields()
		if ok {
			
			err = conn.Dynamo.PutTableEntry(clientPlatformToken, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the ClientPlatformToken.
func CreateUnsafe(clientPlatformToken ClientPlatformToken) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(clientPlatformToken)(conn)
		core.ErrorHandler(err2, "daos/ClientPlatformToken", fmt.Sprintf("Could not create platformID==%s in %s\n", clientPlatformToken.PlatformID, TableName(conn.ClientID)))
	}
}


// Read reads ClientPlatformToken
func Read(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out ClientPlatformToken, err error) {
	return func (conn common.DynamoDBConnection) (out ClientPlatformToken, err error) {
		var outs []ClientPlatformToken
		outs, err = ReadOrEmpty(platformID)(conn)
		if err == nil && len(outs) == 0 {
			err = fmt.Errorf("Not found platformID==%s in %s\n", platformID, TableName(conn.ClientID))
		}
		if len(outs) > 0 {
			out = outs[0]
		}
		return
	}
}


// ReadUnsafe reads the ClientPlatformToken. Panics in case of any errors
func ReadUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) ClientPlatformToken {
	return func (conn common.DynamoDBConnection) ClientPlatformToken {
		out, err2 := Read(platformID)(conn)
		core.ErrorHandler(err2, "daos/ClientPlatformToken", fmt.Sprintf("Error reading platformID==%s in %s\n", platformID, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads ClientPlatformToken
func ReadOrEmpty(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []ClientPlatformToken, err error) {
	return func (conn common.DynamoDBConnection) (out []ClientPlatformToken, err error) {
       out, err = ReadOrEmptyIncludingInactive(platformID)(conn)
       
       
		return
	}
}


// ReadOrEmptyUnsafe reads the ClientPlatformToken. Panics in case of any errors
func ReadOrEmptyUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) []ClientPlatformToken {
	return func (conn common.DynamoDBConnection) []ClientPlatformToken {
		out, err2 := ReadOrEmpty(platformID)(conn)
		core.ErrorHandler(err2, "daos/ClientPlatformToken", fmt.Sprintf("Error while reading platformID==%s in %s\n", platformID, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads ClientPlatformToken
func ReadOrEmptyIncludingInactive(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []ClientPlatformToken, err error) {
	return func (conn common.DynamoDBConnection) (out []ClientPlatformToken, err error) {
		var outOrEmpty ClientPlatformToken
		ids := idParams(platformID)
		var found bool
		found, err = conn.Dynamo.GetItemOrEmptyFromTable(TableName(conn.ClientID), ids, &outOrEmpty)
		if found {
			if outOrEmpty.PlatformID == platformID {
				out = append(out, outOrEmpty)
			} else {
				err = fmt.Errorf("Requested ids: platformID==%s are different from the found ones: platformID==%s", platformID, outOrEmpty.PlatformID) // unexpected error: found ids != ids
			}
		}
		err = errors.Wrapf(err, "ClientPlatformToken DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the ClientPlatformToken. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(platformID common.PlatformID) func (conn common.DynamoDBConnection) []ClientPlatformToken {
	return func (conn common.DynamoDBConnection) []ClientPlatformToken {
		out, err2 := ReadOrEmptyIncludingInactive(platformID)(conn)
		core.ErrorHandler(err2, "daos/ClientPlatformToken", fmt.Sprintf("Error while reading platformID==%s in %s\n", platformID, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the ClientPlatformToken regardless of if it exists.
func CreateOrUpdate(clientPlatformToken ClientPlatformToken) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		
		var olds []ClientPlatformToken
		olds, err = ReadOrEmpty(clientPlatformToken.PlatformID)(conn)
		err = errors.Wrapf(err, "ClientPlatformToken DAO.CreateOrUpdate(id = platformID==%s) couldn't ReadOrEmpty", clientPlatformToken.PlatformID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(clientPlatformToken)(conn)
				err = errors.Wrapf(err, "ClientPlatformToken DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := clientPlatformToken.CollectEmptyFields()
				if ok {
					old := olds[0]
					
					
					key := idParams(old.PlatformID)
					expr, exprAttributes, names := updateExpression(clientPlatformToken, old)
					input := dynamodb.UpdateItemInput{
						ExpressionAttributeValues: exprAttributes,
						TableName:                 aws.String(TableName(conn.ClientID)),
						Key:                       key,
						ReturnValues:              aws.String("UPDATED_NEW"),
						UpdateExpression:          aws.String(expr),
					}
					if names != nil { input.ExpressionAttributeNames = *names } // workaround for a pointer to an empty slice
					if  len(exprAttributes) > 0 { // if there some changes
						err = conn.Dynamo.UpdateItemInternal(input)
					} else {
						// WARN: no changes.
					}
					err = errors.Wrapf(err, "ClientPlatformToken DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the ClientPlatformToken regardless of if it exists.
func CreateOrUpdateUnsafe(clientPlatformToken ClientPlatformToken) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(clientPlatformToken)(conn)
		core.ErrorHandler(err2, "daos/ClientPlatformToken", fmt.Sprintf("could not create or update %v in %s\n", clientPlatformToken, TableName(conn.ClientID)))
	}
}


// Delete removes ClientPlatformToken from db
func Delete(platformID common.PlatformID) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams(platformID))
	}
}


// DeleteUnsafe deletes ClientPlatformToken and panics in case of errors.
func DeleteUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Delete(platformID)(conn)
		core.ErrorHandler(err2, "daos/ClientPlatformToken", fmt.Sprintf("Could not delete platformID==%s in %s\n", platformID, TableName(conn.ClientID)))
	}
}

