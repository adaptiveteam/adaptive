package values

import (
	"github.com/pkg/errors"
	"fmt"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DAO - wrapper around a Dynamo DB table to work with adaptiveValues inside it
type DAO interface {
	Create(adaptiveValue models.AdaptiveValue) error
	ReadOrEmpty(adaptiveValueID string) ([]models.AdaptiveValue, error)
	Read(adaptiveValueID string) (models.AdaptiveValue, bool, error)
	ReadUnsafe(adaptiveValueID string) models.AdaptiveValue
	Update(adaptiveValue models.AdaptiveValue) error
	Delete(adaptiveValueID string) error
	Deactivate(adaptiveValueID string) error

	ForPlatformID(platformID daosCommon.PlatformID) PlatformDAO
}

// PlatformDAO is a set of utilities that work for a fixed `teamID`
type PlatformDAO interface {
	All() ([]models.AdaptiveValue, error)
	AllUnsafe() []models.AdaptiveValue
	Create(adaptiveValue models.AdaptiveValue) error
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	DNS *common.DynamoNamespace
	models.AdaptiveValuesTableSchema
}

// PlatformDAOImpl DAO that will implement PlatformDAO interface
type PlatformDAOImpl struct {
	PlatformID daosCommon.PlatformID
	DAOImpl
}

// NewDAO creates an instance of DAO that will provide access to adaptiveValues table
func NewDAO(dns *common.DynamoNamespace, table string, index string) DAO {
	if table == "" {
		panic(errors.New("Cannot create adaptiveValues DAO without table"))
	}
	if index == "" {
		panic(errors.New("Cannot create adaptiveValues DAO without index"))
	}
	return DAOImpl{DNS: dns, AdaptiveValuesTableSchema: models.AdaptiveValuesTableSchema{Name: table, PlatformIDIndex: index}}
}

// NewDAOFromSchema creates an instance of DAO that will provide access to adaptiveValues table
func NewDAOFromSchema(dns *common.DynamoNamespace, schema models.Schema) DAO {
	return DAOImpl{DNS: dns, AdaptiveValuesTableSchema: schema.AdaptiveValues}
}

// Create creates an adaptiveValue
func (d DAOImpl) Create(adaptiveValue models.AdaptiveValue) error {
	return d.DNS.Dynamo.PutTableEntry(adaptiveValue, d.Name)
}

func (d DAOImpl) ReadOrEmpty(adaptiveValueID string) (out []models.AdaptiveValue, err error) {
	var value models.AdaptiveValue
	var found bool
	value, found, err = d.Read(adaptiveValueID)
	if found {
		out = append(out, value)
	}
	return
}
// Read reads the adaptiveValue
func (d DAOImpl) Read(adaptiveValueID string) (out models.AdaptiveValue, found bool, err error) {
	params := map[string]*dynamodb.AttributeValue{
		"id": daosCommon.DynS(adaptiveValueID),
	}
	found, err = d.DNS.Dynamo.GetItemOrEmptyFromTable(d.Name, params, &out)
	return
}

// ReadUnsafe reads the adaptiveValue. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(adaptiveValueID string) models.AdaptiveValue {
	adaptiveValue, found, err2 := d.Read(adaptiveValueID)
	if !found {
		panic(fmt.Errorf("Competency %s not found", adaptiveValueID))
	}
	core.ErrorHandler(err2, d.DNS.Namespace, fmt.Sprintf("Could not find %s in %s", adaptiveValueID, d.Name))
	return adaptiveValue
}

// Delete removes the adaptiveValue
func (d DAOImpl) Delete(adaptiveValueID string) error {
	userParams := map[string]*dynamodb.AttributeValue{
		"id": daosCommon.DynS(adaptiveValueID),
	}
	return d.DNS.Dynamo.DeleteEntry(d.Name, userParams)
}

// Update updates the adaptiveValue by ID
func (d DAOImpl) Update(adaptiveValue models.AdaptiveValue) error {
	return d.DNS.Dynamo.PutTableEntry(adaptiveValue, d.Name)
}

func (d DAOImpl) Deactivate(adaptiveValueID string) (err error) {
	adaptiveValue, found, err2 := d.Read(adaptiveValueID)
	err = err2
	if err == nil {
		if found {
			adaptiveValue.DeactivatedAt = core.TimestampLayout.Format(time.Now())
			err = d.Update(adaptiveValue)
		}
	}
	return
}

// UpdateGivenAttributes updates the adaptiveValue by ID
// Deprecated: This might be used if we want to update only some of the attributes
func (d DAOImpl) UpdateGivenAttributes(adaptiveValue models.AdaptiveValue) error {
	exprAttributes := map[string]*dynamodb.AttributeValue{
		":value_name": daosCommon.DynS(adaptiveValue.Name),
		":description": daosCommon.DynS(adaptiveValue.Description),
		":value_type": daosCommon.DynS(adaptiveValue.ValueType),
		":platform_id": daosCommon.DynS(string(adaptiveValue.PlatformID)),
	}
	key := map[string]*dynamodb.AttributeValue{
		"id": daosCommon.DynS(adaptiveValue.ID),
	}
	updateExpression := "set value_name = :value_name, description = :description, value_type = :value_type, platform_id = :platform_id"
	input := dynamodb.UpdateItemInput{
		ExpressionAttributeValues: exprAttributes,
		TableName:                 aws.String(d.Name),
		Key:                       key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String(updateExpression),
	}

	return d.DNS.Dynamo.UpdateItemInternal(input)
}

// ForPlatformID creates PlatformDAO that can be used for queries with PlatformID
func (d DAOImpl) ForPlatformID(platformID daosCommon.PlatformID) PlatformDAO {
	return PlatformDAOImpl{
		PlatformID: platformID,
		DAOImpl:    d,
	}
}

// All reads all adaptiveValues for the given PlatformID
// from dynamo table
func (p PlatformDAOImpl) All() (res []models.AdaptiveValue, err error) {
	var values []models.AdaptiveValue
	err = p.DNS.Dynamo.QueryTableWithIndex(p.Name,
		awsutils.DynamoIndexExpression{
			IndexName: p.PlatformIDIndex,
			Condition: "platform_id = :platform_id",
			Attributes: map[string]interface{}{
				":platform_id": p.PlatformID,
			},
		}, map[string]string{}, true, -1, &values)
	res = models.AdaptiveValueFilterActive(values)
	return
}

// AllUnsafe reads all adaptiveValues for PlatformID and panics in case of errors
func (p PlatformDAOImpl) AllUnsafe() []models.AdaptiveValue {
	adaptiveValues, err := p.All()
	core.ErrorHandler(err, p.DNS.Namespace, fmt.Sprintf("Could not query table %s", p.Name))
	return adaptiveValues
}

// Create creates an adaptiveValue making sure that PlatformID is correct
func (p PlatformDAOImpl) Create(adaptiveValue models.AdaptiveValue) error {
	adaptiveValue2 := adaptiveValue
	adaptiveValue2.PlatformID = daosCommon.PlatformID(p.PlatformID)
	return p.DNS.Dynamo.PutTableEntry(adaptiveValue2, p.Name)
}
