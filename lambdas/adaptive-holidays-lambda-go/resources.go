package lambda

import (
	"strings"
	"time"

	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// Subject is the typeof the key for accessing templates
type Subject string

const (
	// HolidaysContext is the key to fetch dialog messages
	HolidaysContext = "engagements/holidays"

	InviteToListOfHolidaysSubject Subject = "InviteToListOfHolidays"
	HolidaysListTitleSubject      Subject = "HolidaysListTitle"
	HolidayItemSubject            Subject = "HolidayItem"
	CreatedHolidayNoticeSubject   Subject = "CreatedHolidayNotice"
	EditButtonSubject             Subject = "EditButton"
	UpdatedHolidayNoticeSubject   Subject = "UpdatedHolidayNotice"
	DeleteButtonSubject           Subject = "DeleteButton"
	DeletedHolidayNoticeSubject   Subject = "DeletedHolidayNotice"
	CancelDeleteSubject           Subject = "CancelDeleteAction"
	CreateButtonSubject           Subject = "CreateButton"
	ModifyListButtonSubject       Subject = "ModifyListButton"
)

const (
	inviteToListOfHolidaysTemplate = "Please find the list of holidays in the thread below :arrow_down:"
	holidaysListTitleTemplate      = "Holidays"
	holidayItemTemplate            = "• {FormatDate(AdHocHoliday.Date, 'MMMM DD')}, *{AdHocHoliday.Name}*\n{AdHocHoliday.Description}\n"
	createdHolidayNoticeTemplate   = "Created holiday"
	editButtonLabel                = "Edit" // "🖉"
	updatedHolidayNoticeTemplate   = "Holiday updated"
	deleteButtonLabel              = "Delete" // "🗑"
	deletedHolidayNoticeTemplate   = "Deleted the holiday."
	cancelDeleteText               = "Holiday deletion has been cancelled"
	createButtonLabel              = "Add another"
	modifyListButtonLabel          = "Modify/Delete one"
)

var (
	templates = map[Subject]string{
		InviteToListOfHolidaysSubject: inviteToListOfHolidaysTemplate,
		HolidaysListTitleSubject:      holidaysListTitleTemplate,
		HolidayItemSubject:            holidayItemTemplate,
		CreatedHolidayNoticeSubject:   createdHolidayNoticeTemplate,
		EditButtonSubject:             editButtonLabel,
		UpdatedHolidayNoticeSubject:   updatedHolidayNoticeTemplate,
		DeleteButtonSubject:           deleteButtonLabel,
		DeletedHolidayNoticeSubject:   deletedHolidayNoticeTemplate,
		CancelDeleteSubject:           cancelDeleteText,
		CreateButtonSubject:           createButtonLabel,
		ModifyListButtonSubject:       modifyListButtonLabel,
	}
)

// RetrieveTemplate returns a dialog message template for
// a given context and subject
func RetrieveTemplate(context string, subject Subject) string {
	if context == HolidaysContext {
		return templates[subject]
	}
	return "No template for context " + context
}

// HolidaysTemplate returns a dialog message template for
// a given subject
func HolidaysTemplate(subject Subject) string {
	return RetrieveTemplate(HolidaysContext, subject)
}

// RenderAdHocHolidayItem returns a formatted string based on the provided data
func RenderAdHocHolidayItem(adHocHoliday models.AdHocHoliday) string {
	t := RetrieveTemplate(HolidaysContext, HolidayItemSubject)
	dateString := adHocHoliday.Date
	d, err := time.Parse("2006-01-02", dateString)
	if err == nil {
		dateString = d.Format("January 02")
	}
	t = strings.Replace(t, "{FormatDate(AdHocHoliday.Date, 'MMMM DD')}", dateString, -1)
	t = strings.Replace(t, "{AdHocHoliday.Name}", adHocHoliday.Name, -1)
	t = strings.Replace(t, "{AdHocHoliday.Description}", adHocHoliday.Description, -1)
	return t
}

// field creates a field view from provided label and value.
func field(label ui.PlainText, value ui.RichText) ebm.AttachmentField {
	return ebm.AttachmentField{
		Title: string(label),
		Value: string(value),
		Short: false,
	}
}

// ShowAdHocHolidayItemAsFields converts ad-hoc holiday.
// Don't forget to enable markdown everywhere!
// MarkDownIn([]ebm.MarkdownField{ebm.MarkdownFieldText, ebm.MarkdownFieldPretext}).
func ShowAdHocHolidayItemAsFields(adHocHoliday models.AdHocHoliday) []ebm.AttachmentField {
	return []ebm.AttachmentField{
		field("Name", ui.RichText(adHocHoliday.Name).Bold()),
		field("Date", ui.RichText(adHocHoliday.Date)),
		field("Description", ui.RichText(adHocHoliday.Description)),
	}
}
