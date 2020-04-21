package lambda

import (
	"github.com/pkg/errors"
	"fmt"
	"strings"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"

	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"

	"github.com/aws/aws-lambda-go/events"
	"golang.org/x/net/context"

	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/slackTeam"
	oauth2 "github.com/adaptiveteam/adaptive/oauth2-reimpl"
	// "golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type CODE string

const StateValue = "BE9DE56F"


var (
	// AdaptiveScopes - Slack authentication scopes. You can obtain the
	// auth URL from
	// https://slack.com/oauth/v2/authorize?client_id=436528929141.622333620304
	// &scope=app_mentions:read,calls:read,calls:write,channels:read,chat:write,dnd:read,files:read,files:write,groups:history,groups:read,groups:write,im:history,im:read,im:write,incoming-webhook,links:write,mpim:history,mpim:read,mpim:write,pins:write,reactions:read,reactions:write,remote_files:read,remote_files:share,remote_files:write,team:read,usergroups:read,usergroups:write,users:read,users:read.email,users:write
	// &user_scope=channels:read,chat:write,files:write,groups:read,groups:write,im:history,users:read

	// https://slack.com/oauth/authorize?access_type=offline&client_id=436528929141.622333620304&response_type=code
	// &scope=calls%3Aread+calls%3Awrite+channels%3Aread+dnd%3Aread+files%3Aread+groups%3Ahistory+groups%3Aread+groups%3Awrite+im%3Ahistory+im%3Aread+im%3Awrite+incoming-webhook+links%3Awrite+mpim%3Ahistory+mpim%3Aread+mpim%3Awrite+pins%3Awrite+reactions%3Aread+reactions%3Awrite+remote_files%3Aread+remote_files%3Ashare+remote_files%3Awrite+team%3Aread+usergroups%3Aread+usergroups%3Awrite+users%3Aread+users%3Aread.email+users%3Awrite
	// &user_scope=channels%3Aread+chat%3Awrite+files%3Awrite+groups%3Aread+groups%3Awrite+im%3Ahistory+users%3Aread&state=state
	// incoming-webhook,
	AdaptiveScopes = "app_mentions:read,calls:read,calls:write,channels:read,chat:write,chat:write.customize,dnd:read,files:read,files:write,groups:history,groups:read,groups:write,im:history,im:read,im:write,links:write,mpim:history,mpim:read,mpim:write,pins:write,reactions:read,reactions:write,team:read,usergroups:read,usergroups:write,users:read,users:read.email,users:write"
	// remote_files:read,remote_files:share,remote_files:write,
	// incoming-webhook - this requires a channel where bot will be posting messages from outside
	UserScopes = "" // "channels:read,chat:write,files:write,groups:read,groups:write,im:history,users:read"
	// https://1vvtp0yc61.execute-api.us-east-1.amazonaws.com/dev/%7Bproxy+%7D?code=436528929141.946109056755.2599abf6300120ff5bc70e3b20a6680e64e05d652dc046b72c6c4d24a9f73c69&state=state
)
//https://slack.com/oauth/v2/authorize?access_type=offline&client_id=436528929141.622333620304&response_type=code&scope=chat:write+chat:write.customize+calls%3Aread+calls%3Awrite+channels%3Aread+dnd%3Aread+files%3Aread+groups%3Ahistory+groups%3Aread+groups%3Awrite+im%3Ahistory+im%3Aread+im%3Awrite+links%3Awrite+mpim%3Ahistory+mpim%3Aread+mpim%3Awrite+pins%3Awrite+reactions%3Aread+reactions%3Awrite+remote_files%3Aread+remote_files%3Ashare+remote_files%3Awrite+team%3Aread+usergroups%3Aread+usergroups%3Awrite+users%3Aread+users%3Aread.email+users%3Awrite&state=BE9DE56F&user_scope=
// SlackOAuthConfig - constructs OAuth2 config with Slack endpoint and AdaptiveScopes
func SlackOAuthConfig(clientID, clientSecret string) oauth2.Config {
	return oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       strings.Split(AdaptiveScopes, ","),

		Endpoint:     oauth2.Endpoint{
			AuthStyle:   oauth2.AuthStyle(endpoints.Slack.AuthStyle),
			AuthURL:     "https://slack.com/oauth/v2/authorize",//endpoints.Slack.AuthURL,
			TokenURL:    "https://slack.com/api/oauth.v2.access",
		},
	}
}

// GlobalSlackOAuthConfig reads environment variables to create configuration for Slack
func GlobalSlackOAuthConfig() oauth2.Config {
	slackClientID := utils.NonEmptyEnv("SLACK_CLIENT_ID")
	slackClientSecret := utils.NonEmptyEnv("SLACK_CLIENT_SECRET")
	return SlackOAuthConfig(slackClientID, slackClientSecret)
}

func generateAuthorizationURL(conf oauth2.Config) (url string) {
	userScope := oauth2.SetAuthURLParam("user_scope", UserScopes) // strings.Join(UserScopes, ","))
	url = conf.AuthCodeURL(StateValue, oauth2.AccessTypeOffline, userScope)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
	return url
}

// ExchangeCodeForAuthenticationToken -
// Use the authorization code that is pushed to the redirect
// URL. Exchange will do the handshake to retrieve the
// initial access token.
func ExchangeCodeForAuthenticationToken(conf oauth2.Config, code CODE) (token oauth2.Token, err error) {
	ctx := context.Background()

	var tokenRef *oauth2.Token
	tokenRef, err = conf.Exchange(ctx, string(code))
	if err == nil {
		token = *tokenRef
	}
	return
}

// GenerateAddToSlackURL - handle Slack platform request to
func GenerateAddToSlackURL(userID, channelID string, teamID models.TeamID) {
	conn := globalConnection(teamID)
	slackAdapter := mapper.SlackAdapterForTeamID(conn)
	url := generateAuthorizationURL(GlobalSlackOAuthConfig())
	slackAdapter.PostSyncUnsafe(platform.Post(platform.ConversationID(userID),
		platform.Message(ui.Sprintf(
			"Send the following link to a user wishing to add Adaptive to their Slack workspace: %s", url)),
	))
}

/*
{
    "ok": true,
    "access_token": "xoxb-17653672481-19874698323-pdFZKVeTuE8sk7oOcBrzbqgy",
    "token_type": "bot",
    "scope": "commands,incoming-webhook",
    "bot_user_id": "U0KRQLJ9H",
    "app_id": "A0KRD7HC3",
    "team": {
        "name": "Slack Softball Team",
        "id": "T9TK3CUKW"
    },
    "enterprise": {
        "name": "slack-sports",
        "id": "E12345678"
    },
    "authed_user": {
        "id": "U1234",
        "scope": "chat:write",
        "access_token": "xoxp-1234",
        "token_type": "user"
    }
}

Instead Slack sends:
{
    "ok": true,
    "access_token": "xoxp-436528929141-612915016337-610914894754-a293bbc150eba8d1a2c8884cccd2beff",
    "scope": "identify,bot,incoming-webhook,groups:history,im:history,mpim:history,channels:read,files:read,groups:read,im:read,mpim:read,reactions:read,team:read,users:read,users:read.email,usergroups:read,dnd:read,chat:write,files:write,groups:write,im:write,mpim:write,reactions:write,users:write,pins:write,usergroups:write,links:write,remote_files:write,remote_files:share,remote_files:read,calls:write,calls:read",
    "user_id": "UJ0SX0G9X",
    "team_id": "TCUFJTB45",
    "enterprise_id": null,
    "team_name": "Adaptive.Team",
    "incoming_webhook": {
        "channel": "@adaptive-dev-ivan",
        "channel_id": "DJAA34V5W",
        "configuration_url": "https://adaptive--team.slack.com/services/BUBQGHEQ7",
        "url": "https://hooks.slack.com/services/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"
	}

		//    "user_id": "UJ0SX0G9X",
    // "team_id": "TCUFJTB45",
    // "enterprise_id": null,
    // "team_name": "Adaptive.Team",

}


They promise to send (https://api.slack.com/docs/oauth)
{
    "access_token": "xoxp-XXXXXXXX-XXXXXXXX-XXXXX",
    "scope": "incoming-webhook,commands,bot",
    "team_name": "Team Installing Your Hook",
    "team_id": "XXXXXXXXXX",
    "incoming_webhook": {
        "url": "https://hooks.slack.com/TXXXXX/BXXXXX/XXXXXXXXXX",
        "channel": "#channel-it-will-post-to",
        "configuration_url": "https://teamname.slack.com/services/BXXXXX"
    },
    "bot":{
        "bot_user_id":"UTTTTTTTTTTR",
        "bot_access_token":"xoxb-XXXXXXXXXXXX-TTTTTTTTTTTTTT"
    }
}
Indeed they send:
{
    "ok": true,
    "access_token": "xoxp-4*",
    "scope": "identify,bot,incoming-webhook,groups:history,im:history,mpim:history,channels:read,files:read,groups:read,im:read,mpim:read,reactions:read,team:read,users:read,users:read.email,usergroups:read,dnd:read,chat:write,files:write,groups:write,im:write,mpim:write,reactions:write,users:write,pins:write,usergroups:write,links:write,remote_files:write,remote_files:share,remote_files:read,calls:write,calls:read",
    "user_id": "U*",
    "team_id": "T*",
    "enterprise_id": null,
    "team_name": "Adaptive.Team",
    "bot": {
        "bot_user_id": "U*",
        "bot_access_token": "xoxb-4*"
    }
}

The new version of Slack answer (2020-04-20):

Body:
{
    "ok": true,
    "app_id": "AJA9TJ88Y",
    "authed_user": {
        "id": "UP0NQ8B9V"
    },
    "scope": "chat:write,chat:write.customize,calls:read,calls:write,channels:read,dnd:read,files:read,groups:history,groups:read,groups:write,im:history,im:read,im:write,links:write,mpim:history,mpim:read,mpim:write,pins:write,reactions:read,reactions:write,remote_files:read,remote_files:share,remote_files:write,team:read,usergroups:read,usergroups:write,users:read,users:read.email,users:write",
    "token_type": "bot",
    "access_token": "xoxb-7...",
    "bot_user_id": "UV4CAUA9Y",
    "team": {
        "id": "TLN1E8P6V",
        "name": "adaptive-test"
    },
    "enterprise": null
}
TODO: when user denies, Slack sends http://yourapp.com/oauth?error=access_denied as described here: https://api.slack.com/docs/oauth
*/

// HandleRedirectURLGetRequest processes GET request from Slack for OAuth2
func HandleRedirectURLGetRequest(conn common.DynamoDBConnection, request events.APIGatewayProxyRequest) (err error) {
	logger.
		WithField("Path", request.Path).
		WithField("Body", request.Body).
		WithField("HTTPMethod", request.HTTPMethod).
		WithField("PathParameters", request.PathParameters).
		WithField("QueryStringParameters", request.QueryStringParameters).
		Infof("HandleRedirectURLGetRequest")
	errText, hasError := request.QueryStringParameters["error"]

	if hasError {
		if errText == "access_denied" { // 
			logger.Warningf("teamWARN User has denied our application access to their workspace")
			// we just exit with err = nil. Outside of this function there will be a redirect
		}
		err = errors.Errorf("Outer web hook has been invoked (GET) (presumably by Slack) with error message: %s", errText)
	} else {
		conf := GlobalSlackOAuthConfig()
		code := CODE(request.QueryStringParameters["code"])
		if code == "" {
			err = errors.Errorf("'code' is empty in the query: %v", request.QueryStringParameters)
		} else {
			var token oauth2.Token
			token, err = ExchangeCodeForAuthenticationToken(conf, code)
			if len(token.AccessToken) > 8 {
				fmt.Printf("Obtained token: %s...\n", token.AccessToken[:8])
			}

			if err != nil {
				return
			}
			var team slackTeam.SlackTeam
			team, err = extractSlackTeamFromToken(token)
			if err == nil {
				dao := slackTeam.NewDAO(conn.Dynamo, "HandleRedirectURLGetRequest", conn.ClientID)
				err = dao.CreateOrUpdate(team)
			}
		}
	}
	return
}

func extractSlackTeamFromToken(token oauth2.Token) (team slackTeam.SlackTeam, err error) {
	enterpriseID := ""
	// if token.Extra("enterprise_id") != nil {
	// 	enterpriseID = token.Extra("enterprise_id").(string)
	// }
	authedUser1 := token.Extra("authed_user")
	userID := ""
	if authedUser1 != nil {
		authedUser := authedUser1.(map[string]interface{})
		userID = authedUser["id"].(string)
	} else {
		fmt.Printf("WARN Slack responded without authed_user")
	}
	if token.TokenType != "bot" {
		err = errors.Errorf("Unexpected token_type=%s", token.TokenType)
	} else {
		botUserID := token.Extra("bot_user_id").(string)
		team1 := token.Extra("team")
		// accessToken := token.Extra("team")
		teamID := ""
		teamName := ""
		if team1 != nil {
			team := team1.(map[string]interface{})
			// fmt.Printf("bot: %v", bot)
			teamID = team["id"].(string)
			teamName = team["name"].(string)
		} else {
			fmt.Printf("WARN Slack responded without `team {id, name}`")
		}
		scopes1 := token.Extra("scope").(string)
		scopes := strings.Split(scopes1, ",")

		team = slackTeam.SlackTeam{
			TeamID:          common.PlatformID(teamID),
			AccessToken:     token.AccessToken,
			UserID:          userID,
			TeamName:        teamName,
			BotAccessToken:  token.AccessToken, // we only get one token these days
			BotUserID:       botUserID,
			EnterpriseID:    enterpriseID,
			Scopes:          scopes,
		}
	}
	return
}
