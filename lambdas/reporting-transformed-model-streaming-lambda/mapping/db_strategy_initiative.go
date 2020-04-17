package mapping

import (
	"encoding/json"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/adaptiveteam/adaptive/daos/common"

)

type DBStrategyInitiative struct {
	ID                    string  `gorm:"primary_key"`
	Advocate              string  `gorm:"type:CHAR(9)"`
	Budget                float64 `gorm:"type:DOUBLE"`
	CapabilityObjectiveID string  `gorm:"type:CHAR(36)"`
	CreatedBy             string  `gorm:"type:CHAR(9)"`
	DefinitionOfVictory   string  `gorm:"type:TEXT"`
	Description           string  `gorm:"type:TEXT"`
	ExpectedEndDate       string  `gorm:"type:DATE"`
	InitiativeCommunityID string  `gorm:"type:CHAR(36)"`
	Name                  string  `gorm:"type:TEXT"`
	PlatformID            common.PlatformID  `gorm:"type:CHAR(9)"`
	CreatedAt             string  `gorm:"type:TEXT"`
	model.DBModel
}

// TableName return table name
func (d DBStrategyInitiative) TableName() string {
	return "strategy_initiative"
}

func initiativeDBMapping(vis models.StrategyInitiative) DBStrategyInitiative {
	return DBStrategyInitiative{
		ID:                    vis.ID,
		Advocate:              vis.Advocate,
		Budget:                stringToFloat(vis.Budget),
		CapabilityObjectiveID: vis.CapabilityObjective,
		CreatedBy:             vis.CreatedBy,
		DefinitionOfVictory:   vis.DefinitionOfVictory,
		Description:           vis.Description,
		ExpectedEndDate:       vis.ExpectedEndDate,
		InitiativeCommunityID: vis.InitiativeCommunityID,
		Name:                  vis.Name,
		PlatformID:            vis.PlatformID,
		CreatedAt:             vis.CreatedAt,
	}
}

func (d DBStrategyInitiative) AsAdd() (op DBStrategyInitiative) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBStrategyInitiative) AsUpdate() (op DBStrategyInitiative) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBStrategyInitiative) AsDelete() (op DBStrategyInitiative) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func InterfaceToInitiativeUnsafe(ip interface{}, logger logger.AdaptiveLogger) interface{} {
	var init models.StrategyInitiative
	js, _ := json.Marshal(ip)
	err := json.Unmarshal(js, &init)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.StrategyInitiative")
	}
	return init
}

func InitiativeStreamEntityHandler(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for initiative")
	conn.AutoMigrate(&DBStrategyInitiative{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newInit = e2.NewEntity.(models.StrategyInitiative)
		dbInit := initiativeDBMapping(newInit).AsAdd()
		conn.Where("id = ?", dbInit.ID).
			Assign(dbInit).
			FirstOrCreate(&dbInit)
	case model.StreamEventEdit:
		var oldInit = e2.OldEntity.(models.StrategyInitiative)
		var newInit = e2.NewEntity.(models.StrategyInitiative)

		dbInit := initiativeDBMapping(newInit).AsUpdate()
		conn.Where("id = ?", oldInit.ID).
			Assign(dbInit).
			FirstOrCreate(&dbInit)
	case model.StreamEventDelete:
		var oldInit = e2.OldEntity.(models.StrategyInitiative)
		var oldDbInit = initiativeDBMapping(oldInit).AsDelete()
		conn.Where("id = ?", oldInit.ID).
			First(&oldDbInit).
			Delete(&oldDbInit)
	}
}
