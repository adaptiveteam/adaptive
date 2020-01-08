package platform

import (
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

// DAO - wrapper around a Dynamo DB table to work with PlatformID -> PlatformToken mapping
type DAO interface {
	Read(platformID models.PlatformID) (models.ClientPlatformToken, error)
	ReadUnsafe(platformID models.PlatformID) models.ClientPlatformToken
	GetPlatformTokenUnsafe(platformID models.PlatformID) string
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	models.ClientPlatformTokenTableSchema
}

// NewDAO creates an instance of DAO that will provide access to ClientPlatformToken table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, table string) DAO {
	if table == "" { panic("Cannot create ClientPlatformToken DAO without table") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		ClientPlatformTokenTableSchema: models.ClientPlatformTokenTableSchema{Name: table},
	}
}

// NewDAOFromSchema creates an instance of DAO that will provide access to adaptiveValues table
func NewDAOFromSchema(dynamo *awsutils.DynamoRequest, namespace string, schema models.Schema) DAO {	
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		ClientPlatformTokenTableSchema: schema.ClientPlatformTokens}
}

// Read reads ClientPlatformToken
func (d DAOImpl) Read(platformID models.PlatformID) (models.ClientPlatformToken, error) {
	params := map[string]*dynamodb.AttributeValue{
		"platform_id": dynString(string(platformID)),
	}
	var out models.ClientPlatformToken
	err := d.Dynamo.QueryTable(d.Name, params, &out)
	return out, err
}

// ReadUnsafe reads the ClientPlatformToken. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(platformID models.PlatformID) models.ClientPlatformToken {
	out, err := d.Read(platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not find %s in %s", platformID, d.Name))
	return out
}

func dynString(str string) (attr *dynamodb.AttributeValue) {
	return &dynamodb.AttributeValue{S: aws.String(str)}
}
// GetPlatformTokenUnsafe reads platform token from database
func (d DAOImpl) GetPlatformTokenUnsafe(platformID models.PlatformID) string {
	return d.ReadUnsafe(platformID).PlatformToken
}