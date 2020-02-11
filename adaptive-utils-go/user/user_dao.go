package user

import (
	"github.com/pkg/errors"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DAO is a wrapper around the _adaptive_users Dynamo DB table to work with adaptive-users table (CRUD)
type DAO interface {
	Read(userID string) (models.User, error)
	ReadUnsafe(userID string) models.User
	ReadByPlatformIDUnsafe(platformID models.PlatformID) (users []models.User)
	Create(user models.User) error
	CreateUnsafe(user models.User)
	UserIDsToDisplayNamesUnsafe(userIDs []string) (res []models.KvPair)
	Delete(userID string) error
	DeleteUnsafe(userID string)
	Update(user models.User) error
	UpdateUnsafe(user models.User)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	models.AdaptiveUsersTableSchema
}

// NewDAO creates an instance of DAO that will provide access to ClientPlatformToken table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, table string) DAO {
	if table == "" {
		panic("Cannot create User DAO without table")
	}
	return DAOImpl{Dynamo: dynamo, Namespace: namespace,
		AdaptiveUsersTableSchema: models.AdaptiveUsersTableSchema{Name: table},
	}
}

// NewDAOFromSchema creates an instance of DAO that will provide access to adaptiveValues table
func NewDAOFromSchema(dynamo *awsutils.DynamoRequest, namespace string, schema models.Schema) DAO {
	return DAOImpl{Dynamo: dynamo, Namespace: namespace,
		AdaptiveUsersTableSchema: schema.AdaptiveUsers}
}

// Read reads User
func (d DAOImpl) Read(userID string) (out models.User, err error) {
	if userID == "" {
		err = errors.Errorf("An attempt to read user with an empty user id")
	} else {
		err = d.Dynamo.GetItemFromTable(d.Name, idParams(userID), &out)
	}
	return
}

// ReadUnsafe reads the User. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(userID string) models.User {
	out, err := d.Read(userID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not find %s in %s", userID, d.Name))
	return out
}

// Create saves the User.
func (d DAOImpl) Create(user models.User) error {
	return d.Dynamo.PutTableEntryWithCondition(user, d.Name,
		"attribute_not_exists(id)")
}

// CreateUnsafe saves the User.
func (d DAOImpl) CreateUnsafe(user models.User) {
	err := d.Create(user)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create %s in %s", user.ID, d.Name))
}

func dynString(str string) (attr *dynamodb.AttributeValue) {
	return &dynamodb.AttributeValue{S: aws.String(str)}
}

// UserIDsToDisplayNamesUnsafe converts a bunch of user ids to their names
// NB! O(n)! TODO: implement a query that returns many users at once.
func (d DAOImpl) UserIDsToDisplayNamesUnsafe(userIDs []string) (res []models.KvPair) {
	if len(userIDs) > 10 {
		fmt.Println("WARN: Very slow user data fetching")
	}
	for _, userID := range userIDs {
		user := d.ReadUnsafe(userID)
		res = append(res, models.KvPair{Key: user.DisplayName, Value: userID})
	}
	return
}

// UserIDsToDisplayNames is a function that fetches user information for ids.
// For each user their display name is put to `Name` and user id to `Value`.
type UserIDsToDisplayNames func([]string) []models.KvPair

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

func (d DAOImpl) ReadByPlatformIDUnsafe(platformID models.PlatformID) (users []models.User) {
	err := d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: d.PlatformIndex,
		// there is no != operator for ConditionExpression
		Condition: "platform_id = :p",
		Attributes: map[string]interface{}{
			":p": platformID,
		},
	}, map[string]string{}, true, -1, &users)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query %s index on %s table",
		d.PlatformIndex, d.Name))
	return
}

// Update saves the changed User.
func (d DAOImpl) Update(user models.User) error {
	return d.Dynamo.PutTableEntry(user, d.Name)
}

// UpdateUnsafe saves the changed User.
func (d DAOImpl) UpdateUnsafe(user models.User) {
	err := d.Update(user)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create %s in %s", user.ID, d.Name))
}

// ConvertUsersToUserProfilesAndRemoveAdaptiveBot converts users to user profiles.
func ConvertUsersToUserProfilesAndRemoveAdaptiveBot(users []models.User) (userProfiles []models.UserProfile) {
	for _, each := range users {
		if !each.IsShared && !each.IsAdaptiveBot && each.DeactivatedAt == "" {
			userProfiles = append(userProfiles,
				models.UserProfile{Id: each.ID,
					DisplayName: each.DisplayName,
					FirstName:   each.FirstName,
					LastName:    each.LastName,
					Timezone:    each.Timezone})
		}
	}
	return
}
