package garminconnect

import (
	"context"
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
	FloorsAscendedGoal               float64 `json:"userFloorsAscendedGoal"`
	MinHeartRate                     int     `json:"minHeartRate"`
	MaxHeartRate                     int     `json:"maxHeartRate"`
	RestingHeartRate                 int     `json:"restingHeartRate"`
	LastSevenDaysAvgRestingHeartRate int     `json:"lastSevenDaysAvgRestingHeartRate"`
	AbnormalHeartRateAlertsCount     int     `json:"abnormalHeartRateAlertsCount"`
	AvgWakingRespirationValue        float64 `json:"avgWakingRespirationValue"`
	HighestRespirationValue          float64 `json:"highestRespirationValue"`
	LowestRespirationValue           float64 `json:"lowestRespirationValue"`
	LatestRespirationValue           float64 `json:"latestRespirationValue"`
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
	UserProfilePK                int               `json:"userProfilePK"`
	CalendarDate                 string            `json:"calendarDate"`
	StartTimestampGMT            string            `json:"startTimestampGMT"`
	EndTimestampGMT              string            `json:"endTimestampGMT"`
	FloorsValueDescriptorDTOList []ValueDescriptor `json:"floorsValueDescriptorDTOList"`
	FloorValuesArray             json.RawMessage   `json:"floorValuesArray"` // [["startTimeGMT", "endTimeGMT", ascended, descended], ...]
}

// HydrationData holds hydration intake for a day.
type HydrationData struct {
	UserProfilePK      int     `json:"userId"`
	CalendarDate       string  `json:"calendarDate"`
	ValueInML          float64 `json:"valueInML"`
	GoalInML           float64 `json:"goalInML"`
	DailyAverageinML   float64 `json:"dailyAverageinML"`
	SweatLossInML      float64 `json:"sweatLossInML"`
	ActivityIntakeInML float64 `json:"activityIntakeInML"`
	// LastEntryTimestampLocal is only set by AddHydration.
	LastEntryTimestampLocal string `json:"lastEntryTimestampLocal"`
}

// RespirationData holds breathing rate data for a day.
type RespirationData struct {
	StartTimestampGMT                  string            `json:"startTimestampGMT"`
	EndTimestampGMT                    string            `json:"endTimestampGMT"`
	StartTimestampLocal                string            `json:"startTimestampLocal"`
	EndTimestampLocal                  string            `json:"endTimestampLocal"`
	TodayAvgWakingRespirationValue     float64           `json:"avgWakingRespirationValue"`
	HighestRespirationValue            float64           `json:"highestRespirationValue"`
	LowestRespirationValue             float64           `json:"lowestRespirationValue"`
	RespirationValueDescriptorsDTOList []ValueDescriptor `json:"respirationValueDescriptorsDTOList"`
	RespirationValuesArray             [][]float64       `json:"respirationValuesArray"` // [timestamp_ms, value]
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

// IntensityMinutesData holds intensity minutes for one day, plus the
// week-to-date aggregates Garmin reports alongside. WeeklyTotal counts
// vigorous minutes double, matching Garmin's goal accounting.
type IntensityMinutesData struct {
	UserProfilePK            int    `json:"userProfilePK"`
	CalendarDate             string `json:"calendarDate"`
	WeeklyGoal               int    `json:"weekGoal"`
	ModerateIntensityMinutes int    `json:"moderateMinutes"`
	VigorousIntensityMinutes int    `json:"vigorousMinutes"`
	WeeklyModerate           int    `json:"weeklyModerate"`
	WeeklyVigorous           int    `json:"weeklyVigorous"`
	WeeklyTotal              int    `json:"weeklyTotal"`
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
	Version          int64  `json:"version"` // identifies the reading for delete/update
	TimestampGMT     string `json:"measurementTimestampGMT"`
	TimestampLocal   string `json:"measurementTimestampLocal"`
	Systolic         int    `json:"systolic"`
	Diastolic        int    `json:"diastolic"`
	Pulse            int    `json:"pulse"`
	MultiMeasurement bool   `json:"multiMeasurement"`
	Notes            string `json:"notes"`
	SourceType       string `json:"sourceType"`
	Category         string `json:"category"`
	CategoryName     string `json:"categoryName"`
}

// BloodPressureDay groups one day's blood pressure readings with their range.
type BloodPressureDay struct {
	StartDate         string                  `json:"startDate"`
	EndDate           string                  `json:"endDate"`
	HighSystolic      int                     `json:"highSystolic"`
	HighDiastolic     int                     `json:"highDiastolic"`
	LowSystolic       int                     `json:"lowSystolic"`
	LowDiastolic      int                     `json:"lowDiastolic"`
	NumOfMeasurements int                     `json:"numOfMeasurements"`
	Category          string                  `json:"category"`
	CategoryName      string                  `json:"categoryName"`
	Measurements      []BloodPressureSnapshot `json:"measurements"`
}

// BloodPressureSummary holds blood pressure readings grouped by day.
type BloodPressureSummary struct {
	From  string             `json:"from"`
	Until string             `json:"until"`
	Days  []BloodPressureDay `json:"measurementSummaries"`
}

// UserSummary returns the daily wellness summary for the given date.
func (c *Client) UserSummary(ctx context.Context, d time.Time) (*UserSummary, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	params := url.Values{"calendarDate": {date(d)}}
	var out UserSummary
	if err := c.get(ctx, "/usersummary-service/usersummary/daily/"+name, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// BodyBattery returns body battery readings between start and end dates.
func (c *Client) BodyBattery(ctx context.Context, start, end time.Time) ([]BodyBatteryEntry, error) {
	params := url.Values{
		"startDate": {date(start)},
		"endDate":   {date(end)},
	}
	var out []BodyBatteryEntry
	if err := c.get(ctx, "/wellness-service/wellness/bodyBattery/reports/daily", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AllDayStress returns all-day stress data for the given date.
func (c *Client) AllDayStress(ctx context.Context, d time.Time) (*StressData, error) {
	var out StressData
	if err := c.get(ctx, fmt.Sprintf("/wellness-service/wellness/dailyStress/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Floors returns floors ascended/descended data for the given date.
func (c *Client) Floors(ctx context.Context, d time.Time) (*FloorsData, error) {
	var out FloorsData
	if err := c.get(ctx, fmt.Sprintf("/wellness-service/wellness/floorsChartData/daily/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Hydration returns hydration intake data for the given date.
func (c *Client) Hydration(ctx context.Context, d time.Time) (*HydrationData, error) {
	var out HydrationData
	if err := c.get(ctx, fmt.Sprintf("/usersummary-service/usersummary/hydration/daily/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Respiration returns breathing rate data for the given date.
func (c *Client) Respiration(ctx context.Context, d time.Time) (*RespirationData, error) {
	var out RespirationData
	if err := c.get(ctx, fmt.Sprintf("/wellness-service/wellness/daily/respiration/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SpO2 returns blood oxygen saturation data for the given date.
func (c *Client) SpO2(ctx context.Context, d time.Time) (*SpO2Data, error) {
	var out SpO2Data
	if err := c.get(ctx, fmt.Sprintf("/wellness-service/wellness/daily/spo2/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// IntensityMinutes returns intensity minutes data for the given date's week.
func (c *Client) IntensityMinutes(ctx context.Context, d time.Time) (*IntensityMinutesData, error) {
	var out IntensityMinutesData
	if err := c.get(ctx, fmt.Sprintf("/wellness-service/wellness/daily/im/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Steps returns step data for the given date in 15-minute intervals.
func (c *Client) Steps(ctx context.Context, d time.Time) ([]StepEntry, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	params := url.Values{"date": {date(d)}}
	var out []StepEntry
	if err := c.get(ctx, "/wellness-service/wellness/dailySummaryChart/"+name, params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BloodPressure returns blood pressure measurements between start and end dates.
func (c *Client) BloodPressure(ctx context.Context, start, end time.Time) (*BloodPressureSummary, error) {
	params := url.Values{"includeAll": {"true"}}
	var out BloodPressureSummary
	if err := c.get(ctx, fmt.Sprintf("/bloodpressure-service/bloodpressure/range/%s/%s", date(start), date(end)), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
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

// BodyBatteryEvents returns body battery charge/drain events for the given date.
func (c *Client) BodyBatteryEvents(ctx context.Context, d time.Time) ([]BodyBatteryEvent, error) {
	var out []BodyBatteryEvent
	if err := c.get(ctx, fmt.Sprintf("/wellness-service/wellness/bodyBattery/events/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AllDayEvents returns wellness events (meals, stress, etc.) for the given date.
func (c *Client) AllDayEvents(ctx context.Context, d time.Time) (map[string]json.RawMessage, error) {
	params := url.Values{"calendarDate": {date(d)}}
	var out map[string]json.RawMessage
	if err := c.get(ctx, "/wellness-service/wellness/dailyEvents", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LifestyleData returns lifestyle logging data for the given date.
func (c *Client) LifestyleData(ctx context.Context, d time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/lifestylelogging-service/dailyLog/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddHydration logs valueML millilitres of fluid intake at the given time; its
// location supplies the local time and calendar day. A negative value
// subtracts, which is the only way to undo an entry.
// It returns the day's updated hydration.
func (c *Client) AddHydration(ctx context.Context, valueML float64, at time.Time) (*HydrationData, error) {
	var out HydrationData
	body := map[string]any{
		"calendarDate":   date(at),
		"valueInML":      valueML,
		"timestampLocal": timestamp(at),
	}
	if err := c.put(ctx, "/usersummary-service/usersummary/hydration/log", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DailyLastBloodPressure returns the last blood pressure reading of each day
// between start and end.
func (c *Client) DailyLastBloodPressure(ctx context.Context, start, end time.Time) ([]BloodPressureSnapshot, error) {
	var out []BloodPressureSnapshot
	if err := c.get(ctx, fmt.Sprintf("/bloodpressure-service/bloodpressure/daily/last/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WeeklyBloodPressure returns blood pressure grouped by week (in Days, one
// entry per week) for the 52 weeks ending at end.
func (c *Client) WeeklyBloodPressure(ctx context.Context, end time.Time) (*BloodPressureSummary, error) {
	var out BloodPressureSummary
	path := fmt.Sprintf("/bloodpressure-service/bloodpressure/weeklyRange/%s/%s", date(weeklyRangeStart(end)), date(end))
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetBloodPressure records a manual blood pressure reading taken at the given
// time; its location supplies the local time Garmin displays.
func (c *Client) SetBloodPressure(ctx context.Context, systolic, diastolic, pulse int, at time.Time, notes string) (*BloodPressureSnapshot, error) {
	// Centiseconds, unlike timestamp(): the format this endpoint was verified with.
	const layout = "2006-01-02T15:04:05.00"
	body := map[string]any{
		"systolic":                  systolic,
		"diastolic":                 diastolic,
		"pulse":                     pulse,
		"measurementTimestampLocal": at.Format(layout),
		"measurementTimestampGMT":   at.UTC().Format(layout),
		"sourceType":                "MANUAL",
		"notes":                     notes,
	}
	var out BloodPressureSnapshot
	if err := c.post(ctx, "/bloodpressure-service/bloodpressure", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetBloodPressureNotes changes the note of a reading from BloodPressure.
// Garmin stores only the note; to correct the values, delete the reading and
// record it again.
func (c *Client) SetBloodPressureNotes(ctx context.Context, m BloodPressureSnapshot, notes string) error {
	m.Notes = notes
	return c.put(ctx, "/bloodpressure-service/bloodpressure", m, nil)
}

// DeleteBloodPressure removes a blood pressure reading by its version and day.
func (c *Client) DeleteBloodPressure(ctx context.Context, d time.Time, version int64) error {
	return c.del(ctx, fmt.Sprintf("/bloodpressure-service/bloodpressure/%s/%d", date(d), version))
}

// DownloadHealthSnapshot returns the Health Snapshot ZIP for the given date.
func (c *Client) DownloadHealthSnapshot(ctx context.Context, d time.Time) ([]byte, error) {
	return c.getBytes(ctx, "/download-service/files/wellness/"+date(d), nil)
}
