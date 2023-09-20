package handlers

import (
	"bytes"
	"fmt"

	"github.com/games4l/pkg/errors"
)

func (h *QuestionHandlers) uploadImage(id int32, ext string, fallbackExts []string, buf []byte) error {
	err := h.bucket.Store(
		fmt.Sprintf("%s%d.%s", h.iap, id, ext),
		"image/"+ext,
		bytes.NewReader(buf),
	)

	if err != nil {
		dbc, err := h.qs.GetInstance()
		if err != nil {
			return err
		}

		dbc.DeleteById(id)

		if len(fallbackExts) > 0 {
			objects := make([]string, len(fallbackExts))

			for i, v := range fallbackExts {
				objects[i] = fmt.Sprintf("%s%d.%s", h.iap, id, v)
			}

			err = h.bucket.DestroyMany(objects)
			if err != nil {
				h.logger.Error("Failed to delete already uploaded s3 images: " + err.Error())
			}
		}

		return errors.ErrInternalServerError
	}

	return nil
}
