package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"slices"
	"strings"
)

// GolfScorecard summarises a round of golf.
type GolfScorecard struct {
	ScorecardID string `json:"scorecardPk"`
	CourseName  string `json:"courseName"`
	EventDate   string `json:"eventDate"`
	TotalScore  int    `json:"totalScore"`
	ToPar       int    `json:"differentialToPar"`
}

// GolfSummary returns a paginated list of scorecard summaries.
// The API returns a pagination wrapper whose scorecard field name is not
// stable across accounts; this function unwraps it by locating the array
// value (scanning keys in sorted order for determinism) and returns the
// scorecard slice.
func (c *Client) GolfSummary(ctx context.Context, start, limit int) ([]GolfScorecard, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var raw map[string]json.RawMessage
	if err := c.get(ctx, "/gcs-golfcommunity/api/v2/scorecard/summary", params, &raw); err != nil {
		return nil, err
	}
	var firstErr error
	for _, k := range slices.Sorted(maps.Keys(raw)) {
		v := raw[k]
		if len(v) == 0 || v[0] != '[' {
			continue
		}
		var out []GolfScorecard
		if err := json.Unmarshal(v, &out); err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("decode scorecard summary %q: %w", k, err)
			}
			continue
		}
		return out, nil
	}
	return nil, firstErr
}

// GolfScorecard returns full details for a single scorecard.
func (c *Client) GolfScorecard(ctx context.Context, scorecardID string) (map[string]json.RawMessage, error) {
	params := url.Values{"uuid": {scorecardID}}
	var out map[string]json.RawMessage
	if err := c.get(ctx, "/gcs-golfcommunity/api/v2/scorecard/detail", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GolfShotData returns shot-level data for the given scorecard and hole numbers.
func (c *Client) GolfShotData(ctx context.Context, scorecardID string, holeNumbers []int) (map[string]json.RawMessage, error) {
	holes := make([]string, len(holeNumbers))
	for i, h := range holeNumbers {
		holes[i] = fmt.Sprintf("%d", h)
	}
	params := url.Values{"holeNumbers": {strings.Join(holes, ",")}}
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/gcs-golfcommunity/api/v2/shot/scorecard/%s/hole", scorecardID), params, &out); err != nil {
		return nil, err
	}
	return out, nil
}
