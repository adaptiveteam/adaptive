package capabilityCommunity
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


// Create saves the CapabilityCommunity.
func Create(capabilityCommunity CapabilityCommunity) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := capabilityCommunity.CollectEmptyFields()
		if ok {
			capabilityCommunity.ModifiedAt = core.CurrentRFCTimestamp()
	capabilityCommunity.CreatedAt = capabilityCommunity.ModifiedAt
	
			err = conn.Dynamo.PutTableEntry(capabilityCommunity, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the CapabilityCommunity.
func CreateUnsafe(capabilityCommunity CapabilityCommunity) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(capabilityCommunity)(conn)
		core.ErrorHandler(err2, "daos/CapabilityCommunity", fmt.Sprintf("Could not create platformID==%s, id==%s in %s\n", capabilityCommunity.PlatformID, capabilityCommunity.ID, TableName(conn.ClientID)))
	}
}


// Read reads CapabilityCommunity
func Read(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) (out CapabilityCommunity, err error) {
	return func (conn common.DynamoDBConnection) (out CapabilityCommunity, err error) {
		var outs []CapabilityCommunity
		outs, err = ReadOrEmpty(platformID, id)(conn)
		if err == nil && len(outs) == 0 {
			err = fmt.Errorf("Not found platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID))
		}
		if len(outs) > 0 {
			out = outs[0]
		}
		return
	}
}


// ReadUnsafe reads the CapabilityCommunity. Panics in case of any errors
func ReadUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) CapabilityCommunity {
	return func (conn common.DynamoDBConnection) CapabilityCommunity {
		out, err2 := Read(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/CapabilityCommunity", fmt.Sprintf("Error reading platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads CapabilityCommunity
func ReadOrEmpty(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) (out []CapabilityCommunity, err error) {
	return func (conn common.DynamoDBConnection) (out []CapabilityCommunity, err error) {
		var outOrEmpty CapabilityCommunity
		ids := idParams(platformID, id)
		var found bool
		found, err = conn.Dynamo.GetItemOrEmptyFromTable(TableName(conn.ClientID), ids, &outOrEmpty)
		if found {
			if outOrEmpty.PlatformID == platformID && outOrEmpty.ID == id {
				out = append(out, outOrEmpty)
			} else {
				err = fmt.Errorf("Requested ids: platformID==%s, id==%s are different from the found ones: platformID==%s, id==%s", platformID, id, outOrEmpty.PlatformID, outOrEmpty.ID) // unexpected error: found ids != ids
			}
		}
		err = errors.Wrapf(err, "CapabilityCommunity DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyUnsafe reads the CapabilityCommunity. Panics in case of any errors
func ReadOrEmptyUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) []CapabilityCommunity {
	return func (conn common.DynamoDBConnection) []CapabilityCommunity {
		out, err2 := ReadOrEmpty(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/CapabilityCommunity", fmt.Sprintf("Error while reading platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the CapabilityCommunity regardless of if it exists.
func CreateOrUpdate(capabilityCommunity CapabilityCommunity) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		capabilityCommunity.ModifiedAt = core.CurrentRFCTimestamp()
	if capabilityCommunity.CreatedAt == "" { capabilityCommunity.CreatedAt = capabilityCommunity.ModifiedAt }
	
		var olds []CapabilityCommunity
		olds, err = ReadOrEmpty(capabilityCommunity.PlatformID, capabilityCommunity.ID)(conn)
		err = errors.Wrapf(err, "CapabilityCommunity DAO.CreateOrUpdate(id = platformID==%s, id==%s) couldn't ReadOrEmpty", capabilityCommunity.PlatformID, capabilityCommunity.ID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(capabilityCommunity)(conn)
				err = errors.Wrapf(err, "CapabilityCommunity DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := capabilityCommunity.CollectEmptyFields()
				if ok {
					old := olds[0]
					capabilityCommunity.CreatedAt  = old.CreatedAt
					capabilityCommunity.ModifiedAt = core.CurrentRFCTimestamp()
					key := idParams(old.PlatformID, old.ID)
					expr, exprAttributes, names := updateExpression(capabilityCommunity, old)
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
					err = errors.Wrapf(err, "CapabilityCommunity DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the CapabilityCommunity regardless of if it exists.
func CreateOrUpdateUnsafe(capabilityCommunity CapabilityCommunity) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(capabilityCommunity)(conn)
		core.ErrorHandler(err2, "daos/CapabilityCommunity", fmt.Sprintf("could not create or update %v in %s\n", capabilityCommunity, TableName(conn.ClientID)))
	}
}


// Delete removes CapabilityCommunity from db
func Delete(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams(platformID, id))
	}
}


// DeleteUnsafe deletes CapabilityCommunity and panics in case of errors.
func DeleteUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Delete(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/CapabilityCommunity", fmt.Sprintf("Could not delete platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
	}
}


func ReadByPlatformID(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []CapabilityCommunity, err error) {
	return func (conn common.DynamoDBConnection) (out []CapabilityCommunity, err error) {
		var instances []CapabilityCommunity
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "PlatformIDIndex",
			Condition: "platform_id = :a0",
			Attributes: map[string]interface{}{
				":a0": platformID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByPlatformIDUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []CapabilityCommunity) {
	return func (conn common.DynamoDBConnection) (out []CapabilityCommunity) {
		out, err2 := ReadByPlatformID(platformID)(conn)
		core.ErrorHandler(err2, "daos/CapabilityCommunity", fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}

