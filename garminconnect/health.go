package garminconnect

import (
	"encoding/json"
	"fmt"
	"time"
)

// MenstrualData returns menstrual cycle data for the given date.
func (c *Client) MenstrualData(d time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/women-health-service/menstrualcycle/dayview/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// MenstrualCalendar returns menstrual cycle data between start and end dates.
func (c *Client) MenstrualCalendar(start, end time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/women-health-service/menstrualcycle/calendar/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PregnancySummary returns the current pregnancy snapshot for the authenticated user.
func (c *Client) PregnancySummary() (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get("/women-health-service/pregnancy/snapshot", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
