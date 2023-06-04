package main

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
	"github.com/goccy/go-json"
)

type SigInBody struct {
	Credential string `json:"credential,omitempty" validate:"required"`
	Password   string `json:"password,omitempty" validate:"required"`
}

func HandleSignIn(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	bodyP := SigInBody{}

	if err := json.Unmarshal([]byte(req.Body), &bodyP); err != nil {
		return nil, utils.NewStatusCodeErr(
			"failed to decode body",
			httpcodes.StatusBadRequest,
		)
	}

	if err := validate.Struct(bodyP); err != nil {
		return nil, utils.NewStatusCodeErr(
			"invalid body format",
			httpcodes.StatusBadRequest,
		)
	}

	if err := Connect(); err != nil {
		logger.Error(err.Error())

		return nil, utils.NewStatusCodeErr(
			"failed to connect to database",
			httpcodes.StatusInternalServerError,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, err := dba.SignInUser(ctx, bodyP.Credential, bodyP.Password)

	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: MarshalJSON(JSON{
			"message": "success",
			"token":   token,
		}),
	}, nil
}
