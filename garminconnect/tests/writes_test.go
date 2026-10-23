package garminconnect_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
