package objectives

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"time"
)

const (
	Comments = "Comments"

	ObjectiveCommentsNowActionLabel   = "Yes, right now"
	ObjectiveCommentsSkipActionLabel  = "Skip this, please"
)

func commentsEngAttachmentActions(mc models.MessageCallback, uObj *models.UserObjective) []ebm.AttachmentAction {
	return []ebm.AttachmentAction{
		*models.SimpleAttachAction(*mc.WithTarget(uObj.ID).WithTopic(Comments),
			models.Now,
			ObjectiveCommentsNowActionLabel),
		*models.SimpleAttachAction(mc, models.Ignore, ObjectiveCommentsSkipActionLabel)}
}

func CommentsEng(table string, mc models.MessageCallback, title, fallback string, obj *models.UserObjective, urgent bool,
	dns common.DynamoNamespace, check models.UserEngagementCheckWithValue) {
	actions := commentsEngAttachmentActions(mc, obj)
	utils.AddChatEngagement(mc, title, core.EmptyString, fallback, mc.Source, actions, []ebm.AttachmentField{},
		obj.PlatformID, urgent, table, dns.Dynamo, dns.Namespace, time.Now().Unix(), check)
}
