package playht

import (
	"encoding/json"
	"errors"
)

var (
	// ErrTooManyRequests is returned when the cient hits rate limit.
	ErrTooManyRequests = errors.New("too many requests")
	// ErrUnexpectedStatusCode is returned when an unexpected status is returned from the API.
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

// APIErrGen is a generic API error.
type APIErrGen struct {
	ID      string `json:"error_id"`
	Message string `json:"error_message"`
}

// Error implements error interface.
func (e APIErrGen) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "unknown error"
	}
	return string(b)
}

// APIErrInternal is an error returned
// when the API responds with 50x status code.
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
