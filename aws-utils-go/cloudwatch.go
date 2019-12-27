package aws_utils_go

import (
	"fmt"
	"github.com/adaptiveteam/core-utils-go/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type CloudWatchSchedule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enable      bool   `json:"enable"`
	Expression  string `json:"expression"`
}

type CloudWatchRequest struct {
	events *cloudwatchevents.CloudWatchEvents
	log    *logger.Logger
}

func NewCloudWatch(region, endpoint, namespace string) *CloudWatchRequest {
	session, config := sess(region, endpoint)
	return &CloudWatchRequest{
		events: cloudwatchevents.New(session, config),
		log:    logger.WithNamespace(fmt.Sprintf("adaptive.cloudwatch.%s", namespace)),
	}
}

func (c *CloudWatchRequest) errorLog(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		// cloudwatch logs error codes
		case cloudwatchlogs.ErrCodeInvalidParameterException:
			c.log.Error(cloudwatchlogs.ErrCodeInvalidParameterException, aerr.Error())
		case cloudwatchlogs.ErrCodeResourceNotFoundException:
			c.log.Error(cloudwatchlogs.ErrCodeResourceNotFoundException, aerr.Error())
		case cloudwatchlogs.ErrCodeServiceUnavailableException:
			c.log.Error(cloudwatchlogs.ErrCodeServiceUnavailableException, aerr.Error())

		// cloudwatch event error codes
		case cloudwatchevents.ErrCodeInvalidEventPatternException:
			c.log.Error(cloudwatchlogs.ErrCodeServiceUnavailableException, aerr.Error())
		case cloudwatchevents.ErrCodeLimitExceededException:
			c.log.Error(cloudwatchevents.ErrCodeLimitExceededException, aerr.Error())
		case cloudwatchevents.ErrCodeConcurrentModificationException:
			c.log.Error(cloudwatchevents.ErrCodeConcurrentModificationException, aerr.Error())
		case cloudwatchevents.ErrCodeInternalException:
			c.log.Error(cloudwatchevents.ErrCodeInternalException, aerr.Error())
		default:
			c.log.Error(aerr.Error())
		}
	} else {
		c.log.Error(err.Error())
	}
}

func (c *CloudWatchRequest) CreateOrUpdateSchedule(sc CloudWatchSchedule) (res string, err error) {
	c.log.Printf("Create schedule for cloudwatch, name: %s, cron: %s...\n", sc.Name, sc.Expression)
	state := "DISABLED"
	if sc.Enable {
		state = "ENABLED"
	}

	input := &cloudwatchevents.PutRuleInput{
		Description:        aws.String(sc.Description),
		Name:               aws.String(sc.Name),
		ScheduleExpression: aws.String(sc.Expression),
		State:              aws.String(state),
	}
	print(input, true)
	result, err2 := c.events.PutRule(input)
	err = err2
	if err != nil {
		c.errorLog(err)
		return
	}
	print(result, true)
	res = *result.RuleArn
	c.log.Info("Schedule event created successfully")
	return
}

func (c *CloudWatchRequest) DeleteSchedule(name string) error {
	c.log.Printf("Delete schedule from cloudwatch, name %s...\n", name)
	input := &cloudwatchevents.DeleteRuleInput{
		Name: aws.String(name),
	}
	print(input, true)
	result, err2 := c.events.DeleteRule(input)
	if err2 != nil {
		c.errorLog(err2)
		return err2
	}
	print(result, true)
	c.log.Info("Schedule event deleted successfully")
	return err2
}

// functionArn is the lambda function ARN
func (c *CloudWatchRequest) PutTarget(scheduleName, functionArn, id, json string) error {
	c.log.Printf("Putting schedule target of lambda function to %s...\n", scheduleName)
	input := &cloudwatchevents.PutTargetsInput{
		Rule: aws.String(scheduleName),
		Targets: []*cloudwatchevents.Target{
			{
				Arn:   aws.String(functionArn),
				Id:    aws.String(id),
				Input: aws.String(json),
			},
		},
	}
	print(input, true)
	result, err2 := c.events.PutTargets(input)
	if err2 != nil {
		c.errorLog(err2)
		return err2
	}
	print(result, true)
	c.log.Info("Put schedule target successfully")
	return err2
}
