package handler

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/auth"
	"github.com/games4l/internal/question/handlers"
	"github.com/games4l/internal/question/repository"
	"github.com/games4l/internal/utils"
	"github.com/games4l/internal/utils/s3u"
	"github.com/games4l/pkg/errors"
	"github.com/games4l/pkg/ffmpeg"
)

func NewEnvServer() *Server {
	qs := repository.NewPostgresSingleton(os.Getenv("DATABASE_URL"))
	sc := s3u.NewS3Singleton("sa-east-1")
	fmp := ffmpeg.NewProvider("/tmp/", "ffmpeg", "ffprobe")

	h := handlers.NewQuestionHandlers(qs, sc, fmp, os.Getenv("APP_QUESTION_BUCKET_NAME"), "question/images/")

	ap := auth.NewAuthProvider([]byte(os.Getenv("WEBHOOK_SIG")), []byte(os.Getenv("JWT_SIG")))

	return &Server{
		h:      h,
		ap:     ap,
		prefix: os.Getenv("API_GATEWAY_PREFIX"),
	}
}

type Server struct {
	h      *handlers.QuestionHandlers
	ap     *auth.AuthProvider
	prefix string
}

func (s *Server) RequestHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		fErr error
		res  *events.APIGatewayProxyResponse
	)

	if strings.HasPrefix(req.Path, "/question") || strings.HasPrefix(req.Path, "/"+s.prefix+"/question") {
		if req.HTTPMethod == "GET" {
			if _, ok := req.PathParameters["id"]; ok {
				res, fErr = s.HandleGetOne(req)
			} else {
				res, fErr = s.HandleGetMany(req)
			}
		} else if req.HTTPMethod == "POST" {
			res, fErr = s.HandlePostOne(req)
		} else if req.HTTPMethod == "PATCH" || req.HTTPMethod == "PUT" {
			fErr = errors.ErrMethodNotAllowed
			// update one route
		} else {
			fErr = errors.ErrMethodNotAllowed
		}
	} else {
		fErr = errors.ErrNoSuchRoute
	}

	if fErr != nil {
		errInfo := errors.GetStatusErr(fErr)
		errBody := errors.ErrorBody{
			Message:   errInfo.Message(),
			ErrorCode: errInfo.CustomCode(),
		}

		res = &events.APIGatewayProxyResponse{
			StatusCode:      errInfo.HttpCode(),
			Headers:         applicationJsonHeader,
			IsBase64Encoded: false,
			Body:            utils.B2S(errBody.Marshal()),
		}
	}

	return *res, nil
}

var applicationJsonHeader = map[string]string{
	"Content-Type": "application/json",
}
