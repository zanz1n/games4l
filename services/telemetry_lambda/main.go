package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
)

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Path == "/telemetry" {
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
			return events.APIGatewayProxyResponse{
				StatusCode:      fErr.Status(),
				Body:            fErr.Error(),
				IsBase64Encoded: false,
			}, nil
		}

		return *res, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusMethodNotAllowed,
		Body:            "Method not allowed",
		IsBase64Encoded: false,
	}, nil
}

func main() {
	os.Setenv("NO_COLOR", "1")
	logger.Init()
	lambda.Start(Handler)
}
