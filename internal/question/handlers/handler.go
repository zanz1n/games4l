package handlers

import (
	"github.com/games4l/internal/bucket"
	log "github.com/games4l/internal/logger"
	"github.com/games4l/internal/question/repository"
	"github.com/games4l/pkg/ffmpeg"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func NewQuestionHandlers(
	qs repository.Singleton,
	bucket bucket.FileStorer,
	fmpg *ffmpeg.Provider,
	imageAssetPrefix string,
) *QuestionHandlers {
	return &QuestionHandlers{
		qs:     qs,
		bucket: bucket,
		fmpg:   fmpg,
		iap:    imageAssetPrefix,
		logger: log.NewLogger("question_handlers"),
	}
}

type QuestionHandlers struct {
	qs     repository.Singleton
	bucket bucket.FileStorer
	// Image asset prefix. The folder of the bucket that the image assets will
	// be stored
	iap    string
	logger log.Logger
	fmpg   *ffmpeg.Provider
}
