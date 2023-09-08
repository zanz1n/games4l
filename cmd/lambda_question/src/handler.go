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
			fErr = errors.DefaultErrorList.MethodNotAllowed
		}
	} else {
		fErr = errors.DefaultErrorList.NoSuchRoute
	}

	if fErr != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode:      fErr.Status(),
			Headers:         applicationJsonHeader,
			IsBase64Encoded: false,
			Body: utils.MarshalJSON(JSON{
				"error": utils.FirstUpper(fErr.Error()),
			}),
		}
	}

	return *res, nil
}
