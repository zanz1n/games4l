package main

import (
	"strings"

	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
	"github.com/go-playground/validator"
	"github.com/goccy/go-json"
)

var (
	applicationJsonHeader = map[string]string{
		"Content-Type": "application/json",
	}
	validate = validator.New()
)

type JSON map[string]interface{}

func MarshalJSON(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "{\"error\":\"the real message to be shown could not be encoded, this is not the intended one\"}"
	}

	return string(bytes)
}

func AuthBySig(header string, body string) utils.StatusCodeErr {
	authHeaderS := strings.Split(header, " ")

	if len(authHeaderS) < 3 {
		return utils.NewStatusCodeErr(
			"this route requires admin authorization",
			httpcodes.StatusBadRequest,
		)
	}

	if authHeaderS[0] != "Signature" {
		return utils.NewStatusCodeErr(
			"invalid auth strategy "+authHeaderS[0], httpcodes.StatusBadRequest,
		)
	}

	encodingS := auth.ByteEncoding(authHeaderS[1])

	if encodingS != auth.ByteEncodingBase64 && encodingS != auth.ByteEncodingHex {
		return utils.NewStatusCodeErr(
			"invalid encoding strategy "+authHeaderS[1],
			httpcodes.StatusBadRequest,
		)
	}

	err := ap.ValidateSignature(encodingS, []byte(body), []byte(authHeaderS[2]))

	if err != nil {
		return err
	}

	return nil
}
