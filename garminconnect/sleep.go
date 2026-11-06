package garminconnect

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// SleepMovement is a single movement reading during sleep.
type SleepMovement struct {
	StartGMT      string  `json:"startGMT"`
	EndGMT        string  `json:"endGMT"`
	ActivityLevel float64 `json:"activityLevel"`
}

// SleepLevel is a sleep stage interval.
type SleepLevel struct {
	StartGMT      string  `json:"startGMT"`
	EndGMT        string  `json:"endGMT"`
	ActivityLevel float64 `json:"activityLevel"`
}

// DailySleepDTO holds the top-level nightly sleep statistics.
type DailySleepDTO struct {
	ID                       int64   `json:"id"`
	UserProfilePK            int     `json:"userProfilePK"`
	CalendarDate             string  `json:"calendarDate"`
	SleepTimeSeconds         int     `json:"sleepTimeSeconds"`
	NapTimeSeconds           int     `json:"napTimeSeconds"`
	UnmeasurableSleepSeconds int     `json:"unmeasurableSleepSeconds"`
	DeepSleepSeconds         int     `json:"deepSleepSeconds"`
	LightSleepSeconds        int     `json:"lightSleepSeconds"`
	REMSleepSeconds          int     `json:"remSleepSeconds"`
	AwakeSeconds             int     `json:"awakeSleepSeconds"`
	AverageRespirationValue  float64 `json:"averageRespirationValue"`
	LowestRespirationValue   float64 `json:"lowestRespirationValue"`
	HighestRespirationValue  float64 `json:"highestRespirationValue"`
	AvgSleepStress           float64 `json:"avgSleepStress"`
	SpO2AvgReadingPercent    float64 `json:"averageSpO2Value"`
	SleepScoreFeedback       string  `json:"sleepScoreFeedback"`
	SleepScoreInsight        string  `json:"sleepScoreInsight"`
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
	HRVValue         int    `json:"hrvValue"`
	ReadingTimeGMT   string `json:"readingTimeGMT"`
	ReadingTimeLocal string `json:"readingTimeLocal"`
}

// HRVBaseline holds the HRV baseline range used to classify nightly status.
type HRVBaseline struct {
	LowUpper      int     `json:"lowUpper"`
	BalancedLow   int     `json:"balancedLow"`
	BalancedUpper int     `json:"balancedUpper"`
	MarkerValue   float64 `json:"markerValue"`
}

// HRVSummary holds nightly HRV statistics.
type HRVSummary struct {
	WeeklyAvg         int         `json:"weeklyAvg"`
	LastNight         int         `json:"lastNightAvg"`
	LastNight5MinHigh int         `json:"lastNight5MinHigh"`
	Baseline          HRVBaseline `json:"baseline"`
	Status            string      `json:"status"`
	FeedbackPhrase    string      `json:"feedbackPhrase"`
	CalendarDate      string      `json:"calendarDate"`
	CreateTimestamp   string      `json:"createTimeStamp"`
}

// HRVData is the full HRV response.
type HRVData struct {
	HRVSummary        HRVSummary   `json:"hrvSummary"`
	HRVReadings       []HRVReading `json:"hrvReadings"`
	UserProfilePK     int64        `json:"userProfilePk"`
	StartTimestampGMT string       `json:"startTimestampGMT"`
	EndTimestampGMT   string       `json:"endTimestampGMT"`
}

// HRVDataRange returns HRV data for every day between start and end.
// HRVData covers a single day; this queries the same service's range endpoint
// so a span can be fetched in one request.
func (c *Client) HRVDataRange(ctx context.Context, start, end time.Time) ([]HRVSummary, error) {
	var out struct {
		HRVSummaries []HRVSummary `json:"hrvSummaries"`
	}
	if err := c.get(ctx, fmt.Sprintf("/hrv-service/hrv/daily/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return out.HRVSummaries, nil
}

// SleepData returns detailed sleep data for the given date.
// The date should be the morning date (day you woke up).
func (c *Client) SleepData(ctx context.Context, d time.Time) (*SleepData, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	params := url.Values{"date": {date(d)}}
	var out SleepData
	if err := c.get(ctx, "/wellness-service/wellness/dailySleepData/"+name, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DailySleepData returns sleep data for the given date from the sleep-service
// endpoint (distinct from SleepData, which uses the wellness-service path).
// The response shares SleepData's structure (dailySleepDTO, sleepLevels,
// sleepMovement).
func (c *Client) DailySleepData(ctx context.Context, d time.Time) (*SleepData, error) {
	params := url.Values{"date": {date(d)}}
	var out SleepData
	if err := c.get(ctx, "/sleep-service/sleep/dailySleepData", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SleepDailyStats are sleep averages over a period with a per-night series.
type SleepDailyStats struct {
	Overall SleepAverages `json:"overallStats"`
	Days    []SleepDay    `json:"individualStats"`
}

// SleepDay is one night's sleep summary. Durations are in seconds; the
// *InMillis fields are epoch milliseconds.
type SleepDay struct {
	CalendarDate string `json:"calendarDate"`
	Values       struct {
		TotalSleepTimeInSeconds     int      `json:"totalSleepTimeInSeconds"`
		DeepTime                    int      `json:"deepTime"`
		LightTime                   int      `json:"lightTime"`
		RemTime                     int      `json:"remTime"`
		AwakeTime                   int      `json:"awakeTime"`
		SleepNeed                   int      `json:"sleepNeed"`
		SleepScore                  int      `json:"sleepScore"`
		SleepScoreQuality           string   `json:"sleepScoreQuality"`
		RestingHeartRate            int      `json:"restingHeartRate"`
		AvgHeartRate                float64  `json:"avgHeartRate"`
		Respiration                 float64  `json:"respiration"`
		SpO2                        float64  `json:"spO2"`
		BodyBatteryChange           int      `json:"bodyBatteryChange"`
		AvgOvernightHRV             float64  `json:"avgOvernightHrv"`
		HRV7dAverage                float64  `json:"hrv7dAverage"`
		HRVStatus                   string   `json:"hrvStatus"`
		SkinTempC                   *float64 `json:"skinTempC"`
		SkinTempF                   *float64 `json:"skinTempF"`
		LocalSleepStartTimeInMillis int64    `json:"localSleepStartTimeInMillis"`
		LocalSleepEndTimeInMillis   int64    `json:"localSleepEndTimeInMillis"`
		GMTSleepStartTimeInMillis   int64    `json:"gmtSleepStartTimeInMillis"`
		GMTSleepEndTimeInMillis     int64    `json:"gmtSleepEndTimeInMillis"`
	} `json:"values"`
}

// SleepStats returns sleep averages and a per-night series between start and
// end dates. The range must not exceed 28 days.
func (c *Client) SleepStats(ctx context.Context, start, end time.Time) (*SleepDailyStats, error) {
	var out SleepDailyStats
	if err := c.get(ctx, fmt.Sprintf("/sleep-service/stats/sleep/daily/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// HRVData returns heart rate variability data for the given date.
func (c *Client) HRVData(ctx context.Context, d time.Time) (*HRVData, error) {
	var out HRVData
	if err := c.get(ctx, fmt.Sprintf("/hrv-service/hrv/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SleepAverages are average sleep figures over a period. Some fields are only
// set by SleepStatsOverall (SpO2, respiration, heart rate, Body Battery) or
// SleepStatsWeekly (low/high score).
type SleepAverages struct {
	AverageSleepSeconds        float64 `json:"averageSleepSeconds"`
	AverageSleepNeed           float64 `json:"averageSleepNeed"`
	AverageSleepScore          float64 `json:"averageSleepScore"`
	AverageLowSleepScore       float64 `json:"averageLowSleepScore"`
	AverageHighSleepScore      float64 `json:"averageHighSleepScore"`
	AverageLocalSleepStartTime float64 `json:"averageLocalSleepStartTime"`
	AverageLocalSleepEndTime   float64 `json:"averageLocalSleepEndTime"`
	AverageSpO2                float64 `json:"averageSpO2"`
	AverageRespiration         float64 `json:"averageRespiration"`
	AverageRestingHeartRate    float64 `json:"averageRestingHeartRate"`
	MeanAvgHeartRate           float64 `json:"meanAvgHeartRate"`
	AverageBodyBatteryChange   float64 `json:"averageBodyBatteryChange"`
}

// SleepWeek is one week's sleep averages.
type SleepWeek struct {
	WeekStartDate string `json:"weekStartDate"`
	WeekEndDate   string `json:"weekEndDate"`
	Values        struct {
		AverageSleepSeconds float64 `json:"averageSleepSeconds"`
		AverageSleepScore   float64 `json:"averageSleepScore"`
		AverageSleepNeed    float64 `json:"averageSleepNeed"`
		DeepTime            float64 `json:"deepTime"`
		LightTime           float64 `json:"lightTime"`
		RemTime             float64 `json:"remTime"`
		AwakeTime           float64 `json:"awakeTime"`
		SleepScoreQuality   string  `json:"sleepScoreQuality"`
		SleepDataDaysCount  int     `json:"sleepDataDaysCount"`
	} `json:"values"`
}

// SleepPeriodStats are sleep averages over a period, with a per-week
// breakdown for SleepStatsWeekly.
type SleepPeriodStats struct {
	Overall SleepAverages `json:"overallStats"`
	Weeks   []SleepWeek   `json:"individualStats"`
}

// SleepStatsWeekly returns sleep averages for the weeks weeks ending at end.
func (c *Client) SleepStatsWeekly(ctx context.Context, end time.Time, weeks int) (*SleepPeriodStats, error) {
	var out SleepPeriodStats
	if err := c.get(ctx, fmt.Sprintf("/sleep-service/stats/sleep/weekly/%s/%d", date(end), weeks), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SleepStatsOverall returns sleep averages between start and end.
func (c *Client) SleepStatsOverall(ctx context.Context, start, end time.Time) (*SleepPeriodStats, error) {
	var out SleepPeriodStats
	if err := c.get(ctx, fmt.Sprintf("/sleep-service/stats/sleep/overall/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
