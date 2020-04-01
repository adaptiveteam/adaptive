package workbooks

import (
	"github.com/adaptiveteam/adaptive/daos/dialogEntry"
	"fmt"

	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive/adaptive-reports/models"
	"github.com/adaptiveteam/adaptive/adaptive-reports/queries"
	"github.com/adaptiveteam/adaptive/adaptive-reports/utilities"
	"github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/ido"
	"github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/styles"
	utilities2 "github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/utilities"
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
				queries.SelectIDOsByUserID,
				userID,
			)

		var queryResults utilities.QueryResultMap
		queryResults, err = utilities.RunQueries(db, &qm)

		if err == nil {
			allIDOs := models.CreateIDOs(
				queryResults["ido"].GetTable(),
				len(queryResults["ido"].GetRows()),
			)
			instructions := ido.CreateIDOWorksheet(f, "instructions")
			summaryIDO := ido.CreateIDOWorksheet(f, "IDO Summary")
			detailsIDO := ido.CreateIDOWorksheet(f, "IDO Details")
			activeIDOs := ido.CreateIDOWorksheet(f, "Active IDO Updates")
			closedIDOs := ido.CreateIDOWorksheet(f, "Closed IDO Updates")

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

			var dialog dialogEntry.DialogEntry
			dialog, err = dialogDAO.FetchByAlias("report-instructions", "instructions", "ido")
			err = errors.Wrap(err, "error querying instructions for IDO performance report")
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
