package communityUser

import (
	"fmt"

	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/adaptiveteam/adaptive/daos/common"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DAO is a CRUD wrapper around the _community_users Dynamo DB table
type DAO interface {
	IsUserInCommunity(teamID models.TeamID, channelID string, userID string) bool
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	models.CommunityUsersTableSchema
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace string,
	table models.CommunityUsersTableSchema) DAO {
	return DAOImpl{Dynamo: dynamo, Namespace: namespace,
		CommunityUsersTableSchema: table,
	}
}

// NewDAOFromSchema creates an instance of DAO that will provide access to the table
func NewDAOFromSchema(dynamo *awsutils.DynamoRequest, namespace string, schema models.Schema) DAO {
	return NewDAO(dynamo, namespace, schema.CommunityUsers)
}

// DeactivateUserFromCommunity deletes a user from community
func DeactivateUserFromCommunity(teamID models.TeamID, channelID string, userID string) func (conn common.DynamoDBConnection) (err error) {
	return func (conn common.DynamoDBConnection) (err error) {
		err = adaptiveCommunityUser.Deactivate(channelID, userID)(conn)
		return wrapError(err, "DeactivateUserFromCommunity("+userID+","+channelID+")")
	}
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

// IsUserInCommunity checks if a user is part of an Adaptive Community
func (d DAOImpl) IsUserInCommunity(teamID models.TeamID, communityID string, userID string) bool {
	connGen := common.CreateConnectionGenFromEnv()
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	acus, err2 := adaptiveCommunityUser.ReadByUserIDCommunityID(communityID, userID)(conn)
	core.ErrorHandlerf(err2, d.Namespace, "ReadByUserIDCommunityID(communityID=%s, userID=%s", communityID, userID)

	return len(acus) > 0
}

func DeactivateAllCommunityMembers(teamID models.TeamID, channelID string) func (conn common.DynamoDBConnection) (err error) {
	return func (conn common.DynamoDBConnection) (err error) {
		commUsers, err := adaptiveCommunityUser.ReadByChannelID(channelID)(conn)
		if err == nil {
			for _, each := range commUsers {
				err := DeactivateUserFromCommunity(teamID, channelID, each.UserID)
				if err != nil {
					break
				}
			}
		}
		return wrapError(err, "removeCommunityMembers("+channelID+")")
	}
}

func DeactivateAllCommunityMembersUnsafe(teamID models.TeamID, channelID string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := DeactivateAllCommunityMembers(teamID, channelID)(conn)
		core.ErrorHandler(err2, "DeactivateAllCommunityMembersUnsafe", "DeactivateAllCommunityMembersUnsafe channelID=" + channelID)
	}
}
