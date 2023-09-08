package src

import (
	"github.com/games4l/internal/auth"
	"github.com/games4l/internal/telemetry"
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
