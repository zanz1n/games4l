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
