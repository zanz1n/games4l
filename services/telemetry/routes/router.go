package routes

import (
	"context"
	"time"

	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/config"
	"github.com/games4l/backend/libs/logger"
	"github.com/games4l/backend/libs/telemetry"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	mongoClient *mongo.Client
)

func StartRouter(app *fiber.App) {
	cfg := config.GetConfig()

	var err error

	mongoConnStartTime := time.Now()

	serverApi := options.ServerAPI(options.ServerAPIVersion1)

	ctx, cancel := context.WithTimeout(context.Background(), 16*time.Second)

	defer cancel()

	opts := options.Client().ApplyURI(cfg.MongoUri).SetServerAPIOptions(serverApi)

	mongoClient, err = mongo.Connect(ctx, opts)

	if err != nil {
		logger.Fatal(err)
	}

	err = mongoClient.Ping(ctx, readpref.Primary())

	telemetryService := telemetry.NewTelemetryDataService(mongoClient, cfg)

	if err != nil {
		logger.Fatal(err)
	}

	authProvider := auth.NewAuthProvider([]byte(cfg.WebhookSig))

	app.Post(cfg.RoutePrefix+"/telemetry", PostTelemetry(telemetryService))
	app.Get(cfg.RoutePrefix+"/telemetry/:id", GetTelemetryUnit(telemetryService, authProvider))

	logger.Info("Connected to mongodb, handshake took %v", time.Since(mongoConnStartTime))
}

func ShutdownRouter(app *fiber.App) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)

	if err := mongoClient.Disconnect(ctx); err != nil {
		logger.Fatal(err)
	}
}
