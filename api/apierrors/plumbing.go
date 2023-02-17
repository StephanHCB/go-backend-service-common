package apierrors

type AnnotatedErrorImpl struct {
	VApiError   ApiErrorV2
	VHttpStatus int
	VWrapped    error
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

func (e *AnnotatedErrorImpl) ApiError() ApiErrorV2 {
	return e.VApiError
}

func (e *AnnotatedErrorImpl) HttpStatus() int {
	return e.VHttpStatus
}

func (e *AnnotatedErrorImpl) Wrapped() error {
	return e.VWrapped
}

func create(message string, status int, wrapped error) AnnotatedError {
	errCode := int32(status)
	return &AnnotatedErrorImpl{
		VApiError: ApiErrorV2{
			ErrorCode: &errCode,
			Message:   &message,
		},
		VHttpStatus: status,
		VWrapped:    wrapped,
	}
}

func isAnnotatedErrorWithStatus(err error, status int) bool {
	ann, ok := err.(AnnotatedError)
	if !ok {
		return false
	}
	return ann.HttpStatus() == status
}
