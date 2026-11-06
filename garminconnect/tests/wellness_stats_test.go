package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWellnessAverages(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	d := dataDate(t, "GARMIN_SUMMARY_DATE")

	out, err := c.WellnessAverages(t.Context(), d, d)
	require.NoError(t, err)
	assert.NotZero(t, out.AverageSteps)
}

func TestDailyMovement(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.DailyMovement(t.Context(), dataDate(t, "GARMIN_SUMMARY_DATE"))
	require.NoError(t, err)
	assert.NotEmpty(t, out.Descriptors)
	assert.NotEmpty(t, out.Values)
}

func TestIntensityMinutesStats(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.IntensityMinutesStats(t.Context(), dataDate(t, "GARMIN_SUMMARY_DATE"), 1)
	require.NoError(t, err)
	assert.NotZero(t, out.WeeklyGoal)
	assert.NotEmpty(t, out.Days)
}

func TestMonthlySteps(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.MonthlySteps(t.Context(), dataDate(t, "GARMIN_SUMMARY_DATE"), 1)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].Values.TotalSteps)
}

func TestDailySteps(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, 0, -7)
	entries, err := c.DailySteps(t.Context(), start, testDate)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestWeeklyStress(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.WeeklyStress(t.Context(), testDate, 4)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].AvgStressLevel)
	assert.NotNil(t, out)
}

func TestWeeklyIntensityMinutes(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, 0, -7)
	out, err := c.WeeklyIntensityMinutes(t.Context(), start, testDate)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].ModerateIntensityMinutes)
	assert.NotZero(t, out[0].VigorousIntensityMinutes)
	assert.NotNil(t, out)
}

func TestCaloriesDaily(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// One day: the sanitizer collapses every recorded date after testDate to
	// testDate, so a multi-day range would not replay.
	d := dataDate(t, "GARMIN_SUMMARY_DATE")
	out, err := c.CaloriesDaily(t.Context(), d, d)
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	for _, d := range out {
		assert.NotEmpty(t, d.CalendarDate)
		assert.InDelta(t, d.Active+d.Resting, d.Total, 0.001)
	}
}

// TestWeeklySteps only checks that the call succeeds: the recording account
// has no step data, so the list is empty.
func TestWeeklySteps(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	_, err := c.WeeklySteps(t.Context(), testDate, 2)
	skipAPIError(t, err)
	require.NoError(t, err)
}
