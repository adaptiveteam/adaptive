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

type DBUser struct {
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

// TableName return table name
func (d DBUser) TableName() string {
	return "user"
}

func userDBMapping(user models.User) DBUser {
	return DBUser{
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

func (d DBUser) AsAdd() (op DBUser) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBUser) AsUpdate() (op DBUser) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBUser) AsDelete() (op DBUser) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func (d DBUser) ParseUnsafe(js []byte, logger logger2.AdaptiveLogger) interface{} {
	var user models.User
	err := json.Unmarshal(js, &user)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.User")
	}
	return user
}

func (d DBUser) HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger2.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for user")
	conn.AutoMigrate(&DBUser{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newUser = e2.NewEntity.(models.User)
		DBUser := userDBMapping(newUser).AsAdd()
		conn.Where("id = ?", DBUser.ID).
			Assign(DBUser).
			FirstOrCreate(&DBUser)
	case model.StreamEventEdit:
		var oldUser = e2.OldEntity.(models.User)
		var newUser = e2.NewEntity.(models.User)

		DBUser := userDBMapping(newUser).AsUpdate()
		conn.Where("id = ?", oldUser.ID).
			Assign(DBUser).
			FirstOrCreate(&DBUser)
	case model.StreamEventDelete:
		var oldUser = e2.OldEntity.(models.User)
		var oldDBUser = userDBMapping(oldUser).AsDelete()
		conn.Where("id = ?", oldUser.ID).
			First(&oldDBUser).
			Delete(&oldDBUser)
	}
}
