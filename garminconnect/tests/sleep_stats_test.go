package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDailySleepData(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// Replay queries testDate. To record a night with actual sleep, set
	// GARMIN_SLEEP_DATE to a recent date; the sanitizer rewrites that real date
	// back to testDate in the cassette URL so it still replays.
	out, err := c.DailySleepData(t.Context(), dataDate(t, "GARMIN_SLEEP_DATE"))
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.NotEmpty(t, out.DailySleepDTO.CalendarDate)
}

func TestSleepStats(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, 0, -27)
	out, err := c.SleepStats(t.Context(), start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestHRVDataRange(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// One day: the sanitizer collapses every recorded date after testDate to
	// testDate, so a multi-day range would not replay.
	d := dataDate(t, "GARMIN_SLEEP_DATE")
	out, err := c.HRVDataRange(t.Context(), d, d)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}
