package lambda

import (
	"context"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	// mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/nlopes/slack"
	"sync"
)

func addSlackUser(user slack.User, event models.ClientPlatformRequest, platformID models.PlatformID) (err error) {
	now := core.CurrentRFCTimestamp()
	deactivatedAt := ""
	if user.Deleted {
		deactivatedAt = now
	}
	item := models.User{
		// UserProfile: models.UserProfile{
		ID:             user.ID,
		DisplayName:    user.RealName,
		FirstName:      user.Profile.FirstName,
		LastName:       user.Profile.LastName,
		Timezone:       user.TZ,
		TimezoneOffset: user.TZOffset,
		// },
		PlatformID: event.PlatformID, 
		IsAdmin: user.IsAdmin, 
		DeactivatedAt: deactivatedAt,
		CreatedAt: now, IsShared: false}
	item.IsAdaptiveBot = user.IsBot && user.Profile.ApiAppID == string(platformID)

	// Check if the user already exists
	var users []models.User
	users, err = userDao.ReadOrEmpty(user.ID)
	if err == nil {
		// Id not-empty meaning user exists
		for _, existingUser := range users {
			// Preserving the scheduled time
			item.AdaptiveScheduledTime = existingUser.AdaptiveScheduledTime
			item.AdaptiveScheduledTimeInUTC = existingUser.AdaptiveScheduledTimeInUTC
		}
		err = userDao.CreateOrUpdate(item)
	}
	return
}

func syncCommunityUserAsync(commUserID string, api *slack.Client,
	event models.ClientPlatformRequest, wg *sync.WaitGroup, ec chan error, platformID models.PlatformID) {
	defer wg.Done()
	// Get user info from Slack
	slackUser, err := api.GetUserInfo(commUserID)
	if err == nil {
		if slackUser != nil {
			if (!slackUser.IsBot && slackUser.Name != "slackbot") ||
				(slackUser.IsBot && slackUser.Profile.ApiAppID == string(platformID)) {
				err = addSlackUser(*slackUser, event, platformID)
			}
			if err == nil {
				logger.Infof("Updated %s's information in the table", slackUser.ID)
			} else {
				logger.WithField("error", err).Errorf("Error adding user to table %s", commUserID)
				ec <- err
			}
		}
	} else {
		fmt.Printf("Error retrieving user from Slack %s: %v\n", commUserID, err)
		ec <- err
	}
}

func collectErrors(ec chan error) (errors []error) {
	for e := range ec {
		errors = append(errors, e)
	}
	return
}

func syncCommunityUserProfiles(users []string, api *slack.Client, event models.ClientPlatformRequest, platformID models.PlatformID) []error {
	// Set up a wait group and a channel to handle any errors
	wg := &sync.WaitGroup{}
	ec := make(chan error, len(users))

	for _, user := range users {
		// Add adaptive user
		wg.Add(1)
		go syncCommunityUserAsync(user, api, event, wg, ec, platformID)
	}

	// Wait for all of the users to be added. 
	// NB! Blocking main thread.
	wg.Wait()
	// after all goroutines completed, we can close channel as no more errors can appear in it
	close(ec)
	return collectErrors(ec)
}

func deactivateUserAsync(userID string, wg *sync.WaitGroup, ec chan error) {
	defer wg.Done()
	logger.Infof("Removing %s from users", userID)
	err := userDao.Deactivate(userID)
	if err != nil {
		logger.WithField("error", err).Errorf("Error deactivating %s user", userID)
		ec <- err
	}
}

func deactivateUsers(userIDs []string) []error {
	wg := &sync.WaitGroup{}
	ec := make(chan error, len(userIDs))

	for _, userID := range userIDs {
		wg.Add(1)
		go deactivateUserAsync(userID, wg, ec)
	}
	wg.Wait()
	close(ec)
	return collectErrors(ec)
}

func platformCommunities(platformID models.PlatformID) (comms []models.AdaptiveCommunity, err error) {
	err = d.QueryTableWithIndex(userCommunityTable, awsutils.DynamoIndexExpression{
		IndexName: userCommunityPlatformIndex,
		Condition: "platform_id = :pi",
		Attributes: map[string]interface{}{
			":pi": platformID,
		},
	}, map[string]string{}, true, -1, &comms)
	return
}

func readCommMemberIDs(commID string, platformID models.PlatformID) (ids []string, err error) {
	defer core.RecoverToErrorVar("readCommMemberIDs", &err)
	// Get community members by querying community users table based on platform id and community id
	dbMembers := community.CommunityMembers(communityUsersTable, commID, platformID, communityUsersCommunityIndex)
	for _, m := range dbMembers {
		ids = append(ids, m.UserId)
	}
	return
}

func obtainMemberIDsForCommunity(comm models.AdaptiveCommunity,
	api *slack.Client, platformID models.PlatformID) (refreshIDs, removeIDs, addIDs []string, err error) {
	defer core.RecoverToErrorVar("synchronizeCommunity", &err)
	var slackMemberIDs, dbMemberIDs []string
	dbMemberIDs, err = readCommMemberIDs(comm.ID, platformID)
	slackMemberIDs, _,  err = api.GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: comm.ChannelID,
		Cursor: "",
		Limit: 999,
	})//, comm.ChannelID)
	refreshIDs = core.InAAndB   (dbMemberIDs, slackMemberIDs)
	removeIDs  = core.InAButNotB(dbMemberIDs, slackMemberIDs)
	addIDs     = core.InAButNotB(slackMemberIDs, dbMemberIDs)
	return
}
func addCommunityUser(comm models.AdaptiveCommunity, userID string) (err error) {
	acu := adaptiveCommunityUser.AdaptiveCommunityUser{
		ChannelID: comm.ChannelID,
		UserID: userID,
		PlatformID: comm.PlatformID,
		CommunityID: comm.ID,
	}
	err = adaptiveCommunityUserDAO.CreateOrUpdate(acu)
	return
}
func removeCommunityUser(comm models.AdaptiveCommunity, userID string) (err error) {
	return adaptiveCommunityUserDAO.Delete(comm.ChannelID, userID)
}
func allUserIDs(platformID models.PlatformID)(ids []string, err error) {
	allUsers, err := userDao.ReadByPlatformID(platformID)
	for _, u := range allUsers {
		ids = append(ids, u.ID)
	}
	return
}
// HandleRequest handles request from user query lambda
func HandleRequest(ctx context.Context, event models.ClientPlatformRequest) {
	defer core.RecoverAsLogError("slack-user-query-lambda")
	var allRefreshOrAddIDs, allRemoveIDs, allAddIDs []string
	// Get all the user communities for the platform
	platformID := event.PlatformID
	cliPlatformToken := platformTokenDao.GetPlatformTokenUnsafe(event.PlatformID)
	logger.Infof("Retrieved token for org: %s", event.PlatformID)
	api := slack.New(cliPlatformToken)
	communities, err2 := platformCommunities(platformID)
	if err2 == nil {
		for _, comm := range communities {
			refreshIDs, removeIDs, addIDs, err3 := obtainMemberIDsForCommunity(comm, api, platformID)
			if err3 == nil {
				allRefreshOrAddIDs = append(allRefreshOrAddIDs, refreshIDs ...)
				allRefreshOrAddIDs = append(allRefreshOrAddIDs, addIDs ...)
				allRemoveIDs = append(allRemoveIDs, removeIDs ...)
				allAddIDs = append(allAddIDs, addIDs ...)
				for _, id := range allAddIDs {
					err2 := addCommunityUser(comm, id)
					if(err2 != nil) {
						logger.Errorf("Couldn't add user %s to community %s: %+v", id, comm.ChannelID, err2)
					}
				}
				for _, id := range allRemoveIDs {
					err2 := removeCommunityUser(comm, id)
					if(err2 != nil) {
						logger.Errorf("Couldn't remove user %s from community %s: %+v", id, comm.ChannelID, err2)
					}
				}
			} else {
				logger.Errorf("Failure for channelID=%s: %+v", comm.ChannelID, err3)
			}
		}

		logger.Infof("Retrieved users from Slack for %s", platformID)

		distinctRefreshMemberIDs := core.Distinct(allRefreshOrAddIDs)
		logger.Infof("Synchronizing user profiles")

		allErrors := syncCommunityUserProfiles(distinctRefreshMemberIDs, api, event, platformID)
		logger.Infof("Removing non-community members from users")

		ids, err4 := allUserIDs(event.PlatformID)
		if err4 == nil {
			usersToRemove := core.InAButNotB(ids, distinctRefreshMemberIDs)
			errors2 := deactivateUsers(usersToRemove)
			allErrors = append(allErrors, errors2...)
		}
		// if there is an error in the error channel, just return the first one
		if len(allErrors) == 0 {
			logger.Infof("Successfully updated/deleted user(s)")
		} else {
			logger.WithField("errors", allErrors).Errorf("Could not update/delete user(s)")
		}
	} else {
		logger.Errorf("Could not query %s table on %s index: %+v", userCommunityTable, userCommunityPlatformIndex, err2)
	}
	return
}
