package test_engagements

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"log"

	bt "github.com/adaptiveteam/adaptive/business-time"
)

/*
------------------------------------------------------------------------------------
IDO Creation reminders
------------------------------------------------------------------------------------
*/

// EnableDebugPrint should be false to reduce log bloating.
var EnableDebugPrint = false

// Println prints if debug is enabled
func Println(s string) {
	if EnableDebugPrint {
		log.Println(s)
	}
}

// IDOCreateReminder is meant to trigger the engagements that
// reminds the user to create personal improvement objects in the event that they have
// not created any.
func IDOCreateReminder(teamID models.TeamID, date bt.Date, target string) {
	Println("Remember to create your IDO!")
}

/*
------------------------------------------------------------------------------------
Update Reminders
------------------------------------------------------------------------------------
*/

// IDOUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale individual improvement
func IDOUpdateReminder(teamID models.TeamID, date bt.Date, target string) {
	Println("Update your IDO!")
}

// ObjectiveUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale objectives
func ObjectiveUpdateReminder(teamID models.TeamID, date bt.Date, target string) {
	Println("Update you Objectives")
}

// InitiativeUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale initiatives
func InitiativeUpdateReminder(teamID models.TeamID, date bt.Date, target string) {
	Println("Update you Initiatives")
}

/*
------------------------------------------------------------------------------------
Closeout reminders
------------------------------------------------------------------------------------
*/

// IDOCloseoutReminder is meant to trigger engagements that reminds users
// that they have an IDO due in the coming week and to close it out
func IDOCloseoutReminder(teamID models.TeamID, date bt.Date, target string) {
	Println("Close out your IDO's")
}

// InitiativeCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Initiative due in the coming week and to close it out
func InitiativeCloseoutReminder(teamID models.TeamID, date bt.Date, target string) {
	Println("Close out your Initiatives")
}

// ObjectiveCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Objective due in the coming week and to close it out
func ObjectiveCloseoutReminder(teamID models.TeamID, date bt.Date, target string) {
	Println("Close out your Objectives")
}

/*
------------------------------------------------------------------------------------
Due date reminders
------------------------------------------------------------------------------------
*/

// IDOReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week
func IDOReminderOfDueDateInMonth(teamID models.TeamID, date bt.Date, target string) {
	Println("IDO is due in a month")
}

// IDOReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week, month, quarter
func IDOReminderOfDueDateInQaurter(teamID models.TeamID, date bt.Date, target string) {
	Println("IDO is due in a quarter")
}

// InitiativeReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week
func InitiativeReminderOfDueDateInMonth(teamID models.TeamID, date bt.Date, target string) {
	Println("Initative is due in a month")
}

// InitiativeReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week, month, quarter
func InitiativeReminderOfDueDateInQaurter(teamID models.TeamID, date bt.Date, target string) {
	Println("Initative is due in a quarter")
}

// ObjectiveReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week
func ObjectiveReminderOfDueDateInMonth(teamID models.TeamID, date bt.Date, target string) {
	Println("Objective is due in a month")
}

// ObjectiveReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week, month, quarter
func ObjectiveReminderOfDueDateInQaurter(teamID models.TeamID, date bt.Date, target string) {
	Println("Objective is due in a quater")
}

/*
------------------------------------------------------------------------------------
Coaching feedback reminders
------------------------------------------------------------------------------------
*/

// ReminderToProvideCoachingFeedback is meant to trigger engagements at increasingly rates
// until the end of the quarter to maximize the amount of feedback.
func ReminderToProvideCoachingFeedback(teamID models.TeamID, date bt.Date, target string) {
	Println("Provide your colleagues feedback")
}

/*
------------------------------------------------------------------------------------
Report reminders
------------------------------------------------------------------------------------
*/

// ProduceIndividualReports is meant to trigger the engagements that
// sends out a the individual coaching reports to each users.
func ProduceIndividualReports(teamID models.TeamID, date bt.Date, target string) {
	Println("Your report s ready")
}
