package src

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/auth"
	"github.com/games4l/internal/errors"
	"github.com/games4l/internal/utils"
)

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ap == nil {
		ap = auth.NewAuthProvider([]byte(os.Getenv("WEBHOOK_SIG")), []byte(os.Getenv("JWT_SIG")))
	}

	prefix := os.Getenv("API_GATEWAY_PREFIX")

	var (
		fErr errors.StatusCodeErr
		res  *events.APIGatewayProxyResponse
	)

	if req.HTTPMethod == "POST" {
		if req.Path == "/"+prefix+"/auth/signin" || req.Path == "/auth/signin" {
			res, fErr = HandleSignIn(req)
		} else if req.Path == "/"+prefix+"/user" || req.Path == "/user" {
			res, fErr = HandleUserCreation(req)
		} else {
			fErr = errors.DefaultErrorList.NoSuchRoute
		}
	} else {
		fErr = errors.DefaultErrorList.MethodNotAllowed
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
