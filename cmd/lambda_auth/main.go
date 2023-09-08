package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/cmd/lambda_auth/src"
	"github.com/games4l/internal/logger"
)

func main() {
	logger.DefaultConfig.Colors = false
	lambda.Start(src.Handler)
}
