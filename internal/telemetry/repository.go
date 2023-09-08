package telemetry

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/games4l/internal/errors"
	"github.com/games4l/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/transform"
)

type Config struct {
	ProjectEpoch int64
	MongoDbName  string
}

type TelemetryService struct {
	client *mongo.Client
	db     *mongo.Database
	col    *mongo.Collection
	cfg    *Config
}

func NewTelemetryService(c *mongo.Client, cfg *Config) *TelemetryService {
	db := c.Database(cfg.MongoDbName)

	col := db.Collection("telemetry_data")

	return &TelemetryService{
		client: c,
		db:     db,
		col:    col,
		cfg:    cfg,
	}
}

func (ds *TelemetryService) FindByIdWithCtx(ctx context.Context, id string) (*TelemetryUnit, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	done := make(chan findOneResult)

	go func() {
		resultR, errR := ds.FindById(id)

		done <- findOneResult{
			res: resultR,
			err: errR,
		}
	}()

	select {
	case <-ctx.Done():
		return nil, errors.ErrServerOperationTookTooLong
	case result := <-done:
		return result.res, result.err
	}
}

func (ds *TelemetryService) FindById(id string) (*TelemetryUnit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tu := TelemetryUnit{}

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, errors.ErrInvalidObjectID
	}

	err = ds.col.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&tu)

	if err != nil {
		return nil, errors.ErrEntityNotFound
	}

	return &tu, nil
}

func (ds *TelemetryService) Create(data *CreateTelemetryUnitData) (*TelemetryUnit, error) {
	err := validate.Struct(*data)

	if err != nil {
		return nil, errors.ErrInvalidRequestEntity
	}

	for _, answred := range data.Answereds {
		if answred > 4 || answred < 1 {
			return nil, errors.ErrInvalidRequestEntity
		}
	}

	if data.DoneAt < ds.cfg.ProjectEpoch {
		return nil, errors.ErrInvalidRequestEntity
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	normalizedName, _, err := transform.String(normalizer, strings.ToLower(data.PacientName))

	if err != nil {
		return nil, errors.ErrInternalServerError
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
		return nil, errors.ErrInternalServerError
	}

	return &tu, nil
}

func (ds *TelemetryService) FindSimilarNameWithCtx(ctx context.Context, name string) ([]TelemetryUnit, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	spl := strings.Split(name, " ")

	if len(spl) < 2 || len(spl) > 9 {
		return nil, errors.ErrSurnameSearchInvalid
	}

	deadline, ok := ctx.Deadline()

	if !ok {
		return nil, errors.ErrInternalServerError
	}

	done := make(chan similarNameResult)

	go func() {
		res, err := ds.FindSimilarName(deadline, name)

		done <- similarNameResult{
			err: err,
			res: res,
		}
	}()

	select {
	case result := <-done:
		return result.res, result.err
	case <-ctx.Done():
		return nil, errors.ErrServerOperationTookTooLong
	}
}

func (ds *TelemetryService) FindSimilarName(deadline time.Time, name string) ([]TelemetryUnit, error) {
	spl := strings.Split(name, " ")

	queryStart := time.Now()

	firstName, surnames := spl[0], spl[1:]

	resultChan := make(chan struct{})

	maxOpDelay := time.Duration(math.Floor(float64(18)/float64(len(surnames)))) * time.Second

	results := []TelemetryUnit{}
	resultsM := sync.Mutex{}

	for i_, surname_ := range surnames {
		i := i_
		surname := surname_

		go func() {
			logger.Info("Op started - (%v, %s) after %v", i, surname, time.Since(queryStart))
			opStart := time.Now()

			ctx, cancel := context.WithTimeout(context.Background(), maxOpDelay)
			defer cancel()

			filter := bson.M{
				"pacient_name": bson.M{
					"$regex": fmt.Sprintf("(?<=%s)(.*)(?=%s)", firstName, surname),
				},
			}

			cursor, err := ds.col.Find(ctx, filter)

			if err != nil {
				logger.Info("Op (%v, %s) failed in (*mongo.Collection).Find() in %v", i, surname, time.Since(opStart))
				resultChan <- struct{}{}
				return
			}

			ctx2, cancel2 := context.WithTimeout(context.Background(), maxOpDelay)
			defer cancel2()

			res := []TelemetryUnit{}

			if err = cursor.All(ctx2, &res); err != nil {
				logger.Info("Op (%v, %s) failed in (*mongo.Cursor).All() in %v", i, surname, time.Since(opStart))
				resultChan <- struct{}{}
				return
			}

			resultsM.Lock()
			results = append(results, res...)
			resultsM.Unlock()

			logger.Info("Op completed - (%v, %s) in %v", i, surname, time.Since(opStart))
			resultChan <- struct{}{}
		}()
	}

	retCh := make(chan struct{})
	timeoutTicker := time.NewTicker(deadline.Sub(queryStart.Add(24 * time.Millisecond)))

	go func() {
		ti := 1
		for {
			<-resultChan
			ti++
			if ti == len(surnames) {
				retCh <- struct{}{}
				break
			}
		}
	}()

	select {
	case <-timeoutTicker.C:
		logger.Info("Query timed out in %v", time.Since(queryStart))
	case <-retCh:
		logger.Info("Query completed successfully in %v", time.Since(queryStart))
	}

	results = eliminateDuplicates(results)

	return results, nil
}
