package domain

import (
	"fmt"

	"github.com/Rizz404/inventory-api/internal/utils"
)

type AppError struct {
	Code       int
	Message    string
	MessageKey utils.MessageKey // * For i18n support
	Params     []string         // * Parameters for message formatting
	Err        error
}

func (e *AppError) Error() string {
	return e.Message
}

// * GetLocalizedMessage returns the localized error message
func (e *AppError) GetLocalizedMessage(langCode string) string {
	if e.MessageKey != "" {
		return utils.GetLocalizedMessage(e.MessageKey, langCode, e.Params...)
	}
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

// * NewAppErrorWithKey creates a new AppError with i18n support
func NewAppErrorWithKey(code int, messageKey utils.MessageKey, params []string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    utils.GetLocalizedMessage(messageKey, "en-US", params...), // * Default English message
		MessageKey: messageKey,
		Params:     params,
		Err:        err,
	}
}

func ErrBadRequest(message string) *AppError {
	return NewAppError(400, message, nil)
}

// * ErrBadRequestWithKey creates a bad request error with i18n support
func ErrBadRequestWithKey(messageKey utils.MessageKey, params ...string) *AppError {
	return NewAppErrorWithKey(400, messageKey, params, nil)
}

func ErrUnauthorized(message string) *AppError {
	return NewAppError(401, message, nil)
}

// * ErrUnauthorizedWithKey creates an unauthorized error with i18n support
func ErrUnauthorizedWithKey(messageKey utils.MessageKey, params ...string) *AppError {
	return NewAppErrorWithKey(401, messageKey, params, nil)
}

func ErrForbidden(message string) *AppError {
	return NewAppError(403, message, nil)
}

// * ErrForbiddenWithKey creates a forbidden error with i18n support
func ErrForbiddenWithKey(messageKey utils.MessageKey, params ...string) *AppError {
	return NewAppErrorWithKey(403, messageKey, params, nil)
}

func ErrNotFound(entity string) *AppError {
	return NewAppError(404, fmt.Sprintf("%s not found", entity), nil)
}

// * ErrNotFoundWithKey creates a not found error with i18n support
func ErrNotFoundWithKey(messageKey utils.MessageKey, params ...string) *AppError {
	return NewAppErrorWithKey(404, messageKey, params, nil)
}

func ErrConflict(message string) *AppError {
	return NewAppError(409, message, nil)
}

// * ErrConflictWithKey creates a conflict error with i18n support
func ErrConflictWithKey(messageKey utils.MessageKey, params ...string) *AppError {
	return NewAppErrorWithKey(409, messageKey, params, nil)
}

func ErrInternal(err error) *AppError {
	return NewAppError(500, "an unexpected internal error occured", err)
}

// * ErrInternalWithKey creates an internal error with i18n support
func ErrInternalWithKey(messageKey utils.MessageKey, err error, params ...string) *AppError {
	return NewAppErrorWithKey(500, messageKey, params, err)
}
