package ido

import (
excel "github.com/360EntSecGroup-Skylar/excelize/v2"
"github.com/adaptiveteam/adaptive/adaptive-reports/models"
)

func CreateIDOSummary(
	f *excel.File,
	summaryIDO *models.Sheet,
	allIDOs models.IDOs,
	styles map[string]string,
) {
	// Constants
	const idoIndexColumn = 1
	const idoNameColumn = 2
	const updateDateColumn = 3
	const idoCompleted = 4
	const advocateStatusColumn = 5
	const coachStatusColumn = 6
	const coachNameColumn = 7
	const titleRow = 1
	const headingRow = 2
	const startingRow = 3
	const titleHeight = 50
	const coachNameWidth = 20
	const completedWidth = 12
	const indexWidth = 5
	const dateWidth = 12
	const idoNameWidth = 40
	const statusWidth = 12
	const farRight = 7
	const farLeft = 1

	summarySheet := summaryIDO
	// Create Title
	summarySheet.NewMergedCell(
		farLeft, titleRow,
		farRight, titleRow,
	).Value(allIDOs[0].Advocate() +" IDO Summary").
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

	// Create Objectives index column
	summarySheet.NewCell(
		coachNameColumn, headingRow,
	).Value("Coach").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(coachNameWidth)

	// Create Objectives index column
	summarySheet.NewCell(
		idoCompleted, headingRow,
	).Value("State").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(completedWidth)

	// Create Objectives status column
	summarySheet.NewCell(
		advocateStatusColumn, headingRow,
	).Value("Last Advocate Status").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(statusWidth)

	// Create Objectives status column
	summarySheet.NewCell(
		coachStatusColumn, headingRow,
	).Value("Last Coach's Status").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(statusWidth)

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

		// Create Initiatives index column
	summarySheet.NewCell(
		updateDateColumn, headingRow,
	).Value("Last Updated").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(dateWidth)

	// Now post all of the initiatives and objectives
	for idoIndex := 0; idoIndex < len(allIDOs); idoIndex++ {
		// Post the updates
		currentIDO := allIDOs[idoIndex]

		// Post the objective index and merge the rows
		summarySheet.NewCell(
			idoIndexColumn, idoIndex+startingRow,
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
			idoNameColumn, idoIndex+startingRow,
		).Value(currentIDO.Name()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Vertical Center").
					GetStyle("Link Font").
					GetStyle("Normal Background").
					ToNewStyle(f),
			).SetLink(
				"IDO Details",
				"B",
				idoIndex+startingRow,
			)

		summarySheet.NewCell(
			coachNameColumn, idoIndex+startingRow,
		).Value(currentIDO.Coach()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Centered").
					GetStyle("Normal Font").
					GetStyle("Normal Background").
					ToNewStyle(f),
			)

		completedString := completedString(currentIDO.Completed())
		summarySheet.NewCell(
			idoCompleted, idoIndex+startingRow,
		).Value(completedString).
			Style(
				models.NewStyle(styles).
					GetStyle(completedString+" Background").
					GetStyle("Black Borders").
					GetStyle("Centered").
					GetStyle(completedString + " Font").
					ToNewStyle(f),
			)

		summarySheet.NewCell(
			advocateStatusColumn, idoIndex+startingRow,
		).Value(currentIDO.Updates()[0].AdvocateStatus()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Centered").
					GetStyle("Italics Font").
					GetStyle("Normal Background").
					GetStyle(
						getStatusStyle(
							currentIDO.Completed(),
							currentIDO.Updates()[0].AdvocateStatus(),
						),
					).ToNewStyle(f),
			)

		summarySheet.NewCell(
			coachStatusColumn, idoIndex+startingRow,
		).Value(currentIDO.Updates()[0].CoachStatus()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Centered").
					GetStyle("Italics Font").
					GetStyle("Normal Background").
					GetStyle(
						getStatusStyle(
							currentIDO.Completed(),
							currentIDO.Updates()[0].CoachStatus(),
						),
					).ToNewStyle(f),
			)

		summarySheet.NewCell(
			updateDateColumn, idoIndex+startingRow,
		).Value(currentIDO.Updates()[0].UpdateDate()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Centered").
					GetStyle("Link Font").
					GetStyle("Normal Background").
					ToNewStyle(f),
			).SetLink(
					completedString+" IDO Updates",
					"B",
					getUpdateRow(currentIDO,allIDOs)+startingRow,
			)
	}
	summarySheet.AddFilter(idoNameColumn, headingRow, coachNameColumn, headingRow)
}

