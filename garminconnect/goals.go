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
)

// Goal represents a user-defined fitness goal.
type Goal struct {
	GoalID    int64  `json:"goalId"`
	GoalType  struct {
		TypeKey string `json:"typeKey"`
	} `json:"goalType"`
	GoalValueInMetric float64 `json:"goalValueInMetric"`
	StartDate         string  `json:"startDate"`
	EndDate           string  `json:"endDate"`
	Status            string  `json:"status"`
}

// Badge represents an earned or available achievement badge.
type Badge struct {
	BadgeID              int64  `json:"badgeId"`
	BadgeName            string `json:"badgeName"`
	BadgeKey             string `json:"badgeKey"`
	BadgeCategoryTypeKey string `json:"badgeCategoryTypeKey"`
	EarnedDate           string `json:"earnedDate"`
	BadgePoints          int    `json:"badgePoints"`
}

// Goals returns goals filtered by status (e.g. "active", "completed").
// Pass an empty string to return all goals.
func (c *Client) Goals(status string, start, limit int) ([]Goal, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	if status != "" {
		params.Set("status", status)
	}
	var out []Goal
	if err := c.get("/goal-service/goal/goals", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EarnedBadges returns all badges the user has earned.
func (c *Client) EarnedBadges() ([]Badge, error) {
	var out []Badge
	if err := c.get("/badge-service/badge/earned", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AvailableBadges returns badges the user has not yet earned.
func (c *Client) AvailableBadges() ([]Badge, error) {
	var out []Badge
	if err := c.get("/badge-service/badge/available", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AdHocChallenges returns historical ad-hoc challenge data.
func (c *Client) AdHocChallenges(start, limit int) (map[string]json.RawMessage, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out map[string]json.RawMessage
	if err := c.get("/adhocchallenge-service/adHocChallenge/historical", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BadgeChallenges returns completed badge challenges.
func (c *Client) BadgeChallenges(start, limit int) (map[string]json.RawMessage, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out map[string]json.RawMessage
	if err := c.get("/badgechallenge-service/badgeChallenge/completed", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AvailableBadgeChallenges returns badge challenges the user can join.
func (c *Client) AvailableBadgeChallenges(start, limit int) (map[string]json.RawMessage, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out map[string]json.RawMessage
	if err := c.get("/badgechallenge-service/badgeChallenge/available", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// InProgressVirtualChallenges returns currently active virtual challenges.
func (c *Client) InProgressVirtualChallenges(start, limit int) (map[string]json.RawMessage, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out map[string]json.RawMessage
	if err := c.get("/badgechallenge-service/virtualChallenge/inProgress", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}
