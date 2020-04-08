package lambda

import (
	"context"
	"fmt"
	"log"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

func HandleRequest(ctx context.Context, engage models.UserEngage) (uToken models.UserToken, err error) {
	defer core.RecoverAsLogError("profile-lambda")
	defer func() {
		if err != nil {
			log.Printf("Error in user-profile-lambda:%+v\n", err)
			err = nil
		}
	}()
	log.Printf("HandleRequest UserID='%s', TeamID='%s'\n", engage.UserID, engage.TeamID.ToString())
	uToken = models.UserToken{}
	// this is used for keeping this lambda warm
	// we send a request with an empty user every 30 min
	if engage.UserID == "" {
		return
	}
	teamID := engage.TeamID
	profile, teamIDFromDB, found, err := readUserProfile(engage.UserID)
	if err != nil {
		err = wrapError(err, "Couldn't read user profile for "+engage.UserID)
		return
	}
	if !found {
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
	token, err3 := platform.GetToken(teamID)(connGen.ForPlatformID(teamID.ToPlatformID()))
	if err3 != nil {
		err = wrapError(err3, "Couldn't query table "+confTable)
		return
	}
	uToken = models.UserToken{
		UserProfile:           profile,
		ClientPlatform:        models.ClientPlatform{PlatformName: models.SlackPlatform, PlatformToken: token},
		ClientPlatformRequest: models.ClientPlatformRequest{TeamID: teamID, Org: ""},
	}
	uToken.ClientPlatformRequest.TeamID = teamID

	return
}

func readUserProfile(userID string) (profile models.UserProfile, teamID models.TeamID, found bool, err error) {
	var users []models.User
	users, err = userDao.ReadOrEmpty(userID)
	found = len(users) > 0
	if found {
		profile = convertUserToProfile(users[0])
		teamID = models.ParseTeamID(users[0].PlatformID)
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

func refreshUserCache(userID string, teamID models.TeamID) (profile models.UserProfile, err error) {
	if teamID.IsEmpty() {
		panic(errors.New("refreshUserCache: teamID is empty when querying " + userID))
	}
	var token string
	token, err = platform.GetToken(teamID)(connGen.ForPlatformID(teamID.ToPlatformID()))
	if err == nil {
		api := slack.New(token)
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
