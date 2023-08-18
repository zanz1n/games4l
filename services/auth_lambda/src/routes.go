package src

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/user"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
	"github.com/goccy/go-json"
)

type SigInBody struct {
	Credential string `json:"credential,omitempty" validate:"required"`
	Password   string `json:"password,omitempty" validate:"required"`
}

type CreateUserBody struct {
	Username string        `json:"username,omitempty" validate:"required"`
	Email    string        `json:"email,omitempty" validate:"required"`
	Password string        `json:"password,omitempty" validate:"required"`
	Role     auth.UserRole `json:"role,omitempty" validate:"required"`
}

func HandleSignIn(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	bodyP := SigInBody{}

	if err := json.Unmarshal([]byte(req.Body), &bodyP); err != nil {
		return nil, utils.DefaultErrorList.MalformedOrTooBigBody
	}

	if err := validate.Struct(bodyP); err != nil {
		return nil, utils.DefaultErrorList.InvalidRequestEntity
	}

	if err := Connect(); err != nil {
		logger.Error(err.Error())

		return nil, utils.DefaultErrorList.InternalServerError
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

func HandleUserCreation(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, utils.StatusCodeErr) {
	bodyP := CreateUserBody{}

	if err := json.Unmarshal([]byte(req.Body), &bodyP); err != nil {
		return nil, utils.DefaultErrorList.MalformedOrTooBigBody
	}

	if err := validate.Struct(bodyP); err != nil {
		return nil, utils.DefaultErrorList.MalformedOrTooBigBody
	}

	if err := AuthBySig(req.Headers["authorization"], req.Body); err != nil {
		return nil, err
	}

	if err := Connect(); err != nil {
		logger.Error(err.Error())

		return nil, utils.DefaultErrorList.InternalServerError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := dba.CreateUser(ctx, bodyP.Role, &user.CreateUserData{
		Username: bodyP.Username,
		Email:    bodyP.Email,
		Password: bodyP.Password,
	})

	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: MarshalJSON(JSON{
			"message": "success",
			"data":    result.ToEncodable(),
		}),
	}, nil
}
