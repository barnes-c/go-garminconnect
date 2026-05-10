package garminconnect_test

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
)

func readMultipart(t *testing.T, r *http.Request) (filename string, data []byte) {
	t.Helper()
	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(mediaType, "multipart/"))
	mr := multipart.NewReader(r.Body, params["boundary"])
	part, err := mr.NextPart()
	require.NoError(t, err)
	body, err := io.ReadAll(part)
	require.NoError(t, err)
	return part.FileName(), body
}

func TestUploadActivity_FIT(t *testing.T) {
	fitData := []byte{0x0e, 0x10, 0x00, 0x00} // minimal FIT header bytes
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/upload-service/upload.fit", r.URL.Path)
		fname, body := readMultipart(t, r)
		assert.Equal(t, "activity.fit", fname)
		assert.Equal(t, fitData, body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{"detailedImportResult": map[string]any{}})
	}))
	c := newServerClient(t, srv)

	out, err := c.UploadActivity(fitData, "activity.fit")
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestUploadActivity_GPX(t *testing.T) {
	gpxData := []byte(`<?xml version="1.0"?>`)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/upload-service/upload.gpx", r.URL.Path)
		fname, _ := readMultipart(t, r)
		assert.Equal(t, "track.gpx", fname)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{})
	}))
	c := newServerClient(t, srv)

	_, err := c.UploadActivity(gpxData, "track.gpx")
	require.NoError(t, err)
}

func TestUploadActivity_NoExtension(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("server should not be called for invalid input")
	}))
	c := newServerClient(t, srv)

	_, err := c.UploadActivity([]byte("data"), "noextension")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "extension")
}

func TestUploadWorkout(t *testing.T) {
	fitData := []byte{0x0e, 0x10}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/upload-service/upload", r.URL.Path)
		fname, body := readMultipart(t, r)
		assert.Equal(t, "workout.fit", fname)
		assert.Equal(t, fitData, body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{})
	}))
	c := newServerClient(t, srv)

	out, err := c.UploadWorkout(fitData, "workout.fit")
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestUploadActivity_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	c := newServerClient(t, srv)

	_, err := c.UploadActivity([]byte("data"), "file.fit")
	assert.ErrorIs(t, err, gc.ErrUnauthorized)
}
