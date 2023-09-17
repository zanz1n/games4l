package repository

import (
	log "github.com/games4l/internal/logger"
	"github.com/games4l/pkg/errors"
)

type Singleton interface {
	GetInstance() (QuestionRepository, error)
}

func NewPostgresSingleton(url string) Singleton {
	return &questionRepositorySingleton{
		r:      nil,
		url:    url,
		logger: log.NewLogger("question_singleton"),
	}
}

type questionRepositorySingleton struct {
	r      QuestionRepository
	url    string
	logger log.Logger
}

func (s *questionRepositorySingleton) GetInstance() (QuestionRepository, error) {
	if s.r != nil {
		return s.r, nil
	}

	var err error
	s.r, err = NewPostgres(s.url)
	if err != nil {
		s.logger.Error("Failed to connect to postgres: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	return s.r, nil
}
