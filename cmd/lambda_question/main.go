package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/cmd/lambda_question/handler"
)

func main() {
	s := handler.NewEnvServer()

	lambda.Start(s.RequestHandler)
}
