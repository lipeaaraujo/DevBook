package apierrors

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func New(status int, message string, err error) *APIError {
	return &APIError{Status: status, Message: message, Err: err}
}

func UnprocessableEntity(msg string) *APIError {
	return New(http.StatusUnprocessableEntity, msg, nil)
}

func NotFound(resource string) *APIError {
	return New(http.StatusNotFound, fmt.Sprintf("%s not found", resource), nil)
}

func BadRequest(msg string) *APIError {
	return New(http.StatusBadRequest, msg, nil)
}

func Internal(err error) *APIError {
	return New(http.StatusInternalServerError, "Unexpected error", err)
}

func Forbidden(msg string) *APIError {
	return New(http.StatusForbidden, msg, nil)
}

func Unauthorized(msg string) *APIError {
	return New(http.StatusUnauthorized, msg, nil)
}

func ValidationError(msg string) *APIError {
	return New(http.StatusUnprocessableEntity, msg, nil)
}

func ResourceAlreadyExists(resource string) *APIError {
	return New(http.StatusBadRequest, fmt.Sprintf("%s already exists", resource), nil)
}
