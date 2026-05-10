// Copyright Christopher Barnes
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package garminconnect

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// Activity represents a single Garmin Connect activity summary.
type Activity struct {
	ActivityID   int64  `json:"activityId"`
	ActivityName string `json:"activityName"`
	ActivityType struct {
		TypeKey string `json:"typeKey"`
	} `json:"activityType"`
	StartTimeGMT   string  `json:"startTimeGMT"`
	StartTimeLocal string  `json:"startTimeLocal"`
	Duration       float64 `json:"duration"`        // seconds
	ElapsedDuration float64 `json:"elapsedDuration"` // seconds
	MovingDuration float64 `json:"movingDuration"`  // seconds
	Distance       float64 `json:"distance"`        // meters
	Calories       float64 `json:"calories"`
	AverageHR      float64 `json:"averageHR"`
	MaxHR          float64 `json:"maxHR"`
	AverageSpeed   float64 `json:"averageSpeed"` // meters/second
	MaxSpeed       float64 `json:"maxSpeed"`     // meters/second
	ElevationGain  float64 `json:"elevationGain"`
	ElevationLoss  float64 `json:"elevationLoss"`
	Steps          int64   `json:"steps"`
	TrainingEffect float64 `json:"trainingEffect"`
	AnaerobicTrainingEffect float64 `json:"anaerobicTrainingEffect"`
	AerobicTrainingEffectMessage  string `json:"aerobicTrainingEffectMessage"`
	AverageRunningCadenceInStepsPerMinute float64 `json:"averageRunningCadenceInStepsPerMinute"`
	VO2MaxValue    float64 `json:"vO2MaxValue"`
	LocationName   string  `json:"locationName"`
	OwnerId        int64   `json:"ownerId"`
	HasPolyline    bool    `json:"hasPolyline"`
}

// PersonalRecord represents a personal best for a given activity type.
type PersonalRecord struct {
	ID           int64   `json:"id"`
	TypeID       int64   `json:"typeId"`
	ActivityID   int64   `json:"activityId"`
	StartTimeGMT string  `json:"startTimeGMT"`
	Value        float64 `json:"value"`
	PrTypeLabelKey string `json:"prTypeLabelKey"`
}

// Activities returns the most recent limit activities.
func (c *Client) Activities(limit int) ([]Activity, error) {
	params := url.Values{
		"start": {"0"},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out []Activity
	if err := c.get("/activitylist-service/activities/search/activities", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivitiesByDate returns activities between start and end dates, optionally
// filtered by activity type key (e.g. "running", "cycling"). Pass an empty
// string for activityType to return all types.
func (c *Client) ActivitiesByDate(start, end time.Time, activityType string) ([]Activity, error) {
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
	if err := c.get("/activitylist-service/activities/search/activities", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LastActivity returns the single most recent activity.
func (c *Client) LastActivity() (*Activity, error) {
	acts, err := c.Activities(1)
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
func (c *Client) ActivityDetail(id int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/activity-service/activity/%d", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityExerciseSets returns strength training exercise sets for an activity.
func (c *Client) ActivityExerciseSets(id int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/activity-service/activity/%d/exerciseSets", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityWeather returns the weather recorded during an activity.
func (c *Client) ActivityWeather(id int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/activity-service/activity/%d/weather", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PersonalRecords returns personal records for the authenticated user.
func (c *Client) PersonalRecords() ([]PersonalRecord, error) {
	var out []PersonalRecord
	if err := c.get(fmt.Sprintf("/personalrecord-service/personalrecord/prs/%s", c.displayName), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
