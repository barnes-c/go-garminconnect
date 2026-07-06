package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"time"
)

// Activity represents a single Garmin Connect activity summary.
type Activity struct {
	ActivityID                            int64        `json:"activityId"`
	ActivityName                          string       `json:"activityName"`
	Description                           string       `json:"description"`
	ActivityType                          ActivityType `json:"activityType"`
	StartTimeGMT                          string       `json:"startTimeGMT"`
	StartTimeLocal                        string       `json:"startTimeLocal"`
	EndTimeGMT                            string       `json:"endTimeGMT"`
	Duration                              float64      `json:"duration"`        // seconds
	ElapsedDuration                       float64      `json:"elapsedDuration"` // seconds
	MovingDuration                        float64      `json:"movingDuration"`  // seconds
	Distance                              float64      `json:"distance"`        // meters
	Calories                              float64      `json:"calories"`
	AverageHR                             float64      `json:"averageHR"`
	MaxHR                                 float64      `json:"maxHR"`
	AverageSpeed                          float64      `json:"averageSpeed"` // meters/second
	MaxSpeed                              float64      `json:"maxSpeed"`     // meters/second
	ElevationGain                         float64      `json:"elevationGain"`
	ElevationLoss                         float64      `json:"elevationLoss"`
	StartLatitude                         *float64     `json:"startLatitude"` // nil for indoor activities
	StartLongitude                        *float64     `json:"startLongitude"`
	EndLatitude                           *float64     `json:"endLatitude"`
	EndLongitude                          *float64     `json:"endLongitude"`
	Steps                                 int64        `json:"steps"`
	TrainingEffect                        float64      `json:"trainingEffect"`
	AnaerobicTrainingEffect               float64      `json:"anaerobicTrainingEffect"`
	AerobicTrainingEffectMessage          string       `json:"aerobicTrainingEffectMessage"`
	AverageRunningCadenceInStepsPerMinute float64      `json:"averageRunningCadenceInStepsPerMinute"`
	VO2MaxValue                           float64      `json:"vO2MaxValue"`
	LocationName                          string       `json:"locationName"`
	OwnerID                               int64        `json:"ownerId"`
	HasPolyline                           bool         `json:"hasPolyline"`
}

// ActivityType identifies the sport type of an activity.
type ActivityType struct {
	TypeID       int    `json:"typeId"`
	TypeKey      string `json:"typeKey"`
	ParentTypeID int    `json:"parentTypeId"`
	IsHidden     bool   `json:"isHidden"`
	Restricted   bool   `json:"restricted"`
	Trimmable    bool   `json:"trimmable"`
}

// PersonalRecord represents a personal best for a given activity type.
type PersonalRecord struct {
	ID             int64   `json:"id"`
	TypeID         int64   `json:"typeId"`
	ActivityID     int64   `json:"activityId"`
	ActivityName   string  `json:"activityName"`
	ActivityType   string  `json:"activityType"`
	StartTimeGMT   string  `json:"startTimeGMT"`
	PrStartTimeGmt int64   `json:"prStartTimeGmt"` // epoch milliseconds
	Value          float64 `json:"value"`
	PrTypeLabelKey string  `json:"prTypeLabelKey"`
}

// Split is a single lap or segment summary within an activity.
type Split struct {
	StartTimeGMT   string  `json:"startTimeGMT"`
	Distance       float64 `json:"distance"`
	Duration       float64 `json:"duration"`
	MovingDuration float64 `json:"movingDuration"`
	ElevationGain  float64 `json:"elevationGain"`
	ElevationLoss  float64 `json:"elevationLoss"`
	AverageSpeed   float64 `json:"averageSpeed"`
	AverageHR      float64 `json:"averageHR"`
	MaxHR          float64 `json:"maxHR"`
	AveragePower   float64 `json:"averagePower"`
	MaxPower       float64 `json:"maxPower"`
	Calories       float64 `json:"calories"`
}

// SplitsResponse wraps the activity splits endpoint response. Regular laps
// (auto-lap or manual) arrive in LapDTOs; SplitSummaries is only populated
// for structured workouts.
type SplitsResponse struct {
	ActivityID     int64   `json:"activityId"`
	SplitSummaries []Split `json:"splitSummaries"`
	LapDTOs        []Split `json:"lapDTOs"`
}

// HRZone holds heart rate time-in-zone data for an activity.
type HRZone struct {
	ZoneNumber  int     `json:"zoneNumber"`
	SecsInZone  float64 `json:"secsInZone"`
	ZoneLowBPM  int     `json:"zoneLowBoundary"`
	ZoneHighBPM int     `json:"zoneHighBoundary"`
}

// PowerZone holds power time-in-zone data for an activity.
type PowerZone struct {
	ZoneNumber    int     `json:"zoneNumber"`
	SecsInZone    float64 `json:"secsInZone"`
	ZoneLowWatts  int     `json:"zoneLowBoundary"`
	ZoneHighWatts int     `json:"zoneHighBoundary"`
}

// Activities returns the most recent limit activities.
// Skips the first start activities.
func (c *Client) Activities(ctx context.Context, start, limit int) ([]Activity, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out []Activity
	if err := c.get(ctx, "/activitylist-service/activities/search/activities", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivitiesByDate returns activities between start and end dates, optionally
// filtered by activity type key (e.g. "running", "cycling"). Pass an empty
// string for activityType to return all types.
func (c *Client) ActivitiesByDate(ctx context.Context, start, end time.Time, activityType string) ([]Activity, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
		"start":     {"0"},
		"limit":     {"999"},
	}
	if activityType != "" {
		params.Set("activityType", activityType)
	}
	var out []Activity
	if err := c.get(ctx, "/activitylist-service/activities/search/activities", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LastActivity returns the single most recent activity.
func (c *Client) LastActivity(ctx context.Context) (*Activity, error) {
	acts, err := c.Activities(ctx, 0, 1)
	if err != nil {
		return nil, err
	}
	if len(acts) == 0 {
		return nil, ErrNoData
	}
	return &acts[0], nil
}

// ActivityDetail returns the full detail for a single activity. The response
// structure varies by activity type, so the raw JSON is returned as a map.
func (c *Client) ActivityDetail(ctx context.Context, id int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/activity-service/activity/%d", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityDetails returns the full detail object for a single activity,
// including metric descriptors and per-sample metrics. The structure varies
// by activity type, so the raw JSON is returned as a map.
func (c *Client) ActivityDetails(ctx context.Context, id int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/activity-service/activity/%d/details", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityTypes returns the catalog of supported activity types
// (running, cycling, kayaking_v2, ...).
func (c *Client) ActivityTypes(ctx context.Context) ([]ActivityType, error) {
	var out []ActivityType
	if err := c.get(ctx, "/activity-service/activity/activityTypes", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivitiesForDailySummary returns the activities that make up the daily
// summary for the given date.
func (c *Client) ActivitiesForDailySummary(ctx context.Context, d time.Time) ([]Activity, error) {
	params := url.Values{"calendarDate": {date(d)}}
	var out []Activity
	if err := c.get(ctx, fmt.Sprintf("/activitylist-service/activities/fordailysummary/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityExerciseSets returns strength training exercise sets for an activity.
func (c *Client) ActivityExerciseSets(ctx context.Context, id int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/activity-service/activity/%d/exerciseSets", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityWeather returns the weather recorded during an activity.
func (c *Client) ActivityWeather(ctx context.Context, id int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/activity-service/activity/%d/weather", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PersonalRecords returns personal records for the authenticated user.
func (c *Client) PersonalRecords(ctx context.Context) ([]PersonalRecord, error) {
	var out []PersonalRecord
	if err := c.get(ctx, fmt.Sprintf("/personalrecord-service/personalrecord/prs/%s", c.displayName), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityCount returns the total number of activities for the authenticated user.
func (c *Client) ActivityCount(ctx context.Context) (int, error) {
	var out struct {
		Count int `json:"totalCount"`
	}
	if err := c.get(ctx, "/activitylist-service/activities/count", nil, &out); err != nil {
		return 0, err
	}
	return out.Count, nil
}

// ActivitySplits returns lap/split summaries for the given activity.
func (c *Client) ActivitySplits(ctx context.Context, id int64) (*SplitsResponse, error) {
	var out SplitsResponse
	if err := c.get(ctx, fmt.Sprintf("/activity-service/activity/%d/splits", id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ActivityTypedSplits returns typed split data (varies by sport type).
func (c *Client) ActivityTypedSplits(ctx context.Context, id int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/activity-service/activity/%d/typedsplits", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivitySplitSummaries returns split summary statistics for the given activity.
func (c *Client) ActivitySplitSummaries(ctx context.Context, id int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/activity-service/activity/%d/split_summaries", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityHRZones returns time spent in each heart rate zone for the given activity.
func (c *Client) ActivityHRZones(ctx context.Context, id int64) ([]HRZone, error) {
	var out []HRZone
	if err := c.get(ctx, fmt.Sprintf("/activity-service/activity/%d/hrTimeInZones", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityPowerZones returns time spent in each power zone for the given activity.
func (c *Client) ActivityPowerZones(ctx context.Context, id int64) ([]PowerZone, error) {
	var out []PowerZone
	if err := c.get(ctx, fmt.Sprintf("/activity-service/activity/%d/powerTimeInZones", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetActivityName renames the given activity.
func (c *Client) SetActivityName(ctx context.Context, id int64, name string) error {
	return c.put(ctx, fmt.Sprintf("/activity-service/activity/%d", id), map[string]any{
		"activityId":   id,
		"activityName": name,
	}, nil)
}

// SetActivityType changes the sport type of the given activity.
func (c *Client) SetActivityType(ctx context.Context, id, typeID, parentTypeID int64, typeKey string) error {
	return c.put(ctx, fmt.Sprintf("/activity-service/activity/%d", id), map[string]any{
		"activityId": id,
		"activityType": map[string]any{
			"typeId":       typeID,
			"typeKey":      typeKey,
			"parentTypeId": parentTypeID,
		},
	}, nil)
}

// DeleteActivity permanently deletes the given activity.
func (c *Client) DeleteActivity(ctx context.Context, id int64) error {
	return c.del(ctx, fmt.Sprintf("/activity-service/activity/%d", id))
}

// DownloadFormat selects the file format for activity downloads.
type DownloadFormat string

const (
	FormatOriginal DownloadFormat = "original" // FIT (device native)
	FormatTCX      DownloadFormat = "tcx"
	FormatGPX      DownloadFormat = "gpx"
	FormatKML      DownloadFormat = "kml"
	FormatCSV      DownloadFormat = "csv"
)

// DownloadActivity returns the raw bytes of an activity file. FormatOriginal
// returns the device-native FIT file; other formats are server-converted.
func (c *Client) DownloadActivity(ctx context.Context, id int64, format DownloadFormat) ([]byte, error) {
	var path string
	if format == FormatOriginal {
		path = fmt.Sprintf("/download-service/files/activity/%d", id)
	} else {
		path = fmt.Sprintf("/download-service/export/%s/activity/%d", format, id)
	}
	return c.getBytes(ctx, path, nil)
}

// UploadActivity uploads a FIT, GPX, or TCX file and returns the server response.
// The filename extension determines the format (.fit, .gpx, .tcx).
func (c *Client) UploadActivity(ctx context.Context, data []byte, filename string) (map[string]json.RawMessage, error) {
	ext := filepath.Ext(filename)
	if ext == "" {
		return nil, fmt.Errorf("filename must include an extension (.fit, .gpx, or .tcx)")
	}
	var out map[string]json.RawMessage
	if err := c.upload(ctx, "/upload-service/upload"+ext, data, filename, &out); err != nil {
		return nil, err
	}
	return out, nil
}
