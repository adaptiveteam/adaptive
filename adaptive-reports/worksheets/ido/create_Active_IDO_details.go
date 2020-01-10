package ido

import (
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive/adaptive-reports/models"
)

func CreateActiveIDOSummary(
	f *excel.File,
	activeIDOSummary *models.Sheet,
	allIDOs models.IDOs,
	styles map[string]string,
) {
	// Constants
	const idoIndexColumn = 1
	const idoNameColumn = 2
	const advocateStatusColumn = 3
	const coachStatusColumn = 4
	const updateDateColumn = 5
	const advocateUpdateColumn = 6
	const coachUpdateColumn = 7
	const titleRow = 1
	const titleHeight = 50
	const startingRow = 2
	const indexWidth = 5
	const dateWidth = 12
	const idoNameWidth = 30
	const statusWidth = 12
	const updateWidth = 60
	const farRight = 7
	const farLeft = 1

	summarySheet := activeIDOSummary
	// Create Title
	summarySheet.NewMergedCell(
		farLeft, titleRow,
		farRight, titleRow,
	).Value("IDO Performance").
		Style(
			models.NewStyle(styles).
				GetStyle("Title Font").
				GetStyle("Heading Background").
				GetStyle("Centered").
				GetStyle("White Bottom Border").
				ToNewStyle(f),
		).Height(titleHeight)

	// Create Objectives index column
	summarySheet.NewCell(
		idoIndexColumn, startingRow,
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
	summarySheet.NewCell(
		advocateStatusColumn, startingRow,
	).Value("Advocate Status").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		)

	// Create Objectives status column
	summarySheet.NewCell(
		coachStatusColumn, startingRow,
	).Value("Coach's Status").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		)

	// Create Objectives column
	summarySheet.NewCell(
		idoNameColumn, startingRow,
	).Value("IDO").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(idoNameWidth)

	// Create Initiatives status column
	summarySheet.NewCell(
		advocateStatusColumn, startingRow,
	).Value("Your Status").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(statusWidth)

	// Create Initiatives status column
	summarySheet.NewCell(
		coachStatusColumn, startingRow,
	).Value("Coach's Status").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(statusWidth)

	// Create Initiatives index column
	summarySheet.NewCell(
		updateDateColumn, startingRow,
	).Value("Updated On").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(dateWidth)

	// Create Initiatives column
	summarySheet.NewCell(
		advocateUpdateColumn, startingRow,
	).Value("Your Updates").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(updateWidth)

	summarySheet.NewCell(
		coachUpdateColumn, startingRow,
	).Value("Coach's Updates").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(updateWidth)

	// Now post all of the initiatives and objectives
	updateStart := startingRow + 1
	updateEnd := startingRow + 1

	for idoIndex := 0; idoIndex < len(allIDOs); idoIndex++ {
		// Post the updates
		currentIDO := allIDOs[idoIndex]
		for updateIndex := 0; updateIndex < len(currentIDO.Updates()); updateIndex++ {
			currentUpdate := &currentIDO.Updates()[updateIndex]

			// Post the objective index and merge the rows
			summarySheet.NewCell(
				updateDateColumn,
				updateEnd,
			).Value((*currentUpdate).UpdateDate()).
				Style(
					models.NewStyle(styles).
						GetStyle("Index Font").
						GetStyle("Heading Background").
						GetStyle("White Borders").
						GetStyle("Centered").
						ToNewStyle(f),
				)

			summarySheet.NewCell(
				advocateStatusColumn,
				updateEnd,
			).Value((*currentUpdate).AdvocateStatus()).
				Style(
					models.NewStyle(styles).
						GetStyle((*currentUpdate).AdvocateStatus()).
						GetStyle("Black Borders").
						GetStyle("Italics Font").
						GetStyle("Vertical Center").
						ToNewStyle(f),
				)

			summarySheet.NewCell(
				coachStatusColumn,
				updateEnd,
			).Value((*currentUpdate).CoachStatus()).
				Style(
					models.NewStyle(styles).
						GetStyle((*currentUpdate).CoachStatus()).
						GetStyle("Black Borders").
						GetStyle("Italics Font").
						GetStyle("Vertical Center").
						ToNewStyle(f),
				)


			summarySheet.NewCell(
				advocateUpdateColumn,
				updateEnd,
			).Value((*currentUpdate).AdvocateComments()).
				Style(
					models.NewStyle(styles).
						GetStyle("Black Borders").
						GetStyle("Normal Font").
						GetStyle("Normal Background").
						GetStyle("Vertical Center").
						ToNewStyle(f),
			)

			summarySheet.NewCell(
				coachUpdateColumn,
				updateEnd,
			).Value((*currentUpdate).CoachComments()).
				Style(
					models.NewStyle(styles).
						GetStyle("Black Borders").
						GetStyle("Normal Font").
						GetStyle("Normal Background").
						GetStyle("Vertical Center").
						ToNewStyle(f),
				)

			updateEnd++
		}

		// Post the objective index and merge the rows
		summarySheet.NewMergedCell(
			idoIndexColumn, updateStart,
			idoIndexColumn, updateEnd-1,
		).IntValue(idoIndex+1).
			Style(
				models.NewStyle(styles).
					GetStyle("Index Font").
					GetStyle("Heading Background").
					GetStyle("White Borders").
					GetStyle("Centered").
					ToNewStyle(f),
			).Width(indexWidth)

		summarySheet.NewMergedCell(
			idoNameColumn, updateStart,
			idoNameColumn, updateEnd-1,
		).Value(currentIDO.Name()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Vertical Center").
					GetStyle("Normal Font").
					GetStyle("Normal Background").
					ToNewStyle(f),
			)

		updateStart = updateEnd

	}
}
