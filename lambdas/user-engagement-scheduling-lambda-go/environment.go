package userEngagementScheduling

import (
	"encoding/json"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"

	// rename tables
	_ "github.com/adaptiveteam/adaptive/daos"
)

type Config struct {
	namespace               string
	engScriptingLambdaArn   string
	engSchedulerLambdaName  string
	usersTable              string
	communityUsersTable     string
	usersScheduledTimeIndex string
	usersZoneOffsetIndex    string
	communityUsersUserIndex string
	clientConfigTable       string
	region                  string
	cw                      *awsutils.CloudWatchRequest
	l                       *awsutils.LambdaRequest
	d                       *awsutils.DynamoRequest
	clientID                string
	connGen                 daosCommon.DynamoDBConnectionGen
}

func readConfigFromEnvironment() Config {
	region := utils.NonEmptyEnv("AWS_REGION")
	namespace := utils.NonEmptyEnv("LOG_NAMESPACE")
	d := awsutils.NewDynamo(region, "", namespace)
	clientID := utils.NonEmptyEnv("CLIENT_ID")
	return Config{
		namespace:                  namespace,
		engScriptingLambdaArn:      utils.NonEmptyEnv("USER_ENGAGEMENT_SCRIPTING_LAMBDA_ARN"),
		engSchedulerLambdaName:     utils.NonEmptyEnv("USER_ENGAGEMENT_SCHEDULER_LAMBDA_NAME"),
		usersTable:                 utils.NonEmptyEnv("USERS_TABLE_NAME"),
		communityUsersTable:        utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME"),
		usersScheduledTimeIndex:    "PlatformIDAdaptiveScheduledTimeInUTCIndex",
		usersZoneOffsetIndex:       "PlatformIDTimezoneOffsetIndex",
		communityUsersUserIndex:    utils.NonEmptyEnv("COMMUNITY_USERS_USER_INDEX"),
		clientConfigTable:          utils.NonEmptyEnv("CLIENT_CONFIG_TABLE_NAME"),
		region:                     region,
		cw:                         awsutils.NewCloudWatch(region, "", namespace),
		l:                          awsutils.NewLambda(region, "", namespace),
		d:                          d,
		clientID: 				    clientID,
		connGen:                    daosCommon.DynamoDBConnectionGen{
			                            Dynamo:     d,
			                            TableNamePrefix:   clientID,
		},
 	}
}

func invokeScriptingLambda(engage models.UserEngage, config Config) (err error) {
	payloadJSONBytes, _ := json.Marshal(engage)
	_, err = config.l.InvokeFunction(config.engScriptingLambdaArn, payloadJSONBytes, true)
	return
}

func invokeSchedulerLambda(engage models.UserEngage, config Config) (err error) {
	payloadJSONBytes, _ := json.Marshal(engage)
	_, err = config.l.InvokeFunction(config.engSchedulerLambdaName, payloadJSONBytes, true)
	return
}

func postToCommunity(message platform.MessageContent, comm community.AdaptiveCommunity) func(config Config, teamID models.TeamID) (err error) {
	return func(config Config, teamID models.TeamID) (err error) {
		conn := config.connGen.ForPlatformID(teamID.ToPlatformID())
		var communities []adaptiveCommunity.AdaptiveCommunity
		communities, err = adaptiveCommunity.ReadOrEmpty(teamID.ToPlatformID(), string(comm))(conn)
		if err == nil {
			if len(communities) > 0 {
				userComm := communities[0]
				slackAdapter := mapper.SlackAdapterForTeamID(conn)
				post := platform.Post(
					platform.ConversationID(userComm.ChannelID),
					message,
				)
				logger.Infof("Posting to %s: %v", userComm.ChannelID, post)
				_, err = slackAdapter.PostSync(post)
			} else {
				logger.Warnf("%s community not found for team %s", comm, teamID.ToString())
			}
		}
		return
	}
}
