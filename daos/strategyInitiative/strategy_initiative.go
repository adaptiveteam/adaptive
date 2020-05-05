package strategyInitiative
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

type StrategyInitiative struct  {
	PlatformID common.PlatformID `json:"platform_id"`
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	DefinitionOfVictory string `json:"definition_of_victory"`
	Advocate string `json:"advocate"`
	InitiativeCommunityID string `json:"initiative_community_id"`
	Budget string `json:"budget"`
	ExpectedEndDate string `json:"expected_end_date"`
	CapabilityObjective string `json:"capability_objective"`
	CreatedBy string `json:"created_by,omitempty"`
	ModifiedBy string `json:"modified_by,omitempty"`
	// Automatically maintained field
	CreatedAt string `json:"created_at"`
	// Automatically maintained field
	ModifiedAt string `json:"modified_at,omitempty"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (strategyInitiative StrategyInitiative)CollectEmptyFields() (emptyFields []string, ok bool) {
	if strategyInitiative.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if strategyInitiative.ID == "" { emptyFields = append(emptyFields, "ID")}
	if strategyInitiative.Name == "" { emptyFields = append(emptyFields, "Name")}
	if strategyInitiative.Description == "" { emptyFields = append(emptyFields, "Description")}
	if strategyInitiative.DefinitionOfVictory == "" { emptyFields = append(emptyFields, "DefinitionOfVictory")}
	if strategyInitiative.Advocate == "" { emptyFields = append(emptyFields, "Advocate")}
	if strategyInitiative.InitiativeCommunityID == "" { emptyFields = append(emptyFields, "InitiativeCommunityID")}
	if strategyInitiative.Budget == "" { emptyFields = append(emptyFields, "Budget")}
	if strategyInitiative.ExpectedEndDate == "" { emptyFields = append(emptyFields, "ExpectedEndDate")}
	if strategyInitiative.CapabilityObjective == "" { emptyFields = append(emptyFields, "CapabilityObjective")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (strategyInitiative StrategyInitiative) ToJSON() (string, error) {
	b, err := json.Marshal(strategyInitiative)
	return string(b), err
}

type DAO interface {
	Create(strategyInitiative StrategyInitiative) error
	CreateUnsafe(strategyInitiative StrategyInitiative)
	Read(platformID common.PlatformID, id string) (strategyInitiative StrategyInitiative, err error)
	ReadUnsafe(platformID common.PlatformID, id string) (strategyInitiative StrategyInitiative)
	ReadOrEmpty(platformID common.PlatformID, id string) (strategyInitiative []StrategyInitiative, err error)
	ReadOrEmptyUnsafe(platformID common.PlatformID, id string) (strategyInitiative []StrategyInitiative)
	CreateOrUpdate(strategyInitiative StrategyInitiative) error
	CreateOrUpdateUnsafe(strategyInitiative StrategyInitiative)
	Delete(platformID common.PlatformID, id string) error
	DeleteUnsafe(platformID common.PlatformID, id string)
	ReadByPlatformID(platformID common.PlatformID) (strategyInitiative []StrategyInitiative, err error)
	ReadByPlatformIDUnsafe(platformID common.PlatformID) (strategyInitiative []StrategyInitiative)
	ReadByInitiativeCommunityID(initiativeCommunityID string) (strategyInitiative []StrategyInitiative, err error)
	ReadByInitiativeCommunityIDUnsafe(initiativeCommunityID string) (strategyInitiative []StrategyInitiative)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	ConnGen   common.DynamoDBConnectionGen
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic(errors.New("Cannot create StrategyInitiative.DAO without clientID")) }
	return DAOImpl{
		ConnGen:   common.DynamoDBConnectionGen{
			Dynamo: dynamo, 
			TableNamePrefix: clientID,
		},
	}
}
// TableNameSuffixVar is a global variable that contains table name suffix.
// After renaming all tables this may be made `const`.
var TableNameSuffixVar = "_strategy_initiative"

// TableName concatenates table name prefix and suffix and returns table name
func TableName(prefix string) string {
	return prefix + TableNameSuffixVar
}

// Create saves the StrategyInitiative.
func (d DAOImpl) Create(strategyInitiative StrategyInitiative) (err error) {
	emptyFields, ok := strategyInitiative.CollectEmptyFields()
	if ok {
		strategyInitiative.ModifiedAt = core.CurrentRFCTimestamp()
	strategyInitiative.CreatedAt = strategyInitiative.ModifiedAt
	err = d.ConnGen.Dynamo.PutTableEntry(strategyInitiative, TableName(d.ConnGen.TableNamePrefix))
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the StrategyInitiative.
func (d DAOImpl) CreateUnsafe(strategyInitiative StrategyInitiative) {
	err2 := d.Create(strategyInitiative)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not create platformID==%s, id==%s in %s\n", strategyInitiative.PlatformID, strategyInitiative.ID, TableName(d.ConnGen.TableNamePrefix)))
}


// Read reads StrategyInitiative
func (d DAOImpl) Read(platformID common.PlatformID, id string) (out StrategyInitiative, err error) {
	var outs []StrategyInitiative
	outs, err = d.ReadOrEmpty(platformID, id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found platformID==%s, id==%s in %s\n", platformID, id, TableName(d.ConnGen.TableNamePrefix))
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the StrategyInitiative. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(platformID common.PlatformID, id string) StrategyInitiative {
	out, err2 := d.Read(platformID, id)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error reading platformID==%s, id==%s in %s\n", platformID, id, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// ReadOrEmpty reads StrategyInitiative
func (d DAOImpl) ReadOrEmpty(platformID common.PlatformID, id string) (out []StrategyInitiative, err error) {
	var outOrEmpty StrategyInitiative
	ids := idParams(platformID, id)
	var found bool
	found, err = d.ConnGen.Dynamo.GetItemOrEmptyFromTable(TableName(d.ConnGen.TableNamePrefix), ids, &outOrEmpty)
	if found {
		if outOrEmpty.PlatformID == platformID && outOrEmpty.ID == id {
			out = append(out, outOrEmpty)
		} else {
			err = fmt.Errorf("Requested ids: platformID==%s, id==%s are different from the found ones: platformID==%s, id==%s", platformID, id, outOrEmpty.PlatformID, outOrEmpty.ID) // unexpected error: found ids != ids
		}
	}
	err = errors.Wrapf(err, "StrategyInitiative DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(d.ConnGen.TableNamePrefix))
	return
}


// ReadOrEmptyUnsafe reads the StrategyInitiative. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(platformID common.PlatformID, id string) []StrategyInitiative {
	out, err2 := d.ReadOrEmpty(platformID, id)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error while reading platformID==%s, id==%s in %s\n", platformID, id, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// CreateOrUpdate saves the StrategyInitiative regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(strategyInitiative StrategyInitiative) (err error) {
	strategyInitiative.ModifiedAt = core.CurrentRFCTimestamp()
	if strategyInitiative.CreatedAt == "" { strategyInitiative.CreatedAt = strategyInitiative.ModifiedAt }
	
	var olds []StrategyInitiative
	olds, err = d.ReadOrEmpty(strategyInitiative.PlatformID, strategyInitiative.ID)
	err = errors.Wrapf(err, "StrategyInitiative DAO.CreateOrUpdate(id = platformID==%s, id==%s) couldn't ReadOrEmpty", strategyInitiative.PlatformID, strategyInitiative.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(strategyInitiative)
			err = errors.Wrapf(err, "StrategyInitiative DAO.CreateOrUpdate couldn't Create in table %s", TableName(d.ConnGen.TableNamePrefix))
		} else {
			emptyFields, ok := strategyInitiative.CollectEmptyFields()
			if ok {
				old := olds[0]
				strategyInitiative.CreatedAt  = old.CreatedAt
				strategyInitiative.ModifiedAt = core.CurrentRFCTimestamp()
				key := idParams(old.PlatformID, old.ID)
				expr, exprAttributes, names := updateExpression(strategyInitiative, old)
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
				err = errors.Wrapf(err, "StrategyInitiative DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(d.ConnGen.TableNamePrefix), expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the StrategyInitiative regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(strategyInitiative StrategyInitiative) {
	err2 := d.CreateOrUpdate(strategyInitiative)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("could not create or update %v in %s\n", strategyInitiative, TableName(d.ConnGen.TableNamePrefix)))
}


// Delete removes StrategyInitiative from db
func (d DAOImpl)Delete(platformID common.PlatformID, id string) error {
	return d.ConnGen.Dynamo.DeleteEntry(TableName(d.ConnGen.TableNamePrefix), idParams(platformID, id))
}


// DeleteUnsafe deletes StrategyInitiative and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(platformID common.PlatformID, id string) {
	err2 := d.Delete(platformID, id)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not delete platformID==%s, id==%s in %s\n", platformID, id, TableName(d.ConnGen.TableNamePrefix)))
}


func (d DAOImpl)ReadByPlatformID(platformID common.PlatformID) (out []StrategyInitiative, err error) {
	var instances []StrategyInitiative
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
		IndexName: "PlatformIDIndex",
		Condition: "platform_id = :a0",
		Attributes: map[string]interface{}{
			":a0": platformID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID common.PlatformID) (out []StrategyInitiative) {
	out, err2 := d.ReadByPlatformID(platformID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
	return
}


func (d DAOImpl)ReadByInitiativeCommunityID(initiativeCommunityID string) (out []StrategyInitiative, err error) {
	var instances []StrategyInitiative
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
		IndexName: "InitiativeCommunityIDIndex",
		Condition: "initiative_community_id = :a0",
		Attributes: map[string]interface{}{
			":a0": initiativeCommunityID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByInitiativeCommunityIDUnsafe(initiativeCommunityID string) (out []StrategyInitiative) {
	out, err2 := d.ReadByInitiativeCommunityID(initiativeCommunityID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query InitiativeCommunityIDIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
	return
}

func idParams(platformID common.PlatformID, id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"platform_id": common.DynS(string(platformID)),
		"id": common.DynS(id),
	}
	return params
}
func allParams(strategyInitiative StrategyInitiative, old StrategyInitiative) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if strategyInitiative.PlatformID != old.PlatformID { params[":a0"] = common.DynS(string(strategyInitiative.PlatformID)) }
	if strategyInitiative.ID != old.ID { params[":a1"] = common.DynS(strategyInitiative.ID) }
	if strategyInitiative.Name != old.Name { params[":a2"] = common.DynS(strategyInitiative.Name) }
	if strategyInitiative.Description != old.Description { params[":a3"] = common.DynS(strategyInitiative.Description) }
	if strategyInitiative.DefinitionOfVictory != old.DefinitionOfVictory { params[":a4"] = common.DynS(strategyInitiative.DefinitionOfVictory) }
	if strategyInitiative.Advocate != old.Advocate { params[":a5"] = common.DynS(strategyInitiative.Advocate) }
	if strategyInitiative.InitiativeCommunityID != old.InitiativeCommunityID { params[":a6"] = common.DynS(strategyInitiative.InitiativeCommunityID) }
	if strategyInitiative.Budget != old.Budget { params[":a7"] = common.DynS(strategyInitiative.Budget) }
	if strategyInitiative.ExpectedEndDate != old.ExpectedEndDate { params[":a8"] = common.DynS(strategyInitiative.ExpectedEndDate) }
	if strategyInitiative.CapabilityObjective != old.CapabilityObjective { params[":a9"] = common.DynS(strategyInitiative.CapabilityObjective) }
	if strategyInitiative.CreatedBy != old.CreatedBy { params[":a10"] = common.DynS(strategyInitiative.CreatedBy) }
	if strategyInitiative.ModifiedBy != old.ModifiedBy { params[":a11"] = common.DynS(strategyInitiative.ModifiedBy) }
	if strategyInitiative.CreatedAt != old.CreatedAt { params[":a12"] = common.DynS(strategyInitiative.CreatedAt) }
	if strategyInitiative.ModifiedAt != old.ModifiedAt { params[":a13"] = common.DynS(strategyInitiative.ModifiedAt) }
	return
}
func updateExpression(strategyInitiative StrategyInitiative, old StrategyInitiative) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if strategyInitiative.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a0"); params[":a0"] = common.DynS(string(strategyInitiative.PlatformID));  }
	if strategyInitiative.ID != old.ID { updateParts = append(updateParts, "id = :a1"); params[":a1"] = common.DynS(strategyInitiative.ID);  }
	if strategyInitiative.Name != old.Name { updateParts = append(updateParts, "#name = :a2"); params[":a2"] = common.DynS(strategyInitiative.Name); fldName := "name"; names["#name"] = &fldName }
	if strategyInitiative.Description != old.Description { updateParts = append(updateParts, "description = :a3"); params[":a3"] = common.DynS(strategyInitiative.Description);  }
	if strategyInitiative.DefinitionOfVictory != old.DefinitionOfVictory { updateParts = append(updateParts, "definition_of_victory = :a4"); params[":a4"] = common.DynS(strategyInitiative.DefinitionOfVictory);  }
	if strategyInitiative.Advocate != old.Advocate { updateParts = append(updateParts, "advocate = :a5"); params[":a5"] = common.DynS(strategyInitiative.Advocate);  }
	if strategyInitiative.InitiativeCommunityID != old.InitiativeCommunityID { updateParts = append(updateParts, "initiative_community_id = :a6"); params[":a6"] = common.DynS(strategyInitiative.InitiativeCommunityID);  }
	if strategyInitiative.Budget != old.Budget { updateParts = append(updateParts, "budget = :a7"); params[":a7"] = common.DynS(strategyInitiative.Budget);  }
	if strategyInitiative.ExpectedEndDate != old.ExpectedEndDate { updateParts = append(updateParts, "expected_end_date = :a8"); params[":a8"] = common.DynS(strategyInitiative.ExpectedEndDate);  }
	if strategyInitiative.CapabilityObjective != old.CapabilityObjective { updateParts = append(updateParts, "capability_objective = :a9"); params[":a9"] = common.DynS(strategyInitiative.CapabilityObjective);  }
	if strategyInitiative.CreatedBy != old.CreatedBy { updateParts = append(updateParts, "created_by = :a10"); params[":a10"] = common.DynS(strategyInitiative.CreatedBy);  }
	if strategyInitiative.ModifiedBy != old.ModifiedBy { updateParts = append(updateParts, "modified_by = :a11"); params[":a11"] = common.DynS(strategyInitiative.ModifiedBy);  }
	if strategyInitiative.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a12"); params[":a12"] = common.DynS(strategyInitiative.CreatedAt);  }
	if strategyInitiative.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a13"); params[":a13"] = common.DynS(strategyInitiative.ModifiedAt);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
