package bucket

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	log "github.com/games4l/internal/logger"
	"github.com/games4l/internal/utils/s3u"
	"github.com/games4l/pkg/errors"
)

func NewS3(bucket, region string) FileStorer {
	return &s3FileStorer{
		bucket: bucket,
		client: s3u.NewS3Singleton(region),
		logger: log.NewLogger("s3_file_storer"),
	}
}

type s3FileStorer struct {
	bucket string
	client *s3u.S3Singleton
	logger log.Logger
}

func (s *s3FileStorer) Store(key string, mime string, file io.Reader) error {
	c, err := s.client.GetInstance()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = c.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		ContentType: &mime,
		Body:        file,
	})

	if err != nil {
		s.logger.Error("Failed to put s3 object: " + err.Error())
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *s3FileStorer) Fetch(key string) (io.ReadCloser, error) {
	c, err := s.client.GetInstance()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	out, err := c.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})

	if err != nil {
		s.logger.Warn("Failed to get s3 object: " + err.Error())
		return nil, errors.ErrFileObjectNotFound
	}

	return out.Body, nil
}

func (s *s3FileStorer) Destroy(key string) error {
	c, err := s.client.GetInstance()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = c.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})

	if err != nil {
		s.logger.Warn("Failed to delete s3 object: " + err.Error())
		return errors.ErrFileObjectNotFound
	}

	return nil
}

func (s *s3FileStorer) DestroyMany(keys []string) error {
	c, err := s.client.GetInstance()
	if err != nil {
		return err
	}

	objs := make([]types.ObjectIdentifier, len(keys))
	for i := range objs {
		objs[i] = types.ObjectIdentifier{
			Key: &keys[i],
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = c.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: &s.bucket,
		Delete: &types.Delete{
			Objects: objs,
		},
	})

	if err != nil {
		s.logger.Warn("Failed to delete s3 object: " + err.Error())
		return errors.ErrInternalServerError
	}

	return nil
}
