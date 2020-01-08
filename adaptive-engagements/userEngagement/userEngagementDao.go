package userEngagement

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	// "time"
	// awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
	// awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

// DAO is a CRUD interface for working with engagements
type DAO interface {
	CreateUnsafe(eng models.UserEngagement)
}

// TableConfig contains information about table name and indices
type TableConfig struct {
	Name string
}

// DAOImpl is an implementation of DAO
type DAOImpl struct {
	DNS   *common.DynamoNamespace
	TableConfig
}

// CreateUnsafe creates user engagement
func (impl DAOImpl)CreateUnsafe(eng models.UserEngagement) {
	err := impl.DNS.Dynamo.PutTableEntry(eng, impl.Name)
	core.ErrorHandler(err, impl.DNS.Namespace, fmt.Sprintf("Could not write to %s table", impl.Name))
}

// NewDAO creates dao for working with userEngagementTable
func NewDAO(dns common.DynamoNamespace, userEngagementTableName string) DAO {
	return DAOImpl{
		DNS: &dns,
		TableConfig: TableConfig{
			Name: userEngagementTableName,
		},
	}
}
