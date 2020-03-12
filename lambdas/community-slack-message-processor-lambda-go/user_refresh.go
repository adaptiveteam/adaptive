package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/pkg/errors"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	// core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/nlopes/slack"
	"log"
)

func convertUserToProfile(user models.User) (profile models.UserProfile) {
	profile = models.UserProfile{
		Id:             user.ID,
		DisplayName:    user.DisplayName,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Timezone:       user.Timezone,
		TimezoneOffset: user.TimezoneOffset,
	}
	return
}

func readUserProfile(userID string) (profile models.UserProfile, teamID models.TeamID, err error) {
	var users []models.User
	users, err = userDAO.ReadOrEmpty(userID)
	if err == nil {
		user := models.User{}
		if len(users) > 0 {
			user = users[0]
		} else {
			logger.Infof("readUserProfile: Not found in users id=%s", userID)
		}
		profile = convertUserToProfile(user)
		teamID = models.ParseTeamID(user.PlatformID)
	}
	return
}

func refreshUserCache(userID string, teamID models.TeamID) (profile models.UserProfile, isAdaptiveBot bool, err error) {
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
				mUser := utilsUser.ConvertSlackUserToUser(*sUser, teamID)
				var previousUsers [] models.User
				previousUsers, err = userDAO.ReadOrEmpty(mUser.ID)
				if err == nil {
					for _, u := range previousUsers {
						mUser.CreatedAt = u.CreatedAt
						mUser.PlatformOrg = u.PlatformOrg
						mUser.IsAdmin = u.IsAdmin
						mUser.IsAdaptiveBot = u.IsAdaptiveBot
						mUser.AdaptiveScheduledTime = u.AdaptiveScheduledTime
						mUser.AdaptiveScheduledTimeInUTC = u.AdaptiveScheduledTimeInUTC
					}
					err = userDAO.CreateOrUpdate(mUser)
					logger.Infof("refreshUserCache: Created/updated user id=%s", mUser.ID)
						
					profile = convertUserToProfile(mUser)
					
					isAdaptiveBot = false // because mUser.IsAdaptiveBot is never initialized
				}
			}
		}
	}
	return
}

func addUserProfileForCommunityUser(userID string, teamID models.TeamID) (err error) {
	profile, _, err := readUserProfile(userID)
	if err == nil && profile.Id == "" {
		log.Printf("%s user not existing, adding now", userID)
		profile, _, err = refreshUserCache(userID, teamID)
	}
	return
}
