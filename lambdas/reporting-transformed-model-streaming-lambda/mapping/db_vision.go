package mapping

import (
	"encoding/json"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	logger2 "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/jinzhu/gorm"
	"github.com/adaptiveteam/adaptive/daos/common"

	"time"
)

type DBVision struct {
	ID          string `gorm:"primary_key"`
	Vision      string `gorm:"type:TEXT"`
	PlatformID  common.PlatformID `gorm:"type:TEXT"`
	Advocate    string `gorm:"type:TEXT"`
	CreatedTime string `gorm:"type:TIMESTAMP"`
	CreatedBy   string `gorm:"type:TEXT"`
	CreatedAt   string `gorm:"type:TEXT"`
	model.DBModel
}

// TableName return table name
func (d DBVision) TableName() string {
	return "vision"
}

func visionDBMapping(vis models.VisionMission) DBVision {
	return DBVision{
		ID:          vis.ID,
		Vision:      vis.Vision,
		PlatformID:  vis.PlatformID,
		Advocate:    vis.Advocate,
		CreatedTime: vis.CreatedAt,
		CreatedBy:   vis.CreatedBy,
		CreatedAt:   vis.CreatedAt,
	}
}

func (d DBVision) AsAdd() (op DBVision) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBVision) AsUpdate() (op DBVision) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBVision) AsDelete() (op DBVision) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func (d DBVision) ParseUnsafe(js []byte, logger logger2.AdaptiveLogger) interface{} {
	var vision models.VisionMission
	err := json.Unmarshal(js, &vision)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.VisionMission")
	}
	return vision
}

func (d DBVision) HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger2.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for vision")
	conn.AutoMigrate(&DBVision{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newVision = e2.NewEntity.(models.VisionMission)
		dbVision := visionDBMapping(newVision).AsAdd()
		conn.Where("id = ?", dbVision.ID).
			Assign(dbVision).
			FirstOrCreate(&dbVision)
	case model.StreamEventEdit:
		var oldVision = e2.OldEntity.(models.VisionMission)
		var newVision = e2.NewEntity.(models.VisionMission)

		dbVision := visionDBMapping(newVision).AsUpdate()
		conn.Where("id = ?", oldVision.ID).
			Assign(dbVision).
			FirstOrCreate(&dbVision)
	case model.StreamEventDelete:
		var oldVision = e2.OldEntity.(models.VisionMission)
		var oldDBVision = visionDBMapping(oldVision).AsDelete()
		conn.Where("id = ?", oldVision.ID).
			First(&oldDBVision).
			Delete(&oldDBVision)
	}
}
