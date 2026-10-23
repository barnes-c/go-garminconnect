package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Gear represents a piece of equipment tracked in Garmin Connect.
type Gear struct {
	GearPK          int64   `json:"gearPk"`
	UUID            string  `json:"uuid"`
	GearTypeName    string  `json:"gearTypeName"`
	DisplayName     string  `json:"displayName"`
	CustomMakeModel string  `json:"customMakeModel"`
	MaxMeters       float64 `json:"maximumMeters"`
	NotifiedMeters  int     `json:"notifiedAtMeters"`
	DateBegin       string  `json:"dateBegin"`
	DateEnd         string  `json:"dateEnd"`
	GearStatusName  string  `json:"gearStatusName"`
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

// GearStats holds lifetime usage statistics for a piece of gear.
type GearStats struct {
	GearPK          int64   `json:"gearPk"`
	UUID            string  `json:"uuid"`
	TotalDistance   float64 `json:"totalDistance"` // meters
	TotalActivities int     `json:"totalActivities"`
	IsProcessing    bool    `json:"isProcessing"`
	Processing      bool    `json:"processing"`
	CreateDate      int64   `json:"createDate"`
	UpdateDate      int64   `json:"updateDate"`
}

// GearStats returns lifetime usage statistics for the given gear UUID.
func (c *Client) GearStats(ctx context.Context, gearUUID string) (*GearStats, error) {
	var out GearStats
	if err := c.get(ctx, fmt.Sprintf("/gear-service/gear/stats/%s", gearUUID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
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

// GearStatus is the lifecycle state of a piece of gear.
type GearStatus string

const (
	GearActive  GearStatus = "ACTIVE"
	GearRetired GearStatus = "RETIRED"
)

// SetGearStatus retires or re-activates the given gear.
func (c *Client) SetGearStatus(ctx context.Context, gearUUID string, status GearStatus) error {
	return c.put(ctx, fmt.Sprintf("/gear-service/gear/v2/%s/status/%s", gearUUID, status), nil, nil)
}

// DeleteGear permanently deletes the given gear.
func (c *Client) DeleteGear(ctx context.Context, gearUUID string) error {
	return c.del(ctx, fmt.Sprintf("/gear-service/gear/v2/%s", gearUUID))
}

// AddGearToActivity links the given gear to an activity.
func (c *Client) AddGearToActivity(ctx context.Context, gearUUID string, activityID int64) error {
	return c.put(ctx, fmt.Sprintf("/gear-service/gear/link/%s/activity/%d", gearUUID, activityID), nil, nil)
}

// RemoveGearFromActivity unlinks the given gear from an activity.
func (c *Client) RemoveGearFromActivity(ctx context.Context, gearUUID string, activityID int64) error {
	return c.put(ctx, fmt.Sprintf("/gear-service/gear/unlink/%s/activity/%d", gearUUID, activityID), nil, nil)
}

// GearDetail is the full (v2) record of a piece of gear, as used by
// CreateGear, GearDetails and UpdateGear.
type GearDetail struct {
	UUID                    string             `json:"uuid,omitempty"`
	GearType                string             `json:"gearType"` // e.g. "SHOES"
	Brand                   string             `json:"brand"`
	Model                   string             `json:"model"`
	Name                    string             `json:"name"`
	FirstUseDate            string             `json:"firstUseDate"` // YYYY-MM-DD
	UsageType               string             `json:"usageType"`    // "DISTANCE" or "TIME"
	MaxUsageDistanceMeters  float64            `json:"maxUsageDistanceMeters"`
	MaxUsageDurationSeconds int                `json:"maxUsageDurationSeconds"`
	Notes                   string             `json:"notes,omitempty"` // Garmin rejects an empty note
	AssociatedActivityTypes []GearActivityType `json:"associatedActivityTypes,omitempty"`

	// Read-only usage fields, set by Garmin.
	Status              GearStatus `json:"status,omitempty"`
	NumActivitiesLinked int        `json:"numActivitiesLinked,omitempty"`
	DistanceUsedMeters  float64    `json:"distanceUsedMeters,omitempty"`
	DurationUsedSeconds int        `json:"durationUsedSeconds,omitempty"`
	DaysUsed            int        `json:"daysUsed,omitempty"`
	CreatedDate         string     `json:"createdDate,omitempty"`
}

// GearActivityType associates gear with an activity type, e.g. to make it the
// default gear for "running".
type GearActivityType struct {
	ActivityTypeKey string `json:"activityTypeKey"`
	DefaultGear     bool   `json:"defaultGear"`
	PreferredGear   bool   `json:"preferredGear"`
}

// GearDetails returns the full record of the given gear.
func (c *Client) GearDetails(ctx context.Context, gearUUID string) (*GearDetail, error) {
	var out GearDetail
	if err := c.get(ctx, fmt.Sprintf("/gear-service/gear/v2/%s", gearUUID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateGear adds a new piece of gear and returns the created record.
func (c *Client) CreateGear(ctx context.Context, g GearDetail) (*GearDetail, error) {
	var out GearDetail
	if err := c.post(ctx, "/gear-service/gear/v2", g, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateGear saves g, identified by g.UUID. Garmin rejects partial records, so
// start from GearDetails and change the fields you need.
func (c *Client) UpdateGear(ctx context.Context, g *GearDetail) error {
	return c.put(ctx, fmt.Sprintf("/gear-service/gear/v2/%s", g.UUID), g, nil)
}

// GearType is a kind of gear Garmin tracks (bike, shoes, ...).
type GearType struct {
	GearTypePK   int    `json:"gearTypePk"`
	GearTypeName string `json:"gearTypeName"` // e.g. "Shoes"
}

// GearTypes returns every gear type.
func (c *Client) GearTypes(ctx context.Context) ([]GearType, error) {
	var out []GearType
	if err := c.get(ctx, "/gear-service/gear/types", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UserGearTypes returns the gear types the user has gear for (e.g. "SHOES").
func (c *Client) UserGearTypes(ctx context.Context) ([]string, error) {
	var out struct {
		ActiveGearTypes []string `json:"activeGearTypes"`
	}
	if err := c.get(ctx, "/gear-service/gear/v2/user-gear-types", nil, &out); err != nil {
		return nil, err
	}
	return out.ActiveGearTypes, nil
}

// GearMakeModel is a product from Garmin's gear catalog.
type GearMakeModel struct {
	MakeModelName string `json:"makeModelName"` // e.g. "Nike Air Zoom Pegasus"
	Make          struct {
		GearMakePK   int    `json:"gearMakePk"`
		GearMakeName string `json:"gearMakeName"`
	} `json:"gearMakeDTO"`
	Model struct {
		GearModelPK   int    `json:"gearModelPk"`
		GearTypePK    int    `json:"gearTypePk"`
		GearModelName string `json:"gearModelName"`
	} `json:"gearModelDTO"`
}

// SearchGearMakeModels searches Garmin's gear catalog.
func (c *Client) SearchGearMakeModels(ctx context.Context, keywords string) ([]GearMakeModel, error) {
	var out []GearMakeModel
	if err := c.get(ctx, "/gear-service/gear/makemodels", url.Values{"keywords": {keywords}}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GearList returns all of the user's gear as full records.
func (c *Client) GearList(ctx context.Context) ([]GearDetail, error) {
	var out []GearDetail
	if err := c.get(ctx, "/gear-service/gear/v2/list", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityGear returns the gear linked to an activity.
func (c *Client) ActivityGear(ctx context.Context, activityID int64) ([]GearDetail, error) {
	var out []GearDetail
	if err := c.get(ctx, fmt.Sprintf("/gear-service/activity/v2/%d", activityID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
