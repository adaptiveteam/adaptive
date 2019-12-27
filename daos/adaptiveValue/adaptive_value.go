package adaptiveValue
// This file has been automatically generated by `adaptive-platform/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"time"
	"github.com/adaptiveteam/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/aws-utils-go"
	common "github.com/adaptiveteam/daos/common"
	core "github.com/adaptiveteam/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

// AdaptiveValue is a value for a client (Reliability, Skill, Contribution, and Productivity)
type AdaptiveValue struct  {
	ID string `json:"id"`
	PlatformID models.PlatformID `json:"platform_id"`
	Name string `json:"value_name"`
	ValueType string `json:"value_type"`
	Description string `json:"description"`
	DeactivatedOn string `json:"deactivated_on"`
}

// AdaptiveValueFilterActive removes deactivated values
func AdaptiveValueFilterActive(in []AdaptiveValue) (res []AdaptiveValue) {
	for _, i := range in {
		if i.DeactivatedOn == "" {
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
	ReadByID(id string) (adaptiveValue []AdaptiveValue, err error)
	ReadByIDUnsafe(id string) (adaptiveValue []AdaptiveValue)
	ReadByPlatformID(platformID models.PlatformID) (adaptiveValue []AdaptiveValue, err error)
	ReadByPlatformIDUnsafe(platformID models.PlatformID) (adaptiveValue []AdaptiveValue)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	Name      string                  `json:"name"`
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic("Cannot create DAO without clientID") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: clientID + "_adaptive_value",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the AdaptiveValue.
func (d DAOImpl) Create(adaptiveValue AdaptiveValue) error {
	emptyFields, ok := adaptiveValue.CollectEmptyFields()
	if !ok {return fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)}
	return d.Dynamo.PutTableEntry(adaptiveValue, d.Name)
}


// CreateUnsafe saves the AdaptiveValue.
func (d DAOImpl) CreateUnsafe(adaptiveValue AdaptiveValue) {
	err := d.Create(adaptiveValue)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create id==%s in %s\n", adaptiveValue.ID, d.Name))
}


// Read reads AdaptiveValue
func (d DAOImpl) Read(id string) (out AdaptiveValue, err error) {
	var outs []AdaptiveValue
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, d.Name)
	}
	return
}


// ReadUnsafe reads the AdaptiveValue. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) AdaptiveValue {
	out, err := d.Read(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading id==%s in %s\n", id, d.Name))
	return out
}


// ReadOrEmpty reads AdaptiveValue
func (d DAOImpl) ReadOrEmpty(id string) (out []AdaptiveValue, err error) {
	var outOrEmpty AdaptiveValue
	ids := idParams(id)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.ID == id {
		out = append(out, outOrEmpty)
	}
	err = errors.Wrapf(err, "AdaptiveValue DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the AdaptiveValue. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []AdaptiveValue {
	out, err := d.ReadOrEmpty(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading id==%s in %s\n", id, d.Name))
	return out
}


// CreateOrUpdate saves the AdaptiveValue regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(adaptiveValue AdaptiveValue) (err error) {
	
	var olds []AdaptiveValue
	olds, err = d.ReadOrEmpty(adaptiveValue.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(adaptiveValue)
			err = errors.Wrapf(err, "AdaptiveValue DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			old := olds[0]
			ids := idParams(old.ID)
			err = d.Dynamo.UpdateTableEntry(
				allParams(adaptiveValue, old),
				ids,
				updateExpression(adaptiveValue, old),
				d.Name,
			)
			err = errors.Wrapf(err, "AdaptiveValue DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s", ids, d.Name)
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the AdaptiveValue regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(adaptiveValue AdaptiveValue) {
	err := d.CreateOrUpdate(adaptiveValue)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", adaptiveValue, d.Name))
}


// Deactivate "removes" AdaptiveValue. 
// The mechanism is adding timestamp to `DeactivatedOn` field. 
// Then, if this field is not empty, the instance is considered to be "active"
func (d DAOImpl)Deactivate(id string) error {
	instance, err := d.Read(id)
	if err == nil {
		instance.DeactivatedOn = core.ISODateLayout.Format(time.Now())
		err = d.CreateOrUpdate(instance)
	}
	return err
}


// DeactivateUnsafe "deletes" AdaptiveValue and panics in case of errors.
func (d DAOImpl)DeactivateUnsafe(id string) {
	err := d.Deactivate(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not deactivate id==%s in %s\n", id, d.Name))
}


func (d DAOImpl)ReadByID(id string) (out []AdaptiveValue, err error) {
	var instances []AdaptiveValue
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "IDIndex",
		Condition: "id = :a0",
		Attributes: map[string]interface{}{
			":a0": id,
		},
	}, map[string]string{}, true, -1, &instances)
	out = AdaptiveValueFilterActive(instances)
	return
}


func (d DAOImpl)ReadByIDUnsafe(id string) (out []AdaptiveValue) {
	out, err := d.ReadByID(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query IDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByPlatformID(platformID models.PlatformID) (out []AdaptiveValue, err error) {
	var instances []AdaptiveValue
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "PlatformIDIndex",
		Condition: "platform_id = :a0",
		Attributes: map[string]interface{}{
			":a0": platformID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = AdaptiveValueFilterActive(instances)
	return
}


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID models.PlatformID) (out []AdaptiveValue) {
	out, err := d.ReadByPlatformID(platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", d.Name))
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
		if adaptiveValue.ID != old.ID { params["a0"] = common.DynS(adaptiveValue.ID) }
		if adaptiveValue.PlatformID != old.PlatformID { params["a1"] = common.DynS(string(adaptiveValue.PlatformID)) }
		if adaptiveValue.Name != old.Name { params["a2"] = common.DynS(adaptiveValue.Name) }
		if adaptiveValue.ValueType != old.ValueType { params["a3"] = common.DynS(adaptiveValue.ValueType) }
		if adaptiveValue.Description != old.Description { params["a4"] = common.DynS(adaptiveValue.Description) }
		if adaptiveValue.DeactivatedOn != old.DeactivatedOn { params["a5"] = common.DynS(adaptiveValue.DeactivatedOn) }
	return
}
func updateExpression(adaptiveValue AdaptiveValue, old AdaptiveValue) string {
	var updateParts []string
	
		
			
		if adaptiveValue.ID != old.ID { updateParts = append(updateParts, "id = :a0") }
		if adaptiveValue.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a1") }
		if adaptiveValue.Name != old.Name { updateParts = append(updateParts, "value_name = :a2") }
		if adaptiveValue.ValueType != old.ValueType { updateParts = append(updateParts, "value_type = :a3") }
		if adaptiveValue.Description != old.Description { updateParts = append(updateParts, "description = :a4") }
		if adaptiveValue.DeactivatedOn != old.DeactivatedOn { updateParts = append(updateParts, "deactivated_on = :a5") }
	return strings.Join(updateParts, " and ")
}
