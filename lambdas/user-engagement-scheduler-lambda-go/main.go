package lambda

import (
	"github.com/sirupsen/logrus"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
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

type EngSchedule struct {
	Target string `json:"target"`
	Date   string `json:"date"`
}

var (
	clientID = utils.NonEmptyEnv("CLIENT_ID")
	schema   = models.SchemaForClientID(clientID)
	userDao  = utilsUser.NewDAOFromSchema(d, namespace, schema)
)

func handleError(err *error){
	if *err != nil {
		alog.LambdaLogger(logrus.InfoLevel).WithError(*err).Errorf("Ordinary error in user-engagement-scheduler-lambda")
	}
	
	if err2 := recover(); err2 != nil {
		err3, ok := err2.(error)
		err4 := fmt.Errorf("panic-error in user-engagement-scheduler-lambda %+v", err2)
		err = &err4
		if ok {
			alog.LambdaLogger(logrus.InfoLevel).WithError(err3).Errorf("Panic error in user-engagement-scheduler-lambda")
		} else {
			alog.LambdaLogger(logrus.InfoLevel).WithField("non-error", err2).Errorf("Unusual panic error in user-engagement-scheduler-lambda")
		}
	}
}
func HandleRequest(ctx context.Context, event EngSchedule) (err error) {
	defer handleError(&err)
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
		userDao.ReadUnsafe(event.Target).PlatformID,
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
		event.Target,
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
