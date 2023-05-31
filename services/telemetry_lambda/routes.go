package main

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/telemetry"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
	"github.com/goccy/go-json"
)

func HandlePost(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	telemetryData := telemetry.CreateTelemetryUnitData{}

	json.Unmarshal([]byte(req.Body), &telemetryData)

	err := validate.Struct(telemetryData)

	if err != nil {
		return nil, utils.NewStatusCodeErr("malformed body", httpcodes.StatusBadRequest)
	}

	Connect()

	result, fErr := dba.Create(&telemetryData)

	if fErr != nil {
		return nil, fErr
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: httpcodes.StatusCreated,
		Headers:    applicationJsonHeader,
		Body: MarshalJSON(JSON{
			"message": "data created",
			"data":    result,
		}),
		IsBase64Encoded: false,
	}, nil
}

func HandleGetByID(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	idParam, ok := req.PathParameters["id"]

	if !ok || idParam == "" {
		return nil, utils.NewStatusCodeErr(
			"id path param not provided",
			httpcodes.StatusBadRequest,
		)
	}

	Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	result, err := dba.FindByIdWithCtx(ctx, idParam)

	if err != nil {
		return nil, err
	}

	if err = AuthBySig(req.Headers["authorization"], req.Body); err != nil {
		result.PacientName = "<OMITTED>"
	}

	res := events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: MarshalJSON(JSON{
			"data":    result,
			"message": "",
		}),
	}

	return &res, nil
}

func HandleGetByName(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	if err := AuthBySig(req.Headers["authorization"], req.Body); err != nil {
		return nil, err
	}

	nameParam, ok := req.QueryStringParameters["name"]

	if !ok {
		return nil, utils.NewStatusCodeErr(
			"name query param must be provided",
			httpcodes.StatusBadRequest,
		)
	}

	Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	result, err := dba.FindSimilarNameWithCtx(ctx, nameParam)

	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: MarshalJSON(JSON{
			"message": "success",
			"data":    result,
		}),
	}, nil
}
