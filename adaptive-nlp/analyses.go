package nlp

import (
	"time"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
)
// Deprecated: Uses global connections. Instead, connect in lambda and use those connections.
func GetImprovements(text string, lc LanguageCode) (improvements []ImprovementID, errList []error) {
	ensureGlobalConnectionsAreOpen()
	return globalConnections.GetImprovements(text, lc)
}

func (c Connections)actionableGo(text string, lc LanguageCode, improvementsChannel chan ImprovementID, errorsChannel chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	pos, err := c.GetPartOfSpeech(text, lc)
	if err == nil {
		if IsActionable(pos) {
			improvementsChannel <- Actionable
		}
	} else {
		errorsChannel <- err
	}
}

func (c Connections)sentimentGo(text string, lc LanguageCode, improvementsChannel chan ImprovementID, errorsChannel chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	sentiment, err := GetTextSentimentText(text, lc)
	if err == nil {
		if sentiment.GetConfidence() > 80 {
			if sentiment.GetIrony() { improvementsChannel <- Irony }
			if !sentiment.GetSubjectivity() { improvementsChannel <- Subjectivity }
			if IsNonPositive(sentiment) { improvementsChannel <- NotPositive }
		}
	} else {
		errorsChannel <- err
	}
}

var defaultTimeout = 5000*time.Millisecond

func waitWithTimeout(wg *sync.WaitGroup, dur time.Duration) {
	readiness := make(chan struct{}, 2)
	go func(){
		wg.Wait()
		readiness <- struct{}{}
	}()
	go func(){
		timer := time.NewTimer(dur)
		<- timer.C
		readiness <- struct{}{}
	}()
	<- readiness
}

func (c Connections)GetImprovements(text string, lc LanguageCode) ([]ImprovementID, []error) {
	ensureGlobalConnectionsAreOpen()
	// improvements = make([]ImprovementID, 0)
	const checksCount = 5
	improvementsChannel := make(chan ImprovementID, checksCount)
	errorsChannel := make(chan error, checksCount)
	tooShort := optionalToGoRoutine(predicateToOptional(GetTooShort, TooShort), text)
	tooLong := optionalToGoRoutine(predicateToOptional(GetTooLong, TooLong), text)
	wordsToAvoid := optionalToGoRoutine(predicateToOptional(GetWordsToAvoid, WordsToAvoid), text)
	if len(text) > 0 {
		wg := &sync.WaitGroup{}
		// There are totally checksCount goroutines from below
		wg.Add(checksCount)
		go tooShort(improvementsChannel, wg)
		go tooLong(improvementsChannel, wg)
		go wordsToAvoid(improvementsChannel, wg)
		go c.actionableGo(text, lc, improvementsChannel, errorsChannel, wg)
		go c.sentimentGo(text, lc, improvementsChannel, errorsChannel, wg)
		waitWithTimeout(wg, defaultTimeout)
	}
	close(errorsChannel)
	close(improvementsChannel)
	improvements := collectImprovementIDs(improvementsChannel)
	getDysregulated(&improvements)
	fmt.Println(improvements)
	return improvements, collectErrors(errorsChannel)
}

func collectErrors(ec chan error) (errors []error) {
	for e := range ec {
		errors = append(errors, e)
	}
	return
}

func collectImprovementIDs(improvementsChannel chan ImprovementID) (improvements []ImprovementID) {
	for improvement := range improvementsChannel {
		improvements = append(improvements, improvement)
	}
	return
}

// TextSentiment is an interface for working with sentiment data like
// sentiment, confidence, agreement, subjectivity, and irony.
type TextSentiment interface {
	GetSentiment() string
	GetConfidence() int
	GetAgreement() bool
	GetSubjectivity() bool
	GetIrony() bool
}

// GetSentiment returns the overall sentiment of the content.
func (s textSentiment) GetSentiment() string { return s.Sentiment }

// GetConfidence returns an int between 0 and 100 representing the confidence in the sentiment rating.
func (s textSentiment) GetConfidence() int { return s.Confidence }

// GetAgreement determines if all of the content has the same sentiment.
func (s textSentiment) GetAgreement() bool { return s.Agreement }

// GetSubjectivity determines if the content is subjective in nature versus objective.
func (s textSentiment) GetSubjectivity() bool { return s.Subjectivity }

// GetIrony determines if the content contains any irony.
func (s textSentiment) GetIrony() bool { return s.Irony }

func IsNonPositive(text TextSentiment) bool {
	sentiment := text.GetSentiment()
	return sentiment == "negative" || sentiment == "very negative"
}

// TextCategories is an interface for working with deep categorizations in text.
type TextCategories interface {
	GetLabel() string
	GetRelevance() int
	GetAbsRelevance() int
	GetSentiment() string
}

// GetLabel returns the label of a category found in the text.
func (c textCategory) GetLabel() string { return c.Label }

// GetAbsRelevance returns the absolute relevance of the category.
func (c textCategory) GetAbsRelevance() int { return c.AbsRelevance }

// GetRelevance returns the relative relevance value of the category, a number in the 0-100% range.
// It's computed with respect to the top ranked result (for generic models) and with respect to the
// top ranked result in the same dimension (for dimension models).
func (c textCategory) GetRelevance() int { return c.Relevance }

// Returns the sentiment of the category relative to the text.
func (c textCategory) GetSentiment() string { return c.Sentiment }

// GetTextCategoriesURL returns deep categories found in the web page found at the provided URL.
// The language code, LanguageCode, must also be provided.
func GetTextCategoriesURL(url *url.URL, l LanguageCode) ([]TextCategories, error) {
	ensureGlobalConnectionsAreOpen()
	return globalConnections.MeaningCloud.GetTextCategories("url", url.String(), l)
}

// GetTextCategoriesText returns deep categories found in the provided text.
// The max length of the text must be less than 500 characters.
// The language code, LanguageCode, must also be provided.
func GetTextCategoriesText(text string, l LanguageCode) ([]TextCategories, error) {
	ensureGlobalConnectionsAreOpen()
	return globalConnections.MeaningCloud.GetTextCategories("txt", text, l)
}

// GetTextSentimentURL returns sentiment data found in the web page found at the provided URL.
// The language code, LanguageCode, must also be provided.
func GetTextSentimentURL(url *url.URL, l LanguageCode) (TextSentiment, error) {
	ensureGlobalConnectionsAreOpen()
	return globalConnections.MeaningCloud.GetTextSentiment("url", url.String(), l)
}

// GetTextSentimentText returns sentiment data found in the provided text.
// The max length of the text must be less than 500 characters.
// The language code, LanguageCode, must also be provided.
func GetTextSentimentText(text string, l LanguageCode) (TextSentiment, error) {
	ensureGlobalConnectionsAreOpen()
	return globalConnections.MeaningCloud.GetTextSentiment("txt", text, l)
}

// GetSummaryURL provides a summary of numSentences long of the content
// found in the web page found at the provided URL.
// The language code, LanguageCode, must also be provided.
func GetSummaryURL(numSentences int, url *url.URL, l LanguageCode) (string, error) {
	ensureGlobalConnectionsAreOpen()
	return globalConnections.MeaningCloud.GetSummary(numSentences, "url", url.String(), l)
}

// GetSummaryText provides a summary of numSentences long of the content
// found in the provided text. The max length of the text must be less than 500 characters.
// The language code, LanguageCode, must also be provided.
func GetSummaryText(numSentences int, text string, l LanguageCode) (string, error) {
	ensureGlobalConnectionsAreOpen()
	return globalConnections.MeaningCloud.GetSummary(numSentences, "txt", text, l)
}
// IsActionable determines if the content contains action oriented language
func IsActionable(pos map[string]int)bool {
	return !(pos["VERB"] > 3 && (pos["ADV"] > 0 || pos["ADJ"] > 0))
}
// IsPersonal determines if the content contains pronouns
func IsPersonal(pos map[string]int)bool {
	return !(pos["PROPN"] > 0)
}
// // GetActionable determines if the content contains action oriented language
// func (c Connections)GetActionable(pos map[string]int) (bool, error) {
// //	pos, err := c.GetPartOfSpeech(text, lc)
// 	result := true
// 	//if err == nil {
// 		IsActionable(pos)
// 	//}
// 	return result, err
// }

// // GetPersonal determines if the content contains pronouns
// func (c Connections)GetPersonal(text string, lc LanguageCode) (bool, error) {
// 	pos, err := c.GetPartOfSpeech(text, lc)
// 	result := true
// 	if err == nil {
// 		if pos["PROPN"] > 0 {
// 			result = false
// 		}
// 	}
// 	return result, err
// }

// GetTooShort determines if the content is too short
func GetTooShort(text string) bool {
	// The length of most effective line content is  70 characters
	return len(text) <= 70
}

// GetTooShort determines if the content is too long
func GetTooLong(text string) bool {
	return len(text) >= 700
}

// GetPartOfSpeech returns a map of a simplified part of speech analysis.
// Each part of speech in the map contains the count of that part in the provided text.
// The max length of the text must be less than 5000 characters.
// The language code, LanguageCode, must also be provided.
func (c Connections) GetPartOfSpeech(text string, lc LanguageCode) (map[string]int, error) {
	ts, err := c.Comprehend.DetectSyntax(*lc.String(), text)

	posMap := map[string]int{
		"ADJ":   0,
		"ADP":   0,
		"ADV":   0,
		"AUX":   0,
		"CCONJ": 0,
		"DET":   0,
		"INTJ":  0,
		"NOUN":  0,
		"NUM":   0,
		"O":     0,
		"PART":  0,
		"PRON":  0,
		"PROPN": 0,
		"PUNCT": 0,
		"SCONJ": 0,
		"SYM":   0,
		"VERB":  0,
	}

	if err == nil {
		for i := range ts.SyntaxTokens {
			pos := ts.SyntaxTokens[i].PartOfSpeech.Tag
			posCount := posMap[*pos]
			posMap[*pos] = posCount + 1
		}
	}
	return posMap, err
}

// GetWordsToAvoid looks for a list of words that you should typically not used at the start
// of sentences in text.
func GetWordsToAvoid(text string) (found bool) {
	wordsToAvoid := []string{
		"no",
		"but",
		"however",
	}

	var re = regexp.MustCompile(`(?m)(?:^|(?:[.!?]\s))(\w+)`)
	matches := re.FindAllString(strings.ToLower(text), -1)
	var cleanedMatches []string
	for m := range matches {
		cleanedMatches = append(cleanedMatches, strings.Trim(matches[m], ". \t"))
	}
	for i := 0; i < len(cleanedMatches) && !found; i++ {
		for j := 0; j < len(wordsToAvoid) && !found; j++ {
			found = wordsToAvoid[j] == cleanedMatches[i]
		}
	}
	return found
}

// String() returns the string representation of the enumeration
func (lc LanguageCode) String() *string {
	lcs := string(lc)
	return &lcs
}

// GetTranslation returns a translation of the provided text from the source language to the target language.
// The max length of the text must be less than 5000 characters.
func (c Connections) GetTranslation(text string, sourceLanguage LanguageCode, targetLanguage LanguageCode) (string, error) {
	output, err := c.Translate.TranslateText(text, *sourceLanguage.String(), *targetLanguage.String())

	if err != nil {
		log.Println("[Error] failed to aws translate translation message: " + text, err)
		return "", err
	}

	return *output.TranslatedText, nil
}

type LanguageCode string

// Platform names
const (
	Arabic     LanguageCode = "ar"
	Chinese    LanguageCode = "zh"
	Czech      LanguageCode = "cs"
	Danish     LanguageCode = "da"
	Dutch      LanguageCode = "nl"
	English    LanguageCode = "en"
	Finnish    LanguageCode = "fi"
	French     LanguageCode = "fr"
	German     LanguageCode = "de"
	Hebrew     LanguageCode = "he"
	Indonesian LanguageCode = "id"
	Italian    LanguageCode = "it"
	Japanese   LanguageCode = "ja"
	Korean     LanguageCode = "ko"
	Polish     LanguageCode = "pl"
	Portuguese LanguageCode = "pt"
	Russian    LanguageCode = "ru"
	Spanish    LanguageCode = "es"
	Swedish    LanguageCode = "sv"
	Turkish    LanguageCode = "tr"
)

