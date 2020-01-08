package userAttribute
// This file has been automatically generated by `adaptive-platform/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

// UserAttribute encapsulates key-value setting for a user
type UserAttribute struct  {
	// UserId is the Id of the user to send an engagement to
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
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	Name      string                  `json:"name"`
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic("Cannot create DAO without clientID") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: clientID + "_user_attribute",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the UserAttribute.
func (d DAOImpl) Create(userAttribute UserAttribute) (err error) {
	emptyFields, ok := userAttribute.CollectEmptyFields()
	if ok {
		err = d.Dynamo.PutTableEntry(userAttribute, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the UserAttribute.
func (d DAOImpl) CreateUnsafe(userAttribute UserAttribute) {
	err := d.Create(userAttribute)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create userID==%s, attrKey==%s in %s\n", userAttribute.UserID, userAttribute.AttrKey, d.Name))
}


// Read reads UserAttribute
func (d DAOImpl) Read(userID string, attrKey string) (out UserAttribute, err error) {
	var outs []UserAttribute
	outs, err = d.ReadOrEmpty(userID, attrKey)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found userID==%s, attrKey==%s in %s\n", userID, attrKey, d.Name)
	}
	return
}


// ReadUnsafe reads the UserAttribute. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(userID string, attrKey string) UserAttribute {
	out, err := d.Read(userID, attrKey)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading userID==%s, attrKey==%s in %s\n", userID, attrKey, d.Name))
	return out
}


// ReadOrEmpty reads UserAttribute
func (d DAOImpl) ReadOrEmpty(userID string, attrKey string) (out []UserAttribute, err error) {
	var outOrEmpty UserAttribute
	ids := idParams(userID, attrKey)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.UserID == userID && outOrEmpty.AttrKey == attrKey {
		out = append(out, outOrEmpty)
	} else if err != nil && strings.HasPrefix(err.Error(), "[NOT FOUND]") {
		err = nil // expected not-found error	
	}
	err = errors.Wrapf(err, "UserAttribute DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the UserAttribute. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(userID string, attrKey string) []UserAttribute {
	out, err := d.ReadOrEmpty(userID, attrKey)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading userID==%s, attrKey==%s in %s\n", userID, attrKey, d.Name))
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
			err = errors.Wrapf(err, "UserAttribute DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			emptyFields, ok := userAttribute.CollectEmptyFields()
			if ok {
				old := olds[0]
				
				key := idParams(old.UserID, old.AttrKey)
				expr, exprAttributes, names := updateExpression(userAttribute, old)
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
				err = errors.Wrapf(err, "UserAttribute DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the UserAttribute regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(userAttribute UserAttribute) {
	err := d.CreateOrUpdate(userAttribute)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", userAttribute, d.Name))
}


// Delete removes UserAttribute from db
func (d DAOImpl)Delete(userID string, attrKey string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(userID, attrKey))
}


// DeleteUnsafe deletes UserAttribute and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(userID string, attrKey string) {
	err := d.Delete(userID, attrKey)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete userID==%s, attrKey==%s in %s\n", userID, attrKey, d.Name))
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
