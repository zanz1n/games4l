package sqli

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createQuestion = `-- name: CreateQuestion :one
INSERT INTO
    "question" (
        "question",
        "answer_1",
        "answer_2",
        "answer_3",
        "answer_4",
        "correct_answer",
        "type",
        "style",
        "difficulty"
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at, question, answer_1, answer_2, answer_3, answer_4, correct_answer, type, style, difficulty
`

type CreateQuestionParams struct {
	Question      string
	Answer1       string
	Answer2       string
	Answer3       pgtype.Text
	Answer4       pgtype.Text
	CorrectAnswer int16
	Type          QuestionType
	Style         QuestionStyle
	Difficulty    int16
}

func (q *Queries) CreateQuestion(ctx context.Context, arg *CreateQuestionParams) (*Question, error) {
	row := q.db.QueryRow(ctx, createQuestion,
		arg.Question,
		arg.Answer1,
		arg.Answer2,
		arg.Answer3,
		arg.Answer4,
		arg.CorrectAnswer,
		arg.Type,
		arg.Style,
		arg.Difficulty,
	)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Question,
		&i.Answer1,
		&i.Answer2,
		&i.Answer3,
		&i.Answer4,
		&i.CorrectAnswer,
		&i.Type,
		&i.Style,
		&i.Difficulty,
	)
	return &i, err
}

const getManyQuestions = `-- name: GetManyQuestions :many
SELECT id, created_at, updated_at, question, answer_1, answer_2, answer_3, answer_4, correct_answer, type, style, difficulty FROM "question" ORDER BY "id" LIMIT $1
`

func (q *Queries) GetManyQuestions(ctx context.Context, limit int32) ([]*Question, error) {
	rows, err := q.db.Query(ctx, getManyQuestions, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Question{}
	for rows.Next() {
		var i Question
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Question,
			&i.Answer1,
			&i.Answer2,
			&i.Answer3,
			&i.Answer4,
			&i.CorrectAnswer,
			&i.Type,
			&i.Style,
			&i.Difficulty,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQuestionById = `-- name: GetQuestionById :one
SELECT id, created_at, updated_at, question, answer_1, answer_2, answer_3, answer_4, correct_answer, type, style, difficulty FROM "question" WHERE "id" = $1
`

func (q *Queries) GetQuestionById(ctx context.Context, id int32) (*Question, error) {
	row := q.db.QueryRow(ctx, getQuestionById, id)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Question,
		&i.Answer1,
		&i.Answer2,
		&i.Answer3,
		&i.Answer4,
		&i.CorrectAnswer,
		&i.Type,
		&i.Style,
		&i.Difficulty,
	)
	return &i, err
}

const updateQuestionById = `-- name: UpdateQuestionById :one
UPDATE "question"
SET
    "updated_at" = CURRENT_TIMESTAMP,
    "question" = $1,
    "answer_1" = $2,
    "answer_2" = $3,
    "answer_3" = $4,
    "answer_4" = $5,
    "correct_answer" = $6,
    "type" = $7
WHERE "id" = $8 RETURNING id, created_at, updated_at, question, answer_1, answer_2, answer_3, answer_4, correct_answer, type, style, difficulty
`

type UpdateQuestionByIdParams struct {
	Question      string
	Answer1       string
	Answer2       string
	Answer3       pgtype.Text
	Answer4       pgtype.Text
	CorrectAnswer int16
	Type          QuestionType
	ID            int32
}

func (q *Queries) UpdateQuestionById(ctx context.Context, arg *UpdateQuestionByIdParams) (*Question, error) {
	row := q.db.QueryRow(ctx, updateQuestionById,
		arg.Question,
		arg.Answer1,
		arg.Answer2,
		arg.Answer3,
		arg.Answer4,
		arg.CorrectAnswer,
		arg.Type,
		arg.ID,
	)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Question,
		&i.Answer1,
		&i.Answer2,
		&i.Answer3,
		&i.Answer4,
		&i.CorrectAnswer,
		&i.Type,
		&i.Style,
		&i.Difficulty,
	)
	return &i, err
}
