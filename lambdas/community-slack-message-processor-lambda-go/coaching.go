package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	//eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/nlopes/slack"
)

func onRequestCoachClicked(request slack.InteractionCallback, mc models.MessageCallback, teamID models.TeamID) platform.Response {
	// Get coaching community members
	commMembers := communityUserDAO.ReadCommunityMembersUnsafe(string(community.Coaching), teamID)
	var userIDs []string
	for _, each := range commMembers {
		// Self user checking
		if each.UserID != request.User.ID {
			userIDs = append(userIDs, each.UserID)
		}
	}
	mc2 := *mc.WithTopic(CoachingName).WithAction(RequestCoach)
	users := userDAO.ReadByPlatformIDUnsafe(teamID.ToPlatformID())
	userProfiles := utilsUser.ConvertUsersToUserProfilesAndRemoveAdaptiveBot(users)
	filteredProfiles := user.UserProfilesIntersect(userProfiles, userIDs)
	attachmentActions := user.SelectUserTemplateActions(mc2, filteredProfiles)
	
	return platform.OverrideByURL(platform.ResponseURLMessageID{ResponseURL: request.ResponseURL},
		platform.MessageContent{
			Message: ListOfCoachesWelcomeMessage,
			Attachments: []ebm.Attachment{ebm.Attachment{
				Text: string(ListOfCoachesWelcomeMessage),
				Fallback: fmt.Sprintf("Select one of the users for %s", CoachingName),
				Actions: attachmentActions,
			}},
	})
}

func communityNamespaceCoachingDialogSubmissionHandler(dialog slack.InteractionCallback, msgState MsgState, mc models.MessageCallback, form map[string]string) {
	switch mc.Action {
	case CoachConfirm:
		// This is the dialog for when a coachee doesn't accept coach
		rejectionComments := ui.RichText(form[CommentsName])
		attachs := viewCommentsAttachment(mc,
			CoachingRequestRejectionReasonTitleToCoachee,
			rejectionComments)
		publish(models.PlatformSimpleNotification{UserId: dialog.User.ID, Channel: dialog.Channel.ID, Ts: msgState.ThreadTs, Attachments: attachs})
		err := d.PutTableEntry(models.CoachingRejection{Id: mc.ToCallbackID(), CoachRequested: true, CoacheeRejected: true,
			Comments: string(rejectionComments)}, coachingRejectionsTable)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not add entry to %s table", coachingRejectionsTable))
	case RequestCoach:
		// TODO: Remove this
		// This is the dialog for when a coach doesn't accept coachee
		rejectionComments := ui.RichText(form[CommentsName])
		attachs := viewCommentsAttachment(mc,
			CoachingRequestRejectionReasonTitleToCoach,
			rejectionComments)
		publish(models.PlatformSimpleNotification{UserId: dialog.User.ID, Channel: dialog.Channel.ID, Ts: msgState.ThreadTs, Attachments: attachs})
		err := d.PutTableEntry(models.CoachingRejection{Id: mc.ToCallbackID(), CoacheeRequested: true, CoachRejected: true,
			Comments: string(rejectionComments)}, coachingRejectionsTable)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not add entry to %s table", coachingRejectionsTable))
		// Updating engagement as answered
		utils.UpdateEngAsAnswered(mc.Target, mc.ToCallbackID(), engagementTable, d, namespace)
	}
}
