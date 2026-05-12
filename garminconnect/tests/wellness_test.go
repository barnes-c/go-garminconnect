package garminconnect_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSummary(t *testing.T) {
	c, stop := newVCRClient(t, "user_summary")
	defer stop()

	s, err := c.UserSummary(testDate)
	require.NoError(t, err)

	assert.Equal(t, 127516254, s.UserProfileID)
	assert.Equal(t, 1405, s.TotalSteps)
	assert.Equal(t, 6730, s.DailyStepGoal)
	assert.Equal(t, 1347.0, s.TotalKilocalories)
	assert.Equal(t, 46, s.RestingHeartRate)
	assert.Equal(t, 69, s.BodyBatteryMostRecentValue)
}

func TestAllDayStress(t *testing.T) {
	c, stop := newVCRClient(t, "all_day_stress")
	defer stop()

	s, err := c.AllDayStress(testDate)
	require.NoError(t, err)

	assert.Equal(t, 127516254, s.UserProfilePK)
	assert.Equal(t, 17, s.AvgStressLevel)
	assert.Equal(t, 91, s.MaxStressLevel)
	assert.NotEmpty(t, s.StressValuesArray)
}

func TestBodyBattery(t *testing.T) {
	c, stop := newVCRClient(t, "body_battery")
	defer stop()

	start := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	entries, err := c.BodyBattery(start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, entries)
}

func TestFloors(t *testing.T) {
	c, stop := newVCRClient(t, "floors")
	defer stop()

	f, err := c.Floors(testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, f)
}

func TestHydration(t *testing.T) {
	c, stop := newVCRClient(t, "hydration")
	defer stop()

	h, err := c.Hydration(testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, h)
}

func TestRespiration(t *testing.T) {
	c, stop := newVCRClient(t, "respiration")
	defer stop()

	r, err := c.Respiration(testDate)
	require.NoError(t, err)

	assert.Equal(t, 14.0, r.TodayAvgWakingRespirationValue)
	assert.Equal(t, 21.0, r.HighestRespirationValue)
	assert.Equal(t, 7.0, r.LowestRespirationValue)
	assert.NotEmpty(t, r.RespirationValuesArray)
}

func TestSpO2(t *testing.T) {
	c, stop := newVCRClient(t, "spo2")
	defer stop()

	s, err := c.SpO2(testDate)
	require.NoError(t, err)

	assert.Equal(t, 94.0, s.AverageSpO2)
	assert.Equal(t, 85.0, s.LowestSpO2)
	assert.InDelta(t, 94.714, s.LastSevenDaysAvgSpO2, 0.001)
	assert.NotEmpty(t, s.SpO2HourlyAverages)
}

func TestSteps(t *testing.T) {
	c, stop := newVCRClient(t, "steps")
	defer stop()

	steps, err := c.Steps(testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, steps)
}

func TestRestingHeartRate(t *testing.T) {
	c, stop := newVCRClient(t, "resting_heart_rate")
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	out, err := c.RestingHeartRate(start, testDate)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestDailySteps(t *testing.T) {
	c, stop := newVCRClient(t, "daily_steps")
	defer stop()

	start := testDate.AddDate(0, 0, -7)
	entries, err := c.DailySteps(start, testDate)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestWeeklyStress(t *testing.T) {
	c, stop := newVCRClient(t, "weekly_stress")
	defer stop()

	out, err := c.WeeklyStress(testDate, 4)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestWeeklyIntensityMinutes(t *testing.T) {
	c, stop := newVCRClient(t, "weekly_intensity_minutes")
	defer stop()

	start := testDate.AddDate(0, 0, -7)
	out, err := c.WeeklyIntensityMinutes(start, testDate)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestBloodPressure(t *testing.T) {
	c, stop := newVCRClient(t, "blood_pressure")
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	out, err := c.BloodPressure(start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}
