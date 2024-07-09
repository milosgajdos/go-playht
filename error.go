package playht

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrUnknown is returned when an unknown error occurrs.
	ErrUnknown = errors.New("unknown error")
)

// APIError is a pseudo-sum type API error type.
type APIError struct {
	Generic        *ErrGeneric
	Internal       *ErrInternal
	RateLimit      *ErrRateLimit
	UnexpecedError json.RawMessage
}

func (e *APIError) Error() string {
	if e.Generic != nil {
		return e.Generic.Error()
	}
	if e.Internal != nil {
		return e.Internal.Error()
	}
	if e.RateLimit != nil {
		return e.RateLimit.Error()
	}
	if len(e.UnexpecedError) > 0 {
		return string(e.UnexpecedError)
	}
	return ErrUnknown.Error()
}

func (e *APIError) UnmarshalJSON(data []byte) error {
	if strings.Contains(string(data), "Rate limit exceeded") {
		e.RateLimit = &ErrRateLimit{Message: string(data)}
		return nil
	}

	var genError ErrGeneric
	err := json.Unmarshal(data, &genError)
	if err == nil && genError.Message != "" {
		e.Generic = &genError
		return nil
	}

	var internalErr ErrInternal
	err = json.Unmarshal(data, &internalErr)
	if err == nil && len(internalErr.Message) > 0 {
		e.Internal = &internalErr
		return nil
	}

	e.UnexpecedError = data

	return fmt.Errorf("failed decoding API error: %s", string(data))
}

// ErrGeneric is a generic API error.
type ErrGeneric struct {
	ID      string `json:"error_id"`
	Message string `json:"error_message"`
}

// Error implements error interface.
func (e ErrGeneric) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "unknown error"
	}
	return string(b)
}

// ErrInternal is an error returned
// when the API responds with 50x status code.
type ErrInternal struct {
	Message string `json:"message"`
	Err     string `json:"error"`
}

// Error implements error interface.
func (e ErrInternal) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "unknown error"
	}
	return string(b)
}

// ErrRateLimit is an error returned when the API rate limit is exceeded
type ErrRateLimit struct {
	Message string
}

// Error implements error interface.
func (e ErrRateLimit) Error() string {
	return e.Message
}
