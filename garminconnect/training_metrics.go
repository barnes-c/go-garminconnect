package garminconnect

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// TrainingReadinessRange returns training readiness for every day between
// start and end.
func (c *Client) TrainingReadinessRange(ctx context.Context, start, end time.Time) ([]TrainingReadiness, error) {
	var out []TrainingReadiness
	if err := c.get(ctx, fmt.Sprintf("/metrics-service/metrics/trainingreadiness/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TrainingStatusHistory returns the daily training status between start and
// end, keyed by device ID.
func (c *Client) TrainingStatusHistory(ctx context.Context, start, end time.Time) (map[string][]PerDeviceTrainingStatus, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	params := url.Values{"fromCalendarDate": {date(start)}, "toCalendarDate": {date(end)}}
	var out struct {
		ReportData map[string][]PerDeviceTrainingStatus `json:"reportData"`
	}
	if err := c.get(ctx, "/metrics-service/metrics/trainingstatus/weekly/"+name, params, &out); err != nil {
		return nil, err
	}
	return out.ReportData, nil
}

// HeatAltitudeAcclimation is how acclimated the user is to heat and altitude.
type HeatAltitudeAcclimation struct {
	CalendarDate              string `json:"calendarDate"`
	HeatAcclimationPercentage int    `json:"heatAcclimationPercentage"`
	HeatAcclimationDate       string `json:"heatAcclimationDate"`
	AltitudeAcclimation       int    `json:"altitudeAcclimation"` // meters
	AltitudeAcclimationDate   string `json:"altitudeAcclimationDate"`
	CurrentAltitude           int    `json:"currentAltitude"` // meters
	AcclimationPercentage     int    `json:"acclimationPercentage"`
}

// HeatAltitudeAcclimation returns the latest heat and altitude acclimation as
// of d.
func (c *Client) HeatAltitudeAcclimation(ctx context.Context, d time.Time) (*HeatAltitudeAcclimation, error) {
	var out HeatAltitudeAcclimation
	if err := c.get(ctx, "/metrics-service/metrics/heataltitudeacclimation/latest/"+date(d), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// BiometricSample is one measured biometric value and the day it was taken.
type BiometricSample struct {
	SampleID int64   `json:"sampleId"`
	BioDate  string  `json:"bioDate"`
	Value    float64 `json:"value"`
}

// CurrentBiometrics are the user's latest threshold values. A sample is nil
// when Garmin has none.
type CurrentBiometrics struct {
	LactateThresholdSpeed            *BiometricSample `json:"lactateThresholdSpeed"` // meters per second
	LactateThresholdHeartRate        *BiometricSample `json:"lactateThresholdHeartRate"`
	LactateThresholdHeartRateCycling *BiometricSample `json:"lactateThresholdHeartRateCycling"`
	FunctionalThresholdPower         *BiometricSample `json:"functionalThresholdPower"` // watts
}

// CurrentBiometrics returns the user's latest lactate threshold and FTP.
func (c *Client) CurrentBiometrics(ctx context.Context) (*CurrentBiometrics, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	var out CurrentBiometrics
	if err := c.get(ctx, "/userprofile-service/biometric/current/"+name, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FitnessAgeDay is the fitness age and its inputs on one day.
type FitnessAgeDay struct {
	CalendarDate string `json:"calendarDate"`
	Values       struct {
		FitnessAge           float64 `json:"fitnessAge"`
		AchievableFitnessAge float64 `json:"achievableFitnessAge"`
		VigorousDaysAvg      float64 `json:"vigorousDaysAvg"`
		RHR                  int     `json:"rhr"`
		BMI                  float64 `json:"bmi"`
	} `json:"values"`
}

// FitnessAgeHistory returns the fitness age for every day between start and
// end.
func (c *Client) FitnessAgeHistory(ctx context.Context, start, end time.Time) ([]FitnessAgeDay, error) {
	var out []FitnessAgeDay
	if err := c.get(ctx, fmt.Sprintf("/fitnessage-service/stats/daily/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PowerCurvePoint is the best average power held for Duration seconds, and the
// activity it came from.
type PowerCurvePoint struct {
	Duration     int    `json:"duration"` // seconds
	Power        int    `json:"power"`    // watts
	ActivityID   int64  `json:"activityId"`
	ActivityDate string `json:"activityDate"`
}

// PowerCurve returns the best power for a range of durations over activities
// of the given sport (e.g. "cycling") between start and end.
func (c *Client) PowerCurve(ctx context.Context, start, end time.Time, sport string) ([]PowerCurvePoint, error) {
	params := url.Values{"startDate": {date(start)}, "endDate": {date(end)}, "sport": {sport}}
	var out struct {
		Entries []PowerCurvePoint `json:"entries"`
	}
	if err := c.get(ctx, "/fitnessstats-service/powerCurve", params, &out); err != nil {
		return nil, err
	}
	return out.Entries, nil
}

// CyclingAbility rates the user's cycling energy systems (each 0-100) and
// derives a rider profile from them.
type CyclingAbility struct {
	CalendarDate        string `json:"calendarDate"`
	AerobicEndurance    int    `json:"aerobicEndurance"`
	AerobicCapacity     int    `json:"aerobicCapacity"`
	AnaerobicCapacity   int    `json:"anaerobicCapacity"`
	ProfileType         string `json:"profileType"`
	ProfileTypeFeedback string `json:"profileTypeFeedback"`
}

// CyclingAbility returns the latest cycling ability as of d.
func (c *Client) CyclingAbility(ctx context.Context, d time.Time) (*CyclingAbility, error) {
	var out CyclingAbility
	params := url.Values{"calendarDate": {date(d)}}
	if err := c.get(ctx, "/metrics-service/metrics/cycling-ability/latest", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
