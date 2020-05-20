package userAttribute


import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	// awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

// DAO - wrapper around a Dynamo DB table to work with holidays inside it
type DAO interface {
	ForUserID(userID string) UserDAO
}

// UserDAO is a dao-interface that is specific to user id.
type UserDAO interface {
	All()([]models.UserAttribute, error)
	AllUnsafe()[]models.UserAttribute
	SetUnsafe(attribute string, value string, isDefault bool)
	GetUnsafe(attribute string) models.UserAttribute
	CheckIfAllAttribsAreSet(attributes ...string) bool
}

type TableConfig struct {
	Name string
}

type DAOImpl struct {
	DNS   *common.DynamoNamespace
	TableConfig
}
type UserDAOImpl struct {
	DAO DAOImpl
	UserID string
}

func NewDAO(dns common.DynamoNamespace, tableConfig TableConfig) DAO {
	return DAOImpl {
		DNS: &dns,
		TableConfig: tableConfig,
	}
}

func (d DAOImpl)ForUserID(userID string) UserDAO {
	return UserDAOImpl {
		DAO: d,
		UserID: userID,
	}
} 

// GetUnsafe retrieves one user attribute if it's present. 
// It returns empty UserAttribute if the attribute not set.
func (d UserDAOImpl)GetUnsafe(attribute string) (userAttr models.UserAttribute) {
	queryExpr := "#uid = :uid AND #ak = :ak"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":uid": dynString(d.UserID),
		":ak": dynString(attribute),
	}
	expressionAttributeNames := map[string]string{
		"#uid": "user_id",
		"#ak": "attr_key",
	}
	var res []models.UserAttribute
	err := d.DAO.DNS.Dynamo.QueryTableWithExpr(d.DAO.TableConfig.Name, queryExpr, 
		expressionAttributeNames, 
		expressionAttributeValues, true, -1, &res)

	core.ErrorHandler(err, d.DAO.DNS.Namespace, fmt.Sprintf("Could not query %s table", d.DAO.TableConfig.Name))
	if len(res) > 0 {
		userAttr = res[0]
	}
	return 
}

// All retrieves all user attributes.
func (d UserDAOImpl)All()(res []models.UserAttribute, err error) {
	queryExpr := "#uid = :uid"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":uid": dynString(d.UserID),
	}
	expressionAttributeNames := map[string]string{
		"#uid": "user_id",
		"#ak": "attr_key",
	}
	err = d.DAO.DNS.Dynamo.QueryTableWithExpr(d.DAO.TableConfig.Name, 
		queryExpr, 
		expressionAttributeNames, 
		expressionAttributeValues, 
		true, -1, &res)
	wrapError(err, "UserDAOImpl.All")
	return 
}

// AllUnsafe returns all user attributes
func (d UserDAOImpl)AllUnsafe()[]models.UserAttribute {
	all, err := d.All()
	core.ErrorHandler(err, d.DAO.DNS.Namespace, fmt.Sprintf("Could not query %s table", d.DAO.TableConfig.Name))	
	return all
}

// SetUnsafe sets attribute value
func (d UserDAOImpl)SetUnsafe(attribute string, value string, isDefault bool) {
	userAttr := models.UserAttribute{
		UserID:    d.UserID,
		AttrKey:   attribute,
		AttrValue: value,
		Default:   isDefault,
	}
	err2 := d.DAO.DNS.Dynamo.PutTableEntry(userAttr, d.DAO.TableConfig.Name)
	core.ErrorHandler(err2, d.DAO.DNS.Namespace, fmt.Sprintf("Could not save %s=%s to %s table for user %s", attribute, value, d.DAO.TableConfig.Name, d.UserID))
}

// CheckIfAllAttribsAreSet checks that all given attributes are present
func (d UserDAOImpl)CheckIfAllAttribsAreSet(attributes ...string) bool {
	actualAttributes := d.AllUnsafe()
	m := make(map[string]bool)
	for _, attrib := range actualAttributes {
		m[attrib.AttrKey] = true
	}
	for _, attrib := range attributes {
		if !m[attrib] {
			return false
		}
	}
	return true
}

func dynString(str string) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{S: aws.String(str)}
	return &attr
}

func wrapError(err error, name string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("{%s: %v}", name, err)
}
