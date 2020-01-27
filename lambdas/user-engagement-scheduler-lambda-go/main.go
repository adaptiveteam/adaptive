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

type EngSchedule struct {
	Target string `json:"target"`
	Date   string `json:"date"`
}

var (
	clientID = utils.NonEmptyEnv("CLIENT_ID")
	schema   = models.SchemaForClientID(clientID)
	userDao  = utilsUser.NewDAOFromSchema(d, namespace, schema)
)

// HandleRequest -
func HandleRequest(ctx context.Context, event EngSchedule) (err error) {
	defer core.RecoverToErrorVar("user-engagement-scheduler-lambda-go", &err)
	if event.Target == "" {
		return runScript(ctx)
	}
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


func runScript(ctx context.Context) (err error) {
	userIDs := []string{
		"UR4T3ULGP",
		"ULVV98V38",
		"ULVV98X2S",
		"UMADDK909",
		"ULTQQDB9N",
		"UN90UKW3C",
		"ULTQQF3MW",
		"UNL5DTE1F",
		"ULGCU5XL2",
		"UNMPN3RA4",
		"ULTB1N7MJ",
		"ULMETS9R7",
		"ULV1EPHNU",
		"ULTB1NB0U",
		"ULDA7AK4G",
		"ULVLLHH9D",
		"ULTRB1E2Z",
		"UPLBVFJQL",
		"ULTQQE6AU",
		"ULMETRDS5",
		"UQF1A5HCZ",
		"UR31V81T2",
		"ULWSQTGMV",
		"ULTQQFDEU",
		"ULTRB26UV",
		"ULVV9AQP8",
		"UMGD08VA5",
		"UL3TBJ5GS",
		"UMK02N9FC",
		"ULTRB2D7F",
		"UMK02N44A",
	}
	for _, userID := range userIDs {
		s := EngSchedule{
			Target: userID,
			Date: "2019-12-23",
		}
		HandleRequest(ctx, s)
	}
	return
}
