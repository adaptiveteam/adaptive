package strategy

import (
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive-reports/models"
)

func CreateCompetencies(
	f *excel.File,
	competenciesSheet *models.Sheet,
	competencies []models.Competency,
	styles map[string]string,
) {

	const backToSummaryColumn = 1
	const titleRow = 1
	const titleColumn = 2
	const titleHeight = 50
	const title = "Team Competencies"
	const competencyNameWidth = 40
	const competencyHeight = 400
	const competencyDescriptionWidth = 100
	const competencyNameColumn = 1
	const competencyDescriptionColumn = 2

	// Create Title
	competenciesSheet.NewMergedCell(
		competencyNameColumn, titleRow,
		competencyDescriptionColumn, titleRow,
	).Value(title).
		Style(
			models.NewStyle(styles).
				GetStyle("Title Font").
				GetStyle("Heading Background").
				GetStyle("Centered").
				GetStyle("White Bottom Border").
				ToNewStyle(f),
		).Height(titleHeight)

	for i := 0; i < len(competencies); i++ {
		// Label for the Update
		competenciesSheet.NewCell(
			competencyNameColumn,
			titleRow+i+1,
		).Value(
			competencies[i].Name+"\n("+competencies[i].Type+")",
		).Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Vertical Center").
				GetStyle("White Borders").
				ToNewStyle (f),
		).Width(competencyNameWidth).Height(competencyHeight)

		// The actual update
		competenciesSheet.NewCell(
			competencyDescriptionColumn,
			titleRow+i+1,
		).Value(
			competencies[i].Description,
		).Style(
			models.NewStyle(styles).
				GetStyle("Black Borders").
				GetStyle("Normal Font").
				GetStyle("Normal Background").
				GetStyle("Vertical Center").
				ToNewStyle (f),
		).Width(competencyDescriptionWidth)
	}
}
