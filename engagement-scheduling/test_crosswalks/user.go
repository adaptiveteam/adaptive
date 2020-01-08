package test_crosswalks

import (
	models "github.com/adaptiveteam/adaptive/engagement-scheduling-models"
	engagements "github.com/adaptiveteam/adaptive/engagement-scheduling/test_engagements"
	schedules "github.com/adaptiveteam/adaptive/engagement-scheduling/test_schedules"
)

func UserCrosswalk() (rv []models.CrossWalk) {
	return userCrosswalk
}

var userCrosswalk = []models.CrossWalk{
	/*
		------------------------------------------------------------------------------------
		IDO Creation reminders
		------------------------------------------------------------------------------------
	*/
	models.NewCrossWalk(schedules.IDOCreateReminder, engagements.IDOCreateReminder),
	/*
	   ------------------------------------------------------------------------------------
	   Update Reminders
	   ------------------------------------------------------------------------------------
	*/
	models.NewCrossWalk(schedules.IDOUpdateReminder, engagements.IDOUpdateReminder),
	models.NewCrossWalk(schedules.ObjectiveUpdateReminder, engagements.ObjectiveUpdateReminder),
	models.NewCrossWalk(schedules.InitiativeUpdateReminder, engagements.InitiativeUpdateReminder),
	/*
		------------------------------------------------------------------------------------
		Closeout reminders
		------------------------------------------------------------------------------------
	*/
	models.NewCrossWalk(schedules.IDOCloseoutReminder, engagements.IDOCloseoutReminder),
	models.NewCrossWalk(schedules.InitiativeCloseoutReminder, engagements.InitiativeCloseoutReminder),
	models.NewCrossWalk(schedules.ObjectiveCloseoutReminder, engagements.ObjectiveCloseoutReminder),
	/*
		------------------------------------------------------------------------------------
		Due date reminders
		------------------------------------------------------------------------------------
	*/
	models.NewCrossWalk(schedules.IDOReminderOfDueDateInMonth, engagements.IDOReminderOfDueDateInMonth),
	models.NewCrossWalk(schedules.IDOReminderOfDueDateInQaurter, engagements.IDOReminderOfDueDateInQaurter),
	models.NewCrossWalk(schedules.InitiativeReminderOfDueDateInMonth, engagements.InitiativeReminderOfDueDateInMonth),
	models.NewCrossWalk(schedules.InitiativeReminderOfDueDateInQaurter, engagements.InitiativeReminderOfDueDateInQaurter),
	models.NewCrossWalk(schedules.ObjectiveReminderOfDueDateInMonth, engagements.ObjectiveReminderOfDueDateInMonth),
	models.NewCrossWalk(schedules.ObjectiveReminderOfDueDateInQaurter, engagements.ObjectiveReminderOfDueDateInQaurter),
	/*
		------------------------------------------------------------------------------------
		Coaching feedback reminders
		------------------------------------------------------------------------------------
	*/
	models.NewCrossWalk(schedules.ReminderToProvideCoachingFeedback, engagements.ReminderToProvideCoachingFeedback),
	/*
		------------------------------------------------------------------------------------
		Report reminders
		------------------------------------------------------------------------------------
	*/
	models.NewCrossWalk(schedules.ProduceIndividualReports, engagements.ProduceIndividualReports),
}
