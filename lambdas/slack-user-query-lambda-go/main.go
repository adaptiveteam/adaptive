package lambda

import (
	"context"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/nlopes/slack"
	"sync"
)

func addSlackUser(user slack.User, event models.ClientPlatformRequest, platformID models.PlatformID) (err error) {
	item := models.User{
		// UserProfile: models.UserProfile{
		ID:             user.ID,
		DisplayName:    user.RealName,
		FirstName:      user.Profile.FirstName,
		LastName:       user.Profile.LastName,
		Timezone:       user.TZ,
		TimezoneOffset: user.TZOffset,
		// },
		PlatformID: models.PlatformID(event.Id), 
		PlatformOrg: event.Org, IsAdmin: user.IsAdmin, Deleted: user.Deleted,
		CreatedAt: core.CurrentRFCTimestamp(), IsShared: false}
	item.IsAdaptiveBot = user.IsBot && user.Profile.ApiAppID == string(platformID)

	// Check if the user already exists
	existingUser, err := userDao.Read(user.ID)
	if err == nil {
		// Id not-empty meaning user exists
		if existingUser.ID == "" {
			err = userDao.Create(item)
		} else {
			// Preserving the scheduled time
			item.AdaptiveScheduledTime = existingUser.AdaptiveScheduledTime
			item.AdaptiveScheduledTimeInUTC = existingUser.AdaptiveScheduledTimeInUTC
			err = userDao.Update(item)
		}
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

func removeUserAsync(communityUsersIDs []string, userID string, wg *sync.WaitGroup, ec chan error) {
	defer wg.Done()
	// If the user is not part of any community, delete the user
	if !core.ListContainsString(communityUsersIDs, userID) {
		logger.Infof("Removing %s from users", userID)
		err := userDao.Delete(userID)
		if err != nil {
			logger.WithField("error", err).Errorf("Error removing %s user", userID)
			ec <- err
		}
	}
}

func removeNonCommunityUsers(communityUserIDs []string, platformID models.PlatformID) []error {
	allUsers := userDao.ReadByPlatformIDUnsafe(platformID)
	wg := &sync.WaitGroup{}
	ec := make(chan error, len(allUsers))

	for _, each := range allUsers {
		wg.Add(1)
		go removeUserAsync(communityUserIDs, each.ID, wg, ec)
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

func HandleRequest(ctx context.Context, event models.ClientPlatformRequest) {
	var allCommunitiesMemberIDs []string
	// Get all the user communities for the platform
	platformID := models.PlatformID(event.Id)
	communities, err := platformCommunities(platformID)
	if err == nil {
		for _, each := range communities {
			// Get community members by querying community users table based on platform id and community id
			members := community.CommunityMembers(communityUsersTable, each.ID, event.Id, communityUsersCommunityIndex)
			for _, member := range members {
				allCommunitiesMemberIDs = append(allCommunitiesMemberIDs, member.UserId)
			}
		}

		cliPlatformToken := platformTokenDao.GetPlatformTokenUnsafe(models.PlatformID(event.Id))
		logger.Infof("Retrieved token for org: %s", event.Id)

		api := slack.New(cliPlatformToken)
		logger.Infof("Retrieved users from Slack for %s", event.Org)

		distinctMemberIDs := core.Distinct(allCommunitiesMemberIDs)
		logger.Infof("Synchronizing user profiles")

		errors1 := syncCommunityUserProfiles(distinctMemberIDs, api, event, platformID)
		logger.Infof("Removing non-community members from users")

		errors2 := removeNonCommunityUsers(distinctMemberIDs, models.PlatformID(event.Id))
		allErrors := append(errors1, errors2...)

		// if there is an error in the error channel, just return the first one
		if len(allErrors) == 0 {
			logger.Infof("Successfully updated/deleted user(s)")
		} else {
			logger.WithField("errors", allErrors).Errorf("Could not update/delete user(s)")
		}
	} else {
		logger.Errorf("Could not query %s table on %s index", userCommunityTable, userCommunityPlatformIndex)
	}
	return
}
