package cron_test

import (
	"github.com/stretchr/testify/assert"
	"time"
	"testing"
	"github.com/adaptiveteam/adaptive/cron"
)

func Test_GetStartEnd(t *testing.T) {
	a := time.Date(2020, 4, 8, 0,0,0,0,time.UTC)
	ws, we := cron.FullWeek.StartEnd(a)
	assert.Equal(t, time.Date(2020, 4, 5, 0,0,0,0,time.UTC), ws)
	assert.Equal(t, time.Date(2020, 4, 12, 0,0,0,0,time.UTC), we)

	aws, awe := cron.AnyWeek.StartEnd(a)
	assert.Equal(t, time.Date(2020, 4, 5, 0,0,0,0,time.UTC), aws)
	assert.Equal(t, time.Date(2020, 4, 12, 0,0,0,0,time.UTC), awe)

	ms, me := cron.Month.StartEnd(a)
	assert.Equal(t, time.Date(2020, 4, 1, 0,0,0,0,time.UTC), ms)
	assert.Equal(t, time.Date(2020, 5, 1, 0,0,0,0,time.UTC), me)

	ds, de := cron.Day.StartEnd(a)
	assert.Equal(t, time.Date(2020, 4, 8, 0,0,0,0,time.UTC), ds)
	assert.Equal(t, time.Date(2020, 4, 9, 0,0,0,0,time.UTC), de)
}
