package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/cmd/lambda_question/src"
	"github.com/games4l/internal/errors"
)

func main() {
	errors.DefaultErrorList.Apply(errors.PtBtMessages)
	lambda.Start(src.Handler)
}
