package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/cmd/lambda_auth/src"
	"github.com/games4l/internal/errors"
	"github.com/games4l/internal/logger"
)

func main() {
	errors.DefaultErrorList.Apply(errors.PtBtMessages)
	logger.DefaultConfig.Colors = false
	lambda.Start(src.Handler)
}
