package userFeedback
// This file has been automatically generated by `adaptive-platform/scripts`
// The changes will be overridden by the next automatic generation.
import (
	awsutils "github.com/adaptiveteam/aws-utils-go"
	common "github.com/adaptiveteam/daos/common"
	core "github.com/adaptiveteam/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

type UserFeedback struct  {
	ID string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Dimension string `json:"dimension"`
	ConfidenceFactor string `json:"confidence_factor"`
	Feedback string `json:"feedback"`
	QuarterYear string `json:"quarter_year"`
	// Channel, if any, to engage user in response to the feedback
	// This is useful to reply to an event with no knowledge of the previous context
	Channel string `json:"channel"`
	// A reference to the original timestamp that can be used to reply via threading
	MsgTimestamp string `json:"msg_timestamp"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (userFeedback UserFeedback)CollectEmptyFields() (emptyFields []string, ok bool) {
	if userFeedback.ID == "" { emptyFields = append(emptyFields, "ID")}
	if userFeedback.Source == "" { emptyFields = append(emptyFields, "Source")}
	if userFeedback.Target == "" { emptyFields = append(emptyFields, "Target")}
	if userFeedback.Dimension == "" { emptyFields = append(emptyFields, "Dimension")}
	if userFeedback.ConfidenceFactor == "" { emptyFields = append(emptyFields, "ConfidenceFactor")}
	if userFeedback.Feedback == "" { emptyFields = append(emptyFields, "Feedback")}
	if userFeedback.QuarterYear == "" { emptyFields = append(emptyFields, "QuarterYear")}
	if userFeedback.Channel == "" { emptyFields = append(emptyFields, "Channel")}
	if userFeedback.MsgTimestamp == "" { emptyFields = append(emptyFields, "MsgTimestamp")}
	ok = len(emptyFields) == 0
	return
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
	ReadByID(id string) (userFeedback []UserFeedback, err error)
	ReadByIDUnsafe(id string) (userFeedback []UserFeedback)
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
	if clientID == "" { panic("Cannot create DAO without clientID") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: clientID + "_user_feedback",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the UserFeedback.
func (d DAOImpl) Create(userFeedback UserFeedback) error {
	emptyFields, ok := userFeedback.CollectEmptyFields()
	if !ok {return fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)}
	return d.Dynamo.PutTableEntry(userFeedback, d.Name)
}


// CreateUnsafe saves the UserFeedback.
func (d DAOImpl) CreateUnsafe(userFeedback UserFeedback) {
	err := d.Create(userFeedback)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create id==%s in %s\n", userFeedback.ID, d.Name))
}


// Read reads UserFeedback
func (d DAOImpl) Read(id string) (out UserFeedback, err error) {
	var outs []UserFeedback
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, d.Name)
	}
	return
}


// ReadUnsafe reads the UserFeedback. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) UserFeedback {
	out, err := d.Read(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading id==%s in %s\n", id, d.Name))
	return out
}


// ReadOrEmpty reads UserFeedback
func (d DAOImpl) ReadOrEmpty(id string) (out []UserFeedback, err error) {
	var outOrEmpty UserFeedback
	ids := idParams(id)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.ID == id {
		out = append(out, outOrEmpty)
	}
	err = errors.Wrapf(err, "UserFeedback DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the UserFeedback. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []UserFeedback {
	out, err := d.ReadOrEmpty(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading id==%s in %s\n", id, d.Name))
	return out
}


// CreateOrUpdate saves the UserFeedback regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(userFeedback UserFeedback) (err error) {
	
	var olds []UserFeedback
	olds, err = d.ReadOrEmpty(userFeedback.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(userFeedback)
			err = errors.Wrapf(err, "UserFeedback DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			old := olds[0]
			ids := idParams(old.ID)
			err = d.Dynamo.UpdateTableEntry(
				allParams(userFeedback, old),
				ids,
				updateExpression(userFeedback, old),
				d.Name,
			)
			err = errors.Wrapf(err, "UserFeedback DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s", ids, d.Name)
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the UserFeedback regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(userFeedback UserFeedback) {
	err := d.CreateOrUpdate(userFeedback)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", userFeedback, d.Name))
}


// Delete removes UserFeedback from db
func (d DAOImpl)Delete(id string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(id))
}


// DeleteUnsafe deletes UserFeedback and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(id string) {
	err := d.Delete(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete id==%s in %s\n", id, d.Name))
}


func (d DAOImpl)ReadByID(id string) (out []UserFeedback, err error) {
	var instances []UserFeedback
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


func (d DAOImpl)ReadByIDUnsafe(id string) (out []UserFeedback) {
	out, err := d.ReadByID(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query IDIndex on %s table\n", d.Name))
	return
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
	out, err := d.ReadByQuarterYearSource(quarterYear, source)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query QuarterYearSourceIndex on %s table\n", d.Name))
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
	out, err := d.ReadByQuarterYearTarget(quarterYear, target)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query QuarterYearTargetIndex on %s table\n", d.Name))
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
		if userFeedback.ID != old.ID { params["a0"] = common.DynS(userFeedback.ID) }
		if userFeedback.Source != old.Source { params["a1"] = common.DynS(userFeedback.Source) }
		if userFeedback.Target != old.Target { params["a2"] = common.DynS(userFeedback.Target) }
		if userFeedback.Dimension != old.Dimension { params["a3"] = common.DynS(userFeedback.Dimension) }
		if userFeedback.ConfidenceFactor != old.ConfidenceFactor { params["a4"] = common.DynS(userFeedback.ConfidenceFactor) }
		if userFeedback.Feedback != old.Feedback { params["a5"] = common.DynS(userFeedback.Feedback) }
		if userFeedback.QuarterYear != old.QuarterYear { params["a6"] = common.DynS(userFeedback.QuarterYear) }
		if userFeedback.Channel != old.Channel { params["a7"] = common.DynS(userFeedback.Channel) }
		if userFeedback.MsgTimestamp != old.MsgTimestamp { params["a8"] = common.DynS(userFeedback.MsgTimestamp) }
	return
}
func updateExpression(userFeedback UserFeedback, old UserFeedback) string {
	var updateParts []string
	
		
			
		if userFeedback.ID != old.ID { updateParts = append(updateParts, "id = :a0") }
		if userFeedback.Source != old.Source { updateParts = append(updateParts, "#source = :a1") }
		if userFeedback.Target != old.Target { updateParts = append(updateParts, "target = :a2") }
		if userFeedback.Dimension != old.Dimension { updateParts = append(updateParts, "dimension = :a3") }
		if userFeedback.ConfidenceFactor != old.ConfidenceFactor { updateParts = append(updateParts, "confidence_factor = :a4") }
		if userFeedback.Feedback != old.Feedback { updateParts = append(updateParts, "feedback = :a5") }
		if userFeedback.QuarterYear != old.QuarterYear { updateParts = append(updateParts, "quarter_year = :a6") }
		if userFeedback.Channel != old.Channel { updateParts = append(updateParts, "channel = :a7") }
		if userFeedback.MsgTimestamp != old.MsgTimestamp { updateParts = append(updateParts, "msg_timestamp = :a8") }
	return strings.Join(updateParts, " and ")
}
