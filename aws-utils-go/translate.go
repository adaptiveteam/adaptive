package aws_utils_go

import (
	"fmt"
	"github.com/adaptiveteam/core-utils-go/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/translate"
)

type TranslateRequest struct {
	svc *translate.Translate
	log *logger.Logger
}

func NewTranslate(region, endpoint, namespace string) *TranslateRequest {
	session, config := sess(region, endpoint)
	return &TranslateRequest{
		svc: translate.New(session, config),
		log: logger.WithNamespace(fmt.Sprintf("adaptive.translate.%s", namespace)),
	}
}

func (s *TranslateRequest) errorLog(err error) {
	s.log.Error(err.Error())
}

func (s *TranslateRequest) TranslateText(text string, sl, tl string) (*translate.TextOutput, error) {
	input := &translate.TextInput{
		// source language
		SourceLanguageCode: aws.String(sl),
		// target language
		TargetLanguageCode: aws.String(tl),
		Text:               aws.String(text),
	}
	return s.svc.Text(input)
}
