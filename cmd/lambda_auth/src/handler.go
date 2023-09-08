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
		fErr error
		res  *events.APIGatewayProxyResponse
	)

	if req.HTTPMethod == "POST" {
		if req.Path == "/"+prefix+"/auth/signin" || req.Path == "/auth/signin" {
			res, fErr = HandleSignIn(req)
		} else if req.Path == "/"+prefix+"/user" || req.Path == "/user" {
			res, fErr = HandleUserCreation(req)
		} else {
			fErr = errors.ErrNoSuchRoute
		}
	} else {
		fErr = errors.ErrMethodNotAllowed
	}

	if fErr != nil {
		errInfo := errors.GetStatusErr(fErr)
		errBody := errors.ErrorBody{
			Message:   errInfo.Message(),
			ErrorCode: errInfo.CustomCode(),
		}

		res = &events.APIGatewayProxyResponse{
			StatusCode:      errInfo.HttpCode(),
			Headers:         applicationJsonHeader,
			IsBase64Encoded: false,
			Body:            utils.B2S(errBody.Marshal()),
		}
	}

	return *res, nil
}
