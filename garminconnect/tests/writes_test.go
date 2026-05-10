package garminconnect_test

import (
	"encoding/json"
	"fmt"
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
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "Morning Run", body["activityName"])
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.SetActivityName(7, "Morning Run")
	require.NoError(t, err)
}

func TestDeleteActivity(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/activity-service/activity/42", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.DeleteActivity(42)
	require.NoError(t, err)
}

func TestAddWeighIn(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/weight-service/user-weight", r.URL.Path)
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.InDelta(t, 75.5, body["value"], 0.001)
		assert.Equal(t, "kg", body["unitKey"])
		assert.Equal(t, "2026-05-10T08:00:00Z", body["dateTimestamp"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"detailedWeightInResult": map[string]any{}})
	}))
	c := newServerClient(t, srv)

	out, err := c.AddWeighIn(75.5, "kg", "2026-05-10T08:00:00Z")
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestDeleteWeighIn(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/weight-service/weight/2026-05-10/byversion/12345", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.DeleteWeighIn("2026-05-10", 12345)
	require.NoError(t, err)
}

func TestScheduleWorkout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, fmt.Sprintf("/calendar-service/schedule/workout/%d", 55), r.URL.Path)
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "2026-05-15", body["date"])
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"scheduledWorkoutId": 999})
	}))
	c := newServerClient(t, srv)

	out, err := c.ScheduleWorkout(55, "2026-05-15")
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestDeleteWorkout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/workout-service/workout/77", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.DeleteWorkout(77)
	require.NoError(t, err)
}

func TestSetActivityType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/activity-service/activity/10", r.URL.Path)
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		at := body["activityType"].(map[string]any)
		assert.Equal(t, "running", at["typeKey"])
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.SetActivityType(10, 1, 0, "running")
	require.NoError(t, err)
}

func TestUnscheduleWorkout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/calendar-service/schedule/workout/200", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	c := newServerClient(t, srv)

	err := c.UnscheduleWorkout(200)
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

	err := c.SetActivityName(1, "test")
	require.NoError(t, err)
}
