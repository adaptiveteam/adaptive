package schedules

import (
	"github.com/adaptiveteam/adaptive/adaptive-checks"
	bt "github.com/adaptiveteam/adaptive/business-time"
	utils "github.com/adaptiveteam/adaptive/engagement-scheduling"
)

// DebugReminder is meant to always trigger the debug engagements
func DebugReminder(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	return "Debug reminder"
}

/*
------------------------------------------------------------------------------------
IDO Creation reminders
------------------------------------------------------------------------------------
*/

// IDOCreateReminder is meant to trigger the engagements that
// reminds the user to create personal improvement objects in the event that they have
// not created any.
func IDOCreateReminder(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	WeekOne := date.GetDayOfWeekInQuarter(2, bt.Monday) == date
	WeekTwo := date.GetDayOfWeekInQuarter(3, bt.Monday) == date

	begin := date.GetDayOfWeekInQuarter(4, bt.Monday)
	end := date.GetDayOfWeekInQuarter(4, bt.Friday)
	Daily := date.DateAfter(begin, true) && date.DateBefore(end, true)

	rv = utils.ScheduleEntry(
			"Create Individual Development Objectives",
			(Daily || WeekOne || WeekTwo) &&
			fc.TeamValuesExist() &&
			fc.CanBeNudgedForIDO() &&
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
func IDOUpdateReminder(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	if date.DayOfWeek() == bt.Friday {
		rv = utils.ScheduleEntry(
			"Update Individual Development Objectives",
			fc.StaleIDOsExistForMe(),
		)
	}
	return rv
}

// ObjectiveUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale objectives
func ObjectiveUpdateReminder(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	if date.DayOfWeek() == bt.Friday {
		rv = utils.ScheduleEntry(
			"Update Objectives",
			fc.StaleObjectivesExistForMe(),
		)
	}
	return rv
}

// InitiativeUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale initiatives
func InitiativeUpdateReminder(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	if date.DayOfWeek() == bt.Friday {
		rv = utils.ScheduleEntry(
			"Update Initiatives",
			fc.StaleInitiativesExistForMe(),
		)
	}
	return rv
}

/*
------------------------------------------------------------------------------------
Closeout reminders
------------------------------------------------------------------------------------
*/

// IDOCloseoutReminder is meant to trigger engagements that reminds users
// that they have an IDO due in the coming week and to close it out
func IDOCloseoutReminder(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	if date.IsWeekDay() {
		rv = utils.ScheduleEntry(
			"Closeout individual development objective",
			fc.IDOsDueWithinTheWeek(),
		)
	}
	return rv
}

// InitiativeCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Initiative due in the coming week and to close it out
func InitiativeCloseoutReminder(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	if date.IsWeekDay() {
		rv = utils.ScheduleEntry(
			"Closeout Initiative",
			fc.InitiativesDueWithinTheWeek(),
		)
	}
	return rv
}

// ObjectiveCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Objective due in the coming week and to close it out
func ObjectiveCloseoutReminder(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	if date.IsWeekDay() {
		rv = utils.ScheduleEntry(
			"Closeout Objective",
			fc.ObjectivesDueWithinTheWeek(),
		)
	}
	return rv
}

/*
------------------------------------------------------------------------------------
Due date reminders
------------------------------------------------------------------------------------
*/

// IDOReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week
func IDOReminderOfDueDateInMonth(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			"Individual development objective is due within a month",
			fc.IDOsDueWithinTheMonth(),
		)
	}
	return rv
}

// IDOReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week, month, quarter
func IDOReminderOfDueDateInQaurter(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			"Individual development objective is due within a qaurter",
			fc.IDOsDueWithinTheQuarter(),
		)
	}
	return rv
}

// InitiativeReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week
func InitiativeReminderOfDueDateInMonth(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			"Individual development objective is due within a month",
			fc.InitiativesDueWithinTheMonth(),
		)
	}
	return rv
}

// InitiativeReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week, month, quarter
func InitiativeReminderOfDueDateInQaurter(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			"Initiative is due within a quarter",
			fc.InitiativesDueWithinTheQuarter(),
		)
	}
	return rv
}

// ObjectiveReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week
func ObjectiveReminderOfDueDateInMonth(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			"Individual development objective is due within a month",
			fc.ObjectivesDueWithinTheMonth(),
		)
	}
	return rv
}

// ObjectiveReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week, month, quarter
func ObjectiveReminderOfDueDateInQaurter(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			"Objective is due within a qaurter",
			fc.ObjectivesDueWithinTheQuarter(),
		)
	}
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
	WeekFour := date.GetDayOfWeekInQuarter(-1, bt.Monday) == date

	feedbackCycle := WeekFour || WeekOne || WeekTwo || WeekThree

	feedbackNotGivenYet := utils.ScheduleEntry(
		"Reminder to provide coaching feedback to colleagues",
		feedbackCycle &&
		fc.TeamValuesExist() &&
		!fc.FeedbackGivenThisQuarter(),
	)

	feedbackAlreadyGiven := utils.ScheduleEntry(
		"Reminder to provide coaching feedback to additional colleagues",
		feedbackCycle &&
		fc.TeamValuesExist() &&
		fc.FeedbackGivenThisQuarter(),
	)

	if len(feedbackNotGivenYet) > 0 {
		rv = feedbackNotGivenYet
	} else if len(feedbackAlreadyGiven) > 0 {
		rv = feedbackAlreadyGiven
	}

	return rv
}

/*
------------------------------------------------------------------------------------
Report reminders
------------------------------------------------------------------------------------
*/

// GenerateIndividualReports is meant to trigger the engagements that
// generate an individual coaching reports to each users.
func GenerateIndividualReports(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	if date == date.GetFirstDayOfQuarter() {
		rv = utils.ScheduleEntry(
			"Generate individual reports if there is some feedback",
			fc.FeedbackForThePreviousQuarterExists(),
		)
	}
	return
}

// NotifyOnAbsentFeedback - notifies user if they haven't received any feedback for the previous quarter.
func NotifyOnAbsentFeedback(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	if date == date.GetFirstDayOfQuarter() {
		rv = utils.ScheduleEntry(
			"Notify if there is no feedback",
			!fc.FeedbackForThePreviousQuarterExists(),
		)
	}
	return
}

// DeliverIndividualReports is meant to trigger the engagements that
// sends out the individual coaching reports to each user.
func DeliverIndividualReports(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	// Last business day of the first full week of the quarter
	if date.GetDayOfSundayWeek1InQuarter(1, bt.Friday) == date {
		rv = utils.ScheduleEntry(
			"Deliver individual reports",
			// we perform all the checks inside that function.
			fc.CollaborationReportExists(),
		)
	}

	return rv
}

// EveryQuarterOnSecondWeekOnFriday - 
func EveryQuarterOnSecondWeekOnFriday(fc adaptive_checks.TypedProfile, date bt.Date) (rv string) {
	// Last business day of the first full week of the quarter
	if date.GetDayOfSundayWeek1InQuarter(1, bt.Friday) == date {
		rv = utils.ScheduleEntry(
			"Last business day of the first full week of the quarter", // "Every quarter on the second week on Friday",
			true,
		)
	}

	return rv
}
