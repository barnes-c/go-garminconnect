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

	assert.Equal(t, 987654321, s.UserProfileID)
	assert.Equal(t, 12543, s.TotalSteps)
	assert.Equal(t, 10000, s.DailyStepGoal)
	assert.Equal(t, 2850.0, s.TotalKilocalories)
	assert.Equal(t, 52, s.RestingHeartRate)
	assert.Equal(t, 72, s.BodyBatteryMostRecentValue)
}

func TestAllDayStress(t *testing.T) {
	c, stop := newVCRClient(t, "all_day_stress")
	defer stop()

	s, err := c.AllDayStress(testDate)
	require.NoError(t, err)

	assert.Equal(t, 987654321, s.UserProfilePK)
	assert.Equal(t, 28, s.AvgStressLevel)
	assert.Equal(t, 72, s.MaxStressLevel)
	assert.Len(t, s.StressValuesArray, 4)
}

func TestBodyBattery(t *testing.T) {
	c, stop := newVCRClient(t, "body_battery")
	defer stop()

	start := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	entries, err := c.BodyBattery(start, testDate)
	require.NoError(t, err)

	require.Len(t, entries, 2)
	assert.Equal(t, "2026-05-01T00:00:00.0", entries[0].StartTimestampGMT)
	require.Len(t, entries[0].BodyBatteryValues, 3)
	assert.Equal(t, 85, entries[0].BodyBatteryValues[0].Value)
	assert.Equal(t, "CHARGED", entries[0].BodyBatteryValues[0].Status)
}

func TestFloors(t *testing.T) {
	c, stop := newVCRClient(t, "floors")
	defer stop()

	f, err := c.Floors(testDate)
	require.NoError(t, err)

	assert.Equal(t, 987654321, f.UserProfilePK)
	assert.Equal(t, "2026-05-10", f.CalendarDate)
	assert.Len(t, f.FloorValuesArray, 3)
}

func TestHydration(t *testing.T) {
	c, stop := newVCRClient(t, "hydration")
	defer stop()

	h, err := c.Hydration(testDate)
	require.NoError(t, err)

	assert.Equal(t, 2100.0, h.ValueInML)
	assert.Equal(t, 2500.0, h.GoalInML)
	assert.Equal(t, 800.0, h.SweatLossInML)
}

func TestRespiration(t *testing.T) {
	c, stop := newVCRClient(t, "respiration")
	defer stop()

	r, err := c.Respiration(testDate)
	require.NoError(t, err)

	assert.Equal(t, 14.8, r.TodayAvgWakingRespirationValue)
	assert.Equal(t, 20.0, r.HighestRespirationValue)
	assert.Equal(t, 11.0, r.LowestRespirationValue)
	assert.Len(t, r.RespirationValuesArray, 3)
}

func TestSpO2(t *testing.T) {
	c, stop := newVCRClient(t, "spo2")
	defer stop()

	s, err := c.SpO2(testDate)
	require.NoError(t, err)

	assert.Equal(t, 96.5, s.AverageSpO2)
	assert.Equal(t, 93.0, s.LowestSpO2)
	assert.Equal(t, 95.8, s.LastSevenDaysAvgSpO2)
	assert.Len(t, s.SpO2HourlyAverages, 3)
}

func TestSteps(t *testing.T) {
	c, stop := newVCRClient(t, "steps")
	defer stop()

	steps, err := c.Steps(testDate)
	require.NoError(t, err)

	require.Len(t, steps, 3)
	assert.Equal(t, 350, steps[0].Steps)
	assert.Equal(t, "active", steps[0].PrimaryActivityLevel)
	assert.Equal(t, 80, steps[2].Steps)
	assert.Equal(t, "sedentary", steps[2].PrimaryActivityLevel)
}
