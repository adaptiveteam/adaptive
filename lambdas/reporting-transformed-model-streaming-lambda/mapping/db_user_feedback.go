package mapping

import (
	"encoding/json"
	"strings"
	"time"

	logger2 "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	"github.com/jinzhu/gorm"
)

type DBUserFeedback struct {
	ID               string            `gorm:"primary_key"`
	CompetencyID     string            `gorm:"type:CHAR(36)"`
	Source           string            `gorm:"type:TEXT"`
	Target           string            `gorm:"type:TEXT"`
	Channel          string            `gorm:"type:TEXT"`
	Quarter          int               `gorm:"type:SMALLINT"`
	Year             int               `gorm:"type:SMALLINT"`
	ConfidenceFactor int               `gorm:"type:SMALLINT"`
	Feedback         string            `gorm:"type:TEXT"`
	PlatformID       common.PlatformID `gorm:"type:TEXT"`
	MsgTimestamp     string            `gorm:"type:TEXT"`
	model.DBModel
}

// TableName return table name
func (d DBUserFeedback) TableName() string {
	return "user_feedback"
}

func UserFeedbackDBMapping(feedback models.UserFeedback) DBUserFeedback {
	qySplits := strings.Split(feedback.QuarterYear, ":")
	q, y := stringToInt(qySplits[0]), stringToInt(qySplits[1])
	return DBUserFeedback{
		ID:               feedback.ID,
		CompetencyID:     feedback.ValueID,
		Source:           feedback.Source,
		Target:           feedback.Target,
		Quarter:          q,
		Year:             y,
		Channel:          feedback.ChannelID,
		MsgTimestamp:     feedback.MsgTimestamp,
		ConfidenceFactor: stringToInt(feedback.ConfidenceFactor),
		Feedback:         feedback.Feedback,
		PlatformID:       feedback.PlatformID,
	}
}

func (d DBUserFeedback) AsAdd() (op DBUserFeedback) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBUserFeedback) AsUpdate() (op DBUserFeedback) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBUserFeedback) AsDelete() (op DBUserFeedback) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func (d DBUserFeedback) ParseUnsafe(js []byte, logger logger2.AdaptiveLogger) interface{} {
	var feedback models.UserFeedback
	err := json.Unmarshal(js, &feedback)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.UserFeedback")
	}
	return feedback
}

func (d DBUserFeedback) HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger2.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for user feedback")
	conn.AutoMigrate(&DBUserFeedback{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newUserFeedback = e2.NewEntity.(models.UserFeedback)
		DBUserFeedback := UserFeedbackDBMapping(newUserFeedback).AsAdd()
		conn.Where("id = ?", DBUserFeedback.ID).
			Assign(DBUserFeedback).
			FirstOrCreate(&DBUserFeedback)
	case model.StreamEventEdit:
		var oldUserFeedback = e2.OldEntity.(models.UserFeedback)
		var newUserFeedback = e2.NewEntity.(models.UserFeedback)

		DBUserFeedback := UserFeedbackDBMapping(newUserFeedback).AsUpdate()
		conn.Where("id = ?", oldUserFeedback.ID).
			Assign(DBUserFeedback).
			FirstOrCreate(&DBUserFeedback)
	case model.StreamEventDelete:
		var oldUserFeedback = e2.OldEntity.(models.UserFeedback)
		var oldDBUserFeedback = UserFeedbackDBMapping(oldUserFeedback).AsDelete()
		conn.Where("id = ?", oldUserFeedback.ID).
			First(&oldDBUserFeedback).
			Delete(&oldDBUserFeedback)
	}
}
