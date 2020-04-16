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

type DBUserObjectiveProgress struct {
	ID                      string                      `gorm:"primary_key"`
	ObjectiveID             string                      `gorm:"type:CHAR(36)"`
	UserID                  string                      `gorm:"type:CHAR(9)"`
	Comments                string                      `gorm:"type:TEXT"`
	Closeout                bool                        `gorm:"type:BOOLEAN"`
	PercentTimeLapsed       float64                     `gorm:"type:DOUBLE"`
	StatusColor             models.ObjectiveStatusColor `gorm:"type:VARCHAR(6)"`
	PartnerID               string                      `gorm:"type:CHAR(9)"`
	ReviewedByPartner       bool                        `gorm:"type:BOOLEAN"`
	PartnerComments         string                      `gorm:"type:TEXT"`
	PartnerReportedProgress string                      `gorm:"type:TEXT"`
	PlatformID              common.PlatformID           `gorm:"type:CHAR(9)"`
	CreatedDate             string                      `gorm:"type:TEXT"`
	model.DBModel
}

func compositeKey(obj models.UserObjectiveProgress) string {
	return obj.ID + ":" + obj.CreatedOn
}

func userObjectiveProgressDBMapping(obj models.UserObjectiveProgress) DBUserObjectiveProgress {
	return DBUserObjectiveProgress{
		ID:                      compositeKey(obj),
		ObjectiveID:             obj.ID,
		UserID:                  obj.UserID,
		Comments:                obj.Comments,
		Closeout:                intToBoolean(obj.Closeout),
		PercentTimeLapsed:       stringToFloat(obj.PercentTimeLapsed),
		StatusColor:             obj.StatusColor,
		PartnerID:               obj.PartnerID,
		ReviewedByPartner:       obj.ReviewedByPartner,
		PartnerComments:         obj.PartnerComments,
		PartnerReportedProgress: obj.PartnerReportedProgress,
		PlatformID:              obj.PlatformID,
		CreatedDate:             obj.CreatedOn,
	}
}

func (d DBUserObjectiveProgress) AsAdd() (op DBUserObjectiveProgress) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBUserObjectiveProgress) AsUpdate() (op DBUserObjectiveProgress) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBUserObjectiveProgress) AsDelete() (op DBUserObjectiveProgress) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func InterfaceToUserObjectiveProgressUnsafe(ip interface{}, logger logger.AdaptiveLogger) interface{} {
	var uObj models.UserObjectiveProgress
	js, _ := json.Marshal(ip)
	err := json.Unmarshal(js, &uObj)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.UserObjectiveProgress")
	}
	return uObj
}

func UserObjectiveProgressStreamEntityHandler(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for user objective progress")
	conn.AutoMigrate(&DBUserObjectiveProgress{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newUObj = e2.NewEntity.(models.UserObjectiveProgress)
		dbUObj := userObjectiveProgressDBMapping(newUObj).AsAdd()
		conn.Where("id = ?", compositeKey(newUObj)).
			Assign(dbUObj).
			FirstOrCreate(&dbUObj)
	case model.StreamEventEdit:
		var oldUObj = e2.OldEntity.(models.UserObjectiveProgress)
		var newUObj = e2.NewEntity.(models.UserObjectiveProgress)

		dbUObj := userObjectiveProgressDBMapping(newUObj).AsUpdate()
		conn.Where("id = ?", compositeKey(oldUObj)).
			Assign(dbUObj).
			FirstOrCreate(&dbUObj)
	case model.StreamEventDelete:
		var oldUObj = e2.OldEntity.(models.UserObjectiveProgress)
		var oldDbUObj = userObjectiveProgressDBMapping(oldUObj).AsDelete()
		conn.Where("id = ?", compositeKey(oldUObj)).
			First(&oldDbUObj).
			Delete(&oldDbUObj)
	}
}
