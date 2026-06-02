package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Goal represents a user-defined fitness goal.
type Goal struct {
	GoalID   int64 `json:"goalId"`
	GoalType struct {
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
func (c *Client) Goals(ctx context.Context, status string, start, limit int) ([]Goal, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	if status != "" {
		params.Set("status", status)
	}
	var out []Goal
	if err := c.get(ctx, "/goal-service/goal/goals", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EarnedBadges returns all badges the user has earned.
func (c *Client) EarnedBadges(ctx context.Context) ([]Badge, error) {
	var out []Badge
	if err := c.get(ctx, "/badge-service/badge/earned", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AvailableBadges returns badges the user has not yet earned.
func (c *Client) AvailableBadges(ctx context.Context) ([]Badge, error) {
	var out []Badge
	if err := c.get(ctx, "/badge-service/badge/available", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AdHocChallenges returns historical ad-hoc challenge data.
func (c *Client) AdHocChallenges(ctx context.Context, start, limit int) (map[string]json.RawMessage, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out map[string]json.RawMessage
	if err := c.get(ctx, "/adhocchallenge-service/adHocChallenge/historical", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BadgeChallenges returns completed badge challenges.
func (c *Client) BadgeChallenges(ctx context.Context, start, limit int) (map[string]json.RawMessage, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out map[string]json.RawMessage
	if err := c.get(ctx, "/badgechallenge-service/badgeChallenge/completed", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AvailableBadgeChallenges returns badge challenges the user can join.
func (c *Client) AvailableBadgeChallenges(ctx context.Context, start, limit int) (map[string]json.RawMessage, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out map[string]json.RawMessage
	if err := c.get(ctx, "/badgechallenge-service/badgeChallenge/available", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// InProgressVirtualChallenges returns currently active virtual challenges.
func (c *Client) InProgressVirtualChallenges(ctx context.Context, start, limit int) (map[string]json.RawMessage, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out map[string]json.RawMessage
	if err := c.get(ctx, "/badgechallenge-service/virtualChallenge/inProgress", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}
