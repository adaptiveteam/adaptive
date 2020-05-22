package mapping

import (
	"encoding/json"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	logger2 "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/adaptiveteam/adaptive/daos/common"
)

type DBCompetency struct {
	ID             string `gorm:"primary_key"`
	Name           string `gorm:"type:TEXT"`
	Description    string `gorm:"type:TEXT"`
	PlatformID     common.PlatformID `gorm:"type:TEXT"`
	CompetencyType string `gorm:"type:TEXT"`
	DeactivatedOn  string `json:"deactivated_on"`
	model.DBModel
}

// TableName return table name
func (d DBCompetency) TableName() string {
	return "competency"
}

func competencyDBMapping(comp models.AdaptiveValue) DBCompetency {
	return DBCompetency{
		ID:             comp.ID,
		Name:           comp.Name,
		Description:    comp.Description,
		PlatformID:     comp.PlatformID,
		CompetencyType: comp.ValueType,
		DeactivatedOn:  comp.DeactivatedAt,
	}
}

func (d DBCompetency) AsAdd() (op DBCompetency) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBCompetency) AsUpdate() (op DBCompetency) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBCompetency) AsDelete() (op DBCompetency) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func (d DBCompetency) ParseUnsafe(js []byte, logger logger2.AdaptiveLogger) interface{} {
	var competency models.AdaptiveValue
	err := json.Unmarshal(js, &competency)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.AdaptiveValue")
	}
	return competency
}

func (d DBCompetency) HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger2.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for competency")
	conn.AutoMigrate(&DBCompetency{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newCompetency = e2.NewEntity.(models.AdaptiveValue)
		dbCompetency := competencyDBMapping(newCompetency).AsAdd()
		conn.Where("id = ?", dbCompetency.ID).
			Assign(dbCompetency).
			FirstOrCreate(&dbCompetency)
	case model.StreamEventEdit:
		var oldCompetency = e2.OldEntity.(models.AdaptiveValue)
		var newCompetency = e2.NewEntity.(models.AdaptiveValue)

		dbCompetency := competencyDBMapping(newCompetency).AsUpdate()
		conn.Where("id = ?", oldCompetency.ID).
			Assign(dbCompetency).
			FirstOrCreate(&dbCompetency)
	case model.StreamEventDelete:
		var oldCompetency = e2.OldEntity.(models.AdaptiveValue)
		var oldDBCompetency = competencyDBMapping(oldCompetency).AsDelete()
		conn.Where("id = ?", oldCompetency.ID).
			First(&oldDBCompetency).
			Delete(&oldDBCompetency)
	}
}
