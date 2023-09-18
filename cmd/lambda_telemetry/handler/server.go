package handler

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/auth"
	"github.com/games4l/internal/telemetry/handlers"
	"github.com/games4l/internal/telemetry/repository"
	"github.com/games4l/internal/utils"
	"github.com/games4l/pkg/errors"
)

func NewEnvServer() *Server {
	ap := auth.NewAuthProvider([]byte(os.Getenv("WEBHOOK_SIG")), []byte(os.Getenv("JWT_SIG")))

	uri := os.Getenv("MONGO_URI")
	dbname := os.Getenv("MONGO_DATABASE_NAME")

	ts := repository.NewMongoSingleton(uri, dbname)

	return &Server{
		h:      handlers.NewTelemetryHandlers(ts),
		ap:     ap,
		prefix: os.Getenv("API_GATEWAY_PREFIX"),
	}
}

type Server struct {
	h      *handlers.TelemetryHandlers
	ap     *auth.AuthProvider
	prefix string
}

func (s *Server) RequestHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		fErr error
		res  *events.APIGatewayProxyResponse
	)

	if req.Path == "/"+s.prefix+"/telemetry" || req.Path == "/telemetry" {
		if req.HTTPMethod == "POST" {
			res, fErr = s.HandlePostOne(req)
		} else if req.HTTPMethod == "GET" {
			res, fErr = s.HandleGetByName(req)
		} else {
			fErr = errors.ErrMethodNotAllowed
		}
	} else if hPrx(req.Path, "/"+s.prefix+"/telemetry/") || hPrx(req.Path, "/telemetry") {
		if req.HTTPMethod == "GET" {
			res, fErr = s.HandleGetById(req)
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

func hPrx(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

var applicationJsonHeader = map[string]string{
	"Content-Type": "application/json",
}
