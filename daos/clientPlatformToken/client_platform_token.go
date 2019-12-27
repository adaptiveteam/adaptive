package clientPlatformToken
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

type PlatformName string

type ClientPlatformToken struct  {
	PlatformID models.PlatformID `json:"platform_id"`
	Org string `json:"org"`
	// should be slack or ms-teams
	PlatformName PlatformName `json:"platform_name"`
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
	if clientPlatformToken.PlatformToken == "" { emptyFields = append(emptyFields, "PlatformToken")}
	if clientPlatformToken.ContactFirstName == "" { emptyFields = append(emptyFields, "ContactFirstName")}
	if clientPlatformToken.ContactLastName == "" { emptyFields = append(emptyFields, "ContactLastName")}
	if clientPlatformToken.ContactMail == "" { emptyFields = append(emptyFields, "ContactMail")}
	ok = len(emptyFields) == 0
	return
}

type DAO interface {
	Create(clientPlatformToken ClientPlatformToken) error
	CreateUnsafe(clientPlatformToken ClientPlatformToken)
	Read(platformID models.PlatformID) (clientPlatformToken ClientPlatformToken, err error)
	ReadUnsafe(platformID models.PlatformID) (clientPlatformToken ClientPlatformToken)
	ReadOrEmpty(platformID models.PlatformID) (clientPlatformToken []ClientPlatformToken, err error)
	ReadOrEmptyUnsafe(platformID models.PlatformID) (clientPlatformToken []ClientPlatformToken)
	CreateOrUpdate(clientPlatformToken ClientPlatformToken) error
	CreateOrUpdateUnsafe(clientPlatformToken ClientPlatformToken)
	Delete(platformID models.PlatformID) error
	DeleteUnsafe(platformID models.PlatformID)
	ReadByPlatformID(platformID models.PlatformID) (clientPlatformToken []ClientPlatformToken, err error)
	ReadByPlatformIDUnsafe(platformID models.PlatformID) (clientPlatformToken []ClientPlatformToken)
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
		Name: clientID + "_client_platform_token",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the ClientPlatformToken.
func (d DAOImpl) Create(clientPlatformToken ClientPlatformToken) error {
	emptyFields, ok := clientPlatformToken.CollectEmptyFields()
	if !ok {return fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)}
	return d.Dynamo.PutTableEntry(clientPlatformToken, d.Name)
}


// CreateUnsafe saves the ClientPlatformToken.
func (d DAOImpl) CreateUnsafe(clientPlatformToken ClientPlatformToken) {
	err := d.Create(clientPlatformToken)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create platformID==%s in %s\n", clientPlatformToken.PlatformID, d.Name))
}


// Read reads ClientPlatformToken
func (d DAOImpl) Read(platformID models.PlatformID) (out ClientPlatformToken, err error) {
	var outs []ClientPlatformToken
	outs, err = d.ReadOrEmpty(platformID)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found platformID==%s in %s\n", platformID, d.Name)
	}
	return
}


// ReadUnsafe reads the ClientPlatformToken. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(platformID models.PlatformID) ClientPlatformToken {
	out, err := d.Read(platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading platformID==%s in %s\n", platformID, d.Name))
	return out
}


// ReadOrEmpty reads ClientPlatformToken
func (d DAOImpl) ReadOrEmpty(platformID models.PlatformID) (out []ClientPlatformToken, err error) {
	var outOrEmpty ClientPlatformToken
	ids := idParams(platformID)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.PlatformID == platformID {
		out = append(out, outOrEmpty)
	}
	err = errors.Wrapf(err, "ClientPlatformToken DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the ClientPlatformToken. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(platformID models.PlatformID) []ClientPlatformToken {
	out, err := d.ReadOrEmpty(platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading platformID==%s in %s\n", platformID, d.Name))
	return out
}


// CreateOrUpdate saves the ClientPlatformToken regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(clientPlatformToken ClientPlatformToken) (err error) {
	
	var olds []ClientPlatformToken
	olds, err = d.ReadOrEmpty(clientPlatformToken.PlatformID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(clientPlatformToken)
			err = errors.Wrapf(err, "ClientPlatformToken DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			old := olds[0]
			ids := idParams(old.PlatformID)
			err = d.Dynamo.UpdateTableEntry(
				allParams(clientPlatformToken, old),
				ids,
				updateExpression(clientPlatformToken, old),
				d.Name,
			)
			err = errors.Wrapf(err, "ClientPlatformToken DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s", ids, d.Name)
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the ClientPlatformToken regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(clientPlatformToken ClientPlatformToken) {
	err := d.CreateOrUpdate(clientPlatformToken)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", clientPlatformToken, d.Name))
}


// Delete removes ClientPlatformToken from db
func (d DAOImpl)Delete(platformID models.PlatformID) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(platformID))
}


// DeleteUnsafe deletes ClientPlatformToken and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(platformID models.PlatformID) {
	err := d.Delete(platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete platformID==%s in %s\n", platformID, d.Name))
}


func (d DAOImpl)ReadByPlatformID(platformID models.PlatformID) (out []ClientPlatformToken, err error) {
	var instances []ClientPlatformToken
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


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID models.PlatformID) (out []ClientPlatformToken) {
	out, err := d.ReadByPlatformID(platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", d.Name))
	return
}

func idParams(platformID models.PlatformID) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"platform_id": common.DynS(string(platformID)),
	}
	return params
}
func allParams(clientPlatformToken ClientPlatformToken, old ClientPlatformToken) (params map[string]*dynamodb.AttributeValue) {
	
		params = map[string]*dynamodb.AttributeValue{}
		if clientPlatformToken.PlatformID != old.PlatformID { params["a0"] = common.DynS(string(clientPlatformToken.PlatformID)) }
		if clientPlatformToken.Org != old.Org { params["a1"] = common.DynS(clientPlatformToken.Org) }
		if clientPlatformToken.PlatformName != old.PlatformName { params["a2"] = common.DynS(string(clientPlatformToken.PlatformName)) }
		if clientPlatformToken.PlatformToken != old.PlatformToken { params["a3"] = common.DynS(clientPlatformToken.PlatformToken) }
		if clientPlatformToken.ContactFirstName != old.ContactFirstName { params["a4"] = common.DynS(clientPlatformToken.ContactFirstName) }
		if clientPlatformToken.ContactLastName != old.ContactLastName { params["a5"] = common.DynS(clientPlatformToken.ContactLastName) }
		if clientPlatformToken.ContactMail != old.ContactMail { params["a6"] = common.DynS(clientPlatformToken.ContactMail) }
	return
}
func updateExpression(clientPlatformToken ClientPlatformToken, old ClientPlatformToken) string {
	var updateParts []string
	
		
			
		if clientPlatformToken.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a0") }
		if clientPlatformToken.Org != old.Org { updateParts = append(updateParts, "org = :a1") }
		if clientPlatformToken.PlatformName != old.PlatformName { updateParts = append(updateParts, "platform_name = :a2") }
		if clientPlatformToken.PlatformToken != old.PlatformToken { updateParts = append(updateParts, "platform_token = :a3") }
		if clientPlatformToken.ContactFirstName != old.ContactFirstName { updateParts = append(updateParts, "contact_first_name = :a4") }
		if clientPlatformToken.ContactLastName != old.ContactLastName { updateParts = append(updateParts, "contact_last_name = :a5") }
		if clientPlatformToken.ContactMail != old.ContactMail { updateParts = append(updateParts, "contact_mail = :a6") }
	return strings.Join(updateParts, " and ")
}
