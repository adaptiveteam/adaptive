package userEngagementScheduling

import (
	sqlconnector "github.com/adaptiveteam/adaptive/adaptive-utils-go/sql-connector"
	"github.com/pkg/errors"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-reports/stats"
	"github.com/adaptiveteam/adaptive/adaptive-reports/utilities"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

func generateFeedbackStatsReport(config Config, teamID models.TeamID) (report ui.RichText, err error) {
	logger.Infof("generateFeedbackStatsReport(%s)", teamID.ToString())
	year, quarter := core.CurrentYearQuarter()
	rdsConfig := sqlconnector.ReadRDSConfigFromEnv()
	sqlConn := rdsConfig.SQLOpenUnsafe()
	defer utilities.CloseUnsafe(sqlConn)
	var stat stats.FeedbackStats
	stat, err = stats.QueryFeedbackStats(teamID, year, quarter)(sqlConn)
	if err == nil {
		report = ui.Sprintf(
`People who have given feedback - %0.2f%%
People who have received feedback - %0.2f%%`,
			stat.Given, stat.Received)
		logger.Info(report)
	}
	return
}

func sendFeedbackStatsReport(config Config, teamID models.TeamID) (err error) {
	var report ui.RichText
	report, err = generateFeedbackStatsReport(config, teamID)
	if err == nil {
		err = postToCommunity(platform.Message(report), community.User)(config, teamID)
	}
	return errors.Wrap(err, "sendFeedbackStatsReport")
}
