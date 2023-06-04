package main

import (
	"github.com/goccy/go-json"

	"github.com/go-playground/validator"
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
