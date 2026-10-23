package garminconnect_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetActivityName(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/activity-service/activity/7", r.URL.Path)
		var body map[string]any
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "Morning Run", body["activityName"])
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.SetActivityName(t.Context(), 7, "Morning Run")
	require.NoError(t, err)
}

func TestDeleteActivity(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/activity-service/activity/42", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.DeleteActivity(t.Context(), 42)
	require.NoError(t, err)
}

func TestDeleteWorkout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/workout-service/workout/77", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.DeleteWorkout(t.Context(), 77)
	require.NoError(t, err)
}

func TestSetActivityType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/activity-service/activity/10", r.URL.Path)
		var body map[string]any
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		at := body["activityType"].(map[string]any)
		assert.Equal(t, "running", at["typeKey"])
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.SetActivityType(t.Context(), 10, 1, 0, "running")
	require.NoError(t, err)
}

func TestSetActivityName_BearerToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer test", auth)
		body, _ := io.ReadAll(r.Body)
		_ = body
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.SetActivityName(t.Context(), 1, "test")
	require.NoError(t, err)
}
