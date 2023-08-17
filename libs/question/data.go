package question

import (
	"context"
	"time"

	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/utils"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate = validator.New()

type Config struct {
	MongoDbName string
}

type QuestionService struct {
	client *mongo.Client
	db     *mongo.Database
	col    *mongo.Collection
	cfg    *Config
}

type QuestionDbData struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	NID       int                `json:"numeric_id,omitempty" bson:"n_id" validate:"required"`
	CreatedAt primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty" validate:"required"`
	UpdatedAt primitive.DateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty" validate:"required"`
	Data      Question           `json:"data,omitempty" bson:"data,omitempty" validate:"required"`
}

func NewQuestionService(c *mongo.Client, cfg *Config) *QuestionService {
	db := c.Database(cfg.MongoDbName)

	col := db.Collection("question_data")

	return &QuestionService{
		col:    col,
		client: c,
		db:     db,
		cfg:    cfg,
	}
}

func (s *QuestionService) GetMany(ctx context.Context, limit int64) ([]QuestionDbData, utils.StatusCodeErr) {
	maxExecTime := 20 * time.Second

	if deadLine, ok := ctx.Deadline(); ok {
		maxExecTime = time.Until(deadLine)
	}

	ctx, cancel := context.WithTimeout(ctx, maxExecTime/2)
	defer cancel()

	cursor, err := s.col.Find(ctx, bson.D{}, &options.FindOptions{Limit: &limit})
	if err != nil {
		return nil, utils.DefaultErrorList.EntityNotFound
	}

	result := []QuestionDbData{}

	ctx2, cancel2 := context.WithTimeout(ctx, maxExecTime/2)
	defer cancel2()

	if err := cursor.All(ctx2, &result); err != nil {
		return nil, utils.DefaultErrorList.InternalServerError
	}

	return result, nil
}

func (s *QuestionService) Create(ctx context.Context, numId int, data *Question) (*QuestionDbData, utils.StatusCodeErr) {
	if err := validate.Struct(*data); err != nil {
		return nil, utils.DefaultErrorList.InvalidRequestEntity
	}

	if !data.IsValid() {
		return nil, utils.DefaultErrorList.InvalidRequestEntity
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	now := primitive.NewDateTimeFromTime(time.Now())

	insertData := QuestionDbData{
		ID:        primitive.NewObjectID(),
		NID:       numId,
		CreatedAt: now,
		UpdatedAt: now,
		Data:      *data,
	}

	_, err := s.col.InsertOne(ctx, insertData)

	if err != nil {
		return nil, utils.DefaultErrorList.EntityAlreadyExists
	}

	return &insertData, nil
}

func (s *QuestionService) Update(ctx context.Context, hexID string, data *QuestionUpdateData) utils.StatusCodeErr {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(hexID)

	if err != nil {
		return utils.DefaultErrorList.InvalidObjectID
	}

	if _, err = s.col.UpdateByID(ctx, oid, *data); err != nil {
		logger.Error("%s", err.Error())
		return utils.DefaultErrorList.InternalServerError
	}

	return nil
}

func (s *QuestionService) GetByNumID(ctx context.Context, numId int) (*QuestionDbData, utils.StatusCodeErr) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	result := QuestionDbData{}

	err := s.col.FindOne(ctx, bson.D{{Key: "n_id", Value: numId}}).Decode(&result)

	if err != nil {
		return nil, utils.DefaultErrorList.EntityNotFound
	}

	if err = validate.Struct(result); err != nil {
		return nil, utils.DefaultErrorList.EntityNotFound

		// Implement here: delete the result after it's invalidated
	}

	return &result, nil
}

func (s *QuestionService) GetByID(ctx context.Context, hexID string) (*QuestionDbData, utils.StatusCodeErr) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(hexID)

	if err != nil {
		return nil, utils.DefaultErrorList.InvalidObjectID
	}

	result := QuestionDbData{}

	err = s.col.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&result)

	if err != nil {
		return nil, utils.DefaultErrorList.EntityNotFound
	}

	if err = validate.Struct(result); err != nil {
		return nil, utils.DefaultErrorList.EntityNotFound

		// Implement here: delete the result after it's invalidated
	}

	return &result, nil
}
