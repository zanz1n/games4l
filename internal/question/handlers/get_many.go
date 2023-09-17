package handlers

import (
	"strconv"

	"github.com/games4l/internal/question"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/httpcodes"
)

func (h *QuestionHandlers) GetMany(limit int) (*utils.DataResponse[[]question.Question], error) {
	if limit > 2000 || 2 > limit {
		limit = 1000
	}

	dbc, err := h.qs.GetInstance()
	if err != nil {
		return nil, err
	}

	qs, err := dbc.GetMany(int32(limit))
	if err != nil {
		return nil, err
	}

	var msg string

	if l := len(qs); l == 0 {
		msg = "No questions returned"
	} else {
		msg = strconv.Itoa(l) + " questions fetched successfully"
	}

	return &utils.DataResponse[[]question.Question]{
		Message:    msg,
		Data:       qs,
		StatusCode: httpcodes.StatusOK,
	}, nil
}
