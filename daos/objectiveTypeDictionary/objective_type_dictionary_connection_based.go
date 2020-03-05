package objectiveTypeDictionary
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


// Create saves the ObjectiveTypeDictionary.
func Create(objectiveTypeDictionary ObjectiveTypeDictionary) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := objectiveTypeDictionary.CollectEmptyFields()
		if ok {
			objectiveTypeDictionary.ModifiedAt = core.CurrentRFCTimestamp()
	objectiveTypeDictionary.CreatedAt = objectiveTypeDictionary.ModifiedAt
	
			err = conn.Dynamo.PutTableEntry(objectiveTypeDictionary, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the ObjectiveTypeDictionary.
func CreateUnsafe(objectiveTypeDictionary ObjectiveTypeDictionary) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(objectiveTypeDictionary)(conn)
		core.ErrorHandler(err2, "daos/ObjectiveTypeDictionary", fmt.Sprintf("Could not create id==%s in %s\n", objectiveTypeDictionary.ID, TableName(conn.ClientID)))
	}
}


// Read reads ObjectiveTypeDictionary
func Read(id string) func (conn common.DynamoDBConnection) (out ObjectiveTypeDictionary, err error) {
	return func (conn common.DynamoDBConnection) (out ObjectiveTypeDictionary, err error) {
		var outs []ObjectiveTypeDictionary
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


// ReadUnsafe reads the ObjectiveTypeDictionary. Panics in case of any errors
func ReadUnsafe(id string) func (conn common.DynamoDBConnection) ObjectiveTypeDictionary {
	return func (conn common.DynamoDBConnection) ObjectiveTypeDictionary {
		out, err2 := Read(id)(conn)
		core.ErrorHandler(err2, "daos/ObjectiveTypeDictionary", fmt.Sprintf("Error reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads ObjectiveTypeDictionary
func ReadOrEmpty(id string) func (conn common.DynamoDBConnection) (out []ObjectiveTypeDictionary, err error) {
	return func (conn common.DynamoDBConnection) (out []ObjectiveTypeDictionary, err error) {
		var outOrEmpty ObjectiveTypeDictionary
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
		err = errors.Wrapf(err, "ObjectiveTypeDictionary DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyUnsafe reads the ObjectiveTypeDictionary. Panics in case of any errors
func ReadOrEmptyUnsafe(id string) func (conn common.DynamoDBConnection) []ObjectiveTypeDictionary {
	return func (conn common.DynamoDBConnection) []ObjectiveTypeDictionary {
		out, err2 := ReadOrEmpty(id)(conn)
		core.ErrorHandler(err2, "daos/ObjectiveTypeDictionary", fmt.Sprintf("Error while reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the ObjectiveTypeDictionary regardless of if it exists.
func CreateOrUpdate(objectiveTypeDictionary ObjectiveTypeDictionary) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		objectiveTypeDictionary.ModifiedAt = core.CurrentRFCTimestamp()
	if objectiveTypeDictionary.CreatedAt == "" { objectiveTypeDictionary.CreatedAt = objectiveTypeDictionary.ModifiedAt }
	
		var olds []ObjectiveTypeDictionary
		olds, err = ReadOrEmpty(objectiveTypeDictionary.ID)(conn)
		err = errors.Wrapf(err, "ObjectiveTypeDictionary DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", objectiveTypeDictionary.ID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(objectiveTypeDictionary)(conn)
				err = errors.Wrapf(err, "ObjectiveTypeDictionary DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := objectiveTypeDictionary.CollectEmptyFields()
				if ok {
					old := olds[0]
					objectiveTypeDictionary.CreatedAt  = old.CreatedAt
					objectiveTypeDictionary.ModifiedAt = core.CurrentRFCTimestamp()
					key := idParams(old.ID)
					expr, exprAttributes, names := updateExpression(objectiveTypeDictionary, old)
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
					err = errors.Wrapf(err, "ObjectiveTypeDictionary DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the ObjectiveTypeDictionary regardless of if it exists.
func CreateOrUpdateUnsafe(objectiveTypeDictionary ObjectiveTypeDictionary) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(objectiveTypeDictionary)(conn)
		core.ErrorHandler(err2, "daos/ObjectiveTypeDictionary", fmt.Sprintf("could not create or update %v in %s\n", objectiveTypeDictionary, TableName(conn.ClientID)))
	}
}


// Deactivate "removes" ObjectiveTypeDictionary. 
// The mechanism is adding timestamp to `DeactivatedOn` field. 
// Then, if this field is not empty, the instance is considered to be "active"
func Deactivate(id string) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		instance, err2 := Read(id)(conn)
		if err2 == nil {
			instance.DeactivatedAt = core.CurrentRFCTimestamp()
			err2 = CreateOrUpdate(instance)(conn)
		}
		return err2
	}
}


// DeactivateUnsafe "deletes" ObjectiveTypeDictionary and panics in case of errors.
func DeactivateUnsafe(id string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Deactivate(id)(conn)
		core.ErrorHandler(err2, "daos/ObjectiveTypeDictionary", fmt.Sprintf("Could not deactivate id==%s in %s\n", id, TableName(conn.ClientID)))
	}
}


func ReadByPlatformID(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []ObjectiveTypeDictionary, err error) {
	return func (conn common.DynamoDBConnection) (out []ObjectiveTypeDictionary, err error) {
		var instances []ObjectiveTypeDictionary
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "PlatformIDIndex",
			Condition: "platform_id = :a0",
			Attributes: map[string]interface{}{
				":a0": platformID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = ObjectiveTypeDictionaryFilterActive(instances)
		return
	}
}


func ReadByPlatformIDUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []ObjectiveTypeDictionary) {
	return func (conn common.DynamoDBConnection) (out []ObjectiveTypeDictionary) {
		out, err2 := ReadByPlatformID(platformID)(conn)
		core.ErrorHandler(err2, "daos/ObjectiveTypeDictionary", fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}
