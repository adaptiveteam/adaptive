package adaptiveCommunity
// This file has been automatically generated by `adaptive-platform/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/adaptiveteam/adaptive-utils-go/models"
	"time"
	awsutils "github.com/adaptiveteam/aws-utils-go"
	common "github.com/adaptiveteam/daos/common"
	core "github.com/adaptiveteam/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

type AdaptiveCommunity struct  {
	ID string `json:"id"`
	PlatformID models.PlatformID `json:"platform_id"`
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
	ReadByIDPlatformID(id string, platformID models.PlatformID) (adaptiveCommunity []AdaptiveCommunity, err error)
	ReadByIDPlatformIDUnsafe(id string, platformID models.PlatformID) (adaptiveCommunity []AdaptiveCommunity)
	ReadByChannel(channel string) (adaptiveCommunity []AdaptiveCommunity, err error)
	ReadByChannelUnsafe(channel string) (adaptiveCommunity []AdaptiveCommunity)
	ReadByPlatformID(platformID models.PlatformID) (adaptiveCommunity []AdaptiveCommunity, err error)
	ReadByPlatformIDUnsafe(platformID models.PlatformID) (adaptiveCommunity []AdaptiveCommunity)
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
func (d DAOImpl) Create(adaptiveCommunity AdaptiveCommunity) error {
	emptyFields, ok := adaptiveCommunity.CollectEmptyFields()
	if !ok {return fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)}
	adaptiveCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())
	adaptiveCommunity.CreatedAt = adaptiveCommunity.ModifiedAt
	return d.Dynamo.PutTableEntry(adaptiveCommunity, d.Name)
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
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(adaptiveCommunity)
			err = errors.Wrapf(err, "AdaptiveCommunity DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			old := olds[0]
			adaptiveCommunity.ModifiedAt = core.TimestampLayout.Format(time.Now())
ids := idParams(old.ID)
			err = d.Dynamo.UpdateTableEntry(
				allParams(adaptiveCommunity, old),
				ids,
				updateExpression(adaptiveCommunity, old),
				d.Name,
			)
			err = errors.Wrapf(err, "AdaptiveCommunity DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s", ids, d.Name)
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


func (d DAOImpl)ReadByIDPlatformID(id string, platformID models.PlatformID) (out []AdaptiveCommunity, err error) {
	var instances []AdaptiveCommunity
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "IDPlatformIDIndex",
		Condition: "id = :a0 and platform_id = :a1",
		Attributes: map[string]interface{}{
			":a0": id,
			":a1": platformID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByIDPlatformIDUnsafe(id string, platformID models.PlatformID) (out []AdaptiveCommunity) {
	out, err := d.ReadByIDPlatformID(id, platformID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query IDPlatformIDIndex on %s table\n", d.Name))
	return
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


func (d DAOImpl)ReadByPlatformID(platformID models.PlatformID) (out []AdaptiveCommunity, err error) {
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


func (d DAOImpl)ReadByPlatformIDUnsafe(platformID models.PlatformID) (out []AdaptiveCommunity) {
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
		if adaptiveCommunity.ID != old.ID { params["a0"] = common.DynS(adaptiveCommunity.ID) }
		if adaptiveCommunity.PlatformID != old.PlatformID { params["a1"] = common.DynS(string(adaptiveCommunity.PlatformID)) }
		if adaptiveCommunity.Channel != old.Channel { params["a2"] = common.DynS(adaptiveCommunity.Channel) }
		if adaptiveCommunity.Active != old.Active { params["a3"] = common.DynBOOL(adaptiveCommunity.Active) }
		if adaptiveCommunity.RequestedBy != old.RequestedBy { params["a4"] = common.DynS(adaptiveCommunity.RequestedBy) }
		if adaptiveCommunity.CreatedAt != old.CreatedAt { params["a5"] = common.DynS(adaptiveCommunity.CreatedAt) }
		if adaptiveCommunity.ModifiedAt != old.ModifiedAt { params["a6"] = common.DynS(adaptiveCommunity.ModifiedAt) }
	return
}
func updateExpression(adaptiveCommunity AdaptiveCommunity, old AdaptiveCommunity) string {
	var updateParts []string
	
		
			
		if adaptiveCommunity.ID != old.ID { updateParts = append(updateParts, "id = :a0") }
		if adaptiveCommunity.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a1") }
		if adaptiveCommunity.Channel != old.Channel { updateParts = append(updateParts, "channel = :a2") }
		if adaptiveCommunity.Active != old.Active { updateParts = append(updateParts, "active = :a3") }
		if adaptiveCommunity.RequestedBy != old.RequestedBy { updateParts = append(updateParts, "requested_by = :a4") }
		if adaptiveCommunity.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a5") }
		if adaptiveCommunity.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a6") }
	return strings.Join(updateParts, " and ")
}
