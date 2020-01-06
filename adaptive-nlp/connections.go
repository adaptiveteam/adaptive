package nlp

import (
	"github.com/adaptiveteam/adaptive/aws-utils-go"
)
// Connections contains connections to AWS services
type Connections struct {
	Comprehend *aws_utils_go.ComprehendRequest
	Translate  *aws_utils_go.TranslateRequest
	MeaningCloud MeaningCloud

}

// OpenConnections opens connections to AWS NLP services
func OpenConnections(region string, meaningCloudKey string) (connections Connections) {
	return Connections{
		Translate: aws_utils_go.NewTranslate(region, "", "nlp-translate"),
		Comprehend: aws_utils_go.NewComprehend(region, "", "nlp-comprehend"),
		MeaningCloud: NewMeaningCloud(meaningCloudKey),
	}
}
