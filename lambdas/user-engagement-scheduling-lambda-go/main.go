package lambda

import (
	"context"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/sirupsen/logrus"
	"log"
	"math"
	"strings"
	"time"
)

var (
	defaultMeetingTime = business_time.MeetingTime(9, 0)
	logger             = alog.LambdaLogger(logrus.InfoLevel)
)

func usersWithinScheduledPeriod(config Config, platformID string, startUTCTime, endUTCTime string) (users []models.User) {
	err := config.d.QueryTableWithIndex(config.usersTable, awsutils.DynamoIndexExpression{
		IndexName: config.usersScheduledTimeIndex,
		Condition: "platform_id = :pi AND adaptive_scheduled_time_in_utc BETWEEN :t1 AND :t2",
		Attributes: map[string]interface{}{
			":pi": platformID,
			":t1": startUTCTime,
			":t2": endUTCTime,
		},
	}, map[string]string{}, true, -1, &users)
	core.ErrorHandler(err, config.namespace, fmt.Sprintf("Could not query %s table on %s index",
		config.usersTable, config.usersScheduledTimeIndex))
	return
}

func usersWithinOffsetRange(config Config, platformID string, startOffset, endOffset int) (users []models.User) {
	err := config.d.QueryTableWithIndex(config.usersTable, awsutils.DynamoIndexExpression{
		IndexName: config.usersZoneOffsetIndex,
		Condition: "platform_id = :pi AND timezone_offset BETWEEN :t1 AND :t2",
		Attributes: map[string]interface{}{
			":pi": platformID,
			":t1": startOffset,
			":t2": endOffset,
		},
	}, map[string]string{}, true, -1, &users)
	core.ErrorHandler(err, config.namespace, fmt.Sprintf("Could not query %s table on %s index",
		config.usersTable, config.usersScheduledTimeIndex))
	return
}

func timeInHrMin(ip time.Time) string {
	return ip.Format("15") + ip.Format("04")
}

func absoluteOffsetFromUTC(nowTime time.Time, defaultTime time.Time) int {
	defaultTimeSeconds := defaultTime.Hour() * 3600 // 32400 for 9 AM
	return defaultTimeSeconds - (nowTime.Hour()*3600 + nowTime.Minute()*60)
}

func HandleRequest(ctx context.Context) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("error in user-engagement-scheduling-lambda-go %v", err2)
		}
	}()
	logger = logger.WithLambdaContext(ctx)
	config := readConfigFromEnvironment()
	// Query all the client configs
	var clientConfigs []models.ClientPlatformToken
	err = config.d.ScanTable(config.clientConfigTable, &clientConfigs)
	if err == nil {
		for _, clientConfig := range clientConfigs {
			// Query users for a client based on platform id
			platformID := clientConfig.PlatformID

			now := time.Now().UTC()
			// rounded to nearest hour quarter
			nowMinute := math.Floor(float64(now.Minute() / 15))
			// Dynamo index 'INBETWEEN' takes both inclusive, so we are reducing 1 minute on end, and taking to 59 seconds
			// an example would be, 10:00:00 AM to 10:14:59 AM
			hourQuarterStartMinute := int(15 * nowMinute)
			hourQuarterEndMinute := int(15*(nowMinute+1)) - 1
			startUTCTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), hourQuarterStartMinute,
				0, 0, time.UTC)
			endUTCTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), hourQuarterEndMinute,
				59, 0, time.UTC)

			fmt.Println(fmt.Sprintf("Invoking user schedules for %s UTC", startUTCTime.String()))
			// Querying for users who explicitly set the time
			usersToAsk1 := usersWithinScheduledPeriod(config, string(platformID), timeInHrMin(startUTCTime), timeInHrMin(endUTCTime))
			log.Println(fmt.Sprintf("No. of users invoked with scheduled time: %d", len(usersToAsk1)))
			// Querying for others that have the default time and should be invoked now
			offsetEnd := absoluteOffsetFromUTC(startUTCTime, time.Time(defaultMeetingTime))
			offsetStart := absoluteOffsetFromUTC(endUTCTime, time.Time(defaultMeetingTime))
			usersToAsk2 := usersWithinOffsetRange(config, string(platformID), offsetStart, offsetEnd)
			var usersToAsk2Filtered []models.User
			// Filtering out the users who have explicitly set scheduled time same as default time since these are part of `usersToAsk1`
			for _, each := range usersToAsk2 {
				// also filtering out the users whose display name begins with 'adaptive' since they are community related
				if each.AdaptiveScheduledTimeInUTC == "" && !strings.HasPrefix(each.DisplayName, "adaptive") {
					usersToAsk2Filtered = append(usersToAsk2Filtered, each)
				}
			}
			log.Println(fmt.Sprintf("Offset range to look for users with no scheduled time: [%d, %d]", offsetStart, offsetEnd))
			log.Println(fmt.Sprintf("No. of users invoked with default time: %d", len(usersToAsk2Filtered)))

			for _, user := range append(usersToAsk1, usersToAsk2Filtered...) {
				// Engage with the user only if the user is a part of an Adaptive community (exludes Admin Community)
				userCommunities := strategy.QueryCommunityUserIndex(user.ID, config.communityUsersTable, config.communityUsersUserIndex)
				if len(userCommunities) == 1 && userCommunities[0].CommunityId == string(community.Admin) {
					logger.Infof("%s user belongs only to Admin Community, not invoking schedules for this user", user.ID)
				} else if len(userCommunities) > 0 {
					engage := models.UserEngage{
						UserId:     user.ID,
						PlatformID: platformID,
					}
					err = invokeScriptingLambda(engage, config)
					if err != nil {
						logger.WithError(err).Errorf("Could not invoke scripting lambda for %s user in %v platform", engage.UserId, platformID)
					}
					// err = triggerPostponedEvents(engage, config)
					// if err != nil {
					// 	logger.WithError(err).Errorf("Could not TriggerAllPostponedEvents for %s user in %v platform", engage.UserId, platformID)
					// }
					switch platformID {
					case EmbursePlatformID, GeigsenPlatformID:
						emulateDates(EmburseDateShiftConfig, time.Now(), user.ID, platformID, config)
					case IvanPlatformID, StagingPlatformID:
						emulateDates(TestDateShiftConfig, time.Now(), user.ID, platformID, config)
					default:
					}
				}
			}
		}
	}
	return
}
