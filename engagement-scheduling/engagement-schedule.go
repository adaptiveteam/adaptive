package engagement_scheduling

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	bt "github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/checks"
	models "github.com/adaptiveteam/adaptive/engagement-scheduling-models"

	core_utils_go "github.com/adaptiveteam/adaptive/core-utils-go"
)

// ActivateEngagementsOnDay will run through the provided crosswalk to determine
// if there are any test_engagements to activate on the provided day.
// The function returns a list of all of the engagement descriptions for which
// test_engagements were generated.
func ActivateEngagementsOnDay(
	//  The map of function check results
	checkFunctionMap checks.CheckFunctionMap,
	// The day to check
	date bt.Date,
	// The specific crosswalk funtion to use
	scheduledEngagements func() []models.CrossWalk,
	// The holidays currently specified by the customer
	holidays bt.Holidays,
	// The location of the user
	location *time.Location,
	// Ths specific user or channel name to send the test_engagements to
	target string,
) (rv []string) {
	// Only activate engagements if the date is a business day
	var wg sync.WaitGroup
	var allEngagements models.ScheduledEngagementList

	allEngagements = GenerateScheduleOfEngagements(
		checkFunctionMap,
		date,
		target,
		scheduledEngagements,
		holidays,
		location,
		0,
	)
	// log.Print("All engagements: ", allEngagements)

	if len(allEngagements) > 0 {
		// All we really care about is the last day because that is the business day
		businessDay := allEngagements[len(allEngagements)-1]
		wg.Add(len(businessDay.Engagements))
		for _, e := range businessDay.Engagements {
			e2 := e // this creates a new variable to be used in closure. e itself is reused in iterations!
			core_utils_go.Go("Engagement: "+e.Name, func() {
				defer wg.Done()
				e2.Functions.Engagement(date, target)
			})
		}

		wg.Wait()

		rv = GetEngagementNames(businessDay.Engagements)
	} else {
		rv = make([]string, 0)
	}
	// log.Print("Activated engagements for day ", date.DateToString(string(core_utils_go.ISODateLayout)), rv)
	return rv
}

// GenerateScheduleOfEngagements will generate a list of all of the test_engagements
// that will be activated from the start date to the number of days out.
func GenerateScheduleOfEngagements(
	checkFunctionMap checks.CheckFunctionMap,
	date bt.Date,
	target string,
	scheduledEngagements func() []models.CrossWalk,
	holidays bt.Holidays,
	location *time.Location,
	daysOut int,
) (schedule models.ScheduledEngagementList) {
	schedule = make(models.ScheduledEngagementList, 0)

	// If the day before date  wasn't a business day then we need
	// to look for makeup engagements during the previous holiday.
	// To do this we need to adjust the date and the days out to the start of the
	// stretch of business days.
	if date.PreviousBusinessDay(holidays, location) != date.AddTime(0, 0, -1) {
		// plus/minus one is to avoid double counting the last business day
		daysOut = daysOut + date.PreviousBusinessDay(holidays, location).DaysBetween(date) - 1
		date = date.PreviousBusinessDay(holidays, location).AddTime(0, 0, 1)
	}

	endDate := date.AddTime(0, 0, daysOut)
	constructedDays := make(chan models.ScheduledEngagement, daysOut+1)
	wg := &sync.WaitGroup{}
	for i := 0; i <= daysOut; i++ {
		wg.Add(1)
		core_utils_go.Go(fmt.Sprintf("%d: runDay(date=%v)", i, date), func() {
			runDay(
				checkFunctionMap,
				date,
				endDate,
				scheduledEngagements,
				holidays,
				location,
				target,
				constructedDays,
				wg,
			)
		})
		date = date.AddTime(0, 0, 1)
	}

	expandedSchedule := gatherDays(
		constructedDays,
		wg,
	)

	// merge and de-dup engagements on each day
	for _, expandedDay := range expandedSchedule {
		found := false
		for i := 0; i < len(schedule); i++ {
			// Iterate through the consolidated schedule to see
			// if this day already exists and must be consolidated
			if expandedDay.ScheduledDate == schedule[i].ScheduledDate {
				// If we find the day in the final schedule, add the events from
				// the expanded schedule into just one day.
				// first, collect up all of the engagements for the given day
				for _, s := range expandedDay.Engagements {
					schedule[i].Engagements = append(schedule[i].Engagements, s)
				}

				// Now remove any duplicates
				schedule[i].Engagements = removeDupedEngagenements(schedule[i].Engagements)
				found = true
			}
		}
		if !found {
			schedule = append(schedule, expandedDay)
		}
	}

	return schedule
}

// GetEngagementsOnDay just returns a list of the engagement descriptions
// on a given day.
func GetEngagementsOnDay(
	checkResultMap checks.CheckResultMap,
	date bt.Date,
	scheduledEngagements func() []models.CrossWalk,
) (rv models.CrossWalkNameList) {
	engagements := scheduledEngagements()
	rv = make(models.CrossWalkNameList, 0)
	engagementChannel := make(chan models.CrossWalkName)
	wg := &sync.WaitGroup{}
	wg.Add(len(engagements))
	for _, s := range engagements {
		core_utils_go.Go("checkSchedule", func() { checkSchedule(checkResultMap, date, s, wg, engagementChannel) })
	}

	core_utils_go.Go("monitorScheduleCheck", func() { monitorScheduleCheck(wg, engagementChannel) })

	for e := range engagementChannel {
		rv = append(rv, e)
	}
	return rv
}

// GetEngagements generates a sort list of the engagement names
func GetEngagementNames(cw models.CrossWalkNameList) (rv []string) {
	rv = make([]string, 0)
	for _, en := range cw {
		rv = append(rv, en.Name)
	}
	sort.Strings(rv)
	return rv
}

// AddScheduleFunctionCheck enables the developer to chain these checks together to more closely draft
// the Adaptive Dynamic Menu.  The function checks the function named in the ScheduleCheck struct and
// if it it the same as the expected value, passes the message through to the next step.  Otherwise it
// will zero out the message.
func (s ScheduleCheck) AddScheduleFunctionCheck(functionCheck string, expectedValue bool) (rv ScheduleCheck) {
	rv.profile = s.profile

	if len(s.Message) > 0 {
		checkResult, err := s.profile.CheckResult(functionCheck, expectedValue)
		if err == nil {
			if checkResult {
				rv.Message = s.Message
			}
		} else {
			log.Printf("check function in %v is not in list", functionCheck)
		}
	}
	return rv
}

// AddScheduleBooleanCheck enables the developer to chain these checks together to more closely draft
// the Adaptive Dynamic Menu.  The function checks the expectedValue boolean and if it the same as
// the expected value, passes the message through to the next step.  Otherwise it will zero out the message.
func (s ScheduleCheck) AddScheduleBooleanCheck(booleanCheck, expectedValue bool) (rv ScheduleCheck) {
	rv.profile = s.profile

	if len(s.Message) > 0 && expectedValue == booleanCheck {
		rv.Message = s.Message
	} else {
		rv.Message = ""
	}
	return rv
}

// ScheduleEntry is the beginning of the AddScheduleFunctionCheck & AddScheduleBooleanCheck chain
// it defaults to true and then passes the ScheduleCheck through the pipeline of checks.
func ScheduleEntry(fc checks.CheckResultMap, message string) (rv ScheduleCheck) {
	rv.profile = fc
	rv.Message = message

	return rv
}

// Produces a "pretty print" version of the schedule
func PrettyPrintSchedule(
	allUsersSchedule models.ScheduledEngagementList,
) (rv []string) {
	rv = make([]string, 0)
	format := "Monday, January _2"
	for i, s := range allUsersSchedule {
		day := ""
		if s.RescheduledFrom == nil {
			day = day + fmt.Sprintf(
				"*%v) %s*",
				i+1,
				s.ScheduledDate.DateToString(format),
			)
		} else {
			day = day + fmt.Sprintf(
				"*%v) %s* _(Engagements rescheduled from %s)_",
				i+1,
				s.ScheduledDate.DateToString(format),
				s.RescheduledFrom.DateToString(format),
			)
		}

		if len(s.Engagements) > 0 {
			day = day + "\n"
			for i, e := range GetEngagementNames(s.Engagements) {
				day = day + fmt.Sprintf("\t:black_small_square: %s", e)
				if i < len(s.Engagements)-1 {
					day = day + "\n"
				}
			}
		}
		rv = append(rv, day)
	}
	return rv
}
