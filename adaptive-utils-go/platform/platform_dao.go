package platform

import (
	"github.com/adaptiveteam/adaptive/core-utils-go"
	"log"

	"github.com/ReneKroon/ttlcache"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/adaptiveteam/adaptive/daos/clientPlatformToken"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/slackTeam"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func dynString(str string) (attr *dynamodb.AttributeValue) {
	return &dynamodb.AttributeValue{S: aws.String(str)}
}

var globalTokenCache *ttlcache.Cache

// GetToken retrieves the token from the cache or database.
func GetToken(teamID models.TeamID) func(common.DynamoDBConnection) (string, error) {
	return func(conn common.DynamoDBConnection) (token string, err error) {
		if globalTokenCache != nil {
			value, existsInCache := globalTokenCache.Get(teamID.ToString())
			if existsInCache {
				token = value.(string)
			}
		}
		if token == "" {
			if teamID.TeamID != "" {
				log.Printf("Reading SlackTeam %s\n", teamID.TeamID)
				dao := slackTeam.NewDAO(conn.Dynamo, "GetToken", conn.ClientID)
				var teams []slackTeam.SlackTeam
				teams, err = dao.ReadOrEmpty(teamID.TeamID)
				if err != nil {
					log.Printf("Failed to read SlackTeam %s: %v\n", teamID.TeamID, err)
					err = nil
				}
				if len(teams) > 0 {
					token = teams[0].BotAccessToken
				}
			}
			if token == "" {
				if teamID.AppID != "" {
					log.Printf("Reading by AppID %s\n", teamID.AppID)
					var teams []clientPlatformToken.ClientPlatformToken
					teams, err = clientPlatformToken.ReadOrEmpty(teamID.AppID)(conn)
					if err != nil {
						log.Printf("Failed to read ClientPlatformToken %s: %v\n", teamID.AppID, err)
						return
					}
					if len(teams) > 0 {
						token = teams[0].PlatformToken
					}
				}
			}
			if token != "" {
				globalTokenCache.Set(teamID.ToString(), token)
			}
		}
		return
	}
}

// GetAdaptiveBotIDOptional attempts to read BotUserID of the current SlackTeam
func GetAdaptiveBotIDOptional(conn common.DynamoDBConnection) (adaptiveBotID string, err error) {
	var teams []slackTeam.SlackTeam
	teams, err = slackTeam.ReadOrEmpty(conn.PlatformID)(conn)
	for _, team := range teams {
		adaptiveBotID = team.BotUserID
	}
	return
}

// GetTokenForUser searches token for the given user
func GetTokenForUser(dynamo *awsutils.DynamoRequest, clientID string, userID string) (token string, err error) {
	var teamID models.TeamID
	teamID, err = GetTeamIDForUser(dynamo, clientID, userID)
	conn := common.DynamoDBConnection{
		Dynamo:     dynamo,
		ClientID:   clientID,
		PlatformID: teamID.ToPlatformID(),
	}
	if err == nil {
		token, err = GetToken(teamID)(conn)
	}
	return
}

// GetTeamIDForUser -
// Deprecated: We should always provide team id because user id is not unique.
// see https://github.com/adaptiveteam/adaptive/issues/318
// and https://api.slack.com/methods/users.identity, https://stackoverflow.com/questions/39260512/slack-user-id-and-access-token-unique-across-teams-or-users
func GetTeamIDForUser(dynamo *awsutils.DynamoRequest, clientID string, userID string) (teamID models.TeamID, err error) {
	var connGen = common.DynamoDBConnectionGen{
		Dynamo: dynamo,
		TableNamePrefix: clientID,
	}

	// NB! below Read doesn't use platform id from connection at the moment.
	fakeConnWithArbitraryPlatformID := connGen.ForPlatformID("YET-UNKNOWN-PLATFORM-ID")
	var u user.User
	u, err = user.Read(userID)(fakeConnWithArbitraryPlatformID)
	if err == nil {
		teamID = models.ParseTeamID(u.PlatformID)
	}
	return
}

// GetConnectionForUserFromEnv reads environment variables
// and retrieves team id for the user
func GetConnectionForUserFromEnv(userID string) (conn common.DynamoDBConnection, err error) {
	connGen := common.CreateConnectionGenFromEnv()
	var teamID models.TeamID
	teamID, err = GetTeamIDForUser(connGen.Dynamo, connGen.TableNamePrefix, userID)
	if err == nil {
		conn = connGen.ForPlatformID(teamID.ToPlatformID())
	}
	return
}
// GetConnectionForUserFromEnvUnsafe reads environment variables
// and retrieves team id for the user
func GetConnectionForUserFromEnvUnsafe(userID string) (conn common.DynamoDBConnection) {
	var err error
	conn, err = GetConnectionForUserFromEnv(userID)
	core_utils_go.ErrorHandlerf(err, "Couldn't GetConnectionForUserFromEnvUnsafe(userID=%s", userID)
	return
}

var _ = func() int {
	globalTokenCache = InitLocalCache(globalTokenCache)
	return 0
}()
