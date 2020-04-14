package collaboration_report

import (
	"log"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	fetch_dialog "github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/unidoc/unipdf/v3/creator"
)

// BuildReportWithCustomValuesTyped is an entry point for this project.
func BuildReportWithCustomValuesTyped(
	// The last year of feedback received
	received CoachingList,
	// The users name (e.g., Chris Creel)
	UserName string,
	// The quarter for which this report was produced
	Quarter int,
	// The year for which this report was produced
	Year int,
	// Name and location for where to store the file.
	FileName string,
	dialogDao fetch_dialog.DAO,
	logger logger.AdaptiveLogger,
) (tags map[string]string, err error) {
	SetUniDocGlobalLicenseIfAvailable()
	var pdf *creator.Creator
	pdf, tags, err = createPdfReport(received, //given,
		UserName,
		Quarter,
		Year,
		dialogDao,
		logger,
	)
	if err == nil {
		err = pdf.WriteToFile(FileName)
	}
	if err != nil {
		log.Println("Error building report "+FileName, err)
	}
	return tags, err
}
