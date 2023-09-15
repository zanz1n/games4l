package src

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/logger"
	"github.com/games4l/internal/telemetry"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/httpcodes"
	"github.com/games4l/pkg/errors"
	"github.com/goccy/go-json"
)

func HandlePost(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	telemetryData := telemetry.CreateTelemetryUnitData{}

	json.Unmarshal([]byte(req.Body), &telemetryData)

	err := validate.Struct(telemetryData)

	if err != nil {
		return nil, errors.ErrMalformedOrTooBigBody
	}

	Connect()

	result, fErr := dba.Create(&telemetryData)

	if fErr != nil {
		return nil, fErr
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: httpcodes.StatusCreated,
		Headers:    applicationJsonHeader,
		Body: utils.MarshalJSON(JSON{
			"message": "data created",
			"data":    result,
		}),
		IsBase64Encoded: false,
	}, nil
}

func HandleGetByID(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	idParam, ok := req.PathParameters["id"]

	if !ok || idParam == "" {
		return nil, errors.ErrInvalidRequestEntity
	}

	if err := Connect(); err != nil {
		logger.Error("Connect call failed: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	result, err := dba.FindByIdWithCtx(ctx, idParam)

	if err != nil {
		return nil, err
	}
	err = ap.AuthenticateAdminHeader(req.Headers["authorization"], utils.S2B(req.Body))
	if err != nil {
		result.PacientName = "<OMITTED>"
	}

	res := events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: utils.MarshalJSON(JSON{
			"data":    result,
			"message": "",
		}),
	}

	return &res, nil
}

func HandleGetByName(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if err := ap.AuthenticateAdminHeader(
		req.Headers["authorization"],
		utils.S2B(req.Body),
	); err != nil {
		return nil, err
	}

	nameParam, ok := req.QueryStringParameters["name"]

	if !ok {
		return nil, errors.ErrInvalidRequestEntity
	}

	if err := Connect(); err != nil {
		logger.Error("Connect call failed: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

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
		Body: utils.MarshalJSON(JSON{
			"message": "success",
			"data":    result,
		}),
	}, nil
}
