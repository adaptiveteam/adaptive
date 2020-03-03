package holidays

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"time"
)

const (
	NameLabel              ui.PlainText = "Holiday name"
	NamePlaceholder        ui.PlainText = ebm.EmptyPlaceholder
	DescriptionLabel       ui.PlainText = "Description"
	DescriptionPlaceholder ui.PlainText = ebm.EmptyPlaceholder
	DateLabel              ui.PlainText = "Date of the holiday (YYYY-MM-DD)"
	DatePlaceholder        ui.PlainText = ebm.EmptyPlaceholder
)

// CreateAdHocHolidayEng This is an engagement prompting a user to create 
// a new ad-hoc holiday
func CreateAdHocHolidayEng(mc models.MessageCallback,
	fallback, learnTrailPath, namespace string, check models.UserEngagementCheckWithValue,
	teamID models.TeamID) models.UserEngagement {
	title := "Create an ad-hoc holiday"
	actions := models.AppendOptionalAction(
		[]ebm.AttachmentAction{
			*models.DialogAttachAction(mc, models.Now, "SimpleCreate")},
		models.LearnMoreAction(learnTrailPath))
	fields := []ebm.AttachmentField{}
	urgent := true
	return utils.MakeUserEngagement(mc, title, core.EmptyString, fallback,
		mc.Source, actions, fields, urgent, namespace, time.Now().Unix(), check, teamID)
}

// EditAdHocHolidayForm creates a form for modifying/constructing an ad-hoc holiday
func EditAdHocHolidayForm(obj *models.AdHocHoliday) []ebm.AttachmentActionTextElement {
	if obj == nil {
		obj = &models.AdHocHoliday{}
	}
	return []ebm.AttachmentActionTextElement{
		ebm.NewTextBox("Name", NameLabel, NamePlaceholder, ui.PlainText(obj.Name)),
		ebm.NewTextArea("Description", DescriptionLabel, DescriptionPlaceholder, ui.PlainText(obj.Description)),
		ebm.NewDateInput("Date", DateLabel, DatePlaceholder, ui.PlainText(obj.Date)),
	}
}
