package lambda

import (
	"fmt"
	aesc "github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/schedules"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	es "github.com/adaptiveteam/adaptive/engagement-scheduling"
	esmodels "github.com/adaptiveteam/adaptive/engagement-scheduling-models"
	"strings"
	"time"
	"github.com/adaptiveteam/adaptive/daos/adHocHoliday"
)

func ScheduleOfEngagements(target, date, holidaysTable string,
	crosswalks func() []esmodels.CrossWalk, days int, teamID models.TeamID) ([]esmodels.ScheduledEngagement, error) {
	var y, m, d = time.Now().Date()
	if date != core.EmptyString {
		// When passed date is not empty, parse that to date
		parsedTime, err := time.Parse(string(core.ISODateLayout), date)
		fmt.Println(parsedTime)
		if err == nil {
			y, m, d = parsedTime.Date()
		} else {
			return nil, err
		}
	}
	location, _ := time.LoadLocation("UTC")
	holidaysList := schedules.LoadHolidays(time.Date(y, m, d, 0, 0, 0, 0, location), teamID,
		holidaysTable, string(adHocHoliday.PlatformIDDateIndex))

	day := business_time.NewDate(y, int(m), d)
	return es.GenerateScheduleOfEngagements(
		aesc.ProductionProfile,
		day,
		target,
		crosswalks,
		holidaysList,
		location,
		days,
	), nil
}

// allSchedules returns all the applicable schedules for a user for the next n days
func allSchedules(date business_time.Date, userID string, days int, teamID models.TeamID) []esmodels.ScheduledEngagement {
	userSchedules, err := ScheduleOfEngagements(userID, core.ISODateLayout.Format(date.DateToTimeMidnight()),
		adHocHolidaysTable, func() []esmodels.CrossWalk {
			return allCrosswalks
		}, days, teamID)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not get schedules for %s", userID))
	return userSchedules
}

func schedulesSummary(userID string, teamID models.TeamID, quarter, year int) ui.RichText {
	quarterStart := business_time.NewDateFromQuarter(quarter, year)
	quarterEnd := quarterStart.GetLastDayOfQuarter()

	userSchedules := allSchedules(quarterStart, userID, quarterEnd.DaysBetween(quarterStart), teamID)
	// fmt.Println(userSchedules)
	// fmt.Println(len(userSchedules))
	// byt, _ := json.Marshal(es.PrettyPrintSchedule(userSchedules))
	// fmt.Println("### userSchedules: " + string(byt))
	return ui.RichText(strings.Join(es.PrettyPrintSchedule(userSchedules), "\n"))
}
