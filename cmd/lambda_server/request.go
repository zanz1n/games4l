package main

import (
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/games4l/internal/logger"
)

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
		Headers:                         ConvertMultivalueMap(r.Header, true),
		QueryStringParameters:           ConvertMultivalueMap(query, false),
		MultiValueQueryStringParameters: query,
		PathParameters:                  make(map[string]string),
		RequestContext:                  *new(events.APIGatewayProxyRequestContext),
		StageVariables:                  make(map[string]string),
	}, nil
}

func ConvertMultivalueMap(multivalue map[string][]string, lowerfy bool) map[string]string {
	m := make(map[string]string)

	for k, v := range multivalue {
		s := ""

		for _, sv := range v {
			s += sv
		}

		if lowerfy {
			m[strings.ToLower(k)] = s
		} else {
			m[k] = s
		}
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
		logger.Error("Failed to parse request ip: " + err.Error())
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
