package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/telemetry"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
	"github.com/go-playground/validator"
	"github.com/goccy/go-json"
)

var (
	applicationJsonHeader = map[string]string{
		"Content-Type": "application/json",
	}
	validate = validator.New()
)

type JSON map[string]interface{}

func HandlePost(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	telemetryData := telemetry.CreateTelemetryUnitData{}

	json.Unmarshal([]byte(req.Body), &telemetryData)

	err := validate.Struct(telemetryData)

	if err != nil {
		return nil, utils.NewStatusCodeErr("malformed body", httpcodes.StatusBadGateway)
	}

	result, fErr := dba.Create(&telemetryData)

	if fErr != nil {
		return nil, fErr
	}

	resultEnc, err := json.Marshal(JSON{
		"message": "data created",
		"data":    result,
	})

	if err != nil {
		return nil, utils.NewStatusCodeErr(
			"internal server error when creating the data, try again later",
			httpcodes.StatusInternalServerError,
		)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusCreated,
		Headers:         applicationJsonHeader,
		Body:            string(resultEnc),
		IsBase64Encoded: false,
	}, nil
}

func HandleGetByName(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusMethodNotAllowed,
		Body:            "{\"error\":\"method not allowed\"}",
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
	}, nil
}
