package src

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/entity/entityconv"
	"github.com/games4l/internal/entity/question"
	"github.com/games4l/internal/errors"
	"github.com/games4l/internal/httpcodes"
	"github.com/games4l/internal/utils"
)

func HandleGetMany(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	err := Connect()
	if err != nil {
		return nil, err
	}

	var limit int
	limitS, ok := req.PathParameters["l"]
	if ok {
		if limit, err = strconv.Atoi(limitS); err != nil {
			limit = 1000
		} else {
			if limit > 1000 || limit < 10 {
				limit = 1000
			}
		}
	} else {
		limit = 1000
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbQuests, err := dba.GetManyQuestions(ctx, int32(limit))
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	newArr := make([]*question.Question, len(dbQuests))

	for i, q := range dbQuests {
		clone := *q
		newArr[i] = entityconv.QuestionToApiEntity(&clone)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: utils.MarshalJSON(JSON{
			"message": "Success",
			"data":    newArr,
		}),
	}, nil
}
