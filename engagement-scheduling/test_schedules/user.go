package test_schedules

import (
	"github.com/adaptiveteam/adaptive/adaptive-checks"
	bt "github.com/adaptiveteam/adaptive/business-time"
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
func IDOCreateReminder(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		"Reminder to create Individual Development Objectives",
		!fc.IDOsExistForMe(),
	)
	return rv
}

/*
------------------------------------------------------------------------------------
Update Reminders
------------------------------------------------------------------------------------
*/

// IDOUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale individual improvement
func IDOUpdateReminder(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	rv = utils.ScheduleEntry(
		"Reminder to update Individual Development Objectives",
		fc.StaleIDOsExistForMe(),
	)
	return rv
}

// ObjectiveUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale objectives
func ObjectiveUpdateReminder(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	rv = utils.ScheduleEntry(
		"Reminder to update Objectives",
		fc.StaleObjectivesExistForMe(),
	)
	return rv
}

// InitiativeUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale initiatives
func InitiativeUpdateReminder(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	rv = utils.ScheduleEntry(
		"Reminder to update Initiatives",
		fc.StaleInitiativesExistForMe(),
	)
	return rv
}

/*
------------------------------------------------------------------------------------
Closeout reminders
------------------------------------------------------------------------------------
*/

// IDOCloseoutReminder is meant to trigger engagements that reminds users
// that they have an IDO due in the coming week and to close it out
func IDOCloseoutReminder(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		"Reminder to closeout individual development objective",
		fc.IDOsDueWithinTheWeek(),
	)
	return rv
}

// InitiativeCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Initiative due in the coming week and to close it out
func InitiativeCloseoutReminder(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		"Reminder to closeout Initiative",
		fc.InitiativesDueWithinTheWeek(),
	)
	return rv
}

// ObjectiveCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Objective due in the coming week and to close it out
func ObjectiveCloseoutReminder(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		"Reminder to closeout Objective",
		fc.ObjectivesDueWithinTheWeek(),
	)
	return rv
}

/*
------------------------------------------------------------------------------------
Due date reminders
------------------------------------------------------------------------------------
*/

// IDOReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week
func IDOReminderOfDueDateInMonth(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		"Reminder that an individual development objective is due within a month",
		fc.IDOsDueWithinTheMonth(),
	)
	return rv
}

// IDOReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week, month, quarter
func IDOReminderOfDueDateInQaurter(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		"Reminder that an IDO is due within a qaurter",
		fc.IDOsDueWithinTheQuarter(),
	)
	return rv
}

// InitiativeReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week
func InitiativeReminderOfDueDateInMonth(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		"Reminder that an individual development objective is due within a month",
		fc.InitiativesDueWithinTheMonth(),
	)
	return rv
}

// InitiativeReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week, month, quarter
func InitiativeReminderOfDueDateInQaurter(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		"Reminder that an Initiative is due within a quarter",
		fc.InitiativesDueWithinTheQuarter(),
	)
	return rv
}

// ObjectiveReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week
func ObjectiveReminderOfDueDateInMonth(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		"Reminder that an individual development objective is due within a month",
		fc.ObjectivesDueWithinTheMonth(),
	)
	return rv
}

// ObjectiveReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week, month, quarter
func ObjectiveReminderOfDueDateInQaurter(fc adaptive_checks.TypedProfile, _ bt.Date) (rv string) {
	rv = utils.ScheduleEntry(
		"Reminder that an Objective is due within a qaurter",
		fc.ObjectivesDueWithinTheQuarter(),
	)
	return rv
}

/*
------------------------------------------------------------------------------------
Coaching feedback reminders
------------------------------------------------------------------------------------
*/

// ReminderToProvideCoachingFeedback is meant to trigger engagements at increasingly rates
// until the end of the quarter to maximize the amount of feedback.
func ReminderToProvideCoachingFeedback(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	WeekOne := date.GetDayOfWeekInQuarter(-4, bt.Monday) == date
	WeekTwo := date.GetDayOfWeekInQuarter(-3, bt.Monday) == date
	WeekThree := date.GetDayOfWeekInQuarter(-2, bt.Monday) == date

	begin := date.GetDayOfWeekInQuarter(-1, bt.Monday)
	end := date.GetLastDayOfQuarter().GetWeekDay(false)
	Daily := date.DateAfter(begin, true) && date.DateBefore(end, true)
	feedbackCycle := Daily || WeekOne || WeekTwo || WeekThree
	if feedbackCycle {
		rv := utils.ScheduleEntry(
			"Reminder to provide coaching feedback to colleagues",
			fc.FeedbackGivenThisQuarter(),
		)

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
func ProduceIndividualReports(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	// Last business day of the first week of the first quarter
	rv = utils.ScheduleEntry(
		"Produce and deliver individual reports",
		date.GetDayOfWeekInQuarter(1, bt.Friday) == date &&
		fc.CollaborationReportExists(),
	)
	return rv
}
