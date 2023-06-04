package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/utils"
)

func HandleSignIn(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	return nil, utils.NewStatusCodeErr("not implemented (TODO)", 500)
}

func HandleSignUp(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	return nil, utils.NewStatusCodeErr("not implemented (TODO)", 500)
}
