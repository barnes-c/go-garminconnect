package garminconnect_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/barnes-c/go-garminconnect/garminconnect"
)

func TestUserSummary(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	s, err := c.UserSummary(t.Context(), testDate)
	require.NoError(t, err)
	assert.NotZero(t, s.FloorsAscendedGoal)

	assert.NotZero(t, s.UserProfileID)
	assert.NotZero(t, s.DailyStepGoal)
	assert.NotZero(t, s.TotalStressDuration)
	assert.NotZero(t, s.AverageStressLevel)
	assert.NotEmpty(t, s.StressQualifier)
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

func TestBodyBatteryEvents(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	events, err := c.BodyBatteryEvents(t.Context(), dataDate(t, "GARMIN_SUMMARY_DATE"))
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, events)
	assert.NotEmpty(t, events[0].Event.EventType)
	assert.NotEmpty(t, events[0].Event.EventStartTimeGMT)
	assert.NotZero(t, events[0].Event.DurationInMilliseconds)
	assert.NotEmpty(t, events[0].BodyBatteryValueDescriptorsDTOList)
	assert.NotEmpty(t, events[0].BodyBatteryValuesArray)
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

	r, err := c.Respiration(t.Context(), dataDate(t, "GARMIN_SLEEP_DATE"))
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

	s, err := c.SpO2(t.Context(), dataDate(t, "GARMIN_SLEEP_DATE"))
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

func TestBloodPressure(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	out, err := c.BloodPressure(t.Context(), start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
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
	assert.NotEmpty(t, out.CalendarDate)
	_, err = c.AddHydration(ctx, -250, at.Add(time.Minute))
	require.NoError(t, err)
}

func TestHydrationDetails(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.HydrationDetails(t.Context(), testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out.CalendarDate)
	assert.NotEmpty(t, out.HydrationMeasurementUnit)
	assert.NotZero(t, out.GoalInML)
}

// TestUserSummaryDecodesDecimals guards the fields that were only ever null in
// the recordings: a real account sends averageSpo2 as a decimal, which failed
// the whole UserSummary call when the field was an int.
func TestUserSummaryDecodesDecimals(t *testing.T) {
	var s garminconnect.UserSummary
	err := json.Unmarshal([]byte(`{"averageSpo2":95.0,"lowestSpo2":88.0,"latestSpo2":96.0,
		"maxStressLevel":93.0,"minAvgHeartRate":49.0,"maxAvgHeartRate":141.0,
		"bodyBatteryAtWakeTime":78.0,"bodyBatteryDuringSleep":52.0}`), &s)
	require.NoError(t, err)
	assert.InDelta(t, 95.0, s.AverageSpO2, 0.001)
	assert.InDelta(t, 78.0, s.BodyBatteryAtWakeTime, 0.001)
}
