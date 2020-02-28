package tests

import (
	"fmt"
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive/adaptive-reports/utilities"
	"github.com/adaptiveteam/adaptive/adaptive-reports/workbooks"
	utilities2 "github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/utilities"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	fetch_dialog "github.com/adaptiveteam/adaptive/dialog-fetcher"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"testing"
)

func createStrategySummary(
	db *utilities.Database,
	dao fetch_dialog.DAO,
	teamID string,
) (file *excel.File, err error) {
	file, err = workbooks.CreateStrategyWorkbook(
		teamID,
		db,
		dao,
		utilities2.CreateDocumentProperties(
			"Strategy",
			"How is the strategy performing?",
			[]string{"Strategy"},
			"Strategic Performance Report",
			"Strategic Performance",
		),
	)

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
	dynamo := awsutils.NewDynamo(utils.NonEmptyEnv("AWS_REGION"), "", "dialog")
	dialogTableName := "f" //utils.NonEmptyEnv("DIALOG_TABLE")
	dialogDAO := fetch_dialog.NewDAO(dynamo, dialogTableName)

	db := utilities.NewDatabase(
		os.Getenv("driver"),
		os.Getenv("end_point"),
		os.Getenv("user"),
		os.Getenv("password"),
		os.Getenv("port"),
		os.Getenv("database"),
	)
	teamID := "ANT7U58AG"
	userIDs := map[string]string {
		"April":"U38KRFVTQ",
		"Morgan":"UMA5H21FZ",
		"Erin":"UMAE07SR4",
		"Courtney":"ULWS36GP5",
		"Thomas":"ULTRB2D7F",
		"Michael":"ULCRWKDPE",
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
