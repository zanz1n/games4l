package handlers

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"strings"

	"github.com/games4l/internal/question"
	"github.com/games4l/internal/utils"
	"github.com/games4l/pkg/errors"
)

func multipartToApiQuestion(
	mp *multipart.Form,
) (*question.QuestionCreateData, io.ReadCloser, error) {
	files, ok := mp.File["file"]
	if !ok || len(files) != 1 {
		return nil, nil, errors.ErrInvalidFormMedia
	}

	rawData, ok := mp.Value["data"]

	q := question.QuestionCreateData{}

	err := json.Unmarshal(
		utils.S2B(strings.Join(rawData, "")),
		&q,
	)
	if err != nil {
		return nil, nil, errors.ErrInvalidRequestEntity
	}

	if err = validate.Struct(&q); err != nil {
		return nil, nil, errors.ErrInvalidRequestEntity
	}
	if !q.IsValid() {
		return nil, nil, errors.ErrInvalidRequestEntity
	}

	file, err := files[0].Open()
	if err != nil {
		return nil, nil, errors.ErrInvalidFormMedia
	}

	return &q, file, nil
}
