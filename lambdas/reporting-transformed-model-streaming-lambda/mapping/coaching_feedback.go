package mapping

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	logger2 "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/jinzhu/gorm"
)

type DbCoachingFeedback struct {
	ID               string            `gorm:"primary_key"`
	CompetencyID     string            `gorm:"type:CHAR(36)"`
	Source           string            `gorm:"type:CHAR(9)"`
	Target           string            `gorm:"type:CHAR(9)"`
	Channel          string            `gorm:"type:CHAR(9)"`
	Quarter          int               `gorm:"type:SMALLINT"`
	Year             int               `gorm:"type:SMALLINT"`
	ConfidenceFactor int               `gorm:"type:SMALLINT"`
	Feedback         string            `gorm:"type:TEXT"`
	PlatformID       common.PlatformID `gorm:"type:CHAR(9)"`
	MsgTimestamp     string            `gorm:"type:TEXT"`
	model.DBModel
}

func coachingFeedbackDBMapping(feedback models.UserFeedback) DbCoachingFeedback {
	qySplits := strings.Split(feedback.QuarterYear, ":")
	q, y := stringToInt(qySplits[0]), stringToInt(qySplits[1])
	return DbCoachingFeedback{
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

func (d DbCoachingFeedback) AsAdd() (op DbCoachingFeedback) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DbCoachingFeedback) AsUpdate() (op DbCoachingFeedback) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DbCoachingFeedback) AsDelete() (op DbCoachingFeedback) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func InterfaceToCoachingFeedbackUnsafe(ip interface{}, logger logger2.AdaptiveLogger) interface{} {
	var feedback models.UserFeedback
	js, _ := json.Marshal(ip)
	err := json.Unmarshal(js, &feedback)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to models.UserFeedback")
	}
	return feedback
}

func CoachingFeedbackStreamEntityHandler(e2 model.StreamEntity, conn *gorm.DB, logger logger2.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for user feedback")
	conn.AutoMigrate(&DbCoachingFeedback{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newCoachingFeedback = e2.NewEntity.(models.UserFeedback)
		dbCoachingFeedback := coachingFeedbackDBMapping(newCoachingFeedback).AsAdd()
		conn.Where("id = ?", dbCoachingFeedback.ID).
			Assign(dbCoachingFeedback).
			FirstOrCreate(&dbCoachingFeedback)
	case model.StreamEventEdit:
		var oldCoachingFeedback = e2.OldEntity.(models.UserFeedback)
		var newCoachingFeedback = e2.NewEntity.(models.UserFeedback)

		dbCoachingFeedback := coachingFeedbackDBMapping(newCoachingFeedback).AsUpdate()
		conn.Where("id = ?", oldCoachingFeedback.ID).
			Assign(dbCoachingFeedback).
			FirstOrCreate(&dbCoachingFeedback)
	case model.StreamEventDelete:
		var oldCoachingFeedback = e2.OldEntity.(models.UserFeedback)
		var oldDbCoachingFeedback = coachingFeedbackDBMapping(oldCoachingFeedback).AsDelete()
		conn.Where("id = ?", oldCoachingFeedback.ID).
			First(&oldDbCoachingFeedback).
			Delete(&oldDbCoachingFeedback)
	}
}
