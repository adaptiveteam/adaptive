package core_utils_go

import (
	"log"
	"strconv"
	"time"
)

type AdaptiveDateLayout string
type AdaptiveDateLocation string

const (
	// ISODateLayout = "2006-01-02" - YYYY-mm-DD
	ISODateLayout   AdaptiveDateLayout = "2006-01-02"
	// USDateLayout = "January 2, 2006"
	USDateLayout    AdaptiveDateLayout = "January 2, 2006"
	// TimestampLayout = "2006-01-02T15:04:05Z07:00"
	TimestampLayout AdaptiveDateLayout = time.RFC3339
	// TimestampWithoutTimeZoneLayout = "2006-01-02T15:04:05"
	TimestampWithoutTimeZoneLayout AdaptiveDateLayout = "2006-01-02T15:04:05"
	UTC AdaptiveDateLocation = "UTC"
)

func LocalToUtc(t time.Time) time.Time {
	// utc life
	loc, _ := time.LoadLocation("UTC")
	return t.In(loc)
}

func MonthStrToQuarter(month string) int {
	m, _ := strconv.Atoi(month)
	return MonthToQuarter(m)
}

func CurrentYearMonth() (int, int) {
	year, month, _ := time.Now().Date()
	return year, int(month)
}

func CurrentYearQuarter() (int, int) {
	year, month, _ := time.Now().Date()
	return year, MonthToQuarter(int(month))
}

// MonthToQuarter converts month [1..12] to quarter [1..4]
func MonthToQuarter(m int) int {
	return ((m - 1) / 3) + 1
}

// Bod returns the beginning of the day time for a given time
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func (a AdaptiveDateLayout) Parse(str string) (time.Time, error) {
	if str == "" {
		return time.Time{}, nil
	}
	return time.Parse(string(a), str)
}

func (a AdaptiveDateLayout) Format(t time.Time) string {
	return t.Format(string(a))
}

func (a AdaptiveDateLayout) ParseInLocation(str string, location AdaptiveDateLocation) (time.Time, error) {
	if str == "" {
		return time.Time{}, nil
	}
	l, err := time.LoadLocation(string(location))
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(string(a), str, l)
}

func (a AdaptiveDateLayout) ChangeLayout(str string, layout AdaptiveDateLayout) (string, error) {
	if str == "" {
		return "", nil
	}
	t, err := time.Parse(string(a), str)
	if err != nil {
		return "", err
	}
	return layout.Format(t), nil
}

func ISODateStringToUSDateString(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	t, err := ISODateLayout.Parse(str)
	if err != nil {
		return "", err
	}
	return USDateLayout.Format(t), nil
}

func CurrentRFCTimestamp() string {
	return TimestampLayout.Format(time.Now())
}

// ParseDateOrTimestamp handles the date/timestamp field value from database.
// It tries to parse 2006-01-02 or timestamp.
// if failed, logs failure and returns empty time with false flag.
func ParseDateOrTimestamp(dateOrTimestampOrEmpty string) (t time.Time, isDefined bool) {
	isDefined = true
	if dateOrTimestampOrEmpty == "" {
		isDefined = false
	} else {
		var err1 error
		var err2 error
		var err3 error
		t, err1 = ISODateLayout.Parse(dateOrTimestampOrEmpty)
		if err1 != nil {
			t, err2 = TimestampLayout.Parse(dateOrTimestampOrEmpty)
			if err2 != nil {
				t, err3 = TimestampWithoutTimeZoneLayout.Parse(dateOrTimestampOrEmpty)
				if err3 != nil {
					isDefined = false
					log.Printf("Couldn't parse %s as date (%v) or timestamp (%v) or timestamp-without-zone (%v)", dateOrTimestampOrEmpty, err1, err2, err3)
				}
			}
		}
	}
	return
}

// NormalizeTimestamp makes sure that the value is correct timestamp
func NormalizeTimestamp(dateOrTimestampOrEmpty string) (timestampOrEmpty string) {
	t, isDefined := ParseDateOrTimestamp(dateOrTimestampOrEmpty)
	if isDefined {
		timestampOrEmpty = TimestampLayout.Format(t)
	} else {
		timestampOrEmpty = ""
	}
	return
}

// NormalizeDate makes sure that the value is a correct date
func NormalizeDate(dateOrTimestampOrEmpty string) (dateOrEmpty string) {
	t, isDefined := ParseDateOrTimestamp(dateOrTimestampOrEmpty)
	if isDefined {
		dateOrEmpty = ISODateLayout.Format(t)
	} else {
		dateOrEmpty = ""
	}
	return
}
