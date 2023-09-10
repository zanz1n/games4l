package src

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/entity/entityconv"
	"github.com/games4l/internal/entity/question"
	"github.com/games4l/internal/errors"
	"github.com/games4l/internal/httpcodes"
	"github.com/games4l/internal/utils"
	"github.com/goccy/go-json"
)

func HandlePostOne(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	buf := utils.S2B(req.Body)

	err := ap.AuthenticateAdminHeader(req.Headers["authorization"], buf)
	if err != nil {
		return nil, err
	}

	body := question.QuestionCreateData{}

	if err = json.Unmarshal(buf, &body); err != nil {
		return nil, errors.ErrMalformedOrTooBigBody
	}
	if err = validate.Struct(&body); err != nil {
		return nil, errors.ErrInvalidRequestEntity
	}
	if !body.IsValid() {
		return nil, errors.ErrInvalidRequestEntity
	}

	if err = Connect(); err != nil {
		return nil, errors.ErrInternalServerError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	creation, err := dba.CreateQuestion(ctx, entityconv.CreateQuestionParamsToDbEntity(&body))
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: utils.MarshalJSON(JSON{
			"message": "Created successfully",
			"data":    entityconv.QuestionToApiEntity(creation),
		}),
	}, nil
}
