package garminconnect_test

import (
	"testing"
	"time"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJetLagTripLifecycle creates, edits and deletes a throwaway trip. Garmin
// refuses trips before the previous month, so the trip is in the future; its
// dates only appear in request bodies, which replay doesn't match on.
func TestJetLagTripLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	day := time.Now().AddDate(0, 0, 30).Format(time.DateOnly)
	trip, err := c.CreateJetLagTrip(ctx, gc.JetLagTrip{
		StartDate: day,
		EndDate:   day,
		Flights: []gc.JetLagFlight{{
			OriginLocation:       "ZRH",
			OriginTimezone:       "+01:00",
			OriginTimestamp:      day + "T10:00:00.00",
			DestinationLocation:  "JFK",
			DestinationTimezone:  "-05:00",
			DestinationTimestamp: day + "T13:00:00.00",
		}},
	})
	require.NoError(t, err)
	require.NotZero(t, trip.TripID)
	require.NotEmpty(t, trip.Flights)
	assert.NotZero(t, trip.Flights[0].FlightID)

	trips, err := c.JetLagTrips(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, trips)

	trip.Flights[0].DestinationLocation = "EWR"
	require.NoError(t, c.UpdateJetLagTrip(ctx, trip))
	require.NoError(t, c.DeleteJetLagTrip(ctx, trip.TripID))
}

func TestNutritionStatus(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.NutritionStatus(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, out.CurrentStatus)
}

// The tests below need a watch, so they are recorded against a real account
// on a day with data (GARMIN_SUMMARY_DATE / GARMIN_SLEEP_DATE). Ranges query a
// single day; see dataDate.

func TestSleepStatsWeekly(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.SleepStatsWeekly(t.Context(), dataDate(t, "GARMIN_SLEEP_DATE"), 2)
	require.NoError(t, err)
	assert.NotZero(t, out.Overall.AverageSleepSeconds)
	require.NotEmpty(t, out.Weeks)
	assert.NotZero(t, out.Weeks[0].Values.DeepTime)
}

func TestSleepStatsOverall(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	d := dataDate(t, "GARMIN_SLEEP_DATE")

	out, err := c.SleepStatsOverall(t.Context(), d, d)
	require.NoError(t, err)
	assert.NotZero(t, out.Overall.AverageSleepSeconds)
	assert.NotZero(t, out.Overall.AverageSpO2)
}

func TestActivityStatistics(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.ActivityStatistics(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotEmpty(t, out[0].ActivityType.TypeKey)
	assert.NotZero(t, out[0].TotalActivities)
}

func TestWellnessAverages(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	d := dataDate(t, "GARMIN_SUMMARY_DATE")

	out, err := c.WellnessAverages(t.Context(), d, d)
	require.NoError(t, err)
	assert.NotZero(t, out.AverageSteps)
}

func TestDailyMovement(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.DailyMovement(t.Context(), dataDate(t, "GARMIN_SUMMARY_DATE"))
	require.NoError(t, err)
	assert.NotEmpty(t, out.Descriptors)
	assert.NotEmpty(t, out.Values)
}

func TestIntensityMinutesStats(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.IntensityMinutesStats(t.Context(), dataDate(t, "GARMIN_SUMMARY_DATE"), 1)
	require.NoError(t, err)
	assert.NotZero(t, out.WeeklyGoal)
	assert.NotEmpty(t, out.Days)
}

func TestMonthlySteps(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.MonthlySteps(t.Context(), dataDate(t, "GARMIN_SUMMARY_DATE"), 1)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].Values.TotalSteps)
}

func TestWeekAverageRestingHR(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.WeekAverageRestingHR(t.Context(), dataDate(t, "GARMIN_SUMMARY_DATE"))
	require.NoError(t, err)
	assert.NotZero(t, out)
}

func TestLastSyncTimestamp(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.LastSyncTimestamp(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

// TestMenstrualCycleLifecycle sets up cycle tracking, confirms a period, logs a
// day and checks the reads, then clears the log and the period again.
func TestMenstrualCycleLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	start := testDate.AddDate(0, 0, -12)
	logDay := start.AddDate(0, 0, 2)
	rangeStart := testDate.AddDate(0, 0, -31)

	err := c.UpdateUserSettings(ctx, map[string]any{
		"userMenstrualCycleSettings": map[string]any{"menstrualCycleType": "REGULAR"},
	})
	skipAPIError(t, err)
	require.NoError(t, err)

	cycle, err := c.InitMenstrualCycle(ctx, start, 5, 28)
	require.NoError(t, err)
	assert.NotEmpty(t, cycle.PeriodStartDate)
	cycle, err = c.ConfirmMenstrualPeriod(ctx, start, 5, 28, false)
	require.NoError(t, err)
	assert.NotZero(t, cycle.CycleLength)

	flow, notes := "LIGHT", "note"
	logged, err := c.SetMenstrualDayLog(ctx, gc.MenstrualDayLog{
		CalendarDate: logDay.Format(time.DateOnly),
		Symptoms:     []string{"CRAMPS"},
		Moods:        []string{"HAPPY"},
		Flow:         &flow,
		Notes:        &notes,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, logged.Symptoms)

	view, err := c.MenstrualDayView(ctx, logDay)
	require.NoError(t, err)
	require.NotNil(t, view.DaySummary)
	assert.NotZero(t, view.DaySummary.DayInCycle)
	require.NotNil(t, view.DayLog)
	assert.NotEmpty(t, view.DayLog.Moods)

	cal, err := c.MenstrualCalendar(ctx, rangeStart, testDate)
	require.NoError(t, err)
	assert.NotEmpty(t, cal.CycleSummaries)
	assert.NotEmpty(t, cal.LoggedSymptomDays)

	day, err := c.MenstrualCycleDay(ctx, testDate)
	require.NoError(t, err)
	require.NotNil(t, day)
	assert.NotEmpty(t, day.CycleType)

	last, err := c.LastConfirmedMenstrualCycle(ctx, testDate)
	require.NoError(t, err)
	require.NotNil(t, last)
	assert.NotZero(t, last.PeriodLength)

	report, err := c.MenstrualReport(ctx, testDate, 1)
	require.NoError(t, err)
	assert.NotEmpty(t, report.CycleSummaries)
	assert.NotEmpty(t, report.DailyLogs)

	empty := ""
	_, err = c.SetMenstrualDayLog(ctx, gc.MenstrualDayLog{CalendarDate: logDay.Format(time.DateOnly), Notes: &empty})
	require.NoError(t, err)
	require.NoError(t, c.SetMenstrualPeriods(ctx, rangeStart, testDate, nil))
	cal, err = c.MenstrualCalendar(ctx, rangeStart, testDate)
	require.NoError(t, err)
	assert.Empty(t, cal.CycleSummaries)
}

func TestPregnancySummary(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// The recording account isn't tracking a pregnancy: Garmin answers 204.
	_, err := c.PregnancySummary(t.Context())
	require.ErrorIs(t, err, gc.ErrNoData)
}
