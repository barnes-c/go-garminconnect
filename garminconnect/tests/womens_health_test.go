package garminconnect_test

import (
	"testing"
	"time"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The tests below need a watch, so they are recorded against a real account
// on a day with data (GARMIN_SUMMARY_DATE / GARMIN_SLEEP_DATE). Ranges query a
// single day; see dataDate.

// TestMenstrualCycleLifecycle sets up cycle tracking, confirms a period, logs a
// day and checks the reads, then clears the log and the period again.
func TestMenstrualCycleLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	start := testDate.AddDate(0, 0, -12)
	logDay := start.AddDate(0, 0, 2)
	rangeStart := testDate.AddDate(0, 0, -31)

	err := c.UpdateUserSettings(ctx, map[string]any{
		"userMenstrualCycleSettings": map[string]any{"menstrualCycleType": "REGULAR"},
	})
	skipAPIError(t, err)
	require.NoError(t, err)

	cycle, err := c.InitMenstrualCycle(ctx, start, 5, 28)
	require.NoError(t, err)
	assert.NotEmpty(t, cycle.PeriodStartDate)
	cycle, err = c.ConfirmMenstrualPeriod(ctx, start, 5, 28, false)
	require.NoError(t, err)
	assert.NotZero(t, cycle.CycleLength)

	flow, notes := "LIGHT", "note"
	logged, err := c.SetMenstrualDayLog(ctx, gc.MenstrualDayLog{
		CalendarDate: logDay.Format(time.DateOnly),
		Symptoms:     []string{"CRAMPS"},
		Moods:        []string{"HAPPY"},
		Flow:         &flow,
		Notes:        &notes,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, logged.Symptoms)

	view, err := c.MenstrualDayView(ctx, logDay)
	require.NoError(t, err)
	require.NotNil(t, view.DaySummary)
	assert.NotZero(t, view.DaySummary.DayInCycle)
	require.NotNil(t, view.DayLog)
	assert.NotEmpty(t, view.DayLog.Moods)

	cal, err := c.MenstrualCalendar(ctx, rangeStart, testDate)
	require.NoError(t, err)
	assert.NotEmpty(t, cal.CycleSummaries)
	assert.NotEmpty(t, cal.LoggedSymptomDays)

	day, err := c.MenstrualCycleDay(ctx, testDate)
	require.NoError(t, err)
	require.NotNil(t, day)
	assert.NotEmpty(t, day.CycleType)

	last, err := c.LastConfirmedMenstrualCycle(ctx, testDate)
	require.NoError(t, err)
	require.NotNil(t, last)
	assert.NotZero(t, last.PeriodLength)

	report, err := c.MenstrualReport(ctx, testDate, 1)
	require.NoError(t, err)
	assert.NotEmpty(t, report.CycleSummaries)
	assert.NotEmpty(t, report.DailyLogs)

	empty := ""
	_, err = c.SetMenstrualDayLog(ctx, gc.MenstrualDayLog{CalendarDate: logDay.Format(time.DateOnly), Notes: &empty})
	require.NoError(t, err)
	require.NoError(t, c.SetMenstrualPeriods(ctx, rangeStart, testDate, nil))
	cal, err = c.MenstrualCalendar(ctx, rangeStart, testDate)
	require.NoError(t, err)
	assert.Empty(t, cal.CycleSummaries)
}

func TestPregnancySummary(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// The recording account isn't tracking a pregnancy: Garmin answers 204.
	_, err := c.PregnancySummary(t.Context())
	require.ErrorIs(t, err, gc.ErrNoData)
}
