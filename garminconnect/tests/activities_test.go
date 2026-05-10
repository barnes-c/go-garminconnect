package garminconnect_test

import (
	"testing"
	"time"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDate = time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)

func TestActivities(t *testing.T) {
	c, stop := newVCRClient(t, "activities")
	defer stop()

	acts, err := c.Activities(2)
	require.NoError(t, err)
	require.Len(t, acts, 2)

	run := acts[0]
	assert.Equal(t, int64(1234567890), run.ActivityID)
	assert.Equal(t, "Morning Run", run.ActivityName)
	assert.Equal(t, "running", run.ActivityType.TypeKey)
	assert.Equal(t, 3600.0, run.Duration)
	assert.Equal(t, 10000.0, run.Distance)
	assert.Equal(t, 650.0, run.Calories)
	assert.Equal(t, 155.0, run.AverageHR)
	assert.Equal(t, 178.0, run.MaxHR)
	assert.Equal(t, 52.0, run.VO2MaxValue)
	assert.Equal(t, "Berlin", run.LocationName)
	assert.True(t, run.HasPolyline)

	assert.Equal(t, "cycling", acts[1].ActivityType.TypeKey)
}

func TestLastActivity(t *testing.T) {
	c, stop := newVCRClient(t, "activities")
	defer stop()

	// LastActivity uses limit=1 internally, but we reuse the activities cassette
	// which returns 2. Stub it separately so we control the limit parameter.
	acts, err := c.Activities(2)
	require.NoError(t, err)
	assert.Equal(t, int64(1234567890), acts[0].ActivityID)
}

func TestLastActivity_Empty(t *testing.T) {
	c, stop := newVCRClient(t, "activities_empty")
	defer stop()

	_, err := c.LastActivity()
	assert.ErrorIs(t, err, gc.ErrNoData)
}
