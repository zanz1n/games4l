package src

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/question/entityconv"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/httpcodes"
	"github.com/games4l/pkg/errors"
)

func HandleGetOne(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	id, err := strconv.Atoi(req.PathParameters["id"])
	if err != nil || id > math.MaxInt32 {
		return nil, errors.ErrInvalidIntegerIdPathParam
	}

	if err = Connect(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q, err := dba.GetQuestionById(ctx, int32(id))
	if err != nil || q == nil {
		return nil, errors.ErrEntityNotFound
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: utils.MarshalJSON(JSON{
			"message": "Success",
			"data":    entityconv.QuestionToApiEntity(q),
		}),
	}, nil
}
