package adaptiveValue
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

// AdaptiveValue is a value for a client (Reliability, Skill, Contribution, and Productivity)
type AdaptiveValue struct  {
	ID string `json:"id"`
	PlatformID common.PlatformID `json:"platform_id"`
	Name string `json:"value_name"`
	ValueType string `json:"value_type"`
	Description string `json:"description"`
	DeactivatedAt string `json:"deactivated_at,omitempty"`
}

// AdaptiveValueFilterActive removes deactivated values
func AdaptiveValueFilterActive(in []AdaptiveValue) (res []AdaptiveValue) {
	for _, i := range in {
		if i.DeactivatedAt == "" {
			res = append(res, i)
		}
	}
	return
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (adaptiveValue AdaptiveValue)CollectEmptyFields() (emptyFields []string, ok bool) {
	if adaptiveValue.ID == "" { emptyFields = append(emptyFields, "ID")}
	if adaptiveValue.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if adaptiveValue.Name == "" { emptyFields = append(emptyFields, "Name")}
	if adaptiveValue.ValueType == "" { emptyFields = append(emptyFields, "ValueType")}
	if adaptiveValue.Description == "" { emptyFields = append(emptyFields, "Description")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (adaptiveValue AdaptiveValue) ToJSON() (string, error) {
	b, err := json.Marshal(adaptiveValue)
	return string(b), err
}

type DAO interface {
	Create(adaptiveValue AdaptiveValue) error
	CreateUnsafe(adaptiveValue AdaptiveValue)
	Read(id string) (adaptiveValue AdaptiveValue, err error)
	ReadUnsafe(id string) (adaptiveValue AdaptiveValue)
	ReadOrEmpty(id string) (adaptiveValue []AdaptiveValue, err error)
	ReadOrEmptyUnsafe(id string) (adaptiveValue []AdaptiveValue)
	CreateOrUpdate(adaptiveValue AdaptiveValue) error
	CreateOrUpdateUnsafe(adaptiveValue AdaptiveValue)
	Deactivate(id string) error
	DeactivateUnsafe(id string)
	ReadByPlatformID(platformID common.PlatformID) (adaptiveValue []AdaptiveValue, err error)
	ReadByPlatformIDUnsafe(platformID common.PlatformID) (adaptiveValue []AdaptiveValue)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	ConnGen   common.DynamoDBConnectionGen
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic(errors.New("Cannot create AdaptiveValue.DAO without clientID")) }
	return DAOImpl{
		ConnGen:   common.DynamoDBConnectionGen{
			Dynamo: dynamo, 
			TableNamePrefix: clientID,
		},
	}
}
// TableNameSuffixVar is a global variable that contains table name suffix.
// After renaming all tables this may be made `const`.
var TableNameSuffixVar = "_adaptive_value"

// TableName concatenates table name prefix and suffix and returns table name
func TableName(prefix string) string {
	return prefix + TableNameSuffixVar
}

// Create saves the AdaptiveValue.
func (d DAOImpl) Create(adaptiveValue AdaptiveValue) (err error) {
	emptyFields, ok := adaptiveValue.CollectEmptyFields()
	if ok {
		err = d.ConnGen.Dynamo.PutTableEntry(adaptiveValue, TableName(d.ConnGen.TableNamePrefix))
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the AdaptiveValue.
func (d DAOImpl) CreateUnsafe(adaptiveValue AdaptiveValue) {
	err2 := d.Create(adaptiveValue)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not create id==%s in %s\n", adaptiveValue.ID, TableName(d.ConnGen.TableNamePrefix)))
}


// Read reads AdaptiveValue
func (d DAOImpl) Read(id string) (out AdaptiveValue, err error) {
	var outs []AdaptiveValue
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, TableName(d.ConnGen.TableNamePrefix))
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the AdaptiveValue. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) AdaptiveValue {
	out, err2 := d.Read(id)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error reading id==%s in %s\n", id, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// ReadOrEmpty reads AdaptiveValue
func (d DAOImpl) ReadOrEmpty(id string) (out []AdaptiveValue, err error) {
	var outOrEmpty AdaptiveValue
	ids := idParams(id)
	var found bool
	found, err = d.ConnGen.Dynamo.GetItemOrEmptyFromTable(TableName(d.ConnGen.TableNamePrefix), ids, &outOrEmpty)
	if found {
		if outOrEmpty.ID == id {
			out = append(out, outOrEmpty)
		} else {
			err = fmt.Errorf("Requested ids: id==%s are different from the found ones: id==%s", id, outOrEmpty.ID) // unexpected error: found ids != ids
		}
	}
	err = errors.Wrapf(err, "AdaptiveValue DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(d.ConnGen.TableNamePrefix))
	return
}


// ReadOrEmptyUnsafe reads the AdaptiveValue. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []AdaptiveValue {
	out, err2 := d.ReadOrEmpty(id)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error while reading id==%s in %s\n", id, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// CreateOrUpdate saves the AdaptiveValue regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(adaptiveValue AdaptiveValue) (err error) {
	
	var olds []AdaptiveValue
	olds, err = d.ReadOrEmpty(adaptiveValue.ID)
	err = errors.Wrapf(err, "AdaptiveValue DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", adaptiveValue.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(adaptiveValue)
			err = errors.Wrapf(err, "AdaptiveValue DAO.CreateOrUpdate couldn't Create in table %s", TableName(d.ConnGen.TableNamePrefix))
		} else {
			emptyFields, ok := adaptiveValue.CollectEmptyFields()
			if ok {
				old := olds[0]
				
				
				key := idParams(old.ID)
				expr, exprAttributes, names := updateExpression(adaptiveValue, old)
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
				err = errors.Wrapf(err, "AdaptiveValue DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(d.ConnGen.TableNamePrefix), expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the AdaptiveValue regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(adaptiveValue AdaptiveValue) {
	err2 := d.CreateOrUpdate(adaptiveValue)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("could not create or update %v in %s\n", adaptiveValue, TableName(d.ConnGen.TableNamePrefix)))
}


// Deactivate "removes" AdaptiveValue. 
// The mechanism is adding timestamp to `DeactivatedOn` field. 
// Then, if this field is not empty, the instance is considered to be "active"
func (d DAOImpl)Deactivate(id string) error {
	instance, err2 := d.Read(id)
	if err2 == nil {
		instance.DeactivatedAt = core.CurrentRFCTimestamp()
		err2 = d.CreateOrUpdate(instance)
	}
	return err2
}


// DeactivateUnsafe "deletes" AdaptiveValue and panics in case of errors.
func (d DAOImpl)DeactivateUnsafe(id string) {
	err2 := d.Deactivate(id)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not deactivate id==%s in %s\n", id, TableName(d.ConnGen.TableNamePrefix)))
}


func (d DAOImpl)ReadByPlatformID(platformID common.PlatformID) (out []AdaptiveValue, err error) {
	var instances []AdaptiveValue
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
		IndexName: "PlatformIDIndex",
		Condition: "platform_id = :a0",
		Attributes: map[string]interface{}{
			":a0": platformID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = AdaptiveValueFilterActive(instances)
	return
}


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID common.PlatformID) (out []AdaptiveValue) {
	out, err2 := d.ReadByPlatformID(platformID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
	return
}

func idParams(id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"id": common.DynS(id),
	}
	return params
}
func allParams(adaptiveValue AdaptiveValue, old AdaptiveValue) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if adaptiveValue.ID != old.ID { params[":a0"] = common.DynS(adaptiveValue.ID) }
	if adaptiveValue.PlatformID != old.PlatformID { params[":a1"] = common.DynS(string(adaptiveValue.PlatformID)) }
	if adaptiveValue.Name != old.Name { params[":a2"] = common.DynS(adaptiveValue.Name) }
	if adaptiveValue.ValueType != old.ValueType { params[":a3"] = common.DynS(adaptiveValue.ValueType) }
	if adaptiveValue.Description != old.Description { params[":a4"] = common.DynS(adaptiveValue.Description) }
	if adaptiveValue.DeactivatedAt != old.DeactivatedAt { params[":a5"] = common.DynS(adaptiveValue.DeactivatedAt) }
	return
}
func updateExpression(adaptiveValue AdaptiveValue, old AdaptiveValue) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if adaptiveValue.ID != old.ID { updateParts = append(updateParts, "id = :a0"); params[":a0"] = common.DynS(adaptiveValue.ID);  }
	if adaptiveValue.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a1"); params[":a1"] = common.DynS(string(adaptiveValue.PlatformID));  }
	if adaptiveValue.Name != old.Name { updateParts = append(updateParts, "value_name = :a2"); params[":a2"] = common.DynS(adaptiveValue.Name);  }
	if adaptiveValue.ValueType != old.ValueType { updateParts = append(updateParts, "value_type = :a3"); params[":a3"] = common.DynS(adaptiveValue.ValueType);  }
	if adaptiveValue.Description != old.Description { updateParts = append(updateParts, "description = :a4"); params[":a4"] = common.DynS(adaptiveValue.Description);  }
	if adaptiveValue.DeactivatedAt != old.DeactivatedAt { updateParts = append(updateParts, "deactivated_at = :a5"); params[":a5"] = common.DynS(adaptiveValue.DeactivatedAt);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
