package userObjective
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"time"
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

type DevelopmentObjectiveType string
const (
	IndividualDevelopmentObjective DevelopmentObjectiveType = "individual"
	StrategyDevelopmentObjective DevelopmentObjectiveType = "strategy"
)

type AlignedStrategyType string
const (
	ObjectiveStrategyObjectiveAlignment AlignedStrategyType = "strategy_objective"
	ObjectiveStrategyInitiativeAlignment AlignedStrategyType = "strategy_initiative"
	ObjectiveCompetencyAlignment AlignedStrategyType = "competency"
	ObjectiveNoStrategyAlignment AlignedStrategyType = "none"
)

type UserObjective struct  {
	ID string `json:"id"`
	PlatformID common.PlatformID `json:"platform_id"`
	// UserId is the Id of the user to send an engagement to
	// This usually corresponds to the platform user id
	UserID string `json:"user_id"`
	Name string `json:"name"`
	Description string `json:"description"`
	AccountabilityPartner string `json:"accountability_partner"`
	// 1 for true, 0 for false
	Accepted int `json:"accepted"`
	ObjectiveType DevelopmentObjectiveType `json:"type"`
	StrategyAlignmentEntityID string `json:"strategy_alignment_entity_id,omitempty"`
	StrategyAlignmentEntityType AlignedStrategyType `json:"strategy_alignment_entity_type"`
	Quarter int `json:"quarter"`
	Year int `json:"year"`
	// Deprecated, use CreatedAt automated field
	CreatedDate string `json:"created_date"`
	ExpectedEndDate string `json:"expected_end_date"`
	// 1 for true, 0 for false
	Completed int `json:"completed"`
	PartnerVerifiedCompletion bool `json:"partner_verified_completion"`
	CompletedDate string `json:"completed_date,omitempty"`
	PartnerVerifiedCompletionDate string `json:"partner_verified_completion_date,omitempty"`
	Comments string `json:"comments,omitempty"`
	// 1 for true, 0 for false
	Cancelled int `json:"cancelled"`
	// Automatically maintained field
	CreatedAt string `json:"created_at"`
	// Automatically maintained field
	ModifiedAt string `json:"modified_at"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (userObjective UserObjective)CollectEmptyFields() (emptyFields []string, ok bool) {
	if userObjective.ID == "" { emptyFields = append(emptyFields, "ID")}
	if userObjective.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if userObjective.UserID == "" { emptyFields = append(emptyFields, "UserID")}
	if userObjective.Name == "" { emptyFields = append(emptyFields, "Name")}
	if userObjective.Description == "" { emptyFields = append(emptyFields, "Description")}
	if userObjective.AccountabilityPartner == "" { emptyFields = append(emptyFields, "AccountabilityPartner")}
	if userObjective.CreatedDate == "" { emptyFields = append(emptyFields, "CreatedDate")}
	if userObjective.ExpectedEndDate == "" { emptyFields = append(emptyFields, "ExpectedEndDate")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (userObjective UserObjective) ToJSON() (string, error) {
	b, err := json.Marshal(userObjective)
	return string(b), err
}

type DAO interface {
	Create(userObjective UserObjective) error
	CreateUnsafe(userObjective UserObjective)
	Read(id string) (userObjective UserObjective, err error)
	ReadUnsafe(id string) (userObjective UserObjective)
	ReadOrEmpty(id string) (userObjective []UserObjective, err error)
	ReadOrEmptyUnsafe(id string) (userObjective []UserObjective)
	CreateOrUpdate(userObjective UserObjective) error
	CreateOrUpdateUnsafe(userObjective UserObjective)
	Delete(id string) error
	DeleteUnsafe(id string)
	ReadByUserIDCompleted(userID string, completed int) (userObjective []UserObjective, err error)
	ReadByUserIDCompletedUnsafe(userID string, completed int) (userObjective []UserObjective)
	ReadByAccepted(accepted int) (userObjective []UserObjective, err error)
	ReadByAcceptedUnsafe(accepted int) (userObjective []UserObjective)
	ReadByAccountabilityPartner(accountabilityPartner string) (userObjective []UserObjective, err error)
	ReadByAccountabilityPartnerUnsafe(accountabilityPartner string) (userObjective []UserObjective)
	ReadByUserIDType(userID string, objectiveType DevelopmentObjectiveType) (userObjective []UserObjective, err error)
	ReadByUserIDTypeUnsafe(userID string, objectiveType DevelopmentObjectiveType) (userObjective []UserObjective)
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
		Name: clientID + "_user_objective",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the UserObjective.
func (d DAOImpl) Create(userObjective UserObjective) (err error) {
	emptyFields, ok := userObjective.CollectEmptyFields()
	if ok {
		userObjective.ModifiedAt = core.TimestampLayout.Format(time.Now())
	userObjective.CreatedAt = userObjective.ModifiedAt
	err = d.Dynamo.PutTableEntry(userObjective, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the UserObjective.
func (d DAOImpl) CreateUnsafe(userObjective UserObjective) {
	err := d.Create(userObjective)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create id==%s in %s\n", userObjective.ID, d.Name))
}


// Read reads UserObjective
func (d DAOImpl) Read(id string) (out UserObjective, err error) {
	var outs []UserObjective
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, d.Name)
	}
	return
}


// ReadUnsafe reads the UserObjective. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) UserObjective {
	out, err := d.Read(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading id==%s in %s\n", id, d.Name))
	return out
}


// ReadOrEmpty reads UserObjective
func (d DAOImpl) ReadOrEmpty(id string) (out []UserObjective, err error) {
	var outOrEmpty UserObjective
	ids := idParams(id)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.ID == id {
		out = append(out, outOrEmpty)
	} else if err != nil && strings.HasPrefix(err.Error(), "[NOT FOUND]") {
		err = nil // expected not-found error	
	}
	err = errors.Wrapf(err, "UserObjective DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the UserObjective. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []UserObjective {
	out, err := d.ReadOrEmpty(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading id==%s in %s\n", id, d.Name))
	return out
}


// CreateOrUpdate saves the UserObjective regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(userObjective UserObjective) (err error) {
	userObjective.ModifiedAt = core.TimestampLayout.Format(time.Now())
	if userObjective.CreatedAt == "" { userObjective.CreatedAt = userObjective.ModifiedAt }
	
	var olds []UserObjective
	olds, err = d.ReadOrEmpty(userObjective.ID)
	err = errors.Wrapf(err, "UserObjective DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", userObjective.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(userObjective)
			err = errors.Wrapf(err, "UserObjective DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			emptyFields, ok := userObjective.CollectEmptyFields()
			if ok {
				old := olds[0]
				userObjective.ModifiedAt = core.TimestampLayout.Format(time.Now())

				key := idParams(old.ID)
				expr, exprAttributes, names := updateExpression(userObjective, old)
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
				err = errors.Wrapf(err, "UserObjective DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the UserObjective regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(userObjective UserObjective) {
	err := d.CreateOrUpdate(userObjective)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", userObjective, d.Name))
}


// Delete removes UserObjective from db
func (d DAOImpl)Delete(id string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(id))
}


// DeleteUnsafe deletes UserObjective and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(id string) {
	err := d.Delete(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete id==%s in %s\n", id, d.Name))
}


func (d DAOImpl)ReadByUserIDCompleted(userID string, completed int) (out []UserObjective, err error) {
	var instances []UserObjective
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "UserIDCompletedIndex",
		Condition: "user_id = :a0 and completed = :a1",
		Attributes: map[string]interface{}{
			":a0": userID,
			":a1": completed,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByUserIDCompletedUnsafe(userID string, completed int) (out []UserObjective) {
	out, err := d.ReadByUserIDCompleted(userID, completed)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query UserIDCompletedIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByAccepted(accepted int) (out []UserObjective, err error) {
	var instances []UserObjective
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "AcceptedIndex",
		Condition: "accepted = :a0",
		Attributes: map[string]interface{}{
			":a0": accepted,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByAcceptedUnsafe(accepted int) (out []UserObjective) {
	out, err := d.ReadByAccepted(accepted)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query AcceptedIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByAccountabilityPartner(accountabilityPartner string) (out []UserObjective, err error) {
	var instances []UserObjective
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "AccountabilityPartnerIndex",
		Condition: "accountability_partner = :a0",
		Attributes: map[string]interface{}{
			":a0": accountabilityPartner,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByAccountabilityPartnerUnsafe(accountabilityPartner string) (out []UserObjective) {
	out, err := d.ReadByAccountabilityPartner(accountabilityPartner)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query AccountabilityPartnerIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByUserIDType(userID string, objectiveType DevelopmentObjectiveType) (out []UserObjective, err error) {
	var instances []UserObjective
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "UserIDTypeIndex",
		Condition: "user_id = :a0 and #type = :a1",
		Attributes: map[string]interface{}{
			":a0": userID,
			":a1": objectiveType,
		},
	}, map[string]string{"#type": "type"}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByUserIDTypeUnsafe(userID string, objectiveType DevelopmentObjectiveType) (out []UserObjective) {
	out, err := d.ReadByUserIDType(userID, objectiveType)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query UserIDTypeIndex on %s table\n", d.Name))
	return
}

func idParams(id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"id": common.DynS(id),
	}
	return params
}
func allParams(userObjective UserObjective, old UserObjective) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if userObjective.ID != old.ID { params[":a0"] = common.DynS(userObjective.ID) }
	if userObjective.PlatformID != old.PlatformID { params[":a1"] = common.DynS(string(userObjective.PlatformID)) }
	if userObjective.UserID != old.UserID { params[":a2"] = common.DynS(userObjective.UserID) }
	if userObjective.Name != old.Name { params[":a3"] = common.DynS(userObjective.Name) }
	if userObjective.Description != old.Description { params[":a4"] = common.DynS(userObjective.Description) }
	if userObjective.AccountabilityPartner != old.AccountabilityPartner { params[":a5"] = common.DynS(userObjective.AccountabilityPartner) }
	if userObjective.Accepted != old.Accepted { params[":a6"] = common.DynN(userObjective.Accepted) }
	if userObjective.ObjectiveType != old.ObjectiveType { params[":a7"] = common.DynS(string(userObjective.ObjectiveType)) }
	if userObjective.StrategyAlignmentEntityID != old.StrategyAlignmentEntityID { params[":a8"] = common.DynS(userObjective.StrategyAlignmentEntityID) }
	if userObjective.StrategyAlignmentEntityType != old.StrategyAlignmentEntityType { params[":a9"] = common.DynS(string(userObjective.StrategyAlignmentEntityType)) }
	if userObjective.Quarter != old.Quarter { params[":a10"] = common.DynN(userObjective.Quarter) }
	if userObjective.Year != old.Year { params[":a11"] = common.DynN(userObjective.Year) }
	if userObjective.CreatedDate != old.CreatedDate { params[":a12"] = common.DynS(userObjective.CreatedDate) }
	if userObjective.ExpectedEndDate != old.ExpectedEndDate { params[":a13"] = common.DynS(userObjective.ExpectedEndDate) }
	if userObjective.Completed != old.Completed { params[":a14"] = common.DynN(userObjective.Completed) }
	if userObjective.PartnerVerifiedCompletion != old.PartnerVerifiedCompletion { params[":a15"] = common.DynBOOL(userObjective.PartnerVerifiedCompletion) }
	if userObjective.CompletedDate != old.CompletedDate { params[":a16"] = common.DynS(userObjective.CompletedDate) }
	if userObjective.PartnerVerifiedCompletionDate != old.PartnerVerifiedCompletionDate { params[":a17"] = common.DynS(userObjective.PartnerVerifiedCompletionDate) }
	if userObjective.Comments != old.Comments { params[":a18"] = common.DynS(userObjective.Comments) }
	if userObjective.Cancelled != old.Cancelled { params[":a19"] = common.DynN(userObjective.Cancelled) }
	if userObjective.CreatedAt != old.CreatedAt { params[":a20"] = common.DynS(userObjective.CreatedAt) }
	if userObjective.ModifiedAt != old.ModifiedAt { params[":a21"] = common.DynS(userObjective.ModifiedAt) }
	return
}
func updateExpression(userObjective UserObjective, old UserObjective) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if userObjective.ID != old.ID { updateParts = append(updateParts, "id = :a0"); params[":a0"] = common.DynS(userObjective.ID);  }
	if userObjective.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a1"); params[":a1"] = common.DynS(string(userObjective.PlatformID));  }
	if userObjective.UserID != old.UserID { updateParts = append(updateParts, "user_id = :a2"); params[":a2"] = common.DynS(userObjective.UserID);  }
	if userObjective.Name != old.Name { updateParts = append(updateParts, "#name = :a3"); params[":a3"] = common.DynS(userObjective.Name); fldName := "name"; names["#name"] = &fldName }
	if userObjective.Description != old.Description { updateParts = append(updateParts, "description = :a4"); params[":a4"] = common.DynS(userObjective.Description);  }
	if userObjective.AccountabilityPartner != old.AccountabilityPartner { updateParts = append(updateParts, "accountability_partner = :a5"); params[":a5"] = common.DynS(userObjective.AccountabilityPartner);  }
	if userObjective.Accepted != old.Accepted { updateParts = append(updateParts, "accepted = :a6"); params[":a6"] = common.DynN(userObjective.Accepted);  }
	if userObjective.ObjectiveType != old.ObjectiveType { updateParts = append(updateParts, "#type = :a7"); params[":a7"] = common.DynS(string(userObjective.ObjectiveType)); fldName := "type"; names["#type"] = &fldName }
	if userObjective.StrategyAlignmentEntityID != old.StrategyAlignmentEntityID { updateParts = append(updateParts, "strategy_alignment_entity_id = :a8"); params[":a8"] = common.DynS(userObjective.StrategyAlignmentEntityID);  }
	if userObjective.StrategyAlignmentEntityType != old.StrategyAlignmentEntityType { updateParts = append(updateParts, "strategy_alignment_entity_type = :a9"); params[":a9"] = common.DynS(string(userObjective.StrategyAlignmentEntityType));  }
	if userObjective.Quarter != old.Quarter { updateParts = append(updateParts, "quarter = :a10"); params[":a10"] = common.DynN(userObjective.Quarter);  }
	if userObjective.Year != old.Year { updateParts = append(updateParts, "#year = :a11"); params[":a11"] = common.DynN(userObjective.Year); fldName := "year"; names["#year"] = &fldName }
	if userObjective.CreatedDate != old.CreatedDate { updateParts = append(updateParts, "created_date = :a12"); params[":a12"] = common.DynS(userObjective.CreatedDate);  }
	if userObjective.ExpectedEndDate != old.ExpectedEndDate { updateParts = append(updateParts, "expected_end_date = :a13"); params[":a13"] = common.DynS(userObjective.ExpectedEndDate);  }
	if userObjective.Completed != old.Completed { updateParts = append(updateParts, "completed = :a14"); params[":a14"] = common.DynN(userObjective.Completed);  }
	if userObjective.PartnerVerifiedCompletion != old.PartnerVerifiedCompletion { updateParts = append(updateParts, "partner_verified_completion = :a15"); params[":a15"] = common.DynBOOL(userObjective.PartnerVerifiedCompletion);  }
	if userObjective.CompletedDate != old.CompletedDate { updateParts = append(updateParts, "completed_date = :a16"); params[":a16"] = common.DynS(userObjective.CompletedDate);  }
	if userObjective.PartnerVerifiedCompletionDate != old.PartnerVerifiedCompletionDate { updateParts = append(updateParts, "partner_verified_completion_date = :a17"); params[":a17"] = common.DynS(userObjective.PartnerVerifiedCompletionDate);  }
	if userObjective.Comments != old.Comments { updateParts = append(updateParts, "comments = :a18"); params[":a18"] = common.DynS(userObjective.Comments);  }
	if userObjective.Cancelled != old.Cancelled { updateParts = append(updateParts, "cancelled = :a19"); params[":a19"] = common.DynN(userObjective.Cancelled);  }
	if userObjective.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a20"); params[":a20"] = common.DynS(userObjective.CreatedAt);  }
	if userObjective.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a21"); params[":a21"] = common.DynS(userObjective.ModifiedAt);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
