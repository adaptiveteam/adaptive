package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
)

// DoesUserHaveWriteAccessToHolidays check access policy for user access to add/edit holidays
func DoesUserHaveWriteAccessToHolidays(userID string) bool {
	return IsUserInCommunity(userID, community.HR)
}

// IsUserInCommunity checks if the user in the community
func IsUserInCommunity(userID string, aCommunity community.AdaptiveCommunity) bool {
	return community.IsUserInCommunity(userID, communityUsersTable, communityUsersUserCommunityIndex, aCommunity)
}
