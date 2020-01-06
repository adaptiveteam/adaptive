package aws_utils_go

import (
	"encoding/json"
	"fmt"
	"github.com/adaptiveteam/adaptive/core-utils-go/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SnsRequest struct {
	svc *sns.SNS
	log *logger.Logger
}

func NewSNS(region, endpoint, namespace string) *SnsRequest {
	session, config := sess(region, endpoint)
	return &SnsRequest{
		svc: sns.New(session, config),
		log: logger.WithNamespace(fmt.Sprintf("adaptive.sns.%s", namespace)),
	}
}

func (s *SnsRequest) errorLog(err error) {
	s.log.Errorf(err.Error())
}

func (s *SnsRequest) ListTopics(Next *string) ([]string, *string, error) {
	var arns []string
	input := &sns.ListTopicsInput{
		NextToken: Next,
	}
	op, err2 := s.svc.ListTopics(input)
	if err2 != nil {
		s.errorLog(err2)
		return nil, nil, err2
	}
	for _, each := range op.Topics {
		arns = append(arns, *each.TopicArn)
	}
	return arns, op.NextToken, nil
}

func (s *SnsRequest) CreateTopic(name string, attribs map[string]string) (*string, error) {
	awsAttribs := make(map[string]*string)
	for k, v := range attribs {
		awsAttribs[k] = aws.String(v)
	}
	input := &sns.CreateTopicInput{
		Name:       aws.String(name),
		Attributes: awsAttribs,
	}
	op, err2 := s.svc.CreateTopic(input)
	if err2 != nil {
		s.errorLog(err2)
		return nil, err2
	}
	return op.TopicArn, nil
}

func (s *SnsRequest) DeleteTopic(arn string) error {
	input := &sns.DeleteTopicInput{
		TopicArn: aws.String(arn),
	}
	_, err2 := s.svc.DeleteTopic(input)
	if err2 != nil {
		s.errorLog(err2)
	}
	return err2
}

func (s *SnsRequest) Publish(item interface{}, topicArn string) (*string, error) {
	bytes, err2 := json.Marshal(item)
	if err2 != nil {
		s.errorLog(err2)
		return nil, err2
	}
	input := &sns.PublishInput{
		Message:  aws.String(string(bytes)),
		TopicArn: aws.String(topicArn),
	}
	resp, err2 := s.svc.Publish(input) //Call to puclish the message

	if err2 != nil {
		s.errorLog(err2)
		return nil, err2
	}
	return resp.MessageId, nil
}
