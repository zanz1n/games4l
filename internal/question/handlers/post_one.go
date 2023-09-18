package handlers

import (
	"mime/multipart"

	"github.com/games4l/internal/question"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/httpcodes"
	"github.com/games4l/pkg/errors"
)

func (h *QuestionHandlers) PostOne(mp *multipart.Form) (*utils.DataResponse[*question.Question], error) {
	data, file, err := multipartToApiQuestion(mp)
	if err != nil {
		return nil, err
	}

	if file != nil {
		defer file.Close()
	} else {
		if data.Style == question.QuestionStyleImage {
			return nil, errors.ErrInvalidFormMedia
		}
	}

	dbc, err := h.qs.GetInstance()
	if err != nil {
		return nil, err
	}

	q, err := dbc.Create(data)
	if err != nil {
		return nil, err
	}

	if data.Style == question.QuestionStyleImage {
		encodedImages, err := h.encodeImageFormats(q.ID, file)
		if err != nil {
			return nil, err
		}

		if err = h.uploadImages(encodedImages); err != nil {
			return nil, err
		}
	}

	return &utils.DataResponse[*question.Question]{
		Message:    "Created successfully",
		Data:       q,
		StatusCode: httpcodes.StatusCreated,
	}, nil
}
