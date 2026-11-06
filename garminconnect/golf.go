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

// GolfClub is one of the user's golf clubs with its distance statistics.
// Distances use the user's golf distance unit.
type GolfClub struct {
	ID               int64   `json:"id"`
	ClubTypeID       int     `json:"clubTypeId"`
	ShaftLength      float64 `json:"shaftLength"`
	FlexTypeID       string  `json:"flexTypeId"`
	AverageDistance  float64 `json:"averageDistance"`
	AdviceDistance   float64 `json:"adviceDistance"`
	Retired          bool    `json:"retired"`
	Deleted          bool    `json:"deleted"`
	LastModifiedTime string  `json:"lastModifiedTime"`
}

// GolfClubStats returns the user's golf clubs with their distance statistics.
// limit caps the number of clubs returned.
func (c *Client) GolfClubStats(ctx context.Context, limit int) ([]GolfClub, error) {
	params := url.Values{
		"per-page":      {fmt.Sprintf("%d", limit)},
		"include-stats": {"true"},
	}
	var out []GolfClub
	if err := c.get(ctx, "/gcs-golfcommunity/api/v2/club/player", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GolfUserStats returns the user's overall golf statistics, such as handicap
// and strokes gained.
func (c *Client) GolfUserStats(ctx context.Context) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, "/gcs-golfcommunity/api/v2/player/stats", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GolfClubType is one of Garmin's golf club types with its default specs.
// Lengths are in inches, angles in degrees.
type GolfClubType struct {
	Value        int     `json:"value"`
	Name         string  `json:"name"`
	ShaftLength  float64 `json:"shaftLength"`
	LoftAngle    float64 `json:"loftAngle"`
	LieAngle     float64 `json:"lieAngle"`
	DisplayRange string  `json:"displayRange"`
	Valid        bool    `json:"valid"`
}

// GolfClubTypes returns Garmin's catalog of golf club types.
func (c *Client) GolfClubTypes(ctx context.Context) ([]GolfClubType, error) {
	var out []GolfClubType
	if err := c.get(ctx, "/gcs-golfcommunity/api/v2/club/types", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
