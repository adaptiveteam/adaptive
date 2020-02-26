package lambda

import (
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"fmt"
	"strings"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"

	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"

	"github.com/aws/aws-lambda-go/events"
	"golang.org/x/net/context"

	// oauth2 "github.com/adaptiveteam/adaptive/lambdas/slack-message-processor-lambda-go/oauth2-reimpl"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/slackTeam"
	"golang.org/x/oauth2"
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
	AdaptiveScopes = "calls:read,calls:write,channels:read,dnd:read,files:read,groups:history,groups:read,groups:write,im:history,im:read,im:write,incoming-webhook,links:write,mpim:history,mpim:read,mpim:write,pins:write,reactions:read,reactions:write,remote_files:read,remote_files:share,remote_files:write,team:read,usergroups:read,usergroups:write,users:read,users:read.email,users:write"
	// app_mentions:read,chat:write,files:write,
	UserScopes = "" // "channels:read,chat:write,files:write,groups:read,groups:write,im:history,users:read"
	// https://1vvtp0yc61.execute-api.us-east-1.amazonaws.com/dev/%7Bproxy+%7D?code=436528929141.946109056755.2599abf6300120ff5bc70e3b20a6680e64e05d652dc046b72c6c4d24a9f73c69&state=state
)

// SlackOAuthConfig - constructs OAuth2 config with Slack endpoint and AdaptiveScopes
func SlackOAuthConfig(clientID, clientSecret string) oauth2.Config {
	return oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       strings.Split(AdaptiveScopes, ","),

		// Endpoint:     oauth2.Endpoint{
		// 	AuthStyle:   oauth2.AuthStyle(endpoints.Slack.AuthStyle),
		// 	AuthURL:     endpoints.Slack.AuthURL,
		// 	TokenURL:    endpoints.Slack.TokenURL,
		// },
		Endpoint: endpoints.Slack,
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
	conn := globalConnection(teamID.ToPlatformID())
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
	conf := GlobalSlackOAuthConfig()
	code := CODE(request.QueryStringParameters["code"])
	var token oauth2.Token
	token, err = ExchangeCodeForAuthenticationToken(conf, code)
	if len(token.AccessToken) > 8 {
		fmt.Printf("Obtained token: %s...\n", token.AccessToken[:8])
	}

	if err != nil {
		return
	}
	enterpriseID := ""
	if token.Extra("enterprise_id") != nil {
		enterpriseID = token.Extra("enterprise_id").(string)
	}
	team := slackTeam.SlackTeam{
		TeamID:       common.PlatformID(token.Extra("team_id").(string)),
		AccessToken:  token.AccessToken,
		UserID:       token.Extra("user_id").(string),
		TeamName:     token.Extra("team_name").(string),
		EnterpriseID: enterpriseID,
	}
	dao := slackTeam.NewDAO(conn.Dynamo, "HandleRedirectURLGetRequest", conn.ClientID)
	err = dao.CreateOrUpdate(team)
	return
}
