package engagement_scheduling_models

import (
	"github.com/adaptiveteam/adaptive/adaptive-checks"
	bt "github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/cron"
)

type ScheduleFunction func(resultMap adaptive_checks.TypedProfile, date bt.Date) (rv string)
type EngagementFunction func(date bt.Date, target string)
type ScheduledEngagementList []ScheduledEngagement
type CrossWalkNameList []CrossWalkName

type CrossWalk struct {
	Schedule   ScheduleFunction
	Engagement EngagementFunction
}

type CrossWalkName struct {
	Name      string
	Functions CrossWalk
}

func NewCrossWalk(schedule ScheduleFunction, engagement EngagementFunction) CrossWalk {
	return CrossWalk{Schedule: schedule, Engagement: engagement}
}

// CrontabLine constructs a CrossWalk from cron-like schedule
func CrontabLine(schedule cron.Schedule, engagementName string, engagement EngagementFunction) CrossWalk {
	return CrossWalk{
		Schedule: func(resultMap adaptive_checks.TypedProfile, date bt.Date) (rv string){
			if schedule.IsOnSchedule(date.DateToTimeMidnight()) {
				rv = engagementName
			}
			return
		}, 
		Engagement: engagement,
	}
}

// ScheduledEngagement is the the structure used to capture the egagements on a given day
type ScheduledEngagement struct {
	// This is the list of all engagements on a given day
	Engagements CrossWalkNameList
	// This is the date on which the engagements are scheduled
	// Note that this includes holidays which could have caused the
	// engagement to be rescheduled from the original day
	ScheduledDate bt.Date
	// This is the date, if any, the engagements were rescheduled from
	// in the event the original date (scheduledDate) fell on a holiday
	RescheduledFrom bt.Date
	// This is the holiday, if any, for which the engagement was rescheduled.
	RescheduledFor []string
}

func (a ScheduledEngagementList) Len() int {
	return len(a)
}

func (a ScheduledEngagementList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ScheduledEngagementList) Less(i, j int) bool {
	return a[i].ScheduledDate.DateToTimeMidnight().Unix() < a[j].ScheduledDate.DateToTimeMidnight().Unix()
}
