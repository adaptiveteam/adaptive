package platform

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core_utils_go "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/slack-go/slack"
)

// GetSlackClientUnsafe reads token and creates slack api
func GetSlackClientUnsafe(conn common.DynamoDBConnection) *slack.Client {
	ut, err2 := GetToken(models.ParseTeamID(conn.PlatformID))(conn)
	core_utils_go.ErrorHandler(err2, "GetSlackClientUnsafe", "GetToken")
	return slack.New(ut)
}
