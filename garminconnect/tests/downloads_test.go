package garminconnect_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadActivity_URLFormats(t *testing.T) {
	cases := []struct {
		format  gc.DownloadFormat
		wantURL string
	}{
		{gc.FormatOriginal, "/download-service/files/activity/42"},
		{gc.FormatTCX, "/download-service/export/tcx/activity/42"},
		{gc.FormatGPX, "/download-service/export/gpx/activity/42"},
		{gc.FormatKML, "/download-service/export/kml/activity/42"},
		{gc.FormatCSV, "/download-service/export/csv/activity/42"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(string(tc.format), func(t *testing.T) {
			var gotPath string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				w.Write([]byte("fitdata"))
			}))
			c := newServerClient(t, srv)

			_, err := c.DownloadActivity(42, tc.format)
			require.NoError(t, err)
			assert.Equal(t, tc.wantURL, gotPath)
		})
	}
}

func TestDownloadActivity_ReturnsBody(t *testing.T) {
	body := []byte("FIT binary content")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	}))
	c := newServerClient(t, srv)

	got, err := c.DownloadActivity(1, gc.FormatOriginal)
	require.NoError(t, err)
	assert.Equal(t, body, got)
}

func TestDownloadActivity_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	c := newServerClient(t, srv)

	_, err := c.DownloadActivity(1, gc.FormatOriginal)
	assert.ErrorIs(t, err, gc.ErrUnauthorized)
}

func TestDownloadWorkout(t *testing.T) {
	body := []byte("workout FIT data")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, fmt.Sprintf("/download-service/files/workout/%d", 99), r.URL.Path)
		w.Write(body)
	}))
	c := newServerClient(t, srv)

	got, err := c.DownloadWorkout(99)
	require.NoError(t, err)
	assert.Equal(t, body, got)
}

func TestDownloadActivity_BearerToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "Bearer "))
		w.Write([]byte("ok"))
	}))
	c := newServerClient(t, srv)

	_, err := c.DownloadActivity(1, gc.FormatGPX)
	require.NoError(t, err)
}
