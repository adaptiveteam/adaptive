package strategy

import (
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive-reports/models"
	"sort"
	"strconv"
)

func CreateStrategySummary(
	f *excel.File,
	strategicSummary *models.Sheet,
	vision string,
	allObjectives *models.Objectives,
	styles map[string]string,
) {
	// Constants
	const objectiveIndexColumn = 1
	const objectiveNameColumn = 2
	const objectiveStatusColumn = 3
	const initiativeIndexColumn = 4
	const initiativeNameColumn = 5
	const initiativeStatusColumn = 6
	const objectiveTypeColumnLabel = 1
	const objectiveTypeColumnValue = 2
	const titleRow = 1
	const titleHeight = 50
	const visionRow = 2
	const visionHeight = 75
	const startingRow = 3
	const labelWidth = 12
	const indexWidth = 8
	const nameWidth = 37
	const farRightColumn = 6

	// Create Title
	strategicSummary.NewMergedCell(
		objectiveIndexColumn, titleRow,
		initiativeStatusColumn, titleRow,
	).Value("Strategy Performance").
		Style(
			models.NewStyle(styles).
				GetStyle("Title Font").
				GetStyle("Heading Background").
				GetStyle("Centered").
				GetStyle("White Bottom Border").
				ToNewStyle(f),
		).Height(titleHeight)

	// Create Vision
	strategicSummary.NewMergedCell(
		objectiveIndexColumn, visionRow,
		farRightColumn, visionRow,
	).Value(vision).
		Style(
			models.NewStyle(styles).
				GetStyle("Black Borders").
				GetStyle("Normal Font").
				GetStyle("Vertical Center").
				ToNewStyle (f),
		).Height(visionHeight)

	// Create Objectives index column
	strategicSummary.NewCell(
		objectiveIndexColumn, startingRow,
	).Value("#").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(indexWidth)

	// Create Objectives status column
	strategicSummary.NewCell(
		objectiveStatusColumn, startingRow,
	).Value("Status").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		)

	// Create Objectives column
	strategicSummary.NewCell(
		objectiveNameColumn, startingRow,
	).Value("Objectives").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		)

	// Create Initiatives status column
	strategicSummary.NewCell(
		initiativeStatusColumn, startingRow,
	).Value("Status").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(labelWidth)

	// Create Initiatives column
	strategicSummary.NewCell(
		initiativeNameColumn, startingRow,
	).Value("Initiatives").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		)

	// Create Initiatives index column
	strategicSummary.NewCell(
		initiativeIndexColumn, startingRow,
	).Value("#").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(indexWidth)

	// Now post all of the initiatives and objectives
	initiativeStart := startingRow + 1
	initiativeEnd := startingRow + 1
	initiativeIndex := 1

	// we want the objectives in sorted order
	keys := make([]string, 0, len(*allObjectives))
	for k := range *allObjectives {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	objectiveIndex := 0
	for ; objectiveIndex < len(*allObjectives); objectiveIndex++ {
		// Post the initiatives
		currentObjective := (*allObjectives)[keys[objectiveIndex]]
		currentObjective.SetIndex(objectiveIndex+1)
		for i := 0; i < len(currentObjective.Initiatives); i++ {
			currentInitiative := &currentObjective.Initiatives[i]
			(*currentInitiative).SetIndex(initiativeIndex)
			strategicSummary.NewCell(
				initiativeStatusColumn,
				initiativeEnd,
			).Value((*currentInitiative).GetStatus()).
				Style(
					models.NewStyle(styles).
						GetStyle((*currentInitiative).GetStatus()).
						GetStyle("Black Borders").
						GetStyle("Italics Font").
						GetStyle("Vertical Center").
						ToNewStyle(f),
				).Width(labelWidth)

			// Post the objective index and merge the rows
			strategicSummary.NewCell(
				initiativeIndexColumn,
				initiativeEnd,
			).IntValue((*currentInitiative).GetIndex()).
				Style(
					models.NewStyle(styles).
						GetStyle("Index Font").
						GetStyle("Heading Background").
						GetStyle("White Borders").
						GetStyle("Centered").
						ToNewStyle(f),
				).Width(indexWidth)

			strategicSummary.NewCell(
				initiativeNameColumn,
				initiativeEnd,
			).Value((*currentInitiative).GetName()).
				Width(nameWidth).
				Style(
					models.NewStyle(styles).
						GetStyle("Black Borders").
						GetStyle("Link Font").
						GetStyle("Normal Background").
						ToNewStyle(f),
				).SetLink(
				"Initiative "+strconv.Itoa((*currentInitiative).GetIndex()),
				"A",
				1,
			)

			if (*currentInitiative).InitiativeName != "No Initiatives" {
				CreateComponentDetails(
					"Initiative",
					f,
					styles,
					*currentInitiative,
				)
			}

			initiativeEnd++
			initiativeIndex++
		}

		// Post the objective index and merge the rows
		strategicSummary.NewMergedCell(
			objectiveIndexColumn, initiativeStart,
			objectiveIndexColumn, initiativeEnd-1,
		).IntValue(currentObjective.GetIndex()).
			Style(
				models.NewStyle(styles).
					GetStyle("Index Font").
					GetStyle("Heading Background").
					GetStyle("White Borders").
					GetStyle("Centered").
					ToNewStyle(f),
			).Width(indexWidth)

		// Post the objectives and merge the rows
		strategicSummary.NewMergedCell(
			objectiveStatusColumn, initiativeStart,
			objectiveStatusColumn, initiativeEnd-1,
		).Value(currentObjective.GetStatus()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Centered").
					GetStyle("Italics Font").
					GetStyle("Normal Background").
					GetStyle(currentObjective.ObjectiveStatus).
					ToNewStyle(f),
			).Width(labelWidth)

		strategicSummary.NewMergedCell(
			objectiveNameColumn, initiativeStart,
			objectiveNameColumn, initiativeEnd-1,
		).Value(currentObjective.GetName()).
			Width(nameWidth).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Vertical Center").
					GetStyle("Link Font").
					GetStyle("Normal Background").
					ToNewStyle(f),
			).SetLink(
			"Objective "+strconv.Itoa(currentObjective.GetIndex()),
			"A",
			1,
		)
		currentObjective.SetIndex(currentObjective.GetIndex())
		initiativeStart = initiativeEnd

		objectiveSheet, bottomRow, farRightColumn := CreateComponentDetails(
			"Objective",
			f,
			styles,
			currentObjective,
		)

		// Now we need to decorate this sheet with Objective specific data
		// starting with the objective type label.
		objectiveSheet.NewCell(
			objectiveTypeColumnLabel,
			bottomRow+1,
		).Value(
			"Objective Type",
		).Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Vertical Center").
				ToNewStyle (f),
		)

		// Now add the actual objective type
		objectiveSheet.NewMergedCell(
			objectiveTypeColumnValue,
			bottomRow+1,
			farRightColumn,
			bottomRow+1,
		).Value(
			currentObjective.ObjectiveType,
		).Style(
			models.NewStyle(styles).
				GetStyle("Black Borders").
				GetStyle("Normal Font").
				GetStyle("Normal Background").
				ToNewStyle (f),
		)
	}
	strategicSummary.AddFilter(1, startingRow, farRightColumn, startingRow)
}