package telemetry

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/games4l/backend/libs/logger"
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

type similarNameResult struct {
	res []TelemetryUnit
	err utils.StatusCodeErr
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

func eliminateDuplicates(arr []TelemetryUnit) []TelemetryUnit {
	cache := make(map[string]struct{})

	newArr := []TelemetryUnit{}

	var (
		ok bool
		v  TelemetryUnit
	)
	for _, v = range arr {
		if _, ok = cache[v.ID.Hex()]; ok {
			continue
		}

		newArr = append(newArr, v)
		cache[v.ID.Hex()] = struct{}{}
	}

	return newArr
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
			return nil, utils.NewStatusCodeErr(fmt.Sprintf("invalid answered range on %v° item", i+1), httpcodes.StatusBadRequest)
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

func (ds *TelemetryService) FindSimilarName(ctx context.Context, name string) ([]TelemetryUnit, utils.StatusCodeErr) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	spl := strings.Split(name, " ")

	if len(spl) < 2 || len(spl) > 9 {
		return nil, utils.NewStatusCodeErr("at least one surname must be provided", httpcodes.StatusBadRequest)
	}

	deadline, ok := ctx.Deadline()

	if !ok {
		return nil, utils.NewStatusCodeErr("invalid context was provided", httpcodes.StatusInternalServerError)
	}

	done := make(chan similarNameResult)

	go func() {
		res, err := ds.findSimilarName(deadline, name)

		done <- similarNameResult{
			err: err,
			res: res,
		}
	}()

	select {
	case result := <-done:
		return result.res, result.err
	case <-ctx.Done():
		return nil, utils.NewStatusCodeErr("query timeout exceded", httpcodes.StatusRequestTimeout)
	}
}

func (ds *TelemetryService) findSimilarName(deadline time.Time, name string) ([]TelemetryUnit, utils.StatusCodeErr) {
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

	results = eliminateDuplicates(results)

	select {
	case <-timeoutTicker.C:
		logger.Info("Query timed out in %v", time.Since(queryStart))
	case <-retCh:
		logger.Info("Query completed successfully in %v", time.Since(queryStart))
	}

	return results, nil
}
