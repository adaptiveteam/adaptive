package cron_test

import (
	"github.com/stretchr/testify/assert"
	"time"
	"testing"
	"github.com/adaptiveteam/adaptive/cron"
)

func Test_GetPosition(t *testing.T) {
	a := time.Date(2020, 4, 8, 0,0,0,0,time.UTC)
	m := cron.Month.Interval(a)
	assert.Equal(t, 8, m.GetPosition(a, cron.Day))
	assert.Equal(t, 1, m.GetPosition(a, cron.FullWeek))
	assert.Equal(t, 2, m.GetPosition(a, cron.AnyWeek))
	assert.Equal(t, 0, m.GetPosition(a, cron.Year))
}

func Test_GetIntervalByPosition(t *testing.T) {
	a := time.Date(2020, 4, 8, 0,0,0,0,time.UTC)
	m := cron.Month.Interval(a)
	d := m.GetIntervalByPosition(cron.Day, 8)
	assert.Equal(t, 8, d.StartInclusive.Day())
	assert.Equal(t, 9, d.EndExclusive.Day())
	w := m.GetIntervalByPosition(cron.FullWeek, 2)
	assert.Equal(t, 12, w.StartInclusive.Day())
	assert.Equal(t, 19, w.EndExclusive.Day())
	aw := m.GetIntervalByPosition(cron.AnyWeek, 2)
	assert.Equal(t, 5, aw.StartInclusive.Day())
	assert.Equal(t, 12, aw.EndExclusive.Day())
}
