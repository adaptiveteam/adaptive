package values

import (
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
	Read(adaptiveValueID string) (models.AdaptiveValue, error)
	ReadUnsafe(adaptiveValueID string) models.AdaptiveValue
	Update(adaptiveValue models.AdaptiveValue) error
	Delete(adaptiveValueID string) error
	Deactivate(adaptiveValueID string) error

	ForPlatformID(platformID string) PlatformDAO
}

// PlatformDAO is a set of utilities that work for a fixed `platformID`
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
	PlatformID string
	DAOImpl
}

// NewDAO creates an instance of DAO that will provide access to adaptiveValues table
func NewDAO(dns *common.DynamoNamespace, table string, index string) DAO {
	if table == "" {
		panic("Cannot create adaptiveValues DAO without table")
	}
	if index == "" {
		panic("Cannot create adaptiveValues DAO without index")
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

// Read reads the adaptiveValue
func (d DAOImpl) Read(adaptiveValueID string) (models.AdaptiveValue, error) {
	params := map[string]*dynamodb.AttributeValue{
		"id": {S: aws.String(adaptiveValueID)},
	}
	var out models.AdaptiveValue
	err := d.DNS.Dynamo.QueryTable(d.Name, params, &out)
	return out, err
}

// ReadUnsafe reads the adaptiveValue. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(adaptiveValueID string) models.AdaptiveValue {
	adaptiveValue, err := d.Read(adaptiveValueID)
	core.ErrorHandler(err, d.DNS.Namespace, fmt.Sprintf("Could not find %s in %s", adaptiveValueID, d.Name))
	return adaptiveValue
}

// Delete removes the adaptiveValue
func (d DAOImpl) Delete(adaptiveValueID string) error {
	userParams := map[string]*dynamodb.AttributeValue{
		"id": {S: aws.String(adaptiveValueID)},
	}
	return d.DNS.Dynamo.DeleteEntry(d.Name, userParams)
}

// Update updates the adaptiveValue by ID
func (d DAOImpl) Update(adaptiveValue models.AdaptiveValue) error {
	return d.DNS.Dynamo.PutTableEntry(adaptiveValue, d.Name)
}

func (d DAOImpl) Deactivate(adaptiveValueID string) (err error) {
	adaptiveValue, err := d.Read(adaptiveValueID)
	if err == nil {
		adaptiveValue.DeactivatedOn = core.ISODateLayout.Format(time.Now())
		err = d.Update(adaptiveValue)
	}
	return
}

// UpdateGivenAttributes updates the adaptiveValue by ID
// deprecated. This might be used if we want to update only some of the attributes
func (d DAOImpl) UpdateGivenAttributes(adaptiveValue models.AdaptiveValue) error {
	exprAttributes := map[string]*dynamodb.AttributeValue{
		":value_name": {
			S: aws.String(adaptiveValue.Name),
		},
		":description": {
			S: aws.String(adaptiveValue.Description),
		},
		":value_type": {
			S: aws.String(adaptiveValue.ValueType),
		},
		":platform_id": {
			S: aws.String(string(adaptiveValue.PlatformID)),
		},
	}
	key := map[string]*dynamodb.AttributeValue{
		"id": {S: aws.String(adaptiveValue.ID)},
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
func (d DAOImpl) ForPlatformID(platformID string) PlatformDAO {
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
