package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/user"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/pkg/errors"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	// core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/slack-go/slack"
	"log"
)

func readUserProfile(userID string, conn common.DynamoDBConnection) (profile models.UserProfile, teamID models.TeamID, err error) {
	var users []models.User
	users, err = user.ReadOrEmpty(conn.PlatformID, userID)(conn)
	if err == nil {
		user := models.User{}
		if len(users) > 0 {
			user = users[0]
		} else {
			logger.Infof("readUserProfile: Not found in users id=%s", userID)
		}
		profile = models.ConvertUserToProfile(user)
		teamID = models.ParseTeamID(user.PlatformID)
	}
	return
}

func refreshUserCache(userID string, conn common.DynamoDBConnection) (profile models.UserProfile, isAdaptiveBot bool, err error) {
	teamID := models.ParseTeamID(conn.PlatformID)
	if teamID.IsEmpty() {
		err = errors.New("teamID is empty when refreshing user " + userID)
	} else {
		conn := connGen.ForPlatformID(teamID.ToPlatformID())
		token, err2 := platform.GetToken(teamID)(conn)

		err = err2
		if err == nil {
			api := slack.New(token)
			var sUser *slack.User
			sUser, err = api.GetUserInfo(userID)
			if err == nil {
				var adaptiveBotID string
				adaptiveBotID, err = platform.GetAdaptiveBotIDOptional(conn)
				if err == nil {
					mUser := utilsUser.ConvertSlackUserToUser(*sUser, teamID, adaptiveBotID)
					var previousUsers [] models.User
					previousUsers, err = user.ReadOrEmpty(conn.PlatformID, mUser.ID)(conn)
					if err == nil {
						for _, u := range previousUsers {
							mUser.CreatedAt = u.CreatedAt
							mUser.PlatformOrg = u.PlatformOrg
							mUser.IsAdmin = u.IsAdmin
							mUser.IsAdaptiveBot = u.IsAdaptiveBot
							mUser.AdaptiveScheduledTime = u.AdaptiveScheduledTime
							mUser.AdaptiveScheduledTimeInUTC = u.AdaptiveScheduledTimeInUTC
						}
						err = user.CreateOrUpdate(mUser)(conn)
						logger.Infof("refreshUserCache: Created/updated user id=%s", mUser.ID)
						
						profile = models.ConvertUserToProfile(mUser)
					
						isAdaptiveBot = false // because mUser.IsAdaptiveBot is never initialized
					}
				}
			}
		}
	}
	return
}

func addUserProfileForCommunityUser(userID string, conn common.DynamoDBConnection) (err error) {
	profile, _, err := readUserProfile(userID, conn)
	if err == nil && profile.Id == "" {
		log.Printf("%s user not existing, adding now", userID)
		profile, _, err = refreshUserCache(userID, conn)
	}
	return
}
