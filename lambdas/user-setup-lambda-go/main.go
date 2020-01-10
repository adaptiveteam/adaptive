package lambda

import (
	"context"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/userSetup"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

var (
	usersTableName = utils.NonEmptyEnv("USERS_TABLE_NAME")
)

func HandleRequest(ctx context.Context, event models.UserEngage) (string, error) {
	return userSetup.HandleUserSetupRequest(platform, userEngagementDao, event, usersTableName)
}
