package strategy

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
)

var (
	namespace = func() string { return utils.NonEmptyEnv("LOG_NAMESPACE") }
	region    = func() string { return utils.NonEmptyEnv("AWS_REGION") }
	d         = func() *awsutils.DynamoRequest { return awsutils.NewDynamo(region(), "", namespace()) }
	clientID  = func() string { return utils.NonEmptyEnv("CLIENT_ID") }
	schema    = func() models.Schema { return models.SchemaForClientID(clientID()) }
	userDAO   = func() utilsUser.DAO { return utilsUser.NewDAOFromSchema(d(), namespace(), schema()) }
)

// UserIDToTeamID converts userID to teamID using
// globally available variables.
func UserIDToTeamID(userDAO utilsUser.DAO) func(string) models.TeamID {
	return func(userID string) (teamID models.TeamID) {
		return models.ParseTeamID(userDAO.ReadUnsafe(userID).PlatformID)
	}
}

// UserIDToPlatformID converts userID to teamID using
// globally available variables.
func UserIDToPlatformID(userDAO utilsUser.DAO) func(string) daosCommon.PlatformID {
	return func(userID string) daosCommon.PlatformID {
		return userDAO.ReadUnsafe(userID).PlatformID
	}
}
