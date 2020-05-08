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

type DBObjective struct {
	ID           string `gorm:"primary_key"`
	Advocate     string `gorm:"type:TEXT"`
	AsMeasuredBy string `gorm:"type:TEXT"`
	// TODO: Look into JSON field type here
	// https://github.com/jinzhu/gorm/issues/1935
	CapabilityCommunityIDs string `gorm:"type:TEXT"`
	CreatedAt              string `gorm:"type:TEXT"`
	CreatedBy              string `gorm:"type:TEXT"`
	Description            string `gorm:"type:TEXT"`
	ExpectedEndDate        string `gorm:"type:DATE"`
	Name                   string `gorm:"type:TEXT"`
	PlatformID             common.PlatformID `gorm:"type:TEXT"`
	Targets                string `gorm:"type:TEXT"`
	Type                   string `gorm:"type:TEXT"`
	model.DBModel
}

// TableName return table name
func (d DBObjective) TableName() string {
	return "objective"
}

func objectiveDBMapping(so models.StrategyObjective) DBObjective {
	capComms, _ := json.Marshal(so.CapabilityCommunityIDs)
	return DBObjective{
		ID:                     so.ID,
		Advocate:               so.Advocate,
		AsMeasuredBy:           so.AsMeasuredBy,
		CapabilityCommunityIDs: string(capComms),
		CreatedAt:              so.CreatedAt,
		CreatedBy:              so.CreatedBy,
		Description:            so.Description,
		ExpectedEndDate:        so.ExpectedEndDate,
		Name:                   so.Name,
		PlatformID:             so.PlatformID,
		Targets:                so.Targets,
		Type:                   string(so.ObjectiveType),
	}
}

func (d DBObjective) AsAdd() (op DBObjective) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBObjective) AsUpdate() (op DBObjective) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBObjective) AsDelete() (op DBObjective) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func (d DBObjective) ParseUnsafe(js []byte, logger logger.AdaptiveLogger) interface{} {
	var so models.StrategyObjective
	err := json.Unmarshal(js, &so)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.StrategyObjective")
	}
	return so
}

func (d DBObjective) HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for objective")
	conn.AutoMigrate(&DBObjective{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newObjective = e2.NewEntity.(models.StrategyObjective)
		dbObjective := objectiveDBMapping(newObjective).AsAdd()
		conn.Where("id = ?", dbObjective.ID).
			Assign(dbObjective).
			FirstOrCreate(&dbObjective)
	case model.StreamEventEdit:
		var oldObjective = e2.OldEntity.(models.StrategyObjective)
		var newObjective = e2.NewEntity.(models.StrategyObjective)

		dbObjective := objectiveDBMapping(newObjective).AsUpdate()
		conn.Where("id = ?", oldObjective.ID).
			Assign(dbObjective).
			FirstOrCreate(&dbObjective)
	case model.StreamEventDelete:
		var oldObjective = e2.OldEntity.(models.StrategyObjective)
		var oldDbObjective = objectiveDBMapping(oldObjective).AsDelete()
		conn.Where("id = ?", oldObjective.ID).
			First(&oldDbObjective).
			Delete(&oldDbObjective)
	}
}
