package user

import (
	"github.com/adaptiveteam/adaptive/daos/common"
	"fmt"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosUser "github.com/adaptiveteam/adaptive/daos/user"
	"github.com/slack-go/slack"
)

// DAO is a wrapper around the _adaptive_users Dynamo DB table to work with adaptive-users table (CRUD)
type DAO = daosUser.DAO

// UserIDsToDisplayNamesConnUnsafe converts a bunch of user ids to their names
// NB! O(n)! TODO: implement a query that returns many users at once.
func UserIDsToDisplayNamesConnUnsafe(conn common.DynamoDBConnection) func(userIDs []string) (res []models.KvPair) {
	return func(userIDs []string) (res []models.KvPair) {
		if len(userIDs) > 10 {
			fmt.Println("WARN: Very slow user data fetching")
		}
		for _, userID := range userIDs {
			user := daosUser.ReadUnsafe(conn.PlatformID, userID)(conn)
			res = append(res, models.KvPair{Key: user.DisplayName, Value: userID})
		}
		return
	}
}


// ConvertUsersToUserProfilesAndRemoveAdaptiveBot converts users to user profiles.
func ConvertUsersToUserProfilesAndRemoveAdaptiveBot(users []models.User) (userProfiles []models.UserProfile) {
	users = daosUser.UserFilterActive(users)
	for _, each := range users {
		if !each.IsShared && !each.IsAdaptiveBot {
			userProfiles = append(userProfiles, models.ConvertUserToProfile(each))
		}
	}
	return
}

// ConvertSlackUserToUser -
func ConvertSlackUserToUser(slackUser slack.User, 
	teamID models.TeamID,
	adaptiveBotID string,
) (mUser models.User) {
	now := core.CurrentRFCTimestamp()
	deactivatedAt := ""
	if slackUser.Deleted {
		deactivatedAt = now
	}
	return models.User{
		ID:             slackUser.ID,
		DisplayName:    slackUser.RealName,
		FirstName:      slackUser.Profile.FirstName,
		LastName:       slackUser.Profile.LastName,
		Timezone:       slackUser.TZ,
		TimezoneOffset: slackUser.TZOffset,
		PlatformID:     teamID.ToPlatformID(),
		IsAdmin:        slackUser.IsAdmin,
		IsAdaptiveBot:  slackUser.ID == adaptiveBotID || (adaptiveBotID == "" && slackUser.IsBot),
		DeactivatedAt:  deactivatedAt,
		CreatedAt:      now,
		ModifiedAt:     now,
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
