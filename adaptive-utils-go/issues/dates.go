package issues

import (
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// NormalizeTimestampInPlace makes sure that the value is correct timestamp
func NormalizeTimestampInPlace(dateOrTimestampOrEmpty *string) {
	v := core.NormalizeTimestamp(*dateOrTimestampOrEmpty)
	dateOrTimestampOrEmpty = &v
}

// NormalizeDateInPlace makes sure that the value is correct date
func NormalizeDateInPlace(dateOrTimestampOrEmpty *string) {
	if *dateOrTimestampOrEmpty == "indefinite" {
		// do not change
	} else {
		v := core.NormalizeDate(*dateOrTimestampOrEmpty)
		dateOrTimestampOrEmpty = &v
	}
}

// NormalizeIssueDateTimes makes sure that all date/timestamp fields are in correct format
func (issue *Issue)NormalizeIssueDateTimes() {
	NormalizeDateInPlace(&issue.UserObjective.CompletedDate)
	NormalizeTimestampInPlace(&issue.UserObjective.CreatedAt)
	NormalizeDateInPlace(&issue.UserObjective.CreatedDate)
	NormalizeDateInPlace(&issue.UserObjective.ExpectedEndDate)
	NormalizeTimestampInPlace(&issue.UserObjective.ModifiedAt)

	NormalizeTimestampInPlace(&issue.StrategyObjective.CreatedAt)
	NormalizeDateInPlace(&issue.StrategyObjective.ExpectedEndDate)

	NormalizeTimestampInPlace(&issue.StrategyInitiative.CreatedAt)
	NormalizeDateInPlace(&issue.StrategyInitiative.ExpectedEndDate)
	NormalizeTimestampInPlace(&issue.StrategyInitiative.ModifiedAt)
}

