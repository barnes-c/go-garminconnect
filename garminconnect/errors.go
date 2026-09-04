package garminconnect

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrRateLimit       = errors.New("rate limit exceeded")
	ErrMFARequired     = errors.New("MFA required but no prompt configured")
	ErrCaptchaRequired = errors.New("CAPTCHA required")
	ErrNoData          = errors.New("no data")
)

// maxErrorBody caps how much of a failed response body is retained.
const maxErrorBody = 512

// APIError is returned for unexpected HTTP status codes.
//
// Test for a specific condition with [errors.Is] against the sentinels above
// rather than comparing StatusCode: 401 and 429 unwrap to ErrUnauthorized and
// ErrRateLimit respectively.
type APIError struct {
	StatusCode int
	Path       string
	// Body is a bounded snippet of the error response, empty when there was
	// none. It is not guaranteed to be JSON: Garmin returns structured
	// errors, its edge returns HTML challenges, and some tiers return
	// nothing. Kept out of Error() to keep messages short — read it
	// directly when diagnosing.
	Body string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("garmin API %d: %s", e.StatusCode, e.Path)
}

// Unwrap maps the status codes callers act on to sentinel errors, so a single
// returned type stays matchable with errors.Is.
func (e *APIError) Unwrap() error {
	switch e.StatusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusTooManyRequests:
		return ErrRateLimit
	}
	return nil
}

// newAPIError builds an APIError from a failed response, consuming up to
// maxErrorBody of the body. Only call it on a response already known to have
// failed.
func newAPIError(resp *http.Response, path string) *APIError {
	data, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
	return &APIError{
		StatusCode: resp.StatusCode,
		Path:       path,
		Body:       strings.TrimSpace(string(data)),
	}
}
