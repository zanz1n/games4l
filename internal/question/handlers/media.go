package handlers

import "io"

type questionImageUploadData struct {
	ID        int32
	PngImage  []byte
	WebpImage []byte
	AvifImage []byte
}

func (h *QuestionHandlers) uploadImages(data *questionImageUploadData) error {
	err := h.uploadImage(data.ID, "png", []string{}, data.PngImage)
	if err != nil {
		return err
	}

	err = h.uploadImage(data.ID, "webp", []string{"png"}, data.WebpImage)
	if err != nil {
		return err
	}

	err = h.uploadImage(data.ID, "avif", []string{"png", "webp"}, data.AvifImage)
	if err != nil {
		return err
	}

	return nil
}

func (h *QuestionHandlers) encodeImageFormats(id int32, image io.Reader) (*questionImageUploadData, error) {
	var err error

	ud := questionImageUploadData{
		ID: id,
	}

	if ud.PngImage, err = h.fmpg.ScaleDownPng(image); err != nil {
		h.logger.Info("Failed to scale down image to png: " + err.Error())
		return nil, err
	}

	if ud.WebpImage, err = h.fmpg.PngToWebp(ud.PngImage); err != nil {
		h.logger.Error("Failed to convert image to webp: " + err.Error())
		return nil, err
	}

	if ud.AvifImage, err = h.fmpg.PngToAvif(ud.PngImage); err != nil {
		h.logger.Error("Failed to convert image to avif: " + err.Error())
		return nil, err
	}

	return &ud, nil
}
