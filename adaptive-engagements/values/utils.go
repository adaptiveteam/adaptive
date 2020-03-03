package values

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

func valuesDAO() DAO {
	clientID               := utils.NonEmptyEnv("CLIENT_ID")
	schema                 := models.SchemaForClientID(clientID)
	dns                    := common.DeprecatedGetGlobalDns()
	valuesDao              := NewDAOFromSchema(&dns, schema)
	return valuesDao
}

func PlatformValues(teamID models.TeamID) []models.AdaptiveValue {
	return valuesDAO().ForPlatformID(teamID.ToPlatformID()).AllUnsafe()
}
