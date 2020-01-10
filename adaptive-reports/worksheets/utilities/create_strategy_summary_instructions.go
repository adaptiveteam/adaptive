package utilities

import (
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive-reports/models"
)

func CreateInstructions(
	f *excel.File,
	instructions *models.Sheet,
	text string,
styles map[string]string,
) {
	const instructionsHeight = 75
	const instructionsWidth = 100
	const instructionsColumn = 1
	const startingRow = 1

	// Now add some instructions
	instructions.NewCell(
		instructionsColumn,
		startingRow,
	).Value(
		text,
	).Style(
		models.NewStyle(styles).
			GetStyle("Black Borders").
			GetStyle("Normal Font").
			GetStyle("Instructions Background").
			GetStyle("Vertical Center").
			ToNewStyle (f),
	).Height(
		instructionsHeight,
	).Width(instructionsWidth)
}
