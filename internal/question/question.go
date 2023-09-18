package question

import "time"

type Question struct {
	ID int32 `json:"id,omitempty" validate:"required"`

	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`

	Question string   `json:"question,omitempty" validate:"required"`
	Answers  []string `json:"answers,omitempty" validate:"required"`

	CorrectAnswer int32 `json:"correct_answer" validate:"gte=0,lte=3"`

	Type  QuestionType  `json:"type" validate:"required"`
	Style QuestionStyle `json:"style" validate:"required"`

	Difficulty uint8 `json:"difficulty"`
}

func (q *Question) IsValid() bool {
	if q.CorrectAnswer <= 0 {
		return false
	}

	if q.Style != QuestionStyleAudio &&
		q.Style != QuestionStyleImage &&
		q.Style != QuestionStyleText {
		return false
	}

	if q.Type == QuestionType4Alt {
		if len(q.Answers) != 4 || q.CorrectAnswer > 4 {
			return false
		}
	} else if q.Type == QuestionType2Alt {
		if len(q.Answers) != 2 || q.CorrectAnswer > 2 {
			return false
		}
	} else {
		return false
	}

	for i := range q.Answers {
		if len(q.Answers[i]) > 64 {
			return false
		}
	}

	if len(q.Question) > 200 {
		return false
	}

	return true
}

func (q *Question) IntoQuestionAlternative() *QuestionAlternativeFmt {
	nq := QuestionAlternativeFmt{
		ID:            q.ID,
		CreatedAt:     q.CreatedAt,
		UpdatedAt:     q.UpdatedAt,
		Question:      q.Question,
		CorrectAnswer: q.CorrectAnswer,
		Type:          q.Type,
		Style:         q.Style,
		Difficulty:    q.Difficulty,
	}

	nq.Answer1 = q.Answers[0]
	nq.Answer2 = q.Answers[1]

	if nq.Type == QuestionType4Alt {
		nq.Answer3 = &q.Answers[2]
		nq.Answer4 = &q.Answers[3]
	}

	return &nq
}
