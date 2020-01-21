package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	plat "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/core-utils-go/mmap"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"golang.org/x/sync/errgroup"
)

func postToUser(eng models.UserEngagement, userID string, api mapper.PlatformAPI) (err error) {
	var postMsg ebm.Message
	err = json.Unmarshal([]byte(eng.Script), &postMsg)
	if err == nil {
		_, err = api.PostSync(plat.PostToThread(plat.ThreadID{
			ConversationID: plat.ConversationID(userID),
		}, plat.MessageContent{
			Attachments: postMsg.Attachments,
		}))
	}
	return
}


func publish(msg models.PlatformSimpleNotification) {
	_, err := sns.Publish(msg, platformNotificationTopic)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not pusblish message to %s topic", platformNotificationTopic))
}

func postEngs(engs []models.UserEngagement, userID string, api mapper.PlatformAPI) (delivered, undelivered []models.UserEngagement) {
	// TODO: Implement a more sophisticated scripting algorithm
	m := mmap.NewMultiMap()
	for _, each := range engs {
		// Grouping engagements by target id
		m.Put(each.TargetID, each)
	}
	// Taking one target user at a time
	for idx := range m.KeySet() {
		list, _ := m.Get(m.KeySet()[idx])
		for _, each := range list {
			eng := each.(models.UserEngagement)
			err := postToUser(eng, userID, api)
			if err == nil {
				delivered = append(delivered, eng)

			} else {
				logger.Errorf("Could not deliver engagement with id %s to %s user", eng.ID, userID)
				undelivered = append(undelivered, eng)
			}
		}
	}
	return
}

func greeting() ui.PlainText {
	return defaultGreeting
}

type EngSchedule struct {
	Target string `json:"target"`
	Date   string `json:"date"`
}

func HandleRequest(ctx context.Context, engage models.UserEngageWithCheckValues) {
	logger = logger.WithLambdaContext(ctx)

	// Not invoking scheduler lambda for on-demand asking for engagements
	if !engage.OnDemand {
		// Invoke user-engagement-scheduler lambda to do necessary checks for the user
		schedulerPayload, _ := json.Marshal(EngSchedule{Target: engage.UserId})
		// Wait until all checks are done and engagements are added to the stream
		_, err := l.InvokeFunction(userEngagementSchedulerLambda, schedulerPayload, false)
		if err != nil {
			logger.WithField("error", err).Errorf("Could not invoke %s lambda", userEngagementSchedulerLambda)
		}
	}

	engs := user.NotPostedUnansweredNotIgnoredEngagements(engage.UserId, engagementTable, engagementAnsweredIndex)
	logger.WithField("engagements", &engs).Info("Queried engagements")
	logger.Infof("Queried all the engagements for user %s, total: %d", engage.UserId, len(engs))

	// Check if there are any un-answered high priority engagements
	var urgentEngs []models.UserEngagement
	var nonUrgentEngs []models.UserEngagement

	// Collect urgent and non-urgent engagements
	for _, eng := range engs {
		if eng.Priority == models.UrgentPriority && eng.Ignored == 0 {
			urgentEngs = append(urgentEngs, eng)
		} else if eng.Ignored == 0 {
			nonUrgentEngs = append(nonUrgentEngs, eng)
		}
	}

	allValidEngagements := append(urgentEngs, nonUrgentEngs...)

	if len(allValidEngagements) > 0 {
		slackAdapter := platformAdapter.ForPlatformID(engage.PlatformID)
		if !engage.OnDemand {
			slackAdapter.PostSyncUnsafe(plat.Post(plat.ConversationID(engage.UserId), plat.MessageContent{
				Message:     ui.RichText(fmt.Sprintf("You have %d engagements for today", len(allValidEngagements))),
				Attachments: nil,
			}))
		}
		var delivered []models.UserEngagement
		// First post urgent engagements and then post non-urgent engagements
		if len(urgentEngs) > 0 {
			delivered1, _ := postEngs(urgentEngs, engage.UserId, slackAdapter)
			delivered = append(delivered, delivered1...)
		}
		logger.Infof("Posted urgent engagements for user %s", engage.UserId)
		if len(nonUrgentEngs) > 0 {
			delivered2, _ := postEngs(nonUrgentEngs, engage.UserId, slackAdapter)
			delivered = append(delivered, delivered2...)
		}
		logger.Infof("Posted non-urgent engagements for user %s", engage.UserId)
		// Deleting only delivered engagements
		updateEngagementsAsPostedAsync(engage.UserId, delivered)

	} else if engage.OnDemand {
		publish(models.PlatformSimpleNotification{UserId: engage.UserId, Message: string(greeting())})
	}
}

func updateEngagementsAsPostedAsync(userID string, engagements []models.UserEngagement) {
	var g errgroup.Group
	for _, each := range engagements {
		e := each // https://golang.org/doc/faq#closures_and_goroutines
		postedTime := core.CurrentRFCTimestamp()
		g.Go(func() error {
			params := map[string]*dynamodb.AttributeValue{
				"user_id": dynString(e.UserID),
				"id":      dynString(e.ID),
			}

			exprAttributes := map[string]*dynamodb.AttributeValue{
				":pt": dynString(postedTime),
			}
			updateExpression := "set posted_at = :pt"
			return d.UpdateTableEntry(exprAttributes, params, updateExpression, engagementTable)
		})
	}

	// Wait for all deletions to complete
	if err := g.Wait(); err == nil {
		logger.Infof("Deleted scheduled engagements from table for %s after posting", userID)
	} else {
		logger.Infof("Could not delete scheduled engagements for %s after posting: %v", userID, err)
	}
}

func dynString(str string) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{S: aws.String(str)}
	return &attr
}
