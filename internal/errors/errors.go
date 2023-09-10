package errors

import (
	"github.com/goccy/go-json"
)

func GetStatusErr(key error) StatusError {
	v, ok := mpe[key]

	if !ok {
		return &statusErrorImpl{
			code:     50000,
			httpCode: 500,
			message:  "Unknown err: " + key.Error(),
		}
	}

	return v
}

type statusErrorImpl struct {
	code     uint
	httpCode int
	message  string
}

func (e *statusErrorImpl) Message() string {
	return e.message
}

func (e *statusErrorImpl) CustomCode() uint {
	return e.code
}

func (e *statusErrorImpl) HttpCode() int {
	return e.httpCode
}

type StatusError interface {
	Message() string
	CustomCode() uint
	HttpCode() int
}

type errorImpl struct {
	m string
}

func (e *errorImpl) Error() string {
	return e.m
}

func New(text string) error {
	return &errorImpl{m: text}
}

type ErrorBody struct {
	Message   string `json:"message"`
	ErrorCode uint   `json:"errorCode"`
}

func (e *ErrorBody) Marshal() []byte {
	buf, err := json.Marshal(e)

	if err != nil {
		return []byte("{\"message\":\"Failed to marshal response body\",\"errorCode\":5000}")
	}

	return buf
}
