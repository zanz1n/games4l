package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/cmd/lambda_telemetry/src"
)

func main() {
	lambda.Start(src.Handler)
}
