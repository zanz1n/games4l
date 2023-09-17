package s3u

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/games4l/internal/logger"
	"github.com/games4l/pkg/errors"
)

func NewS3Singleton(region string) *S3Singleton {
	return &S3Singleton{
		region: region,
		logger: log.NewLogger("s3_singleton"),
		r:      nil,
	}
}

type S3Singleton struct {
	region string
	r      *s3.Client
	logger log.Logger
}

func (s *S3Singleton) GetInstance() (*s3.Client, error) {
	if s.r != nil {
		return s.r, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("sa-east-1"))
	if err != nil {
		s.logger.Error("Failed to create aws config: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	s.r = s3.NewFromConfig(cfg)

	return s.r, nil
}
