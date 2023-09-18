package repository

import (
	"github.com/games4l/internal/telemetry"
)

type TelemetryRepository interface {
	GetById(id string) (*telemetry.Registry, error)
	Create(data *telemetry.CreateRegistryData) (*telemetry.Registry, error)
	GetBySimilarName(name string) ([]telemetry.Registry, error)
	DeleteById(id string) (*telemetry.Registry, error)
}
