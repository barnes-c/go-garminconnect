package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Gear represents a piece of equipment tracked in Garmin Connect.
type Gear struct {
	GearPK          int64  `json:"gearPk"`
	UUID            string `json:"uuid"`
	GearTypeName    string `json:"gearTypeName"`
	DisplayName     string `json:"displayName"`
	CustomMakeModel string `json:"customMakeModel"`
	MaxMeters       int    `json:"maxMeters"`
	NotifiedMeters  int    `json:"notifiedAtMeters"`
	DateBegin       string `json:"dateBegin"`
	DateEnd         string `json:"dateEnd"`
	GearStatusName  string `json:"gearStatusName"`
}

// Gear returns all gear registered to the given user profile number.
func (c *Client) Gear(ctx context.Context, userProfileNumber int) ([]Gear, error) {
	params := url.Values{"userProfilePk": {fmt.Sprintf("%d", userProfileNumber)}}
	var out []Gear
	if err := c.get(ctx, "/gear-service/gear/filterGear", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GearStats returns usage statistics for the given gear UUID.
func (c *Client) GearStats(ctx context.Context, gearUUID string) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/gear-service/gear/stats/%s", gearUUID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GearActivities returns activities that used the given gear.
func (c *Client) GearActivities(ctx context.Context, gearUUID string, start, limit int) (map[string]json.RawMessage, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/gear-service/gear/%s/activities", gearUUID), params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GearDefaults returns the default gear assigned per activity type for the user.
func (c *Client) GearDefaults(ctx context.Context, userProfileNumber int) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/gear-service/gear/user/%d/activityTypes", userProfileNumber), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
