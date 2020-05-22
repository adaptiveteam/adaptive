package nlp

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

var mappings = map[string]string{
	"P+":   "very positive",
	"P":    "positive",
	"NEU":  "neutral",
	"N":    "negative",
	"N+":   "very negative",
	"NONE": "neutral",
}

func appendError(err error, errList []error) {
	if err != nil {
		errList = append(errList, err)
	}
}

func mcError(m McResponseStatus) error {
	return errors.New(fmt.Sprintf("Received %s status from MeaningCloud with message: %s", m.Code, m.Msg))
}

// GetTextCategories calls the following Meaning Cloud service:
// https://www.meaningcloud.com/developer/deep-categorization/doc/1.0/what-is-deep-categorization
// The result is a deep categorization of the contents of the provided content.
// The key/value parameters are used to drive the source to evaluate - URL or text.
// If the the key is 'txt' then the value must be the text to evaluate.
// If the key is 'url' then the value must be a web page to evaluate.
func (m MeaningCloud) GetTextCategories(key string, value string, l LanguageCode) (rv []TextCategories, err error) {
	if key == "txt" || key == "url" {
		if len(value) > 0 {
			parameters := map[string]string{
				"of":       ofJson,
				"lang":     *l.String(),
				"model":    "VoE-Performance_en",
				"polarity": "y",
				"verbose":  "y",
				key:        value,
			}
			response, err := m.HitMeaningCloudService(deepCategorizationURL, parameters)
			if err == nil && response.statusCode == 200 {
				retries := 0
				for rateLimitError := "104"; rateLimitError == "104" && retries < 30; time.Sleep(1 * time.Second) {
					var m meaningCloudDeepCategorizationResponse
					err = json.Unmarshal(response.payload, &m)
					if m.Status.Code != "0" {
						rateLimitError = m.Status.Code
						err = mcError(m.Status)
						retries++
					} else if err == nil {
						rv, err = m.toTextCategories()
						rateLimitError = ""
					}
				}
			} else {
				fmt.Println("Error in GetTextCategories key="+key+", value="+value, response, err)
			}
		}
	} else {
		err = errors.New("GetTextCategories: Expected txt or url but instead got key=" + key)
	}

	return rv, err
}

// getTextSentiment calls the following Meaning Cloud service:
// https://www.meaningcloud.com/developer/sentiment-analysis/doc/2.1/what-is-sentiment-analysis
// The result is several pieces of sentiment information accessible by the TextSentiment interface
// The key/value parameters are used to drive the source to evaluate - URL or text.
// If the the key is 'txt' then the value must be the text to evaluate.
// If the key is 'url' then the value must be a web page to evaluate.
func (m MeaningCloud) GetTextSentiment(key string, value string, l LanguageCode) (rv TextSentiment, err error) {
	if key == "txt" || key == "url" {
		// sentimentAnalysisURL := "https://api.meaningcloud.com/sentiment-2.1/post"
		if len(value) > 0 {
			parameters := map[string]string{
				"of":    ofJson,
				"lang":  *l.String(),
				"model": "general",
				"egp":   "n",
				"rt":    "n",
				"uw":    "n",
				"dm":    "s",
				"sdg":   "l",
				key:     value,
			}
			var response McRestResponse
			response, err = m.HitMeaningCloudService(sentimentAnalysisURL, parameters)
			if err == nil && response.statusCode == 200 {
				retries := 0
				for rateLimitError := "104"; rateLimitError == "104" && retries < 30; time.Sleep(1 * time.Second) {
					var m meaningCloudSentimentAnalysisResponse
					err = json.Unmarshal(response.payload, &m)
					if m.Status.Code != "0" {
						rateLimitError = m.Status.Code
						err = mcError(m.Status)
						retries++
					} else if err == nil {
						rv, err = m.toTextSentiment()
						rateLimitError = ""
					}
				}
			} else {
				fmt.Printf("ERROR in MeaningCloud) GetTextSentiment: response=%v\nerr=\n%+v\n", response, err)
				if err == nil {
					err = errors.Errorf("Invalid response code from MeaningCloud: %d", response.statusCode)
				}
			}
		} else {
			rv = textSentiment{
				Sentiment:    "none",
				Confidence:   100,
				Agreement:    false,
				Subjectivity: false,
				Irony:        false,
			}
		}
	} else {
		return nil, errors.New("TextSentiment: Expected txt or url but instead got key=" + key)
	}
	return rv, err
}

// getSummary calls the following Meaning Cloud service:
// https://www.meaningcloud.com/developer/summarization/doc
// The result is a numSentences long summary of the provided content.
// The key/value parameters are used to drive the source to evaluate - URL or text.
// If the the key is 'txt' then the value must be the text to evaluate.
// If the key is 'url' then the value must be a web page to evaluate.
func (m MeaningCloud) GetSummary(numSentences int, key string, value string, l LanguageCode) (rv string, err error) {
	if key == "txt" || key == "url" {
		if len(value) > 0 {
			parameters := map[string]string{
				"of":        ofJson,
				"lang":      *l.String(),
				"sentences": strconv.Itoa(numSentences),
				key:         value,
			}
			response, err := m.HitMeaningCloudService(summaryURL, parameters)
			if err == nil && response.statusCode == 200 {
				retries := 0
				for rateLimitError := "104"; rateLimitError == "104" && retries < 30; time.Sleep(1 * time.Second) {
					var m meaningCloudSummarizationResponse
					err = json.Unmarshal(response.payload, &m)
					if m.Status.Code != "0" {
						rateLimitError = m.Status.Code
						err = mcError(m.Status)
						retries++
					} else if err == nil {
						rv = m.Summary
						rateLimitError = ""
					}
				}
			} else {
				fmt.Println("Error in GetSummary", response, err)
			}
		}
	} else {
		err = errors.New("Summary: Expected txt or url but instead got key=" + key)
	}
	return rv, err
}

// toTextCategories converts the meaningCloudDeepCategorizationResponse that represents the
// JSON response from the deep categorization service from Meaning Cloud:
// https://www.meaningcloud.com/developer/deep-categorization/doc/1.0/response
// This method also maps the Meaning Cloud specific sentiment language to an Adaptive standard
func (m *meaningCloudDeepCategorizationResponse) toTextCategories() ([]TextCategories, error) {
	var rv []TextCategories = nil
	var err error = nil
	for i := 0; i < len(m.CategoryList) && err == nil; i++ {
		each := m.CategoryList[i]
		relevance, relevanceErr := strconv.Atoi(each.Relevance)
		absRelevance, absRelevanceErr := strconv.Atoi(each.AbsRelevance)
		if relevanceErr == nil || absRelevanceErr == nil {
			rv = append(rv, &textCategory{
				Label:        each.Label,
				Relevance:    relevance,
				AbsRelevance: absRelevance,
				Sentiment:    mappings[each.Polarity],
			})
		} else {
			err = errors.New("relevanceErr " + relevanceErr.Error() + ", " + "absRelevance " + absRelevanceErr.Error())
		}
	}
	// time.Sleep(1 * time.Second)
	return rv, err
}

// toTextSentiment converts the meaningCloudSentimentAnalyisResponse that represents the
// JSON response from the sentiment analysis service from Meaning Cloud:
// https://www.meaningcloud.com/developer/sentiment-analysis/doc
// This method also maps the Meaning Cloud specific sentiment language to an Adaptive standard
func (m *meaningCloudSentimentAnalysisResponse) toTextSentiment() (TextSentiment, error) {
	var rv textSentiment
	confidence, err := strconv.Atoi(m.Confidence)

	if err == nil {
		rv.Confidence = confidence
		rv.Sentiment = mappings[m.ScoreTag]
		rv.Agreement = m.Agreement == "AGREEMENT"
		rv.Irony = !(m.Irony == "NONIRONIC")
		rv.Subjectivity = m.Subjectivity == "SUBJECTIVE"
	}
	return &rv, err
}
