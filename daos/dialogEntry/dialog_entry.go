package dialogEntry
// This file has been automatically generated by `adaptive-platform/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/aws-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

// DialogEntry stores all of the  relevant information for a piece of dialog including:
type DialogEntry struct  {
	// This is an immutable UUID that developers can use
	DialogID string `json:"dialog_id"`
	// This is the context path for the piece of dialog
	Context string `json:"context"`
	// This is the dialog subject
	Subject string `json:"subject"`
	// This was when the dialog was last updated
	Updated string `json:"updated"`
	// These are the dialog options
	Dialog []string `json:"dialog"`
	// Comments to help cultivators understand the dialog intent
	Comments []string `json:"comments"`
	// This the link to the LearnMore page
	LearnMoreLink string `json:"learn_more_link"`
	// This is the actual content from the LearnMore page
	LearnMoreContent string `json:"learn_more_content"`
	BuildBranch string `json:"build_branch"`
	CultivationBranch string `json:"cultivation_branch"`
	MasterBranch string `json:"master_branch"`
	BuildID string `json:"build_id"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (dialogEntry DialogEntry)CollectEmptyFields() (emptyFields []string, ok bool) {
	if dialogEntry.DialogID == "" { emptyFields = append(emptyFields, "DialogID")}
	if dialogEntry.Context == "" { emptyFields = append(emptyFields, "Context")}
	if dialogEntry.Subject == "" { emptyFields = append(emptyFields, "Subject")}
	if dialogEntry.Updated == "" { emptyFields = append(emptyFields, "Updated")}
	if dialogEntry.Dialog == nil { emptyFields = append(emptyFields, "Dialog")}
	if dialogEntry.Comments == nil { emptyFields = append(emptyFields, "Comments")}
	if dialogEntry.LearnMoreLink == "" { emptyFields = append(emptyFields, "LearnMoreLink")}
	if dialogEntry.LearnMoreContent == "" { emptyFields = append(emptyFields, "LearnMoreContent")}
	if dialogEntry.BuildBranch == "" { emptyFields = append(emptyFields, "BuildBranch")}
	if dialogEntry.CultivationBranch == "" { emptyFields = append(emptyFields, "CultivationBranch")}
	if dialogEntry.MasterBranch == "" { emptyFields = append(emptyFields, "MasterBranch")}
	if dialogEntry.BuildID == "" { emptyFields = append(emptyFields, "BuildID")}
	ok = len(emptyFields) == 0
	return
}

type DAO interface {
	Create(dialogEntry DialogEntry) error
	CreateUnsafe(dialogEntry DialogEntry)
	Read(dialogID string) (dialogEntry DialogEntry, err error)
	ReadUnsafe(dialogID string) (dialogEntry DialogEntry)
	ReadOrEmpty(dialogID string) (dialogEntry []DialogEntry, err error)
	ReadOrEmptyUnsafe(dialogID string) (dialogEntry []DialogEntry)
	CreateOrUpdate(dialogEntry DialogEntry) error
	CreateOrUpdateUnsafe(dialogEntry DialogEntry)
	Delete(dialogID string) error
	DeleteUnsafe(dialogID string)
	ReadByContextSubject(context string, subject string) (dialogEntry []DialogEntry, err error)
	ReadByContextSubjectUnsafe(context string, subject string) (dialogEntry []DialogEntry)
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
		Name: clientID + "_dialog_entry",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the DialogEntry.
func (d DAOImpl) Create(dialogEntry DialogEntry) error {
	emptyFields, ok := dialogEntry.CollectEmptyFields()
	if !ok {return fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)}
	return d.Dynamo.PutTableEntry(dialogEntry, d.Name)
}


// CreateUnsafe saves the DialogEntry.
func (d DAOImpl) CreateUnsafe(dialogEntry DialogEntry) {
	err := d.Create(dialogEntry)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create dialogID==%s in %s\n", dialogEntry.DialogID, d.Name))
}


// Read reads DialogEntry
func (d DAOImpl) Read(dialogID string) (out DialogEntry, err error) {
	var outs []DialogEntry
	outs, err = d.ReadOrEmpty(dialogID)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found dialogID==%s in %s\n", dialogID, d.Name)
	}
	return
}


// ReadUnsafe reads the DialogEntry. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(dialogID string) DialogEntry {
	out, err := d.Read(dialogID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading dialogID==%s in %s\n", dialogID, d.Name))
	return out
}


// ReadOrEmpty reads DialogEntry
func (d DAOImpl) ReadOrEmpty(dialogID string) (out []DialogEntry, err error) {
	var outOrEmpty DialogEntry
	ids := idParams(dialogID)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.DialogID == dialogID {
		out = append(out, outOrEmpty)
	} else if err != nil && strings.HasPrefix(err.Error(), "In table ") {
		err = nil // expected not-found error	
	}
	err = errors.Wrapf(err, "DialogEntry DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the DialogEntry. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(dialogID string) []DialogEntry {
	out, err := d.ReadOrEmpty(dialogID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading dialogID==%s in %s\n", dialogID, d.Name))
	return out
}


// CreateOrUpdate saves the DialogEntry regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(dialogEntry DialogEntry) (err error) {
	
	var olds []DialogEntry
	olds, err = d.ReadOrEmpty(dialogEntry.DialogID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(dialogEntry)
			err = errors.Wrapf(err, "DialogEntry DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			old := olds[0]
			
			key := idParams(old.DialogID)
			expr, exprAttributes, names := updateExpression(dialogEntry, old)
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
			err = errors.Wrapf(err, "DialogEntry DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s", key, d.Name)
			return
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the DialogEntry regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(dialogEntry DialogEntry) {
	err := d.CreateOrUpdate(dialogEntry)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", dialogEntry, d.Name))
}


// Delete removes DialogEntry from db
func (d DAOImpl)Delete(dialogID string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(dialogID))
}


// DeleteUnsafe deletes DialogEntry and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(dialogID string) {
	err := d.Delete(dialogID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete dialogID==%s in %s\n", dialogID, d.Name))
}


func (d DAOImpl)ReadByContextSubject(context string, subject string) (out []DialogEntry, err error) {
	var instances []DialogEntry
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "ContextSubjectIndex",
		Condition: "context = :a0 and subject = :a1",
		Attributes: map[string]interface{}{
			":a0": context,
			":a1": subject,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByContextSubjectUnsafe(context string, subject string) (out []DialogEntry) {
	out, err := d.ReadByContextSubject(context, subject)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query ContextSubjectIndex on %s table\n", d.Name))
	return
}

func idParams(dialogID string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"dialog_id": common.DynS(dialogID),
	}
	return params
}
func allParams(dialogEntry DialogEntry, old DialogEntry) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if dialogEntry.DialogID != old.DialogID { params[":a0"] = common.DynS(dialogEntry.DialogID) }
	if dialogEntry.Context != old.Context { params[":a1"] = common.DynS(dialogEntry.Context) }
	if dialogEntry.Subject != old.Subject { params[":a2"] = common.DynS(dialogEntry.Subject) }
	if dialogEntry.Updated != old.Updated { params[":a3"] = common.DynS(dialogEntry.Updated) }
	if !common.StringArraysEqual(dialogEntry.Dialog, old.Dialog) { params[":a4"] = common.DynSS(dialogEntry.Dialog) }
	if !common.StringArraysEqual(dialogEntry.Comments, old.Comments) { params[":a5"] = common.DynSS(dialogEntry.Comments) }
	if dialogEntry.LearnMoreLink != old.LearnMoreLink { params[":a6"] = common.DynS(dialogEntry.LearnMoreLink) }
	if dialogEntry.LearnMoreContent != old.LearnMoreContent { params[":a7"] = common.DynS(dialogEntry.LearnMoreContent) }
	if dialogEntry.BuildBranch != old.BuildBranch { params[":a8"] = common.DynS(dialogEntry.BuildBranch) }
	if dialogEntry.CultivationBranch != old.CultivationBranch { params[":a9"] = common.DynS(dialogEntry.CultivationBranch) }
	if dialogEntry.MasterBranch != old.MasterBranch { params[":a10"] = common.DynS(dialogEntry.MasterBranch) }
	if dialogEntry.BuildID != old.BuildID { params[":a11"] = common.DynS(dialogEntry.BuildID) }
	return
}
func updateExpression(dialogEntry DialogEntry, old DialogEntry) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if dialogEntry.DialogID != old.DialogID { updateParts = append(updateParts, "dialog_id = :a0"); params[":a0"] = common.DynS(dialogEntry.DialogID);  }
	if dialogEntry.Context != old.Context { updateParts = append(updateParts, "context = :a1"); params[":a1"] = common.DynS(dialogEntry.Context);  }
	if dialogEntry.Subject != old.Subject { updateParts = append(updateParts, "subject = :a2"); params[":a2"] = common.DynS(dialogEntry.Subject);  }
	if dialogEntry.Updated != old.Updated { updateParts = append(updateParts, "updated = :a3"); params[":a3"] = common.DynS(dialogEntry.Updated);  }
	if !common.StringArraysEqual(dialogEntry.Dialog, old.Dialog) { updateParts = append(updateParts, "dialog = :a4"); params[":a4"] = common.DynSS(dialogEntry.Dialog);  }
	if !common.StringArraysEqual(dialogEntry.Comments, old.Comments) { updateParts = append(updateParts, "comments = :a5"); params[":a5"] = common.DynSS(dialogEntry.Comments);  }
	if dialogEntry.LearnMoreLink != old.LearnMoreLink { updateParts = append(updateParts, "learn_more_link = :a6"); params[":a6"] = common.DynS(dialogEntry.LearnMoreLink);  }
	if dialogEntry.LearnMoreContent != old.LearnMoreContent { updateParts = append(updateParts, "learn_more_content = :a7"); params[":a7"] = common.DynS(dialogEntry.LearnMoreContent);  }
	if dialogEntry.BuildBranch != old.BuildBranch { updateParts = append(updateParts, "build_branch = :a8"); params[":a8"] = common.DynS(dialogEntry.BuildBranch);  }
	if dialogEntry.CultivationBranch != old.CultivationBranch { updateParts = append(updateParts, "cultivation_branch = :a9"); params[":a9"] = common.DynS(dialogEntry.CultivationBranch);  }
	if dialogEntry.MasterBranch != old.MasterBranch { updateParts = append(updateParts, "master_branch = :a10"); params[":a10"] = common.DynS(dialogEntry.MasterBranch);  }
	if dialogEntry.BuildID != old.BuildID { updateParts = append(updateParts, "build_id = :a11"); params[":a11"] = common.DynS(dialogEntry.BuildID);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
