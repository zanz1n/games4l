package handlers

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/games4l/pkg/errors"
)

func (h *QuestionHandlers) uploadImage(id int32, ext string, fallbackExts []string, buf []byte) error {
	s3c, err := h.ss.GetInstance()
	if err != nil {
		return err
	}

	err = h.upload(
		"questions/images/"+strconv.Itoa(int(id))+"."+ext,
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
			objects := make([]types.ObjectIdentifier, len(fallbackExts))

			for i, v := range fallbackExts {
				objects[i] = types.ObjectIdentifier{
					Key: aws.String("questions/images/" + strconv.Itoa(int(id)) + "." + v),
				}
			}

			_, err = s3c.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
				Bucket: &h.s3Bucket,
				Delete: &types.Delete{
					Objects: objects,
				},
			})
			if err != nil {
				h.logger.Error("Failed to delete already uploaded s3 images: " + err.Error())
			}
		}

		return errors.ErrInternalServerError
	}

	return nil
}

func (h *QuestionHandlers) upload(key, mime string, r io.Reader) error {
	s3c, err := h.ss.GetInstance()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = s3c.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &h.s3Bucket,
		Key:         &key,
		ContentType: &mime,
		Body:        r,
	})
	if err != nil {
		h.logger.Error("Failed to upload s3 object: " + err.Error())
		return errors.ErrInternalServerError
	}

	return nil
}
