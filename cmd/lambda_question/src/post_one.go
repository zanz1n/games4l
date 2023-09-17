package src

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/games4l/internal/logger"
	"github.com/games4l/internal/question"
	"github.com/games4l/internal/question/entityconv"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/apigwt"
	"github.com/games4l/internal/utils/httpcodes"
	"github.com/games4l/pkg/errors"
	"github.com/games4l/pkg/ffmpeg"
	"github.com/goccy/go-json"
)

func upload(id int32, ext string, fallbackExts []string, buf []byte) error {
	var err error

	ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel1()

	_, err = s3c.PutObject(ctx1, &s3.PutObjectInput{
		Bucket:      aws.String(s3Bucket),
		Key:         aws.String("questions/images/" + strconv.Itoa(int(id)) + "." + ext),
		ContentType: aws.String("image/" + ext),
		Body:        bytes.NewReader(buf),
	})
	if err != nil {
		logger.Error("Failed to upload to s3 bucket: " + err.Error())

		ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel2()

		dba.DeleteQuestionById(ctx2, id)

		if len(fallbackExts) > 0 {
			objects := make([]types.ObjectIdentifier, len(fallbackExts))

			for i, v := range fallbackExts {
				objects[i] = types.ObjectIdentifier{
					Key: aws.String("questions/images/" + strconv.Itoa(int(id)) + "." + v),
				}
			}

			_, err = s3c.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
				Bucket: aws.String(s3Bucket),
				Delete: &types.Delete{
					Objects: objects,
				},
			})
			if err != nil {
				logger.Error("Failed to delete already uploaded s3 images: " + err.Error())
			}
		}

		return errors.ErrInternalServerError
	}

	return nil
}

var p = ffmpeg.NewProvider("/tmp", "ffmpeg", "ffprobe")

func HandlePostOne(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	err := ap.AuthenticateAdminHeader(
		req.Headers["authorization"],
		utils.S2B(req.Body),
	)
	if err != nil {
		return nil, err
	}

	// 10MB = 1024 * 1024 * 10 bytes = 2^10 * 2^10 * 10 = 10^20 * 10
	mp, err := apigwt.ParseMultipartForm(&req, (2<<20)*10)
	if err != nil {
		return nil, err
	}

	data, file, err := MultipartToApiQuestion(mp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err = Connect(); err != nil {
		return nil, err
	}
	if err = ConnectS3(); err != nil {
		return nil, err
	}

	png, err := p.ScaleDownPng(file)
	if err != nil {
		logger.Debug("Failed to scale down png")
		return nil, err
	}

	webp, err := p.PngToWebp(png)
	if err != nil {
		logger.Debug("Failed to convert to webp")
		return nil, err
	}

	avif, err := p.PngToAvif(png)
	if err != nil {
		logger.Debug("Failed to convert to avif")
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := dba.CreateQuestion(
		ctx,
		entityconv.CreateQuestionParamsToDbEntity(data),
	)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	if err = upload(result.ID, "png", []string{}, png); err != nil {
		return nil, err
	}
	if err = upload(result.ID, "webp", []string{"png"}, webp); err != nil {
		return nil, err
	}
	if err = upload(result.ID, "avif", []string{"png", "webp"}, avif); err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      httpcodes.StatusOK,
		Headers:         applicationJsonHeader,
		IsBase64Encoded: false,
		Body: utils.MarshalJSON(JSON{
			"message": "Created successfully",
			"data":    entityconv.QuestionToApiEntity(result),
		}),
	}, nil
}

func MultipartToApiQuestion(
	mp *multipart.Form,
) (*question.QuestionCreateData, io.ReadCloser, error) {
	files, ok := mp.File["file"]
	if !ok || len(files) != 1 {
		return nil, nil, errors.ErrInvalidFormMedia
	}

	rawData, ok := mp.Value["data"]

	q := question.QuestionCreateData{}

	err := json.Unmarshal(
		utils.S2B(strings.Join(rawData, "")),
		&q,
	)
	if err != nil {
		return nil, nil, errors.ErrInvalidRequestEntity
	}

	if err = validate.Struct(&q); err != nil {
		return nil, nil, errors.ErrInvalidRequestEntity
	}
	if !q.IsValid() {
		return nil, nil, errors.ErrInvalidRequestEntity
	}

	file, err := files[0].Open()
	if err != nil {
		return nil, nil, errors.ErrInvalidFormMedia
	}

	return &q, file, nil
}
