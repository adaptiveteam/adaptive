package community

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"fmt"
)

// DAO is a CRUD wrapper around the _org_community (_community_users) Dynamo DB table 
type DAO interface {
	ReadByID(teamID models.TeamID, communityID string) (comm models.AdaptiveCommunity, err error)
	ReadByIDUnsafe(teamID models.TeamID, communityID string) (comm models.AdaptiveCommunity)
	ReadByChannelID(channelID string) ([]models.AdaptiveCommunity, error)
	ReadAll(teamID models.TeamID) ([]models.AdaptiveCommunity, error)
	Delete(teamID models.TeamID, communityID string) (err error)
	DeleteUnsafe(teamID models.TeamID, communityID string)
	Create(community models.AdaptiveCommunity) error
	CreateUnsafe(community models.AdaptiveCommunity)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	models.AdaptiveCommunityTableSchema
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace string, 
	table models.AdaptiveCommunityTableSchema) DAO {
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		AdaptiveCommunityTableSchema: table,
	}
}

// NewDAOFromSchema creates an instance of DAO that will provide access to the table
func NewDAOFromSchema(dynamo *awsutils.DynamoRequest, namespace string, schema models.Schema) DAO {	
	return NewDAO(dynamo, namespace, schema.AdaptiveCommunity)
}
// Create saves the AdaptiveCommunity.
func (d DAOImpl) Create(community models.AdaptiveCommunity) error {
	return d.Dynamo.PutTableEntryWithCondition(community, d.Name, 
		"attribute_not_exists(id)")
}
// CreateUnsafe saves the User.
func (d DAOImpl) CreateUnsafe(community models.AdaptiveCommunity) {
	err := d.Create(community)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create %s in %s",
		community.ID, d.Name))
}

func (d DAOImpl) ReadByID(teamID models.TeamID, communityID string) (comm models.AdaptiveCommunity, err error) {
	// Querying for admin community
	params := map[string]*dynamodb.AttributeValue{
		"id":          dynString(communityID),
		"platform_id": dynString(teamID.ToString()),
	}
	err = d.Dynamo.GetItemFromTable(d.Name, params, &comm)
	return
}

func (d DAOImpl) ReadByIDUnsafe(teamID models.TeamID, communityID string) (comm models.AdaptiveCommunity) {
	comm, err := d.ReadByID(teamID, communityID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query %s table", 
		d.Name))
	return
}


func (d DAOImpl) ReadByChannelID(channelID string)  (comms []models.AdaptiveCommunity, err error) {
	err = d.Dynamo.QueryTableWithIndex(
		d.Name, awsutils.DynamoIndexExpression{
		IndexName: d.ChannelIndex,
		Condition: "channel = :c",
		Attributes: map[string]interface{}{
			":c": channelID,
		},
	}, map[string]string{}, true, -1, &comms)
	err = wrapError(err, "subscribedCommunities")
	return
}

func (d DAOImpl) ReadAll(teamID models.TeamID) (comms []models.AdaptiveCommunity, err error) {
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: d.PlatformIndex,
		Condition: "platform_id = :p",
		Attributes: map[string]interface{}{
			":p": teamID,
		},
	}, map[string]string{}, true, -1, &comms)
	return
}

// Delete entry from communities table
func (d DAOImpl) Delete(teamID models.TeamID, communityID string) (err error) {
	commParams := idAndPlatformIDParams(communityID, teamID)
	err = d.Dynamo.DeleteEntry(d.Name, commParams)
	return
}

// DeleteUnsafe delete&panic
func (d DAOImpl) DeleteUnsafe(teamID models.TeamID, communityID string) {
	// Delete entry from communities table
	err := d.Delete(teamID, communityID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete entry from %s table", d.Name))
}

func dynString(str string) (attr *dynamodb.AttributeValue) {
	return &dynamodb.AttributeValue{S: aws.String(str)}
}

func idParams(id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue{
		"id": dynString(id),
	}
	return params
}

func wrapError(err error, name string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("{%s: %v}", name, err)
}

func idAndPlatformIDParams(id string, teamID models.TeamID) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue{
		"id":          dynString(id),
		"platform_id": dynString(teamID.ToString()),
	}
	return params
}
