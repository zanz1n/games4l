package main

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/question"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
	"github.com/goccy/go-json"
)

func HandleGetMany(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	limitParam, ok := req.QueryStringParameters["limit"]

	limit := 1000

	if ok {
		l, err := strconv.Atoi(limitParam)
		if err != nil || l <= 0 {
			return nil, utils.DefaultErrorList.InvalidRequestEntity
		}
		limit = l
	}

	if err := Connect(); err != nil {
		return nil, utils.DefaultErrorList.InternalServerError
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

func HandlePost(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	authHeader, ok := req.Headers["authorization"]

	if !ok {
		return nil, utils.DefaultErrorList.RouteRequiresAdminAuth
	}

	if err := AuthAdmin(authHeader, req.Body); err != nil {
		return nil, err
	}

	qfmt, ok := req.QueryStringParameters["fmt"]

	qb := &question.Question{}

	if !ok || qfmt == "new" {
		bodyP := question.Question{}

		if err := json.Unmarshal([]byte(req.Body), &bodyP); err != nil {
			return nil, utils.DefaultErrorList.MalformedOrTooBigBody
		}

		if err := validate.Struct(bodyP); err != nil || !bodyP.IsValid() {
			return nil, utils.DefaultErrorList.InvalidRequestEntity
		}

		qb = &bodyP
	} else if qfmt == "old" {
		bodyP := question.QuestionAlternativeFmt{}

		if err := json.Unmarshal([]byte(req.Body), &bodyP); err != nil {
			return nil, utils.DefaultErrorList.MalformedOrTooBigBody
		}

		if err := validate.Struct(bodyP); err != nil || !bodyP.IsValid() {
			return nil, utils.DefaultErrorList.InvalidRequestEntity
		}

		qb = bodyP.Parse()
	} else {
		return nil, utils.DefaultErrorList.InvalidFMTQueryParam
	}

	nidParam := req.QueryStringParameters["nid"]

	nid, err := strconv.Atoi(nidParam)

	if err != nil || 0 >= nid {
		return nil, utils.DefaultErrorList.InvalidNIDQueryParam
	}

	if err := Connect(); err != nil {
		return nil, utils.DefaultErrorList.InternalServerError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, fErr := dba.Create(ctx, nid, qb)

	if err != nil {
		return nil, fErr
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: MarshalJSON(JSON{
			"message": "created",
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
			return nil, utils.DefaultErrorList.InvalidRequestEntity
		}

		id, err := strconv.Atoi(idParam)

		if err != nil {
			return nil, utils.DefaultErrorList.InvalidNIDQueryParam
		}

		if err = Connect(); err != nil {
			return nil, utils.DefaultErrorList.InternalServerError
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
		return nil, utils.DefaultErrorList.InternalServerError
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
