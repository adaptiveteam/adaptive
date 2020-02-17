package userEngagement
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
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

// UserEngagement encapsulates an engagement we want to provide to a user
type UserEngagement struct  {
	// UserID is the ID of the user to send an engagement to
	// This usually corresponds to the platform user id
	UserID string `json:"user_id"`
	// A unique id to identify the engagement
	ID string `json:"id"`
	// PlatformID is the identifier of the platform.
	// It's used to get platform token required to send message to Slack/Teams.
	PlatformID common.PlatformID `json:"platform_id"`
	// TargetID is the ID of the user for whom this is related to
	TargetID string `json:"target_id"`
	// Namespace for the engagement
	Namespace string `json:"namespace"`
	// Check identifier for the engagement
	CheckIdentifier string `json:"check_identifier,omitempty"`
	CheckValue bool `json:"check_value,omitempty"`
	// Script that should be sent to a user to start engaging.
	// It's a serialized ebm.Message
	// deprecated. Use `Message` directly.
	Script string `json:"script"`
	// Message is the message we want to send to user
	Message ebm.Message `json:"message"`
	// Priority of the engagement
	// Urgent priority engagements are immediately sent to a user
	// Other priority engagements are queued up in the order of priority to be sent to user in next window
	Priority common.PriorityValue `json:"priority"`
	// A boolean flag indicating if it's optional
	Optional bool `json:"optional"`
	// Answered is a flag indicating that a user has responded to the engagement: 1 for answered, 0 for un-answered. 
	// This is required because, we need to keep the engagement even after a user has answered it. 
	// If the user wants to edit later, we will refer to the same engagement to post to user, like getting survey information
	// So, we need a way to differentiate between answered and unanswered engagements
	Answered int `json:"answered"`
	// Flag indicating if an engagement is ignored, 1 for yes, 0 for no
	Ignored int `json:"ignored"`
	EffectiveStartDate string `json:"effective_start_date,omitempty"`
	EffectiveEndDate string `json:"effective_end_date,omitempty"`
	PostedAt string `json:"posted_at,omitempty"`
	// Re-scheduled timestamp for the engagement, if any
	RescheduledFrom string `json:"rescheduled_from"`
	// Automatically maintained field
	CreatedAt string `json:"created_at"`
	// Automatically maintained field
	ModifiedAt string `json:"modified_at,omitempty"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (userEngagement UserEngagement)CollectEmptyFields() (emptyFields []string, ok bool) {
	if userEngagement.UserID == "" { emptyFields = append(emptyFields, "UserID")}
	if userEngagement.ID == "" { emptyFields = append(emptyFields, "ID")}
	if userEngagement.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if userEngagement.TargetID == "" { emptyFields = append(emptyFields, "TargetID")}
	if userEngagement.Namespace == "" { emptyFields = append(emptyFields, "Namespace")}
	if userEngagement.Script == "" { emptyFields = append(emptyFields, "Script")}
	if userEngagement.Priority == "" { emptyFields = append(emptyFields, "Priority")}
	if userEngagement.RescheduledFrom == "" { emptyFields = append(emptyFields, "RescheduledFrom")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (userEngagement UserEngagement) ToJSON() (string, error) {
	b, err := json.Marshal(userEngagement)
	return string(b), err
}

type DAO interface {
	Create(userEngagement UserEngagement) error
	CreateUnsafe(userEngagement UserEngagement)
	Read(userID string, id string) (userEngagement UserEngagement, err error)
	ReadUnsafe(userID string, id string) (userEngagement UserEngagement)
	ReadOrEmpty(userID string, id string) (userEngagement []UserEngagement, err error)
	ReadOrEmptyUnsafe(userID string, id string) (userEngagement []UserEngagement)
	CreateOrUpdate(userEngagement UserEngagement) error
	CreateOrUpdateUnsafe(userEngagement UserEngagement)
	Delete(userID string, id string) error
	DeleteUnsafe(userID string, id string)
	ReadByUserIDAnswered(userID string, answered int) (userEngagement []UserEngagement, err error)
	ReadByUserIDAnsweredUnsafe(userID string, answered int) (userEngagement []UserEngagement)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	Name      string                  `json:"name"`
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic(errors.New("Cannot create UserEngagement.DAO without clientID")) }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: TableName(clientID),
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic(errors.New("Cannot create UserEngagement.DAO without tableName")) }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}
func TableName(clientID string) string {
	return clientID + "_user_engagement"
}

// Create saves the UserEngagement.
func (d DAOImpl) Create(userEngagement UserEngagement) (err error) {
	emptyFields, ok := userEngagement.CollectEmptyFields()
	if ok {
		userEngagement.ModifiedAt = core.CurrentRFCTimestamp()
	userEngagement.CreatedAt = userEngagement.ModifiedAt
	err = d.Dynamo.PutTableEntry(userEngagement, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the UserEngagement.
func (d DAOImpl) CreateUnsafe(userEngagement UserEngagement) {
	err2 := d.Create(userEngagement)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not create userID==%s, id==%s in %s\n", userEngagement.UserID, userEngagement.ID, d.Name))
}


// Read reads UserEngagement
func (d DAOImpl) Read(userID string, id string) (out UserEngagement, err error) {
	var outs []UserEngagement
	outs, err = d.ReadOrEmpty(userID, id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found userID==%s, id==%s in %s\n", userID, id, d.Name)
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the UserEngagement. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(userID string, id string) UserEngagement {
	out, err2 := d.Read(userID, id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Error reading userID==%s, id==%s in %s\n", userID, id, d.Name))
	return out
}


// ReadOrEmpty reads UserEngagement
func (d DAOImpl) ReadOrEmpty(userID string, id string) (out []UserEngagement, err error) {
	var outOrEmpty UserEngagement
	ids := idParams(userID, id)
	var found bool
	found, err = d.Dynamo.GetItemOrEmptyFromTable(d.Name, ids, &outOrEmpty)
	if found {
		if outOrEmpty.UserID == userID && outOrEmpty.ID == id {
			out = append(out, outOrEmpty)
		} else {
			err = fmt.Errorf("Requested ids: userID==%s, id==%s are different from the found ones: userID==%s, id==%s", userID, id, outOrEmpty.UserID, outOrEmpty.ID) // unexpected error: found ids != ids
		}
	}
	err = errors.Wrapf(err, "UserEngagement DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the UserEngagement. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(userID string, id string) []UserEngagement {
	out, err2 := d.ReadOrEmpty(userID, id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Error while reading userID==%s, id==%s in %s\n", userID, id, d.Name))
	return out
}


// CreateOrUpdate saves the UserEngagement regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(userEngagement UserEngagement) (err error) {
	userEngagement.ModifiedAt = core.CurrentRFCTimestamp()
	if userEngagement.CreatedAt == "" { userEngagement.CreatedAt = userEngagement.ModifiedAt }
	
	var olds []UserEngagement
	olds, err = d.ReadOrEmpty(userEngagement.UserID, userEngagement.ID)
	err = errors.Wrapf(err, "UserEngagement DAO.CreateOrUpdate(id = userID==%s, id==%s) couldn't ReadOrEmpty", userEngagement.UserID, userEngagement.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(userEngagement)
			err = errors.Wrapf(err, "UserEngagement DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			emptyFields, ok := userEngagement.CollectEmptyFields()
			if ok {
				old := olds[0]
				userEngagement.ModifiedAt = core.CurrentRFCTimestamp()

				key := idParams(old.UserID, old.ID)
				expr, exprAttributes, names := updateExpression(userEngagement, old)
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
				err = errors.Wrapf(err, "UserEngagement DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the UserEngagement regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(userEngagement UserEngagement) {
	err2 := d.CreateOrUpdate(userEngagement)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", userEngagement, d.Name))
}


// Delete removes UserEngagement from db
func (d DAOImpl)Delete(userID string, id string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(userID, id))
}


// DeleteUnsafe deletes UserEngagement and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(userID string, id string) {
	err2 := d.Delete(userID, id)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not delete userID==%s, id==%s in %s\n", userID, id, d.Name))
}


func (d DAOImpl)ReadByUserIDAnswered(userID string, answered int) (out []UserEngagement, err error) {
	var instances []UserEngagement
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "UserIDAnsweredIndex",
		Condition: "user_id = :a0 and answered = :a1",
		Attributes: map[string]interface{}{
			":a0": userID,
			":a1": answered,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByUserIDAnsweredUnsafe(userID string, answered int) (out []UserEngagement) {
	out, err2 := d.ReadByUserIDAnswered(userID, answered)
	core.ErrorHandler(err2, d.Namespace, fmt.Sprintf("Could not query UserIDAnsweredIndex on %s table\n", d.Name))
	return
}

func idParams(userID string, id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"user_id": common.DynS(userID),
		"id": common.DynS(id),
	}
	return params
}
func allParams(userEngagement UserEngagement, old UserEngagement) (params map[string]*dynamodb.AttributeValue) {
	panic(errors.New("struct fields are not supported in UserEngagement.CreateOrUpdate/allParams"))
	return
}
func updateExpression(userEngagement UserEngagement, old UserEngagement) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if userEngagement.UserID != old.UserID { updateParts = append(updateParts, "user_id = :a0"); params[":a0"] = common.DynS(userEngagement.UserID);  }
	if userEngagement.ID != old.ID { updateParts = append(updateParts, "id = :a1"); params[":a1"] = common.DynS(userEngagement.ID);  }
	if userEngagement.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a2"); params[":a2"] = common.DynS(string(userEngagement.PlatformID));  }
	if userEngagement.TargetID != old.TargetID { updateParts = append(updateParts, "target_id = :a3"); params[":a3"] = common.DynS(userEngagement.TargetID);  }
	if userEngagement.Namespace != old.Namespace { updateParts = append(updateParts, "namespace = :a4"); params[":a4"] = common.DynS(userEngagement.Namespace);  }
	if userEngagement.CheckIdentifier != old.CheckIdentifier { updateParts = append(updateParts, "check_identifier = :a5"); params[":a5"] = common.DynS(userEngagement.CheckIdentifier);  }
	if userEngagement.CheckValue != old.CheckValue { updateParts = append(updateParts, "check_value = :a6"); params[":a6"] = common.DynBOOL(userEngagement.CheckValue);  }
	if userEngagement.Script != old.Script { updateParts = append(updateParts, "script = :a7"); params[":a7"] = common.DynS(userEngagement.Script);  }
	panic(errors.New("struct fields are not supported in UserEngagement.CreateOrUpdate/updateExpression Message"))
	if userEngagement.Priority != old.Priority { updateParts = append(updateParts, "priority = :a9"); params[":a9"] = common.DynS(string(userEngagement.Priority));  }
	if userEngagement.Optional != old.Optional { updateParts = append(updateParts, "optional = :a10"); params[":a10"] = common.DynBOOL(userEngagement.Optional);  }
	if userEngagement.Answered != old.Answered { updateParts = append(updateParts, "answered = :a11"); params[":a11"] = common.DynN(userEngagement.Answered);  }
	if userEngagement.Ignored != old.Ignored { updateParts = append(updateParts, "ignored = :a12"); params[":a12"] = common.DynN(userEngagement.Ignored);  }
	if userEngagement.EffectiveStartDate != old.EffectiveStartDate { updateParts = append(updateParts, "effective_start_date = :a13"); params[":a13"] = common.DynS(userEngagement.EffectiveStartDate);  }
	if userEngagement.EffectiveEndDate != old.EffectiveEndDate { updateParts = append(updateParts, "effective_end_date = :a14"); params[":a14"] = common.DynS(userEngagement.EffectiveEndDate);  }
	if userEngagement.PostedAt != old.PostedAt { updateParts = append(updateParts, "posted_at = :a15"); params[":a15"] = common.DynS(userEngagement.PostedAt);  }
	if userEngagement.RescheduledFrom != old.RescheduledFrom { updateParts = append(updateParts, "rescheduled_from = :a16"); params[":a16"] = common.DynS(userEngagement.RescheduledFrom);  }
	if userEngagement.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a17"); params[":a17"] = common.DynS(userEngagement.CreatedAt);  }
	if userEngagement.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a18"); params[":a18"] = common.DynS(userEngagement.ModifiedAt);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
