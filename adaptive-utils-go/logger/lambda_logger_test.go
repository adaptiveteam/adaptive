package logger

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/google/uuid"

	. "github.com/onsi/ginkgo"
	"github.com/sirupsen/logrus"
)

type mockLambdaContext struct{}

func (ctx mockLambdaContext) Deadline() (deadline time.Time, ok bool) {
	return deadline, ok
}

func (ctx mockLambdaContext) Done() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func (ctx mockLambdaContext) Err() error {
	return context.DeadlineExceeded
}

func (ctx mockLambdaContext) Value(key interface{}) interface{} {
	return &lambdacontext.LambdaContext{
		AwsRequestID: uuid.New().String(),
	}
}

type mockContext struct{}

func (ctx mockContext) Deadline() (deadline time.Time, ok bool) {
	return deadline, ok
}

func (ctx mockContext) Done() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func (ctx mockContext) Err() error {
	return context.DeadlineExceeded
}

func (ctx mockContext) Value(key interface{}) interface{} {
	return ""
}

var _ = Describe("LambdaLogger Tests", func() {
	baseLogger := LambdaLogger(logrus.InfoLevel)
	Context("Lambda Context", func() {
		It("should log with lambda context", func() {
			lambdaLogger := baseLogger.WithLambdaContext(mockLambdaContext{})
			lambdaLogger.Infof("test logging with context")
		})

		It("should log without lambda context", func() {
			lambdaLogger := baseLogger.WithLambdaContext(mockContext{})
			lambdaLogger.Infof("test logging with no context")
		})
	})
})
