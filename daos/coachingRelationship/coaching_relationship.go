package coachingRelationship
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"encoding/json"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
)

type CoachingRelationship struct  {
	CoachQuarterYear string `json:"coach_quarter_year"`
	CoacheeQuarterYear string `json:"coachee_quarter_year"`
	Coachee string `json:"coachee"`
	Quarter int `json:"quarter"`
	Year int `json:"year"`
	CoachRequested bool `json:"coach_requested"`
	CoacheeRequested bool `json:"coachee_requested"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (coachingRelationship CoachingRelationship)CollectEmptyFields() (emptyFields []string, ok bool) {
	if coachingRelationship.CoachQuarterYear == "" { emptyFields = append(emptyFields, "CoachQuarterYear")}
	if coachingRelationship.CoacheeQuarterYear == "" { emptyFields = append(emptyFields, "CoacheeQuarterYear")}
	if coachingRelationship.Coachee == "" { emptyFields = append(emptyFields, "Coachee")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (coachingRelationship CoachingRelationship) ToJSON() (string, error) {
	b, err := json.Marshal(coachingRelationship)
	return string(b), err
}

type DAO interface {
	Create(coachingRelationship CoachingRelationship) error
	CreateUnsafe(coachingRelationship CoachingRelationship)
	Read(coachQuarterYear string) (coachingRelationship CoachingRelationship, err error)
	ReadUnsafe(coachQuarterYear string) (coachingRelationship CoachingRelationship)
	ReadOrEmpty(coachQuarterYear string) (coachingRelationship []CoachingRelationship, err error)
	ReadOrEmptyUnsafe(coachQuarterYear string) (coachingRelationship []CoachingRelationship)
	CreateOrUpdate(coachingRelationship CoachingRelationship) error
	CreateOrUpdateUnsafe(coachingRelationship CoachingRelationship)
	Delete(coachQuarterYear string) error
	DeleteUnsafe(coachQuarterYear string)
	ReadByCoachQuarterYear(coachQuarterYear string) (coachingRelationship []CoachingRelationship, err error)
	ReadByCoachQuarterYearUnsafe(coachQuarterYear string) (coachingRelationship []CoachingRelationship)
	ReadByQuarterYear(quarter int, year int) (coachingRelationship []CoachingRelationship, err error)
	ReadByQuarterYearUnsafe(quarter int, year int) (coachingRelationship []CoachingRelationship)
	ReadByCoacheeQuarterYear(coacheeQuarterYear string) (coachingRelationship []CoachingRelationship, err error)
	ReadByCoacheeQuarterYearUnsafe(coacheeQuarterYear string) (coachingRelationship []CoachingRelationship)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	ConnGen   common.DynamoDBConnectionGen
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic(errors.New("Cannot create CoachingRelationship.DAO without clientID")) }
	return DAOImpl{
		ConnGen:   common.DynamoDBConnectionGen{
			Dynamo: dynamo, 
			TableNamePrefix: clientID,
		},
	}
}
// TableNameSuffixVar is a global variable that contains table name suffix.
// After renaming all tables this may be made `const`.
var TableNameSuffixVar = "_coaching_relationship"

// TableName concatenates table name prefix and suffix and returns table name
func TableName(prefix string) string {
	return prefix + TableNameSuffixVar
}

// Create saves the CoachingRelationship.
func (d DAOImpl) Create(coachingRelationship CoachingRelationship) (err error) {
	emptyFields, ok := coachingRelationship.CollectEmptyFields()
	if ok {
		err = d.ConnGen.Dynamo.PutTableEntry(coachingRelationship, TableName(d.ConnGen.TableNamePrefix))
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the CoachingRelationship.
func (d DAOImpl) CreateUnsafe(coachingRelationship CoachingRelationship) {
	err2 := d.Create(coachingRelationship)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not create coachQuarterYear==%s in %s\n", coachingRelationship.CoachQuarterYear, TableName(d.ConnGen.TableNamePrefix)))
}


// Read reads CoachingRelationship
func (d DAOImpl) Read(coachQuarterYear string) (out CoachingRelationship, err error) {
	var outs []CoachingRelationship
	outs, err = d.ReadOrEmpty(coachQuarterYear)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found coachQuarterYear==%s in %s\n", coachQuarterYear, TableName(d.ConnGen.TableNamePrefix))
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the CoachingRelationship. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(coachQuarterYear string) CoachingRelationship {
	out, err2 := d.Read(coachQuarterYear)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error reading coachQuarterYear==%s in %s\n", coachQuarterYear, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// ReadOrEmpty reads CoachingRelationship
func (d DAOImpl) ReadOrEmpty(coachQuarterYear string) (out []CoachingRelationship, err error) {
	var outOrEmpty CoachingRelationship
	ids := idParams(coachQuarterYear)
	var found bool
	found, err = d.ConnGen.Dynamo.GetItemOrEmptyFromTable(TableName(d.ConnGen.TableNamePrefix), ids, &outOrEmpty)
	if found {
		if outOrEmpty.CoachQuarterYear == coachQuarterYear {
			out = append(out, outOrEmpty)
		} else {
			err = fmt.Errorf("Requested ids: coachQuarterYear==%s are different from the found ones: coachQuarterYear==%s", coachQuarterYear, outOrEmpty.CoachQuarterYear) // unexpected error: found ids != ids
		}
	}
	err = errors.Wrapf(err, "CoachingRelationship DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(d.ConnGen.TableNamePrefix))
	return
}


// ReadOrEmptyUnsafe reads the CoachingRelationship. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(coachQuarterYear string) []CoachingRelationship {
	out, err2 := d.ReadOrEmpty(coachQuarterYear)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error while reading coachQuarterYear==%s in %s\n", coachQuarterYear, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// CreateOrUpdate saves the CoachingRelationship regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(coachingRelationship CoachingRelationship) (err error) {
	
	var olds []CoachingRelationship
	olds, err = d.ReadOrEmpty(coachingRelationship.CoachQuarterYear)
	err = errors.Wrapf(err, "CoachingRelationship DAO.CreateOrUpdate(id = coachQuarterYear==%s) couldn't ReadOrEmpty", coachingRelationship.CoachQuarterYear)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(coachingRelationship)
			err = errors.Wrapf(err, "CoachingRelationship DAO.CreateOrUpdate couldn't Create in table %s", TableName(d.ConnGen.TableNamePrefix))
		} else {
			emptyFields, ok := coachingRelationship.CollectEmptyFields()
			if ok {
				old := olds[0]
				
				
				key := idParams(old.CoachQuarterYear)
				expr, exprAttributes, names := updateExpression(coachingRelationship, old)
				input := dynamodb.UpdateItemInput{
					ExpressionAttributeValues: exprAttributes,
					TableName:                 aws.String(TableName(d.ConnGen.TableNamePrefix)),
					Key:                       key,
					ReturnValues:              aws.String("UPDATED_NEW"),
					UpdateExpression:          aws.String(expr),
				}
				if names != nil { input.ExpressionAttributeNames = *names } // workaround for a pointer to an empty slice
				if  len(exprAttributes) > 0 { // if there some changes
					err = d.ConnGen.Dynamo.UpdateItemInternal(input)
				} else {
					// WARN: no changes.
				}
				err = errors.Wrapf(err, "CoachingRelationship DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(d.ConnGen.TableNamePrefix), expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the CoachingRelationship regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(coachingRelationship CoachingRelationship) {
	err2 := d.CreateOrUpdate(coachingRelationship)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("could not create or update %v in %s\n", coachingRelationship, TableName(d.ConnGen.TableNamePrefix)))
}


// Delete removes CoachingRelationship from db
func (d DAOImpl)Delete(coachQuarterYear string) error {
	return d.ConnGen.Dynamo.DeleteEntry(TableName(d.ConnGen.TableNamePrefix), idParams(coachQuarterYear))
}


// DeleteUnsafe deletes CoachingRelationship and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(coachQuarterYear string) {
	err2 := d.Delete(coachQuarterYear)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not delete coachQuarterYear==%s in %s\n", coachQuarterYear, TableName(d.ConnGen.TableNamePrefix)))
}


func (d DAOImpl)ReadByCoachQuarterYear(coachQuarterYear string) (out []CoachingRelationship, err error) {
	var instances []CoachingRelationship
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
		IndexName: "CoachQuarterYearIndex",
		Condition: "coach_quarter_year = :a0",
		Attributes: map[string]interface{}{
			":a0": coachQuarterYear,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByCoachQuarterYearUnsafe(coachQuarterYear string) (out []CoachingRelationship) {
	out, err2 := d.ReadByCoachQuarterYear(coachQuarterYear)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query CoachQuarterYearIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
	return
}


func (d DAOImpl)ReadByQuarterYear(quarter int, year int) (out []CoachingRelationship, err error) {
	var instances []CoachingRelationship
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
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


func (d DAOImpl)ReadByQuarterYearUnsafe(quarter int, year int) (out []CoachingRelationship) {
	out, err2 := d.ReadByQuarterYear(quarter, year)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query QuarterYearIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
	return
}


func (d DAOImpl)ReadByCoacheeQuarterYear(coacheeQuarterYear string) (out []CoachingRelationship, err error) {
	var instances []CoachingRelationship
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
		IndexName: "CoacheeQuarterYearIndex",
		Condition: "coachee_quarter_year = :a0",
		Attributes: map[string]interface{}{
			":a0": coacheeQuarterYear,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByCoacheeQuarterYearUnsafe(coacheeQuarterYear string) (out []CoachingRelationship) {
	out, err2 := d.ReadByCoacheeQuarterYear(coacheeQuarterYear)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query CoacheeQuarterYearIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
	return
}

func idParams(coachQuarterYear string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"coach_quarter_year": common.DynS(coachQuarterYear),
	}
	return params
}
func allParams(coachingRelationship CoachingRelationship, old CoachingRelationship) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if coachingRelationship.CoachQuarterYear != old.CoachQuarterYear { params[":a0"] = common.DynS(coachingRelationship.CoachQuarterYear) }
	if coachingRelationship.CoacheeQuarterYear != old.CoacheeQuarterYear { params[":a1"] = common.DynS(coachingRelationship.CoacheeQuarterYear) }
	if coachingRelationship.Coachee != old.Coachee { params[":a2"] = common.DynS(coachingRelationship.Coachee) }
	if coachingRelationship.Quarter != old.Quarter { params[":a3"] = common.DynN(coachingRelationship.Quarter) }
	if coachingRelationship.Year != old.Year { params[":a4"] = common.DynN(coachingRelationship.Year) }
	if coachingRelationship.CoachRequested != old.CoachRequested { params[":a5"] = common.DynBOOL(coachingRelationship.CoachRequested) }
	if coachingRelationship.CoacheeRequested != old.CoacheeRequested { params[":a6"] = common.DynBOOL(coachingRelationship.CoacheeRequested) }
	return
}
func updateExpression(coachingRelationship CoachingRelationship, old CoachingRelationship) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if coachingRelationship.CoachQuarterYear != old.CoachQuarterYear { updateParts = append(updateParts, "coach_quarter_year = :a0"); params[":a0"] = common.DynS(coachingRelationship.CoachQuarterYear);  }
	if coachingRelationship.CoacheeQuarterYear != old.CoacheeQuarterYear { updateParts = append(updateParts, "coachee_quarter_year = :a1"); params[":a1"] = common.DynS(coachingRelationship.CoacheeQuarterYear);  }
	if coachingRelationship.Coachee != old.Coachee { updateParts = append(updateParts, "coachee = :a2"); params[":a2"] = common.DynS(coachingRelationship.Coachee);  }
	if coachingRelationship.Quarter != old.Quarter { updateParts = append(updateParts, "quarter = :a3"); params[":a3"] = common.DynN(coachingRelationship.Quarter);  }
	if coachingRelationship.Year != old.Year { updateParts = append(updateParts, "#year = :a4"); params[":a4"] = common.DynN(coachingRelationship.Year); fldName := "year"; names["#year"] = &fldName }
	if coachingRelationship.CoachRequested != old.CoachRequested { updateParts = append(updateParts, "coach_requested = :a5"); params[":a5"] = common.DynBOOL(coachingRelationship.CoachRequested);  }
	if coachingRelationship.CoacheeRequested != old.CoacheeRequested { updateParts = append(updateParts, "coachee_requested = :a6"); params[":a6"] = common.DynBOOL(coachingRelationship.CoacheeRequested);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
