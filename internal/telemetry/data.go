package telemetry

import (
	"time"
)

type Registry struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	DoneAt       time.Time `json:"done_at"`
	CompleteTime int64     `json:"complete_time"`
	Answereds    []int8    `json:"answereds"`
	QuestionID   int32     `json:"question_id"`
	PacientName  string    `json:"pacient_name"`
}

type CreateRegistryData struct {
	DoneAt       int64  `json:"done_at" validate:"gt=1693710000"` // At least 3/9/2023
	CompleteTime int64  `json:"complete_time" validate:"gt=200"`
	Answereds    []int8 `json:"answereds" validate:"required"`
	QuestionID   int32  `json:"question_id" validate:"gte=1"`
	PacientName  string `json:"pacient_name" validate:"required"`
}

func (d *CreateRegistryData) IsValid() bool {
	for _, a := range d.Answereds {
		if a > 3 {
			return false
		}
	}

	return true
}
