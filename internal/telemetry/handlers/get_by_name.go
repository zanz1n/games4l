package handlers

import (
	"strconv"

	"github.com/games4l/internal/telemetry"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/httpcodes"
)

func (h *TelemetryHandlers) GetByName(name string) (*utils.DataResponse[[]telemetry.Registry], error) {
	dbc, err := h.ts.GetInstance()
	if err != nil {
		return nil, err
	}

	r, err := dbc.GetBySimilarName(name)
	if err != nil {
		return nil, err
	}

	var msg string
	if l := len(r); l == 0 {
		msg = "No entries returned"
	} else {
		msg = strconv.Itoa(l) + " entries fetched successfully"
	}

	return &utils.DataResponse[[]telemetry.Registry]{
		Message:    msg,
		Data:       r,
		StatusCode: httpcodes.StatusOK,
	}, nil
}
