package domain

import "fmt"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func ErrBadRequest(message string) *AppError {
	return NewAppError(400, message, nil)
}

func ErrUnauthorized(message string) *AppError {
	return NewAppError(401, message, nil)
}

func ErrForbidden(message string) *AppError {
	return NewAppError(403, message, nil)
}

func ErrNotFound(entity string) *AppError {
	return NewAppError(404, fmt.Sprintf("%s not found", entity), nil)
}

func ErrConflict(message string) *AppError {
	return NewAppError(409, message, nil)
}

func ErrInternal(err error) *AppError {
	return NewAppError(500, "an unexpected internal error occured", err)
}
