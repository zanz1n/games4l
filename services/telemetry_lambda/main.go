package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
)

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	prefix := os.Getenv("API_GATEWAY_PREFIX")

	if req.Path == "/"+prefix+"/telemetry" {
		var (
			fErr utils.StatusCodeErr
			res  *events.APIGatewayProxyResponse
		)

		if req.HTTPMethod == "POST" || req.HTTPMethod == "GET" {
			Connect()
		}

		if req.HTTPMethod == "POST" {
			res, fErr = HandlePost(req)
		} else if req.HTTPMethod == "GET" {
			res, fErr = HandleGetByName(req)
		}

		if fErr != nil {
			errBody, _ := json.Marshal(JSON{
				"error": fErr.Error(),
			})

			return events.APIGatewayProxyResponse{
				StatusCode:      fErr.Status(),
				Headers:         applicationJsonHeader,
				Body:            string(errBody),
				IsBase64Encoded: false,
			}, nil
		}

		return *res, nil
	}

	errBody, _ := json.Marshal(JSON{
		"error": "method not allowed",
	})

	return events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusMethodNotAllowed,
		Headers:         applicationJsonHeader,
		Body:            string(errBody),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	os.Setenv("NO_COLOR", "1")
	logger.Init()
	lambda.Start(Handler)
}
