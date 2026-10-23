package garminconnect

import (
	"context"
	"fmt"
)

// JetLagFlight is one flight of a Jet Lag Adviser trip. Timestamps are local
// times in the form "2006-01-02T15:04:05.00" (Garmin rejects them without the
// fraction); timezones are UTC offsets like "+01:00".
type JetLagFlight struct {
	FlightID             int64  `json:"flightId,omitempty"`
	TripID               int64  `json:"tripId,omitempty"`
	OriginLocation       string `json:"originLocation"` // e.g. "ZRH"
	OriginTimezone       string `json:"originTimezone"`
	OriginTimestamp      string `json:"originTimestamp"`
	DestinationLocation  string `json:"destinationLocation"`
	DestinationTimezone  string `json:"destinationTimezone"`
	DestinationTimestamp string `json:"destinationTimestamp"`
}

// JetLagTrip is a Jet Lag Adviser trip. Trips may not start before the first
// day of the previous month.
type JetLagTrip struct {
	TripID    int64          `json:"tripId,omitempty"`
	UserID    int64          `json:"userId,omitempty"` // required by UpdateJetLagTrip; set by Garmin
	StartDate string         `json:"startDate"`        // YYYY-MM-DD
	EndDate   string         `json:"endDate"`
	Flights   []JetLagFlight `json:"flights"`
}

// JetLagTrips returns all of the user's Jet Lag Adviser trips.
func (c *Client) JetLagTrips(ctx context.Context) ([]JetLagTrip, error) {
	var out []JetLagTrip
	if err := c.get(ctx, "/jetlag-service/jetlag/trip/all", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateJetLagTrip adds a trip and returns it with its trip and flight IDs.
func (c *Client) CreateJetLagTrip(ctx context.Context, trip JetLagTrip) (*JetLagTrip, error) {
	var out JetLagTrip
	if err := c.post(ctx, "/jetlag-service/jetlag/trip", trip, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateJetLagTrip replaces a trip, identified by trip.TripID. Send every
// flight, keeping their FlightIDs; flights left out are removed.
func (c *Client) UpdateJetLagTrip(ctx context.Context, trip *JetLagTrip) error {
	return c.put(ctx, "/jetlag-service/jetlag/trip", trip, nil)
}

// DeleteJetLagTrip removes a trip.
func (c *Client) DeleteJetLagTrip(ctx context.Context, id int64) error {
	return c.del(ctx, fmt.Sprintf("/jetlag-service/jetlag/trip/%d", id))
}
