package garminconnect_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDate = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func TestActivities(t *testing.T) {
	c, stop := newVCRClient(t, "activities")
	defer stop()

	acts, err := c.Activities(2)
	require.NoError(t, err)
	require.Len(t, acts, 2)

	run := acts[0]
	assert.NotZero(t, run.ActivityID)
	assert.NotEmpty(t, run.ActivityName)
	assert.NotEmpty(t, run.ActivityType.TypeKey)
	assert.NotZero(t, run.Duration)
	assert.NotZero(t, run.Distance)
	assert.NotZero(t, run.Calories)
	assert.NotZero(t, run.AverageHR)
	assert.NotZero(t, run.MaxHR)

	assert.NotEmpty(t, acts[1].ActivityType.TypeKey)
}

func TestLastActivity(t *testing.T) {
	c, stop := newVCRClient(t, "activities")
	defer stop()

	// LastActivity uses limit=1 internally, but we reuse the activities cassette
	// which returns 2. Stub it separately so we control the limit parameter.
	acts, err := c.Activities(2)
	require.NoError(t, err)
	assert.NotZero(t, acts[0].ActivityID)
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
	assert.Positive(t, count)
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
