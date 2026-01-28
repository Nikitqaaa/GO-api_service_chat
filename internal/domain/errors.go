package domain

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("already exists")
)

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func NewAPIError(status int, message string, err error) APIError {
	apiErr := APIError{
		Status:  status,
		Message: message,
	}
	if err != nil {
		apiErr.Error = err.Error()
	}
	return apiErr
}
