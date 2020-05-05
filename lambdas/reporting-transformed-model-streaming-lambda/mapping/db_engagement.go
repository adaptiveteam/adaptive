package mapping

import (
	"encoding/json"
	"time"

	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/jinzhu/gorm"
)

type DBEngagement struct {
	ID         string `gorm:"primary_key"`
	Answered   bool   `gorm:"type:BOOLEAN"`
	Ignored    bool   `gorm:"type:BOOLEAN"`
	PlatformID string `gorm:"type:TEXT"`
	Priority   string `gorm:"type:CHAR(6)"`
	TargetId   string `gorm:"type:TEXT"`
	UserID     string `gorm:"type:TEXT"`
	CreatedAt  string `gorm:"type:TEXT"`
	model.DBModel
}

// TableName return table name
func (d DBEngagement) TableName() string {
	return "engagement"
}

func engagementDBMapping(eng models.UserEngagement) DBEngagement {
	return DBEngagement{
		ID:         eng.ID,
		Answered:   intToBoolean(eng.Answered),
		Ignored:    intToBoolean(eng.Ignored),
		PlatformID: string(eng.PlatformID),
		Priority:   string(eng.Priority),
		TargetId:   eng.TargetID,
		UserID:     eng.UserID,
		CreatedAt:  eng.CreatedAt,
	}
}

func (d DBEngagement) AsAdd() (op DBEngagement) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBEngagement) AsUpdate() (op DBEngagement) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBEngagement) AsDelete() (op DBEngagement) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func (d DBEngagement) ParseUnsafe(js []byte, logger logger.AdaptiveLogger) interface{} {
	var eng models.UserEngagement
	err := json.Unmarshal(js, &eng)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.UserEngagement")
	}
	return eng
}

func (d DBEngagement) HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for engagement")
	conn.AutoMigrate(&DBEngagement{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newEng = e2.NewEntity.(models.UserEngagement)
		dbEng := engagementDBMapping(newEng).AsAdd()
		conn.Where("id = ?", dbEng.ID).
			Assign(dbEng).
			FirstOrCreate(&dbEng)
	case model.StreamEventEdit:
		var oldEng = e2.OldEntity.(models.UserEngagement)
		var newEng = e2.NewEntity.(models.UserEngagement)

		dbEng := engagementDBMapping(newEng).AsUpdate()
		conn.Where("id = ?", oldEng.ID).
			Assign(dbEng).
			FirstOrCreate(&dbEng)
	case model.StreamEventDelete:
		var oldEng = e2.OldEntity.(models.UserEngagement)
		var oldDbEng = engagementDBMapping(oldEng).AsDelete()
		conn.Where("id = ?", oldEng.ID).
			First(&oldDbEng).
			Delete(&oldDbEng)
	}
}
