package mapping

import (
	"encoding/json"
	"time"

	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/jinzhu/gorm"
)

type DBClientConfig struct {
	PlatformID   common.PlatformID `gorm:"primary_key"`
	PlatformName string            `gorm:"type:CHAR(9)"`
	PlatformOrg  string            `gorm:"type:TEXT"`
	model.DBModel
}

func clientConfigDBMapping(c models.ClientPlatformToken) DBClientConfig {
	return DBClientConfig{
		PlatformID:   c.PlatformID,
		PlatformName: string(c.PlatformName),
		PlatformOrg:  c.Org,
	}
}

func (d DBClientConfig) AsAdd() (op DBClientConfig) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBClientConfig) AsUpdate() (op DBClientConfig) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBClientConfig) AsDelete() (op DBClientConfig) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func InterfaceToClientConfigUnsafe(ip interface{}, logger logger.AdaptiveLogger) interface{} {
	var cpt models.ClientPlatformToken
	js, _ := json.Marshal(ip)
	err := json.Unmarshal(js, &cpt)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.ClientPlatformToken")
	}
	return cpt
}

func ClientConfigStreamEntityHandler(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for client config")
	conn.AutoMigrate(&DBClientConfig{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newCpt = e2.NewEntity.(models.ClientPlatformToken)
		dbCpt := clientConfigDBMapping(newCpt).AsAdd()
		conn.Where("platform_id = ?", newCpt.PlatformID).
			Assign(dbCpt).
			FirstOrCreate(&dbCpt)
	case model.StreamEventEdit:
		var oldCpt = e2.OldEntity.(models.ClientPlatformToken)
		var newCpt = e2.NewEntity.(models.ClientPlatformToken)

		dbCpt := clientConfigDBMapping(newCpt).AsUpdate()
		conn.Where("platform_id = ?", oldCpt.PlatformID).
			Assign(dbCpt).
			FirstOrCreate(&dbCpt)
	case model.StreamEventDelete:
		var oldCpt = e2.OldEntity.(models.ClientPlatformToken)
		var oldDbCpt = clientConfigDBMapping(oldCpt).AsDelete()
		conn.Where("platform_id = ?", oldCpt.PlatformID).
			First(&oldDbCpt).
			Delete(&oldDbCpt)
	}
}
