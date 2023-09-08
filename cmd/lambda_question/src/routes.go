package src

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/errors"
	"github.com/games4l/internal/httpcodes"
	"github.com/games4l/internal/logger"
	"github.com/games4l/internal/question"
	"github.com/games4l/internal/utils"
	"github.com/goccy/go-json"
)

func HandleGetMany(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	limitParam, ok := req.QueryStringParameters["limit"]

	limit := 1000

	if ok {
		l, err := strconv.Atoi(limitParam)
		if err != nil || l <= 0 {
			return nil, errors.ErrInvalidRequestEntity
		}
		limit = l
	}

	if err := Connect(); err != nil {
		logger.Error("Connect call failed: " + err.Error())
		return nil, errors.ErrInternalServerError
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
		Body: utils.MarshalJSON(JSON{
			"message": "success",
			"data":    result,
		}),
	}, nil
}

func HandlePost(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	authHeader, ok := req.Headers["authorization"]

	if !ok {
		return nil, errors.ErrRouteRequiresAdminAuth
	}

	if err := ap.AuthenticateAdminHeader(
		authHeader,
		utils.S2B(req.Body),
	); err != nil {
		return nil, err
	}

	qfmt, ok := req.QueryStringParameters["fmt"]

	qb := &question.Question{}

	if !ok || qfmt == "new" {
		bodyP := question.Question{}

		if err := json.Unmarshal([]byte(req.Body), &bodyP); err != nil {
			return nil, errors.ErrMalformedOrTooBigBody
		}

		if err := validate.Struct(&bodyP); err != nil || !bodyP.IsValid() {
			return nil, errors.ErrInvalidRequestEntity
		}

		qb = &bodyP
	} else if qfmt == "old" {
		bodyP := question.QuestionAlternativeFmt{}

		if err := json.Unmarshal([]byte(req.Body), &bodyP); err != nil {
			return nil, errors.ErrMalformedOrTooBigBody
		}

		if err := validate.Struct(bodyP); err != nil || !bodyP.IsValid() {
			return nil, errors.ErrInvalidRequestEntity
		}

		qb = bodyP.Parse()
	} else {
		return nil, errors.ErrInvalidFMTQueryParam
	}

	nidParam := req.QueryStringParameters["nid"]

	nid, err := strconv.Atoi(nidParam)

	if err != nil || 0 >= nid {
		return nil, errors.ErrInvalidNIDQueryParam
	}

	if err := Connect(); err != nil {
		logger.Error("Connect call failed: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, fErr := dba.Create(ctx, nid, qb)

	if fErr != nil {
		return nil, fErr
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: utils.MarshalJSON(JSON{
			"message": "created",
			"data":    result,
		}),
	}, nil
}

func HandleGetByID(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
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
			return nil, errors.ErrInvalidRequestEntity
		}

		id, err := strconv.Atoi(idParam)

		if err != nil {
			return nil, errors.ErrInvalidNIDQueryParam
		}

		if err := Connect(); err != nil {
			logger.Error("Connect call failed: " + err.Error())
			return nil, errors.ErrInternalServerError
		}

		result, fErr := dba.GetByNumID(ctx, id)

		if fErr != nil {
			return nil, fErr
		}

		res.Body = utils.MarshalJSON(JSON{
			"message": "success",
			"data":    result,
		})

		return &res, nil
	}

	if err := Connect(); err != nil {
		logger.Error("Connect call failed: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	result, err := dba.GetByID(ctx, idParam)

	if err != nil {
		return nil, err
	}

	res.Body = utils.MarshalJSON(JSON{
		"message": "success",
		"data":    result,
	})

	return &res, nil
}
