package logger

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

// AdaptiveLogger is a container for logrus logger
type AdaptiveLogger struct {
	*logrus.Entry
}

func LambdaLogger(level logrus.Level) AdaptiveLogger {
	baseLogger := logrus.New()
	baseLogger.SetLevel(level)
	baseLogger.SetFormatter(&logrus.JSONFormatter{})
	return AdaptiveLogger{baseLogger.WithFields(logrus.Fields{})}
}

// WithLambdaContext sets lambda context for the logger and logs lambda_request_id each time
func (a AdaptiveLogger) WithLambdaContext(ctx context.Context) AdaptiveLogger {
	var requestID string

	if lc, ok := lambdacontext.FromContext(ctx); ok {
		requestID = lc.AwsRequestID
	} else {
		a.WithField("context", ctx).Warn("Fail to extract lambda context")
		requestID = "N/A"
	}

	return AdaptiveLogger{a.
		WithField("lambda_request_id", requestID).
		WithField("file", fileInfo(2))}
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
