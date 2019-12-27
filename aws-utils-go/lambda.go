package aws_utils_go

import (
	"fmt"
	"github.com/adaptiveteam/core-utils-go/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type LambdaFunction struct {
	Name       string  `json:"name"`
	Handler    string  `json:"handler"`
	Role       string  `json:"role"`
	MemorySize int64   `json:"memory_size"`
	Timeout    int64   `json:"timeout"`
	Schedule   *string `json:"schedule"`
}

type LambdaRequest struct {
	svc *lambda.Lambda
	log *logger.Logger
}

func NewLambda(region, endpoint, namespace string) *LambdaRequest {
	session, config := sess(region, endpoint)
	return &LambdaRequest{
		svc: lambda.New(session, config),
		log: logger.WithNamespace(fmt.Sprintf("adaptive.lambda.%s", namespace)),
	}
}

func (l *LambdaRequest) errorLog(err error, skipCodes ...string) {
	if aerr, ok := err.(awserr.Error); ok {
		code := aerr.Code()
		for _, c := range skipCodes {
			if c == code {
				return
			}
		}
		switch code {
		case lambda.ErrCodeServiceException:
			l.log.Error(lambda.ErrCodeServiceException, aerr.Error())
		case lambda.ErrCodeResourceNotFoundException:
			l.log.Error(lambda.ErrCodeResourceNotFoundException, aerr.Error())
		case lambda.ErrCodeInvalidRequestContentException:
			l.log.Error(lambda.ErrCodeInvalidRequestContentException, aerr.Error())
		case lambda.ErrCodeRequestTooLargeException:
			l.log.Error(lambda.ErrCodeRequestTooLargeException, aerr.Error())
		case lambda.ErrCodeUnsupportedMediaTypeException:
			l.log.Error(lambda.ErrCodeUnsupportedMediaTypeException, aerr.Error())
		case lambda.ErrCodeTooManyRequestsException:
			l.log.Error(lambda.ErrCodeTooManyRequestsException, aerr.Error())
		case lambda.ErrCodeInvalidParameterValueException:
			l.log.Error(lambda.ErrCodeInvalidParameterValueException, aerr.Error())
		case lambda.ErrCodeEC2UnexpectedException:
			l.log.Error(lambda.ErrCodeEC2UnexpectedException, aerr.Error())
		case lambda.ErrCodeSubnetIPAddressLimitReachedException:
			l.log.Error(lambda.ErrCodeSubnetIPAddressLimitReachedException, aerr.Error())
		case lambda.ErrCodeENILimitReachedException:
			l.log.Error(lambda.ErrCodeENILimitReachedException, aerr.Error())
		case lambda.ErrCodeEC2ThrottledException:
			l.log.Error(lambda.ErrCodeEC2ThrottledException, aerr.Error())
		case lambda.ErrCodeEC2AccessDeniedException:
			l.log.Error(lambda.ErrCodeEC2AccessDeniedException, aerr.Error())
		case lambda.ErrCodeInvalidSubnetIDException:
			l.log.Error(lambda.ErrCodeInvalidSubnetIDException, aerr.Error())
		case lambda.ErrCodeInvalidSecurityGroupIDException:
			l.log.Error(lambda.ErrCodeInvalidSecurityGroupIDException, aerr.Error())
		case lambda.ErrCodeInvalidZipFileException:
			l.log.Error(lambda.ErrCodeInvalidZipFileException, aerr.Error())
		case lambda.ErrCodeKMSDisabledException:
			l.log.Error(lambda.ErrCodeKMSDisabledException, aerr.Error())
		case lambda.ErrCodeKMSInvalidStateException:
			l.log.Error(lambda.ErrCodeKMSInvalidStateException, aerr.Error())
		case lambda.ErrCodeKMSAccessDeniedException:
			l.log.Error(lambda.ErrCodeKMSAccessDeniedException, aerr.Error())
		case lambda.ErrCodeKMSNotFoundException:
			l.log.Error(lambda.ErrCodeKMSNotFoundException, aerr.Error())
		case lambda.ErrCodeInvalidRuntimeException:
			l.log.Error(lambda.ErrCodeInvalidRuntimeException, aerr.Error())
		default:
			l.log.Error(aerr.Error())
		}
	} else {
		l.log.Errorf(err.Error())
	}
}

func (l *LambdaRequest) FunctionExists(name string) bool {
	_, err2 := l.svc.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(name),
	})
	return err2 == nil
}

func (l *LambdaRequest) CreateFunction(fn *LambdaFunction, zipBytes []byte) (string, error) {
	input := &lambda.CreateFunctionInput{
		Code: &lambda.FunctionCode{
			ZipFile: zipBytes,
		},
		FunctionName: aws.String(fn.Name),
		Handler:      aws.String(fn.Handler),
		Role:         aws.String(fn.Role),
		MemorySize:   aws.Int64(fn.MemorySize),
		Publish:      aws.Bool(true),
		Timeout:      aws.Int64(fn.Timeout),
		Runtime:      aws.String("go1.x"),
	}
	print(input, true)
	result, err2 := l.svc.CreateFunction(input)
	if err2 != nil {
		l.errorLog(err2)
		return "", err2
	}
	print(result, true)
	return *result.FunctionArn, nil
}

// eventArn is the cloudwatch rule arn
func (l *LambdaRequest) AddCloudWatchPermission(name, eventArn string) error {
	l.log.Printf("Add CloudWatch permission for %s...\n", name)
	input := &lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),
		Principal:    aws.String("events.amazonaws.com"),
		FunctionName: aws.String(name),
		StatementId:  aws.String(generateStatementId("cloudwatch")),
		SourceArn:    aws.String(eventArn),
	}
	_, err2 := l.svc.AddPermission(input)
	if err2 != nil {
		l.errorLog(err2)
	} else {
		l.log.Info("Permission added successfully")
	}
	return err2
}

func (l *LambdaRequest) InvokeFunction(name string, payload []byte, async bool) (*lambda.InvokeOutput, error) {
	var invType = "RequestResponse"
	if async {
		invType = "Event"
	}
	input := &lambda.InvokeInput{
		FunctionName:   aws.String(name),
		Payload:        payload,
		InvocationType: aws.String(invType),
	}
	result, err2 := l.svc.Invoke(input)
	if err2 != nil {
		l.errorLog(err2)
		return nil, err2
	}
	if result.FunctionError != nil {
		l.log.Errorf("Function invoked on %s and handed error: %s\n", name, *result.FunctionError)
	} else {
		l.log.Infof("Function invoked on %s and succeeded\n", name)
	}
	return result, err2
}

func (l *LambdaRequest) ListLambdas() (*lambda.ListFunctionsOutput, error) {
	return l.svc.ListFunctions(nil)
}
