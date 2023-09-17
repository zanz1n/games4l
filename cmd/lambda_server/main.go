package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	questionh "github.com/games4l/cmd/lambda_question/handler"
	"github.com/games4l/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var questionServer *questionh.Server

func init() {
	if os.Getenv("APP_ENV") == "" {
		if err := godotenv.Load(); err != nil {
			logger.Error("Failed to load .env file: " + err.Error())
		} else {
			logger.Info("Loaded environment configuration from .env file")
		}
	} else {
		logger.Info("Loaded environment configuration from shell variables")
	}
}

func main() {
	questionServer = questionh.NewEnvServer()

	endCh := make(chan os.Signal, 1)
	signal.Notify(endCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	r := chi.NewRouter()
	cr := cors.AllowAll()

	r.Use(LoggerMiddleware)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(middleware.RealIP)
	r.Use(middleware.RedirectSlashes)
	r.Use(cr.Handler)

	r.Get("/telemetry", HandleRequest(HandlerTypeTelemetry))
	r.Get("/telemetry/{id}", HandleRequest(HandlerTypeTelemetryParams))
	r.Post("/telemetry", HandleRequest(HandlerTypeTelemetry))

	r.Post("/auth/signin", HandleRequest(HandlerTypeAuth))
	r.Post("/user", HandleRequest(HandlerTypeAuth))

	r.Get("/question/{id}", HandleRequest(HandlerTypeQuestionParams))
	r.Get("/question", HandleRequest(HandlerTypeQuestion))
	r.Put("/question", HandleRequest(HandlerTypeQuestion))
	r.Patch("/question", HandleRequest(HandlerTypeQuestion))
	r.Post("/question", HandleRequest(HandlerTypeQuestion))

	go func() {
		addr := os.Getenv("APP_LNADDR")
		if addr == "" {
			addr = ":8080"
		}

		logger.Info("Listening on addr " + addr)

		if err := http.ListenAndServe(addr, r); err != nil {
			logger.Fatal(err)
		}
	}()

	<-endCh
	logger.Info("Stopping ...")
}
