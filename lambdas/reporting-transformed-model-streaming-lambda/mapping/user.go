package mapping

import (
	"encoding/json"
	"time"

	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	logger2 "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/jinzhu/gorm"
	"github.com/adaptiveteam/adaptive/daos/common"

)

type DbUser struct {
	ID                         string `gorm:"primary_key"`
	PlatformID                 common.PlatformID `gorm:"type:CHAR(9)"`
	FirstName                  string `gorm:"type:TEXT"`
	LastName                   string `gorm:"type:TEXT"`
	DisplayName                string `gorm:"type:TEXT"`
	Timezone                   string `gorm:"type:TEXT"`
	TimezoneOffset             int32  `gorm:"type:INTEGER"`
	IsAdmin                    bool   `gorm:"type:BOOLEAN"`
	IsShared                   bool   `gorm:"type:BOOLEAN"`
	AdaptiveScheduledTime      string `gorm:"type:CHAR(4)"`
	AdaptiveScheduledTimeInUTC string `gorm:"type:CHAR(4)"`
	CreatedAt                  string `gorm:"type:TEXT"`
	model.DBModel
}

func userDBMapping(user models.User) DbUser {
	return DbUser{
		ID:                         user.ID,
		PlatformID:                 user.PlatformID,
		FirstName:                  user.FirstName,
		LastName:                   user.LastName,
		DisplayName:                user.DisplayName,
		Timezone:                   user.Timezone,
		TimezoneOffset:             int32(user.TimezoneOffset),
		IsAdmin:                    user.IsAdmin,
		IsShared:                   user.IsShared,
		AdaptiveScheduledTime:      user.AdaptiveScheduledTime,
		AdaptiveScheduledTimeInUTC: user.AdaptiveScheduledTimeInUTC,
		CreatedAt:                  user.CreatedAt,
	}
}

func (d DbUser) AsAdd() (op DbUser) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DbUser) AsUpdate() (op DbUser) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DbUser) AsDelete() (op DbUser) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func InterfaceToUserUnsafe(ip interface{}, logger logger2.AdaptiveLogger) interface{} {
	var user models.User
	js, _ := json.Marshal(ip)
	err := json.Unmarshal(js, &user)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.User")
	}
	return user
}

func UserStreamEntityHandler(e2 model.StreamEntity, conn *gorm.DB, logger logger2.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for user")
	conn.AutoMigrate(&DbUser{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newUser = e2.NewEntity.(models.User)
		dbUser := userDBMapping(newUser).AsAdd()
		conn.Where("id = ?", dbUser.ID).
			Assign(dbUser).
			FirstOrCreate(&dbUser)
	case model.StreamEventEdit:
		var oldUser = e2.OldEntity.(models.User)
		var newUser = e2.NewEntity.(models.User)

		dbUser := userDBMapping(newUser).AsUpdate()
		conn.Where("id = ?", oldUser.ID).
			Assign(dbUser).
			FirstOrCreate(&dbUser)
	case model.StreamEventDelete:
		var oldUser = e2.OldEntity.(models.User)
		var oldDbUser = userDBMapping(oldUser).AsDelete()
		conn.Where("id = ?", oldUser.ID).
			First(&oldDbUser).
			Delete(&oldDbUser)
	}
}