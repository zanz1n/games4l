package question

type QuestionStyle string

const (
	QuestionStyleImage QuestionStyle = "image"
	QuestionStyleAudio QuestionStyle = "audio"
	QuestionStyleText  QuestionStyle = "text"
)

type QuestionType string

const (
	QuestionType2Alt QuestionType = "2Alt"
	QuestionType4Alt QuestionType = "4Alt"
)
