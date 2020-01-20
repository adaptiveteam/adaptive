package lambda

import (
	"fmt"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

const (
	EmbursePlatformID = "ANT7U58AG"
)

func emulateDatesForEmburse(date time.Time, userID string, platformID models.PlatformID, config Config) {
	// Starting with January 27th
	startDate := time.Date(2020, 1, 27, 0, 0, 0, 0, time.UTC)

	startEmulateDate := time.Date(2019, 12, 23, 0, 0, 0, 0, time.UTC)
	endEmulateDate := time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC)

	if !date.Before(startDate) {
		days := int(date.Sub(startDate).Hours() / 24)
		dateToEmulate := startEmulateDate.AddDate(0, 0, days)
		if !dateToEmulate.After(endEmulateDate) {
			emulatedDateStr := core.ISODateLayout.Format(dateToEmulate)
			dateStr := core.ISODateLayout.Format(date)
			fmt.Println(fmt.Sprintf("user %s: Emulating date %s --> %s", userID, dateStr, emulatedDateStr))
			engage := models.UserEngage{
				UserId:     userID,
				PlatformID: platformID,
				Date:       emulatedDateStr,
			}
			logger.Infof("Emulating %s date for %s user on %s", emulatedDateStr, userID, dateStr)
			invokeScriptingLambda(engage, models.PlatformID(platformID), config)
		}
	}
}
