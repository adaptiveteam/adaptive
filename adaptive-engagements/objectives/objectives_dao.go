package objectives

import (
	// "fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	// "time"
	// awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	// core "github.com/adaptiveteam/adaptive/core-utils-go"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
)

// DAO - wrapper around a Dynamo DB table to work with objectives inside it
type DAO interface {
	// All() ([]models.UserObjective, error)
	// AllUnsafe() []models.UserObjective
	Create(objective models.UserObjective) error
}

// TableConfig is a structure that contains all configuration information from environment
type TableConfig struct {
	Table                  string
	ProgressTable          string
	UserIDIndex            string
	ProgressIDIndex        string
	ProgressCreatedOnIndex string
	AcceptedIndex          string
	PartnerIndex           string
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	DNS *common.DynamoNamespace
	TableConfig
}

// PlatformDAOImpl DAO that will implement PlatformDAO interface
type PlatformDAOImpl struct {
	PlatformID string
	DAOImpl
}

const (
	userObjectivesTableKey                  = "USER_OBJECTIVES_TABLE"
	userObjectivesProgressTableKey          = "USER_OBJECTIVES_PROGRESS_TABLE"
	userObjectivesUserIDIndexKey            = "USER_OBJECTIVES_USER_ID_INDEX"
	userObjectivesProgressIDIndexKey        = "USER_OBJECTIVES_PROGRESS_ID_INDEX"
	userObjectivesProgressCreatedOnIndexKey = "USER_OBJECTIVES_PROGRESS_CREATED_ON_INDEX"
	userObjectivesAcceptedIndexKey          = "USER_OBJECTIVES_ACCEPTED_INDEX"
	userObjectivesPartnerIndexKey           = "USER_OBJECTIVES_PARTNER_INDEX"
)

func nonEmpty(env func(string) string) func(string) string {
	return func(key string) string {
		value := env(key)
		if value == "" {
			panic("Key " + key + " is not defined")
		}
		return value
	}
}

// ReadTableConfigUnsafe reads configuration values from environment
func ReadTableConfigUnsafe(env func(string) string) TableConfig {
	senv := nonEmpty(env)
	userObjectivesTable := senv(userObjectivesTableKey)
	userObjectivesProgressTable := senv(userObjectivesProgressTableKey)
	userObjectivesUserIDIndex := senv(userObjectivesUserIDIndexKey)
	userObjectivesProgressIDIndex := senv(userObjectivesProgressIDIndexKey)
	userObjectivesProgressCreatedOnIndex := senv(userObjectivesProgressCreatedOnIndexKey)
	userObjectivesAcceptedIndex := senv(userObjectivesAcceptedIndexKey)
	userObjectivesPartnerIndex := senv(userObjectivesPartnerIndexKey)
	return TableConfig{
		Table:                  userObjectivesTable,
		ProgressTable:          userObjectivesProgressTable,
		UserIDIndex:            userObjectivesUserIDIndex,
		ProgressIDIndex:        userObjectivesProgressIDIndex,
		ProgressCreatedOnIndex: userObjectivesProgressCreatedOnIndex,
		AcceptedIndex:          userObjectivesAcceptedIndex,
		PartnerIndex:           userObjectivesPartnerIndex,
	}
}

// NewDAO creates an instance of DAO that will provide access to objectives table
func NewDAO(dns *common.DynamoNamespace, tableConfig TableConfig) DAO {
	return DAOImpl{DNS: dns, TableConfig: tableConfig}
}

// // All reads all objectives from dynamo table
// func (d DAOImpl) All() ([]models.UserObjective, error) {
//	var objectives []models.UserObjective
//	err := d.DNS.Dynamo.ScanTable(d.Table, &objectives)
//	return objectives, err
// }
//
// // AllUnsafe reads all objectives from dynamo table. Panics in case of errors
// func (d DAOImpl) AllUnsafe() []models.UserObjective {
//	objectives, err := d.All()
//	core.ErrorHandler(err, d.DNS.Namespace, fmt.Sprintf("Could not query table %s", d.Table))
//	return objectives
//
// }

// Create puts the objective into the table
func (d DAOImpl) Create(objective models.UserObjective) error {
	return d.DNS.Dynamo.PutTableEntry(objective, d.Table)
}
