package providers

import (
	"context"
	"strings"
	"time"

	"github.com/games4l/telemetry-service/logger"
	nanoid "github.com/matoous/go-nanoid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateTelemetryUnitData struct {
	DoneAt       int64  `json:"done_at,omitempty" validate:"required"`
	CompleteTime int64  `json:"complete_time,omitempty" validate:"required"`
	Answereds    []int8 `json:"answereds,omitempty" validate:"required"`
	QuestionID   uint   `json:"quest_id,omitempty" validate:"required"`
	PacientName  string `json:"pacient_name,omitempty" validate:"required"`
}

type TelemetryUnit struct {
	ID           string `bson:"_id,omitempty" json:"id,omitempty" validate:"required"`
	CreatedAt    int64  `bson:"created_at,omitempty" json:"created_at,omitempty" validate:"required"`

	DoneAt       int64  `bson:"done_at,omitempty" json:"done_at,omitempty" validate:"required"`
	CompleteTime int64  `bson:"complete_time,omitempty" json:"complete_time,omitempty" validate:"required"`
	Answereds    []int8 `bson:"answereds,omitempty" json:"answered,omitempty" validate:"required"`
	QuestionID   uint   `bson:"quest_id,omitempty" json:"quest_id,omitempty" validate:"required"`
	PacientName  string `bson:"pacient_name,omitempty" json:"pacient_name,omitempty" validate:"required"`
}

type TelemetryService struct {
	client *mongo.Client
	db     *mongo.Database
	col    *mongo.Collection
}

func NewTelemetryDataService(c *mongo.Client) *TelemetryService {
	config := GetConfig()

	db := c.Database(config.MongoDbName)
	col := db.Collection("telemetry_data")

	return &TelemetryService{
		client: c,
		db:     db,
		col:    col,
	}
}

func (ds *TelemetryService) FindById(id string) (tu *TelemetryUnit, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tu = &TelemetryUnit{}

	err = ds.col.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(tu)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	return
}

func (ds *TelemetryService) Create(data *CreateTelemetryUnitData) (tu *TelemetryUnit, err error) {
	err = validate.Struct(*data)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tu = &TelemetryUnit{
		DoneAt:       data.DoneAt,
		CompleteTime: data.CompleteTime,
		Answereds:    data.Answereds,
		QuestionID:   data.QuestionID,
		PacientName:  strings.ToLower(data.PacientName),
	}

	tu.ID, err = nanoid.New(18)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	tu.CreatedAt = time.Now().UnixMilli()

	_, err = ds.col.InsertOne(ctx, tu)

	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return tu, nil
}
