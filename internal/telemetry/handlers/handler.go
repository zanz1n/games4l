package handlers

import (
	log "github.com/games4l/internal/logger"
	"github.com/games4l/internal/telemetry/repository"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func NewTelemetryHandlers(ts repository.Singleton) *TelemetryHandlers {
	return &TelemetryHandlers{
		ts:     ts,
		logger: log.NewLogger("telemetry_handlers"),
	}
}

type TelemetryHandlers struct {
	ts     repository.Singleton
	logger log.Logger
}
