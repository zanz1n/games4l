package main

import (
	"context"
	"os"
	"time"

	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/telemetry"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dba *telemetry.TelemetryService

func Connect() error {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DATABASE_NAME")
	WebhookSig := os.Getenv("WEBHOOK_SIG")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	mongoConnStartTime := time.Now()

	serverApi := options.ServerAPI(options.ServerAPIVersion1)

	defer cancel()

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return err
	}

	dba = telemetry.NewTelemetryDataService(client, &telemetry.Config{
		ProjectEpoch: 1684542947161,
		MongoDbName: dbName,
	})

	_ = WebhookSig

	logger.Info("Connected to mongodb, handshake took %v", time.Since(mongoConnStartTime))

	return nil
}
