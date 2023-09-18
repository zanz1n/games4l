package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/games4l/internal/logger"
	"github.com/games4l/internal/telemetry"
	"github.com/games4l/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewMongo(url string, dbName string) (TelemetryRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	c, err := mongo.Connect(ctx)
	if err != nil {
		return nil, err
	}

	col := c.Database(dbName).Collection("telemetry_data")

	return newTelemetryService(col), nil
}

func newTelemetryService(col *mongo.Collection) *telemetryService {
	return &telemetryService{
		col:    col,
		logger: log.NewLogger("telemetry_repository"),
	}
}

type telemetryService struct {
	col    *mongo.Collection
	logger log.Logger
}

func (s *telemetryService) GetById(id string) (*telemetry.Registry, error) {
	objectid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.ErrInvalidObjectID
	}

	r := mongoRegistryData{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = s.col.FindOne(ctx, bson.D{{Key: "_id", Value: objectid}}).Decode(&r)
	if err != nil {
		s.logger.Info("Failed to get by id: " + err.Error())
		return nil, errors.ErrEntityNotFound
	}

	return telemetryToApiEntity(&r), nil
}

func (s *telemetryService) Create(data *telemetry.CreateRegistryData) (*telemetry.Registry, error) {
	r, err := newMongoRegistryData(data)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err = s.col.InsertOne(ctx, r); err != nil {
		s.logger.Error("Failed to create: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	return telemetryToApiEntity(r), nil
}

func (s *telemetryService) GetBySimilarName(name string) ([]telemetry.Registry, error) {
	spl := strings.Split(name, " ")

	firstName, surnames := spl[0], spl[1:]

	if 1 > len(surnames) {
		return nil, errors.ErrSurnameSearchInvalid
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	res := []mongoRegistryData{}

	var err2 error
	for _, sn := range surnames {
		wg.Add(1)

		surname := sn
		go func() {
			arr, err := s.findByName(firstName, surname)
			wg.Done()

			if err != nil {
				err2 = err
			} else {
				mu.Lock()
				res = append(res, arr...)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if err2 != nil {
		return nil, err2
	}

	res = eliminateDuplicates(res)
	newRes := make([]telemetry.Registry, len(res))

	for i := range res {
		newRes[i] = *telemetryToApiEntity(&res[i])
	}

	return newRes, nil
}

func (s *telemetryService) findByName(firstname, surname string) ([]mongoRegistryData, error) {
	filter := bson.M{
		"pacient_name": bson.M{
			"$regex": fmt.Sprintf("(?<=%s)(.*)(?=%s)", firstname, surname),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := s.col.Find(ctx, filter)
	if err != nil {
		s.logger.Error("Surname search failed: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	r := []mongoRegistryData{}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()

	if err = cursor.All(ctx2, &r); err != nil {
		s.logger.Error("Surname search cursor.All failed: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	return r, nil
}

func (s *telemetryService) DeleteById(id string) (*telemetry.Registry, error) {
	objectid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.ErrInvalidObjectID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r := mongoRegistryData{}

	err = s.col.FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: objectid}}).Decode(&r)
	if err != nil {
		s.logger.Info("Failed to delete: " + err.Error())
		return nil, errors.ErrEntityNotFound
	}

	return telemetryToApiEntity(&r), nil
}
