package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
)

var (
	ap *auth.AuthProvider
)

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	prefix := os.Getenv("API_GATEWAY_PREFIX")

	var (
		fErr utils.StatusCodeErr
		res  *events.APIGatewayProxyResponse
	)

	if req.HTTPMethod == "POST" {
		if req.Path == "/"+prefix+"/auth/signin" {
			res, fErr = HandleSignIn(req)
		} else if req.Path == "/"+prefix+"/auth/signup" {
			res, fErr = HandleSignUp(req)
		} else {
			fErr = utils.NewStatusCodeErr(
				"no such route "+req.Path,
				httpcodes.StatusMethodNotAllowed,
			)
		}
	} else {
		fErr = utils.NewStatusCodeErr(
			"method not allowed",
			httpcodes.StatusMethodNotAllowed,
		)
	}

	if fErr != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode:      fErr.Status(),
			Headers:         applicationJsonHeader,
			IsBase64Encoded: false,
			Body: MarshalJSON(JSON{
				"error": fErr.Error(),
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
