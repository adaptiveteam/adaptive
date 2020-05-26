package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-checks"
	"github.com/adaptiveteam/adaptive/daos/common"
	"context"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/crosswalks"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/schedules"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	es "github.com/adaptiveteam/adaptive/engagement-scheduling"
	"time"
)

var (
	clientID = utils.NonEmptyEnv("CLIENT_ID")
	schema   = models.SchemaForClientID(clientID)
)

// HandleRequest -
func HandleRequest(ctx context.Context, event models.UserEngage) (err error) {
	defer core.RecoverToErrorVar("user-engagement-scheduler-lambda-go", &err)
	var t time.Time
	if event.Date != "" {
		fmt.Printf("Date is present in UserEngage.Date=%s", event.Date)
		t, err = core.ISODateLayout.Parse(event.Date)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse %s as date", event.Date))
	} else {
		t = time.Now()
		fmt.Printf("Date not present in UserEngage, using  date of current time %v", t)
	}
	var y, m, d = t.Date()
	fmt.Println("### business time date: " + business_time.NewDate(y, int(m), d).DateToString(string(core.ISODateLayout)))
	// TODO: Take date from eng
	location, _ := time.LoadLocation("UTC")
	holidaysList := schedules.LoadHolidays(time.Date(y, m, d, 0, 0, 0, 0, location),
		event.TeamID,
		adHocHolidaysTable, adHocHolidaysPlatformDateIndex)
	// allCrosswalks := func() []esmodels.CrossWalk {
	// 	return concatAppend([][]esmodels.CrossWalk{crosswalks.UserCrosswalk()})
	// }
	day := business_time.NewDate(y, int(m), d)
	connGen := common.CreateConnectionGenFromEnv()
	conn := connGen.ForPlatformID(event.TeamID.ToPlatformID())
	es.ActivateEngagementsOnDay(
		adaptive_checks.EvalProfile,
		day,
		crosswalks.UserCrosswalk,
		holidaysList,
		location,
		event.TeamID, 
		event.UserID,
		conn,
	)
	return
}

// func concatAppend(slices [][]esmodels.CrossWalk) []esmodels.CrossWalk {
// 	var tmp []esmodels.CrossWalk
// 	for _, s := range slices {
// 		tmp = append(tmp, s...)
// 	}
// 	return tmp
// }
