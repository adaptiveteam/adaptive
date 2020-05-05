package mapping

import (
	"encoding/json"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	"github.com/jinzhu/gorm"
	"time"
)

type DBCommunityUser struct {
	ID          string `gorm:"primary_key"`
	ChannelID   string `gorm:"type:TEXT"`
	CommunityID string `gorm:"type:TEXT"`
	PlatformID  string `gorm:"type:TEXT"`
	UserID      string `gorm:"type:TEXT"`
	model.DBModel
}

// TableName return table name
func (d DBCommunityUser) TableName() string {
	return "community_user"
}

func communityUserCompositeKey(acu models.AdaptiveCommunityUser3) string {
	return acu.ChannelID + ":" + acu.UserID
}

func communityUserBMapping(acu models.AdaptiveCommunityUser3) DBCommunityUser {
	return DBCommunityUser{
		ID:          communityUserCompositeKey(acu),
		ChannelID:   acu.ChannelID,
		CommunityID: acu.CommunityID,
		PlatformID:  string(acu.PlatformID),
		UserID:      acu.UserID,
	}
}

func (d DBCommunityUser) AsAdd() (op DBCommunityUser) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBCommunityUser) AsUpdate() (op DBCommunityUser) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBCommunityUser) AsDelete() (op DBCommunityUser) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func (d DBCommunityUser) ParseUnsafe(js []byte, logger logger.AdaptiveLogger) interface{} {
	var acu models.AdaptiveCommunityUser3
	err := json.Unmarshal(js, &acu)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.AdaptiveCommunityUser3")
	}
	return acu
}

func (d DBCommunityUser) HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for community user")
	conn.AutoMigrate(&DBCommunityUser{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newCommunityUser = e2.NewEntity.(models.AdaptiveCommunityUser3)
		dbCommunityUser := communityUserBMapping(newCommunityUser).AsAdd()
		conn.Where("id = ?", dbCommunityUser.ID).
			Assign(dbCommunityUser).
			FirstOrCreate(&dbCommunityUser)
	case model.StreamEventEdit:
		var oldCommunityUser = e2.OldEntity.(models.AdaptiveCommunityUser3)
		var newCommunityUser = e2.NewEntity.(models.AdaptiveCommunityUser3)

		dbCommunityUser := communityUserBMapping(newCommunityUser).AsUpdate()
		conn.Where("id = ?", communityUserCompositeKey(oldCommunityUser)).
			Assign(dbCommunityUser).
			FirstOrCreate(&dbCommunityUser)
	case model.StreamEventDelete:
		var oldCommunityUser = e2.OldEntity.(models.AdaptiveCommunityUser3)
		var oldDbCommunityUser = communityUserBMapping(oldCommunityUser).AsDelete()
		conn.Where("id = ?", communityUserCompositeKey(oldCommunityUser)).
			First(&oldDbCommunityUser).
			Delete(&oldDbCommunityUser)
	}
}
