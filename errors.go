package linkwarden

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Machine-readable error codes Linkwarden sends in APIError.Code.
const (
	// ErrorCodeEmailNotVerified: the credentials are right but the account's email
	// isn't verified. Linkwarden emails a code to send back in
	// CreateSessionRequest.Code.
	ErrorCodeEmailNotVerified = "EMAIL_NOT_VERIFIED"
	// ErrorCodeInvalidVerificationCode: the email verification code is wrong or has
	// expired.
	ErrorCodeInvalidVerificationCode = "INVALID_VERIFICATION_CODE"
	// ErrorCodeTooManyAttempts: too many wrong passwords for the account. Wait
	// APIError.RetryAfter before trying again.
	ErrorCodeTooManyAttempts = "TOO_MANY_ATTEMPTS"
)

// APIError represents a non-2xx response from the Linkwarden API.
type APIError struct {
	StatusCode int
	Message    string
	// Code is the machine-readable error code, when the API sends one (see the
	// ErrorCode constants). It's empty for most errors.
	Code string
	// RetryAfter is how long to wait before retrying, from the Retry-After header of a
	// 429 response. It's zero when the response has none.
	RetryAfter time.Duration
	RawBody    []byte
}

func newAPIError(statusCode int, body []byte) *APIError {
	msg := extractErrorMessage(body)
	if msg == "" {
		msg = http.StatusText(statusCode)
	}

	var withCode struct {
		Code string `json:"code"`
	}
	_ = json.Unmarshal(body, &withCode)

	return &APIError{
		StatusCode: statusCode,
		Message:    msg,
		Code:       withCode.Code,
		RawBody:    body,
	}
}

// parseRetryAfter reads a Retry-After header given in seconds, which is the form
// Linkwarden sends. Anything else is treated as absent.
func parseRetryAfter(v string) time.Duration {
	secs, err := strconv.Atoi(v)
	if err != nil || secs < 0 {
		return 0
	}

	return time.Duration(secs) * time.Second
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

// IsTooManyRequests reports whether err is an APIError with a 429 status code, such
// as SessionService.Create after too many wrong passwords. APIError.RetryAfter says
// how long to wait.
func IsTooManyRequests(err error) bool {
	return statusCodeIs(err, http.StatusTooManyRequests)
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
