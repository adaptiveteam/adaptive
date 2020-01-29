package lambda

import (
	"fmt"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

const (
	EmbursePlatformID = "ANT7U58AG"
	GeigsenPlatformID = "AGEGG1U7J"
	IvanPlatformID = "AJA9TJ88Y"
	StagingPlatformID = "ALV1A59GR"
)
// DateShiftConfig configures dates emulation
type DateShiftConfig struct {
	StartDate time.Time
	StartEmulateDate time.Time
	EndEmulateDate time.Time
}

var EmburseDateShiftConfig = DateShiftConfig{
	// Starting with January 27th
	StartDate: time.Date(2020, 1, 27, 0, 0, 0, 0, time.UTC),
	StartEmulateDate: time.Date(2019, 12, 23, 0, 0, 0, 0, time.UTC),
	EndEmulateDate: time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC),
}

var TestDateShiftConfig = DateShiftConfig{
	// Starting with January 21th
	StartDate: time.Date(2020, 1, 27, 0, 0, 0, 0, time.UTC),
	StartEmulateDate: time.Date(2019, 12, 23, 0, 0, 0, 0, time.UTC),
	EndEmulateDate: time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC),
}

func emulateDates(dateShiftConfig DateShiftConfig, date time.Time, userID string, platformID models.PlatformID, config Config) {

	if !date.Before(dateShiftConfig.StartDate) {
		days := int(date.Sub(dateShiftConfig.StartDate).Hours() / 24)
		dateToEmulate := dateShiftConfig.StartEmulateDate.AddDate(0, 0, days)
		if !dateToEmulate.After(dateShiftConfig.EndEmulateDate) {
			emulatedDateStr := core.ISODateLayout.Format(dateToEmulate)
			dateStr := core.ISODateLayout.Format(date)
			fmt.Println(fmt.Sprintf("user %s: Emulating date %s --> %s", userID, dateStr, emulatedDateStr))
			engage := models.UserEngage{
				UserId:     userID,
				PlatformID: platformID,
				Date:       "2019-12-23",
			}
			logger.Infof("Emulating %s date for %s user on %s", emulatedDateStr, userID, dateStr)
			invokeSchedulerLambda(engage, config)
		}
	}
}
