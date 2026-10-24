package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Goal represents a user-defined fitness goal.
type Goal struct {
	GoalID            int64    `json:"goalId"`
	GoalType          GoalType `json:"goalType"`
	GoalValueInMetric float64  `json:"goalValueInMetric"`
	StartDate         string   `json:"startDate"`
	EndDate           string   `json:"endDate"`
	Status            string   `json:"status"`
}

// GoalType identifies the kind of a goal.
type GoalType struct {
	TypeKey string `json:"typeKey"`
}

// Badge represents an earned or available achievement badge.
type Badge struct {
	BadgeID              int64  `json:"badgeId"`
	BadgeName            string `json:"badgeName"`
	BadgeKey             string `json:"badgeKey"`
	BadgeCategoryTypeKey string `json:"badgeCategoryTypeKey"`
	EarnedDate           string `json:"badgeEarnedDate"`
	BadgePoints          int    `json:"badgePoints"`
	BadgeCategoryID      int    `json:"badgeCategoryId"`   // see BadgeAttributes
	BadgeDifficultyID    int    `json:"badgeDifficultyId"` // see BadgeAttributes
	BadgeTypeIDs         []int  `json:"badgeTypeIds"`      // see BadgeAttributes
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

// BadgeDetail returns a single badge.
func (c *Client) BadgeDetail(ctx context.Context, badgeID int64) (*Badge, error) {
	var out Badge
	if err := c.get(ctx, fmt.Sprintf("/badge-service/badge/detail/v3/%d", badgeID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LatestEarnedBadges returns up to limit of the user's most recently earned
// badges.
func (c *Client) LatestEarnedBadges(ctx context.Context, limit int) ([]Badge, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	var out []Badge
	params := url.Values{"limit": {fmt.Sprint(limit)}}
	if err := c.get(ctx, "/badge-service/badge/earned/latest/"+name, params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DailyBadges returns the badges earned on day d.
func (c *Client) DailyBadges(ctx context.Context, d time.Time) ([]Badge, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	var out []Badge
	if err := c.get(ctx, fmt.Sprintf("/badge-service/badge/%s/earned/dailysummary/%s", name, date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityBadges returns the badges an activity earned.
func (c *Client) ActivityBadges(ctx context.Context, activityID int64) ([]Badge, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	var out []Badge
	if err := c.get(ctx, fmt.Sprintf("/badge-service/badge/%s/earned/activity/%d", name, activityID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BadgeAttribute is one entry of a badge lookup table.
type BadgeAttribute struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

// BadgeAttributes are the lookup tables for the IDs on a Badge.
type BadgeAttributes struct {
	Types        []BadgeAttribute
	Categories   []BadgeAttribute
	Difficulties []BadgeAttribute
	Units        []BadgeAttribute
	AssocTypes   []BadgeAttribute
}

// BadgeAttributes returns the badge type, category, difficulty, unit and
// association lookup tables.
func (c *Client) BadgeAttributes(ctx context.Context) (*BadgeAttributes, error) {
	// Each table names its id and key fields after itself (badgeTypeId,
	// badgeTypeKey, ...), so decode generically and normalize.
	var raw map[string][]map[string]json.RawMessage
	if err := c.get(ctx, "/badge-service/badge/attributes", nil, &raw); err != nil {
		return nil, err
	}
	table := func(name, prefix string) ([]BadgeAttribute, error) {
		var out []BadgeAttribute
		for _, row := range raw[name] {
			var a BadgeAttribute
			if err := json.Unmarshal(row[prefix+"Id"], &a.ID); err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
			if err := json.Unmarshal(row[prefix+"Key"], &a.Key); err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
			out = append(out, a)
		}
		return out, nil
	}
	var out BadgeAttributes
	var err error
	for _, t := range []struct {
		dst          *[]BadgeAttribute
		name, prefix string
	}{
		{&out.Types, "badgeTypes", "badgeType"},
		{&out.Categories, "badgeCategories", "badgeCategory"},
		{&out.Difficulties, "badgeDifficulties", "badgeDifficulty"},
		{&out.Units, "badgeUnits", "badgeUnit"},
		{&out.AssocTypes, "badgeAssocTypes", "badgeAssocType"},
	} {
		if *t.dst, err = table(t.name, t.prefix); err != nil {
			return nil, err
		}
	}
	return &out, nil
}

// VirtualChallenge is a virtual challenge (e.g. a distance goal) the user can
// join.
type VirtualChallenge struct {
	UUID             string  `json:"uuid"`
	Name             string  `json:"badgeChallengeName"`
	BadgeID          int64   `json:"badgeId"`
	BadgeKey         string  `json:"badgeKey"`
	BadgePoints      int     `json:"badgePoints"`
	BadgeUnitID      int     `json:"badgeUnitId"` // see BadgeAttributes
	BadgeTargetValue float64 `json:"badgeTargetValue"`
	UserJoined       bool    `json:"userJoined"`
}

// AvailableVirtualChallenges returns virtual challenges the user can join,
// starting at the 1-based position start.
func (c *Client) AvailableVirtualChallenges(ctx context.Context, start, limit int) ([]VirtualChallenge, error) {
	params := url.Values{"start": {fmt.Sprint(start)}, "limit": {fmt.Sprint(limit)}}
	var out []VirtualChallenge
	if err := c.get(ctx, "/badgechallenge-service/virtualChallenge/available", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PersonalRecordType is a kind of personal record (e.g. fastest 5K).
type PersonalRecordType struct {
	ID       int     `json:"id"`
	Key      string  `json:"key"`
	Sport    string  `json:"sport"`
	Visible  bool    `json:"visible"`
	MinValue float64 `json:"minValue"`
	MaxValue float64 `json:"maxValue"`
}

// PersonalRecordTypes returns every personal record type.
func (c *Client) PersonalRecordTypes(ctx context.Context) ([]PersonalRecordType, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	var out []PersonalRecordType
	if err := c.get(ctx, "/personalrecord-service/personalrecordtype/prtypes/"+name, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityPersonalRecords returns the personal records set by an activity.
func (c *Client) ActivityPersonalRecords(ctx context.Context, activityID int64) ([]PersonalRecord, error) {
	var out []PersonalRecord
	if err := c.get(ctx, fmt.Sprintf("/personalrecord-service/personalrecord/prByActivityId/%d", activityID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DailyGoal is a daily goal, e.g. GoalType "INTENSITY_MINUTES" or
// "FLOORS_CLIMBED".
type DailyGoal struct {
	GoalType  string  `json:"goalType"`
	GoalValue float64 `json:"goalValue"`
}

// EffectiveGoals returns the daily goals in effect on day d.
func (c *Client) EffectiveGoals(ctx context.Context, d time.Time) ([]DailyGoal, error) {
	var out []DailyGoal
	if err := c.get(ctx, "/goal-service/goal/user/effective/"+date(d), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetDailyGoal sets a daily goal from day d onwards.
func (c *Client) SetDailyGoal(ctx context.Context, d time.Time, goal DailyGoal) error {
	slug := goal.GoalType
	if slug != "WEIGHT_GRAMS" {
		slug = strings.ReplaceAll(strings.ToLower(slug), "_", "-")
	}
	return c.put(ctx, fmt.Sprintf("/goal-service/goal/user/effective/%s/%s", slug, date(d)), goal, nil)
}

// StepsGoal is the daily step goal and whether the device uses it.
type StepsGoal struct {
	UserGoal       int    `json:"userGoal"`
	DeviceGoal     int    `json:"deviceGoal"`
	SyncedToDevice bool   `json:"syncedToDevice"`
	GoalType       string `json:"goalType"`
	Effective      string `json:"effective"`
}

// StepsGoal returns the step goal in effect on day d.
func (c *Client) StepsGoal(ctx context.Context, d time.Time) (*StepsGoal, error) {
	var out StepsGoal
	if err := c.get(ctx, "/wellness-service/wellness/wellness-goals/consolidated/steps/"+date(d), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
