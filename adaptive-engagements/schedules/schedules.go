package schedules

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/holidays"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"time"
)

func LoadHolidays(time time.Time, teamID models.TeamID, holidaysTable, holidaysPlatformDateIndex string) business_time.Holidays {
	var res []models.AdHocHoliday
	namespace := common.DeprecatedGetGlobalDns().Namespace
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(holidaysTable, awsutils.DynamoIndexExpression{
		IndexName: holidaysPlatformDateIndex,
		// there is no != operator for ConditionExpression
		Condition: "platform_id = :pl AND #date >= :target_date",
		Attributes: map[string]interface{}{
			":pl":          teamID.ToString(),
			":target_date": time.Format(models.AdHocHolidayDateFormat),
		},
	}, map[string]string{"#date": "date"}, true, -1, &res)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not query %s index on %s table", holidaysPlatformDateIndex,
		holidaysTable))
	return holidays.ConvertHolidaysArray(res)
}
