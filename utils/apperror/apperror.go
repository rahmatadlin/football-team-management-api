package apperror

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func BadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func Internal(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg}
}

func HTTPStatus(err error) int {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae.Code
	}
	return http.StatusInternalServerError
}

func Message(err error) string {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae.Message
	}
	return "internal server error"
}
