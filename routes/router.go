package routes

import (
	"context"
	"time"

	"github.com/games4l/telemetria/logger"
	"github.com/games4l/telemetria/providers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
)

func StartRouter(app *fiber.App) {
	config := providers.GetConfig()

	var err error

	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoUri))

	if err != nil {
		logger.Fatal(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"Hello": "World",
		})
	})
}

func ShutdownRouter(app *fiber.App) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)

	if err := mongoClient.Disconnect(ctx); err != nil {
		logger.Fatal(err)
	}
}
