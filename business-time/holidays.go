package business_time

import (
	"time"
	"encoding/json"
	"sort"
)

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
