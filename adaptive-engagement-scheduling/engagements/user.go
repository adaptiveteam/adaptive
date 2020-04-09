package engagements

import (
	"github.com/adaptiveteam/adaptive/workflows"
	"github.com/adaptiveteam/adaptive/workflows/exchange"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	. "github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	bt "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	issuesUtils "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
)

/*
------------------------------------------------------------------------------------
IDO Creation reminders
------------------------------------------------------------------------------------
*/

func globalConnectionGen() daosCommon.DynamoDBConnectionGen {
	return daosCommon.DynamoDBConnectionGen{
		Dynamo: D,
		TableNamePrefix: ClientID,
	}
}

func readUser(userID string) (models.User, error) {
	conn := globalConnectionGen()
	dao := utilsUser.DAOFromConnectionGen(conn)
	return dao.Read(userID)
}
// IDOCreateReminder is meant to trigger the engagements that
// reminds the user to create personal improvement objects in the event that they have
// not created any.
func IDOCreateReminder(date bt.Date, target string) {
	log.Println(fmt.Sprintf("Checking IDOCreateReminder for user: %s", target))
	user, err2 := readUser(target)
	core.ErrorHandler(err2, "IDOCreateReminder", "readUser")
	AddObjective(target, models.ParseTeamID(user.PlatformID), false, utilsUser.UserIDsToDisplayNamesUnsafe(UserDAO))

}

/*
------------------------------------------------------------------------------------
Update Reminders
------------------------------------------------------------------------------------
*/

// slowlyCreateConnection - internally requests DB, reads user and retrieves PlatformID
func slowlyCreateConnection(userID string) daosCommon.DynamoDBConnection {
	return daosCommon.DynamoDBConnection {
		Dynamo: common.DeprecatedGetGlobalDns().Dynamo,
		ClientID: utils.NonEmptyEnv("CLIENT_ID"),
		PlatformID: strategy.UserIDToPlatformID(UserDAO)(userID),
	}
}
// IDOUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale individual improvement
func IDOUpdateReminder(_ bt.Date, targetUserID string) {
	conn := slowlyCreateConnection(targetUserID)
	err2 := workflows.InvokeWorkflowByPath(exchange.PromptStaleIssues(targetUserID, issuesUtils.IDO, 7))(conn)
	if err2 != nil {
		log.Printf("ERROR IDOUpdateReminder: %+v", err2)
	}
	// ObjectiveProgressUpdateAsk(targetUserID, user.IDOUpdateDueWithinWeek,
		// "You have Individual Development Objective(s) that haven't been updated in last 7 days. Would you like to update them?")
}

// ObjectiveUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale objectives
func ObjectiveUpdateReminder(date bt.Date, targetUserID string) {
	conn := slowlyCreateConnection(targetUserID)
	err2 := workflows.InvokeWorkflowByPath(exchange.PromptStaleIssues(targetUserID, issuesUtils.SObjective, 30))(conn)
	if err2 != nil {
		log.Printf("ERROR ObjectiveUpdateReminder: %+v", err2)
	}
	// ObjectiveProgressUpdateAsk(target, user.CapabilityObjectiveUpdateDueWithinMonth,
	// 	"You have Capability Objective(s) that haven't been updated in last 30 days. Would you like to update them?")
}

// InitiativeUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale initiatives
func InitiativeUpdateReminder(date bt.Date, targetUserID string) {
	conn := slowlyCreateConnection(targetUserID)
	err2 := workflows.InvokeWorkflowByPath(exchange.PromptStaleIssues(targetUserID, issuesUtils.Initiative, 7))(conn)
	if err2 != nil {
		log.Printf("ERROR InitiativeUpdateReminder: %+v", err2)
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
func IDOCloseoutReminder(date bt.Date, target string) {
	// TODO: dont ask for update, ask for closeout
	UserCloseObjectivesAsk(target, "You have IDOs that are due in a week. Would you like to close them out?")
}

// InitiativeCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Initiative due in the coming week and to close it out
func InitiativeCloseoutReminder(date bt.Date, target string) {
	UserCloseObjectivesAsk(target, "You have Initiatives that are due in a week. Would you like to close them out?")
}

// ObjectiveCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Objective due in the coming week and to close it out
func ObjectiveCloseoutReminder(date bt.Date, target string) {
	UserCloseObjectivesAsk(target, "You have Capability Objectives that are due in a week. Would you like to close them out?")
}

/*
------------------------------------------------------------------------------------
Due date reminders
------------------------------------------------------------------------------------
*/

// IDOReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week
func IDOReminderOfDueDateInMonth(date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(target, user.IDOUpdateDueWithinMonth,
		"You have some Individual Development Objectives that are due in 30 days from now. Would you like to update them?")
}

// IDOReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week, month, quarter
func IDOReminderOfDueDateInQaurter(date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(target, user.IDOUpdateDueWithinQuarter,
		"You have some Individual Development Objectives that are due in 90 days from now. Would you like to update them?")
}

// InitiativeReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week
func InitiativeReminderOfDueDateInMonth(date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(target, user.InitiativeUpdateDueWithinMonth,
		"You have some Initiatives that are due in 30 days from now. Would you like to update them?")
}

// InitiativeReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week, month, quarter
func InitiativeReminderOfDueDateInQaurter(date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(target, user.InitiativeUpdateDueWithinQuarter,
		"You have some Initiatives that are due in 90 days from now. Would you like to update them?")
}

// ObjectiveReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week
func ObjectiveReminderOfDueDateInMonth(date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(target, user.CapabilityObjectiveUpdateDueWithinMonth,
		"You have some Capability Objectives that are due within a month, but haven't been updated since last week. Would you like to update them?")
}

// ObjectiveReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week, month, quarter
func ObjectiveReminderOfDueDateInQaurter(date bt.Date, target string) {
	ObjectiveProgressUpdateAsk(target, user.CapabilityObjectiveUpdateDueWithinQuarter,
		"You have some Capability Objectives that are due in a quarter, but haven't been updated since last month. Would you like to update them?")
}

/*
------------------------------------------------------------------------------------
Coaching feedback reminders
------------------------------------------------------------------------------------
*/

// ReminderToProvideCoachingFeedback is meant to trigger engagements at increasingly rates
// until the end of the quarter to maximize the amount of feedback.
func ReminderToProvideCoachingFeedback(date bt.Date, target string) {
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
	teamID := strategy.UserIDToTeamID(UserDAO)(target)
	UserConfirmEng(teamID, mc, string(text), "", false, Dns)
}

func UserSelectEngagement(mc models.MessageCallback, users, filter []string, userID string,
	text ui.RichText, context string) {
	teamID := models.ParseTeamID(UserDAO.ReadUnsafe(userID).PlatformID)
	user.UserSelectEng(userID, EngagementTable, teamID, UserDAO, mc,
		users, filter, string(text), context, models.UserEngagementCheckWithValue{})
}

/*
------------------------------------------------------------------------------------
Report reminders
------------------------------------------------------------------------------------
*/

// GenerateIndividualReports generates an individual coaching report
func GenerateIndividualReports(date bt.Date, target string) {
	// year := date.GetPreviousQuarterYear()
	// month := date.GetMonth()
	userID := target
	// mc := models.MessageCallback{Module: "feedback", Source: userID, Topic: "report", Action: ViewCollaborationReport,
	// 	Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	userEngageByt, _ := json.Marshal(models.UserEngage{UserID: userID, 
		IsNew: false, Update: true, Date: date.DateToString(time.RFC3339)})
	_, err := L.InvokeFunction(FeedbackReportLambda, userEngageByt, true)
	core.ErrorHandler(err, Namespace, fmt.Sprintf("Could not invoke %s lambda", FeedbackReportLambda))
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
			models.TeamID(teamID))
	}
}

// EngagePersonalImprovementObjectiveUpdate
func ObjectiveProgressUpdateAsk(userID, action, message string) { // ViewObjectivesNoProgressThisWeek
	// "There are objectives with no progress this week. Would you like to update them?"
	teamID := strategy.UserIDToTeamID(UserDAO)(userID)
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
// func DeliverIndividualReports(date bt.Date, target string) {
// 	DeliverIndividualReports(target, date)
// }


// DeliverIndividualReports - deliver individual reports. This function should only be
// triggered if the report already exists.
func DeliverIndividualReports(date bt.Date, userID string) {
	teamID := strategy.UserIDToTeamID(UserDAO)(userID)
	// conn := daosCommon.CreateConnectionGenFromEnv().ForPlatformID(teamID.ToPlatformID())
	// feedbackExist, err2 := isFeedbackExistsForUser(userID, date)(conn)
	// if err2 == nil {
	year := date.GetPreviousQuarterYear()
	month := date.GetMonth()
	mc := models.MessageCallback{Module: "feedback", Source: userID, Topic: "report", Action: ViewCollaborationReport,
		Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	// if feedbackExist {
	engage := models.UserEngage{UserID: userID, 
		IsNew: false, Update: true, Date: date.DateToString(time.RFC3339)}
	userEngageByt, err2 := json.Marshal(engage)
	if err2 == nil {
		_, err2 = L.InvokeFunction(FeedbackReportLambda, userEngageByt, true)
		if err2 == nil {
			UserConfirmEng(teamID, mc,
				fmt.Sprintf("It's that time of the year to review your collabaration report. Would you like to see it?"),
				"", false, Dns)
		}
	}
	if err2 != nil {
		log.Printf("IGNORING ERROR in DeliverIndividualReports: %+v\n", err2)
	}
	return
}

// NotifyOnAbsentFeedback -
func NotifyOnAbsentFeedback(date bt.Date, userID string) {
	teamID := strategy.UserIDToTeamID(UserDAO)(userID)
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
			log.Printf("WARN: NotifyOnAbsentFeedback but feedback exists for user %s", userID)
		}
	}	
	if err2 != nil {
		log.Printf("IGNORING ERROR in NotifyOnAbsentFeedback: %+v\n", err2)
	}
}

func UserCloseObjectivesAsk(userID, text string) {
	if len(userOpenObjectives(userID)) > 0 {
		teamID := strategy.UserIDToTeamID(UserDAO)(userID)
		// Create an objective
		year, month := core.CurrentYearMonth()
		mc := models.MessageCallback{Module: "objectives", Source: userID, Topic: "init", Action: ViewOpenObjectives,
			Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
		UserConfirmEng(teamID, mc, text, "", true, Dns)
	}
}

func userOpenObjectives(userID string) []models.UserObjective {
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
func ProduceAndDeliverIndividualReportsOrNotifyOnAbsentFeedback(date bt.Date, userID string) {
	teamID := strategy.UserIDToTeamID(UserDAO)(userID)
	conn := daosCommon.CreateConnectionGenFromEnv().ForPlatformID(teamID.ToPlatformID())
	feedbackExist, err2 := isFeedbackExistsForUser(userID, date)(conn)
	if err2 == nil {
		if feedbackExist {
			GenerateIndividualReports(date, userID)
			DeliverIndividualReports(date, userID)
		} else {
			NotifyOnAbsentFeedback(date, userID)
		}
	}
	if err2 != nil {
		log.Printf("IGNORING ERROR in ProduceAndDeliverIndividualReportsOrNotifyOnAbsentFeedback: %+v\n", err2)
	}
}
