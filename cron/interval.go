package cron

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

// Interval of time [StartInclusive, EndExclusive)
type Interval struct {
	StartInclusive time.Time
	EndExclusive   time.Time
}

// IntervalsFlatMap - applies function to each interval and concatenates all slices.
func IntervalsFlatMap(intervals []Interval, f func(Interval)[]Interval)(res []Interval) {
	for _, i := range intervals {
		res = append(res, f(i)...)
	}
	return
}
// GetIntervalByPosition returns subinterval of the length `period`
// that has index `position`.
// `position`  is numbered from 1. It's the first full period after 
// the start of the interval. One may use position 0 and negative
// to get intervals that start before the outer interval.
// Interval end is not respected in any way.
func (i Interval)GetIntervalByPosition(period Period, position int) (res Interval) {
	t := i.StartInclusive
	switch period {
	case Epoch:
		start := time.Date(1 + 10000 * position, 1, 1, 0, 0, 0, 0, time.UTC)
		res = Interval{
			StartInclusive: start,
			EndExclusive: start.AddDate(10000,0,0),
		}
	case Century:
		start := time.Date(t.Year()/ 100 * 100 + 1 + 100 * position, 1, 1, 0, 0, 0, 0, time.UTC)
		res = Interval{
			StartInclusive: start,
			EndExclusive: start.AddDate(100,0,0),
		}
	case Year:
		if t == time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC) { // if it's exactly on year boundary
			position = position - 1
		}
		start := time.Date(t.Year() + position, 1, 1, 0, 0, 0, 0, time.UTC)
		res = Interval{
			StartInclusive: start,
			EndExclusive: start.AddDate(1,0,0),
		}
	case Month:
		if t == time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC) { // if it's exactly on month boundary
			position = position - 1
		}
		start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC).
			AddDate(0, position, 0)
		res = Interval{
			StartInclusive: start,
			EndExclusive: start.AddDate(0,1,0),
		}
	case FullWeek:
		if t == time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC) && 
		t.Weekday() == time.Sunday { // if it's exactly on week boundary
			position = position - 1
		}
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).
			AddDate(0, 0, position * 7 - int(t.Weekday()))
		res = Interval{
			StartInclusive: start,
			EndExclusive: start.AddDate(0,0,7),
		}
	case AnyWeek:
		position = position - 1
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).
			AddDate(0, 0, position * 7 - int(t.Weekday()))
		res = Interval{
			StartInclusive: start,
			EndExclusive: start.AddDate(0,0,7),
		}
	case Day:
		if t == time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC) { // if it's exactly on day boundary
			position = position - 1
		}
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).
			AddDate(0, 0, position)
		res = Interval{
			StartInclusive: start,
			EndExclusive: start.AddDate(0,0,1),
		}
	default:
		panic(errors.New(fmt.Sprintf("GetIntervalByPosition(%v, %s, %d) is invalid", t, period, position)))	
	}
	return
}

// GetPosition returns position of the period that includes the given time moment.
func (i Interval)GetPosition(t time.Time, period Period) (position int) {
	start := i.StartInclusive
	// end   := i.EndExclusive
	tStart, _ := period.StartEnd(t)
	switch period {
	case Epoch:
		position = 0
	case Century:
		tAbsPos := (tStart.Year() - 1) / 100
		iAbsPos := (start.Year() - 1) / 100

		position = tAbsPos - iAbsPos
		if i.StartInclusive == time.Date(iAbsPos * 100 + 1, 1, 1, 0, 0, 0, 0, time.UTC) {
			position = position + 1
		}
	case Year:
		tAbsPos := tStart.Year()
		iAbsPos := start.Year()

		position = tAbsPos - iAbsPos
		if start == time.Date(start.Year(), 1, 1, 0, 0, 0, 0, time.UTC) {
			position = position + 1
		}
		
	case Month:
		tAbsPos := tStart.Year() * 12 + int(tStart.Month())
		iAbsPos :=  start.Year() * 12 + int( start.Month())

		position = tAbsPos - iAbsPos
		if start == time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC) {
			position = position + 1
		}
		
	case FullWeek:
		tAbsPos := (tStart.Unix() / 3600 / 24) / 7
		iAbsPos := ( start.Unix() / 3600 / 24 - int64(start.Weekday())) / 7

		position = int(tAbsPos - iAbsPos)
		if start == time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC) && 
		start.Weekday() == time.Sunday {
			position = position + 1
		}		
	case AnyWeek:
		tAbsPos := (tStart.Unix() / 3600 / 24) / 7
		iAbsPos := ( start.Unix() / 3600 / 24 - int64(start.Weekday())) / 7

		position = int(tAbsPos - iAbsPos)
		if start == time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC) && 
		start.Weekday() == time.Sunday {
			position = position + 2
		} else {
			position = position + 1
		}	
	case Day:
		tAbsPos := tStart.Unix() / 3600 / 24
		iAbsPos :=  start.Unix() / 3600 / 24

		position = int(tAbsPos - iAbsPos)
		if start == time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC) {
			position = position + 1
		}
	default:
		panic(errors.New(fmt.Sprintf("Interval.GetPosition(%v, %s) is invalid", t, period)))
	}
	return
}
