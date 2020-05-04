package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"context"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	// mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/slack-go/slack"
	"sync"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/user"
)

func updateSlackUser(slackUser slack.User, event models.ClientPlatformRequest, 
	teamID models.TeamID,
	conn daosCommon.DynamoDBConnection) (err error) {
	now := core.CurrentRFCTimestamp()
	deactivatedAt := ""
	if slackUser.Deleted {
		deactivatedAt = now
	}
	item := models.User{
		ID:             slackUser.ID,
		DisplayName:    slackUser.RealName,
		FirstName:      slackUser.Profile.FirstName,
		LastName:       slackUser.Profile.LastName,
		Timezone:       slackUser.TZ,
		TimezoneOffset: slackUser.TZOffset,
		PlatformID:     event.TeamID.ToPlatformID(), 
		IsAdmin:        slackUser.IsAdmin,
		DeactivatedAt:  deactivatedAt,
		CreatedAt:      now,
		IsShared:       false}
	item.IsAdaptiveBot = slackUser.IsBot && slackUser.Profile.ApiAppID == teamID.ToString()

	// Check if the user already exists
	var users []models.User
	users, err = user.ReadOrEmpty(slackUser.ID)(conn)
	if err == nil {
		// Id not-empty meaning user exists
		for _, existingUser := range users {
			// Preserving the scheduled time
			item.AdaptiveScheduledTime = existingUser.AdaptiveScheduledTime
			item.AdaptiveScheduledTimeInUTC = existingUser.AdaptiveScheduledTimeInUTC
			item.CreatedAt = existingUser.CreatedAt
		}
		err = user.CreateOrUpdate(item)(conn)
	}
	return
}

func syncCommunityUserAsync(commUserID string, api *slack.Client,
	event models.ClientPlatformRequest, wg *sync.WaitGroup, ec chan error, 
	teamID models.TeamID,
	conn daosCommon.DynamoDBConnection) {
	defer wg.Done()
	// Get user info from Slack
	slackUser, err := api.GetUserInfo(commUserID)
	if err == nil {
		if slackUser != nil {
			if (!slackUser.IsBot && slackUser.Name != "slackbot") ||
				(slackUser.IsBot && slackUser.Profile.ApiAppID == teamID.ToString()) {
				err = updateSlackUser(*slackUser, event, teamID, conn)
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

func syncCommunityUserProfiles(users []string, api *slack.Client, event models.ClientPlatformRequest, 
	teamID models.TeamID,
	conn daosCommon.DynamoDBConnection) []error {
	// Set up a wait group and a channel to handle any errors
	wg := &sync.WaitGroup{}
	ec := make(chan error, len(users))

	for _, userID := range users {
		// Add adaptive user
		wg.Add(1)
		core.Go("syncCommunityUserAsync", func (){syncCommunityUserAsync(userID, api, event, wg, ec, teamID, conn)})
	}

	// Wait for all of the users to be added. 
	// NB! Blocking main thread.
	wg.Wait()
	// after all goroutines completed, we can close channel as no more errors can appear in it
	close(ec)
	return collectErrors(ec)
}

func deactivateUserAsync(userID string, wg *sync.WaitGroup, ec chan error,
	conn daosCommon.DynamoDBConnection) {
	defer wg.Done()
	logger.Infof("Deactivating user %s", userID)
	err := user.Deactivate(userID)(conn)
	if err != nil {
		logger.WithField("error", err).Errorf("Error deactivating %s user", userID)
		ec <- err
	}
}

func deactivateUsers(userIDs []string,
	conn daosCommon.DynamoDBConnection) []error {
	wg := &sync.WaitGroup{}
	ec := make(chan error, len(userIDs))

	for _, userID := range userIDs {
		wg.Add(1)
		core.Go("deactivateUserAsync", func(){ deactivateUserAsync(userID, wg, ec, conn)})
	}
	wg.Wait()
	close(ec)
	return collectErrors(ec)
}

func platformCommunities(teamID models.TeamID) (comms []models.AdaptiveCommunity, err error) {
	err = d.QueryTableWithIndex(userCommunityTable, awsutils.DynamoIndexExpression{
		IndexName: userCommunityPlatformIndex,
		Condition: "platform_id = :pi",
		Attributes: map[string]interface{}{
			":pi": teamID.ToString(),
		},
	}, map[string]string{}, true, -1, &comms)
	return
}

func readCommMemberIDs(commID string, teamID models.TeamID) (ids []string, err error) {
	defer core.RecoverToErrorVar("readCommMemberIDs", &err)
	// Get community members by querying community users table based on platform id and community id
	dbMembers := community.CommunityMembers(communityUsersTable, commID, teamID)
	for _, m := range dbMembers {
		ids = append(ids, m.UserId)
	}
	return
}

func obtainMemberIDsForCommunity(comm models.AdaptiveCommunity,
	api *slack.Client, teamID models.TeamID) (refreshIDs, removeIDs, addIDs []string, err error) {
	defer core.RecoverToErrorVar("synchronizeCommunity", &err)
	var slackMemberIDs, dbMemberIDs []string
	dbMemberIDs, err = readCommMemberIDs(comm.ID, teamID)
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
func createOrUpdateCommunityUser(comm models.AdaptiveCommunity, userID string) func (conn daosCommon.DynamoDBConnection) (err error) {
	logger.Infof("Adding user %s to community channelID=%s", userID, comm.ChannelID)
	acu := adaptiveCommunityUser.AdaptiveCommunityUser{
		ChannelID: comm.ChannelID,
		UserID: userID,
		PlatformID: comm.PlatformID,
		CommunityID: comm.ID,
	}
	return adaptiveCommunityUser.CreateOrUpdate(acu)
}

func removeCommunityUser(comm models.AdaptiveCommunity, userID string) func (conn daosCommon.DynamoDBConnection) (err error) {
	logger.Infof("Removing user %s from community channelID=%s", userID, comm.ChannelID)
	return adaptiveCommunityUser.Delete(comm.ChannelID, userID)
}

func allUserIDs(conn daosCommon.DynamoDBConnection)(ids []string, err error) {
	var allUsers []user.User
	allUsers, err = user.ReadByPlatformID(conn.PlatformID)(conn)
	for _, u := range allUsers {
		ids = append(ids, u.ID)
	}
	return
}
// HandleRequest handles request from user query lambda
func HandleRequest(ctx context.Context, event models.ClientPlatformRequest) {
	defer core.RecoverAsLogError("slack-user-query-lambda")
	var err error
	// Get all the user communities for the platform
	teamID := event.TeamID
	connGen := daosCommon.CreateConnectionGenFromEnv()
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	var cliPlatformToken string
	cliPlatformToken, err = platform.GetToken(teamID)(conn)
	if err == nil {
		logger.Infof("Retrieved token for org: %s", teamID.ToString())
		api := slack.New(cliPlatformToken)
		communities, err2 := platformCommunities(teamID)
		if err2 == nil {
			var allRefreshOrAddIDs, allRemoveIDs, allAddIDs []string
			for _, comm := range communities {
				refreshIDs, removeIDs, addIDs, err3 := obtainMemberIDsForCommunity(comm, api, teamID)
				if err3 == nil {
					allRefreshOrAddIDs = append(allRefreshOrAddIDs, refreshIDs ...)
					allRefreshOrAddIDs = append(allRefreshOrAddIDs, addIDs ...)
					allRemoveIDs = append(allRemoveIDs, removeIDs ...)
					allAddIDs = append(allAddIDs, addIDs ...)
					for _, id := range allAddIDs {
						err2 := createOrUpdateCommunityUser(comm, id)(conn)
						if(err2 != nil) {
							logger.Errorf("Couldn't add user %s to community %s: %+v", id, comm.ChannelID, err2)
						}
					}
					for _, id := range allRemoveIDs {
						err2 := removeCommunityUser(comm, id)(conn)
						if(err2 != nil) {
							logger.Errorf("Couldn't remove user %s from community %s: %+v", id, comm.ChannelID, err2)
						}
					}
				} else {
					logger.Errorf("Failure for channelID=%s: %+v", comm.ChannelID, err3)
				}
			}

			logger.Infof("Retrieved users from Slack for %s", teamID.ToString())

			distinctRefreshMemberIDs := core.Distinct(allRefreshOrAddIDs)
			logger.Infof("Synchronizing user profiles")

			allErrors := syncCommunityUserProfiles(distinctRefreshMemberIDs, api, event, teamID, conn)
			logger.Infof("Removing non-community members from users")

			ids, err4 := allUserIDs(conn)
			if err4 == nil {
				usersToRemove := core.InAButNotB(ids, distinctRefreshMemberIDs)
				errors2 := deactivateUsers(usersToRemove, conn)
				allErrors = append(allErrors, errors2...)
			}
			// if there is an error in the error channel, just return the first one
			if len(allErrors) == 0 {
				logger.Infof("Successfully updated/deleted user(s)")
			} else {
				logger.WithField("errors", allErrors).Errorf("Could not update/delete user(s)")
			}
		} 
	} 
	if err != nil {
		logger.Errorf("Could not query %s table on %s index: %+v", userCommunityTable, userCommunityPlatformIndex, err)
	}
	return
}
