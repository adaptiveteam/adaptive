package cron_test

import (
	"github.com/stretchr/testify/assert"
	"time"
	"testing"
	"github.com/adaptiveteam/adaptive/cron"
)

func Test_GetPositionA(t *testing.T) {
	a := time.Date(2020, 4, 7, 0,0,0,0,time.UTC)
	a_day, a_mMax := cron.GetPosition(a, cron.Day, cron.Month)
	assert.Equal(t, 7, a_day)
	assert.Equal(t, 30, a_mMax)

	a_qDay, a_qDays := cron.GetPosition(a, cron.Day, cron.Quarter)
	assert.Equal(t, 7, a_qDay)
	assert.Equal(t, 91, a_qDays)

	a_wDay, a_wDays := cron.GetPosition(a, cron.Day, cron.FullWeek)
	assert.Equal(t, 3, a_wDay)
	assert.Equal(t, 7, a_wDays)

	a_mWeek, a_mwMax := cron.GetPosition(a, cron.FullWeek, cron.Month)
	assert.Equal(t, 1, a_mWeek)
	assert.Equal(t, 3, a_mwMax)

	a_mAnyWeek, a_mawMax := cron.GetPosition(a, cron.AnyWeek, cron.Month)
	assert.Equal(t, 2, a_mAnyWeek)
	assert.Equal(t, 5, a_mawMax)

	a_yWeek, a_myMax := cron.GetPosition(a, cron.FullWeek, cron.Year)
	assert.Equal(t, 14, a_yWeek)
	assert.Equal(t, 51, a_myMax)
}

func Test_GetPositionB(t *testing.T) {
	b := time.Date(2020, 2, 7, 0,0,0,0,time.UTC)
	b_day, b_mMax := cron.GetPosition(b, cron.Day, cron.Month)
	assert.Equal(t, 7, b_day)
	assert.Equal(t, 29, b_mMax)

	b_qDay, b_qDays := cron.GetPosition(b, cron.Day, cron.Quarter)
	assert.Equal(t, 31 + 7, b_qDay)
	assert.Equal(t, 91, b_qDays)
}

func Test_GetPositionC(t *testing.T) {
	с := time.Date(2020, 8, 31, 0,0,0,0,time.UTC)
	с_day, с_mMax := cron.GetPosition(с, cron.Day, cron.Month)
	assert.Equal(t, 31, с_day)
	assert.Equal(t, 31, с_mMax)

	с_mWeek, с_mwMax := cron.GetPosition(с, cron.FullWeek, cron.Month)
	assert.Equal(t, 5, с_mWeek)
	assert.Equal(t, 4, с_mwMax)

	с_mAnyWeek, с_mawMax := cron.GetPosition(с, cron.AnyWeek, cron.Month)
	assert.Equal(t, 6, с_mAnyWeek)
	assert.Equal(t, 6, с_mawMax)
}

func Test_IsOnSchedule(t *testing.T) {
	a := time.Date(2020, 4, 7, 0,0,0,0,time.UTC)
	
	assert.True(t, cron.S().IsOnSchedule(a))
	assert.True(t, cron.S().Every(cron.Day).IsOnSchedule(a))
	assert.True(t, cron.S().Every(cron.Month).InRange(cron.Day, 1, 10).IsOnSchedule(a))
	assert.False(t, cron.S().Every(cron.Month).InRange(cron.Day, 10, 20).IsOnSchedule(a))
	assert.True(t, cron.S().Every(cron.Month).InRange(cron.Day, 1, 10).EveryMN(cron.Day, 0, 7).IsOnSchedule(a))

	assert.True(t, cron.S().Every(cron.Month).InRange(cron.FullWeek, 1, 1).InRange(cron.WeekDay, 2, 2).IsOnSchedule(a))
	assert.False(t, cron.S().Every(cron.Month).InRange(cron.AnyWeek, 1, 1).InRange(cron.WeekDay, 2, 2).IsOnSchedule(a))
	assert.True(t, cron.S().Every(cron.Month).InRange(cron.AnyWeek, 2, 2).InRange(cron.WeekDay, 2, 2).IsOnSchedule(a))

}
