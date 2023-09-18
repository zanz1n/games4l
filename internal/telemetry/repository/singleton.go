package repository

import (
	log "github.com/games4l/internal/logger"
	"github.com/games4l/pkg/errors"
)

func NewMongoSingleton(url, dbname string) Singleton {
	return &telemetryRepositorySingleton{
		r:      nil,
		url:    url,
		dbname: dbname,
		logger: log.NewLogger("telemetry_singleton"),
	}
}

type Singleton interface {
	GetInstance() (TelemetryRepository, error)
}

type telemetryRepositorySingleton struct {
	r      TelemetryRepository
	url    string
	dbname string
	logger log.Logger
}

func (s *telemetryRepositorySingleton) GetInstance() (TelemetryRepository, error) {
	if s.r != nil {
		return s.r, nil
	}

	var err error
	s.r, err = NewMongo(s.url, s.dbname)
	if err != nil {
		s.logger.Error("Failed to connect to mongodb: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	return s.r, nil
}
