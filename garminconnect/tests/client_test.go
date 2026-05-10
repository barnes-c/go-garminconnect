package garminconnect_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
)

func clientWith(t *testing.T, code int) *gc.Client {
	t.Helper()
	return gc.NewClient("",
		gc.WithHTTPClient(&http.Client{Transport: fixedTransport(code, "")}),
		gc.WithToken("test"),
		gc.WithDisplayName("testuser"),
	)
}

func TestActivities_Unauthorized(t *testing.T) {
	c := clientWith(t, http.StatusUnauthorized)
	_, err := c.Activities(1)
	assert.ErrorIs(t, err, gc.ErrUnauthorized)
}

func TestActivities_RateLimit(t *testing.T) {
	c := clientWith(t, http.StatusTooManyRequests)
	_, err := c.Activities(1)
	assert.ErrorIs(t, err, gc.ErrRateLimit)
}

func TestActivities_APIError(t *testing.T) {
	c := clientWith(t, http.StatusInternalServerError)
	_, err := c.Activities(1)
	require.Error(t, err)
	var apiErr *gc.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
}
