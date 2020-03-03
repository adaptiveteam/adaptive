package communityUser

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"fmt"
)

// DAO is a CRUD wrapper around the _community_users Dynamo DB table 
type DAO interface {
	Create(user models.AdaptiveCommunityUser3) error
	CreateUnsafe(user models.AdaptiveCommunityUser3)
	// Read reads user by id, returns zero or one results
	Read(userID string) ([]models.AdaptiveCommunityUser3, error)
	ReadUnsafe(userID string) []models.AdaptiveCommunityUser3
	ReadCommunityUsers(channelID string) (users []models.AdaptiveCommunityUser3, err error)
	ReadCommunityMembers(channelID string, teamID models.TeamID) (users []models.AdaptiveCommunityUser3, err error)
	ReadCommunityMembersUnsafe(channelID string, teamID models.TeamID) (users []models.AdaptiveCommunityUser3)
	ReadAnyCommunityUsers(teamID models.TeamID) (users []models.AdaptiveCommunityUser3, err error)
	ReadAnyCommunityUsersUnsafe(teamID models.TeamID) (users []models.AdaptiveCommunityUser3)
	ReadCommunityUserOptional(channelID string, userID string) (user []models.AdaptiveCommunityUser3, err error)
	ReadCommunityUserOptionalUnsafe(channelID string, userID string) (user []models.AdaptiveCommunityUser3)
	DeleteUserFromCommunity(channelID string, userID string) (err error)
	DeleteAllCommunityMembers(channelID string) (err error)
	DeleteAllCommunityMembersUnsafe(channelID string)
	// Delete(userID string) error
	// DeleteUnsafe(userID string)
	IsUserInCommunity(channelID string, userID string) bool
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
// Create saves the User.
func (d DAOImpl) Create(user models.AdaptiveCommunityUser3) error {
	return d.Dynamo.PutTableEntryWithCondition(user, d.Name, 
		"attribute_not_exists(id)")
}
// CreateUnsafe saves the User.
func (d DAOImpl) CreateUnsafe(user models.AdaptiveCommunityUser3) {
	err := d.Create(user)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create %s in %s", user.UserID, d.Name))

}
// Read reads user by id, returns zero or one results
func (d DAOImpl) Read(userID string) (out []models.AdaptiveCommunityUser3, err error) {
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: d.UserIndex,
		Condition: "user_id = :u",
		Attributes: map[string]interface{}{
			":u": userID,
		},
	}, map[string]string{}, true, -1, &out)
	return
}
// ReadUnsafe reads data. Panics in case of errors
func (d DAOImpl) ReadUnsafe(userID string) []models.AdaptiveCommunityUser3 {
	out, err2 := d.Read(userID)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not find %s in %s using index %s", userID, d.Name, d.UserIndex))
	return out
}
// ReadCommunityUsers reads users of the channel
// NB! Use another method with PlatformID argument.
func (d DAOImpl) ReadCommunityUsers(channelID string) (users []models.AdaptiveCommunityUser3, err error) {
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: d.ChannelIndex,
		Condition: "channel_id = :c",
		Attributes: map[string]interface{}{
			":c": channelID,
		},
	}, map[string]string{}, true, -1, &users)
	return
}
// ReadCommunityUsersUnsafe reads&panics
func (d DAOImpl) ReadCommunityUsersUnsafe(channelID string) (users []models.AdaptiveCommunityUser3) {
	users, err := d.ReadCommunityUsers(channelID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query %s table on %s index", d.Name, d.UserCommunityIndex))
	return
}
// ReadCommunityMembers reads members using teamID
func (d DAOImpl) ReadCommunityMembers(channelID string, teamID models.TeamID) (users []models.AdaptiveCommunityUser3, err error) {
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: d.UserCommunityIndex,
		Condition: "platform_id = :pi AND community_id = :c",
		Attributes: map[string]interface{}{
			":c":  channelID,
			":pi": teamID,
		},
	}, map[string]string{}, true, -1, &users)
	return
}
// ReadCommunityMembersUnsafe read and panic
func (d DAOImpl) ReadCommunityMembersUnsafe(channelID string, teamID models.TeamID) (users []models.AdaptiveCommunityUser3) {
	users, err := d.ReadCommunityMembers(channelID, teamID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query %s table on %s index",
		d.Name, d.UserCommunityIndex))
	return
}

func (d DAOImpl) ReadAnyCommunityUsers(teamID models.TeamID) (users []models.AdaptiveCommunityUser3, err error) {
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: d.CommunityIndex,
		Condition: "platform_id = :pi",
		Attributes: map[string]interface{}{
			":pi": teamID,
		},
	}, map[string]string{}, true, -1, &users)
	return
}

func (d DAOImpl) ReadAnyCommunityUsersUnsafe(teamID models.TeamID) (users []models.AdaptiveCommunityUser3) {
	users, err := d.ReadAnyCommunityUsers(teamID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query %s table on %s index",
		d.Name, d.CommunityIndex))
	return
}

// DeleteUserFromCommunity deletes a user from community
func (d DAOImpl) DeleteUserFromCommunity(channelID string, userID string) (err error) {
	commUserParams := map[string]*dynamodb.AttributeValue{
		"channel_id": dynString(channelID),
		"user_id":    dynString(userID),
	}
	err = d.Dynamo.DeleteEntry(d.Name, commUserParams)
	return wrapError(err, "deleteUserFromCommunity("+userID+","+channelID+")")
}

func dynString(str string) (attr *dynamodb.AttributeValue) {
	return &dynamodb.AttributeValue{S: aws.String(str)}
}

// Delete removes user from db
func (d DAOImpl)Delete(userID string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(userID))
}

// DeleteUnsafe deletes user and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(userID string) {
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
// ReadCommunityUserOptional checks if the user is in the channel
func (d DAOImpl)ReadCommunityUserOptional(channelID string, userID string) (user []models.AdaptiveCommunityUser3, err error) {
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: d.UserCommunityIndex,
		Condition: "user_id = :u and community_id = :c",
		Attributes: map[string]interface{}{
			":u": userID,
			":c": channelID,
		},
	}, map[string]string{}, true, -1, &user)
	return
}
// ReadCommunityUserOptionalUnsafe read&panic
func (d DAOImpl)ReadCommunityUserOptionalUnsafe(channelID string, userID string) (user []models.AdaptiveCommunityUser3) {
	user, err := d.ReadCommunityUserOptional(channelID, userID) 
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query %s table on %s index",
		d.Name, d.UserCommunityIndex))
	return
}
// IsUserInCommunity checks if a user is part of an Adaptive Community
func (d DAOImpl)IsUserInCommunity(channelID string, userID string) bool {
	return len(d.ReadCommunityUserOptionalUnsafe(channelID, userID)) > 0
}

func (d DAOImpl)DeleteAllCommunityMembers(channelID string) (err error) {
	commUsers, err := d.ReadCommunityUsers(channelID)
	if err == nil {
		for _, each := range commUsers {
			err := d.DeleteUserFromCommunity(channelID, each.UserID)
			if err != nil {
				break
			}
		}
	}
	return wrapError(err, "removeCommunityMembers("+channelID+")")
}

func (d DAOImpl)DeleteAllCommunityMembersUnsafe(channelID string) {
	err := d.DeleteAllCommunityMembers(channelID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("removeCommunityMembersUnsafe: Could not query %s table on %s index", 
		d.Name, d.ChannelIndex))
}
