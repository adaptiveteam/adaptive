package mapping

import (
	"encoding/json"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	"github.com/jinzhu/gorm"
)

type DBCommunity struct {
	ID          string            `gorm:"primary_key"`
	Active      bool              `gorm:"type:BOOLEAN"`
	Channel     string            `gorm:"type:TEXT"`
	PlatformID  common.PlatformID `gorm:"type:TEXT"`
	RequestedBy string            `gorm:"type:TEXT"`
	CreatedAt   string            `gorm:"type:TEXT"`
	model.DBModel
}

// TableName return table name
func (d DBCommunity) TableName() string {
	return "community"
}

func communityDBMapping(ac models.AdaptiveCommunity) DBCommunity {
	return DBCommunity{
		ID:          ac.ID,
		Active:      ac.Active,
		Channel:     ac.ChannelID,
		PlatformID:  ac.PlatformID,
		RequestedBy: ac.RequestedBy,
		CreatedAt:   ac.CreatedAt,
	}
}

func (d DBCommunity) AsAdd() (op DBCommunity) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBCommunity) AsUpdate() (op DBCommunity) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBCommunity) AsDelete() (op DBCommunity) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func (d DBCommunity) ParseUnsafe(js []byte, logger logger.AdaptiveLogger) interface{} {
	var ac models.AdaptiveCommunity
	err := json.Unmarshal(js, &ac)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.AdaptiveCommunity")
	}
	return ac
}

func (d DBCommunity) HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for community")
	conn.AutoMigrate(&DBCommunity{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newCommunity = e2.NewEntity.(models.AdaptiveCommunity)
		dbCommunity := communityDBMapping(newCommunity).AsAdd()
		conn.Where("id = ?", dbCommunity.ID).
			Assign(dbCommunity).
			FirstOrCreate(&dbCommunity)
	case model.StreamEventEdit:
		var oldCommunity = e2.OldEntity.(models.AdaptiveCommunity)
		var newCommunity = e2.NewEntity.(models.AdaptiveCommunity)

		dbCommunity := communityDBMapping(newCommunity).AsUpdate()
		conn.Where("id = ?", oldCommunity.ID).
			Assign(dbCommunity).
			FirstOrCreate(&dbCommunity)
	case model.StreamEventDelete:
		var oldCommunity = e2.OldEntity.(models.AdaptiveCommunity)
		var oldDbCommunity = communityDBMapping(oldCommunity).AsDelete()
		conn.Where("id = ?", oldCommunity.ID).
			First(&oldDbCommunity).
			Delete(&oldDbCommunity)
	}
}
