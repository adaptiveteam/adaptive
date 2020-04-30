package user

import (
	"fmt"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
)

// UserProfilesIntersect only keeps the given users.
func UserProfilesIntersect(userProfiles []models.UserProfile, userIDs []string) (userProfilesOut []models.UserProfile) {
	for _, each := range userProfiles {
		if core.ListContainsString(userIDs, each.Id) {
			userProfilesOut = append(userProfilesOut, each)
		}
	}
	return
}

// UserProfilesSubtract removes the given users from the list
// NB! O(N*M)
func UserProfilesSubtract(userProfiles []models.UserProfile, userIDs []string) (userProfilesOut []models.UserProfile) {
	for _, each := range userProfiles {
		if !core.ListContainsString(userIDs, each.Id) {
			userProfilesOut = append(userProfilesOut, each)
		}
	}
	return
}

// MapUserProfilesToMenuOptions converts user profiles to menu options
func MapUserProfilesToMenuOptions(userProfiles []models.UserProfile) []ebm.MenuOption {
	var menuOptions []ebm.MenuOption
	for _, user := range userProfiles {
		menuOptions = append(menuOptions, ebm.MenuOption{Text: user.DisplayName, Value: user.Id})
	}
	return menuOptions
}

// SelectUserAction renders the given list of users as a select action
func SelectUserAction(mc models.MessageCallback, userProfiles []models.UserProfile) (action model.AttachmentAction) {
	options := MapUserProfilesToMenuOptions(userProfiles)

	action = *models.SelectAttachAction(mc, models.Now,
		"Choose user ...", options, models.EmptyActionMenuOptionGroups())
	return
}

// SelectUserTemplateActions renders the given list of users as select action
// and adds Cancel action as well
func SelectUserTemplateActions(mc models.MessageCallback, userProfiles []models.UserProfile) []model.AttachmentAction {
	action1 := SelectUserAction(mc, userProfiles)
	action2 := models.SimpleAttachAction(mc, models.Cancel, "Not now") // TODO: Danger?
	return []model.AttachmentAction{action1, *action2}
}

// UserSelectAttachments reads users, filters them twice, then renders options as attachments.
// Deprecated: Breaks SRP. Inline instead
func UserSelectAttachments(mc models.MessageCallback, userIDs, toFilterOutUserIDs []string,
	conn daosCommon.DynamoDBConnection) []model.AttachmentAction {
	userProfiles := ReadAllUserProfiles(conn)
	if len(userIDs) > 0 {
		// If users are passed, use them directly
		userProfiles = UserProfilesIntersect(userProfiles, userIDs)
	}
	if len(toFilterOutUserIDs) > 0 {
		// If filter is passed,exclude those from all users
		userProfiles = UserProfilesSubtract(userProfiles, toFilterOutUserIDs)
	}
	return SelectUserTemplateActions(mc, userProfiles)
}

// UserSelectEng reads users, filters them twice, then renders options as attachments,
// then creates engagement.
// Deprecated: Breaks SRP. Inline.
func UserSelectEng(userID, engagementsTable string,
	conn daosCommon.DynamoDBConnection,
	mc models.MessageCallback, users, toFilterUsers []string,
	text, context string, check models.UserEngagementCheckWithValue) {
	teamID := models.ParseTeamID(conn.PlatformID)
	attachs := UserSelectAttachments(mc, users, toFilterUsers, conn)
	dns := common.DeprecatedGetGlobalDns()

	utils.AddChatEngagement(mc, "", text, fmt.Sprintf("Select one of the users for %s", context), userID,
		attachs, []ebm.AttachmentField{}, teamID, true, engagementsTable, dns.Dynamo,
		dns.Namespace, time.Now().Unix(), check)
}
