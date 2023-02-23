package apierrors

import (
	"net/http"
)

type ApiErrorV2 struct {
	ErrorCode *int32  `json:"errorCode,omitempty"`
	Message   *string `json:"message,omitempty"`
}

type AnnotatedError interface {
	error
	ApiError() ApiErrorV2
	HttpStatus() int
	Wrapped() error
}

func NewInternalServerError(message string, wrapped error) AnnotatedError {
	return create(message, http.StatusInternalServerError, wrapped)
}

func IsInternalServerError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusInternalServerError)
}

func NewGatewayTimeoutError(message string, wrapped error) AnnotatedError {
	return create(message, http.StatusGatewayTimeout, wrapped)
}

func IsGatewayTimeoutError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusGatewayTimeout)
}

func NewBadGatewayError(message string, wrapped error) AnnotatedError {
	return create(message, http.StatusBadGateway, wrapped)
}

func IsBadGatewayError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusBadGateway)
}

func NewBadRequestError(message string, wrapped error) AnnotatedError {
	return create(message, http.StatusBadRequest, wrapped)
}

func IsBadRequestError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusBadRequest)
}
