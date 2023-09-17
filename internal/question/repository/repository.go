package repository

import (
	"github.com/games4l/internal/question"
)

type QuestionRepository interface {
	GetById(id int32) (*question.Question, error)
	Create(data *question.QuestionCreateData) (*question.Question, error)
	GetMany(limit int32) ([]question.Question, error)
	UpdateById(id int32, data *question.QuestionUpdateData) (*question.Question, error)
	DeleteById(id int32) (*question.Question, error)
}
