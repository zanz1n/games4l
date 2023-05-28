package main

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/auth"
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
		return nil, utils.NewStatusCodeErr("malformed body", httpcodes.StatusBadGateway)
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

func HandleGetByName(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	authHeader, ok := req.Headers["authorization"]

	if !ok {
		return nil, utils.NewStatusCodeErr(
			"this route requires the 'Authorization' header",
			httpcodes.StatusBadRequest,
		)
	}

	authHeaderS := strings.Split(authHeader, " ")

	if len(authHeaderS) < 3 {
		return nil, utils.NewStatusCodeErr(
			"this route requires admin authorization",
			httpcodes.StatusBadRequest,
		)
	}

	if authHeaderS[0] != "Signature" {
		return nil, utils.NewStatusCodeErr(
			"invalid auth strategy "+authHeaderS[0], httpcodes.StatusBadRequest,
		)
	}

	encodingS := auth.ByteEncoding(authHeaderS[1])

	if encodingS != auth.ByteEncodingBase64 && encodingS != auth.ByteEncodingHex {
		return nil, utils.NewStatusCodeErr(
			"invalid encoding strategy "+authHeaderS[1],
			httpcodes.StatusBadRequest,
		)
	}

	err := ap.ValidateSignature(encodingS, []byte(req.Body), []byte(authHeaderS[2]))

	if err != nil {
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
