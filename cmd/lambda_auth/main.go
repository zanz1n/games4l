package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/cmd/auth_lambda/src"
	"github.com/games4l/internal/errors"
)

func main() {
	errors.DefaultErrorList.Apply(errors.PtBtMessages)
	os.Setenv("NO_COLOR", "1")
	lambda.Start(src.Handler)
}
