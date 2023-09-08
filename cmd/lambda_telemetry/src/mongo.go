package src

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/games4l/internal/telemetry"
	"github.com/games4l/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() error {
	if dba != nil {
		return nil
	}

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DATABASE_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	mongoConnStartTime := time.Now()

	serverApi := options.ServerAPI(options.ServerAPIVersion1)

	defer cancel()

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return errors.New("failed to connect to mongodb: " + err.Error())
	}

	dba = telemetry.NewTelemetryService(client, &telemetry.Config{
		ProjectEpoch: 1684542947161,
		MongoDbName:  dbName,
	})

	logger.Info("Connected to mongodb, handshake took %v", time.Since(mongoConnStartTime))

	return nil
}
