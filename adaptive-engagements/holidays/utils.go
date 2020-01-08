package holidays

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"time"
)

func ConvertHolidaysArray(adHocHolidays []models.AdHocHoliday) business_time.Holidays {
	holidaysInterface := business_time.NewHolidayList()
	for _, h := range adHocHolidays {
		holidaysInterface.AddHoliday(h.Name, parseDate(h.Date), *time.UTC)
	}
	return holidaysInterface
}

func parseDate(d string) business_time.Date {
	date, err := business_time.DateFromYMDString(d)
	core.ErrorHandler(err, "holidays", "Couldn't parse date")
	return date
}
