package garminconnect

import (
	"context"
	"net/url"
	"time"
)

// HeartRates holds daily heart rate data.
type HeartRates struct {
	UserProfilePK                    int               `json:"userProfilePK"`
	CalendarDate                     string            `json:"calendarDate"`
	StartTimestampGMT                string            `json:"startTimestampGMT"`
	EndTimestampGMT                  string            `json:"endTimestampGMT"`
	StartTimestampLocal              string            `json:"startTimestampLocal"`
	EndTimestampLocal                string            `json:"endTimestampLocal"`
	RestingHeartRate                 int               `json:"restingHeartRate"`
	MinHeartRate                     int               `json:"minHeartRate"`
	MaxHeartRate                     int               `json:"maxHeartRate"`
	LastSevenDaysAvgRestingHeartRate int               `json:"lastSevenDaysAvgRestingHeartRate"`
	HeartRateValueDescriptors        []ValueDescriptor `json:"heartRateValueDescriptors"`
	HeartRateValues                  [][]int64         `json:"heartRateValues"` // [timestamp_ms, bpm]
}

// RestingHeartRateEntry is a single resting heart rate data point.
type RestingHeartRateEntry struct {
	UserProfilePK       int     `json:"userProfilePK"`
	StatisticsStartDate string  `json:"statisticsStartDate"`
	CalendarDate        string  `json:"calendarDate"`
	Value               float64 `json:"value"`
	WeeklyAvg           int     `json:"weeklyAvg"`
}

// RestingHeartRateResponse wraps the resting heart rate API response.
type RestingHeartRateResponse struct {
	AllMetrics RestingHeartRateMetrics `json:"allMetrics"`
}

// RestingHeartRateMetrics holds the metrics map of a resting heart rate response.
type RestingHeartRateMetrics struct {
	MetricsMap RestingHeartRateMetricsMap `json:"metricsMap"`
}

// RestingHeartRateMetricsMap maps metric keys to their resting heart rate entries.
type RestingHeartRateMetricsMap struct {
	WellnessRestingHeartRate []RestingHeartRateEntry `json:"WELLNESS_RESTING_HEART_RATE"`
}

// ValueDescriptor labels a column in a wellness values array by key and index.
type ValueDescriptor struct {
	Key   string `json:"key"`
	Index int    `json:"index"`
}

// HeartRates returns heart rate data for the given date.
func (c *Client) HeartRates(ctx context.Context, d time.Time) (*HeartRates, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	params := url.Values{"date": {date(d)}}
	var out HeartRates
	if err := c.get(ctx, "/wellness-service/wellness/dailyHeartRate/"+name, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RestingHeartRate returns resting heart rate history between start and end dates.
func (c *Client) RestingHeartRate(ctx context.Context, start, end time.Time) (*RestingHeartRateResponse, error) {
	params := url.Values{
		"fromDate":  {date(start)},
		"untilDate": {date(end)},
		"metricId":  {"60"},
	}
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	var out RestingHeartRateResponse
	if err := c.get(ctx, "/userstats-service/wellness/daily/"+name, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
