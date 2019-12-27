// Package businessTime provides more advanced time functions that
// align with common business cycles
package business_time

import (
	"bytes"
	"container/ring"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/adaptiveteam/adaptive-utils-go/models"
)

func EqualBytes(s1, s2 []byte) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal(s1, &o1)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 1 : %s", err.Error())
	}
	err = json.Unmarshal(s2, &o2)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 2 : %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

func Test_date_Year(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want int
	}{
		{name: "normal one",
			d: NewDate(
				2019,
				12,
				1),
			want: 2019},
		{name: "normal two",
			d: NewDate(
				2013,
				12,
				12),
			want: 2013},
		{name: "abnormal one",
			d: NewDate(
				1969,
				9,
				12),
			want: 1969},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetYear(); got != tt.want {
				t.Errorf("date.GetYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_EndDateList(t *testing.T) {
	endDates := make([]models.KvPair, 0)
	today := TimeToDate(time.Now().UTC())

	endOfQuarter := strings.Split(today.GetLastDayOfQuarter().DateToString(time.RFC3339), "T")[0]
	endOfQuarter += "(" + strconv.Itoa(today.GetLastDayOfQuarter().DaysBetween(today)) + " days)"
	endDates = append(endDates, models.KvPair{
		"End of Quarter ",
		endOfQuarter,
	})

	ninetyDays := strings.Split(today.AddTime(0, 0, 90).DateToString(time.RFC3339), "T")[0]
	endDates = append(endDates, models.KvPair{
		"Ninety Days",
		ninetyDays,
	})

	sixtyDays := strings.Split(today.AddTime(0, 0, 60).DateToString(time.RFC3339), "T")[0]
	endDates = append(endDates, models.KvPair{
		"Sixty Days",
		sixtyDays,
	})

	FortyFiveDays := strings.Split(today.AddTime(0, 0, 45).DateToString(time.RFC3339), "T")[0]
	endDates = append(endDates, models.KvPair{
		"Forty Five Days",
		FortyFiveDays,
	})

	thirtyDays := strings.Split(today.AddTime(0, 0, 30).DateToString(time.RFC3339), "T")[0]
	endDates = append(endDates, models.KvPair{
		"Thirty Days",
		thirtyDays,
	})

	fifteenDays := strings.Split(today.AddTime(0, 0, 15).DateToString(time.RFC3339), "T")[0]
	endDates = append(endDates, models.KvPair{
		"Fifteen Days",
		fifteenDays,
	})

	fmt.Println(endDates)
}

func Test_date_Month(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want int
	}{
		{name: "normal one",
			d: NewDate(
				2019,
				12,
				1),
			want: 12},
		{name: "normal two",
			d: NewDate(
				1969,
				12,
				12),
			want: 12},
		{name: "abnormal one",
			d: NewDate(
				1969,
				10,
				12),
			want: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetMonth(); got != tt.want {
				t.Errorf("date.GetMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_Day(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want int
	}{
		{name: "normal one",
			d: NewDate(
				2019,
				12,
				1),
			want: 1},
		{name: "normal two",
			d: NewDate(
				1969,
				12,
				12),
			want: 12},
		{name: "abnormal one",
			d: NewDate(
				1969,
				10,
				24),
			want: 24},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetDay(); got != tt.want {
				t.Errorf("date.GetDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NewDateFromQuarter(t *testing.T) {
	d := NewDateFromQuarter(1, 2019)
	if d.GetMonth() != 1 {
		t.Errorf("Got month %v, wanted month %v", d.GetMonth(), 1)
	}
	d = NewDateFromQuarter(2, 2019)
	if d.GetMonth() != 4 {
		t.Errorf("Got month %v, wanted month %v", d.GetMonth(), 4)
	}
	d = NewDateFromQuarter(3, 2019)
	if d.GetMonth() != 7 {
		t.Errorf("Got month %v, wanted month %v", d.GetMonth(), 7)
	}
	d = NewDateFromQuarter(4, 2019)
	if d.GetMonth() != 10 {
		t.Errorf("Got month %v, wanted month %v", d.GetMonth(), 10)
	}
}

func Test_date_Quarter(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want int
	}{
		{name: "January",
			d: NewDate(
				2019,
				1,
				30),
			want: 1},
		{name: "February",
			d: NewDate(
				1969,
				1,
				15),
			want: 1},
		{name: "March",
			d: NewDate(
				1969,
				3,
				24),
			want: 1},
		{name: "April",
			d: NewDate(
				2019,
				4,
				2),
			want: 2},
		{name: "May",
			d: NewDate(
				1969,
				5,
				15),
			want: 2},
		{name: "June",
			d: NewDate(
				1969,
				6,
				21),
			want: 2},
		{name: "July",
			d: NewDate(
				2019,
				7,
				30),
			want: 3},
		{name: "August",
			d: NewDate(
				1969,
				8,
				15),
			want: 3},
		{name: "September",
			d: NewDate(
				1969,
				9,
				28),
			want: 3},
		{name: "October",
			d: NewDate(
				2019,
				10,
				30),
			want: 4},
		{name: "November",
			d: NewDate(
				1969,
				11,
				15),
			want: 4},
		{name: "December",
			d: NewDate(
				1969,
				12,
				3),
			want: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetQuarter(); got != tt.want {
				t.Errorf("date.Quarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_DayOfWeek(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want int
	}{
		{name: "Sunday 1974-04-07",
			d: NewDate(
				1974,
				4,
				7,
			),
			want: 0},
		{name: "Monday 1949-01-24",
			d: NewDate(
				1949,
				1,
				24,
			),
			want: 1},
		{name: "Tuesday 1967-01-03",
			d: NewDate(
				1967,
				1,
				03),
			want: 2},
		{name: "Wednesday 1952-03-05",
			d: NewDate(
				1952,
				3,
				5),
			want: 3},
		{name: "Thursday 2009-09-24",
			d: NewDate(
				2009,
				9,
				24),
			want: 4},
		{name: "Friday 1982-06-18",
			d: NewDate(
				1982,
				6,
				18),
			want: 5},
		{name: "Saturday 2003-09-06",
			d: NewDate(
				2003,
				9,
				6),
			want: 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.DayOfWeek(); got != tt.want {
				t.Errorf("date.DayOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_DateBefore(t *testing.T) {
	var dates []Date
	for i := 0; i < 5; i++ {
		nd := NewDate(2019, 12, i+1)
		dates = append(dates, nd)
	}

	type args struct {
		d2        Date
		inclusive bool
	}

	tests := []struct {
		name string
		d    Date
		args args
		want bool
	}{
		{name: "Same day test",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				d2:        dates[0],
				inclusive: false,
			},
			want: false},
		{name: "Day before test",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				d2:        dates[1],
				inclusive: false,
			}, want: true},
		{name: "two days before",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				d2:        dates[2],
				inclusive: false,
			},
			want: true},
		{name: "One day before",
			d: NewDate(
				dates[2].GetYear(),
				dates[2].GetMonth(),
				dates[2].GetDay()),
			args: args{
				d2:        dates[3],
				inclusive: false,
			},
			want: true},
		{name: "Days ahead",
			d: NewDate(
				dates[4].GetYear(),
				dates[4].GetMonth(),
				dates[4].GetDay()),
			args: args{
				d2:        dates[0],
				inclusive: false,
			},
			want: false},
		{name: "same day",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				d2:        dates[0],
				inclusive: true,
			},
			want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.DateBefore(tt.args.d2, tt.args.inclusive); got != tt.want {
				t.Errorf("date.DateBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_DateAfter(t *testing.T) {
	var dates []Date
	for i := 0; i < 5; i++ {
		nd := NewDate(2019, 12, i+1)
		dates = append(dates, nd)
	}

	type args struct {
		d2        Date
		inclusive bool
	}
	tests := []struct {
		name string
		d    Date
		args args
		want bool
	}{
		{name: "Same day test",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				d2:        dates[0],
				inclusive: false,
			},
			want: false},
		{name: "Day before test",
			d: NewDate(
				dates[1].GetYear(),
				dates[1].GetMonth(),
				dates[1].GetDay()),
			args: args{
				d2:        dates[0],
				inclusive: false,
			},
			want: true},
		{name: "two days before",
			d: NewDate(
				dates[2].GetYear(),
				dates[2].GetMonth(),
				dates[2].GetDay()),
			args: args{
				d2:        dates[0],
				inclusive: false,
			},
			want: true},
		{name: "One day before",
			d: NewDate(
				dates[3].GetYear(),
				dates[3].GetMonth(),
				dates[3].GetDay()),
			args: args{
				d2:        dates[2],
				inclusive: false,
			},
			want: true},
		{name: "Days ahead",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				d2:        dates[4],
				inclusive: false,
			},
			want: false},
		{name: "same day",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				d2:        dates[0],
				inclusive: true,
			},
			want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.DateAfter(tt.args.d2, tt.args.inclusive); got != tt.want {
				t.Errorf("date.DateAfter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_AddTime(t *testing.T) {
	var dates []Date
	for i := 0; i < 5; i++ {
		nd := NewDate(2019, 12, i+1)
		dates = append(dates, nd)
	}
	type args struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name string
		d    Date
		args args
		want Date
	}{
		{name: "Add nothing",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				Year:  0,
				Month: 0,
				Day:   0,
			},
			want: dates[0]},
		{name: "Add one day",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				Year:  0,
				Month: 0,
				Day:   1,
			},
			want: dates[1]},
		{name: "Add nothing",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				Year:  0,
				Month: 0,
				Day:   0,
			},
			want: dates[0]},
		{name: "Add three days",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				Year:  0,
				Month: 0,
				Day:   3,
			},
			want: dates[3]},
		{name: "Subtract a day, same month",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()+1),
			args: args{
				Year:  0,
				Month: 0,
				Day:   -1,
			},
			want: dates[0]},
		{name: "Subtract a day, different month",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				Year:  0,
				Month: 0,
				Day:   -1,
			},
			want: TimeToDate(dates[0].DateToTimeMidnight().AddDate(0, 0, -1))},
		{name: "Add a month",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			args: args{
				Year:  0,
				Month: 1,
				Day:   0,
			},
			want: TimeToDate(dates[0].DateToTimeMidnight().AddDate(0, 1, 0))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.AddTime(tt.args.Year, tt.args.Month, tt.args.Day); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("date.AddTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_DateToTimeMidnight(t *testing.T) {
	var dates []Date
	for i := 0; i < 3; i++ {
		nd := NewDate(2019, 12, i+1)
		dates = append(dates, nd)
	}
	tests := []struct {
		name string
		d    Date
		want time.Time
	}{
		{name: "normal, #1",
			d: NewDate(
				dates[0].GetYear(),
				dates[0].GetMonth(),
				dates[0].GetDay()),
			want: time.Date(dates[0].GetYear(), time.Month(dates[0].GetMonth()), dates[0].GetDay(), 0, 0, 0, 0, time.UTC)},
		{name: "normal, #2",
			d: NewDate(
				dates[1].GetYear(),
				dates[1].GetMonth(),
				dates[1].GetDay()),
			want: time.Date(dates[1].GetYear(), time.Month(dates[1].GetMonth()), dates[1].GetDay(), 0, 0, 0, 0, time.UTC)},
		{name: "errant one",
			d: NewDate(
				2019,
				12,
				32),
			want: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.DateToTimeMidnight(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("date.DateToTimeMidnight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_DateToString(t *testing.T) {
	sourceDate := NewDate(2019, 12, 25)

	tests := []struct {
		name string
		d    Date
		want string
	}{
		{name: "normal",
			d:    sourceDate,
			want: sourceDate.DateToTimeMidnight().Format(time.RFC3339),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.DateToString(time.RFC3339); got != tt.want {
				t.Errorf("date.DateToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDate(t *testing.T) {
	normalDate := NewDate(2019, 12, 1)
	errorDate := NewDate(2019, 12, 1)
	type args struct {
		Year  int
		Month int
		Day   int
		Zone  string
	}
	tests := []struct {
		name string
		args args
		want Date
	}{
		{name: "Normal create",
			args: args{
				Year:  2019,
				Month: 12,
				Day:   1,
				Zone:  "Local",
			},
			want: normalDate,
		},
		{name: "Normal error",
			args: args{
				Year:  2019,
				Month: 12,
				Day:   1,
				Zone:  "Dodge a wrench",
			},
			want: errorDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDate(tt.args.Year, tt.args.Month, tt.args.Day)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToday(t *testing.T) {
	testTime := time.Now().UTC()
	type args struct {
		location string
	}
	tests := []struct {
		name string
		args args
		want Date
	}{
		{name: "Normal",
			args: args{

				location: "Local",
			},
			want: TimeToDate(testTime.In(time.Local)),
		},
		{name: "Error",
			args: args{

				location: "Dodge a wrench",
			},
			want: TimeToDate(testTime),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Today(time.UTC)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Today() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeToDate(t *testing.T) {
	var dates []Date
	for i := 0; i < 10; i++ {
		nd := NewDate(2019, 12, i+25)
		dates = append(dates, nd)
	}

	l, _ := time.LoadLocation("Local")
	var times []time.Time
	for i := 0; i < 10; i++ {
		nt := time.Date(2019, time.Month(12), i+25, 0, 0, 0, 0, l)
		times = append(times, nt)
	}

	type args struct {
		t time.Time
	}
	var tests []struct {
		name string
		args args
		want Date
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, struct {
			name string
			args args
			want Date
		}{
			name: strconv.Itoa(i),
			args: struct {
				t time.Time
			}{
				times[i],
			},
			want: dates[i],
		},
		)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeToDate(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeToDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateFromString(t *testing.T) {
	var dates []Date
	for i := 0; i < 10; i++ {
		nd := NewDate(2019, 12, i+25)
		dates = append(dates, nd)
	}

	var dateStrings []string
	for i := 0; i < 10; i++ {
		dateStrings = append(dateStrings, dates[i].DateToString(time.RFC3339))
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	var datesError []Date
	for i := 0; i < 10; i++ {
		nd := NewDate(2019, 12, i+25).DateToString(time.RFC3339)
		damaged := nd[:len(nd)-r1.Intn(5)]
		damagedDate := DateFromString(damaged)
		datesError = append(datesError, damagedDate)
	}

	type args struct {
		date string
	}
	var tests []struct {
		name string
		args args
		want Date
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, struct {
			name string
			args args
			want Date
		}{
			name: strconv.Itoa(i),
			args: struct {
				date string
			}{
				date: dateStrings[i],
			},
			want: dates[i]},
		)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DateFromString(tt.args.date)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateFromYMDString(t *testing.T) {
	type args struct {
		date string
	}
	type test struct {
		name string
		args args
		want Date
	}
	var tests [10]test

	for i := 0; i < 10; i++ {
		d := NewDate(2019, 12, i+25)
		tests[i].name = d.DateToString("2006-01-02")
		tests[i].want = d
		tests[i].args.date = d.DateToString("2006-01-02")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DateFromYMDString(tt.args.date)
			if err != nil {
				t.Errorf("DateFromYMDString() err = %v", err)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateFromYMDString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_holidays_HolidaysToJSON(t *testing.T) {
	hl := NewHolidayList()
	christmas := NewDate(2019, 12, 25)
	thanksgiving := NewDate(2019, 11, 28)
	diwali := NewDate(2019, 10, 27)
	hl.AddHoliday("Christmas", christmas, *time.UTC)
	hl.AddHoliday("Thanksgiving", thanksgiving, *time.UTC)
	hl.AddHoliday("Diwali", diwali, *time.UTC)

	targetJSON := []byte(`
[
{
		"what":"Diwali",
		"when":{
			"Year":2019,
			"Month":10,
			"Day":27
		},
		"where":"UTC"
},
{
		"what":"Thanksgiving",
		"when":{
			"Year":2019,
			"Month":11,
			"Day":28
		},
		"where":"UTC"
},
{
		"what":"Christmas",
		"when":{
			"Year":2019,
			"Month":12,
			"Day":25
		},
		"where":"UTC"
}
]`)

	buffer := new(bytes.Buffer)
	// json.Compact removes removes insignificant space characters
	if err := json.Compact(buffer, targetJSON); err != nil {
		fmt.Println(err)
	}

	tests := []struct {
		name string
		h    Holidays
		want []byte
	}{
		{
			name: "Initial",
			h:    hl,
			want: buffer.Bytes(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.h.HolidaysToJSON()
			res, err := EqualBytes(got, tt.want)
			if err != nil {
				t.Errorf("Could not compare the JSON strings")
			}
			if !res {
				t.Errorf("holidays.ToJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_holidays_HolidaysFromJSON(t *testing.T) {
	hl1 := NewHolidayList()
	christmas := NewDate(2019, 12, 25)
	thanksgiving := NewDate(2019, 11, 28)
	diwali := NewDate(2019, 10, 27)
	hl1.AddHoliday("Christmas", christmas, *time.UTC)
	hl1.AddHoliday("Thanksgiving", thanksgiving, *time.UTC)
	hl1.AddHoliday("Diwali", diwali, *time.UTC)

	var holidayJSON []byte
	hl2 := NewHolidayList()
	holidayJSON = hl1.HolidaysToJSON()
	err := HolidaysFromJSON(hl2, holidayJSON)
	if err == nil {
		holidaysArray := hl1.GetHolidays()

		for v := range holidaysArray {
			hn := holidaysArray[v]
			if !hl2.HolidayRegistered(hn) {
				t.Errorf("Holiday %v is not registered", hn)
			}
			d1 := hl1.GetHolidayDate(hn)
			d2 := hl2.GetHolidayDate(hn)
			if d1 != d2 {
				t.Errorf("Holiday %v is has incorrect date", hn)
			}
		}
	} else {
		t.Errorf("Wasn't able to unmarshall with error %v", err)
	}
}

func Test_holidays_AddHoliday(t *testing.T) {
	hl := NewHolidayList()
	testHolidays := [3]string{"Christmas", "Thanksgiving", "Diwali"}
	testDates := [3]Date{
		NewDate(2019, 12, 25),
		NewDate(2019, 11, 28),
		NewDate(2019, 10, 27),
	}

	for i := range testHolidays {
		hl.AddHoliday(testHolidays[i], testDates[i], *time.UTC)
		if !hl.HolidayRegistered(testHolidays[i]) {
			t.Errorf("Holiday %v is not registered", testHolidays[i])
			d := hl.GetHolidayDate(testHolidays[i])
			if d != testDates[i] {
				t.Errorf("Holiday %v doesn't have the correct date of %v", testHolidays[i], testDates[i])
			}
		}
	}

	for i := range testHolidays {
		if hl.AddHoliday(testHolidays[i], testDates[i], *time.UTC) {
			t.Errorf("Holiday %v was no registered correctly", testHolidays[i])
		}
	}

	if hl.HolidayRegistered("Kwanza") {
		t.Errorf("Never registered Kwanza!")
	}
}

func Test_holidays_DeleteHoliday(t *testing.T) {
	h := NewHolidayList()
	testHolidays := [3]string{"Christmas", "Thanksgiving", "Diwali"}
	testDates := [3]Date{
		NewDate(2019, 12, 25),
		NewDate(2019, 11, 28),
		NewDate(2019, 10, 27),
	}
	for i := range testHolidays {
		h.AddHoliday(testHolidays[i], testDates[i], *time.UTC)
		if !h.HolidayRegistered(testHolidays[i]) {
			t.Errorf("Holiday %v is not registered", testHolidays[i])
			d := h.GetHolidayDate(testHolidays[i])
			if d != testDates[i] {
				t.Errorf("Holiday %v doesn't have the correct date of %v", testHolidays[i], testDates[i])
			}
		}
	}

	for i := range testHolidays {
		l1 := h.Len()
		if !h.DeleteHoliday(testHolidays[i]) {
			t.Errorf("Holiday %v was not deleted correctly", testHolidays[i])
		}
		if h.DeleteHoliday(testHolidays[i]) {
			t.Errorf("Duplicate deletion not working correcting")
		}
		l2 := h.Len()
		if h.HolidayRegistered(testHolidays[i]) {
			t.Errorf("Holiday %v was not deleted", testHolidays[i])
		}
		if l1-l2 != 1 {
			t.Errorf("Length before and after don't make sense")
		}
	}

	if h.Len() != 0 {
		t.Errorf("Not all holidays deleted")
	}
}

func Test_holidays_HolidaysOnDate(t *testing.T) {
	h := NewHolidayList()
	testHolidays := []string{"Christmas", "Diwali", "Thanksgiving"}
	testDates := []Date{
		NewDate(2019, 12, 25),
		NewDate(2019, 12, 25),
		NewDate(2019, 12, 25),
	}

	for i := range testHolidays {
		h.AddHoliday(testHolidays[i], testDates[i], *time.UTC)
		if !h.HolidayRegistered(testHolidays[i]) {
			t.Errorf("Holiday %v is not registered", testHolidays[i])
			d := h.GetHolidayDate(testHolidays[i])
			if d != testDates[i] {
				t.Errorf("Holiday %v doesn't have the correct date of %v", testHolidays[i], testDates[i])
			}
		}
	}

	targetHolidays := h.HolidaysOnDate(NewDate(2019, 12, 25), time.UTC)
	sort.Strings(targetHolidays)
	sort.Strings(testHolidays)
	if !reflect.DeepEqual(targetHolidays, testHolidays) {
		t.Errorf("Wanted %v but got %v", targetHolidays, testHolidays)
	}
}

func Test_holidays_GetHolidayDate(t *testing.T) {
	hl := NewHolidayList()
	testHolidays := [3]string{"Christmas", "Thanksgiving", "Diwali"}
	testDates := [3]Date{
		NewDate(2019, 12, 25),
		NewDate(2019, 11, 28),
		NewDate(2019, 10, 27),
	}
	for i := range testHolidays {
		hl.AddHoliday(testHolidays[i], testDates[i], *time.UTC)
		if !hl.HolidayRegistered(testHolidays[i]) {
			t.Errorf("Holiday %v is not registered", testHolidays[i])
			d := hl.GetHolidayDate(testHolidays[i])
			if d != testDates[i] {
				t.Errorf("Holiday %v doesn't have the correct date of %v", testHolidays[i], testDates[i])
			}
		}
	}

	de := hl.GetHolidayDate("This is not a holiday")
	if de != nil {
		t.Errorf("Bad holiday given but didn't get nil in return")
	}

	type args struct {
		hn string
	}
	tests := []struct {
		name string
		h    Holidays
		args args
		want Date
	}{
		{
			name: "Christmas",
			h:    hl,
			args: struct {
				hn string
			}{
				hn: "Christmas",
			},
			want: NewDate(2019, 12, 25),
		},
		{
			name: "Thanksgiving",
			h:    hl,
			args: struct {
				hn string
			}{
				hn: "Thanksgiving",
			},
			want: NewDate(2019, 11, 28),
		},
		{
			name: "Diwali",
			h:    hl,
			args: struct {
				hn string
			}{
				hn: "Diwali",
			},
			want: NewDate(2019, 10, 27),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.GetHolidayDate(tt.args.hn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("holidays.GetHolidayDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBusinessDay(t *testing.T) {

	hl := NewHolidayList()
	hl.AddHoliday("Hump Day", NewDate(2019, 1, 30), *time.UTC)

	type args struct {
		d Date
		h Holidays
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Sunday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 27),
				h: hl,
			},
			want: false,
		},
		{
			name: "Monday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 28),
				h: hl,
			},
			want: true,
		},
		{
			name: "Tuesday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 29),
				h: hl,
			},
			want: true,
		},
		{
			name: "Wednesday with holiday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 30),
				h: hl,
			},
			want: false,
		},
		{
			name: "Thursday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 31),
				h: hl,
			},
			want: true,
		},
		{
			name: "Friday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 2, 1),
				h: hl,
			},
			want: true,
		},
		{
			name: "Saturday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 2, 2),
				h: hl,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.d.IsBusinessDay(tt.args.h, time.UTC); got != tt.want {
				t.Errorf("IsBusinessDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPreviousBusinessDay(t *testing.T) {
	hl := NewHolidayList()
	hl.AddHoliday("Day before New Years", NewDate(2018, 12, 28), *time.UTC)
	hl.AddHoliday("New Years", NewDate(2019, 1, 1), *time.UTC)

	type args struct {
		d Date
		h Holidays
	}
	tests := []struct {
		name string
		args args
		want Date
	}{
		{
			name: "Sunday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2018, 12, 30),
				h: hl,
			},
			want: NewDate(2018, 12, 27),
		},
		{
			name: "Monday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2018, 12, 31),
				h: hl,
			},
			want: NewDate(2018, 12, 27),
		},
		{
			name: "Tuesday is Holiday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 1),
				h: hl,
			},
			want: NewDate(2018, 12, 31),
		},
		{
			name: "Wednesday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 2),
				h: hl,
			},
			want: NewDate(2018, 12, 31),
		},
		{
			name: "Thursday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 3),
				h: hl,
			},
			want: NewDate(2019, 1, 2),
		},
		{
			name: "Friday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 4),
				h: hl,
			},
			want: NewDate(2019, 1, 3),
		},
		{
			name: "Saturday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 5),
				h: hl,
			},
			want: NewDate(2019, 1, 4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.d.PreviousBusinessDay(tt.args.h, time.UTC); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PreviousBusinessDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextBusinessDay(t *testing.T) {
	hl := NewHolidayList()
	hl.AddHoliday("Day before New Years", NewDate(2018, 12, 28), *time.UTC)
	hl.AddHoliday("New Years", NewDate(2019, 1, 1), *time.UTC)

	type args struct {
		d Date
		h Holidays
	}
	tests := []struct {
		name string
		args args
		want Date
	}{
		{
			name: "Sunday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2018, 12, 30),
				h: hl,
			},
			want: NewDate(2018, 12, 31),
		},
		{
			name: "Monday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2018, 12, 31),
				h: hl,
			},
			want: NewDate(2019, 1, 2),
		},
		{
			name: "Tuesday is Holiday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 1),
				h: hl,
			},
			want: NewDate(2019, 1, 2),
		},
		{
			name: "Wednesday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 2),
				h: hl,
			},
			want: NewDate(2019, 1, 3),
		},
		{
			name: "Thursday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 3),
				h: hl,
			},
			want: NewDate(2019, 1, 4),
		},
		{
			name: "Friday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 4),
				h: hl,
			},
			want: NewDate(2019, 1, 7),
		},
		{
			name: "Saturday",
			args: struct {
				d Date
				h Holidays
			}{
				d: NewDate(2019, 1, 5),
				h: hl,
			},
			want: NewDate(2019, 1, 7),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.d.NextBusinessDay(tt.args.h, time.UTC); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NextBusinessDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetQuarter(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want int
	}{
		{
			name: "Q119",
			d:    NewDate(2019, 1, 15),
			want: 1,
		},
		{
			name: "Q119",
			d:    NewDate(2019, 2, 3),
			want: 1,
		},
		{
			name: "Q119",
			d:    NewDate(2019, 3, 18),
			want: 1,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 4, 6),
			want: 2,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 5, 5),
			want: 2,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 6, 28),
			want: 2,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 7, 7),
			want: 3,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 8, 19),
			want: 3,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 9, 22),
			want: 3,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 10, 5),
			want: 4,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 11, 30),
			want: 4,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 12, 14),
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetQuarter(); got != tt.want {
				t.Errorf("date.GetQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetPreviousQuarter(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want int
	}{
		{
			name: "Q119",
			d:    NewDate(2019, 1, 15),
			want: 4,
		},
		{
			name: "Q119",
			d:    NewDate(2019, 2, 3),
			want: 4,
		},
		{
			name: "Q119",
			d:    NewDate(2019, 3, 18),
			want: 4,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 4, 6),
			want: 1,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 5, 5),
			want: 1,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 6, 28),
			want: 1,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 7, 7),
			want: 2,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 8, 19),
			want: 2,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 9, 22),
			want: 2,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 10, 5),
			want: 3,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 11, 30),
			want: 3,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 12, 14),
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetPreviousQuarter(); got != tt.want {
				t.Errorf("date.GetPreviousQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetNextQuarter(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want int
	}{
		{
			name: "Q119",
			d:    NewDate(2019, 1, 15),
			want: 2,
		},
		{
			name: "Q119",
			d:    NewDate(2019, 2, 3),
			want: 2,
		},
		{
			name: "Q119",
			d:    NewDate(2019, 3, 18),
			want: 2,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 4, 6),
			want: 3,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 5, 5),
			want: 3,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 6, 28),
			want: 3,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 7, 7),
			want: 4,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 8, 19),
			want: 4,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 9, 22),
			want: 4,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 10, 5),
			want: 1,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 11, 30),
			want: 1,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 12, 14),
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetNextQuarter(); got != tt.want {
				t.Errorf("date.GetNextQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetPreviousQuarterYear(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want int
	}{
		{
			name: "Q119",
			d:    NewDate(2019, 1, 15),
			want: 2018,
		},
		{
			name: "Q119",
			d:    NewDate(2019, 2, 3),
			want: 2018,
		},
		{
			name: "Q119",
			d:    NewDate(2019, 3, 18),
			want: 2018,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 4, 6),
			want: 2019,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 5, 5),
			want: 2019,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 6, 28),
			want: 2019,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 7, 7),
			want: 2019,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 8, 19),
			want: 2019,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 9, 22),
			want: 2019,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 10, 5),
			want: 2019,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 11, 30),
			want: 2019,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 12, 14),
			want: 2019,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetPreviousQuarterYear(); got != tt.want {
				t.Errorf("date.GetPreviousQuarterYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetNextQuarterYear(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want int
	}{
		{
			name: "Q119",
			d:    NewDate(2019, 1, 15),
			want: 2019,
		},
		{
			name: "Q119",
			d:    NewDate(2019, 2, 3),
			want: 2019,
		},
		{
			name: "Q119",
			d:    NewDate(2019, 3, 18),
			want: 2019,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 4, 6),
			want: 2019,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 5, 5),
			want: 2019,
		},
		{
			name: "Q219",
			d:    NewDate(2019, 6, 28),
			want: 2019,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 7, 7),
			want: 2019,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 8, 19),
			want: 2019,
		},
		{
			name: "Q319",
			d:    NewDate(2019, 9, 22),
			want: 2019,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 10, 5),
			want: 2020,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 11, 30),
			want: 2020,
		},
		{
			name: "Q419",
			d:    NewDate(2019, 12, 14),
			want: 2020,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetNextQuarterYear(); got != tt.want {
				t.Errorf("date.GetNextQuarterYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetLastDayOfQuarter(t *testing.T) {

	tests := []struct {
		name string
		d    Date
		want Date
	}{
		{
			name: "Q119",
			d:    NewDate(2019, 1, 15),
			want: NewDate(2019, 3, 31),
		},
		{
			name: "Q119",
			d:    NewDate(2019, 2, 3),
			want: NewDate(2019, 3, 31),
		},
		{
			name: "Q119",
			d:    NewDate(2019, 3, 18),
			want: NewDate(2019, 3, 31),
		},
		{
			name: "Q219",
			d:    NewDate(2019, 4, 6),
			want: NewDate(2019, 6, 30),
		},
		{
			name: "Q219",
			d:    NewDate(2019, 5, 5),
			want: NewDate(2019, 6, 30),
		},
		{
			name: "Q219",
			d:    NewDate(2019, 6, 28),
			want: NewDate(2019, 6, 30),
		},
		{
			name: "Q319",
			d:    NewDate(2019, 7, 7),
			want: NewDate(2019, 9, 30),
		},
		{
			name: "Q319",
			d:    NewDate(2019, 8, 19),
			want: NewDate(2019, 9, 30),
		},
		{
			name: "Q319",
			d:    NewDate(2019, 9, 22),
			want: NewDate(2019, 9, 30),
		},
		{
			name: "Q419",
			d:    NewDate(2019, 10, 5),
			want: NewDate(2019, 12, 31),
		},
		{
			name: "Q419",
			d:    NewDate(2019, 11, 30),
			want: NewDate(2019, 12, 31),
		},
		{
			name: "Q419",
			d:    NewDate(2019, 12, 14),
			want: NewDate(2019, 12, 31),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetLastDayOfQuarter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("date.GetLastDayOfQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetLastBusinessDayOfQuarter(t *testing.T) {
	holidays := NewHolidayList()
	holidays.AddHoliday("fake vacation Q1", NewDate(2019, 3, 29), *time.UTC)
	holidays.AddHoliday("fake vacation Q2", NewDate(2019, 6, 28), *time.UTC)

	type args struct {
		h Holidays
	}
	tests := []struct {
		name string
		d    Date
		args args
		want Date
	}{
		{
			name: "1, Q119",
			d:    NewDate(2019, 1, 21),
			want: NewDate(2019, 3, 28),
		},
		{
			name: "2, Q119",
			d:    NewDate(2019, 2, 3),
			want: NewDate(2019, 3, 28),
		},
		{
			name: "3, Q119",
			d:    NewDate(2019, 3, 18),
			want: NewDate(2019, 3, 28),
		},
		{
			name: "4, Q219",
			d:    NewDate(2019, 4, 6),
			want: NewDate(2019, 6, 27),
		},
		{
			name: "5, Q219",
			d:    NewDate(2019, 5, 5),
			want: NewDate(2019, 6, 27),
		},
		{
			name: "6, Q219",
			d:    NewDate(2019, 6, 28),
			want: NewDate(2019, 6, 27),
		},
		{
			name: "7, Q319",
			d:    NewDate(2019, 7, 7),
			want: NewDate(2019, 9, 30),
		},
		{
			name: "8, Q319",
			d:    NewDate(2019, 8, 19),
			want: NewDate(2019, 9, 30),
		},
		{
			name: "9, Q319",
			d:    NewDate(2019, 9, 22),
			want: NewDate(2019, 9, 30),
		},
		{
			name: "10, Q419",
			d:    NewDate(2019, 10, 5),
			want: NewDate(2019, 12, 31),
		},
		{
			name: "11, Q419",
			d:    NewDate(2019, 11, 30),
			want: NewDate(2019, 12, 31),
		},
		{
			name: "12, Q419",
			d:    NewDate(2019, 12, 14),
			want: NewDate(2019, 12, 31),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetLastBusinessDayOfQuarter(holidays, time.UTC); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("date.GetLastBusinessDayOfQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetFirstDayOfQuarter(t *testing.T) {

	tests := []struct {
		name string
		d    Date
		want Date
	}{
		{
			name: "Q119 - January",
			d:    NewDate(2019, 1, 21),
			want: NewDate(2019, 1, 1),
		},
		{
			name: "Q119 - February",
			d:    NewDate(2019, 2, 3),
			want: NewDate(2019, 1, 1),
		},
		{
			name: "Q119 - March",
			d:    NewDate(2019, 3, 18),
			want: NewDate(2019, 1, 1),
		},
		{
			name: "Q219 - April",
			d:    NewDate(2019, 4, 6),
			want: NewDate(2019, 4, 1),
		},
		{
			name: "Q219 - May",
			d:    NewDate(2019, 5, 5),
			want: NewDate(2019, 4, 1),
		},
		{
			name: "Q219 - June",
			d:    NewDate(2019, 6, 28),
			want: NewDate(2019, 4, 1),
		},
		{
			name: "Q319 - July",
			d:    NewDate(2019, 7, 7),
			want: NewDate(2019, 7, 1),
		},
		{
			name: "Q319 - August",
			d:    NewDate(2019, 8, 19),
			want: NewDate(2019, 7, 1),
		},
		{
			name: "Q319 - September",
			d:    NewDate(2019, 9, 22),
			want: NewDate(2019, 7, 1),
		},
		{
			name: "Q419 - October",
			d:    NewDate(2019, 10, 5),
			want: NewDate(2019, 10, 1),
		},
		{
			name: "Q419 - November",
			d:    NewDate(2019, 11, 30),
			want: NewDate(2019, 10, 1),
		},
		{
			name: "Q419 - December",
			d:    NewDate(2019, 12, 14),
			want: NewDate(2019, 10, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetFirstDayOfQuarter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("date.GetFirstBusinessDayOfQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDateFromQuarter(t *testing.T) {
	type args struct {
		quarter int
		year    int
	}
	tests := []struct {
		name string
		args args
		want Date
	}{
		{
			name: "Q119 - January",
			args: args{
				year:    2019,
				quarter: 1,
			},
			want: NewDate(2019, 1, 1),
		},
		{
			name: "Q219 - April",
			args: args{
				year:    2019,
				quarter: 2,
			},
			want: NewDate(2019, 4, 1),
		},
		{
			name: "Q319 - August",
			args: args{
				year:    2019,
				quarter: 3,
			},
			want: NewDate(2019, 7, 1),
		},
		{
			name: "Q119 - January",
			args: args{
				year:    2019,
				quarter: 4,
			},
			want: NewDate(2019, 10, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDateFromQuarter(tt.args.quarter, tt.args.year); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDateFromQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetDayOfWeekInQuarter(t *testing.T) {
	timeNow := time.Now().In(time.UTC)
	type fields struct {
		Year  int
		Month int
		Day   int
		Time  time.Time
	}
	type args struct {
		week int
		day  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantRv Date
	}{
		{
			name: "First Monday, first week week",
			fields: fields{
				Year:  2019,
				Month: 1,
				Day:   1,
				Time:  timeNow,
			},
			args: args{
				week: 0,
				day:  1,
			},
			wantRv: nil,
		},
		{
			name: "First Monday, second week, last day of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   31,
				Time:  timeNow,
			},
			args: args{
				week: 2,
				day:  1,
			},
			wantRv: NewDate(2019, 1, 14),
		},
		{
			name: "First Friday, second week, last day of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   31,
				Time:  timeNow,
			},
			args: args{
				week: 2,
				day:  5,
			},
			wantRv: NewDate(2019, 1, 18),
		},
		{
			name: "First Wednesday, second week, last day of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   31,
				Time:  timeNow,
			},
			args: args{
				week: 2,
				day:  3,
			},
			wantRv: NewDate(2019, 1, 16),
		},
		{
			name: "First Monday, second week, last Thursday of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   28,
				Time:  timeNow,
			},
			args: args{
				week: 2,
				day:  1,
			},
			wantRv: NewDate(2019, 1, 14),
		},
		{
			name: "First Friday, second week, last day of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   28,
				Time:  timeNow,
			},
			args: args{
				week: 2,
				day:  5,
			},
			wantRv: NewDate(2019, 1, 18),
		},
		{
			name: "First Wednesday, second week, last day of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   28,
				Time:  timeNow,
			},
			args: args{
				week: 2,
				day:  3,
			},
			wantRv: NewDate(2019, 1, 16),
		},
		{
			name: "Second to last Monday, second week, last Thursday of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   28,
				Time:  timeNow,
			},
			args: args{
				week: -2,
				day:  1,
			},
			wantRv: NewDate(2019, 3, 18),
		},
		{
			name: "Second to last  Friday, second week, last day of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   28,
				Time:  timeNow,
			},
			args: args{
				week: -2,
				day:  5,
			},
			wantRv: NewDate(2019, 3, 22),
		},
		{
			name: "Second to last  Wednesday, second week, last day of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   28,
				Time:  timeNow,
			},
			args: args{
				week: -2,
				day:  3,
			},
			wantRv: NewDate(2019, 3, 20),
		},
		{
			name: "Second to last  Monday, last day of March, 2019, starting from Saturday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   23,
				Time:  timeNow,
			},
			args: args{
				week: -2,
				day:  1,
			},
			wantRv: NewDate(2019, 3, 18),
		},
		{
			name: "Second to last  Wednesday, last day of March, 2019, starting from Saturda",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   23,
				Time:  timeNow,
			},
			args: args{
				week: -2,
				day:  3,
			},
			wantRv: NewDate(2019, 3, 20),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if gotRv := d.GetDayOfWeekInQuarter(tt.args.week, tt.args.day); !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("date.GetDayOfWeekInQuarter() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func Test_date_GetLastWeekDayOfQuarter(t *testing.T) {

	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name   string
		fields fields
		want   Date
	}{
		{
			name: "Quarter #1, February",
			fields: fields{
				Year:  2019,
				Month: 2,
				Day:   15,
			},
			want: NewDateFromQuarter(1, 2019).GetLastDayOfQuarter().GetWeekDay(false),
		},
		{
			name: "Quarter #2, May",
			fields: fields{
				Year:  2019,
				Month: 5,
				Day:   31,
			},
			want: NewDateFromQuarter(2, 2019).GetLastDayOfQuarter().GetWeekDay(false),
		},
		{
			name: "Quarter #3, June",
			fields: fields{
				Year:  2019,
				Month: 7,
				Day:   28,
			},
			want: NewDateFromQuarter(3, 2019).GetLastDayOfQuarter().GetWeekDay(false),
		},
		{
			name: "Quarter #4, October",
			fields: fields{
				Year:  2019,
				Month: 10,
				Day:   1,
			},
			want: NewDateFromQuarter(4, 2019).GetLastDayOfQuarter().GetWeekDay(false),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if got := d.GetLastWeekDayOfQuarter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("date.GetLastWeekDayOfQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_IsWeekDay(t *testing.T) {

	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Sunday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   24,
			},
			want: false,
		},
		{
			name: "Monday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   25,
			},
			want: true,
		},
		{
			name: "Tuesday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   26,
			},
			want: true,
		},
		{
			name: "Wednesday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   27,
			},
			want: true,
		},
		{
			name: "Thursday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   28,
			},
			want: true,
		},
		{
			name: "Friday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   29,
			},
			want: true,
		},
		{
			name: "Saturday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   30,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if got := d.IsWeekDay(); got != tt.want {
				t.Errorf("date.IsWeekDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetWeekDay(t *testing.T) {

	type fields struct {
		Year  int
		Month int
		Day   int
	}
	type args struct {
		forward bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantRv Date
	}{
		{
			name: "Sunday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   24,
			},
			args: args{
				forward: false,
			},
			wantRv: NewDate(2019, 3, 22),
		},
		{
			name: "Monday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   25,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 3, 25),
		},
		{
			name: "Tuesday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   26,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 3, 26),
		},
		{
			name: "Wednesday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   27,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 3, 27),
		},
		{
			name: "Thursday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   28,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 3, 28),
		},
		{
			name: "Friday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   29,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 3, 29),
		},
		{
			name: "Saturday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   30,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 4, 1),
		},
		{
			name: "Sunday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   24,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 3, 25),
		},
		{
			name: "Monday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   25,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 3, 25),
		},
		{
			name: "Tuesday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   26,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 3, 26),
		},
		{
			name: "Wednesday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   27,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 3, 27),
		},
		{
			name: "Thursday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   28,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 3, 28),
		},
		{
			name: "Friday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   29,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 3, 29),
		},
		{
			name: "Saturday",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   30,
			},
			args: args{
				forward: true,
			},
			wantRv: NewDate(2019, 4, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if gotRv := d.GetWeekDay(tt.args.forward); !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("date.GetWeekDay() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func Test_date_GetFirstDayOfWeek(t *testing.T) {

	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name   string
		fields fields
		want   Date
	}{
		{
			name: "1st week of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   2,
			},
			// This will be the previous month
			want: NewDate(2019, 2, 24),
		},
		{
			name: "2nd week of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   7,
			},
			want: NewDate(2019, 3, 3),
		},
		{
			name: "3rd week of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   13,
			},
			want: NewDate(2019, 3, 10),
		},
		{
			name: "4th week of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   21,
			},
			want: NewDate(2019, 3, 17),
		},
		{
			name: "5th week of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   29,
			},
			want: NewDate(2019, 3, 24),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if got := d.GetFirstDayOfWeek(); got != tt.want {
				t.Errorf("date.GetFirstDayOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetLastDayOfWeek(t *testing.T) {

	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name   string
		fields fields
		want   Date
	}{
		{
			name: "1st week of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   2,
			},
			// This will be the previous month
			want: NewDate(2019, 3, 2),
		},
		{
			name: "2nd week of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   7,
			},
			want: NewDate(2019, 3, 9),
		},
		{
			name: "3rd week of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   13,
			},
			want: NewDate(2019, 3, 16),
		},
		{
			name: "4th week of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   21,
			},
			want: NewDate(2019, 3, 23),
		},
		{
			name: "5th week of March, 2019",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   29,
			},
			want: NewDate(2019, 3, 30),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if got := d.GetLastDayOfWeek(); got != tt.want {
				t.Errorf("date.GetLastDayOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetBusinessDay(t *testing.T) {
	hl := NewHolidayList()
	christmas := NewDate(2019, 12, 25)
	thanksgiving := NewDate(2019, 11, 28)
	diwali := NewDate(2019, 10, 27)
	hl.AddHoliday("Christmas", christmas, *time.UTC)
	hl.AddHoliday("Thanksgiving", thanksgiving, *time.UTC)
	hl.AddHoliday("Diwali", diwali, *time.UTC)

	type fields struct {
		Year  int
		Month int
		Day   int
	}
	type args struct {
		h       Holidays
		forward bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantRv Date
	}{
		{
			name: "Christmas Day, forward",
			fields: fields{
				Year:  2019,
				Month: 12,
				Day:   25,
			},
			args: args{
				h:       hl,
				forward: true,
			},
			wantRv: NewDate(2019, 12, 26),
		},
		{
			name: "Thanksgiving, forward",
			fields: fields{
				Year:  2019,
				Month: 11,
				Day:   28,
			},
			args: args{
				h:       hl,
				forward: true,
			},
			wantRv: NewDate(2019, 11, 29),
		},
		{
			name: "Diwali, forward",
			fields: fields{
				Year:  2019,
				Month: 10,
				Day:   27,
			},
			args: args{
				h:       hl,
				forward: true,
			},
			wantRv: NewDate(2019, 10, 28),
		},
		{
			name: "No holiday, forward",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   27,
			},
			args: args{
				h:       hl,
				forward: true,
			},
			wantRv: NewDate(2019, 3, 27),
		},
		{
			name: "Christmas Day, backward",
			fields: fields{
				Year:  2019,
				Month: 12,
				Day:   25,
			},
			args: args{
				h:       hl,
				forward: false,
			},
			wantRv: NewDate(2019, 12, 24),
		},
		{
			name: "Thanksgiving, backward",
			fields: fields{
				Year:  2019,
				Month: 11,
				Day:   28,
			},
			args: args{
				h:       hl,
				forward: false,
			},
			wantRv: NewDate(2019, 11, 27),
		},
		{
			name: "Diwali, backward",
			fields: fields{
				Year:  2019,
				Month: 10,
				Day:   27,
			},
			args: args{
				h:       hl,
				forward: false,
			},
			wantRv: NewDate(2019, 10, 25),
		},
		{
			name: "No holiday, backward",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   27,
			},
			args: args{
				h:       hl,
				forward: false,
			},
			wantRv: NewDate(2019, 3, 27),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if gotRv := d.GetBusinessDay(tt.args.h, time.UTC, tt.args.forward); !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("date.GetBusinessDay() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func Test_date_GetDateInLocation(t *testing.T) {
	indianaZone, _ := time.LoadLocation("US/East-Indiana")
	parisZone, _ := time.LoadLocation("Europe/Paris")
	indiaZone, _ := time.LoadLocation("Asia/Kolkata")

	type fields struct {
		Year  int
		Month int
		Day   int
	}
	type args struct {
		l *time.Location
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Date
	}{
		{
			name: "test one",
			fields: fields{
				2020,
				4,
				13,
			},
			args: args{
				indiaZone,
			},
			want: TimeToDate(time.Date(2020, 4, 13, 0, 0, 0, 0, indiaZone)),
		},
		{
			name: "test two",
			fields: fields{
				2020,
				4,
				13,
			},
			args: args{
				indianaZone,
			},
			want: TimeToDate(time.Date(2020, 4, 13, 0, 0, 0, 0, indianaZone)),
		},
		{
			name: "test three",
			fields: fields{
				2020,
				4,
				13,
			},
			args: args{
				parisZone,
			},
			want: TimeToDate(time.Date(2020, 4, 13, 0, 0, 0, 0, parisZone)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if got := d.GetDateInLocation(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("date.GetDateInLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_DaysBetween(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	type args struct {
		d2 Date
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "ten days",
			fields: fields{
				2019,
				4,
				15,
			},
			args: args{
				d2: NewDate(2019, 4, 25),
			},
			want: 10,
		},
		{
			name: "ten days",
			fields: fields{
				2019,
				4,
				25,
			},
			args: args{
				d2: NewDate(2019, 5, 8),
			},
			want: 13,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if got := d.DaysBetween(tt.args.d2); got != tt.want {
				t.Errorf("date.DaysBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_holidays_GetHolidayLocation(t *testing.T) {
	indiaZone, _ := time.LoadLocation("Asia/Kolkata")

	type args struct {
		name string
	}
	tests := []struct {
		name   string
		h      holidays
		args   args
		wantRv *time.Location
	}{
		{
			name: "test one",
			h: holidays{
				"Diwali": {
					When: date{
						Year:  2019,
						Month: 10,
						Day:   27,
					},
					Where: indiaZone.String(),
					What:  "Diwali",
				},
			},
			args: args{
				"Diwali",
			},
			wantRv: indiaZone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRv := tt.h.GetHolidayLocation(tt.args.name); !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("holidays.GetHolidayLocation() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func Test_holidays_GetHolidays(t *testing.T) {
	indiaZone, _ := time.LoadLocation("Asia/Kolkata")
	tests := []struct {
		name string
		h    holidays
		want []string
	}{
		{
			name: "test one",
			h: holidays{
				"Diwali": {
					When: date{
						Year:  2019,
						Month: 10,
						Day:   27,
					},
					Where: indiaZone.String(),
					What:  "Diwali",
				},
			},
			want: []string{"Diwali"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.GetHolidays(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("holidays.GetHolidays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_holidays_GetHoliday(t *testing.T) {
	indiaZone, _ := time.LoadLocation("Asia/Kolkata")
	diwali := NewDate(2019, 10, 27)
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		h      holidays
		args   args
		wantRv Date
	}{
		{
			name: "test one",
			h: holidays{
				"Diwali": {
					When: date{
						Year:  2019,
						Month: 10,
						Day:   27,
					},
					Where: indiaZone.String(),
					What:  "Diwali",
				},
			},
			args: args{
				"Diwali",
			},
			wantRv: diwali,
		},
		{
			name: "failure",
			h: holidays{
				"Diwali": {
					When: date{
						Year:  2019,
						Month: 10,
						Day:   27,
					},
					Where: indiaZone.String(),
					What:  "Diwali",
				},
			},
			args: args{
				"Dude, where's my car?",
			},
			wantRv: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRv := tt.h.GetHoliday(tt.args.name); !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("holidays.GetHoliday() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func Test_date_GetFirstWeekDayOfQuarter(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name   string
		fields fields
		want   Date
	}{
		{
			name: "quarter one",
			fields: fields{
				Year:  2019,
				Month: 1,
				Day:   1,
			},
			want: NewDate(2019, 1, 1),
		},
		{
			name: "quarter two",
			fields: fields{
				Year:  2019,
				Month: 4,
				Day:   1,
			},
			want: NewDate(2019, 4, 1),
		},
		{
			name: "quarter three",
			fields: fields{
				Year:  2019,
				Month: 7,
				Day:   1,
			},
			want: NewDate(2019, 7, 1),
		},
		{
			name: "quarter four",
			fields: fields{
				Year:  2019,
				Month: 10,
				Day:   1,
			},
			want: NewDate(2019, 10, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if got := d.GetFirstWeekDayOfQuarter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("date.GetFirstWeekDayOfQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetFirstBusinessDayOfQuarter(t *testing.T) {
	hl := NewHolidayList()
	holidayOne := NewDate(2019, 1, 1)
	holidayTwo := NewDate(2019, 4, 1)
	holidayThree := NewDate(2019, 10, 1)
	hl.AddHoliday("holidayOne", holidayOne, *time.UTC)
	hl.AddHoliday("holidayTwo", holidayTwo, *time.UTC)
	hl.AddHoliday("holidayThree", holidayThree, *time.UTC)
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	type args struct {
		h        Holidays
		location *time.Location
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Date
	}{
		{
			name: "quarter one",
			fields: fields{
				Year:  2019,
				Month: 1,
				Day:   1,
			},
			want: NewDate(2019, 1, 2),
			args: args{
				hl,
				time.UTC,
			},
		},
		{
			name: "quarter two",
			fields: fields{
				Year:  2019,
				Month: 4,
				Day:   1,
			},
			want: NewDate(2019, 4, 2),
			args: args{
				hl,
				time.UTC,
			},
		},
		{
			name: "quarter three",
			fields: fields{
				Year:  2019,
				Month: 7,
				Day:   1,
			},
			want: NewDate(2019, 7, 1),
			args: args{
				hl,
				time.UTC,
			},
		},
		{
			name: "quarter four",
			fields: fields{
				Year:  2019,
				Month: 10,
				Day:   1,
			},
			want: NewDate(2019, 10, 2),
			args: args{
				hl,
				time.UTC,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if got := d.GetFirstBusinessDayOfQuarter(tt.args.h, tt.args.location); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("date.GetFirstBusinessDayOfQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_date_GetFirstDayOfMonth(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name   string
		fields fields
		wantRv Date
	}{
		{
			name: "January",
			fields: fields{
				Year:  2019,
				Month: 1,
				Day:   15,
			},
			wantRv: NewDate(2019, 1, 1),
		},
		{
			name: "February",
			fields: fields{
				Year:  2019,
				Month: 2,
				Day:   15,
			},
			wantRv: NewDate(2019, 2, 1),
		},
		{
			name: "March",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   15,
			},
			wantRv: NewDate(2019, 3, 1),
		},
		{
			name: "April",
			fields: fields{
				Year:  2019,
				Month: 4,
				Day:   15,
			},
			wantRv: NewDate(2019, 4, 1),
		},
		{
			name: "May",
			fields: fields{
				Year:  2019,
				Month: 5,
				Day:   15,
			},
			wantRv: NewDate(2019, 5, 1),
		},
		{
			name: "June",
			fields: fields{
				Year:  2019,
				Month: 6,
				Day:   15,
			},
			wantRv: NewDate(2019, 6, 1),
		},
		{
			name: "July",
			fields: fields{
				Year:  2019,
				Month: 7,
				Day:   15,
			},
			wantRv: NewDate(2019, 7, 1),
		},
		{
			name: "August",
			fields: fields{
				Year:  2019,
				Month: 8,
				Day:   15,
			},
			wantRv: NewDate(2019, 8, 1),
		},
		{
			name: "September",
			fields: fields{
				Year:  2019,
				Month: 9,
				Day:   15,
			},
			wantRv: NewDate(2019, 9, 1),
		},
		{
			name: "October",
			fields: fields{
				Year:  2019,
				Month: 10,
				Day:   15,
			},
			wantRv: NewDate(2019, 10, 1),
		},
		{
			name: "November",
			fields: fields{
				Year:  2019,
				Month: 11,
				Day:   15,
			},
			wantRv: NewDate(2019, 11, 1),
		},
		{
			name: "December",
			fields: fields{
				Year:  2019,
				Month: 12,
				Day:   15,
			},
			wantRv: NewDate(2019, 12, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if gotRv := d.GetFirstDayOfMonth(); !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("date.GetFirstDayOfMonth() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func Test_date_GetLastDayOfMonth(t *testing.T) {
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	tests := []struct {
		name   string
		fields fields
		wantRv Date
	}{
		{
			name: "January",
			fields: fields{
				Year:  2019,
				Month: 1,
				Day:   15,
			},
			wantRv: NewDate(2019, 1, 31),
		},
		{
			name: "February",
			fields: fields{
				Year:  2019,
				Month: 2,
				Day:   15,
			},
			wantRv: NewDate(2019, 2, 28),
		},
		{
			name: "March",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   15,
			},
			wantRv: NewDate(2019, 3, 31),
		},
		{
			name: "April",
			fields: fields{
				Year:  2019,
				Month: 4,
				Day:   15,
			},
			wantRv: NewDate(2019, 4, 30),
		},
		{
			name: "May",
			fields: fields{
				Year:  2019,
				Month: 5,
				Day:   15,
			},
			wantRv: NewDate(2019, 5, 31),
		},
		{
			name: "June",
			fields: fields{
				Year:  2019,
				Month: 6,
				Day:   15,
			},
			wantRv: NewDate(2019, 6, 30),
		},
		{
			name: "July",
			fields: fields{
				Year:  2019,
				Month: 7,
				Day:   15,
			},
			wantRv: NewDate(2019, 7, 31),
		},
		{
			name: "August",
			fields: fields{
				Year:  2019,
				Month: 8,
				Day:   15,
			},
			wantRv: NewDate(2019, 8, 31),
		},
		{
			name: "September",
			fields: fields{
				Year:  2019,
				Month: 9,
				Day:   15,
			},
			wantRv: NewDate(2019, 9, 30),
		},
		{
			name: "October",
			fields: fields{
				Year:  2019,
				Month: 10,
				Day:   15,
			},
			wantRv: NewDate(2019, 10, 31),
		},
		{
			name: "November",
			fields: fields{
				Year:  2019,
				Month: 11,
				Day:   15,
			},
			wantRv: NewDate(2019, 11, 30),
		},
		{
			name: "December",
			fields: fields{
				Year:  2019,
				Month: 12,
				Day:   15,
			},
			wantRv: NewDate(2019, 12, 31),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if gotRv := d.GetLastDayOfMonth(); !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("date.GetLastDayOfMonth() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func Test_date_GetFirstBusinessDayOfMonth(t *testing.T) {
	hl := NewHolidayList()
	holidayOne := NewDate(2019, 1, 1)
	holidayTwo := NewDate(2019, 4, 1)
	holidayThree := NewDate(2019, 10, 1)
	hl.AddHoliday("holidayOne", holidayOne, *time.UTC)
	hl.AddHoliday("holidayTwo", holidayTwo, *time.UTC)
	hl.AddHoliday("holidayThree", holidayThree, *time.UTC)
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	type args struct {
		h Holidays
		l *time.Location
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantRv Date
	}{
		{
			name: "January",
			fields: fields{
				Year:  2019,
				Month: 1,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 1, 2),
		},
		{
			name: "February",
			fields: fields{
				Year:  2019,
				Month: 2,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 2, 1),
		},
		{
			name: "March",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 3, 1),
		},
		{
			name: "April",
			fields: fields{
				Year:  2019,
				Month: 4,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 4, 2),
		},
		{
			name: "May",
			fields: fields{
				Year:  2019,
				Month: 5,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 5, 1),
		},
		{
			name: "June",
			fields: fields{
				Year:  2019,
				Month: 6,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 6, 3),
		},
		{
			name: "July",
			fields: fields{
				Year:  2019,
				Month: 7,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 7, 1),
		},
		{
			name: "August",
			fields: fields{
				Year:  2019,
				Month: 8,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 8, 1),
		},
		{
			name: "September",
			fields: fields{
				Year:  2019,
				Month: 9,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 9, 2),
		},
		{
			name: "October",
			fields: fields{
				Year:  2019,
				Month: 10,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 10, 2),
		},
		{
			name: "November",
			fields: fields{
				Year:  2019,
				Month: 11,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 11, 1),
		},
		{
			name: "December",
			fields: fields{
				Year:  2019,
				Month: 12,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 12, 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if gotRv := d.GetFirstBusinessDayOfMonth(tt.args.h, tt.args.l); !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("date.GetFirstBusinessDayOfMonth() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func Test_date_GetLastBusinessDayOfMonth(t *testing.T) {
	hl := NewHolidayList()
	holidayOne := NewDate(2019, 1, 31)
	holidayTwo := NewDate(2019, 4, 30)
	holidayThree := NewDate(2019, 10, 30)
	hl.AddHoliday("holidayOne", holidayOne, *time.UTC)
	hl.AddHoliday("holidayTwo", holidayTwo, *time.UTC)
	hl.AddHoliday("holidayThree", holidayThree, *time.UTC)
	type fields struct {
		Year  int
		Month int
		Day   int
	}
	type args struct {
		h Holidays
		l *time.Location
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantRv Date
	}{
		{
			name: "January",
			fields: fields{
				Year:  2019,
				Month: 1,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 1, 30),
		},
		{
			name: "February",
			fields: fields{
				Year:  2019,
				Month: 2,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 2, 28),
		},
		{
			name: "March",
			fields: fields{
				Year:  2019,
				Month: 3,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 3, 29),
		},
		{
			name: "April",
			fields: fields{
				Year:  2019,
				Month: 4,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 4, 29),
		},
		{
			name: "May",
			fields: fields{
				Year:  2019,
				Month: 5,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 5, 31),
		},
		{
			name: "June",
			fields: fields{
				Year:  2019,
				Month: 6,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 6, 28),
		},
		{
			name: "July",
			fields: fields{
				Year:  2019,
				Month: 7,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 7, 31),
		},
		{
			name: "August",
			fields: fields{
				Year:  2019,
				Month: 8,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 8, 30),
		},
		{
			name: "September",
			fields: fields{
				Year:  2019,
				Month: 9,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 9, 30),
		},
		{
			name: "October",
			fields: fields{
				Year:  2019,
				Month: 10,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 10, 31),
		},
		{
			name: "November",
			fields: fields{
				Year:  2019,
				Month: 11,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 11, 29),
		},
		{
			name: "December",
			fields: fields{
				Year:  2019,
				Month: 12,
				Day:   15,
			},
			args: args{
				h: hl,
				l: time.UTC,
			},
			wantRv: NewDate(2019, 12, 31),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := date{
				Year:  tt.fields.Year,
				Month: tt.fields.Month,
				Day:   tt.fields.Day,
			}
			if gotRv := d.GetLastBusinessDayOfMonth(tt.args.h, tt.args.l); !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("date.GetLastBusinessDayOfMonth() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func Test_date_GetDayOfWeekInMonth(t *testing.T) {

	r := ring.New(7)
	for i := 0; i < r.Len(); i++ {
		r.Value = i
		r = r.Next()
	}
	r.Next()
	d := NewDate(2019,1,1)
	r = r.Move(d.DayOfWeek())
	countDown := d.DayOfWeek()
	countUp := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 7; j++ {
			if countDown > 0 || countUp > 0 {
				if d.GetDayOfWeekInMonth( i,j) != nil {
					t.Errorf("counting up to first day failed")
				}
				countDown = countDown - 1
			} else {
				if d.GetDayOfWeekInMonth(i,j) == d.GetLastDayOfMonth() {
					countUp = countUp+1
				}
				if d.GetDayOfWeekInMonth(i,j).DayOfWeek() != r.Value {
					t.Errorf("date.GetDayOfWeekInMonth() = %v, want %v", d.DayOfWeek(), r.Value)
				}
				r = r.Next()
			}
		}
	}
}
