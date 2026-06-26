package garminconnect

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// TrainingReadiness holds the training readiness score for a day.
type TrainingReadiness struct {
	UserProfilePK     int     `json:"userProfilePK"`
	CalendarDate      string  `json:"calendarDate"`
	Score             int     `json:"score"`
	ScoreQualifier    string  `json:"scoreQualifier"`
	SleepScore        int     `json:"sleepScore"`
	RecoveryTime      int     `json:"recoveryTime"` // minutes
	AcuteLoad         float64 `json:"acuteLoad"`
	TrainingLoad      float64 `json:"trainingLoad"`
	SleepHistoryScore int     `json:"sleepHistoryScore"`
	HrvWeeklyAverage  int     `json:"hrvWeeklyAverage"`
	FeedbackPhrase    string  `json:"feedbackPhrase"`
}

// TrainingStatusEntry is a single day of training status data.
type TrainingStatusEntry struct {
	CalendarDate       string  `json:"calendarDate"`
	TrainingStatusType string  `json:"trainingStatusType"`
	TrainingLoadType   string  `json:"trainingLoadType"`
	WorkoutGoal        float64 `json:"workoutGoal"`
	AtpPlanLowLoad     float64 `json:"atpPlanLowLoad"`
	AtpPlanHighLoad    float64 `json:"atpPlanHighLoad"`
}

// MaxMetricsEntry holds a VO2 Max data point.
type MaxMetricsEntry struct {
	Generic *struct {
		CalendarDate   string  `json:"calendarDate"`
		VO2MaxValue    float64 `json:"vo2MaxValue"`
		FitnessAge     int     `json:"fitnessAge"`
		FitnessAgeDesc string  `json:"fitnessAgeDescription"`
	} `json:"generic"`
	Cycling *struct {
		CalendarDate string  `json:"calendarDate"`
		VO2MaxValue  float64 `json:"vo2MaxValue"`
		FitnessAge   int     `json:"fitnessAge"`
	} `json:"cycling"`
}

// EnduranceScoreEntry holds an endurance score data point.
type EnduranceScoreEntry struct {
	CalendarDate string  `json:"calendarDate"`
	Score        float64 `json:"overallScore"`
	Contributors []struct {
		ActivityTypeID *int    `json:"activityTypeId"`
		Contribution   float64 `json:"contribution"`
	} `json:"contributors"`
}

// LatestRacePredictions holds the most recent predicted finish times for standard distances.
type LatestRacePredictions struct {
	UserID           int    `json:"userId"`
	CalendarDate     string `json:"calendarDate"`
	Time5K           int    `json:"time5K"`           // seconds
	Time10K          int    `json:"time10K"`          // seconds
	TimeHalfMarathon int    `json:"timeHalfMarathon"` // seconds
	TimeMarathon     int    `json:"timeMarathon"`     // seconds
}

// HillScoreEntry holds a hill score data point.
type HillScoreEntry struct {
	CalendarDate   string  `json:"calendarDate"`
	HillScore      float64 `json:"overallScore"`
	StrengthScore  float64 `json:"strengthScore"`
	EnduranceScore float64 `json:"enduranceScore"`
}

// TrainingReadiness returns the training readiness scores for the given date.
// The API returns an array (typically after-wakeup and realtime entries).
func (c *Client) TrainingReadiness(ctx context.Context, d time.Time) ([]TrainingReadiness, error) {
	var out []TrainingReadiness
	if err := c.get(ctx, fmt.Sprintf("/metrics-service/metrics/trainingreadiness/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TrainingStatusResponse is the aggregated training status response.
type TrainingStatusResponse struct {
	UserID                        int64                     `json:"userId"`
	MostRecentVO2Max              *TrainingStatusVO2Max     `json:"mostRecentVO2Max"`
	MostRecentTrainingLoadBalance *TrainingLoadBalance      `json:"mostRecentTrainingLoadBalance"`
	MostRecentTrainingStatus      *MostRecentTrainingStatus `json:"mostRecentTrainingStatus"`
}

// TrainingStatusVO2Max holds VO2 max data within the training status response.
type TrainingStatusVO2Max struct {
	Generic *TrainingStatusGenericVO2Max `json:"generic"`
	Cycling *TrainingStatusCyclingVO2Max `json:"cycling"`
}

// TrainingStatusGenericVO2Max holds generic (running) VO2 max fields.
type TrainingStatusGenericVO2Max struct {
	CalendarDate       string  `json:"calendarDate"`
	VO2MaxPreciseValue float64 `json:"vo2MaxPreciseValue"`
	FitnessAge         *int    `json:"fitnessAge"`
}

// TrainingStatusCyclingVO2Max holds cycling VO2 max fields.
type TrainingStatusCyclingVO2Max struct {
	CalendarDate       string  `json:"calendarDate"`
	VO2MaxPreciseValue float64 `json:"vo2MaxPreciseValue"`
}

// TrainingLoadBalance holds monthly training load balance data keyed by device ID.
type TrainingLoadBalance struct {
	UserID     int64                                   `json:"userId"`
	MetricsMap map[string]TrainingLoadBalancePerDevice `json:"metricsTrainingLoadBalanceDTOMap"`
}

// TrainingLoadBalancePerDevice holds load balance data for one device.
type TrainingLoadBalancePerDevice struct {
	MonthlyLoadAerobicLow  float64 `json:"monthlyLoadAerobicLow"`
	MonthlyLoadAerobicHigh float64 `json:"monthlyLoadAerobicHigh"`
	MonthlyLoadAnaerobic   float64 `json:"monthlyLoadAnaerobic"`
	PrimaryTrainingDevice  bool    `json:"primaryTrainingDevice"`
}

// MostRecentTrainingStatus holds the latest training status data keyed by device ID.
type MostRecentTrainingStatus struct {
	UserID                   int64                              `json:"userId"`
	LatestTrainingStatusData map[string]PerDeviceTrainingStatus `json:"latestTrainingStatusData"`
}

// PerDeviceTrainingStatus holds training status for a single device.
type PerDeviceTrainingStatus struct {
	CalendarDate          string             `json:"calendarDate"`
	TrainingStatus        int                `json:"trainingStatus"`
	WeeklyTrainingLoad    *float64           `json:"weeklyTrainingLoad"`
	FitnessTrend          int                `json:"fitnessTrend"`
	PrimaryTrainingDevice bool               `json:"primaryTrainingDevice"`
	AcuteTrainingLoad     *AcuteTrainingLoad `json:"acuteTrainingLoadDTO"`
}

// AcuteTrainingLoad holds acute:chronic workload ratio data.
type AcuteTrainingLoad struct {
	ACWRPercent                    int     `json:"acwrPercent"`
	DailyTrainingLoadAcute         int     `json:"dailyTrainingLoadAcute"`
	DailyTrainingLoadChronic       int     `json:"dailyTrainingLoadChronic"`
	DailyAcuteChronicWorkloadRatio float64 `json:"dailyAcuteChronicWorkloadRatio"`
}

// TrainingStatus returns training status metrics for the given date.
func (c *Client) TrainingStatus(ctx context.Context, d time.Time) (*TrainingStatusResponse, error) {
	var out TrainingStatusResponse
	if err := c.get(ctx, fmt.Sprintf("/metrics-service/metrics/trainingstatus/aggregated/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MaxMetrics returns VO2 Max metrics between start and end dates.
func (c *Client) MaxMetrics(ctx context.Context, start, end time.Time) ([]MaxMetricsEntry, error) {
	var out []MaxMetricsEntry
	if err := c.get(ctx, fmt.Sprintf("/metrics-service/metrics/maxmet/daily/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EnduranceScore returns endurance score data between start and end dates.
// The API returns a single current entry wrapped in {"enduranceScoreDTO": {...}}.
func (c *Client) EnduranceScore(ctx context.Context, start, end time.Time) ([]EnduranceScoreEntry, error) {
	params := url.Values{
		"startDate":   {date(start)},
		"endDate":     {date(end)},
		"aggregation": {"weekly"},
	}
	var wrapper struct {
		Entry *EnduranceScoreEntry `json:"enduranceScoreDTO"`
	}
	if err := c.get(ctx, "/metrics-service/metrics/endurancescore/stats", params, &wrapper); err != nil {
		return nil, err
	}
	if wrapper.Entry == nil {
		return nil, nil
	}
	return []EnduranceScoreEntry{*wrapper.Entry}, nil
}

// RacePredictions returns the latest predicted finish times for the user.
func (c *Client) RacePredictions(ctx context.Context) (*LatestRacePredictions, error) {
	var out LatestRacePredictions
	if err := c.get(ctx, fmt.Sprintf("/metrics-service/metrics/racepredictions/latest/%s", c.displayName), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// HillScore returns hill score data between start and end dates.
// The API returns entries wrapped in {"hillScoreDTOList": [...]}.
func (c *Client) HillScore(ctx context.Context, start, end time.Time) ([]HillScoreEntry, error) {
	params := url.Values{
		"startDate":   {date(start)},
		"endDate":     {date(end)},
		"aggregation": {"daily"},
	}
	var wrapper struct {
		Entries []HillScoreEntry `json:"hillScoreDTOList"`
	}
	if err := c.get(ctx, "/metrics-service/metrics/hillscore/stats", params, &wrapper); err != nil {
		return nil, err
	}
	return wrapper.Entries, nil
}

// LactateThresholdEntry holds a single lactate threshold measurement.
// The HearRate field name preserves the API typo.
type LactateThresholdEntry struct {
	UserProfilePK    int64    `json:"userProfilePK"`
	CalendarDate     string   `json:"calendarDate"`
	Speed            *float64 `json:"speed"`            // running LT speed m/s
	HearRate         *int     `json:"hearRate"`         // running LT HR bpm
	HeartRateCycling *int     `json:"heartRateCycling"` // cycling LT HR bpm
	RowSpeed         *float64 `json:"rowSpeed"`
	HeartRateRowing  *int     `json:"heartRateRowing"`
}

// LactateThreshold returns the latest lactate threshold measurement.
func (c *Client) LactateThreshold(ctx context.Context) ([]LactateThresholdEntry, error) {
	var out []LactateThresholdEntry
	if err := c.get(ctx, "/biometric-service/biometric/latestLactateThreshold", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FitnessAge holds fitness-age data for a single date. The headline ages are
// floats; Garmin reports ChronologicalAge as a whole number.
type FitnessAge struct {
	ChronologicalAge     int                            `json:"chronologicalAge"`
	FitnessAge           float64                        `json:"fitnessAge"`
	AchievableFitnessAge float64                        `json:"achievableFitnessAge"`
	PreviousFitnessAge   float64                        `json:"previousFitnessAge"`
	LastUpdated          string                         `json:"lastUpdated"`
	Components           map[string]FitnessAgeComponent `json:"components"`
}

// FitnessAgeComponent is one contributing factor to the fitness age (e.g. bmi,
// rhr, vigorousMinutesAvg). Optional fields are pointers because they vary by
// factor; resting heart rate, for example, carries only Value and Stale.
type FitnessAgeComponent struct {
	Value               float64  `json:"value"`
	TargetValue         *float64 `json:"targetValue"`
	ImprovementValue    *float64 `json:"improvementValue"`
	PotentialAge        *float64 `json:"potentialAge"`
	Priority            *int     `json:"priority"`
	Stale               bool     `json:"stale"`
	NumOfWeeksForIM     *int     `json:"numOfWeeksForIm"`
	LastMeasurementDate string   `json:"lastMeasurementDate"`
}

// FitnessAge returns fitness age data for the given date.
func (c *Client) FitnessAge(ctx context.Context, d time.Time) (*FitnessAge, error) {
	var out FitnessAge
	if err := c.get(ctx, fmt.Sprintf("/fitnessage-service/fitnessage/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RunningToleranceEntry holds a single day of running tolerance data.
// Field names are inferred — verify against a live response if the API returns data.
type RunningToleranceEntry struct {
	CalendarDate string  `json:"calendarDate"`
	Score        float64 `json:"score"`
	Level        int     `json:"level"`
}

// RunningTolerance returns running tolerance statistics between start and end dates.
func (c *Client) RunningTolerance(ctx context.Context, start, end time.Time) ([]RunningToleranceEntry, error) {
	params := url.Values{
		"startDate":   {date(start)},
		"endDate":     {date(end)},
		"aggregation": {"daily"},
	}
	var out []RunningToleranceEntry
	if err := c.get(ctx, "/metrics-service/metrics/runningtolerance/stats", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CyclingFTP holds the latest cycling FTP (functional threshold power)
// estimate. FunctionalThresholdPower is the FTP in watts, or nil when the
// account has no estimate (e.g. no power data recorded).
type CyclingFTP struct {
	UserProfilePK            int64    `json:"userProfilePK"`
	FunctionalThresholdPower *float64 `json:"functionalThresholdPower"`
	Sport                    *string  `json:"sport"`
	CalendarDate             *string  `json:"calendarDate"`
	BiometricSourceType      *string  `json:"biometricSourceType"`
}

// CyclingFTP returns the latest cycling FTP (functional threshold power) estimate.
func (c *Client) CyclingFTP(ctx context.Context) (*CyclingFTP, error) {
	var out CyclingFTP
	if err := c.get(ctx, "/biometric-service/biometric/latestFunctionalThresholdPower/CYCLING", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
