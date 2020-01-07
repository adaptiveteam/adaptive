package strategy

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"time"
)

const (
	CreateNowActionLabel   = "Yes, right now"
)

func CreateStrategyEntityEng(table string, mc models.MessageCallback, title, fallback string, urgent bool,
	check models.UserEngagementCheckWithValue, platformID models.PlatformID) {
	actions := []ebm.AttachmentAction{
		*models.DialogAttachAction(mc, models.Now, CreateNowActionLabel),
	}
	utils.AddChatEngagement(mc, title, core.EmptyString, fallback, mc.Source, actions, []ebm.AttachmentField{},
		platformID, urgent, table, common.DeprecatedGetGlobalDns().Dynamo, common.DeprecatedGetGlobalDns().Namespace, time.Now().Unix(), check)
}
