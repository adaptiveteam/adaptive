package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/nlopes/slack"
	"time"
)

func generateReportMenuHandler(request slack.InteractionCallback, mc models.MessageCallback, plaformID models.TeamID) {
	userID := request.User.ID
	channelID := request.Channel.ID
	mc.Set("Topic", "reports")
	mc.Set("Action", GenerateReportHR)
	// Posting user confirmation engagement
	// User id here should be channel since we are posting into a channel
	actions := user.UserSelectAttachments(mc, []string{}, []string{}, plaformID, userDAO)
	attach := utils.ChatAttachment(string(UserForReportSelectionPrompt),
		"", "", mc.ToCallbackID(), actions,
		[]ebm.AttachmentField{}, time.Now().Unix())
	// Delete the original engagement
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Attachments: []model.Attachment{*attach},
		AsUser: true, Ts: request.MessageTs})
}

func fetchReportMenuHandler(request slack.InteractionCallback, mc models.MessageCallback, teamID models.TeamID) {
	userID := request.User.ID
	channelID := request.Channel.ID
	mc.Set("Topic", "reports")
	mc.Set("Action", FetchReportHR)
	// Posting user confirmation engagement
	// User id here should be channel since we are posting into a channel
	actions := user.UserSelectAttachments(mc, []string{}, []string{}, teamID, userDAO)
	attach := utils.ChatAttachment(string(UserForReportSelectionPrompt),
		"", "", mc.ToCallbackID(), actions, []ebm.AttachmentField{}, time.Now().Unix())
	// Delete the original engagement
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Attachments: []model.Attachment{*attach},
		AsUser: true, Ts: request.MessageTs})
}
