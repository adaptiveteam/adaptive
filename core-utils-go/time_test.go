package core_utils_go

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLocalToUtc(t *testing.T) {
	t1 := time.Now()
	fmt.Println(t1)
	fmt.Println(LocalToUtc(t1))
}

func TestDateUtils(t *testing.T) {
	monthMap := map[string]int{
		"1":  1,
		"2":  1,
		"3":  1,
		"4":  2,
		"5":  2,
		"6":  2,
		"7":  3,
		"8":  3,
		"9":  3,
		"10": 4,
		"11": 4,
		"12": 4,
	}
	for k, v := range monthMap {
		assert.True(t, MonthStrToQuarter(k) == v)
	}
}

func TestParseDateOrTimestamp(t *testing.T) {
	_, isDefined := ParseDateOrTimestamp("")
	assert.False(t, isDefined)
	_, isDefined2 := ParseDateOrTimestamp("019450-9")
	assert.False(t, isDefined2)
	date, isDefined3 := ParseDateOrTimestamp("2019-11-23")
	assert.True(t, isDefined3)
	assert.Equal(t, time.Date(2019,11,23,0,0,0,0, time.UTC), date)
}

func TestFormat(t *testing.T) {
	str := TimestampLayout.Format(time.Date(2019,11,23,0,0,0,0, time.UTC))
	assert.Equal(t, "2019-11-23T00:00:00Z",str)
	est, err := time.LoadLocation("EST")
	assert.Nil(t, err)
	str2 := TimestampLayout.Format(time.Date(2019,11,23,0,0,0,0, est))
	assert.Equal(t, "2019-11-23T00:00:00-05:00",str2)
}

func TestNormalizeDate(t *testing.T) {
	e := NormalizeDate("")
	assert.Equal(t, "", e)
	e2 := NormalizeDate("123")
	assert.Equal(t, "", e2)
	e3 := NormalizeDate("2019-11-23")
	assert.Equal(t, "2019-11-23", e3)
	e4 := NormalizeDate("2019-11-23T12:34:56Z")
	assert.Equal(t, "2019-11-23", e4)
	e5 := NormalizeDate("2019-11-23T12:34:56-03:00")
	assert.Equal(t, "2019-11-23", e5)
}

func TestNormalizeTimestamp(t *testing.T) {
	e := NormalizeTimestamp("")
	assert.Equal(t, "", e)
	e2 := NormalizeTimestamp("123")
	assert.Equal(t, "", e2)
	e3 := NormalizeTimestamp("2019-11-23")
	assert.Equal(t, "2019-11-23T00:00:00Z", e3)
	e4 := NormalizeTimestamp("2019-11-23T12:34:56Z")
	assert.Equal(t, "2019-11-23T12:34:56Z", e4)
	e5 := NormalizeTimestamp("2019-11-23T12:34:56-03:00")
	assert.Equal(t, "2019-11-23T12:34:56-03:00", e5)
	e6 := NormalizeTimestamp("2019-11-23T12:34:56")
	assert.Equal(t, "2019-11-23T12:34:56Z", e6)
}
