package adaptiveValue
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


// Create saves the AdaptiveValue.
func Create(adaptiveValue AdaptiveValue) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := adaptiveValue.CollectEmptyFields()
		if ok {
			
			err = conn.Dynamo.PutTableEntry(adaptiveValue, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the AdaptiveValue.
func CreateUnsafe(adaptiveValue AdaptiveValue) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(adaptiveValue)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveValue", fmt.Sprintf("Could not create id==%s in %s\n", adaptiveValue.ID, TableName(conn.ClientID)))
	}
}


// Read reads AdaptiveValue
func Read(id string) func (conn common.DynamoDBConnection) (out AdaptiveValue, err error) {
	return func (conn common.DynamoDBConnection) (out AdaptiveValue, err error) {
		var outs []AdaptiveValue
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


// ReadUnsafe reads the AdaptiveValue. Panics in case of any errors
func ReadUnsafe(id string) func (conn common.DynamoDBConnection) AdaptiveValue {
	return func (conn common.DynamoDBConnection) AdaptiveValue {
		out, err2 := Read(id)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveValue", fmt.Sprintf("Error reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads AdaptiveValue
func ReadOrEmpty(id string) func (conn common.DynamoDBConnection) (out []AdaptiveValue, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveValue, err error) {
       out, err = ReadOrEmptyIncludingInactive(id)(conn)
       out = AdaptiveValueFilterActive(out)
       
		return
	}
}


// ReadOrEmptyUnsafe reads the AdaptiveValue. Panics in case of any errors
func ReadOrEmptyUnsafe(id string) func (conn common.DynamoDBConnection) []AdaptiveValue {
	return func (conn common.DynamoDBConnection) []AdaptiveValue {
		out, err2 := ReadOrEmpty(id)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveValue", fmt.Sprintf("Error while reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads AdaptiveValue
func ReadOrEmptyIncludingInactive(id string) func (conn common.DynamoDBConnection) (out []AdaptiveValue, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveValue, err error) {
		var outOrEmpty AdaptiveValue
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
		err = errors.Wrapf(err, "AdaptiveValue DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the AdaptiveValue. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(id string) func (conn common.DynamoDBConnection) []AdaptiveValue {
	return func (conn common.DynamoDBConnection) []AdaptiveValue {
		out, err2 := ReadOrEmptyIncludingInactive(id)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveValue", fmt.Sprintf("Error while reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the AdaptiveValue regardless of if it exists.
func CreateOrUpdate(adaptiveValue AdaptiveValue) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		
		var olds []AdaptiveValue
		olds, err = ReadOrEmpty(adaptiveValue.ID)(conn)
		err = errors.Wrapf(err, "AdaptiveValue DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", adaptiveValue.ID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(adaptiveValue)(conn)
				err = errors.Wrapf(err, "AdaptiveValue DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := adaptiveValue.CollectEmptyFields()
				if ok {
					old := olds[0]
					
					
					key := idParams(old.ID)
					expr, exprAttributes, names := updateExpression(adaptiveValue, old)
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
					err = errors.Wrapf(err, "AdaptiveValue DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the AdaptiveValue regardless of if it exists.
func CreateOrUpdateUnsafe(adaptiveValue AdaptiveValue) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(adaptiveValue)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveValue", fmt.Sprintf("could not create or update %v in %s\n", adaptiveValue, TableName(conn.ClientID)))
	}
}


// Deactivate "removes" AdaptiveValue. 
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


// DeactivateUnsafe "deletes" AdaptiveValue and panics in case of errors.
func DeactivateUnsafe(id string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Deactivate(id)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveValue", fmt.Sprintf("Could not deactivate id==%s in %s\n", id, TableName(conn.ClientID)))
	}
}


func ReadByPlatformID(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []AdaptiveValue, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveValue, err error) {
		var instances []AdaptiveValue
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "PlatformIDIndex",
			Condition: "platform_id = :a0",
			Attributes: map[string]interface{}{
				":a0": platformID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = AdaptiveValueFilterActive(instances)
		return
	}
}


func ReadByPlatformIDUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []AdaptiveValue) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveValue) {
		out, err2 := ReadByPlatformID(platformID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveValue", fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}

