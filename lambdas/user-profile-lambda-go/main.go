package lambda

import (
	"context"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core_utils_go "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/nlopes/slack"
	"log"
)

func HandleRequest(ctx context.Context, engage models.UserEngage) (uToken models.UserToken, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("error in user-profile-lambda %v", err2)
		}
	}()
	fmt.Printf("HandleRequest UserID='%s', PlatformID='%s'\n", engage.UserId, engage.PlatformID)
	uToken = models.UserToken{}
	// this is used for keeping this lambda warm
	// we send a request with an empty user every 30 min
	if engage.UserId == "" {
		return
	}
	platformID := engage.PlatformID
	profile, platformIDFromDB, err := readUserProfile(engage.UserId)
	if err != nil {
		err = wrapError(err, "Couldn't read user profile for "+engage.UserId)
		return
	}
	if profile.Id == "" {
		log.Printf("Cache missing for %s: %v\n", engage.UserId, err)
		profile, err = refreshUserCache(engage.UserId, platformID)
		if err != nil {
			err = wrapError(err, "Couldn't even refresh cache for user "+engage.UserId)
			return
		}
	}

	if platformID == "" {
		platformID = platformIDFromDB
	}
	platform, err2 := platformTokenDao.Read(platformID)
	if err2 != nil {
		err = wrapError(err2, "Couldn't query table "+confTable)
		return
	}
	uToken = models.UserToken{
		UserProfile:           profile,
		ClientPlatform:        models.ClientPlatform{PlatformName: platform.PlatformName, PlatformToken: platform.PlatformToken},
		ClientPlatformRequest: models.ClientPlatformRequest{Id: string(platform.PlatformID), Org: platform.Org},
	}
	uToken.ClientPlatformRequest.Id = string(platformID)

	return
}

func readUserProfile(userID string) (profile models.UserProfile, platformID models.PlatformID, err error) {
	user, err := userDao.Read(userID)
	profile = convertUserToProfile(user)
	platformID = user.PlatformID
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

func convertSlackUserToUser(user slack.User, platformID models.PlatformID) (mUser models.User) {
	return models.User{
		// UserProfile: models.UserProfile{
		ID:             user.ID,
		DisplayName:    user.RealName,
		FirstName:      user.Profile.FirstName,
		LastName:       user.Profile.LastName,
		Timezone:       user.TZ,
		TimezoneOffset: user.TZOffset,
		// },
		PlatformID: platformID,
		IsAdmin:    user.IsAdmin,
		Deleted:    user.Deleted,
		CreatedAt:  core_utils_go.CurrentRFCTimestamp(),
		IsShared:   false}
}

func wrapError(err error, name string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("{%s: %v}", name, err)
}

func refreshUserCache(userID string, platformID models.PlatformID) (profile models.UserProfile, err error) {
	if platformID == "" {
		panic("platformID is empty when querying " + userID)
	}
	platform, err := platformTokenDao.Read(platformID)
	if err == nil {
		api := slack.New(platform.PlatformToken)
		user, err2 := api.GetUserInfo(userID)
		err = err2
		mUser := convertSlackUserToUser(*user, platformID)
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
