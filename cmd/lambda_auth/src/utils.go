package src

import (
	"strings"

	"github.com/games4l/internal/auth"
	"github.com/games4l/internal/errors"
	"github.com/goccy/go-json"

	"github.com/games4l/internal/user"
	"github.com/go-playground/validator/v10"
)

var (
	applicationJsonHeader = map[string]string{
		"Content-Type": "application/json",
	}
	validate = validator.New()
	ap       *auth.AuthProvider
	dba      *user.UserService
)

type JSON map[string]interface{}

func MarshalJSON(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "{\"error\":\"the real message to be shown could not be encoded, this is not the intended one\"}"
	}

	return string(bytes)
}

func AuthBySig(header string, body string) errors.StatusCodeErr {
	authHeaderS := strings.Split(header, " ")

	if len(authHeaderS) < 3 {
		return errors.DefaultErrorList.RouteRequiresAdminAuth
	}

	if authHeaderS[0] != "Signature" {
		return errors.DefaultErrorList.InvalidAuthStrategy
	}

	encodingS := auth.ByteEncoding(authHeaderS[1])

	if encodingS != auth.ByteEncodingBase64 && encodingS != auth.ByteEncodingHex {
		return errors.DefaultErrorList.InvalidAuthSignatureEncodingMethod
	}

	err := ap.ValidateSignature(encodingS, []byte(body), []byte(authHeaderS[2]))

	if err != nil {
		return err
	}

	return nil
}
