package engagements

import (
	"github.com/adaptiveteam/adaptive/workflows"
	// wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
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

// CheckIDOCreateReminder is meant to trigger the engagements that
// reminds the user to create personal improvement objects in the event that they have
// not created any.
func IDOCreateReminder(date bt.Date, target string) {
	log.Println(fmt.Sprintf("Checking IDOCreateReminder for user: %s", target))
	ut := UserToken(target)
	AddObjective(target, ut.PlatformIDUnsafe(), false, utilsUser.UserIDsToDisplayNamesUnsafe(UserDAO))

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
	platformID := strategy.UserIDToPlatformID(UserDAO)(target)
	UserConfirmEng(EngagementTable, platformID, mc, string(text), "", false, Dns)
}

func UserSelectEngagement(mc models.MessageCallback, users, filter []string, userID string,
	text ui.RichText, context string) {
	platformID := UserDAO.ReadUnsafe(userID).PlatformID
	user.UserSelectEng(userID, EngagementTable, models.PlatformID(platformID), UserDAO, mc,
		users, filter, string(text), context, models.UserEngagementCheckWithValue{})
}

/*
------------------------------------------------------------------------------------
Report reminders
------------------------------------------------------------------------------------
*/

// ProduceIndividualReports is meant to trigger the engagements that
// sends out a the individual coaching reports to each users.
func ProduceIndividualReports(date bt.Date, target string) {
	UserCollaborationReportAsk(target, date)
}

// utility methods

func AddObjective(userID string, platformID models.PlatformID, urgent bool, userIDsToDisplayNames common.UserIDsToDisplayNames) {
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
			objectives.IDOCoaches(userID, platformID, CommunityUsersTable, CommunityUsersCommunityIndex, userIDsToDisplayNames),
			objectives.DevelopmentObjectiveDates(Namespace, core.EmptyString),
			[]ebm.AttachmentActionElementOptionGroup{},
			"You do not have any Individual Development Objectives defined. Do you want to add one?",
			"Create Individual Development Objective", urgent, Dns, models.UserEngagementCheckWithValue{},
			models.PlatformID(platformID))
	}
}

// EngagePersonalImprovementObjectiveUpdate
func ObjectiveProgressUpdateAsk(userID, action, message string) { // ViewObjectivesNoProgressThisWeek
	// "There are objectives with no progress this week. Would you like to update them?"
	platformID := strategy.UserIDToPlatformID(UserDAO)(userID)
	year, month := core.CurrentYearMonth()
	mc := models.MessageCallback{Module: "objectives", Source: userID, Topic: "init", Action: action,
		Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	UserConfirmEng(EngagementTable, platformID, mc, message, "Update Objectives", false, Dns)
}

func UserConfirmEng(table string, platformID models.PlatformID, mc models.MessageCallback, title, fallback string,
	urgent bool, dns common.DynamoNamespace) {
	eng := utils.MakeUserEngagement(mc, title, "", fallback, mc.Source,
		userConfirmAttachmentActions(mc),
		[]ebm.AttachmentField{}, urgent, dns.Namespace,
		time.Now().Unix(), models.UserEngagementCheckWithValue{},
		platformID)
	utils.AddEng(eng, table, dns.Dynamo, dns.Namespace)
}

func userConfirmAttachmentActions(mc models.MessageCallback) []ebm.AttachmentAction {
	return []ebm.AttachmentAction{
		*models.SimpleAttachAction(mc, models.Now, models.YesLabel),
		*models.SimpleAttachAction(mc, models.Ignore, models.DefaultSkipThisTemplate),
	}
}

// EngageProduceIndividualReports
func UserCollaborationReportAsk(userID string, date bt.Date) {
	platformID := strategy.UserIDToPlatformID(UserDAO)(userID)
	userEngageByt, _ := json.Marshal(models.UserEngage{UserId: userID, IsNew: false, Update: true, Date: date.DateToString(time.RFC3339)})
	_, err := L.InvokeFunction(FeedbackReportLambda, userEngageByt, true)
	core.ErrorHandler(err, Namespace, fmt.Sprintf("Could not invoke %s lambda", FeedbackReportLambda))

	year := date.GetPreviousQuarterYear()
	month := date.GetMonth()
	mc := models.MessageCallback{Module: "feedback", Source: userID, Topic: "report", Action: ViewCollaborationReport,
		Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	UserConfirmEng(EngagementTable, platformID, mc,
		fmt.Sprintf("It's that time of the year to review your collabaration report. Would you like to see it?"),
		"", false, Dns)
}

func UserCloseObjectivesAsk(userID, text string) {
	if len(userOpenObjectives(userID)) > 0 {
		platformID := strategy.UserIDToPlatformID(UserDAO)(userID)
		// Create an objective
		year, month := core.CurrentYearMonth()
		mc := models.MessageCallback{Module: "objectives", Source: userID, Topic: "init", Action: ViewOpenObjectives,
			Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
		UserConfirmEng(EngagementTable, platformID, mc, text, "", true, Dns)
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
