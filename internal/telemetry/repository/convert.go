package repository

import (
	"strings"
	"time"

	"github.com/games4l/internal/telemetry"
	"github.com/games4l/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/transform"
)

type mongoRegistryData struct {
	ID           primitive.ObjectID `bson:"_id" validate:"required"`
	CreatedAt    primitive.DateTime `bson:"created_at" validate:"required"`
	DoneAt       primitive.DateTime `bson:"done_at" validate:"required"`
	CompleteTime int64              `bson:"complete_time" validate:"gt=200"`
	Answereds    []int8             `bson:"answereds" validate:"required"`
	QuestionID   int32              `bson:"question_id" validate:"gte=0"`
	PacientName  string             `bson:"pacient_name" validate:"required"`
}

func newMongoRegistryData(data *telemetry.CreateRegistryData) (*mongoRegistryData, error) {
	normalizedName, _, err := transform.String(
		normalizer,
		strings.ToLower(data.PacientName),
	)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	return &mongoRegistryData{
		ID:           primitive.NewObjectID(),
		CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
		DoneAt:       primitive.NewDateTimeFromTime(time.UnixMilli(data.DoneAt)),
		CompleteTime: data.CompleteTime,
		Answereds:    data.Answereds,
		QuestionID:   data.QuestionID,
		PacientName:  normalizedName,
	}, nil
}

func telemetryToDbEntity(t *telemetry.Registry) (*mongoRegistryData, error) {
	objectid, err := primitive.ObjectIDFromHex(t.ID)
	if err != nil {
		return nil, errors.ErrInvalidObjectID
	}

	return &mongoRegistryData{
		ID:           objectid,
		CreatedAt:    primitive.NewDateTimeFromTime(t.CreatedAt),
		DoneAt:       primitive.NewDateTimeFromTime(t.DoneAt),
		CompleteTime: t.CompleteTime,
		Answereds:    t.Answereds,
		QuestionID:   t.QuestionID,
		PacientName:  t.PacientName,
	}, nil
}

func telemetryToApiEntity(t *mongoRegistryData) *telemetry.Registry {
	return &telemetry.Registry{
		ID:           t.ID.Hex(),
		CreatedAt:    t.CreatedAt.Time(),
		DoneAt:       t.DoneAt.Time(),
		CompleteTime: t.CompleteTime,
		Answereds:    t.Answereds,
		QuestionID:   t.QuestionID,
		PacientName:  t.PacientName,
	}
}
