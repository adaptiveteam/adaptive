package user
// This file has been automatically generated by `adaptive-platform/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/adaptiveteam/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/aws-utils-go"
	common "github.com/adaptiveteam/daos/common"
	core "github.com/adaptiveteam/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

type User struct  {
	ID string `json:"id"`
	DisplayName string `json:"display_name"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	Timezone string `json:"timezone"`
	TimezoneOffset int `json:"timezone_offset"`
	// in 24 hr format, localtime
	AdaptiveScheduledTime string `json:"adaptive_scheduled_time,omitempty"`
	AdaptiveScheduledTimeInUTC string `json:"adaptive_scheduled_time_in_utc,omitempty"`
	PlatformID models.PlatformID `json:"platform_id"`
	PlatformOrg string `json:"platform_org"`
	IsAdmin bool `json:"is_admin"`
	Deleted bool `json:"deleted"`
	CreatedAt string `json:"created_at"`
	ModifiedAt string `json:"modified_at,omitempty"`
	IsShared bool `json:"is_shared"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (user User)CollectEmptyFields() (emptyFields []string, ok bool) {
	if user.ID == "" { emptyFields = append(emptyFields, "ID")}
	if user.DisplayName == "" { emptyFields = append(emptyFields, "DisplayName")}
	if user.Timezone == "" { emptyFields = append(emptyFields, "Timezone")}
	if user.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if user.PlatformOrg == "" { emptyFields = append(emptyFields, "PlatformOrg")}
	if user.CreatedAt == "" { emptyFields = append(emptyFields, "CreatedAt")}
	ok = len(emptyFields) == 0
	return
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
	Delete(id string) error
	DeleteUnsafe(id string)
	ReadByID(id string) (user []User, err error)
	ReadByIDUnsafe(id string) (user []User)
	ReadByPlatformID(platformID models.PlatformID) (user []User, err error)
	ReadByPlatformIDUnsafe(platformID models.PlatformID) (user []User)
	ReadByPlatformIDTimezoneOffset(platformID models.PlatformID, timezoneOffset int) (user []User, err error)
	ReadByPlatformIDTimezoneOffsetUnsafe(platformID models.PlatformID, timezoneOffset int) (user []User)
	ReadByPlatformIDAdaptiveScheduledTimeInUTC(platformID models.PlatformID, adaptiveScheduledTimeInUTC string) (user []User, err error)
	ReadByPlatformIDAdaptiveScheduledTimeInUTCUnsafe(platformID models.PlatformID, adaptiveScheduledTimeInUTC string) (user []User)
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
		Name: clientID + "_user",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the User.
func (d DAOImpl) Create(user User) error {
	emptyFields, ok := user.CollectEmptyFields()
	if !ok {return fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)}
	return d.Dynamo.PutTableEntry(user, d.Name)
}


// CreateUnsafe saves the User.
func (d DAOImpl) CreateUnsafe(user User) {
	err := d.Create(user)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create id==%s in %s\n", user.ID, d.Name))
}


// Read reads User
func (d DAOImpl) Read(id string) (out User, err error) {
	var outs []User
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, d.Name)
	}
	return
}


// ReadUnsafe reads the User. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) User {
	out, err := d.Read(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading id==%s in %s\n", id, d.Name))
	return out
}


// ReadOrEmpty reads User
func (d DAOImpl) ReadOrEmpty(id string) (out []User, err error) {
	var outOrEmpty User
	ids := idParams(id)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.ID == id {
		out = append(out, outOrEmpty)
	}
	err = errors.Wrapf(err, "User DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the User. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []User {
	out, err := d.ReadOrEmpty(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading id==%s in %s\n", id, d.Name))
	return out
}


// CreateOrUpdate saves the User regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(user User) (err error) {
	
	var olds []User
	olds, err = d.ReadOrEmpty(user.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(user)
			err = errors.Wrapf(err, "User DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			old := olds[0]
			ids := idParams(old.ID)
			err = d.Dynamo.UpdateTableEntry(
				allParams(user, old),
				ids,
				updateExpression(user, old),
				d.Name,
			)
			err = errors.Wrapf(err, "User DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s", ids, d.Name)
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the User regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(user User) {
	err := d.CreateOrUpdate(user)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", user, d.Name))
}


// Delete removes User from db
func (d DAOImpl)Delete(id string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(id))
}


// DeleteUnsafe deletes User and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(id string) {
	err := d.Delete(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete id==%s in %s\n", id, d.Name))
}


func (d DAOImpl)ReadByID(id string) (out []User, err error) {
	var instances []User
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "IDIndex",
		Condition: "id = :a0",
		Attributes: map[string]interface{}{
			":a0": id,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByIDUnsafe(id string) (out []User) {
	out, err := d.ReadByID(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query IDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByPlatformID(platformID models.PlatformID) (out []User, err error) {
	var instances []User
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "PlatformIDIndex",
		Condition: "platform_id = :a0",
		Attributes: map[string]interface{}{
			":a0": platformID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID models.PlatformID) (out []User) {
	out, err := d.ReadByPlatformID(platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByPlatformIDTimezoneOffset(platformID models.PlatformID, timezoneOffset int) (out []User, err error) {
	var instances []User
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "PlatformIDTimezoneOffsetIndex",
		Condition: "platform_id = :a0 and timezone_offset = :a1",
		Attributes: map[string]interface{}{
			":a0": platformID,
			":a1": timezoneOffset,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByPlatformIDTimezoneOffsetUnsafe(platformID models.PlatformID, timezoneOffset int) (out []User) {
	out, err := d.ReadByPlatformIDTimezoneOffset(platformID, timezoneOffset)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query PlatformIDTimezoneOffsetIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByPlatformIDAdaptiveScheduledTimeInUTC(platformID models.PlatformID, adaptiveScheduledTimeInUTC string) (out []User, err error) {
	var instances []User
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "PlatformIDAdaptiveScheduledTimeInUTCIndex",
		Condition: "platform_id = :a0 and adaptive_scheduled_time_in_utc = :a1",
		Attributes: map[string]interface{}{
			":a0": platformID,
			":a1": adaptiveScheduledTimeInUTC,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByPlatformIDAdaptiveScheduledTimeInUTCUnsafe(platformID models.PlatformID, adaptiveScheduledTimeInUTC string) (out []User) {
	out, err := d.ReadByPlatformIDAdaptiveScheduledTimeInUTC(platformID, adaptiveScheduledTimeInUTC)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query PlatformIDAdaptiveScheduledTimeInUTCIndex on %s table\n", d.Name))
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
		if user.ID != old.ID { params["a0"] = common.DynS(user.ID) }
		if user.DisplayName != old.DisplayName { params["a1"] = common.DynS(user.DisplayName) }
		if user.FirstName != old.FirstName { params["a2"] = common.DynS(user.FirstName) }
		if user.LastName != old.LastName { params["a3"] = common.DynS(user.LastName) }
		if user.Timezone != old.Timezone { params["a4"] = common.DynS(user.Timezone) }
		if user.TimezoneOffset != old.TimezoneOffset { params["a5"] = common.DynN(user.TimezoneOffset) }
		if user.AdaptiveScheduledTime != old.AdaptiveScheduledTime { params["a6"] = common.DynS(user.AdaptiveScheduledTime) }
		if user.AdaptiveScheduledTimeInUTC != old.AdaptiveScheduledTimeInUTC { params["a7"] = common.DynS(user.AdaptiveScheduledTimeInUTC) }
		if user.PlatformID != old.PlatformID { params["a8"] = common.DynS(string(user.PlatformID)) }
		if user.PlatformOrg != old.PlatformOrg { params["a9"] = common.DynS(user.PlatformOrg) }
		if user.IsAdmin != old.IsAdmin { params["a10"] = common.DynBOOL(user.IsAdmin) }
		if user.Deleted != old.Deleted { params["a11"] = common.DynBOOL(user.Deleted) }
		if user.CreatedAt != old.CreatedAt { params["a12"] = common.DynS(user.CreatedAt) }
		if user.ModifiedAt != old.ModifiedAt { params["a13"] = common.DynS(user.ModifiedAt) }
		if user.IsShared != old.IsShared { params["a14"] = common.DynBOOL(user.IsShared) }
	return
}
func updateExpression(user User, old User) string {
	var updateParts []string
	
		
			
		if user.ID != old.ID { updateParts = append(updateParts, "id = :a0") }
		if user.DisplayName != old.DisplayName { updateParts = append(updateParts, "display_name = :a1") }
		if user.FirstName != old.FirstName { updateParts = append(updateParts, "first_name = :a2") }
		if user.LastName != old.LastName { updateParts = append(updateParts, "last_name = :a3") }
		if user.Timezone != old.Timezone { updateParts = append(updateParts, "#timezone = :a4") }
		if user.TimezoneOffset != old.TimezoneOffset { updateParts = append(updateParts, "timezone_offset = :a5") }
		if user.AdaptiveScheduledTime != old.AdaptiveScheduledTime { updateParts = append(updateParts, "adaptive_scheduled_time = :a6") }
		if user.AdaptiveScheduledTimeInUTC != old.AdaptiveScheduledTimeInUTC { updateParts = append(updateParts, "adaptive_scheduled_time_in_utc = :a7") }
		if user.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a8") }
		if user.PlatformOrg != old.PlatformOrg { updateParts = append(updateParts, "platform_org = :a9") }
		if user.IsAdmin != old.IsAdmin { updateParts = append(updateParts, "is_admin = :a10") }
		if user.Deleted != old.Deleted { updateParts = append(updateParts, "deleted = :a11") }
		if user.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a12") }
		if user.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a13") }
		if user.IsShared != old.IsShared { updateParts = append(updateParts, "is_shared = :a14") }
	return strings.Join(updateParts, " and ")
}
