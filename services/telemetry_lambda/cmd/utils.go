package cmd

import (
	"strings"

	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/services/telemetry_lambda/repository"
	"github.com/go-playground/validator/v10"
)

var (
	applicationJsonHeader = map[string]string{
		"Content-Type": "application/json",
	}
	validate = validator.New()
	dba      *repository.TelemetryService
	ap       *auth.AuthProvider
)

type JSON map[string]interface{}

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
