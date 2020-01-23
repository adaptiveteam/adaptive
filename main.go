package main

import (
	"log"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	platformNotification "github.com/adaptiveteam/adaptive/lambdas/adaptive-platform-notification-lambda-go"
	platformEngagementScheduler "github.com/adaptiveteam/adaptive/lambdas/platform-engagement-scheduler-lambda-go"
	slackMessageProcessor "github.com/adaptiveteam/adaptive/lambdas/slack-message-processor-lambda-go"
	slackUserQuery "github.com/adaptiveteam/adaptive/lambdas/slack-user-query-lambda-go"
	userEngagement "github.com/adaptiveteam/adaptive/lambdas/user-engagement-lambda-go"
	userEngagementScheduler "github.com/adaptiveteam/adaptive/lambdas/user-engagement-scheduler-lambda-go"
	userEngagementScheduling "github.com/adaptiveteam/adaptive/lambdas/user-engagement-scheduling-lambda-go"
	userEngagementScripting "github.com/adaptiveteam/adaptive/lambdas/user-engagement-scripting-lambda-go"
	userProfile "github.com/adaptiveteam/adaptive/lambdas/user-profile-lambda-go"
	userQuery "github.com/adaptiveteam/adaptive/lambdas/user-query-lambda-go"
	userSettings "github.com/adaptiveteam/adaptive/lambdas/user-settings-lambda-go"
	userSetup "github.com/adaptiveteam/adaptive/lambdas/user-setup-lambda-go"
	ls "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	role := utils.NonEmptyEnv("LAMBDA_ROLE")
	switch role {
	case "platform-notification": ls.Start(platformNotification.HandleRequest)
	case "platform-engagement-scheduler": ls.Start(platformEngagementScheduler.HandleRequest)
	case "slack-message-processor": ls.Start(slackMessageProcessor.HandleRequest)
	case "slack-user-query": ls.Start(slackUserQuery.HandleRequest)
	case "user-engagement": ls.Start(userEngagement.HandleRequest)
	case "user-engagement-scheduler": ls.Start(userEngagementScheduler.HandleRequest)
	case "user-engagement-scheduling": ls.Start(userEngagementScheduling.HandleRequest)
	case "user-engagement-scripting": ls.Start(userEngagementScripting.HandleRequest)
	case "user-profile": ls.Start(userProfile.HandleRequest)
	case "user-query": ls.Start(userQuery.HandleRequest)
	case "user-settings": ls.Start(userSettings.HandleRequest)
	case "user-setup": ls.Start(userSetup.HandleRequest)
	default:
		log.Printf("Unknown role %s", role)
	}
}
