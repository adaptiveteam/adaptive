package cron

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

// Schedule is similar to CRON. It allows to use the following periods:
// year, 
// quarter, 
// month, 
// week,
// day,
// . 
// And numbering of the periods might be started from any of the above longer periods:
// month-in-year, month-in-quarter, 
// week-in-year, week-in-quarter, week-in-month,
// day-in-year, day-in-quarter, day-in-month, day-in-week.
// For each period we may start from the beginning of the longer period (+1), or
// from the end of the longer period (-1).
// We may introduce simple sparse periodic - every month, every second week, every 10-th day in month
// We may also limit range of numbers
// Example:
//   Every second year 
//     last quarter
//     every week-in-quarter in range [-3,-1]
//     every day-in-week in range [Monday, Friday]
// Here is the respective code:
//  schedules.S().
//            EveryN(Year, 2).
//            Last(Quarter).
//            InRange(Week, -3, -1).Every(Week).
//            InRange(Day, Monday, Friday).Every(Day)

// Schedule is the recursive structured cron-like schedule.
type Schedule struct {
	Parent *Schedule
	Period
	M int // phase within [0, N-1]
	N int
	RangeStart int
	RangeEnd int
}

// S constructs an empty schedule 
func S() Schedule {
	return Schedule{
		Parent: nil,
		Period: Epoch,
		N: 1,
	}
}

// Q constructs a schedule that will trigger every quarter
func Q() Schedule {
	return S().Every(Quarter)
}
// BooleanSchedule checks if the given time moment belongs to the schedule
type BooleanSchedule func (time.Time) bool
// EveryMN makes a schedule that will be valid every n-th period on the m-th phase.
func (s Schedule)EveryMN(period Period, m, n int) Schedule {
	if period == WeekDay {
		period = Day
		m = m + 1
	}
	return Schedule{
		Parent: &s,
		Period: period,
		M: m,
		N: n,
	}
}

// EveryN makes a schedule that will be valid every n-th period
func (s Schedule)EveryN(period Period, n int) Schedule {
	return s.EveryMN(period, 0, n)
}

// Every makes a schedule that is valid every period
func (s Schedule)Every(period Period) Schedule {
	return s.EveryMN(period, 0, 1)
}

// InRange makes a schedule that is valid for periods within the given range.
// start and/or end might be negative. In this case they are calculated from the end of the previous period.
// if period is WeekDay, then start and end are incremented.
func (s Schedule)InRange(period Period, start, end int) (res Schedule) {
	res = Schedule{
		Parent: &s,
		Period: period,
	}
	switch period {
	case WeekDay:
		res.RangeStart = start + 1
		res.RangeEnd = end + 1
		res.Period = Day
	default:
		res.RangeStart = start
		res.RangeEnd = end
	}
	return
}

// InRange0 makes a schedule that is valid for periods within the given range.
// start0 and/or end0 might be negative. In this case they are calculated from the end of the previous period.
func (s Schedule)InRange0(period Period, start0, end0 int) (res Schedule) {
	res = Schedule{
		Parent: &s,
		Period: period,
	}
	switch period {
	default:
		if start0 >= 0 {
			res.RangeStart = start0 + 1
		} else {
			res.RangeStart = start0
		}
		if end0 >= 0 {
			res.RangeEnd = end0 + 1
		} else {
			res.RangeEnd = end0
		}
	}
	return
}

// InDayRange creates a range of week days
func (s Schedule)InDayRange(from, to time.Weekday) (res Schedule) {
	return s.InRange(WeekDay, int(from), int(to))
}
// OnWeekDay is on particular week day 
func (s Schedule)OnWeekDay(day time.Weekday) (res Schedule) {
	return s.InDayRange(day, day)
}

// toSliceOfParentsOrdered unwraps parents and returns the list of
// schedules starting from the very first schedule
func (s Schedule) toSliceOfParentsOrdered() (slice []Schedule) {
	if s.Parent == nil {
		slice = []Schedule{s}
	} else {
		slice = append(s.Parent.toSliceOfParentsOrdered(), s)
	}
	return
}

// DivMod returns quotinent and modulus making sure that modulus is not negative
func DivMod(x, y int) (d int, m int) {
	d, m = x / y, x % y
	if (m < 0) {
		m = m + y
		d = d - 1
	}
	return
}

// GetPosition returns the position of the period of time moment 
// within the parent period.
// position for the first complete period is equal to 1.
// max will contain the maximum index of the period within the given period.
// Week is a special period. It's boundary is not aligned with the boundaries of
// outer periods. So we count weeks in two ways - FullWeek and AnyWeek.
// first FullWeek is the first week of the outer period that contains Sunday (first day).
// last FullWeek is the last week of the outer period that contains Saturday (last day).
// first AnyWeek is the first week of the outer period that contains Saturday (at least one day).
// last AnyWeek is the last week of the outer period that contains Sunday (at least day).
func GetPosition(t time.Time, period, parent Period) (position, max int) {
	invalidInput := func(){ panic(errors.New(fmt.Sprintf("GetPosition(%v, %s, %s) is invalid", t, period, parent)))}

	if parent == period {
		position = 1
		max = 1
	} else {
		switch parent {
		case Epoch:
			position = t.Year()
			max = 10000
		case Century:
			position = (t.Year() - 1) % 100 + 1
			max = 100
		case Year:
			dec31 := time.Date(t.Year(), 12, 31, 0, 0, 0, 0, time.UTC)
			jan1 := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		
			firstYearWeekLength := int(7 - jan1.Weekday()) % 7
			yearDay := t.YearDay()
		
			// lastYearWeekLength := int(dec31.Weekday() + 1) % 7
		
			switch period {
			case Quarter:
				position = (int(t.Month()) - 1) % 3 + 1
				max = 4
			case Month:
				position = int(t.Month())
				max = 12
			case FullWeek:
				lastYearWeekLength := int(dec31.Weekday() + 1) % 7
				position, _ =  DivMod(yearDay   - firstYearWeekLength + 7, 7)
				max = (dec31.YearDay() - firstYearWeekLength - lastYearWeekLength - 1) / 7 + 1
			case AnyWeek:
				countOfShortWeeksInTheBeginningOfYear := (firstYearWeekLength + 6) / 7 // 0 or 1
				// countOfShortWeeksInTheEndOfYear       := ( lastYearWeekLength + 6) / 7 // 0 or 1
				position, _ =  DivMod(yearDay   - firstYearWeekLength + 7 + countOfShortWeeksInTheBeginningOfYear*7, 7)
				max =  (dec31.YearDay() - firstYearWeekLength) / 7 + 1 + countOfShortWeeksInTheBeginningOfYear
			case Day:
				position = t.YearDay()
				max =  dec31.YearDay()
			default:
				invalidInput()
			}
		case Quarter:
			quarterStart, quarterEndPlus1 := Quarter.StartEnd(t)
			quarterDay := t.YearDay() - quarterStart.YearDay() + 1
			quarterEndDay := quarterEndPlus1.YearDay() - quarterStart.YearDay()

			firstQuarterWeekLength := int(7 - quarterStart.Weekday()) % 7
		
			switch period {
			case Month:
				position = (int(t.Month()) - 1) % 3 + 1
				max = 3
			case FullWeek:
				lastQuarterWeekLength := int(quarterEndPlus1.Weekday()) % 7
				position = (7 + quarterDay - firstQuarterWeekLength) / 7
				max =   (7 + quarterEndDay - firstQuarterWeekLength - lastQuarterWeekLength - 1) / 7
			case AnyWeek:
				countOfShortWeeksInTheBeginningOfQuarter := (firstQuarterWeekLength + 6) / 7 // 0 or 1
				position =   (7 +    quarterDay - firstQuarterWeekLength) / 7 + countOfShortWeeksInTheBeginningOfQuarter
				max      =   (7 + quarterEndDay - firstQuarterWeekLength) / 7 + countOfShortWeeksInTheBeginningOfQuarter
			case Day:
				position = quarterDay
				max = quarterEndDay
			default:
				invalidInput()
			}
		case Month:

			monthStart, monthEndPlus1 := Month.StartEnd(t)
			monthDay := t.Day()
			monthEndDay := monthEndPlus1.YearDay() - monthStart.YearDay()
			
			firstMonthWeekLength := int(7 - monthStart.Weekday()) % 7

			switch period {
			case FullWeek:
				lastMonthWeekLength := int(monthEndPlus1.Weekday()) % 7
				position =   (7 + monthDay    - firstMonthWeekLength) / 7
				max      =   (7 + monthEndDay - firstMonthWeekLength - lastMonthWeekLength - 1) / 7 
			case AnyWeek:
				countOfShortWeeksInTheBeginningOfMonth := (firstMonthWeekLength + 6) / 7 // 0 or 1
				position = (7 + monthDay    - firstMonthWeekLength) / 7 + countOfShortWeeksInTheBeginningOfMonth
				max      = (7 + monthEndDay - firstMonthWeekLength) / 7 + countOfShortWeeksInTheBeginningOfMonth
			case Day:
				position = t.Day()
				max = monthEndDay
				
			default:
				invalidInput()
			}
		case AnyWeek, FullWeek: // NB! this doesn't work very well for short weeks.
			switch period {
			case Day:
				position = int(t.Weekday()) + 1
				max = 7
			default:
				invalidInput()
			}
		case Day:
			switch period {
			case Hour:
				position = t.Hour() + 1
				max = 24
			case QuarterHour:
				position = t.Hour() * 4 + t.Minute() / 15  + 1
				max = 24 * 4
			case Minute:
				position = t.Hour() * 60 + t.Minute() + 1
				max = 60
			case Second:
				position = (t.Hour() * 60 + t.Minute()) * 60 + t.Second() + 1
				max = 60 * 60
			}
		case Hour:
			switch period {
			case QuarterHour:
				position = t.Minute() / 15 + 1
				max = 4
			case Minute:
				position = t.Minute() + 1
				max = 60
			case Second:
				position = t.Minute() * 60 + t.Second() + 1
				max = 60 * 60
			}
		default:
			invalidInput()
		}
	}
	return
}
// IsOnSchedule checks if the provided time satisfies the schedule
func (s Schedule)IsOnSchedule(t time.Time) (res bool) {
	parentPeriod := Epoch
	if s.Parent != nil {
		res = s.Parent.IsOnSchedule(t)
		if !res {
			return
		}
		p := &s
		for p != nil && p.Period == s.Period {
			p = p.Parent
			if p != nil {
				parentPeriod = p.Period
			}
		}
	}
	p, max := GetPosition(t, s.Period, parentPeriod)
	var start, end int
	if s.RangeEnd == 0 && s.RangeStart == 0 {
		// range is not filtered
		// expand the range by one in both directions
		start = 1 - 1
		end = max + 1
	} else {
		if s.RangeStart < 0 {
			start = max + s.RangeStart + 1
		} else {
			start = s.RangeStart
		}
		if s.RangeEnd < 0 {
			end = max + s.RangeEnd + 1
		} else {
			end = s.RangeEnd
		}
	}
	res = (p >= start && p <= end) && 		// within range
			(s.N == 0 ||                    // no need to check n-th
				((p - start) % s.N == s.M)) // phase is the same as given
	return
}
// ToBooleanSchedule converts schedule to a function from time to boolean
func (s Schedule)ToBooleanSchedule() BooleanSchedule {
	return s.IsOnSchedule
}

// Intervals return all intervals that this schedule produces
func (s Schedule)Intervals(i Interval) []Interval {
	parentIntervals := []Interval{i}
	if s.Parent != nil {
		parentIntervals = s.Parent.Intervals(i)
	}
	return IntervalsFlatMap(parentIntervals, 
		func(i Interval)(intervals []Interval){
			
			return
		},
	)
}
