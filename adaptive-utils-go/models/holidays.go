package models

import (
	"github.com/adaptiveteam/adaptive/daos/adHocHoliday"
)

// AnnualHoliday represents a single holiday that repeats on the same date annualy. 
// Holidays are not universal, they should only apply to certain communities.
type AnnualHoliday struct {
	ID               string `json:"id"`
	PlatformID       string `json:"platform_id"`
	Month            int    `json:"month"`
	DayOfMonth       int    `json:"day_of_month"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ScopeCommunities string `json:"scope_communities"`
}

// AdHocHoliday is a holiday on exact date.
type AdHocHoliday = adHocHoliday.AdHocHoliday
//  struct {
// 	ID               string `json:"id"`
// 	PlatformID       string `json:"platform_id"`
// 	Date             string `json:"date"` // time.Time
// 	Name             string `json:"name"`
// 	Description      string `json:"description"`
// 	ScopeCommunities string `json:"scope_communities"`
// }

const (
	// AdHocHolidayDateFormat is the format that is used in Date field
	// see "time" package for details
	AdHocHolidayDateFormat string = "2006-01-02"
)

// HolidaysTableSchema schema of ad-hoc holidays table
type HolidaysTableSchema struct {
	Name string
	PlatformDateIndex string
}

// HolidaysTableSchemaForClientID creates HolidaysTableSchema for a given clientID
func HolidaysTableSchemaForClientID(clientID string) HolidaysTableSchema {
	return HolidaysTableSchema{
		Name: clientID + "_ad_hoc_holidays",
		PlatformDateIndex: "HolidaysPlatformDateIndex",
	}
}
