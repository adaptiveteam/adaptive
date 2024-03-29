package adaptiveCommunity
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


// Create saves the AdaptiveCommunity.
func Create(adaptiveCommunity AdaptiveCommunity) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := adaptiveCommunity.CollectEmptyFields()
		if ok {
			adaptiveCommunity.ModifiedAt = core.CurrentRFCTimestamp()
	adaptiveCommunity.CreatedAt = adaptiveCommunity.ModifiedAt
	
			err = conn.Dynamo.PutTableEntry(adaptiveCommunity, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the AdaptiveCommunity.
func CreateUnsafe(adaptiveCommunity AdaptiveCommunity) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(adaptiveCommunity)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunity", fmt.Sprintf("Could not create platformID==%s, id==%s in %s\n", adaptiveCommunity.PlatformID, adaptiveCommunity.ID, TableName(conn.ClientID)))
	}
}


// Read reads AdaptiveCommunity
func Read(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) (out AdaptiveCommunity, err error) {
	return func (conn common.DynamoDBConnection) (out AdaptiveCommunity, err error) {
		var outs []AdaptiveCommunity
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


// ReadUnsafe reads the AdaptiveCommunity. Panics in case of any errors
func ReadUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) AdaptiveCommunity {
	return func (conn common.DynamoDBConnection) AdaptiveCommunity {
		out, err2 := Read(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunity", fmt.Sprintf("Error reading platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads AdaptiveCommunity
func ReadOrEmpty(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunity, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunity, err error) {
       out, err = ReadOrEmptyIncludingInactive(platformID, id)(conn)
       out = AdaptiveCommunityFilterActive(out)
       
		return
	}
}


// ReadOrEmptyUnsafe reads the AdaptiveCommunity. Panics in case of any errors
func ReadOrEmptyUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) []AdaptiveCommunity {
	return func (conn common.DynamoDBConnection) []AdaptiveCommunity {
		out, err2 := ReadOrEmpty(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunity", fmt.Sprintf("Error while reading platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads AdaptiveCommunity
func ReadOrEmptyIncludingInactive(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunity, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunity, err error) {
		var outOrEmpty AdaptiveCommunity
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
		err = errors.Wrapf(err, "AdaptiveCommunity DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the AdaptiveCommunity. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) []AdaptiveCommunity {
	return func (conn common.DynamoDBConnection) []AdaptiveCommunity {
		out, err2 := ReadOrEmptyIncludingInactive(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunity", fmt.Sprintf("Error while reading platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the AdaptiveCommunity regardless of if it exists.
func CreateOrUpdate(adaptiveCommunity AdaptiveCommunity) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		adaptiveCommunity.ModifiedAt = core.CurrentRFCTimestamp()
	if adaptiveCommunity.CreatedAt == "" { adaptiveCommunity.CreatedAt = adaptiveCommunity.ModifiedAt }
	
		var olds []AdaptiveCommunity
		olds, err = ReadOrEmpty(adaptiveCommunity.PlatformID, adaptiveCommunity.ID)(conn)
		err = errors.Wrapf(err, "AdaptiveCommunity DAO.CreateOrUpdate(id = platformID==%s, id==%s) couldn't ReadOrEmpty", adaptiveCommunity.PlatformID, adaptiveCommunity.ID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(adaptiveCommunity)(conn)
				err = errors.Wrapf(err, "AdaptiveCommunity DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := adaptiveCommunity.CollectEmptyFields()
				if ok {
					old := olds[0]
					adaptiveCommunity.CreatedAt  = old.CreatedAt
					adaptiveCommunity.ModifiedAt = core.CurrentRFCTimestamp()
					key := idParams(old.PlatformID, old.ID)
					expr, exprAttributes, names := updateExpression(adaptiveCommunity, old)
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
					err = errors.Wrapf(err, "AdaptiveCommunity DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the AdaptiveCommunity regardless of if it exists.
func CreateOrUpdateUnsafe(adaptiveCommunity AdaptiveCommunity) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(adaptiveCommunity)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunity", fmt.Sprintf("could not create or update %v in %s\n", adaptiveCommunity, TableName(conn.ClientID)))
	}
}


// Deactivate "removes" AdaptiveCommunity. 
// The mechanism is adding timestamp to `DeactivatedOn` field. 
// Then, if this field is not empty, the instance is considered to be "active"
func Deactivate(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		instance, err2 := Read(platformID, id)(conn)
		if err2 == nil {
			instance.DeactivatedAt = core.CurrentRFCTimestamp()
			err2 = CreateOrUpdate(instance)(conn)
		}
		return err2
	}
}


// DeactivateUnsafe "deletes" AdaptiveCommunity and panics in case of errors.
func DeactivateUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Deactivate(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunity", fmt.Sprintf("Could not deactivate platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
	}
}


func ReadByHashKeyID(id string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunity, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunity, err error) {
		var instances []AdaptiveCommunity
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			
			Condition: "id = :a",
			Attributes: map[string]interface{}{
				":a" : id,
			},
		}, map[string]string{}, true, -1, &instances)
		out = AdaptiveCommunityFilterActive(instances)
		return
	}
}


func ReadByHashKeyIDUnsafe(id string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunity) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunity) {
		out, err2 := ReadByHashKeyID(id)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunity", fmt.Sprintf("Could not query IDPlatformIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByChannel(channelID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunity, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunity, err error) {
		var instances []AdaptiveCommunity
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "ChannelIndex",
			Condition: "channel = :a0",
			Attributes: map[string]interface{}{
				":a0": channelID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = AdaptiveCommunityFilterActive(instances)
		return
	}
}


func ReadByChannelUnsafe(channelID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunity) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunity) {
		out, err2 := ReadByChannel(channelID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunity", fmt.Sprintf("Could not query ChannelIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByPlatformID(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []AdaptiveCommunity, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunity, err error) {
		var instances []AdaptiveCommunity
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "PlatformIDIndex",
			Condition: "platform_id = :a0",
			Attributes: map[string]interface{}{
				":a0": platformID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = AdaptiveCommunityFilterActive(instances)
		return
	}
}


func ReadByPlatformIDUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []AdaptiveCommunity) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunity) {
		out, err2 := ReadByPlatformID(platformID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunity", fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}

