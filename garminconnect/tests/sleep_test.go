package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSleepData(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	s, err := c.SleepData(t.Context(), dataDate(t, "GARMIN_SLEEP_DATE"))
	require.NoError(t, err)

	assert.NotEmpty(t, s.DailySleepDTO.CalendarDate)
	if s.DailySleepDTO.SleepTimeSeconds == 0 {
		t.Skip("no sleep data for test date")
	}
	assert.NotZero(t, s.DailySleepDTO.SleepTimeSeconds)
	assert.NotZero(t, s.DailySleepDTO.DeepSleepSeconds)
	assert.NotZero(t, s.DailySleepDTO.REMSleepSeconds)
	assert.NotZero(t, s.DailySleepDTO.SpO2AvgReadingPercent)
	assert.NotZero(t, s.DailySleepDTO.LowestSpO2Value)
	assert.NotZero(t, s.DailySleepDTO.HighestSpO2Value)
	assert.NotZero(t, s.DailySleepDTO.SleepScores["overall"].Value)
	require.NotNil(t, s.DailySleepDTO.SleepNeed)
	assert.NotZero(t, s.DailySleepDTO.SleepNeed.Actual)
	assert.NotZero(t, s.RestingHeartRate)
	assert.NotEmpty(t, s.SleepHeartRate)
	require.NotNil(t, s.SpO2Summary)
	assert.NotZero(t, s.SpO2Summary.LowestSpO2)
	assert.NotEmpty(t, s.DailySleepDTO.SleepScoreFeedback)
}

func TestHRVData(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	h, err := c.HRVData(t.Context(), dataDate(t, "GARMIN_SLEEP_DATE"))
	require.NoError(t, err)

	assert.NotZero(t, h.HRVSummary.WeeklyAvg)
	assert.NotEmpty(t, h.HRVSummary.Status)
	if len(h.HRVReadings) == 0 {
		t.Skip("no HRV readings for test date")
	}
	assert.NotZero(t, h.HRVReadings[0].HRVValue)
	assert.NotEmpty(t, h.HRVReadings[0].ReadingTimeGMT)
}

func TestSleepStatsWeekly(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.SleepStatsWeekly(t.Context(), dataDate(t, "GARMIN_SLEEP_DATE"), 2)
	require.NoError(t, err)
	assert.NotZero(t, out.Overall.AverageSleepSeconds)
	require.NotEmpty(t, out.Weeks)
	assert.NotZero(t, out.Weeks[0].Values.DeepTime)
}

func TestSleepStatsOverall(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	d := dataDate(t, "GARMIN_SLEEP_DATE")

	out, err := c.SleepStatsOverall(t.Context(), d, d)
	require.NoError(t, err)
	assert.NotZero(t, out.Overall.AverageSleepSeconds)
	assert.NotZero(t, out.Overall.AverageSpO2)
}

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
	assert.NotZero(t, out.Overall.AverageSleepSeconds)
	require.NotEmpty(t, out.Days)
	assert.NotZero(t, out.Days[0].Values.TotalSleepTimeInSeconds)
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
	require.NotEmpty(t, out)
	assert.NotEmpty(t, out[0].CalendarDate)
}
