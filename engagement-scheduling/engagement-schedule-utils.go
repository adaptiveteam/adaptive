package engagement_scheduling

import (
	"log"
	bt "github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/checks"
	models "github.com/adaptiveteam/adaptive/engagement-scheduling-models"
	"sort"
	"sync"
	"time"
)

// ScheduleCheck enables us to chain together checks.  The message will be passed along
// in full from one check to the next if the indicated check is the same as the expected value.
// If the indicated check is not the same as the expected value then the Message will be zeroed out.
type ScheduleCheck struct {
	profile        checks.CheckResultMap
	Message        string
}

// checkSchedule checks the schedule for a specific day and a specific
// schedule function to determine if there is an engagement to activate.
func checkSchedule(
	checkResultMap checks.CheckResultMap,
	date bt.Date,
	cw models.CrossWalk,
	group *sync.WaitGroup,
	out chan models.CrossWalkName,
) {
	defer group.Done()
	engagement := cw.Schedule(checkResultMap,date)
	if len(engagement) > 0 {
		e := models.CrossWalkName {
			Name:engagement,
			Functions:cw,
		}
		out <- e
	}
}

func gatherDays(
	channel chan models.ScheduledEngagement,
	wg *sync.WaitGroup,
) (rv models.ScheduledEngagementList){
	wg.Wait()
	close(channel)
	rv = make(models.ScheduledEngagementList, 0)
	for d := range channel {
		rv = append(rv,d)
	}
	sort.Sort(rv)
	return rv
}

func runDay(
	checkFunctionMap checks.CheckFunctionMap,
	date bt.Date,
	endDate bt.Date,
	scheduledEngagements func() []models.CrossWalk,
	holidays bt.Holidays,
	location *time.Location,
	target string,
	channel chan models.ScheduledEngagement,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	defer func(){
		if err2 := recover(); err2 != nil {
			log.Printf("GenerateScheduleOfEngagements/runDay recovered for date=%v: %+v\n", date, err2)

		}
	}()
	// log.Printf("GenerateScheduleOfEngagements/runDay for date=%v\n", date)
	day, ok := constructDay(
		checkFunctionMap,
		date,
		endDate,
		scheduledEngagements,
		holidays,
		location,
		target,
	)
	// log.Printf("GenerateScheduleOfEngagements/runDay for date=%v ok=%v\n", date, ok)
	
	if ok {
		channel <- day
	}
}

func constructDay(
	checkFunctionMap checks.CheckFunctionMap,
	date bt.Date,
	endDate bt.Date,
	scheduledEngagements func() []models.CrossWalk,
	holidays bt.Holidays,
	location *time.Location,
	target string,
) (rv models.ScheduledEngagement, ok bool){
	// log.Printf("GenerateScheduleOfEngagements/constructDay 1 for date=%v \n", date)
	checkResultMap := checkFunctionMap.Evaluate(target, date) // this function never returns
	// log.Printf("GenerateScheduleOfEngagements/constructDay 2 for date=%v \n", date)
	engagementsOnDay := GetEngagementsOnDay(checkResultMap, date, scheduledEngagements)
	ok = false
	if len(engagementsOnDay) > 0 {
		// log.Printf("GenerateScheduleOfEngagements/constructDay 3 for date=%v \n", date)
		rescheduledDate := date.GetBusinessDay(
			holidays,
			location,
			true,
		)
		// log.Printf("GenerateScheduleOfEngagements/constructDay 3 for date=%v \n", date)
		dateBefore := rescheduledDate.DateBefore(endDate, true)
		if dateBefore {
			if rescheduledDate != date {
				rv = models.ScheduledEngagement{
						Engagements:     engagementsOnDay,
						ScheduledDate:   rescheduledDate,
						RescheduledFrom: date,
						RescheduledFor:  holidays.HolidaysOnDate(date, location),
					}
				ok = true
			} else {
				rv = models.ScheduledEngagement{
						Engagements:     engagementsOnDay,
						ScheduledDate:   date,
						RescheduledFrom: nil,
						RescheduledFor:  nil,
					}
				ok = true
			}
			// log.Printf("GenerateScheduleOfEngagements/constructDay 4 for date=%v \n", date)

		}
	}
	// log.Printf("GenerateScheduleOfEngagements/constructDay 5 for date=%v \n", date)

	return rv, ok
}

func monitorScheduleCheck(group *sync.WaitGroup,out chan models.CrossWalkName,) {
	group.Wait()
	close(out)
}

func removeDupedEngagenements(list models.CrossWalkNameList) (rv models.CrossWalkNameList) {
	rv = make(models.CrossWalkNameList,0)
	for _,v := range list {
		found := false
		for i := 0; i < len(rv) && !found; i++ {
			if rv[i].Name == v.Name {
				found = true
			}
		}
		if found == false {
			rv = append(rv,v)
		}
	}
	return rv
}