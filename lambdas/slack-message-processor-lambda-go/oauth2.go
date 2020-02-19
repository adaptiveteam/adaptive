package lambda

import (
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"fmt"
	"strings"

	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"

	"github.com/aws/aws-lambda-go/events"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type CODE string

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

		Endpoint:     endpoints.Slack,
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
	url = conf.AuthCodeURL("state", oauth2.AccessTypeOffline, userScope)
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
func GenerateAddToSlackURL(userID, channelID string, platformID models.PlatformID) {
	slackAdapter := platformAdapter.ForPlatformID(platformID)
	url := generateAuthorizationURL(GlobalSlackOAuthConfig())
	slackAdapter.PostSyncUnsafe(platform.Post(platform.ConversationID(userID), 
		platform.Message(ui.Sprintf(
			"Send the following link to a user wishing to add Adaptive to their Slack workspace: %s", url)),
	))
}

// HandleRedirectURLGetRequest processes GET request from Slack for OAuth2
func HandleRedirectURLGetRequest(request events.APIGatewayProxyRequest) (err error) {
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
		fmt.Printf("Obtained token: %s...", token.AccessToken[:8])
	}
	// token. - save to DB
	return
}
