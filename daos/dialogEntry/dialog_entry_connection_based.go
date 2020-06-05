package dialogEntry
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
)


// Create saves the DialogEntry.
func Create(dialogEntry DialogEntry) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := dialogEntry.CollectEmptyFields()
		if ok {
			
			err = conn.Dynamo.PutTableEntry(dialogEntry, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the DialogEntry.
func CreateUnsafe(dialogEntry DialogEntry) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(dialogEntry)(conn)
		core.ErrorHandler(err2, "daos/DialogEntry", fmt.Sprintf("Could not create dialogID==%s in %s\n", dialogEntry.DialogID, TableName(conn.ClientID)))
	}
}


// Read reads DialogEntry
func Read(dialogID string) func (conn common.DynamoDBConnection) (out DialogEntry, err error) {
	return func (conn common.DynamoDBConnection) (out DialogEntry, err error) {
		var outs []DialogEntry
		outs, err = ReadOrEmpty(dialogID)(conn)
		if err == nil && len(outs) == 0 {
			err = fmt.Errorf("Not found dialogID==%s in %s\n", dialogID, TableName(conn.ClientID))
		}
		if len(outs) > 0 {
			out = outs[0]
		}
		return
	}
}


// ReadUnsafe reads the DialogEntry. Panics in case of any errors
func ReadUnsafe(dialogID string) func (conn common.DynamoDBConnection) DialogEntry {
	return func (conn common.DynamoDBConnection) DialogEntry {
		out, err2 := Read(dialogID)(conn)
		core.ErrorHandler(err2, "daos/DialogEntry", fmt.Sprintf("Error reading dialogID==%s in %s\n", dialogID, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads DialogEntry
func ReadOrEmpty(dialogID string) func (conn common.DynamoDBConnection) (out []DialogEntry, err error) {
	return func (conn common.DynamoDBConnection) (out []DialogEntry, err error) {
       out, err = ReadOrEmptyIncludingInactive(dialogID)(conn)
       
       
		return
	}
}


// ReadOrEmptyUnsafe reads the DialogEntry. Panics in case of any errors
func ReadOrEmptyUnsafe(dialogID string) func (conn common.DynamoDBConnection) []DialogEntry {
	return func (conn common.DynamoDBConnection) []DialogEntry {
		out, err2 := ReadOrEmpty(dialogID)(conn)
		core.ErrorHandler(err2, "daos/DialogEntry", fmt.Sprintf("Error while reading dialogID==%s in %s\n", dialogID, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads DialogEntry
func ReadOrEmptyIncludingInactive(dialogID string) func (conn common.DynamoDBConnection) (out []DialogEntry, err error) {
	return func (conn common.DynamoDBConnection) (out []DialogEntry, err error) {
		var outOrEmpty DialogEntry
		ids := idParams(dialogID)
		var found bool
		found, err = conn.Dynamo.GetItemOrEmptyFromTable(TableName(conn.ClientID), ids, &outOrEmpty)
		if found {
			if outOrEmpty.DialogID == dialogID {
				out = append(out, outOrEmpty)
			} else {
				err = fmt.Errorf("Requested ids: dialogID==%s are different from the found ones: dialogID==%s", dialogID, outOrEmpty.DialogID) // unexpected error: found ids != ids
			}
		}
		err = errors.Wrapf(err, "DialogEntry DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the DialogEntry. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(dialogID string) func (conn common.DynamoDBConnection) []DialogEntry {
	return func (conn common.DynamoDBConnection) []DialogEntry {
		out, err2 := ReadOrEmptyIncludingInactive(dialogID)(conn)
		core.ErrorHandler(err2, "daos/DialogEntry", fmt.Sprintf("Error while reading dialogID==%s in %s\n", dialogID, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the DialogEntry regardless of if it exists.
func CreateOrUpdate(dialogEntry DialogEntry) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		
		var olds []DialogEntry
		olds, err = ReadOrEmpty(dialogEntry.DialogID)(conn)
		err = errors.Wrapf(err, "DialogEntry DAO.CreateOrUpdate(id = dialogID==%s) couldn't ReadOrEmpty", dialogEntry.DialogID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(dialogEntry)(conn)
				err = errors.Wrapf(err, "DialogEntry DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := dialogEntry.CollectEmptyFields()
				if ok {
					old := olds[0]
					
					
					key := idParams(old.DialogID)
					expr, exprAttributes, names := updateExpression(dialogEntry, old)
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
					err = errors.Wrapf(err, "DialogEntry DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the DialogEntry regardless of if it exists.
func CreateOrUpdateUnsafe(dialogEntry DialogEntry) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(dialogEntry)(conn)
		core.ErrorHandler(err2, "daos/DialogEntry", fmt.Sprintf("could not create or update %v in %s\n", dialogEntry, TableName(conn.ClientID)))
	}
}


// Delete removes DialogEntry from db
func Delete(dialogID string) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams(dialogID))
	}
}


// DeleteUnsafe deletes DialogEntry and panics in case of errors.
func DeleteUnsafe(dialogID string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Delete(dialogID)(conn)
		core.ErrorHandler(err2, "daos/DialogEntry", fmt.Sprintf("Could not delete dialogID==%s in %s\n", dialogID, TableName(conn.ClientID)))
	}
}


func ReadByContextSubject(context string, subject string) func (conn common.DynamoDBConnection) (out []DialogEntry, err error) {
	return func (conn common.DynamoDBConnection) (out []DialogEntry, err error) {
		var instances []DialogEntry
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "ContextSubjectIndex",
			Condition: "context = :a0 and subject = :a1",
			Attributes: map[string]interface{}{
				":a0": context,
			":a1": subject,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByContextSubjectUnsafe(context string, subject string) func (conn common.DynamoDBConnection) (out []DialogEntry) {
	return func (conn common.DynamoDBConnection) (out []DialogEntry) {
		out, err2 := ReadByContextSubject(context, subject)(conn)
		core.ErrorHandler(err2, "daos/DialogEntry", fmt.Sprintf("Could not query ContextSubjectIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByHashKeyContext(context string) func (conn common.DynamoDBConnection) (out []DialogEntry, err error) {
	return func (conn common.DynamoDBConnection) (out []DialogEntry, err error) {
		var instances []DialogEntry
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: string(ContextSubjectIndex),
			Condition: "context = :a",
			Attributes: map[string]interface{}{
				":a" : context,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByHashKeyContextUnsafe(context string) func (conn common.DynamoDBConnection) (out []DialogEntry) {
	return func (conn common.DynamoDBConnection) (out []DialogEntry) {
		out, err2 := ReadByHashKeyContext(context)(conn)
		core.ErrorHandler(err2, "daos/DialogEntry", fmt.Sprintf("Could not query ContextSubjectIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}

