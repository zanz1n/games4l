package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	authsrc "github.com/games4l/cmd/lambda_auth/src"
	"github.com/go-chi/chi/v5"
)

func HandleRequest(ht HandleType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := ConvertRequest(r)
		if err != nil {
			HandleError(err, w)
			return
		}

		var response events.APIGatewayProxyResponse

		switch ht {
		case HandlerTypeAuth:
			response, err = authsrc.Handler(*req)

		case HandlerTypeQuestion:
			response, err = questionServer.RequestHandler(*req)

		case HandlerTypeQuestionParams:
			req.PathParameters["id"] = chi.URLParam(r, "id")
			response, err = questionServer.RequestHandler(*req)

		case HandlerTypeTelemetry:
			response, err = telemetryServer.RequestHandler(*req)

		case HandlerTypeTelemetryParams:
			req.PathParameters["id"] = chi.URLParam(r, "id")
			response, err = telemetryServer.RequestHandler(*req)

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

		for k, v := range response.Headers {
			w.Header().Add(k, v)
		}
	}
}
