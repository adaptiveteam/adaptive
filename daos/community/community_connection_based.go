package community
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.
import (
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"fmt"
)


// Create saves the Community.
func Create(community Community) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		emptyFields, ok := community.CollectEmptyFields()
		if ok {
			community.ModifiedAt = core.CurrentRFCTimestamp()
	community.CreatedAt = community.ModifiedAt
	
			err = conn.Dynamo.PutTableEntry(community, TableName(conn.ClientID))
		} else {
			err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
		}
		return
	}
}


// CreateUnsafe saves the Community.
func CreateUnsafe(community Community) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Create(community)(conn)
		core.ErrorHandler(err2, "daos/Community", fmt.Sprintf("Could not create platformID==%s, id==%s in %s\n", community.PlatformID, community.ID, TableName(conn.ClientID)))
	}
}


// Read reads Community
func Read(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) (out Community, err error) {
	return func (conn common.DynamoDBConnection) (out Community, err error) {
		var outs []Community
		outs, err = ReadOrEmpty(platformID, id)(conn)
		if err == nil && len(outs) == 0 {
			err = fmt.Errorf("Not found platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID))
		}
		if len(outs) > 0 {
			out = outs[0]
		}
		return
	}
}


// ReadUnsafe reads the Community. Panics in case of any errors
func ReadUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) Community {
	return func (conn common.DynamoDBConnection) Community {
		out, err2 := Read(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/Community", fmt.Sprintf("Error reading platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmpty reads Community
func ReadOrEmpty(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) (out []Community, err error) {
	return func (conn common.DynamoDBConnection) (out []Community, err error) {
       out, err = ReadOrEmptyIncludingInactive(platformID, id)(conn)
       out = CommunityFilterActive(out)
       
		return
	}
}


// ReadOrEmptyUnsafe reads the Community. Panics in case of any errors
func ReadOrEmptyUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) []Community {
	return func (conn common.DynamoDBConnection) []Community {
		out, err2 := ReadOrEmpty(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/Community", fmt.Sprintf("Error while reading platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
		return out
	}
}


// ReadOrEmptyIncludingInactive reads Community
func ReadOrEmptyIncludingInactive(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) (out []Community, err error) {
	return func (conn common.DynamoDBConnection) (out []Community, err error) {
		var outOrEmpty Community
		ids := idParams(platformID, id)
		var found bool
		found, err = conn.Dynamo.GetItemOrEmptyFromTable(TableName(conn.ClientID), ids, &outOrEmpty)
		if found {
			if outOrEmpty.PlatformID == platformID && outOrEmpty.ID == id {
				out = append(out, outOrEmpty)
			} else {
				err = fmt.Errorf("Requested ids: platformID==%s, id==%s are different from the found ones: platformID==%s, id==%s", platformID, id, outOrEmpty.PlatformID, outOrEmpty.ID) // unexpected error: found ids != ids
			}
		}
		err = errors.Wrapf(err, "Community DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(conn.ClientID))
		return
	}
}


// ReadOrEmptyIncludingInactiveUnsafe reads the Community. Panics in case of any errors
func ReadOrEmptyIncludingInactiveUnsafeIncludingInactive(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) []Community {
	return func (conn common.DynamoDBConnection) []Community {
		out, err2 := ReadOrEmptyIncludingInactive(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/Community", fmt.Sprintf("Error while reading platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
		return out
	}
}


// CreateOrUpdate saves the Community regardless of if it exists.
func CreateOrUpdate(community Community) common.ConnectionProc {
	return func (conn common.DynamoDBConnection) (err error) {
		community.ModifiedAt = core.CurrentRFCTimestamp()
	if community.CreatedAt == "" { community.CreatedAt = community.ModifiedAt }
	
		var olds []Community
		olds, err = ReadOrEmpty(community.PlatformID, community.ID)(conn)
		err = errors.Wrapf(err, "Community DAO.CreateOrUpdate(id = platformID==%s, id==%s) couldn't ReadOrEmpty", community.PlatformID, community.ID)
		if err == nil {
			if len(olds) == 0 {
				err = Create(community)(conn)
				err = errors.Wrapf(err, "Community DAO.CreateOrUpdate couldn't Create in table %s", TableName(conn.ClientID))
			} else {
				emptyFields, ok := community.CollectEmptyFields()
				if ok {
					old := olds[0]
					community.CreatedAt  = old.CreatedAt
					community.ModifiedAt = core.CurrentRFCTimestamp()
					key := idParams(old.PlatformID, old.ID)
					expr, exprAttributes, names := updateExpression(community, old)
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
					err = errors.Wrapf(err, "Community DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(conn.ClientID), expr)
				} else {
					err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
				}
			}
		}
		return 
	}
}


// CreateOrUpdateUnsafe saves the Community regardless of if it exists.
func CreateOrUpdateUnsafe(community Community) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := CreateOrUpdate(community)(conn)
		core.ErrorHandler(err2, "daos/Community", fmt.Sprintf("could not create or update %v in %s\n", community, TableName(conn.ClientID)))
	}
}


// Deactivate "removes" Community. 
// The mechanism is adding timestamp to `DeactivatedOn` field. 
// Then, if this field is not empty, the instance is considered to be "active"
func Deactivate(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) error {
	return func (conn common.DynamoDBConnection) error {
		instance, err2 := Read(platformID, id)(conn)
		if err2 == nil {
			instance.DeactivatedAt = core.CurrentRFCTimestamp()
			err2 = CreateOrUpdate(instance)(conn)
		}
		return err2
	}
}


// DeactivateUnsafe "deletes" Community and panics in case of errors.
func DeactivateUnsafe(platformID common.PlatformID, id string) func (conn common.DynamoDBConnection) {
	return func (conn common.DynamoDBConnection) {
		err2 := Deactivate(platformID, id)(conn)
		core.ErrorHandler(err2, "daos/Community", fmt.Sprintf("Could not deactivate platformID==%s, id==%s in %s\n", platformID, id, TableName(conn.ClientID)))
	}
}


func ReadByHashKeyID(id string) func (conn common.DynamoDBConnection) (out []Community, err error) {
	return func (conn common.DynamoDBConnection) (out []Community, err error) {
		var instances []Community
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			
			Condition: "id = :a",
			Attributes: map[string]interface{}{
				":a" : id,
			},
		}, map[string]string{}, true, -1, &instances)
		out = CommunityFilterActive(instances)
		return
	}
}


func ReadByHashKeyIDUnsafe(id string) func (conn common.DynamoDBConnection) (out []Community) {
	return func (conn common.DynamoDBConnection) (out []Community) {
		out, err2 := ReadByHashKeyID(id)(conn)
		core.ErrorHandler(err2, "daos/Community", fmt.Sprintf("Could not query IDPlatformIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByChannelIDPlatformID(channelID string, platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []Community, err error) {
	return func (conn common.DynamoDBConnection) (out []Community, err error) {
		var instances []Community
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "ChannelIDPlatformIDIndex",
			Condition: "channel_id = :a0 and platform_id = :a1",
			Attributes: map[string]interface{}{
				":a0": channelID,
			":a1": platformID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = CommunityFilterActive(instances)
		return
	}
}


func ReadByChannelIDPlatformIDUnsafe(channelID string, platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []Community) {
	return func (conn common.DynamoDBConnection) (out []Community) {
		out, err2 := ReadByChannelIDPlatformID(channelID, platformID)(conn)
		core.ErrorHandler(err2, "daos/Community", fmt.Sprintf("Could not query ChannelIDPlatformIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByPlatformIDCommunityKind(platformID common.PlatformID, communityKind common.CommunityKind) func (conn common.DynamoDBConnection) (out []Community, err error) {
	return func (conn common.DynamoDBConnection) (out []Community, err error) {
		var instances []Community
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: "PlatformIDCommunityKindIndex",
			Condition: "platform_id = :a0 and community_kind = :a1",
			Attributes: map[string]interface{}{
				":a0": platformID,
			":a1": communityKind,
			},
		}, map[string]string{}, true, -1, &instances)
		out = CommunityFilterActive(instances)
		return
	}
}


func ReadByPlatformIDCommunityKindUnsafe(platformID common.PlatformID, communityKind common.CommunityKind) func (conn common.DynamoDBConnection) (out []Community) {
	return func (conn common.DynamoDBConnection) (out []Community) {
		out, err2 := ReadByPlatformIDCommunityKind(platformID, communityKind)(conn)
		core.ErrorHandler(err2, "daos/Community", fmt.Sprintf("Could not query PlatformIDCommunityKindIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByHashKeyPlatformID(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []Community, err error) {
	return func (conn common.DynamoDBConnection) (out []Community, err error) {
		var instances []Community
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: string(PlatformIDCommunityKindIndex),
			Condition: "platform_id = :a",
			Attributes: map[string]interface{}{
				":a" : platformID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = CommunityFilterActive(instances)
		return
	}
}


func ReadByHashKeyPlatformIDUnsafe(platformID common.PlatformID) func (conn common.DynamoDBConnection) (out []Community) {
	return func (conn common.DynamoDBConnection) (out []Community) {
		out, err2 := ReadByHashKeyPlatformID(platformID)(conn)
		core.ErrorHandler(err2, "daos/Community", fmt.Sprintf("Could not query PlatformIDCommunityKindIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}


func ReadByHashKeyChannelID(channelID string) func (conn common.DynamoDBConnection) (out []Community, err error) {
	return func (conn common.DynamoDBConnection) (out []Community, err error) {
		var instances []Community
		err = conn.Dynamo.QueryTableWithIndex(TableName(conn.ClientID), awsutils.DynamoIndexExpression{
			IndexName: string(ChannelIDPlatformIDIndex),
			Condition: "channel_id = :a",
			Attributes: map[string]interface{}{
				":a" : channelID,
			},
		}, map[string]string{}, true, -1, &instances)
		out = CommunityFilterActive(instances)
		return
	}
}


func ReadByHashKeyChannelIDUnsafe(channelID string) func (conn common.DynamoDBConnection) (out []Community) {
	return func (conn common.DynamoDBConnection) (out []Community) {
		out, err2 := ReadByHashKeyChannelID(channelID)(conn)
		core.ErrorHandler(err2, "daos/Community", fmt.Sprintf("Could not query ChannelIDPlatformIDIndex on %s table\n", TableName(conn.ClientID)))
		return
	}
}

