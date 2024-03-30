package playht

import (
	"encoding/json"
	"errors"
)

var (
	ErrTooManyRequests = errors.New("too many requests")
)

// APIError is an API error.
type APIError struct {
	ID      string `json:"error_id"`
	Message string `json:"error_message"`
}

// Error implements error interface.
func (e APIError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "unknown error"
	}
	return string(b)
}

// APIErrInternal is an internal API error.
type APIErrInternal struct {
	Message string `json:"message"`
	Err     string `json:"error"`
}

// Error implements error interface.
func (e APIErrInternal) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "unknown error"
	}
	return string(b)
}
