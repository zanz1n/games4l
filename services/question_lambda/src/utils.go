package src

import (
	"strings"

	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/question"
	"github.com/games4l/backend/libs/utils"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
)

var (
	applicationJsonHeader = map[string]string{
		"Content-Type": "application/json",
	}
	validate = validator.New()
	ap       *auth.AuthProvider
	dba      *question.QuestionService
)

type JSON map[string]interface{}

func MarshalJSON(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "{\"error\":\"the real message to be shown could not be encoded, this is not the intended one\"}"
	}

	return string(bytes)
}

func AuthAdmin(header string, body string) utils.StatusCodeErr {
	authHeaderS := strings.Split(header, " ")
	l := len(authHeaderS)

	if l < 2 || l > 3 {
		return utils.DefaultErrorList.RouteRequiresAdminAuth
	}

	if authHeaderS[0] == "Signature" {
		if err := AuthBySig(header, body); err != nil {
			return err
		}
	} else if authHeaderS[0] == "Bearer" {
		jwtToken := authHeaderS[1]

		// This may change depending on the jwt algorithm used
		if len(jwtToken) < 100 || l > 2 {
			return utils.DefaultErrorList.InvalidJwtTokenFormat
		}

		user, err := ap.AuthUser(jwtToken)

		if err != nil {
			return err
		}

		if user.Role != auth.UserRoleAdmin {
			return utils.DefaultErrorList.RouteRequiresAdminAuth
		}
	} else {
		return utils.DefaultErrorList.InvalidAuthStrategy
	}

	return nil
}

func AuthBySig(header string, body string) utils.StatusCodeErr {
	authHeaderS := strings.Split(header, " ")

	if len(authHeaderS) < 3 {
		return utils.DefaultErrorList.RouteRequiresAdminAuth
	}

	if authHeaderS[0] != "Signature" {
		return utils.DefaultErrorList.InvalidAuthStrategy
	}

	encodingS := auth.ByteEncoding(authHeaderS[1])

	if encodingS != auth.ByteEncodingBase64 && encodingS != auth.ByteEncodingHex {
		return utils.DefaultErrorList.InvalidAuthSignatureEncodingMethod
	}

	err := ap.ValidateSignature(encodingS, []byte(body), []byte(authHeaderS[2]))

	if err != nil {
		return err
	}

	return nil
}
