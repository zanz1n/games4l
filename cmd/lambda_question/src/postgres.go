package src

import (
	"context"
	"os"
	"time"

	"github.com/games4l/internal/logger"
	"github.com/games4l/internal/sqli"
	"github.com/games4l/pkg/errors"
	"github.com/jackc/pgx/v5"
)

func Connect() error {
	if dba != nil {
		return nil
	}

	url := os.Getenv("DATABASE_URL")

	if url == "" {
		logger.Error("DATABASE_URL environment variable not provided")
		return errors.ErrInternalServerError
	}

	connCtx, connCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer connCancel()

	conn, err := pgx.Connect(connCtx, url)
	if err != nil {
		logger.Error("Failed to connect to postgres: " + err.Error())
		return errors.ErrInternalServerError
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingCancel()

	if err = conn.Ping(pingCtx); err != nil {
		logger.Error("Failed to ping postgres instance: " + err.Error())
		return errors.ErrInternalServerError
	}

	dba = sqli.New(conn)

	return nil
}
