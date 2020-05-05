package mapping

import (
	"encoding/json"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/jinzhu/gorm"
	"time"
)

type DBPartnershipRejection struct {
	ID                      string `gorm:"primary_key"`
	ObjectiveID             string `gorm:"type:TEXT"`
	AccountabilityPartnerID string `gorm:"type:TEXT"`
	Comments                string `gorm:"type:TEXT"`
	UserID                  string `gorm:"type:TEXT"`
	CreatedDate             string `gorm:"type:TEXT"`
	model.DBModel
}

// TableName return table name
func (d DBPartnershipRejection) TableName() string {
	return "partnership_rejection"
}

func partnershipRejectionCompositeKey(apr models.AccountabilityPartnerShipRejection) string {
	return apr.ObjectiveID + ":" + apr.CreatedOn
}

func partnershipRejectionDBMapping(apr models.AccountabilityPartnerShipRejection) DBPartnershipRejection {
	return DBPartnershipRejection{
		ID:                      partnershipRejectionCompositeKey(apr),
		ObjectiveID:             apr.ObjectiveID,
		AccountabilityPartnerID: apr.AccountabilityPartnerID,
		Comments:                apr.Comments,
		UserID:                  apr.UserID,
		CreatedDate:             apr.CreatedOn,
	}
}

func (d DBPartnershipRejection) AsAdd() (op DBPartnershipRejection) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBPartnershipRejection) AsUpdate() (op DBPartnershipRejection) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBPartnershipRejection) AsDelete() (op DBPartnershipRejection) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func (d DBPartnershipRejection) ParseUnsafe(js []byte, logger logger.AdaptiveLogger) interface{} {
	var apr models.AccountabilityPartnerShipRejection
	err := json.Unmarshal(js, &apr)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.AccountabilityPartnerShipRejection")
	}
	return apr
}

func (d DBPartnershipRejection) HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for partnership rejection")
	conn.AutoMigrate(&DBPartnershipRejection{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newPR = e2.NewEntity.(models.AccountabilityPartnerShipRejection)
		dbPR := partnershipRejectionDBMapping(newPR).AsAdd()
		conn.Where("id = ?", dbPR.ID).
			Assign(dbPR).
			FirstOrCreate(&dbPR)
	case model.StreamEventEdit:
		var oldPR = e2.OldEntity.(models.AccountabilityPartnerShipRejection)
		var newPR = e2.NewEntity.(models.AccountabilityPartnerShipRejection)

		dbPR := partnershipRejectionDBMapping(newPR).AsUpdate()
		conn.Where("id = ?", partnershipRejectionCompositeKey(oldPR)).
			Assign(dbPR).
			FirstOrCreate(&dbPR)
	case model.StreamEventDelete:
		var oldPR = e2.OldEntity.(models.AccountabilityPartnerShipRejection)
		var oldDbPR = partnershipRejectionDBMapping(oldPR).AsDelete()
		conn.Where("id = ?", partnershipRejectionCompositeKey(oldPR)).
			First(&oldDbPR).
			Delete(&oldDbPR)
	}
}
