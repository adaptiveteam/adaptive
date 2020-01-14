package holidays

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"time"
	// awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

// DAO - wrapper around a Dynamo DB table to work with holidays inside it
type DAO interface {
	AddAdHocHoliday(holiday models.AdHocHoliday) error

	Create(holiday models.AdHocHoliday) error
	Read(holidayID string) (models.AdHocHoliday, error)
	ReadUnsafe(holidayID string) models.AdHocHoliday
	Update(holiday models.AdHocHoliday) error
	Delete(holidayID string) error

	ForPlatformID(platformID models.PlatformID) PlatformDAO
}
// PlatformDAO is a set of utilities that work for a fixed `platformID`
type PlatformDAO interface {
	All() ([]models.AdHocHoliday, error)
	AllUnsafe() []models.AdHocHoliday
	SelectNotEarlierThan(time time.Time) ([]models.AdHocHoliday, error)
	SelectNotEarlierThanUnsafe(time time.Time) []models.AdHocHoliday
	Create(holiday models.AdHocHoliday) error
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
	PlatformID models.PlatformID
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
			panic("Key " + key + " is not defined")
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
		panic("Cannot create Holidays DAO without table and index")
	}
	return DAOImpl{DNS: dns, TableConfig: TableConfig{Table: table, PlatformDateIndex: index}}
}

// Create creates an ad-hoc holiday
func (d DAOImpl) Create(holiday models.AdHocHoliday) error {
	return d.DNS.Dynamo.PutTableEntry(holiday, d.Table)
}

// AddAdHocHoliday creates an ad-hoc holiday
// deprecated. Use .Create
func (d DAOImpl) AddAdHocHoliday(holiday models.AdHocHoliday) error {
	return d.Create(holiday)
}

// Read reads the ad-hoc holiday
func (d DAOImpl) Read(holidayID string) (models.AdHocHoliday, error) {
	params := map[string]*dynamodb.AttributeValue{
		"id": {S: aws.String(holidayID)},
	}
	var out models.AdHocHoliday
	err := d.DNS.Dynamo.QueryTable(d.Table, params, &out)
	//if len(out) < 1 { return models.AdHocHoliday{}, errors.New("NotFound AdHocHoliday#ID=" + holidayID) }
	//if len(out) > 1 { return models.AdHocHoliday{}, errors.New("Found many AdHocHoliday#ID=" + holidayID) }
	return out, err
}

// ReadUnsafe reads the ad-hoc holiday. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(holidayID string) models.AdHocHoliday {
	holiday, err := d.Read(holidayID)
	core.ErrorHandler(err, d.DNS.Namespace, fmt.Sprintf("Could not find %s in %s", holidayID, d.Table))
	return holiday
}

// Delete removes the ad-hoc holiday
func (d DAOImpl) Delete(holidayID string) error {
	userParams := map[string]*dynamodb.AttributeValue{
		"id": {S: aws.String(holidayID)},
	}
	return d.DNS.Dynamo.DeleteEntry(d.Table, userParams)
}

// Update updates the ad-hoc holiday by ID
func (d DAOImpl) Update(holiday models.AdHocHoliday) error {
	exprAttributes := map[string]*dynamodb.AttributeValue{
		":name1": {
			S: aws.String(holiday.Name),
		},
		":description": {
			S: aws.String(holiday.Description),
		},
		":date1": {
			S: aws.String(holiday.Date),
		},
		":location": {
			S: aws.String(holiday.ScopeCommunities),
		},
		":platform_id": {
			S: aws.String(string(holiday.PlatformID)),
		},
	}
	key := map[string]*dynamodb.AttributeValue{
		"id": {S: aws.String(holiday.ID)},
	}
	n := "name"
	da := "date"
	updateExpression := "set #n = :name1, description = :description, #da = :date1, scope_communities = :location, platform_id = :platform_id"
	input := dynamodb.UpdateItemInput{
		ExpressionAttributeValues: exprAttributes,
		TableName:                 aws.String(d.Table),
		Key:                       key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  map[string]*string{"#n": &n, "#da": &da},
	}

	return d.DNS.Dynamo.UpdateItemInternal(input)
}

// ForPlatformID creates PlatformDAO that can be used for queries with PlatformID
func (d DAOImpl) ForPlatformID(platformID models.PlatformID) PlatformDAO {
	return PlatformDAOImpl{
		PlatformID: platformID,
		DAOImpl: d,
	}
}

// All reads all ad-hoc holidays for the given PlatformID 
// from dynamo table
func (p PlatformDAOImpl)All() ([]models.AdHocHoliday, error){
	var res []models.AdHocHoliday
	err := p.DNS.Dynamo.QueryTableWithIndex(p.Table, 
		awsutils.DynamoIndexExpression{
		IndexName: p.PlatformDateIndex,
		Condition: "platform_id = :platform_id",
		Attributes: map[string]interface{}{
			":platform_id": p.PlatformID,
		},
	}, map[string]string{}, true, -1, &res)
	return res,err
}
// AllUnsafe reads all ad-hoc holidays for PlatformID and panics in case of errors
func (p PlatformDAOImpl)AllUnsafe() []models.AdHocHoliday{
	holidays, err := p.All()
	core.ErrorHandler(err, p.DNS.Namespace, fmt.Sprintf("Could not query table %s", p.Table))
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
			":platform_id": p.PlatformID,
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
// Create creates an ad-hoc holiday making sure that PlatformID is correct
func (p PlatformDAOImpl)Create(holiday models.AdHocHoliday) error {
	holiday2 := holiday
	holiday2.PlatformID = p.PlatformID
	return p.DNS.Dynamo.PutTableEntry(holiday2, p.Table)
}
