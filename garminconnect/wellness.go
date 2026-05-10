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

// UserSummary holds the daily wellness summary for a user.
type UserSummary struct {
	UserProfileID              int     `json:"userProfileId"`
	TotalKilocalories          float64 `json:"totalKilocalories"`
	ActiveKilocalories         float64 `json:"activeKilocalories"`
	BmrKilocalories            float64 `json:"bmrKilocalories"`
	TotalSteps                 int     `json:"totalSteps"`
	DailyStepGoal              int     `json:"dailyStepGoal"`
	TotalDistanceMeters        float64 `json:"totalDistanceMeters"`
	WellnessStartTimeLocal     string  `json:"wellnessStartTimeLocal"`
	WellnessEndTimeLocal       string  `json:"wellnessEndTimeLocal"`
	DurationInMilliseconds     int64   `json:"durationInMilliseconds"`
	HighlyActiveSeconds        int     `json:"highlyActiveSeconds"`
	ActiveSeconds              int     `json:"activeSeconds"`
	SedentarySeconds           int     `json:"sedentarySeconds"`
	SleepingSeconds            int     `json:"sleepingSeconds"`
	ModerateDurationMinutes    int     `json:"moderateIntensityMinutes"`
	VigorousDurationMinutes    int     `json:"vigorousIntensityMinutes"`
	FloorsAscended             float64 `json:"floorsAscended"`
	FloorsDescended            float64 `json:"floorsDescended"`
	FloorsAscendedGoal         float64 `json:"floorsAscendedGoal"`
	MinHeartRate               int     `json:"minHeartRate"`
	MaxHeartRate               int     `json:"maxHeartRate"`
	RestingHeartRate           int     `json:"restingHeartRate"`
	LastSevenDaysAvgRestingHeartRate int `json:"lastSevenDaysAvgRestingHeartRate"`
	AbnormalHeartRateAlertsCount int   `json:"abnormalHeartRateAlertsCount"`
	AvgWakingRespirationValue  float64 `json:"avgWakingRespirationValue"`
	HighestRespirationValue    float64 `json:"highestRespirationValue"`
	LowestRespirationValue     float64 `json:"lowestRespirationValue"`
	LatestRespirationValue     float64 `json:"latestRespirationValue"`
	AvgStressDuration          int     `json:"avgStressDuration"`
	HighStressDuration         int     `json:"highStressDuration"`
	LowStressDuration          int     `json:"lowStressDuration"`
	RestStressDuration         int     `json:"restStressDuration"`
	BodyBatteryChargedValue    int     `json:"bodyBatteryChargedValue"`
	BodyBatteryDrainedValue    int     `json:"bodyBatteryDrainedValue"`
	BodyBatteryHighestValue    int     `json:"bodyBatteryHighestValue"`
	BodyBatteryLowestValue     int     `json:"bodyBatteryLowestValue"`
	BodyBatteryMostRecentValue int     `json:"bodyBatteryMostRecentValue"`
}

// BodyBatteryEntry is a single body battery reading.
type BodyBatteryEntry struct {
	StartTimestampGMT   string `json:"startTimestampGMT"`
	EndTimestampGMT     string `json:"endTimestampGMT"`
	StartTimestampLocal string `json:"startTimestampLocal"`
	EndTimestampLocal   string `json:"endTimestampLocal"`
	BodyBatteryValues   []struct {
		Timestamp int64  `json:"timestamp"` // unix ms
		Value     int    `json:"value"`
		Status    string `json:"status"`
	} `json:"bodyBatteryValuesArray"`
}

// StressData holds the all-day stress data for a single day.
type StressData struct {
	UserProfilePK int    `json:"userProfilePK"`
	CalendarDate  string `json:"calendarDate"`
	StartTimestampGMT   string `json:"startTimestampGMT"`
	EndTimestampGMT     string `json:"endTimestampGMT"`
	AvgStressLevel       int    `json:"avgStressLevel"`
	MaxStressLevel       int    `json:"maxStressLevel"`
	StressChartValueOffset int  `json:"stressChartValueOffset"`
	StressChartYAxisOrigin int  `json:"stressChartYAxisOrigin"`
	StressValuesArray     [][]int64 `json:"stressValuesArray"` // [timestamp_ms, stress_level]
	BodyBatteryValuesArray [][]int64 `json:"bodyBatteryValuesArray"`
}

// FloorsData holds floors ascended/descended for a day.
type FloorsData struct {
	UserProfilePK    int    `json:"userProfilePK"`
	CalendarDate     string `json:"calendarDate"`
	StartTimestampGMT string `json:"startTimestampGMT"`
	EndTimestampGMT  string `json:"endTimestampGMT"`
	FloorsValueDescriptorDTOList []struct {
		Key   string `json:"key"`
		Index int    `json:"index"`
	} `json:"floorValuesDescriptor"`
	FloorValuesArray [][]int64 `json:"floorValuesArray"` // [timestamp_ms, ascended, descended]
}

// HydrationData holds hydration intake for a day.
type HydrationData struct {
	UserProfilePK        int     `json:"userProfilePK"`
	CalendarDate         string  `json:"calendarDate"`
	ValueInML            float64 `json:"valueInML"`
	GoalInML             float64 `json:"goalInML"`
	DailyAverageinML     float64 `json:"dailyAverageinML"`
	SweatLossInML        float64 `json:"sweatLossInML"`
	ActivityIntakeInML   float64 `json:"activityIntakeInML"`
}

// RespirationData holds breathing rate data for a day.
type RespirationData struct {
	StartTimestampGMT        string  `json:"startTimestampGMT"`
	EndTimestampGMT          string  `json:"endTimestampGMT"`
	StartTimestampLocal      string  `json:"startTimestampLocal"`
	EndTimestampLocal        string  `json:"endTimestampLocal"`
	TodayAvgWakingRespirationValue float64 `json:"todayAvgWakingRespirationValue"`
	HighestRespirationValue  float64 `json:"highestRespirationValue"`
	LowestRespirationValue   float64 `json:"lowestRespirationValue"`
	RespirationValueDescriptorsDTOList []struct {
		Key   string `json:"key"`
		Index int    `json:"index"`
	} `json:"respirationValueDescriptorsDTOList"`
	RespirationValuesArray [][]float64 `json:"respirationValuesArray"` // [timestamp_ms, value]
}

// SpO2Data holds blood oxygen saturation data for a day.
type SpO2Data struct {
	UserProfilePK   int    `json:"userProfilePK"`
	CalendarDate    string `json:"calendarDate"`
	StartTimestampGMT string `json:"startTimestampGMT"`
	EndTimestampGMT   string `json:"endTimestampGMT"`
	AverageSpO2     float64 `json:"averageSpO2"`
	LowestSpO2      float64 `json:"lowestSpO2"`
	LastSevenDaysAvgSpO2 float64 `json:"lastSevenDaysAvgSpO2"`
	SpO2HourlyAverages []struct {
		StartTimestampGMT string  `json:"startTimestampGMT"`
		Value             float64 `json:"value"`
	} `json:"spO2HourlyAverages"`
}

// IntensityMinutesData holds weekly intensity minutes.
type IntensityMinutesData struct {
	UserProfilePK          int    `json:"userProfilePK"`
	CalendarDate           string `json:"calendarDate"`
	WeeklyGoal             int    `json:"weeklyGoal"`
	ModerateIntensityMinutes int  `json:"moderateIntensityMinutes"`
	VigorousIntensityMinutes int  `json:"vigorousIntensityMinutes"`
}

// StepEntry is a single steps reading for a time interval.
type StepEntry struct {
	StartTimestampGMT string `json:"startGMT"`
	EndTimestampGMT   string `json:"endGMT"`
	Steps             int    `json:"steps"`
	Pushes            int    `json:"pushes"`
	PrimaryActivityLevel string `json:"primaryActivityLevel"`
}

// BloodPressureSnapshot is a single blood pressure reading.
type BloodPressureSnapshot struct {
	TimestampGMT    string  `json:"measurementTimestampGMT"`
	TimestampLocal  string  `json:"measurementTimestampLocal"`
	Systolic        int     `json:"systolic"`
	Diastolic       int     `json:"diastolic"`
	Pulse           int     `json:"pulse"`
	Notes           string  `json:"notes"`
}

// BloodPressureSummary wraps blood pressure readings and aggregate stats.
type BloodPressureSummary struct {
	UserProfilePK int                      `json:"userProfilePK"`
	Measurements  []BloodPressureSnapshot  `json:"measurementSummaries"`
}

// UserSummary returns the daily wellness summary for the given date.
func (c *Client) UserSummary(d time.Time) (*UserSummary, error) {
	params := url.Values{"calendarDate": {date(d)}}
	var out UserSummary
	if err := c.get(fmt.Sprintf("/usersummary-service/usersummary/daily/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// BodyBattery returns body battery readings between start and end dates.
func (c *Client) BodyBattery(start, end time.Time) ([]BodyBatteryEntry, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out []BodyBatteryEntry
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/dailyBodyBattery/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AllDayStress returns all-day stress data for the given date.
func (c *Client) AllDayStress(d time.Time) (*StressData, error) {
	var out StressData
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/dailyStress/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Floors returns floors ascended/descended data for the given date.
func (c *Client) Floors(d time.Time) (*FloorsData, error) {
	var out FloorsData
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/dailyFloor/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Hydration returns hydration intake data for the given date.
func (c *Client) Hydration(d time.Time) (*HydrationData, error) {
	var out HydrationData
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/hydration/dailyReport/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Respiration returns breathing rate data for the given date.
func (c *Client) Respiration(d time.Time) (*RespirationData, error) {
	var out RespirationData
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/daily/respiration/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SpO2 returns blood oxygen saturation data for the given date.
func (c *Client) SpO2(d time.Time) (*SpO2Data, error) {
	var out SpO2Data
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/daily/spo2/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// IntensityMinutes returns intensity minutes data for the given date's week.
func (c *Client) IntensityMinutes(d time.Time) (*IntensityMinutesData, error) {
	params := url.Values{
		"startDate": {date(d)},
		"endDate":   {date(d)},
	}
	var out IntensityMinutesData
	if err := c.get(fmt.Sprintf("/usersummary-service/usersummary/intensity_minutes/daily/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Steps returns step data for the given date in 15-minute intervals.
func (c *Client) Steps(d time.Time) ([]StepEntry, error) {
	var out []StepEntry
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/dailySteps/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BloodPressure returns blood pressure measurements between start and end dates.
func (c *Client) BloodPressure(start, end time.Time) (*BloodPressureSummary, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out BloodPressureSummary
	if err := c.get(fmt.Sprintf("/bloodpressure-service/bloodpressure/range/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
