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

type DBUserObjective struct {
	ID                          string                          `gorm:"primary_key"`
	UserID                      string                          `gorm:"type:CHAR(9)"`
	Name                        string                          `gorm:"type:TEXT"`
	Description                 string                          `gorm:"type:TEXT"`
	AccountabilityPartner       string                          `gorm:"type:CHAR(9)"`
	Accepted                    bool                            `gorm:"type:BOOLEAN"`
	Type                        models.DevelopmentObjectiveType `gorm:"type:TEXT"`
	StrategyAlignmentEntityID   string                          `gorm:"type:CHAR(36)"`
	StrategyAlignmentEntityType models.AlignedStrategyType      `gorm:"type:TEXT"`
	Quarter                     uint
	Year                        uint
	CreatedDate                 string            `gorm:"type:DATE"`
	ExpectedEndDate             string            `gorm:"type:DATE"`
	Completed                   bool              `gorm:"type:BOOLEAN"`
	PartnerVerifiedCompletion   bool              `gorm:"type:BOOLEAN"`
	Comments                    string            `gorm:"type:TEXT"`
	Cancelled                   bool              `gorm:"type:BOOLEAN"`
	PlatformID                  common.PlatformID `gorm:"type:CHAR(9)"`
	model.DBModel
}

func userObjectiveDBMapping(obj models.UserObjective) DBUserObjective {
	return DBUserObjective{
		ID:                          obj.ID,
		UserID:                      obj.UserID,
		Name:                        obj.Name,
		Description:                 obj.Description,
		AccountabilityPartner:       obj.AccountabilityPartner,
		Accepted:                    intToBoolean(obj.Accepted),
		Type:                        obj.ObjectiveType,
		StrategyAlignmentEntityID:   obj.StrategyAlignmentEntityID,
		StrategyAlignmentEntityType: obj.StrategyAlignmentEntityType,
		Quarter:                     uint(obj.Quarter),
		Year:                        uint(obj.Year),
		CreatedDate:                 obj.CreatedDate,
		ExpectedEndDate:             obj.ExpectedEndDate,
		Completed:                   intToBoolean(obj.Completed),
		PartnerVerifiedCompletion:   obj.PartnerVerifiedCompletion,
		Comments:                    obj.Comments,
		Cancelled:                   intToBoolean(obj.Cancelled),
		PlatformID:                  obj.PlatformID,
	}
}

func (d DBUserObjective) AsAdd() (op DBUserObjective) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBUserObjective) AsUpdate() (op DBUserObjective) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBUserObjective) AsDelete() (op DBUserObjective) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func InterfaceToUserObjectiveUnsafe(ip interface{}, logger logger.AdaptiveLogger) interface{} {
	var uObj models.UserObjective
	js, _ := json.Marshal(ip)
	err := json.Unmarshal(js, &uObj)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.UserObjective")
	}
	return uObj
}

func UserObjectiveStreamEntityHandler(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for UserObjective")
	conn.AutoMigrate(&DBUserObjective{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newUObj = e2.NewEntity.(models.UserObjective)
		dbUObj := userObjectiveDBMapping(newUObj).AsAdd()
		conn.Where("id = ?", dbUObj.ID).
			Assign(dbUObj).
			FirstOrCreate(&dbUObj)
	case model.StreamEventEdit:
		var oldUObj = e2.OldEntity.(models.UserObjective)
		var newUObj = e2.NewEntity.(models.UserObjective)

		dbUObj := userObjectiveDBMapping(newUObj).AsUpdate()
		conn.Where("id = ?", oldUObj.ID).
			Assign(dbUObj).
			FirstOrCreate(&dbUObj)
	case model.StreamEventDelete:
		var oldUObj = e2.OldEntity.(models.UserObjective)
		var oldDbUObj = userObjectiveDBMapping(oldUObj).AsDelete()
		conn.Where("id = ?", oldUObj.ID).
			First(&oldDbUObj).
			Delete(&oldDbUObj)
	}
}
