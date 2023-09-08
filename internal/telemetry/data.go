package telemetry

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/games4l/internal/errors"
)

type CreateTelemetryUnitData struct {
	DoneAt       int64  `json:"done_at,omitempty" validate:"required"`
	CompleteTime int64  `json:"complete_time,omitempty" validate:"required"`
	Answereds    []int8 `json:"answereds,omitempty" validate:"required"`
	QuestionID   uint   `json:"quest_id,omitempty" validate:"required"`
	PacientName  string `json:"pacient_name,omitempty" validate:"required"`
}

type TelemetryUnit struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" validate:"required"`
	CreatedAt    primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty" validate:"required"`
	DoneAt       primitive.DateTime `bson:"done_at,omitempty" json:"done_at,omitempty" validate:"required"`
	CompleteTime int64              `bson:"complete_time,omitempty" json:"complete_time,omitempty" validate:"required"`
	Answereds    []int8             `bson:"answereds,omitempty" json:"answered,omitempty" validate:"required"`
	QuestionID   uint               `bson:"quest_id,omitempty" json:"quest_id,omitempty" validate:"required"`
	PacientName  string             `bson:"pacient_name,omitempty" json:"pacient_name,omitempty" validate:"required"`
}

type similarNameResult struct {
	res []TelemetryUnit
	err errors.StatusCodeErr
}

type findOneResult struct {
	res *TelemetryUnit
	err errors.StatusCodeErr
}
