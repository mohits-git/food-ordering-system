package utils

type AppError struct {
	Code    AppErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

type AppErrorCode int

const (
	ErrNotFound AppErrorCode = iota + 1
	ErrUnauthorized
	ErrForbidden
	ErrConflict
	ErrInvalid
	ErrInternal
	ErrTimeout
)

func NewAppError(code AppErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
