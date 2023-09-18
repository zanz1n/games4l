package question

type QuestionCreateData struct {
	Question string   `json:"question" validate:"required"`
	Answers  []string `json:"answers" validate:"required"`

	CorrectAnswer int32 `json:"correct_answer" validate:"gte=0,lte=3"`

	Style QuestionStyle `json:"style" validate:"required"`

	Difficulty uint8 `json:"difficulty" validate:"gte=1,lte=255"`
}

func (q *QuestionCreateData) IsValid() bool {
	if len(q.Question) > 200 {
		return false
	}

	if q.Style != QuestionStyleAudio &&
		q.Style != QuestionStyleImage &&
		q.Style != QuestionStyleText {
		return false
	}

	if q.CorrectAnswer != 0 &&
		q.CorrectAnswer != 1 &&
		q.CorrectAnswer != 2 &&
		q.CorrectAnswer != 3 {
		return false
	}

	switch len(q.Answers) {
	case 2:
		if q.CorrectAnswer != 0 && q.CorrectAnswer != 1 {
			return false
		}
	case 4:
	default:
		return false
	}

	for i := range q.Answers {
		if len(q.Answers[i]) > 64 {
			return false
		}
	}

	return true
}
