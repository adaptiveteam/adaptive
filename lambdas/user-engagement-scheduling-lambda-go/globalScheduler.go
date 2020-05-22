package userEngagementScheduling

import (
	"github.com/pkg/errors"
	"github.com/adaptiveteam/adaptive/cron"
	"time"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

var (
	everyQuarterLastMonthEveryFridayAt12 = cron.Q().
		EveryN(cron.Month, 3).
		Every(cron.FullWeek).
		OnWeekDay(time.Friday).
		InRange0(cron.Hour, 12, 12). 
		InRange0(cron.QuarterHour, 0, 0)
)

func runGlobalScheduleForTeam(config Config, teamID models.TeamID, t time.Time) (err error) {
	if everyQuarterLastMonthEveryFridayAt12.IsOnSchedule(t) {
		err = sendFeedbackStatsReport(config, teamID)
	}
	return errors.Wrapf(err, "runGlobalScheduleForTeam(config, %v, %v)", teamID, t)
}
