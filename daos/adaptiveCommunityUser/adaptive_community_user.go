package adaptiveCommunityUser
// This file has been automatically generated by `adaptive-platform/scripts`
// The changes will be overridden by the next automatic generation.
import (
	"github.com/adaptiveteam/adaptive-utils-go/models"
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/aws-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

type AdaptiveCommunityUser struct  {
	ChannelID string `json:"channel_id"`
	UserID string `json:"user_id"`
	PlatformID models.PlatformID `json:"platform_id"`
	CommunityID string `json:"community_id"`
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (adaptiveCommunityUser AdaptiveCommunityUser)CollectEmptyFields() (emptyFields []string, ok bool) {
	if adaptiveCommunityUser.ChannelID == "" { emptyFields = append(emptyFields, "ChannelID")}
	if adaptiveCommunityUser.UserID == "" { emptyFields = append(emptyFields, "UserID")}
	if adaptiveCommunityUser.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if adaptiveCommunityUser.CommunityID == "" { emptyFields = append(emptyFields, "CommunityID")}
	ok = len(emptyFields) == 0
	return
}

type DAO interface {
	Create(adaptiveCommunityUser AdaptiveCommunityUser) error
	CreateUnsafe(adaptiveCommunityUser AdaptiveCommunityUser)
	Read(channelID string, userID string) (adaptiveCommunityUser AdaptiveCommunityUser, err error)
	ReadUnsafe(channelID string, userID string) (adaptiveCommunityUser AdaptiveCommunityUser)
	ReadOrEmpty(channelID string, userID string) (adaptiveCommunityUser []AdaptiveCommunityUser, err error)
	ReadOrEmptyUnsafe(channelID string, userID string) (adaptiveCommunityUser []AdaptiveCommunityUser)
	CreateOrUpdate(adaptiveCommunityUser AdaptiveCommunityUser) error
	CreateOrUpdateUnsafe(adaptiveCommunityUser AdaptiveCommunityUser)
	Delete(channelID string, userID string) error
	DeleteUnsafe(channelID string, userID string)
	ReadByChannelID(channelID string) (adaptiveCommunityUser []AdaptiveCommunityUser, err error)
	ReadByChannelIDUnsafe(channelID string) (adaptiveCommunityUser []AdaptiveCommunityUser)
	ReadByUserIDCommunityID(userID string, communityID string) (adaptiveCommunityUser []AdaptiveCommunityUser, err error)
	ReadByUserIDCommunityIDUnsafe(userID string, communityID string) (adaptiveCommunityUser []AdaptiveCommunityUser)
	ReadByUserID(userID string) (adaptiveCommunityUser []AdaptiveCommunityUser, err error)
	ReadByUserIDUnsafe(userID string) (adaptiveCommunityUser []AdaptiveCommunityUser)
	ReadByPlatformIDCommunityID(platformID models.PlatformID, communityID string) (adaptiveCommunityUser []AdaptiveCommunityUser, err error)
	ReadByPlatformIDCommunityIDUnsafe(platformID models.PlatformID, communityID string) (adaptiveCommunityUser []AdaptiveCommunityUser)
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
		Name: clientID + "_adaptive_community_user",
	}
}

// NewDAOByTableName creates an instance of DAO that will provide access to the table
func NewDAOByTableName(dynamo *awsutils.DynamoRequest, namespace, tableName string) DAO {
	if tableName == "" { panic("Cannot create DAO without tableName") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		Name: tableName,
	}
}

// Create saves the AdaptiveCommunityUser.
func (d DAOImpl) Create(adaptiveCommunityUser AdaptiveCommunityUser) error {
	emptyFields, ok := adaptiveCommunityUser.CollectEmptyFields()
	if !ok {return fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)}
	return d.Dynamo.PutTableEntry(adaptiveCommunityUser, d.Name)
}


// CreateUnsafe saves the AdaptiveCommunityUser.
func (d DAOImpl) CreateUnsafe(adaptiveCommunityUser AdaptiveCommunityUser) {
	err := d.Create(adaptiveCommunityUser)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create channelID==%s, userID==%s in %s\n", adaptiveCommunityUser.ChannelID, adaptiveCommunityUser.UserID, d.Name))
}


// Read reads AdaptiveCommunityUser
func (d DAOImpl) Read(channelID string, userID string) (out AdaptiveCommunityUser, err error) {
	var outs []AdaptiveCommunityUser
	outs, err = d.ReadOrEmpty(channelID, userID)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found channelID==%s, userID==%s in %s\n", channelID, userID, d.Name)
	}
	return
}


// ReadUnsafe reads the AdaptiveCommunityUser. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(channelID string, userID string) AdaptiveCommunityUser {
	out, err := d.Read(channelID, userID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error reading channelID==%s, userID==%s in %s\n", channelID, userID, d.Name))
	return out
}


// ReadOrEmpty reads AdaptiveCommunityUser
func (d DAOImpl) ReadOrEmpty(channelID string, userID string) (out []AdaptiveCommunityUser, err error) {
	var outOrEmpty AdaptiveCommunityUser
	ids := idParams(channelID, userID)
	err = d.Dynamo.QueryTable(d.Name, ids, &outOrEmpty)
	if outOrEmpty.ChannelID == channelID && outOrEmpty.UserID == userID {
		out = append(out, outOrEmpty)
	} else if err != nil && strings.HasPrefix(err.Error(), "In table ") {
		err = nil // expected not-found error	
	}
	err = errors.Wrapf(err, "AdaptiveCommunityUser DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, d.Name)
	return
}


// ReadOrEmptyUnsafe reads the AdaptiveCommunityUser. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(channelID string, userID string) []AdaptiveCommunityUser {
	out, err := d.ReadOrEmpty(channelID, userID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Error while reading channelID==%s, userID==%s in %s\n", channelID, userID, d.Name))
	return out
}


// CreateOrUpdate saves the AdaptiveCommunityUser regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(adaptiveCommunityUser AdaptiveCommunityUser) (err error) {
	
	var olds []AdaptiveCommunityUser
	olds, err = d.ReadOrEmpty(adaptiveCommunityUser.ChannelID, adaptiveCommunityUser.UserID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(adaptiveCommunityUser)
			err = errors.Wrapf(err, "AdaptiveCommunityUser DAO.CreateOrUpdate couldn't Create in table %s", d.Name)
		} else {
			old := olds[0]
			
			key := idParams(old.ChannelID, old.UserID)
			expr, exprAttributes, names := updateExpression(adaptiveCommunityUser, old)
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
			err = errors.Wrapf(err, "AdaptiveCommunityUser DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s", key, d.Name)
			return
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the AdaptiveCommunityUser regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(adaptiveCommunityUser AdaptiveCommunityUser) {
	err := d.CreateOrUpdate(adaptiveCommunityUser)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("could not create or update %v in %s\n", adaptiveCommunityUser, d.Name))
}


// Delete removes AdaptiveCommunityUser from db
func (d DAOImpl)Delete(channelID string, userID string) error {
	return d.Dynamo.DeleteEntry(d.Name, idParams(channelID, userID))
}


// DeleteUnsafe deletes AdaptiveCommunityUser and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(channelID string, userID string) {
	err := d.Delete(channelID, userID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not delete channelID==%s, userID==%s in %s\n", channelID, userID, d.Name))
}


func (d DAOImpl)ReadByChannelID(channelID string) (out []AdaptiveCommunityUser, err error) {
	var instances []AdaptiveCommunityUser
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "ChannelIDIndex",
		Condition: "channel_id = :a0",
		Attributes: map[string]interface{}{
			":a0": channelID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByChannelIDUnsafe(channelID string) (out []AdaptiveCommunityUser) {
	out, err := d.ReadByChannelID(channelID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query ChannelIDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByUserIDCommunityID(userID string, communityID string) (out []AdaptiveCommunityUser, err error) {
	var instances []AdaptiveCommunityUser
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "UserIDCommunityIDIndex",
		Condition: "user_id = :a0 and community_id = :a1",
		Attributes: map[string]interface{}{
			":a0": userID,
			":a1": communityID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByUserIDCommunityIDUnsafe(userID string, communityID string) (out []AdaptiveCommunityUser) {
	out, err := d.ReadByUserIDCommunityID(userID, communityID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query UserIDCommunityIDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByUserID(userID string) (out []AdaptiveCommunityUser, err error) {
	var instances []AdaptiveCommunityUser
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "UserIDIndex",
		Condition: "user_id = :a0",
		Attributes: map[string]interface{}{
			":a0": userID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByUserIDUnsafe(userID string) (out []AdaptiveCommunityUser) {
	out, err := d.ReadByUserID(userID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query UserIDIndex on %s table\n", d.Name))
	return
}


func (d DAOImpl)ReadByPlatformIDCommunityID(platformID models.PlatformID, communityID string) (out []AdaptiveCommunityUser, err error) {
	var instances []AdaptiveCommunityUser
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: "PlatformIDCommunityIDIndex",
		Condition: "platform_id = :a0 and community_id = :a1",
		Attributes: map[string]interface{}{
			":a0": platformID,
			":a1": communityID,
		},
	}, map[string]string{}, true, -1, &instances)
	out = instances
	return
}


func (d DAOImpl)ReadByPlatformIDCommunityIDUnsafe(platformID models.PlatformID, communityID string) (out []AdaptiveCommunityUser) {
	out, err := d.ReadByPlatformIDCommunityID(platformID, communityID)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query PlatformIDCommunityIDIndex on %s table\n", d.Name))
	return
}

func idParams(channelID string, userID string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"channel_id": common.DynS(channelID),
		"user_id": common.DynS(userID),
	}
	return params
}
func allParams(adaptiveCommunityUser AdaptiveCommunityUser, old AdaptiveCommunityUser) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if adaptiveCommunityUser.ChannelID != old.ChannelID { params[":a0"] = common.DynS(adaptiveCommunityUser.ChannelID) }
	if adaptiveCommunityUser.UserID != old.UserID { params[":a1"] = common.DynS(adaptiveCommunityUser.UserID) }
	if adaptiveCommunityUser.PlatformID != old.PlatformID { params[":a2"] = common.DynS(string(adaptiveCommunityUser.PlatformID)) }
	if adaptiveCommunityUser.CommunityID != old.CommunityID { params[":a3"] = common.DynS(adaptiveCommunityUser.CommunityID) }
	return
}
func updateExpression(adaptiveCommunityUser AdaptiveCommunityUser, old AdaptiveCommunityUser) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if adaptiveCommunityUser.ChannelID != old.ChannelID { updateParts = append(updateParts, "channel_id = :a0"); params[":a0"] = common.DynS(adaptiveCommunityUser.ChannelID);  }
	if adaptiveCommunityUser.UserID != old.UserID { updateParts = append(updateParts, "user_id = :a1"); params[":a1"] = common.DynS(adaptiveCommunityUser.UserID);  }
	if adaptiveCommunityUser.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a2"); params[":a2"] = common.DynS(string(adaptiveCommunityUser.PlatformID));  }
	if adaptiveCommunityUser.CommunityID != old.CommunityID { updateParts = append(updateParts, "community_id = :a3"); params[":a3"] = common.DynS(adaptiveCommunityUser.CommunityID);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}
