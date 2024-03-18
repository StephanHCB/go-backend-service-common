package apierrors

import (
	"github.com/StephanHCB/go-backend-service-common/api"
	"net/http"
	"time"
)

type AnnotatedError interface {
	error
	ApiError() api.ErrorDto
	ResponseObject() any
	HttpStatus() int
	Wrapped() error
}

func NewInternalServerError(message string, details string, wrapped error, timestamp time.Time) AnnotatedError {
	return create(message, details, http.StatusInternalServerError, wrapped, timestamp)
}

func IsInternalServerError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusInternalServerError)
}

func NewGatewayTimeoutError(message string, details string, wrapped error, timestamp time.Time) AnnotatedError {
	return create(message, details, http.StatusGatewayTimeout, wrapped, timestamp)
}

func IsGatewayTimeoutError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusGatewayTimeout)
}

func NewBadGatewayError(message string, details string, wrapped error, timestamp time.Time) AnnotatedError {
	return create(message, details, http.StatusBadGateway, wrapped, timestamp)
}

func IsBadGatewayError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusBadGateway)
}

func NewBadRequestError(message string, details string, wrapped error, timestamp time.Time) AnnotatedError {
	return create(message, details, http.StatusBadRequest, wrapped, timestamp)
}

func IsBadRequestError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusBadRequest)
}

func NewConflictError(message string, details string, wrapped error, timestamp time.Time) AnnotatedError {
	return create(message, details, http.StatusConflict, wrapped, timestamp)
}

func NewConflictErrorWithResponse(message string, details string, wrapped error, response any, timestamp time.Time) AnnotatedError {
	return createWithResponse(message, details, http.StatusConflict, wrapped, response, timestamp)
}

func IsConflictError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusConflict)
}

func NewNotFoundError(message string, details string, wrapped error, timestamp time.Time) AnnotatedError {
	return create(message, details, http.StatusNotFound, wrapped, timestamp)
}

func IsNotFoundError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusNotFound)
}

func NewUnprocessableEntity(message string, details string, wrapped error, timestamp time.Time) AnnotatedError {
	return create(message, details, http.StatusUnprocessableEntity, wrapped, timestamp)
}

func IsUnprocessableEntity(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusUnprocessableEntity)
}

func NewUnauthorisedError(message string, details string, wrapped error, timestamp time.Time) AnnotatedError {
	return create(message, details, http.StatusUnauthorized, wrapped, timestamp)
}

func IsUnauthorisedError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusUnauthorized)
}

func NewForbiddenError(message string, details string, wrapped error, timestamp time.Time) AnnotatedError {
	return create(message, details, http.StatusForbidden, wrapped, timestamp)
}

func IsForbiddenError(err error) bool {
	return isAnnotatedErrorWithStatus(err, http.StatusForbidden)
}
