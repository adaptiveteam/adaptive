package postponedEvent
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
)


// Create saves the PostponedEvent.
func Create(postponedEvent PostponedEvent) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := postponedEvent.CollectEmptyFields()
		if ok {
			postponedEvent.ModifiedAt = core.CurrentRFCTimestamp()
	postponedEvent.CreatedAt = postponedEvent.ModifiedAt
	
			err = conn.Dynamo.PutTableEntry(postponedEvent, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the PostponedEvent.
func CreateUnsafe(postponedEvent PostponedEvent) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(postponedEvent)(conn)
		core.ErrorHandler(err2, "daos/PostponedEvent", fmt.Sprintf("Could not create id==%s in %s\n", postponedEvent.ID, TableName(conn.ClientID)))
	}
}


// Read reads PostponedEvent
func Read(id string) func (conn common.DynamoDBConnection) (out PostponedEvent, err error) {
	return func (conn common.DynamoDBConnection) (out PostponedEvent, err error) {
		var outs []PostponedEvent
		outs, err = ReadOrEmpty(id)(conn)
		if err == nil && len(outs) == 0 {
			err = fmt.Errorf("Not found id==%s in %s\n", id, TableName(conn.ClientID))
		}
		if len(outs) > 0 {
			out = outs[0]
		}
		return
	}
}


// ReadUnsafe reads the PostponedEvent. Panics in case of any errors
func ReadUnsafe(id string) func (conn common.DynamoDBConnection) PostponedEvent {
	return func (conn common.DynamoDBConnection) PostponedEvent {
		out, err2 := Read(id)(conn)
		core.ErrorHandler(err2, "daos/PostponedEvent", fmt.Sprintf("Error reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads PostponedEvent
func ReadOrEmpty(id string) func (conn common.DynamoDBConnection) (out []PostponedEvent, err error) {
	return func (conn common.DynamoDBConnection) (out []PostponedEvent, err error) {
       out, err = ReadOrEmptyIncludingInactive(id)(conn)
       
       
		return
	}
}


// ReadOrEmptyUnsafe reads the PostponedEvent. Panics in case of any errors
func ReadOrEmptyUnsafe(id string) func (conn common.DynamoDBConnection) []PostponedEvent {
	return func (conn common.DynamoDBConnection) []PostponedEvent {
		out, err2 := ReadOrEmpty(id)(conn)
		core.ErrorHandler(err2, "daos/PostponedEvent", fmt.Sprintf("Error while reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads PostponedEvent
func ReadOrEmptyIncludingInactive(id string) func (conn common.DynamoDBConnection) (out []PostponedEvent, err error) {
	return func (conn common.DynamoDBConnection) (out []PostponedEvent, err error) {
		var outOrEmpty PostponedEvent
		ids := idParams(id)
		var found bool
		found, err = conn.Dynamo.GetItemOrEmptyFromTable(TableName(conn.ClientID), ids, &outOrEmpty)
		if found {
			if outOrEmpty.ID == id {
				out = append(out, outOrEmpty)
			} else {
				err = fmt.Errorf("Requested ids: id==%s are different from the found ones: id==%s", id, outOrEmpty.ID) // unexpected error: found ids != ids
			}
		}
		err = errors.Wrapf(err, "PostponedEvent DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the PostponedEvent. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(id string) func (conn common.DynamoDBConnection) []PostponedEvent {
	return func (conn common.DynamoDBConnection) []PostponedEvent {
		out, err2 := ReadOrEmptyIncludingInactive(id)(conn)
		core.ErrorHandler(err2, "daos/PostponedEvent", fmt.Sprintf("Error while reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the PostponedEvent regardless of if it exists.
func CreateOrUpdate(postponedEvent PostponedEvent) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		postponedEvent.ModifiedAt = core.CurrentRFCTimestamp()
	if postponedEvent.CreatedAt == "" { postponedEvent.CreatedAt = postponedEvent.ModifiedAt }
	
		var olds []PostponedEvent
		olds, err = ReadOrEmpty(postponedEvent.ID)(conn)
		err = errors.Wrapf(err, "PostponedEvent DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", postponedEvent.ID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(postponedEvent)(conn)
				err = errors.Wrapf(err, "PostponedEvent DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := postponedEvent.CollectEmptyFields()
				if ok {
					old := olds[0]
					postponedEvent.CreatedAt  = old.CreatedAt
					postponedEvent.ModifiedAt = core.CurrentRFCTimestamp()
					key := idParams(old.ID)
					expr, exprAttributes, names := updateExpression(postponedEvent, old)
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
					err = errors.Wrapf(err, "PostponedEvent DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the PostponedEvent regardless of if it exists.
func CreateOrUpdateUnsafe(postponedEvent PostponedEvent) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(postponedEvent)(conn)
		core.ErrorHandler(err2, "daos/PostponedEvent", fmt.Sprintf("could not create or update %v in %s\n", postponedEvent, TableName(conn.ClientID)))
	}
}


// Delete removes PostponedEvent from db
func Delete(id string) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams(id))
	}
}


// DeleteUnsafe deletes PostponedEvent and panics in case of errors.
func DeleteUnsafe(id string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Delete(id)(conn)
		core.ErrorHandler(err2, "daos/PostponedEvent", fmt.Sprintf("Could not delete id==%s in %s\n", id, TableName(conn.ClientID)))
	}
}


func ReadByPlatformIDUserID(platformID common.PlatformID, userID string) func (conn common.DynamoDBConnection) (out []PostponedEvent, err error) {
	return func (conn common.DynamoDBConnection) (out []PostponedEvent, err error) {
		var instances []PostponedEvent
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "PlatformIDUserIDIndex",
			Condition: "platform_id = :a0 and user_id = :a1",
			Attributes: map[string]interface{}{
				":a0": platformID,
			":a1": userID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByPlatformIDUserIDUnsafe(platformID common.PlatformID, userID string) func (conn common.DynamoDBConnection) (out []PostponedEvent) {
	return func (conn common.DynamoDBConnection) (out []PostponedEvent) {
		out, err2 := ReadByPlatformIDUserID(platformID, userID)(conn)
		core.ErrorHandler(err2, "daos/PostponedEvent", fmt.Sprintf("Could not query PlatformIDUserIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByUserID(userID string) func (conn common.DynamoDBConnection) (out []PostponedEvent, err error) {
	return func (conn common.DynamoDBConnection) (out []PostponedEvent, err error) {
		var instances []PostponedEvent
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "UserIDIndex",
			Condition: "user_id = :a0",
			Attributes: map[string]interface{}{
				":a0": userID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByUserIDUnsafe(userID string) func (conn common.DynamoDBConnection) (out []PostponedEvent) {
	return func (conn common.DynamoDBConnection) (out []PostponedEvent) {
		out, err2 := ReadByUserID(userID)(conn)
		core.ErrorHandler(err2, "daos/PostponedEvent", fmt.Sprintf("Could not query UserIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}

