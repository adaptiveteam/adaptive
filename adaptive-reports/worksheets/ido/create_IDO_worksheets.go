package ido

import (
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive-reports/models"
)

func CreateIDOWorksheets(
	f *excel.File,
	instructionsTitle string,
	idoSummaryTitle string,
	idoDetailsTitle string,
	activeIDOTitle string,
	closedIDOTitle string,
) (
	instructions *models.Sheet,
	idoSummary *models.Sheet,
	idoDetails *models.Sheet,
	activeIDO *models.Sheet,
	closedIDO *models.Sheet,
) {
	// Create a new instructions sheet
	instructions = models.NewSheet(
		f,
		instructionsTitle,
	).Landscape()

	// Create a new strategicSummary for the summary
	idoSummary = models.NewSheet(
		f,
		idoSummaryTitle,
	).Landscape()

	// Create a new worksheet for the details
	idoDetails = models.NewSheet(
		f,
		idoDetailsTitle,
	).Landscape()

	// Create another strategicSummary for the alignment
	// Create a new strategicSummary for the summary
	activeIDO = models.NewSheet(
		f,
		activeIDOTitle,
	).Landscape()

	// Create another competencySummary
	closedIDO = models.NewSheet(
		f,
		closedIDOTitle,
	).Landscape()

	return instructions,  idoSummary,  idoDetails, activeIDO, closedIDO
}

