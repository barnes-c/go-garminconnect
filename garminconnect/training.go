package garminconnect

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// TrainingReadiness holds the training readiness score for a day.
type TrainingReadiness struct {
	UserProfilePK     int     `json:"userProfilePK"`
	CalendarDate      string  `json:"calendarDate"`
	Score             int     `json:"score"`
	ScoreQualifier    string  `json:"scoreQualifier"`
	SleepScore        int     `json:"sleepScore"`
	RecoveryTime      int     `json:"recoveryTime"` // hours
	HRVFactorPercent  float64 `json:"acuteLoad"`
	AcuteLoad         float64 `json:"trainingLoad"`
	SleepHistoryScore int     `json:"sleepHistoryScore"`
	HrvWeeklyAverage  int     `json:"hrvWeeklyAverage"`
	FeedbackPhrase    string  `json:"feedbackPhrase"`
}

// TrainingStatusEntry is a single day of training status data.
type TrainingStatusEntry struct {
	CalendarDate       string  `json:"calendarDate"`
	TrainingStatusType string  `json:"trainingStatusType"`
	TrainingLoadType   string  `json:"trainingLoadType"`
	WorkoutGoal        float64 `json:"workoutGoal"`
	AtpPlanLowLoad     float64 `json:"atpPlanLowLoad"`
	AtpPlanHighLoad    float64 `json:"atpPlanHighLoad"`
}

// MaxMetricsEntry holds a VO2 Max data point.
type MaxMetricsEntry struct {
	Generic *struct {
		CalendarDate   string  `json:"calendarDate"`
		VO2MaxValue    float64 `json:"vo2MaxValue"`
		FitnessAge     int     `json:"fitnessAge"`
		FitnessAgeDesc string  `json:"fitnessAgeDescription"`
	} `json:"generic"`
	Cycling *struct {
		CalendarDate string  `json:"calendarDate"`
		VO2MaxValue  float64 `json:"vo2MaxValue"`
		FitnessAge   int     `json:"fitnessAge"`
	} `json:"cycling"`
}

// EnduranceScoreEntry holds an endurance score data point.
type EnduranceScoreEntry struct {
	CalendarDate string  `json:"calendarDate"`
	Score        float64 `json:"score"`
	Level        string  `json:"level"`
	Contributors []struct {
		ActivityType string  `json:"activityType"`
		Contribution float64 `json:"contribution"`
	} `json:"contributors"`
}

// RacePrediction holds an estimated finish time for a race distance.
type RacePrediction struct {
	RaceDistance        string `json:"raceDistance"`  // e.g. "RACE_5K"
	TimePredicted       int    `json:"timePredicted"` // seconds
	TimeUncertainty     int    `json:"timeUncertainty"`
	PredictionAvailable bool   `json:"predictionAvailable"`
}

// HillScoreEntry holds a hill score data point.
type HillScoreEntry struct {
	CalendarDate string  `json:"calendarDate"`
	HillScore    float64 `json:"hillScore"`
	Level        string  `json:"level"`
}

// TrainingReadiness returns the training readiness score for the given date.
func (c *Client) TrainingReadiness(d time.Time) (*TrainingReadiness, error) {
	var out TrainingReadiness
	if err := c.get(fmt.Sprintf("/trainingreadiness-service/trainingreadiness/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TrainingStatus returns training status metrics for the given date.
func (c *Client) TrainingStatus(d time.Time) ([]TrainingStatusEntry, error) {
	params := url.Values{"calendarDate": {date(d)}}
	var out []TrainingStatusEntry
	if err := c.get("/fitnessstats-service/fitness/statistics/training-status", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// MaxMetrics returns VO2 Max metrics between start and end dates.
func (c *Client) MaxMetrics(start, end time.Time) ([]MaxMetricsEntry, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out []MaxMetricsEntry
	if err := c.get(fmt.Sprintf("/metrics-service/metrics/maxmet/daily/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EnduranceScore returns endurance score data between start and end dates.
func (c *Client) EnduranceScore(start, end time.Time) ([]EnduranceScoreEntry, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out []EnduranceScoreEntry
	if err := c.get("/endurancescore-service/endurancescore/stats", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RacePredictions returns estimated race finish times for the user.
func (c *Client) RacePredictions() ([]RacePrediction, error) {
	var out []RacePrediction
	if err := c.get(fmt.Sprintf("/race-predictor-service/race-predictor/races/%s", c.displayName), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// HillScore returns hill score data between start and end dates.
func (c *Client) HillScore(start, end time.Time) ([]HillScoreEntry, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out []HillScoreEntry
	if err := c.get("/hillscore-service/hillscore", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LactateThreshold returns the latest lactate threshold measurement.
func (c *Client) LactateThreshold() (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.get("/biometric-service/biometric/latestLactateThreshold", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FitnessAge returns fitness age data for the given date.
func (c *Client) FitnessAge(d time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/fitnessage-service/fitnessage/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RunningTolerance returns running tolerance statistics between start and end dates.
func (c *Client) RunningTolerance(start, end time.Time) (map[string]json.RawMessage, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out map[string]json.RawMessage
	if err := c.get("/metrics-service/metrics/runningtolerance/stats", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CyclingFTP returns the latest cycling FTP (functional threshold power) estimate.
func (c *Client) CyclingFTP(start, end time.Time) (map[string]json.RawMessage, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out map[string]json.RawMessage
	if err := c.get("/metrics-service/metrics/cyclingftp/latest", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}
