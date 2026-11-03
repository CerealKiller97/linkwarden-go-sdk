package linkwarden

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// APIError represents a non-2xx response from the Linkwarden API.
type APIError struct {
	StatusCode int
	Message    string
	RawBody    []byte
}

func newAPIError(statusCode int, body []byte) *APIError {
	msg := extractErrorMessage(body)
	if msg == "" {
		msg = http.StatusText(statusCode)
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    msg,
		RawBody:    body,
	}
}

func extractErrorMessage(body []byte) string {
	var withMessage struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &withMessage); err == nil && withMessage.Message != "" {
		return withMessage.Message
	}

	var withResponse struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &withResponse); err == nil && withResponse.Response != "" {
		return withResponse.Response
	}

	return ""
}

func (e *APIError) Error() string {
	return fmt.Sprintf("linkwarden: %d %s", e.StatusCode, e.Message)
}

// IsNotFound reports whether err is an APIError with a 404 status code.
func IsNotFound(err error) bool {
	return statusCodeIs(err, http.StatusNotFound)
}

// IsConflict reports whether err is an APIError with a 409 status code.
func IsConflict(err error) bool {
	return statusCodeIs(err, http.StatusConflict)
}

func statusCodeIs(err error, code int) bool {
	apiErr, ok := errors.AsType[*APIError](err)
	if !ok {
		return false
	}

	return apiErr.StatusCode == code
}
