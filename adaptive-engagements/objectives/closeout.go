package objectives

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"time"
)

const (
	ObjectiveCloseout                       ui.PlainText = "Objective Closeout"
	ObjectiveCloseoutNowActionLabel         ui.PlainText = "I agree"
	ObjectiveCloseoutDisagreeActionLabel    ui.PlainText = "I tend to disagree"
	ObjectiveCloseoutWhyDisagreeSurveyLabel ui.PlainText = "Why are you disagreeing with closeout?"

	ObjectiveCloseoutComment = "objective_closeout_comment"
)

// This engagement is to closeout an objective
func ObjectiveCloseoutEng(table string, mc models.MessageCallback, coach, title, text, fallback, learnLink string, urgent bool,
	dns common.DynamoNamespace, check models.UserEngagementCheckWithValue, platformID models.PlatformID) {
	// Setting action to closeout
	utils.AddChatEngagement(*mc.WithAction(string(Closeout)), title, text, fallback, coach, closeoutAttachmentActions(mc, learnLink),
		[]ebm.AttachmentField{}, platformID, urgent, table, dns.Dynamo, dns.Namespace, time.Now().Unix(), check)
}

func closeoutAttachmentActions(mc models.MessageCallback, learnTrailPath string) []ebm.AttachmentAction {
	return models.AppendOptionalAction(
		[]ebm.AttachmentAction{
			*models.SimpleAttachAction(mc, models.Now, ObjectiveCloseoutNowActionLabel),
			*models.GenAttachAction(mc, No, string(ObjectiveCloseoutDisagreeActionLabel),
				models.EmptyActionConfirm(), true),
		},
		models.LearnMoreAction(learnTrailPath),
	)
}
