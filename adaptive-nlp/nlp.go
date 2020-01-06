package nlp

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/core-utils-go"
	fetch_dialog "github.com/adaptiveteam/dialog-fetcher"
)

// ImprovementID is a enum type
type ImprovementID string

// EQ Analysis Subjects
const (
	// If the text contains 3 verbs
	Actionable   ImprovementID = "actionable"
	// If the text has a sentence that starts with one of the stop-words: but, however, never
	WordsToAvoid ImprovementID = "avoid-words"
	// If the text has meaning cloud sentiment of !NONIRONIC
	Irony        ImprovementID = "irony"

	NotPositive  ImprovementID = "not-positive"
	// If the text has meaning cloud sentiment of SUBJECTIVE
	Subjectivity ImprovementID = "subjective"
	// TooShort if text is shorter than 70 characters
	TooShort     ImprovementID = "too-short"
	// TooLong if text is longer than 700 characters
	TooLong      ImprovementID = "too-long"
	// Dysregulated = !positive & !subjective
	Dysregulated ImprovementID = "dysregulated"
	// GoodDescriptionSubject = none of the above recommendations
	GoodDescriptionSubject ImprovementID = "good"
)


// GetDialog requests text analysis from Meaning Cloud.
// And then converts them to recommendations using dialog.
// deprecated: breaks SRP, the name is misleading. Inline instead.
func GetDialog(
	textToAnalyze string,
	lc LanguageCode,
	dialogFetcherDAO fetch_dialog.DAO,
	context string,
) (dialog map[ImprovementID]string, errList []error) {
	ensureGlobalConnectionsAreOpen()
	recommendations, errList := globalConnections.GetImprovements(textToAnalyze, lc)
	dialog, errList2 := FetchDialogForImprovementsOrGood(dialogFetcherDAO, context, recommendations)
	return  dialog,append(errList, errList2...)
}

// FetchDialogForImprovements fetches dialog messages for the list of recommendations.
// deprecated. Use FetchDialogForImprovementsOrGood
func FetchDialogForImprovements(
	dialogFetcherDAO fetch_dialog.DAO,
	context string,
	recommendations []ImprovementID,
) (dialog map[ImprovementID]string, errList []error) {
	return FetchDialogForImprovementsOrGood(dialogFetcherDAO, context, recommendations)
}

// FetchDialogForImprovementsOrGood fetches dialog messages for the list of recommendations.
func FetchDialogForImprovementsOrGood(
	dialogFetcherDAO fetch_dialog.DAO,
	context string,
	recommendations []ImprovementID,
) (dialog map[ImprovementID]string, errList []error) {
	if len(recommendations) == 0 {
		recommendations = append(recommendations, GoodDescriptionSubject)
	}
	return FetchDialogForListOfImprovements(dialogFetcherDAO, context, recommendations)
}

// FetchDialogForListOfImprovements fetches dialog messages for the list of recommendations.
func FetchDialogForListOfImprovements(
	dialogFetcherDAO fetch_dialog.DAO,
	context string,
	recommendations []ImprovementID,
) (dialog map[ImprovementID]string, errList []error) {
	
 	dialog = make(map[ImprovementID]string, 0)

	for i := 0; i < len(recommendations) && len(errList) == 0; i++ {
		subjectDialog, err := dialogFetcherDAO.FetchByContextSubject(context, string(recommendations[i]))
		dialog, errList = appendToMapAndErrors(dialog, errList, recommendations[i], subjectDialog, err)
	}
	return  dialog,errList
}

// FetchDialogForGood fetches dialog messages for the list of recommendations.
func FetchDialogForGood(
	dialogFetcherDAO fetch_dialog.DAO,
	context string,
) (dialog map[ImprovementID]string, errList []error) {
	dialog = make(map[ImprovementID]string, 0)
	subjectDialog, err := dialogFetcherDAO.FetchByContextSubject(context, string(GoodDescriptionSubject))
	dialog, errList = appendToMapAndErrors(dialog, errList, GoodDescriptionSubject, subjectDialog, err)
	return  dialog,errList
}

func appendToMapAndErrors(dialog map[ImprovementID]string, errList []error, 
	improvement ImprovementID, subjectDialog fetch_dialog.DialogEntry, err error) (dialogOut map[ImprovementID]string, errListOut []error) {
	dialogOut = dialog
	errListOut = errList
	if err == nil && len(subjectDialog.Dialog) > 0 {
		dialogOut[improvement] = core_utils_go.RandomString(subjectDialog.Dialog)
	} else {
		errListOut = append(errList,err)
	}
	return 
}

func wrapDialogError(context string, improvementID ImprovementID, err error) (errOut error) {
	if err != nil {
		errOut = fmt.Errorf(
			"Couldn't find dialog for %s:%s due to %v",
			context,
			improvementID,
			err,
		)
	}
	return
}