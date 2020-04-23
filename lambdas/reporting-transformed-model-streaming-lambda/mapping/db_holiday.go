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

type DBHoliday struct {
	ID               string `gorm:"primary_key"`
	Date             string `gorm:"type:DATE"`
	Description      string `gorm:"type:TEXT"`
	Name             string `gorm:"type:TEXT"`
	PlatformID       common.PlatformID `gorm:"type:CHAR(9)"`
	ScopeCommunities string `gorm:"type:TEXT"`
	model.DBModel
}

// TableName return table name
func (d DBHoliday) TableName() string {
	return "holiday"
}

func holidayDBMapping(h models.AdHocHoliday) DBHoliday {
	return DBHoliday{
		ID:               h.ID,
		Date:             h.Date,
		Name:             h.Name,
		PlatformID:       h.PlatformID,
		ScopeCommunities: h.ScopeCommunities,
	}
}

func (d DBHoliday) AsAdd() (op DBHoliday) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBHoliday) AsUpdate() (op DBHoliday) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBHoliday) AsDelete() (op DBHoliday) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func (d DBHoliday) ParseUnsafe(js []byte, logger logger.AdaptiveLogger) interface{} {
	var holiday models.AdHocHoliday
	err := json.Unmarshal(js, &holiday)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.AdHocHoliday")
	}
	return holiday
}

func (d DBHoliday) HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for holiday")
	conn.AutoMigrate(&DBHoliday{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newHoliday = e2.NewEntity.(models.AdHocHoliday)
		dbHoliday := holidayDBMapping(newHoliday).AsAdd()
		conn.Where("id = ?", dbHoliday.ID).
			Assign(dbHoliday).
			FirstOrCreate(&dbHoliday)
	case model.StreamEventEdit:
		var oldHoliday = e2.OldEntity.(models.AdHocHoliday)
		var newHoliday = e2.NewEntity.(models.AdHocHoliday)

		dbHoliday := holidayDBMapping(newHoliday).AsUpdate()
		conn.Where("id = ?", oldHoliday.ID).
			Assign(dbHoliday).
			FirstOrCreate(&dbHoliday)
	case model.StreamEventDelete:
		var oldHoliday = e2.OldEntity.(models.AdHocHoliday)
		var oldDbHoliday = holidayDBMapping(oldHoliday).AsDelete()
		conn.Where("id = ?", oldHoliday.ID).
			First(&oldDbHoliday).
			Delete(&oldDbHoliday)
	}
}
