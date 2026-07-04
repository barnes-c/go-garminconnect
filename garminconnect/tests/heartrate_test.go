package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeartRates(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	hr, err := c.HeartRates(t.Context(), testDate)
	require.NoError(t, err)

	assert.NotEmpty(t, hr.CalendarDate)
	if hr.RestingHeartRate == 0 {
		t.Skip("no heart rate data for test date")
	}
	assert.NotZero(t, hr.RestingHeartRate)
	assert.NotZero(t, hr.MinHeartRate)
	assert.NotZero(t, hr.MaxHeartRate)
	assert.NotZero(t, hr.LastSevenDaysAvgRestingHeartRate)
	if assert.NotEmpty(t, hr.HeartRateValues) {
		assert.NotZero(t, hr.HeartRateValues[0][0])
		assert.NotZero(t, hr.HeartRateValues[0][1])
	}
}
