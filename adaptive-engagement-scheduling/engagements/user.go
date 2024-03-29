package engagements

import (
	"fmt"
	"log"
	"strconv"
	"time"

	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/workflows"
	"github.com/adaptiveteam/adaptive/workflows/exchange"

	. "github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	issuesUtils "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	bt "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	feedbackReportingLambda "github.com/adaptiveteam/adaptive/lambdas/feedback-reporting-lambda-go"
)

/*
------------------------------------------------------------------------------------
IDO Creation reminders
------------------------------------------------------------------------------------
*/

var globalConnectionGen = daosCommon.CreateConnectionGenFromEnv

// IDOCreateReminder is meant to trigger the engagements that
// reminds the user to create personal improvement objects in the event that they have
// not created any.
func IDOCreateReminder(teamID models.TeamID, date bt.Date, target string) {
	log.Println(fmt.Sprintf("Checking IDOCreateReminder for user: %s", target))
	conn := globalConnectionGen().ForPlatformID(teamID.ToPlatformID())
	AddObjective(target, models.ParseTeamID(conn.PlatformID), false, utilsUser.UserIDsToDisplayNamesConnUnsafe(conn))

}

/*
------------------------------------------------------------------------------------
Update Reminders
------------------------------------------------------------------------------------
*/

// slowlyCreateConnection - 
// Deprecated: use connGen approach
func slowlyCreateConnection(teamID models.TeamID) daosCommon.DynamoDBConnection {
	return daosCommon.DynamoDBConnection {
		Dynamo: common.DeprecatedGetGlobalDns().Dynamo,
		ClientID: utils.NonEmptyEnv("CLIENT_ID"),
		PlatformID: teamID.ToPlatformID(),
	}
}
// IDOUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale individual improvement
func IDOUpdateReminder(teamID models.TeamID, _ bt.Date, userID string) {
	conn := slowlyCreateConnection(teamID)
	err2 := workflows.InvokeWorkflowByPath(exchange.PromptStaleIssues(teamID, userID, issuesUtils.IDO, 7))(conn)
	if err2 != nil {
		log.Printf("ERROR IDOUpdateReminder: %+v\n", err2)
	}
	// ObjectiveProgressUpdateAsk(userID, user.IDOUpdateDueWithinWeek,
		// "You have Individual Development Objective(s) that haven't been updated in last 7 days. Would you like to update them?")
}

// ObjectiveUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale objectives
func ObjectiveUpdateReminder(teamID models.TeamID, date bt.Date, userID string) {
	conn := slowlyCreateConnection(teamID)
	err2 := workflows.InvokeWorkflowByPath(exchange.PromptStaleIssues(teamID, userID, issuesUtils.SObjective, 30))(conn)
	if err2 != nil {
		log.Printf("ERROR ObjectiveUpdateReminder: %+v\n", err2)
	}
	// ObjectiveProgressUpdateAsk(target, user.CapabilityObjectiveUpdateDueWithinMonth,
	// 	"You have Capability Objective(s) that haven't been updated in last 30 days. Would you like to update them?")
}

// InitiativeUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale initiatives
func InitiativeUpdateReminder(teamID models.TeamID, date bt.Date, userID string) {
	conn := slowlyCreateConnection(teamID)
	err2 := workflows.InvokeWorkflowByPath(exchange.PromptStaleIssues(teamID, userID, issuesUtils.Initiative, 7))(conn)
	if err2 != nil {
		log.Printf("ERROR InitiativeUpdateReminder: %+v\n", err2)
	}
	// ObjectiveProgressUpdateAsk(target, user.InitiativeUpdateDueWithinWeek,
	// 	"You have Initiative(s) that haven't been updated in last 7 days. Would you like to update them?")
}

/*
------------------------------------------------------------------------------------
Closeout reminders
------------------------------------------------------------------------------------
*/

// IDOCloseoutReminder is meant to trigger engagements that reminds users
// that they have an IDO due in the coming week and to close it out
func IDOCloseoutReminder(teamID models.TeamID, date bt.Date, target string) {
	conn := slowlyCreateConnection(teamID)
	err2 := workflows.InvokeWorkflowByPath(exchange.PromptStaleIssues(teamID, target, issues.IDO, 7))(conn)
	if err2 != nil {
		log.Printf("ERROR IDOCloseoutReminder: %+v\n", err2)
	}
	// // TODO: dont ask for update, ask for closeout
	// UserCloseObjectivesAsk(teamID, target, "You have IDOs that are due in a week. Would you like to close them out?")
}

// InitiativeCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Initiative due in the coming week and to close it out
func InitiativeCloseoutReminder(teamID models.TeamID, date bt.Date, target string) {
	conn := slowlyCreateConnection(teamID)
	err2 := workflows.InvokeWorkflowByPath(exchange.PromptStaleIssues(teamID, target, issues.Initiative, 7))(conn)
	if err2 != nil {
		log.Printf("ERROR InitiativeCloseoutReminder: %+v\n", err2)
	}
	// UserCloseObjectivesAsk(teamID, target, "You have Initiatives that are due in a week. Would you like to close them out?")
}

// ObjectiveCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Objective due in the coming week and to close it out
func ObjectiveCloseoutReminder(teamID models.TeamID, date bt.Date, target string) {
	conn := slowlyCreateConnection(teamID)
	err2 := workflows.InvokeWorkflowByPath(exchange.PromptStaleIssues(teamID, target, issues.SObjective, 7))(conn)
	if err2 != nil {
		log.Printf("ERROR ObjectiveCloseoutReminder: %+v\n", err2)
	}
	// UserCloseObjectivesAsk(teamID, target, "You have Capability Objectives that are due in a week. Would you like to close them out?")
}

/*
------------------------------------------------------------------------------------
Due date reminders
------------------------------------------------------------------------------------
*/

// IDOReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week
func IDOReminderOfDueDateInMonth(teamID models.TeamID, date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(teamID, target, user.IDOUpdateDueWithinMonth,
		"You have some Individual Development Objectives that are due in 30 days from now. Would you like to update them?")
}

// IDOReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week, month, quarter
func IDOReminderOfDueDateInQaurter(teamID models.TeamID, date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(teamID, target, user.IDOUpdateDueWithinQuarter,
		"You have some Individual Development Objectives that are due in 90 days from now. Would you like to update them?")
}

// InitiativeReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week
func InitiativeReminderOfDueDateInMonth(teamID models.TeamID, date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(teamID, target, user.InitiativeUpdateDueWithinMonth,
		"You have some Initiatives that are due in 30 days from now. Would you like to update them?")
}

// InitiativeReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week, month, quarter
func InitiativeReminderOfDueDateInQaurter(teamID models.TeamID, date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(teamID, target, user.InitiativeUpdateDueWithinQuarter,
		"You have some Initiatives that are due in 90 days from now. Would you like to update them?")
}

// ObjectiveReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week
func ObjectiveReminderOfDueDateInMonth(teamID models.TeamID, date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(teamID, target, user.CapabilityObjectiveUpdateDueWithinMonth,
		"You have some Capability Objectives that are due within a month, but haven't been updated since last week. Would you like to update them?")
}

// ObjectiveReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week, month, quarter
func ObjectiveReminderOfDueDateInQaurter(teamID models.TeamID, date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(teamID, target, user.CapabilityObjectiveUpdateDueWithinQuarter,
		"You have some Capability Objectives that are due in a quarter, but haven't been updated since last month. Would you like to update them?")
}

/*
------------------------------------------------------------------------------------
Coaching feedback reminders
------------------------------------------------------------------------------------
*/

// ReminderToProvideCoachingFeedback is meant to trigger engagements at increasingly rates
// until the end of the quarter to maximize the amount of feedback.
func ReminderToProvideCoachingFeedback(teamID models.TeamID, date bt.Date, target string) {
	year := date.GetYear()
	month := date.GetMonth()
	quarter := date.GetQuarter()

	var userIDs []string // We currently do not save feedback requestsS
	userGivenFeedback, err := coaching.FeedbackGivenForTheQuarter(target, quarter, year, FeedbackTableName,
		FeedbackSourceQuarterYearIndex)
	if err != nil {
		log.Panicf("Could not query user given feedback: %+v", err)
	}
	for _, each := range userGivenFeedback {
		userIDs = append(userIDs, each.Target)
	}

	text := ProvidedFeedbackConfirmationAndSuggestProvidingMoreTemplate(core.Distinct(userIDs))

	// Setting action as 'select'
	mc := models.MessageCallback{
		Module: "coaching",
		Source: target,
		Topic:  "user_feedback",
		Action: "confirm",
		Month:  strconv.Itoa(int(month)),
		Year:   strconv.Itoa(year),
	}
	UserConfirmEng(teamID, mc, string(text), "", false, Dns)
}

func UserSelectEngagement(mc models.MessageCallback, users, filter []string, userID string,
	text ui.RichText, context string, conn daosCommon.DynamoDBConnection) {
	user.UserSelectEng(userID, EngagementTable, conn, mc,
		users, filter, string(text), context, models.UserEngagementCheckWithValue{})
}

/*
------------------------------------------------------------------------------------
Report reminders
------------------------------------------------------------------------------------
*/

// GenerateIndividualReports generates an individual coaching report
func GenerateIndividualReports(teamID models.TeamID, date bt.Date, userID string) {
	err2 := feedbackReportingLambda.GeneratePerformanceReportAndPostToUserAsync(teamID, userID, date.DateToTimeMidnight())
	core.ErrorHandler(err2, Namespace, fmt.Sprintf("Could not invoke %s lambda", FeedbackReportingLambdaName))
}

// utility methods

func AddObjective(userID string, teamID models.TeamID, urgent bool, userIDsToDisplayNames common.UserIDsToDisplayNames) {
	allObjs := objectives.AllUserObjectives(userID, UserObjectivesTable, UserObjectivesUserIDIndex,
		models.IndividualDevelopmentObjective, 0)
	if len(allObjs) == 0 {
		// Create an objective
		year, month := core.CurrentYearMonth()
		mc := models.MessageCallback{Module: "objectives", Source: userID, Topic: "init", Action: "ask",
			Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
		// Posting engagement to user prompting to create an objective
		// TODO: Take InitsAndObjectives from strategy once it's moved there
		objectives.CreateObjectiveEng(EngagementTable, mc,
			objectives.IDOCoaches(userID, teamID, CommunityUsersTable, CommunityUsersCommunityIndex, userIDsToDisplayNames),
			objectives.DevelopmentObjectiveDates(Namespace, core.EmptyString),
			[]ebm.AttachmentActionElementOptionGroup{},
			"You do not have any Individual Development Objectives defined. Do you want to add one?",
			"Create Individual Development Objective", urgent, Dns, models.UserEngagementCheckWithValue{},
			teamID)
	}
}

// EngagePersonalImprovementObjectiveUpdate
func ObjectiveProgressUpdateAsk(teamID models.TeamID, userID, action, message string) { // ViewObjectivesNoProgressThisWeek
	// "There are objectives with no progress this week. Would you like to update them?"
	year, month := core.CurrentYearMonth()
	mc := models.MessageCallback{Module: "objectives", Source: userID, Topic: "init", Action: action,
		Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	UserConfirmEng(teamID, mc, message, "Update Objectives", false, Dns)
}

func UserConfirmEng(teamID models.TeamID, mc models.MessageCallback, title, fallback string,
	urgent bool, dns common.DynamoNamespace) {
	eng := utils.MakeUserEngagement(mc, title, "", fallback, mc.Source,
		userConfirmAttachmentActions(mc),
		[]ebm.AttachmentField{}, urgent, dns.Namespace,
		time.Now().Unix(), models.UserEngagementCheckWithValue{},
		teamID)
	utils.AddEng(eng, EngagementTable, dns.Dynamo, dns.Namespace)
}

func userConfirmAttachmentActions(mc models.MessageCallback) []ebm.AttachmentAction {
	return []ebm.AttachmentAction{
		*models.SimpleAttachAction(mc, models.Now, models.YesLabel),
		*models.SimpleAttachAction(mc, models.Ignore, models.DefaultSkipThisTemplate),
	}
}

func isFeedbackExistsForUser(userID string, date bt.Date) func (daosCommon.DynamoDBConnection) (bool, error) {
	return  func (conn daosCommon.DynamoDBConnection) (res bool, err error) {
		q := date.GetPreviousQuarter()
		y := date.GetPreviousQuarterYear()
		var feedbacks []models.UserFeedback
		feedbacks, err = coaching.FeedbackReceivedForTheQuarter(userID, q, y)(conn)
		res = len(feedbacks) > 0
		return
	}
}

// // DeliverIndividualReports is meant to trigger the engagements that
// // sends out an individual coaching reports to each users.
// func DeliverIndividualReports(teamID models.TeamID, date bt.Date, target string) {
// 	DeliverIndividualReports(target, date)
// }


// DeliverIndividualReports - deliver individual reports. This function should only be
// triggered if the report already exists.
func DeliverIndividualReports(teamID models.TeamID, date bt.Date, userID string) {
	// conn := daosCommon.CreateConnectionGenFromEnv().ForPlatformID(teamID.ToPlatformID())
	// feedbackExist, err2 := isFeedbackExistsForUser(userID, date)(conn)
	// if err2 == nil {
	year := date.GetPreviousQuarterYear()
	month := date.GetMonth()
	mc := models.MessageCallback{Module: "feedback", Source: userID, Topic: "report", Action: ViewCollaborationReport,
		Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	// if feedbackExist {
	GenerateIndividualReports(teamID, date, userID)
	UserConfirmEng(teamID, mc,
		fmt.Sprintf("It's that time of the year to review your collabaration report. Would you like to see it?"),
		"", false, Dns)
	return
}

// NotifyOnAbsentFeedback -
func NotifyOnAbsentFeedback(teamID models.TeamID, date bt.Date, userID string) {
	conn := daosCommon.CreateConnectionGenFromEnv().ForPlatformID(teamID.ToPlatformID())
	feedbackExist, err2 := isFeedbackExistsForUser(userID, date)(conn)
	if err2 == nil {
		if !feedbackExist {
			year := date.GetPreviousQuarterYear()
			month := date.GetMonth()

			mc := models.MessageCallback{Module: "feedback", Source: userID, Topic: "report", Action: ViewCollaborationReport,
				Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
			eng := utils.MakeUserEngagement(mc, 
				"You did not receive any feedback last quarter so I do not have a report for you. "+
				"Talk with your colleagues and consider using the Request Feedback feature this quarter.",
				"", "No feedback last quarter", mc.Source,
				[]ebm.AttachmentAction{}, []ebm.AttachmentField{}, 
				false, "NotifyOnAbsentFeedback",
				time.Now().Unix(), models.UserEngagementCheckWithValue{},
				teamID,
			)
			utils.AddEng(eng, EngagementTable, conn.Dynamo, "NotifyOnAbsentFeedback")
		} else {
			log.Printf("WARN: NotifyOnAbsentFeedback but feedback exists for user %s\n", userID)
		}
	}	
	if err2 != nil {
		log.Printf("IGNORING ERROR in NotifyOnAbsentFeedback: %+v\n", err2)
	}
}

// func UserCloseObjectivesAsk(teamID models.TeamID, userID, text string) {
// 	if len(openIDOsForUser(userID)) > 0 {
// 		// View IDOs

// 		// Create an objective
// 		year, month := core.CurrentYearMonth()
// 		mc := models.MessageCallback{Module: "objectives", Source: userID,
// 			Topic: "init", Action: ViewOpenObjectives,
// 			Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
// 		UserConfirmEng(teamID, mc, text, "", true, Dns)
// 	}
// }

func openIDOsForUser(userID string) []models.UserObjective {
	// TODO: Include quarter, year in the query
	allObjs := objectives.AllUserObjectives(userID, UserObjectivesTable, UserObjectivesUserIDIndex,
		models.IndividualDevelopmentObjective, 0)
	var openObjs []models.UserObjective
	for _, each := range allObjs {
		if each.Completed == 0 {
			openObjs = append(openObjs, each)
		}
	}
	return openObjs
}

// ProduceAndDeliverIndividualReportsOrNotifyOnAbsentFeedback -
func ProduceAndDeliverIndividualReportsOrNotifyOnAbsentFeedback(teamID models.TeamID, date bt.Date, userID string) {
	conn := daosCommon.CreateConnectionGenFromEnv().ForPlatformID(teamID.ToPlatformID())
	feedbackExist, err2 := isFeedbackExistsForUser(userID, date)(conn)
	if err2 == nil {
		if feedbackExist {
			GenerateIndividualReports(teamID, date, userID)
			DeliverIndividualReports(teamID, date, userID)
		} else {
			NotifyOnAbsentFeedback(teamID, date, userID)
		}
	}
	if err2 != nil {
		log.Printf("IGNORING ERROR in ProduceAndDeliverIndividualReportsOrNotifyOnAbsentFeedback: %+v\n", err2)
	}
}

// DeliverIndividualReportsOrNotifyOnAbsentFeedback -
func DeliverIndividualReportsOrNotifyOnAbsentFeedback(teamID models.TeamID, date bt.Date, userID string) {
	conn := daosCommon.CreateConnectionGenFromEnv().ForPlatformID(teamID.ToPlatformID())
	feedbackExist, err2 := isFeedbackExistsForUser(userID, date)(conn)
	if err2 == nil {
		if feedbackExist {
			DeliverIndividualReports(teamID, date, userID)
		} else {
			NotifyOnAbsentFeedback(teamID, date, userID)
		}
	}
	if err2 != nil {
		log.Printf("IGNORING ERROR in ProduceAndDeliverIndividualReportsOrNotifyOnAbsentFeedback: %+v\n", err2)
	}
}
