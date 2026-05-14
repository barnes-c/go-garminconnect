package garminconnect

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// UserSummary holds the daily wellness summary for a user.
type UserSummary struct {
	UserProfileID                    int     `json:"userProfileId"`
	TotalKilocalories                float64 `json:"totalKilocalories"`
	ActiveKilocalories               float64 `json:"activeKilocalories"`
	BmrKilocalories                  float64 `json:"bmrKilocalories"`
	TotalSteps                       int     `json:"totalSteps"`
	DailyStepGoal                    int     `json:"dailyStepGoal"`
	TotalDistanceMeters              float64 `json:"totalDistanceMeters"`
	WellnessStartTimeLocal           string  `json:"wellnessStartTimeLocal"`
	WellnessEndTimeLocal             string  `json:"wellnessEndTimeLocal"`
	DurationInMilliseconds           int64   `json:"durationInMilliseconds"`
	HighlyActiveSeconds              int     `json:"highlyActiveSeconds"`
	ActiveSeconds                    int     `json:"activeSeconds"`
	SedentarySeconds                 int     `json:"sedentarySeconds"`
	SleepingSeconds                  int     `json:"sleepingSeconds"`
	ModerateDurationMinutes          int     `json:"moderateIntensityMinutes"`
	VigorousDurationMinutes          int     `json:"vigorousIntensityMinutes"`
	FloorsAscended                   float64 `json:"floorsAscended"`
	FloorsDescended                  float64 `json:"floorsDescended"`
	FloorsAscendedGoal               float64 `json:"floorsAscendedGoal"`
	MinHeartRate                     int     `json:"minHeartRate"`
	MaxHeartRate                     int     `json:"maxHeartRate"`
	RestingHeartRate                 int     `json:"restingHeartRate"`
	LastSevenDaysAvgRestingHeartRate int     `json:"lastSevenDaysAvgRestingHeartRate"`
	AbnormalHeartRateAlertsCount     int     `json:"abnormalHeartRateAlertsCount"`
	AvgWakingRespirationValue        float64 `json:"avgWakingRespirationValue"`
	HighestRespirationValue          float64 `json:"highestRespirationValue"`
	LowestRespirationValue           float64 `json:"lowestRespirationValue"`
	LatestRespirationValue           float64 `json:"latestRespirationValue"`
	AvgStressDuration                int     `json:"avgStressDuration"`
	HighStressDuration               int     `json:"highStressDuration"`
	LowStressDuration                int     `json:"lowStressDuration"`
	RestStressDuration               int     `json:"restStressDuration"`
	BodyBatteryChargedValue          int     `json:"bodyBatteryChargedValue"`
	BodyBatteryDrainedValue          int     `json:"bodyBatteryDrainedValue"`
	BodyBatteryHighestValue          int     `json:"bodyBatteryHighestValue"`
	BodyBatteryLowestValue           int     `json:"bodyBatteryLowestValue"`
	BodyBatteryMostRecentValue       int     `json:"bodyBatteryMostRecentValue"`
}

// BodyBatteryEntry is a single body battery reading.
type BodyBatteryEntry struct {
	StartTimestampGMT   string          `json:"startTimestampGMT"`
	EndTimestampGMT     string          `json:"endTimestampGMT"`
	StartTimestampLocal string          `json:"startTimestampLocal"`
	EndTimestampLocal   string          `json:"endTimestampLocal"`
	BodyBatteryValues   json.RawMessage `json:"bodyBatteryValuesArray"` // [[timestamp_ms, level|null], ...]
}

// StressData holds the all-day stress data for a single day.
type StressData struct {
	UserProfilePK          int             `json:"userProfilePK"`
	CalendarDate           string          `json:"calendarDate"`
	StartTimestampGMT      string          `json:"startTimestampGMT"`
	EndTimestampGMT        string          `json:"endTimestampGMT"`
	AvgStressLevel         int             `json:"avgStressLevel"`
	MaxStressLevel         int             `json:"maxStressLevel"`
	StressChartValueOffset int             `json:"stressChartValueOffset"`
	StressChartYAxisOrigin int             `json:"stressChartYAxisOrigin"`
	StressValuesArray      [][]int64       `json:"stressValuesArray"`      // [timestamp_ms, stress_level]
	BodyBatteryValuesArray json.RawMessage `json:"bodyBatteryValuesArray"` // mixed-type rows: [ts, status, level, version]
}

// FloorsData holds floors ascended/descended for a day.
type FloorsData struct {
	UserProfilePK                int    `json:"userProfilePK"`
	CalendarDate                 string `json:"calendarDate"`
	StartTimestampGMT            string `json:"startTimestampGMT"`
	EndTimestampGMT              string `json:"endTimestampGMT"`
	FloorsValueDescriptorDTOList []struct {
		Key   string `json:"key"`
		Index int    `json:"index"`
	} `json:"floorsValueDescriptorDTOList"`
	FloorValuesArray json.RawMessage `json:"floorValuesArray"` // [["startTimeGMT", "endTimeGMT", ascended, descended], ...]
}

// HydrationData holds hydration intake for a day.
type HydrationData struct {
	UserProfilePK      int     `json:"userProfilePK"`
	CalendarDate       string  `json:"calendarDate"`
	ValueInML          float64 `json:"valueInML"`
	GoalInML           float64 `json:"goalInML"`
	DailyAverageinML   float64 `json:"dailyAverageinML"`
	SweatLossInML      float64 `json:"sweatLossInML"`
	ActivityIntakeInML float64 `json:"activityIntakeInML"`
}

// RespirationData holds breathing rate data for a day.
type RespirationData struct {
	StartTimestampGMT                  string  `json:"startTimestampGMT"`
	EndTimestampGMT                    string  `json:"endTimestampGMT"`
	StartTimestampLocal                string  `json:"startTimestampLocal"`
	EndTimestampLocal                  string  `json:"endTimestampLocal"`
	TodayAvgWakingRespirationValue     float64 `json:"avgWakingRespirationValue"`
	HighestRespirationValue            float64 `json:"highestRespirationValue"`
	LowestRespirationValue             float64 `json:"lowestRespirationValue"`
	RespirationValueDescriptorsDTOList []struct {
		Key   string `json:"key"`
		Index int    `json:"index"`
	} `json:"respirationValueDescriptorsDTOList"`
	RespirationValuesArray [][]float64 `json:"respirationValuesArray"` // [timestamp_ms, value]
}

// SpO2Data holds blood oxygen saturation data for a day.
type SpO2Data struct {
	UserProfilePK        int         `json:"userProfilePK"`
	CalendarDate         string      `json:"calendarDate"`
	StartTimestampGMT    string      `json:"startTimestampGMT"`
	EndTimestampGMT      string      `json:"endTimestampGMT"`
	AverageSpO2          float64     `json:"averageSpO2"`
	LowestSpO2           float64     `json:"lowestSpO2"`
	LastSevenDaysAvgSpO2 float64     `json:"lastSevenDaysAvgSpO2"`
	SpO2HourlyAverages   [][]float64 `json:"spO2HourlyAverages"` // [timestamp_ms, value]
}

// IntensityMinutesData holds weekly intensity minutes.
type IntensityMinutesData struct {
	UserProfilePK            int    `json:"userProfilePK"`
	CalendarDate             string `json:"calendarDate"`
	WeeklyGoal               int    `json:"weeklyGoal"`
	ModerateIntensityMinutes int    `json:"moderateIntensityMinutes"`
	VigorousIntensityMinutes int    `json:"vigorousIntensityMinutes"`
}

// StepEntry is a single steps reading for a time interval.
type StepEntry struct {
	StartTimestampGMT    string `json:"startGMT"`
	EndTimestampGMT      string `json:"endGMT"`
	Steps                int    `json:"steps"`
	Pushes               int    `json:"pushes"`
	PrimaryActivityLevel string `json:"primaryActivityLevel"`
}

// BloodPressureSnapshot is a single blood pressure reading.
type BloodPressureSnapshot struct {
	TimestampGMT   string `json:"measurementTimestampGMT"`
	TimestampLocal string `json:"measurementTimestampLocal"`
	Systolic       int    `json:"systolic"`
	Diastolic      int    `json:"diastolic"`
	Pulse          int    `json:"pulse"`
	Notes          string `json:"notes"`
}

// BloodPressureSummary wraps blood pressure readings and aggregate stats.
type BloodPressureSummary struct {
	UserProfilePK int                     `json:"userProfilePK"`
	Measurements  []BloodPressureSnapshot `json:"measurementSummaries"`
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
	if err := c.get("/wellness-service/wellness/bodyBattery/reports/daily", params, &out); err != nil {
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
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/floorsChartData/daily/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Hydration returns hydration intake data for the given date.
func (c *Client) Hydration(d time.Time) (*HydrationData, error) {
	var out HydrationData
	if err := c.get(fmt.Sprintf("/usersummary-service/usersummary/hydration/daily/%s", date(d)), nil, &out); err != nil {
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
	var out IntensityMinutesData
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/daily/im/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Steps returns step data for the given date in 15-minute intervals.
func (c *Client) Steps(d time.Time) ([]StepEntry, error) {
	params := url.Values{"date": {date(d)}}
	var out []StepEntry
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/dailySummaryChart/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BloodPressure returns blood pressure measurements between start and end dates.
func (c *Client) BloodPressure(start, end time.Time) (*BloodPressureSummary, error) {
	params := url.Values{"includeAll": {"true"}}
	var out BloodPressureSummary
	if err := c.get(fmt.Sprintf("/bloodpressure-service/bloodpressure/range/%s/%s", date(start), date(end)), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DailyStepStat holds aggregated step statistics for a single calendar date.
type DailyStepStat struct {
	CalendarDate  string `json:"calendarDate"`
	TotalSteps    int    `json:"totalSteps"`
	TotalDistance int    `json:"totalDistance"`
	StepGoal      int    `json:"stepGoal"`
}

// WeeklyStepStat holds weekly aggregated step statistics.
type WeeklyStepStat struct {
	CalendarDate  string `json:"calendarDate"`
	TotalSteps    int    `json:"totalSteps"`
	TotalDistance int    `json:"totalDistance"`
}

// WeeklyStressStat holds weekly aggregated stress statistics.
type WeeklyStressStat struct {
	CalendarDate     string `json:"calendarDate"`
	AvgStressLevel   int    `json:"avgStressLevel"`
	MaxStressLevel   int    `json:"maxStressLevel"`
	StressDuration   int    `json:"stressDuration"`
	RestDuration     int    `json:"restDuration"`
	ActivityDuration int    `json:"activityDuration"`
}

// BodyBatteryEvent is a single body battery charge or drain event.
type BodyBatteryEvent struct {
	EventTimestamp    string `json:"eventTimestamp"`
	Event             string `json:"event"` // "CHARGE" or "DRAIN"
	DurationInMS      int64  `json:"durationInMS"`
	BodyBatteryImpact int    `json:"bodyBatteryImpact"`
	FeedbackType      string `json:"feedbackType"`
	FeedbackShortType string `json:"feedbackShortType"`
}

// WeeklyIMStat holds weekly intensity minutes statistics.
type WeeklyIMStat struct {
	CalendarDate             string `json:"calendarDate"`
	WeeklyGoal               int    `json:"weeklyGoal"`
	ModerateIntensityMinutes int    `json:"moderateIntensityMinutes"`
	VigorousIntensityMinutes int    `json:"vigorousIntensityMinutes"`
}

// StepsData returns intraday step data in 15-minute intervals via the summary chart endpoint.
func (c *Client) StepsData(d time.Time) ([]StepEntry, error) {
	params := url.Values{"calendarDate": {date(d)}}
	var out []StepEntry
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/dailySummaryChart/%s", c.displayName), params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DailySteps returns daily step totals between start and end dates.
func (c *Client) DailySteps(start, end time.Time) ([]DailyStepStat, error) {
	var out []DailyStepStat
	if err := c.get(fmt.Sprintf("/usersummary-service/stats/steps/daily/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeeklySteps returns weekly step totals ending on end for the given number of weeks.
func (c *Client) WeeklySteps(end time.Time, weeks int) ([]WeeklyStepStat, error) {
	var out []WeeklyStepStat
	if err := c.get(fmt.Sprintf("/usersummary-service/stats/steps/weekly/%s/%d", date(end), weeks), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeeklyStress returns weekly stress aggregates ending on end for the given number of weeks.
func (c *Client) WeeklyStress(end time.Time, weeks int) ([]WeeklyStressStat, error) {
	var out []WeeklyStressStat
	if err := c.get(fmt.Sprintf("/usersummary-service/stats/stress/weekly/%s/%d", date(end), weeks), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BodyBatteryEvents returns body battery charge/drain events for the given date.
func (c *Client) BodyBatteryEvents(d time.Time) ([]BodyBatteryEvent, error) {
	var out []BodyBatteryEvent
	if err := c.get(fmt.Sprintf("/wellness-service/wellness/bodyBattery/events/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeeklyIntensityMinutes returns weekly intensity minutes between start and end dates.
func (c *Client) WeeklyIntensityMinutes(start, end time.Time) ([]WeeklyIMStat, error) {
	var out []WeeklyIMStat
	if err := c.get(fmt.Sprintf("/usersummary-service/stats/im/weekly/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AllDayEvents returns wellness events (meals, stress, etc.) for the given date.
func (c *Client) AllDayEvents(d time.Time) (map[string]json.RawMessage, error) {
	params := url.Values{"calendarDate": {date(d)}}
	var out map[string]json.RawMessage
	if err := c.get("/wellness-service/wellness/dailyEvents", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LifestyleData returns lifestyle logging data for the given date.
func (c *Client) LifestyleData(d time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/lifestylelogging-service/dailyLog/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddHydration logs a hydration intake entry.
// valueML is the amount in millilitres; timestamp should be an RFC3339 timestamp.
func (c *Client) AddHydration(valueML float64, timestamp, cdate string) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	body := map[string]any{
		"calendarDate":   cdate,
		"valueInML":      valueML,
		"userProfilePK":  0, // filled server-side
		"timestampGMT":   timestamp,
		"timestampLocal": timestamp,
	}
	if err := c.put("/usersummary-service/usersummary/hydration/log", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetBloodPressure records a blood pressure measurement.
func (c *Client) SetBloodPressure(systolic, diastolic, pulse int, timestamp, notes string) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	body := map[string]any{
		"systolic":                systolic,
		"diastolic":               diastolic,
		"pulse":                   pulse,
		"measurementTimestampGMT": timestamp,
		"notes":                   notes,
	}
	if err := c.post("/bloodpressure-service/bloodpressure", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteBloodPressure removes a blood pressure reading by its version and date.
func (c *Client) DeleteBloodPressure(cdate string, version int) error {
	return c.del(fmt.Sprintf("/bloodpressure-service/bloodpressure/%s/%d", cdate, version))
}
