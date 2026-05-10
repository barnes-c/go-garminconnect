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

// HeartRates holds daily heart rate data.
type HeartRates struct {
	UserProfilePK            int       `json:"userProfilePK"`
	CalendarDate             string    `json:"calendarDate"`
	StartTimestampGMT        string    `json:"startTimestampGMT"`
	EndTimestampGMT          string    `json:"endTimestampGMT"`
	StartTimestampLocal      string    `json:"startTimestampLocal"`
	EndTimestampLocal        string    `json:"endTimestampLocal"`
	RestingHeartRate         int       `json:"restingHeartRate"`
	MinHeartRate             int       `json:"minHeartRate"`
	MaxHeartRate             int       `json:"maxHeartRate"`
	LastSevenDaysAvgRestingHeartRate int `json:"lastSevenDaysAvgRestingHeartRate"`
	HeartRateValueDescriptors []struct {
		Key   string `json:"key"`
		Index int    `json:"index"`
	} `json:"heartRateValueDescriptors"`
	HeartRateValues [][]int64 `json:"heartRateValues"` // [timestamp_ms, bpm]
}

// RestingHeartRateEntry is a single resting heart rate data point.
type RestingHeartRateEntry struct {
	UserProfilePK    int    `json:"userProfilePK"`
	StatisticsStartDate string `json:"statisticsStartDate"`
	CalendarDate     string `json:"calendarDate"`
	Value            int    `json:"value"`
	WeeklyAvg        int    `json:"weeklyAvg"`
}

// RestingHeartRateResponse wraps the resting heart rate API response.
type RestingHeartRateResponse struct {
	AllMetrics struct {
		MetricsMap struct {
			WELLNESS_RESTING_HEART_RATE []RestingHeartRateEntry `json:"WELLNESS_RESTING_HEART_RATE"`
		} `json:"metricsMap"`
	} `json:"allMetrics"`
}

// HeartRates returns heart rate data for the given date.
func (c *Client) HeartRates(d time.Time) (*HeartRates, error) {
	params := url.Values{"date": {date(d)}}
	var out HeartRates
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/dailyHeartRate/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RestingHeartRate returns resting heart rate history between start and end dates.
func (c *Client) RestingHeartRate(start, end time.Time) (*RestingHeartRateResponse, error) {
	params := url.Values{
		"fromDate":  {date(start)},
		"untilDate": {date(end)},
		"metricId":  {"60"},
	}
	var out RestingHeartRateResponse
	if err := c.get(fmt.Sprintf("/userstats-service/wellness/daily/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
