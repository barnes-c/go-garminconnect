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

	assert.Equal(t, "2026-05-10", hr.CalendarDate)
	assert.Equal(t, 52, hr.RestingHeartRate)
	assert.Equal(t, 48, hr.MinHeartRate)
	assert.Equal(t, 178, hr.MaxHeartRate)
	assert.Equal(t, 54, hr.LastSevenDaysAvgRestingHeartRate)
	assert.Len(t, hr.HeartRateValues, 3)
	assert.Equal(t, int64(1746835200000), hr.HeartRateValues[0][0])
	assert.Equal(t, int64(52), hr.HeartRateValues[0][1])
}
