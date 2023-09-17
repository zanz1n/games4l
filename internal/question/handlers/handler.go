package handlers

import (
	log "github.com/games4l/internal/logger"
	"github.com/games4l/internal/question/repository"
	"github.com/games4l/internal/utils/s3u"
	"github.com/games4l/pkg/ffmpeg"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func NewQuestionHandlers(
	qs repository.Singleton,
	ss *s3u.S3Singleton,
	fmpg *ffmpeg.Provider,
	s3Bucket string,
	s3ImagesPath string,
) *QuestionHandlers {
	return &QuestionHandlers{
		qs:           qs,
		ss:           ss,
		fmpg:         fmpg,
		s3Bucket:     s3Bucket,
		s3ImagesPath: s3ImagesPath,
		logger:       log.NewLogger("question_handlers"),
	}
}

type QuestionHandlers struct {
	qs           repository.Singleton
	ss           *s3u.S3Singleton
	s3ImagesPath string
	logger       log.Logger
	s3Bucket     string
	fmpg         *ffmpeg.Provider
}
