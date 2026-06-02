package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSleepData(t *testing.T) {
	c, stop := newVCRClient(t, "sleep_data")
	defer stop()

	s, err := c.SleepData(t.Context(), testDate)
	require.NoError(t, err)

	assert.NotEmpty(t, s.DailySleepDTO.CalendarDate)
	if s.DailySleepDTO.SleepTimeSeconds == 0 {
		t.Skip("no sleep data for test date")
	}
	assert.NotZero(t, s.DailySleepDTO.SleepTimeSeconds)
	assert.NotZero(t, s.DailySleepDTO.DeepSleepSeconds)
	assert.NotZero(t, s.DailySleepDTO.REMSleepSeconds)
	assert.NotZero(t, s.DailySleepDTO.SpO2AvgReadingPercent)
	assert.NotEmpty(t, s.DailySleepDTO.SleepScoreFeedback)
}

func TestHRVData(t *testing.T) {
	c, stop := newVCRClient(t, "hrv_data")
	defer stop()

	h, err := c.HRVData(t.Context(), testDate)
	require.NoError(t, err)

	assert.NotZero(t, h.HRVSummary.WeeklyAvg)
	assert.NotEmpty(t, h.HRVSummary.Status)
	if len(h.HRVReadings) == 0 {
		t.Skip("no HRV readings for test date")
	}
	assert.NotZero(t, h.HRVReadings[0].HRVValue)
}
