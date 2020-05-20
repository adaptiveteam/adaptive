package schedules

import (
	fcn "github.com/adaptiveteam/adaptive/adaptive-check-function-names"
	bt "github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/checks"
	utils "github.com/adaptiveteam/adaptive/engagement-scheduling"
)

// DebugReminder is meant to always trigger the debug engagements
func DebugReminder(fc checks.CheckResultMap, date bt.Date) (rv string) {
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
func IDOCreateReminder(fc checks.CheckResultMap, date bt.Date) (rv string) {
	WeekOne := date.GetDayOfWeekInQuarter(2, bt.Monday) == date
	WeekTwo := date.GetDayOfWeekInQuarter(3, bt.Monday) == date

	begin := date.GetDayOfWeekInQuarter(4, bt.Monday)
	end := date.GetDayOfWeekInQuarter(4, bt.Friday)
	Daily := date.DateAfter(begin, true) && date.DateBefore(end, true)

	rv = utils.ScheduleEntry(
		fc,
		"Create Individual Development Objectives",
	).AddScheduleBooleanCheck(
		Daily || WeekOne || WeekTwo,
		true,
	).AddScheduleFunctionCheck(
		fcn.TeamValuesExist,
		true,
	).AddScheduleFunctionCheck(
		fcn.CanBeNudgedForIDO,
		true,
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
func IDOUpdateReminder(fc checks.CheckResultMap, date bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	if date.DayOfWeek() == bt.Friday {
		rv = utils.ScheduleEntry(
			fc,
			"Update Individual Development Objectives",
		).AddScheduleFunctionCheck(
			fcn.StaleIDOsExistForMe,
			true,
		).Message
	}
	return rv
}

// ObjectiveUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale objectives
func ObjectiveUpdateReminder(fc checks.CheckResultMap, date bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	if date.DayOfWeek() == bt.Friday {
		rv = utils.ScheduleEntry(
			fc,
			"Update Objectives",
		).AddScheduleFunctionCheck(
			fcn.StaleObjectivesExistForMe,
			true,
		).Message
	}
	return rv
}

// InitiativeUpdateReminder is meant to trigger the engagements that
// reminds the user to update stale initiatives
func InitiativeUpdateReminder(fc checks.CheckResultMap, date bt.Date) (rv string) {
	// Starting the last business day of the second week of the first quarter with  a preference for earlier
	if date.DayOfWeek() == bt.Friday {
		rv = utils.ScheduleEntry(
			fc,
			"Update Initiatives",
		).AddScheduleFunctionCheck(
			fcn.StaleInitiativesExistForMe,
			true,
		).Message
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
func IDOCloseoutReminder(fc checks.CheckResultMap, date bt.Date) (rv string) {
	if date.IsWeekDay() {
		rv = utils.ScheduleEntry(
			fc,
			"Closeout individual development objective",
		).AddScheduleFunctionCheck(
			fcn.IDOsDueWithinTheWeek,
			true,
		).Message
	}
	return rv
}

// InitiativeCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Initiative due in the coming week and to close it out
func InitiativeCloseoutReminder(fc checks.CheckResultMap, date bt.Date) (rv string) {
	if date.IsWeekDay() {
		rv = utils.ScheduleEntry(
			fc,
			"Closeout Initiative",
		).AddScheduleFunctionCheck(
			fcn.InitiativesDueWithinTheWeek,
			true,
		).Message
	}
	return rv
}

// ObjectiveCloseoutReminder is meant to trigger engagements that reminds users
// that they have an Objective due in the coming week and to close it out
func ObjectiveCloseoutReminder(fc checks.CheckResultMap, date bt.Date) (rv string) {
	if date.IsWeekDay() {
		rv = utils.ScheduleEntry(
			fc,
			"Closeout Objective",
		).AddScheduleFunctionCheck(
			fcn.ObjectivesDueWithinTheWeek,
			true,
		).Message
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
func IDOReminderOfDueDateInMonth(fc checks.CheckResultMap, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			fc,
			"Individual development objective is due within a month",
		).AddScheduleFunctionCheck(
			fcn.IDOsDueWithinTheMonth,
			true,
		).Message
	}
	return rv
}

// IDOReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an IDO due in the coming week, month, quarter
func IDOReminderOfDueDateInQaurter(fc checks.CheckResultMap, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			fc,
			"Individual development objective is due within a qaurter",
		).AddScheduleFunctionCheck(
			fcn.IDOsDueWithinTheQuarter,
			true,
		).Message
	}
	return rv
}

// InitiativeReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week
func InitiativeReminderOfDueDateInMonth(fc checks.CheckResultMap, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			fc,
			"Individual development objective is due within a month",
		).AddScheduleFunctionCheck(
			fcn.InitiativesDueWithinTheMonth,
			true,
		).Message
	}
	return rv
}

// InitiativeReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Initiative due in the coming week, month, quarter
func InitiativeReminderOfDueDateInQaurter(fc checks.CheckResultMap, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			fc,
			"Initiative is due within a quarter",
		).AddScheduleFunctionCheck(
			fcn.InitiativesDueWithinTheQuarter,
			true,
		).Message
	}
	return rv
}

// ObjectiveReminderOfDueDateInMonth is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week
func ObjectiveReminderOfDueDateInMonth(fc checks.CheckResultMap, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			fc,
			"Individual development objective is due within a month",
		).AddScheduleFunctionCheck(
			fcn.ObjectivesDueWithinTheMonth,
			true,
		).Message
	}
	return rv
}

// ObjectiveReminderOfDueDateInQaurter is meant to trigger  engagements that reminds users
// that they have an Objective due in the coming week, month, quarter
func ObjectiveReminderOfDueDateInQaurter(fc checks.CheckResultMap, date bt.Date) (rv string) {
	if date.DayOfWeek() == bt.Monday {
		rv = utils.ScheduleEntry(
			fc,
			"Objective is due within a qaurter",
		).AddScheduleFunctionCheck(
			fcn.ObjectivesDueWithinTheQuarter,
			true,
		).Message
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
func ReminderToProvideCoachingFeedback(fc checks.CheckResultMap, date bt.Date) (rv string) {
	WeekOne := date.GetDayOfWeekInQuarter(-4, bt.Monday) == date
	WeekTwo := date.GetDayOfWeekInQuarter(-3, bt.Monday) == date
	WeekThree := date.GetDayOfWeekInQuarter(-2, bt.Monday) == date
	WeekFour := date.GetDayOfWeekInQuarter(-1, bt.Monday) == date

	feedbackCycle := WeekFour || WeekOne || WeekTwo || WeekThree

	feedbackNotGivenYet := utils.ScheduleEntry(
		fc,
		"Reminder to provide coaching feedback to colleagues",
	).AddScheduleBooleanCheck(
		feedbackCycle,
		true,
	).AddScheduleFunctionCheck(
		fcn.TeamValuesExist,
		true,
	).AddScheduleFunctionCheck(
		fcn.FeedbackGivenThisQuarter,
		false,
	).Message

	feedbackAlreadyGiven := utils.ScheduleEntry(
		fc,
		"Reminder to provide coaching feedback to additional colleagues",
	).AddScheduleBooleanCheck(
		feedbackCycle,
		true,
	).AddScheduleFunctionCheck(
		fcn.TeamValuesExist,
		true,
	).AddScheduleFunctionCheck(
		fcn.FeedbackGivenThisQuarter,
		true,
	).Message

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
func GenerateIndividualReports(fc checks.CheckResultMap, date bt.Date) (rv string) {
	if date == date.GetFirstDayOfQuarter() {
		rv = utils.ScheduleEntry(
			fc,
			"Generate individual reports if there is some feedback",
		).AddScheduleFunctionCheck(
			fcn.FeedbackForThePreviousQuarterExists,
			true,
		).Message
	}
	return
}

// NotifyOnAbsentFeedback - notifies user if they haven't received any feedback for the previous quarter.
func NotifyOnAbsentFeedback(fc checks.CheckResultMap, date bt.Date) (rv string) {
	if date == date.GetFirstDayOfQuarter() {
		rv = utils.ScheduleEntry(
			fc,
			"Notify if there is no feedback",
		).AddScheduleFunctionCheck(
			fcn.FeedbackForThePreviousQuarterExists,
			false,
		).Message
	}
	return
}

// DeliverIndividualReports is meant to trigger the engagements that
// sends out the individual coaching reports to each user.
func DeliverIndividualReports(fc checks.CheckResultMap, date bt.Date) (rv string) {
	// Last business day of the first full week of the quarter
	if date.GetDayOfSundayWeek1InQuarter(1, bt.Friday) == date {
		rv = utils.ScheduleEntry(
			fc,
			"Deliver individual reports",
			// we perform all the checks inside that function.
		).AddScheduleFunctionCheck(
			fcn.CollaborationReportExists,
			true,
		).Message
	}

	return rv
}

// EveryQuarterOnSecondWeekOnFriday - 
func EveryQuarterOnSecondWeekOnFriday(fc checks.CheckResultMap, date bt.Date) (rv string) {
	// Last business day of the first full week of the quarter
	if date.GetDayOfSundayWeek1InQuarter(1, bt.Friday) == date {
		rv = utils.ScheduleEntry(
			fc,
			"Last business day of the first full week of the quarter", // "Every quarter on the second week on Friday",
		).Message
	}

	return rv
}
