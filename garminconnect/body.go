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
	"fmt"
	"net/url"
	"time"
)

// BodyComposition holds body composition metrics over a date range.
type BodyComposition struct {
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	TotalAverage struct {
		From                  string  `json:"from"`
		Until                 string  `json:"until"`
		WeightTimestamp       string  `json:"weightTimestamp"`
		Weight                float64 `json:"weight"`
		Bmi                   float64 `json:"bmi"`
		BodyFat               float64 `json:"bodyFat"`
		BodyWater             float64 `json:"bodyWater"`
		BoneMass              float64 `json:"boneMass"`
		MuscleMass            float64 `json:"muscleMass"`
		VisceralFat           float64 `json:"visceralFat"`
		MetabolicAge          float64 `json:"metabolicAge"`
		PhysiqueRating        float64 `json:"physiqueRating"`
	} `json:"totalAverage"`
	DateWeightList []struct {
		CalendarDate          string  `json:"calendarDate"`
		Weight                float64 `json:"weight"`
		Bmi                   float64 `json:"bmi"`
		BodyFat               float64 `json:"bodyFat"`
		BodyWater             float64 `json:"bodyWater"`
		BoneMass              float64 `json:"boneMass"`
		MuscleMass            float64 `json:"muscleMass"`
		VisceralFat           float64 `json:"visceralFat"`
		MetabolicAge          float64 `json:"metabolicAge"`
	} `json:"dateWeightList"`
}

// WeighIn represents a single weigh-in measurement.
type WeighIn struct {
	SamplePK            int64   `json:"samplePk"`
	Date                string  `json:"date"`
	CalendarDate        string  `json:"calendarDate"`
	Weight              float64 `json:"weight"` // grams
	Bmi                 float64 `json:"bmi"`
	BodyFat             float64 `json:"bodyFat"`
	BodyWater           float64 `json:"bodyWater"`
	BoneMass            float64 `json:"boneMass"`
	MuscleMass          float64 `json:"muscleMass"`
	VisceralFat         float64 `json:"visceralFat"`
	MetabolicAge        float64 `json:"metabolicAge"`
	SourceType          string  `json:"sourceType"`
}

// WeighInsResponse wraps the weigh-ins API response.
type WeighInsResponse struct {
	StartDate           string    `json:"startDate"`
	EndDate             string    `json:"endDate"`
	DateWeightList      []WeighIn `json:"dateWeightList"`
	TotalCount          int       `json:"totalCount"`
}

// BodyComposition returns body composition data between start and end dates.
func (c *Client) BodyComposition(start, end time.Time) (*BodyComposition, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out BodyComposition
	if err := c.get(fmt.Sprintf("/weight-service/user-summary/period/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// WeighIns returns all weigh-in measurements between start and end dates.
func (c *Client) WeighIns(start, end time.Time) (*WeighInsResponse, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out WeighInsResponse
	if err := c.get(fmt.Sprintf("/weight-service/weight/range/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DailyWeighIns returns weigh-in measurements for a specific date.
func (c *Client) DailyWeighIns(d time.Time) (*WeighInsResponse, error) {
	params := url.Values{"calendarDate": {date(d)}}
	var out WeighInsResponse
	if err := c.get(fmt.Sprintf("/weight-service/weight/dayview/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
