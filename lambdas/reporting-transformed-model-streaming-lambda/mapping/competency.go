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

type DbCompetency struct {
	ID             string `gorm:"primary_key"`
	Name           string `gorm:"type:TEXT"`
	Description    string `gorm:"type:TEXT"`
	PlatformID     common.PlatformID `gorm:"type:CHAR(9)"`
	CompetencyType string `gorm:"type:TEXT"`
	DeactivatedOn  string `json:"deactivated_on"`
	model.DBModel
}

func competencyDBMapping(comp models.AdaptiveValue) DbCompetency {
	return DbCompetency{
		ID:             comp.ID,
		Name:           comp.Name,
		Description:    comp.Description,
		PlatformID:     comp.PlatformID,
		CompetencyType: comp.ValueType,
		DeactivatedOn:  comp.DeactivatedAt,
	}
}

func (d DbCompetency) AsAdd() (op DbCompetency) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DbCompetency) AsUpdate() (op DbCompetency) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DbCompetency) AsDelete() (op DbCompetency) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func InterfaceToCompetencyUnsafe(ip interface{}, logger logger2.AdaptiveLogger) interface{} {
	var competency models.AdaptiveValue
	js, _ := json.Marshal(ip)
	err := json.Unmarshal(js, &competency)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.AdaptiveValue")
	}
	return competency
}

func CompetencyStreamEntityHandler(e2 model.StreamEntity, conn *gorm.DB, logger logger2.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for competency")
	conn.AutoMigrate(&DbCompetency{})

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
		var oldDbCompetency = competencyDBMapping(oldCompetency).AsDelete()
		conn.Where("id = ?", oldCompetency.ID).
			First(&oldDbCompetency).
			Delete(&oldDbCompetency)
	}
}
