package lambda

import (
	"sort"
	"context"
	"encoding/json"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"

	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	plat "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"golang.org/x/sync/errgroup"
	workflows "github.com/adaptiveteam/adaptive/workflows"
	"github.com/adaptiveteam/adaptive/daos/common"
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
	for _, eng := range engs {
		err2 := postToUser(eng, userID, api)
		if err2 == nil {
			delivered = append(delivered, eng)
		} else {
			logger.WithError(err2).Errorf("Could not deliver engagement with id %s to %s user", eng.ID, userID)
			undelivered = append(undelivered, eng)
		}
	}
	return
}

func greeting() ui.PlainText {
	return defaultGreeting
}

// HandleRequest -
func HandleRequest(ctx context.Context, engage models.UserEngage) {
	defer core.RecoverAsLogError("user-engagement-scripting-lambda-go.HandleRequest")
	logger = logger.WithLambdaContext(ctx)

	// Not invoking scheduler lambda for on-demand asking for engagements
	if !engage.OnDemand {
		// Invoke user-engagement-scheduler lambda to do necessary checks for the user
		schedulerPayload, _ := json.Marshal(engage)
		// Wait until all checks are done and engagements are added to the stream
		_, err := l.InvokeFunction(userEngagementSchedulerLambda, schedulerPayload, false)
		if err != nil {
			logger.WithField("error", err).Errorf("Could not invoke %s lambda", userEngagementSchedulerLambda)
		}
	}

	engs := user.NotPostedUnansweredNotIgnoredEngagements(engage.UserID, engagementTable, engagementAnsweredIndex)
	allValidEngagements := filterEngagements(engs, isNotIgnored)
	logger.WithField("not ignored engagements", &allValidEngagements).Info("Queried engagements")
	logger.Infof("Queried all not ignored engagements for user %s, total: %d", engage.UserID, len(engs))
	showEngagements(engage, allValidEngagements)
	totalCount := len(allValidEngagements)
	conn := common.DynamoDBConnection{Dynamo: d, ClientID: clientID, PlatformID: engage.TeamID.ToPlatformID()}
	count, err2 := workflows.TriggerAllPostponedEvents(engage.TeamID, engage.UserID)(conn)
	core.ErrorHandler(err2, "TriggerAllPostponedEvents", "TriggerAllPostponedEvents")
	totalCount += count
	if totalCount == 0 && engage.OnDemand {
		publish(models.PlatformSimpleNotification{UserId: engage.UserID, Message: string(greeting())})
	}
}

func filterEngagements(engs[]models.UserEngagement, f func(models.UserEngagement)bool ) (filtered []models.UserEngagement) {
	for _, eng := range engs {
		if f(eng) {
			filtered = append(filtered, eng)
		}
	}
	return
}
// UserEngagementSortedByPriorityAndTarget sorter for engagements
type UserEngagementSortedByPriorityAndTarget []models.UserEngagement

func (a UserEngagementSortedByPriorityAndTarget) Len() int           { return len(a) }
func (a UserEngagementSortedByPriorityAndTarget) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UserEngagementSortedByPriorityAndTarget) Less(i, j int) bool { 
	urgentI := 0
	if a[i].Priority != models.UrgentPriority { urgentI = 1}
	urgentJ := 0
	if a[j].Priority != models.UrgentPriority { urgentJ = 1}

	return urgentI < urgentJ || (urgentI == urgentJ && a[i].TargetID < a[j].TargetID) 
}

// Collect urgent and non-urgent engagements
func sortEngagements(engs[]models.UserEngagement) {
	sort.Sort(UserEngagementSortedByPriorityAndTarget(engs))
	// urgentEngs := filterEngagements(engs, isUrgent)
	// nonUrgentEngs := filterEngagements(engs, isNotUrgent)
	// orderedEngs = append(urgentEngs, nonUrgentEngs...)
	return
}
func isNotIgnored(eng models.UserEngagement) bool {
	return eng.Ignored == 0
}
func isUrgent(eng models.UserEngagement) bool {
	return eng.Priority == models.UrgentPriority
}
func isNotUrgent(eng models.UserEngagement) bool {
	return !isUrgent(eng)
}

func showEngagements(engage models.UserEngage, allValidEngagements[]models.UserEngagement) {
	// Check if there are any un-answered high priority engagements
	// First post urgent engagements and then post non-urgent engagements. Also sort by TargetID to improve locality
	sortEngagements(allValidEngagements)

	conn := daosCommon.DynamoDBConnection{
		Dynamo: d,
		ClientID: clientID,
		PlatformID: engage.TeamID.ToPlatformID(),
	}
	slackAdapter := mapper.SlackAdapterForTeamID(conn)
	if !engage.OnDemand && len(allValidEngagements) > 0 {
		slackAdapter.PostSyncUnsafe(plat.Post(plat.ConversationID(engage.UserID), plat.MessageContent{
			Message:     ui.RichText(fmt.Sprintf("You have %d engagements for today", len(allValidEngagements))),
			Attachments: nil,
		}))
	}
	delivered, _ := postEngs(allValidEngagements, engage.UserID, slackAdapter)
	logger.Infof("Posted engagements for user %s", engage.UserID)
	// Deleting only delivered engagements
	updateEngagementsAsPostedAsync(engage.UserID, delivered)
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
