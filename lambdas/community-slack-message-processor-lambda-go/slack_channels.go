package lambda

import (
	"github.com/adaptiveteam/adaptive/pagination"
	"github.com/slack-go/slack"
)

// SlackChannelsToInterfaces convert []channel to []interface{}
func SlackChannelsToInterfaces(channels []slack.Channel) (sl pagination.InterfaceSlice) {
	for _, ch := range channels {
		sl = append(sl, ch)
	}
	return
}

// SlackGetConversationsPager returns pager that will retrieve all conversations according to the request
func SlackGetConversationsPager(api *slack.Client, conversationParams slack.GetConversationsParameters) pagination.InterfacePager {
	return func() (sl pagination.InterfaceSlice, ip pagination.InterfacePager, err error) {
		var channels []slack.Channel
		var cursor string
		channels, cursor, err = api.GetConversations(&conversationParams)
		if err == nil {
			sl = SlackChannelsToInterfaces(channels)
			conversationParams.Cursor = cursor
			if cursor == "" {
				ip = pagination.InterfacePagerPure()
			} else {
				ip = SlackGetConversationsPager(api, conversationParams)
			}
		}
		return
	}
}

// SlackGetUsersInConversationPager retrieves users from conversation
func SlackGetUsersInConversationPager(api *slack.Client, params slack.GetUsersInConversationParameters)pagination.InterfacePager {
	return func() (sl pagination.InterfaceSlice, ip pagination.InterfacePager, err error) {
		var users []string
		var cursor string
		users, cursor, err = api.GetUsersInConversation(&params)
		if err == nil {
			sl = pagination.StringsToInterfaceSlice(users)
			params.Cursor = cursor
			if cursor == "" {
				ip = pagination.InterfacePagerPure()
			} else {
				ip = SlackGetUsersInConversationPager(api, params)
			}
		}
		return
	}
}
