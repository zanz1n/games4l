package sqli

import (
	"context"
)

type Querier interface {
	CreateQuestion(ctx context.Context, arg *CreateQuestionParams) (*Question, error)
	GetManyQuestions(ctx context.Context, limit int32) ([]*Question, error)
	GetQuestionById(ctx context.Context, id int32) (*Question, error)
	UpdateQuestionById(ctx context.Context, arg *UpdateQuestionByIdParams) (*Question, error)
}

var _ Querier = (*Queries)(nil)
