package tests

import (
	"github.com/adaptiveteam/adaptive/engagement-scheduling/test_engagements"
	"github.com/adaptiveteam/adaptive/checks"
	es "github.com/adaptiveteam/adaptive/engagement-scheduling"
	"github.com/adaptiveteam/adaptive/engagement-scheduling/test_checks"
	"github.com/adaptiveteam/adaptive/engagement-scheduling/test_crosswalks"
	"reflect"
	"sort"
	"testing"
	"time"

	bt "github.com/adaptiveteam/adaptive/business-time"
	models "github.com/adaptiveteam/adaptive/engagement-scheduling-models"
)

func getTestHolidays() bt.Holidays {
	holidays := bt.NewHolidayList()
	newYears := bt.NewDate(2019, 12, 31)
	christmas := bt.NewDate(2019, 12, 25)
	thanksgiving := bt.NewDate(2019, 11, 28)
	diwali := bt.NewDate(2019, 10, 28)
	newHalloween := bt.NewDate(2019, 10, 14)
	holidays.AddHoliday("Christmas", christmas, *time.UTC)
	holidays.AddHoliday("New Years", newYears, *time.UTC)
	holidays.AddHoliday("Thanksgiving", thanksgiving, *time.UTC)
	holidays.AddHoliday("Diwali", diwali, *time.UTC)
	holidays.AddHoliday("New Halloween", newHalloween, *time.UTC)
	return holidays
}

func TestCreateQuarterlySchedule(t *testing.T) {
	target := "ctcreel"
	holidays := getTestHolidays()
	quarterStart := bt.NewDateFromQuarter(3, 2019)

	quarterEnd := quarterStart.GetLastDayOfQuarter()

	allUsersSchedule := es.GenerateScheduleOfEngagements(
		test_checks.AllTrueTestProfile,
		quarterStart,
		target,
		test_crosswalks.UserCrosswalk,
		holidays,
		time.UTC,
		quarterEnd.DaysBetween(quarterStart),
	)

	test_engagements.Println("Schedule of Engagements for all users.")
	schedule := es.PrettyPrintSchedule(allUsersSchedule)
	for _,s := range schedule {
		test_engagements.Println(s)
	}
}

func TestActivateEngagements(t *testing.T) {
	holidays := getTestHolidays()
	target := "ctcreel"

	type args struct {
		date       bt.Date
		holidays   bt.Holidays
		target     string
		targetType string
	}
	tests := []struct {
		name   string
		args   args
		wantRv []string
	}{
		{
			name: "basic",
			args: args{
				date:       bt.NewDateFromQuarter(4, 2019).GetDayOfWeekInQuarter(1, bt.Monday),
				holidays:   holidays,
				target:     target,
			},
			wantRv: []string{
				"Reminder to update Individual Development Objectives",
				"Reminder to update Initiatives",
				"Reminder to update Objectives",
			},
		},
		{
			name: "basic holiday",
			args: args{
				date:       bt.NewDateFromQuarter(4, 2019).GetDayOfWeekInQuarter(2, bt.Monday),
				holidays:   holidays,
				target:     target,
			},
			wantRv: []string{},
		},
		{
			name: "holiday",
			args: args{
				date:       bt.NewDate(2019, 12, 25),
				holidays:   holidays,
				target:     target,
			},
			wantRv: []string{},
		},
		{
			name: "Testing make-up test_engagements",
			args: args{
				date:       bt.NewDate(2019, 10, 15),
				holidays:   holidays,
				target:     target,
				targetType: "user",
			},
			wantRv: []string{
				"Reminder to update Individual Development Objectives",
				"Reminder to update Initiatives",
				"Reminder to update Objectives",
			},
		},
	}

	daysInQuarter := bt.NewDateFromQuarter(1, 2019).GetFirstDayOfQuarter().DaysBetween(bt.NewDateFromQuarter(4, 2019).GetLastDayOfQuarter())
	date := bt.NewDateFromQuarter(1, 2019)
	for i := 0; i < daysInQuarter; i++ {
		if date.IsBusinessDay(holidays,time.UTC) {
			engagements := es.GenerateScheduleOfEngagements(
				test_checks.AllTrueTestProfile,
				date,
				target,
				test_crosswalks.UserCrosswalk,
				holidays,
				time.UTC,
				0,
			)
			if len(engagements) >= 1 {
				tests = append(
					tests,
					struct {
						name   string
						args   args
						wantRv []string
					}{
						name: date.DateToString("2006-01-02"),
						args: args{
							date:       date,
							holidays:   holidays,
							target:     target,
							targetType: "user",
						},
						wantRv: es.GetEngagementNames(
							engagements[0].Engagements,
						),
					},
				)
			} else {
				t.Errorf(
					"Expected one engagement but got %v for date %v",
					len(engagements),
					date.DateToString("2006-01-02"),
				)
			}
		}
		date = date.AddTime(0, 0, 1)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRv := es.ActivateEngagementsOnDay(
				test_checks.AllTrueTestProfile,
				tt.args.date,
				test_crosswalks.UserCrosswalk,
				tt.args.holidays,
				time.UTC,
				tt.args.target,
			)
			sort.Strings(gotRv)
			sort.Strings(tt.wantRv)
			if !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("ActivateEngagements() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func TestGenerateScheduleOfEngagements(t *testing.T) {
	holidays := getTestHolidays()
	target := "ctcreel"

	type args struct {
		checkFunctionMap       checks.CheckFunctionMap
		date                 bt.Date
		target               string
		scheduledEngagements func() []models.CrossWalk
		holidays             bt.Holidays
		location             *time.Location
		daysOut              int
	}
	tests := []struct {
		name         string
		args         args
		wantSchedule models.ScheduledEngagementList
	}{
		{
			name:"on holiday",
			args: args{
				checkFunctionMap:test_checks.AllTrueTestProfile,
				date:bt.NewDate(2019, 10, 14),
				target:target,
				scheduledEngagements:test_crosswalks.UserCrosswalk,
				holidays:holidays,
				location:time.UTC,
				daysOut:0,
			},
			wantSchedule: models.ScheduledEngagementList{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSchedule := es.GenerateScheduleOfEngagements(
				tt.args.checkFunctionMap,
				tt.args.date,
				tt.args.target,
				tt.args.scheduledEngagements,
				tt.args.holidays,
				tt.args.location,
				tt.args.daysOut,
			); !reflect.DeepEqual(gotSchedule, tt.wantSchedule) {
				t.Errorf("GenerateScheduleOfEngagements() = %v, want %v", gotSchedule, tt.wantSchedule)
			}
		})
	}
}
