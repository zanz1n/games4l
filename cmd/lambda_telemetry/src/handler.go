package src

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/auth"
	"github.com/games4l/internal/errors"
	"github.com/games4l/internal/utils"
)

func hPrx(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ap == nil {
		ap = auth.NewAuthProvider([]byte(os.Getenv("WEBHOOK_SIG")), []byte(os.Getenv("JWT_SIG")))
	}

	prefix := os.Getenv("API_GATEWAY_PREFIX")

	var (
		fErr error
		res  *events.APIGatewayProxyResponse
	)

	if req.Path == "/"+prefix+"/telemetry" || req.Path == "/telemetry" {
		if req.HTTPMethod == "POST" {
			res, fErr = HandlePost(req)
		} else if req.HTTPMethod == "GET" {
			res, fErr = HandleGetByName(req)
		} else {
			fErr = errors.ErrMethodNotAllowed
		}
	} else if hPrx(req.Path, "/"+prefix+"/telemetry/") || hPrx(req.Path, "/telemetry") {
		if req.HTTPMethod == "GET" {
			res, fErr = HandleGetByID(req)
		} else {
			fErr = errors.ErrMethodNotAllowed
		}
	} else {
		fErr = errors.ErrNoSuchRoute
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
