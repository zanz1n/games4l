package src

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/utils"
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

func init() {
	ap = auth.NewAuthProvider([]byte(os.Getenv("WEBHOOK_SIG")), []byte(os.Getenv("JWT_SIG")))
}
