package strategy

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

var (
	namespace = func ()string { return utils.NonEmptyEnv("LOG_NAMESPACE")}
	region    = func ()string { return utils.NonEmptyEnv("AWS_REGION") }
	d         = func ()*awsutils.DynamoRequest { return awsutils.NewDynamo(region(), "", namespace()) }
	clientID  = func ()string { return utils.NonEmptyEnv("CLIENT_ID") }
	schema    = func ()models.Schema { return models.SchemaForClientID(clientID()) }
	userDAO   = func ()utilsUser.DAO { return utilsUser.NewDAOFromSchema(d(), namespace(), schema()) }
)

// UserIDToPlatformID converts userID to platformID using
// globally available variables.
func UserIDToPlatformID(userDAO utilsUser.DAO) func(string) models.PlatformID {
	return func(userID string) (platformID models.PlatformID) {
		return models.PlatformID(userDAO.ReadUnsafe(userID).PlatformId)
	}
}
