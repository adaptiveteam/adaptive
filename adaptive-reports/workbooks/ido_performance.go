package workbooks

import (
	"fmt"
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive-reports/models"
	"github.com/adaptiveteam/adaptive-reports/queries"
	"github.com/adaptiveteam/adaptive-reports/utilities"
	"github.com/adaptiveteam/adaptive-reports/worksheets/ido"
	"github.com/adaptiveteam/adaptive-reports/worksheets/styles"
	utilities2 "github.com/adaptiveteam/adaptive-reports/worksheets/utilities"
	fetch_dialog "github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/pkg/errors"
)

func CreateIDOWorkbook(
	userID string,
	db *utilities.Database,
	dialogDAO fetch_dialog.DAO,
	properties *excel.DocProperties,
) (f *excel.File, err error) {

	f = excel.NewFile()
	err = f.SetDocProps(properties)
	if err == nil {
		qm := utilities.NewQueryMap().
			AddToQuery("ido",
				queries.IDOs,
				userID,
			)

		var queryResults utilities.QueryResultMap
		queryResults, err = utilities.RunQueries(db, &qm)

		if err == nil {
			allIDOs := models.CreateIDOs(
				queryResults["ido"].GetTable(),
				len(queryResults["ido"].GetRows()),
			)
			instructions,
			summaryIDO,
			detailsIDO,
			activeIDOs,
			closedIDOs := ido.CreateIDOWorksheets(
				f,
				"instructions",
				"IDO Summary",
				"IDO Details",
				"Active IDO Updates",
				"Closed IDO Updates",
			)

			ido.CreateIDOUpdates(
				f,
				activeIDOs,
				allIDOs,
				false,
				styles.Styles,
			)

			ido.CreateIDOUpdates(
				f,
				closedIDOs,
				allIDOs,
				true,
				styles.Styles,
			)

			ido.CreateIDOSummary(
				f,
				summaryIDO,
				allIDOs,
				styles.Styles,
			)

			ido.CreateIDODetails(
				f,
				detailsIDO,
				allIDOs,
				styles.Styles,
			)

			dialog,err := dialogDAO.FetchByAlias("report-instructions", "instructions", "ido")
			if err != nil {
				fmt.Println(errors.Wrap(err, "error querying instructions for IDO performance report"))
			}
			utilities2.CreateInstructions(
				f,
				instructions,
				dialog.Dialog[0],
				styles.Styles,
			)
			summaryIDO.Activate()
			err = f.SetSheetVisible("Sheet1", false)
		}
	}

	if err != nil {
		fmt.Println(err)
	}
	return f, err
}
