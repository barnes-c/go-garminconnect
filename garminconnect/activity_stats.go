package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// ActivityLaps returns the laps of an activity.
func (c *Client) ActivityLaps(ctx context.Context, id int64) ([]Split, error) {
	var out struct {
		Laps []Split `json:"lapDTOs"`
	}
	if err := c.get(ctx, fmt.Sprintf("/activity-service/activity/%d/laps", id), nil, &out); err != nil {
		return nil, err
	}
	return out.Laps, nil
}

// Polyline is an activity's GPS track in Google's encoded polyline format.
type Polyline struct {
	NumberOfPoints int     `json:"numberOfPoints"`
	EncodedSamples string  `json:"encodedSamples"`
	MinLat         float64 `json:"minLat"`
	MaxLat         float64 `json:"maxLat"`
	MinLon         float64 `json:"minLon"`
	MaxLon         float64 `json:"maxLon"`
	StartLat       float64 `json:"startLat"`
	StartLon       float64 `json:"startLon"`
	EndLat         float64 `json:"endLat"`
	EndLon         float64 `json:"endLon"`
}

// ActivityPolyline returns the GPS track of an activity.
func (c *Client) ActivityPolyline(ctx context.Context, id int64) (*Polyline, error) {
	var out struct {
		Polyline Polyline `json:"gPolyline"`
	}
	if err := c.get(ctx, fmt.Sprintf("/activity-service/activity/polyline/%d", id), nil, &out); err != nil {
		return nil, err
	}
	return &out.Polyline, nil
}

// ActivityEventType classifies why an activity happened (race, training, ...).
type ActivityEventType struct {
	TypeID    int    `json:"typeId"`
	TypeKey   string `json:"typeKey"` // e.g. "race", "training"
	SortOrder int    `json:"sortOrder"`
}

// ActivityEventTypes returns every activity event type.
func (c *Client) ActivityEventTypes(ctx context.Context) ([]ActivityEventType, error) {
	var out []ActivityEventType
	if err := c.get(ctx, "/activity-service/activity/eventTypes", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// MetricStats summarizes one metric over a set of activities.
type MetricStats struct {
	Count int     `json:"count"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Sum   float64 `json:"sum"`
}

// ActivityProgress is the activity totals for a period.
type ActivityProgress struct {
	Date              string `json:"date"`
	CountOfActivities int    `json:"countOfActivities"`
	// Stats is keyed by group ("all") and then by metric (e.g. "distance").
	Stats map[string]map[string]MetricStats `json:"stats"`
}

// ActivityProgress returns lifetime totals of the given metrics (e.g.
// "distance", "duration"; see AvailableActivityMetrics) for activities between
// start and end.
func (c *Client) ActivityProgress(ctx context.Context, start, end time.Time, metrics ...string) ([]ActivityProgress, error) {
	params := url.Values{
		"startDate":                 {date(start)},
		"endDate":                   {date(end)},
		"aggregation":               {"lifetime"},
		"groupByParentActivityType": {"false"},
		"metric":                    metrics,
	}
	var out []ActivityProgress
	if err := c.get(ctx, "/fitnessstats-service/activity", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AvailableActivityMetrics returns the metric names usable with
// ActivityProgress and ActivityMetrics, grouped by category.
func (c *Client) AvailableActivityMetrics(ctx context.Context) (map[string][]string, error) {
	var out map[string][]string
	if err := c.get(ctx, "/fitnessstats-service/activity/availableMetrics", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityMetricValues holds the requested metrics of one activity.
type ActivityMetricValues struct {
	ActivityID int64
	StartLocal string
	Values     map[string]float64 // keyed by metric name
}

func (a *ActivityMetricValues) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	a.Values = map[string]float64{}
	for k, v := range raw {
		var err error
		switch k {
		case "activityId":
			err = json.Unmarshal(v, &a.ActivityID)
		case "startLocal":
			err = json.Unmarshal(v, &a.StartLocal)
		default:
			var f *float64
			if err = json.Unmarshal(v, &f); err == nil && f != nil {
				a.Values[k] = *f
			}
		}
		if err != nil {
			return fmt.Errorf("%s: %w", k, err)
		}
	}
	return nil
}

// ActivityMetrics returns the given metrics for each activity between start
// and end.
func (c *Client) ActivityMetrics(ctx context.Context, start, end time.Time, metrics ...string) ([]ActivityMetricValues, error) {
	params := url.Values{"startDate": {date(start)}, "endDate": {date(end)}, "metric": metrics}
	var out []ActivityMetricValues
	if err := c.get(ctx, "/fitnessstats-service/activity/all", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}
