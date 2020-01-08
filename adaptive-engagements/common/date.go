package common

import (
	"log"
	"fmt"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"strings"
	"time"
)

const (
	StrategyIndefiniteDateKey   = "Indefinite"
	StrategyIndefiniteDateValue = "indefinite"
	DateFormat                  = core.ISODateLayout
)

// DurationDays calculates the days difference between two dates.
// start and end should be valid dates in the given layout.
// incorrect values are replaced with `time.Now`
// This function does not panic.
func DurationDays(start, end string, layout core.AdaptiveDateLayout, namespace string) int {
	startDate, err2 := time.Parse(string(layout), start)
	if err2 != nil {
		log.Printf("DurationDays: Could not parse start date %s as %s. Using time.Now()", start, layout)
		startDate = time.Now()
	}
	endDate, err2 := time.Parse(string(layout), end)
	if err2 != nil {
		log.Printf("DurationDays: Could not parse end date %s as %s. Using time.Now()", end, layout)
		endDate = time.Now()
	}
	return int(endDate.Sub(startDate).Hours() / 24)
}

// A struct holding start and end dates
type ObjectiveDate struct {
	CreatedDate     string `json:"created_date"`
	ExpectedEndDate string `json:"expected_end_date"`
}

// Render returns formatted progress information using current date and
// two dates - start/end.
// start and end should be valid dates in the given layout.
// incorrect values are shown as is, and replaced with `time.Now` in calculations.
// This function does not panic.
func (od ObjectiveDate) Render(parseLayout, renderLayout core.AdaptiveDateLayout, namespace string) string {
	var op string
	renderedStartDate, err2 := parseLayout.ChangeLayout(od.CreatedDate, renderLayout)
	if err2 != nil {
		log.Printf("ObjectiveDate) Render: Could not convert date %s from %s to %s", od.CreatedDate, parseLayout, renderLayout)
		renderedStartDate = od.CreatedDate
	}
	op += fmt.Sprintf("Created: %s \n", renderedStartDate)

	today := time.Now().Format(string(parseLayout))
	if od.ExpectedEndDate != StrategyIndefiniteDateValue {
		renderedEndDate, err2 := parseLayout.ChangeLayout(od.ExpectedEndDate, renderLayout)
		if err2 != nil {
			log.Printf("ObjectiveDate) Render: Could not convert date %s from %s to %s", od.ExpectedEndDate, parseLayout, renderLayout)
			renderedEndDate = od.ExpectedEndDate
		}

		op += fmt.Sprintf("Expected End: %s \n", renderedEndDate)
		op += fmt.Sprintf("Remaining Days: %d \n", DurationDays(today, od.ExpectedEndDate, parseLayout, namespace))
	}
	op += fmt.Sprintf("Elapsed Days: %d \n", DurationDays(od.CreatedDate, today, parseLayout, namespace))
	return op
}

func FormatDateWithIndefiniteOption(date string, ipLayout, opLayout core.AdaptiveDateLayout, namespace string) (op string, err error) {
	if date == StrategyIndefiniteDateValue {
		return strings.Title(StrategyIndefiniteDateValue), nil
	} else {
		pDate, err := time.Parse(string(ipLayout), date)
		if err == nil {
			op = pDate.Format(string(opLayout))
		}
	}
	return op, err
}
