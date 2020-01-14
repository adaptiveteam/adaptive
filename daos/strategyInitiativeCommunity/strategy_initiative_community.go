package strategyInitiativeCommunity
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"time"
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"encoding/json"
	"strings"
)

type StrategyInitiativeCommunity struct  {
	ID string `json:"id"`
	PlatformID common.PlatformID `json:"platform_id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Advocate string `json:"advocate"`
	CapabilityCommunityID string `json:"capability_community_id"`
	CreatedBy string `json:"created_by"`
	// Automatically maintained field
	CreatedAt string `json:"created_at"`
	// Automatically maintained field
	ModifiedAt string `json:"modified_at"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (strategyInitiativeCommunity StrategyInitiativeCommunity)CollectEmptyFields() (emptyFields []string, ok bool) {
	if strategyInitiativeCommunity.ID == "" { emptyFields = append(emptyFields, "ID")}
	if strategyInitiativeCommunity.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if strategyInitiativeCommunity.Name == "" { emptyFields = append(emptyFields, "Name")}
	if strategyInitiativeCommunity.Description == "" { emptyFields = append(emptyFields, "Description")}
	if strategyInitiativeCommunity.Advocate == "" { emptyFields = append(emptyFields, "Advocate")}
	if strategyInitiativeCommunity.CapabilityCommunityID == "" { emptyFields = append(emptyFields, "CapabilityCommunityID")}
	if strategyInitiativeCommunity.CreatedBy == "" { emptyFields = append(emptyFields, "CreatedBy")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (strategyInitiativeCommunity StrategyInitiativeCommunity) ToJSON() (string, error) {
	b, err := json.Marshal(strategyInitiativeCommunity)
	return string(b), err
}

type DAO interface {
	Create(strategyInitiativeCommunity StrategyInitiativeCommunity) error
	CreateUnsafe(strategyInitiativeCommunity StrategyInitiativeCommunity)
	Read(id string) (strategyInitiativeCommunity StrategyInitiativeCommunity, err error)
	ReadUnsafe(id string) (strategyInitiativeCommunity StrategyInitiativeCommunity)
	ReadOrEmpty(id string) (strategyInitiativeCommunity []StrategyInitiativeCommunity, err error)
	ReadOrEmptyUnsafe(id string) (strategyInitiativeCommunity []StrategyInitiativeCommunity)
	CreateOrUpdate(strategyInitiativeCommunity StrategyInitiativeCommunity) error
	CreateOrUpdateUnsafe(strategyInitiativeCommunity StrategyInitiativeCommunity)
	Delete(id string) error
	DeleteUnsafe(id string)
	ReadByPlatformID(platformID common.PlatformID) (strategyInitiativeCommunity []StrategyInitiativeCommunity, err error)
	ReadByPlatformIDUnsafe(platformID common.PlatformID) (strategyInitiativeCommunity []StrategyInitiativeCommunity)
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
		Name: clientID + "_strategy_initiative_community",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the StrategyInitiativeCommunity.
func (d DAOImpl) Create(strategyInitiativeCommunity StrategyInitiativeCommunity) (err error) {
	emptyFields, ok := strategyInitiativeCommunity.CollectEmptyFields()
	if ok {
		strategyInitiativeCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())
	strategyInitiativeCommunity.CreatedAt = strategyInitiativeCommunity.ModifiedAt
	err = d.Dynamo.PutTableEntry(strategyInitiativeCommunity, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the StrategyInitiativeCommunity.
func (d DAOImpl) CreateUnsafe(strategyInitiativeCommunity StrategyInitiativeCommunity) {
	err := d.Create(strategyInitiativeCommunity)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create id==%s in %s\n", strategyInitiativeCommunity.ID, d.Name))
}


// Read reads StrategyInitiativeCommunity
func (d DAOImpl) Read(id string) (out StrategyInitiativeCommunity, err error) {
	var outs []StrategyInitiativeCommunity
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, d.Name)
	}
	return
}


// ReadUnsafe reads the StrategyInitiativeCommunity. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) StrategyInitiativeCommunity {
	out, err := d.Read(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading id==%s in %s\n", id, d.Name))
	return out
}


// ReadOrEmpty reads StrategyInitiativeCommunity
func (d DAOImpl) ReadOrEmpty(id string) (out []StrategyInitiativeCommunity, err error) {
	var outOrEmpty StrategyInitiativeCommunity
	ids := idParams(id)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.ID == id {
		out = append(out, outOrEmpty)
	} else if err != nil && strings.HasPrefix(err.Error(), "[NOT FOUND]") {
		err = nil // expected not-found error	
	}
	err = errors.Wrapf(err, "StrategyInitiativeCommunity DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the StrategyInitiativeCommunity. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []StrategyInitiativeCommunity {
	out, err := d.ReadOrEmpty(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading id==%s in %s\n", id, d.Name))
	return out
}


// CreateOrUpdate saves the StrategyInitiativeCommunity regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(strategyInitiativeCommunity StrategyInitiativeCommunity) (err error) {
	strategyInitiativeCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())
	if strategyInitiativeCommunity.CreatedAt == "" { strategyInitiativeCommunity.CreatedAt = strategyInitiativeCommunity.ModifiedAt }
	
	var olds []StrategyInitiativeCommunity
	olds, err = d.ReadOrEmpty(strategyInitiativeCommunity.ID)
	err = errors.Wrapf(err, "StrategyInitiativeCommunity DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", strategyInitiativeCommunity.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(strategyInitiativeCommunity)
			err = errors.Wrapf(err, "StrategyInitiativeCommunity DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			emptyFields, ok := strategyInitiativeCommunity.CollectEmptyFields()
			if ok {
				old := olds[0]
				strategyInitiativeCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())

				key := idParams(old.ID)
				expr, exprAttributes, names := updateExpression(strategyInitiativeCommunity, old)
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
				err = errors.Wrapf(err, "StrategyInitiativeCommunity DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the StrategyInitiativeCommunity regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(strategyInitiativeCommunity StrategyInitiativeCommunity) {
	err := d.CreateOrUpdate(strategyInitiativeCommunity)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", strategyInitiativeCommunity, d.Name))
}


// Delete removes StrategyInitiativeCommunity from db
func (d DAOImpl)Delete(id string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(id))
}


// DeleteUnsafe deletes StrategyInitiativeCommunity and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(id string) {
	err := d.Delete(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete id==%s in %s\n", id, d.Name))
}


func (d DAOImpl)ReadByPlatformID(platformID common.PlatformID) (out []StrategyInitiativeCommunity, err error) {
	var instances []StrategyInitiativeCommunity
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


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID common.PlatformID) (out []StrategyInitiativeCommunity) {
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
func allParams(strategyInitiativeCommunity StrategyInitiativeCommunity, old StrategyInitiativeCommunity) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if strategyInitiativeCommunity.ID != old.ID { params[":a0"] = common.DynS(strategyInitiativeCommunity.ID) }
	if strategyInitiativeCommunity.PlatformID != old.PlatformID { params[":a1"] = common.DynS(string(strategyInitiativeCommunity.PlatformID)) }
	if strategyInitiativeCommunity.Name != old.Name { params[":a2"] = common.DynS(strategyInitiativeCommunity.Name) }
	if strategyInitiativeCommunity.Description != old.Description { params[":a3"] = common.DynS(strategyInitiativeCommunity.Description) }
	if strategyInitiativeCommunity.Advocate != old.Advocate { params[":a4"] = common.DynS(strategyInitiativeCommunity.Advocate) }
	if strategyInitiativeCommunity.CapabilityCommunityID != old.CapabilityCommunityID { params[":a5"] = common.DynS(strategyInitiativeCommunity.CapabilityCommunityID) }
	if strategyInitiativeCommunity.CreatedBy != old.CreatedBy { params[":a6"] = common.DynS(strategyInitiativeCommunity.CreatedBy) }
	if strategyInitiativeCommunity.CreatedAt != old.CreatedAt { params[":a7"] = common.DynS(strategyInitiativeCommunity.CreatedAt) }
	if strategyInitiativeCommunity.ModifiedAt != old.ModifiedAt { params[":a8"] = common.DynS(strategyInitiativeCommunity.ModifiedAt) }
	return
}
func updateExpression(strategyInitiativeCommunity StrategyInitiativeCommunity, old StrategyInitiativeCommunity) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if strategyInitiativeCommunity.ID != old.ID { updateParts = append(updateParts, "id = :a0"); params[":a0"] = common.DynS(strategyInitiativeCommunity.ID);  }
	if strategyInitiativeCommunity.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a1"); params[":a1"] = common.DynS(string(strategyInitiativeCommunity.PlatformID));  }
	if strategyInitiativeCommunity.Name != old.Name { updateParts = append(updateParts, "#name = :a2"); params[":a2"] = common.DynS(strategyInitiativeCommunity.Name); fldName := "name"; names["#name"] = &fldName }
	if strategyInitiativeCommunity.Description != old.Description { updateParts = append(updateParts, "description = :a3"); params[":a3"] = common.DynS(strategyInitiativeCommunity.Description);  }
	if strategyInitiativeCommunity.Advocate != old.Advocate { updateParts = append(updateParts, "advocate = :a4"); params[":a4"] = common.DynS(strategyInitiativeCommunity.Advocate);  }
	if strategyInitiativeCommunity.CapabilityCommunityID != old.CapabilityCommunityID { updateParts = append(updateParts, "capability_community_id = :a5"); params[":a5"] = common.DynS(strategyInitiativeCommunity.CapabilityCommunityID);  }
	if strategyInitiativeCommunity.CreatedBy != old.CreatedBy { updateParts = append(updateParts, "created_by = :a6"); params[":a6"] = common.DynS(strategyInitiativeCommunity.CreatedBy);  }
	if strategyInitiativeCommunity.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a7"); params[":a7"] = common.DynS(strategyInitiativeCommunity.CreatedAt);  }
	if strategyInitiativeCommunity.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a8"); params[":a8"] = common.DynS(strategyInitiativeCommunity.ModifiedAt);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
