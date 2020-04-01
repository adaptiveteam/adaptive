package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/slackTeam"
	"github.com/adaptiveteam/adaptive/daos/clientPlatformToken"
	"github.com/adaptiveteam/adaptive/adaptive-reports/utilities"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"

	// "database/sql"
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-reports/stats"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	_ "github.com/adaptiveteam/adaptive/daos"
	"github.com/sirupsen/logrus"
)

var (
	defaultMeetingTime = business_time.MeetingTime(9, 0)
	globalScheduleTime = func() time.Time {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(),
			12, 0, // 12:00 UTC = 08:00 EDT
			0, 0, time.UTC)
	}
	logger = alog.LambdaLogger(logrus.InfoLevel)
)

func usersWithinScheduledPeriod(config Config, teamID models.TeamID, startUTCTime, endUTCTime string) (users []models.User) {
	err := config.d.QueryTableWithIndex(config.usersTable, awsutils.DynamoIndexExpression{
		IndexName: config.usersScheduledTimeIndex,
		Condition: "platform_id = :pi AND adaptive_scheduled_time_in_utc BETWEEN :t1 AND :t2",
		Attributes: map[string]interface{}{
			":pi": teamID.ToString(),
			":t1": startUTCTime,
			":t2": endUTCTime,
		},
	}, map[string]string{}, true, -1, &users)
	core.ErrorHandler(err, config.namespace, fmt.Sprintf("Could not query %s table on %s index",
		config.usersTable, config.usersScheduledTimeIndex))
	return
}

func usersWithinOffsetRange(config Config, teamID models.TeamID, startOffset, endOffset int) (users []models.User) {
	err := config.d.QueryTableWithIndex(config.usersTable, awsutils.DynamoIndexExpression{
		IndexName: config.usersZoneOffsetIndex,
		Condition: "platform_id = :pi AND timezone_offset BETWEEN :t1 AND :t2",
		Attributes: map[string]interface{}{
			":pi": teamID.ToString(),
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
	defer core.RecoverAsLogError("user-engagement-scheduling-lambda-go")
	logger = logger.WithLambdaContext(ctx)
	defer func() {
		if err != nil {
			logger.WithError(err).Errorf("ERROR HandleRequest: %+v\n", err)
			err = nil
		}
	}()
	config := readConfigFromEnvironment()

	var teamIDs []models.TeamID
	teamIDs, err = getTeamIDsFromClientConfigs(config)
	if err == nil {
		var teamIDs2 []models.TeamID
		teamIDs2, err = getTeamIDsFromSlackTeams(config)
		teamIDs = append(teamIDs, teamIDs2...)
		// Query all the client configs
		if err == nil {
			for _, teamID := range teamIDs {
				err = runScheduleForTeam(config, teamID)
				if err != nil {
					logger.WithError(err).Errorf("HandleRequest.runScheduleForTeam")
				}
				err = runGlobalScheduleForTeam(config, teamID)
				if err != nil {
					logger.WithError(err).Errorf("HandleRequest.runGlobalScheduleForTeam")
				}
			}
		}
	}
	return
}

func getTeamIDsFromClientConfigs(config Config) (teamIDs []models.TeamID, err error) {
	// Query all the client configs
	var clientConfigs []models.ClientPlatformToken
	err = config.d.ScanTable(clientPlatformToken.TableName(config.clientID), &clientConfigs)
	if err == nil {
		for _, clientConfig := range clientConfigs {
			teamID := models.ParseTeamID(clientConfig.PlatformID)
			teamIDs = append(teamIDs, teamID)
		}
	}
	return
}
func getTeamIDsFromSlackTeams(config Config) (teamIDs []models.TeamID, err error) {
	var slackTeams []slackTeam.SlackTeam
	err = config.d.ScanTable(slackTeam.TableName(config.clientID), &slackTeams)
	if err == nil {
		for _, slackTeam := range slackTeams {
			teamID := models.ParseTeamID(slackTeam.TeamID)
			teamIDs = append(teamIDs, teamID)
		}
	}
	return
}
func getCurrentQuarterHourInterval() (startUTCTime, endUTCTime time.Time) {
	now := time.Now().UTC()
	// rounded to nearest hour quarter
	hourQuarterStartMinute := int(math.Floor(float64(now.Minute()/15)) * 15)
	// Dynamo index 'INBETWEEN' takes both inclusive, so we are reducing 1 minute on end, and taking to 59 seconds
	// an example would be, 10:00:00 AM to 10:14:59 AM
	hourQuarterEndMinute := hourQuarterStartMinute + 14
	startUTCTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(),
		hourQuarterStartMinute, 0,
		0, time.UTC)
	endUTCTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(),
		hourQuarterEndMinute, 59,
		0, time.UTC)
	return
}

// Query users for a client based on platform id
func runScheduleForTeam(config Config, teamID models.TeamID) (err error) {
	startUTCTime, endUTCTime := getCurrentQuarterHourInterval()

	fmt.Println(fmt.Sprintf("Invoking user schedules for [%s, %s] UTC", startUTCTime.String(), endUTCTime.String()))
	// Querying for users who explicitly set the time
	usersToAsk1 := usersWithinScheduledPeriod(config, teamID, timeInHrMin(startUTCTime), timeInHrMin(endUTCTime))
	log.Println(fmt.Sprintf("No. of users to be invoked with the scheduled time: %d", len(usersToAsk1)))
	// Querying for others that have the default time and should be invoked now
	offsetEnd := absoluteOffsetFromUTC(startUTCTime, time.Time(defaultMeetingTime))
	offsetStart := absoluteOffsetFromUTC(endUTCTime, time.Time(defaultMeetingTime))
	usersToAsk2 := usersWithinOffsetRange(config, teamID, offsetStart, offsetEnd)
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
				UserID: user.ID,
				Date:   "", // current date
				TeamID: teamID,
			}
			switch teamID.AppID {
			//					case EmbursePlatformID, GeigsenPlatformID:
			//						emulateDates(EmburseDateShiftConfig, time.Now(), user.ID, teamID, config)
			case IvanPlatformID, StagingPlatformID:
				emulateDates(TestDateShiftConfig, time.Now(), user.ID, teamID, config)
			default:
			}
			err = invokeScriptingLambda(engage, config)
			if err != nil {
				logger.WithError(err).Errorf("Could not invoke scripting lambda for %s user in %v platform", engage.UserID, teamID)
			}
		}
	}
	return
}

func runGlobalScheduleForTeam(config Config, teamID models.TeamID) (err error) {
	startUTCTime, endUTCTime := getCurrentQuarterHourInterval()
	scheduleTimeForToday := globalScheduleTime()
	if !scheduleTimeForToday.Before(startUTCTime) &&
		!scheduleTimeForToday.After(endUTCTime) {
		logger.Infof("runGlobalScheduleForTeam(%s)", teamID.ToString())
		year, quarter := core.CurrentYearQuarter()
		rdsConfig := utilities.ReadRDSConfigFromEnv()
		sqlConn := rdsConfig.SQLOpenUnsafe()
		defer utilities.CloseUnsafe(sqlConn)
		var stat stats.FeedbackStats
		stat, err = stats.QueryFeedbackStats(teamID, year, quarter)(sqlConn)
		if err == nil {
			message := ui.Sprintf(
				`People who have given feedback - %0.2f%%
People who have received feedback - %0.2f%%`,
				stat.Given, stat.Received)
			logger.Info(message)
			conn := config.connGen.ForPlatformID(teamID.ToPlatformID())
			var communities []adaptiveCommunity.AdaptiveCommunity
			communities, err = adaptiveCommunity.ReadOrEmpty(teamID.ToPlatformID(), string(community.User))(conn)
			if err == nil {
				if len(communities) > 0 {
					userComm := communities[0]
					slackAdapter := mapper.SlackAdapterForTeamID(conn)
					post := platform.Post(
						platform.ConversationID(userComm.ChannelID),
						platform.Message(message),
					)
					logger.Infof("Posting to %s: %v", userComm.ChannelID, post)
					_, err = slackAdapter.PostSync(post)
				} else {
					logger.Warnf("HR community not found for team %s", teamID.ToString())
				}
			}
		}
	} else {
		logger.Infof("runGlobalScheduleForTeam(%s) - skipping. Today it's planned to %v", teamID.ToString(), scheduleTimeForToday)
	}
	return
}
