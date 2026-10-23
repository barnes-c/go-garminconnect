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
