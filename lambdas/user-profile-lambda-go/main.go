package lambda

import (
	"context"
	"fmt"
	"log"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/user"
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
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	profile, teamIDFromDB, found, err := readUserProfile(engage.UserID)(conn)
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

func readUserProfile(userID string) func (conn daosCommon.DynamoDBConnection) (profile models.UserProfile, teamID models.TeamID, found bool, err error) {
	return func (conn daosCommon.DynamoDBConnection) (profile models.UserProfile, teamID models.TeamID, found bool, err error) {
		var users []models.User
		users, err = user.ReadOrEmpty(userID)(conn)
		found = len(users) > 0
		if found {
			profile = models.ConvertUserToProfile(users[0])
			teamID = models.ParseTeamID(users[0].PlatformID)
		}
		return
	}
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
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	token, err = platform.GetToken(teamID)(conn)
	if err == nil {
		api := slack.New(token)
		us, err2 := api.GetUserInfo(userID)
		err = err2
		var adaptiveBotID string
		adaptiveBotID, err = platform.GetAdaptiveBotIDOptional(conn)

		if err == nil {
			mUser := utilsUser.ConvertSlackUserToUser(*us, teamID, adaptiveBotID)
			err = user.Create(mUser)(conn)
			profile = models.ConvertUserToProfile(mUser)
		}
	}
	return
}
