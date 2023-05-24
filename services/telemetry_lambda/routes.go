package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/utils"
)

func HandlePost(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	return events.APIGatewayProxyResponse{}, nil
}

func HandleGetByName(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	return events.APIGatewayProxyResponse{}, nil
}
