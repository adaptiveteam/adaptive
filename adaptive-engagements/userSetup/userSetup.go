// Package userSetup contains engagement code to deal with user-specific settings.
package userSetup

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/userEngagement"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strings"
)

var (
	defaultMeetingTime = business_time.MeetingTime(9, 0)
)

const (
	UserSettingsModule      = "user_settings"
	UserSettingsUpdateTopic = "settings_update"
)

func convertToMenuOption(t business_time.LocalTime) ebm.MenuOption {
	return ebm.Option(t.ID(), ui.PlainText(t.ToUserFriendly()))
}

func mapTimeToMenuOption(times []business_time.LocalTime, f func(business_time.LocalTime) ebm.MenuOption) (menuOptions []ebm.MenuOption) {
	menuOptions = make([]ebm.MenuOption, len(times))
	for i := 0; i < len(times); i++ {
		menuOptions[i] = convertToMenuOption(times[i])
	}
	return
}

func dbGetCurrentMeetingTime(userID, usersTable string) (currentTime business_time.LocalTime) {
	params := map[string]*dynamodb.AttributeValue{
		"id": {S: aws.String(userID)},
	}
	var out models.User
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTable(usersTable, params, &out)

	currentValue := out.AdaptiveScheduledTime
	if currentValue == "" {
		currentValue = defaultMeetingTime.ID()
	}
	currentTime, err = business_time.ParseLocalTimeID(currentValue)
	if err != nil {
		fmt.Printf("Couldn't parse current meeting time %s: %v; Using default value %v", currentValue, err, defaultMeetingTime)
		currentTime = defaultMeetingTime
	}
	return
}

func meetingTime(platform utils.Platform, userID, usersTable string) {
	mc := callback(userID, MeetingTimeUserAttributeID, string(models.Select))
	callbackID := mc.ToCallbackID()
	// Meeting time attachment
	meetingTimeRange := business_time.DefaultMeetingTimeRange()
	menuOptions := mapTimeToMenuOption(meetingTimeRange, convertToMenuOption)
	currentTime := dbGetCurrentMeetingTime(userID, usersTable)

	chooseMeetingTimeNowAction := models.SelectAttachAction(mc, models.Now, string(ChooseMeetingTime),
		menuOptions, []ebm.MenuOptionGroup{})
	cancelAttachAction := models.SimpleAttachAction(mc, models.Cancel, models.CancelLabel)

	attach, _ := eb.NewAttachmentBuilder().
		Title(string(QueryMeetingTime)).
		Text(string(MeetingIsScheduledForCalmNotice(currentTime))).
		Fallback(string(QueryMeetingTime2)).
		Color(models.BlueColorHex).
		AttachmentType(models.DefaultAttachmentType).
		CallbackId(callbackID).
		Actions(
			[]ebm.AttachmentAction{*chooseMeetingTimeNowAction, *cancelAttachAction},
		).
		Build()


	platform.Publish(models.PlatformSimpleNotification{
		UserId:      userID,
		Attachments: []ebm.Attachment{*attach},
	})
}

func settingsUpdateMessage(engDao userEngagement.DAO, prefix string, mc models.MessageCallback, userID, ts string, platform utils.Platform) {
	callbackID := mc.ToCallbackID()
	updateAction := models.SimpleAttachAction(mc, models.Update, ui.PlainText(strings.Title(string(models.Update))))
	cancelAction := models.SimpleAttachAction(mc, models.Cancel, ui.PlainText(strings.Title(string(models.Cancel))))

	attach, _ := eb.NewAttachmentBuilder().
		Title(prefix + " " + string(PromptToUpdateSettings)).
		Fallback(string(PromptToUpdateSettings2)).
		Color(models.BlueColorHex).
		AttachmentType(models.DefaultAttachmentType).
		CallbackId(callbackID).
		Actions([]ebm.AttachmentAction{*updateAction, *cancelAction}).
		Build()
	engagement := eb.NewEngagementBuilder().
		Id(callbackID).
		WithResponseType(models.SlackInChannel).
		WithAttachment(attach).
		Build()
	bytes, err := engagement.ToJson()
	core.ErrorHandler(err, "user-setup-engagement", "Could not convert engagement to JSON")

	engDao.CreateUnsafe(models.UserEngagement{UserID: userID, ID: callbackID, Script: string(bytes),
		Priority: models.UrgentPriority, Answered: 0, CreatedAt: core.CurrentRFCTimestamp()})
	platform.Publish(models.PlatformSimpleNotification{
		UserId: userID,
		Ts:     ts,
	})
}

func callback(userID, topic, action string) models.MessageCallback {
	// We are writing month rather than quarter in engagement because quarter can always be inferred from month
	// This is a global setting. We should not include quarter and year for this. Or, hardcode them.
	mc := models.MessageCallback{Module: UserSettingsModule, Source: userID, Topic: topic, Action: action}
	return mc
}

// HandleUserSetupRequest is a handler that can be used instead of userSetup lambda.
func HandleUserSetupRequest(platform utils.Platform, engDao userEngagement.DAO, event models.UserEngage, usersTable string) (string, error) {
	if event.IsNew {
		// For a new user, we post all the engagements for settings
		meetingTime(platform, event.UserId, usersTable)
	} else {
		if event.Update {
			// Checking if existing user wants to update settings
			settingsUpdateMessage(engDao, "", callback(event.UserId, UserSettingsUpdateTopic, user.AskForEngagements),
				event.UserId, event.ThreadTs, platform)
		}
	}
	return "ok", nil
}
