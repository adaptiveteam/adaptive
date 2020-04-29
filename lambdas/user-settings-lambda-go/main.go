package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/common"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
	"math/rand"
	"strings"
	"time"

	daosUser "github.com/adaptiveteam/adaptive/daos/user"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/userSetup"
)

var (
	usersTableName = utils.NonEmptyEnv("USERS_TABLE_NAME")
)

func HandleRequest(ctx context.Context, e events.SNSEvent) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("error in user-settings-lambda %v", err2)
		}
	}()

	for _, record := range e.Records {
		sns := record.SNS
		if sns.Message == "warmup" {
			fmt.Println("Warmed up the lambda")
		} else {
			var np models.NamespacePayload4
			err = json.Unmarshal([]byte(sns.Message), &np)
			core.ErrorHandler(err, namespace, fmt.Sprintf("Could not unmarshal stream image to NamespacePayload"))
			if np.Namespace == "settings" {
				connGen := common.CreateConnectionGenFromEnv()
				conn := connGen.ForPlatformID(np.TeamID.ToPlatformID())
				switch np.SlackRequest.Type {
				case models.InteractionSlackRequestType:
					byt, _ := json.Marshal(np.SlackRequest)
					log.Println("### Interaction event: " + string(byt))
					log.Println(fmt.Sprintf("### interactive_message event: %v", np.SlackRequest.ToEventsAPIEventUnsafe()))
					request := np.SlackRequest.InteractionCallback
					action := request.ActionCallback.AttachmentActions[0]

					userID := request.User.ID
					channelID := request.Channel.ID

					notes := responses()
					// Handling the init message
					if request.CallbackID == "init_message" {
						selected := action.SelectedOptions[0]
						text := selected.Value
						engage := models.UserEngage{UserID: userID, IsNew: true, Update: true, OnDemand: true, ThreadTs: request.MessageTs}
						switch text {
						case user.UpdateSettings:
							// For a new user, we post all the engagements for settings
							// meetingTime(userAttributeDao, platform, request.User.ID)
							handleUserSetupRequestUnsafe(engage)
						default:
						}
						notes = responses(
							models.PlatformSimpleNotification{UserId: userID,
								Channel: channelID, Ts: request.MessageTs},
						)
					} else {
						// Parse callback Id to messageCallback
						mc, err := utils.ParseToCallback(request.CallbackID)
						core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))

						if strings.HasPrefix(action.Name, string(models.Select)) {
							notes := responses()
							act := strings.TrimPrefix(action.Name, fmt.Sprintf("%s_", models.Select))
							switch act {
							case string(models.Now):
								// Query user table to get timezone
								user := daosUser.ReadUnsafe(userID)(conn)
								notes = selectMeetingTimeHandler(request, user.Timezone)
							case string(models.Cancel):
								// Remove the engagement
								notes = responses(
									models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: request.MessageTs},
								)
							}
							platform.PublishAll(notes)
							fmt.Printf("UpdateEngAsAnswered(...)\n")
							// Be it 'select' or 'cancel', are going to update the engagement as answered and won't remind the user again
							utils.UpdateEngAsAnswered(mc.Source, request.CallbackID, engTable, d, namespace)
							fmt.Printf("UpdateEngAsAnswered(...) completed\n")
						} else if mc.Action == user.AskForEngagements {
							if mc.Topic == userSetup.UserSettingsUpdateTopic {
								switch strings.TrimPrefix(action.Name, fmt.Sprintf("%s:", mc.Action)) {
								case string(models.Update):
									userEngage := models.UserEngage{UserID: userID, IsNew: true, Update: false}
									handleUserSetupRequestUnsafe(userEngage)
								case string(models.Cancel):
								}
								notes = responses(models.PlatformSimpleNotification{UserId: userID, Channel: channelID,
									Ts: request.OriginalMessage.Timestamp})
								// We have now added feedback for a coaching engagement. We can now update the original engagement as answered.
								fmt.Printf("UpdateEngAsAnswered2(...)\n")
								utils.UpdateEngAsAnswered(mc.Source, request.CallbackID, engTable, d, namespace)
								fmt.Printf("UpdateEngAsAnswered2(...) completed\n")
							}
						}
						return nil
					}
					platform.PublishAll(notes)
				default:
					fmt.Printf("### callback event of unsupported type %s: %v\n", np.SlackRequest.Type, np)
				}

			}
		}
	}

	return
}

func resetGlobalRNG() bool {
	rand.Seed(time.Now().Unix())
	return true
}

// initialize global pseudo random generator
// we should do it once per lambda instance, not every request.
var _ = resetGlobalRNG()

func dynString(str string) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{S: aws.String(str)}
	return &attr
}

// https://golang.org/ref/spec#Rune_literals
func SplitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}
	return subs
}

func attrDuration(str string) time.Duration {
	// time string consists of 4 characters, first 2 being hour and last 2 being minutes, in 24 hr format
	strs := SplitSubN(str, 2)
	formattedString := fmt.Sprintf("%sh%sm", strs[0], strs[1])
	duration, err := time.ParseDuration(formattedString)
	core.ErrorHandler(err, "namespace", "Could not parse string to duration")
	return duration
}

func userTimeInUTC(location string, setTime string) (utcTimeStr string, err error) {
	serverLocation, err := time.LoadLocation(location)
	if err == nil {
		// Get the beginning of the day for the location and add user selected time to it. So, you have user selected time for today
		userTime := core.Bod(time.Now().In(serverLocation)).Add(attrDuration(setTime))
		// convert time from user location to UTC
		utcTime := core.LocalToUtc(userTime)
		utcTimeStr = utcTime.Format("15") + utcTime.Format("04")
	}
	return
}

func selectMeetingTimeHandler(request slack.InteractionCallback, location string) []models.PlatformSimpleNotification {
	action := request.ActionCallback.AttachmentActions[0]
	selected := action.SelectedOptions[0]
	meetingTime, err := business_time.ParseLocalTimeID(selected.Value)
	core.ErrorHandler(err, namespace, "Couldn't parse meeting time "+selected.Value)
	// dbWriteUserAttributeUnsafe(request.User.ID, MeetingTimeUserAttributeID, meetingTime.ID(), false)

	keyParams := map[string]*dynamodb.AttributeValue{
		"id": dynString(request.User.ID),
	}
	meetTime := meetingTime.ID()
	utcMeetTime, err := userTimeInUTC(location, meetTime)
	core.ErrorHandler(err, namespace, "Couldn't convert user meeting time to UTC: "+meetTime)

	exprAttributes := map[string]*dynamodb.AttributeValue{
		":ast":  dynString(meetTime),
		":astu": dynString(utcMeetTime),
	}
	updateExpression := "set adaptive_scheduled_time = :ast, adaptive_scheduled_time_in_utc = :astu"
	err = d.UpdateTableEntry(exprAttributes, keyParams, updateExpression, usersTableName)
	core.ErrorHandler(err, namespace, "Couldn't update meeting time in %s table"+usersTableName)

	return responses(
		// Delete the original request
		utils.InteractionCallbackOverrideRequestMessage(request, ""),
		// Publish new request
		utils.InteractionCallbackSimpleResponse(request, string(MeetingIsScheduledFor(meetingTime))),
	)
}

func responses(notifications ...models.PlatformSimpleNotification) []models.PlatformSimpleNotification {
	return notifications
}

func handleUserSetupRequestUnsafe(userEngage models.UserEngage) {
	fmt.Printf("handleUserSetupRequestUnsafe()\n")
	userSetup.HandleUserSetupRequest(platform, userEngagementDao, userEngage, usersTableName)
	fmt.Printf("handleUserSetupRequestUnsafe() completed\n")
}
