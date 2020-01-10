package strategy

import (
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive-reports/models"
	"strconv"
)

func CreateComponentDetails(
	componentName string,
	f *excel.File,
	styles map[string]string,
	objective models.StrategyComponent,
) (objectiveSheet *models.Sheet, bottomRow, farRightColumn int) {

	const backToSummaryColumn = 1
	const titleRow = 1
	const titleColumn = 2
	const titleHeight = 50
	const nameRow = 2
	const nameLabelColumn = 1
	const nameValueColumn = 2
	const advocateRow = 3
	const advocateLabelColumn = 1
	const advocateValueColumn = 2
	const statusRow = 4
	const statusLabelColumn = 1
	const statusValueColumn = 2
	const updatedRow = 4
	const updatedLabelColumn = 3
	const updatedValueColumn = 4
	const createdOnRow = 4
	const createdOnDateLabelColumn = 5
	const createdOnDateValueColumn = 6
	const endDateRow = 4
	const endDateLabelColumn = 7
	const endDateValueColumn = 8
	const timeUsedRow = 4
	const timeUsedDateLabelColumn = 9
	const timeUsedDateValueColumn = 10
	const updateRow = 5
	const updateLabelColumn = 1
	const updateValueColumn = 2
	const updateHeight = 210
	const descriptionRow = 6
	const descriptionLabelColumn = 1
	const descriptionValueColumn = 2
	const descriptionHeight = 200
	const leftColumnWidth = 12
	const labelWidth = 10
	const valueWidth = 12
	farRightColumn = 10
	bottomRow = 6

	objectiveSheet = models.NewSheet(
		f,
		componentName+" "+strconv.Itoa(objective.GetIndex()),
	).Landscape()

	// Create Back Button
	objectiveSheet.NewCell(
		backToSummaryColumn,
		titleRow,
	).Value(
		"Back to Summary",
	).Style(
		models.NewStyle(styles).
			GetStyle("Return Font").
			GetStyle("Return Background").
			GetStyle("Vertical Center").
			GetStyle("Horizontal Center").
			GetStyle("White Borders").
			ToNewStyle (f),
	).SetLink("Performance Summary","A",1)

	// Create Title
	objectiveSheet.NewMergedCell(
		titleColumn, titleRow,
		farRightColumn, titleRow,
	).Value(
		componentName+" Performance",
	).Style(
		models.NewStyle(styles).
			GetStyle("Title Font").
			GetStyle("Heading Background").
			GetStyle("Centered").
			GetStyle("White Bottom Border").
			ToNewStyle (f),
	).Height(titleHeight)

	// Heading for the Objective name
	objectiveSheet.NewCell(
		nameLabelColumn,nameRow,
	).Value(
		"Name",
	).Style(
		models.NewStyle(styles).
			GetStyle("Heading Font").
			GetStyle("Heading Background").
			GetStyle("Vertical Center").
			GetStyle("White Borders").
			ToNewStyle (f),
	).Width(leftColumnWidth)

	// The actual objective name.
	objectiveSheet.NewMergedCell(
		nameValueColumn,
		nameRow,
		farRightColumn,
		nameRow,
	).Value(
		objective.GetName(),
	).Style(
		models.NewStyle(styles).
			GetStyle("Black Borders").
			GetStyle("Normal Font").
			GetStyle("Normal Background").
			ToNewStyle (f),
	)

	// Heading for the Advocate
	objectiveSheet.NewCell(
		advocateLabelColumn,
		advocateRow,
	).Value(
		"Advocate",
	).Style(
		models.NewStyle(styles).
			GetStyle("Heading Font").
			GetStyle("Heading Background").
			GetStyle("Vertical Center").
			ToNewStyle (f),
	).Width(leftColumnWidth)

	// Actual Advocate name
	objectiveSheet.NewMergedCell(
		advocateValueColumn,
		advocateRow,
		farRightColumn,
		advocateRow,
	).Value(
		objective.GetAdvocate(),
	).Style(
		models.NewStyle(styles).
			GetStyle("Black Borders").
			GetStyle("Normal Font").
			GetStyle("Normal Background").
			ToNewStyle (f),
	)

	// Heading for the Status
	objectiveSheet.NewCell(
		statusLabelColumn,
		statusRow,
	).Value(
		"Status",
	).Style(
		models.NewStyle(styles).
			GetStyle("Heading Font").
			GetStyle("Heading Background").
			GetStyle("Vertical Center").
			GetStyle("White Borders").
			ToNewStyle (f),
	).Width(leftColumnWidth)

	// Actual value for the Status
	objectiveSheet.NewCell(
		statusValueColumn,
		statusRow,
	).Value(
		objective.GetStatus(),
	).Style(
		models.NewStyle(styles).
			GetStyle(objective.GetStatus()).
			GetStyle("Black Borders").
			GetStyle("Italics Font").
			GetStyle("Vertical Center").
			ToNewStyle (f),
	).Width(valueWidth)

	// Label for the updated date
	objectiveSheet.NewCell(
		updatedLabelColumn,
		updatedRow,
	).Value(
		"Updated",
	).Style(
		models.NewStyle(styles).
			GetStyle("Heading Font").
			GetStyle("Heading Background").
			GetStyle("Vertical Center").
			ToNewStyle (f),
	).Width(labelWidth)

	// Actual updated date
	objectiveSheet.NewCell(
		updatedValueColumn,
		updatedRow,
	).Value(
		objective.GetUpdatedOn(),
	).Style(
		models.NewStyle(styles).
			GetStyle("Black Borders").
			GetStyle("Normal Font").
			GetStyle("Normal Background").
			ToNewStyle (f),
	).Width(valueWidth)

	// Label for the end date
	objectiveSheet.NewCell(
		endDateLabelColumn,
		endDateRow,
	).Value(
		"End",
	).Style(
		models.NewStyle(styles).
			GetStyle("Heading Font").
			GetStyle("Heading Background").
			GetStyle("Vertical Center").
			ToNewStyle (f),
	).Width(labelWidth)

	// Actual end date
	objectiveSheet.NewCell(
		endDateValueColumn,
		endDateRow,
	).Value(
		objective.GetEndDate(),
	).Style(
		models.NewStyle(styles).
			GetStyle("Black Borders").
			GetStyle("Normal Font").
			GetStyle("Normal Background").
			ToNewStyle (f),
	).Width(valueWidth)

	// Label for created on
	objectiveSheet.NewCell(
		createdOnDateLabelColumn,
		createdOnRow,
	).Value(
		"Start",
	).Style(
		models.NewStyle(styles).
			GetStyle("Heading Font").
			GetStyle("Heading Background").
			GetStyle("Vertical Center").
			ToNewStyle (f),
	).Width(labelWidth)

	// Actual value for the created on date
	objectiveSheet.NewCell(
		createdOnDateValueColumn,
		createdOnRow,
	).Value(
		objective.GetCreatedOn(),
	).Style(
		models.NewStyle(styles).
			GetStyle("Black Borders").
			GetStyle("Normal Font").
			GetStyle("Normal Background").
			ToNewStyle (f),
	).Width(valueWidth)

	// Label for the amount of time used
	objectiveSheet.NewCell(
		timeUsedDateLabelColumn,
		timeUsedRow,
	).Value(
		"Time Left",
	).Style(
		models.NewStyle(styles).
			GetStyle("Heading Font").
			GetStyle("Heading Background").
			GetStyle("Vertical Center").
			ToNewStyle (f),
	).Width(labelWidth)

	// The amount of time used
	objectiveSheet.NewCell(
		timeUsedDateValueColumn,
		timeUsedRow,
	).Value(
		objective.GetPercentTimeLeft(),
	).Style(
		models.NewStyle(styles).
			GetStyle("Black Borders").
			GetStyle("Normal Font").
			GetStyle("Normal Background").
			ToNewStyle (f),
	).Width(valueWidth)

	// Label for the Update
	objectiveSheet.NewCell(
		updateLabelColumn,
		updateRow,
	).Value(
		"Last Update",
	).Style(
		models.NewStyle(styles).
			GetStyle("Heading Font").
			GetStyle("Heading Background").
			GetStyle("Rotated").
			GetStyle("White Borders").
			ToNewStyle (f),
	).Width(leftColumnWidth)

	// The actual update
	objectiveSheet.NewMergedCell(
		updateValueColumn,
		updateRow,
		farRightColumn,
		updateRow,
	).Value(
		objective.GetUpdate(),
	).Style(
		models.NewStyle(styles).
			GetStyle("Black Borders").
			GetStyle("Normal Font").
			GetStyle("Normal Background").
			GetStyle("Vertical Center").
			ToNewStyle (f),
	).Height(updateHeight)

	// Label for the description
	objectiveSheet.NewCell(
		descriptionLabelColumn,
		descriptionRow,
	).Value(
		"Description",
	).Style(
		models.NewStyle(styles).
			GetStyle("Heading Font").
			GetStyle("Heading Background").
			GetStyle("Rotated").
			GetStyle("White Borders").
			ToNewStyle (f),
	).Width(leftColumnWidth)

	// The actual description
	objectiveSheet.NewMergedCell(
		descriptionValueColumn,
		descriptionRow,
		farRightColumn,
		descriptionRow,
	).Value(
		objective.GetDescription(),
	).Style(
		models.NewStyle(styles).
			GetStyle("Black Borders").
			GetStyle("Normal Font").
			GetStyle("Normal Background").
			GetStyle("Vertical Center").
			ToNewStyle (f),
	).Height(descriptionHeight)

	return objectiveSheet, bottomRow, farRightColumn
}

