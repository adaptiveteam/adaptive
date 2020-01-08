package adHocHoliday
// This file has been automatically generated by `adaptive-platform/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

// AdHocHoliday is a holiday on exact date.
type AdHocHoliday struct  {
	ID string `json:"id"`
	PlatformID models.PlatformID `json:"platform_id"`
	Date string `json:"date"`
	Name string `json:"name"`
	Description string `json:"description"`
	ScopeCommunities string `json:"scope_communities"`
	DeactivatedOn string `json:"deactivated_on"`
}

// AdHocHolidayFilterActive removes deactivated values
func AdHocHolidayFilterActive(in []AdHocHoliday) (res []AdHocHoliday) {
	for _, i := range in {
		if i.DeactivatedOn == "" {
			res = append(res, i)
		}
	}
	return
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (adHocHoliday AdHocHoliday)CollectEmptyFields() (emptyFields []string, ok bool) {
	if adHocHoliday.ID == "" { emptyFields = append(emptyFields, "ID")}
	if adHocHoliday.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if adHocHoliday.Date == "" { emptyFields = append(emptyFields, "Date")}
	if adHocHoliday.Name == "" { emptyFields = append(emptyFields, "Name")}
	if adHocHoliday.Description == "" { emptyFields = append(emptyFields, "Description")}
	if adHocHoliday.ScopeCommunities == "" { emptyFields = append(emptyFields, "ScopeCommunities")}
	ok = len(emptyFields) == 0
	return
}

type DAO interface {
	Create(adHocHoliday AdHocHoliday) error
	CreateUnsafe(adHocHoliday AdHocHoliday)
	Read(id string) (adHocHoliday AdHocHoliday, err error)
	ReadUnsafe(id string) (adHocHoliday AdHocHoliday)
	ReadOrEmpty(id string) (adHocHoliday []AdHocHoliday, err error)
	ReadOrEmptyUnsafe(id string) (adHocHoliday []AdHocHoliday)
	CreateOrUpdate(adHocHoliday AdHocHoliday) error
	CreateOrUpdateUnsafe(adHocHoliday AdHocHoliday)
	Deactivate(id string) error
	DeactivateUnsafe(id string)
	ReadByDatePlatformID(date string, platformID models.PlatformID) (adHocHoliday []AdHocHoliday, err error)
	ReadByDatePlatformIDUnsafe(date string, platformID models.PlatformID) (adHocHoliday []AdHocHoliday)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	Name      string                  `json:"name"`
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic("Cannot create DAO without clientID") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: clientID + "_ad_hoc_holiday",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the AdHocHoliday.
func (d DAOImpl) Create(adHocHoliday AdHocHoliday) (err error) {
	emptyFields, ok := adHocHoliday.CollectEmptyFields()
	if ok {
		err = d.Dynamo.PutTableEntry(adHocHoliday, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the AdHocHoliday.
func (d DAOImpl) CreateUnsafe(adHocHoliday AdHocHoliday) {
	err := d.Create(adHocHoliday)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create id==%s in %s\n", adHocHoliday.ID, d.Name))
}


// Read reads AdHocHoliday
func (d DAOImpl) Read(id string) (out AdHocHoliday, err error) {
	var outs []AdHocHoliday
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, d.Name)
	}
	return
}


// ReadUnsafe reads the AdHocHoliday. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) AdHocHoliday {
	out, err := d.Read(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading id==%s in %s\n", id, d.Name))
	return out
}


// ReadOrEmpty reads AdHocHoliday
func (d DAOImpl) ReadOrEmpty(id string) (out []AdHocHoliday, err error) {
	var outOrEmpty AdHocHoliday
	ids := idParams(id)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.ID == id {
		out = append(out, outOrEmpty)
	} else if err != nil && strings.HasPrefix(err.Error(), "In table ") {
		err = nil // expected not-found error	
	}
	err = errors.Wrapf(err, "AdHocHoliday DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the AdHocHoliday. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []AdHocHoliday {
	out, err := d.ReadOrEmpty(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading id==%s in %s\n", id, d.Name))
	return out
}


// CreateOrUpdate saves the AdHocHoliday regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(adHocHoliday AdHocHoliday) (err error) {
	
	var olds []AdHocHoliday
	olds, err = d.ReadOrEmpty(adHocHoliday.ID)
	err = errors.Wrapf(err, "AdHocHoliday DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", adHocHoliday.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(adHocHoliday)
			err = errors.Wrapf(err, "AdHocHoliday DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			emptyFields, ok := adHocHoliday.CollectEmptyFields()
			if ok {
				old := olds[0]
				
				key := idParams(old.ID)
				expr, exprAttributes, names := updateExpression(adHocHoliday, old)
				input := dynamodb.UpdateItemInput{
					ExpressionAttributeValues: exprAttributes,
					TableName:                 aws.String(d.Name),
					Key:                       key,
					ReturnValues:              aws.String("UPDATED_NEW"),
					UpdateExpression:          aws.String(expr),
				}
				if names != nil { input.ExpressionAttributeNames = *names } // workaround for a pointer to an empty slice
				if err == nil {
					err = d.Dynamo.UpdateItemInternal(input)
				}
				err = errors.Wrapf(err, "AdHocHoliday DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the AdHocHoliday regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(adHocHoliday AdHocHoliday) {
	err := d.CreateOrUpdate(adHocHoliday)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", adHocHoliday, d.Name))
}


// Deactivate "removes" AdHocHoliday. 
// The mechanism is adding timestamp to `DeactivatedOn` field. 
// Then, if this field is not empty, the instance is considered to be "active"
func (d DAOImpl)Deactivate(id string) error {
	instance, err := d.Read(id)
	if err == nil {
		instance.DeactivatedOn = core.ISODateLayout.Format(time.Now())
		err = d.CreateOrUpdate(instance)
	}
	return err
}


// DeactivateUnsafe "deletes" AdHocHoliday and panics in case of errors.
func (d DAOImpl)DeactivateUnsafe(id string) {
	err := d.Deactivate(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not deactivate id==%s in %s\n", id, d.Name))
}


func (d DAOImpl)ReadByDatePlatformID(date string, platformID models.PlatformID) (out []AdHocHoliday, err error) {
	var instances []AdHocHoliday
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "DatePlatformIDIndex",
		Condition: "#date = :a0 and platform_id = :a1",
		Attributes: map[string]interface{}{
			":a0": date,
			":a1": platformID,
		},
	}, map[string]string{"#date": "date"}, true, -1, &instances)
	out = AdHocHolidayFilterActive(instances)
	return
}


func (d DAOImpl)ReadByDatePlatformIDUnsafe(date string, platformID models.PlatformID) (out []AdHocHoliday) {
	out, err := d.ReadByDatePlatformID(date, platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query DatePlatformIDIndex on %s table\n", d.Name))
	return
}

func idParams(id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"id": common.DynS(id),
	}
	return params
}
func allParams(adHocHoliday AdHocHoliday, old AdHocHoliday) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if adHocHoliday.ID != old.ID { params[":a0"] = common.DynS(adHocHoliday.ID) }
	if adHocHoliday.PlatformID != old.PlatformID { params[":a1"] = common.DynS(string(adHocHoliday.PlatformID)) }
	if adHocHoliday.Date != old.Date { params[":a2"] = common.DynS(adHocHoliday.Date) }
	if adHocHoliday.Name != old.Name { params[":a3"] = common.DynS(adHocHoliday.Name) }
	if adHocHoliday.Description != old.Description { params[":a4"] = common.DynS(adHocHoliday.Description) }
	if adHocHoliday.ScopeCommunities != old.ScopeCommunities { params[":a5"] = common.DynS(adHocHoliday.ScopeCommunities) }
	if adHocHoliday.DeactivatedOn != old.DeactivatedOn { params[":a6"] = common.DynS(adHocHoliday.DeactivatedOn) }
	return
}
func updateExpression(adHocHoliday AdHocHoliday, old AdHocHoliday) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if adHocHoliday.ID != old.ID { updateParts = append(updateParts, "id = :a0"); params[":a0"] = common.DynS(adHocHoliday.ID);  }
	if adHocHoliday.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a1"); params[":a1"] = common.DynS(string(adHocHoliday.PlatformID));  }
	if adHocHoliday.Date != old.Date { updateParts = append(updateParts, "#date = :a2"); params[":a2"] = common.DynS(adHocHoliday.Date); fldName := "date"; names["#date"] = &fldName }
	if adHocHoliday.Name != old.Name { updateParts = append(updateParts, "#name = :a3"); params[":a3"] = common.DynS(adHocHoliday.Name); fldName := "name"; names["#name"] = &fldName }
	if adHocHoliday.Description != old.Description { updateParts = append(updateParts, "description = :a4"); params[":a4"] = common.DynS(adHocHoliday.Description);  }
	if adHocHoliday.ScopeCommunities != old.ScopeCommunities { updateParts = append(updateParts, "scope_communities = :a5"); params[":a5"] = common.DynS(adHocHoliday.ScopeCommunities);  }
	if adHocHoliday.DeactivatedOn != old.DeactivatedOn { updateParts = append(updateParts, "deactivated_on = :a6"); params[":a6"] = common.DynS(adHocHoliday.DeactivatedOn);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
