package garminconnect

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// BodyComposition holds body composition metrics over a date range.
type BodyComposition struct {
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	TotalAverage struct {
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
	} `json:"totalAverage"`
	DateWeightList []struct {
		CalendarDate string  `json:"calendarDate"`
		Weight       float64 `json:"weight"`
		Bmi          float64 `json:"bmi"`
		BodyFat      float64 `json:"bodyFat"`
		BodyWater    float64 `json:"bodyWater"`
		BoneMass     float64 `json:"boneMass"`
		MuscleMass   float64 `json:"muscleMass"`
		VisceralFat  float64 `json:"visceralFat"`
		MetabolicAge float64 `json:"metabolicAge"`
	} `json:"dateWeightList"`
}

// WeighIn represents a single weigh-in measurement.
type WeighIn struct {
	SamplePK     int64   `json:"samplePk"`
	Date         string  `json:"date"`
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
}

// WeighInsResponse wraps the weigh-ins API response.
type WeighInsResponse struct {
	StartDate      string    `json:"startDate"`
	EndDate        string    `json:"endDate"`
	DateWeightList []WeighIn `json:"dateWeightList"`
	TotalCount     int       `json:"totalCount"`
}

// BodyComposition returns body composition data between start and end dates.
func (c *Client) BodyComposition(start, end time.Time) (*BodyComposition, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out BodyComposition
	if err := c.get("/weight-service/weight/dateRange", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// WeighIns returns all weigh-in measurements between start and end dates.
func (c *Client) WeighIns(start, end time.Time) (*WeighInsResponse, error) {
	params := url.Values{"includeAll": {"true"}}
	var out WeighInsResponse
	if err := c.get(fmt.Sprintf("/weight-service/weight/range/%s/%s", date(start), date(end)), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DailyWeighIns returns weigh-in measurements for a specific date.
func (c *Client) DailyWeighIns(d time.Time) (*WeighInsResponse, error) {
	params := url.Values{"calendarDate": {date(d)}}
	var out WeighInsResponse
	if err := c.get(fmt.Sprintf("/weight-service/weight/dayview/%s", date(d)), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddWeighIn records a new weigh-in. weightKg is in kilograms; timestamp is RFC3339.
// The server stores weight in grams, so weightKg is converted automatically.
func (c *Client) AddWeighIn(weightKg float64, unitKey, timestamp string) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	body := map[string]any{
		"unitKey":       unitKey, // e.g. "kg"
		"value":         weightKg,
		"dateTimestamp": timestamp,
		"gmtTimestamp":  timestamp,
	}
	if err := c.post("/weight-service/user-weight", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteWeighIn removes a weigh-in by its primary key and calendar date.
func (c *Client) DeleteWeighIn(cdate string, weightPK int64) error {
	return c.del(fmt.Sprintf("/weight-service/weight/%s/byversion/%d", cdate, weightPK))
}
