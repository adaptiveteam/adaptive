package user

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	daosUser "github.com/adaptiveteam/adaptive/daos/user"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/nlopes/slack"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
)

// DAO is a wrapper around the _adaptive_users Dynamo DB table to work with adaptive-users table (CRUD)
type DAO = daosUser.DAO 
// interface {
// 	Read(userID string) (models.User, error)
// 	ReadUnsafe(userID string) models.User
// 	ReadByPlatformIDUnsafe(teamID models.TeamID) (users []models.User)
// 	Create(user models.User) error
// 	CreateUnsafe(user models.User)
// 	UserIDsToDisplayNamesUnsafe(userIDs []string) (res []models.KvPair)
// 	Update(user models.User) error
// 	UpdateUnsafe(user models.User)
// }

// // DAOImpl - a container for all information needed to access a DynamoDB table
// type DAOImpl struct {
// 	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
// 	Namespace string                  `json:"namespace"`
// 	models.AdaptiveUsersTableSchema
// }

// NewDAO creates an instance of DAO that will provide access to ClientPlatformToken table
var NewDAOByTableName = daosUser.NewDAOByTableName

// TableName is a function that returns `_user` table name having client id
var TableName = func(clientID string) string { return clientID + "_adaptive_users" }

// DAOFromConnectionGen -
func DAOFromConnectionGen(conn daosCommon.DynamoDBConnectionGen) DAO {
	return NewDAOByTableName(conn.Dynamo, "UserDAO", TableName(conn.TableNamePrefix))
}
// DAOFromConnection -
func DAOFromConnection(conn daosCommon.DynamoDBConnection) DAO {
	return NewDAOByTableName(conn.Dynamo, "UserDAO", TableName(conn.ClientID))
}
// NewDAOFromSchema creates an instance of DAO that will provide access to adaptiveValues table
func NewDAOFromSchema(dynamo *awsutils.DynamoRequest, namespace string, schema models.Schema) DAO {
	return daosUser.NewDAOByTableName(dynamo, namespace, schema.AdaptiveUsers.Name)
	//  DAOImpl{Dynamo: dynamo, Namespace: namespace,
		// AdaptiveUsersTableSchema: schema.AdaptiveUsers}
}

// // Read reads User
// func (d DAOImpl) Read(userID string) (out models.User, err error) {
// 	if userID == "" {
// 		err = errors.Errorf("An attempt to read user with an empty user id")
// 	} else {
// 		err = d.Dynamo.GetItemFromTable(d.Name, idParams(userID), &out)
// 	}
// 	return
// }

// // ReadUnsafe reads the User. Panics in case of any errors
// func (d DAOImpl) ReadUnsafe(userID string) models.User {
// 	out, err := d.Read(userID)
// 	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not find %s in %s", userID, d.Name))
// 	return out
// }

// // Create saves the User.
// func (d DAOImpl) Create(user models.User) error {
// 	return d.Dynamo.PutTableEntryWithCondition(user, d.Name,
// 		"attribute_not_exists(id)")
// }

// // CreateUnsafe saves the User.
// func (d DAOImpl) CreateUnsafe(user models.User) {
// 	err := d.Create(user)
// 	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create %s in %s", user.ID, d.Name))
// }

// func dynString(str string) (attr *dynamodb.AttributeValue) {
// 	return &dynamodb.AttributeValue{S: aws.String(str)}
// }

// UserIDsToDisplayNamesUnsafe converts a bunch of user ids to their names
// NB! O(n)! TODO: implement a query that returns many users at once.
func UserIDsToDisplayNamesUnsafe(dao DAO) func(userIDs []string) (res []models.KvPair) {
	return func (userIDs []string) (res []models.KvPair) {
		if len(userIDs) > 10 {
			fmt.Println("WARN: Very slow user data fetching")
		}
		for _, userID := range userIDs {
			user := dao.ReadUnsafe(userID)
			res = append(res, models.KvPair{Key: user.DisplayName, Value: userID})
		}
		return
	}
}

// // UserIDsToDisplayNames is a function that fetches user information for ids.
// // For each user their display name is put to `Name` and user id to `Value`.
// type UserIDsToDisplayNames func([]string) []models.KvPair

// func idParams(id string) map[string]*dynamodb.AttributeValue {
// 	params := map[string]*dynamodb.AttributeValue{
// 		"id": dynString(id),
// 	}
// 	return params
// }

// func (d DAOImpl) ReadByPlatformIDUnsafe(teamID models.TeamID) (users []models.User) {
// 	err := d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
// 		IndexName: d.PlatformIndex,
// 		// there is no != operator for ConditionExpression
// 		Condition: "platform_id = :p",
// 		Attributes: map[string]interface{}{
// 			":p": teamID,
// 		},
// 	}, map[string]string{}, true, -1, &users)
// 	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query %s index on %s table",
// 		d.PlatformIndex, d.Name))
// 	return
// }

// // Update saves the changed User.
// func (d DAOImpl) Update(user models.User) error {
// 	return d.Dynamo.PutTableEntry(user, d.Name)
// }

// // UpdateUnsafe saves the changed User.
// func (d DAOImpl) UpdateUnsafe(user models.User) {
// 	err := d.Update(user)
// 	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create %s in %s", user.ID, d.Name))
// }

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

// ConvertSlackUserToUser -
func ConvertSlackUserToUser(user slack.User, teamID models.TeamID) (mUser models.User) {
	now := core.CurrentRFCTimestamp()
	deactivatedAt := ""
	if user.Deleted {
		deactivatedAt = now
	}
	return models.User{
		ID:             user.ID,
		DisplayName:    user.RealName,
		FirstName:      user.Profile.FirstName,
		LastName:       user.Profile.LastName,
		Timezone:       user.TZ,
		TimezoneOffset: user.TZOffset,
		PlatformID:     teamID.ToPlatformID(),
		IsAdmin:        user.IsAdmin,
		DeactivatedAt:  deactivatedAt,
		CreatedAt:      now,
		IsShared:       false,
	}
}

const UserID_Requested = "requested"
const UserID_None = "none"

func IsSpecialUserID(userID string) bool {
	return userID == UserID_None || userID == UserID_Requested
}

func IsSpecialOrEmptyUserID(userID string) bool {
	return IsSpecialUserID(userID) || userID == ""
}
