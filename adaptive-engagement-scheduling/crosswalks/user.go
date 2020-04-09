package crosswalks

import (
	"time"
	"github.com/adaptiveteam/adaptive/cron"
	"github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/engagements"
	"github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/schedules"
	models "github.com/adaptiveteam/adaptive/engagement-scheduling-models"
)

func UserCrosswalk() (rv []models.CrossWalk) {
	return userCrosswalk
}

var userCrosswalk = []models.CrossWalk{
	/* DEBUG reminder*/
	// models.NewCrossWalk(schedules.DebugReminder, engagements.IDOUpdateReminder),               // ac.StaleIDOsExist
	/*
		------------------------------------------------------------------------------------
		IDO Creation reminders
		------------------------------------------------------------------------------------
	*/
	models.NewCrossWalk(schedules.IDOCreateReminder, engagements.IDOCreateReminder), // ac.IDOsExistForMe
	/*
	   ------------------------------------------------------------------------------------
	   Update Reminders
	   ------------------------------------------------------------------------------------
	*/
	models.NewCrossWalk(schedules.IDOUpdateReminder, engagements.IDOUpdateReminder),               // ac.StaleIDOsExist
	models.NewCrossWalk(schedules.ObjectiveUpdateReminder, engagements.ObjectiveUpdateReminder),   // ac.StaleObjectivesExist
	models.NewCrossWalk(schedules.InitiativeUpdateReminder, engagements.InitiativeUpdateReminder), // ac.StaleInitiativesExist
	/*
		------------------------------------------------------------------------------------
		Closeout reminders
		------------------------------------------------------------------------------------
	*/
	// TODO: Look into these implementations again
	models.NewCrossWalk(schedules.IDOCloseoutReminder, engagements.IDOCloseoutReminder),               // ac.IDOsDueInAWeek,
	models.NewCrossWalk(schedules.InitiativeCloseoutReminder, engagements.InitiativeCloseoutReminder), // ac.InitiativesDueInAWeek
	models.NewCrossWalk(schedules.ObjectiveCloseoutReminder, engagements.ObjectiveCloseoutReminder),   // ac.CapabilityObjectivesDueInAWeek
	/*
		------------------------------------------------------------------------------------
		Due date reminders
		------------------------------------------------------------------------------------
	*/
	models.NewCrossWalk(schedules.IDOReminderOfDueDateInMonth, engagements.IDOReminderOfDueDateInMonth),                   // ac.IDOsDueInAMonth
	models.NewCrossWalk(schedules.IDOReminderOfDueDateInQaurter, engagements.IDOReminderOfDueDateInQaurter),               // ac.IDOsDueInAQuarter
	models.NewCrossWalk(schedules.InitiativeReminderOfDueDateInMonth, engagements.InitiativeReminderOfDueDateInMonth),     // ac.InitiativesDueWithinTheMonth
	models.NewCrossWalk(schedules.InitiativeReminderOfDueDateInQaurter, engagements.InitiativeReminderOfDueDateInQaurter), // ac.InitiativesDueInAQuarter
	models.NewCrossWalk(schedules.ObjectiveReminderOfDueDateInMonth, engagements.ObjectiveReminderOfDueDateInMonth),       // ac.CapabilityObjectivesDueInAMonth
	models.NewCrossWalk(schedules.ObjectiveReminderOfDueDateInQaurter, engagements.ObjectiveReminderOfDueDateInQaurter),   // ac.ac.CapabilityObjectivesDueInAQuarter
	/*
		------------------------------------------------------------------------------------
		Coaching feedback reminders
		------------------------------------------------------------------------------------
	*/
	// TODO: Look into this implementation
	models.NewCrossWalk(schedules.ReminderToProvideCoachingFeedback, engagements.ReminderToProvideCoachingFeedback),
	/*
		------------------------------------------------------------------------------------
		Report reminders
		------------------------------------------------------------------------------------
	*/
	models.NewCrossWalk(schedules.GenerateIndividualReports, engagements.GenerateIndividualReports),
	models.CrontabLine(
		cron.Q().
			InRange(cron.FullWeek, 1, 1).
			OnWeekDay(time.Thursday),
		"Produce individual reports if feedback is present", 
		engagements.GenerateIndividualReports),
	models.CrontabLine(
		cron.Q().
			InRange(cron.FullWeek, 1, 1).
			OnWeekDay(time.Friday),
		"Deliver individual reports or notify if feedback is absent", 
		engagements.DeliverIndividualReportsOrNotifyOnAbsentFeedback),
	// models.NewCrossWalk(schedules.DeliverIndividualReports, engagements.DeliverIndividualReports),
	
}
