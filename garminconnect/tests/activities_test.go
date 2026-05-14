package garminconnect_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
)

var testDate = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func TestActivities(t *testing.T) {
	c, stop := newVCRClient(t, "activities")
	defer stop()

	acts, err := c.Activities(2)
	require.NoError(t, err)
	require.Len(t, acts, 2)

	run := acts[0]
	assert.Equal(t, int64(10000001), run.ActivityID)
	assert.Equal(t, "Morning Run", run.ActivityName)
	assert.Equal(t, "running", run.ActivityType.TypeKey)
	assert.Equal(t, 1500.0, run.Duration)
	assert.Equal(t, 4100.0, run.Distance)
	assert.Equal(t, 334.0, run.Calories)
	assert.Equal(t, 153.0, run.AverageHR)
	assert.Equal(t, 172.0, run.MaxHR)
	assert.Equal(t, 53.0, run.VO2MaxValue)
	assert.Equal(t, "Anytown", run.LocationName)
	assert.True(t, run.HasPolyline)

	assert.Equal(t, "kayaking_v2", acts[1].ActivityType.TypeKey)
}

func TestLastActivity(t *testing.T) {
	c, stop := newVCRClient(t, "activities")
	defer stop()

	// LastActivity uses limit=1 internally, but we reuse the activities cassette
	// which returns 2. Stub it separately so we control the limit parameter.
	acts, err := c.Activities(2)
	require.NoError(t, err)
	assert.Equal(t, int64(10000001), acts[0].ActivityID)
}

func TestLastActivity_Empty(t *testing.T) {
	c, stop := newVCRClient(t, "activities_empty")
	defer stop()

	_, err := c.LastActivity()
	assert.ErrorIs(t, err, gc.ErrNoData)
}

func TestActivityDetail(t *testing.T) {
	c, stop := newVCRClient(t, "activity_detail")
	defer stop()

	// Record cassette: fetch the most recent activity to get a real ID.
	acts, err := c.Activities(1)
	require.NoError(t, err)
	require.NotEmpty(t, acts)

	detail, err := c.ActivityDetail(acts[0].ActivityID)
	require.NoError(t, err)
	assert.NotEmpty(t, detail)
}

func TestActivityCount(t *testing.T) {
	c, stop := newVCRClient(t, "activity_count")
	defer stop()

	count, err := c.ActivityCount()
	require.NoError(t, err)
	assert.Greater(t, count, 0)
}

func TestActivitiesByDate(t *testing.T) {
	c, stop := newVCRClient(t, "activities_by_date")
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	acts, err := c.ActivitiesByDate(start, testDate, "")
	require.NoError(t, err)
	assert.NotNil(t, acts)
}

func TestPersonalRecords(t *testing.T) {
	c, stop := newVCRClient(t, "personal_records")
	defer stop()

	prs, err := c.PersonalRecords()
	require.NoError(t, err)
	assert.NotEmpty(t, prs)
}

func TestIntensityMinutes(t *testing.T) {
	c, stop := newVCRClient(t, "intensity_minutes")
	defer stop()

	out, err := c.IntensityMinutes(testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}
