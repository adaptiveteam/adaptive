package engagement_scheduling

import (
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/adaptive-checks"
	"fmt"
	"sort"
	"sync"
	"time"

	bt "github.com/adaptiveteam/adaptive/business-time"
	models "github.com/adaptiveteam/adaptive/engagement-scheduling-models"

	core_utils_go "github.com/adaptiveteam/adaptive/core-utils-go"
)

// ActivateEngagementsOnDay will run through the provided crosswalk to determine
// if there are any test_engagements to activate on the provided day.
// The function returns a list of all of the engagement descriptions for which
// test_engagements were generated.
func ActivateEngagementsOnDay(
	typedProfileConstructor adaptive_checks.TypedProfileConstructor,
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
	conn common.DynamoDBConnection,
) (rv []string) {
	// Only activate engagements if the date is a business day
	wg := &sync.WaitGroup{}
	var allEngagements models.ScheduledEngagementList

	allEngagements = GenerateScheduleOfEngagements(
		typedProfileConstructor,
		date,
		target,
		scheduledEngagements,
		holidays,
		location,
		0,
		conn,
	)
	// log.Print("All engagements: ", allEngagements)

	if len(allEngagements) > 0 {
		// All we really care about is the last day because that is the business day
		businessDay := allEngagements[len(allEngagements)-1]
		for _, e := range businessDay.Engagements {
			wg.Add(1)
			e2 := e // this creates a new variable to be used in closure. e itself is reused in iterations!
			core_utils_go.Go("Engagement: "+e2.Name, func() {
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
	typedProfileConstructor adaptive_checks.TypedProfileConstructor,
	date bt.Date,
	target string,
	scheduledEngagements func() []models.CrossWalk,
	holidays bt.Holidays,
	location *time.Location,
	daysOut int,
	conn common.DynamoDBConnection,
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
		date2 := date
		core_utils_go.Go(fmt.Sprintf("%d: runDay(date=%v)", i, date), func() {
			runDay(
				typedProfileConstructor,
				date2,
				endDate,
				scheduledEngagements,
				holidays,
				location,
				target,
				constructedDays,
				wg,
				conn,
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
	checkResultMap adaptive_checks.TypedProfile,
	date bt.Date,
	scheduledEngagements func() []models.CrossWalk,
) (rv models.CrossWalkNameList) {
	engagements := scheduledEngagements()
	rv = make(models.CrossWalkNameList, 0)
	engagementChannel := make(chan models.CrossWalkName)
	wg := &sync.WaitGroup{}
	for _, s := range engagements {
		wg.Add(1)
		s2 := s
		core_utils_go.Go("checkSchedule", func() { checkSchedule(checkResultMap, date, s2, wg, engagementChannel) })
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

// ScheduleEntry is the beginning of the AddScheduleFunctionCheck & AddScheduleBooleanCheck chain
// it defaults to true and then passes the ScheduleCheck through the pipeline of checks.
func ScheduleEntry(message string, flag bool) (rv string) {
	if flag {
		rv = message
	}
	return
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
