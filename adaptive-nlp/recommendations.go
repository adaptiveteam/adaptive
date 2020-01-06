package nlp

import (
	"log"
	"sync"
)

// predicateToOptional converts the given function to another function that will return 
// empty or non empty slice with the given improvement id.
func predicateToOptional(f func(string)bool, improvementID ImprovementID)func(string)[]ImprovementID {
	return func(text string) (res []ImprovementID) {
		found := f(text)
		if found {
			res = []ImprovementID{improvementID}
		}
		return
	}
}

func optionalToGoRoutine(f func(string)[]ImprovementID, text string) func (improvements chan ImprovementID, wg *sync.WaitGroup) {
	return lazyToGoRoutine(func()[]ImprovementID{return f(text)})
	// func (text string, improvementsChannel chan ImprovementID, wg *sync.WaitGroup) {
	// 	defer wg.Done()
	// 	improvements := f(text)
	// 	for _, improvementID := range improvements {
	// 		improvementsChannel <- improvementID
	// 	}
	// }
}

func lazyToGoRoutine(f func()[]ImprovementID) func (improvements chan ImprovementID, wg *sync.WaitGroup) {
	return func (improvementsChannel chan ImprovementID, wg *sync.WaitGroup) {
		defer wg.Done()
		improvements := f()
		for _, improvementID := range improvements {
			improvementsChannel <- improvementID
		}
	}
}

func errPredicateToOptional(f func(string)(bool, error), improvementID ImprovementID)func(string)[]ImprovementID {
	return func(text string) (res []ImprovementID) {
		found, err := f(text)
		if err == nil {
			if found {
				res = []ImprovementID{improvementID}
			}
		} else {
			log.Printf("Error while evaluating improvement %s: %v\n", improvementID, err)
		}
		return
	}
}

// // getWordsToAvoidRecommendations returns a random improvements based on the results of the GetWordsToAvoid function.
// func getWordsToAvoidImprovement(text string, improvements *[]ImprovementID, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	found := GetWordsToAvoid(text)
// 	if found {
// 		*improvements = append(*improvements, WordsToAvoid)
// 	}
// }

// // getTooShortRecommendations returns a random improvements based on the results of the GetTooShort
// func getTooShortImprovement(text string, improvements *[]ImprovementID, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	found := GetTooShort(text)
// 	if found {
// 		*improvements = append(*improvements, TooShort)
// 	}
// }

// // getTooLongRecommendations returns a random improvements based on the results of the GetTooShort
// func getTooLongImprovement(text string, improvements *[]ImprovementID, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	found := GetTooLong(text)
// 	if found {
// 		*improvements = append(*improvements, TooLong)
// 	}
// }

// // getActionableRecommendation returns a random improvements based on the results of the GetActionable function.
// func getActionableImprovement(text string, lc LanguageCode, improvements *[]ImprovementID, wg *sync.WaitGroup) (err error) {
// 	defer wg.Done()
// 	found, err := globalConnections.GetActionable(text, lc)
// 	if err == nil && found {
// 		*improvements = append(*improvements, Actionable)
// 	}
// 	return err
// }

// // GetTextSentimentText function. This function enables recommendations from
// // getIronyRecommendation, getSubjectivityRecommendation, getAgreementRecommendation, and getSentimentRecommendation
// func getSentimentImprovements(text string, lc LanguageCode, improvements *[]ImprovementID, wg *sync.WaitGroup) (err error) {
// 	defer wg.Done()
// 	sentiment, err := GetTextSentimentText(text, lc)
// 	if err == nil && sentiment.GetConfidence() > 80 {
// 		getIronyImprovement(sentiment, improvements, wg)
// 		getSubjectivityImprovement(sentiment, improvements, wg)
// 		getNotPositiveImprovement(sentiment, improvements, wg)
// 	}
// 	return err
// }

// // getIronyRecommendation returns a random improvements based on the results of the GetIrony function
// // from te GetTextSentimentText function.
// func getIronyImprovement(text TextSentiment, improvements *[]ImprovementID, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	found := text.GetIrony()
// 	if found {
// 		*improvements = append(*improvements, Irony)
// 	}
// }

// // getSubjectivityRecommendation returns a random improvements based on the results of the GetSubjectivity function
// // from te GetTextSentimentText function.
// func getSubjectivityImprovement(text TextSentiment, improvements *[]ImprovementID, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	found := text.GetSubjectivity()
// 	if !found {
// 		*improvements = append(*improvements, Subjectivity)
// 	}
// }

// // getSentimentRecommendation returns a random improvements based on the results of the GetSubjectivity function
// // from te GetTextSentimentText function.
// func getNotPositiveImprovement(text TextSentiment, improvements *[]ImprovementID, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	if IsNonPositive(text) {
// 		*improvements = append(*improvements, NotPositive)
// 	}
// }

func getDysregulated(improvements *[]ImprovementID) {
	positive := true
	subjective := true
	for i := 0; i < len(*improvements) && (positive || subjective); i++ {
		if (*improvements)[i] == NotPositive {
			positive = false
		}
		if (*improvements)[i] == Subjectivity {
			subjective = false
		}
	}
	if !positive && !subjective {
		*improvements = append(*improvements, Dysregulated)
	}
}