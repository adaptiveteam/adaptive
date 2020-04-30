package lambda

import (
	"fmt"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/crosswalks"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/common"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/nlopes/slack"
)

var allCrosswalks = concatAppend(
	crosswalks.UserCrosswalk(),
)

func simulateSurvey(userID string,
	conn common.DynamoDBConnection,
	quarter, year int) ebm.AttachmentActionSurvey {
	var scheduleOpts []models.KvPair
	var userOpts []models.KvPair
	quarterStart := business_time.NewDateFromQuarter(quarter, year)
	quarterEnd := quarterStart.GetLastDayOfQuarter()

	allUsersSchedule := allSchedules(quarterStart, userID, quarterEnd.DaysBetween(quarterStart), conn)

	for _, each := range allUsersSchedule {
		date := each.ScheduledDate.DateToString(format)
		scheduleOpts = append(scheduleOpts,
			models.KvPair{Value: date, Key: date})
	}
	// Get user options
	userProfiles := user.ReadAllUserProfiles(conn)
	for _, each := range userProfiles {
		userOpts = append(userOpts, models.KvPair{Key: each.DisplayName, Value: each.Id})
	}

	actionElems := []ebm.AttachmentActionTextElement{
		selectControl(SimulateUserFieldID, SimulateUserLabel, userOpts),
		selectControl(SimulateDateFieldID, SimulateDateLabel, scheduleOpts),
	}
	fmt.Println(actionElems)
	return utils.AttachmentSurvey(string(SimulateDialogTitle), actionElems)
}

func simulateCurrentQuarterMenuHandler(request slack.InteractionCallback, mc models.MessageCallback, 
	conn common.DynamoDBConnection,
) {
	userID := request.User.ID
	callbackID := mc.WithAction(SimulateCurrentQuarterAction).ToCallbackID()
	y, m, d := time.Now().Date()
	bt := business_time.NewDate(y, int(m), d)
	val := simulateSurvey(userID, conn, bt.GetQuarter(), bt.GetYear())
	ut := userTokenSyncUnsafe(request.User.ID)
	api := slack.New(ut)
	err := utils.SlackSurvey(api, request, val, callbackID, surveyState(request, mc.Target))
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", callbackID+":"+request.CallbackID))
}

func simulateNextQuarterMenuHandler(request slack.InteractionCallback, mc models.MessageCallback, conn common.DynamoDBConnection) {
	userID := request.User.ID
	callbackID := mc.WithAction(SimulateNextQuarterAction).ToCallbackID()
	y, m, d := time.Now().Date()
	bt := business_time.NewDate(y, int(m), d)
	val := simulateSurvey(userID, conn, 
		bt.GetNextQuarter(), bt.GetNextQuarterYear())
	ut := userTokenSyncUnsafe(request.User.ID)
	api := slack.New(ut)
	err := utils.SlackSurvey(api, request, val, callbackID, surveyState(request, mc.Target))
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", callbackID+":"+request.CallbackID))
}

func currentQuarterScheduleMenuHandler(request slack.InteractionCallback, conn common.DynamoDBConnection) {
	userID := request.User.ID
	channelID := request.Channel.ID
	y, m, d := time.Now().Date()
	bt := business_time.NewDate(y, int(m), d)
	// Publishing 90 day summary to the user
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
		Message: string(ScheduleForCurrentQuarterTitle), Ts: request.MessageTs})
	// Post the schedules in the thread
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
		Message: string(schedulesSummary(userID, conn, bt.GetQuarter(), bt.GetYear())), ThreadTs: request.MessageTs})
}

func nextQuarterScheduleMenuHandler(request slack.InteractionCallback, conn common.DynamoDBConnection) {
	userID := request.User.ID
	channelID := request.Channel.ID
	y, m, d := time.Now().Date()
	bt := business_time.NewDate(y, int(m), d)
	// Publishing 90 day summary to the user
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
		Message: string(ScheduleForNextQuarterTitle),
		AsUser:  true, Ts: request.MessageTs})
	// Post the schedules in the thread
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
		Message: string(schedulesSummary(userID, conn, bt.GetNextQuarter(), bt.GetNextQuarterYear())),
		AsUser:  true, ThreadTs: request.MessageTs})
}
