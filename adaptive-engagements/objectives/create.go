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
	Objectives                      = "Development Objectives"
	ObjectiveCreateNowActionLabel   = "Yes, right now"
	ObjectiveCreateSkipActionLabel  = "Skip this"
)

// This is an engagement prompting a user to create an objective
// It consists of 4 action labels: create now, create later, skip, learn more
func CreateObjectiveEng(table string, mc models.MessageCallback, coaches, dates []models.KvPair,
	initsAndObjs []ebm.AttachmentActionElementOptionGroup, title,
	fallback string, urgent bool, dns common.DynamoNamespace, check models.UserEngagementCheckWithValue,
	teamID models.TeamID) {
	actions := []ebm.AttachmentAction{
		*models.DialogAttachAction(mc, models.Now, ObjectiveCreateNowActionLabel),
		*models.SimpleAttachAction(mc, models.Ignore, ObjectiveCreateSkipActionLabel),
	}
	utils.AddChatEngagement(mc, title, core.EmptyString, fallback, mc.Source, actions, []ebm.AttachmentField{},
		teamID, urgent, table, dns.Dynamo, dns.Namespace, time.Now().Unix(), check)
}
