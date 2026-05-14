package garminconnect

import (
	"errors"
	"fmt"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrRateLimit    = errors.New("rate limit exceeded")
	ErrNoData       = errors.New("no data")
)

// APIError is returned for unexpected HTTP status codes.
type APIError struct {
	StatusCode int
	Path       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("garmin API %d: %s", e.StatusCode, e.Path)
}
