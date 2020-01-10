package lambda

import (
	"time"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/nlopes/slack"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

func isAdaptiveHiResponse(attachs []slack.Attachment) bool {
	if len(attachs) > 0 {
		if attachs[0].Title == user.AdaptiveHiReply {
			return true
		}
	}
	return false
}

func isThereVeryRecentHiResponse(history *slack.History) (res bool) {
	for _, each := range history.Messages {
		if isAdaptiveHiResponse(each.Attachments) {
			t, isDefined := core.ParseDateOrTimestamp(each.Timestamp)
			if isDefined {
				now := time.Now()
				s := now.Sub(t).Seconds()
				if s <= 1 {
					res = true
				}
			} else {
				logger.Warnf("Couldn't parse Timestamp=%s", each.Timestamp)
			}
		}
	}
	return
}

func getChannelHistory(api *slack.Client, postTo string) (history *slack.History, err error) {
	fmt.Printf("api.GetIMHistory(%s, ...)\n", postTo)
	history, err = api.GetIMHistory(postTo, slack.HistoryParameters{
		Latest: slack.DEFAULT_HISTORY_LATEST,
		Oldest: slack.DEFAULT_HISTORY_OLDEST,
		Count:  10,
	})
	return
}

func cleanEarlierHiMessage(api *slack.Client, postTo string) {
	his, err := getChannelHistory(api, postTo)
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
