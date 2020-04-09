// Package businessTime provides more advanced time functions that
// align with common business cycles
package business_time

import (
	"github.com/pkg/errors"
	"encoding/json"
	"math"
	"sort"
	"time"
)

/****************************************************************************/
/* The following functions are meant to handle basic date functions         */
/****************************************************************************/

// The Date interface creates a simple business date
type Date interface {
	GetYear() int
	GetMonth() int
	GetDay() int
	GetQuarter() int
	GetPreviousQuarter() int
	GetNextQuarter() int
	GetPreviousQuarterYear() int
	GetNextQuarterYear() int
	GetLastDayOfQuarter() Date
	GetFirstDayOfQuarter() Date
	GetLastWeekDayOfQuarter() Date
	GetFirstWeekDayOfQuarter() Date
	GetDayOfSundayWeek1InQuarter(sundayWeek int, day int) (date Date)
	GetDayOfWeekInQuarter(week int, day int) Date
	GetDayOfWeekInMonth(week int, day int) Date
	GetFirstDayOfMonth() Date
	GetLastDayOfMonth() Date
	GetFirstBusinessDayOfMonth(h Holidays, l *time.Location) Date
	GetLastBusinessDayOfMonth(h Holidays, l *time.Location) Date
	IsBusinessDay(h Holidays, l *time.Location) bool
	GetBusinessDay(h Holidays, l *time.Location, forward bool) Date
	PreviousBusinessDay(h Holidays, l *time.Location) Date
	NextBusinessDay(h Holidays, l *time.Location) Date
	GetLastBusinessDayOfQuarter(h Holidays, location *time.Location) Date
	GetFirstBusinessDayOfQuarter(h Holidays, location *time.Location) Date
	DayOfWeek() int
	IsWeekDay() bool
	GetWeekDay(forward bool) Date
	PreviousWeekDay() Date
	NextWeekDay() Date
	GetFirstDayOfWeek() Date
	GetLastDayOfWeek() Date
	DateBefore(d Date, inclusive bool) bool
	DateAfter(d Date, inclusive bool) bool
	AddTime(year int, month int, day int) Date
	DaysBetween(d2 Date) int
	DateToTimeMidnight() time.Time
	DateToString(format string) string
}

// Holidays interface maintains a list of holidays
type Holidays interface {
	HolidaysToJSON() []byte
	// AddHoliday adds a new holiday. 
	// Returns true. In case of duplicate name returns false.
	AddHoliday(name string, date Date, location time.Location) bool
	DeleteHoliday(name string) bool
	GetHoliday(name string) Date
	GetHolidays() []string
	// GetHolidayDate returns holiday date by holiday name
	GetHolidayDate(string) Date
	GetHolidayLocation(string) *time.Location
	HolidayRegistered(name string) bool
	HolidaysOnDate(date Date, location *time.Location) []string
	Len() int
}

// date is meant to be a simple date without the time element.
// The date type is meant to make working with dates of the form YYYY/MM/DD easier
// Note that this is not exported from the package
type date struct {
	Year  int
	Month int
	Day   int
}

const (
	Sunday = 0
	Monday = 1
	Tuesday = 2
	Wednesday = 3
	Thursday = 4
	Friday = 5
	Saturday = 6
)


/***************************/
/* Receivers for date type */
/***************************/

// Year returns the year of the date
func (d date) GetYear() int {
	return d.Year
}

// Month returns the month of the date
func (d date) GetMonth() int {
	return d.Month
}

// Day returns the day of the date
func (d date) GetDay() int {
	return d.Day
}

// Quarter returns the quarter of the date
func (d date) GetQuarter() int {
	return quarters[d.Month]
}

// Returns the current date inthe provided time zone
func (d date) GetDateInLocation(l *time.Location) Date {
	t := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, l).In(l)
	return date{
		Year:  t.Year(),
		Month: int(t.Month()),
		Day:   t.Day(),
	}
}

// DayOfWeek returns the day of the week for the provided date
func (d date) DayOfWeek() int {
	return int(d.DateToTimeMidnight().Weekday())
}

func (d date) IsWeekDay() bool {
	return d.DayOfWeek() >= 1 && d.DayOfWeek() <= 5
}

func (d date) GetWeekDay(forward bool) (rv Date) {
	rv = d
	if !d.IsWeekDay() {
		if forward {
			rv = d.NextWeekDay()
		} else {
			rv = d.PreviousWeekDay()
		}
	}
	return rv
}

func (d date) PreviousWeekDay() (rv Date) {
	if d.DayOfWeek() == 1 {
		rv = d.AddTime(0, 0, -3)
	} else if d.DayOfWeek() == 0 {
		rv = d.AddTime(0, 0, -2)
	} else {
		rv = d.AddTime(0, 0, -1)
	}
	return rv
}

func (d date) NextWeekDay() (rv Date) {
	if d.DayOfWeek() == 5 {
		rv = d.AddTime(0, 0, 3)
	} else if d.DayOfWeek() == 6 {
		rv = d.AddTime(0, 0, 2)
	} else {
		rv = d.AddTime(0, 0, 1)
	}
	return rv
}

// DateBefore returns true if d1 is before d2
func (d date) DateBefore(d2 Date, inclusive bool) (rv bool) {
	if inclusive {
		rv = d.DateToTimeMidnight() == d2.DateToTimeMidnight() || d.DateToTimeMidnight().Before(d2.DateToTimeMidnight())
	} else {
		rv = d.DateToTimeMidnight().Before(d2.DateToTimeMidnight())
	}
	return rv
}

// DateBefore returns true if d1 is after d2
func (d date) DateAfter(d2 Date, inclusive bool) (rv bool) {
	if inclusive {
		rv = d.DateToTimeMidnight() == d2.DateToTimeMidnight() || d.DateToTimeMidnight().After(d2.DateToTimeMidnight())
	} else {
		rv = d.DateToTimeMidnight().After(d2.DateToTimeMidnight())
	}
	return rv
}

func (d date) DaysBetween(d2 Date) int {
	dateOne := d.DateToTimeMidnight()
	dateTwo := d2.DateToTimeMidnight()
	diff := dateOne.Sub(dateTwo)
	return int(math.Abs(diff.Hours() / 24))
}

// AddTime adds the specified amount of time to the date
func (d date) AddTime(year int, month int, day int) Date {
	oldDate := d.DateToTimeMidnight()
	newDate := oldDate.AddDate(year, month, day)
	return TimeToDate(newDate)
}

// DateToTimeMidnight converts a date type to a Time type
func (d date) DateToTimeMidnight() time.Time {
	t := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
	return t
}

func (d date) DateToString(format string) string {
	ds := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC).Format(format)
	return ds
}

/****************************************************************************/
/* The following code is designed to facilitate creating dates.             */
/****************************************************************************/

// CreateSimpleDate creates a date with just the year, month, day, and timezone.
// It provides zero values for all other fields. If the zone is invalid the
// function will use the timezone of the current server (UTC in AWS).
func NewDate(year int, month int, day int) Date {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return date{
		Year:  t.Year(),
		Month: int(t.Month()),
		Day:   t.Day(),
	}
}

// NewDateFromQuarter is a simple helper function to create a new date from just
// quarter and year information.
func NewDateFromQuarter(quarter int, year int) Date {
	return FirstDateOfQuarter(quarter, year)
}

// FirstDateOfQuarter returns the first date of the quarter
func FirstDateOfQuarter(quarter int, year int) Date {
	if quarter < 1 || quarter > 4 {
		panic("FirstDateOfQuarter: quarter should be within [1,4]")
	}
	month1 := quarter * 3 - 2
	day := 1
	return NewDate(year, month1, day)
}

// LastDateOfQuarter returns the last date of the quarter
// It calculates the first date of the next quarter 
// and then subtracts 1 day.
func LastDateOfQuarter(quarter int, year int) Date {
	if quarter < 1 || quarter > 4 {
		panic("LastDateOfQuarter: quarter should be within [1,4]")
	}
	month1Plus3 := quarter * 3 + 1
	day := 1
	return NewDate(year, month1Plus3, day).AddTime(0,0,-1)
}
// Today creates a date with just the year, month, day or the current date in the given timezone.
// If the zone is invalid the function will use the timezone of the current server (UTC in AWS).
func Today(l *time.Location) Date {
	t := time.Now().In(l)
	return TimeToDate(t)
}

// TimeToDate converts a Time type to a date type
func TimeToDate(t time.Time) Date {
	return date{
		Year:  t.Year(),
		Month: int(t.Month()),
		Day:   t.Day(),
	}
}

// DateFromString converts strings in the form of YYYY-MM-DD-Location in to a Date
// If the location is wrong it will default to Local.
// NB: This function parses a different format - "2006-01-02T15:04:05Z07:00"
// Use DateFromISOWithLocationString instead
func DateFromString(ds string) Date {
	t, e := time.Parse(time.RFC3339, ds)
	if e == nil {
		return TimeToDate(t)
	}
	return nil
}

// DateFromYMDString converts strings in the form of YYYY-MM-DD in to a Date
func DateFromYMDString(ds string) (Date, error) {
	t, e := time.Parse("2006-01-02", ds)
	if e == nil {
		return TimeToDate(t), nil
	}
	return date{}, e
}

/****************************************************************************/
/* The following code is designed to facilitate the management of holidays. */
/****************************************************************************/

// holiday is a struct used internally for JSON creation and manipulation for holidays.
type holiday struct {
	What  string `json:"what"`
	When  date   `json:"when"`
	Where string `json:"where"`
}

// holidays is a mapping of dates to holidays to help calcuate business days.
type holidays map[string]holiday

// holidayArray is a struct used internally for JSON creation and manipulation for holidays.
type holidayArray []holiday

func NewHolidayList() Holidays {
	return make(holidays)
}

func (h holidays) HolidaysToJSON() []byte {
	var ha holidayArray

	for k, v := range h {
		ha = append(ha, holiday{
			k,
			date{
				Year:  v.When.GetYear(),
				Month: v.When.GetMonth(),
				Day:   v.When.GetDay(),
			},
			v.Where,
		})
	}

	sort.Sort(ha)
	js, _ := json.Marshal(ha)
	return js
}

// Len is required to make a sortable holiday array to support JSON
func (h holidayArray) Len() int {
	return len(h)
}

// Swap is required to make a sortable holiday array to support JSON
func (h holidayArray) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Less is required to make a sortable holiday array to support JSON
func (h holidayArray) Less(i, j int) bool {
	return h[i].When.DateToTimeMidnight().Unix() < h[j].When.DateToTimeMidnight().Unix()
}

// HolidaysFromJSON reads a JSON file and loads it in to a Holidays map
func HolidaysFromJSON(h Holidays, js []byte) (err error) {
	var ha holidayArray
	var location *time.Location
	err = json.Unmarshal(js, &ha)
	for i := 0; i < len(ha) && err == nil; i++ {
		location, err = time.LoadLocation(ha[i].Where)
		h.AddHoliday(ha[i].What, ha[i].When, *location)
	}
	return err
}

func (h holidays) GetHolidayLocation(name string) (rv *time.Location) {
	rv, _ = time.LoadLocation(h[name].Where)
	return rv
}

// AddHoliday enables a developer to add a holiday to a map of holidays.
// Note that you can add multiple holidays on the same date.
// The function will return an error if the holiday already exists, even on a different date.
func (h holidays) AddHoliday(name string, d Date, location time.Location) bool {
	if h.HolidayRegistered(name) {
		return false
	} else {
		h[name] = holiday{
			What: name,
			When: date{
				d.GetYear(),
				d.GetMonth(),
				d.GetDay(),
			},
			Where: location.String(),
		}
		return true
	}
}

// DeleteHoliday enables a developer to delete a holiday.
// The function returns an error if the holiday doesn't exist
func (h holidays) DeleteHoliday(name string) bool {
	if h.HolidayRegistered(name) {
		delete(h, name)
		return true
	} else {
		return false
	}
}

// GetHolidays returns all of the holidays
func (h holidays) GetHolidays() []string {
	ha := make([]string, 0)
	for k := range h {
		ha = append(ha, k)
	}
	return ha
}

// GetHolidayDate return the date for given holiday.
func (h holidays) GetHolidayDate(hn string) Date {
	if h.HolidayRegistered(hn) {
		return h[hn].When
	}
	return nil
}

// HolidayDateRegistered checks for the existence of a registered holiday
func (h holidays) HolidayRegistered(name string) bool {
	_, ok := h[name]
	return ok
}

// HolidayDateRegistered checks for the existence of a registered holiday
func (h holidays) GetHoliday(name string) (rv Date) {
	rv = nil
	if h.HolidayRegistered(name) {
		rv = h[name].When
	}
	return rv
}

func (h holidays) HolidaysOnDate(d Date, location *time.Location) []string {
	holidays := make([]string, 0)
	for k, v := range h {
		if v.When.DateToTimeMidnight().In(location) == d.DateToTimeMidnight().In(location) {
			holidays = append(holidays, k)
		}
	}
	return holidays
}

func (h holidays) Len() int {
	return len(h)
}

/*****************************************************************************/
/* The following functions are meant to are meant to work with business days */
/*****************************************************************************/

// FirstDateOfQuarter returns the first date of the quarter
func (d date)FirstDateOfQuarter() Date {
	return FirstDateOfQuarter(d.GetQuarter(), d.GetYear())
}

// LastDateOfQuarter returns the last date of the quarter
func (d date)LastDateOfQuarter() Date {
	return LastDateOfQuarter(d.GetQuarter(), d.GetYear())
}

// IsBusinessDay returns whether or not the given day is a designated business day
func (d date) IsBusinessDay(h Holidays, l *time.Location) bool {
	day := d.DayOfWeek()
	holidays := h.HolidaysOnDate(d, l)
	if day == 6 || day == 0 || len(holidays) > 0 {
		return false
	} else {
		return true
	}
}

func (d date) GetBusinessDay(h Holidays, l *time.Location, forward bool) (rv Date) {
	rv = d
	if !d.IsBusinessDay(h, l) {
		if forward {
			rv = d.NextBusinessDay(h, l)
		} else {
			rv = d.PreviousBusinessDay(h, l)
		}
	}
	return rv
}

// PreviousBusinessDay returns the first business day before the given time
func (d date) PreviousBusinessDay(h Holidays, l *time.Location) Date {
	t := d.DateToTimeMidnight()
	for i := t.AddDate(0, 0, -1); !TimeToDate(i).IsBusinessDay(h, l); i = i.AddDate(0, 0, -1) {
		t = i
	}
	newDate := TimeToDate(t.AddDate(0, 0, -1))
	return newDate
}

// NextBusinessDay returns the next business day before the given time
func (d date) NextBusinessDay(h Holidays, l *time.Location) Date {
	t := d.DateToTimeMidnight()
	for i := t.AddDate(0, 0, +1); !TimeToDate(i).IsBusinessDay(h, l); i = i.AddDate(0, 0, +1) {
		t = i
	}
	newDate := TimeToDate(t.AddDate(0, 0, +1))
	return newDate
}

/*****************************************************************************/
/* The following functions are meant to are meant to work business quarters  */
/*****************************************************************************/

// quarters is an unexported mapping from months to business quarters
var quarters = map[int]int{
	1:  1,
	2:  1,
	3:  1,
	4:  2,
	5:  2,
	6:  2,
	7:  3,
	8:  3,
	9:  3,
	10: 4,
	11: 4,
	12: 4,
}

// GetPreviousQuarter returns the quarter before the quarter of the given time
func (d date) GetPreviousQuarter() int {
	q := d.GetQuarter()
	if q == 1 {
		return 4
	} else {
		return q - 1
	}
}

// GetNextQuarter returns the quarter before the quarter of the given time
func (d date) GetNextQuarter() int {
	q := d.GetQuarter()
	if q == 4 {
		return 1
	} else {
		return q + 1
	}
}

// GetPreviousQuarterYear returns the year of the previous quarter
func (d date) GetPreviousQuarterYear() int {
	quarter, _ := quarters[d.GetMonth()]
	if quarter == 1 {
		return d.GetYear() - 1
	} else {
		return d.GetYear()
	}
}

// GetNextQuarterYear get the the year of the next quarter
func (d date) GetNextQuarterYear() int {
	quarter, _ := quarters[d.GetMonth()]
	if quarter == 4 {
		return d.GetYear() + 1
	} else {
		return d.GetYear()
	}
}

// GetLastDayOfQuarter returns the last day of the last month of the quarter.
func (d date) GetLastDayOfQuarter() Date {
	quarter := d.GetQuarter()
	month := quarter*3 + 1
	lastDayOfQuarter := NewDate(d.GetYear(), month, 1).AddTime(0, 0, -1)
	return lastDayOfQuarter
}

// GetFirstDayOfQuarter returns the first day of the first month of the quarter.
func (d date) GetFirstDayOfQuarter() Date {
	quarter := d.GetQuarter()
	month := quarter*3 - 2
	firstDayOfQuarter := NewDate(d.GetYear(), month, 1)
	// if the last day is a Saturday or Sunday then forward to Monday
	if firstDayOfQuarter.DayOfWeek() == 6 {
		firstDayOfQuarter = firstDayOfQuarter.AddTime(0, 0, 2)
	} else if firstDayOfQuarter.DayOfWeek() == 0 {
		firstDayOfQuarter = firstDayOfQuarter.AddTime(0, 0, 1)
	}
	return firstDayOfQuarter
}

func (d date) GetLastWeekDayOfQuarter() Date {
	quarter := d.GetQuarter()
	month := quarter*3 + 1
	lastDayOfQuarter := NewDate(d.GetYear(), month, 1).AddTime(0, 0, -1)
	// if the last day is a Saturday or Sunday then back it up to Friday
	if lastDayOfQuarter.DayOfWeek() == 6 {
		lastDayOfQuarter = lastDayOfQuarter.AddTime(0, 0, -1)
	} else if lastDayOfQuarter.DayOfWeek() == 0 {
		lastDayOfQuarter = lastDayOfQuarter.AddTime(0, 0, -2)
	}
	return lastDayOfQuarter
}

func (d date) GetFirstWeekDayOfQuarter() Date {
	quarter := d.GetQuarter()
	month := quarter*3 - 2
	lastDayOfQuarter := NewDate(d.GetYear(), month, 1)
	return lastDayOfQuarter
}

// GetLastBusinessDayOfQuarter returns the last business day of the quarter.
func (d date) GetLastBusinessDayOfQuarter(h Holidays, location *time.Location) Date {
	lastDayOfQuarter := d.GetLastDayOfQuarter()
	if lastDayOfQuarter.IsBusinessDay(h, location) {
		return lastDayOfQuarter
	} else {
		return lastDayOfQuarter.PreviousBusinessDay(h, location)
	}
}

// GetLastBusinessDayOfQuarter returns the last business day of the quarter.
func (d date) GetFirstBusinessDayOfQuarter(h Holidays, location *time.Location) Date {
	firstDayOfQuarter := d.GetFirstDayOfQuarter()
	if firstDayOfQuarter.IsBusinessDay(h, location) {
		return firstDayOfQuarter
	} else {
		return firstDayOfQuarter.NextBusinessDay(h, location)
	}
}

// GetFirstBusinessDayOfWeekInQuarter returns the first business day of the given
// week in the quarter.  If week is positive it starts from the beginning of
// the quarter.  If the week is negative it starts from the end of the quarter.
// The parameter day should be between 0 and 6, where 0 is Sunday.
// The parameter forward is whether you want to adjust forward or back.
// Deprecated: the first week is 0
func (d date) GetDayOfWeekInQuarter(week int, day int) (rv Date) {
	var dayInWeek Date
	quarter := NewDateFromQuarter(d.GetQuarter(), d.GetYear())
	if week >= 0 {
		dayInWeek = quarter.GetFirstDayOfQuarter().AddTime(0, 0, 7*week)
	} else if week < 0 {
		dayInWeek = quarter.GetLastDayOfQuarter().AddTime(0, 0, 7*week)
	}
	dayOfWeek := dayInWeek.DayOfWeek()
	rv = dayInWeek.AddTime(0, 0, day-dayOfWeek)
	if rv.GetQuarter() < d.GetQuarter() || rv.GetQuarter() > d.GetQuarter() {
		rv = nil
	}
	return rv
}
// GetSundayDateInQuarter1 - returns date of sunday.
// if sunday index is > 0 then it's started from the first sunday (1) of the quarter
// if sunday index is < 0 then it's counted from the last sunday (-1) of the quarter
// if sunday index == 0 then it panics.
func (d date) GetSundayDateInQuarter1(sundayIndex1 int) (date Date) {
	if sundayIndex1 == 0 {
		panic(errors.New("GetSundayDateInQuarter1 does not allow 0"))
	} else if sundayIndex1 > 0 {
		firstDay := d.FirstDateOfQuarter()
		firstDayDayOfWeek := firstDay.DayOfWeek()
		deltaToSunday := 0 
		if firstDayDayOfWeek > 0 {
			deltaToSunday = 7 - firstDayDayOfWeek
		}
		date = firstDay.AddTime(0, 0, deltaToSunday + 7 * (sundayIndex1 - 1))
	} else {
		lastDay := d.LastDateOfQuarter()
		lastDayDayOfWeek := lastDay.DayOfWeek()
		deltaToSunday := -lastDayDayOfWeek 
		date = lastDay.AddTime(0, 0, deltaToSunday + 7 * (sundayIndex1 + 1))
	}
	return
}

// GetDayOfWeek1InQuarter - returns date of the given day in the
// n-th week that includes Sunday. The first week has number 1 (there might be a few days before sunday though).
// One may use negative week indices to get week counted from the end 
// of quarter. In this case the week with index -1 might be incomplete.
func (d date) GetDayOfSundayWeek1InQuarter(sundayWeek int, day int) (date Date) {
	date = d.GetSundayDateInQuarter1(sundayWeek).AddTime(0, 0, day)
	return
}
// GetFirstBusinessDayOfWeekInMonth returns the date of the desired business day in
// the month of the given day.
func (d date) GetDayOfWeekInMonth(week int, day int) (rv Date) {
	var dayInWeek Date
	month := NewDate(d.GetYear(), d.GetMonth(),1)
	if week >= 0 {
		dayInWeek = month.GetFirstDayOfMonth().AddTime(0, 0, 7*week)
	} else if week < 0 {
		dayInWeek = month.GetLastDayOfMonth().AddTime(0, 0, 7*week)
	}
	dayOfWeek := dayInWeek.DayOfWeek()
	rv = dayInWeek.AddTime(0, 0, day-dayOfWeek)
	if rv.GetMonth() < d.GetMonth() || rv.GetMonth() > d.GetMonth() {
		rv = nil
	}
	return rv
}

// GetFirstDayOfWeek returns the first day(to midnight) of the current week
// Week is assumed to start with Sunday and end with Saturday.
func (d date) GetFirstDayOfWeek() Date {
	weekday := d.DayOfWeek()
	return d.AddTime(0, 0, -weekday)
}

// GetLastDayOfWeek returns the last day(to midnight) of the current week
func (d date) GetLastDayOfWeek() Date {
	return d.GetFirstDayOfWeek().AddTime(0, 0, 6)
}

func (d date) GetFirstDayOfMonth() (rv Date) {
	rv = NewDate(d.Year,d.Month,1)
	return rv
}

func (d date) GetLastDayOfMonth() (rv Date) {
	rv = d.GetFirstDayOfMonth().AddTime(0,1,0).GetFirstDayOfMonth().AddTime(0,0,-1)
	return  rv
}

func (d date) GetFirstBusinessDayOfMonth(h Holidays, l *time.Location) (rv Date) {
	rv = d.GetFirstDayOfMonth().GetBusinessDay(h,l,true)
	return  rv
}

func (d date) GetLastBusinessDayOfMonth(h Holidays, l *time.Location) (rv Date) {
	rv = d.GetLastDayOfMonth().GetBusinessDay(h,l,false)
	return  rv
}