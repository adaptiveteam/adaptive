package lambda

import (
	"context"
	"fmt"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	"github.com/pkg/errors"

	// core "github.com/adaptiveteam/adaptive/core-utils-go"
	"log"

	"github.com/nlopes/slack"
)

func HandleRequest(ctx context.Context, engage models.UserEngage) (uToken models.UserToken, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("error in user-profile-lambda %v", err2)
		}
	}()
	fmt.Printf("HandleRequest UserID='%s', PlatformID='%s'\n", engage.UserID, engage.TeamID)
	uToken = models.UserToken{}
	// this is used for keeping this lambda warm
	// we send a request with an empty user every 30 min
	if engage.UserID == "" {
		return
	}
	teamID := engage.TeamID
	profile, teamIDFromDB, err := readUserProfile(engage.UserID)
	if err != nil {
		err = wrapError(err, "Couldn't read user profile for "+engage.UserID)
		return
	}
	if profile.Id == "" {
		log.Printf("Cache missing for %s: %v\n", engage.UserID, err)
		profile, err = refreshUserCache(engage.UserID, teamID)
		if err != nil {
			err = wrapError(err, "Couldn't even refresh cache for user "+engage.UserID)
			return
		}
	}

	if teamID.IsEmpty() {
		teamID = teamIDFromDB
	}
	platform, err2 := platformTokenDao.Read(teamID)
	if err2 != nil {
		err = wrapError(err2, "Couldn't query table "+confTable)
		return
	}
	uToken = models.UserToken{
		UserProfile:           profile,
		ClientPlatform:        models.ClientPlatform{PlatformName: platform.PlatformName, PlatformToken: platform.PlatformToken},
		ClientPlatformRequest: models.ClientPlatformRequest{TeamID: models.ParseTeamID(platform.PlatformID), Org: platform.Org},
	}
	uToken.ClientPlatformRequest.TeamID = teamID

	return
}

func readUserProfile(userID string) (profile models.UserProfile, teamID models.TeamID, err error) {
	user, err := userDao.Read(userID)
	profile = convertUserToProfile(user)
	teamID = models.ParseTeamID(user.PlatformID)
	return
}

func convertUserToProfile(user models.User) (profile models.UserProfile) {
	profile = models.UserProfile{
		Id:          user.ID,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Timezone:    user.Timezone,
	}
	return
}

func wrapError(err error, name string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("{%s: %v}", name, err)
}

func refreshUserCache(userID string, teamID models.TeamID) (profile models.UserProfile, err error) {
	if teamID.IsEmpty() {
		panic(errors.New("refreshUserCache: teamID is empty when querying " + userID))
	}
	platform, err := platformTokenDao.Read(teamID)
	if err == nil {
		api := slack.New(platform.PlatformToken)
		user, err2 := api.GetUserInfo(userID)
		err = err2
		mUser := utilsUser.ConvertSlackUserToUser(*user, teamID)
		err = userDao.Create(mUser)
		profile = models.UserProfile{ //mUser.UserProfile
			Id:             mUser.ID,
			DisplayName:    mUser.DisplayName,
			FirstName:      mUser.FirstName,
			LastName:       mUser.LastName,
			Timezone:       mUser.Timezone,
			TimezoneOffset: mUser.TimezoneOffset,
		}
	}
	return
}
