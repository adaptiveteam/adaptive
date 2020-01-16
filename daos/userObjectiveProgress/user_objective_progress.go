package userObjectiveProgress
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

type UserObjectiveProgress struct  {
	ID string `json:"id"`
	CreatedOn string `json:"created_on"`
	PlatformID common.PlatformID `json:"platform_id"`
	UserID string `json:"user_id"`
	PartnerID string `json:"partner_id"`
	Comments string `json:"comments"`
	// 1 for true, 0 for false
	Closeout int `json:"closeout"`
	PercentTimeLapsed string `json:"percent_time_lapsed"`
	StatusColor common.ObjectiveStatusColor `json:"status_color"`
	ReviewedByPartner bool `json:"reviewed_by_partner"`
	PartnerComments string `json:"partner_comments,omitempty"`
	PartnerReportedProgress string `json:"partner_reported_progress,omitempty"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (userObjectiveProgress UserObjectiveProgress)CollectEmptyFields() (emptyFields []string, ok bool) {
	if userObjectiveProgress.ID == "" { emptyFields = append(emptyFields, "ID")}
	if userObjectiveProgress.CreatedOn == "" { emptyFields = append(emptyFields, "CreatedOn")}
	if userObjectiveProgress.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if userObjectiveProgress.UserID == "" { emptyFields = append(emptyFields, "UserID")}
	if userObjectiveProgress.PartnerID == "" { emptyFields = append(emptyFields, "PartnerID")}
	if userObjectiveProgress.Comments == "" { emptyFields = append(emptyFields, "Comments")}
	if userObjectiveProgress.PercentTimeLapsed == "" { emptyFields = append(emptyFields, "PercentTimeLapsed")}
	if userObjectiveProgress.StatusColor == "" { emptyFields = append(emptyFields, "StatusColor")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (userObjectiveProgress UserObjectiveProgress) ToJSON() (string, error) {
	b, err := json.Marshal(userObjectiveProgress)
	return string(b), err
}

type DAO interface {
	Create(userObjectiveProgress UserObjectiveProgress) error
	CreateUnsafe(userObjectiveProgress UserObjectiveProgress)
	Read(id string, createdOn string) (userObjectiveProgress UserObjectiveProgress, err error)
	ReadUnsafe(id string, createdOn string) (userObjectiveProgress UserObjectiveProgress)
	ReadOrEmpty(id string, createdOn string) (userObjectiveProgress []UserObjectiveProgress, err error)
	ReadOrEmptyUnsafe(id string, createdOn string) (userObjectiveProgress []UserObjectiveProgress)
	CreateOrUpdate(userObjectiveProgress UserObjectiveProgress) error
	CreateOrUpdateUnsafe(userObjectiveProgress UserObjectiveProgress)
	Delete(id string, createdOn string) error
	DeleteUnsafe(id string, createdOn string)
	ReadByID(id string) (userObjectiveProgress []UserObjectiveProgress, err error)
	ReadByIDUnsafe(id string) (userObjectiveProgress []UserObjectiveProgress)
	ReadByCreatedOn(createdOn string) (userObjectiveProgress []UserObjectiveProgress, err error)
	ReadByCreatedOnUnsafe(createdOn string) (userObjectiveProgress []UserObjectiveProgress)
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
	return clientID + "_user_objective_progress"
}

// Create saves the UserObjectiveProgress.
func (d DAOImpl) Create(userObjectiveProgress UserObjectiveProgress) (err error) {
	emptyFields, ok := userObjectiveProgress.CollectEmptyFields()
	if ok {
		err = d.Dynamo.PutTableEntry(userObjectiveProgress, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the UserObjectiveProgress.
func (d DAOImpl) CreateUnsafe(userObjectiveProgress UserObjectiveProgress) {
	err := d.Create(userObjectiveProgress)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create id==%s, createdOn==%s in %s\n", userObjectiveProgress.ID, userObjectiveProgress.CreatedOn, d.Name))
}


// Read reads UserObjectiveProgress
func (d DAOImpl) Read(id string, createdOn string) (out UserObjectiveProgress, err error) {
	var outs []UserObjectiveProgress
	outs, err = d.ReadOrEmpty(id, createdOn)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s, createdOn==%s in %s\n", id, createdOn, d.Name)
	}
	return
}


// ReadUnsafe reads the UserObjectiveProgress. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string, createdOn string) UserObjectiveProgress {
	out, err := d.Read(id, createdOn)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading id==%s, createdOn==%s in %s\n", id, createdOn, d.Name))
	return out
}


// ReadOrEmpty reads UserObjectiveProgress
func (d DAOImpl) ReadOrEmpty(id string, createdOn string) (out []UserObjectiveProgress, err error) {
	var outOrEmpty UserObjectiveProgress
	ids := idParams(id, createdOn)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.ID == id && outOrEmpty.CreatedOn == createdOn {
		out = append(out, outOrEmpty)
	} else if err != nil && strings.HasPrefix(err.Error(), "[NOT FOUND]") {
		err = nil // expected not-found error	
	}
	err = errors.Wrapf(err, "UserObjectiveProgress DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the UserObjectiveProgress. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string, createdOn string) []UserObjectiveProgress {
	out, err := d.ReadOrEmpty(id, createdOn)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading id==%s, createdOn==%s in %s\n", id, createdOn, d.Name))
	return out
}


// CreateOrUpdate saves the UserObjectiveProgress regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(userObjectiveProgress UserObjectiveProgress) (err error) {
	
	var olds []UserObjectiveProgress
	olds, err = d.ReadOrEmpty(userObjectiveProgress.ID, userObjectiveProgress.CreatedOn)
	err = errors.Wrapf(err, "UserObjectiveProgress DAO.CreateOrUpdate(id = id==%s, createdOn==%s) couldn't ReadOrEmpty", userObjectiveProgress.ID, userObjectiveProgress.CreatedOn)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(userObjectiveProgress)
			err = errors.Wrapf(err, "UserObjectiveProgress DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			emptyFields, ok := userObjectiveProgress.CollectEmptyFields()
			if ok {
				old := olds[0]
				
				key := idParams(old.ID, old.CreatedOn)
				expr, exprAttributes, names := updateExpression(userObjectiveProgress, old)
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
				err = errors.Wrapf(err, "UserObjectiveProgress DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the UserObjectiveProgress regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(userObjectiveProgress UserObjectiveProgress) {
	err := d.CreateOrUpdate(userObjectiveProgress)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", userObjectiveProgress, d.Name))
}


// Delete removes UserObjectiveProgress from db
func (d DAOImpl)Delete(id string, createdOn string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(id, createdOn))
}


// DeleteUnsafe deletes UserObjectiveProgress and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(id string, createdOn string) {
	err := d.Delete(id, createdOn)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete id==%s, createdOn==%s in %s\n", id, createdOn, d.Name))
}


func (d DAOImpl)ReadByID(id string) (out []UserObjectiveProgress, err error) {
	var instances []UserObjectiveProgress
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


func (d DAOImpl)ReadByIDUnsafe(id string) (out []UserObjectiveProgress) {
	out, err := d.ReadByID(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query IDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByCreatedOn(createdOn string) (out []UserObjectiveProgress, err error) {
	var instances []UserObjectiveProgress
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "CreatedOnIndex",
		Condition: "created_on = :a0",
		Attributes: map[string]interface{}{
			":a0": createdOn,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByCreatedOnUnsafe(createdOn string) (out []UserObjectiveProgress) {
	out, err := d.ReadByCreatedOn(createdOn)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query CreatedOnIndex on %s table\n", d.Name))
	return
}

func idParams(id string, createdOn string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"id": common.DynS(id),
		"created_on": common.DynS(createdOn),
	}
	return params
}
func allParams(userObjectiveProgress UserObjectiveProgress, old UserObjectiveProgress) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if userObjectiveProgress.ID != old.ID { params[":a0"] = common.DynS(userObjectiveProgress.ID) }
	if userObjectiveProgress.CreatedOn != old.CreatedOn { params[":a1"] = common.DynS(userObjectiveProgress.CreatedOn) }
	if userObjectiveProgress.PlatformID != old.PlatformID { params[":a2"] = common.DynS(string(userObjectiveProgress.PlatformID)) }
	if userObjectiveProgress.UserID != old.UserID { params[":a3"] = common.DynS(userObjectiveProgress.UserID) }
	if userObjectiveProgress.PartnerID != old.PartnerID { params[":a4"] = common.DynS(userObjectiveProgress.PartnerID) }
	if userObjectiveProgress.Comments != old.Comments { params[":a5"] = common.DynS(userObjectiveProgress.Comments) }
	if userObjectiveProgress.Closeout != old.Closeout { params[":a6"] = common.DynN(userObjectiveProgress.Closeout) }
	if userObjectiveProgress.PercentTimeLapsed != old.PercentTimeLapsed { params[":a7"] = common.DynS(userObjectiveProgress.PercentTimeLapsed) }
	if userObjectiveProgress.StatusColor != old.StatusColor { params[":a8"] = common.DynS(string(userObjectiveProgress.StatusColor)) }
	if userObjectiveProgress.ReviewedByPartner != old.ReviewedByPartner { params[":a9"] = common.DynBOOL(userObjectiveProgress.ReviewedByPartner) }
	if userObjectiveProgress.PartnerComments != old.PartnerComments { params[":a10"] = common.DynS(userObjectiveProgress.PartnerComments) }
	if userObjectiveProgress.PartnerReportedProgress != old.PartnerReportedProgress { params[":a11"] = common.DynS(userObjectiveProgress.PartnerReportedProgress) }
	return
}
func updateExpression(userObjectiveProgress UserObjectiveProgress, old UserObjectiveProgress) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if userObjectiveProgress.ID != old.ID { updateParts = append(updateParts, "id = :a0"); params[":a0"] = common.DynS(userObjectiveProgress.ID);  }
	if userObjectiveProgress.CreatedOn != old.CreatedOn { updateParts = append(updateParts, "created_on = :a1"); params[":a1"] = common.DynS(userObjectiveProgress.CreatedOn);  }
	if userObjectiveProgress.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a2"); params[":a2"] = common.DynS(string(userObjectiveProgress.PlatformID));  }
	if userObjectiveProgress.UserID != old.UserID { updateParts = append(updateParts, "user_id = :a3"); params[":a3"] = common.DynS(userObjectiveProgress.UserID);  }
	if userObjectiveProgress.PartnerID != old.PartnerID { updateParts = append(updateParts, "partner_id = :a4"); params[":a4"] = common.DynS(userObjectiveProgress.PartnerID);  }
	if userObjectiveProgress.Comments != old.Comments { updateParts = append(updateParts, "comments = :a5"); params[":a5"] = common.DynS(userObjectiveProgress.Comments);  }
	if userObjectiveProgress.Closeout != old.Closeout { updateParts = append(updateParts, "closeout = :a6"); params[":a6"] = common.DynN(userObjectiveProgress.Closeout);  }
	if userObjectiveProgress.PercentTimeLapsed != old.PercentTimeLapsed { updateParts = append(updateParts, "percent_time_lapsed = :a7"); params[":a7"] = common.DynS(userObjectiveProgress.PercentTimeLapsed);  }
	if userObjectiveProgress.StatusColor != old.StatusColor { updateParts = append(updateParts, "status_color = :a8"); params[":a8"] = common.DynS(string(userObjectiveProgress.StatusColor));  }
	if userObjectiveProgress.ReviewedByPartner != old.ReviewedByPartner { updateParts = append(updateParts, "reviewed_by_partner = :a9"); params[":a9"] = common.DynBOOL(userObjectiveProgress.ReviewedByPartner);  }
	if userObjectiveProgress.PartnerComments != old.PartnerComments { updateParts = append(updateParts, "partner_comments = :a10"); params[":a10"] = common.DynS(userObjectiveProgress.PartnerComments);  }
	if userObjectiveProgress.PartnerReportedProgress != old.PartnerReportedProgress { updateParts = append(updateParts, "partner_reported_progress = :a11"); params[":a11"] = common.DynS(userObjectiveProgress.PartnerReportedProgress);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
