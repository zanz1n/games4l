package src

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/games4l/internal/logger"
	"github.com/games4l/pkg/errors"
)

func ConnectS3() error {
	if s3c != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("sa-east-1"))
	if err != nil {
		logger.Error("Failed to instantiate aws config: " + err.Error())
		return errors.ErrInternalServerError
	}
	s3Bucket = os.Getenv("APP_QUESTION_BUCKET_NAME")
	s3c = s3.NewFromConfig(cfg)

	return nil
}
