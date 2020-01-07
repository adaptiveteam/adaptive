package user

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"time"
)

const (
	NowActionLabel  = "Yes, right now"
	SkipActionLabel = "Skip this, please"
)

func userConfirmAttachmentActions(mc models.MessageCallback, ignoreAction bool) []ebm.AttachmentAction {
	actions := []ebm.AttachmentAction{
		*models.SimpleAttachAction(mc, models.Now, NowActionLabel),
	}
	if ignoreAction {
		actions = append(actions,
			*models.SimpleAttachAction(mc, models.Ignore, SkipActionLabel))
	}
	return actions
}

func UserConfirmEng(table string, mc models.MessageCallback, title, fallback string, urgent bool, dns common.DynamoNamespace,
	ignoreAction bool, check models.UserEngagementCheckWithValue, platformID models.PlatformID) {
	utils.AddChatEngagement(mc, title, core.EmptyString, fallback, mc.Source, userConfirmAttachmentActions(mc, ignoreAction),
		[]ebm.AttachmentField{}, platformID, urgent, table, dns.Dynamo, dns.Namespace, time.Now().Unix(), check)
}
