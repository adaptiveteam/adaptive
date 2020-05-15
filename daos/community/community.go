package community
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

type Community struct  {
	PlatformID common.PlatformID `json:"platform_id"`
	ID string `json:"id"`
	ChannelID string `json:"channel_id,omitempty"`
	CommunityKind common.CommunityKind `json:"community_kind"`
	ParentCommunityID string `json:"parent_community_id,omitempty"`
	Name string `json:"name"`
	Description string `json:"description"`
	// Owner, responsible person
	Advocate string `json:"advocate,omitempty"`
	// Nudging person
	AccountabilityPartner string `json:"accountability_partner,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	ModifiedBy string `json:"modified_by,omitempty"`
	// Automatically maintained field
	CreatedAt string `json:"created_at"`
	// Automatically maintained field
	ModifiedAt string `json:"modified_at,omitempty"`
	DeactivatedAt string `json:"deactivated_at,omitempty"`
}

// CommunityFilterActive removes deactivated values
func CommunityFilterActive(in []Community) (res []Community) {
	for _, i := range in {
		if i.DeactivatedAt == "" {
			res = append(res, i)
		}
	}
	return
}

// CollectEmptyFields returns entity field names that are empty.
// It also returns the boolean ok-flag if the list is empty.
func (community Community)CollectEmptyFields() (emptyFields []string, ok bool) {
	if community.PlatformID == "" { emptyFields = append(emptyFields, "PlatformID")}
	if community.ID == "" { emptyFields = append(emptyFields, "ID")}
	if community.CommunityKind == "" { emptyFields = append(emptyFields, "CommunityKind")}
	if community.Name == "" { emptyFields = append(emptyFields, "Name")}
	if community.Description == "" { emptyFields = append(emptyFields, "Description")}
	ok = len(emptyFields) == 0
	return
}
// ToJSON returns json string
func (community Community) ToJSON() (string, error) {
	b, err := json.Marshal(community)
	return string(b), err
}

type DAO interface {
	Create(community Community) error
	CreateUnsafe(community Community)
	Read(platformID common.PlatformID, id string) (community Community, err error)
	ReadUnsafe(platformID common.PlatformID, id string) (community Community)
	ReadOrEmpty(platformID common.PlatformID, id string) (community []Community, err error)
	ReadOrEmptyUnsafe(platformID common.PlatformID, id string) (community []Community)
	CreateOrUpdate(community Community) error
	CreateOrUpdateUnsafe(community Community)
	Deactivate(platformID common.PlatformID, id string) error
	DeactivateUnsafe(platformID common.PlatformID, id string)
	ReadByChannelIDPlatformID(channelID string, platformID common.PlatformID) (community []Community, err error)
	ReadByChannelIDPlatformIDUnsafe(channelID string, platformID common.PlatformID) (community []Community)
	ReadByPlatformIDCommunityKind(platformID common.PlatformID, communityKind common.CommunityKind) (community []Community, err error)
	ReadByPlatformIDCommunityKindUnsafe(platformID common.PlatformID, communityKind common.CommunityKind) (community []Community)
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	ConnGen   common.DynamoDBConnectionGen
}

// NewDAO creates an instance of DAO that will provide access to the table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, clientID string) DAO {
	if clientID == "" { panic(errors.New("Cannot create Community.DAO without clientID")) }
	return DAOImpl{
		ConnGen:   common.DynamoDBConnectionGen{
			Dynamo: dynamo, 
			TableNamePrefix: clientID,
		},
	}
}
// TableNameSuffixVar is a global variable that contains table name suffix.
// After renaming all tables this may be made `const`.
var TableNameSuffixVar = "_community"

// TableName concatenates table name prefix and suffix and returns table name
func TableName(prefix string) string {
	return prefix + TableNameSuffixVar
}

// Create saves the Community.
func (d DAOImpl) Create(community Community) (err error) {
	emptyFields, ok := community.CollectEmptyFields()
	if ok {
		community.ModifiedAt = core.CurrentRFCTimestamp()
	community.CreatedAt = community.ModifiedAt
	err = d.ConnGen.Dynamo.PutTableEntry(community, TableName(d.ConnGen.TableNamePrefix))
	} else {
		err = fmt.Errorf("Cannot create entity with empty fields: %v", emptyFields)
	}
	return
}


// CreateUnsafe saves the Community.
func (d DAOImpl) CreateUnsafe(community Community) {
	err2 := d.Create(community)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not create platformID==%s, id==%s in %s\n", community.PlatformID, community.ID, TableName(d.ConnGen.TableNamePrefix)))
}


// Read reads Community
func (d DAOImpl) Read(platformID common.PlatformID, id string) (out Community, err error) {
	var outs []Community
	outs, err = d.ReadOrEmpty(platformID, id)
	if err == nil && len(outs) == 0 {
		err = fmt.Errorf("Not found platformID==%s, id==%s in %s\n", platformID, id, TableName(d.ConnGen.TableNamePrefix))
	}
	if len(outs) > 0 {
		out = outs[0]
	}
	return
}


// ReadUnsafe reads the Community. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(platformID common.PlatformID, id string) Community {
	out, err2 := d.Read(platformID, id)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error reading platformID==%s, id==%s in %s\n", platformID, id, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// ReadOrEmpty reads Community
func (d DAOImpl) ReadOrEmpty(platformID common.PlatformID, id string) (out []Community, err error) {
	var outOrEmpty Community
	ids := idParams(platformID, id)
	var found bool
	found, err = d.ConnGen.Dynamo.GetItemOrEmptyFromTable(TableName(d.ConnGen.TableNamePrefix), ids, &outOrEmpty)
	if found {
		if outOrEmpty.PlatformID == platformID && outOrEmpty.ID == id {
			out = append(out, outOrEmpty)
		} else {
			err = fmt.Errorf("Requested ids: platformID==%s, id==%s are different from the found ones: platformID==%s, id==%s", platformID, id, outOrEmpty.PlatformID, outOrEmpty.ID) // unexpected error: found ids != ids
		}
	}
	err = errors.Wrapf(err, "Community DAO.ReadOrEmpty(id = %v) couldn't GetItem in table %s", ids, TableName(d.ConnGen.TableNamePrefix))
	return
}


// ReadOrEmptyUnsafe reads the Community. Panics in case of any errors
func (d DAOImpl) ReadOrEmptyUnsafe(platformID common.PlatformID, id string) []Community {
	out, err2 := d.ReadOrEmpty(platformID, id)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Error while reading platformID==%s, id==%s in %s\n", platformID, id, TableName(d.ConnGen.TableNamePrefix)))
	return out
}


// CreateOrUpdate saves the Community regardless of if it exists.
func (d DAOImpl) CreateOrUpdate(community Community) (err error) {
	community.ModifiedAt = core.CurrentRFCTimestamp()
	if community.CreatedAt == "" { community.CreatedAt = community.ModifiedAt }
	
	var olds []Community
	olds, err = d.ReadOrEmpty(community.PlatformID, community.ID)
	err = errors.Wrapf(err, "Community DAO.CreateOrUpdate(id = platformID==%s, id==%s) couldn't ReadOrEmpty", community.PlatformID, community.ID)
	if err == nil {
		if len(olds) == 0 {
			err = d.Create(community)
			err = errors.Wrapf(err, "Community DAO.CreateOrUpdate couldn't Create in table %s", TableName(d.ConnGen.TableNamePrefix))
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
				err = errors.Wrapf(err, "Community DAO.CreateOrUpdate(id = %v) couldn't UpdateTableEntry in table %s, expression='%s'", key, TableName(d.ConnGen.TableNamePrefix), expr)
			} else {
				err = fmt.Errorf("Cannot update entity with empty required fields: %v", emptyFields)
			}
		}
	}
	return 
}


// CreateOrUpdateUnsafe saves the Community regardless of if it exists.
func (d DAOImpl) CreateOrUpdateUnsafe(community Community) {
	err2 := d.CreateOrUpdate(community)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("could not create or update %v in %s\n", community, TableName(d.ConnGen.TableNamePrefix)))
}


// Deactivate "removes" Community. 
// The mechanism is adding timestamp to `DeactivatedOn` field. 
// Then, if this field is not empty, the instance is considered to be "active"
func (d DAOImpl)Deactivate(platformID common.PlatformID, id string) error {
	instance, err2 := d.Read(platformID, id)
	if err2 == nil {
		instance.DeactivatedAt = core.CurrentRFCTimestamp()
		err2 = d.CreateOrUpdate(instance)
	}
	return err2
}


// DeactivateUnsafe "deletes" Community and panics in case of errors.
func (d DAOImpl)DeactivateUnsafe(platformID common.PlatformID, id string) {
	err2 := d.Deactivate(platformID, id)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not deactivate platformID==%s, id==%s in %s\n", platformID, id, TableName(d.ConnGen.TableNamePrefix)))
}


func (d DAOImpl)ReadByChannelIDPlatformID(channelID string, platformID common.PlatformID) (out []Community, err error) {
	var instances []Community
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
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


func (d DAOImpl)ReadByChannelIDPlatformIDUnsafe(channelID string, platformID common.PlatformID) (out []Community) {
	out, err2 := d.ReadByChannelIDPlatformID(channelID, platformID)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query ChannelIDPlatformIDIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
	return
}


func (d DAOImpl)ReadByPlatformIDCommunityKind(platformID common.PlatformID, communityKind common.CommunityKind) (out []Community, err error) {
	var instances []Community
	err = d.ConnGen.Dynamo.QueryTableWithIndex(TableName(d.ConnGen.TableNamePrefix), awsutils.DynamoIndexExpression{
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


func (d DAOImpl)ReadByPlatformIDCommunityKindUnsafe(platformID common.PlatformID, communityKind common.CommunityKind) (out []Community) {
	out, err2 := d.ReadByPlatformIDCommunityKind(platformID, communityKind)
	core.ErrorHandler(err2, TableNameSuffixVar, fmt.Sprintf("Could not query PlatformIDCommunityKindIndex on %s table\n", TableName(d.ConnGen.TableNamePrefix)))
	return
}

func idParams(platformID common.PlatformID, id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue {
		"platform_id": common.DynS(string(platformID)),
		"id": common.DynS(id),
	}
	return params
}
func allParams(community Community, old Community) (params map[string]*dynamodb.AttributeValue) {
	params = map[string]*dynamodb.AttributeValue{}
	if community.PlatformID != old.PlatformID { params[":a0"] = common.DynS(string(community.PlatformID)) }
	if community.ID != old.ID { params[":a1"] = common.DynS(community.ID) }
	if community.ChannelID != old.ChannelID { params[":a2"] = common.DynS(community.ChannelID) }
	if community.CommunityKind != old.CommunityKind { params[":a3"] = common.DynS(string(community.CommunityKind)) }
	if community.ParentCommunityID != old.ParentCommunityID { params[":a4"] = common.DynS(community.ParentCommunityID) }
	if community.Name != old.Name { params[":a5"] = common.DynS(community.Name) }
	if community.Description != old.Description { params[":a6"] = common.DynS(community.Description) }
	if community.Advocate != old.Advocate { params[":a7"] = common.DynS(community.Advocate) }
	if community.AccountabilityPartner != old.AccountabilityPartner { params[":a8"] = common.DynS(community.AccountabilityPartner) }
	if community.CreatedBy != old.CreatedBy { params[":a9"] = common.DynS(community.CreatedBy) }
	if community.ModifiedBy != old.ModifiedBy { params[":a10"] = common.DynS(community.ModifiedBy) }
	if community.CreatedAt != old.CreatedAt { params[":a11"] = common.DynS(community.CreatedAt) }
	if community.ModifiedAt != old.ModifiedAt { params[":a12"] = common.DynS(community.ModifiedAt) }
	if community.DeactivatedAt != old.DeactivatedAt { params[":a13"] = common.DynS(community.DeactivatedAt) }
	return
}
func updateExpression(community Community, old Community) (expr string, params map[string]*dynamodb.AttributeValue, namesPtr *map[string]*string) {
	var updateParts []string
	params = map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	if community.PlatformID != old.PlatformID { updateParts = append(updateParts, "platform_id = :a0"); params[":a0"] = common.DynS(string(community.PlatformID));  }
	if community.ID != old.ID { updateParts = append(updateParts, "id = :a1"); params[":a1"] = common.DynS(community.ID);  }
	if community.ChannelID != old.ChannelID { updateParts = append(updateParts, "channel_id = :a2"); params[":a2"] = common.DynS(community.ChannelID);  }
	if community.CommunityKind != old.CommunityKind { updateParts = append(updateParts, "community_kind = :a3"); params[":a3"] = common.DynS(string(community.CommunityKind));  }
	if community.ParentCommunityID != old.ParentCommunityID { updateParts = append(updateParts, "parent_community_id = :a4"); params[":a4"] = common.DynS(community.ParentCommunityID);  }
	if community.Name != old.Name { updateParts = append(updateParts, "#name = :a5"); params[":a5"] = common.DynS(community.Name); fldName := "name"; names["#name"] = &fldName }
	if community.Description != old.Description { updateParts = append(updateParts, "description = :a6"); params[":a6"] = common.DynS(community.Description);  }
	if community.Advocate != old.Advocate { updateParts = append(updateParts, "advocate = :a7"); params[":a7"] = common.DynS(community.Advocate);  }
	if community.AccountabilityPartner != old.AccountabilityPartner { updateParts = append(updateParts, "accountability_partner = :a8"); params[":a8"] = common.DynS(community.AccountabilityPartner);  }
	if community.CreatedBy != old.CreatedBy { updateParts = append(updateParts, "created_by = :a9"); params[":a9"] = common.DynS(community.CreatedBy);  }
	if community.ModifiedBy != old.ModifiedBy { updateParts = append(updateParts, "modified_by = :a10"); params[":a10"] = common.DynS(community.ModifiedBy);  }
	if community.CreatedAt != old.CreatedAt { updateParts = append(updateParts, "created_at = :a11"); params[":a11"] = common.DynS(community.CreatedAt);  }
	if community.ModifiedAt != old.ModifiedAt { updateParts = append(updateParts, "modified_at = :a12"); params[":a12"] = common.DynS(community.ModifiedAt);  }
	if community.DeactivatedAt != old.DeactivatedAt { updateParts = append(updateParts, "deactivated_at = :a13"); params[":a13"] = common.DynS(community.DeactivatedAt);  }
	expr = "set " + strings.Join(updateParts, ", ")
	if len(names) == 0 { namesPtr = nil } else { namesPtr = &names } // workaround for ValidationException: ExpressionAttributeNames must not be empty
	return
}