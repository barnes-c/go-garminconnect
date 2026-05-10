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
	assert.Equal(t, 25200, s.DailySleepDTO.SleepTimeSeconds)
	assert.Equal(t, 5400, s.DailySleepDTO.DeepSleepSeconds)
	assert.Equal(t, 6300, s.DailySleepDTO.REMSleepSeconds)
	assert.Equal(t, 96.0, s.DailySleepDTO.SpO2AvgReadingPercent)
	assert.Equal(t, "GOOD", s.DailySleepDTO.SleepScoreFeedback)
	assert.Equal(t, 5, s.RestlessMomentsCount)
}

func TestHRVData(t *testing.T) {
	c, stop := newVCRClient(t, "hrv_data")
	defer stop()

	h, err := c.HRVData(testDate)
	require.NoError(t, err)

	assert.Equal(t, 45, h.HRVSummary.WeeklyAvg)
	assert.Equal(t, 48, h.HRVSummary.LastNight)
	assert.Equal(t, "BALANCED", h.HRVSummary.Status)
	assert.Len(t, h.HRVReadings, 3)
	assert.Equal(t, 48, h.HRVReadings[0].HRVValue)
}
