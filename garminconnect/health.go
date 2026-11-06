package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

const menstrualPath = "/periodichealth-service/menstrualcycle"

// MenstrualCycle is one period start with its period and cycle lengths in
// days, as returned by InitMenstrualCycle, ConfirmMenstrualPeriod and
// LastConfirmedMenstrualCycle.
type MenstrualCycle struct {
	UserProfilePK            int64   `json:"userProfilePk"`
	PeriodStartDate          string  `json:"periodStartDate"`
	PeriodLength             int     `json:"periodLength"`
	CycleLength              int     `json:"cycleLength"`
	ReportTimestamp          string  `json:"reportTimestamp"`
	EducationContentMod      *int    `json:"educationContentMod"`
	HasSpecifiedCycleLength  *bool   `json:"hasSpecifiedCycleLength"`
	HasSpecifiedPeriodLength *bool   `json:"hasSpecifiedPeriodLength"`
	CycleType                *string `json:"cycleType"`
	CycleStatus              *string `json:"cycleStatus"`
	PredictedCycle           bool    `json:"predictedCycle"`
}

// MenstrualCycleDay summarizes where a day falls in its cycle. Day counts are
// 1-based within the cycle.
type MenstrualCycleDay struct {
	StartDate             string `json:"startDate"`
	DayInCycle            int    `json:"dayInCycle"`
	PeriodLength          int    `json:"periodLength"`
	CurrentPhase          int    `json:"currentPhase"`
	LengthOfCurrentPhase  int    `json:"lengthOfCurrentPhase"`
	DaysUntilNextPhase    int    `json:"daysUntilNextPhase"`
	FertileWindowStart    int    `json:"fertileWindowStart"`
	LengthOfFertileWindow int    `json:"lengthOfFertileWindow"`
	PredictedCycleLength  int    `json:"predictedCycleLength"`
	EducationContentMod   int    `json:"educationContentMod"`
	CycleType             string `json:"cycleType"`
	FeedbackPhrase        string `json:"feedbackPhrase"`
	PredictedCycle        bool   `json:"predictedCycle"`
}

// MenstrualDayLog is what the user logged for one day. Symptoms, Moods and
// Discharge take Garmin's enum names, e.g. "CRAMPS", "HAPPY", "SPOTTING";
// Flow is "LIGHT", "MEDIUM" or "HEAVY"; SexDrive "LOW", "AVERAGE" or "HIGH";
// SexualActivity "PROTECTED" or "UNPROTECTED".
type MenstrualDayLog struct {
	UserProfilePK   int64    `json:"userProfilePk,omitempty"`
	CalendarDate    string   `json:"calendarDate"`
	Symptoms        []string `json:"symptoms,omitempty"`
	Moods           []string `json:"moods,omitempty"`
	Discharge       []string `json:"discharge,omitempty"`
	Flow            *string  `json:"flow,omitempty"`
	SexDrive        *string  `json:"sexDrive,omitempty"`
	SexualActivity  *string  `json:"sexualActivity,omitempty"`
	Notes           *string  `json:"notes,omitempty"`
	ReportTimestamp string   `json:"reportTimestamp,omitempty"`
	HasBabyMovement bool     `json:"hasBabyMovement,omitempty"`
	OvulationDay    bool     `json:"ovulationDay"`
}

// MenstrualDayView is a day's cycle position and log. DaySummary is nil
// outside a tracked cycle.
type MenstrualDayView struct {
	DaySummary *MenstrualCycleDay `json:"daySummary"`
	DayLog     *MenstrualDayLog   `json:"dayLog"`
}

// MenstrualCalendar lists the cycles and logged days within a date range.
// Logged days are YYYY-MM-DD dates.
type MenstrualCalendar struct {
	CycleSummaries      []MenstrualCycleSummary `json:"cycleSummaries"`
	LoggedSymptomDays   []string                `json:"loggedSymptomDays"`
	LoggedOvulationDays []string                `json:"loggedOvulationDays"`
	LoggedNoteDays      []string                `json:"loggedNoteDays"`
}

// MenstrualCycleSummary is one cycle in a calendar or report. CycleLength,
// CycleType and GapCycle are only set in reports.
type MenstrualCycleSummary struct {
	StartDate             string `json:"startDate"`
	PeriodLength          int    `json:"periodLength"`
	FertileWindowStart    int    `json:"fertileWindowStart"`
	LengthOfFertileWindow int    `json:"lengthOfFertileWindow"`
	EducationContentMod   int    `json:"educationContentMod"`
	PredictedCycle        bool   `json:"predictedCycle"`
	CycleLength           int    `json:"cycleLength"`
	CycleType             string `json:"cycleType"`
	GapCycle              bool   `json:"gapCycle"`
}

// MenstrualReport averages the cycles in a report and lists their day logs.
type MenstrualReport struct {
	AvgCycleLength  float64                 `json:"avgCycleLength"`
	AvgPeriodLength float64                 `json:"avgPeriodLength"`
	CycleSummaries  []MenstrualCycleSummary `json:"cycleSummaries"`
	DailyLogs       []MenstrualDayLog       `json:"dailyLogs"`
}

// MenstrualPeriod is a confirmed period of Length days starting at Start.
type MenstrualPeriod struct {
	Start  time.Time
	Length int
}

func reportTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000")
}

// MenstrualDayView returns the cycle position and log for a day. Without cycle
// tracking set up it returns ErrNoData.
func (c *Client) MenstrualDayView(ctx context.Context, d time.Time) (*MenstrualDayView, error) {
	var out MenstrualDayView
	if err := c.get(ctx, menstrualPath+"/dayview/"+date(d), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MenstrualCalendar returns the cycles with days between start and end. Garmin
// rejects ranges over 90 days.
func (c *Client) MenstrualCalendar(ctx context.Context, start, end time.Time) (*MenstrualCalendar, error) {
	var out MenstrualCalendar
	if err := c.get(ctx, fmt.Sprintf("%s/calendar/%s/%s", menstrualPath, date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MenstrualCycleDay returns where d falls in its cycle. Outside a tracked
// cycle it returns ErrNoData.
func (c *Client) MenstrualCycleDay(ctx context.Context, d time.Time) (*MenstrualCycleDay, error) {
	var out MenstrualCycleDay
	if err := c.get(ctx, menstrualPath+"/summary/"+date(d), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LastConfirmedMenstrualCycle returns the last confirmed cycle as of d. If
// there is none it returns ErrNoData.
func (c *Client) LastConfirmedMenstrualCycle(ctx context.Context, d time.Time) (*MenstrualCycle, error) {
	var out MenstrualCycle
	if err := c.get(ctx, menstrualPath+"/lastconfirmed/"+date(d), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MenstrualReport returns a report over the last cycles cycles ending at end.
// Garmin accepts cycles of 1, 6 or 12.
func (c *Client) MenstrualReport(ctx context.Context, end time.Time, cycles int) (*MenstrualReport, error) {
	params := url.Values{"next": {"false"}, "reportType": {"CYCLE"}}
	var out MenstrualReport
	path := fmt.Sprintf("/periodichealth-service/reports/menstrualcycle/%d/%s", cycles, date(end))
	if err := c.get(ctx, path, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PregnancySummary returns the current pregnancy snapshot. If the user isn't
// tracking a pregnancy it returns ErrNoData.
func (c *Client) PregnancySummary(ctx context.Context) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, menstrualPath+"/pregnancysnapshot", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// InitMenstrualCycle starts cycle tracking with a first period. Logging days
// also needs the tracking settings, which UpdateUserSettings sets, e.g.
//
//	map[string]any{"userMenstrualCycleSettings": map[string]any{"menstrualCycleType": "REGULAR"}}
func (c *Client) InitMenstrualCycle(ctx context.Context, periodStart time.Time, periodLength, cycleLength int) (*MenstrualCycle, error) {
	body := map[string]any{
		"periodStartDate": date(periodStart),
		"periodLength":    periodLength,
		"cycleLength":     cycleLength,
		"reportTimestamp": reportTimestamp(),
	}
	var out MenstrualCycle
	if err := c.post(ctx, menstrualPath+"/initCycleSetup", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ConfirmMenstrualPeriod confirms a period starting at periodStart. Set
// predicted to confirm a cycle Garmin predicted.
func (c *Client) ConfirmMenstrualPeriod(ctx context.Context, periodStart time.Time, periodLength, cycleLength int, predicted bool) (*MenstrualCycle, error) {
	body := map[string]any{
		"periodStartDate":          date(periodStart),
		"periodLength":             periodLength,
		"cycleLength":              cycleLength,
		"predictedCycle":           predicted,
		"hasSpecifiedCycleLength":  true,
		"hasSpecifiedPeriodLength": true,
		"reportTimestamp":          reportTimestamp(),
	}
	var out MenstrualCycle
	if err := c.post(ctx, menstrualPath+"/"+date(periodStart), body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetMenstrualDayLog replaces the log for log.CalendarDate (YYYY-MM-DD). It
// is a full snapshot: nil fields are cleared, except a nil Notes keeps the
// stored note; set Notes to "" to clear it.
func (c *Client) SetMenstrualDayLog(ctx context.Context, log MenstrualDayLog) (*MenstrualDayLog, error) {
	log.ReportTimestamp = reportTimestamp()
	var out MenstrualDayLog
	if err := c.post(ctx, menstrualPath+"/dailylog/"+url.PathEscape(log.CalendarDate), log, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetMenstrualPeriods replaces the confirmed periods between start and end
// with periods. Periods in the range that aren't passed are removed, so an
// empty list clears the range.
func (c *Client) SetMenstrualPeriods(ctx context.Context, start, end time.Time, periods []MenstrualPeriod) error {
	lists := make([][]string, 0, len(periods))
	for _, p := range periods {
		days := make([]string, p.Length)
		for i := range days {
			days[i] = date(p.Start.AddDate(0, 0, i))
		}
		lists = append(lists, days)
	}
	body := map[string]any{
		"todayCalendarDate": date(time.Now()),
		"startDate":         date(start),
		"endDate":           date(end),
		"cycleDatesLists":   lists,
		"futureEditsByFE":   true,
		"reportTimestamp":   reportTimestamp(),
	}
	return c.post(ctx, menstrualPath+"/calendarupdates", body, nil)
}
