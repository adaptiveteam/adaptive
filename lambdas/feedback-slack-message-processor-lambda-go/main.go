package lambda

import (
	"github.com/adaptiveteam/adaptive/lambdas/feedback-report-posting-lambda-go"
	"fmt"
	"context"
	"encoding/json"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	feedbackReportingLambda "github.com/adaptiveteam/adaptive/lambdas/feedback-reporting-lambda-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/adaptiveteam/adaptive/workflows"
	ls "github.com/aws/aws-lambda-go/lambda"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

func publish(msg models.PlatformSimpleNotification) {
	_, err := sns.Publish(msg, platformNotificationTopic)
	core.ErrorHandler(err, "publish", fmt.Sprintf("publish(): Could not publish message to %s topic", platformNotificationTopic))
}

func respond(teamID models.TeamID, responses ...platform.Response) {
	for _, response := range responses {
		logger.Infof("Respond(,%v): %s", response, response.Type)
		presp := platform.TeamResponse{
			TeamID:   teamID,
			Response: response,
		}
		// publish(presp)
		_, err := sns.Publish(presp, platformNotificationTopic)
		core.ErrorHandler(err, "respond", fmt.Sprintf("respond(): Could not publish message to %s topic", platformNotificationTopic))
		// logger.WithField("error", err).Errorf("Could not publish message to %s topic", platformNotificationTopic)
	}
}

func HandleRequest(ctx context.Context, np models.NamespacePayload4) (err error) {
	defer core.RecoverAsLogError("feedback-slack-message-processor-lambda")

	conn := connGen.ForPlatformID(np.TeamID.ToPlatformID())
	if strings.HasPrefix(np.InteractionCallback.CallbackID, "/") {
		err = workflows.InvokeWorkflow(np, conn)
	} else if np.ID == "" {
		log.Println("Warmed up...")
		return
	} else {
		teamID := np.TeamID
		// This module only looks for payload with 'feedback' namespace
		if np.Namespace == "feedback" {
			switch np.SlackRequest.Type {
			case models.InteractionSlackRequestType:
				message := np.SlackRequest.InteractionCallback
				// This is to handle the hello message
				if message.CallbackID == "init_message" {
					action := message.ActionCallback.AttachmentActions[0]
					selected := action.SelectedOptions[0]
					text := selected.Value
					if text == coaching.GiveFeedback {
						logger.Infof("Handling give feedback event")
						responses := handleFeedbackRequest(message.User.ID, message.Channel.ID,
							GiveFeedbackAction, GiveFeedbackMessage, "give-feedback", message.MessageTs, conn)
						respond(teamID, responses...)
					} else if text == coaching.RequestFeedback {
						logger.Infof("Handling request feedback event")
						responses := handleFeedbackRequest(message.User.ID, message.Channel.ID,
							RequestFeedbackAction, RequestFeedbackMessage, "request-feedback", 
							message.MessageTs, conn)
						respond(teamID, responses...)
					} else if text == user.FetchReport {
						err = feedbackReportPostingLambda.DeliverReportToUserAsync(teamID, message.User.ID, time.Now())
						logger.WithField("error", err).Error("Could not DeliverReportToUserAsync")

						msg := core.IfThenElse(err == nil, FetchingReportMessage, InternalErrorMessage).(ui.RichText)
						// Update original message
						publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID,
							Message: string(msg), Attachments: []model.Attachment{}, Ts: message.MessageTs})
					} else if text == user.GenerateReport {
						err = feedbackReportingLambda.GeneratePerformanceReportAndPostToUserAsync(teamID, message.User.ID, time.Now())
						if err != nil {
							logger.WithField("error", err).
								Error("Could not GeneratePerformanceReportAndPostToUserAsync from feedback-slack-message-processor")
						}
						msg := core.IfThenElse(err == nil, GeneratingReportMessage, InternalErrorMessage).(ui.RichText)
						// Update original message
						publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID,
							Message: string(msg), Ts: message.MessageTs})
					} else if text == coaching.ViewCoachees {
						err = workflows.EnterWorkflow(workflows.ViewCoacheesWorkflow, np, conn, "")//onViewCoacheeIDOs(np)
					} else if text == coaching.ViewAdvocates {
						response := onViewAdvocates(message.User.ID, conn)
						respond(teamID, response,
							platform.DeleteByResponseURL(message.ResponseURL))
					}
				} else {
					passthrough(np)
				}
			default:
				// for interaction and dialog_submission events, we invoke feedback setup lambda
				passthrough(np)
			}
		}
	}
	if err != nil {
		logger.WithError(err).Errorf("Error in HandleRequest: %+v", err)
	}
	return
}

func passthrough(np models.NamespacePayload4) {
	// for interaction and dialog_submission events, we invoke feedback setup lambda
	bytes, _ := json.Marshal(np)
	_, err := l.InvokeFunction(feedbackSetupLambdaName, []byte(bytes), false)
	logger.WithField("error", err).Errorf("Could not invoke %s lambda", feedbackSetupLambdaName)
}

func main() {
	ls.Start(HandleRequest)
}

type SortByRichText []ui.RichText

func (a SortByRichText) Len() int           { return len(a) }
func (a SortByRichText) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByRichText) Less(i, j int) bool { return a[i] < a[j] }

type RichTextGroup struct {
	GroupID  string
	Elements []ui.RichText
}

// Sort sorts the group elements
func (rtg RichTextGroup) Sort() {
	sort.Sort(SortByRichText(rtg.Elements))
}

// JoinInParenthesisTemplate joins the group elements and returns them in ()
func (rtg RichTextGroup) JoinInParenthesisTemplate() ui.RichText {
	return ui.RichText("(") + ui.Join(rtg.Elements, ", ") + ")"
}

// FormatObjective renders an objective
type FormatObjective func(userObjective.UserObjective) ui.RichText

// FormatObjectiveName shows the objective name
func FormatObjectiveName(objective userObjective.UserObjective) ui.RichText {
	return ui.RichText(objective.Name)
}

// FormatObjectiveItem renders objective as an item in list
func FormatObjectiveItem(objective userObjective.UserObjective) ui.RichText {
	return ui.Sprintf(">*%s*\n>%s\n", objective.Name, objective.Description)
}

// GroupByUserID groups objectives by user id and renders each objective as RichText
func GroupByUserID(objectives []userObjective.UserObjective, formatObjective FormatObjective) (groupedObjectives map[string]RichTextGroup) {
	groupedObjectives = make(map[string]RichTextGroup, 0)
	for _, objective := range objectives {
		group, ok := groupedObjectives[objective.UserID]
		if !ok {
			group = RichTextGroup{GroupID: objective.UserID}
		}
		group.Elements = append(group.Elements, formatObjective(objective))
		groupedObjectives[objective.UserID] = group
	}
	return
}

func renderGroupsAsIDAndElementsInParentheses(groups map[string]RichTextGroup) (items []ui.RichText) {
	items = make([]ui.RichText, 0)
	for _, group := range groups {
		group.Sort()
		items = append(items, ui.RichText("<@"+group.GroupID+"> ")+group.JoinInParenthesisTemplate())
	}
	return
}

func renderGroupsAsIDAndElementsAsSubList(groups map[string]RichTextGroup) (items []ui.RichText) {
	items = make([]ui.RichText, 0)
	for _, group := range groups {
		group.Sort()
		items = append(items, ui.RichText("<@"+group.GroupID+">\n")+ui.Join(group.Elements, "\n"))
	}
	return
}

// accountabilityPartnerID is the same as coachID
// it's the userID for the user we are interacting at the moment
func onViewAdvocates(accountabilityPartnerID string, conn daosCommon.DynamoDBConnection) platform.Response {
	objectives := userObjective.ReadByAccountabilityPartnerUnsafe(accountabilityPartnerID)(conn)
	strObjectives := filterObjectivesByObjectiveType(objectives, userObjective.StrategyDevelopmentObjective)
	infos := GroupByUserID(strObjectives, FormatObjectiveName)
	names := renderGroupsAsIDAndElementsInParentheses(infos)
	text := AdvocateListPretext + "\n" + ui.ListItems(names...)
	return platform.Post(platform.ConversationID(accountabilityPartnerID), platform.MessageContent{Message: text})
}

func handleFeedbackRequest(userID, channelID, action, text, context, ts string, conn daosCommon.DynamoDBConnection) (
	responses []platform.Response) {
	teamID := models.ParseTeamID(conn.PlatformID)
	// Getting current year and month
	// We are writing month rather than year in engagement because quarter can always be inferred from month
	year, month, _ := time.Now().Date()
	// Setting action as 'select'
	mc := models.MessageCallback{Module: "coaching", Source: userID,
		Topic: "user_feedback", Action: action, Target: "", Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	communityUserIDs, err := community.PlatformCommunityMemberIDs(teamID, communitiesTable,
		communityPlatformIndex, communityUsersTable, communityUsersCommunityIndex)
	if err == nil {
		usersWithoutSelf := core.InAButNotB(communityUserIDs, []string{userID})
		if len(usersWithoutSelf) > 0 {
			UserSelectEngagement(userID, mc, usersWithoutSelf, []string{}, text, context, conn)
		} else {
			switch action {
			case GiveFeedbackAction:
				responses = append(responses, platform.Post(platform.ConversationID(userID),
					platform.MessageContent{Message: GiveFeedbackNoUsersExistMessage}), )
			case RequestFeedbackAction:
				responses = append(responses, platform.Post(platform.ConversationID(userID),
					platform.MessageContent{Message: RequestFeedbackNoUsersExistMessage}), )
			}
		}
		// Delete the original engagement
		responses = append(responses, platform.Delete(platform.TargetMessageID{
			ConversationID: platform.ConversationID(channelID),
			Ts:             ts,
		}))
	} else {
		logger.WithField("error", err).
			Errorf("Error with retrieving community user ids for %s platform", teamID)
	}
	return
}

func UserSelectEngagement(userID string, mc models.MessageCallback,
	users, filter []string, text, context string, conn daosCommon.DynamoDBConnection) {
	user.UserSelectEng(userID, engagementsTable, conn,
		mc, users, filter, text, context, models.UserEngagementCheckWithValue{})
}
