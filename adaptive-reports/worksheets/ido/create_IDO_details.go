package ido

import (
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive-reports/models"
)

func CreateIDODetails(
	f *excel.File,
	idoDetails *models.Sheet,
	allIDOs models.IDOs,
	styles map[string]string,
) {
	// Constants
	const idoIndexColumn = 1
	const idoNameColumn = 2
	const updateDateColumn = 3
	const coachNameColumn = 4
	const idoCompleted = 5
	const advocateStatusColumn = 6
	const coachStatusColumn = 7
	const alignedWith = 8
	const isA = 9
	const driving = 10
	const idoDescription = 11
	const alignmentDescription = 12
	const titleRow = 1
	const titleHeight = 50
	const headingRow = 2
	const startingRow = 3
	const indexWidth = 5
	const dateWidth = 12
	const coachNameWidth = 20
	const completedWidth = 12
	const idoNameWidth = 30
	const idoDescriptionWidth = 80
	const alignmentDescriptionWidth = 80
	const statusWidth = 12
	const updateWidth = 60
	const alignedWidth = 30
	const whichIsAWidth = 15
	const drivingWidth = 30
	const farRight = 12
	const farLeft = 1

	idoDetails.NewMergedCell(
		farLeft, titleRow,
		farRight, titleRow,
	).Value("IDO Details").
		Style(
			models.NewStyle(styles).
				GetStyle("Title Font").
				GetStyle("Heading Background").
				GetStyle("Centered").
				GetStyle("White Bottom Border").
				ToNewStyle(f),
		).Height(titleHeight)

	// Create Objectives index column
	idoDetails.NewCell(
		idoIndexColumn,
		headingRow,
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
	idoDetails.NewCell(
		advocateStatusColumn,
		headingRow,
	).Value("Advocate Status").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(statusWidth)

	// Create Objectives status column
	idoDetails.NewCell(
		coachStatusColumn,
		headingRow,
	).Value("Coach's Status").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(statusWidth)

	// Create Objectives column
	idoDetails.NewCell(
		idoNameColumn,
		headingRow,
	).Value("IDO").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(idoNameWidth)

	// Create Objectives index column
	idoDetails.NewCell(
		coachNameColumn,
		headingRow,
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
	idoDetails.NewCell(
		idoCompleted,
		headingRow,
	).Value("State").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(completedWidth)

	// Create Initiatives status column
	idoDetails.NewCell(
		advocateStatusColumn,
		headingRow,
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
	idoDetails.NewCell(
		coachStatusColumn,
		headingRow,
	).Value("Coach's Status").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(statusWidth)

	// Create Initiatives status column
	idoDetails.NewCell(
		updateDateColumn,
		headingRow,
	).Value("Last Updated").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(dateWidth)

	// Create Initiatives index column
	idoDetails.NewCell(
		alignedWith,
		headingRow,
	).Value("aligned with").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(alignedWidth)

	// Create Initiatives index column
	idoDetails.NewCell(
		isA,
		headingRow,
	).Value("which is a(n)").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(whichIsAWidth)

	// Create Initiatives index column
	idoDetails.NewCell(
		driving,
		headingRow,
	).Value("that drives").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(drivingWidth)

	// Create Initiatives column
	idoDetails.NewCell(
		idoDescription,
		headingRow,
	).Value("IDO Description").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(updateWidth)

	// Create Initiatives column
	idoDetails.NewCell(
		alignmentDescription,
		headingRow,
	).Value("Strategic Alignment Description").
		Style(
			models.NewStyle(styles).
				GetStyle("Heading Font").
				GetStyle("Heading Background").
				GetStyle("Horizontal Center").
				GetStyle("White Borders").
				ToNewStyle(f),
		).Width(alignmentDescriptionWidth)

	// Now post all of the initiatives and objectives
	updateEnd := startingRow

	for idoIndex := 0; idoIndex < len(allIDOs); idoIndex++ {
		// Post the updates
		currentIDO := allIDOs[idoIndex]
		// Post the objective index and merge the rows
		idoDetails.NewCell(
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

		idoDetails.NewCell(
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

		idoDetails.NewCell(
			coachNameColumn, updateEnd,
		).Value(currentIDO.Coach()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Vertical Center").
					GetStyle("Normal Font").
					GetStyle("Normal Background").
					ToNewStyle(f),
			)

		completedString := completedString(currentIDO.Completed())
		idoDetails.NewCell(
			idoCompleted, updateEnd,
		).Value(completedString).
			Style(
				models.NewStyle(styles).
					GetStyle(completedString+" Background").
					GetStyle("Black Borders").
					GetStyle("Centered").
					GetStyle(completedString + " Font").
					ToNewStyle(f),
			)

		idoDetails.NewCell(
			advocateStatusColumn, updateEnd,
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

		idoDetails.NewCell(
			coachStatusColumn, updateEnd,
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

		idoDetails.NewCell(
			updateDateColumn, updateEnd,
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

		idoDetails.NewCell(
			alignedWith, updateEnd,
		).Value(currentIDO.FocusedOnName()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Vertical Center").
					GetStyle("Normal Font").
					GetStyle("Normal Background").
					ToNewStyle(f),
			)

		idoDetails.NewCell(
			isA, updateEnd,
		).Value(currentIDO.IsA()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Centered").
					GetStyle("Normal Font").
					GetStyle("Normal Background").
					ToNewStyle(f),
			)

		idoDetails.NewCell(
			driving, updateEnd,
		).Value(currentIDO.Driving()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Vertical Center").
					GetStyle("Normal Font").
					GetStyle("Normal Background").
					ToNewStyle(f),
			)

		idoDetails.NewCell(
			idoDescription, updateEnd,
		).Value(currentIDO.Description()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Vertical Center").
					GetStyle("Normal Font").
					GetStyle("Normal Background").
					ToNewStyle(f),
			).Width(idoDescriptionWidth)

		idoDetails.NewCell(
			alignmentDescription, updateEnd,
		).Value(currentIDO.FocusedOnDescription()).
			Style(
				models.NewStyle(styles).
					GetStyle("Black Borders").
					GetStyle("Vertical Center").
					GetStyle("Normal Font").
					GetStyle("Normal Background").
					ToNewStyle(f),
			).Width(alignmentDescriptionWidth)

		updateEnd++
	}
	idoDetails.AddFilter(idoNameColumn, headingRow, driving, headingRow)
}
