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
)


// Create saves the AdaptiveCommunityUser.
func Create(adaptiveCommunityUser AdaptiveCommunityUser) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := adaptiveCommunityUser.CollectEmptyFields()
		if ok {
			
			err = conn.Dynamo.PutTableEntry(adaptiveCommunityUser, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the AdaptiveCommunityUser.
func CreateUnsafe(adaptiveCommunityUser AdaptiveCommunityUser) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(adaptiveCommunityUser)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("Could not create channelID==%s, userID==%s in %s\n", adaptiveCommunityUser.ChannelID, adaptiveCommunityUser.UserID, TableName(conn.ClientID)))
	}
}


// Read reads AdaptiveCommunityUser
func Read(channelID string, userID string) func (conn common.DynamoDBConnection) (out AdaptiveCommunityUser, err error) {
	return func (conn common.DynamoDBConnection) (out AdaptiveCommunityUser, err error) {
		var outs []AdaptiveCommunityUser
		outs, err = ReadOrEmpty(channelID, userID)(conn)
		if err == nil && len(outs) == 0 {
			err = fmt.Errorf("Not found channelID==%s, userID==%s in %s\n", channelID, userID, TableName(conn.ClientID))
		}
		if len(outs) > 0 {
			out = outs[0]
		}
		return
	}
}


// ReadUnsafe reads the AdaptiveCommunityUser. Panics in case of any errors
func ReadUnsafe(channelID string, userID string) func (conn common.DynamoDBConnection) AdaptiveCommunityUser {
	return func (conn common.DynamoDBConnection) AdaptiveCommunityUser {
		out, err2 := Read(channelID, userID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("Error reading channelID==%s, userID==%s in %s\n", channelID, userID, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads AdaptiveCommunityUser
func ReadOrEmpty(channelID string, userID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
       out, err = ReadOrEmptyIncludingInactive(channelID, userID)(conn)
       
       
		return
	}
}


// ReadOrEmptyUnsafe reads the AdaptiveCommunityUser. Panics in case of any errors
func ReadOrEmptyUnsafe(channelID string, userID string) func (conn common.DynamoDBConnection) []AdaptiveCommunityUser {
	return func (conn common.DynamoDBConnection) []AdaptiveCommunityUser {
		out, err2 := ReadOrEmpty(channelID, userID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("Error while reading channelID==%s, userID==%s in %s\n", channelID, userID, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads AdaptiveCommunityUser
func ReadOrEmptyIncludingInactive(channelID string, userID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
		var outOrEmpty AdaptiveCommunityUser
		ids := idParams(channelID, userID)
		var found bool
		found, err = conn.Dynamo.GetItemOrEmptyFromTable(TableName(conn.ClientID), ids, &outOrEmpty)
		if found {
			if outOrEmpty.ChannelID == channelID && outOrEmpty.UserID == userID {
				out = append(out, outOrEmpty)
			} else {
				err = fmt.Errorf("Requested ids: channelID==%s, userID==%s are different from the found ones: channelID==%s, userID==%s", channelID, userID, outOrEmpty.ChannelID, outOrEmpty.UserID) // unexpected error: found ids != ids
			}
		}
		err = errors.Wrapf(err, "AdaptiveCommunityUser DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the AdaptiveCommunityUser. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(channelID string, userID string) func (conn common.DynamoDBConnection) []AdaptiveCommunityUser {
	return func (conn common.DynamoDBConnection) []AdaptiveCommunityUser {
		out, err2 := ReadOrEmptyIncludingInactive(channelID, userID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("Error while reading channelID==%s, userID==%s in %s\n", channelID, userID, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the AdaptiveCommunityUser regardless of if it exists.
func CreateOrUpdate(adaptiveCommunityUser AdaptiveCommunityUser) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		
		var olds []AdaptiveCommunityUser
		olds, err = ReadOrEmpty(adaptiveCommunityUser.ChannelID, adaptiveCommunityUser.UserID)(conn)
		err = errors.Wrapf(err, "AdaptiveCommunityUser DAO.CreateOrUpdate(id = channelID==%s, userID==%s) couldn't ReadOrEmpty", adaptiveCommunityUser.ChannelID, adaptiveCommunityUser.UserID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(adaptiveCommunityUser)(conn)
				err = errors.Wrapf(err, "AdaptiveCommunityUser DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := adaptiveCommunityUser.CollectEmptyFields()
				if ok {
					old := olds[0]
					
					
					key := idParams(old.ChannelID, old.UserID)
					expr, exprAttributes, names := updateExpression(adaptiveCommunityUser, old)
					input := dynamodb.UpdateItemInput{
						ExpressionAttributeValues: exprAttributes,
						TableName:                 aws.String(TableName(conn.ClientID)),
						Key:                       key,
						ReturnValues:              aws.String("UPDATED_NEW"),
						UpdateExpression:          aws.String(expr),
					}
					if names != nil { input.ExpressionAttributeNames = *names } // workaround for a pointer to an empty slice
					if  len(exprAttributes) > 0 { // if there some changes
						err = conn.Dynamo.UpdateItemInternal(input)
					} else {
						// WARN: no changes.
					}
					err = errors.Wrapf(err, "AdaptiveCommunityUser DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the AdaptiveCommunityUser regardless of if it exists.
func CreateOrUpdateUnsafe(adaptiveCommunityUser AdaptiveCommunityUser) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(adaptiveCommunityUser)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("could not create or update %v in %s\n", adaptiveCommunityUser, TableName(conn.ClientID)))
	}
}


// Delete removes AdaptiveCommunityUser from db
func Delete(channelID string, userID string) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		return conn.Dynamo.DeleteEntry(TableName(conn.ClientID), idParams(channelID, userID))
	}
}


// DeleteUnsafe deletes AdaptiveCommunityUser and panics in case of errors.
func DeleteUnsafe(channelID string, userID string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Delete(channelID, userID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("Could not delete channelID==%s, userID==%s in %s\n", channelID, userID, TableName(conn.ClientID)))
	}
}


func ReadByChannelID(channelID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
		var instances []AdaptiveCommunityUser
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "ChannelIDIndex",
			Condition: "channel_id = :a0",
			Attributes: map[string]interface{}{
				":a0": channelID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByChannelIDUnsafe(channelID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
		out, err2 := ReadByChannelID(channelID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("Could not query ChannelIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByUserIDCommunityID(userID string, communityID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
		var instances []AdaptiveCommunityUser
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
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
}


func ReadByUserIDCommunityIDUnsafe(userID string, communityID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
		out, err2 := ReadByUserIDCommunityID(userID, communityID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("Could not query UserIDCommunityIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByUserID(userID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
		var instances []AdaptiveCommunityUser
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "UserIDIndex",
			Condition: "user_id = :a0",
			Attributes: map[string]interface{}{
				":a0": userID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByUserIDUnsafe(userID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
		out, err2 := ReadByUserID(userID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("Could not query UserIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByPlatformIDCommunityID(platformID common.PlatformID, communityID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
		var instances []AdaptiveCommunityUser
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
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
}


func ReadByPlatformIDCommunityIDUnsafe(platformID common.PlatformID, communityID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
		out, err2 := ReadByPlatformIDCommunityID(platformID, communityID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("Could not query PlatformIDCommunityIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByHashKeyPlatformID(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
		var instances []AdaptiveCommunityUser
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "PlatformIDCommunityIDIndex",
			Condition: "platform_id = :a",
			Attributes: map[string]interface{}{
				":a" : platformID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByHashKeyPlatformIDUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
		out, err2 := ReadByHashKeyPlatformID(platformID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("Could not query PlatformIDCommunityIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByHashKeyUserID(userID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser, err error) {
		var instances []AdaptiveCommunityUser
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "UserIDCommunityIDIndex",
			Condition: "user_id = :a",
			Attributes: map[string]interface{}{
				":a" : userID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = instances
		return
	}
}


func ReadByHashKeyUserIDUnsafe(userID string) func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
	return func (conn common.DynamoDBConnection) (out []AdaptiveCommunityUser) {
		out, err2 := ReadByHashKeyUserID(userID)(conn)
		core.ErrorHandler(err2, "daos/AdaptiveCommunityUser", fmt.Sprintf("Could not query UserIDCommunityIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}

