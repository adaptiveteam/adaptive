package lambda

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/aws/aws-lambda-go/events"

)
// IMPORTANT: It should always return this with the empty body. Else actions won't work.
var responseOk = events.APIGatewayProxyResponse{
	StatusCode: 200,
	Headers: map[string]string{
		"Content-type": "application/json; charset=utf-8",
	},
}

func responseOkBody(body string)  (response events.APIGatewayProxyResponse) {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: 200,
	}
}

func responseServerError(err error) (response events.APIGatewayProxyResponse) {
	errNum := random.Random(0, 999_999)
	errorMessage := fmt.Sprintf("Server error. See log for details (%d)", errNum)
	response = events.APIGatewayProxyResponse{
		StatusCode: 500, // Server error
		Body: errorMessage,
		Headers: map[string]string{
			"Content-type": "text/plain; charset=utf-8",
		},
	}
	logger.WithError(err).Errorf(errorMessage)
	return
}

func responsePermanentRedirect(location string) (response events.APIGatewayProxyResponse) {
	response = events.APIGatewayProxyResponse{
		StatusCode: 308, // Permanent Redirect
		Headers: map[string]string{
			"Location": location,
		},
	}
	logger.Infof("Permanent Redirect (308), Location: %s", location)
	return
}
