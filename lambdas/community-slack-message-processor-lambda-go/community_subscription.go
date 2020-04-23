package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// TODO: save engagements.

// Subscribe menu item clicked
func onCommunitySubscribeClicked(
	request slack.InteractionCallback,
	teamID models.TeamID) (response platform.Response) {

	return platform.OverrideByURL(
		platform.ResponseURLMessageID{ResponseURL: request.ResponseURL},
		getSubscribeMessage(platform.ConversationID(request.Channel.ID), teamID, request.User.ID))
}

// A community is selected for subscription.
// communityID - Selected community, one of hr, user, coaching, users
// - it might also be of the form `community-type`:`parent-id`
func onCommunitySubscribeCommunityClicked(
	request slack.InteractionCallback,
	communityID string, //
	mc models.MessageCallback,
	teamID models.TeamID,
	conn common.DynamoDBConnection) {
	communityName := ui.PlainText(communityID)
	logger.Infof("Platform id for %s community: %s", communityID, teamID)
	// Let's add this channel as a new user
	// Get the information about the user who initiated this
	channelID := request.Channel.ID
	err2 := createCommunityFromCreatorUser(request.User.ID, channelID, communityID)
	if err2 != nil {
		if strings.Contains(err2.Error(), "ConditionalCheckFailedException") {
			logger.Infof("User %s already exists, not adding", request.User.ID)
		} else {
			logger.
				WithField("namespace", namespace).
				WithError(err2).
				Errorf("Could not add %s to %s table", request.Channel.ID, usersTable)
		}
	}
	comm := adaptiveCommunity.AdaptiveCommunity{
		ID:         communityID,
		PlatformID: teamID.ToPlatformID(),
		ChannelID:  request.Channel.ID,
		Active:     true, RequestedBy: request.User.ID,
		CreatedAt: core.CurrentRFCTimestamp(),
	}
	// Reading community by ID
	var dbCommunities []adaptiveCommunity.AdaptiveCommunity
	dbCommunities, err2 = adaptiveCommunity.ReadOrEmpty(teamID.ToPlatformID(), communityID)(conn)
	var dbCommunity adaptiveCommunity.AdaptiveCommunity
	if len(dbCommunities) == 0 {
		logger.Infof("%s community not found. It's normal, we gonna create one", communityID)
		dbCommunity = models.AdaptiveCommunity{}
	}
	if err2 == nil {
		if dbCommunity.ID != "" {
			logger.Infof("%s community is already used up", communityID)
			// Selected community already exists, send a message back
			text := CommunityIsUsedErrorMessage(communityName)
			replyReplace(request, teamID, platform.MessageContent{Message: text})
		} else {
			// Create the community
			err2 = communityDAO.Create(comm)
			if err2 != nil {
				logger.WithField("namespace", namespace).WithField("error", err2).
					Errorf("Could not add entry to %s table", orgCommunitiesTable)
			} else {
				// Once a channel/group is subscribed to a community, get all existing users from the channel and add as community users
				// Adding existing channel members
				existingUsers := slackChannelMembers(channelID, teamID)
				logger.Infof("Existing members in %s channel for %s community: %v", channelID, teamID, existingUsers)

				setupCommunityUsers(channelID, communityID, existingUsers, teamID)
				// Checking if the selected value contains ":"
				// this happens for strategy communities where value is `community-type`:`parent-id`
				if strings.Contains(communityID, ":") {
					subscribedToStrategyCommunity(request, mc, teamID, channelID, communityID)
				} else {
					text := CommunitySuccessfullyAddedNotification(communityName)
					deleteOriginalMessage(request, teamID)
					postToChannel(teamID, platform.ConversationID(channelID), platform.MessageContent{Message: text})
				}
				postSubscriptionConfirmationToAdmin(teamID, communityID, request.User.ID)
			}
		}
	} else {
		logger.
			WithField("namespace", namespace).
			WithError(err2).
			Errorf("Error reading community with id %s from with platform %s", comm.ID, teamID)
	}
}

func setupCommunityUsers(channelID, communityID string, communityMemberIDs []string, teamID models.TeamID) {
	hasBeenSubscribedMany := isUserSubscribedToAnyCommunityMany(communityMemberIDs)
	userCommunities := addUsersToCommunity(teamID, channelID, communityID, communityMemberIDs)
	logger.Infof("Added %s users from %s channel to %s community in team %v", communityMemberIDs, channelID, communityID, teamID)
	welcomeAllUsers(teamID, userCommunities)
	// If the user has already subscribed to other channels,  we don't show adaptive scheduled time engagement
	for userID, hasBeenSubscribed := range hasBeenSubscribedMany {
		if !hasBeenSubscribed {
			setupUser(userID)
		}
	}
}

func subscribedToStrategyCommunity(request slack.InteractionCallback,
	mc models.MessageCallback,
	teamID models.TeamID, channelID string, communityID string) {
	communityName := ui.PlainText(communityID)

	// Get information about the objective community
	strategyCommunityID, typ, name := StrategyCommunityIdTypeName(string(communityName), teamID)
	// A channel has been created for a objective community. Update strategy communities with the same
	// Set channel_created and channel_id values
	err := updateStrategyCommunity(channelID, strategyCommunityID)
	if err == nil {
		logger.Infof("%s community is associated with Adaptive. Updated channel information.", strategyCommunityID)

		text := CommunityWithTypeSuccessfullyAddedNotification(ui.PlainText(name), ui.PlainText(typ))
		// Also add this as a user objective for the advocate so it can used to updates like coach-coachee
		// err = d.PutTableEntry(uObj, userObjectivesTable)
		// core.ErrorHandler(err, namespace, fmt.Sprintf("Could not write to %s table", userObjectivesTable))
		deleteOriginalMessage(request, teamID)
		postToChannel(teamID, platform.ConversationID(channelID), platform.MessageContent{Message: text})
		// We have now added feedback for a coaching engagement. We can now update the original engagement as answered.
		// TODO: This engagement is being shown as a notification. Update it into an attachment. Then it can be updated as answered
		commValueSplits := strings.Split(communityID, ":")

		commType := commValueSplits[0]
		commID := commValueSplits[1]

		postCommunityToStrategy(teamID, mc, commType, commID)
		postVisionIfExists(request, teamID, channelID)
	} else {
		logger.WithField("namespace", namespace).WithField("error", err).
			Errorf("Could not update channel_created flag in %s table", strategyCommunitiesTable)
	}
}

func postVisionIfExists(request slack.InteractionCallback, teamID models.TeamID, channelID string) {
	// Also, post vision statement, if it exists, to the new channel
	vision := strategy.StrategyVision(models.TeamID(teamID), strategyVisionTableName)
	if vision != nil {
		response := platform.Post(platform.ConversationID(channelID),
			platform.MessageContent{Message: VisionNotification(ui.RichText(vision.Vision))},
		)
		respond(teamID, response)
	}
}

func postToAdmin(teamID models.TeamID, text ui.RichText) {
	adminComm, err2 := communityDAO.ReadByID(teamID, string(community.Admin))
	if err2 != nil && strings.Contains(err2.Error(), "not found") {
		err2 = nil
		logger.Warnf("(1) No Admin Community found for platform: %s", teamID)
	}
	if adminComm.ChannelID == "" {
		logger.Warnf("(2) No Admin Community found for platform: %s", teamID)
	} else {
		response := platform.Post(platform.ConversationID(adminComm.ChannelID),
			platform.MessageContent{
				Message: text,
			})
		respond(teamID, response)
	}
}

func postToChannel(teamID models.TeamID, channelID platform.ConversationID, message platform.MessageContent) {
	response := platform.Post(channelID, message)
	respond(teamID, response)
}

func postSubscriptionConfirmationToAdmin(teamID models.TeamID, communityID, userID string) {
	// Publish a message only for non-admin channel
	if communityID != string(community.Admin) {
		postToAdmin(teamID, InvitedByUserToCommunityNotification(userID, communityTypeFromID(communityID)))
	}
}

func postSubscriptionRemovalToAdmin(teamID models.TeamID, communityID, userID string) {
	// Publish a message only for non-admin channel
	if communityID != string(community.Admin) {
		postToAdmin(teamID, RemovedByUserFromCommunityNotification(userID, communityTypeFromID(communityID)))
	}
}

// Unsubscribe menu item clicked
func onCommunityUnsubscribeClicked(request slack.InteractionCallback, teamID models.TeamID) platform.Response {
	channelID := request.Channel.ID
	fmt.Printf("Unsubscribing (platform=%s) from channel %s\n", teamID, channelID)
	commIDs := subscribedCommunityIDs(teamID, channelID)
	opts := liftStringToOption(simpleOptionStr)(commIDs)
	message := selectOptionsMessage(
		callback(channelID, "unsubscription", "select"),
		SelectCommunityToLeavePrompt,
		SelectCommunityMenuText,
		SelectCommunityFallbackMenuText,
		opts)
	return platform.OverrideByURL(
		platform.ResponseURLMessageID{ResponseURL: request.ResponseURL},
		message)
}

func onCommunityUnsubscribeCommunityClicked(
	request slack.InteractionCallback,
	communityID string,
	mc models.MessageCallback,
	teamID models.TeamID) (message platform.MessageContent) {
	err2 := channelUnsubscribe(request.Channel.ID, teamID)
	if err2 == nil {
		message = platform.MessageContent{Message: LeavingCommunityNotification(ui.PlainText(communityID))}
		// We have now added feedback for a coaching engagement. We can now update the original engagement as answered.
		utils.UpdateEngAsAnswered(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
		// Posting removal message to Admin
		postSubscriptionRemovalToAdmin(teamID, communityID, request.User.ID)
	} else {
		logger.WithError(err2).Errorf("Couldn't unsubscribe communityID=%s from teamID=%v", communityID, teamID)
		message = platform.MessageContent{Message: UnsubscriptionErrorMessage}
	}
	return
}

func communityTypeFromID(communityID string) string {
	splits := strings.Split(communityID, ":")
	return core.IfThenElse(len(splits) == 2, splits[0], communityID).(string)
}

func onCommunityUnsubscribeCancelled(request slack.InteractionCallback, mc models.MessageCallback) (message platform.MessageContent) {
	// User has decided to cancel this. This means, we remove this engagement
	// For now, we mark this as answered
	// TODO: We need to handle cases where a user ignores an engagement. This is different from not being reminded any more
	utils.UpdateEngAsAnswered(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
	return platform.MessageContent{Message: EngagementRemovedNotification}
}

func onMemberJoinedChannel(slackMsg slackevents.MemberJoinedChannelEvent, teamID models.TeamID) {
	logger.WithField("event", "member_joined_channel").Infof("%s user joined channel", slackMsg.User)
	// Member joined
	// Ensuring profile exists for the user. This could be the first time Adaptive is invited and no user yet exists.
	up, isAdaptiveBot, err := refreshUserCache(slackMsg.User, teamID)
	// Also refresh profile for the inviter as we  would be interacting with them immediately
	_, _, err2 := refreshUserCache(slackMsg.Inviter, teamID)
	if err == nil {
		logger.Infof("Newly joined user profile: %v", up)
		if err2 == nil {
			logger.Infof("Refreshed inviter profile: %v", up)
			if isAdaptiveBot {
				// Check if the member is a bot
				logger.Infof("Adaptive joined %s channel in %s platform", slackMsg.Channel, teamID)
				// We need to send inviter information to get subscribed communities
				onAdaptiveJoinedChannel(platform.ConversationID(slackMsg.Channel), teamID, slackMsg.Inviter)
			} else {
				logger.Infof("%s joined %s channel on invitation in %s platform", slackMsg.User, slackMsg.Channel, teamID)
				// If another user is added
				// Get the subscribed communities
				subComms := subscribedCommunities(teamID, slackMsg.Channel)
				hasBeenSubscribed := isUserSubscribedToAnyCommunity(slackMsg.User)
				userCommunityPairs := addUserToAllCommunities(teamID, slackMsg.User, subComms)
				logger.Infof("Welcoming newly added %s user", slackMsg.User)
				welcomeAllUsers(teamID, userCommunityPairs)
				if !hasBeenSubscribed {
					setupUser(slackMsg.User)
				}
			}
		} else {
			logger.WithField("namespace", namespace).WithField("error", err2).Errorf("Error refreshing user profile for %s", slackMsg.Inviter)
		}
	} else {
		logger.WithField("namespace", namespace).WithField("error", err).Errorf("Error refreshing user profile for %s", slackMsg.User)
	}
}

func isUserSubscribedToAnyCommunity(userID string) bool {
	comms, err2 := communityUserDAO.Read(userID)
	if err2 != nil && strings.Contains(err2.Error(), "not found") {
		logger.Infof("Not found community user %s", userID)
	}
	return err2 == nil && len(comms) > 0
}

func isUserSubscribedToAnyCommunityMany(userIDs []string) (m map[string]bool) {
	m = make(map[string]bool)
	for _, u := range userIDs {
		m[u] = isUserSubscribedToAnyCommunity(u)
	}
	return
}

// func getUserIDs(userCommunityPairs []models.AdaptiveCommunityUser3) (userIDs []string) {
// 	seen := make(map[string] bool)
// 	for _, each := range userCommunityPairs {
// 		if !seen[each.UserID] {
// 			userIDs = append(userIDs, each.UserID)
// 			seen[each.UserID] = true
// 		}
// 	}
// 	return
// }

func welcomeAllUsers(teamID models.TeamID, userCommunityPairs []models.AdaptiveCommunityUser3) {
	b1, _ := json.Marshal(userCommunityPairs)
	fmt.Printf("### welcomeAllUsers: %v", string(b1))
	for _, each := range userCommunityPairs {
		// Ensure the user profile exists
		err := addUserProfileForCommunityUser(each.UserID, models.ParseTeamID(each.PlatformID))
		if err == nil {
			welcomeUserToCommunity(teamID, each.UserID, community.AdaptiveCommunity(each.CommunityID))
		} else {
			log.Println(fmt.Sprintf("There was error with adding %s user: %v", each.UserID, err))
		}
	}
}

func setupUser(userID string) {
	userEngage := models.UserEngage{UserID: userID,
		IsNew: true, Update: false}
	invokeUserSetupLambdaUnsafe(userEngage)
}

func welcomeUserToCommunity(teamID models.TeamID, userID string, communityID community.AdaptiveCommunity) { // commUser models.AdaptiveCommunityUser3) {
	welcomeMessage := map[community.AdaptiveCommunity]ui.RichText{
		community.User:       UserCommunityWelcomeMessage,
		community.HR:         HRCommunityWelcomeMessage,
		community.Coaching:   CoachingCommunityWelcomeMessage,
		community.Admin:      AdminCommunityWelcomeMessage,
		community.Strategy:   StrategyCommunityWelcomeMessage,
		community.Competency: CompetencyCommunityWelcomeMessage,
		community.Capability: CapabilityCommunityWelcomeMessage,
		community.Initiative: InitiativeCommunityWelcomeMessage,
	}
	var commID community.AdaptiveCommunity
	splits := strings.Split(string(communityID), ":")
	if len(splits) == 2 {
		commID = community.AdaptiveCommunity(splits[0])
	} else {
		commID = communityID
	}
	message := WelcomeUserToCommunity(userID) + welcomeMessage[commID]

	directMessageToUser(teamID, userID, simpleMessage(message))
}

func onAdaptiveJoinedChannel(channelID platform.ConversationID, teamID models.TeamID, userID string) {
	// There is no user. We added new user. And post subscribe engagement
	message := getSubscribeMessage(channelID, teamID, userID)
	response := platform.PostEphemeral(userID, channelID, message)
	respond(teamID, response)
}

// A regular user is removed from the channel
func onMemberLeftChannel(slackMsg slack.MemberLeftChannelEvent) {
	err := communityUserDAO.DeleteUserFromCommunity(slackMsg.Channel, slackMsg.User)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not remove entry from %s table", communityUsersTable))
}

func onGroupLeftEvent(cbEvent slackevents.EventsAPICallbackEvent, teamID models.TeamID) {
	logger.Infof("Handling onGroupLeftEvent %v", *cbEvent.InnerEvent)
	var groupLeftEvent models.GroupLeftEvent
	err := json.Unmarshal(*cbEvent.InnerEvent, &groupLeftEvent)
	core.ErrorHandler(err, namespace, "Could not unmarshal raw json to GroupLeftEvent")

	if len(cbEvent.AuthedUsers) > 0 {
		authedUser := cbEvent.AuthedUsers[0]
		us, err2 := userDAO.Read(authedUser)
		if err2 != nil && strings.Contains(err2.Error(), "not found") {
			logger.Infof("Not found user %s", authedUser)
			err2 = nil
		}
		core.ErrorHandler(err2, namespace, "Error reading from users table")
		if us.IsAdaptiveBot {
			removeChannel(groupLeftEvent.ActorId, groupLeftEvent.Channel, teamID)
		} else {
			logger.Warnf("Weird onGroupLeftEvent (1) - %s (%s) not IsAdaptiveBot", authedUser, us.ID)
			err2 := communityUserDAO.DeleteUserFromCommunity(groupLeftEvent.Channel, authedUser)
			core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not remove entry from %s table", communityUsersTable))
		}
	} else {
		logger.Warnf("Weird onGroupLeftEvent (2) - AuthedUsers is empty")
	}

}

func getAllAvailableCommunitiesAsMenuOptions(channelID string, teamID models.TeamID, userID string) []ebm.MenuOption {
	availComms := liftStringToOption(simpleOptionStr)(availableCommunities(teamID))
	availStrComms := liftKvPairToOption(kvPairToMenuOption)(availableStrategyCommunities(teamID, userID))
	opts := append(availComms, availStrComms...)

	return opts
}

func getSubscribeMessage(channelID platform.ConversationID, teamID models.TeamID, userID string) (message platform.MessageContent) {
	mc := callback(string(channelID), "subscription", "select")
	availComms := liftStringToOption(simpleOptionStr)(availableCommunities(teamID))
	availStrComms := liftKvPairToOption(kvPairToMenuOption)(availableStrategyCommunities(teamID, userID))
	opts := append(availComms, availStrComms...)
	logger.Infof("Available communities for Adaptive to join: %s", opts)
	if len(opts) > 0 {
		message = selectOptionsMessage(mc,
			PostSubscribeEngagementTitle,
			SelectCommunityMenuText,
			SelectCommunityFallbackMenuText,
			opts)
		message.Message = InvitationToChannelAcknowledgement
	} else { // no communities left
		message = platform.MessageContent{Message: InvitationToChannelRejection}
	}
	return
}

func postCommunityToStrategy(teamID models.TeamID, mc models.MessageCallback,
	commType, commID string) {

	var attachs []ebm.Attachment
	switch commType {
	case string(community.Capability):
		capComm := strategy.CapabilityCommunityByID(models.TeamID(teamID), commID, capabilityCommunitiesTable)
		attachs = strategy.CapabilityCommunityViewAttachment(mc, &capComm, nil, false)
	case string(community.Initiative):
		initComm := strategy.InitiativeCommunityByID(models.TeamID(teamID), commID, strategyInitiativeCommunitiesTable)
		stratComm := StrategyCommunityByID(initComm.CapabilityCommunityID)
		attachs = strategy.InitiativeCommunityViewAttachmentReadOnly(mc, &initComm, nil, capabilityCommunitiesTable)
		// Also post the update to objective community
		response := platform.Post(platform.ConversationID(stratComm.ChannelID),
			platform.MessageContent{
				Message:     NotifyAboutNewAbilitiesInCommunityNotification(ui.PlainText(commType)),
				Attachments: attachs,
			})
		respond(teamID, response)
	default:
		fmt.Printf("Unknown strategy community: %s\n", commType)
		return
	}
	strategyComm, err2 := communityDAO.ReadByID(teamID, string(community.Strategy))
	if err2 != nil && strings.Contains(err2.Error(), "not found") {
		logger.Warnf("Not found strategy community")
		return
	}
	core.ErrorHandler(err2, namespace, "Error reading Strategy community")
	response := platform.Post(platform.ConversationID(strategyComm.ChannelID),
		platform.MessageContent{
			Message:     NotifyAboutNewAbilitiesInCommunityNotification(ui.PlainText(commType)),
			Attachments: attachs,
		})
	respond(teamID, response)

}
