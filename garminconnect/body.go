package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// BodyComposition holds body composition metrics over a date range.
type BodyComposition struct {
	StartDate      string                  `json:"startDate"`
	EndDate        string                  `json:"endDate"`
	TotalAverage   BodyCompositionAverage  `json:"totalAverage"`
	DateWeightList []BodyCompositionSample `json:"dateWeightList"`
}

// BodyCompositionAverage is the averaged body composition over the date range.
type BodyCompositionAverage struct {
	From            int64   `json:"from"`  // unix ms
	Until           int64   `json:"until"` // unix ms
	WeightTimestamp string  `json:"weightTimestamp"`
	Weight          float64 `json:"weight"`
	Bmi             float64 `json:"bmi"`
	BodyFat         float64 `json:"bodyFat"`
	BodyWater       float64 `json:"bodyWater"`
	BoneMass        float64 `json:"boneMass"`
	MuscleMass      float64 `json:"muscleMass"`
	VisceralFat     float64 `json:"visceralFat"`
	MetabolicAge    float64 `json:"metabolicAge"`
	PhysiqueRating  float64 `json:"physiqueRating"`
}

// BodyCompositionSample is a single day's body composition measurement.
type BodyCompositionSample struct {
	CalendarDate string  `json:"calendarDate"`
	Weight       float64 `json:"weight"`
	Bmi          float64 `json:"bmi"`
	BodyFat      float64 `json:"bodyFat"`
	BodyWater    float64 `json:"bodyWater"`
	BoneMass     float64 `json:"boneMass"`
	MuscleMass   float64 `json:"muscleMass"`
	VisceralFat  float64 `json:"visceralFat"`
	MetabolicAge float64 `json:"metabolicAge"`
}

// WeighIn represents a single weigh-in measurement.
type WeighIn struct {
	SamplePK     int64   `json:"samplePk"`
	Date         int64   `json:"date"` // epoch milliseconds
	CalendarDate string  `json:"calendarDate"`
	Weight       float64 `json:"weight"` // grams
	Bmi          float64 `json:"bmi"`
	BodyFat      float64 `json:"bodyFat"`
	BodyWater    float64 `json:"bodyWater"`
	BoneMass     float64 `json:"boneMass"`
	MuscleMass   float64 `json:"muscleMass"`
	VisceralFat  float64 `json:"visceralFat"`
	MetabolicAge float64 `json:"metabolicAge"`
	SourceType   string  `json:"sourceType"`
	TimestampGMT int64   `json:"timestampGMT"` // epoch milliseconds
}

// DailyWeightSummary groups the weigh-ins recorded on one day.
type DailyWeightSummary struct {
	SummaryDate        string    `json:"summaryDate"`
	NumOfWeightEntries int       `json:"numOfWeightEntries"`
	MinWeight          float64   `json:"minWeight"`
	MaxWeight          float64   `json:"maxWeight"`
	LatestWeight       *WeighIn  `json:"latestWeight"`
	AllWeightMetrics   []WeighIn `json:"allWeightMetrics"`
}

// WeighInsResponse wraps the weigh-ins range API response.
type WeighInsResponse struct {
	DailyWeightSummaries []DailyWeightSummary `json:"dailyWeightSummaries"`
}

// DailyWeighInsResponse wraps the single-day weigh-ins (dayview) response,
// which uses a flat list rather than per-day summaries.
type DailyWeighInsResponse struct {
	StartDate      string    `json:"startDate"`
	EndDate        string    `json:"endDate"`
	DateWeightList []WeighIn `json:"dateWeightList"`
}

// BodyComposition returns body composition data between start and end dates.
func (c *Client) BodyComposition(ctx context.Context, start, end time.Time) (*BodyComposition, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out BodyComposition
	if err := c.get(ctx, "/weight-service/weight/dateRange", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// WeighIns returns all weigh-in measurements between start and end dates.
func (c *Client) WeighIns(ctx context.Context, start, end time.Time) (*WeighInsResponse, error) {
	params := url.Values{"includeAll": {"true"}}
	var out WeighInsResponse
	if err := c.get(ctx, fmt.Sprintf("/weight-service/weight/range/%s/%s", date(start), date(end)), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DailyWeighIns returns weigh-in measurements for a specific date.
func (c *Client) DailyWeighIns(ctx context.Context, d time.Time) (*DailyWeighInsResponse, error) {
	params := url.Values{"calendarDate": {date(d)}}
	var out DailyWeighInsResponse
	if err := c.get(ctx, fmt.Sprintf("/weight-service/weight/dayview/%s", date(d)), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LatestWeightEntry is the most recent weigh-in with optional body-composition
// metrics (the composition fields are null when entered manually).
type LatestWeightEntry struct {
	Date           int64   `json:"date"`         // unix ms
	TimestampGMT   int64   `json:"timestampGMT"` // unix ms
	Version        int64   `json:"version"`
	Weight         float64 `json:"weight"`      // grams
	WeightDelta    float64 `json:"weightDelta"` // grams
	Bmi            float64 `json:"bmi"`
	BodyFat        float64 `json:"bodyFat"`
	BodyWater      float64 `json:"bodyWater"`
	BoneMass       float64 `json:"boneMass"`
	MuscleMass     float64 `json:"muscleMass"`
	VisceralFat    float64 `json:"visceralFat"`
	MetabolicAge   float64 `json:"metabolicAge"`
	PhysiqueRating float64 `json:"physiqueRating"`
	CaloricIntake  float64 `json:"caloricIntake"`
	SourceType     string  `json:"sourceType"`
}

// LatestWeight returns the most recent weigh-in measurement as of the given date.
func (c *Client) LatestWeight(ctx context.Context, d time.Time) (*LatestWeightEntry, error) {
	params := url.Values{"date": {date(d)}}
	var out LatestWeightEntry
	if err := c.get(ctx, "/weight-service/weight/latest", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddWeighIn records a new weigh-in. weightKg is in kilograms; timestamp is RFC3339.
// The server stores weight in grams, so weightKg is converted automatically.
func (c *Client) AddWeighIn(ctx context.Context, weightKg float64, unitKey, timestamp string) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	body := map[string]any{
		"unitKey":       unitKey, // e.g. "kg"
		"value":         weightKg,
		"dateTimestamp": timestamp,
		"gmtTimestamp":  timestamp,
	}
	if err := c.post(ctx, "/weight-service/user-weight", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteWeighIns removes every weigh-in on the given day and returns how many
// were deleted.
func (c *Client) DeleteWeighIns(ctx context.Context, d time.Time) (int, error) {
	day, err := c.DailyWeighIns(ctx, d)
	if err != nil {
		return 0, err
	}
	for i, w := range day.DateWeightList {
		if err := c.DeleteWeighIn(ctx, date(d), w.SamplePK); err != nil {
			return i, err
		}
	}
	return len(day.DateWeightList), nil
}

// DeleteWeighIn removes a weigh-in by its primary key and calendar date.
func (c *Client) DeleteWeighIn(ctx context.Context, cdate string, weightPK int64) error {
	return c.del(ctx, fmt.Sprintf("/weight-service/weight/%s/byversion/%d", cdate, weightPK))
}

// FirstWeight returns the user's earliest weigh-in.
func (c *Client) FirstWeight(ctx context.Context) (*LatestWeightEntry, error) {
	var out LatestWeightEntry
	if err := c.get(ctx, "/weight-service/weight/first", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DailyLastWeights returns the last weigh-in of each day between start and
// end.
func (c *Client) DailyLastWeights(ctx context.Context, start, end time.Time) ([]LatestWeightEntry, error) {
	params := url.Values{"startDate": {date(start)}, "endDate": {date(end)}}
	var out []LatestWeightEntry
	if err := c.get(ctx, "/weight-service/weight/daily/last", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeightAverage is an average weight over a period.
type WeightAverage struct {
	From   int64   `json:"from"`  // unix ms
	Until  int64   `json:"until"` // unix ms
	Weight float64 `json:"weight"`
}

// WeeklyWeight is a year of weekly weight summaries.
type WeeklyWeight struct {
	Weeks []struct {
		StartDate     string        `json:"startDate"`
		EndDate       string        `json:"endDate"`
		MinWeight     float64       `json:"minWeight"` // grams
		MaxWeight     float64       `json:"maxWeight"`
		AverageWeight WeightAverage `json:"averageWeight"`
	} `json:"weeklyWeightSummaries"`
	Average WeightAverage `json:"weightAverageDTO"`
}

// weeklyRangeStart is the start Garmin requires for a weekly range ending at
// end: exactly 52 weeks, counting both ends.
func weeklyRangeStart(end time.Time) time.Time { return end.AddDate(0, 0, -363) }

// WeeklyWeight returns weekly weight summaries for the 52 weeks ending at end.
func (c *Client) WeeklyWeight(ctx context.Context, end time.Time) (*WeeklyWeight, error) {
	var out WeeklyWeight
	path := fmt.Sprintf("/weight-service/weight/weeklyRange/%s/%s", date(weeklyRangeStart(end)), date(end))
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
