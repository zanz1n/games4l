package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/question"
	"github.com/games4l/backend/libs/utils"
)

var (
	ap  *auth.AuthProvider
	dba *question.QuestionService
)

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	prefix := os.Getenv("API_GATEWAY_PREFIX")

	var (
		fErr utils.StatusCodeErr
		res  *events.APIGatewayProxyResponse
	)

	if req.Path == "/"+prefix+"/question" || req.Path == "/question" {
		if req.HTTPMethod == "GET" {
			_, ok1 := req.QueryStringParameters["id"]
			_, ok2 := req.QueryStringParameters["uid"]

			if ok1 || ok2 {
				res, fErr = HandleGetByID(req)
			} else {
				res, fErr = HandleGetMany(req)
			}
		} else if req.HTTPMethod == "POST" {
			res, fErr = HandlePost(req)
		} else {
			fErr = utils.DefaultErrorList.MethodNotAllowed
		}
	} else {
		fErr = utils.DefaultErrorList.NoSuchRoute
	}

	if fErr != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode:      fErr.Status(),
			Headers:         applicationJsonHeader,
			IsBase64Encoded: false,
			Body: MarshalJSON(JSON{
				"error": utils.FirstUpper(fErr.Error()),
			}),
		}
	}

	return *res, nil
}

func main() {
	os.Setenv("NO_COLOR", "1")
	logger.Init()
	ap = auth.NewAuthProvider([]byte(os.Getenv("WEBHOOK_SIG")), []byte(os.Getenv("JWT_SIG")))
	lambda.Start(Handler)
}
