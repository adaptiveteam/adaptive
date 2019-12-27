package business_time

import (
	"fmt"
	"strconv"
	"time"
)

// LocalTime represents local time without time zone
type LocalTime time.Time

const (
	// MeetingTimeFormat is a format string for converting time to string
	MeetingTimeFormat = "03:04 PM"
)

// FormatTimeAsMeetingTime converts given time to user-friendly representation
func FormatTimeAsMeetingTime(t LocalTime) string {
	return time.Time(t).Format(MeetingTimeFormat)
}


// MeetingTime constructs time for hour and minute
func MeetingTime(h, m int) LocalTime {
	return LocalTime(time.Date(0, 0, 0, h, m, 0,0, time.UTC))
}

// MeetingTimeRange generates meeting time range from startHour to endHourExclusive
// with interval of intervalMin minutes.
func MeetingTimeRange(startHour int, endHourExclusive int, intervalMin int) (times []LocalTime) {
	times = make([]LocalTime, (endHourExclusive - startHour) * (60 / intervalMin))
	i := 0
	for h := startHour; h < endHourExclusive; h++ {
		for m := 0; m < 60; m = m + intervalMin {
			t := MeetingTime(h, m)
			times[i] = t
			i ++
		}
	}
	return
}

// DefaultMeetingTimeRange returns time range from 8 AM to 6:45 PM.
func DefaultMeetingTimeRange() []LocalTime {
	return MeetingTimeRange(8, 19, 15)
}

// ID returns an identifier of the local time - like 6:45PM will become 1845. 
func (l LocalTime)ID() string {
	t := time.Time(l)
	return t.Format("15") + t.Format("04")
}

// ToUserFriendly converts LocalTime to 06:45 PM
func (l LocalTime)ToUserFriendly() string {
	return FormatTimeAsMeetingTime(l)
}
// ParseLocalTimeID converts identifier to LocalTime
func ParseLocalTimeID(id string) (res LocalTime, err error) {
	if len(id) == 4 {
		h, err := strconv.Atoi(id[0:2])
		if err == nil {
			m, err := strconv.Atoi(id[2:4])
			if err == nil {
				res = MeetingTime(h, m)
			}
		} 
	} else {
		err = fmt.Errorf("Cannot parse %s as LocalTime", id)
	}
	return
}
// ParseLocalTimeUserFriendly converts identifier to LocalTime
func ParseLocalTimeUserFriendly(text string) (res LocalTime, err error) {
	time, err := time.Parse(MeetingTimeFormat, text)
	return LocalTime(time), err
}
