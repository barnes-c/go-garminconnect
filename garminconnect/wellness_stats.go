package garminconnect

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"time"
)

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

// DailyStepStat holds aggregated step statistics for a single calendar date.
type DailyStepStat struct {
	CalendarDate  string `json:"calendarDate"`
	TotalSteps    int    `json:"totalSteps"`
	TotalDistance int    `json:"totalDistance"`
	StepGoal      int    `json:"stepGoal"`
}

// DailySteps returns daily step totals between start and end dates.
func (c *Client) DailySteps(ctx context.Context, start, end time.Time) ([]DailyStepStat, error) {
	var out []DailyStepStat
	if err := c.get(ctx, fmt.Sprintf("/usersummary-service/stats/steps/daily/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeeklyStepStat holds weekly aggregated step statistics.
type WeeklyStepStat struct {
	CalendarDate  string `json:"calendarDate"`
	TotalSteps    int    `json:"totalSteps"`
	TotalDistance int    `json:"totalDistance"`
}

// WeeklySteps returns weekly step totals ending on end for the given number of weeks.
func (c *Client) WeeklySteps(ctx context.Context, end time.Time, weeks int) ([]WeeklyStepStat, error) {
	var out []WeeklyStepStat
	if err := c.get(ctx, fmt.Sprintf("/usersummary-service/stats/steps/weekly/%s/%d", date(end), weeks), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeeklyStressStat holds weekly aggregated stress statistics.
type WeeklyStressStat struct {
	CalendarDate   string `json:"calendarDate"`
	AvgStressLevel int    `json:"value"`
	MaxStressLevel int    `json:"maxStressLevel"`
	StressDuration int    `json:"stressDuration"`
}

// WeeklyStress returns weekly stress aggregates ending on end for the given number of weeks.
func (c *Client) WeeklyStress(ctx context.Context, end time.Time, weeks int) ([]WeeklyStressStat, error) {
	var out []WeeklyStressStat
	if err := c.get(ctx, fmt.Sprintf("/usersummary-service/stats/stress/weekly/%s/%d", date(end), weeks), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeeklyIMStat holds weekly intensity minutes statistics.
type WeeklyIMStat struct {
	CalendarDate             string `json:"calendarDate"`
	WeeklyGoal               int    `json:"weeklyGoal"`
	ModerateIntensityMinutes int    `json:"moderateValue"`
	VigorousIntensityMinutes int    `json:"vigorousValue"`
}

// WeeklyIntensityMinutes returns weekly intensity minutes between start and end dates.
func (c *Client) WeeklyIntensityMinutes(ctx context.Context, start, end time.Time) ([]WeeklyIMStat, error) {
	var out []WeeklyIMStat
	if err := c.get(ctx, fmt.Sprintf("/usersummary-service/stats/im/weekly/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DailyCalories is one day's active and resting (BMR) calorie totals.
type DailyCalories struct {
	CalendarDate string  `json:"calendarDate"`
	Active       float64 `json:"active"`
	Resting      float64 `json:"resting"`
	Total        float64 `json:"total"`
}

// calorieEntry is one row of a wellness-stats calorie series.
type calorieEntry struct {
	CalendarDate string   `json:"calendarDate"`
	Value        *float64 `json:"value"`
}

// CaloriesDaily returns daily active and resting calories between start and
// end. Both series come from a single request; days present in only one series
// report zero for the other.
func (c *Client) CaloriesDaily(ctx context.Context, start, end time.Time) ([]DailyCalories, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	// 22 = active calories, 23 = BMR calories.
	params := url.Values{
		"fromDate":  {date(start)},
		"untilDate": {date(end)},
		"metricId":  {"22", "23"},
	}
	var raw struct {
		AllMetrics struct {
			MetricsMap struct {
				Active  []calorieEntry `json:"WELLNESS_ACTIVE_CALORIES"`
				Resting []calorieEntry `json:"WELLNESS_BMR_CALORIES"`
			} `json:"metricsMap"`
		} `json:"allMetrics"`
	}
	if err := c.get(ctx, "/userstats-service/wellness/daily/"+name, params, &raw); err != nil {
		return nil, err
	}

	byDate := map[string]*DailyCalories{}
	collect := func(entries []calorieEntry, set func(*DailyCalories, float64)) {
		for _, e := range entries {
			if e.Value == nil || e.CalendarDate == "" {
				continue
			}
			day, ok := byDate[e.CalendarDate]
			if !ok {
				day = &DailyCalories{CalendarDate: e.CalendarDate}
				byDate[e.CalendarDate] = day
			}
			set(day, *e.Value)
		}
	}
	collect(raw.AllMetrics.MetricsMap.Active, func(d *DailyCalories, v float64) { d.Active = v })
	collect(raw.AllMetrics.MetricsMap.Resting, func(d *DailyCalories, v float64) { d.Resting = v })

	out := make([]DailyCalories, 0, len(byDate))
	for _, d := range byDate {
		d.Total = d.Active + d.Resting
		out = append(out, *d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CalendarDate < out[j].CalendarDate })
	return out, nil
}
