package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/core-utils-go"
	"bytes"
	"time"

	"github.com/pkg/errors"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	utilities "github.com/adaptiveteam/adaptive/adaptive-reports/utilities"
	"github.com/adaptiveteam/adaptive/adaptive-reports/workbooks"
	utilities2 "github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/utilities"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	// This import is needed for reports to work
	_ "github.com/go-sql-driver/mysql"
	"github.com/nlopes/slack"
)

func onStrategyPerformanceReport(RDSConfig RDSConfig, teamID models.TeamID) (buf *bytes.Buffer, reportname string, err error) {
	defer recoverToErrorVar("onStrategyPerformanceReport", &err)
	logger.Infof("onStrategyPerformanceReport(teamID=%s", teamID)
	platformID := teamID.ToPlatformID()
	db := utilities.SQLOpenUnsafe(RDSConfig.Driver, RDSConfig.ConnectionString)
	defer utilities.CloseUnsafe(db)
	var file *excelize.File

	// f, err = reports.StrategyPerformanceReport(db, customerID)

	// TODO Please help me set the time-zone Bharath or Arseniy
	loc, _ := time.LoadLocation("America/Indianapolis")
	timeString := time.Now().In(loc).Format(time.RFC3339)
	reportname = "Strategic Performance, "+timeString
	file = excelize.NewFile()
	properties := utilities2.CreateDocumentProperties(
		"Strategy",
		"How is the strategy performing?",
		[]string{"Strategy"},
		"Strategic Performance Report",
		reportname,
	)
	err = file.SetDocProps(properties)
	if err == nil {
		err = workbooks.CreateStrategyWorkbook(
			file,
			platformID,
			utilities.WrapDB(db),
			dialogFetcherDAO,
		)
		if err == nil {
			buf, err = file.WriteToBuffer()
		}
	}
	return
}

func onIDOPerformanceReport(RDSConfig RDSConfig, userID string) (buf *bytes.Buffer, reportname string, err error) {
	defer core_utils_go.RecoverToErrorVar("onIDOPerformanceReport", &err)
	logger.Infof("onIDOPerformanceReport")
	db := utilities.SQLOpenUnsafe(RDSConfig.Driver, RDSConfig.ConnectionString)
	defer utilities.CloseUnsafe(db)
	var file *excelize.File

	// f, err = reports.StrategyPerformanceReport(db, customerID)

	reportname = "IDO Performance"
	file, err = workbooks.CreateIDOWorkbook(
		userID,
		utilities.WrapDB(db),
		dialogFetcherDAO,
		utilities2.CreateDocumentProperties(
			"IDO",
			"How are you doing on your IDO's?",
			[]string{"IDO"},
			"IDO Performance Report",
			reportname,
		),
	)
	if err == nil {
		buf, err = file.WriteToBuffer()
	}
	err = errors.Wrap(err, "onIDOPerformanceReport")
	return
}

// saves the report to s3 bucket
func uploadReportToS3(buf *bytes.Buffer, name string) (err error) {
	return errors.New("Not implemented: uploadReportToS3")
}

func sendReportToUser(
	teamID models.TeamID,
	userID,
	name string,
	buf *bytes.Buffer,
	conn daosCommon.DynamoDBConnection,
) (err error) {
	defer core_utils_go.RecoverToErrorVar("sendReportToUser", &err)
	filename := name + ".xlsx"
	logger.Infof("Sending report %s (size=%d b) to user %s", filename, buf.Len(), userID)
	var token string
	token, err = platform.GetToken(teamID)(conn)
	if err == nil {
		api := slack.New(token)

		params := slack.FileUploadParameters{
			Title:           name + " Report",
			Filename:        filename,
			Reader:          buf,
			Channels:        []string{userID},
			ThreadTimestamp: "",
		}
		var slackFile *slack.File
		slackFile, err = api.UploadFile(params)
		if err == nil {
			logger.Infof("Slack file: %v", slackFile)
		}
	} 
	if err != nil {
		logger.WithError(err).Errorf("Error while uploading file %s", filename)
	}
	return
}
