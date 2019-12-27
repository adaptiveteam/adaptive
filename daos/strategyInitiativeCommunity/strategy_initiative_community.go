package strategyInitiativeCommunity
// This file has been automatically generated by `adaptive-platform/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/adaptiveteam/adaptive-utils-go/models"
	"time"
	awsutils "github.com/adaptiveteam/aws-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

type StrategyInitiativeCommunity struct  {
	ID string `json:"id"`
	PlatformID models.PlatformID `json:"platform_id"`
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
	ReadByIDPlatformID(id string, platformID models.PlatformID) (strategyInitiativeCommunity []StrategyInitiativeCommunity, err error)
	ReadByIDPlatformIDUnsafe(id string, platformID models.PlatformID) (strategyInitiativeCommunity []StrategyInitiativeCommunity)
	ReadByPlatformID(platformID models.PlatformID) (strategyInitiativeCommunity []StrategyInitiativeCommunity, err error)
	ReadByPlatformIDUnsafe(platformID models.PlatformID) (strategyInitiativeCommunity []StrategyInitiativeCommunity)
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
func (d DAOImpl) Create(strategyInitiativeCommunity StrategyInitiativeCommunity) error {
	emptyFields, ok := strategyInitiativeCommunity.CollectEmptyFields()
	if !ok {return fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)}
	strategyInitiativeCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())
	strategyInitiativeCommunity.CreatedAt = strategyInitiativeCommunity.ModifiedAt
	return d.Dynamo.PutTableEntry(strategyInitiativeCommunity, d.Name)
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
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(strategyInitiativeCommunity)
			err = errors.Wrapf(err, "StrategyInitiativeCommunity DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			old := olds[0]
			strategyInitiativeCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())
ids := idParams(old.ID)
			err = d.Dynamo.UpdateTableEntry(
				allParams(strategyInitiativeCommunity, old),
				ids,
				updateExpression(strategyInitiativeCommunity, old),
				d.Name,
			)
			err = errors.Wrapf(err, "StrategyInitiativeCommunity DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s", ids, d.Name)
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


func (d DAOImpl)ReadByIDPlatformID(id string, platformID models.PlatformID) (out []StrategyInitiativeCommunity, err error) {
	var instances []StrategyInitiativeCommunity
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


func (d DAOImpl)ReadByIDPlatformIDUnsafe(id string, platformID models.PlatformID) (out []StrategyInitiativeCommunity) {
	out, err := d.ReadByIDPlatformID(id, platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query IDPlatformIDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByPlatformID(platformID models.PlatformID) (out []StrategyInitiativeCommunity, err error) {
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


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID models.PlatformID) (out []StrategyInitiativeCommunity) {
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
		if strategyInitiativeCommunity.ID != old.ID { params["a0"] = common.DynS(strategyInitiativeCommunity.ID) }
		if strategyInitiativeCommunity.PlatformID != old.PlatformID { params["a1"] = common.DynS(string(strategyInitiativeCommunity.PlatformID)) }
		if strategyInitiativeCommunity.Name != old.Name { params["a2"] = common.DynS(strategyInitiativeCommunity.Name) }
		if strategyInitiativeCommunity.Description != old.Description { params["a3"] = common.DynS(strategyInitiativeCommunity.Description) }
		if strategyInitiativeCommunity.Advocate != old.Advocate { params["a4"] = common.DynS(strategyInitiativeCommunity.Advocate) }
		if strategyInitiativeCommunity.CapabilityCommunityID != old.CapabilityCommunityID { params["a5"] = common.DynS(strategyInitiativeCommunity.CapabilityCommunityID) }
		if strategyInitiativeCommunity.CreatedBy != old.CreatedBy { params["a6"] = common.DynS(strategyInitiativeCommunity.CreatedBy) }
		if strategyInitiativeCommunity.CreatedAt != old.CreatedAt { params["a7"] = common.DynS(strategyInitiativeCommunity.CreatedAt) }
		if strategyInitiativeCommunity.ModifiedAt != old.ModifiedAt { params["a8"] = common.DynS(strategyInitiativeCommunity.ModifiedAt) }
	return
}
func updateExpression(strategyInitiativeCommunity StrategyInitiativeCommunity, old StrategyInitiativeCommunity) string {
	var updateParts []string
	
		
			
		if strategyInitiativeCommunity.ID != old.ID { updateParts = append(updateParts, "id = :a0") }
		if strategyInitiativeCommunity.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a1") }
		if strategyInitiativeCommunity.Name != old.Name { updateParts = append(updateParts, "#name = :a2") }
		if strategyInitiativeCommunity.Description != old.Description { updateParts = append(updateParts, "description = :a3") }
		if strategyInitiativeCommunity.Advocate != old.Advocate { updateParts = append(updateParts, "advocate = :a4") }
		if strategyInitiativeCommunity.CapabilityCommunityID != old.CapabilityCommunityID { updateParts = append(updateParts, "capability_community_id = :a5") }
		if strategyInitiativeCommunity.CreatedBy != old.CreatedBy { updateParts = append(updateParts, "created_by = :a6") }
		if strategyInitiativeCommunity.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a7") }
		if strategyInitiativeCommunity.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a8") }
	return strings.Join(updateParts, " and ")
}
