package src

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/utils"
)

func hPrx(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ap = auth.NewAuthProvider([]byte(os.Getenv("WEBHOOK_SIG")), []byte(os.Getenv("JWT_SIG")))

	prefix := os.Getenv("API_GATEWAY_PREFIX")

	var (
		fErr utils.StatusCodeErr
		res  *events.APIGatewayProxyResponse
	)

	if req.Path == "/"+prefix+"/telemetry" || req.Path == "/telemetry" {
		if req.HTTPMethod == "POST" {
			res, fErr = HandlePost(req)
		} else if req.HTTPMethod == "GET" {
			res, fErr = HandleGetByName(req)
		} else {
			fErr = utils.DefaultErrorList.MethodNotAllowed
		}
	} else if hPrx(req.Path, "/"+prefix+"/telemetry/") || hPrx(req.Path, "/telemetry") {
		if req.HTTPMethod == "GET" {
			res, fErr = HandleGetByID(req)
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