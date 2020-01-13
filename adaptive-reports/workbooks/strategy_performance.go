package workbooks

import (
	"fmt"
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive/adaptive-reports/models"
	"github.com/adaptiveteam/adaptive/adaptive-reports/queries"
	"github.com/adaptiveteam/adaptive/adaptive-reports/utilities"
	"github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/styles"
	utilities2 "github.com/adaptiveteam/adaptive/adaptive-reports/worksheets/utilities"
	"github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/pkg/errors"
)

func CreateStrategyWorkbook(
	platformID string,
	db *utilities.Database,
	dialogDAO fetch_dialog.DAO,
	properties *excel.DocProperties,
) (f *excel.File, err error) {

	f = excel.NewFile()
	err = f.SetDocProps(properties)
	if err == nil {
		qm := utilities.NewQueryMap().
			AddToQuery("strategy",
				queries.StrategyStatus,
				platformID,
			).AddToQuery(
				"alignment",
				queries.AlignmentSummary,
				platformID,
			).AddToQuery(
			"vision",
			queries.Vision,
				platformID,
			).AddToQuery(
			"competencies",
			queries.Competencies,
			platformID,
		)
		queryResults, err := utilities.RunQueries(db, &qm)

		if err == nil {
			allObjectives, allInitiatives := models.CreateObjectives(
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

			dialog,err := dialogDAO.FetchByAlias("report-instructions", "instructions", "strategy")
			if err != nil {
				fmt.Println(errors.Wrap(err, "error querying instructions for strategy performance report"))
			}
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

	return f, err
}
