package handlers

import (
	"github.com/games4l/internal/telemetry"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/httpcodes"
	"github.com/games4l/pkg/errors"
)

func (h *TelemetryHandlers) GetById(id string, authed bool) (*utils.DataResponse[*telemetry.Registry], error) {
	if id == "" || 10 > len(id) {
		return nil, errors.ErrInvalidObjectID
	}

	dbc, err := h.ts.GetInstance()
	if err != nil {
		return nil, err
	}

	r, err := dbc.GetById(id)
	if err != nil {
		return nil, err
	}

	var msg string
	if !authed {
		r.PacientName = "<CENSURADO> - autentique-se para ver"
		msg = "Fetched successfully. Please auth to see the pacient name"
	} else {
		msg = "Fetched successfully"
	}

	return &utils.DataResponse[*telemetry.Registry]{
		Message:    msg,
		Data:       r,
		StatusCode: httpcodes.StatusOK,
	}, nil
}
