package ido

import (
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive-reports/models"
)

func CreateIDOUpdates(
	f *excel.File,
	activeIDOSummary *models.Sheet,
	allIDOs models.IDOs,
	completed bool,
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
	const headingRow = 2
	const startingRow = 3
	const indexWidth = 5
	const dateWidth = 12
	const idoNameWidth = 30
	const statusWidth = 12
	const updateWidth = 60
	const farRight = 7
	const farLeft = 1

	summarySheet := activeIDOSummary

	summarySheet.NewMergedCell(
		farLeft, titleRow,
		farRight, titleRow,
	).Value(completedString(completed) + " IDO Updates").
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
		idoIndexColumn, headingRow,
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
		advocateStatusColumn, headingRow,
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
		coachStatusColumn, headingRow,
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
		idoNameColumn, headingRow,
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
		advocateStatusColumn, headingRow,
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
		coachStatusColumn, headingRow,
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
		updateDateColumn, headingRow,
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
		advocateUpdateColumn, headingRow,
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
		coachUpdateColumn, headingRow,
	).Value("Coach's Response").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(updateWidth)

	// Now post all of the initiatives and objectives
	updateEnd := startingRow

	for idoIndex := 0; idoIndex < len(allIDOs); idoIndex++ {
		// Post the updates
		currentIDO := allIDOs[idoIndex]
		if currentIDO.Completed() == completed {
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

				// Post the objective index and merge the rows
				summarySheet.NewCell(
					idoIndexColumn, updateEnd,
				).IntValue(idoIndex+1).
					Style(
						models.NewStyle(styles).
							GetStyle("Index Font").
							GetStyle("Heading Background").
							GetStyle("White Borders").
							GetStyle("Centered").
							ToNewStyle(f),
					).Width(indexWidth)

				summarySheet.NewCell(
					idoNameColumn, updateEnd,
				).Value(currentIDO.Name()).
					Style(
						models.NewStyle(styles).
							GetStyle("Black Borders").
							GetStyle("Vertical Center").
							GetStyle("Link Font").
							GetStyle("Normal Background").
							ToNewStyle(f),
					).SetLink(
						"IDO Summary",
						"B",
						idoIndex+startingRow,
					)
				updateEnd++
			}
		}
	}
}
