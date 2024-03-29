package coachingRelationship
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


// Create saves the CoachingRelationship.
func Create(coachingRelationship CoachingRelationship) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := coachingRelationship.CollectEmptyFields()
		if ok {
			
			err = conn.Dynamo.PutTableEntry(coachingRelationship, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the CoachingRelationship.
func CreateUnsafe(coachingRelationship CoachingRelationship) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(coachingRelationship)(conn)
		core.ErrorHandler(err2, "daos/CoachingRelationship", fmt.Sprintf("Could not create coachQuarterYear==%s in %s\n", coachingRelationship.CoachQuarterYear, TableName(conn.ClientID)))
	}
}


// Read reads CoachingRelationship
func Read(coachQuarterYear string) func (conn common.DynamoDBConnection) (out CoachingRelationship, err error) {
	return func (conn common.DynamoDBConnection) (out CoachingRelationship, err error) {
		var outs []CoachingRelationship
		outs, err = ReadOrEmpty(coachQuarterYear)(conn)
		if err == nil && len(outs) == 0 {
			err = fmt.Errorf("Not found coachQuarterYear==%s in %s\n", coachQuarterYear, TableName(conn.ClientID))
		}
		if len(outs) > 0 {
			out = outs[0]
		}
		return
	}
}


// ReadUnsafe reads the CoachingRelationship. Panics in case of any errors
func ReadUnsafe(coachQuarterYear string) func (conn common.DynamoDBConnection) CoachingRelationship {
	return func (conn common.DynamoDBConnection) CoachingRelationship {
		out, err2 := Read(coachQuarterYear)(conn)
		core.ErrorHandler(err2, "daos/CoachingRelationship", fmt.Sprintf("Error reading coachQuarterYear==%s in %s\n", coachQuarterYear, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads CoachingRelationship
func ReadOrEmpty(coachQuarterYear string) func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
       out, err = ReadOrEmptyIncludingInactive(coachQuarterYear)(conn)
       
       
		return
	}
}


// ReadOrEmptyUnsafe reads the CoachingRelationship. Panics in case of any errors
func ReadOrEmptyUnsafe(coachQuarterYear string) func (conn common.DynamoDBConnection) []CoachingRelationship {
	return func (conn common.DynamoDBConnection) []CoachingRelationship {
		out, err2 := ReadOrEmpty(coachQuarterYear)(conn)
		core.ErrorHandler(err2, "daos/CoachingRelationship", fmt.Sprintf("Error while reading coachQuarterYear==%s in %s\n", coachQuarterYear, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads CoachingRelationship
func ReadOrEmptyIncludingInactive(coachQuarterYear string) func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
		var outOrEmpty CoachingRelationship
		ids := idParams(coachQuarterYear)
		var found bool
		found, err = conn.Dynamo.GetItemOrEmptyFromTable(TableName(conn.ClientID), ids, &outOrEmpty)
		if found {
			if outOrEmpty.CoachQuarterYear == coachQuarterYear {
				out = append(out, outOrEmpty)
			} else {
				err = fmt.Errorf("Requested ids: coachQuarterYear==%s are different from the found ones: coachQuarterYear==%s", coachQuarterYear, outOrEmpty.CoachQuarterYear) // unexpected error: found ids != ids
			}
		}
		err = errors.Wrapf(err, "CoachingRelationship DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the CoachingRelationship. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(coachQuarterYear string) func (conn common.DynamoDBConnection) []CoachingRelationship {
	return func (conn common.DynamoDBConnection) []CoachingRelationship {
		out, err2 := ReadOrEmptyIncludingInactive(coachQuarterYear)(conn)
		core.ErrorHandler(err2, "daos/CoachingRelationship", fmt.Sprintf("Error while reading coachQuarterYear==%s in %s\n", coachQuarterYear, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the CoachingRelationship regardless of if it exists.
func CreateOrUpdate(coachingRelationship CoachingRelationship) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		
		var olds []CoachingRelationship
		olds, err = ReadOrEmpty(coachingRelationship.CoachQuarterYear)(conn)
		err = errors.Wrapf(err, "CoachingRelationship DAO.CreateOrUpdate(id = coachQuarterYear==%s) couldn't ReadOrEmpty", coachingRelationship.CoachQuarterYear)
		if err == nil {
			if len(olds) == 0 {
				err = Create(coachingRelationship)(conn)
				err = errors.Wrapf(err, "CoachingRelationship DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := coachingRelationship.CollectEmptyFields()
				if ok {
					old := olds[0]
					
					
					key := idParams(old.CoachQuarterYear)
					expr, exprAttributes, names := updateExpression(coachingRelationship, old)
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
					err = errors.Wrapf(err, "CoachingRelationship DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the CoachingRelationship regardless of if it exists.
func CreateOrUpdateUnsafe(coachingRelationship CoachingRelationship) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(coachingRelationship)(conn)
		core.ErrorHandler(err2, "daos/CoachingRelationship", fmt.Sprintf("could not create or update %v in %s\n", coachingRelationship, TableName(conn.ClientID)))
	}
}


// Delete removes CoachingRelationship from db
func Delete(coachQuarterYear string) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams(coachQuarterYear))
	}
}


// DeleteUnsafe deletes CoachingRelationship and panics in case of errors.
func DeleteUnsafe(coachQuarterYear string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Delete(coachQuarterYear)(conn)
		core.ErrorHandler(err2, "daos/CoachingRelationship", fmt.Sprintf("Could not delete coachQuarterYear==%s in %s\n", coachQuarterYear, TableName(conn.ClientID)))
	}
}


func ReadByHashKeyCoachQuarterYear(coachQuarterYear string) func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
		var instances []CoachingRelationship
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			
			Condition: "coach_quarter_year = :a",
			Attributes: map[string]interface{}{
				":a" : coachQuarterYear,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByHashKeyCoachQuarterYearUnsafe(coachQuarterYear string) func (conn common.DynamoDBConnection) (out []CoachingRelationship) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship) {
		out, err2 := ReadByHashKeyCoachQuarterYear(coachQuarterYear)(conn)
		core.ErrorHandler(err2, "daos/CoachingRelationship", fmt.Sprintf("Could not query CoachQuarterYearCoacheeIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByCoachQuarterYear(coachQuarterYear string) func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
		var instances []CoachingRelationship
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "CoachQuarterYearIndex",
			Condition: "coach_quarter_year = :a0",
			Attributes: map[string]interface{}{
				":a0": coachQuarterYear,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByCoachQuarterYearUnsafe(coachQuarterYear string) func (conn common.DynamoDBConnection) (out []CoachingRelationship) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship) {
		out, err2 := ReadByCoachQuarterYear(coachQuarterYear)(conn)
		core.ErrorHandler(err2, "daos/CoachingRelationship", fmt.Sprintf("Could not query CoachQuarterYearIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByQuarterYear(quarter int, year int) func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
		var instances []CoachingRelationship
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "QuarterYearIndex",
			Condition: "quarter = :a0 and #year = :a1",
			Attributes: map[string]interface{}{
				":a0": quarter,
			":a1": year,
			},
		}, map[string]string{"#year": "year"}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByQuarterYearUnsafe(quarter int, year int) func (conn common.DynamoDBConnection) (out []CoachingRelationship) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship) {
		out, err2 := ReadByQuarterYear(quarter, year)(conn)
		core.ErrorHandler(err2, "daos/CoachingRelationship", fmt.Sprintf("Could not query QuarterYearIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByCoacheeQuarterYear(coacheeQuarterYear string) func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
		var instances []CoachingRelationship
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "CoacheeQuarterYearIndex",
			Condition: "coachee_quarter_year = :a0",
			Attributes: map[string]interface{}{
				":a0": coacheeQuarterYear,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByCoacheeQuarterYearUnsafe(coacheeQuarterYear string) func (conn common.DynamoDBConnection) (out []CoachingRelationship) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship) {
		out, err2 := ReadByCoacheeQuarterYear(coacheeQuarterYear)(conn)
		core.ErrorHandler(err2, "daos/CoachingRelationship", fmt.Sprintf("Could not query CoacheeQuarterYearIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByHashKeyQuarter(quarter int) func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship, err error) {
		var instances []CoachingRelationship
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: string(QuarterYearIndex),
			Condition: "quarter = :a",
			Attributes: map[string]interface{}{
				":a" : quarter,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByHashKeyQuarterUnsafe(quarter int) func (conn common.DynamoDBConnection) (out []CoachingRelationship) {
	return func (conn common.DynamoDBConnection) (out []CoachingRelationship) {
		out, err2 := ReadByHashKeyQuarter(quarter)(conn)
		core.ErrorHandler(err2, "daos/CoachingRelationship", fmt.Sprintf("Could not query QuarterYearIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}

