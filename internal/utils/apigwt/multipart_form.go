package apigwt

import (
	"bytes"
	"encoding/base64"
	"mime"
	"mime/multipart"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/logger"
	"github.com/games4l/pkg/errors"
)

func ParseMultipartForm(req *events.APIGatewayProxyRequest, maxMemSize int64) (*multipart.Form, error) {
	ct, ok := req.Headers["Content-Type"]

	if !ok {
		ct, ok = req.Headers["content-type"]
	}

	if !ok || ct == "" {
		return nil, errors.ErrNoMultipartForm
	}

	d, params, err := mime.ParseMediaType(ct)
	if err != nil || d != "multipart/form-data" {
		return nil, errors.ErrNoMultipartForm
	}

	boundary, ok := params["boundary"]
	if !ok {
		return nil, errors.ErrNoMultipartForm
	}

	buf, err := base64.StdEncoding.DecodeString(req.Body)
	if err != nil {
		return nil, errors.ErrNoMultipartForm
	}

	mpr := multipart.NewReader(
		bytes.NewReader(buf),
		boundary,
	)

	mp, err := mpr.ReadForm(maxMemSize)
	if err != nil {
		logger.Debug(err.Error())
		return nil, errors.ErrFailedToReadMultipartForm
	}

	return mp, nil
}
