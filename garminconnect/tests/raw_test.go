package garminconnect_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
)

func TestRawGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/some-service/unwrapped", r.URL.Path)
		assert.Equal(t, "2026-01-01", r.URL.Query().Get("date"))
		assert.Equal(t, "Bearer test", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"value": 42})
	}))
	c := newServerClient(t, srv)

	var out map[string]any
	err := c.Get(t.Context(), "/some-service/unwrapped", url.Values{"date": {"2026-01-01"}}, &out)
	require.NoError(t, err)
	assert.InDelta(t, 42, out["value"], 0.001)
}

func TestRawGetAbsoluteURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/elsewhere", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	// Point the client at a different base to prove the absolute URL wins.
	t.Cleanup(srv.Close)
	c := gc.NewClient("",
		gc.WithBaseURL("https://invalid.example"),
		gc.WithToken("test"),
		gc.WithDisplayName("testuser"),
	)

	var out map[string]any
	err := c.Get(t.Context(), srv.URL+"/elsewhere", nil, &out)
	require.NoError(t, err)
	assert.Equal(t, true, out["ok"])
}

func TestRawGetBytes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/download-service/file", r.URL.Path)
		w.Write([]byte("\x0eFIT-ish bytes"))
	}))
	c := newServerClient(t, srv)

	body, err := c.GetBytes(t.Context(), "/download-service/file", nil)
	require.NoError(t, err)
	assert.Equal(t, "\x0eFIT-ish bytes", string(body))
}

func TestRawPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/some-service/thing", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		var body map[string]any
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "hello", body["name"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": "abc"})
	}))
	c := newServerClient(t, srv)

	var out map[string]any
	err := c.Post(t.Context(), "/some-service/thing", map[string]any{"name": "hello"}, &out)
	require.NoError(t, err)
	assert.Equal(t, "abc", out["id"])
}

func TestRawPut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/some-service/thing/1", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.Put(t.Context(), "/some-service/thing/1", map[string]any{"x": 1}, nil)
	require.NoError(t, err)
}

func TestRawDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/some-service/thing/1", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.Delete(t.Context(), "/some-service/thing/1")
	require.NoError(t, err)
}

func TestRawUpload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/upload-service/upload/.fit", r.URL.Path)
		f, hdr, err := r.FormFile("file")
		if assert.NoError(t, err) {
			defer f.Close()
			assert.Equal(t, "run.fit", hdr.Filename)
			data, err := io.ReadAll(f)
			assert.NoError(t, err)
			assert.Equal(t, "payload", string(data))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"detailedImportResult": map[string]any{}})
	}))
	c := newServerClient(t, srv)

	var out map[string]any
	err := c.Upload(t.Context(), "/upload-service/upload/.fit", []byte("payload"), "run.fit", &out)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestRawGetPropagatesAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	c := newServerClient(t, srv)

	err := c.Get(t.Context(), "/some-service/missing", nil, nil)
	var apiErr *gc.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
}
