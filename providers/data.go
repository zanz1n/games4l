package providers

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/games4l/telemetry-service/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

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

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		err = errors.New("invalid object id format")
		return
	}

	err = ds.col.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(tu)

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

	for _, answred := range data.Answereds {
		if answred > 4 || answred < 1 {
			err = errors.New("invalid answered range")
			return
		}
	}

	if data.DoneAt < GetConfig().ProjectEpoch {
		err = errors.New("invalid done_at prop: timestamp out of accepted range")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	normalizedName, _, err := transform.String(normalizer, strings.ToLower(data.PacientName))

	tu = &TelemetryUnit{
		DoneAt:       primitive.NewDateTimeFromTime(time.UnixMilli(data.DoneAt)),
		CompleteTime: data.CompleteTime,
		Answereds:    data.Answereds,
		QuestionID:   data.QuestionID,
		PacientName:  normalizedName,
	}

	tu.ID = primitive.NewObjectID()

	if err != nil {
		logger.Error(err.Error())
		return
	}

	tu.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = ds.col.InsertOne(ctx, tu)

	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return tu, nil
}
