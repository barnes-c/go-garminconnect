package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeartRates(t *testing.T) {
	c, stop := newVCRClient(t, "heart_rates")
	defer stop()

	hr, err := c.HeartRates(testDate)
	require.NoError(t, err)

	assert.Equal(t, "2026-01-01", hr.CalendarDate)
	assert.Equal(t, 46, hr.RestingHeartRate)
	assert.Equal(t, 42, hr.MinHeartRate)
	assert.Equal(t, 93, hr.MaxHeartRate)
	assert.Equal(t, 48, hr.LastSevenDaysAvgRestingHeartRate)
	assert.NotEmpty(t, hr.HeartRateValues)
	assert.Equal(t, int64(1778364000000), hr.HeartRateValues[0][0])
	assert.Equal(t, int64(56), hr.HeartRateValues[0][1])
}
