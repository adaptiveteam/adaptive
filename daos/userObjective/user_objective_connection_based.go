package userObjective
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


// Create saves the UserObjective.
func Create(userObjective UserObjective) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := userObjective.CollectEmptyFields()
		if ok {
			userObjective.ModifiedAt = core.CurrentRFCTimestamp()
	userObjective.CreatedAt = userObjective.ModifiedAt
	
			err = conn.Dynamo.PutTableEntry(userObjective, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the UserObjective.
func CreateUnsafe(userObjective UserObjective) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(userObjective)(conn)
		core.ErrorHandler(err2, "daos/UserObjective", fmt.Sprintf("Could not create id==%s in %s\n", userObjective.ID, TableName(conn.ClientID)))
	}
}


// Read reads UserObjective
func Read(id string) func (conn common.DynamoDBConnection) (out UserObjective, err error) {
	return func (conn common.DynamoDBConnection) (out UserObjective, err error) {
		var outs []UserObjective
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


// ReadUnsafe reads the UserObjective. Panics in case of any errors
func ReadUnsafe(id string) func (conn common.DynamoDBConnection) UserObjective {
	return func (conn common.DynamoDBConnection) UserObjective {
		out, err2 := Read(id)(conn)
		core.ErrorHandler(err2, "daos/UserObjective", fmt.Sprintf("Error reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads UserObjective
func ReadOrEmpty(id string) func (conn common.DynamoDBConnection) (out []UserObjective, err error) {
	return func (conn common.DynamoDBConnection) (out []UserObjective, err error) {
		var outOrEmpty UserObjective
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
		err = errors.Wrapf(err, "UserObjective DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyUnsafe reads the UserObjective. Panics in case of any errors
func ReadOrEmptyUnsafe(id string) func (conn common.DynamoDBConnection) []UserObjective {
	return func (conn common.DynamoDBConnection) []UserObjective {
		out, err2 := ReadOrEmpty(id)(conn)
		core.ErrorHandler(err2, "daos/UserObjective", fmt.Sprintf("Error while reading id==%s in %s\n", id, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the UserObjective regardless of if it exists.
func CreateOrUpdate(userObjective UserObjective) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		userObjective.ModifiedAt = core.CurrentRFCTimestamp()
	if userObjective.CreatedAt == "" { userObjective.CreatedAt = userObjective.ModifiedAt }
	
		var olds []UserObjective
		olds, err = ReadOrEmpty(userObjective.ID)(conn)
		err = errors.Wrapf(err, "UserObjective DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", userObjective.ID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(userObjective)(conn)
				err = errors.Wrapf(err, "UserObjective DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := userObjective.CollectEmptyFields()
				if ok {
					old := olds[0]
					userObjective.CreatedAt  = old.CreatedAt
					userObjective.ModifiedAt = core.CurrentRFCTimestamp()
					key := idParams(old.ID)
					expr, exprAttributes, names := updateExpression(userObjective, old)
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
					err = errors.Wrapf(err, "UserObjective DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the UserObjective regardless of if it exists.
func CreateOrUpdateUnsafe(userObjective UserObjective) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(userObjective)(conn)
		core.ErrorHandler(err2, "daos/UserObjective", fmt.Sprintf("could not create or update %v in %s\n", userObjective, TableName(conn.ClientID)))
	}
}


// Delete removes UserObjective from db
func Delete(id string) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams(id))
	}
}


// DeleteUnsafe deletes UserObjective and panics in case of errors.
func DeleteUnsafe(id string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Delete(id)(conn)
		core.ErrorHandler(err2, "daos/UserObjective", fmt.Sprintf("Could not delete id==%s in %s\n", id, TableName(conn.ClientID)))
	}
}


func ReadByUserIDCompleted(userID string, completed int) func (conn common.DynamoDBConnection) (out []UserObjective, err error) {
	return func (conn common.DynamoDBConnection) (out []UserObjective, err error) {
		var instances []UserObjective
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "UserIDCompletedIndex",
			Condition: "user_id = :a0 and completed = :a1",
			Attributes: map[string]interface{}{
				":a0": userID,
			":a1": completed,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByUserIDCompletedUnsafe(userID string, completed int) func (conn common.DynamoDBConnection) (out []UserObjective) {
	return func (conn common.DynamoDBConnection) (out []UserObjective) {
		out, err2 := ReadByUserIDCompleted(userID, completed)(conn)
		core.ErrorHandler(err2, "daos/UserObjective", fmt.Sprintf("Could not query UserIDCompletedIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByAccepted(accepted int) func (conn common.DynamoDBConnection) (out []UserObjective, err error) {
	return func (conn common.DynamoDBConnection) (out []UserObjective, err error) {
		var instances []UserObjective
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "AcceptedIndex",
			Condition: "accepted = :a0",
			Attributes: map[string]interface{}{
				":a0": accepted,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByAcceptedUnsafe(accepted int) func (conn common.DynamoDBConnection) (out []UserObjective) {
	return func (conn common.DynamoDBConnection) (out []UserObjective) {
		out, err2 := ReadByAccepted(accepted)(conn)
		core.ErrorHandler(err2, "daos/UserObjective", fmt.Sprintf("Could not query AcceptedIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByAccountabilityPartner(accountabilityPartner string) func (conn common.DynamoDBConnection) (out []UserObjective, err error) {
	return func (conn common.DynamoDBConnection) (out []UserObjective, err error) {
		var instances []UserObjective
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "AccountabilityPartnerIndex",
			Condition: "accountability_partner = :a0",
			Attributes: map[string]interface{}{
				":a0": accountabilityPartner,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByAccountabilityPartnerUnsafe(accountabilityPartner string) func (conn common.DynamoDBConnection) (out []UserObjective) {
	return func (conn common.DynamoDBConnection) (out []UserObjective) {
		out, err2 := ReadByAccountabilityPartner(accountabilityPartner)(conn)
		core.ErrorHandler(err2, "daos/UserObjective", fmt.Sprintf("Could not query AccountabilityPartnerIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByUserIDType(userID string, objectiveType DevelopmentObjectiveType) func (conn common.DynamoDBConnection) (out []UserObjective, err error) {
	return func (conn common.DynamoDBConnection) (out []UserObjective, err error) {
		var instances []UserObjective
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "UserIDTypeIndex",
			Condition: "user_id = :a0 and #type = :a1",
			Attributes: map[string]interface{}{
				":a0": userID,
			":a1": objectiveType,
			},
		}, map[string]string{"#type": "type"}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByUserIDTypeUnsafe(userID string, objectiveType DevelopmentObjectiveType) func (conn common.DynamoDBConnection) (out []UserObjective) {
	return func (conn common.DynamoDBConnection) (out []UserObjective) {
		out, err2 := ReadByUserIDType(userID, objectiveType)(conn)
		core.ErrorHandler(err2, "daos/UserObjective", fmt.Sprintf("Could not query UserIDTypeIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}
