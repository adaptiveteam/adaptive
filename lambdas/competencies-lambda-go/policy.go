package competencies

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
)

// DoesUserHaveWriteAccessToCompetencies check access policy for user access to add/edit competencies
func DoesUserHaveWriteAccessToCompetencies(userID string) bool {
	return IsUserInCommunity(userID, community.Competency)
}

// IsUserInCommunity checks if the user in the community
func IsUserInCommunity(userID string, aCommunity community.AdaptiveCommunity) bool {
	return community.IsUserInCommunity(userID, communityUsersTable, communityUsersUserCommunityIndex, aCommunity)
}
