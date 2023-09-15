package src

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/auth"
	"github.com/games4l/internal/logger"
	"github.com/games4l/internal/user"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/httpcodes"
	"github.com/games4l/pkg/errors"
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

func HandleSignIn(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	bodyP := SigInBody{}

	if err := json.Unmarshal([]byte(req.Body), &bodyP); err != nil {
		return nil, errors.ErrMalformedOrTooBigBody
	}

	if err := validate.Struct(bodyP); err != nil {
		return nil, errors.ErrInvalidRequestEntity
	}

	if err := Connect(); err != nil {
		logger.Error("Connect call failed: " + err.Error())
		return nil, errors.ErrInternalServerError
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
		Body: utils.MarshalJSON(JSON{
			"message": "success",
			"token":   token,
		}),
	}, nil
}

func HandleUserCreation(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		cErr error
		err  error
	)
	bodyP := CreateUserBody{}

	if err = json.Unmarshal([]byte(req.Body), &bodyP); err != nil {
		return nil, errors.ErrMalformedOrTooBigBody
	}

	if err = validate.Struct(bodyP); err != nil {
		return nil, errors.ErrMalformedOrTooBigBody
	}

	cErr = ap.AuthenticateAdminHeader(req.Headers["authorization"], utils.S2B(req.Body))
	if err != nil {
		return nil, cErr
	}

	if err = Connect(); err != nil {
		logger.Error("Connect call failed: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, cErr := dba.CreateUser(ctx, bodyP.Role, &user.CreateUserData{
		Username: bodyP.Username,
		Email:    bodyP.Email,
		Password: bodyP.Password,
	})

	if cErr != nil {
		return nil, cErr
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: utils.MarshalJSON(JSON{
			"message": "success",
			"data":    result.ToEncodable(),
		}),
	}, nil
}
