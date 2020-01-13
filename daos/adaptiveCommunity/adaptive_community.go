package adaptiveCommunity
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

type AdaptiveCommunity struct  {
	ID string `json:"id"`
	PlatformID common.PlatformID `json:"platform_id"`
	Channel string `json:"channel"`
	Active bool `json:"active"`
	RequestedBy string `json:"requested_by"`
	// Automatically maintained field
	CreatedAt string `json:"created_at"`
	// Automatically maintained field
	ModifiedAt string `json:"modified_at"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (adaptiveCommunity AdaptiveCommunity)CollectEmptyFields() (emptyFields []string, ok bool) {
	if adaptiveCommunity.ID == "" { emptyFields = append(emptyFields, "ID")}
	if adaptiveCommunity.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if adaptiveCommunity.Channel == "" { emptyFields = append(emptyFields, "Channel")}
	if adaptiveCommunity.RequestedBy == "" { emptyFields = append(emptyFields, "RequestedBy")}
	ok = len(emptyFields) == 0
	return
}

type DAO interface {
	Create(adaptiveCommunity AdaptiveCommunity) error
	CreateUnsafe(adaptiveCommunity AdaptiveCommunity)
	Read(id string) (adaptiveCommunity AdaptiveCommunity, err error)
	ReadUnsafe(id string) (adaptiveCommunity AdaptiveCommunity)
	ReadOrEmpty(id string) (adaptiveCommunity []AdaptiveCommunity, err error)
	ReadOrEmptyUnsafe(id string) (adaptiveCommunity []AdaptiveCommunity)
	CreateOrUpdate(adaptiveCommunity AdaptiveCommunity) error
	CreateOrUpdateUnsafe(adaptiveCommunity AdaptiveCommunity)
	Delete(id string) error
	DeleteUnsafe(id string)
	ReadByChannel(channel string) (adaptiveCommunity []AdaptiveCommunity, err error)
	ReadByChannelUnsafe(channel string) (adaptiveCommunity []AdaptiveCommunity)
	ReadByPlatformID(platformID common.PlatformID) (adaptiveCommunity []AdaptiveCommunity, err error)
	ReadByPlatformIDUnsafe(platformID common.PlatformID) (adaptiveCommunity []AdaptiveCommunity)
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
		Name: clientID + "_adaptive_community",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the AdaptiveCommunity.
func (d DAOImpl) Create(adaptiveCommunity AdaptiveCommunity) (err error) {
	emptyFields, ok := adaptiveCommunity.CollectEmptyFields()
	if ok {
		adaptiveCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())
	adaptiveCommunity.CreatedAt = adaptiveCommunity.ModifiedAt
	err = d.Dynamo.PutTableEntry(adaptiveCommunity, d.Name)
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the AdaptiveCommunity.
func (d DAOImpl) CreateUnsafe(adaptiveCommunity AdaptiveCommunity) {
	err := d.Create(adaptiveCommunity)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create id==%s in %s\n", adaptiveCommunity.ID, d.Name))
}


// Read reads AdaptiveCommunity
func (d DAOImpl) Read(id string) (out AdaptiveCommunity, err error) {
	var outs []AdaptiveCommunity
	outs, err = d.ReadOrEmpty(id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found id==%s in %s\n", id, d.Name)
	}
	return
}


// ReadUnsafe reads the AdaptiveCommunity. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(id string) AdaptiveCommunity {
	out, err := d.Read(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading id==%s in %s\n", id, d.Name))
	return out
}


// ReadOrEmpty reads AdaptiveCommunity
func (d DAOImpl) ReadOrEmpty(id string) (out []AdaptiveCommunity, err error) {
	var outOrEmpty AdaptiveCommunity
	ids := idParams(id)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.ID == id {
		out = append(out, outOrEmpty)
	} else if err != nil && strings.HasPrefix(err.Error(), "[NOT FOUND]") {
		err = nil // expected not-found error	
	}
	err = errors.Wrapf(err, "AdaptiveCommunity DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the AdaptiveCommunity. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(id string) []AdaptiveCommunity {
	out, err := d.ReadOrEmpty(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading id==%s in %s\n", id, d.Name))
	return out
}


// CreateOrUpdate saves the AdaptiveCommunity regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(adaptiveCommunity AdaptiveCommunity) (err error) {
	adaptiveCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())
	if adaptiveCommunity.CreatedAt == "" { adaptiveCommunity.CreatedAt = adaptiveCommunity.ModifiedAt }
	
	var olds []AdaptiveCommunity
	olds, err = d.ReadOrEmpty(adaptiveCommunity.ID)
	err = errors.Wrapf(err, "AdaptiveCommunity DAO.CreateOrUpdate(id = id==%s) couldn't ReadOrEmpty", adaptiveCommunity.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(adaptiveCommunity)
			err = errors.Wrapf(err, "AdaptiveCommunity DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			emptyFields, ok := adaptiveCommunity.CollectEmptyFields()
			if ok {
				old := olds[0]
				adaptiveCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())

				key := idParams(old.ID)
				expr, exprAttributes, names := updateExpression(adaptiveCommunity, old)
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
				err = errors.Wrapf(err, "AdaptiveCommunity DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, d.Name, expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the AdaptiveCommunity regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(adaptiveCommunity AdaptiveCommunity) {
	err := d.CreateOrUpdate(adaptiveCommunity)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", adaptiveCommunity, d.Name))
}


// Delete removes AdaptiveCommunity from db
func (d DAOImpl)Delete(id string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(id))
}


// DeleteUnsafe deletes AdaptiveCommunity and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(id string) {
	err := d.Delete(id)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete id==%s in %s\n", id, d.Name))
}


func (d DAOImpl)ReadByChannel(channel string) (out []AdaptiveCommunity, err error) {
	var instances []AdaptiveCommunity
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "ChannelIndex",
		Condition: "channel = :a0",
		Attributes: map[string]interface{}{
			":a0": channel,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByChannelUnsafe(channel string) (out []AdaptiveCommunity) {
	out, err := d.ReadByChannel(channel)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query ChannelIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByPlatformID(platformID common.PlatformID) (out []AdaptiveCommunity, err error) {
	var instances []AdaptiveCommunity
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


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID common.PlatformID) (out []AdaptiveCommunity) {
	out, err := d.ReadByPlatformID(platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query PlatformIDIndex on %s table\n", d.Name))
	return
}

func idParams(id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"id": common.DynS(id),
	}
	return params
}
func allParams(adaptiveCommunity AdaptiveCommunity, old AdaptiveCommunity) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if adaptiveCommunity.ID != old.ID { params[":a0"] = common.DynS(adaptiveCommunity.ID) }
	if adaptiveCommunity.PlatformID != old.PlatformID { params[":a1"] = common.DynS(string(adaptiveCommunity.PlatformID)) }
	if adaptiveCommunity.Channel != old.Channel { params[":a2"] = common.DynS(adaptiveCommunity.Channel) }
	if adaptiveCommunity.Active != old.Active { params[":a3"] = common.DynBOOL(adaptiveCommunity.Active) }
	if adaptiveCommunity.RequestedBy != old.RequestedBy { params[":a4"] = common.DynS(adaptiveCommunity.RequestedBy) }
	if adaptiveCommunity.CreatedAt != old.CreatedAt { params[":a5"] = common.DynS(adaptiveCommunity.CreatedAt) }
	if adaptiveCommunity.ModifiedAt != old.ModifiedAt { params[":a6"] = common.DynS(adaptiveCommunity.ModifiedAt) }
	return
}
func updateExpression(adaptiveCommunity AdaptiveCommunity, old AdaptiveCommunity) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if adaptiveCommunity.ID != old.ID { updateParts = append(updateParts, "id = :a0"); params[":a0"] = common.DynS(adaptiveCommunity.ID);  }
	if adaptiveCommunity.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a1"); params[":a1"] = common.DynS(string(adaptiveCommunity.PlatformID));  }
	if adaptiveCommunity.Channel != old.Channel { updateParts = append(updateParts, "channel = :a2"); params[":a2"] = common.DynS(adaptiveCommunity.Channel);  }
	if adaptiveCommunity.Active != old.Active { updateParts = append(updateParts, "active = :a3"); params[":a3"] = common.DynBOOL(adaptiveCommunity.Active);  }
	if adaptiveCommunity.RequestedBy != old.RequestedBy { updateParts = append(updateParts, "requested_by = :a4"); params[":a4"] = common.DynS(adaptiveCommunity.RequestedBy);  }
	if adaptiveCommunity.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a5"); params[":a5"] = common.DynS(adaptiveCommunity.CreatedAt);  }
	if adaptiveCommunity.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a6"); params[":a6"] = common.DynS(adaptiveCommunity.ModifiedAt);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
