package main

import (
	"net/http"
	"time"
	"unsafe"

	"github.com/go-chi/chi/v5/middleware"
)

type HandleType uint8

const (
	HandlerTypeAuth            HandleType = 1 << iota
	HandlerTypeQuestion        HandleType = 1 << iota
	HandlerTypeQuestionParams  HandleType = 1 << iota
	HandlerTypeTelemetry       HandleType = 1 << iota
	HandlerTypeTelemetryParams HandleType = 1 << iota
)

func b2s(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func s2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		h.ServeHTTP(ww, r)

		LogRequest(r, ww.Status(), start)
	})
}
