package lambda

import (
	"github.com/adaptiveteam/adaptive/lambdas/feedback-report-posting-lambda-go"
	"time"
	"github.com/adaptiveteam/adaptive/lambdas/feedback-reporting-lambda-go"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/slack-go/slack/slackevents"
	"log"
	"strings"
)

func dispatchAppMentionSlackEvent(eventsAPIEvent slackevents.EventsAPIEvent, conn daosCommon.DynamoDBConnection) {
	fmt.Printf("InnerEvent: %v\n", eventsAPIEvent.InnerEvent.Type)
	teamID := models.ParseTeamID(conn.PlatformID)
	if eventsAPIEvent.InnerEvent.Type == slackevents.AppMention {
		slackMsg := ParseAsAppMentionEvent(eventsAPIEvent)
		text := core.TrimLower(slackMsg.Text)
		fmt.Printf("Got app_mention with text: %v\n", text)
		// We first check for requestForUserRegex because botMentionRegex is a subset of this
		if requestForUserRegex.MatchString(text) {
			comms := subscribedCommunities(slackMsg.Channel, conn)

			// It consists of 4 elements, 0: original, 1: first group (channel), 2: second group (command), 3: third group (target)
			list := requestForUserRegex.FindStringSubmatch(text)
			command := core.TrimLower(list[2])
			// when a user is mentioned with '@', the id is coming in with smallcase letters, we stored users with upper case letters
			targetUserID := strings.ToUpper(list[3])

			if len(comms) > 0 {
				if doesUserHavePermissionToExecuteCommand(command, comms) {
					if command == "fetch report for" {
						userMentionFetchReportHandler(*slackMsg, teamID, targetUserID)
					} else if command == "generate report for" {
						userMentionGenerateReportHandler(*slackMsg, teamID, targetUserID)
					} else {
						replyInThread(*slackMsg, teamID, simpleMessage(UserCommandUnknownText))
					}
				} else {
					replyInThread(*slackMsg, teamID, simpleMessage(UserCommandUnknownText))
				}
			} else {
				replyInThread(*slackMsg, teamID, simpleMessage(UnsubscribedUserCommandRejectText))
			}
		} else if botMentionRegex.MatchString(text) {
			_, _, err := refreshUserCache(slackMsg.User, conn)
			if err == nil {
				log.Println(fmt.Sprintf("Got app mention from %s, ensuring profile exists", slackMsg.User))
			} else {
				log.Println(fmt.Sprintf("Error refreshing user profile for %s", slackMsg.User))
			}
			// There will be 3 elements, original, botId, text
			matches := botMentionRegex.FindStringSubmatch(text)
			command := core.TrimLower(strings.ToLower(matches[2]))
			response := onBotMentioned(*slackMsg, teamID, command, conn)
			respond(teamID, response)
		}
	}
}

func userMentionFetchReportHandler(slackMsg slackevents.AppMentionEvent, teamID models.TeamID, targetUserID string) {
	// Posting message to the channel in which user requested this
	replyInThread(slackMsg, teamID, simpleMessage(FetchingReportNotification))
	err2 := feedbackReportPostingLambda.DeliverReportToUserAsync(teamID, targetUserID, time.Now())
	core.ErrorHandler(err2, namespace, "Could not DeliverReportToUserAsync")
}

func userMentionGenerateReportHandler(slackMsg slackevents.AppMentionEvent, teamID models.TeamID, targetUserID string) {
	// Posting message to the channel in which user requested this
	replyInThread(slackMsg, teamID, simpleMessage(GeneratingReportNotification))
	err2 := feedbackReportingLambda.GeneratePerformanceReportAndPostToUserAsync(teamID, targetUserID, time.Now())
	core.ErrorHandler(err2, namespace, "Could not invoke GeneratePerformanceReportAndPostToUserAsync from app-mention")
}

func onBotMentioned(slackMsg slackevents.AppMentionEvent, teamID models.TeamID, command string,
	conn daosCommon.DynamoDBConnection,
) (response platform.Response) {
	comms := subscribedCommunities(slackMsg.Channel, conn)
	switch command {
	case "hello", "hi":
		response = onBotMentionedHelloCommand(comms, slackMsg, teamID)
	default:
		response = platform.PostEphemeral(
			slackMsg.User,
			platform.ConversationID(slackMsg.Channel),
			platform.MessageContent{Message: UserCommandUnknownText},
		)
	}
	fmt.Printf("Going to respond: %v\n", response)
	return
}

func onBotMentionedHelloCommand(comms []models.AdaptiveCommunity, slackMsg slackevents.AppMentionEvent,
	teamID models.TeamID) (response platform.Response) {
	// Post initial engagements for the user
	mc := callback(slackMsg.User, "init", "select")
	var message platform.MessageContent
	if len(availableCommunities(teamID)) == 0 &&
		len(comms) == 0 &&
		len(availableStrategyCommunities(teamID, slackMsg.User)) == 0 {
		message = platform.MessageContent{
			Message: UnsubscribedUserAndNoCommunityAvailableCommandRejectText,
		}
	} else {
		message = platform.MessageContent{
			Attachments: CreateCommunityMenu(mc.ToCallbackID(),
				slackMsg.User, teamID, comms),
		}
	}
	response = platform.PostEphemeral(
		slackMsg.User,
		platform.ConversationID(slackMsg.Channel),
		message,
	)
	return
}

var communityCommandPermissionMap = map[community.AdaptiveCommunity][]string{
	community.HR: {"fetch report for", "generate report for"},
}

func doesUserHavePermissionToExecuteCommand(comm string, comms []models.AdaptiveCommunity) bool {
	for _, each := range comms {
		commands := communityCommandPermissionMap[community.AdaptiveCommunity(each.ID)]
		if core.ListContainsString(commands, comm) {
			return true
		}
	}
	return false
}
