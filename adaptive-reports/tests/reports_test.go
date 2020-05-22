package tests

import (
	"github.com/pkg/errors"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/sql-connector"
	"github.com/adaptiveteam/adaptive/daos/common"
	"fmt"
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive/adaptive-reports/utilities"
	"github.com/adaptiveteam/adaptive/adaptive-reports/workbooks"
	utilities2 "github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/utilities"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	fetch_dialog "github.com/adaptiveteam/adaptive/dialog-fetcher"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4"
	"os"
	"testing"
)

func createStrategySummary(
	db *utilities.Database,
	dao fetch_dialog.DAO,
	platformID common.PlatformID,
) (file *excel.File, err error) {
	file = excel.NewFile()

	properties := utilities2.CreateDocumentProperties(
		"Strategy",
		"How is the strategy performing?",
		[]string{"Strategy"},
		"Strategic Performance Report",
		"Strategic Performance",
	)
	err = file.SetDocProps(properties)
	if err == nil {
		err = workbooks.CreateStrategyWorkbook(
			file,
			platformID,
			db,
			dao,
		)
	}
	return file, err
}

func createIDOSummary(
	db *utilities.Database,
	dao fetch_dialog.DAO,
	userID string,
) (file *excel.File, err error) {
	file, err = workbooks.CreateIDOWorkbook(
		userID,
		db,
		dao,
		utilities2.CreateDocumentProperties(
			"IDO",
			"How are you doing on your IDO's?",
			[]string{"IDO"},
			"IDO Performance Report",
			"IDO Performance",
		),
	)
	return file, err
}

func TestCreateWorkbooks(t *testing.T) {

	// TODO: Reimplement test to be executable on travis
	if true { return }
/*
	AWS_REGION=us-east-2;
	DIALOG_TABLE=lexcorp_dialog_content;
	driver=mysql;
	end_point=lexcorp-reporting.chwrqdykifiq.us-east-2.rds.amazonaws.com;
	user=user;
	password=<this is in the AWS console for the reporting Lambda>;
	database=test_report;
	port=3306
*/
	dynamo := awsutils.NewDynamo(utils.NonEmptyEnv("AWS_REGION"), "", "dialog")
	dialogTableName := utils.NonEmptyEnv("DIALOG_TABLE")
	dialogDAO := fetch_dialog.NewDAO(dynamo, dialogTableName)

	conn, err2 := sqlconnector.ReadRDSConfigFromEnv().SQLOpen()
	if err2 != nil {
		panic(errors.Wrap(err2, "ReadRDSConfigFromEnv().SQLOpen()"))
	}
	db := utilities.WrapDB(conn)
	teamID := common.PlatformID("AGEGG1U7J")
	// teamID := common.PlatformID("ANT7U58AG")
	userIDs := map[string]string {
		"Ryan":"UFNLVKFT4",
		/*
		"April":"U38KRFVTQ",
		"Morgan":"UMA5H21FZ",
		"Erin":"UMAE07SR4",
		"Courtney":"ULWS36GP5",
		"Thomas":"ULTRB2D7F",
		"Michael":"ULCRWKDPE",
		*/
	}
	output := os.Getenv("output")
	defer db.CloseDatabase()
	var err error
	var file *excel.File

	file, err = createStrategySummary(db, dialogDAO, teamID)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		err = file.SaveAs(output+"strategy"+".xlsx")
	}

	for k, v := range userIDs {
		file, err = createIDOSummary(db, dialogDAO, v)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			err = file.SaveAs(output+"IDO - "+k+".xlsx")
		}
	}
}
