package handlers

import (
	"math"

	"github.com/games4l/internal/question"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/httpcodes"
	"github.com/games4l/pkg/errors"
)

func (h *QuestionHandlers) GetOne(id int) (*utils.DataResponse[*question.Question], error) {
	if id > math.MaxInt32 {
		return nil, errors.ErrInvalidIntegerIdPathParam
	}

	i, err := h.qs.GetInstance()
	if err != nil {
		return nil, err
	}

	q, err := i.GetById(int32(id))
	if err != nil {
		return nil, err
	}

	return &utils.DataResponse[*question.Question]{
		Message:    "Fetched successfully",
		Data:       q,
		StatusCode: httpcodes.StatusOK,
	}, nil
}
