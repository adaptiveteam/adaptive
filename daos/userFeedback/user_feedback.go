package userFeedback
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

type UserFeedback struct  {
	ID string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	ValueID string `json:"value_id"`
	ConfidenceFactor string `json:"confidence_factor"`
	Feedback string `json:"feedback"`
	QuarterYear string `json:"quarter_year"`
	// ChannelID is a channel identifier. TODO: rename db field `channel` to `channel_id`
	// ChannelID, if any, to engage user in response to the feedback
	// This is useful to reply to an event with no knowledge of the previous context
	ChannelID string `json:"channel"`
	// A reference to the original timestamp that can be used to reply via threading
	MsgTimestamp string `json:"msg_timestamp"`
	PlatformID common.PlatformID `json:"platform_id"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (userFeedback UserFeedback)CollectEmptyFields() (emptyFields []string, ok bool) {
	if userFeedback.ID == "" { emptyFields = append(emptyFields, "ID")}
	if userFeedback.Source == "" { emptyFields = append(emptyFields, "Source")}
	if userFeedback.Target == "" { emptyFields = append(emptyFields, "Target")}
	if userFeedback.ValueID == "" { emptyFields = append(emptyFields, "ValueID")}
	if userFeedback.ConfidenceFactor == "" { emptyFields = append(emptyFields, "ConfidenceFactor")}
	if userFeedback.Feedback == "" { emptyFields = append(emptyFields, "Feedback")}
	if userFeedback.QuarterYear == "" { emptyFields = append(emptyFields, "QuarterYear")}
	if userFeedback.ChannelID == "" { emptyFields = append(emptyFields, "ChannelID")}
	if userFeedback.MsgTimestamp == "" { emptyFields = append(emptyFields, "MsgTimestamp")}
	if userFeedback.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (userFeedback UserFeedback) ToJSON() (string, error) {
	b, err := json.Marshal(userFeedback)
	return string(b), err
}

type DAO interface {
	Create(userFeedback UserFeedback) error
	CreateUnsafe(userFeedback UserFeedback)
	Read(id string) (userFeedback UserFeedback, err error)
	ReadUnsafe(id string) (userFeedback UserFeedback)
	ReadOrEmpty(id string) (userFeedback []UserFeedback, err error)
	ReadOrEmptyUnsafe(id string) (userFeedback []UserFeedback)
	CreateOrUpdate(userFeedback UserFeedback) error
	CreateOrUpdateUnsafe(userFeedback UserFeedback)
	Delete(id string) error
	DeleteUnsafe(id string)
	ReadByQuarterYearSource(quarterYear string, source string) (userFeedback []UserFeedback, err error)
	ReadByQuarterYearSourceUnsafe(quarterYear string, source string) (userFeedback []UserFeedback)
	ReadByQuarterYearTarget(quarterYear string, target string) (userFeedback []UserFeedback, err error)
	ReadByQuarterYearTargetUnsafe(quarterYear string, target string) (userFeedback []UserFeedback)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	Name      string                  `json:"name"`
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic(errors.New("Cannot create UserFeedback.DAO without clientID")) }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: TableName(clientID),
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic(errors.New("Cannot create UserFeedback.DAO without tableName")) }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}
func TableName(clientID string) string {
	return clientID + "_user_feedback"
}

// Create saves the UserFeedback.
func (d DAOImpl) Create(userFeedback UserFeedback) (err error) {
	emptyFields, ok := userFeedback.CollectEmptyFields()
	if ok {
		err = d.Dynamo.PutTableEntry(userFeedback, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the UserFeedback.
func (d DAOImpl) CreateUnsafe(userFeedback UserFeedback) {
	err2 := d.Create(userFeedback)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not create id==%s in %s\n", userFeedback.ID, d.Name))
}


// Read reads UserFeedback
func (d DAOImpl) Read(id string) (out UserFeedback, err error) {
	var outs []UserFeedback
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, d.Name)
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the UserFeedback. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) UserFeedback {
	out, err2 := d.Read(id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Error reading id==%s in %s\n", id, d.Name))
	return out
}


// ReadOrEmpty reads UserFeedback
func (d DAOImpl) ReadOrEmpty(id string) (out []UserFeedback, err error) {
	var outOrEmpty UserFeedback
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
	err = errors.Wrapf(err, "UserFeedback DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the UserFeedback. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []UserFeedback {
	out, err2 := d.ReadOrEmpty(id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Error while reading id==%s in %s\n", id, d.Name))
	return out
}


// CreateOrUpdate saves the UserFeedback regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(userFeedback UserFeedback) (err error) {
	
	var olds []UserFeedback
	olds, err = d.ReadOrEmpty(userFeedback.ID)
	err = errors.Wrapf(err, "UserFeedback DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", userFeedback.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(userFeedback)
			err = errors.Wrapf(err, "UserFeedback DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			emptyFields, ok := userFeedback.CollectEmptyFields()
			if ok {
				old := olds[0]
				
				key := idParams(old.ID)
				expr, exprAttributes, names := updateExpression(userFeedback, old)
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
				err = errors.Wrapf(err, "UserFeedback DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the UserFeedback regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(userFeedback UserFeedback) {
	err2 := d.CreateOrUpdate(userFeedback)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", userFeedback, d.Name))
}


// Delete removes UserFeedback from db
func (d DAOImpl)Delete(id string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(id))
}


// DeleteUnsafe deletes UserFeedback and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(id string) {
	err2 := d.Delete(id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not delete id==%s in %s\n", id, d.Name))
}


func (d DAOImpl)ReadByQuarterYearSource(quarterYear string, source string) (out []UserFeedback, err error) {
	var instances []UserFeedback
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "QuarterYearSourceIndex",
		Condition: "quarter_year = :a0 and #source = :a1",
		Attributes: map[string]interface{}{
			":a0": quarterYear,
			":a1": source,
		},
	}, map[string]string{"#source": "source"}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByQuarterYearSourceUnsafe(quarterYear string, source string) (out []UserFeedback) {
	out, err2 := d.ReadByQuarterYearSource(quarterYear, source)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not query QuarterYearSourceIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByQuarterYearTarget(quarterYear string, target string) (out []UserFeedback, err error) {
	var instances []UserFeedback
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "QuarterYearTargetIndex",
		Condition: "quarter_year = :a0 and target = :a1",
		Attributes: map[string]interface{}{
			":a0": quarterYear,
			":a1": target,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByQuarterYearTargetUnsafe(quarterYear string, target string) (out []UserFeedback) {
	out, err2 := d.ReadByQuarterYearTarget(quarterYear, target)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not query QuarterYearTargetIndex on %s table\n", d.Name))
	return
}

func idParams(id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"id": common.DynS(id),
	}
	return params
}
func allParams(userFeedback UserFeedback, old UserFeedback) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if userFeedback.ID != old.ID { params[":a0"] = common.DynS(userFeedback.ID) }
	if userFeedback.Source != old.Source { params[":a1"] = common.DynS(userFeedback.Source) }
	if userFeedback.Target != old.Target { params[":a2"] = common.DynS(userFeedback.Target) }
	if userFeedback.ValueID != old.ValueID { params[":a3"] = common.DynS(userFeedback.ValueID) }
	if userFeedback.ConfidenceFactor != old.ConfidenceFactor { params[":a4"] = common.DynS(userFeedback.ConfidenceFactor) }
	if userFeedback.Feedback != old.Feedback { params[":a5"] = common.DynS(userFeedback.Feedback) }
	if userFeedback.QuarterYear != old.QuarterYear { params[":a6"] = common.DynS(userFeedback.QuarterYear) }
	if userFeedback.ChannelID != old.ChannelID { params[":a7"] = common.DynS(userFeedback.ChannelID) }
	if userFeedback.MsgTimestamp != old.MsgTimestamp { params[":a8"] = common.DynS(userFeedback.MsgTimestamp) }
	if userFeedback.PlatformID != old.PlatformID { params[":a9"] = common.DynS(string(userFeedback.PlatformID)) }
	return
}
func updateExpression(userFeedback UserFeedback, old UserFeedback) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if userFeedback.ID != old.ID { updateParts = append(updateParts, "id = :a0"); params[":a0"] = common.DynS(userFeedback.ID);  }
	if userFeedback.Source != old.Source { updateParts = append(updateParts, "#source = :a1"); params[":a1"] = common.DynS(userFeedback.Source); fldName := "source"; names["#source"] = &fldName }
	if userFeedback.Target != old.Target { updateParts = append(updateParts, "target = :a2"); params[":a2"] = common.DynS(userFeedback.Target);  }
	if userFeedback.ValueID != old.ValueID { updateParts = append(updateParts, "value_id = :a3"); params[":a3"] = common.DynS(userFeedback.ValueID);  }
	if userFeedback.ConfidenceFactor != old.ConfidenceFactor { updateParts = append(updateParts, "confidence_factor = :a4"); params[":a4"] = common.DynS(userFeedback.ConfidenceFactor);  }
	if userFeedback.Feedback != old.Feedback { updateParts = append(updateParts, "feedback = :a5"); params[":a5"] = common.DynS(userFeedback.Feedback);  }
	if userFeedback.QuarterYear != old.QuarterYear { updateParts = append(updateParts, "quarter_year = :a6"); params[":a6"] = common.DynS(userFeedback.QuarterYear);  }
	if userFeedback.ChannelID != old.ChannelID { updateParts = append(updateParts, "channel = :a7"); params[":a7"] = common.DynS(userFeedback.ChannelID);  }
	if userFeedback.MsgTimestamp != old.MsgTimestamp { updateParts = append(updateParts, "msg_timestamp = :a8"); params[":a8"] = common.DynS(userFeedback.MsgTimestamp);  }
	if userFeedback.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a9"); params[":a9"] = common.DynS(string(userFeedback.PlatformID));  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
