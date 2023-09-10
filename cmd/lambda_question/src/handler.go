package src

import (
	"os"
	"strings"

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

	if strings.HasPrefix(req.Path, "/question") || strings.HasPrefix(req.Path, "/"+prefix+"/question") {
		if req.HTTPMethod == "GET" {
			if _, ok := req.PathParameters["id"]; ok {
				res, fErr = HandleGetOne(req)
			} else {
				res, fErr = HandleGetMany(req)
			}
		} else if req.HTTPMethod == "POST" {
			res, fErr = HandlePostOne(req)
		} else if req.HTTPMethod == "PATCH" || req.HTTPMethod == "PUT" {
			fErr = errors.ErrMethodNotAllowed
			// update one route
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
