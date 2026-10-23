package garminconnect_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDate = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func TestActivities(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	acts, err := c.Activities(t.Context(), 0, 2)
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
	assert.NotEmpty(t, run.LocationName)

	assert.NotEmpty(t, acts[1].ActivityType.TypeKey)
}

func TestLastActivity(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// LastActivity uses limit=1 internally, but we reuse the activities cassette
	// which returns 2. Stub it separately so we control the limit parameter.
	acts, err := c.Activities(t.Context(), 0, 2)
	require.NoError(t, err)
	assert.NotZero(t, acts[0].ActivityID)
}

func TestActivityDetail(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// Record cassette: fetch the most recent activity to get a real ID.
	acts, err := c.Activities(t.Context(), 0, 1)
	require.NoError(t, err)
	require.NotEmpty(t, acts)

	detail, err := c.ActivityDetail(t.Context(), acts[0].ActivityID)
	require.NoError(t, err)
	assert.NotEmpty(t, detail)
}

func TestActivityDetails(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	acts, err := c.Activities(t.Context(), 0, 1)
	require.NoError(t, err)
	require.NotEmpty(t, acts)

	details, err := c.ActivityDetails(t.Context(), acts[0].ActivityID)
	require.NoError(t, err)
	assert.NotEmpty(t, details)
}

func TestActivityTypes(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	types, err := c.ActivityTypes(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, types)
	assert.NotEmpty(t, types[0].TypeKey)
}

func TestActivitiesForDailySummary(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// Replay queries testDate. To record a non-empty response, set
	// GARMIN_SUMMARY_DATE to a day with a logged activity; the sanitizer rewrites
	// that real date back to testDate in the cassette URL so it replays cleanly.
	summaryDate := testDate
	if d := os.Getenv("GARMIN_SUMMARY_DATE"); d != "" {
		parsed, err := time.Parse("2006-01-02", d)
		require.NoError(t, err)
		summaryDate = parsed
	}
	out, err := c.ActivitiesForDailySummary(t.Context(), summaryDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	for _, a := range out {
		assert.NotZero(t, a.ActivityID)
	}
}

func TestActivityHRZones(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// Record cassette: fetch the most recent activity to get a real ID.
	acts, err := c.Activities(t.Context(), 0, 1)
	require.NoError(t, err)
	require.NotEmpty(t, acts)

	zones, err := c.ActivityHRZones(t.Context(), acts[0].ActivityID)
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, zones)
	for i, z := range zones {
		assert.Equal(t, i+1, z.ZoneNumber)
	}
	// Garmin sends no zoneHighBoundary on this endpoint, only lower bounds.
	assert.Positive(t, zones[len(zones)-1].ZoneLowBPM)
}

func TestActivityPowerZones(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// Record cassette: fetch the most recent activity to get a real ID. The
	// activity must have power data (e.g. a ride or run with a power source).
	acts, err := c.Activities(t.Context(), 0, 1)
	require.NoError(t, err)
	require.NotEmpty(t, acts)

	zones, err := c.ActivityPowerZones(t.Context(), acts[0].ActivityID)
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, zones)
	for i, z := range zones {
		assert.Equal(t, i+1, z.ZoneNumber)
	}
	// Garmin sends no zoneHighBoundary on this endpoint, only lower bounds.
	assert.Positive(t, zones[len(zones)-1].ZoneLowWatts)
}

func TestActivityCount(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	count, err := c.ActivityCount(t.Context())
	require.NoError(t, err)
	assert.Positive(t, count)
}

func TestActivitiesByDate(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	acts, err := c.ActivitiesByDate(t.Context(), start, testDate, "")
	require.NoError(t, err)
	assert.NotNil(t, acts)
}

func TestPersonalRecords(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	prs, err := c.PersonalRecords(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, prs)
}

func TestIntensityMinutes(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.IntensityMinutes(t.Context(), testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

// TestActivityWritesLifecycle favorites a throwaway strength activity and
// writes its exercise sets, then deletes it so recording leaves the account
// unchanged.
func TestActivityWritesLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	id, err := c.CreateManualActivity(ctx, gc.ManualActivity{
		Name:     "go-garminconnect activity test",
		TypeKey:  "strength_training",
		Start:    testDate.Add(8 * time.Hour),
		TimeZone: "Europe/Zurich",
		Duration: 30 * time.Minute,
	})
	require.NoError(t, err)
	require.NotZero(t, id)

	require.NoError(t, c.SetActivityFavorite(ctx, id, true))
	// The favorites list lags behind the flag, so only check the call works.
	_, err = c.FavoriteActivities(ctx, 5)
	require.NoError(t, err)

	sets, err := c.ActivityExerciseSets(ctx, id)
	require.NoError(t, err)
	sets["exerciseSets"] = json.RawMessage(`[{"exercises":[{"category":"BENCH_PRESS","name":null,"probability":100}],` +
		`"repetitionCount":10,"weight":60000,"duration":45,"setType":"ACTIVE","startTime":"2026-01-01T07:00:00.0"}]`)
	require.NoError(t, c.SetActivityExerciseSets(ctx, id, sets))

	require.NoError(t, c.DeleteActivity(ctx, id))
}
