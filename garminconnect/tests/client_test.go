package garminconnect_test

import (
	"io"
	"net/http"
	"strings"
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

func TestAutoRefreshOnUnauthorized(t *testing.T) {
	calls := 0
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		switch calls {
		case 1:
			// First API call: Garmin has revoked the access token.
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		case 2:
			// Client exchanges the refresh token for a new access token.
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"new-token","refresh_token":"new-refresh","expires_in":3600}`)),
				Header:     make(http.Header),
			}, nil
		case 3:
			// Retried API call — must carry the refreshed token.
			assert.Equal(t, "Bearer new-token", req.Header.Get("Authorization"))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"displayName":"Test User"}`)),
				Header:     make(http.Header),
			}, nil
		default:
			t.Fatalf("unexpected HTTP call %d to %s", calls, req.URL)
			return nil, nil
		}
	})

	c := gc.NewClient("",
		gc.WithHTTPClient(&http.Client{Transport: transport}),
		gc.WithToken("old-token"),
		gc.WithDisplayName("testuser"),
		gc.WithRefreshToken("some-refresh-token"),
	)

	profile, err := c.UserProfile()
	require.NoError(t, err)
	assert.Equal(t, "Test User", profile.DisplayName)
	assert.Equal(t, 3, calls)
}

func TestRefreshFailureReturnsUnauthorized(t *testing.T) {
	calls := 0
	transport := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		calls++
		code := http.StatusUnauthorized
		if calls == 2 {
			code = http.StatusBadRequest // refresh token is also expired
		}
		return &http.Response{
			StatusCode: code,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
		}, nil
	})

	c := gc.NewClient("",
		gc.WithHTTPClient(&http.Client{Transport: transport}),
		gc.WithToken("old-token"),
		gc.WithDisplayName("testuser"),
		gc.WithRefreshToken("expired-refresh-token"),
	)

	_, err := c.UserProfile()
	require.ErrorIs(t, err, gc.ErrUnauthorized)
	assert.Equal(t, 2, calls)
}
