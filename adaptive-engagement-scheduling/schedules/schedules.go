package schedules

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
//            InRange(Monday,Friday).Every(Day)

// Period - one of Year, Quarter,...Day
type Period string

const (
	Century Period = "Century"
	Year Period = "Year"
	Quarter Period = "Quarter"
	Month Period = "Month"
	Week Period = "Week"
	Day Period = "Day"
)

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
		Period: Century,
		N: 1,
	}
}

// EveryMN makes a schedule that will be valid every n-th period on the m-th phase.
func (s Schedule)EveryMN(period Period, m, n int) Schedule {
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
func (s Schedule)InRange(period Period, start, end int) Schedule {
	return Schedule{
		Parent: &s,
		Period: period,
		RangeStart: start,
		RangeEnd: end,
	}
}
