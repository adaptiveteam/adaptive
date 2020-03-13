package mapping

import (
	"encoding/json"
	"time"

	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/jinzhu/gorm"
	"github.com/adaptiveteam/adaptive/daos/common"

)

type DBObjectiveCommunity struct {
	ID          string   `gorm:"primary_key"`
	Advocate    string   `gorm:"type:CHAR(9)"`
	CreatedBy   string   `gorm:"type:CHAR(9)"`
	Description string   `gorm:"type:TEXT"`
	Name        string   `gorm:"type:TEXT"`
	PlatformID  common.PlatformID `gorm:"type:CHAR(9)"`
	CreatedAt   string   `gorm:"type:TEXT"`
	model.DBModel
}

func objectiveCommunityDBMapping(cc models.CapabilityCommunity) DBObjectiveCommunity {
	return DBObjectiveCommunity{
		ID:          cc.ID,
		Advocate:    cc.Advocate,
		CreatedBy:   cc.CreatedBy,
		Description: cc.Description,
		Name:        cc.Name,
		PlatformID:  cc.PlatformID,
		CreatedAt:   cc.CreatedAt,
	}
}

func (d DBObjectiveCommunity) AsAdd() (op DBObjectiveCommunity) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBObjectiveCommunity) AsUpdate() (op DBObjectiveCommunity) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBObjectiveCommunity) AsDelete() (op DBObjectiveCommunity) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func InterfaceToObjectiveCommunityUnsafe(ip interface{}, logger logger.AdaptiveLogger) interface{} {
	var cc models.CapabilityCommunity
	js, _ := json.Marshal(ip)
	err := json.Unmarshal(js, &cc)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.CapabilityCommunity")
	}
	return cc
}

func ObjectiveCommunityStreamEntityHandler(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for objective community")
	conn.AutoMigrate(&DBObjectiveCommunity{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newObjectiveComm = e2.NewEntity.(models.CapabilityCommunity)
		dbObjectiveComm := objectiveCommunityDBMapping(newObjectiveComm).AsAdd()
		conn.Where("id = ?", dbObjectiveComm.ID).
			Assign(dbObjectiveComm).
			FirstOrCreate(&dbObjectiveComm)
	case model.StreamEventEdit:
		var oldObjectiveComm = e2.OldEntity.(models.CapabilityCommunity)
		var newObjectiveComm = e2.NewEntity.(models.CapabilityCommunity)

		dbObjectiveComm := objectiveCommunityDBMapping(newObjectiveComm).AsUpdate()
		conn.Where("id = ?", oldObjectiveComm.ID).
			Assign(dbObjectiveComm).
			FirstOrCreate(&dbObjectiveComm)
	case model.StreamEventDelete:
		var oldObjectiveComm = e2.OldEntity.(models.CapabilityCommunity)
		var oldDbObjectiveComm = objectiveCommunityDBMapping(oldObjectiveComm).AsDelete()
		conn.Where("id = ?", oldObjectiveComm.ID).
			First(&oldDbObjectiveComm).
			Delete(&oldDbObjectiveComm)
	}
}