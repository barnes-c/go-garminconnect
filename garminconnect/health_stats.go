package garminconnect

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

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

// ActivityTotals are the user's totals for one activity type.
type ActivityTotals struct {
	ActivityType struct {
		TypeKey string `json:"typeKey"` // "all" for the grand total
	} `json:"activityType"`
	TotalActivities    int     `json:"totalActivities"`
	TotalDistance      float64 `json:"totalDistance"` // meters
	TotalDuration      float64 `json:"totalDuration"` // seconds
	TotalCalories      float64 `json:"totalCalories"`
	TotalElevationGain float64 `json:"totalElevationGain"`
}

// ActivityStatistics returns the user's activity totals per activity type.
func (c *Client) ActivityStatistics(ctx context.Context) ([]ActivityTotals, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	var out struct {
		UserMetrics []ActivityTotals `json:"userMetrics"`
	}
	if err := c.get(ctx, "/userstats-service/statistics/"+name, nil, &out); err != nil {
		return nil, err
	}
	return out.UserMetrics, nil
}

// WellnessAverages are daily averages over a period.
type WellnessAverages struct {
	AverageSteps            float64 `json:"averageSteps"`
	AverageRestingHeartRate float64 `json:"averageRestingHeartRate"`
	AverageHydrationInML    float64 `json:"averageHydrationInML"`
}

// WellnessAverages returns daily averages between start and end.
func (c *Client) WellnessAverages(ctx context.Context, start, end time.Time) (*WellnessAverages, error) {
	var out WellnessAverages
	if err := c.get(ctx, fmt.Sprintf("/usersummary-service/stats/averages/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DailyMovement is a day's movement, sampled per minute.
type DailyMovement struct {
	CalendarDate string `json:"calendarDate"`
	// Descriptors name the columns of each Values row: an epoch-ms timestamp
	// and the movement level.
	Descriptors []struct {
		Index int    `json:"index"`
		Key   string `json:"key"`
	} `json:"movementValueDescriptors"`
	Values [][]float64 `json:"movementValues"`
}

// DailyMovement returns a day's per-minute movement.
func (c *Client) DailyMovement(ctx context.Context, d time.Time) (*DailyMovement, error) {
	var out DailyMovement
	params := url.Values{"calendarDate": {date(d)}}
	if err := c.get(ctx, "/wellness-service/wellness/dailyMovement", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// IntensityMinutesStats are the daily intensity minutes over a period.
type IntensityMinutesStats struct {
	WeeklyGoal int `json:"weeklyGoal"`
	Days       []struct {
		CalendarDate    string `json:"calendarDate"`
		ModerateMinutes int    `json:"moderateMinutes"`
		VigorousMinutes int    `json:"vigorousMinutes"`
	} `json:"dailyStats"`
}

// IntensityMinutesStats returns intensity minutes for the days days ending at
// end.
func (c *Client) IntensityMinutesStats(ctx context.Context, end time.Time, days int) (*IntensityMinutesStats, error) {
	var out IntensityMinutesStats
	if err := c.get(ctx, fmt.Sprintf("/wellness-service/stats/hourly/im/%s/%d", date(end), days), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MonthlyStepStat is one month's step totals.
type MonthlyStepStat struct {
	CalendarDate string `json:"calendarDate"`
	Values       struct {
		TotalSteps    int `json:"totalSteps"`
		TotalDistance int `json:"totalDistance"` // meters
		StepsGoal     int `json:"TotalStepsGoal"`
		DaysWithData  int `json:"wellnessDataDaysCount"`
	} `json:"values"`
}

// MonthlySteps returns step totals for the months months ending at end.
func (c *Client) MonthlySteps(ctx context.Context, end time.Time, months int) ([]MonthlyStepStat, error) {
	var out []MonthlyStepStat
	if err := c.get(ctx, fmt.Sprintf("/usersummary-service/stats/steps/monthly/%s/%d", date(end), months), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeekAverageRestingHR returns the average resting heart rate of the week
// ending at d.
func (c *Client) WeekAverageRestingHR(ctx context.Context, d time.Time) (int, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return 0, err
	}
	var out struct {
		WeekAvgRestingHR int `json:"weekAvgRestingHR"`
	}
	params := url.Values{"calendarDate": {date(d)}}
	if err := c.get(ctx, "/usersummary-service/usersummary/weekAverageRestingHR/"+name, params, &out); err != nil {
		return 0, err
	}
	return out.WeekAvgRestingHR, nil
}

// LastSyncTimestamp returns when the user's device last synced, as Garmin
// reports it (e.g. "2026-01-01T07:30:00.000", no time zone).
func (c *Client) LastSyncTimestamp(ctx context.Context) (string, error) {
	var out string
	if err := c.get(ctx, "/wellness-service/wellness/syncTimestamp", nil, &out); err != nil {
		return "", err
	}
	return out, nil
}
