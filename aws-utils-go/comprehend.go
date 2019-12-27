package aws_utils_go

import (
	"fmt"
	"github.com/adaptiveteam/core-utils-go/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehend"
)

type ComprehendRequest struct {
	svc *comprehend.Comprehend
	log *logger.Logger
}

func NewComprehend(region, endpoint, namespace string) *ComprehendRequest {
	session, config := sess(region, endpoint)
	return &ComprehendRequest{
		svc: comprehend.New(session, config),
		log: logger.WithNamespace(fmt.Sprintf("adaptive.comprehend.%s", namespace)),
	}
}

func (s *ComprehendRequest) errorLog(err error) {
	s.log.Error(err.Error())
}

func (c *ComprehendRequest) DetectSyntax(lc, text string) (*comprehend.DetectSyntaxOutput, error) {
	input := &comprehend.DetectSyntaxInput{
		// language code
		LanguageCode: aws.String(lc),
		Text:         aws.String(text),
	}
	return c.svc.DetectSyntax(input)
}
