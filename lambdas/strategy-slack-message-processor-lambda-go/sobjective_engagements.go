package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"time"
)

const (
	CreateNow   = "Yes, right now"
	CreateLater = "Yes, but later"
	SkipCreate  = "Skip this"
)

// This is an engagement prompting a user to create an objective
// It consists of 4 action labels: create now, create later, skip, learn more
func CreateAskEngagement(table string, platformID models.PlatformID, mc models.MessageCallback,
	title, fallback, learnTrailPath string, urgent bool, dns common.DynamoNamespace) {
	actions := models.AppendOptionalAction(
		[]ebm.AttachmentAction{
			*models.SimpleAttachAction(mc, models.Now, CreateNow),
			*models.SimpleAttachAction(mc, models.Ignore, SkipCreate),
		},
		models.LearnMoreAction(learnTrailPath))
	utils.AddChatEngagement(mc, title, core.EmptyString, fallback, mc.Source, actions, []ebm.AttachmentField{},
		platformID, urgent, table, dns.Dynamo, dns.Namespace, time.Now().Unix(), models.UserEngagementCheckWithValue{})
}
