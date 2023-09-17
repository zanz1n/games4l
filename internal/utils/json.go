package utils

import "github.com/goccy/go-json"

type DataResponse[T any] struct {
	Message    string `json:"message"`
	Data       T      `json:"data"`
	StatusCode int    `json:"-"`
}

func MarshalJSON(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "{\"error\":\"the real message to be shown could not be encoded, this is not the intended one\"}"
	}

	return string(bytes)
}
