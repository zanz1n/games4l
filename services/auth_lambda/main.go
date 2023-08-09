package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/backend/libs/auth"
	userlib "github.com/games4l/backend/libs/user"
	"github.com/games4l/backend/libs/utils"
)

var (
	ap  *auth.AuthProvider
	dba *userlib.UserService
)

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	prefix := os.Getenv("API_GATEWAY_PREFIX")

	var (
		fErr utils.StatusCodeErr
		res  *events.APIGatewayProxyResponse
	)

	if req.HTTPMethod == "POST" {
		if req.Path == "/"+prefix+"/auth/signin" || req.Path == "/auth/signin" {
			res, fErr = HandleSignIn(req)
		} else if req.Path == "/"+prefix+"/user" || req.Path == "/user" {
			res, fErr = HandleUserCreation(req)
		} else {
			fErr = utils.DefaultErrorList.NoSuchRoute
		}
	} else {
		fErr = utils.DefaultErrorList.MethodNotAllowed
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
	utils.DefaultErrorList.Apply(utils.PtBtMessages)
	os.Setenv("NO_COLOR", "1")
	ap = auth.NewAuthProvider([]byte(os.Getenv("WEBHOOK_SIG")), []byte(os.Getenv("JWT_SIG")))
	lambda.Start(Handler)
}
