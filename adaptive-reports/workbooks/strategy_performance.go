package workbooks

import (
	"github.com/adaptiveteam/adaptive/daos/dialogEntry"

	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive/adaptive-reports/models"
	"github.com/adaptiveteam/adaptive/adaptive-reports/queries"
	"github.com/adaptiveteam/adaptive/adaptive-reports/utilities"
	"github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/styles"
	utilities2 "github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/utilities"
	"github.com/adaptiveteam/adaptive/daos/common"
	fetch_dialog "github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/pkg/errors"
)

// CreateStrategyWorkbook -
func CreateStrategyWorkbook(
	f *excel.File,
	platformID common.PlatformID,
	db *utilities.Database,
	dialogDAO fetch_dialog.DAO,
) (err error) {
	qm := utilities.InvokeQueries(
		utilities.InvokeQuery("strategy", queries.SelectStrategyStatusByPlatformID, string(platformID)),
		utilities.InvokeQuery("alignment", queries.SelectAlignmentSummaryByPlatformID, string(platformID)),
		utilities.InvokeQuery("vision", queries.SelectVisionByPlatformID, string(platformID)),
		utilities.InvokeQuery("competencies", queries.SelectCompetenciesByPlatformID, string(platformID)),
	)
	var queryResults utilities.QueryResultMap
	queryResults, err = qm.Run(db)

	if err == nil {
		allObjectives, allInitiatives := models.ConvertTableToObjectivesAndInitiatives(
			queryResults["strategy"].GetTable(),
			len(queryResults["strategy"].GetRows()),
		)
		allAlignments := models.CreateStrategyAlignments(
			queryResults["alignment"].GetTable(),
			len(queryResults["alignment"].GetRows()),
		)
		vision := queryResults["vision"].GetTable().GetValue("vision", 0)
		allCompetencies := models.CreateCompetencies(
			queryResults["competencies"].GetTable(),
			len(queryResults["competencies"].GetRows()),
		)

		instructions, summary, alignment, competencies := strategy.CreateStrategyWorksheets(
			f,
			"Instructions",
			"Performance Summary",
			"Alignment",
			"Competencies",
		)

		var dialog dialogEntry.DialogEntry
		dialog, err = dialogDAO.FetchByAlias("report-instructions", "instructions", "strategy")
		err = errors.Wrap(err, "error querying instructions for strategy performance report")
		if err == nil {
			strategy.CreateStrategySummary(f, summary, vision, &allObjectives, styles.Styles)
			strategy.CreateAlignmentSummary(f, alignment, allAlignments, allObjectives, allInitiatives, allCompetencies, styles.Styles)
			strategy.CreateCompetencies(f, competencies, allCompetencies, styles.Styles)
			utilities2.CreateInstructions(
				f,
				instructions,
				dialog.Dialog[0],
				styles.Styles,
			)
			summary.Activate()
			err = f.SetSheetVisible("Sheet1", false)
		}
	}

	return
}
