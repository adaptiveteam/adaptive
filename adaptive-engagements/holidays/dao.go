package holidays

import (
	"github.com/pkg/errors"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/adHocHoliday"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"time"
	// awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

// DAO - wrapper around a Dynamo DB table to work with holidays inside it
type DAO interface {
	ForPlatformID(teamID models.TeamID) PlatformDAO
}
// PlatformDAO is a set of utilities that work for a fixed `teamID`
type PlatformDAO interface {
	SelectNotEarlierThan(time time.Time) ([]models.AdHocHoliday, error)
	SelectNotEarlierThanUnsafe(time time.Time) []models.AdHocHoliday
}

// TableConfig is a structure that contains all configuration information from environment
type TableConfig struct {
	Table string
	PlatformDateIndex string
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	DNS   *common.DynamoNamespace
	TableConfig
}

// PlatformDAOImpl DAO that will implement PlatformDAO interface
type PlatformDAOImpl struct {
	TeamID models.TeamID
	DAOImpl
}
const (
	adHocHolidaysTableKey             = "HOLIDAYS_AD_HOC_TABLE"
	adHocHolidaysPlatformDateIndexKey = "HOLIDAYS_PLATFORM_DATE_INDEX"
)

func nonEmpty(env func(string)string) func(string)string {
	return func(key string)string {
		value := env(key)
		if value == "" {
			panic(errors.New("Key " + key + " is not defined"))
		}
		return value
	}
}
// ReadTableConfigUnsafe reads configuration values from environment
func ReadTableConfigUnsafe(env func(string)string) TableConfig {
	senv := nonEmpty(env)
	adHocHolidaysTable := senv(adHocHolidaysTableKey)
	adHocHolidaysPlatformDateIndex := senv(adHocHolidaysPlatformDateIndexKey)
	return TableConfig{
		Table: adHocHolidaysTable,
		PlatformDateIndex: adHocHolidaysPlatformDateIndex,
	}
}
// NewDAO creates an instance of DAO that will provide access to holidays table
func NewDAO(dns *common.DynamoNamespace, table string, index string) DAO {
	if table == "" || index == "" {
		panic(errors.New("Cannot create Holidays DAO without table and index"))
	}
	return DAOImpl{DNS: dns, TableConfig: TableConfig{Table: table, PlatformDateIndex: index}}
}

// ForPlatformID creates PlatformDAO that can be used for queries with PlatformID
func (d DAOImpl) ForPlatformID(teamID models.TeamID) PlatformDAO {
	return PlatformDAOImpl{
		TeamID: teamID,
		DAOImpl: d,
	}
}

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
func (p PlatformDAOImpl)SelectNotEarlierThan(time time.Time) ([]models.AdHocHoliday, error){
	var res []models.AdHocHoliday
	err := p.DNS.Dynamo.QueryTableWithIndex(p.Table, 
		awsutils.DynamoIndexExpression{
		IndexName: p.PlatformDateIndex,
		// there is no != operator for ConditionExpression
		Condition: "platform_id = :platform_id AND #date >= :target_date",
		Attributes: map[string]interface{}{
			":platform_id": p.TeamID.ToPlatformID(),
			":target_date": aws.String(time.Format(models.AdHocHolidayDateFormat)),
		},
	}, map[string]string{"#date": "date"}, true, -1, &res)
	return res,err
}
// SelectNotEarlierThanUnsafe reads all ad-hoc holidays from dynamo table
// that are later or at the given time moment. Panics in case of any errors
func (p PlatformDAOImpl)SelectNotEarlierThanUnsafe(time time.Time) []models.AdHocHoliday{
	holidays, err := p.SelectNotEarlierThan(time)
	core.ErrorHandler(err, p.DNS.Namespace, fmt.Sprintf("Could not query %s index on %s table", p.PlatformDateIndex, p.Table))
	return holidays
}
