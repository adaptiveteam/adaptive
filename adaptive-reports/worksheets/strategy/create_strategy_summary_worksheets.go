package strategy

import (
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive-reports/models"
)

func CreateStrategyWorksheets(
	f *excel.File,
	instructionsTitle string,
	strategicSummaryTitle string,
	alignmentTitle string,
	competencyTitle string,
) (
	instructions *models.Sheet,
	strategicSummary *models.Sheet,
	strategicAlignment *models.Sheet,
	competencies *models.Sheet,
) {
	// Create a new instructions sheet
	instructions = models.NewSheet(
		f,
		instructionsTitle,
	).Landscape()

	// Create a new strategicSummary for the summary
	strategicSummary = models.NewSheet(
		f,
		strategicSummaryTitle,
	).Landscape()

	// Create another strategicSummary for the alignment
	// Create a new strategicSummary for the summary
	strategicAlignment = models.NewSheet(
		f,
		alignmentTitle,
	).Landscape()

	// Create another competencySummary
	competencies = models.NewSheet(
		f,
		competencyTitle,
	).Landscape()

	return instructions,  strategicSummary,  strategicAlignment, competencies
}
