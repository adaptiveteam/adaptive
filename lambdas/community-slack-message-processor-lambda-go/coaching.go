package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"fmt"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	daosUser "github.com/adaptiveteam/adaptive/daos/user"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/slack-go/slack"
)

func onRequestCoachClicked(request slack.InteractionCallback, mc models.MessageCallback,
	conn daosCommon.DynamoDBConnection,
) platform.Response {
	// Get coaching community members
	commMembers := adaptiveCommunityUser.ReadByPlatformIDCommunityIDUnsafe(conn.PlatformID, string(community.Coaching))(conn)
	var userIDs []string
	for _, each := range commMembers {
		// Self user checking
		if each.UserID != request.User.ID {
			userIDs = append(userIDs, each.UserID)
		}
	}
	mc2 := *mc.WithTopic(CoachingName).WithAction(RequestCoach)
	users := daosUser.ReadByPlatformIDUnsafe(conn.PlatformID)(conn)
	userProfiles := utilsUser.ConvertUsersToUserProfilesAndRemoveAdaptiveBot(users)
	filteredProfiles := user.UserProfilesIntersect(userProfiles, userIDs)
	attachmentActions := user.SelectUserTemplateActions(mc2, filteredProfiles)

	return platform.OverrideByURL(platform.ResponseURLMessageID{ResponseURL: request.ResponseURL},
		platform.MessageContent{
			Message: ListOfCoachesWelcomeMessage,
			Attachments: []ebm.Attachment{ebm.Attachment{
				Text:     string(ListOfCoachesWelcomeMessage),
				Fallback: fmt.Sprintf("Select one of the users for %s", CoachingName),
				Actions:  attachmentActions,
			}},
		})
}
