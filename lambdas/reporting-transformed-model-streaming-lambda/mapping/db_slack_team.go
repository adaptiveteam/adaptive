package mapping

import (
	"github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/slackTeam"
	"encoding/json"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/adaptiveteam/adaptive/daos/common"

)

type SlackTeam struct {
	ID               string `gorm:"primary_key"`
	Date             string `gorm:"type:DATE"`
	Description      string `gorm:"type:TEXT"`
	Name             string `gorm:"type:TEXT"`
	PlatformID       common.PlatformID `gorm:"type:CHAR(9)"`
	ScopeCommunities string `gorm:"type:TEXT"`

	TeamID common.PlatformID `gorm:"primary_key"`
	AccessToken string `gorm:"type:TEXT"`
	TeamName string `gorm:"type:TEXT"`
	// UserID is the ID of the user to send an engagement to
	// This usually corresponds to the platform user id
	UserID string `gorm:"type:TEXT"`
	EnterpriseID string `gorm:"type:TEXT"`
	BotUserID string `gorm:"type:TEXT"`
	// bot_access_token
	BotAccessToken string `gorm:"type:TEXT"`

	model.DBModel
}

// TableName return table name
func (d SlackTeam) TableName() string {
	return "slack_team"
}

func slackteamDBMapping(h slackTeam.SlackTeam) SlackTeam {
	createdAt, createdAtDefined := core_utils_go.ParseDateOrTimestamp(h.CreatedAt)
	if !createdAtDefined {
		createdAt = time.Now()
	}
	updatedAt, updatedAtDefined := core_utils_go.ParseDateOrTimestamp(h.ModifiedAt)
	if !updatedAtDefined {
		updatedAt = time.Now()
	}
	return SlackTeam{
		TeamID: h.TeamID,
		AccessToken: h.AccessToken,
		TeamName: h.TeamName,
		UserID: h.UserID,
		EnterpriseID: h.EnterpriseID,
		BotUserID: h.BotUserID,
		BotAccessToken: h.BotAccessToken,
		DBModel: model.DBModel{
			DBCreatedAt: createdAt,
			DBUpdatedAt: updatedAt,
		},
	}
}

func (d SlackTeam) AsDelete() (op SlackTeam) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}

func InterfaceToSlackTeamUnsafe(ip interface{}, logger logger.AdaptiveLogger) interface{} {
	var slackteam slackTeam.SlackTeam
	js, _ := json.Marshal(ip)
	err := json.Unmarshal(js, &slackteam)
	if err != nil {
		logger.WithField("error", err).Errorf("Could not unmarshal to slackTeam.SlackTeam")
	}
	return slackteam
}

func SlackTeamStreamEntityHandler(e2 model.StreamEntity, conn *gorm.DB, logger logger.AdaptiveLogger) {
	logger.WithField("mapped_event", &e2).Info("Transformed request for slackteam")
	conn.AutoMigrate(&SlackTeam{})

	switch e2.EventType {
	case model.StreamEventAdd:
		var newSlackTeam = e2.NewEntity.(slackTeam.SlackTeam)
		dbSlackTeam := slackteamDBMapping(newSlackTeam)
		conn.Where("team_id = ?", dbSlackTeam.TeamID).
			Assign(dbSlackTeam).
			FirstOrCreate(&dbSlackTeam)
	case model.StreamEventEdit:
		var oldSlackTeam = e2.OldEntity.(slackTeam.SlackTeam)
		var newSlackTeam = e2.NewEntity.(slackTeam.SlackTeam)

		dbSlackTeam := slackteamDBMapping(newSlackTeam).AsUpdate()
		conn.Where("team_id = ?", oldSlackTeam.TeamID).
			Assign(dbSlackTeam).
			FirstOrCreate(&dbSlackTeam)
	case model.StreamEventDelete:
		var oldSlackTeam = e2.OldEntity.(slackTeam.SlackTeam)
		var dbOldSlackTeam = slackteamDBMapping(oldSlackTeam).AsDelete()
		conn.Where("team_id = ?", dbOldSlackTeam.TeamID).
			First(&dbOldSlackTeam).
			Delete(&dbOldSlackTeam)
	}
}
