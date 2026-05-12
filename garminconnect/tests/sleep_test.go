package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSleepData(t *testing.T) {
	c, stop := newVCRClient(t, "sleep_data")
	defer stop()

	s, err := c.SleepData(testDate)
	require.NoError(t, err)

	assert.Equal(t, "2026-05-10", s.DailySleepDTO.CalendarDate)
	assert.Equal(t, 16759, s.DailySleepDTO.SleepTimeSeconds)
	assert.Equal(t, 3000, s.DailySleepDTO.DeepSleepSeconds)
	assert.Equal(t, 4380, s.DailySleepDTO.REMSleepSeconds)
	assert.Equal(t, 96.0, s.DailySleepDTO.SpO2AvgReadingPercent)
	assert.Equal(t, "POSITIVE_SHORT_BUT_REFRESHING", s.DailySleepDTO.SleepScoreFeedback)
	assert.Equal(t, 14, s.RestlessMomentsCount)
}

func TestHRVData(t *testing.T) {
	c, stop := newVCRClient(t, "hrv_data")
	defer stop()

	h, err := c.HRVData(testDate)
	require.NoError(t, err)

	assert.Equal(t, 71, h.HRVSummary.WeeklyAvg)
	assert.Equal(t, 75, h.HRVSummary.LastNight)
	assert.Equal(t, "BALANCED", h.HRVSummary.Status)
	assert.NotEmpty(t, h.HRVReadings)
	assert.Equal(t, 81, h.HRVReadings[0].HRVValue)
}
