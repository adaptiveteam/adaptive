package userFeedback
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


// Create saves the UserFeedback.
func Create(userFeedback UserFeedback) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := userFeedback.CollectEmptyFields()
		if ok {
			
			err = conn.Dynamo.PutTableEntry(userFeedback, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the UserFeedback.
func CreateUnsafe(userFeedback UserFeedback) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(userFeedback)(conn)
		core.ErrorHandler(err2, "daos/UserFeedback", fmt.Sprintf("Could not create id==%s in %s\n", userFeedback.ID, TableName(conn.ClientID)))
	}
}


// Read reads UserFeedback
func Read(id string) func (conn common.DynamoDBConnection) (out UserFeedback, err error) {
	return func (conn common.DynamoDBConnection) (out UserFeedback, err error) {
		var outs []UserFeedback
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


// ReadUnsafe reads the UserFeedback. Panics in case of any errors
func ReadUnsafe(id string) func (conn common.DynamoDBConnection) UserFeedback {
	return func (conn common.DynamoDBConnection) UserFeedback {
		out, err2 := Read(id)(conn)
		core.ErrorHandler(err2, "daos/UserFeedback", fmt.Sprintf("Error reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads UserFeedback
func ReadOrEmpty(id string) func (conn common.DynamoDBConnection) (out []UserFeedback, err error) {
	return func (conn common.DynamoDBConnection) (out []UserFeedback, err error) {
       out, err = ReadOrEmptyIncludingInactive(id)(conn)
       
       
		return
	}
}


// ReadOrEmptyUnsafe reads the UserFeedback. Panics in case of any errors
func ReadOrEmptyUnsafe(id string) func (conn common.DynamoDBConnection) []UserFeedback {
	return func (conn common.DynamoDBConnection) []UserFeedback {
		out, err2 := ReadOrEmpty(id)(conn)
		core.ErrorHandler(err2, "daos/UserFeedback", fmt.Sprintf("Error while reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads UserFeedback
func ReadOrEmptyIncludingInactive(id string) func (conn common.DynamoDBConnection) (out []UserFeedback, err error) {
	return func (conn common.DynamoDBConnection) (out []UserFeedback, err error) {
		var outOrEmpty UserFeedback
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
		err = errors.Wrapf(err, "UserFeedback DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the UserFeedback. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(id string) func (conn common.DynamoDBConnection) []UserFeedback {
	return func (conn common.DynamoDBConnection) []UserFeedback {
		out, err2 := ReadOrEmptyIncludingInactive(id)(conn)
		core.ErrorHandler(err2, "daos/UserFeedback", fmt.Sprintf("Error while reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the UserFeedback regardless of if it exists.
func CreateOrUpdate(userFeedback UserFeedback) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		
		var olds []UserFeedback
		olds, err = ReadOrEmpty(userFeedback.ID)(conn)
		err = errors.Wrapf(err, "UserFeedback DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", userFeedback.ID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(userFeedback)(conn)
				err = errors.Wrapf(err, "UserFeedback DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := userFeedback.CollectEmptyFields()
				if ok {
					old := olds[0]
					
					
					key := idParams(old.ID)
					expr, exprAttributes, names := updateExpression(userFeedback, old)
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
					err = errors.Wrapf(err, "UserFeedback DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the UserFeedback regardless of if it exists.
func CreateOrUpdateUnsafe(userFeedback UserFeedback) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(userFeedback)(conn)
		core.ErrorHandler(err2, "daos/UserFeedback", fmt.Sprintf("could not create or update %v in %s\n", userFeedback, TableName(conn.ClientID)))
	}
}


// Delete removes UserFeedback from db
func Delete(id string) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams(id))
	}
}


// DeleteUnsafe deletes UserFeedback and panics in case of errors.
func DeleteUnsafe(id string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Delete(id)(conn)
		core.ErrorHandler(err2, "daos/UserFeedback", fmt.Sprintf("Could not delete id==%s in %s\n", id, TableName(conn.ClientID)))
	}
}


func ReadByQuarterYearSource(quarterYear string, source string) func (conn common.DynamoDBConnection) (out []UserFeedback, err error) {
	return func (conn common.DynamoDBConnection) (out []UserFeedback, err error) {
		var instances []UserFeedback
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "QuarterYearSourceIndex",
			Condition: "quarter_year = :a0 and #source = :a1",
			Attributes: map[string]interface{}{
				":a0": quarterYear,
			":a1": source,
			},
		}, map[string]string{"#source": "source"}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByQuarterYearSourceUnsafe(quarterYear string, source string) func (conn common.DynamoDBConnection) (out []UserFeedback) {
	return func (conn common.DynamoDBConnection) (out []UserFeedback) {
		out, err2 := ReadByQuarterYearSource(quarterYear, source)(conn)
		core.ErrorHandler(err2, "daos/UserFeedback", fmt.Sprintf("Could not query QuarterYearSourceIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByQuarterYearTarget(quarterYear string, target string) func (conn common.DynamoDBConnection) (out []UserFeedback, err error) {
	return func (conn common.DynamoDBConnection) (out []UserFeedback, err error) {
		var instances []UserFeedback
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "QuarterYearTargetIndex",
			Condition: "quarter_year = :a0 and target = :a1",
			Attributes: map[string]interface{}{
				":a0": quarterYear,
			":a1": target,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByQuarterYearTargetUnsafe(quarterYear string, target string) func (conn common.DynamoDBConnection) (out []UserFeedback) {
	return func (conn common.DynamoDBConnection) (out []UserFeedback) {
		out, err2 := ReadByQuarterYearTarget(quarterYear, target)(conn)
		core.ErrorHandler(err2, "daos/UserFeedback", fmt.Sprintf("Could not query QuarterYearTargetIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByHashKeyQuarterYear(quarterYear string) func (conn common.DynamoDBConnection) (out []UserFeedback, err error) {
	return func (conn common.DynamoDBConnection) (out []UserFeedback, err error) {
		var instances []UserFeedback
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: string(QuarterYearSourceIndex),
			Condition: "quarter_year = :a",
			Attributes: map[string]interface{}{
				":a" : quarterYear,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByHashKeyQuarterYearUnsafe(quarterYear string) func (conn common.DynamoDBConnection) (out []UserFeedback) {
	return func (conn common.DynamoDBConnection) (out []UserFeedback) {
		out, err2 := ReadByHashKeyQuarterYear(quarterYear)(conn)
		core.ErrorHandler(err2, "daos/UserFeedback", fmt.Sprintf("Could not query QuarterYearSourceIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}

