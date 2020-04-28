package holidays

import (
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/adHocHoliday"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

// All reads all ad-hoc holidays for the given PlatformID 
// from dynamo table
func All(conn daosCommon.DynamoDBConnection) ([]models.AdHocHoliday, error){
	var res []models.AdHocHoliday
	err := conn.Dynamo.QueryTableWithIndex(adHocHoliday.TableName(conn.ClientID), 
		awsutils.DynamoIndexExpression{
		IndexName: string(adHocHoliday.PlatformIDDateIndex),
		Condition: "platform_id = :platform_id",
		Attributes: map[string]interface{}{
			":platform_id": conn.PlatformID,
		},
	}, map[string]string{}, true, -1, &res)
	return res,err
}
// AllUnsafe reads all ad-hoc holidays for PlatformID and panics in case of errors
func AllUnsafe(conn daosCommon.DynamoDBConnection) []models.AdHocHoliday{
	holidays, err2 := All(conn)
	core.ErrorHandler(err2, "AllUnsafe", "Could not query table adHocHoliday table")
	return holidays
}
// SelectNotEarlierThan reads all ad-hoc holidays from dynamo table
// that are later or at the given time moment
func SelectNotEarlierThan(time time.Time) func (conn daosCommon.DynamoDBConnection) ([]models.AdHocHoliday, error){
	return func (conn daosCommon.DynamoDBConnection) ([]models.AdHocHoliday, error){
		var res []models.AdHocHoliday
		err := conn.Dynamo.QueryTableWithIndex(adHocHoliday.TableName(conn.ClientID), 
			awsutils.DynamoIndexExpression{
			IndexName: string(adHocHoliday.PlatformIDDateIndex),
			// there is no != operator for ConditionExpression
			Condition: "platform_id = :platform_id AND #date >= :target_date",
			Attributes: map[string]interface{}{
				":platform_id": conn.PlatformID,
				":target_date": aws.String(time.Format(models.AdHocHolidayDateFormat)),
			},
		}, map[string]string{"#date": "date"}, true, -1, &res)
		return res,err
	}
}
// SelectNotEarlierThanUnsafe reads all ad-hoc holidays from dynamo table
// that are later or at the given time moment. Panics in case of any errors
func SelectNotEarlierThanUnsafe(time time.Time)  func (conn daosCommon.DynamoDBConnection) []models.AdHocHoliday{
	return func (conn daosCommon.DynamoDBConnection) []models.AdHocHoliday{
		holidays, err2 := SelectNotEarlierThan(time)(conn)
		core.ErrorHandler(err2, "SelectNotEarlierThanUnsafe", "Could not query holiday by date")
		return holidays
	}
}
