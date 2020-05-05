package adaptiveCommunityUser
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

type AdaptiveCommunityUser struct  {
	ChannelID string `json:"channel_id"`
	// UserID is the ID of the user to send an engagement to
	// This usually corresponds to the platform user id
	UserID string `json:"user_id"`
	PlatformID common.PlatformID `json:"platform_id"`
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
// ToJSON returns json string
func (adaptiveCommunityUser AdaptiveCommunityUser) ToJSON() (string, error) {
	b, err := json.Marshal(adaptiveCommunityUser)
	return string(b), err
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
	ReadByPlatformIDCommunityID(platformID common.PlatformID, communityID string) (adaptiveCommunityUser []AdaptiveCommunityUser, err error)
	ReadByPlatformIDCommunityIDUnsafe(platformID common.PlatformID, communityID string) (adaptiveCommunityUser []AdaptiveCommunityUser)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	ConnGen   common.DynamoDBConnectionGen
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic(errors.New("Cannot create AdaptiveCommunityUser.DAO without clientID")) }
	return DAOImpl{
		ConnGen:   common.DynamoDBConnectionGen{
			Dynamo: dynamo, 
			TableNamePrefix: clientID,
		},
	}
}
// TableNameSuffixVar is a global variable that contains table name suffix.
// After renaming all tables this may be made `const`.
var TableNameSuffixVar = "_adaptive_community_user"

// TableName concatenates table name prefix and suffix and returns table name
func TableName(prefix string) string {
	return prefix + TableNameSuffixVar
}

// Create saves the AdaptiveCommunityUser.
func (d DAOImpl) Create(adaptiveCommunityUser AdaptiveCommunityUser) (err error) {
	emptyFields, ok := adaptiveCommunityUser.CollectEmptyFields()
	if ok {
		err = d.ConnGen.Dynamo.PutTableEntry(adaptiveCommunityUser, TableName(d.ConnGen.TableNamePrefix))
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the AdaptiveCommunityUser.
func (d DAOImpl) CreateUnsafe(adaptiveCommunityUser AdaptiveCommunityUser) {
	err2 := d.Create(adaptiveCommunityUser)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not create channelID==%s, userID==%s in %s\n", adaptiveCommunityUser.ChannelID, adaptiveCommunityUser.UserID, TableName(d.ConnGen.TableNamePrefix)))
}


// Read reads AdaptiveCommunityUser
func (d DAOImpl) Read(channelID string, userID string) (out AdaptiveCommunityUser, err error) {
	var outs []AdaptiveCommunityUser
	outs, err = d.ReadOrEmpty(channelID, userID)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found channelID==%s, userID==%s in %s\n", channelID, userID, TableName(d.ConnGen.TableNamePrefix))
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the AdaptiveCommunityUser. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(channelID string, userID string) AdaptiveCommunityUser {
	out, err2 := d.Read(channelID, userID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error reading channelID==%s, userID==%s in %s\n", channelID, userID, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// ReadOrEmpty reads AdaptiveCommunityUser
func (d DAOImpl) ReadOrEmpty(channelID string, userID string) (out []AdaptiveCommunityUser, err error) {
	var outOrEmpty AdaptiveCommunityUser
	ids := idParams(channelID, userID)
	var found bool
	found, err = d.ConnGen.Dynamo.GetItemOrEmptyFromTable(TableName(d.ConnGen.TableNamePrefix), ids, &outOrEmpty)
	if found {
		if outOrEmpty.ChannelID == channelID && outOrEmpty.UserID == userID {
			out = append(out, outOrEmpty)
		} else {
			err = fmt.Errorf("Requested ids: channelID==%s, userID==%s are different from the found ones: channelID==%s, userID==%s", channelID, userID, outOrEmpty.ChannelID, outOrEmpty.UserID) // unexpected error: found ids != ids
		}
	}
	err = errors.Wrapf(err, "AdaptiveCommunityUser DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(d.ConnGen.TableNamePrefix))
	return
}


// ReadOrEmptyUnsafe reads the AdaptiveCommunityUser. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(channelID string, userID string) []AdaptiveCommunityUser {
	out, err2 := d.ReadOrEmpty(channelID, userID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error while reading channelID==%s, userID==%s in %s\n", channelID, userID, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// CreateOrUpdate saves the AdaptiveCommunityUser regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(adaptiveCommunityUser AdaptiveCommunityUser) (err error) {
	
	var olds []AdaptiveCommunityUser
	olds, err = d.ReadOrEmpty(adaptiveCommunityUser.ChannelID, adaptiveCommunityUser.UserID)
	err = errors.Wrapf(err, "AdaptiveCommunityUser DAO.CreateOrUpdate(id = channelID==%s, userID==%s) couldn't ReadOrEmpty", adaptiveCommunityUser.ChannelID, adaptiveCommunityUser.UserID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(adaptiveCommunityUser)
			err = errors.Wrapf(err, "AdaptiveCommunityUser DAO.CreateOrUpdate couldn't Create in table %s", TableName(d.ConnGen.TableNamePrefix))
		} else {
			emptyFields, ok := adaptiveCommunityUser.CollectEmptyFields()
			if ok {
				old := olds[0]
				
				
				key := idParams(old.ChannelID, old.UserID)
				expr, exprAttributes, names := updateExpression(adaptiveCommunityUser, old)
				input := dynamodb.UpdateItemInput{
					ExpressionAttributeValues: exprAttributes,
					TableName:                 aws.String(TableName(d.ConnGen.TableNamePrefix)),
					Key:                       key,
					ReturnValues:              aws.String("UPDATED_NEW"),
					UpdateExpression:          aws.String(expr),
				}
				if names != nil { input.ExpressionAttributeNames = *names } // workaround for a pointer to an empty slice
				if  len(exprAttributes) > 0 { // if there some changes
					err = d.ConnGen.Dynamo.UpdateItemInternal(input)
				} else {
					// WARN: no changes.
				}
				err = errors.Wrapf(err, "AdaptiveCommunityUser DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(d.ConnGen.TableNamePrefix), expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the AdaptiveCommunityUser regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(adaptiveCommunityUser AdaptiveCommunityUser) {
	err2 := d.CreateOrUpdate(adaptiveCommunityUser)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("could not create or update %v in %s\n", adaptiveCommunityUser, TableName(d.ConnGen.TableNamePrefix)))
}


// Delete removes AdaptiveCommunityUser from db
func (d DAOImpl)Delete(channelID string, userID string) error {
	return d.ConnGen.Dynamo.DeleteEntry(TableName(d.ConnGen.TableNamePrefix), idParams(channelID, userID))
}


// DeleteUnsafe deletes AdaptiveCommunityUser and panics in case of errors.
func (d DAOImpl)DeleteUnsafe(channelID string, userID string) {
	err2 := d.Delete(channelID, userID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not delete channelID==%s, userID==%s in %s\n", channelID, userID, TableName(d.ConnGen.TableNamePrefix)))
}


func (d DAOImpl)ReadByChannelID(channelID string) (out []AdaptiveCommunityUser, err error) {
	var instances []AdaptiveCommunityUser
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
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
	out, err2 := d.ReadByChannelID(channelID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query ChannelIDIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
	return
}


func (d DAOImpl)ReadByUserIDCommunityID(userID string, communityID string) (out []AdaptiveCommunityUser, err error) {
	var instances []AdaptiveCommunityUser
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
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
	out, err2 := d.ReadByUserIDCommunityID(userID, communityID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query UserIDCommunityIDIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
	return
}


func (d DAOImpl)ReadByUserID(userID string) (out []AdaptiveCommunityUser, err error) {
	var instances []AdaptiveCommunityUser
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
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
	out, err2 := d.ReadByUserID(userID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query UserIDIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
	return
}


func (d DAOImpl)ReadByPlatformIDCommunityID(platformID common.PlatformID, communityID string) (out []AdaptiveCommunityUser, err error) {
	var instances []AdaptiveCommunityUser
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
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


func (d DAOImpl)ReadByPlatformIDCommunityIDUnsafe(platformID common.PlatformID, communityID string) (out []AdaptiveCommunityUser) {
	out, err2 := d.ReadByPlatformIDCommunityID(platformID, communityID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query PlatformIDCommunityIDIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
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
