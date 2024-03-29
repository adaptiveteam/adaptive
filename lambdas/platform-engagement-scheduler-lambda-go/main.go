package platform_engagement_scheduler_lambda_go

import (
	"context"
	"fmt"
	"log"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	plat "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core_utils_go "github.com/adaptiveteam/adaptive/core-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	region    = utils.NonEmptyEnv("AWS_REGION")
	namespace = utils.NonEmptyEnv("LOG_NAMESPACE")
	clientID  = utils.NonEmptyEnv("CLIENT_ID")
	schema    = models.SchemaForClientID(clientID)
	config    = Config{
		namespace:                       namespace,
		clientConfigTable:               utils.NonEmptyEnv("CLIENT_CONFIG_TABLE_NAME"),
		region:                          region,
		visionTable:                     utils.NonEmptyEnv("VISION_TABLE_NAME"),
		strategyObjectivesTable:         utils.NonEmptyEnv("STRATEGY_OBJECTIVES_TABLE_NAME"),
		strategyObjectivesPlatformIndex: "PlatformIDIndex",
		userObjectivesTable:             utils.NonEmptyEnv("USER_OBJECTIVES_TABLE_NAME"),
		communitiesTable:                utils.NonEmptyEnv("ADAPTIVE_COMMUNITIES_TABLE"),
		l:                               awsutils.NewLambda(region, "", namespace),
		d:                               awsutils.NewDynamo(region, "", namespace),
	}
	// platformAdapter  = mapper.SlackAdapterForTeamID(platformTokenDao)
)

type Config struct {
	namespace                       string
	clientConfigTable               string
	region                          string
	visionTable                     string
	strategyObjectivesTable         string
	strategyObjectivesPlatformIndex string
	userObjectivesTable             string
	communitiesTable                string
	l                               *awsutils.LambdaRequest
	d                               *awsutils.DynamoRequest
}

func communityChannel(community community.AdaptiveCommunity,
	teamID models.TeamID,
	communitiesTable string,
	d *awsutils.DynamoRequest,
	namespace string) (channel string, err error) {
	params := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(string(community)),
		},
		"platform_id": {
			S: aws.String(teamID.ToString()),
		},
	}
	var comm models.AdaptiveCommunity
	err = d.GetItemFromTable(communitiesTable, params, &comm)
	// core.ErrorHandler(err, namespace, fmt.Sprintf("Could not query %s table", communitiesTable))
	channel = comm.ChannelID
	return
}

func HandleRequest(ctx context.Context) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("error in user-engagement-scheduling-lambda-go %v", err2)
		}
	}()
	// Scan platforms table
	// Query all the client configs
	var clientConfigs []models.ClientPlatformToken
	err = config.d.ScanTable(config.clientConfigTable, &clientConfigs)
	for _, clientConfig := range clientConfigs {
		teamID := models.ParseTeamID(clientConfig.PlatformID)
		conn := daosCommon.DynamoDBConnection{
			Dynamo: config.d,
			ClientID: clientID,
			PlatformID: clientConfig.PlatformID,
		}
		slackAdapter := mapper.SlackAdapterForTeamID(conn)
		// Check vision exists
		vision := strategy.StrategyVision(teamID, config.visionTable)
		HRChannel, err := communityChannel(community.HR, teamID, config.communitiesTable, config.d, config.namespace)
		if err != nil {
			err = fmt.Errorf("error in retrieving channel for %s community", community.HR)
		} else {
			if vision == nil {
				slackAdapter.PostAsync(plat.Post(plat.ConversationID(HRChannel), plat.MessageContent{
					Message:     ui.RichText(fmt.Sprintf("There is no vision defined for the org")),
					Attachments: nil,
				}))
				log.Println(fmt.Sprintf("There is no vision defined for the org"))
				// Vision doesn't exist, post a notification to strategy community
			} else {
				slackAdapter.PostAsync(plat.Post(plat.ConversationID(HRChannel), plat.MessageContent{
					Message:     ui.RichText(fmt.Sprintf("Vision is defined for the org")),
					Attachments: nil,
				}))
				log.Println(fmt.Sprintf("Vision is defined for the org"))
			}

			// get strategy objectives
			stratObjs, err2 := strategy.SelectFromStrategyObjectiveJoinUserObjectiveWhereNotCompleted()(conn)
			core_utils_go.ErrorHandler(err2, "HandleRequest", "SelectFromStrategyObjectiveJoinUserObjectiveWhereNotCompleted")
			if len(stratObjs) == 0 {
				slackAdapter.PostAsync(plat.Post(plat.ConversationID(HRChannel), plat.MessageContent{
					Message:     ui.RichText(fmt.Sprintf("No strategy objectives are defined for the org")),
					Attachments: nil,
				}))
				log.Println(fmt.Sprintf("No strategy objectives are defined for the org"))
			} else {
				slackAdapter.PostAsync(plat.Post(plat.ConversationID(HRChannel), plat.MessageContent{
					Message:     ui.RichText(fmt.Sprintf("Strategy objectives are defined for the org")),
					Attachments: nil,
				}))
				log.Println(fmt.Sprintf("Strategy objectives are defined for the org"))
			}
		}
	}
	return
}
