package cmd

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type JSON map[string]any

func identJson(buf []byte) []byte {
	var err error

	m := make(map[string]any)

	if err = json.Unmarshal(buf, &m); err != nil {
		return buf
	}

	var newBuf []byte

	newBuf, err = json.MarshalIndent(m, "", "  ")
	if err != nil {
		return buf
	}

	return newBuf
}
