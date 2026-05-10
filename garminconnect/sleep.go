package garminconnect

import (
	"fmt"
	"net/url"
	"time"
)

// SleepMovement is a single movement reading during sleep.
type SleepMovement struct {
	StartGMT    string  `json:"startGMT"`
	EndGMT      string  `json:"endGMT"`
	ActivityLevel float64 `json:"activityLevel"`
}

// SleepLevel is a sleep stage interval.
type SleepLevel struct {
	StartGMT   string `json:"startGMT"`
	EndGMT     string `json:"endGMT"`
	ActivityLevel float64 `json:"activityLevel"`
}

// DailySleepDTO holds the top-level nightly sleep statistics.
type DailySleepDTO struct {
	ID                         int64   `json:"id"`
	UserProfilePK              int     `json:"userProfilePK"`
	CalendarDate               string  `json:"calendarDate"`
	SleepTimeSeconds           int     `json:"sleepTimeSeconds"`
	NapTimeSeconds             int     `json:"napTimeSeconds"`
	UnmeasurableSleepSeconds   int     `json:"unmeasurableSleepSeconds"`
	DeepSleepSeconds           int     `json:"deepSleepSeconds"`
	LightSleepSeconds          int     `json:"lightSleepSeconds"`
	REMSleepSeconds            int     `json:"remSleepSeconds"`
	AwakeSeconds               int     `json:"awakeSleepSeconds"`
	AverageRespirationValue    float64 `json:"averageRespirationValue"`
	LowestRespirationValue     float64 `json:"lowestRespirationValue"`
	HighestRespirationValue    float64 `json:"highestRespirationValue"`
	AvgSleepStress             float64 `json:"avgSleepStress"`
	SpO2AvgReadingPercent      float64 `json:"spO2AvgReadingPercent"`
	SpO2LowReadingPercent      float64 `json:"spO2LowReadingPercent"`
	SleepScoreFeedback         string  `json:"sleepScoreFeedback"`
	SleepScoreInsight          string  `json:"sleepScoreInsight"`
}

// SleepData is the full sleep response including stages and movement.
type SleepData struct {
	DailySleepDTO        DailySleepDTO   `json:"dailySleepDTO"`
	SleepMovement        []SleepMovement `json:"sleepMovement"`
	SleepLevels          []SleepLevel    `json:"sleepLevels"`
	RestlessMomentsCount int             `json:"restlessMomentsCount"`
}

// HRVReading is a single HRV measurement.
type HRVReading struct {
	HRVValue  int    `json:"hrvValue"`
	StartGMT  string `json:"startGMT"`
}

// HRVSummary holds nightly HRV statistics.
type HRVSummary struct {
	UserProfilePK   int    `json:"userProfilePK"`
	WeeklyAvg       int    `json:"weeklyAvg"`
	LastNight        int    `json:"lastNight"`
	LastNight5MinHigh int  `json:"lastNight5MinHigh"`
	Baseline        struct {
		LowUpper  int `json:"lowUpper"`
		BalancedLow int `json:"balancedLow"`
		BalancedUpper int `json:"balancedUpper"`
		MarkerValue string `json:"markerValue"`
	} `json:"baseline"`
	Status          string `json:"status"`
	FeedbackPhrase  string `json:"feedbackPhrase"`
	CalendarDate    string `json:"calendarDate"`
	StartTimestampGMT string `json:"startTimestampGMT"`
	EndTimestampGMT   string `json:"endTimestampGMT"`
}

// HRVData is the full HRV response.
type HRVData struct {
	HRVSummary  HRVSummary   `json:"hrvSummary"`
	HRVReadings []HRVReading `json:"hrvReadings"`
	StartTimestampGMT string  `json:"startTimestampGMT"`
	EndTimestampGMT   string  `json:"endTimestampGMT"`
}

// SleepData returns detailed sleep data for the given date.
// The date should be the morning date (day you woke up).
func (c *Client) SleepData(d time.Time) (*SleepData, error) {
	params := url.Values{"date": {date(d)}}
	var out SleepData
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/dailySleepData/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// HRVData returns heart rate variability data for the given date.
func (c *Client) HRVData(d time.Time) (*HRVData, error) {
	var out HRVData
	if err := c.get(fmt.Sprintf("/hrv-service/hrv/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
