package providers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateTelemetryUnitData struct {
	DoneAt        int64 `json:"done_at,omitempty" validate:"required"`
	CompleteTime  int64 `json:"complete_time,omitempty" validate:"required"`
	AnsweredQuest uint8 `json:"answered_quest,omitempty" validate:"required"`
}

type TelemetryUnit struct {
	ID            string `bson:"_id,omitempty" json:"id,omitempty" validate:"required"`
	CreatedAt     int64  `bson:"created_at,omitempty" json:"created_at,omitempty" validate:"required"`
	DoneAt        int64  `bson:"done_at,omitempty" json:"done_at,omitempty" validate:"required"`
	CompleteTime  int64  `bson:"complete_time,omitempty" json:"complete_time,omitempty" validate:"required"`
	AnsweredQuest uint8  `bson:"answered_quest,omitempty" json:"answered_quest,omitempty" validate:"required"`
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

func (ds *TelemetryService) FindById(id string) (*TelemetryUnit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := TelemetryUnit{}

	err := ds.col.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ds *TelemetryService) Create(data *CreateTelemetryUnitData) (*TelemetryUnit, error) {
	err := validate.Struct(*data)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dt := TelemetryUnit{
		DoneAt:        data.DoneAt,
		CompleteTime:  data.CompleteTime,
		AnsweredQuest: data.AnsweredQuest,
	}

	dt.ID = uuid.NewString()
	dt.CreatedAt = time.Now().UnixMilli()

	_, err = ds.col.InsertOne(ctx, dt)

	if err != nil {
		return nil, err
	}

	return &dt, nil
}
