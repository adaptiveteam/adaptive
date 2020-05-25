package strategyObjective
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


// Create saves the StrategyObjective.
func Create(strategyObjective StrategyObjective) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := strategyObjective.CollectEmptyFields()
		if ok {
			strategyObjective.ModifiedAt = core.CurrentRFCTimestamp()
	strategyObjective.CreatedAt = strategyObjective.ModifiedAt
	
			err = conn.Dynamo.PutTableEntry(strategyObjective, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the StrategyObjective.
func CreateUnsafe(strategyObjective StrategyObjective) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(strategyObjective)(conn)
		core.ErrorHandler(err2, "daos/StrategyObjective", fmt.Sprintf("Could not create platformID==%s, id==%s in %s\n", strategyObjective.PlatformID, strategyObjective.ID, TableName(conn.ClientID)))
	}
}


// Read reads StrategyObjective
func Read(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) (out StrategyObjective, err error) {
	return func (conn common.DynamoDBConnection) (out StrategyObjective, err error) {
		var outs []StrategyObjective
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


// ReadUnsafe reads the StrategyObjective. Panics in case of any errors
func ReadUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) StrategyObjective {
	return func (conn common.DynamoDBConnection) StrategyObjective {
		out, err2 := Read(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/StrategyObjective", fmt.Sprintf("Error reading platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads StrategyObjective
func ReadOrEmpty(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) (out []StrategyObjective, err error) {
	return func (conn common.DynamoDBConnection) (out []StrategyObjective, err error) {
       out, err = ReadOrEmptyIncludingInactive(platformID, id)(conn)
       
       
		return
	}
}


// ReadOrEmptyUnsafe reads the StrategyObjective. Panics in case of any errors
func ReadOrEmptyUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) []StrategyObjective {
	return func (conn common.DynamoDBConnection) []StrategyObjective {
		out, err2 := ReadOrEmpty(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/StrategyObjective", fmt.Sprintf("Error while reading platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads StrategyObjective
func ReadOrEmptyIncludingInactive(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) (out []StrategyObjective, err error) {
	return func (conn common.DynamoDBConnection) (out []StrategyObjective, err error) {
		var outOrEmpty StrategyObjective
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
		err = errors.Wrapf(err, "StrategyObjective DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the StrategyObjective. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) []StrategyObjective {
	return func (conn common.DynamoDBConnection) []StrategyObjective {
		out, err2 := ReadOrEmptyIncludingInactive(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/StrategyObjective", fmt.Sprintf("Error while reading platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the StrategyObjective regardless of if it exists.
func CreateOrUpdate(strategyObjective StrategyObjective) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		strategyObjective.ModifiedAt = core.CurrentRFCTimestamp()
	if strategyObjective.CreatedAt == "" { strategyObjective.CreatedAt = strategyObjective.ModifiedAt }
	
		var olds []StrategyObjective
		olds, err = ReadOrEmpty(strategyObjective.PlatformID, strategyObjective.ID)(conn)
		err = errors.Wrapf(err, "StrategyObjective DAO.CreateOrUpdate(id = platformID==%s, id==%s) couldn't ReadOrEmpty", strategyObjective.PlatformID, strategyObjective.ID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(strategyObjective)(conn)
				err = errors.Wrapf(err, "StrategyObjective DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := strategyObjective.CollectEmptyFields()
				if ok {
					old := olds[0]
					strategyObjective.CreatedAt  = old.CreatedAt
					strategyObjective.ModifiedAt = core.CurrentRFCTimestamp()
					key := idParams(old.PlatformID, old.ID)
					expr, exprAttributes, names := updateExpression(strategyObjective, old)
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
					err = errors.Wrapf(err, "StrategyObjective DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the StrategyObjective regardless of if it exists.
func CreateOrUpdateUnsafe(strategyObjective StrategyObjective) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(strategyObjective)(conn)
		core.ErrorHandler(err2, "daos/StrategyObjective", fmt.Sprintf("could not create or update %v in %s\n", strategyObjective, TableName(conn.ClientID)))
	}
}


// Delete removes StrategyObjective from db
func Delete(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams(platformID, id))
	}
}


// DeleteUnsafe deletes StrategyObjective and panics in case of errors.
func DeleteUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Delete(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/StrategyObjective", fmt.Sprintf("Could not delete platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
	}
}


func ReadByHashKeyID(id string) func (conn common.DynamoDBConnection) (out []StrategyObjective, err error) {
	return func (conn common.DynamoDBConnection) (out []StrategyObjective, err error) {
		var instances []StrategyObjective
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			
			Condition: "id = :a",
			Attributes: map[string]interface{}{
				":a" : id,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByHashKeyIDUnsafe(id string) func (conn common.DynamoDBConnection) (out []StrategyObjective) {
	return func (conn common.DynamoDBConnection) (out []StrategyObjective) {
		out, err2 := ReadByHashKeyID(id)(conn)
		core.ErrorHandler(err2, "daos/StrategyObjective", fmt.Sprintf("Could not query IDPlatformIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByPlatformID(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []StrategyObjective, err error) {
	return func (conn common.DynamoDBConnection) (out []StrategyObjective, err error) {
		var instances []StrategyObjective
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


func ReadByPlatformIDUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []StrategyObjective) {
	return func (conn common.DynamoDBConnection) (out []StrategyObjective) {
		out, err2 := ReadByPlatformID(platformID)(conn)
		core.ErrorHandler(err2, "daos/StrategyObjective", fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByCapabilityCommunityIDs(capabilityCommunityIDs []string) func (conn common.DynamoDBConnection) (out []StrategyObjective, err error) {
	return func (conn common.DynamoDBConnection) (out []StrategyObjective, err error) {
		var instances []StrategyObjective
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "CapabilityCommunityIDsIndex",
			Condition: "capability_community_ids = :a0",
			Attributes: map[string]interface{}{
				":a0": capabilityCommunityIDs,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByCapabilityCommunityIDsUnsafe(capabilityCommunityIDs []string) func (conn common.DynamoDBConnection) (out []StrategyObjective) {
	return func (conn common.DynamoDBConnection) (out []StrategyObjective) {
		out, err2 := ReadByCapabilityCommunityIDs(capabilityCommunityIDs)(conn)
		core.ErrorHandler(err2, "daos/StrategyObjective", fmt.Sprintf("Could not query CapabilityCommunityIDsIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}

