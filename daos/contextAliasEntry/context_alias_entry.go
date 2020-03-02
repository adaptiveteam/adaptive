package contextAliasEntry
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

// ContextAliasEntry contains all of the information needed for a context alias
// A context alias is a way to alias  a piece of context without spelling out
// the context path.  If the path changes you can still safely use the alias.
type ContextAliasEntry struct  {
	ApplicationAlias string `json:"application_alias"`
	Context string `json:"context"`
	BuildID string `json:"build_id"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (contextAliasEntry ContextAliasEntry)CollectEmptyFields() (emptyFields []string, ok bool) {
	if contextAliasEntry.ApplicationAlias == "" { emptyFields = append(emptyFields, "ApplicationAlias")}
	if contextAliasEntry.Context == "" { emptyFields = append(emptyFields, "Context")}
	if contextAliasEntry.BuildID == "" { emptyFields = append(emptyFields, "BuildID")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (contextAliasEntry ContextAliasEntry) ToJSON() (string, error) {
	b, err := json.Marshal(contextAliasEntry)
	return string(b), err
}

type DAO interface {
	Create(contextAliasEntry ContextAliasEntry) error
	CreateUnsafe(contextAliasEntry ContextAliasEntry)
	Read(applicationAlias string) (contextAliasEntry ContextAliasEntry, err error)
	ReadUnsafe(applicationAlias string) (contextAliasEntry ContextAliasEntry)
	ReadOrEmpty(applicationAlias string) (contextAliasEntry []ContextAliasEntry, err error)
	ReadOrEmptyUnsafe(applicationAlias string) (contextAliasEntry []ContextAliasEntry)
	CreateOrUpdate(contextAliasEntry ContextAliasEntry) error
	CreateOrUpdateUnsafe(contextAliasEntry ContextAliasEntry)
	Delete(applicationAlias string) error
	DeleteUnsafe(applicationAlias string)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	Name      string                  `json:"name"`
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic(errors.New("Cannot create ContextAliasEntry.DAO without clientID")) }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: TableName(clientID),
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic(errors.New("Cannot create ContextAliasEntry.DAO without tableName")) }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}
func TableName(clientID string) string {
	return clientID + "_context_alias_entry"
}

// Create saves the ContextAliasEntry.
func (d DAOImpl) Create(contextAliasEntry ContextAliasEntry) (err error) {
	emptyFields, ok := contextAliasEntry.CollectEmptyFields()
	if ok {
		err = d.Dynamo.PutTableEntry(contextAliasEntry, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the ContextAliasEntry.
func (d DAOImpl) CreateUnsafe(contextAliasEntry ContextAliasEntry) {
	err2 := d.Create(contextAliasEntry)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not create applicationAlias==%s in %s\n", contextAliasEntry.ApplicationAlias, d.Name))
}


// Read reads ContextAliasEntry
func (d DAOImpl) Read(applicationAlias string) (out ContextAliasEntry, err error) {
	var outs []ContextAliasEntry
	outs, err = d.ReadOrEmpty(applicationAlias)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found applicationAlias==%s in %s\n", applicationAlias, d.Name)
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the ContextAliasEntry. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(applicationAlias string) ContextAliasEntry {
	out, err2 := d.Read(applicationAlias)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Error reading applicationAlias==%s in %s\n", applicationAlias, d.Name))
	return out
}


// ReadOrEmpty reads ContextAliasEntry
func (d DAOImpl) ReadOrEmpty(applicationAlias string) (out []ContextAliasEntry, err error) {
	var outOrEmpty ContextAliasEntry
	ids := idParams(applicationAlias)
	var found bool
	found, err = d.Dynamo.GetItemOrEmptyFromTable(d.Name, ids, &outOrEmpty)
	if found {
		if outOrEmpty.ApplicationAlias == applicationAlias {
			out = append(out, outOrEmpty)
		} else {
			err = fmt.Errorf("Requested ids: applicationAlias==%s are different from the found ones: applicationAlias==%s", applicationAlias, outOrEmpty.ApplicationAlias) // unexpected error: found ids != ids
		}
	}
	err = errors.Wrapf(err, "ContextAliasEntry DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the ContextAliasEntry. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(applicationAlias string) []ContextAliasEntry {
	out, err2 := d.ReadOrEmpty(applicationAlias)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Error while reading applicationAlias==%s in %s\n", applicationAlias, d.Name))
	return out
}


// CreateOrUpdate saves the ContextAliasEntry regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(contextAliasEntry ContextAliasEntry) (err error) {
	
	var olds []ContextAliasEntry
	olds, err = d.ReadOrEmpty(contextAliasEntry.ApplicationAlias)
	err = errors.Wrapf(err, "ContextAliasEntry DAO.CreateOrUpdate(id = applicationAlias==%s) couldn't ReadOrEmpty", contextAliasEntry.ApplicationAlias)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(contextAliasEntry)
			err = errors.Wrapf(err, "ContextAliasEntry DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			emptyFields, ok := contextAliasEntry.CollectEmptyFields()
			if ok {
				old := olds[0]
				
				
				key := idParams(old.ApplicationAlias)
				expr, exprAttributes, names := updateExpression(contextAliasEntry, old)
				input := dynamodb.UpdateItemInput{
					ExpressionAttributeValues: exprAttributes,
					TableName:                 aws.String(d.Name),
					Key:                       key,
					ReturnValues:              aws.String("UPDATED_NEW"),
					UpdateExpression:          aws.String(expr),
				}
				if names != nil { input.ExpressionAttributeNames = *names } // workaround for a pointer to an empty slice
				if  len(exprAttributes) > 0 { // if there some changes
					err = d.Dynamo.UpdateItemInternal(input)
				} else {
					// WARN: no changes.
				}
				err = errors.Wrapf(err, "ContextAliasEntry DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the ContextAliasEntry regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(contextAliasEntry ContextAliasEntry) {
	err2 := d.CreateOrUpdate(contextAliasEntry)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", contextAliasEntry, d.Name))
}


// Delete removes ContextAliasEntry from db
func (d DAOImpl)Delete(applicationAlias string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(applicationAlias))
}


// DeleteUnsafe deletes ContextAliasEntry and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(applicationAlias string) {
	err2 := d.Delete(applicationAlias)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not delete applicationAlias==%s in %s\n", applicationAlias, d.Name))
}

func idParams(applicationAlias string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"application_alias": common.DynS(applicationAlias),
	}
	return params
}
func allParams(contextAliasEntry ContextAliasEntry, old ContextAliasEntry) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if contextAliasEntry.ApplicationAlias != old.ApplicationAlias { params[":a0"] = common.DynS(contextAliasEntry.ApplicationAlias) }
	if contextAliasEntry.Context != old.Context { params[":a1"] = common.DynS(contextAliasEntry.Context) }
	if contextAliasEntry.BuildID != old.BuildID { params[":a2"] = common.DynS(contextAliasEntry.BuildID) }
	return
}
func updateExpression(contextAliasEntry ContextAliasEntry, old ContextAliasEntry) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if contextAliasEntry.ApplicationAlias != old.ApplicationAlias { updateParts = append(updateParts, "application_alias = :a0"); params[":a0"] = common.DynS(contextAliasEntry.ApplicationAlias);  }
	if contextAliasEntry.Context != old.Context { updateParts = append(updateParts, "context = :a1"); params[":a1"] = common.DynS(contextAliasEntry.Context);  }
	if contextAliasEntry.BuildID != old.BuildID { updateParts = append(updateParts, "build_id = :a2"); params[":a2"] = common.DynS(contextAliasEntry.BuildID);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
