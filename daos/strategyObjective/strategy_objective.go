package strategyObjective
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
	"encoding/json"
	"strings"
)

type StrategyObjectiveType string

type StrategyObjective struct  {
	// range key
	PlatformID common.PlatformID `json:"platform_id"`
	// hash
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	AsMeasuredBy string `json:"as_measured_by"`
	Targets string `json:"targets"`
	ObjectiveType StrategyObjectiveType `json:"objective_type"`
	Advocate string `json:"advocate"`
	// community id not require d for customer/financial objectives
	CapabilityCommunityIDs []string `json:"capability_community_id"`
	ExpectedEndDate string `json:"expected_end_date"`
	CreatedBy string `json:"created_by"`
	// Automatically maintained field
	CreatedAt string `json:"created_at"`
	// Automatically maintained field
	ModifiedAt string `json:"modified_at,omitempty"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (strategyObjective StrategyObjective)CollectEmptyFields() (emptyFields []string, ok bool) {
	if strategyObjective.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if strategyObjective.ID == "" { emptyFields = append(emptyFields, "ID")}
	if strategyObjective.Name == "" { emptyFields = append(emptyFields, "Name")}
	if strategyObjective.Description == "" { emptyFields = append(emptyFields, "Description")}
	if strategyObjective.AsMeasuredBy == "" { emptyFields = append(emptyFields, "AsMeasuredBy")}
	if strategyObjective.Targets == "" { emptyFields = append(emptyFields, "Targets")}
	if strategyObjective.Advocate == "" { emptyFields = append(emptyFields, "Advocate")}
	if strategyObjective.CapabilityCommunityIDs == nil { emptyFields = append(emptyFields, "CapabilityCommunityIDs")}
	if strategyObjective.ExpectedEndDate == "" { emptyFields = append(emptyFields, "ExpectedEndDate")}
	if strategyObjective.CreatedBy == "" { emptyFields = append(emptyFields, "CreatedBy")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (strategyObjective StrategyObjective) ToJSON() (string, error) {
	b, err := json.Marshal(strategyObjective)
	return string(b), err
}

type DAO interface {
	Create(strategyObjective StrategyObjective) error
	CreateUnsafe(strategyObjective StrategyObjective)
	Read(platformID common.PlatformID, id string) (strategyObjective StrategyObjective, err error)
	ReadUnsafe(platformID common.PlatformID, id string) (strategyObjective StrategyObjective)
	ReadOrEmpty(platformID common.PlatformID, id string) (strategyObjective []StrategyObjective, err error)
	ReadOrEmptyUnsafe(platformID common.PlatformID, id string) (strategyObjective []StrategyObjective)
	CreateOrUpdate(strategyObjective StrategyObjective) error
	CreateOrUpdateUnsafe(strategyObjective StrategyObjective)
	Delete(platformID common.PlatformID, id string) error
	DeleteUnsafe(platformID common.PlatformID, id string)
	ReadByPlatformID(platformID common.PlatformID) (strategyObjective []StrategyObjective, err error)
	ReadByPlatformIDUnsafe(platformID common.PlatformID) (strategyObjective []StrategyObjective)
	ReadByCapabilityCommunityID(capabilityCommunityIDs []string) (strategyObjective []StrategyObjective, err error)
	ReadByCapabilityCommunityIDUnsafe(capabilityCommunityIDs []string) (strategyObjective []StrategyObjective)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	Name      string                  `json:"name"`
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic(errors.New("Cannot create StrategyObjective.DAO without clientID")) }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: TableName(clientID),
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic(errors.New("Cannot create StrategyObjective.DAO without tableName")) }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}
func TableName(clientID string) string {
	return clientID + "_strategy_objective"
}

// Create saves the StrategyObjective.
func (d DAOImpl) Create(strategyObjective StrategyObjective) (err error) {
	emptyFields, ok := strategyObjective.CollectEmptyFields()
	if ok {
		strategyObjective.ModifiedAt = core.CurrentRFCTimestamp()
	strategyObjective.CreatedAt = strategyObjective.ModifiedAt
	err = d.Dynamo.PutTableEntry(strategyObjective, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the StrategyObjective.
func (d DAOImpl) CreateUnsafe(strategyObjective StrategyObjective) {
	err2 := d.Create(strategyObjective)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not create platformID==%s, id==%s in %s\n", strategyObjective.PlatformID, strategyObjective.ID, d.Name))
}


// Read reads StrategyObjective
func (d DAOImpl) Read(platformID common.PlatformID, id string) (out StrategyObjective, err error) {
	var outs []StrategyObjective
	outs, err = d.ReadOrEmpty(platformID, id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found platformID==%s, id==%s in %s\n", platformID, id, d.Name)
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the StrategyObjective. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(platformID common.PlatformID, id string) StrategyObjective {
	out, err2 := d.Read(platformID, id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Error reading platformID==%s, id==%s in %s\n", platformID, id, d.Name))
	return out
}


// ReadOrEmpty reads StrategyObjective
func (d DAOImpl) ReadOrEmpty(platformID common.PlatformID, id string) (out []StrategyObjective, err error) {
	var outOrEmpty StrategyObjective
	ids := idParams(platformID, id)
	var found bool
	found, err = d.Dynamo.GetItemOrEmptyFromTable(d.Name, ids, &outOrEmpty)
	if found {
		if outOrEmpty.PlatformID == platformID && outOrEmpty.ID == id {
			out = append(out, outOrEmpty)
		} else {
			err = fmt.Errorf("Requested ids: platformID==%s, id==%s are different from the found ones: platformID==%s, id==%s", platformID, id, outOrEmpty.PlatformID, outOrEmpty.ID) // unexpected error: found ids != ids
		}
	}
	err = errors.Wrapf(err, "StrategyObjective DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the StrategyObjective. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(platformID common.PlatformID, id string) []StrategyObjective {
	out, err2 := d.ReadOrEmpty(platformID, id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Error while reading platformID==%s, id==%s in %s\n", platformID, id, d.Name))
	return out
}


// CreateOrUpdate saves the StrategyObjective regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(strategyObjective StrategyObjective) (err error) {
	strategyObjective.ModifiedAt = core.CurrentRFCTimestamp()
	if strategyObjective.CreatedAt == "" { strategyObjective.CreatedAt = strategyObjective.ModifiedAt }
	
	var olds []StrategyObjective
	olds, err = d.ReadOrEmpty(strategyObjective.PlatformID, strategyObjective.ID)
	err = errors.Wrapf(err, "StrategyObjective DAO.CreateOrUpdate(id = platformID==%s, id==%s) couldn't ReadOrEmpty", strategyObjective.PlatformID, strategyObjective.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(strategyObjective)
			err = errors.Wrapf(err, "StrategyObjective DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
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
					TableName:                 aws.String(d.Name),
					Key:                       key,
					ReturnValues:              aws.String("UPDATED_NEW"),
					UpdateExpression:          aws.String(expr),
				}
				if names != nil { input.ExpressionAttributeNames = *names } // workaround for a pointer to an empty slice
				if  len(exprAttributes) > 0 { // if there some changes
					err = d.Dynamo.UpdateItemInternal(input)
				} else {
					// WARN: no changes.
				}
				err = errors.Wrapf(err, "StrategyObjective DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the StrategyObjective regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(strategyObjective StrategyObjective) {
	err2 := d.CreateOrUpdate(strategyObjective)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", strategyObjective, d.Name))
}


// Delete removes StrategyObjective from db
func (d DAOImpl)Delete(platformID common.PlatformID, id string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(platformID, id))
}


// DeleteUnsafe deletes StrategyObjective and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(platformID common.PlatformID, id string) {
	err2 := d.Delete(platformID, id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not delete platformID==%s, id==%s in %s\n", platformID, id, d.Name))
}


func (d DAOImpl)ReadByPlatformID(platformID common.PlatformID) (out []StrategyObjective, err error) {
	var instances []StrategyObjective
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


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID common.PlatformID) (out []StrategyObjective) {
	out, err2 := d.ReadByPlatformID(platformID)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByCapabilityCommunityID(capabilityCommunityIDs []string) (out []StrategyObjective, err error) {
	var instances []StrategyObjective
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "CapabilityCommunityIDIndex",
		Condition: "capability_community_id = :a0",
		Attributes: map[string]interface{}{
			":a0": capabilityCommunityIDs,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByCapabilityCommunityIDUnsafe(capabilityCommunityIDs []string) (out []StrategyObjective) {
	out, err2 := d.ReadByCapabilityCommunityID(capabilityCommunityIDs)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not query CapabilityCommunityIDIndex on %s table\n", d.Name))
	return
}

func idParams(platformID common.PlatformID, id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"platform_id": common.DynS(string(platformID)),
		"id": common.DynS(id),
	}
	return params
}
func allParams(strategyObjective StrategyObjective, old StrategyObjective) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if strategyObjective.PlatformID != old.PlatformID { params[":a0"] = common.DynS(string(strategyObjective.PlatformID)) }
	if strategyObjective.ID != old.ID { params[":a1"] = common.DynS(strategyObjective.ID) }
	if strategyObjective.Name != old.Name { params[":a2"] = common.DynS(strategyObjective.Name) }
	if strategyObjective.Description != old.Description { params[":a3"] = common.DynS(strategyObjective.Description) }
	if strategyObjective.AsMeasuredBy != old.AsMeasuredBy { params[":a4"] = common.DynS(strategyObjective.AsMeasuredBy) }
	if strategyObjective.Targets != old.Targets { params[":a5"] = common.DynS(strategyObjective.Targets) }
	if strategyObjective.ObjectiveType != old.ObjectiveType { params[":a6"] = common.DynS(string(strategyObjective.ObjectiveType)) }
	if strategyObjective.Advocate != old.Advocate { params[":a7"] = common.DynS(strategyObjective.Advocate) }
	if !common.StringArraysEqual(strategyObjective.CapabilityCommunityIDs, old.CapabilityCommunityIDs) { params[":a8"] = common.DynSS(strategyObjective.CapabilityCommunityIDs) }
	if strategyObjective.ExpectedEndDate != old.ExpectedEndDate { params[":a9"] = common.DynS(strategyObjective.ExpectedEndDate) }
	if strategyObjective.CreatedBy != old.CreatedBy { params[":a10"] = common.DynS(strategyObjective.CreatedBy) }
	if strategyObjective.CreatedAt != old.CreatedAt { params[":a11"] = common.DynS(strategyObjective.CreatedAt) }
	if strategyObjective.ModifiedAt != old.ModifiedAt { params[":a12"] = common.DynS(strategyObjective.ModifiedAt) }
	return
}
func updateExpression(strategyObjective StrategyObjective, old StrategyObjective) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if strategyObjective.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a0"); params[":a0"] = common.DynS(string(strategyObjective.PlatformID));  }
	if strategyObjective.ID != old.ID { updateParts = append(updateParts, "id = :a1"); params[":a1"] = common.DynS(strategyObjective.ID);  }
	if strategyObjective.Name != old.Name { updateParts = append(updateParts, "#name = :a2"); params[":a2"] = common.DynS(strategyObjective.Name); fldName := "name"; names["#name"] = &fldName }
	if strategyObjective.Description != old.Description { updateParts = append(updateParts, "description = :a3"); params[":a3"] = common.DynS(strategyObjective.Description);  }
	if strategyObjective.AsMeasuredBy != old.AsMeasuredBy { updateParts = append(updateParts, "as_measured_by = :a4"); params[":a4"] = common.DynS(strategyObjective.AsMeasuredBy);  }
	if strategyObjective.Targets != old.Targets { updateParts = append(updateParts, "targets = :a5"); params[":a5"] = common.DynS(strategyObjective.Targets);  }
	if strategyObjective.ObjectiveType != old.ObjectiveType { updateParts = append(updateParts, "objective_type = :a6"); params[":a6"] = common.DynS(string(strategyObjective.ObjectiveType));  }
	if strategyObjective.Advocate != old.Advocate { updateParts = append(updateParts, "advocate = :a7"); params[":a7"] = common.DynS(strategyObjective.Advocate);  }
	if !common.StringArraysEqual(strategyObjective.CapabilityCommunityIDs, old.CapabilityCommunityIDs) { updateParts = append(updateParts, "capability_community_id = :a8"); params[":a8"] = common.DynSS(strategyObjective.CapabilityCommunityIDs);  }
	if strategyObjective.ExpectedEndDate != old.ExpectedEndDate { updateParts = append(updateParts, "expected_end_date = :a9"); params[":a9"] = common.DynS(strategyObjective.ExpectedEndDate);  }
	if strategyObjective.CreatedBy != old.CreatedBy { updateParts = append(updateParts, "created_by = :a10"); params[":a10"] = common.DynS(strategyObjective.CreatedBy);  }
	if strategyObjective.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a11"); params[":a11"] = common.DynS(strategyObjective.CreatedAt);  }
	if strategyObjective.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a12"); params[":a12"] = common.DynS(strategyObjective.ModifiedAt);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
