package question

import "time"

type QuestionAlternativeFmt struct {
	ID int32 `json:"id"`

	CreatedAt time.Time `json:"created_at,omitempty" validate:"required"`
	UpdatedAt time.Time `json:"updated_at,omitempty" validate:"required"`

	Question string `json:"quest,omitempty" validate:"required"`

	Answer1 string  `json:"a1,omitempty" validate:"required"`
	Answer2 string  `json:"a2,omitempty" validate:"required"`
	Answer3 *string `json:"a3"` // Nullable
	Answer4 *string `json:"a4"` // Nullable

	CorrectAnswer int32 `json:"ccr"`

	Type  QuestionType  `json:"type,omitempty" validate:"required"`
	Style QuestionStyle `json:"style,omitempty" validate:"required"`

	Difficulty uint8 `json:"difficulty"`
}

func (q *QuestionAlternativeFmt) IsValid() bool {
	if q.CorrectAnswer <= 0 {
		return false
	}

	if q.Style != QuestionStyleAudio &&
		q.Style != QuestionStyleImage &&
		q.Style != QuestionStyleText {
		return false
	}

	if len(q.Answer1) > 64 ||
		len(q.Answer2) > 64 {
		return false
	}

	if q.Type == QuestionType4Alt {
		if q.Answer3 == nil || q.Answer4 == nil || q.CorrectAnswer > 4 {
			return false
		}

		if len(*q.Answer3) > 64 ||
			len(*q.Answer4) > 64 {
			return false
		}
	} else if q.Type == QuestionType2Alt {
		if q.CorrectAnswer > 2 {
			return false
		}
	} else {
		return false
	}

	return true
}

func (q *QuestionAlternativeFmt) IntoQuestion() *Question {
	nq := Question{
		ID:            q.ID,
		CreatedAt:     q.CreatedAt,
		UpdatedAt:     q.UpdatedAt,
		Question:      q.Question,
		Answers:       []string{q.Answer1, q.Answer2},
		CorrectAnswer: q.CorrectAnswer,
		Type:          q.Type,
		Style:         q.Style,
		Difficulty:    q.Difficulty,
	}

	if q.Type == QuestionType4Alt {
		nq.Answers = append(nq.Answers, *q.Answer3, *q.Answer4)
	}

	return &nq
}
