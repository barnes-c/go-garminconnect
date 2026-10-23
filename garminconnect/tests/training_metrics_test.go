package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These need a watch (and a power meter for the cycling ones), so they are
// recorded against a real account on a day with data. Ranges query a single
// day; see dataDate.

func TestTrainingReadinessRange(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	d := dataDate(t, "GARMIN_SUMMARY_DATE")

	out, err := c.TrainingReadinessRange(t.Context(), d, d)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].Score)
}

func TestTrainingStatusHistory(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	d := dataDate(t, "GARMIN_SUMMARY_DATE")

	out, err := c.TrainingStatusHistory(t.Context(), d, d)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	for _, days := range out {
		require.NotEmpty(t, days)
		assert.NotEmpty(t, days[0].CalendarDate)
	}
}

func TestHeatAltitudeAcclimation(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.HeatAltitudeAcclimation(t.Context(), dataDate(t, "GARMIN_SUMMARY_DATE"))
	require.NoError(t, err)
	assert.NotEmpty(t, out.CalendarDate)
	assert.NotZero(t, out.CurrentAltitude)
}

func TestCurrentBiometrics(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.CurrentBiometrics(t.Context())
	require.NoError(t, err)
	require.NotNil(t, out.LactateThresholdHeartRate)
	assert.NotZero(t, out.LactateThresholdHeartRate.Value)
}

func TestFitnessAgeHistory(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	d := dataDate(t, "GARMIN_SUMMARY_DATE")

	out, err := c.FitnessAgeHistory(t.Context(), d, d)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].Values.FitnessAge)
}

func TestPowerCurve(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	d := dataDate(t, "GARMIN_SUMMARY_DATE")

	out, err := c.PowerCurve(t.Context(), d, d, "cycling")
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].Duration)
	assert.NotZero(t, out[0].Power)
}

func TestCyclingAbility(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.CyclingAbility(t.Context(), dataDate(t, "GARMIN_SUMMARY_DATE"))
	require.NoError(t, err)
	assert.NotZero(t, out.AerobicEndurance)
	assert.NotEmpty(t, out.ProfileType)
}
