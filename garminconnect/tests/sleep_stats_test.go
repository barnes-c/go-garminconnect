package garminconnect_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDailySleepData(t *testing.T) {
	c, stop := newVCRClient(t, "daily_sleep_data")
	defer stop()

	// Replay queries testDate. To record a night with actual sleep, set
	// GARMIN_SLEEP_DATE to a recent date; the sanitizer rewrites that real date
	// back to testDate in the cassette URL so it still replays.
	d := testDate
	if v := os.Getenv("GARMIN_SLEEP_DATE"); v != "" {
		parsed, err := time.Parse("2006-01-02", v)
		require.NoError(t, err)
		d = parsed
	}
	out, err := c.DailySleepData(t.Context(), d)
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.NotEmpty(t, out.DailySleepDTO.CalendarDate)
}

func TestSleepStats(t *testing.T) {
	c, stop := newVCRClient(t, "sleep_stats")
	defer stop()

	start := testDate.AddDate(0, 0, -27)
	out, err := c.SleepStats(t.Context(), start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}
