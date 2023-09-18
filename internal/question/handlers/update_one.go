package handlers

import (
	"math"

	"github.com/games4l/internal/question"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/httpcodes"
	"github.com/games4l/pkg/errors"
	"github.com/goccy/go-json"
)

func (h *QuestionHandlers) UpdateOne(id int, body []byte) (*utils.DataResponse[*question.Question], error) {
	if id > math.MaxInt32 {
		return nil, errors.ErrInvalidIntegerIdPathParam
	}

	data := question.QuestionUpdateData{}

	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.ErrMalformedOrTooBigBody
	}

	if err = validate.Struct(&data); err != nil {
		return nil, errors.ErrInvalidRequestEntity
	} else if !data.IsValid() {
		return nil, errors.ErrInvalidRequestEntity
	}

	dbc, err := h.qs.GetInstance()
	if err != nil {
		return nil, err
	}

	q, err := dbc.UpdateById(int32(id), &data)
	if err != nil {
		return nil, err
	}

	return &utils.DataResponse[*question.Question]{
		Message:    "Updated successfully",
		Data:       q,
		StatusCode: httpcodes.StatusOK,
	}, nil
}
