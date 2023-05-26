package telemetry

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/games4l/backend/libs/utils"
	"github.com/games4l/backend/libs/utils/httpcodes"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	validate   = validator.New()
)

type Config struct {
	ProjectEpoch int64
	MongoDbName  string
}

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
	cfg    *Config
}

func NewTelemetryDataService(c *mongo.Client, cfg *Config) *TelemetryService {
	db := c.Database(cfg.MongoDbName)

	col := db.Collection("telemetry_data")

	return &TelemetryService{
		client: c,
		db:     db,
		col:    col,
		cfg:    cfg,
	}
}

func (ds *TelemetryService) FindById(id string) (*TelemetryUnit, utils.StatusCodeErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tu := TelemetryUnit{}

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, utils.NewStatusCodeErr("invalid object id format", httpcodes.StatusBadRequest)
	}

	err = ds.col.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&tu)

	if err != nil {
		return nil, utils.NewStatusCodeErr("not find", httpcodes.StatusNotFound)
	}

	return &tu, nil
}

func (ds *TelemetryService) Create(data *CreateTelemetryUnitData) (*TelemetryUnit, utils.StatusCodeErr) {
	err := validate.Struct(*data)

	if err != nil {
		return nil, utils.NewStatusCodeErr("invalid body schema", httpcodes.StatusBadRequest)
	}

	for i, answred := range data.Answereds {
		if answred > 4 || answred < 1 {
			return nil, utils.NewStatusCodeErr(fmt.Sprintf("invalid answered range on %v idx", i), httpcodes.StatusBadRequest)
		}
	}

	if data.DoneAt < ds.cfg.ProjectEpoch {
		return nil, utils.NewStatusCodeErr("timestamp out of accepted range on done_at", httpcodes.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	normalizedName, _, err := transform.String(normalizer, strings.ToLower(data.PacientName))

	if err != nil {
		return nil, utils.NewStatusCodeErr("something went wrong", httpcodes.StatusInternalServerError)
	}

	tu := TelemetryUnit{
		DoneAt:       primitive.NewDateTimeFromTime(time.UnixMilli(data.DoneAt)),
		CompleteTime: data.CompleteTime,
		Answereds:    data.Answereds,
		QuestionID:   data.QuestionID,
		PacientName:  normalizedName,
	}

	tu.ID = primitive.NewObjectID()

	tu.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = ds.col.InsertOne(ctx, tu)

	if err != nil {
		return nil, utils.NewStatusCodeErr("something went wrong", httpcodes.StatusInternalServerError)
	}

	return &tu, nil
}
