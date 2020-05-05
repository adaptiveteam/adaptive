package mapping

import (
	"encoding/json"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/jinzhu/gorm"
	"time"
)

type DBInitiativeCommunity struct {
	ID                    string `gorm:"primary_key"`
	Advocate              string `gorm:"type:TEXT"`
	CapabilityCommunityID string `gorm:"type:TEXT"`
	CreatedBy             string `gorm:"type:TEXT"`
	Description           string `gorm:"type:TEXT"`
	Name                  string `gorm:"type:TEXT"`
	PlatformID            common.PlatformID `gorm:"type:TEXT"`
	CreatedAt             string `gorm:"type:TEXT"`
	model.DBModel
}

// TableName return table name
func (d DBInitiativeCommunity) TableName() string {
	return "initiative_community"
}

func initiativeCommDBMapping(sic models.StrategyInitiativeCommunity) DBInitiativeCommunity {
	return DBInitiativeCommunity{
		ID:                    sic.ID,
		Advocate:              sic.Advocate,
		CapabilityCommunityID: sic.CapabilityCommunityID,
		CreatedBy:             sic.CreatedBy,
		Description:           sic.Description,
		Name:                  sic.Name,
		PlatformID:            sic.PlatformID,
		CreatedAt:             sic.CreatedAt,
	}
}

func (d DBInitiativeCommunity) AsAdd() (op DBInitiativeCommunity) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBInitiativeCommunity) AsUpdate() (op DBInitiativeCommunity) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBInitiativeCommunity) AsDelete() (op DBInitiativeCommunity) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func (d DBInitiativeCommunity) ParseUnsafe(js []byte, logger logger.AdaptiveLogger) interface{} {
	var sic models.StrategyInitiativeCommunity
	err := json.Unmarshal(js, &sic)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.StrategyInitiativeCommunity")
	}
	return sic
}

func (d DBInitiativeCommunity) HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for initiative community")
	conn.AutoMigrate(&DBInitiativeCommunity{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newInitComm = e2.NewEntity.(models.StrategyInitiativeCommunity)
		dbInitComm := initiativeCommDBMapping(newInitComm).AsAdd()
		conn.Where("id = ?", dbInitComm.ID).
			Assign(dbInitComm).
			FirstOrCreate(&dbInitComm)
	case model.StreamEventEdit:
		var oldInitComm = e2.OldEntity.(models.StrategyInitiativeCommunity)
		var newInitComm = e2.NewEntity.(models.StrategyInitiativeCommunity)

		dbInitComm := initiativeCommDBMapping(newInitComm).AsUpdate()
		conn.Where("id = ?", oldInitComm.ID).
			Assign(dbInitComm).
			FirstOrCreate(&dbInitComm)
	case model.StreamEventDelete:
		var oldInitComm = e2.OldEntity.(models.StrategyInitiativeCommunity)
		var oldDbInitComm = initiativeCommDBMapping(oldInitComm).AsDelete()
		conn.Where("id = ?", oldInitComm.ID).
			First(&oldDbInitComm).
			Delete(&oldDbInitComm)
	}
}
