package competencies

import (
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"strconv"
)

func responses(notifications ...models.PlatformSimpleNotification) []models.PlatformSimpleNotification {
	return notifications
}

func callbackID(userID string, action string) models.MessageCallback {
	year, month := core.CurrentYearMonth()
	return models.MessageCallback{
		Module: AdaptiveValuesNamespace,
		Source: userID,
		Topic:  "AdaptiveValueManagement",
		Action: action,
		Month:  strconv.Itoa(int(month)),
		Year:   strconv.Itoa(year),
	}

}