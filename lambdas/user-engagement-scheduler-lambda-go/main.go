package lambda

import (
	"context"
	"fmt"
	aesc "github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/crosswalks"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/schedules"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	es "github.com/adaptiveteam/adaptive/engagement-scheduling"
	esmodels "github.com/adaptiveteam/adaptive/engagement-scheduling-models"
	"time"
)

// type EngSchedule struct {
// 	Target string `json:"target"`
// 	Date   string `json:"date"`
// }

var (
	clientID = utils.NonEmptyEnv("CLIENT_ID")
	schema   = models.SchemaForClientID(clientID)
	userDao  = utilsUser.NewDAOFromSchema(d, namespace, schema)
)

// HandleRequest -
func HandleRequest(ctx context.Context, event models.UserEngage) (err error) {
	defer core.RecoverToErrorVar("user-engagement-scheduler-lambda-go", &err)
	var t time.Time
	if event.Date != "" {
		fmt.Printf("Date is present in EngSchedule")
		t, err = core.ISODateLayout.Parse(event.Date)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse %s as date", event.Date))
	} else {
		fmt.Printf("Date not present in EngSchedule, using current date")
		t = time.Now()
	}
	var y, m, d = t.Date()
	fmt.Println("### business time date: " + business_time.NewDate(y, int(m), d).DateToString(string(core.ISODateLayout)))
	// TODO: Take date from eng
	location, _ := time.LoadLocation("UTC")
	holidaysList := schedules.LoadHolidays(time.Date(y, m, d, 0, 0, 0, 0, location),
		userDao.ReadUnsafe(event.UserId).PlatformID,
		adHocHolidaysTable, adHocHolidaysPlatformDateIndex)
	allCrosswalks := func() []esmodels.CrossWalk {
		return concatAppend([][]esmodels.CrossWalk{crosswalks.UserCrosswalk()})
	}
	day := business_time.NewDate(y, int(m), d)
	es.ActivateEngagementsOnDay(
		aesc.ProductionProfile,
		day,
		allCrosswalks,
		holidaysList,
		location,
		event.UserId,
	)
	return
}

func concatAppend(slices [][]esmodels.CrossWalk) []esmodels.CrossWalk {
	var tmp []esmodels.CrossWalk
	for _, s := range slices {
		tmp = append(tmp, s...)
	}
	return tmp
}
