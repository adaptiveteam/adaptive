package strategyInitiative
// This file has been automatically generated by `adaptive-platform/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/adaptiveteam/adaptive-utils-go/models"
	"time"
	awsutils "github.com/adaptiveteam/aws-utils-go"
	common "github.com/adaptiveteam/daos/common"
	core "github.com/adaptiveteam/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

type StrategyInitiative struct  {
	ID string `json:"id"`
	PlatformID models.PlatformID `json:"platform_id"`
	Name string `json:"name"`
	Description string `json:"description"`
	DefinitionOfVictory string `json:"definition_of_victory"`
	Advocate string `json:"advocate"`
	InitiativeCommunityID string `json:"initiative_community_id"`
	Budget string `json:"budget"`
	ExpectedEndDate string `json:"expected_end_date"`
	CapabilityObjective string `json:"capability_objective"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	// Automatically maintained field
	CreatedAt string `json:"created_at"`
	// Automatically maintained field
	ModifiedAt string `json:"modified_at"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (strategyInitiative StrategyInitiative)CollectEmptyFields() (emptyFields []string, ok bool) {
	if strategyInitiative.ID == "" { emptyFields = append(emptyFields, "ID")}
	if strategyInitiative.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if strategyInitiative.Name == "" { emptyFields = append(emptyFields, "Name")}
	if strategyInitiative.Description == "" { emptyFields = append(emptyFields, "Description")}
	if strategyInitiative.DefinitionOfVictory == "" { emptyFields = append(emptyFields, "DefinitionOfVictory")}
	if strategyInitiative.Advocate == "" { emptyFields = append(emptyFields, "Advocate")}
	if strategyInitiative.InitiativeCommunityID == "" { emptyFields = append(emptyFields, "InitiativeCommunityID")}
	if strategyInitiative.Budget == "" { emptyFields = append(emptyFields, "Budget")}
	if strategyInitiative.ExpectedEndDate == "" { emptyFields = append(emptyFields, "ExpectedEndDate")}
	if strategyInitiative.CapabilityObjective == "" { emptyFields = append(emptyFields, "CapabilityObjective")}
	if strategyInitiative.CreatedBy == "" { emptyFields = append(emptyFields, "CreatedBy")}
	if strategyInitiative.ModifiedBy == "" { emptyFields = append(emptyFields, "ModifiedBy")}
	ok = len(emptyFields) == 0
	return
}

type DAO interface {
	Create(strategyInitiative StrategyInitiative) error
	CreateUnsafe(strategyInitiative StrategyInitiative)
	Read(id string) (strategyInitiative StrategyInitiative, err error)
	ReadUnsafe(id string) (strategyInitiative StrategyInitiative)
	ReadOrEmpty(id string) (strategyInitiative []StrategyInitiative, err error)
	ReadOrEmptyUnsafe(id string) (strategyInitiative []StrategyInitiative)
	CreateOrUpdate(strategyInitiative StrategyInitiative) error
	CreateOrUpdateUnsafe(strategyInitiative StrategyInitiative)
	Delete(id string) error
	DeleteUnsafe(id string)
	ReadByIDPlatformID(id string, platformID models.PlatformID) (strategyInitiative []StrategyInitiative, err error)
	ReadByIDPlatformIDUnsafe(id string, platformID models.PlatformID) (strategyInitiative []StrategyInitiative)
	ReadByPlatformID(platformID models.PlatformID) (strategyInitiative []StrategyInitiative, err error)
	ReadByPlatformIDUnsafe(platformID models.PlatformID) (strategyInitiative []StrategyInitiative)
	ReadByInitiativeCommunityID(initiativeCommunityID string) (strategyInitiative []StrategyInitiative, err error)
	ReadByInitiativeCommunityIDUnsafe(initiativeCommunityID string) (strategyInitiative []StrategyInitiative)
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
		Name: clientID + "_strategy_initiative",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the StrategyInitiative.
func (d DAOImpl) Create(strategyInitiative StrategyInitiative) error {
	emptyFields, ok := strategyInitiative.CollectEmptyFields()
	if !ok {return fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)}
	strategyInitiative.ModifiedAt = core.TimestampLayout.Format(time.Now())
	strategyInitiative.CreatedAt = strategyInitiative.ModifiedAt
	return d.Dynamo.PutTableEntry(strategyInitiative, d.Name)
}


// CreateUnsafe saves the StrategyInitiative.
func (d DAOImpl) CreateUnsafe(strategyInitiative StrategyInitiative) {
	err := d.Create(strategyInitiative)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create id==%s in %s\n", strategyInitiative.ID, d.Name))
}


// Read reads StrategyInitiative
func (d DAOImpl) Read(id string) (out StrategyInitiative, err error) {
	var outs []StrategyInitiative
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, d.Name)
	}
	return
}


// ReadUnsafe reads the StrategyInitiative. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) StrategyInitiative {
	out, err := d.Read(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading id==%s in %s\n", id, d.Name))
	return out
}


// ReadOrEmpty reads StrategyInitiative
func (d DAOImpl) ReadOrEmpty(id string) (out []StrategyInitiative, err error) {
	var outOrEmpty StrategyInitiative
	ids := idParams(id)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.ID == id {
		out = append(out, outOrEmpty)
	}
	err = errors.Wrapf(err, "StrategyInitiative DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the StrategyInitiative. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []StrategyInitiative {
	out, err := d.ReadOrEmpty(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading id==%s in %s\n", id, d.Name))
	return out
}


// CreateOrUpdate saves the StrategyInitiative regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(strategyInitiative StrategyInitiative) (err error) {
	strategyInitiative.ModifiedAt = core.TimestampLayout.Format(time.Now())
	if strategyInitiative.CreatedAt == "" { strategyInitiative.CreatedAt = strategyInitiative.ModifiedAt }
	
	var olds []StrategyInitiative
	olds, err = d.ReadOrEmpty(strategyInitiative.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(strategyInitiative)
			err = errors.Wrapf(err, "StrategyInitiative DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			old := olds[0]
			strategyInitiative.ModifiedAt = core.TimestampLayout.Format(time.Now())
ids := idParams(old.ID)
			err = d.Dynamo.UpdateTableEntry(
				allParams(strategyInitiative, old),
				ids,
				updateExpression(strategyInitiative, old),
				d.Name,
			)
			err = errors.Wrapf(err, "StrategyInitiative DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s", ids, d.Name)
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the StrategyInitiative regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(strategyInitiative StrategyInitiative) {
	err := d.CreateOrUpdate(strategyInitiative)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", strategyInitiative, d.Name))
}


// Delete removes StrategyInitiative from db
func (d DAOImpl)Delete(id string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(id))
}


// DeleteUnsafe deletes StrategyInitiative and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(id string) {
	err := d.Delete(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete id==%s in %s\n", id, d.Name))
}


func (d DAOImpl)ReadByIDPlatformID(id string, platformID models.PlatformID) (out []StrategyInitiative, err error) {
	var instances []StrategyInitiative
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "IDPlatformIDIndex",
		Condition: "id = :a0 and platform_id = :a1",
		Attributes: map[string]interface{}{
			":a0": id,
			":a1": platformID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByIDPlatformIDUnsafe(id string, platformID models.PlatformID) (out []StrategyInitiative) {
	out, err := d.ReadByIDPlatformID(id, platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query IDPlatformIDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByPlatformID(platformID models.PlatformID) (out []StrategyInitiative, err error) {
	var instances []StrategyInitiative
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


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID models.PlatformID) (out []StrategyInitiative) {
	out, err := d.ReadByPlatformID(platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByInitiativeCommunityID(initiativeCommunityID string) (out []StrategyInitiative, err error) {
	var instances []StrategyInitiative
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
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
	out, err := d.ReadByInitiativeCommunityID(initiativeCommunityID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query InitiativeCommunityIDIndex on %s table\n", d.Name))
	return
}

func idParams(id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"id": common.DynS(id),
	}
	return params
}
func allParams(strategyInitiative StrategyInitiative, old StrategyInitiative) (params map[string]*dynamodb.AttributeValue) {
	
		params = map[string]*dynamodb.AttributeValue{}
		if strategyInitiative.ID != old.ID { params["a0"] = common.DynS(strategyInitiative.ID) }
		if strategyInitiative.PlatformID != old.PlatformID { params["a1"] = common.DynS(string(strategyInitiative.PlatformID)) }
		if strategyInitiative.Name != old.Name { params["a2"] = common.DynS(strategyInitiative.Name) }
		if strategyInitiative.Description != old.Description { params["a3"] = common.DynS(strategyInitiative.Description) }
		if strategyInitiative.DefinitionOfVictory != old.DefinitionOfVictory { params["a4"] = common.DynS(strategyInitiative.DefinitionOfVictory) }
		if strategyInitiative.Advocate != old.Advocate { params["a5"] = common.DynS(strategyInitiative.Advocate) }
		if strategyInitiative.InitiativeCommunityID != old.InitiativeCommunityID { params["a6"] = common.DynS(strategyInitiative.InitiativeCommunityID) }
		if strategyInitiative.Budget != old.Budget { params["a7"] = common.DynS(strategyInitiative.Budget) }
		if strategyInitiative.ExpectedEndDate != old.ExpectedEndDate { params["a8"] = common.DynS(strategyInitiative.ExpectedEndDate) }
		if strategyInitiative.CapabilityObjective != old.CapabilityObjective { params["a9"] = common.DynS(strategyInitiative.CapabilityObjective) }
		if strategyInitiative.CreatedBy != old.CreatedBy { params["a10"] = common.DynS(strategyInitiative.CreatedBy) }
		if strategyInitiative.ModifiedBy != old.ModifiedBy { params["a11"] = common.DynS(strategyInitiative.ModifiedBy) }
		if strategyInitiative.CreatedAt != old.CreatedAt { params["a12"] = common.DynS(strategyInitiative.CreatedAt) }
		if strategyInitiative.ModifiedAt != old.ModifiedAt { params["a13"] = common.DynS(strategyInitiative.ModifiedAt) }
	return
}
func updateExpression(strategyInitiative StrategyInitiative, old StrategyInitiative) string {
	var updateParts []string
	
		
			
		if strategyInitiative.ID != old.ID { updateParts = append(updateParts, "id = :a0") }
		if strategyInitiative.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a1") }
		if strategyInitiative.Name != old.Name { updateParts = append(updateParts, "#name = :a2") }
		if strategyInitiative.Description != old.Description { updateParts = append(updateParts, "description = :a3") }
		if strategyInitiative.DefinitionOfVictory != old.DefinitionOfVictory { updateParts = append(updateParts, "definition_of_victory = :a4") }
		if strategyInitiative.Advocate != old.Advocate { updateParts = append(updateParts, "advocate = :a5") }
		if strategyInitiative.InitiativeCommunityID != old.InitiativeCommunityID { updateParts = append(updateParts, "initiative_community_id = :a6") }
		if strategyInitiative.Budget != old.Budget { updateParts = append(updateParts, "budget = :a7") }
		if strategyInitiative.ExpectedEndDate != old.ExpectedEndDate { updateParts = append(updateParts, "expected_end_date = :a8") }
		if strategyInitiative.CapabilityObjective != old.CapabilityObjective { updateParts = append(updateParts, "capability_objective = :a9") }
		if strategyInitiative.CreatedBy != old.CreatedBy { updateParts = append(updateParts, "created_by = :a10") }
		if strategyInitiative.ModifiedBy != old.ModifiedBy { updateParts = append(updateParts, "modified_by = :a11") }
		if strategyInitiative.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a12") }
		if strategyInitiative.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a13") }
	return strings.Join(updateParts, " and ")
}
