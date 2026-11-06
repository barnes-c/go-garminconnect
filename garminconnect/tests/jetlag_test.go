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
