package user

import (
	"fmt"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	// daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	daosUser "github.com/adaptiveteam/adaptive/daos/user"
	"github.com/nlopes/slack"
)

// DAO is a wrapper around the _adaptive_users Dynamo DB table to work with adaptive-users table (CRUD)
type DAO = daosUser.DAO

// NewDAOFromSchema creates an instance of DAO that will provide access to adaptiveValues table
func NewDAOFromSchema(dynamo *awsutils.DynamoRequest, namespace string, schema models.Schema) DAO {
	return daosUser.NewDAOByTableName(dynamo, namespace, schema.AdaptiveUsers.Name)
	//  DAOImpl{Dynamo: dynamo, Namespace: namespace,
	// AdaptiveUsersTableSchema: schema.AdaptiveUsers}
}


// UserIDsToDisplayNamesUnsafe converts a bunch of user ids to their names
// NB! O(n)! TODO: implement a query that returns many users at once.
func UserIDsToDisplayNamesUnsafe(dao DAO) func(userIDs []string) (res []models.KvPair) {
	return func(userIDs []string) (res []models.KvPair) {
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
