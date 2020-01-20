package main

import (
	lambda "github.com/adaptiveteam/adaptive/lambdas/user-engagement-lambda-go"
	ls "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ls.Start(lambda.HandleRequest)
}

