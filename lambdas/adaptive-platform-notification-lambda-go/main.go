package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ReneKroon/ttlcache"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	plat "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	cache  *ttlcache.Cache
	logger = alog.LambdaLogger(logrus.InfoLevel)
)

func isAdaptiveHiResponse(attachs []slack.Attachment) bool {
	if len(attachs) > 0 {
		if attachs[0].Title == user.AdaptiveHiReply {
			return true
		}
	}
	return false
}

func cleanEarlierHiMessage(api *slack.Client, postTo string) {
	fmt.Printf("api.GetIMHistory(%s, ...)\n", postTo)
	his, err := api.GetIMHistory(postTo, slack.HistoryParameters{
		Latest: slack.DEFAULT_HISTORY_LATEST,
		Oldest: slack.DEFAULT_HISTORY_OLDEST,
		Count:  10,
	})
	if err == nil {
		fmt.Printf("api.GetIMHistory(%s, ...) completed\n", postTo)
		for _, each := range his.Messages {
			if isAdaptiveHiResponse(each.Attachments) {
				fmt.Printf("api.DeleteMessage(%s, ...) (hi response)\n", postTo)
				_, _, _ = api.DeleteMessage(postTo, each.Timestamp)
				fmt.Printf("api.DeleteMessage(%s, ...) (hi response) completed\n", postTo)
			}
		}
	} else {
		logger.WithError(err).Errorf("Unable to retrieve IM history for %s channel", postTo)
	}
}
func platformToken(userID string) plat.UserPlatformToken {
	u := userDao.ReadUnsafe(userID)
	pt := platformTokenDao.ReadUnsafe(models.PlatformID(u.PlatformId))
	return plat.UserPlatformToken{
		PlatformName:  pt.ClientPlatform.PlatformName,
		PlatformToken: pt.ClientPlatform.PlatformToken}
}

func HandleRequest(ctx context.Context, e events.SNSEvent) {
	cache = plat.InitLocalCache(cache)
	for _, record := range e.Records {
		sns := record.SNS
		fmt.Printf("HandleRequest: %s\n", sns.Message)
		var psn models.PlatformSimpleNotification
		err := json.Unmarshal([]byte(sns.Message), &psn)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not unmarshal to PlatformSimpleNotification (%s)", sns.Message))
		var platformResponse plat.PlatformResponse
		err2 := json.Unmarshal([]byte(sns.Message), &platformResponse)
		if platformResponse.PlatformID != "" {
			if err2 == nil {
				messageID, err := platformAdapter.PostSync(platformResponse)
				if err != nil {
					logger.WithField("namespace", namespace).WithField("error", err).
						Errorf("Could not post Slack message with messageID: %s", messageID)
				}
				// fmt.Printf("Waiting for MessageID...\n")
				// fmt.Printf("MessageID=%s\n", <-messageID)
			} else {
				logger.Warnf("Couldn't parse PlatformResponse (%s): %v\n", sns.Message, err2)
			}
		} else {
			fmt.Printf("Parsed as PlatformSimpleNotification: %v\n", psn)
			upt := plat.UserPlatformTokenFromCache(psn.UserId, cache, platformToken, 300*time.Second)
			// upt := platformToken(psn.UserId) - without cache
			if upt.PlatformName == models.SlackPlatform {
				// Slack token and post to slack
				api := slack.New(upt.PlatformToken)
				// Converting generic attachments to slack attachments
				var msgOption slack.MsgOption
				if len(psn.Attachments) > 0 {
					msgOption = slack.MsgOptionAttachments(mapper.Attachments(psn.Attachments)...)
				} else {
					msgOption = slack.MsgOptionAttachments(slack.Attachment{})
				}
				// base message param
				var msgParams = []slack.MsgOption{
					slack.MsgOptionText(psn.Message, false),
					slack.MsgOptionAsUser(true),
					msgOption,
				}

				if psn.ThreadTs != "" {
					msgParams = append(msgParams, slack.MsgOptionTS(psn.ThreadTs))
				}

				var postTo = psn.UserId
				if psn.Channel != "" {
					// When channel is set, post to that. Else, post to the user.
					postTo = psn.Channel
				}

				if psn.Ts == "" {
					// Post new notification to slack
					if isAdaptiveHiResponse(mapper.Attachments(psn.Attachments)) {
						cleanEarlierHiMessage(api, postTo)
					}
					fmt.Printf("api.PostMessage(%s, ...)\n", postTo)
					_, _, err := api.PostMessage(postTo, msgParams...)
					fmt.Printf("api.PostMessage(%s, ...) completed\n", postTo)
					core.ErrorHandler(err, namespace, "Could not post message to slack")
				} else {
					// We update the message when timestamp is not empty and message is not empty or attachments non-empty
					if psn.Message != "" || len(psn.Attachments) > 0 {
						// Update existing message in slack when message is non-empty
						fmt.Printf("api.UpdateMessage(%s, ...)\n", postTo)
						_, _, _, err := api.UpdateMessage(postTo, psn.Ts, msgParams...)
						fmt.Printf("api.UpdateMessage(%s, ...) completed\n", postTo)
						core.ErrorHandler(err, namespace, "Could not update message in slack")
					} else if psn.Message == "" && len(psn.Attachments) == 0 {
						// Delete existing message in slack when message is empty and no attachments
						fmt.Printf("api.DeleteMessage(%s, ...)\n", postTo)
						_, _, err := api.DeleteMessage(postTo, psn.Ts)
						fmt.Printf("api.DeleteMessage(%s, ...) completed\n", postTo)
						core.ErrorHandler(err, namespace, "Could not delete message in slack")
					}
				}
			} else { // if upt.PlatformName == models.MsTeamsPlatform {
				// Handle posting to Teams here
				panic("Unsupported platform " + upt.PlatformName)
			}
		}
	}
	return
}
