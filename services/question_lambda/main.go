package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/question"
)

var (
	ap  *auth.AuthProvider
	dba *question.QuestionService
)

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	panic("")
}

func main() {
	os.Setenv("NO_COLOR", "1")
	logger.Init()
	ap = auth.NewAuthProvider([]byte(os.Getenv("WEBHOOK_SIG")), []byte(os.Getenv("JWT_SIG")))
	lambda.Start(Handler)
}
