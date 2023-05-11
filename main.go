package main

import (
	"github.com/games4l/telemetria/logger"
	"github.com/games4l/telemetria/providers"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config := providers.GetConfig()

	app := fiber.New(fiber.Config{
		Prefork:               false,
		ServerHeader:          "Fiber",
		CaseSensitive:         true,
		StrictRouting:         false,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	app.Hooks().OnListen(func() error {
		logger.Info("Listenning on port %v with %v handlers", config.Port, app.HandlersCount())
		return nil
	})

	app.Hooks().OnShutdown(func() error {
		logger.Info("Shutting down ...")
		return nil
	})

	app.Use(logger.NewFiberMiddleware())

	app.Use(recover.New())
	app.Use(cors.New())

	app.Listen(":3333")
}
