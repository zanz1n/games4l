package repository

import (
	"github.com/games4l/internal/question"
	"github.com/games4l/internal/sqli"
	"github.com/jackc/pgx/v5/pgtype"
)

func updateQuestionParamsToDbEntity(id int32, q *question.QuestionUpdateData) *sqli.UpdateQuestionByIdParams {
	nq := sqli.UpdateQuestionByIdParams{
		ID:            id,
		Question:      q.Question,
		Answer1:       q.Answers[0],
		Answer2:       q.Answers[1],
		CorrectAnswer: int16(q.CorrectAnswer),
	}

	if len(q.Answers) == 4 {
		nq.Type = sqli.QuestionType4Alt

		nq.Answer3 = pgtype.Text{
			Valid:  true,
			String: q.Answers[2],
		}
		nq.Answer4 = pgtype.Text{
			Valid:  true,
			String: q.Answers[3],
		}
	} else {
		nq.Type = sqli.QuestionType2Alt

		nq.Answer3 = pgtype.Text{
			Valid:  false,
			String: "",
		}
		nq.Answer4 = pgtype.Text{
			Valid:  false,
			String: "",
		}
	}

	return &nq
}

func createQuestionParamsToDbEntity(q *question.QuestionCreateData) *sqli.CreateQuestionParams {
	nq := sqli.CreateQuestionParams{
		Question:      q.Question,
		Answer1:       q.Answers[0],
		Answer2:       q.Answers[1],
		CorrectAnswer: int16(q.CorrectAnswer),
		Style:         sqli.QuestionStyle(q.Style),
		Difficulty:    int16(q.Difficulty),
	}

	switch len(q.Answers) {
	case 2:
		nq.Type = sqli.QuestionType2Alt

		nq.Answer3 = pgtype.Text{
			String: "",
			Valid:  false,
		}

		nq.Answer4 = pgtype.Text{
			String: "",
			Valid:  false,
		}
	case 4:
		nq.Type = sqli.QuestionType4Alt

		nq.Answer3 = pgtype.Text{
			String: q.Answers[2],
			Valid:  true,
		}

		nq.Answer4 = pgtype.Text{
			String: q.Answers[3],
			Valid:  true,
		}
	}

	return &nq
}

func questionToApiEntity(m *sqli.Question) *question.Question {
	q := question.Question{
		ID:            m.ID,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
		Question:      m.Question,
		Answers:       []string{m.Answer1, m.Answer2},
		Type:          question.QuestionType(m.Type),
		Style:         question.QuestionStyle(m.Style),
		CorrectAnswer: int32(m.CorrectAnswer),
		Difficulty:    uint8(m.Difficulty),
	}

	if m.Answer3.Valid {
		q.Answers = append(q.Answers, m.Answer3.String)
	}

	if m.Answer4.Valid {
		q.Answers = append(q.Answers, m.Answer4.String)
	}

	return &q
}
