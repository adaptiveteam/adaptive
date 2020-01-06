package nlp

type McRestResponse struct {
	payload    []byte
	statusCode int
}

type McResponseStatus struct {
	//	* 0: OK
	//	* 100: Operation denied
	//	* 101: License expired
	//	* 102: Credits per subscription exceeded
	//	* 103: Request too large
	//	* 104: Request rate limit exceeded
	//	* 200: Missing required parameter(s) - [name of the parameter]
	//	* 202: Engine internal error
	//	* 203: Cannot connect to service
	//	* 205: Language not supported
	//	* 212: No content to analyze
	//	* 215: Timeout exceeded for service response
	Code             string `json:"code"`
	Msg              string `json:"msg"`
	Credits          string `json:"credits"`
	RemainingCredits string `json:"remaining_credits"`
}

type McResultList struct {
	Text string `json:"text"`
	Type string `json:"type"`
	// Initial position of the result, starting on 0
	Inip string `json:"inip"`
	// End position of the result
	Endp string `json:"endp"`
	// Level of the error
	Level      string `json:"level"`
	Bop        string `json:"bop"`
	Confidence string `json:"confidence"`
	ScoreTag   string `json:"score_tag"`
	Form       string `json:"form"`
	ID         string `json:"id"`
	Variant    string `json:"variant"`
}

// meaningCloudSentimentAnalyisResponse was automatically generated from this site:
// https://mholt.github.io/json-to-go/
// for this JSON response:
// https://www.meaningcloud.com/developer/sentiment-analysis/doc/2.1/response
type meaningCloudSentimentAnalysisResponse struct {
	Status       McResponseStatus `json:"status"`
	Model        string           `json:"model"`
	ScoreTag     string           `json:"score_tag"`
	Agreement    string           `json:"agreement"`
	Subjectivity string           `json:"subjectivity"`
	Confidence   string           `json:"confidence"`
	Irony        string           `json:"irony"`
	SentenceList []struct {
		McResultList
		Agreement   string `json:"agreement"`
		SegmentList []struct {
			SegmentType string `json:"segment_type"`
			McResultList
			PolarityTermList []struct {
				McResultList
				SentimentedConceptList []McResultList `json:"sentimented_concept_list"`
			} `json:"polarity_term_list"`
			SentimentedEntityList []McResultList `json:"sentimented_entity_list"`
		} `json:"segment_list"`
		SentimentedEntityList  []McResultList `json:"sentimented_entity_list"`
		SentimentedConceptList []McResultList `json:"sentimented_concept_list"`
	} `json:"sentence_list"`
	SentimentedEntityList  []McResultList `json:"sentimented_entity_list"`
	SentimentedConceptList []McResultList `json:"sentimented_concept_list"`
}

// meaningCloudSummarizationResponse was automatically generated from this site:
// https://mholt.github.io/json-to-go/
// for this JSON response:
// https://www.meaningcloud.com/developer/summarization/doc/1.0/response
type meaningCloudSummarizationResponse struct {
	Status  McResponseStatus `json:"status"`
	Summary string           `json:"summary"`
}

// meaningCloudDeepCategorizationResponse was automatically generated from this site:
// https://mholt.github.io/json-to-go/
// for this JSON response:
// https://www.meaningcloud.com/developer/deep-categorization/doc/1.0/response
type meaningCloudDeepCategorizationResponse struct {
	Status       McResponseStatus `json:"status"`
	CategoryList []struct {
		Code         string `json:"code"`
		Label        string `json:"label"`
		AbsRelevance string `json:"abs_relevance"`
		Relevance    string `json:"relevance"`
		Polarity     string `json:"polarity"`
		TermList     []struct {
			Form         string `json:"form"`
			AbsRelevance string `json:"abs_relevance"`
			OffsetList   []struct {
				Inip string `json:"inip"`
				Endp string `json:"endp"`
			} `json:"offset_list"`
		} `json:"term_list"`
	} `json:"category_list"`
}

// textSentiment contains sentiment information derived from the
// Meaning Cloud sentiment analysis service.
type textSentiment struct {
	Sentiment    string
	Confidence   int
	Agreement    bool
	Subjectivity bool
	Irony        bool
}

// textCategory contains deep classification category information from the
// Meaning Cloud deep classification service
type textCategory struct {
	Label        string
	Relevance    int
	AbsRelevance int
	Sentiment    string
}
