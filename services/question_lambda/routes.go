package main

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
)

func HandleGetMany(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	limitParam, ok := req.QueryStringParameters["limit"]

	limit := 1000

	if ok {
		l, err := strconv.Atoi(limitParam)
		if err != nil || l <= 0 {
			return nil, utils.NewStatusCodeErr(
				"limit query param must be null or a valid unsigned integer",
				httpcodes.StatusBadRequest,
			)
		}
		limit = l
	}

	if err := Connect(); err != nil {
		return nil, utils.NewStatusCodeErr(
			"failed to connect to database",
			httpcodes.StatusInternalServerError,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	result, err := dba.GetMany(ctx, int64(limit))

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

func HandleGetByID(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	idParam, ok := req.QueryStringParameters["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
	}

	if !ok {
		idParam, ok = req.QueryStringParameters["nid"]

		if !ok {
			return nil, utils.NewStatusCodeErr(
				"the id/nid query parameter must be provided",
				httpcodes.StatusBadRequest,
			)
		}

		id, err := strconv.Atoi(idParam)

		if err != nil {
			return nil, utils.NewStatusCodeErr(
				"the nid query parameter must be a valid integer",
				httpcodes.StatusBadRequest,
			)
		}

		if err = Connect(); err != nil {
			return nil, utils.NewStatusCodeErr(
				"failed to connect to the database",
				httpcodes.StatusInternalServerError,
			)
		}

		result, fErr := dba.GetByNumID(ctx, id)

		if fErr != nil {
			return nil, fErr
		}

		res.Body = MarshalJSON(JSON{
			"message": "success",
			"data":    result,
		})

		return &res, nil
	}

	if err := Connect(); err != nil {
		return nil, utils.NewStatusCodeErr(
			"failed to connect to the database",
			httpcodes.StatusInternalServerError,
		)
	}

	result, err := dba.GetByID(ctx, idParam)

	if err != nil {
		return nil, err
	}

	res.Body = MarshalJSON(JSON{
		"message": "success",
		"data":    result,
	})

	return &res, nil
}
