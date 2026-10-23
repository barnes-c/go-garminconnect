package garminconnect_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSummary(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	s, err := c.UserSummary(t.Context(), testDate)
	require.NoError(t, err)
	assert.NotZero(t, s.FloorsAscendedGoal)

	assert.NotZero(t, s.UserProfileID)
	assert.NotZero(t, s.DailyStepGoal)
}

func TestAllDayStress(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	s, err := c.AllDayStress(t.Context(), testDate)
	require.NoError(t, err)

	assert.NotZero(t, s.UserProfilePK)
	assert.NotEmpty(t, s.StressValuesArray)
}

func TestBodyBattery(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	entries, err := c.BodyBattery(t.Context(), testDate, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, entries)
}

func TestFloors(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	f, err := c.Floors(t.Context(), testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, f)
}

func TestHydration(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	h, err := c.Hydration(t.Context(), testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotZero(t, h.UserProfilePK)
	assert.NotNil(t, h)
}

func TestRespiration(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	r, err := c.Respiration(t.Context(), testDate)
	require.NoError(t, err)

	if r.TodayAvgWakingRespirationValue == 0 {
		t.Skip("no respiration data for test date")
	}
	assert.NotZero(t, r.TodayAvgWakingRespirationValue)
	assert.NotZero(t, r.HighestRespirationValue)
	assert.NotZero(t, r.LowestRespirationValue)
	assert.NotEmpty(t, r.RespirationValuesArray)
}

func TestSpO2(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	s, err := c.SpO2(t.Context(), testDate)
	require.NoError(t, err)

	if s.AverageSpO2 == 0 {
		t.Skip("no SpO2 data for test date")
	}
	assert.NotZero(t, s.AverageSpO2)
	assert.NotZero(t, s.LowestSpO2)
	assert.NotZero(t, s.LastSevenDaysAvgSpO2)
	assert.NotEmpty(t, s.SpO2HourlyAverages)
}

func TestSteps(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	steps, err := c.Steps(t.Context(), testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, steps)
}

func TestRestingHeartRate(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	out, err := c.RestingHeartRate(t.Context(), start, testDate)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestDailySteps(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, 0, -7)
	entries, err := c.DailySteps(t.Context(), start, testDate)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestWeeklyStress(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.WeeklyStress(t.Context(), testDate, 4)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].AvgStressLevel)
	assert.NotNil(t, out)
}

func TestWeeklyIntensityMinutes(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, 0, -7)
	out, err := c.WeeklyIntensityMinutes(t.Context(), start, testDate)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].ModerateIntensityMinutes)
	assert.NotZero(t, out[0].VigorousIntensityMinutes)
	assert.NotNil(t, out)
}

func TestBloodPressure(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	out, err := c.BloodPressure(t.Context(), start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestCaloriesDaily(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, 0, -6)
	out, err := c.CaloriesDaily(t.Context(), start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	for _, d := range out {
		assert.NotEmpty(t, d.CalendarDate)
		assert.InDelta(t, d.Active+d.Resting, d.Total, 0.001)
	}
}

// TestAddHydration logs an intake and then subtracts it again; hydration
// entries can't be deleted, so this is how recording leaves the day unchanged.
func TestAddHydration(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	zurich, err := time.LoadLocation("Europe/Zurich")
	require.NoError(t, err)
	at := time.Date(testDate.Year(), testDate.Month(), testDate.Day(), 8, 0, 0, 0, zurich)

	out, err := c.AddHydration(ctx, 250, at)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	_, err = c.AddHydration(ctx, -250, at.Add(time.Minute))
	require.NoError(t, err)
}
