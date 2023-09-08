package src

import (
	"strings"

	"github.com/games4l/internal/auth"
	"github.com/games4l/internal/telemetry"
	"github.com/games4l/internal/errors"
	"github.com/go-playground/validator/v10"
)

var (
	applicationJsonHeader = map[string]string{
		"Content-Type": "application/json",
	}
	validate = validator.New()
	dba      *telemetry.TelemetryService
	ap       *auth.AuthProvider
)

type JSON map[string]interface{}

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
