package user

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

func ReadAllUserProfiles(userDAO user.DAO, platformID models.PlatformID) []models.UserProfile {
	users := userDAO.ReadByPlatformIDUnsafe(platformID)
	return user.ConvertUsersToUserProfilesAndRemoveAdaptiveBot(users)
}

// Engagements retrieves engagements for a user, passing 0 retrieves not-yet-answered engagements
func Engagements(userID string, userEngagementsTable, userEngagementsAnsweredIndex string, answered int) []models.UserEngagement {
	var engs []models.UserEngagement
	dns := common.DeprecatedGetGlobalDns()
	// Query engagements table for the user's engagements
	err := dns.Dynamo.QueryTableWithIndex(userEngagementsTable, awsutils.DynamoIndexExpression{
		IndexName: userEngagementsAnsweredIndex,
		// there is no != operator for ConditionExpression
		Condition: "user_id = :u AND answered = :a",
		Attributes: map[string]interface{}{
			":u": userID,
			":a": answered,
		},
	}, map[string]string{}, true, -1, &engs)
	core.ErrorHandler(err, dns.Namespace, fmt.Sprintf("Could not query %s index on %s table",
		userEngagementsAnsweredIndex, userEngagementsTable))
	return engs
}

// NotPostedUnansweredEngagements retrieves not-answered and not-ignored engagements that aren't yet posted to the user
func NotPostedUnansweredNotIgnoredEngagements(
	userID string, userEngagementsTable, userEngagementsAnsweredIndex string) (res []models.UserEngagement) {
	engs := Engagements(userID, userEngagementsTable, userEngagementsAnsweredIndex, 0)
	for _, each := range engs {
		if each.PostedAt == "" && each.Ignored == 0 {
			res = append(res, each)
		}
	}
	return
}
