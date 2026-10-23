package garminconnect_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestActivityExtras uploads a throwaway GPS activity (dated just before
// testDate), reads its laps, track and stats, then deletes it.
func TestActivityExtras(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	_, err := c.UploadActivity(ctx, syntheticGPX("go-garminconnect extras", true), "extras.gpx")
	require.NoError(t, err)
	var id int64
	for range 10 {
		acts, err := c.Activities(ctx, 0, 1)
		require.NoError(t, err)
		if len(acts) > 0 {
			id = acts[0].ActivityID
			break
		}
		if os.Getenv("GARMIN_TOKEN") != "" {
			time.Sleep(2 * time.Second)
		}
	}
	require.NotZero(t, id)

	laps, err := c.ActivityLaps(ctx, id)
	require.NoError(t, err)
	require.NotEmpty(t, laps)
	assert.NotZero(t, laps[0].Distance)

	track, err := c.ActivityPolyline(ctx, id)
	require.NoError(t, err)
	assert.NotZero(t, track.NumberOfPoints)
	assert.NotEmpty(t, track.EncodedSamples)

	types, err := c.ActivityEventTypes(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, types)

	metrics, err := c.AvailableActivityMetrics(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, metrics)

	start := testDate.AddDate(0, 0, -7)
	progress, err := c.ActivityProgress(ctx, start, testDate, "distance", "duration")
	require.NoError(t, err)
	require.NotEmpty(t, progress)
	assert.NotZero(t, progress[0].Stats["all"]["distance"].Count)

	values, err := c.ActivityMetrics(ctx, start, testDate, "distance", "duration")
	require.NoError(t, err)
	require.NotEmpty(t, values)
	assert.NotZero(t, values[0].ActivityID)
	assert.Contains(t, values[0].Values, "distance")

	require.NoError(t, c.DeleteActivity(ctx, id))
}
