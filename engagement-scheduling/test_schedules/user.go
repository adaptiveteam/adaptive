package test_schedules

import (
	fcn "github.com/adaptiveteam/adaptive/adaptive-check-function-names"
	bt "github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/checks"
	utils "github.com/adaptiveteam/adaptive/engagement-scheduling"
)

/*
------------------------------------------------------------------------------------
IDO Creation reminders
------------------------------------------------------------------------------------
*/

// CheckIDOCreateReminder is meant to trigger the engagements that
// reminds the user to create personal improvement objects in the event that they have
// not created any.
func IDOCreateReminder(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		fc,
		"Reminder to create Individual Improvement Objectives",
		).AddScheduleFunctionCheck(
		fcn.IDOsExistForMe,
		false,
		).Message
	return rv
}

/*
------------------------------------------------------------------------------------
Update Reminders
------------------------------------------------------------------------------------
*/

// IDOUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale individual improvement
func IDOUpdateReminder(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	rv = utils.ScheduleEntry(
		fc,
		"Reminder to update Individual Improvement Objectives",
	).AddScheduleFunctionCheck(
		fcn.StaleIDOsExistForMe,
		true,
	).Message
	return rv
}

// ObjectiveUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale objectives
func ObjectiveUpdateReminder(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	rv = utils.ScheduleEntry(
		fc,
		"Reminder to update Objectives",
	).AddScheduleFunctionCheck(
		fcn.StaleObjectivesExistForMe,
		true,
	).Message
	return rv
}

// InitiativeUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale initiatives
func InitiativeUpdateReminder(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	rv = utils.ScheduleEntry(
		fc,
		"Reminder to update Initiatives",
	).AddScheduleFunctionCheck(
		fcn.StaleInitiativesExistForMe,
		true,
	).Message
	return rv
}

/*
------------------------------------------------------------------------------------
Closeout reminders
------------------------------------------------------------------------------------
*/

// IDOCloseoutReminder is meant to trigger engagements that reminds users
// that they have an IDO due in the coming week and to close it out
func IDOCloseoutReminder(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		fc,
		"Reminder to closeout individual development objective",
	).AddScheduleFunctionCheck(
		fcn.IDOsDueWithinTheWeek,
		true,
	).Message
	return rv
}

// InitiativeCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Initiative due in the coming week and to close it out
func InitiativeCloseoutReminder(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		fc,
		"Reminder to closeout Initiative",
	).AddScheduleFunctionCheck(
		fcn.InitiativesDueWithinTheWeek,
		true,
	).Message
	return rv
}

// ObjectiveCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Objective due in the coming week and to close it out
func ObjectiveCloseoutReminder(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		fc,
		"Reminder to closeout Objective",
	).AddScheduleFunctionCheck(
		fcn.ObjectivesDueWithinTheWeek,
		true,
	).Message
	return rv
}

/*
------------------------------------------------------------------------------------
Due date reminders
------------------------------------------------------------------------------------
*/

// IDOReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week
func IDOReminderOfDueDateInMonth(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		fc,
		"Reminder that an individual development objective is due within a month",
	).AddScheduleFunctionCheck(
		fcn.IDOsDueWithinTheMonth,
		true,
	).Message
	return rv
}

// IDOReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week, month, quarter
func IDOReminderOfDueDateInQaurter(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		fc,
		"Reminder that an IDO is due within a qaurter",
	).AddScheduleFunctionCheck(
		fcn.IDOsDueWithinTheQuarter,
		true,
	).Message
	return rv
}

// InitiativeReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week
func InitiativeReminderOfDueDateInMonth(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		fc,
		"Reminder that an individual development objective is due within a month",
	).AddScheduleFunctionCheck(
		fcn.InitiativesDueWithinTheMonth,
		true,
	).Message
	return rv
}

// InitiativeReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week, month, quarter
func InitiativeReminderOfDueDateInQaurter(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		fc,
		"Reminder that an Initiative is due within a quarter",
	).AddScheduleFunctionCheck(
		fcn.InitiativesDueWithinTheQuarter,
		true,
	).Message
	return rv
}

// ObjectiveReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week
func ObjectiveReminderOfDueDateInMonth(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		fc,
		"Reminder that an individual development objective is due within a month",
	).AddScheduleFunctionCheck(
		fcn.ObjectivesDueWithinTheMonth,
		true,
	).Message
	return rv
}

// ObjectiveReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week, month, quarter
func ObjectiveReminderOfDueDateInQaurter(fc checks.CheckResultMap, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		fc,
		"Reminder that an Objective is due within a qaurter",
	).AddScheduleFunctionCheck(
		fcn.ObjectivesDueWithinTheQuarter,
		true,
	).Message
	return rv
}

/*
------------------------------------------------------------------------------------
Coaching feedback reminders
------------------------------------------------------------------------------------
*/

// ReminderToProvideCoachingFeedback is meant to trigger engagements at increasingly rates
// until the end of the quarter to maximize the amount of feedback.
func ReminderToProvideCoachingFeedback(fc checks.CheckResultMap, date bt.Date) (rv string) {
	WeekOne := date.GetDayOfWeekInQuarter(-4, bt.Monday) == date
	WeekTwo := date.GetDayOfWeekInQuarter(-3, bt.Monday) == date
	WeekThree := date.GetDayOfWeekInQuarter(-2, bt.Monday) == date

	begin := date.GetDayOfWeekInQuarter(-1, bt.Monday)
	end := date.GetLastDayOfQuarter().GetWeekDay(false)
	Daily := date.DateAfter(begin, true) && date.DateBefore(end, true)
	feedbackCycle := Daily || WeekOne || WeekTwo || WeekThree
	if feedbackCycle {
		rv := utils.ScheduleEntry(
			fc,
			"Reminder to provide coaching feedback to colleagues",
		).AddScheduleFunctionCheck(
			fcn.FeedbackGivenThisQuarter,
			true,
		).Message

		if len(rv) == 0 {
			rv = "Reminder to provide coaching feedback to additional colleagues"
		}
	}
	return rv
}

/*
------------------------------------------------------------------------------------
Report reminders
------------------------------------------------------------------------------------
*/

// ProduceIndividualReports is meant to trigger the engagements that
// sends out a the individual coaching reports to each users.
func ProduceIndividualReports(fc checks.CheckResultMap, date bt.Date) (rv string) {
	// Last business day of the first week of the first quarter
	rv = utils.ScheduleEntry(
		fc,
		"Produce and deliver individual reports",
	).AddScheduleBooleanCheck(
		date.GetDayOfWeekInQuarter(1, bt.Friday) == date,
		true,
	).AddScheduleFunctionCheck(
		fcn.CollaborationReportExists,
		true,
	).Message
	return rv
}
