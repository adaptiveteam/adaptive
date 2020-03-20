package fetch_dialog

import (
	"fmt"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strings"
)

// DAO - wrapper around a Dynamo DB table to work with holidays inside it
type DAO interface {
	FetchByContextSubject(
		context string,
		subject string,
	) (rv DialogEntry, err error)
	FetchByDialogID(dialogID string) (result DialogEntry, found bool, err error)
	FetchByAlias(
		packageName,
		contextAlias,
		subject string,
	) (rv DialogEntry, err error)
	Create(dialogEntry DialogEntry) error
	CreateAlias(aliasEntry ContextAliasEntry) error 
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo *awsutils.DynamoRequest
	Table string
	AliasesTable string
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, table string) DAO {
	return DAOImpl{
		Dynamo: dynamo, 
		Table: table, 
		AliasesTable: table + "_alias",
	}
}
// ConvertContextPathToHash converts normal context to hashes
func ConvertContextPathToHash(context string) string {
	return strings.Replace(strings.Trim(context,"/ "), "/", "#", -1) + "#"
}
// FetchByContextSubject fetches a piece of dialog by context and subject.
func (d DAOImpl)FetchByContextSubject(
	context string,
	subject string,
) (rv DialogEntry, err error) {
	contextHash := ConvertContextPathToHash(context)

	result := make([]DialogEntry, 0)
	err = d.Dynamo.QueryTableWithIndex(
		d.Table,
		awsutils.DynamoIndexExpression{
			IndexName: "ContextSubjectIndex",
			Condition: "context = :c and subject = :s",
			Attributes: map[string]interface{}{
				":c": contextHash,
				":s": subject,
			},
		}, map[string]string{}, true, -1, &result)
		
	if len(result) == 1 {
		rv = result[0]
	} else if err == nil {
		err = fmt.Errorf("FetchByContextSubject(context=%s, subject=%s): expected one result but got %v", context, subject, len(result))
	}

	return rv, err
}

// FetchByDialogID fetches a piece of dialog using a unique UUID associated with the dialog
func (d DAOImpl)FetchByDialogID(dialogID string) (result DialogEntry, found bool, err error) {
	params := map[string]*dynamodb.AttributeValue{
		"dialog_id": daosCommon.DynS(dialogID),
	}
	found, err = d.Dynamo.GetItemOrEmptyFromTable(d.Table, params, &result)
	return
}

// FetchByAlias fetches a piece of dialog using a application/package ID, context alias, and subject
// https://github.com/adaptiveteam/dialog-library/tree/cultivate/aliases
func (d DAOImpl)FetchByAlias(
	packageName,
	contextAlias,
	subject string,
) (rv DialogEntry, err error) {
	var contextAliasEntry ContextAliasEntry
	applicationAlias := packageName+"#"+contextAlias
	params := map[string]*dynamodb.AttributeValue{
		"application_alias": daosCommon.DynS(applicationAlias),
	}

	err = d.Dynamo.GetItemFromTable(d.AliasesTable, params, &contextAliasEntry)

	if err == nil {
		rv, err = d.FetchByContextSubject(contextAliasEntry.Context, subject)
	}
	return rv, err
}
// Create s a new item in the dialog table
func (d DAOImpl)Create(dialogEntry DialogEntry) error {
	return d.Dynamo.PutTableEntry(dialogEntry, d.Table)
}
// CreateAlias creates a new alias in the aliases table
func (d DAOImpl)CreateAlias(aliasEntry ContextAliasEntry) error {
	return d.Dynamo.PutTableEntry(aliasEntry, d.AliasesTable)
}

// FetchDialogImpl a wrapper for a function that implements FetchDialog interface
type FetchDialogImpl struct {
	FetchDialogFunc func(subject string) (dialog DialogEntry, err error)
}
// FetchDialog implementation
func (i FetchDialogImpl)FetchDialog(subject string) (dialog DialogEntry, err error) {
	return i.FetchDialogFunc(subject)
}
