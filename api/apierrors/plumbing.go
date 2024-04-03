package apierrors

import (
	"errors"
	"github.com/StephanHCB/go-backend-service-common/api"
	"time"
)

type AnnotatedErrorImpl struct {
	VApiError       api.ErrorDto
	VResponseObject any
	VHttpStatus     int
	VWrapped        error
}

func (e *AnnotatedErrorImpl) ResponseObject() any {
	return e.VResponseObject
}

//goland:noinspection GoUnusedFunction
func implementsInterfaces() (error, AnnotatedError) {
	return &AnnotatedErrorImpl{}, &AnnotatedErrorImpl{}
}

func (e *AnnotatedErrorImpl) Error() string {
	if e.VApiError.Message == nil {
		return "error message not provided - this is an implementation error"
	}
	return *e.VApiError.Message
}

func (e *AnnotatedErrorImpl) ApiError() api.ErrorDto {
	return e.VApiError
}

func (e *AnnotatedErrorImpl) HttpStatus() int {
	return e.VHttpStatus
}

func (e *AnnotatedErrorImpl) Wrapped() error {
	return e.VWrapped
}

func create(message string, details string, status int, wrapped error, timestamp time.Time) AnnotatedError {
	return &AnnotatedErrorImpl{
		VApiError: api.ErrorDto{
			Details:   &details,
			Message:   &message,
			Timestamp: &timestamp,
		},
		VHttpStatus: status,
		VWrapped:    wrapped,
	}
}

func createWithResponse(message string, details string, status int, wrapped error, response any, timestamp time.Time) AnnotatedError {
	return &AnnotatedErrorImpl{
		VApiError: api.ErrorDto{
			Details:   &details,
			Message:   &message,
			Timestamp: &timestamp,
		},
		VResponseObject: response,
		VHttpStatus:     status,
		VWrapped:        wrapped,
	}
}

func isAnnotatedErrorWithStatus(err error, status int) bool {
	var ann AnnotatedError
	ok := errors.As(err, &ann)
	if !ok {
		return false
	}
	return ann.HttpStatus() == status
}
