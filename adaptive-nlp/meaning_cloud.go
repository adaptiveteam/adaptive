package nlp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type MeaningCloud struct {
	Key string
	BaseURL string
}

const (
	baseURL string = "https://api.meaningcloud.com/"
	// TODO: replace these with Terraform provided environment variables
	deepCategorizationURL = baseURL + "/deepcategorization-1.0/post"
	sentimentAnalysisURL  = baseURL + "/sentiment-2.1/post"
	summaryURL            = baseURL + "/summarization-1.0"
	ofJson                = `json`
	globalMeaningCloudKey = "d677938b7bf41ba18b7fbffb18ad6730"
)

// NewMeaningCloud creates a connection to meaning cloud
func NewMeaningCloud(key string) MeaningCloud {
	return MeaningCloud{Key: key, BaseURL: baseURL}
}

// HitMeaningCloudService hits a Meaning Cloud service at the given URL:
// https://www.meaningcloud.com/developer/documentation
// The parameters are passed in as a map. The response is a standard HTTP response plus an error.
func (m MeaningCloud)HitMeaningCloudService(url string, parameters map[string]string) (rv McRestResponse, err error) {
	parameters["key"] = m.Key
	jsonValue, _ := json.Marshal(parameters)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	rv = McRestResponse{body, response.StatusCode}
	return rv, err
}
