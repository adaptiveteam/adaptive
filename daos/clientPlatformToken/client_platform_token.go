package clientPlatformToken
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

type ClientPlatformToken struct  {
	PlatformID common.PlatformID `json:"platform_id"`
	Org string `json:"org"`
	// should be slack or ms-teams
	PlatformName common.PlatformName `json:"platform_name"`
	PlatformToken string `json:"platform_token"`
	ContactFirstName string `json:"contact_first_name"`
	ContactLastName string `json:"contact_last_name"`
	ContactMail string `json:"contact_mail"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (clientPlatformToken ClientPlatformToken)CollectEmptyFields() (emptyFields []string, ok bool) {
	if clientPlatformToken.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if clientPlatformToken.Org == "" { emptyFields = append(emptyFields, "Org")}
	if clientPlatformToken.PlatformName == "" { emptyFields = append(emptyFields, "PlatformName")}
	if clientPlatformToken.PlatformToken == "" { emptyFields = append(emptyFields, "PlatformToken")}
	if clientPlatformToken.ContactFirstName == "" { emptyFields = append(emptyFields, "ContactFirstName")}
	if clientPlatformToken.ContactLastName == "" { emptyFields = append(emptyFields, "ContactLastName")}
	if clientPlatformToken.ContactMail == "" { emptyFields = append(emptyFields, "ContactMail")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (clientPlatformToken ClientPlatformToken) ToJSON() (string, error) {
	b, err := json.Marshal(clientPlatformToken)
	return string(b), err
}

type DAO interface {
	Create(clientPlatformToken ClientPlatformToken) error
	CreateUnsafe(clientPlatformToken ClientPlatformToken)
	Read(platformID common.PlatformID) (clientPlatformToken ClientPlatformToken, err error)
	ReadUnsafe(platformID common.PlatformID) (clientPlatformToken ClientPlatformToken)
	ReadOrEmpty(platformID common.PlatformID) (clientPlatformToken []ClientPlatformToken, err error)
	ReadOrEmptyUnsafe(platformID common.PlatformID) (clientPlatformToken []ClientPlatformToken)
	CreateOrUpdate(clientPlatformToken ClientPlatformToken) error
	CreateOrUpdateUnsafe(clientPlatformToken ClientPlatformToken)
	Delete(platformID common.PlatformID) error
	DeleteUnsafe(platformID common.PlatformID)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	ConnGen   common.DynamoDBConnectionGen
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic(errors.New("Cannot create ClientPlatformToken.DAO without clientID")) }
	return DAOImpl{
		ConnGen:   common.DynamoDBConnectionGen{
			Dynamo: dynamo, 
			TableNamePrefix: clientID,
		},
	}
}

// // NewDAOByTableName creates an instance of DAO that will provide access to the table
// func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
// 	if tableName == "" { panic(errors.New("Cannot create ClientPlatformToken.DAO without tableName")) }
// 	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
// 		Name: tableName,
// 	}
// }
// TableNameSuffixVar is a global variable that contains table name suffix.
// After renaming all tables this may be made `const`.
var TableNameSuffixVar = "_client_platform_token"

// TableName concatenates table name prefix and suffix and returns table name
func TableName(prefix string) string {
	return prefix + TableNameSuffixVar
}

// Create saves the ClientPlatformToken.
func (d DAOImpl) Create(clientPlatformToken ClientPlatformToken) (err error) {
	emptyFields, ok := clientPlatformToken.CollectEmptyFields()
	if ok {
		err = d.ConnGen.Dynamo.PutTableEntry(clientPlatformToken, TableName(d.ConnGen.TableNamePrefix))
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the ClientPlatformToken.
func (d DAOImpl) CreateUnsafe(clientPlatformToken ClientPlatformToken) {
	err2 := d.Create(clientPlatformToken)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not create platformID==%s in %s\n", clientPlatformToken.PlatformID, TableName(d.ConnGen.TableNamePrefix)))
}


// Read reads ClientPlatformToken
func (d DAOImpl) Read(platformID common.PlatformID) (out ClientPlatformToken, err error) {
	var outs []ClientPlatformToken
	outs, err = d.ReadOrEmpty(platformID)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found platformID==%s in %s\n", platformID, TableName(d.ConnGen.TableNamePrefix))
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the ClientPlatformToken. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(platformID common.PlatformID) ClientPlatformToken {
	out, err2 := d.Read(platformID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error reading platformID==%s in %s\n", platformID, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// ReadOrEmpty reads ClientPlatformToken
func (d DAOImpl) ReadOrEmpty(platformID common.PlatformID) (out []ClientPlatformToken, err error) {
	var outOrEmpty ClientPlatformToken
	ids := idParams(platformID)
	var found bool
	found, err = d.ConnGen.Dynamo.GetItemOrEmptyFromTable(TableName(d.ConnGen.TableNamePrefix), ids, &outOrEmpty)
	if found {
		if outOrEmpty.PlatformID == platformID {
			out = append(out, outOrEmpty)
		} else {
			err = fmt.Errorf("Requested ids: platformID==%s are different from the found ones: platformID==%s", platformID, outOrEmpty.PlatformID) // unexpected error: found ids != ids
		}
	}
	err = errors.Wrapf(err, "ClientPlatformToken DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(d.ConnGen.TableNamePrefix))
	return
}


// ReadOrEmptyUnsafe reads the ClientPlatformToken. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(platformID common.PlatformID) []ClientPlatformToken {
	out, err2 := d.ReadOrEmpty(platformID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error while reading platformID==%s in %s\n", platformID, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// CreateOrUpdate saves the ClientPlatformToken regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(clientPlatformToken ClientPlatformToken) (err error) {
	
	var olds []ClientPlatformToken
	olds, err = d.ReadOrEmpty(clientPlatformToken.PlatformID)
	err = errors.Wrapf(err, "ClientPlatformToken DAO.CreateOrUpdate(id = platformID==%s) couldn't ReadOrEmpty", clientPlatformToken.PlatformID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(clientPlatformToken)
			err = errors.Wrapf(err, "ClientPlatformToken DAO.CreateOrUpdate couldn't Create in table %s", TableName(d.ConnGen.TableNamePrefix))
		} else {
			emptyFields, ok := clientPlatformToken.CollectEmptyFields()
			if ok {
				old := olds[0]
				
				
				key := idParams(old.PlatformID)
				expr, exprAttributes, names := updateExpression(clientPlatformToken, old)
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
				err = errors.Wrapf(err, "ClientPlatformToken DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(d.ConnGen.TableNamePrefix), expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the ClientPlatformToken regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(clientPlatformToken ClientPlatformToken) {
	err2 := d.CreateOrUpdate(clientPlatformToken)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("could not create or update %v in %s\n", clientPlatformToken, TableName(d.ConnGen.TableNamePrefix)))
}


// Delete removes ClientPlatformToken from db
func (d DAOImpl)Delete(platformID common.PlatformID) error {
	return d.ConnGen.Dynamo.DeleteEntry(TableName(d.ConnGen.TableNamePrefix), idParams(platformID))
}


// DeleteUnsafe deletes ClientPlatformToken and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(platformID common.PlatformID) {
	err2 := d.Delete(platformID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not delete platformID==%s in %s\n", platformID, TableName(d.ConnGen.TableNamePrefix)))
}

func idParams(platformID common.PlatformID) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"platform_id": common.DynS(string(platformID)),
	}
	return params
}
func allParams(clientPlatformToken ClientPlatformToken, old ClientPlatformToken) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if clientPlatformToken.PlatformID != old.PlatformID { params[":a0"] = common.DynS(string(clientPlatformToken.PlatformID)) }
	if clientPlatformToken.Org != old.Org { params[":a1"] = common.DynS(clientPlatformToken.Org) }
	if clientPlatformToken.PlatformName != old.PlatformName { params[":a2"] = common.DynS(string(clientPlatformToken.PlatformName)) }
	if clientPlatformToken.PlatformToken != old.PlatformToken { params[":a3"] = common.DynS(clientPlatformToken.PlatformToken) }
	if clientPlatformToken.ContactFirstName != old.ContactFirstName { params[":a4"] = common.DynS(clientPlatformToken.ContactFirstName) }
	if clientPlatformToken.ContactLastName != old.ContactLastName { params[":a5"] = common.DynS(clientPlatformToken.ContactLastName) }
	if clientPlatformToken.ContactMail != old.ContactMail { params[":a6"] = common.DynS(clientPlatformToken.ContactMail) }
	return
}
func updateExpression(clientPlatformToken ClientPlatformToken, old ClientPlatformToken) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if clientPlatformToken.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a0"); params[":a0"] = common.DynS(string(clientPlatformToken.PlatformID));  }
	if clientPlatformToken.Org != old.Org { updateParts = append(updateParts, "org = :a1"); params[":a1"] = common.DynS(clientPlatformToken.Org);  }
	if clientPlatformToken.PlatformName != old.PlatformName { updateParts = append(updateParts, "platform_name = :a2"); params[":a2"] = common.DynS(string(clientPlatformToken.PlatformName));  }
	if clientPlatformToken.PlatformToken != old.PlatformToken { updateParts = append(updateParts, "platform_token = :a3"); params[":a3"] = common.DynS(clientPlatformToken.PlatformToken);  }
	if clientPlatformToken.ContactFirstName != old.ContactFirstName { updateParts = append(updateParts, "contact_first_name = :a4"); params[":a4"] = common.DynS(clientPlatformToken.ContactFirstName);  }
	if clientPlatformToken.ContactLastName != old.ContactLastName { updateParts = append(updateParts, "contact_last_name = :a5"); params[":a5"] = common.DynS(clientPlatformToken.ContactLastName);  }
	if clientPlatformToken.ContactMail != old.ContactMail { updateParts = append(updateParts, "contact_mail = :a6"); params[":a6"] = common.DynS(clientPlatformToken.ContactMail);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
