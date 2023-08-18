package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/backend/libs/logger"
	authsrc "github.com/games4l/backend/services/auth_lambda/src"
	questionsrc "github.com/games4l/backend/services/question_lambda/src"
	telemetrysrc "github.com/games4l/backend/services/telemetry_lambda/src"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("APP_ENV") == "" {
		if err := godotenv.Load(); err != nil {
			logger.Error("Failed to load .env file: " + err.Error())
		} else {
			logger.Info("Loaded environment configuration from .env file")
		}
	} else {
		logger.Info("Loaded environment configuration from shell variables")
	}
}

func main() {
	server := Server{}

	http.ListenAndServe(":8080", &server)
}

type Server struct{}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	req, err := ConvertRequest(r)
	if err != nil {
		HandleError(err, w)
		return
	}

	var response events.APIGatewayProxyResponse

	switch {
	case strings.HasPrefix(r.URL.Path, "/telemetry"):
		response, err = telemetrysrc.Handler(*req)
	case strings.HasPrefix(r.URL.Path, "/auth"):
		response, err = authsrc.Handler(*req)
	case strings.HasPrefix(r.URL.Path, "/user"):
		response, err = authsrc.Handler(*req)
	case strings.HasPrefix(r.URL.Path, "/question"):
		response, err = questionsrc.Handler(*req)
	default:
		w.WriteHeader(404)
		w.Write(s2b("Not Found"))
		return
	}

	if err != nil {
		HandleError(err, w)
		return
	}

	w.WriteHeader(response.StatusCode)
	w.Write(s2b(response.Body))

	for k, v := range response.MultiValueHeaders {
		s := ""
		for _, part := range v {
			s += ", " + part
		}

		w.Header().Add(k, s)
	}

	for k, v := range response.Headers {
		w.Header().Add(k, v)
	}

	LogRequest(r, response.StatusCode, start)
}

func ConvertRequest(r *http.Request) (*events.APIGatewayProxyRequest, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	query := r.URL.Query()

	return &events.APIGatewayProxyRequest{
		Body:                            b2s(bodyBytes),
		Path:                            r.URL.Path,
		Resource:                        r.URL.Host,
		HTTPMethod:                      r.Method,
		IsBase64Encoded:                 false,
		MultiValueHeaders:               r.Header,
		Headers:                         ConvertMultivalueMap(r.Header),
		QueryStringParameters:           ConvertMultivalueMap(query),
		MultiValueQueryStringParameters: query,
		PathParameters:                  make(map[string]string),
		RequestContext:                  *new(events.APIGatewayProxyRequestContext),
		StageVariables:                  make(map[string]string),
	}, nil
}

func ConvertMultivalueMap(multivalue map[string][]string) map[string]string {
	m := make(map[string]string)

	for k, v := range multivalue {
		s := ""

		for _, sv := range v {
			s += ", " + sv
		}

		m[k] = s
	}

	return m
}

func HandleError(e error, w http.ResponseWriter) {
	w.Write(s2b(e.Error()))
	w.WriteHeader(500)
}

func LogRequest(r *http.Request, status int, start time.Time) {
	addr, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
	if err != nil {
		return
	}

	logger.LogRequest(&logger.RequestInfo{
		Addr:       addr,
		Method:     r.Method,
		Path:       r.URL.Path,
		StatusCode: status,
		Duration:   time.Since(start),
	})
}

func b2s(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func s2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
