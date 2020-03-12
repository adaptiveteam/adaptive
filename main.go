package main

import (
	"log"

	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	platformNotification "github.com/adaptiveteam/adaptive/lambdas/adaptive-platform-notification-lambda-go"
	communitySlackMessageProcessor "github.com/adaptiveteam/adaptive/lambdas/community-slack-message-processor-lambda-go"
	competencies "github.com/adaptiveteam/adaptive/lambdas/competencies-lambda-go"
	holidays "github.com/adaptiveteam/adaptive/lambdas/holidays-lambda-go"
	platformEngagementScheduler "github.com/adaptiveteam/adaptive/lambdas/platform-engagement-scheduler-lambda-go"
	slackMessageProcessor "github.com/adaptiveteam/adaptive/lambdas/slack-message-processor-lambda-go"
	slackUserQuery "github.com/adaptiveteam/adaptive/lambdas/slack-user-query-lambda-go"
	strategySlackMessageProcessor "github.com/adaptiveteam/adaptive/lambdas/strategy-slack-message-processor-lambda-go"
	userEngagement "github.com/adaptiveteam/adaptive/lambdas/user-engagement-lambda-go"
	userEngagementScheduler "github.com/adaptiveteam/adaptive/lambdas/user-engagement-scheduler-lambda-go"
	userEngagementScheduling "github.com/adaptiveteam/adaptive/lambdas/user-engagement-scheduling-lambda-go"
	userEngagementScripting "github.com/adaptiveteam/adaptive/lambdas/user-engagement-scripting-lambda-go"
	userObjectives "github.com/adaptiveteam/adaptive/lambdas/user-objectives-lambda-go"
	userProfile "github.com/adaptiveteam/adaptive/lambdas/user-profile-lambda-go"
	userQuery "github.com/adaptiveteam/adaptive/lambdas/user-query-lambda-go"
	userSettings "github.com/adaptiveteam/adaptive/lambdas/user-settings-lambda-go"
	userSetup "github.com/adaptiveteam/adaptive/lambdas/user-setup-lambda-go"
	reportingTransformedModelStreaming"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda"
	entityBootstrap "github.com/adaptiveteam/adaptive/lambdas/entity-bootstrap-lambda"
	entityStreaming "github.com/adaptiveteam/adaptive/lambdas/entity-streaming-lambda"
	
	_ "github.com/adaptiveteam/adaptive/daos" // call init to rename table suffixes
	ls "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	role := utils.NonEmptyEnv("LAMBDA_ROLE")
	switch role {
	case "platform-notification":
		ls.Start(platformNotification.HandleRequest)
	case "platform-engagement-scheduler":
		ls.Start(platformEngagementScheduler.HandleRequest)
	case "slack-message-processor":
		ls.Start(slackMessageProcessor.HandleRequest)
	case "slack-user-query":
		ls.Start(slackUserQuery.HandleRequest)
	case "user-engagement":
		ls.Start(userEngagement.HandleRequest)
	case "user-engagement-scheduler":
		ls.Start(userEngagementScheduler.HandleRequest)
	case "user-engagement-scheduling":
		ls.Start(userEngagementScheduling.HandleRequest)
	case "user-engagement-scripting":
		ls.Start(userEngagementScripting.HandleRequest)
	case "user-profile":
		ls.Start(userProfile.HandleRequest)
	case "user-query":
		ls.Start(userQuery.HandleRequest)
	case "user-settings":
		ls.Start(userSettings.HandleRequest)
	case "user-setup":
		ls.Start(userSetup.HandleRequest)
	case "strategy-slack-message-processor":
		ls.Start(strategySlackMessageProcessor.HandleRequest)
	case "holidays":
		holidays.LambdaRouting.StartHandler()
	case "competencies":
		ls.Start(competencies.HandleRequest)
	case "community-slack-message-processor":
		ls.Start(communitySlackMessageProcessor.HandleRequest)
	case "user-objectives":
		ls.Start(userObjectives.HandleRequest)
	case "stream-event-mapping":
		ls.Start(reportingTransformedModelStreaming.HandleRequest)
	case "entity-bootstrapping":
		ls.Start(entityBootstrap.HandleRequest)
	case "entity-streaming":
		ls.Start(entityStreaming.HandleRequest)	
	default:
		log.Printf("Unknown role %s", role)
	}
}
