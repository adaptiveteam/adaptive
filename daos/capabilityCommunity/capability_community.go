package capabilityCommunity
// This file has been automatically generated by `adaptive-platform/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

type CapabilityCommunity struct  {
	ID string `json:"id"`
	PlatformID models.PlatformID `json:"platform_id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Advocate string `json:"advocate"`
	CreatedBy string `json:"created_by"`
	// Automatically maintained field
	CreatedAt string `json:"created_at"`
	// Automatically maintained field
	ModifiedAt string `json:"modified_at"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (capabilityCommunity CapabilityCommunity)CollectEmptyFields() (emptyFields []string, ok bool) {
	if capabilityCommunity.ID == "" { emptyFields = append(emptyFields, "ID")}
	if capabilityCommunity.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if capabilityCommunity.Name == "" { emptyFields = append(emptyFields, "Name")}
	if capabilityCommunity.Description == "" { emptyFields = append(emptyFields, "Description")}
	if capabilityCommunity.Advocate == "" { emptyFields = append(emptyFields, "Advocate")}
	if capabilityCommunity.CreatedBy == "" { emptyFields = append(emptyFields, "CreatedBy")}
	ok = len(emptyFields) == 0
	return
}

type DAO interface {
	Create(capabilityCommunity CapabilityCommunity) error
	CreateUnsafe(capabilityCommunity CapabilityCommunity)
	Read(id string) (capabilityCommunity CapabilityCommunity, err error)
	ReadUnsafe(id string) (capabilityCommunity CapabilityCommunity)
	ReadOrEmpty(id string) (capabilityCommunity []CapabilityCommunity, err error)
	ReadOrEmptyUnsafe(id string) (capabilityCommunity []CapabilityCommunity)
	CreateOrUpdate(capabilityCommunity CapabilityCommunity) error
	CreateOrUpdateUnsafe(capabilityCommunity CapabilityCommunity)
	Delete(id string) error
	DeleteUnsafe(id string)
	ReadByPlatformID(platformID models.PlatformID) (capabilityCommunity []CapabilityCommunity, err error)
	ReadByPlatformIDUnsafe(platformID models.PlatformID) (capabilityCommunity []CapabilityCommunity)
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
		Name: clientID + "_capability_community",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the CapabilityCommunity.
func (d DAOImpl) Create(capabilityCommunity CapabilityCommunity) error {
	emptyFields, ok := capabilityCommunity.CollectEmptyFields()
	if !ok {return fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)}
	capabilityCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())
	capabilityCommunity.CreatedAt = capabilityCommunity.ModifiedAt
	return d.Dynamo.PutTableEntry(capabilityCommunity, d.Name)
}


// CreateUnsafe saves the CapabilityCommunity.
func (d DAOImpl) CreateUnsafe(capabilityCommunity CapabilityCommunity) {
	err := d.Create(capabilityCommunity)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create id==%s in %s\n", capabilityCommunity.ID, d.Name))
}


// Read reads CapabilityCommunity
func (d DAOImpl) Read(id string) (out CapabilityCommunity, err error) {
	var outs []CapabilityCommunity
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, d.Name)
	}
	return
}


// ReadUnsafe reads the CapabilityCommunity. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) CapabilityCommunity {
	out, err := d.Read(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading id==%s in %s\n", id, d.Name))
	return out
}


// ReadOrEmpty reads CapabilityCommunity
func (d DAOImpl) ReadOrEmpty(id string) (out []CapabilityCommunity, err error) {
	var outOrEmpty CapabilityCommunity
	ids := idParams(id)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.ID == id {
		out = append(out, outOrEmpty)
	} else if err != nil && strings.HasPrefix(err.Error(), "In table ") {
		err = nil // expected not-found error	
	}
	err = errors.Wrapf(err, "CapabilityCommunity DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the CapabilityCommunity. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []CapabilityCommunity {
	out, err := d.ReadOrEmpty(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading id==%s in %s\n", id, d.Name))
	return out
}


// CreateOrUpdate saves the CapabilityCommunity regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(capabilityCommunity CapabilityCommunity) (err error) {
	capabilityCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())
	if capabilityCommunity.CreatedAt == "" { capabilityCommunity.CreatedAt = capabilityCommunity.ModifiedAt }
	
	var olds []CapabilityCommunity
	olds, err = d.ReadOrEmpty(capabilityCommunity.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(capabilityCommunity)
			err = errors.Wrapf(err, "CapabilityCommunity DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			old := olds[0]
			capabilityCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())

			key := idParams(old.ID)
			expr, exprAttributes, names := updateExpression(capabilityCommunity, old)
			input := dynamodb.UpdateItemInput{
				ExpressionAttributeValues: exprAttributes,
				TableName:                 aws.String(d.Name),
				Key:                       key,
				ReturnValues:              aws.String("UPDATED_NEW"),
				UpdateExpression:          aws.String(expr),
			}
			if names != nil { input.ExpressionAttributeNames = *names } // workaround for a pointer to an empty slice
			if err == nil {
				err = d.Dynamo.UpdateItemInternal(input)
			}
			err = errors.Wrapf(err, "CapabilityCommunity DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s", key, d.Name)
			return
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the CapabilityCommunity regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(capabilityCommunity CapabilityCommunity) {
	err := d.CreateOrUpdate(capabilityCommunity)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", capabilityCommunity, d.Name))
}


// Delete removes CapabilityCommunity from db
func (d DAOImpl)Delete(id string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(id))
}


// DeleteUnsafe deletes CapabilityCommunity and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(id string) {
	err := d.Delete(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete id==%s in %s\n", id, d.Name))
}


func (d DAOImpl)ReadByPlatformID(platformID models.PlatformID) (out []CapabilityCommunity, err error) {
	var instances []CapabilityCommunity
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "PlatformIDIndex",
		Condition: "platform_id = :a0",
		Attributes: map[string]interface{}{
			":a0": platformID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID models.PlatformID) (out []CapabilityCommunity) {
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
func allParams(capabilityCommunity CapabilityCommunity, old CapabilityCommunity) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if capabilityCommunity.ID != old.ID { params[":a0"] = common.DynS(capabilityCommunity.ID) }
	if capabilityCommunity.PlatformID != old.PlatformID { params[":a1"] = common.DynS(string(capabilityCommunity.PlatformID)) }
	if capabilityCommunity.Name != old.Name { params[":a2"] = common.DynS(capabilityCommunity.Name) }
	if capabilityCommunity.Description != old.Description { params[":a3"] = common.DynS(capabilityCommunity.Description) }
	if capabilityCommunity.Advocate != old.Advocate { params[":a4"] = common.DynS(capabilityCommunity.Advocate) }
	if capabilityCommunity.CreatedBy != old.CreatedBy { params[":a5"] = common.DynS(capabilityCommunity.CreatedBy) }
	if capabilityCommunity.CreatedAt != old.CreatedAt { params[":a6"] = common.DynS(capabilityCommunity.CreatedAt) }
	if capabilityCommunity.ModifiedAt != old.ModifiedAt { params[":a7"] = common.DynS(capabilityCommunity.ModifiedAt) }
	return
}
func updateExpression(capabilityCommunity CapabilityCommunity, old CapabilityCommunity) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if capabilityCommunity.ID != old.ID { updateParts = append(updateParts, "id = :a0"); params[":a0"] = common.DynS(capabilityCommunity.ID);  }
	if capabilityCommunity.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a1"); params[":a1"] = common.DynS(string(capabilityCommunity.PlatformID));  }
	if capabilityCommunity.Name != old.Name { updateParts = append(updateParts, "#name = :a2"); params[":a2"] = common.DynS(capabilityCommunity.Name); fldName := "name"; names["#name"] = &fldName }
	if capabilityCommunity.Description != old.Description { updateParts = append(updateParts, "description = :a3"); params[":a3"] = common.DynS(capabilityCommunity.Description);  }
	if capabilityCommunity.Advocate != old.Advocate { updateParts = append(updateParts, "advocate = :a4"); params[":a4"] = common.DynS(capabilityCommunity.Advocate);  }
	if capabilityCommunity.CreatedBy != old.CreatedBy { updateParts = append(updateParts, "created_by = :a5"); params[":a5"] = common.DynS(capabilityCommunity.CreatedBy);  }
	if capabilityCommunity.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a6"); params[":a6"] = common.DynS(capabilityCommunity.CreatedAt);  }
	if capabilityCommunity.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a7"); params[":a7"] = common.DynS(capabilityCommunity.ModifiedAt);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
