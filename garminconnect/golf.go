package garminconnect

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// GolfScorecard summarises a round of golf.
type GolfScorecard struct {
	ScorecardID  string `json:"scorecardPk"`
	CourseName   string `json:"courseName"`
	EventDate    string `json:"eventDate"`
	TotalScore   int    `json:"totalScore"`
	ToPar        int    `json:"differentialToPar"`
}

// GolfSummary returns a paginated list of scorecard summaries.
func (c *Client) GolfSummary(start, limit int) ([]GolfScorecard, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out []GolfScorecard
	if err := c.get("/gcs-golfcommunity/api/v2/scorecard/summary", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GolfScorecard returns full details for a single scorecard.
func (c *Client) GolfScorecard(scorecardID string) (map[string]json.RawMessage, error) {
	params := url.Values{"uuid": {scorecardID}}
	var out map[string]json.RawMessage
	if err := c.get("/gcs-golfcommunity/api/v2/scorecard/detail", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GolfShotData returns shot-level data for the given scorecard and hole numbers.
func (c *Client) GolfShotData(scorecardID string, holeNumbers []int) (map[string]json.RawMessage, error) {
	holes := make([]string, len(holeNumbers))
	for i, h := range holeNumbers {
		holes[i] = fmt.Sprintf("%d", h)
	}
	params := url.Values{"holeNumbers": {strings.Join(holes, ",")}}
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/gcs-golfcommunity/api/v2/shot/scorecard/%s/hole", scorecardID), params, &out); err != nil {
		return nil, err
	}
	return out, nil
}
