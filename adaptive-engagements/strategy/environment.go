package strategy

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

var (
	namespace = utils.NonEmptyEnv("LOG_NAMESPACE")
	region    = utils.NonEmptyEnv("AWS_REGION")
	d         = awsutils.NewDynamo(region, "", namespace)
	clientID  = utils.NonEmptyEnv("CLIENT_ID")
	schema    = models.SchemaForClientID(clientID)
	userDAO   = utilsUser.NewDAOFromSchema(d, namespace, schema)
)

// UserIDToPlatformID converts userID to platformID using
// globally available variables.
func UserIDToPlatformID(userDAO utilsUser.DAO) func(string) models.PlatformID {
	return func(userID string) (platformID models.PlatformID) {
		return models.PlatformID(userDAO.ReadUnsafe(userID).PlatformId)
	}
}
