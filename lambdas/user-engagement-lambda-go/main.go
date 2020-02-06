package lambda

import (
	"context"
	"encoding/json"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core_utils_go "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	. "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

func postToUser(userId string, eng models.UserEngagement, namespace string) (err error) {
	var postMsg ebm.Message
	err = json.Unmarshal([]byte(eng.Script), &postMsg)
	if err == nil {
		slackEng := SlackEngagement{Message: postMsg}
		options := slackEng.MsgOptions()
		ut, err := utils.UserToken(userId, profileLambda, region, namespace)
		if err == nil {
			api := slack.New(ut.PlatformToken)
			_, _, err = api.PostMessage(userId, options...)
		}
	}
	return err
}

func updateEngagementAsPosted(userID, engID string) (err error) {
	params := map[string]*dynamodb.AttributeValue{
		"user_id": dynString(userID),
		"id":      dynString(engID),
	}

	exprAttributes := map[string]*dynamodb.AttributeValue{
		":pt": dynString(core_utils_go.CurrentRFCTimestamp()),
	}
	updateExpression := "set posted_at = :pt"
	return d.UpdateTableEntry(exprAttributes, params, updateExpression, engagementTable)
}

var (
	namespace                  = utils.NonEmptyEnv("LOG_NAMESPACE")
	region                     = utils.NonEmptyEnv("AWS_REGION")
	profileLambda              = utils.NonEmptyEnv("USER_PROFILE_LAMBDA_NAME")
	communityUsersTable        = utils.NonEmptyEnv("COMMUNITY_USERS_TABLE_NAME")
	communityUsersUserIndex    = utils.NonEmptyEnv("COMMUNITY_USERS_USER_INDEX")
	communityUsersChannelIndex = "ChannelIDIndex"
	engagementTable            = utils.NonEmptyEnv("USER_ENGAGEMENTS_TABLE_NAME")
	d                          = awsutils.NewDynamo(utils.NonEmptyEnv("AWS_REGION"), "", namespace)

	logger = alog.LambdaLogger(logrus.InfoLevel)
)

func queryCommunityChannelIndex(channelId string) ([]interface{}, error) {
	var rels []interface{}
	err := d.QueryTableWithIndex(communityUsersTable, awsutils.DynamoIndexExpression{
		IndexName: communityUsersChannelIndex,
		Condition: "channel_id = :c",
		Attributes: map[string]interface{}{
			":c": channelId,
		},
	}, map[string]string{}, true, -1, &rels)
	return rels, err
}

func dynString(str string) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{S: aws.String(str)}
	return &attr
}

func HandleRequest(ctx context.Context, e events.DynamoDBEvent) {
	for _, record := range e.Records {
		logger.Infof("Processing request data for event ID %s, type %s", record.EventID, record.EventName)
		// sometimes, existing engagements could be modified, like coaching feedback
		// we should treat scenario same as a new event
		if record.EventName == string(events.DynamoDBOperationTypeInsert) || record.EventName == string(events.DynamoDBOperationTypeModify) {
			var eng models.UserEngagement
			// Print new values for attributes of type String
			err := awsutils.UnmarshalStreamImage(record.Change.NewImage, &eng)
			logger.WithField("entity", &eng).Infof("Received %s event", record.EventName)
			if err == nil {
				if eng.PostedAt == "" {
					userRels := strategy.QueryCommunityUserIndex(eng.UserID, communityUsersTable, communityUsersUserIndex)
					// For community related engagements, userId is the channel. So, we also check if the channel exists in community-users table
					channelRels, err := queryCommunityChannelIndex(eng.UserID)

					if err == nil {
						// Post engagement to user only when a user is part of some community
						if len(userRels) > 0 || len(channelRels) > 0 {
							if eng.Priority == models.UrgentPriority && eng.Answered == 0 && eng.Ignored == 0 {
								err = postToUser(eng.UserID, eng, namespace)
								if err == nil {
									logger.Infof("Posted urgent engagement with id %s for user %s", eng.ID, eng.UserID)
									// After posting the engagement to the user, delete it from the table
									err = updateEngagementAsPosted(eng.UserID, eng.ID)
								}
							}
						}
					} else {
						logger.WithField("error", err).Errorf("Could not query %s table on %s index", communityUsersTable, communityUsersChannelIndex)
					}
				}
			} else {
				logger.WithField("error", err).Errorf("Could not unmarshal stream image")
			}
		}
	}
	return
}
