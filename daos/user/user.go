package user
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"encoding/json"
	"strings"
)

type User struct  {
	ID string `json:"id"`
	DisplayName string `json:"display_name"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	Timezone string `json:"timezone"`
	IsAdaptiveBot bool `json:"is_adaptive_bot,omitempty"`
	TimezoneOffset int `json:"timezone_offset"`
	// in 24 hr format, localtime
	AdaptiveScheduledTime string `json:"adaptive_scheduled_time,omitempty"`
	AdaptiveScheduledTimeInUTC string `json:"adaptive_scheduled_time_in_utc,omitempty"`
	PlatformID common.PlatformID `json:"platform_id"`
	PlatformOrg string `json:"platform_org"`
	IsAdmin bool `json:"is_admin"`
	IsShared bool `json:"is_shared"`
	DeactivatedAt string `json:"deactivated_at,omitempty"`
	// Automatically maintained field
	CreatedAt string `json:"created_at"`
	// Automatically maintained field
	ModifiedAt string `json:"modified_at,omitempty"`
}

// UserFilterActive removes deactivated values
func UserFilterActive(in []User) (res []User) {
	for _, i := range in {
		if i.DeactivatedAt == "" {
			res = append(res, i)
		}
	}
	return
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (user User)CollectEmptyFields() (emptyFields []string, ok bool) {
	if user.ID == "" { emptyFields = append(emptyFields, "ID")}
	if user.DisplayName == "" { emptyFields = append(emptyFields, "DisplayName")}
	if user.Timezone == "" { emptyFields = append(emptyFields, "Timezone")}
	if user.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if user.PlatformOrg == "" { emptyFields = append(emptyFields, "PlatformOrg")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (user User) ToJSON() (string, error) {
	b, err := json.Marshal(user)
	return string(b), err
}

type DAO interface {
	Create(user User) error
	CreateUnsafe(user User)
	Read(id string) (user User, err error)
	ReadUnsafe(id string) (user User)
	ReadOrEmpty(id string) (user []User, err error)
	ReadOrEmptyUnsafe(id string) (user []User)
	CreateOrUpdate(user User) error
	CreateOrUpdateUnsafe(user User)
	Deactivate(id string) error
	DeactivateUnsafe(id string)
	ReadByPlatformID(platformID common.PlatformID) (user []User, err error)
	ReadByPlatformIDUnsafe(platformID common.PlatformID) (user []User)
	ReadByPlatformIDTimezoneOffset(platformID common.PlatformID, timezoneOffset int) (user []User, err error)
	ReadByPlatformIDTimezoneOffsetUnsafe(platformID common.PlatformID, timezoneOffset int) (user []User)
	ReadByPlatformIDAdaptiveScheduledTimeInUTC(platformID common.PlatformID, adaptiveScheduledTimeInUTC string) (user []User, err error)
	ReadByPlatformIDAdaptiveScheduledTimeInUTCUnsafe(platformID common.PlatformID, adaptiveScheduledTimeInUTC string) (user []User)
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
		Name: TableName(clientID),
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}
func TableName(clientID string) string {
	return clientID + "_user"
}

// Create saves the User.
func (d DAOImpl) Create(user User) (err error) {
	emptyFields, ok := user.CollectEmptyFields()
	if ok {
		user.ModifiedAt = core.CurrentRFCTimestamp()
	user.CreatedAt = user.ModifiedAt
	err = d.Dynamo.PutTableEntry(user, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the User.
func (d DAOImpl) CreateUnsafe(user User) {
	err2 := d.Create(user)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not create id==%s in %s\n", user.ID, d.Name))
}


// Read reads User
func (d DAOImpl) Read(id string) (out User, err error) {
	var outs []User
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, d.Name)
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the User. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) User {
	out, err2 := d.Read(id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Error reading id==%s in %s\n", id, d.Name))
	return out
}


// ReadOrEmpty reads User
func (d DAOImpl) ReadOrEmpty(id string) (out []User, err error) {
	var outOrEmpty User
	ids := idParams(id)
	var found bool
	found, err = d.Dynamo.GetItemOrEmptyFromTable(d.Name, ids, &outOrEmpty)
	if found {
		if outOrEmpty.ID == id {
			out = append(out, outOrEmpty)
		} else {
			err = fmt.Errorf("Requested ids: id==%s are different from the found ones: id==%s", id, outOrEmpty.ID) // unexpected error: found ids != ids
		}
	}
	err = errors.Wrapf(err, "User DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the User. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []User {
	out, err2 := d.ReadOrEmpty(id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Error while reading id==%s in %s\n", id, d.Name))
	return out
}


// CreateOrUpdate saves the User regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(user User) (err error) {
	user.ModifiedAt = core.CurrentRFCTimestamp()
	if user.CreatedAt == "" { user.CreatedAt = user.ModifiedAt }
	
	var olds []User
	olds, err = d.ReadOrEmpty(user.ID)
	err = errors.Wrapf(err, "User DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", user.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(user)
			err = errors.Wrapf(err, "User DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			emptyFields, ok := user.CollectEmptyFields()
			if ok {
				old := olds[0]
				user.ModifiedAt = core.CurrentRFCTimestamp()

				key := idParams(old.ID)
				expr, exprAttributes, names := updateExpression(user, old)
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
				err = errors.Wrapf(err, "User DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the User regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(user User) {
	err2 := d.CreateOrUpdate(user)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", user, d.Name))
}


// Deactivate "removes" User. 
// The mechanism is adding timestamp to `DeactivatedOn` field. 
// Then, if this field is not empty, the instance is considered to be "active"
func (d DAOImpl)Deactivate(id string) error {
	instance, err2 := d.Read(id)
	if err2 == nil {
		instance.DeactivatedAt = core.CurrentRFCTimestamp()
		err2 = d.CreateOrUpdate(instance)
	}
	return err2
}


// DeactivateUnsafe "deletes" User and panics in case of errors.
func (d DAOImpl)DeactivateUnsafe(id string) {
	err2 := d.Deactivate(id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not deactivate id==%s in %s\n", id, d.Name))
}


func (d DAOImpl)ReadByPlatformID(platformID common.PlatformID) (out []User, err error) {
	var instances []User
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "PlatformIDIndex",
		Condition: "platform_id = :a0",
		Attributes: map[string]interface{}{
			":a0": platformID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = UserFilterActive(instances)
	return
}


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID common.PlatformID) (out []User) {
	out, err2 := d.ReadByPlatformID(platformID)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByPlatformIDTimezoneOffset(platformID common.PlatformID, timezoneOffset int) (out []User, err error) {
	var instances []User
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "PlatformIDTimezoneOffsetIndex",
		Condition: "platform_id = :a0 and timezone_offset = :a1",
		Attributes: map[string]interface{}{
			":a0": platformID,
			":a1": timezoneOffset,
		},
	}, map[string]string{}, true, -1, &instances)
	out = UserFilterActive(instances)
	return
}


func (d DAOImpl)ReadByPlatformIDTimezoneOffsetUnsafe(platformID common.PlatformID, timezoneOffset int) (out []User) {
	out, err2 := d.ReadByPlatformIDTimezoneOffset(platformID, timezoneOffset)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not query PlatformIDTimezoneOffsetIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByPlatformIDAdaptiveScheduledTimeInUTC(platformID common.PlatformID, adaptiveScheduledTimeInUTC string) (out []User, err error) {
	var instances []User
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "PlatformIDAdaptiveScheduledTimeInUTCIndex",
		Condition: "platform_id = :a0 and adaptive_scheduled_time_in_utc = :a1",
		Attributes: map[string]interface{}{
			":a0": platformID,
			":a1": adaptiveScheduledTimeInUTC,
		},
	}, map[string]string{}, true, -1, &instances)
	out = UserFilterActive(instances)
	return
}


func (d DAOImpl)ReadByPlatformIDAdaptiveScheduledTimeInUTCUnsafe(platformID common.PlatformID, adaptiveScheduledTimeInUTC string) (out []User) {
	out, err2 := d.ReadByPlatformIDAdaptiveScheduledTimeInUTC(platformID, adaptiveScheduledTimeInUTC)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not query PlatformIDAdaptiveScheduledTimeInUTCIndex on %s table\n", d.Name))
	return
}

func idParams(id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"id": common.DynS(id),
	}
	return params
}
func allParams(user User, old User) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if user.ID != old.ID { params[":a0"] = common.DynS(user.ID) }
	if user.DisplayName != old.DisplayName { params[":a1"] = common.DynS(user.DisplayName) }
	if user.FirstName != old.FirstName { params[":a2"] = common.DynS(user.FirstName) }
	if user.LastName != old.LastName { params[":a3"] = common.DynS(user.LastName) }
	if user.Timezone != old.Timezone { params[":a4"] = common.DynS(user.Timezone) }
	if user.IsAdaptiveBot != old.IsAdaptiveBot { params[":a5"] = common.DynBOOL(user.IsAdaptiveBot) }
	if user.TimezoneOffset != old.TimezoneOffset { params[":a6"] = common.DynN(user.TimezoneOffset) }
	if user.AdaptiveScheduledTime != old.AdaptiveScheduledTime { params[":a7"] = common.DynS(user.AdaptiveScheduledTime) }
	if user.AdaptiveScheduledTimeInUTC != old.AdaptiveScheduledTimeInUTC { params[":a8"] = common.DynS(user.AdaptiveScheduledTimeInUTC) }
	if user.PlatformID != old.PlatformID { params[":a9"] = common.DynS(string(user.PlatformID)) }
	if user.PlatformOrg != old.PlatformOrg { params[":a10"] = common.DynS(user.PlatformOrg) }
	if user.IsAdmin != old.IsAdmin { params[":a11"] = common.DynBOOL(user.IsAdmin) }
	if user.IsShared != old.IsShared { params[":a12"] = common.DynBOOL(user.IsShared) }
	if user.DeactivatedAt != old.DeactivatedAt { params[":a13"] = common.DynS(user.DeactivatedAt) }
	if user.CreatedAt != old.CreatedAt { params[":a14"] = common.DynS(user.CreatedAt) }
	if user.ModifiedAt != old.ModifiedAt { params[":a15"] = common.DynS(user.ModifiedAt) }
	return
}
func updateExpression(user User, old User) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if user.ID != old.ID { updateParts = append(updateParts, "id = :a0"); params[":a0"] = common.DynS(user.ID);  }
	if user.DisplayName != old.DisplayName { updateParts = append(updateParts, "display_name = :a1"); params[":a1"] = common.DynS(user.DisplayName);  }
	if user.FirstName != old.FirstName { updateParts = append(updateParts, "first_name = :a2"); params[":a2"] = common.DynS(user.FirstName);  }
	if user.LastName != old.LastName { updateParts = append(updateParts, "last_name = :a3"); params[":a3"] = common.DynS(user.LastName);  }
	if user.Timezone != old.Timezone { updateParts = append(updateParts, "#timezone = :a4"); params[":a4"] = common.DynS(user.Timezone); fldName := "timezone"; names["#timezone"] = &fldName }
	if user.IsAdaptiveBot != old.IsAdaptiveBot { updateParts = append(updateParts, "is_adaptive_bot = :a5"); params[":a5"] = common.DynBOOL(user.IsAdaptiveBot);  }
	if user.TimezoneOffset != old.TimezoneOffset { updateParts = append(updateParts, "timezone_offset = :a6"); params[":a6"] = common.DynN(user.TimezoneOffset);  }
	if user.AdaptiveScheduledTime != old.AdaptiveScheduledTime { updateParts = append(updateParts, "adaptive_scheduled_time = :a7"); params[":a7"] = common.DynS(user.AdaptiveScheduledTime);  }
	if user.AdaptiveScheduledTimeInUTC != old.AdaptiveScheduledTimeInUTC { updateParts = append(updateParts, "adaptive_scheduled_time_in_utc = :a8"); params[":a8"] = common.DynS(user.AdaptiveScheduledTimeInUTC);  }
	if user.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a9"); params[":a9"] = common.DynS(string(user.PlatformID));  }
	if user.PlatformOrg != old.PlatformOrg { updateParts = append(updateParts, "platform_org = :a10"); params[":a10"] = common.DynS(user.PlatformOrg);  }
	if user.IsAdmin != old.IsAdmin { updateParts = append(updateParts, "is_admin = :a11"); params[":a11"] = common.DynBOOL(user.IsAdmin);  }
	if user.IsShared != old.IsShared { updateParts = append(updateParts, "is_shared = :a12"); params[":a12"] = common.DynBOOL(user.IsShared);  }
	if user.DeactivatedAt != old.DeactivatedAt { updateParts = append(updateParts, "deactivated_at = :a13"); params[":a13"] = common.DynS(user.DeactivatedAt);  }
	if user.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a14"); params[":a14"] = common.DynS(user.CreatedAt);  }
	if user.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a15"); params[":a15"] = common.DynS(user.ModifiedAt);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
