package handlers

import (
	"github.com/games4l/internal/telemetry"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/httpcodes"
	"github.com/games4l/pkg/errors"
	"github.com/goccy/go-json"
)

func (h *TelemetryHandlers) PostOne(payload []byte) (*utils.DataResponse[*telemetry.Registry], error) {
	dbc, err := h.ts.GetInstance()
	if err != nil {
		return nil, err
	}

	cd := telemetry.CreateRegistryData{}

	if err = json.Unmarshal(payload, &cd); err != nil {
		return nil, errors.ErrMalformedOrTooBigBody
	}
	if err = validate.Struct(&cd); err != nil {
		return nil, errors.ErrInvalidRequestEntity
	}
	if !cd.IsValid() {
		return nil, errors.ErrInvalidRequestEntity
	}

	r, err := dbc.Create(&cd)
	if err != nil {
		return nil, err
	}

	return &utils.DataResponse[*telemetry.Registry]{
		Message:    "Created successfully",
		Data:       r,
		StatusCode: httpcodes.StatusCreated,
	}, nil
}
