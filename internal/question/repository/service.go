package repository

import (
	"context"
	"time"

	log "github.com/games4l/internal/logger"
	"github.com/games4l/internal/question"
	"github.com/games4l/internal/sqli"
	"github.com/games4l/pkg/errors"
	"github.com/jackc/pgx/v5"
)

func Connect(url string) (QuestionRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}

	return newQuestionService(sqli.New(conn)), nil
}

func newQuestionService(dba *sqli.Queries) *questionService {
	return &questionService{
		dba:    dba,
		logger: log.NewLogger("question_repository"),
	}
}

type questionService struct {
	dba    *sqli.Queries
	logger log.Logger
}

func (s *questionService) GetById(id int32) (*question.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q, err := s.dba.GetQuestionById(ctx, id)
	if err != nil {
		s.logger.Info("Failed to get by id: " + err.Error())
		return nil, errors.ErrEntityNotFound
	}

	return questionToApiEntity(q), nil
}

func (s *questionService) Create(
	data *question.QuestionCreateData,
) (*question.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q, err := s.dba.CreateQuestion(ctx, createQuestionParamsToDbEntity(data))
	if err != nil {
		s.logger.Error("Failed to create: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	return questionToApiEntity(q), nil
}

func (s *questionService) GetMany(limit int32) ([]question.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.dba.GetManyQuestions(ctx, limit)
	if err != nil {
		s.logger.Error("Failed to get many: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	nq := make([]question.Question, len(res))

	for i, q := range res {
		nq[i] = *questionToApiEntity(q)
	}

	return nq, nil
}

func (s *questionService) UpdateById(
	id int32,
	data *question.QuestionUpdateData,
) (*question.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q, err := s.dba.UpdateQuestionById(ctx, updateQuestionParamsToDbEntity(id, data))
	if err != nil {
		s.logger.Info("Failed to update by id: " + err.Error())
		return nil, errors.ErrEntityNotFound
	}

	return questionToApiEntity(q), nil
}

func (s *questionService) DeleteById(id int32) (*question.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q, err := s.dba.DeleteQuestionById(ctx, id)
	if err != nil {
		s.logger.Info("Failed to delete by id: " + err.Error())
		return nil, errors.ErrEntityNotFound
	}

	return questionToApiEntity(q), nil
}
