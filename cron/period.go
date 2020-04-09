package cron

import (
	"github.com/pkg/errors"
	"time"
)

// Period - one of Year, Quarter,...Day
type Period string

const (
	Epoch Period = "Epoch"
	Century Period = "Century"
	Year Period = "Year"
	Quarter Period = "Quarter"
	Month Period = "Month"
	FullWeek Period = "FullWeek"
	AnyWeek Period = "AnyWeek" // this week may have at least 1 day
	Day Period = "Day"
	WeekDay Period = "WeekDay" // a special period. Used only in combination with Sunday, ... Saturday, because Sunday == 0, while all numbering starts with 1.
	Hour Period = "Hour"
	QuarterHour Period = "QuarterHour" // 15 minutes
	Minute Period = "Minute"
	Second Period = "Second"
)

// StartEnd returns time moments of the beginning of the period and 
// the beginning of the next period.
func (period Period)StartEnd(t time.Time) (start, end time.Time) {
	switch period {
	case Epoch:
		start = time.Date(t.Year() / 10000 * 10000 + 1, 1, 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(10000, 0, 0)
	case Century:
		start = time.Date(t.Year() / 100 * 100 + 1, 1, 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(100, 0, 0)
	case Year:
		start = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(1, 0, 0)
	case Quarter:
		start = time.Date(t.Year(), (t.Month() - 1) / 3 * 3 + 1, 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 3, 0)
	case Month:
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 1, 0)
	case AnyWeek, FullWeek: // NB! This doesn't work for short weeks
		weekDay := t.Weekday()
		start = time.Date(t.Year(), t.Month(), t.Day() - int(weekDay), 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 0, 7)
	case Day:
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 0, 1)
	case Hour:
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
		end = start.Add(time.Hour)
	case QuarterHour:
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute() / 15 * 15, 0, 0, time.UTC)
		end = start.Add(15 * time.Minute)
	case Minute:
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
		end = start.Add(time.Minute)
	case Second:
		start = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
		end = start.Add(time.Second)
	default:
		panic(errors.New("GetPeriodStartEnd has received an unknown period " + string(period)))
	}
	return
}

// Interval returns time moments of the beginning of the period and 
// the beginning of the next period.
func (period Period)Interval(t time.Time) Interval {
	start, end := period.StartEnd(t)
	return Interval{
		StartInclusive: start,
		EndExclusive: end,
	}
}
