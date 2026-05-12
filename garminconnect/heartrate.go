package garminconnect

import (
	"fmt"
	"net/url"
	"time"
)

// HeartRates holds daily heart rate data.
type HeartRates struct {
	UserProfilePK                    int    `json:"userProfilePK"`
	CalendarDate                     string `json:"calendarDate"`
	StartTimestampGMT                string `json:"startTimestampGMT"`
	EndTimestampGMT                  string `json:"endTimestampGMT"`
	StartTimestampLocal              string `json:"startTimestampLocal"`
	EndTimestampLocal                string `json:"endTimestampLocal"`
	RestingHeartRate                 int    `json:"restingHeartRate"`
	MinHeartRate                     int    `json:"minHeartRate"`
	MaxHeartRate                     int    `json:"maxHeartRate"`
	LastSevenDaysAvgRestingHeartRate int    `json:"lastSevenDaysAvgRestingHeartRate"`
	HeartRateValueDescriptors        []struct {
		Key   string `json:"key"`
		Index int    `json:"index"`
	} `json:"heartRateValueDescriptors"`
	HeartRateValues [][]int64 `json:"heartRateValues"` // [timestamp_ms, bpm]
}

// RestingHeartRateEntry is a single resting heart rate data point.
type RestingHeartRateEntry struct {
	UserProfilePK       int    `json:"userProfilePK"`
	StatisticsStartDate string `json:"statisticsStartDate"`
	CalendarDate        string `json:"calendarDate"`
	Value               float64 `json:"value"`
	WeeklyAvg           int    `json:"weeklyAvg"`
}

// RestingHeartRateResponse wraps the resting heart rate API response.
type RestingHeartRateResponse struct {
	AllMetrics struct {
		MetricsMap struct {
			WellnessRestingHeartRate []RestingHeartRateEntry `json:"WELLNESS_RESTING_HEART_RATE"`
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
