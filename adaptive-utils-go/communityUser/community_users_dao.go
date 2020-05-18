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
	DeactivateAllCommunityMembers(teamID models.TeamID, channelID string) (err error)
	DeactivateAllCommunityMembersUnsafe(teamID models.TeamID, channelID string)
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

// Delete removes user from db
func (d DAOImpl) Delete(userID string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(userID))
}

// DeleteUnsafe deletes user and panics in case of errors.
func (d DAOImpl) DeleteUnsafe(userID string) {
	err := d.Delete(userID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete %s in %s", userID, d.Name))
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

func (d DAOImpl) DeactivateAllCommunityMembers(teamID models.TeamID, channelID string) (err error) {
	connGen := common.CreateConnectionGenFromEnv()
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
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

func (d DAOImpl) DeactivateAllCommunityMembersUnsafe(teamID models.TeamID, channelID string) {
	err := d.DeactivateAllCommunityMembers(teamID, channelID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("removeCommunityMembersUnsafe: Could not query %s table on %s index",
		d.Name, d.ChannelIndex))
}
