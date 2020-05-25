package userAttribute
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"encoding/json"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
)

// UserAttribute encapsulates key-value setting for a user
type UserAttribute struct  {
	// UserID is the ID of the user to send an engagement to
	// This usually corresponds to the platform user id
	UserID string `json:"user_id"`
	// Key of the setting
	AttrKey string `json:"attr_key"`
	// Value of the setting
	AttrValue string `json:"attr_value"`
	// A flag that tells whether setting is default or is explicitly set
	// Every user will have default settings
	Default bool `json:"default"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (userAttribute UserAttribute)CollectEmptyFields() (emptyFields []string, ok bool) {
	if userAttribute.UserID == "" { emptyFields = append(emptyFields, "UserID")}
	if userAttribute.AttrKey == "" { emptyFields = append(emptyFields, "AttrKey")}
	if userAttribute.AttrValue == "" { emptyFields = append(emptyFields, "AttrValue")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (userAttribute UserAttribute) ToJSON() (string, error) {
	b, err := json.Marshal(userAttribute)
	return string(b), err
}

type DAO interface {
	Create(userAttribute UserAttribute) error
	CreateUnsafe(userAttribute UserAttribute)
	Read(userID string, attrKey string) (userAttribute UserAttribute, err error)
	ReadUnsafe(userID string, attrKey string) (userAttribute UserAttribute)
	ReadOrEmpty(userID string, attrKey string) (userAttribute []UserAttribute, err error)
	ReadOrEmptyUnsafe(userID string, attrKey string) (userAttribute []UserAttribute)
	CreateOrUpdate(userAttribute UserAttribute) error
	CreateOrUpdateUnsafe(userAttribute UserAttribute)
	Delete(userID string, attrKey string) error
	DeleteUnsafe(userID string, attrKey string)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	ConnGen   common.DynamoDBConnectionGen
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic(errors.New("Cannot create UserAttribute.DAO without clientID")) }
	return DAOImpl{
		ConnGen:   common.DynamoDBConnectionGen{
			Dynamo: dynamo, 
			TableNamePrefix: clientID,
		},
	}
}
// TableNameSuffixVar is a global variable that contains table name suffix.
// After renaming all tables this may be made `const`.
var TableNameSuffixVar = "_user_attribute"

// TableName concatenates table name prefix and suffix and returns table name
func TableName(prefix string) string {
	return prefix + TableNameSuffixVar
}

// Create saves the UserAttribute.
func (d DAOImpl) Create(userAttribute UserAttribute) (err error) {
	emptyFields, ok := userAttribute.CollectEmptyFields()
	if ok {
		err = d.ConnGen.Dynamo.PutTableEntry(userAttribute, TableName(d.ConnGen.TableNamePrefix))
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the UserAttribute.
func (d DAOImpl) CreateUnsafe(userAttribute UserAttribute) {
	err2 := d.Create(userAttribute)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not create userID==%s, attrKey==%s in %s\n", userAttribute.UserID, userAttribute.AttrKey, TableName(d.ConnGen.TableNamePrefix)))
}


// Read reads UserAttribute
func (d DAOImpl) Read(userID string, attrKey string) (out UserAttribute, err error) {
	var outs []UserAttribute
	outs, err = d.ReadOrEmpty(userID, attrKey)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found userID==%s, attrKey==%s in %s\n", userID, attrKey, TableName(d.ConnGen.TableNamePrefix))
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the UserAttribute. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(userID string, attrKey string) UserAttribute {
	out, err2 := d.Read(userID, attrKey)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error reading userID==%s, attrKey==%s in %s\n", userID, attrKey, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// ReadOrEmpty reads UserAttribute
func (d DAOImpl) ReadOrEmpty(userID string, attrKey string) (out []UserAttribute, err error) {
	var outOrEmpty UserAttribute
	ids := idParams(userID, attrKey)
	var found bool
	found, err = d.ConnGen.Dynamo.GetItemOrEmptyFromTable(TableName(d.ConnGen.TableNamePrefix), ids, &outOrEmpty)
	if found {
		if outOrEmpty.UserID == userID && outOrEmpty.AttrKey == attrKey {
			out = append(out, outOrEmpty)
		} else {
			err = fmt.Errorf("Requested ids: userID==%s, attrKey==%s are different from the found ones: userID==%s, attrKey==%s", userID, attrKey, outOrEmpty.UserID, outOrEmpty.AttrKey) // unexpected error: found ids != ids
		}
	}
	err = errors.Wrapf(err, "UserAttribute DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(d.ConnGen.TableNamePrefix))
	return
}


// ReadOrEmptyUnsafe reads the UserAttribute. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(userID string, attrKey string) []UserAttribute {
	out, err2 := d.ReadOrEmpty(userID, attrKey)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error while reading userID==%s, attrKey==%s in %s\n", userID, attrKey, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// CreateOrUpdate saves the UserAttribute regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(userAttribute UserAttribute) (err error) {
	
	var olds []UserAttribute
	olds, err = d.ReadOrEmpty(userAttribute.UserID, userAttribute.AttrKey)
	err = errors.Wrapf(err, "UserAttribute DAO.CreateOrUpdate(id = userID==%s, attrKey==%s) couldn't ReadOrEmpty", userAttribute.UserID, userAttribute.AttrKey)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(userAttribute)
			err = errors.Wrapf(err, "UserAttribute DAO.CreateOrUpdate couldn't Create in table %s", TableName(d.ConnGen.TableNamePrefix))
		} else {
			emptyFields, ok := userAttribute.CollectEmptyFields()
			if ok {
				old := olds[0]
				
				
				key := idParams(old.UserID, old.AttrKey)
				expr, exprAttributes, names := updateExpression(userAttribute, old)
				input := dynamodb.UpdateItemInput{
					ExpressionAttributeValues: exprAttributes,
					TableName:                 aws.String(TableName(d.ConnGen.TableNamePrefix)),
					Key:                       key,
					ReturnValues:              aws.String("UPDATED_NEW"),
					UpdateExpression:          aws.String(expr),
				}
				if names != nil { input.ExpressionAttributeNames = *names } // workaround for a pointer to an empty slice
				if  len(exprAttributes) > 0 { // if there some changes
					err = d.ConnGen.Dynamo.UpdateItemInternal(input)
				} else {
					// WARN: no changes.
				}
				err = errors.Wrapf(err, "UserAttribute DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(d.ConnGen.TableNamePrefix), expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the UserAttribute regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(userAttribute UserAttribute) {
	err2 := d.CreateOrUpdate(userAttribute)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("could not create or update %v in %s\n", userAttribute, TableName(d.ConnGen.TableNamePrefix)))
}


// Delete removes UserAttribute from db
func (d DAOImpl)Delete(userID string, attrKey string) error {
	return d.ConnGen.Dynamo.DeleteEntry(TableName(d.ConnGen.TableNamePrefix), idParams(userID, attrKey))
}


// DeleteUnsafe deletes UserAttribute and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(userID string, attrKey string) {
	err2 := d.Delete(userID, attrKey)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not delete userID==%s, attrKey==%s in %s\n", userID, attrKey, TableName(d.ConnGen.TableNamePrefix)))
}

func idParams(userID string, attrKey string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"user_id": common.DynS(userID),
		"attr_key": common.DynS(attrKey),
	}
	return params
}
func allParams(userAttribute UserAttribute, old UserAttribute) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if userAttribute.UserID != old.UserID { params[":a0"] = common.DynS(userAttribute.UserID) }
	if userAttribute.AttrKey != old.AttrKey { params[":a1"] = common.DynS(userAttribute.AttrKey) }
	if userAttribute.AttrValue != old.AttrValue { params[":a2"] = common.DynS(userAttribute.AttrValue) }
	if userAttribute.Default != old.Default { params[":a3"] = common.DynBOOL(userAttribute.Default) }
	return
}
func updateExpression(userAttribute UserAttribute, old UserAttribute) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if userAttribute.UserID != old.UserID { updateParts = append(updateParts, "user_id = :a0"); params[":a0"] = common.DynS(userAttribute.UserID);  }
	if userAttribute.AttrKey != old.AttrKey { updateParts = append(updateParts, "attr_key = :a1"); params[":a1"] = common.DynS(userAttribute.AttrKey);  }
	if userAttribute.AttrValue != old.AttrValue { updateParts = append(updateParts, "attr_value = :a2"); params[":a2"] = common.DynS(userAttribute.AttrValue);  }
	if userAttribute.Default != old.Default { updateParts = append(updateParts, "#default = :a3"); params[":a3"] = common.DynBOOL(userAttribute.Default); fldName := "default"; names["#default"] = &fldName }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
