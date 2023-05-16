package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/games4l/telemetria/logger"
	"github.com/games4l/telemetria/providers"
	"github.com/games4l/telemetria/routes"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var endCh = make(chan os.Signal)

func init() {
	logger.Init()

	if os.Getenv("APP_CONFIG") != "" {
		providers.AcquireFromEnv()
	} else {
		providers.AcquireFromFile(os.Getenv("APP_CONFIG_FILE"))
	}
}

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

	app.Use(logger.NewFiberMiddleware())

	app.Use(recover.New())
	app.Use(cors.New())

	routes.StartRouter(app)

	signal.Notify(endCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go app.Listen(fmt.Sprintf(":%v", config.Port))

	<-endCh
	logger.Info("Shutting down ...")
	routes.ShutdownRouter(app)
}
