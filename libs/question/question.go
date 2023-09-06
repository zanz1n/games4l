package question

type QuestionStyle string

const (
	QuestionStyleImage QuestionStyle = "image"
	QuestionStyleAudio QuestionStyle = "audio"
	QuestionStyleVideo QuestionStyle = "video"
	QuestionStyleText  QuestionStyle = "text"
)

type QuestionType string

const (
	QuestionType2Alt QuestionType = "2Alt"
	QuestionType4Alt QuestionType = "4Alt"
)

type QuestionUpdateData struct {
	Question      *string       `json:"question"`
	Answers       []string      `json:"answers"`
	CorrectAnswer *string       `json:"correct_answer"`
	Type          *QuestionType `json:"type"`
	File          *string       `json:"file"`
	ImageWidth    *int          `json:"image_width"`
	ImageHeight   *int          `json:"image_height"`
}

type Question struct {
	Question string   `json:"question,omitempty" bson:"question,omitempty" validate:"required"`
	Answers  []string `json:"answers,omitempty" bson:"answers,omitempty" validate:"required"`

	CorrectAnswer *int `json:"correct_answer,omitempty" bson:"correct_answer,omitempty" validate:"required"`

	Type  QuestionType  `json:"type,omitempty" bson:"type,omitempty" validate:"required"`
	Style QuestionStyle `json:"style,omitempty" bson:"style,omitempty" validate:"required"`

	File *string `json:"file" bson:"file"`

	ImageWidth  *int `json:"image_width" bson:"image_width"`   // Nullable
	ImageHeight *int `json:"image_height" bson:"image_height"` // Nullable
}

func (q *Question) IsValid() bool {
	if *q.CorrectAnswer <= 0 {
		return false
	}

	if q.Type == QuestionType4Alt {
		if len(q.Answers) != 4 || *q.CorrectAnswer > 4 {
			return false
		}
	} else if q.Type == QuestionType2Alt {
		if len(q.Answers) != 2 || *q.CorrectAnswer > 2 {
			return false
		}
	} else {
		return false
	}

	if q.Style != QuestionStyleText && q.File == nil {
		return false
	}

	if (q.ImageHeight != nil && q.ImageWidth == nil) ||
		(q.ImageWidth != nil && q.ImageHeight == nil) {
		return false
	}

	return true
}

func (q *Question) Parse() *QuestionAlternativeFmt {
	nq := QuestionAlternativeFmt{
		Question:      q.Question,
		CorrectAnswer: *q.CorrectAnswer,
		Type:          q.Type,
		Style:         q.Style,
		Audio:         nil,
		File:          nil,
		ImageWidth:    nil,
		ImageHeight:   nil,
	}

	if q.ImageHeight != nil && q.ImageWidth != nil {
		nq.ImageHeight = q.ImageHeight
		nq.ImageWidth = q.ImageWidth
	}

	if q.File != nil {
		if q.Style == QuestionStyleAudio {
			nq.Audio = q.File
		} else {
			nq.File = q.File
		}
	}

	nq.Answer1 = q.Answers[0]
	nq.Answer2 = q.Answers[1]

	if nq.Type == QuestionType4Alt {
		nq.Answer3 = &q.Answers[2]
		nq.Answer4 = &q.Answers[3]
	}

	return &nq
}

type QuestionAlternativeFmt struct {
	Question string `json:"quest,omitempty" validate:"required"`

	Answer1 string  `json:"a1,omitempty" validate:"required"`
	Answer2 string  `json:"a2,omitempty" validate:"required"`
	Answer3 *string `json:"a3,omitempty"` // Nullable
	Answer4 *string `json:"a4,omitempty"` // Nullable

	CorrectAnswer int `json:"ccr,omitempty" validate:"required"`

	Audio *string      `json:"audio,omitempty"` // Nullable
	Type  QuestionType `json:"type,omitempty" validate:"required"`

	File *string `json:"file,omitempty"` // Nullable

	Style QuestionStyle `json:"style,omitempty" validate:"required"`

	ImageWidth  *int `json:"x"` // Nullable
	ImageHeight *int `json:"y"` // Nullable
}

func (q *QuestionAlternativeFmt) IsValid() bool {
	if q.CorrectAnswer <= 0 {
		return false
	}

	if q.Type == QuestionType4Alt {
		if q.Answer3 == nil || q.Answer4 == nil || q.CorrectAnswer > 4 {
			return false
		}
	} else if q.Type == QuestionType2Alt {
		if q.CorrectAnswer > 2 {
			return false
		}
	} else {
		return false
	}

	if q.Style == QuestionStyleAudio {
		if q.Audio == nil {
			return false
		}
	} else {
		if q.Audio != nil {
			return false
		}
	}

	if q.Style != QuestionStyleText {
		if q.Audio == nil && q.File == nil {
			return false
		}
	}

	if (q.ImageHeight != nil && q.ImageWidth == nil) ||
		(q.ImageWidth != nil && q.ImageHeight == nil) {
		return false
	}

	return true
}

func (q *QuestionAlternativeFmt) Parse() *Question {
	nq := Question{
		Question:      q.Question,
		Answers:       []string{q.Answer1, q.Answer2},
		CorrectAnswer: &q.CorrectAnswer,
		Type:          q.Type,
		Style:         q.Style,
		ImageWidth:    q.ImageWidth,
		ImageHeight:   q.ImageHeight,
	}

	if q.Type == QuestionType4Alt {
		nq.Answers = append(nq.Answers, *q.Answer3, *q.Answer4)
	}

	if q.Style == QuestionStyleAudio {
		nq.File = q.Audio
	} else {
		nq.File = q.File
	}

	return &nq
}
