package visionMission
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/aws/aws-sdk-go/aws"
	common "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
)


// Create saves the VisionMission.
func Create(visionMission VisionMission) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := visionMission.CollectEmptyFields()
		if ok {
			visionMission.ModifiedAt = core.CurrentRFCTimestamp()
	visionMission.CreatedAt = visionMission.ModifiedAt
	
			err = conn.Dynamo.PutTableEntry(visionMission, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the VisionMission.
func CreateUnsafe(visionMission VisionMission) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(visionMission)(conn)
		core.ErrorHandler(err2, "daos/VisionMission", fmt.Sprintf("Could not create platformID==%s in %s\n", visionMission.PlatformID, TableName(conn.ClientID)))
	}
}


// Read reads VisionMission
func Read(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out VisionMission, err error) {
	return func (conn common.DynamoDBConnection) (out VisionMission, err error) {
		var outs []VisionMission
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


// ReadUnsafe reads the VisionMission. Panics in case of any errors
func ReadUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) VisionMission {
	return func (conn common.DynamoDBConnection) VisionMission {
		out, err2 := Read(platformID)(conn)
		core.ErrorHandler(err2, "daos/VisionMission", fmt.Sprintf("Error reading platformID==%s in %s\n", platformID, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads VisionMission
func ReadOrEmpty(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []VisionMission, err error) {
	return func (conn common.DynamoDBConnection) (out []VisionMission, err error) {
       out, err = ReadOrEmptyIncludingInactive(platformID)(conn)
       
       
		return
	}
}


// ReadOrEmptyUnsafe reads the VisionMission. Panics in case of any errors
func ReadOrEmptyUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) []VisionMission {
	return func (conn common.DynamoDBConnection) []VisionMission {
		out, err2 := ReadOrEmpty(platformID)(conn)
		core.ErrorHandler(err2, "daos/VisionMission", fmt.Sprintf("Error while reading platformID==%s in %s\n", platformID, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads VisionMission
func ReadOrEmptyIncludingInactive(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []VisionMission, err error) {
	return func (conn common.DynamoDBConnection) (out []VisionMission, err error) {
		var outOrEmpty VisionMission
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
		err = errors.Wrapf(err, "VisionMission DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the VisionMission. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(platformID common.PlatformID) func (conn common.DynamoDBConnection) []VisionMission {
	return func (conn common.DynamoDBConnection) []VisionMission {
		out, err2 := ReadOrEmptyIncludingInactive(platformID)(conn)
		core.ErrorHandler(err2, "daos/VisionMission", fmt.Sprintf("Error while reading platformID==%s in %s\n", platformID, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the VisionMission regardless of if it exists.
func CreateOrUpdate(visionMission VisionMission) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		visionMission.ModifiedAt = core.CurrentRFCTimestamp()
	if visionMission.CreatedAt == "" { visionMission.CreatedAt = visionMission.ModifiedAt }
	
		var olds []VisionMission
		olds, err = ReadOrEmpty(visionMission.PlatformID)(conn)
		err = errors.Wrapf(err, "VisionMission DAO.CreateOrUpdate(id = platformID==%s) couldn't ReadOrEmpty", visionMission.PlatformID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(visionMission)(conn)
				err = errors.Wrapf(err, "VisionMission DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := visionMission.CollectEmptyFields()
				if ok {
					old := olds[0]
					visionMission.CreatedAt  = old.CreatedAt
					visionMission.ModifiedAt = core.CurrentRFCTimestamp()
					key := idParams(old.PlatformID)
					expr, exprAttributes, names := updateExpression(visionMission, old)
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
					err = errors.Wrapf(err, "VisionMission DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the VisionMission regardless of if it exists.
func CreateOrUpdateUnsafe(visionMission VisionMission) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(visionMission)(conn)
		core.ErrorHandler(err2, "daos/VisionMission", fmt.Sprintf("could not create or update %v in %s\n", visionMission, TableName(conn.ClientID)))
	}
}


// Delete removes VisionMission from db
func Delete(platformID common.PlatformID) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams(platformID))
	}
}


// DeleteUnsafe deletes VisionMission and panics in case of errors.
func DeleteUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Delete(platformID)(conn)
		core.ErrorHandler(err2, "daos/VisionMission", fmt.Sprintf("Could not delete platformID==%s in %s\n", platformID, TableName(conn.ClientID)))
	}
}

