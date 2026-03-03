package apperror

import (
	"fmt"
	"net/http"
)

// Code represents a machine-readable error code.
type Code string

const (
	CodeNotFound          Code = "NOT_FOUND"
	CodeConflict          Code = "CONFLICT"
	CodeValidation        Code = "VALIDATION_ERROR"
	CodeUnauthorized      Code = "UNAUTHORIZED"
	CodeForbidden         Code = "FORBIDDEN"
	CodeInternal          Code = "INTERNAL_ERROR"
	CodeBadRequest        Code = "BAD_REQUEST"
	CodeConcurrencyConflict Code = "CONCURRENCY_CONFLICT"
)

// Error is the application-level error type used across all layers.
type Error struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// HTTPStatus maps error codes to HTTP status codes.
func (e *Error) HTTPStatus() int {
	switch e.Code {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict, CodeConcurrencyConflict:
		return http.StatusConflict
	case CodeValidation, CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func NewNotFound(entity, id string) *Error {
	return &Error{
		Code:    CodeNotFound,
		Message: fmt.Sprintf("%s with id '%s' not found", entity, id),
	}
}

func NewValidation(message string) *Error {
	return &Error{
		Code:    CodeValidation,
		Message: message,
	}
}

func NewConflict(message string) *Error {
	return &Error{
		Code:    CodeConflict,
		Message: message,
	}
}

func NewUnauthorized(message string) *Error {
	return &Error{
		Code:    CodeUnauthorized,
		Message: message,
	}
}

func NewForbidden(message string) *Error {
	return &Error{
		Code:    CodeForbidden,
		Message: message,
	}
}

func NewInternal(err error) *Error {
	return &Error{
		Code:    CodeInternal,
		Message: "an internal error occurred",
		Err:     err,
	}
}

func NewBadRequest(message string) *Error {
	return &Error{
		Code:    CodeBadRequest,
		Message: message,
	}
}

func NewConcurrencyConflict(entity string) *Error {
	return &Error{
		Code:    CodeConcurrencyConflict,
		Message: fmt.Sprintf("%s has been modified by another user; please reload and retry", entity),
	}
}
