package slackTeam
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/aws/aws-sdk-go/aws"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
)


// Create saves the SlackTeam.
func Create(slackTeam SlackTeam) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := slackTeam.CollectEmptyFields()
		if ok {
			slackTeam.ModifiedAt = core.CurrentRFCTimestamp()
	slackTeam.CreatedAt = slackTeam.ModifiedAt
	
			err = conn.Dynamo.PutTableEntry(slackTeam, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the SlackTeam.
func CreateUnsafe(slackTeam SlackTeam) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(slackTeam)(conn)
		core.ErrorHandler(err2, "daos/SlackTeam", fmt.Sprintf("Could not create teamID==%s in %s\n", slackTeam.TeamID, TableName(conn.ClientID)))
	}
}


// Read reads SlackTeam
func Read(teamID common.PlatformID) func (conn common.DynamoDBConnection) (out SlackTeam, err error) {
	return func (conn common.DynamoDBConnection) (out SlackTeam, err error) {
		var outs []SlackTeam
		outs, err = ReadOrEmpty(teamID)(conn)
		if err == nil && len(outs) == 0 {
			err = fmt.Errorf("Not found teamID==%s in %s\n", teamID, TableName(conn.ClientID))
		}
		if len(outs) > 0 {
			out = outs[0]
		}
		return
	}
}


// ReadUnsafe reads the SlackTeam. Panics in case of any errors
func ReadUnsafe(teamID common.PlatformID) func (conn common.DynamoDBConnection) SlackTeam {
	return func (conn common.DynamoDBConnection) SlackTeam {
		out, err2 := Read(teamID)(conn)
		core.ErrorHandler(err2, "daos/SlackTeam", fmt.Sprintf("Error reading teamID==%s in %s\n", teamID, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads SlackTeam
func ReadOrEmpty(teamID common.PlatformID) func (conn common.DynamoDBConnection) (out []SlackTeam, err error) {
	return func (conn common.DynamoDBConnection) (out []SlackTeam, err error) {
       out, err = ReadOrEmptyIncludingInactive(teamID)(conn)
       
       
		return
	}
}


// ReadOrEmptyUnsafe reads the SlackTeam. Panics in case of any errors
func ReadOrEmptyUnsafe(teamID common.PlatformID) func (conn common.DynamoDBConnection) []SlackTeam {
	return func (conn common.DynamoDBConnection) []SlackTeam {
		out, err2 := ReadOrEmpty(teamID)(conn)
		core.ErrorHandler(err2, "daos/SlackTeam", fmt.Sprintf("Error while reading teamID==%s in %s\n", teamID, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads SlackTeam
func ReadOrEmptyIncludingInactive(teamID common.PlatformID) func (conn common.DynamoDBConnection) (out []SlackTeam, err error) {
	return func (conn common.DynamoDBConnection) (out []SlackTeam, err error) {
		var outOrEmpty SlackTeam
		ids := idParams(teamID)
		var found bool
		found, err = conn.Dynamo.GetItemOrEmptyFromTable(TableName(conn.ClientID), ids, &outOrEmpty)
		if found {
			if outOrEmpty.TeamID == teamID {
				out = append(out, outOrEmpty)
			} else {
				err = fmt.Errorf("Requested ids: teamID==%s are different from the found ones: teamID==%s", teamID, outOrEmpty.TeamID) // unexpected error: found ids != ids
			}
		}
		err = errors.Wrapf(err, "SlackTeam DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the SlackTeam. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(teamID common.PlatformID) func (conn common.DynamoDBConnection) []SlackTeam {
	return func (conn common.DynamoDBConnection) []SlackTeam {
		out, err2 := ReadOrEmptyIncludingInactive(teamID)(conn)
		core.ErrorHandler(err2, "daos/SlackTeam", fmt.Sprintf("Error while reading teamID==%s in %s\n", teamID, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the SlackTeam regardless of if it exists.
func CreateOrUpdate(slackTeam SlackTeam) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		slackTeam.ModifiedAt = core.CurrentRFCTimestamp()
	if slackTeam.CreatedAt == "" { slackTeam.CreatedAt = slackTeam.ModifiedAt }
	
		var olds []SlackTeam
		olds, err = ReadOrEmpty(slackTeam.TeamID)(conn)
		err = errors.Wrapf(err, "SlackTeam DAO.CreateOrUpdate(id = teamID==%s) couldn't ReadOrEmpty", slackTeam.TeamID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(slackTeam)(conn)
				err = errors.Wrapf(err, "SlackTeam DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := slackTeam.CollectEmptyFields()
				if ok {
					old := olds[0]
					slackTeam.CreatedAt  = old.CreatedAt
					slackTeam.ModifiedAt = core.CurrentRFCTimestamp()
					key := idParams(old.TeamID)
					expr, exprAttributes, names := updateExpression(slackTeam, old)
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
					err = errors.Wrapf(err, "SlackTeam DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the SlackTeam regardless of if it exists.
func CreateOrUpdateUnsafe(slackTeam SlackTeam) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(slackTeam)(conn)
		core.ErrorHandler(err2, "daos/SlackTeam", fmt.Sprintf("could not create or update %v in %s\n", slackTeam, TableName(conn.ClientID)))
	}
}


// Delete removes SlackTeam from db
func Delete(teamID common.PlatformID) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams(teamID))
	}
}


// DeleteUnsafe deletes SlackTeam and panics in case of errors.
func DeleteUnsafe(teamID common.PlatformID) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Delete(teamID)(conn)
		core.ErrorHandler(err2, "daos/SlackTeam", fmt.Sprintf("Could not delete teamID==%s in %s\n", teamID, TableName(conn.ClientID)))
	}
}

