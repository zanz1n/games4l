package repository

import (
	log "github.com/games4l/internal/logger"
	"github.com/games4l/pkg/errors"
)

func NewSingleton(url string) *QuestionRepositorySingleton {
	return &QuestionRepositorySingleton{
		r:      nil,
		url:    url,
		logger: log.NewLogger("question_singleton"),
	}
}

type QuestionRepositorySingleton struct {
	r      QuestionRepository
	url    string
	logger log.Logger
}

func (s *QuestionRepositorySingleton) GetInstance() (QuestionRepository, error) {
	if s.r != nil {
		return s.r, nil
	}

	var err error
	s.r, err = Connect(s.url)
	if err != nil {
		s.logger.Error("Failed to connect to postgres: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	return s.r, nil
}
