package lambda

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	utilities "github.com/adaptiveteam/adaptive/adaptive-reports/utilities"
	"github.com/adaptiveteam/adaptive/adaptive-reports/workbooks"
	utilities2 "github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/utilities"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	// This import is needed for reports to work
	_ "github.com/go-sql-driver/mysql"
	"github.com/nlopes/slack"
)

func onStrategyPerformanceReport(RDSConfig RDSConfig, platformID models.PlatformID) (buf *bytes.Buffer, reportname string, err error) {
	defer recoverToErrorVar("onStrategyPerformanceReport", &err)
	customerID := string(platformID)
	db := utilities.SqlOpenUnsafe(RDSConfig.Driver, RDSConfig.ConnectionString)
	defer utilities.CloseUnsafe(db)
	var file *excelize.File

	// f, err = reports.StrategyPerformanceReport(db, customerID)

	reportname = "Strategic Performance"
	file, err = workbooks.CreateStrategyWorkbook(
		customerID,
		utilities.WrapDB(db),
		dialogFetcherDAO,
		utilities2.CreateDocumentProperties(
			"Strategy",
			"How is the strategy performing?",
			[]string{"Strategy"},
			"Strategic Performance Report",
			reportname,
		),
	)
	if err == nil {
		buf, err = file.WriteToBuffer()
	}
	return
}

func onIDOPerformanceReport(RDSConfig RDSConfig, userID string) (buf *bytes.Buffer, reportname string, err error) {
	defer recoverToErrorVar("onIDOPerformanceReport", &err)
	db := utilities.SqlOpenUnsafe(RDSConfig.Driver, RDSConfig.ConnectionString)
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
	return
}

// saves the report to s3 bucket
func uploadReportToS3(buf *bytes.Buffer, name string) (err error) {
	return errors.New("Not implemented: uploadReportToS3")
}

func sendReportToUser(
	platformID models.PlatformID,
	userID,
	name string,
	buf *bytes.Buffer,
) (err error) {
	defer recoverToErrorVar("sendReportToUser", &err)
	token := platformTokenDAO.GetPlatformTokenUnsafe(platformID)
	api := slack.New(token)

	params := slack.FileUploadParameters{
		Title:           name + " Report",
		Filename:        name + ".xlsx",
		Reader:          buf,
		Channels:        []string{userID},
		ThreadTimestamp: "",
	}
	_, err = api.UploadFile(params)
	return
}
