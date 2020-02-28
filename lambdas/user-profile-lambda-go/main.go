package lambda

import (
	"github.com/pkg/errors"
	"context"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	// daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/nlopes/slack"
	"log"
)

func HandleRequest(ctx context.Context, engage models.UserEngage) (uToken models.UserToken, err error) {
	defer core.RecoverToErrorVar("profile-lambda", &err)
	defer func() {
		if err != nil {
			log.Printf("Error in user-profile-lambda:%v\n", err)
		}
	}()
	log.Printf("HandleRequest UserID='%s', PlatformID='%s'\n", engage.UserId, engage.PlatformID)
	uToken = models.UserToken{}
	// this is used for keeping this lambda warm
	// we send a request with an empty user every 30 min
	if engage.UserId == "" {
		return
	}
	platformID := engage.PlatformID
	profile, platformIDFromDB, found, err2 := readUserProfile(engage.UserId)
	if err2 != nil {
		err = wrapError(err, "Couldn't read user profile for "+engage.UserId)
		return
	}
	if !found {
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
	platform, found, err3 := platformTokenDao.Read(platformID)
	if err3 != nil {
		err = wrapError(err3, "Couldn't query table "+confTable)
		return
	}
	if !found {
		err = fmt.Errorf("not found platformID=%s",platformID)
	}
	uToken = models.UserToken{
		UserProfile:           profile,
		ClientPlatform:        models.ClientPlatform{PlatformName: platform.PlatformName, PlatformToken: platform.PlatformToken},
		ClientPlatformRequest: models.ClientPlatformRequest{PlatformID: platform.PlatformID, Org: platform.Org},
	}
	uToken.ClientPlatformRequest.PlatformID = platformID

	return
}

func readUserProfile(userID string) (profile models.UserProfile, platformID models.PlatformID, found bool, err error) {
	var users []models.User
	users, err = userDao.ReadOrEmpty(userID)
	found = len(users)>0
	if found {
		profile = convertUserToProfile(users[0])
		platformID = users[0].PlatformID
	}
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

func refreshUserCache(userID string, platformID models.PlatformID) (profile models.UserProfile, err error) {
	if platformID == "" {
		err = errors.New("refreshUserCache: teamID is empty when querying " + userID)
		return
	}
	platform, found, err2 := platformTokenDao.Read(platformID)
	err = err2
	if !found {
		err = fmt.Errorf("Token not found for platformID=%s", platformID)
	}
	if err == nil {
		api := slack.New(platform.PlatformToken)
		user, err3 := api.GetUserInfo(userID)
		err = err3
		mUser := utilsUser.ConvertSlackUserToUser(*user, platformID)
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
